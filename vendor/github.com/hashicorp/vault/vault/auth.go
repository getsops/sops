package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// coreAuthConfigPath is used to store the auth configuration.
	// Auth configuration is protected within the Vault itself, which means it
	// can only be viewed or modified after an unseal.
	coreAuthConfigPath = "core/auth"

	// coreLocalAuthConfigPath is used to store credential configuration for
	// local (non-replicated) mounts
	coreLocalAuthConfigPath = "core/local-auth"

	// credentialBarrierPrefix is the prefix to the UUID used in the
	// barrier view for the credential backends.
	credentialBarrierPrefix = "auth/"

	// credentialRoutePrefix is the mount prefix used for the router
	credentialRoutePrefix = "auth/"

	// credentialTableType is the value we expect to find for the credential
	// table and corresponding entries
	credentialTableType = "auth"
)

var (
	// errLoadAuthFailed if loadCredentials encounters an error
	errLoadAuthFailed = errors.New("failed to setup auth table")

	// credentialAliases maps old backend names to new backend names, allowing us
	// to move/rename backends but maintain backwards compatibility
	credentialAliases = map[string]string{"aws-ec2": "aws"}
)

// enableCredential is used to enable a new credential backend
func (c *Core) enableCredential(ctx context.Context, entry *MountEntry) error {
	return c.enableCredentialInternal(ctx, entry, MountTableUpdateStorage)
}

// enableCredential is used to enable a new credential backend
func (c *Core) enableCredentialInternal(ctx context.Context, entry *MountEntry, updateStorage bool) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(entry.Path, "/") {
		entry.Path += "/"
	}

	// Ensure there is a name
	if entry.Path == "/" {
		return fmt.Errorf("backend path must be specified")
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	entry.NamespaceID = ns.ID
	entry.namespace = ns

	// Populate cache
	NamespaceByID(ctx, ns.ID, c)

	// Basic check for matching names
	for _, ent := range c.auth.Entries {
		if ns.ID == ent.NamespaceID {
			switch {
			// Existing is oauth/github/ new is oauth/ or
			// existing is oauth/ and new is oauth/github/
			case strings.HasPrefix(ent.Path, entry.Path):
				fallthrough
			case strings.HasPrefix(entry.Path, ent.Path):
				return logical.CodedError(409, fmt.Sprintf("path is already in use at %s", ent.Path))
			}
		}
	}

	// Ensure the token backend is a singleton
	if entry.Type == "token" {
		return fmt.Errorf("token credential backend cannot be instantiated")
	}

	// Check for conflicts according to the router
	if conflict := c.router.MountConflict(ctx, credentialRoutePrefix+entry.Path); conflict != "" {
		return logical.CodedError(409, fmt.Sprintf("existing mount at %s", conflict))
	}

	// Generate a new UUID and view
	if entry.UUID == "" {
		entryUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		entry.UUID = entryUUID
	}
	if entry.BackendAwareUUID == "" {
		bUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		entry.BackendAwareUUID = bUUID
	}
	if entry.Accessor == "" {
		accessor, err := c.generateMountAccessor("auth_" + entry.Type)
		if err != nil {
			return err
		}
		entry.Accessor = accessor
	}
	// Sync values to the cache
	entry.SyncCache()

	viewPath := entry.ViewPath()
	view := NewBarrierView(c.barrier, viewPath)

	nilMount, err := preprocessMount(c, entry, view)
	if err != nil {
		return err
	}
	origViewReadOnlyErr := view.getReadOnlyErr()

	// Mark the view as read-only until the mounting is complete and
	// ensure that it is reset after. This ensures that there will be no
	// writes during the construction of the backend.
	view.setReadOnlyErr(logical.ErrSetupReadOnly)
	defer view.setReadOnlyErr(origViewReadOnlyErr)

	var backend logical.Backend
	// Create the new backend
	sysView := c.mountEntrySysView(entry)
	backend, err = c.newCredentialBackend(ctx, entry, sysView, view)
	if err != nil {
		return err
	}
	if backend == nil {
		return fmt.Errorf("nil backend returned from %q factory", entry.Type)
	}

	// Check for the correct backend type
	backendType := backend.Type()
	if backendType != logical.TypeCredential {
		return fmt.Errorf("cannot mount %q of type %q as an auth backend", entry.Type, backendType)
	}

	addPathCheckers(c, entry, backend, viewPath)

	// If the mount is filtered or we are on a DR secondary we don't want to
	// keep the actual backend running, so we clean it up and set it to nil
	// so the router does not have a pointer to the object.
	if nilMount {
		backend.Cleanup(ctx)
		backend = nil
	}

	// Update the auth table
	newTable := c.auth.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if updateStorage {
		if err := c.persistAuth(ctx, newTable, &entry.Local); err != nil {
			if err == logical.ErrReadOnly && c.perfStandby {
				return err
			}
			return errors.New("failed to update auth table")
		}
	}

	c.auth = newTable

	if err := c.router.Mount(backend, credentialRoutePrefix+entry.Path, entry, view); err != nil {
		return err
	}

	if !nilMount {
		// restore the original readOnlyErr, so we can write to the view in
		// Initialize() if necessary
		view.setReadOnlyErr(origViewReadOnlyErr)
		// initialize, using the core's active context.
		err := backend.Initialize(c.activeContext, &logical.InitializationRequest{Storage: view})
		if err != nil {
			return err
		}
	}

	if c.logger.IsInfo() {
		c.logger.Info("enabled credential backend", "path", entry.Path, "type", entry.Type)
	}
	return nil
}

// disableCredential is used to disable an existing credential backend; the
// boolean indicates if it existed
func (c *Core) disableCredential(ctx context.Context, path string) error {
	// Ensure we end the path in a slash
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Ensure the token backend is not affected
	if path == "token/" {
		return fmt.Errorf("token credential backend cannot be disabled")
	}

	return c.disableCredentialInternal(ctx, path, MountTableUpdateStorage)
}

func (c *Core) disableCredentialInternal(ctx context.Context, path string, updateStorage bool) error {
	path = credentialRoutePrefix + path

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	// Verify exact match of the route
	match := c.router.MatchingMount(ctx, path)
	if match == "" || ns.Path+path != match {
		return fmt.Errorf("no matching mount")
	}

	// Store the view for this backend
	view := c.router.MatchingStorageByAPIPath(ctx, path)
	if view == nil {
		return fmt.Errorf("no matching backend %q", path)
	}

	// Get the backend/mount entry for this path, used to remove ignored
	// replication prefixes
	backend := c.router.MatchingBackend(ctx, path)
	entry := c.router.MatchingMountEntry(ctx, path)

	// Mark the entry as tainted
	if err := c.taintCredEntry(ctx, path, updateStorage); err != nil {
		return err
	}

	// Taint the router path to prevent routing
	if err := c.router.Taint(ctx, path); err != nil {
		return err
	}

	if c.expiration != nil && backend != nil {
		// Revoke credentials from this path
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return err
		}
		revokeCtx := namespace.ContextWithNamespace(c.activeContext, ns)
		if err := c.expiration.RevokePrefix(revokeCtx, path, true); err != nil {
			return err
		}
	}

	if backend != nil {
		// Call cleanup function if it exists
		backend.Cleanup(ctx)
	}

	viewPath := entry.ViewPath()
	switch {
	case !updateStorage:
		// Don't attempt to clear data, replication will handle this
	case c.IsDRSecondary():
		// If we are a dr secondary we want to clear the view, but the provided
		// view is marked as read only. We use the barrier here to get around
		// it.

		if err := logical.ClearViewWithLogging(ctx, NewBarrierView(c.barrier, viewPath), c.logger.Named("auth.deletion").With("namespace", ns.ID, "path", path)); err != nil {
			c.logger.Error("failed to clear view for path being unmounted", "error", err, "path", path)
			return err
		}

	case entry.Local, !c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary):
		// Have writable storage, remove the whole thing
		if err := logical.ClearViewWithLogging(ctx, view, c.logger.Named("auth.deletion").With("namespace", ns.ID, "path", path)); err != nil {
			c.logger.Error("failed to clear view for path being unmounted", "error", err, "path", path)
			return err
		}

	case !entry.Local && c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary):
		if err := clearIgnoredPaths(ctx, c, backend, viewPath); err != nil {
			return err
		}
	}

	// Remove the mount table entry
	if err := c.removeCredEntry(ctx, strings.TrimPrefix(path, credentialRoutePrefix), updateStorage); err != nil {
		return err
	}

	// Unmount the backend
	if err := c.router.Unmount(ctx, path); err != nil {
		return err
	}

	removePathCheckers(c, entry, viewPath)

	if c.logger.IsInfo() {
		c.logger.Info("disabled credential backend", "path", path)
	}

	return nil
}

// removeCredEntry is used to remove an entry in the auth table
func (c *Core) removeCredEntry(ctx context.Context, path string, updateStorage bool) error {
	c.authLock.Lock()
	defer c.authLock.Unlock()

	// Taint the entry from the auth table
	newTable := c.auth.shallowClone()
	entry, err := newTable.remove(ctx, path)
	if err != nil {
		return err
	}
	if entry == nil {
		c.logger.Error("nil entry found removing entry in auth table", "path", path)
		return logical.CodedError(500, "failed to remove entry in auth table")
	}

	if updateStorage {
		// Update the auth table
		if err := c.persistAuth(ctx, newTable, &entry.Local); err != nil {
			if err == logical.ErrReadOnly && c.perfStandby {
				return err
			}

			return errors.New("failed to update auth table")
		}
	}

	c.auth = newTable

	return nil
}

// remountCredEntryForce takes a copy of the mount entry for the path and fully
// unmounts and remounts the backend to pick up any changes, such as filtered
// paths
func (c *Core) remountCredEntryForce(ctx context.Context, path string) error {
	fullPath := credentialRoutePrefix + path
	me := c.router.MatchingMountEntry(ctx, fullPath)
	if me == nil {
		return fmt.Errorf("cannot find mount for path %q", path)
	}

	me, err := me.Clone()
	if err != nil {
		return err
	}

	if err := c.disableCredential(ctx, path); err != nil {
		return err
	}
	return c.enableCredential(ctx, me)
}

// taintCredEntry is used to mark an entry in the auth table as tainted
func (c *Core) taintCredEntry(ctx context.Context, path string, updateStorage bool) error {
	c.authLock.Lock()
	defer c.authLock.Unlock()

	// Taint the entry from the auth table
	// We do this on the original since setting the taint operates
	// on the entries which a shallow clone shares anyways
	entry, err := c.auth.setTaint(ctx, strings.TrimPrefix(path, credentialRoutePrefix), true)
	if err != nil {
		return err
	}

	// Ensure there was a match
	if entry == nil {
		return fmt.Errorf("no matching backend")
	}

	if updateStorage {
		// Update the auth table
		if err := c.persistAuth(ctx, c.auth, &entry.Local); err != nil {
			if err == logical.ErrReadOnly && c.perfStandby {
				return err
			}
			return errors.New("failed to update auth table")
		}
	}

	return nil
}

// loadCredentials is invoked as part of postUnseal to load the auth table
func (c *Core) loadCredentials(ctx context.Context) error {
	// Load the existing mount table
	raw, err := c.barrier.Get(ctx, coreAuthConfigPath)
	if err != nil {
		c.logger.Error("failed to read auth table", "error", err)
		return errLoadAuthFailed
	}
	rawLocal, err := c.barrier.Get(ctx, coreLocalAuthConfigPath)
	if err != nil {
		c.logger.Error("failed to read local auth table", "error", err)
		return errLoadAuthFailed
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

	if raw != nil {
		authTable, err := c.decodeMountTable(ctx, raw.Value)
		if err != nil {
			c.logger.Error("failed to decompress and/or decode the auth table", "error", err)
			return err
		}
		c.auth = authTable
	}

	var needPersist bool
	if c.auth == nil {
		c.auth = c.defaultAuthTable()
		needPersist = true
	}

	if rawLocal != nil {
		localAuthTable, err := c.decodeMountTable(ctx, rawLocal.Value)
		if err != nil {
			c.logger.Error("failed to decompress and/or decode the local mount table", "error", err)
			return err
		}
		if localAuthTable != nil && len(localAuthTable.Entries) > 0 {
			c.auth.Entries = append(c.auth.Entries, localAuthTable.Entries...)
		}
	}

	// Upgrade to typed auth table
	if c.auth.Type == "" {
		c.auth.Type = credentialTableType
		needPersist = true
	}

	// Upgrade to table-scoped entries
	for _, entry := range c.auth.Entries {
		if entry.Table == "" {
			entry.Table = c.auth.Type
			needPersist = true
		}
		if entry.Accessor == "" {
			accessor, err := c.generateMountAccessor("auth_" + entry.Type)
			if err != nil {
				return err
			}
			entry.Accessor = accessor
			needPersist = true
		}
		if entry.BackendAwareUUID == "" {
			bUUID, err := uuid.GenerateUUID()
			if err != nil {
				return err
			}
			entry.BackendAwareUUID = bUUID
			needPersist = true
		}

		if entry.NamespaceID == "" {
			entry.NamespaceID = namespace.RootNamespaceID
			needPersist = true
		}
		ns, err := NamespaceByID(ctx, entry.NamespaceID, c)
		if err != nil {
			return err
		}
		if ns == nil {
			return namespace.ErrNoNamespace
		}
		entry.namespace = ns

		// Sync values to the cache
		entry.SyncCache()
	}

	if !needPersist {
		return nil
	}

	if err := c.persistAuth(ctx, c.auth, nil); err != nil {
		c.logger.Error("failed to persist auth table", "error", err)
		return errLoadAuthFailed
	}

	return nil
}

// persistAuth is used to persist the auth table after modification
func (c *Core) persistAuth(ctx context.Context, table *MountTable, local *bool) error {
	if table.Type != credentialTableType {
		c.logger.Error("given table to persist has wrong type", "actual_type", table.Type, "expected_type", credentialTableType)
		return fmt.Errorf("invalid table type given, not persisting")
	}

	for _, entry := range table.Entries {
		if entry.Table != table.Type {
			c.logger.Error("given entry to persist in auth table has wrong table value", "path", entry.Path, "entry_table_type", entry.Table, "actual_type", table.Type)
			return fmt.Errorf("invalid auth entry found, not persisting")
		}
	}

	nonLocalAuth := &MountTable{
		Type: credentialTableType,
	}

	localAuth := &MountTable{
		Type: credentialTableType,
	}

	for _, entry := range table.Entries {
		if entry.Local {
			localAuth.Entries = append(localAuth.Entries, entry)
		} else {
			nonLocalAuth.Entries = append(nonLocalAuth.Entries, entry)
		}
	}

	writeTable := func(mt *MountTable, path string) error {
		// Encode the mount table into JSON and compress it (lzw).
		compressedBytes, err := jsonutil.EncodeJSONAndCompress(mt, nil)
		if err != nil {
			c.logger.Error("failed to encode or compress auth mount table", "error", err)
			return err
		}

		// Create an entry
		entry := &logical.StorageEntry{
			Key:   path,
			Value: compressedBytes,
		}

		// Write to the physical backend
		if err := c.barrier.Put(ctx, entry); err != nil {
			c.logger.Error("failed to persist auth mount table", "error", err)
			return err
		}
		return nil
	}

	var err error
	switch {
	case local == nil:
		// Write non-local mounts
		err := writeTable(nonLocalAuth, coreAuthConfigPath)
		if err != nil {
			return err
		}

		// Write local mounts
		err = writeTable(localAuth, coreLocalAuthConfigPath)
		if err != nil {
			return err
		}
	case *local:
		err = writeTable(localAuth, coreLocalAuthConfigPath)
	default:
		err = writeTable(nonLocalAuth, coreAuthConfigPath)
	}

	return err
}

// setupCredentials is invoked after we've loaded the auth table to
// initialize the credential backends and setup the router
func (c *Core) setupCredentials(ctx context.Context) error {
	var persistNeeded bool

	c.authLock.Lock()
	defer c.authLock.Unlock()

	for _, entry := range c.auth.sortEntriesByPathDepth().Entries {
		var backend logical.Backend

		// Create a barrier view using the UUID
		viewPath := entry.ViewPath()

		// Singleton mounts cannot be filtered on a per-secondary basis
		// from replication
		if strutil.StrListContains(singletonMounts, entry.Type) {
			addFilterablePath(c, viewPath)
		}

		view := NewBarrierView(c.barrier, viewPath)

		// Determining the replicated state of the mount
		nilMount, err := preprocessMount(c, entry, view)
		if err != nil {
			return err
		}
		origViewReadOnlyErr := view.getReadOnlyErr()

		// Mark the view as read-only until the mounting is complete and
		// ensure that it is reset after. This ensures that there will be no
		// writes during the construction of the backend.
		view.setReadOnlyErr(logical.ErrSetupReadOnly)
		if strutil.StrListContains(singletonMounts, entry.Type) {
			defer view.setReadOnlyErr(origViewReadOnlyErr)
		} else {
			c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
				view.setReadOnlyErr(origViewReadOnlyErr)
			})
		}

		// Initialize the backend
		sysView := c.mountEntrySysView(entry)

		backend, err = c.newCredentialBackend(ctx, entry, sysView, view)
		if err != nil {
			c.logger.Error("failed to create credential entry", "path", entry.Path, "error", err)
			if !c.builtinRegistry.Contains(entry.Type, consts.PluginTypeCredential) {
				// If we encounter an error instantiating the backend due to an error,
				// skip backend initialization but register the entry to the mount table
				// to preserve storage and path.
				c.logger.Warn("skipping plugin-based credential entry", "path", entry.Path)
				goto ROUTER_MOUNT
			}
			return errLoadAuthFailed
		}
		if backend == nil {
			return fmt.Errorf("nil backend returned from %q factory", entry.Type)
		}

		{
			// Check for the correct backend type
			backendType := backend.Type()
			if backendType != logical.TypeCredential {
				return fmt.Errorf("cannot mount %q of type %q as an auth backend", entry.Type, backendType)
			}

			addPathCheckers(c, entry, backend, viewPath)
		}

		// If the mount is filtered or we are on a DR secondary we don't want to
		// keep the actual backend running, so we clean it up and set it to nil
		// so the router does not have a pointer to the object.
		if nilMount {
			backend.Cleanup(ctx)
			backend = nil
		}

	ROUTER_MOUNT:
		// Mount the backend
		path := credentialRoutePrefix + entry.Path
		err = c.router.Mount(backend, path, entry, view)
		if err != nil {
			c.logger.Error("failed to mount auth entry", "path", entry.Path, "error", err)
			return errLoadAuthFailed
		}

		if c.logger.IsInfo() {
			c.logger.Info("successfully enabled credential backend", "type", entry.Type, "path", entry.Path)
		}

		// Ensure the path is tainted if set in the mount table
		if entry.Tainted {
			c.router.Taint(ctx, path)
		}

		// Check if this is the token store
		if entry.Type == "token" {
			c.tokenStore = backend.(*TokenStore)

			// At some point when this isn't beta we may persist this but for
			// now always set it on mount
			entry.Config.TokenType = logical.TokenTypeDefaultService

			// this is loaded *after* the normal mounts, including cubbyhole
			c.router.tokenStoreSaltFunc = c.tokenStore.Salt
			if !c.IsDRSecondary() {
				c.tokenStore.cubbyholeBackend = c.router.MatchingBackend(ctx, cubbyholeMountPath).(*CubbyholeBackend)
			}
		}

		// Populate cache
		NamespaceByID(ctx, entry.NamespaceID, c)

		// Initialize
		if !nilMount {
			// Bind locally
			localEntry := entry
			c.postUnsealFuncs = append(c.postUnsealFuncs, func() {
				if backend == nil {
					c.logger.Error("skipping initialization on nil backend", "path", localEntry.Path)
					return
				}

				err := backend.Initialize(ctx, &logical.InitializationRequest{Storage: view})
				if err != nil {
					c.logger.Error("failed to initialize auth entry", "path", localEntry.Path, "error", err)
				}
			})
		}
	}

	if persistNeeded {
		// persist non-local auth
		return c.persistAuth(ctx, c.auth, nil)
	}

	return nil
}

// teardownCredentials is used before we seal the vault to reset the credential
// backends to their unloaded state. This is reversed by loadCredentials.
func (c *Core) teardownCredentials(ctx context.Context) error {
	c.authLock.Lock()
	defer c.authLock.Unlock()

	if c.auth != nil {
		authTable := c.auth.shallowClone()
		for _, e := range authTable.Entries {
			backend := c.router.MatchingBackend(namespace.ContextWithNamespace(ctx, e.namespace), credentialRoutePrefix+e.Path)
			if backend != nil {
				backend.Cleanup(ctx)
			}

			viewPath := e.ViewPath()
			removePathCheckers(c, e, viewPath)
		}
	}

	c.auth = nil
	c.tokenStore = nil
	return nil
}

// newCredentialBackend is used to create and configure a new credential backend by name
func (c *Core) newCredentialBackend(ctx context.Context, entry *MountEntry, sysView logical.SystemView, view logical.Storage) (logical.Backend, error) {
	t := entry.Type
	if alias, ok := credentialAliases[t]; ok {
		t = alias
	}

	f, ok := c.credentialBackends[t]
	if !ok {
		f = plugin.Factory
	}

	// Set up conf to pass in plugin_name
	conf := make(map[string]string, len(entry.Options)+1)
	for k, v := range entry.Options {
		conf[k] = v
	}

	switch {
	case entry.Type == "plugin":
		conf["plugin_name"] = entry.Config.PluginName
	default:
		conf["plugin_name"] = t
	}

	conf["plugin_type"] = consts.PluginTypeCredential.String()

	authLogger := c.baseLogger.Named(fmt.Sprintf("auth.%s.%s", t, entry.Accessor))
	c.AddLogger(authLogger)
	config := &logical.BackendConfig{
		StorageView: view,
		Logger:      authLogger,
		Config:      conf,
		System:      sysView,
		BackendUUID: entry.BackendAwareUUID,
	}

	b, err := f(ctx, config)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// defaultAuthTable creates a default auth table
func (c *Core) defaultAuthTable() *MountTable {
	table := &MountTable{
		Type: credentialTableType,
	}
	tokenUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not generate UUID for default auth table token entry: %v", err))
	}
	tokenAccessor, err := c.generateMountAccessor("auth_token")
	if err != nil {
		panic(fmt.Sprintf("could not generate accessor for default auth table token entry: %v", err))
	}
	tokenBackendUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(fmt.Sprintf("could not create identity backend UUID: %v", err))
	}
	tokenAuth := &MountEntry{
		Table:            credentialTableType,
		Path:             "token/",
		Type:             "token",
		Description:      "token based credentials",
		UUID:             tokenUUID,
		Accessor:         tokenAccessor,
		BackendAwareUUID: tokenBackendUUID,
	}
	table.Entries = append(table.Entries, tokenAuth)
	return table
}
