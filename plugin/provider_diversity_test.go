package plugin_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keyservice"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
)

const (
	noMetadataPluginName      = "no-metadata"
	fileAESPluginName         = "file-aes"
	agentPluginName           = "key-agent"
	secretReferencePluginName = "secret-reference"
	multiRecipientPluginName  = "multi-recipient"
	contractPluginName        = "contract-fixture"
	mockAgentToken            = "test-agent-credential"
	mockSecretReferenceToken  = "test-secret-reference-credential"
)

func TestNoMetadataPluginRoundTripsThroughKeyService(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/nometa")
	setAllowlist(t, binary, noMetadataPluginName)

	key, err := sopsplugin.NewMasterKey(noMetadataPluginName, map[string]interface{}{})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.Empty(t, key.ProviderMetadata)
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
}

func TestFileAESPluginUsesConfiguredFileAndDetectsKeyChanges(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/fileaesfixture")
	setAllowlist(t, binary, fileAESPluginName)
	keyFile := filepath.Join(t.TempDir(), "master.key")
	firstKey := []byte("0123456789abcdef0123456789abcdef")
	require.NoError(t, os.WriteFile(keyFile, firstKey, 0600))

	key, err := sopsplugin.NewMasterKey(fileAESPluginName, map[string]interface{}{
		"key_file": keyFile,
		"key_id":   "local-master-v1",
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.Equal(t, "local-master-v1", key.ProviderMetadata["key_id"])
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)

	require.NoError(t, os.WriteFile(keyFile, []byte("abcdef0123456789abcdef0123456789"), 0600))
	_, err = key.Decrypt()
	require.ErrorContains(t, err, "decrypt data key")
}

func TestAgentPluginKeepsPrivateKeyOutsidePluginProcess(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/agentfixture")
	privateKey, publicKeyPEM, privateKeyPEM := generateTestRSAKeypair(t)
	agent := &fakePrivateKeyAgent{token: mockAgentToken, privateKey: privateKey}
	server := httptest.NewServer(agent)
	t.Cleanup(server.Close)
	setAllowlist(t, binary, agentPluginName)
	t.Setenv("SOPS_TEST_AGENT_TOKEN", mockAgentToken)

	key, err := sopsplugin.NewMasterKey(agentPluginName, map[string]interface{}{
		"agent_endpoint": server.URL,
		"key_id":         "hardware-key-1",
		"public_key":     string(publicKeyPEM),
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
	require.Equal(t, 1, agent.decryptions())

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), string(privateKeyPEM))
}

func TestSecretReferencePluginStoresOnlyAnOpaqueReferenceInSOPS(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/secretreffixture")
	server := newFakeSecretReferenceStore(t, mockSecretReferenceToken)
	setAllowlist(t, binary, secretReferencePluginName)
	t.Setenv("SOPS_TEST_SECRET_REF_TOKEN", mockSecretReferenceToken)

	key, err := sopsplugin.NewMasterKey(secretReferencePluginName, map[string]interface{}{
		"endpoint":  server.URL(),
		"namespace": "sops-tests",
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.NotEmpty(t, key.ProviderMetadata["secret_id"])
	require.NotEqual(t, dataKey, key.EncryptedKey)
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
	require.Equal(t, 2, server.wrapsAndUnwraps())

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), base64.StdEncoding.EncodeToString(dataKey))
	require.NotContains(t, string(encoded), mockSecretReferenceToken)
}

func TestMultiRecipientPluginFallsBackToAnotherConfiguredIdentity(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/multirecipientfixture")
	setAllowlist(t, binary, multiRecipientPluginName)
	_, publicA, privateBytesA := generateTestRSAKeypair(t)
	privateB, publicB, _ := generateTestRSAKeypair(t)

	keyDir := t.TempDir()
	privatePathA := filepath.Join(keyDir, "recipient-a.pem")
	privatePathB := filepath.Join(keyDir, "recipient-b.pem")
	require.NoError(t, os.WriteFile(privatePathA, privateBytesA, 0600))
	privateBytesB, err := x509.MarshalPKCS8PrivateKey(privateB)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(privatePathB, pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytesB,
	}), 0600))

	key, err := sopsplugin.NewMasterKey(multiRecipientPluginName, map[string]interface{}{
		"recipients": map[string]interface{}{
			"recipient-a": string(publicA),
			"recipient-b": string(publicB),
		},
		"private_key_files": map[string]interface{}{
			"recipient-a": privatePathA,
			"recipient-b": privatePathB,
		},
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.Equal(t, []interface{}{"recipient-a", "recipient-b"}, key.ProviderMetadata["recipient_ids"])
	require.NoError(t, os.Remove(privatePathA))
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)

	require.NoError(t, os.Remove(privatePathB))
	_, err = key.Decrypt()
	require.ErrorContains(t, err, "no configured recipient could unwrap the data key")
}

func TestContractFixtureErrorsAndInvalidResponsesAreSurfaced(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/contractfixture")
	setAllowlist(t, binary, contractPluginName)
	dataKey := []byte("01234567890123456789012345678901")

	_, err := sopsplugin.WrapKey(context.Background(), contractPluginName, sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{"mode": "error"},
		DataKey: dataKey,
	})
	require.ErrorContains(t, err, "fixture wrap failure")

	_, err = sopsplugin.WrapKey(context.Background(), contractPluginName, sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{"mode": "empty-wrap"},
		DataKey: dataKey,
	})
	require.ErrorContains(t, err, "returned an empty wrapped data key")

	_, err = sopsplugin.WrapKey(context.Background(), contractPluginName, sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{"mode": "invalid-metadata"},
		DataKey: dataKey,
	})
	require.ErrorContains(t, err, "unsupported type")

	_, err = sopsplugin.UnwrapKey(context.Background(), contractPluginName, sopsplugin.UnwrapKeyRequest{
		Config:           map[string]interface{}{"mode": "error"},
		EncryptedDataKey: []byte("ciphertext"),
	})
	require.ErrorContains(t, err, "fixture unwrap failure")

	_, err = sopsplugin.UnwrapKey(context.Background(), contractPluginName, sopsplugin.UnwrapKeyRequest{
		Config:           map[string]interface{}{"mode": "empty-unwrap"},
		EncryptedDataKey: []byte("ciphertext"),
	})
	require.ErrorContains(t, err, "returned an empty data key")

	_, err = sopsplugin.UnwrapKey(context.Background(), contractPluginName, sopsplugin.UnwrapKeyRequest{
		Config:           map[string]interface{}{"mode": "short-unwrap"},
		EncryptedDataKey: []byte("ciphertext"),
	})
	require.ErrorContains(t, err, "expected 32 or 33 bytes, got 5")
}

func generateTestRSAKeypair(t *testing.T) (*rsa.PrivateKey, []byte, []byte) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	return privateKey, pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}), pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		})
}

type fakePrivateKeyAgent struct {
	token      string
	privateKey *rsa.PrivateKey
	mu         sync.Mutex
	count      int
}

func (agent *fakePrivateKeyAgent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/decrypt" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	if r.Header.Get("Authorization") != "Bearer "+agent.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var request struct {
		Ciphertext string `json:"ciphertext_b64"`
	}
	if !decodeJSONRequest(w, r, &request) {
		return
	}
	ciphertext, err := base64.StdEncoding.DecodeString(request.Ciphertext)
	if err != nil {
		http.Error(w, "invalid ciphertext", http.StatusBadRequest)
		return
	}
	dataKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, agent.privateKey, ciphertext, nil)
	if err != nil {
		http.Error(w, "unable to unwrap data key", http.StatusBadRequest)
		return
	}
	agent.mu.Lock()
	agent.count++
	agent.mu.Unlock()
	writeJSONResponse(w, map[string]string{"data_key_b64": base64.StdEncoding.EncodeToString(dataKey)})
}

func (agent *fakePrivateKeyAgent) decryptions() int {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	return agent.count
}

type fakeSecretReferenceStore struct {
	token      string
	mu         sync.Mutex
	nextID     int
	secrets    map[string][]byte
	operations int
	httpServer *httptest.Server
}

func newFakeSecretReferenceStore(t *testing.T, token string) *fakeSecretReferenceStore {
	t.Helper()
	store := &fakeSecretReferenceStore{
		token:   token,
		secrets: make(map[string][]byte),
	}
	store.httpServer = httptest.NewServer(store)
	t.Cleanup(store.httpServer.Close)
	return store
}

func (store *fakeSecretReferenceStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+store.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/v1/secrets" {
		var request struct {
			Namespace string `json:"namespace"`
			DataKey   string `json:"data_key_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		dataKey, err := base64.StdEncoding.DecodeString(request.DataKey)
		if err != nil || request.Namespace == "" {
			http.Error(w, "invalid secret request", http.StatusBadRequest)
			return
		}
		store.mu.Lock()
		store.nextID++
		secretID := fmt.Sprintf("%s-%03d", request.Namespace, store.nextID)
		store.secrets[secretID] = append([]byte(nil), dataKey...)
		store.operations++
		store.mu.Unlock()
		writeJSONResponse(w, map[string]string{"secret_id": secretID, "version": "v1"})
		return
	}
	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/secrets/") {
		secretID := strings.TrimPrefix(r.URL.Path, "/v1/secrets/")
		namespace := r.URL.Query().Get("namespace")
		version := r.URL.Query().Get("version")
		store.mu.Lock()
		dataKey, ok := store.secrets[secretID]
		if ok && version == "v1" && strings.HasPrefix(secretID, namespace+"-") {
			store.operations++
		}
		store.mu.Unlock()
		if !ok || version != "v1" || !strings.HasPrefix(secretID, namespace+"-") {
			http.NotFound(w, r)
			return
		}
		writeJSONResponse(w, map[string]string{"data_key_b64": base64.StdEncoding.EncodeToString(dataKey)})
		return
	}
	http.NotFound(w, r)
}

func (store *fakeSecretReferenceStore) wrapsAndUnwraps() int {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.operations
}

func (store *fakeSecretReferenceStore) URL() string {
	return store.httpServer.URL
}
