package raft

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	snapshot "github.com/hashicorp/raft-snapshot"
	raftboltdb "github.com/hashicorp/vault/physical/raft/logstore"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/seal"

	"github.com/hashicorp/vault/sdk/physical"
)

// Verify RaftBackend satisfies the correct interfaces
var _ physical.Backend = (*RaftBackend)(nil)
var _ physical.Transactional = (*RaftBackend)(nil)

var (
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512

	raftState         = "raft/"
	peersFileName     = "peers.json"
	snapshotsRetained = 2

	// Set a max size of 512kb
	maxCommandSizeBytes = 512 * 1024

	// ErrCommandTooLarge is returned when the backend tries to apply a log
	// greater than the max allowed size.
	ErrCommandTooLarge = fmt.Errorf("%s: exceeds %d byte limit", physical.ErrValueTooLarge, maxCommandSizeBytes)

	restoreOpDelayDuration = 5 * time.Second
)

// RaftBackend implements the backend interfaces and uses the raft protocol to
// persist writes to the FSM.
type RaftBackend struct {
	logger log.Logger
	conf   map[string]string
	l      sync.RWMutex

	// fsm is the state store for vault's data
	fsm *FSM

	// raft is the instance of raft we will operate on.
	raft *raft.Raft

	// raftNotifyCh is used to receive updates about leadership changes
	// regarding this node.
	raftNotifyCh chan bool

	// streamLayer is the network layer used to connect the nodes in the raft
	// cluster.
	streamLayer *raftLayer

	// raftTransport is the transport layer that the raft library uses for RPC
	// communication.
	raftTransport raft.Transport

	// snapStore is our snapshot mechanism.
	snapStore raft.SnapshotStore

	// logStore is used by the raft library to store the raft logs in durable
	// storage.
	logStore raft.LogStore

	// stableStore is used by the raft library to store additional metadata in
	// durable storage.
	stableStore raft.StableStore

	// bootstrapConfig is only set when this node needs to be bootstrapped upon
	// startup.
	bootstrapConfig *raft.Configuration

	// dataDir is the location on the local filesystem that raft and FSM data
	// will be stored.
	dataDir string

	// localID is the ID for this node. This can either be configured in the
	// config file, via a file on disk, or is otherwise randomly generated.
	localID string

	// serverAddressProvider is used to map server IDs to addresses.
	serverAddressProvider raft.ServerAddressProvider
}

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Create the FSM.
	var err error
	fsm, err := NewFSM(conf, logger.Named("fsm"))
	if err != nil {
		return nil, fmt.Errorf("failed to create fsm: %v", err)
	}

	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	// Build an all in-memory setup for dev mode, otherwise prepare a full
	// disk-based setup.
	var log raft.LogStore
	var stable raft.StableStore
	var snap raft.SnapshotStore
	var devMode bool
	if devMode {
		store := raft.NewInmemStore()
		stable = store
		log = store
		snap = raft.NewInmemSnapshotStore()
	} else {
		// Create the base raft path.
		path := filepath.Join(path, raftState)
		if err := EnsurePath(path, true); err != nil {
			return nil, err
		}

		// Create the backend raft store for logs and stable storage.
		store, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft.db"))
		if err != nil {
			return nil, err
		}
		stable = store

		// Wrap the store in a LogCache to improve performance.
		cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
		if err != nil {
			return nil, err
		}
		log = cacheStore

		// Create the snapshot store.
		snapshots, err := NewBoltSnapshotStore(path, snapshotsRetained, logger.Named("snapshot"), fsm)
		if err != nil {
			return nil, err
		}
		snap = snapshots
	}

	var localID string
	{
		// Determine the local node ID
		localID = conf["node_id"]

		// If not set in the config check the "node-id" file.
		if len(localID) == 0 {
			localIDRaw, err := ioutil.ReadFile(filepath.Join(path, "node-id"))
			switch {
			case err == nil:
				if len(localIDRaw) > 0 {
					localID = string(localIDRaw)
				}
			case os.IsNotExist(err):
			default:
				return nil, err
			}
		}

		// If the file didn't exist generate a UUID and persist it to tne
		// "node-id" file.
		if len(localID) == 0 {
			id, err := uuid.GenerateUUID()
			if err != nil {
				return nil, err
			}

			if err := ioutil.WriteFile(filepath.Join(path, "node-id"), []byte(id), 0600); err != nil {
				return nil, err
			}

			localID = id
		}
	}

	return &RaftBackend{
		logger:      logger,
		fsm:         fsm,
		conf:        conf,
		logStore:    log,
		stableStore: stable,
		snapStore:   snap,
		dataDir:     path,
		localID:     localID,
	}, nil
}

// RaftServer has information about a server in the Raft configuration
type RaftServer struct {
	// NodeID is the name of the server
	NodeID string `json:"node_id"`

	// Address is the IP:port of the server, used for Raft communications
	Address string `json:"address"`

	// Leader is true if this server is the current cluster leader
	Leader bool `json:"leader"`

	// Protocol version is the raft protocol version used by the server
	ProtocolVersion string `json:"protocol_version"`

	// Voter is true if this server has a vote in the cluster. This might
	// be false if the server is staging and still coming online.
	Voter bool `json:"voter"`
}

// RaftConfigurationResponse is returned when querying for the current Raft
// configuration.
type RaftConfigurationResponse struct {
	// Servers has the list of servers in the Raft configuration.
	Servers []*RaftServer `json:"servers"`

	// Index has the Raft index of this configuration.
	Index uint64 `json:"index"`
}

// Peer defines the ID and Adress for a given member of the raft cluster.
type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

// NodeID returns the identifier of the node
func (b *RaftBackend) NodeID() string {
	return b.localID
}

// Initialized tells if raft is running or not
func (b *RaftBackend) Initialized() bool {
	b.l.RLock()
	init := b.raft != nil
	b.l.RUnlock()
	return init
}

// SetTLSKeyring is used to install a new keyring. If the active key has changed
// it will also close any network connections or streams forcing a reconnect
// with the new key.
func (b *RaftBackend) SetTLSKeyring(keyring *RaftTLSKeyring) error {
	b.l.RLock()
	err := b.streamLayer.setTLSKeyring(keyring)
	b.l.RUnlock()

	return err
}

// SetServerAddressProvider sets a the address provider for determining the raft
// node addresses. This is currently only used in tests.
func (b *RaftBackend) SetServerAddressProvider(provider raft.ServerAddressProvider) {
	b.l.Lock()
	b.serverAddressProvider = provider
	b.l.Unlock()
}

// Bootstrap prepares the given peers to be part of the raft cluster
func (b *RaftBackend) Bootstrap(ctx context.Context, peers []Peer) error {
	b.l.Lock()
	defer b.l.Unlock()

	hasState, err := raft.HasExistingState(b.logStore, b.stableStore, b.snapStore)
	if err != nil {
		return err
	}

	if hasState {
		return errors.New("error bootstrapping cluster: cluster already has state")
	}

	raftConfig := &raft.Configuration{
		Servers: make([]raft.Server, len(peers)),
	}

	for i, p := range peers {
		raftConfig.Servers[i] = raft.Server{
			ID:      raft.ServerID(p.ID),
			Address: raft.ServerAddress(p.Address),
		}
	}

	// Store the config for later use
	b.bootstrapConfig = raftConfig
	return nil
}

// SetRestoreCallback sets the callback to be used when a restoreCallbackOp is
// processed through the FSM.
func (b *RaftBackend) SetRestoreCallback(restoreCb restoreCallback) {
	b.fsm.l.Lock()
	b.fsm.restoreCb = restoreCb
	b.fsm.l.Unlock()
}

func (b *RaftBackend) applyConfigSettings(config *raft.Config) error {
	config.Logger = b.logger
	multiplierRaw, ok := b.conf["performance_multiplier"]
	multiplier := 5
	if ok {
		var err error
		multiplier, err = strconv.Atoi(multiplierRaw)
		if err != nil {
			return err
		}
	}
	config.ElectionTimeout = config.ElectionTimeout * time.Duration(multiplier)
	config.HeartbeatTimeout = config.HeartbeatTimeout * time.Duration(multiplier)
	config.LeaderLeaseTimeout = config.LeaderLeaseTimeout * time.Duration(multiplier)

	snapThresholdRaw, ok := b.conf["snapshot_threshold"]
	if ok {
		var err error
		snapThreshold, err := strconv.Atoi(snapThresholdRaw)
		if err != nil {
			return err
		}
		config.SnapshotThreshold = uint64(snapThreshold)
	}

	trailingLogsRaw, ok := b.conf["trailing_logs"]
	if ok {
		var err error
		trailingLogs, err := strconv.Atoi(trailingLogsRaw)
		if err != nil {
			return err
		}
		config.TrailingLogs = uint64(trailingLogs)
	}

	return nil
}

// SetupCluster starts the raft cluster and enables the networking needed for
// the raft nodes to communicate.
func (b *RaftBackend) SetupCluster(ctx context.Context, raftTLSKeyring *RaftTLSKeyring, clusterListener cluster.ClusterHook) error {
	b.logger.Trace("setting up raft cluster")

	b.l.Lock()
	defer b.l.Unlock()

	// We are already unsealed
	if b.raft != nil {
		b.logger.Debug("raft already started, not setting up cluster")
		return nil
	}

	if len(b.localID) == 0 {
		return errors.New("no local node id configured")
	}

	// Setup the raft config
	raftConfig := raft.DefaultConfig()
	if err := b.applyConfigSettings(raftConfig); err != nil {
		return err
	}

	switch {
	case raftTLSKeyring == nil && clusterListener == nil:
		// If we don't have a provided network we use an in-memory one.
		// This allows us to bootstrap a node without bringing up a cluster
		// network. This will be true during bootstrap, tests and dev modes.
		_, b.raftTransport = raft.NewInmemTransportWithTimeout(raft.ServerAddress(b.localID), time.Second)
	case raftTLSKeyring == nil:
		return errors.New("no keyring provided")
	case clusterListener == nil:
		return errors.New("no cluster listener provided")
	default:
		// Load the base TLS config from the cluster listener.
		baseTLSConfig, err := clusterListener.TLSConfig(ctx)
		if err != nil {
			return err
		}

		// Set the local address and localID in the streaming layer and the raft config.
		streamLayer, err := NewRaftLayer(b.logger.Named("stream"), raftTLSKeyring, clusterListener.Addr(), baseTLSConfig)
		if err != nil {
			return err
		}
		transConfig := &raft.NetworkTransportConfig{
			Stream:                streamLayer,
			MaxPool:               3,
			Timeout:               10 * time.Second,
			ServerAddressProvider: b.serverAddressProvider,
		}
		transport := raft.NewNetworkTransportWithConfig(transConfig)

		b.streamLayer = streamLayer
		b.raftTransport = transport
	}

	raftConfig.LocalID = raft.ServerID(b.localID)

	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = raftNotifyCh

	// If we have a bootstrapConfig set we should bootstrap now.
	if b.bootstrapConfig != nil {
		bootstrapConfig := b.bootstrapConfig
		// Unset the bootstrap config
		b.bootstrapConfig = nil

		// Bootstrap raft with our known cluster members.
		if err := raft.BootstrapCluster(raftConfig, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *bootstrapConfig); err != nil {
			return err
		}
		// If we are the only node we should start as the leader.
		if len(bootstrapConfig.Servers) == 1 {
			raftConfig.StartAsLeader = true
		}
	}

	// Setup the Raft store.
	b.fsm.SetNoopRestore(true)

	raftPath := filepath.Join(b.dataDir, raftState)
	peersFile := filepath.Join(raftPath, peersFileName)
	_, err := os.Stat(peersFile)
	if err == nil {
		b.logger.Info("raft recovery initiated", "recovery_file", peersFileName)

		recoveryConfig, err := raft.ReadConfigJSON(peersFile)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed to parse peers.json: {{err}}", err)
		}

		b.logger.Info("raft recovery: found new config", "config", recoveryConfig)
		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, recoveryConfig)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed: {{err}}", err)
		}

		err = os.Remove(peersFile)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed to delete peers.json; please delete manually: {{err}}", err)
		}
		b.logger.Info("raft recovery deleted peers.json")
	}

	raftObj, err := raft.NewRaft(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport)
	b.fsm.SetNoopRestore(false)
	if err != nil {
		return err
	}
	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	if b.streamLayer != nil {
		// Add Handler to the cluster.
		clusterListener.AddHandler(consts.RaftStorageALPN, b.streamLayer)

		// Add Client to the cluster.
		clusterListener.AddClient(consts.RaftStorageALPN, b.streamLayer)
	}

	return nil
}

// TeardownCluster shuts down the raft cluster
func (b *RaftBackend) TeardownCluster(clusterListener cluster.ClusterHook) error {
	if clusterListener != nil {
		clusterListener.StopHandler(consts.RaftStorageALPN)
		clusterListener.RemoveClient(consts.RaftStorageALPN)
	}

	b.l.Lock()
	future := b.raft.Shutdown()
	b.raft = nil
	b.l.Unlock()

	return future.Error()
}

// AppliedIndex returns the latest index applied to the FSM
func (b *RaftBackend) AppliedIndex() uint64 {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return 0
	}

	return b.raft.AppliedIndex()
}

// RemovePeer removes the given peer ID from the raft cluster. If the node is
// ourselves we will give up leadership.
func (b *RaftBackend) RemovePeer(ctx context.Context, peerID string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	future := b.raft.RemoveServer(raft.ServerID(peerID), 0, 0)

	return future.Error()
}

func (b *RaftBackend) GetConfiguration(ctx context.Context) (*RaftConfigurationResponse, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage is not initialized")
	}

	future := b.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}

	config := &RaftConfigurationResponse{
		Index: future.Index(),
	}

	for _, server := range future.Configuration().Servers {
		entry := &RaftServer{
			NodeID:          string(server.ID),
			Address:         string(server.Address),
			Leader:          server.Address == b.raft.Leader(),
			Voter:           server.Suffrage == raft.Voter,
			ProtocolVersion: strconv.Itoa(raft.ProtocolVersionMax),
		}
		config.Servers = append(config.Servers, entry)
	}

	return config, nil
}

// AddPeer adds a new server to the raft cluster
func (b *RaftBackend) AddPeer(ctx context.Context, peerID, clusterAddr string) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	b.logger.Debug("adding raft peer", "node_id", peerID, "cluster_addr", clusterAddr)

	future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)

	return future.Error()
}

// Peers returns all the servers present in the raft cluster
func (b *RaftBackend) Peers(ctx context.Context) ([]Peer, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage backend is not initialized")
	}

	future := b.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}

	ret := make([]Peer, len(future.Configuration().Servers))
	for i, s := range future.Configuration().Servers {
		ret[i] = Peer{
			ID:      string(s.ID),
			Address: string(s.Address),
		}
	}

	return ret, nil
}

// Snapshot takes a raft snapshot, packages it into a archive file and writes it
// to the provided writer. Seal access is used to encrypt the SHASUM file so we
// can validate the snapshot was taken using the same master keys or not.
func (b *RaftBackend) Snapshot(out io.Writer, access seal.Access) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage backend is sealed")
	}

	// If we have access to the seal create a sealer object
	var s snapshot.Sealer
	if access != nil {
		s = &sealer{
			access: access,
		}
	}

	snap, err := snapshot.NewWithSealer(b.logger.Named("snapshot"), b.raft, s)
	if err != nil {
		return err
	}
	defer snap.Close()

	_, err = io.Copy(out, snap)
	if err != nil {
		return err
	}

	return nil
}

// WriteSnapshotToTemp reads a snapshot archive off the provided reader,
// extracts the data and writes the snapshot to a temporary file. The seal
// access is used to decrypt the SHASUM file in the archive to ensure this
// snapshot has the same master key as the running instance. If the provided
// access is nil then it will skip that validation.
func (b *RaftBackend) WriteSnapshotToTemp(in io.ReadCloser, access seal.Access) (*os.File, func(), raft.SnapshotMeta, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	var metadata raft.SnapshotMeta
	if b.raft == nil {
		return nil, nil, metadata, errors.New("raft storage backend is sealed")
	}

	// If we have access to the seal create a sealer object
	var s snapshot.Sealer
	if access != nil {
		s = &sealer{
			access: access,
		}
	}

	snap, cleanup, err := snapshot.WriteToTempFileWithSealer(b.logger.Named("snapshot"), in, &metadata, s)
	return snap, cleanup, metadata, err
}

// RestoreSnapshot applies the provided snapshot metadata and snapshot data to
// raft.
func (b *RaftBackend) RestoreSnapshot(ctx context.Context, metadata raft.SnapshotMeta, snap io.Reader) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	if err := b.raft.Restore(&metadata, snap, 0); err != nil {
		b.logger.Named("snapshot").Error("failed to restore snapshot", "error", err)
		return err
	}

	// Apply a log that tells the follower nodes to run the restore callback
	// function. This is done after the restore call so we can be sure the
	// snapshot applied to a quorum of nodes.
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: restoreCallbackOp,
			},
		},
	}

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()

	// Do a best-effort attempt to let the standbys apply the restoreCallbackOp
	// before we continue.
	time.Sleep(restoreOpDelayDuration)
	return err
}

// Delete inserts an entry in the log to delete the given path
func (b *RaftBackend) Delete(ctx context.Context, path string) error {
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: deleteOp,
				Key:    path,
			},
		},
	}

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()
	return err
}

// Get returns the value corresponding to the given path from the fsm
func (b *RaftBackend) Get(ctx context.Context, path string) (*physical.Entry, error) {
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	return b.fsm.Get(ctx, path)
}

// Put inserts an entry in the log for the put operation
func (b *RaftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: putOp,
				Key:    entry.Key,
				Value:  entry.Value,
			},
		},
	}

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()
	return err
}

// List enumerates all the items under the prefix from the fsm
func (b *RaftBackend) List(ctx context.Context, prefix string) ([]string, error) {
	if b.fsm == nil {
		return nil, errors.New("raft: fsm not configured")
	}

	return b.fsm.List(ctx, prefix)
}

// Transaction applies all the given operations into a single log and
// applies it.
func (b *RaftBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	command := &LogData{
		Operations: make([]*LogOperation, len(txns)),
	}
	for i, txn := range txns {
		op := &LogOperation{}
		switch txn.Operation {
		case physical.PutOperation:
			op.OpType = putOp
			op.Key = txn.Entry.Key
			op.Value = txn.Entry.Value
		case physical.DeleteOperation:
			op.OpType = deleteOp
			op.Key = txn.Entry.Key
		default:
			return fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
		}

		command.Operations[i] = op
	}

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()
	return err
}

// applyLog will take a given log command and apply it to the raft log. applyLog
// doesn't return until the log has been applied to a quorum of servers and is
// persisted to the local FSM. Caller should hold the backend's read lock.
func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	if b.raft == nil {
		return errors.New("raft storage backend is not initialized")
	}

	commandBytes, err := proto.Marshal(command)
	if err != nil {
		return err
	}

	// Restrict the value to maxCommandSizeBytes in length
	if len(commandBytes) > maxCommandSizeBytes {
		return ErrCommandTooLarge
	}

	applyFuture := b.raft.Apply(commandBytes, 0)
	err = applyFuture.Error()
	if err != nil {
		return err
	}

	if resp, ok := applyFuture.Response().(*FSMApplyResponse); !ok || !resp.Success {
		return errors.New("could not apply data")
	}

	return nil
}

// HAEnabled is the implemention of the HABackend interface
func (b *RaftBackend) HAEnabled() bool { return true }

// HAEnabled is the implemention of the HABackend interface
func (b *RaftBackend) LockWith(key, value string) (physical.Lock, error) {
	return &RaftLock{
		key:   key,
		value: []byte(value),
		b:     b,
	}, nil
}

// RaftLock implements the physical Lock interface and enables HA for this
// backend. The Lock uses the raftNotifyCh for receiving leadership edge
// triggers. Vault's active duty matches raft's leadership.
type RaftLock struct {
	key   string
	value []byte

	b *RaftBackend
}

// monitorLeadership waits until we receive an update on the raftNotifyCh and
// closes the leaderLost channel.
func (l *RaftLock) monitorLeadership(stopCh <-chan struct{}, leaderNotifyCh <-chan bool) <-chan struct{} {
	leaderLost := make(chan struct{})
	go func() {
		select {
		case <-leaderNotifyCh:
			close(leaderLost)
		case <-stopCh:
		}
	}()
	return leaderLost
}

// Lock blocks until we become leader or are shutdown. It returns a channel that
// is closed when we detect a loss of leadership.
func (l *RaftLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.b.l.RLock()

	// Cache the notifyCh locally
	leaderNotifyCh := l.b.raftNotifyCh

	// Check to see if we are already leader.
	if l.b.raft.State() == raft.Leader {
		err := l.b.applyLog(context.Background(), &LogData{
			Operations: []*LogOperation{
				&LogOperation{
					OpType: putOp,
					Key:    l.key,
					Value:  l.value,
				},
			},
		})
		l.b.l.RUnlock()
		if err != nil {
			return nil, err
		}

		return l.monitorLeadership(stopCh, leaderNotifyCh), nil
	}
	l.b.l.RUnlock()

	for {
		select {
		case isLeader := <-leaderNotifyCh:
			if isLeader {
				// We are leader, set the key
				l.b.l.RLock()
				err := l.b.applyLog(context.Background(), &LogData{
					Operations: []*LogOperation{
						&LogOperation{
							OpType: putOp,
							Key:    l.key,
							Value:  l.value,
						},
					},
				})
				l.b.l.RUnlock()
				if err != nil {
					return nil, err
				}

				return l.monitorLeadership(stopCh, leaderNotifyCh), nil
			}
		case <-stopCh:
			return nil, nil
		}
	}

	return nil, nil
}

// Unlock gives up leadership.
func (l *RaftLock) Unlock() error {
	return l.b.raft.LeadershipTransfer().Error()
}

// Value reads the value of the lock. This informs us who is currently leader.
func (l *RaftLock) Value() (bool, string, error) {
	e, err := l.b.Get(context.Background(), l.key)
	if err != nil {
		return false, "", err
	}
	if e == nil {
		return false, "", nil
	}

	value := string(e.Value)
	// TODO: how to tell if held?
	return true, value, nil
}

// sealer implements the snapshot.Sealer interface and is used in the snapshot
// process for encrypting/decrypting the SHASUM file in snapshot archives.
type sealer struct {
	access seal.Access
}

// Seal encrypts the data with using the seal access object.
func (s sealer) Seal(ctx context.Context, pt []byte) ([]byte, error) {
	if s.access == nil {
		return nil, errors.New("no seal access available")
	}
	eblob, err := s.access.Encrypt(ctx, pt)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(eblob)
}

// Open decrypts the data using the seal access object.
func (s sealer) Open(ctx context.Context, ct []byte) ([]byte, error) {
	if s.access == nil {
		return nil, errors.New("no seal access available")
	}

	var eblob physical.EncryptedBlobInfo
	err := proto.Unmarshal(ct, &eblob)
	if err != nil {
		return nil, err
	}

	return s.access.Decrypt(ctx, &eblob)
}
