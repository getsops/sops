package testhelpers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/cluster"

	log "github.com/hashicorp/go-hclog"
	raftlib "github.com/hashicorp/raft"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	testing "github.com/mitchellh/go-testing-interface"
)

type ReplicatedTestClusters struct {
	PerfPrimaryCluster     *vault.TestCluster
	PerfSecondaryCluster   *vault.TestCluster
	PerfPrimaryDRCluster   *vault.TestCluster
	PerfSecondaryDRCluster *vault.TestCluster
}

func (r *ReplicatedTestClusters) Cleanup() {
	r.PerfPrimaryCluster.Cleanup()
	r.PerfSecondaryCluster.Cleanup()
	if r.PerfPrimaryDRCluster != nil {
		r.PerfPrimaryDRCluster.Cleanup()
	}
	if r.PerfSecondaryDRCluster != nil {
		r.PerfSecondaryDRCluster.Cleanup()
	}
}

func (r *ReplicatedTestClusters) Primary() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfPrimaryCluster, r.PerfPrimaryCluster.Cores[0], r.PerfPrimaryCluster.Cores[0].Client
}

func (r *ReplicatedTestClusters) Secondary() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfSecondaryCluster, r.PerfSecondaryCluster.Cores[0], r.PerfSecondaryCluster.Cores[0].Client
}

func (r *ReplicatedTestClusters) PrimaryDR() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfPrimaryDRCluster, r.PerfPrimaryDRCluster.Cores[0], r.PerfPrimaryDRCluster.Cores[0].Client
}

func (r *ReplicatedTestClusters) SecondaryDR() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfSecondaryDRCluster, r.PerfSecondaryDRCluster.Cores[0], r.PerfSecondaryDRCluster.Cores[0].Client
}

// Generates a root token on the target cluster.
func GenerateRoot(t testing.T, cluster *vault.TestCluster, drToken bool) string {
	token, err := GenerateRootWithError(t, cluster, drToken)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRootWithError(t testing.T, cluster *vault.TestCluster, drToken bool) (string, error) {
	// If recovery keys supported, use those to perform root token generation instead
	var keys [][]byte
	if cluster.Cores[0].SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}

	client := cluster.Cores[0].Client
	f := client.Sys().GenerateRootInit
	if drToken {
		f = client.Sys().GenerateDROperationTokenInit
	}
	status, err := f("", "")
	if err != nil {
		return "", err
	}

	if status.Required > len(keys) {
		return "", fmt.Errorf("need more keys than have, need %d have %d", status.Required, len(keys))
	}

	otp := status.OTP

	for i, key := range keys {
		if i >= status.Required {
			break
		}
		f := client.Sys().GenerateRootUpdate
		if drToken {
			f = client.Sys().GenerateDROperationTokenUpdate
		}
		status, err = f(base64.StdEncoding.EncodeToString(key), status.Nonce)
		if err != nil {
			return "", err
		}
	}
	if !status.Complete {
		return "", errors.New("generate root operation did not end successfully")
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(status.EncodedToken)
	if err != nil {
		return "", err
	}
	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

// RandomWithPrefix is used to generate a unique name with a prefix, for
// randomizing names in acceptance tests
func RandomWithPrefix(name string) string {
	return fmt.Sprintf("%s-%d", name, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func EnsureCoresSealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		EnsureCoreSealed(t, core)
	}
}

func EnsureCoreSealed(t testing.T, core *vault.TestClusterCore) error {
	core.Seal(t)
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			return fmt.Errorf("timeout waiting for core to seal")
		}
		if core.Core.Sealed() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	return nil
}

func EnsureCoresUnsealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		EnsureCoreUnsealed(t, c, core)
	}
}
func EnsureCoreUnsealed(t testing.T, c *vault.TestCluster, core *vault.TestClusterCore) {
	if !core.Sealed() {
		return
	}

	client := core.Client
	client.Sys().ResetUnsealProcess()
	for j := 0; j < len(c.BarrierKeys); j++ {
		statusResp, err := client.Sys().Unseal(base64.StdEncoding.EncodeToString(c.BarrierKeys[j]))
		if err != nil {
			// Sometimes when we get here it's already unsealed on its own
			// and then this fails for DR secondaries so check again
			if core.Sealed() {
				t.Fatal(err)
			}
			break
		}
		if statusResp == nil {
			t.Fatal("nil status response during unseal")
		}
		if !statusResp.Sealed {
			break
		}
	}
	if core.Sealed() {
		t.Fatal("core is still sealed")
	}
}

func EnsureCoreIsPerfStandby(t testing.T, client *api.Client) {
	t.Helper()
	start := time.Now()
	for {
		health, err := client.Sys().Health()
		if err != nil {
			t.Fatal(err)
		}
		if health.PerformanceStandby {
			break
		}

		t.Log("waiting for performance standby", fmt.Sprintf("%+v", health))
		time.Sleep(time.Millisecond * 500)
		if time.Now().After(start.Add(time.Second * 60)) {
			t.Fatal("did not become a perf standby")
		}
	}
}

func WaitForReplicationState(t testing.T, c *vault.Core, state consts.ReplicationState) {
	timeout := time.Now().Add(10 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatalf("timeout waiting for core to have state %d", uint32(state))
		}
		state := c.ReplicationState()
		if state.HasState(state) {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

type PassthroughWithLocalPaths struct {
	logical.Backend
}

func (p *PassthroughWithLocalPaths) SpecialPaths() *logical.Paths {
	return &logical.Paths{
		LocalStorage: []string{
			"*",
		},
	}
}

func PassthroughWithLocalPathsFactory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b, err := vault.PassthroughBackendFactory(ctx, c)
	if err != nil {
		return nil, err
	}

	return &PassthroughWithLocalPaths{b}, nil
}

func ConfClusterAndCore(t testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) (*vault.TestCluster, *vault.TestClusterCore) {
	if conf.Physical != nil || conf.HAPhysical != nil {
		t.Fatalf("conf.Physical and conf.HAPhysical cannot be specified")
	}
	if opts.Logger == nil {
		t.Fatalf("opts.Logger must be specified")
	}

	inm, err := inmem.NewTransactionalInmem(nil, opts.Logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, opts.Logger)
	if err != nil {
		t.Fatal(err)
	}

	coreConfig := *conf
	coreConfig.Physical = inm
	coreConfig.HAPhysical = inmha.(physical.HABackend)
	coreConfig.CredentialBackends = map[string]logical.Factory{
		"approle":  credAppRole.Factory,
		"userpass": credUserpass.Factory,
	}
	vault.AddNoopAudit(&coreConfig)
	cluster := vault.NewTestCluster(t, &coreConfig, opts)
	cluster.Start()

	cores := cluster.Cores
	core := cores[0]

	vault.TestWaitActive(t, core.Core)

	return cluster, core
}

func GetPerfReplicatedClusters(t testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	ret := &ReplicatedTestClusters{}

	var logger hclog.Logger
	if opts != nil {
		logger = opts.Logger
	}
	if logger == nil {
		logger = log.New(&log.LoggerOptions{
			Mutex: &sync.Mutex{},
			Level: log.Trace,
		})
	}

	// Set this lower so that state populates quickly to standby nodes
	cluster.HeartbeatInterval = 2 * time.Second

	numCores := opts.NumCores
	if numCores == 0 {
		numCores = vault.DefaultNumCores
	}

	localopts := *opts
	localopts.Logger = logger.Named("perf-pri")
	ret.PerfPrimaryCluster, _ = ConfClusterAndCore(t, conf, &localopts)

	localopts.Logger = logger.Named("perf-sec")
	localopts.FirstCoreNumber += numCores
	ret.PerfSecondaryCluster, _ = ConfClusterAndCore(t, conf, &localopts)

	SetupTwoClusterPerfReplication(t, ret.PerfPrimaryCluster, ret.PerfSecondaryCluster)

	return ret
}

func GetFourReplicatedClusters(t testing.T, handlerFunc func(*vault.HandlerProperties) http.Handler) *ReplicatedTestClusters {
	return GetFourReplicatedClustersWithConf(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: handlerFunc,
	})
}

func GetFourReplicatedClustersWithConf(t testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	ret := &ReplicatedTestClusters{}

	logger := log.New(&log.LoggerOptions{
		Mutex: &sync.Mutex{},
		Level: log.Trace,
	})
	// Set this lower so that state populates quickly to standby nodes
	cluster.HeartbeatInterval = 2 * time.Second

	numCores := opts.NumCores
	if numCores == 0 {
		numCores = vault.DefaultNumCores
	}

	localopts := *opts
	localopts.Logger = logger.Named("perf-pri")
	ret.PerfPrimaryCluster, _ = ConfClusterAndCore(t, conf, &localopts)

	localopts.Logger = logger.Named("perf-sec")
	localopts.FirstCoreNumber += numCores
	ret.PerfSecondaryCluster, _ = ConfClusterAndCore(t, conf, &localopts)

	localopts.Logger = logger.Named("perf-pri-dr")
	localopts.FirstCoreNumber += numCores
	ret.PerfPrimaryDRCluster, _ = ConfClusterAndCore(t, conf, &localopts)

	localopts.Logger = logger.Named("perf-sec-dr")
	localopts.FirstCoreNumber += numCores
	ret.PerfSecondaryDRCluster, _ = ConfClusterAndCore(t, conf, &localopts)

	builder := &ReplicatedTestClustersBuilder{clusters: ret}
	builder.setupFourClusterReplication(t)

	// Wait until poison pills have been read
	time.Sleep(45 * time.Second)
	EnsureCoresUnsealed(t, ret.PerfPrimaryCluster)
	EnsureCoresUnsealed(t, ret.PerfSecondaryCluster)
	EnsureCoresUnsealed(t, ret.PerfPrimaryDRCluster)
	EnsureCoresUnsealed(t, ret.PerfSecondaryDRCluster)

	return ret
}

type ReplicatedTestClustersBuilder struct {
	clusters               *ReplicatedTestClusters
	perfToken              string
	drToken                string
	perfSecondaryRootToken string
	perfSecondaryDRToken   string
}

func SetupTwoClusterPerfReplication(t testing.T, pri, sec *vault.TestCluster) {
	clusters := &ReplicatedTestClusters{
		PerfPrimaryCluster:   pri,
		PerfSecondaryCluster: sec,
	}
	builder := &ReplicatedTestClustersBuilder{clusters: clusters}
	builder.setupTwoClusterReplication(t)
}

func (r *ReplicatedTestClustersBuilder) setupTwoClusterReplication(t testing.T) {
	t.Log("enabling perf primary")
	r.enablePerfPrimary(t)
	WaitForActiveNode(t, r.clusters.PerfPrimaryCluster)
	r.getPerformanceToken(t)
	t.Log("enabling perf secondary")
	r.enablePerformanceSecondary(t)
}

func SetupFourClusterReplication(t testing.T, pri, sec, pridr, secdr *vault.TestCluster) {
	clusters := &ReplicatedTestClusters{
		PerfPrimaryCluster:     pri,
		PerfSecondaryCluster:   sec,
		PerfPrimaryDRCluster:   pridr,
		PerfSecondaryDRCluster: secdr,
	}
	builder := &ReplicatedTestClustersBuilder{clusters: clusters}
	builder.setupFourClusterReplication(t)
}

func (r *ReplicatedTestClustersBuilder) setupFourClusterReplication(t testing.T) {
	t.Log("enabling perf primary")
	r.enablePerfPrimary(t)
	r.getPerformanceToken(t)

	t.Log("enabling dr primary")
	enableDrPrimary(t, r.clusters.PerfPrimaryCluster)
	r.drToken = getDrToken(t, r.clusters.PerfPrimaryCluster, "primary-dr-secondary")
	WaitForActiveNode(t, r.clusters.PerfPrimaryCluster)
	time.Sleep(1 * time.Second)

	t.Log("enabling perf secondary")
	r.enablePerformanceSecondary(t)
	enableDrPrimary(t, r.clusters.PerfSecondaryCluster)
	r.perfSecondaryDRToken = getDrToken(t, r.clusters.PerfSecondaryCluster, "secondary-dr-secondary")

	t.Log("enabling dr secondary on primary dr cluster")
	r.enableDrSecondary(t, r.clusters.PerfPrimaryDRCluster, r.drToken, r.clusters.PerfPrimaryCluster.CACertPEMFile)
	r.clusters.PerfPrimaryDRCluster.Cores[0].Client.SetToken(r.clusters.PerfPrimaryCluster.Cores[0].Client.Token())
	WaitForActiveNode(t, r.clusters.PerfPrimaryDRCluster)
	time.Sleep(1 * time.Second)

	t.Log("enabling dr secondary on secondary dr cluster")
	r.enableDrSecondary(t, r.clusters.PerfSecondaryDRCluster, r.perfSecondaryDRToken, r.clusters.PerfSecondaryCluster.CACertPEMFile)
	r.clusters.PerfSecondaryDRCluster.Cores[0].Client.SetToken(r.perfSecondaryRootToken)
	WaitForActiveNode(t, r.clusters.PerfSecondaryDRCluster)
}

func (r *ReplicatedTestClustersBuilder) enablePerfPrimary(t testing.T) {
	c := r.clusters.PerfPrimaryCluster.Cores[0]
	_, err := c.Client.Logical().Write("sys/replication/performance/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, c.Core, consts.ReplicationPerformancePrimary)
}

func (r *ReplicatedTestClustersBuilder) getPerformanceToken(t testing.T) {
	client := r.clusters.PerfPrimaryCluster.Cores[0].Client
	req := map[string]interface{}{
		"id": "perf-secondary",
	}
	secret, err := client.Logical().Write("sys/replication/performance/primary/secondary-token", req)
	if err != nil {
		t.Fatal(err)
	}
	r.perfToken = secret.WrapInfo.Token
}

func enableDrPrimary(t testing.T, tc *vault.TestCluster) {
	c := tc.Cores[0]
	_, err := c.Client.Logical().Write("sys/replication/dr/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, c.Core, consts.ReplicationDRPrimary)
}

func getDrToken(t testing.T, tc *vault.TestCluster, id string) string {
	req := map[string]interface{}{
		"id": id,
	}
	secret, err := tc.Cores[0].Client.Logical().Write("sys/replication/dr/primary/secondary-token", req)
	if err != nil {
		t.Fatal(err)
	}
	return secret.WrapInfo.Token
}

func (r *ReplicatedTestClustersBuilder) enablePerformanceSecondary(t testing.T) {
	c := r.clusters.PerfSecondaryCluster.Cores[0]
	postData := map[string]interface{}{
		"token":   r.perfToken,
		"ca_file": r.clusters.PerfPrimaryCluster.CACertPEMFile,
	}
	if r.clusters.PerfPrimaryCluster.ClientAuthRequired {
		p := r.clusters.PerfPrimaryCluster.Cores[0]
		postData["client_cert_pem"] = string(p.ServerCertPEM)
		postData["client_key_pem"] = string(p.ServerKeyPEM)
	}
	_, err := c.Client.Logical().Write("sys/replication/performance/secondary/enable", postData)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, c.Core, consts.ReplicationPerformanceSecondary)

	r.clusters.PerfSecondaryCluster.BarrierKeys = r.clusters.PerfPrimaryCluster.BarrierKeys

	// We want to make sure we unseal all the nodes so we first need to wait
	// until two of the nodes seal due to the poison pill being written
	WaitForNCoresSealed(t, r.clusters.PerfSecondaryCluster, 2)
	EnsureCoresUnsealed(t, r.clusters.PerfSecondaryCluster)
	WaitForActiveNode(t, r.clusters.PerfSecondaryCluster)

	r.perfSecondaryRootToken = GenerateRoot(t, r.clusters.PerfSecondaryCluster, false)
	for _, core := range r.clusters.PerfSecondaryCluster.Cores {
		core.Client.SetToken(r.perfSecondaryRootToken)
	}
}

func (r *ReplicatedTestClustersBuilder) enableDrSecondary(t testing.T, tc *vault.TestCluster, token, ca_file string) {
	_, err := tc.Cores[0].Client.Logical().Write("sys/replication/dr/secondary/enable", map[string]interface{}{
		"token":   token,
		"ca_file": ca_file,
	})
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, tc.Cores[0].Core, consts.ReplicationDRSecondary)
	tc.BarrierKeys = r.clusters.PerfPrimaryCluster.BarrierKeys

	// We want to make sure we unseal all the nodes so we first need to wait
	// until two of the nodes seal due to the poison pill being written
	WaitForNCoresSealed(t, tc, 2)
	EnsureCoresUnsealed(t, tc)
}

func EnsureStableActiveNode(t testing.T, cluster *vault.TestCluster) {
	activeCore := DeriveActiveCore(t, cluster)

	for i := 0; i < 30; i++ {
		leaderResp, err := activeCore.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			t.Fatal("unstable active node")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func DeriveActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	for i := 0; i < 10; i++ {
		for _, core := range cluster.Cores {
			leaderResp, err := core.Client.Sys().Leader()
			if err != nil {
				t.Fatal(err)
			}
			if leaderResp.IsSelf {
				return core
			}
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("could not derive the active core")
	return nil
}

func DeriveStandbyCores(t testing.T, cluster *vault.TestCluster) []*vault.TestClusterCore {
	cores := make([]*vault.TestClusterCore, 0, 2)
	for _, core := range cluster.Cores {
		leaderResp, err := core.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			cores = append(cores, core)
		}
	}

	return cores
}

func WaitForNCoresUnsealed(t testing.T, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 30; i++ {
		unsealed := 0
		for _, core := range cluster.Cores {
			if !core.Core.Sealed() {
				unsealed++
			}
		}

		if unsealed >= n {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("%d cores were not sealed", n)
}

func WaitForNCoresSealed(t testing.T, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 30; i++ {
		sealed := 0
		for _, core := range cluster.Cores {
			if core.Core.Sealed() {
				sealed++
			}
		}

		if sealed >= n {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("%d cores were not sealed", n)
}

func WaitForActiveNode(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	for i := 0; i < 30; i++ {
		for _, core := range cluster.Cores {
			if standby, _ := core.Core.Standby(); !standby {
				return core
			}
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not become active")
	return nil
}

func WaitForMatchingMerkleRoots(t testing.T, endpoint string, primary, secondary *api.Client) {
	getRoot := func(mode string, cli *api.Client) string {
		status, err := cli.Logical().Read(endpoint + "status")
		if err != nil {
			t.Fatal(err)
		}
		if status == nil || status.Data == nil || status.Data["mode"] == nil {
			t.Fatal("got nil secret or data")
		}
		if status.Data["mode"].(string) != mode {
			t.Fatalf("expected mode=%s, got %s", mode, status.Data["mode"].(string))
		}
		return status.Data["merkle_root"].(string)
	}

	t.Helper()
	for i := 0; i < 30; i++ {
		secRoot := getRoot("secondary", secondary)
		priRoot := getRoot("primary", primary)

		if reflect.DeepEqual(priRoot, secRoot) {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("roots did not become equal")
}

func WaitForWAL(t testing.T, c *vault.TestClusterCore, wal uint64) {
	timeout := time.Now().Add(3 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("timeout waiting for WAL")
		}
		if vault.LastRemoteWAL(c.Core) >= wal {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func RekeyCluster(t testing.T, cluster *vault.TestCluster) {
	client := cluster.Cores[0].Client

	init, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	var statusResp *api.RekeyUpdateResponse
	for j := 0; j < len(cluster.BarrierKeys); j++ {
		statusResp, err = client.Sys().RekeyUpdate(base64.StdEncoding.EncodeToString(cluster.BarrierKeys[j]), init.Nonce)
		if err != nil {
			t.Fatal(err)
		}
		if statusResp == nil {
			t.Fatal("nil status response during unseal")
		}
		if statusResp.Complete {
			break
		}
	}

	if len(statusResp.KeysB64) != 5 {
		t.Fatal("wrong number of keys")
	}

	newBarrierKeys := make([][]byte, 5)
	for i, key := range statusResp.KeysB64 {
		newBarrierKeys[i], err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	cluster.BarrierKeys = newBarrierKeys
}

func CreateRaftBackend(t testing.T, logger hclog.Logger, nodeID string) (physical.Backend, func(), error) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)
	cleanupFunc := func() {
		os.RemoveAll(raftDir)
	}

	logger.Info("raft dir", "dir", raftDir)

	conf := map[string]string{
		"path":                   raftDir,
		"node_id":                nodeID,
		"performance_multiplier": "8",
	}

	backend, err := raft.NewRaftBackend(conf, logger)
	if err != nil {
		cleanupFunc()
		t.Fatal(err)
	}

	return backend, cleanupFunc, nil
}

type TestRaftServerAddressProvider struct {
	Cluster *vault.TestCluster
}

func (p *TestRaftServerAddressProvider) ServerAddr(id raftlib.ServerID) (raftlib.ServerAddress, error) {
	for _, core := range p.Cluster.Cores {
		if core.NodeID == string(id) {
			parsed, err := url.Parse(core.ClusterAddr())
			if err != nil {
				return "", err
			}

			return raftlib.ServerAddress(parsed.Host), nil
		}
	}

	return "", errors.New("could not find cluster addr")
}

func RaftClusterJoinNodes(t testing.T, cluster *vault.TestCluster) {
	addressProvider := &TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()
	atomic.StoreUint32(&vault.UpdateClusterAddrForTests, 1)

	// Seal the leader so we can install an address provider
	{
		EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	// Join core1
	{
		core := cluster.Cores[1]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderAPI, leaderCore.TLSConfig, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	// Join core2
	{
		core := cluster.Cores[2]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderAPI, leaderCore.TLSConfig, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	WaitForNCoresUnsealed(t, cluster, 3)
}
