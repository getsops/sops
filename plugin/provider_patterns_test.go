package plugin_test

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keyservice"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
)

const (
	mockKMSToken       = "test-kms-credential"
	mockBarbicanToken  = "test-barbican-credential"
	mockKMSPluginName  = "mock-kms"
	mockBarbicanPlugin = "mock-barbican"
)

func TestRemoteKMSPluginSupportsMultipleShamirKeyGroups(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/remotekms")
	server := newFakeRemoteKMS(t, mockKMSToken)
	setAllowlist(t, binary, mockKMSPluginName)
	t.Setenv("SOPS_TEST_KMS_TOKEN", mockKMSToken)

	firstKey, err := sopsplugin.NewMasterKey(mockKMSPluginName, map[string]interface{}{
		"endpoint": server.URL(),
		"key_id":   "alpha",
	})
	require.NoError(t, err)
	secondKey, err := sopsplugin.NewMasterKey(mockKMSPluginName, map[string]interface{}{
		"endpoint": server.URL(),
		"key_id":   "beta",
	})
	require.NoError(t, err)
	metadata := sops.Metadata{
		ShamirThreshold: 2,
		KeyGroups:       []sops.KeyGroup{{firstKey}, {secondKey}},
	}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.Equal(t, map[string]interface{}{"key_id": "alpha", "version": "v1"}, firstKey.ProviderMetadata)
	require.Equal(t, map[string]interface{}{"key_id": "beta", "version": "v1"}, secondKey.ProviderMetadata)

	// UpdateMasterKeys caches the plaintext data key; clear it to exercise the
	// plugin Unwrap path and Shamir reconstruction rather than the cache.
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
	require.Equal(t, map[string]int{"alpha": 1, "beta": 1}, server.wrapCounts())
	require.Equal(t, map[string]int{"alpha": 1, "beta": 1}, server.unwrapCounts())

	for _, key := range []*sopsplugin.MasterKey{firstKey, secondKey} {
		encoded, err := json.Marshal(key.ToMap())
		require.NoError(t, err)
		require.NotContains(t, string(encoded), mockKMSToken)
	}
}

func TestRemoteKMSPluginPropagatesCredentialErrors(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/remotekms")
	server := newFakeRemoteKMS(t, mockKMSToken)
	setAllowlist(t, binary, mockKMSPluginName)
	t.Setenv("SOPS_TEST_KMS_TOKEN", "wrong-credential")

	key, err := sopsplugin.NewMasterKey(mockKMSPluginName, map[string]interface{}{
		"endpoint": server.URL(),
		"key_id":   "alpha",
	})
	require.NoError(t, err)
	err = key.Encrypt([]byte("01234567890123456789012345678901"))
	require.ErrorContains(t, err, "HTTP 401")
}

func TestRemoteKMSPluginHonorsContextCancellation(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/remotekms")
	requestStarted := make(chan struct{}, 1)
	requestCanceled := make(chan struct{})
	releaseResponse := make(chan struct{})
	var releaseOnce sync.Once
	release := func() {
		releaseOnce.Do(func() { close(releaseResponse) })
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			return
		}
		requestStarted <- struct{}{}
		select {
		case <-r.Context().Done():
			close(requestCanceled)
		case <-releaseResponse:
		}
		http.Error(w, "request released", http.StatusServiceUnavailable)
	}))
	t.Cleanup(func() {
		release()
		server.Close()
	})
	setAllowlist(t, binary, mockKMSPluginName)
	t.Setenv("SOPS_TEST_KMS_TOKEN", mockKMSToken)

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := sopsplugin.WrapKey(ctx, mockKMSPluginName, sopsplugin.WrapKeyRequest{
			Config: map[string]interface{}{
				"endpoint": server.URL,
				"key_id":   "alpha",
			},
			DataKey: []byte("01234567890123456789012345678901"),
		})
		result <- err
	}()

	select {
	case <-requestStarted:
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("plugin did not send its remote KMS request")
	}
	cancel()
	select {
	case <-requestCanceled:
	case <-time.After(5 * time.Second):
		release()
		t.Fatal("plugin did not cancel its remote KMS request")
	}
	release()
	select {
	case err := <-result:
		require.Error(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("plugin call did not stop after context cancellation")
	}
}

func TestSOPSKeyServicePropagatesContextCancellationToPlugin(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/remotekms")
	requestStarted := make(chan struct{}, 1)
	requestCanceled := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			return
		}
		requestStarted <- struct{}{}
		<-r.Context().Done()
		close(requestCanceled)
		http.Error(w, "request canceled", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)
	setAllowlist(t, binary, mockKMSPluginName)
	t.Setenv("SOPS_TEST_KMS_TOKEN", mockKMSToken)

	key, err := sopsplugin.NewMasterKey(mockKMSPluginName, map[string]interface{}{
		"endpoint": server.URL,
		"key_id":   "alpha",
	})
	require.NoError(t, err)
	key.ProviderMetadata = map[string]interface{}{"key_id": "alpha", "version": "v1"}
	key.EncryptedKey = []byte("wrapped data key")
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	client := keyservice.NewLocalClient()
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := metadata.GetDataKeyWithKeyServicesContext(ctx, []keyservice.KeyServiceClient{client}, nil)
		result <- err
	}()

	select {
	case <-requestStarted:
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("SOPS key service did not send its remote KMS request")
	}
	cancel()
	select {
	case <-requestCanceled:
	case <-time.After(5 * time.Second):
		t.Fatal("plugin did not cancel its remote KMS request")
	}
	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(10 * time.Second):
		t.Fatal("SOPS key-service request did not stop after context cancellation")
	}
}

func TestBarbicanStylePluginRetrievesVersionedSecretAndWrapsLocally(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/barbicanfixture")
	server := newFakeSecretStore(t, mockBarbicanToken)
	setAllowlist(t, binary, mockBarbicanPlugin)
	t.Setenv("SOPS_TEST_BARBICAN_TOKEN", mockBarbicanToken)

	key, err := sopsplugin.NewMasterKey(mockBarbicanPlugin, map[string]interface{}{
		"endpoint":  server.URL(),
		"secret_id": "sops-data-key",
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	firstDataKey := []byte("01234567890123456789012345678901")
	secondDataKey := []byte("abcdefghijabcdefghijabcdefghij12")

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(firstDataKey, services))
	require.Equal(t, "v1", key.ProviderMetadata["secret_version"])
	require.NotEqual(t, firstDataKey, key.EncryptedKey)
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, firstDataKey, decoded)

	server.setCurrentVersion("v2")
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(secondDataKey, services))
	require.Equal(t, "v2", key.ProviderMetadata["secret_version"])
	metadata.DataKey = nil
	decoded, err = metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, secondDataKey, decoded)
	require.Equal(t, map[string]int{"current": 2, "v1": 1, "v2": 1}, server.requestCounts())

	encoded, err := json.Marshal(key.ToMap())
	require.NoError(t, err)
	require.NotContains(t, string(encoded), base64.StdEncoding.EncodeToString(server.secret("v1")))
	require.NotContains(t, string(encoded), mockBarbicanToken)
}

type fakeRemoteKMS struct {
	token      string
	keys       map[string][]byte
	mu         sync.Mutex
	wraps      map[string]int
	unwraps    map[string]int
	httpServer *httptest.Server
}

func newFakeRemoteKMS(t *testing.T, token string) *fakeRemoteKMS {
	t.Helper()
	server := &fakeRemoteKMS{
		token:   token,
		keys:    map[string][]byte{"alpha": bytes.Repeat([]byte{0x11}, 32), "beta": bytes.Repeat([]byte{0x22}, 32)},
		wraps:   make(map[string]int),
		unwraps: make(map[string]int),
	}
	server.httpServer = httptest.NewServer(server)
	t.Cleanup(server.httpServer.Close)
	return server
}

func (server *fakeRemoteKMS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+server.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.URL.Path {
	case "/v1/wrap":
		var request struct {
			KeyID      string `json:"key_id"`
			DataKeyB64 string `json:"data_key_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		key, ok := server.keys[request.KeyID]
		if !ok {
			http.Error(w, "unknown key", http.StatusNotFound)
			return
		}
		dataKey, err := base64.StdEncoding.DecodeString(request.DataKeyB64)
		if err != nil {
			http.Error(w, "invalid data key", http.StatusBadRequest)
			return
		}
		gcm, err := newTestGCM(key)
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
		server.mu.Lock()
		server.wraps[request.KeyID]++
		server.mu.Unlock()
		writeJSONResponse(w, map[string]string{
			"ciphertext": base64.StdEncoding.EncodeToString(ciphertext),
			"version":    "v1",
		})
	case "/v1/unwrap":
		var request struct {
			KeyID         string `json:"key_id"`
			Version       string `json:"version"`
			CiphertextB64 string `json:"ciphertext_b64"`
		}
		if !decodeJSONRequest(w, r, &request) {
			return
		}
		if request.Version != "v1" {
			http.Error(w, "unknown key version", http.StatusNotFound)
			return
		}
		key, ok := server.keys[request.KeyID]
		if !ok {
			http.Error(w, "unknown key", http.StatusNotFound)
			return
		}
		ciphertext, err := base64.StdEncoding.DecodeString(request.CiphertextB64)
		if err != nil {
			http.Error(w, "invalid ciphertext", http.StatusBadRequest)
			return
		}
		gcm, err := newTestGCM(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(ciphertext) < gcm.NonceSize() {
			http.Error(w, "short ciphertext", http.StatusBadRequest)
			return
		}
		nonce := ciphertext[:gcm.NonceSize()]
		dataKey, err := gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], []byte(request.KeyID))
		if err != nil {
			http.Error(w, "unable to unwrap key", http.StatusBadRequest)
			return
		}
		server.mu.Lock()
		server.unwraps[request.KeyID]++
		server.mu.Unlock()
		writeJSONResponse(w, map[string]string{"data_key_b64": base64.StdEncoding.EncodeToString(dataKey)})
	default:
		http.NotFound(w, r)
	}
}

func (server *fakeRemoteKMS) wrapCounts() map[string]int {
	server.mu.Lock()
	defer server.mu.Unlock()
	return copyCounts(server.wraps)
}

func (server *fakeRemoteKMS) unwrapCounts() map[string]int {
	server.mu.Lock()
	defer server.mu.Unlock()
	return copyCounts(server.unwraps)
}

func (server *fakeRemoteKMS) URL() string {
	return server.httpServer.URL
}

type fakeSecretStore struct {
	token          string
	currentVersion string
	secrets        map[string][]byte
	mu             sync.Mutex
	requests       map[string]int
	httpServer     *httptest.Server
}

func newFakeSecretStore(t *testing.T, token string) *fakeSecretStore {
	t.Helper()
	server := &fakeSecretStore{
		token:          token,
		currentVersion: "v1",
		secrets: map[string][]byte{
			"v1": bytes.Repeat([]byte{0x33}, 32),
			"v2": bytes.Repeat([]byte{0x44}, 32),
		},
		requests: make(map[string]int),
	}
	server.httpServer = httptest.NewServer(server)
	t.Cleanup(server.httpServer.Close)
	return server
}

func (server *fakeSecretStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+server.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodGet || r.URL.Path != "/v1/secrets/sops-data-key" {
		http.NotFound(w, r)
		return
	}
	requestedVersion := r.URL.Query().Get("version")
	server.mu.Lock()
	version := requestedVersion
	if version == "current" {
		version = server.currentVersion
	}
	secret, ok := server.secrets[version]
	if ok {
		server.requests[requestedVersion]++
	}
	server.mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSONResponse(w, map[string]string{
		"secret_b64": base64.StdEncoding.EncodeToString(secret),
		"version":    version,
	})
}

func (server *fakeSecretStore) setCurrentVersion(version string) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.currentVersion = version
}

func (server *fakeSecretStore) requestCounts() map[string]int {
	server.mu.Lock()
	defer server.mu.Unlock()
	return copyCounts(server.requests)
}

func (server *fakeSecretStore) secret(version string) []byte {
	server.mu.Lock()
	defer server.mu.Unlock()
	return append([]byte(nil), server.secrets[version]...)
}

func (server *fakeSecretStore) URL() string {
	return server.httpServer.URL
}

func newTestGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func decodeJSONRequest(w http.ResponseWriter, r *http.Request, destination interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(destination); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON request: %v", err), http.StatusBadRequest)
		return false
	}
	return true
}

func writeJSONResponse(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, "encode response", http.StatusInternalServerError)
	}
}

func copyCounts(counts map[string]int) map[string]int {
	copy := make(map[string]int, len(counts))
	for key, count := range counts {
		copy[key] = count
	}
	return copy
}
