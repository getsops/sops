package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
)

type ctxKeyForwardedRequestMountAccessor struct{}

func (c ctxKeyForwardedRequestMountAccessor) String() string {
	return "forwarded-req-mount-accessor"
}

type dynamicSystemView struct {
	core       *Core
	mountEntry *MountEntry
}

type extendedSystemView struct {
	dynamicSystemView
}

func (e extendedSystemView) Auditor() logical.Auditor {
	return genericAuditor{
		mountType: e.mountEntry.Type,
		namespace: e.mountEntry.Namespace(),
		c:         e.core,
	}
}

func (e extendedSystemView) ForwardGenericRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// Forward the request if allowed
	if couldForward(e.core) {
		ctx = namespace.ContextWithNamespace(ctx, e.mountEntry.Namespace())
		ctx = context.WithValue(ctx, ctxKeyForwardedRequestMountAccessor{}, e.mountEntry.Accessor)
		return forward(ctx, e.core, req)
	}

	return nil, logical.ErrReadOnly
}

func (d dynamicSystemView) DefaultLeaseTTL() time.Duration {
	def, _ := d.fetchTTLs()
	return def
}

func (d dynamicSystemView) MaxLeaseTTL() time.Duration {
	_, max := d.fetchTTLs()
	return max
}

func (d dynamicSystemView) SudoPrivilege(ctx context.Context, path string, token string) bool {
	// Resolve the token policy
	te, err := d.core.tokenStore.Lookup(ctx, token)
	if err != nil {
		d.core.logger.Error("failed to lookup token", "error", err)
		return false
	}

	// Ensure the token is valid
	if te == nil {
		d.core.logger.Error("entry not found for given token")
		return false
	}

	policies := make(map[string][]string)
	// Add token policies
	policies[te.NamespaceID] = append(policies[te.NamespaceID], te.Policies...)

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, d.core)
	if err != nil {
		d.core.logger.Error("failed to lookup token namespace", "error", err)
		return false
	}
	if tokenNS == nil {
		d.core.logger.Error("failed to lookup token namespace", "error", namespace.ErrNoNamespace)
		return false
	}

	// Add identity policies from all the namespaces
	entity, identityPolicies, err := d.core.fetchEntityAndDerivedPolicies(ctx, tokenNS, te.EntityID)
	if err != nil {
		d.core.logger.Error("failed to fetch identity policies", "error", err)
		return false
	}
	for nsID, nsPolicies := range identityPolicies {
		policies[nsID] = append(policies[nsID], nsPolicies...)
	}

	tokenCtx := namespace.ContextWithNamespace(ctx, tokenNS)

	// Construct the corresponding ACL object. Derive and use a new context that
	// uses the req.ClientToken's namespace
	acl, err := d.core.policyStore.ACL(tokenCtx, entity, policies)
	if err != nil {
		d.core.logger.Error("failed to retrieve ACL for token's policies", "token_policies", te.Policies, "error", err)
		return false
	}

	// The operation type isn't important here as this is run from a path the
	// user has already been given access to; we only care about whether they
	// have sudo
	req := new(logical.Request)
	req.Operation = logical.ReadOperation
	req.Path = path
	authResults := acl.AllowOperation(ctx, req, true)
	return authResults.RootPrivs
}

// TTLsByPath returns the default and max TTLs corresponding to a particular
// mount point, or the system default
func (d dynamicSystemView) fetchTTLs() (def, max time.Duration) {
	def = d.core.defaultLeaseTTL
	max = d.core.maxLeaseTTL

	if d.mountEntry != nil {
		if d.mountEntry.Config.DefaultLeaseTTL != 0 {
			def = d.mountEntry.Config.DefaultLeaseTTL
		}
		if d.mountEntry.Config.MaxLeaseTTL != 0 {
			max = d.mountEntry.Config.MaxLeaseTTL
		}
	}

	return
}

// Tainted indicates that the mount is in the process of being removed
func (d dynamicSystemView) Tainted() bool {
	return d.mountEntry.Tainted
}

// CachingDisabled indicates whether to use caching behavior
func (d dynamicSystemView) CachingDisabled() bool {
	return d.core.cachingDisabled || (d.mountEntry != nil && d.mountEntry.Config.ForceNoCache)
}

func (d dynamicSystemView) LocalMount() bool {
	return d.mountEntry != nil && d.mountEntry.Local
}

// Checks if this is a primary Vault instance. Caller should hold the stateLock
// in read mode.
func (d dynamicSystemView) ReplicationState() consts.ReplicationState {
	state := d.core.ReplicationState()
	if d.core.perfStandby {
		state |= consts.ReplicationPerformanceStandby
	}
	return state
}

func (d dynamicSystemView) HasFeature(feature license.Features) bool {
	return d.core.HasFeature(feature)
}

// ResponseWrapData wraps the given data in a cubbyhole and returns the
// token used to unwrap.
func (d dynamicSystemView) ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "sys/wrapping/wrap",
	}

	resp := &logical.Response{
		WrapInfo: &wrapping.ResponseWrapInfo{
			TTL: ttl,
		},
		Data: data,
	}

	if jwt {
		resp.WrapInfo.Format = "jwt"
	}

	_, err := d.core.wrapInCubbyhole(ctx, req, resp, nil)
	if err != nil {
		return nil, err
	}

	return resp.WrapInfo, nil
}

// LookupPlugin looks for a plugin with the given name in the plugin catalog. It
// returns a PluginRunner or an error if no plugin was found.
func (d dynamicSystemView) LookupPlugin(ctx context.Context, name string, pluginType consts.PluginType) (*pluginutil.PluginRunner, error) {
	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.pluginCatalog == nil {
		return nil, fmt.Errorf("system view core plugin catalog is nil")
	}
	r, err := d.core.pluginCatalog.Get(ctx, name, pluginType)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("{{err}}: %s", name), ErrPluginNotFound)
	}

	return r, nil
}

// MlockEnabled returns the configuration setting for enabling mlock on plugins.
func (d dynamicSystemView) MlockEnabled() bool {
	return d.core.enableMlock
}

func (d dynamicSystemView) EntityInfo(entityID string) (*logical.Entity, error) {
	// Requests from token created from the token backend will not have entity information.
	// Return missing entity instead of error when requesting from MemDB.
	if entityID == "" {
		return nil, nil
	}

	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.identityStore == nil {
		return nil, fmt.Errorf("system view identity store is nil")
	}

	// Retrieve the entity from MemDB
	entity, err := d.core.identityStore.MemDBEntityByID(entityID, false)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	// Return a subset of the data
	ret := &logical.Entity{
		ID:       entity.ID,
		Name:     entity.Name,
		Disabled: entity.Disabled,
	}

	if entity.Metadata != nil {
		ret.Metadata = make(map[string]string, len(entity.Metadata))
		for k, v := range entity.Metadata {
			ret.Metadata[k] = v
		}
	}

	aliases := make([]*logical.Alias, len(entity.Aliases))
	for i, a := range entity.Aliases {
		alias := &logical.Alias{
			MountAccessor: a.MountAccessor,
			Name:          a.Name,
		}
		// MountType is not stored with the entity and must be looked up
		if mount := d.core.router.validateMountByAccessor(a.MountAccessor); mount != nil {
			alias.MountType = mount.MountType
		}

		if a.Metadata != nil {
			alias.Metadata = make(map[string]string, len(a.Metadata))
			for k, v := range a.Metadata {
				alias.Metadata[k] = v
			}
		}

		aliases[i] = alias
	}
	ret.Aliases = aliases

	return ret, nil
}

func (d dynamicSystemView) PluginEnv(_ context.Context) (*logical.PluginEnvironment, error) {
	return &logical.PluginEnvironment{
		VaultVersion: version.GetVersion().Version,
	}, nil
}
