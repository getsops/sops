package vault

import (
	"regexp"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// Storage prefixes
	entityPrefix = "entity/"
)

var (
	// metaKeyFormatRegEx checks if a metadata key string is valid
	metaKeyFormatRegEx = regexp.MustCompile(`^[a-zA-Z0-9=/+_-]+$`).MatchString
)

const (
	// The meta key prefix reserved for Vault's internal use
	metaKeyReservedPrefix = "vault-"

	// The maximum number of metadata key pairs allowed to be registered
	metaMaxKeyPairs = 64

	// The maximum allowed length of a metadata key
	metaKeyMaxLength = 128

	// The maximum allowed length of a metadata value
	metaValueMaxLength = 512
)

// IdentityStore is composed of its own storage view and a MemDB which
// maintains active in-memory replicas of the storage contents indexed by
// multiple fields.
type IdentityStore struct {
	// IdentityStore is a secret backend in Vault
	*framework.Backend

	// view is the storage sub-view where all the artifacts of identity store
	// gets persisted
	view logical.Storage

	// db is the in-memory database where the storage artifacts gets replicated
	// to enable richer queries based on multiple indexes.
	db *memdb.MemDB

	// locks to make sure things are consistent
	lock     sync.RWMutex
	oidcLock sync.RWMutex

	// groupLock is used to protect modifications to group entries
	groupLock sync.RWMutex

	// oidcCache stores common response data as well as when the periodic func needs
	// to run. This is conservatively managed, and most writes to the OIDC endpoints
	// will invalidate the cache.
	oidcCache *oidcCache

	// logger is the server logger copied over from core
	logger log.Logger

	// entityPacker is used to pack multiple entity storage entries into 256
	// buckets
	entityPacker *storagepacker.StoragePacker

	// groupPacker is used to pack multiple group storage entries into 256
	// buckets
	groupPacker *storagepacker.StoragePacker

	// core is the pointer to Vault's core
	core *Core

	// disableLowerCaseNames indicates whether or not identity artifacts are
	// operated case insensitively
	disableLowerCasedNames bool
}

type groupDiff struct {
	New        []*identity.Group
	Deleted    []*identity.Group
	Unmodified []*identity.Group
}
