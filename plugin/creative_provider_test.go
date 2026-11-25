package plugin_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"filippo.io/age"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keyservice"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
)

func TestAgePluginUsesPublicRecipientAndLocalIdentity(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/agefixture")
	setAllowlist(t, binary, "age-test")
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)
	identityPath := filepath.Join(t.TempDir(), "identity.txt")
	require.NoError(t, os.WriteFile(identityPath, []byte(identity.String()+"\n"), 0600))

	key, err := sopsplugin.NewMasterKey("age-test", map[string]interface{}{
		"recipient":     identity.Recipient().String(),
		"identity_file": identityPath,
	})
	require.NoError(t, err)
	metadata := roundTripPluginDataKey(t, key, []byte("01234567890123456789012345678901"))
	require.Equal(t, identity.Recipient().String(), key.ProviderMetadata["recipient"])
	require.NotNil(t, metadata.DataKey)

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), identity.String())
}

func TestOpenPGPPluginUsesPureGoKeyrings(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/pgpfixture")
	setAllowlist(t, binary, "pgp-test")
	entity, err := openpgp.NewEntity("SOPS plugin test", "", "plugin@example.invalid", nil)
	require.NoError(t, err)
	var publicKey, privateKey bytes.Buffer
	require.NoError(t, entity.Serialize(&publicKey))
	require.NoError(t, entity.SerializePrivate(&privateKey, nil))
	keyDir := t.TempDir()
	publicPath := filepath.Join(keyDir, "public.pgp")
	privatePath := filepath.Join(keyDir, "private.pgp")
	require.NoError(t, os.WriteFile(publicPath, publicKey.Bytes(), 0600))
	require.NoError(t, os.WriteFile(privatePath, privateKey.Bytes(), 0600))

	key, err := sopsplugin.NewMasterKey("pgp-test", map[string]interface{}{
		"key_id":           "openpgp-test-key",
		"public_key_file":  publicPath,
		"private_key_file": privatePath,
	})
	require.NoError(t, err)
	roundTripPluginDataKey(t, key, []byte("01234567890123456789012345678901"))
	require.Equal(t, "openpgp-test-key", key.ProviderMetadata["key_id"])

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), privateKey.String())
}

func TestPassphrasePluginUsesScryptSaltMetadata(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/passwordfixture")
	setAllowlist(t, binary, "passphrase-test")
	t.Setenv("SOPS_TEST_PLUGIN_PASSWORD", "a long test-only passphrase")
	key, err := sopsplugin.NewMasterKey("passphrase-test", map[string]interface{}{"key_id": "recovery-code"})
	require.NoError(t, err)
	roundTripPluginDataKey(t, key, []byte("01234567890123456789012345678901"))
	require.Equal(t, "scrypt", key.ProviderMetadata["kdf"])
	require.NotEmpty(t, key.ProviderMetadata["salt"])

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "a long test-only passphrase")
	t.Setenv("SOPS_TEST_PLUGIN_PASSWORD", "wrong passphrase")
	_, err = key.Decrypt()
	require.ErrorContains(t, err, "decrypt with passphrase")
}

func TestVersionedKeyringPluginRetainsOldVersionsAfterRewrap(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/keyringfixture")
	setAllowlist(t, binary, "keyring-test")
	keyringDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(keyringDir, "current"), []byte("v1\n"), 0600))
	firstKey := []byte("0123456789abcdef0123456789abcdef")
	secondKey := []byte("abcdef0123456789abcdef0123456789")
	require.NoError(t, os.WriteFile(filepath.Join(keyringDir, "v1.key"), firstKey, 0600))
	require.NoError(t, os.WriteFile(filepath.Join(keyringDir, "v2.key"), secondKey, 0600))
	key, err := sopsplugin.NewMasterKey("keyring-test", map[string]interface{}{"keyring_dir": keyringDir})
	require.NoError(t, err)
	metadata := roundTripPluginDataKey(t, key, []byte("01234567890123456789012345678901"))
	oldCiphertext := append([]byte(nil), key.EncryptedKey...)
	oldMetadata := map[string]interface{}{"key_version": key.ProviderMetadata["key_version"]}
	require.Equal(t, "v1", key.ProviderMetadata["key_version"])

	require.NoError(t, os.WriteFile(filepath.Join(keyringDir, "current"), []byte("v2\n"), 0600))
	secondDataKey := []byte("abcdefghijabcdefghijabcdefghij12")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(secondDataKey, services))
	require.Equal(t, "v2", key.ProviderMetadata["key_version"])
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, secondDataKey, decoded)

	oldDataKey, err := sopsplugin.UnwrapKey(context.Background(), "keyring-test", sopsplugin.UnwrapKeyRequest{
		Config:           key.Config,
		Metadata:         oldMetadata,
		EncryptedDataKey: oldCiphertext,
	})
	require.NoError(t, err)
	require.Equal(t, []byte("01234567890123456789012345678901"), oldDataKey.DataKey)
}

func TestFailoverKMSPluginRecordsAndUsesSuccessfulRegion(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/failoverfixture")
	setAllowlist(t, binary, "failover-test")
	primary := &fakeFailoverRegion{failWrap: true, key: bytes.Repeat([]byte{0x51}, 32)}
	secondary := &fakeFailoverRegion{key: bytes.Repeat([]byte{0x52}, 32)}
	primaryServer := httptest.NewServer(primary)
	secondaryServer := httptest.NewServer(secondary)
	t.Cleanup(primaryServer.Close)
	t.Cleanup(secondaryServer.Close)
	key, err := sopsplugin.NewMasterKey("failover-test", map[string]interface{}{
		"key_id":             "regional-key",
		"primary_endpoint":   primaryServer.URL,
		"secondary_endpoint": secondaryServer.URL,
	})
	require.NoError(t, err)
	roundTripPluginDataKey(t, key, []byte("01234567890123456789012345678901"))
	require.Equal(t, "secondary", key.ProviderMetadata["region"])
	require.Equal(t, 1, primary.wrapCount())
	require.Equal(t, 1, secondary.wrapCount())
	require.Equal(t, 0, primary.unwrapCount())
	require.Equal(t, 1, secondary.unwrapCount())
}

func TestThresholdPluginRecoversWithQuorumAndFailsBelowIt(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/thresholdfixture")
	setAllowlist(t, binary, "quorum-test")
	shareStores := []*fakeShareStore{{online: true}, {online: true}, {online: true}}
	servers := make([]*httptest.Server, len(shareStores))
	endpoints := make([]interface{}, len(shareStores))
	for index, store := range shareStores {
		servers[index] = httptest.NewServer(store)
		t.Cleanup(servers[index].Close)
		endpoints[index] = servers[index].URL
	}
	key, err := sopsplugin.NewMasterKey("quorum-test", map[string]interface{}{
		"share_endpoints": endpoints,
		"threshold":       2,
		"generation":      json.Number("9007199254740993"),
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.Equal(t, json.Number("9007199254740993"), key.ProviderMetadata["generation"])

	shareStores[0].setOnline(false)
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)

	shareStores[1].setOnline(false)
	_, err = key.Decrypt()
	require.ErrorContains(t, err, "need 2 shares, retrieved 1")
}

func TestAuditPluginLogsOnlyOperationalMetadata(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/auditfixture")
	setAllowlist(t, binary, "audit-test")
	keyFile := filepath.Join(t.TempDir(), "audit.key")
	require.NoError(t, os.WriteFile(keyFile, bytes.Repeat([]byte{0x61}, 32), 0600))
	audit := &fakeAuditCollector{}
	auditServer := httptest.NewServer(audit)
	t.Cleanup(auditServer.Close)
	key, err := sopsplugin.NewMasterKey("audit-test", map[string]interface{}{
		"key_file":       keyFile,
		"key_id":         "audit-key-1",
		"audit_endpoint": auditServer.URL,
	})
	require.NoError(t, err)
	dataKey := []byte("01234567890123456789012345678901")
	roundTripPluginDataKey(t, key, dataKey)

	events := audit.eventBodies()
	require.Len(t, events, 2)
	require.Contains(t, events[0], `"operation":"wrap"`)
	require.Contains(t, events[1], `"operation":"unwrap"`)
	require.NotContains(t, strings.Join(events, "\n"), base64.StdEncoding.EncodeToString(dataKey))
	require.NotContains(t, strings.Join(events, "\n"), string(dataKey))
}

func TestMultiTenantPluginIsolatesKeysByConfiguration(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/multitenantfixture")
	setAllowlist(t, binary, "tenant-test")
	tenantSecret := "test-only-root-secret-for-tenant-derivation"
	t.Setenv("SOPS_TEST_TENANT_MASTER_SECRET", tenantSecret)
	first, err := sopsplugin.NewMasterKey("tenant-test", map[string]interface{}{
		"tenant_id": "tenant-north",
		"key_id":    "kms-key",
	})
	require.NoError(t, err)
	second, err := sopsplugin.NewMasterKey("tenant-test", map[string]interface{}{
		"tenant_id": "tenant-south",
		"key_id":    "kms-key",
	})
	require.NoError(t, err)
	metadata := sops.Metadata{
		ShamirThreshold: 2,
		KeyGroups:       []sops.KeyGroup{{first}, {second}},
	}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.Equal(t, "tenant-north", first.ProviderMetadata["tenant_id"])
	require.Equal(t, "tenant-south", second.ProviderMetadata["tenant_id"])
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)

	for _, key := range []*sopsplugin.MasterKey{first, second} {
		encoded, err := json.Marshal(key.ToMap())
		require.NoError(t, err)
		require.NotContains(t, string(encoded), tenantSecret)
	}
}

func TestApprovalPluginRequiresPolicyApprovalForUnwrap(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/approvalfixture")
	setAllowlist(t, binary, "approval-test")
	approval := &fakeApprovalPolicy{token: "test-approval-token"}
	server := httptest.NewServer(approval)
	t.Cleanup(server.Close)
	t.Setenv("SOPS_TEST_APPROVAL_KEY", "0123456789abcdef0123456789abcdef")
	t.Setenv("SOPS_TEST_APPROVAL_TOKEN", "test-approval-token")
	key, err := sopsplugin.NewMasterKey("approval-test", map[string]interface{}{
		"key_id":            "production-key",
		"approval_endpoint": server.URL,
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	_, err = key.Decrypt()
	require.ErrorContains(t, err, "unwrap approval denied")

	approval.setApproved(true)
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
	require.Equal(t, 2, approval.requestCount())
}

func TestOIDCPluginExchangesEphemeralCredentialsForKMS(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/oidcfixture")
	setAllowlist(t, binary, "oidc-test")
	service := &fakeOIDCKMS{
		clientSecret: "test-only-oidc-client-secret",
		token:        "short-lived-test-access-token",
		key:          bytes.Repeat([]byte{0x71}, 32),
	}
	server := httptest.NewServer(service)
	t.Cleanup(server.Close)
	t.Setenv("SOPS_TEST_OIDC_CLIENT_SECRET", service.clientSecret)
	key, err := sopsplugin.NewMasterKey("oidc-test", map[string]interface{}{
		"issuer_endpoint": server.URL,
		"kms_endpoint":    server.URL + "/kms",
		"subject":         "sops-test-workload",
		"key_id":          "oidc-kms-key",
	})
	require.NoError(t, err)
	roundTripPluginDataKey(t, key, []byte("01234567890123456789012345678901"))
	require.Equal(t, "oidc-kms-key", key.ProviderMetadata["key_id"])
	require.Equal(t, 2, service.tokenExchangeCount())
	require.Equal(t, 1, service.wrapCount())
	require.Equal(t, 1, service.unwrapCount())

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), service.clientSecret)
	require.NotContains(t, string(encoded), service.token)
}

func roundTripPluginDataKey(t *testing.T, key *sopsplugin.MasterKey, dataKey []byte) sops.Metadata {
	t.Helper()
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
	return metadata
}

type fakeFailoverRegion struct {
	failWrap bool
	key      []byte
	mu       sync.Mutex
	wraps    int
	unwraps  int
}

func (region *fakeFailoverRegion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/wrap":
		region.mu.Lock()
		region.wraps++
		region.mu.Unlock()
		if region.failWrap {
			http.Error(w, "region unavailable", http.StatusServiceUnavailable)
			return
		}
		var request struct {
			KeyID   string `json:"key_id"`
			DataKey string `json:"data_key_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		dataKey, err := base64.StdEncoding.DecodeString(request.DataKey)
		if err != nil {
			http.Error(w, "invalid data key", http.StatusBadRequest)
			return
		}
		gcm, err := newTestGCM(region.key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ciphertext := append(nonce, gcm.Seal(nil, nonce, dataKey, []byte(request.KeyID))...)
		writeJSONResponse(w, map[string]string{
			"ciphertext_b64": base64.StdEncoding.EncodeToString(ciphertext),
			"version":        "v1",
		})
	case "/unwrap":
		var request struct {
			KeyID         string `json:"key_id"`
			Version       string `json:"version"`
			CiphertextB64 string `json:"ciphertext_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		if request.Version != "v1" {
			http.Error(w, "unknown version", http.StatusNotFound)
			return
		}
		ciphertext, err := base64.StdEncoding.DecodeString(request.CiphertextB64)
		if err != nil {
			http.Error(w, "invalid ciphertext", http.StatusBadRequest)
			return
		}
		gcm, err := newTestGCM(region.key)
		if err != nil || len(ciphertext) < gcm.NonceSize() {
			http.Error(w, "invalid ciphertext", http.StatusBadRequest)
			return
		}
		nonce := ciphertext[:gcm.NonceSize()]
		dataKey, err := gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], []byte(request.KeyID))
		if err != nil {
			http.Error(w, "unable to unwrap", http.StatusBadRequest)
			return
		}
		region.mu.Lock()
		region.unwraps++
		region.mu.Unlock()
		writeJSONResponse(w, map[string]string{"data_key_b64": base64.StdEncoding.EncodeToString(dataKey)})
	default:
		http.NotFound(w, r)
	}
}

func (region *fakeFailoverRegion) wrapCount() int {
	region.mu.Lock()
	defer region.mu.Unlock()
	return region.wraps
}

func (region *fakeFailoverRegion) unwrapCount() int {
	region.mu.Lock()
	defer region.mu.Unlock()
	return region.unwraps
}

type fakeShareStore struct {
	mu     sync.Mutex
	online bool
	nextID int
	shares map[string][]byte
}

func (store *fakeShareStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store.mu.Lock()
	online := store.online
	store.mu.Unlock()
	if !online {
		http.Error(w, "share service unavailable", http.StatusServiceUnavailable)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/shares" {
		var request struct {
			Share string `json:"share_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		share, err := base64.StdEncoding.DecodeString(request.Share)
		if err != nil {
			http.Error(w, "invalid share", http.StatusBadRequest)
			return
		}
		store.mu.Lock()
		store.nextID++
		id := fmt.Sprintf("share-%d", store.nextID)
		if store.shares == nil {
			store.shares = make(map[string][]byte)
		}
		store.shares[id] = append([]byte(nil), share...)
		store.mu.Unlock()
		writeJSONResponse(w, map[string]string{"id": id})
		return
	}
	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/shares/") {
		id := strings.TrimPrefix(r.URL.Path, "/shares/")
		store.mu.Lock()
		share, ok := store.shares[id]
		store.mu.Unlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSONResponse(w, map[string]string{"share_b64": base64.StdEncoding.EncodeToString(share)})
		return
	}
	http.NotFound(w, r)
}

func (store *fakeShareStore) setOnline(online bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.online = online
}

type fakeAuditCollector struct {
	mu     sync.Mutex
	events []string
}

func (collector *fakeAuditCollector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.Path != "/events" {
		http.NotFound(w, r)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		http.Error(w, "read event", http.StatusBadRequest)
		return
	}
	collector.mu.Lock()
	collector.events = append(collector.events, string(body))
	collector.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (collector *fakeAuditCollector) eventBodies() []string {
	collector.mu.Lock()
	defer collector.mu.Unlock()
	return append([]string(nil), collector.events...)
}

type fakeApprovalPolicy struct {
	token    string
	mu       sync.Mutex
	approved bool
	requests int
}

func (policy *fakeApprovalPolicy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.Path != "/approve" {
		http.NotFound(w, r)
		return
	}
	if r.Header.Get("Authorization") != "Bearer "+policy.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var request struct {
		KeyID     string `json:"key_id"`
		Operation string `json:"operation"`
	}
	if !decodeJSONRequest(w, r, &request) || request.KeyID == "" || request.Operation != "unwrap" {
		http.Error(w, "invalid approval request", http.StatusBadRequest)
		return
	}
	policy.mu.Lock()
	policy.requests++
	approved := policy.approved
	policy.mu.Unlock()
	if !approved {
		http.Error(w, "approval required", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (policy *fakeApprovalPolicy) setApproved(approved bool) {
	policy.mu.Lock()
	defer policy.mu.Unlock()
	policy.approved = approved
}

func (policy *fakeApprovalPolicy) requestCount() int {
	policy.mu.Lock()
	defer policy.mu.Unlock()
	return policy.requests
}

type fakeOIDCKMS struct {
	clientSecret   string
	token          string
	key            []byte
	mu             sync.Mutex
	tokenExchanges int
	wraps          int
	unwraps        int
}

func (service *fakeOIDCKMS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/token":
		var request struct {
			Subject      string `json:"subject"`
			ClientSecret string `json:"client_secret"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		if request.Subject != "sops-test-workload" || request.ClientSecret != service.clientSecret {
			http.Error(w, "invalid client credentials", http.StatusUnauthorized)
			return
		}
		service.mu.Lock()
		service.tokenExchanges++
		service.mu.Unlock()
		writeJSONResponse(w, map[string]string{"access_token": service.token})
	case "/kms/wrap":
		if !service.authorized(w, r) {
			return
		}
		var request struct {
			KeyID   string `json:"key_id"`
			DataKey string `json:"data_key_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		dataKey, err := base64.StdEncoding.DecodeString(request.DataKey)
		if err != nil {
			http.Error(w, "invalid data key", http.StatusBadRequest)
			return
		}
		gcm, err := newTestGCM(service.key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ciphertext := append(nonce, gcm.Seal(nil, nonce, dataKey, []byte(request.KeyID))...)
		service.mu.Lock()
		service.wraps++
		service.mu.Unlock()
		writeJSONResponse(w, map[string]string{
			"ciphertext_b64": base64.StdEncoding.EncodeToString(ciphertext),
			"version":        "v1",
		})
	case "/kms/unwrap":
		if !service.authorized(w, r) {
			return
		}
		var request struct {
			KeyID         string `json:"key_id"`
			Version       string `json:"version"`
			CiphertextB64 string `json:"ciphertext_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		ciphertext, err := base64.StdEncoding.DecodeString(request.CiphertextB64)
		if err != nil {
			http.Error(w, "invalid ciphertext", http.StatusBadRequest)
			return
		}
		gcm, err := newTestGCM(service.key)
		if err != nil || len(ciphertext) < gcm.NonceSize() {
			http.Error(w, "invalid ciphertext", http.StatusBadRequest)
			return
		}
		nonce := ciphertext[:gcm.NonceSize()]
		dataKey, err := gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], []byte(request.KeyID))
		if err != nil {
			http.Error(w, "unable to unwrap", http.StatusBadRequest)
			return
		}
		service.mu.Lock()
		service.unwraps++
		service.mu.Unlock()
		writeJSONResponse(w, map[string]string{"data_key_b64": base64.StdEncoding.EncodeToString(dataKey)})
	default:
		http.NotFound(w, r)
	}
}

func (service *fakeOIDCKMS) authorized(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Authorization") != "Bearer "+service.token {
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return false
	}
	return true
}

func (service *fakeOIDCKMS) tokenExchangeCount() int {
	service.mu.Lock()
	defer service.mu.Unlock()
	return service.tokenExchanges
}

func (service *fakeOIDCKMS) wrapCount() int {
	service.mu.Lock()
	defer service.mu.Unlock()
	return service.wraps
}

func (service *fakeOIDCKMS) unwrapCount() int {
	service.mu.Lock()
	defer service.mu.Unlock()
	return service.unwraps
}
