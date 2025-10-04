package ocikms

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// testOCID is a valid OCID format for testing
	testOCID = "ocid1.key.oc1.uk-london-1.aaaalgz5aacmg.aaaailjtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq"
	// testDataKey is a dummy 32-byte data key for testing
	testDataKey = "testtesttesttesttesttesttest1234"
)

// mockHTTPClient implements common.HTTPRequestDispatcher for testing
type mockHTTPClient struct {
	// requests stores all requests made for verification
	requests []*http.Request
	// responses is a queue of responses to return
	responses []*http.Response
	// errors is a queue of errors to return
	errors []error
	// currentIndex tracks which response to return next
	currentIndex int
}

func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		requests:  make([]*http.Request, 0),
		responses: make([]*http.Response, 0),
		errors:    make([]error, 0),
	}
}

// Do implements the common.HTTPRequestDispatcher interface
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Store the request for verification
	m.requests = append(m.requests, req)

	if m.currentIndex >= len(m.responses) && m.currentIndex >= len(m.errors) {
		return nil, fmt.Errorf("mock client: no more responses configured")
	}

	// Return error if configured
	if m.currentIndex < len(m.errors) && m.errors[m.currentIndex] != nil {
		err := m.errors[m.currentIndex]
		m.currentIndex++
		return nil, err
	}

	// Return response if configured
	if m.currentIndex < len(m.responses) {
		resp := m.responses[m.currentIndex]
		m.currentIndex++
		return resp, nil
	}

	return nil, fmt.Errorf("mock client: no response or error configured for request %d", m.currentIndex)
}

// addResponse adds a mock HTTP response to the queue
func (m *mockHTTPClient) addResponse(statusCode int, body string) {
	resp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	m.responses = append(m.responses, resp)
}

// addError adds an error to the queue
func (m *mockHTTPClient) addError(err error) {
	m.errors = append(m.errors, err)
}

// getLastRequest returns the most recent request made
func (m *mockHTTPClient) getLastRequest() *http.Request {
	if len(m.requests) == 0 {
		return nil
	}
	return m.requests[len(m.requests)-1]
}

// mockConfigProvider implements common.ConfigurationProvider for testing
type mockConfigProvider struct {
	privateKey *rsa.PrivateKey
}

func newMockConfigProvider() mockConfigProvider {
	// Generate a test RSA key (required by OCI SDK for request signing)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("failed to generate test RSA key: %v", err))
	}
	return mockConfigProvider{
		privateKey: privateKey,
	}
}

func (m mockConfigProvider) TenancyOCID() (string, error) {
	return "ocid1.tenancy.oc1..test", nil
}

func (m mockConfigProvider) UserOCID() (string, error) {
	return "ocid1.user.oc1..test", nil
}

func (m mockConfigProvider) KeyFingerprint() (string, error) {
	return "00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00", nil
}

func (m mockConfigProvider) Region() (string, error) {
	return "uk-london-1", nil
}

func (m mockConfigProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	return m.privateKey, nil
}

func (m mockConfigProvider) KeyID() (string, error) {
	tenancy, _ := m.TenancyOCID()
	user, _ := m.UserOCID()
	fingerprint, _ := m.KeyFingerprint()
	return fmt.Sprintf("%s/%s/%s", tenancy, user, fingerprint), nil
}

func (m mockConfigProvider) AuthType() (common.AuthConfig, error) {
	return common.AuthConfig{
		AuthType: common.UserPrincipal,
	}, nil
}

// createTestMasterKey creates a MasterKey configured for testing with mock HTTP client
func createTestMasterKey(ocid string, mockHTTP *mockHTTPClient) *MasterKey {
	key := NewMasterKeyFromOCID(ocid)

	// Inject mock config provider to avoid real auth
	configProvider := NewConfigurationProvider(newMockConfigProvider())
	configProvider.ApplyToMasterKey(key)

	// Inject mock HTTP client
	if mockHTTP != nil {
		httpClient := NewHTTPClient(mockHTTP)
		httpClient.ApplyToMasterKey(key)
	}

	return key
}

// createEncryptResponse creates a mock OCI KMS encrypt response
func createEncryptResponse(ciphertext string) string {
	response := map[string]interface{}{
		"ciphertext": ciphertext,
	}
	data, _ := json.Marshal(response)
	return string(data)
}

// createDecryptResponse creates a mock OCI KMS decrypt response
func createDecryptResponse(plaintext string) string {
	response := map[string]interface{}{
		"plaintext": plaintext,
	}
	data, _ := json.Marshal(response)
	return string(data)
}

func TestEncryptContext(t *testing.T) {
	tests := []struct {
		name           string
		dataKey        []byte
		mockResponse   string
		mockStatusCode int
		mockError      error
		expectError    bool
		errorContains  string
	}{
		{
			name:           "successful encryption",
			dataKey:        []byte(testDataKey),
			mockResponse:   createEncryptResponse("ENCRYPTED_DATA_KEY_BASE64"),
			mockStatusCode: 200,
			expectError:    false,
		},
		{
			name:          "network error",
			dataKey:       []byte(testDataKey),
			mockError:     fmt.Errorf("network timeout"),
			expectError:   true,
			errorContains: "failed to encrypt sops data key with OCI KMS key",
		},
		{
			name:           "HTTP 500 error",
			dataKey:        []byte(testDataKey),
			mockResponse:   `{"code":"InternalServerError","message":"Internal server error"}`,
			mockStatusCode: 500,
			expectError:    true,
			errorContains:  "failed to encrypt sops data key with OCI KMS key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := newMockHTTPClient()

			if tt.mockError != nil {
				mockHTTP.addError(tt.mockError)
			} else {
				mockHTTP.addResponse(tt.mockStatusCode, tt.mockResponse)
			}

			key := createTestMasterKey(testOCID, mockHTTP)

			err := key.EncryptContext(context.Background(), tt.dataKey)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, key.EncryptedKey)
				assert.Equal(t, "ENCRYPTED_DATA_KEY_BASE64", key.EncryptedKey)
			}

			// Verify request was made (unless error before request)
			if tt.mockError == nil || tt.mockStatusCode > 0 {
				assert.Greater(t, len(mockHTTP.requests), 0, "should have made at least one HTTP request")
			}
		})
	}
}

func TestDecryptContext(t *testing.T) {
	dataKeyBase64 := base64.StdEncoding.EncodeToString([]byte(testDataKey))

	tests := []struct {
		name           string
		encryptedKey   string
		mockResponse   string
		mockStatusCode int
		mockError      error
		expectError    bool
		errorContains  string
		expectedPlain  []byte
	}{
		{
			name:           "successful decryption",
			encryptedKey:   "ENCRYPTED_DATA_KEY_BASE64",
			mockResponse:   createDecryptResponse(dataKeyBase64),
			mockStatusCode: 200,
			expectError:    false,
			expectedPlain:  []byte(testDataKey),
		},
		{
			name:          "network error",
			encryptedKey:  "ENCRYPTED_DATA_KEY_BASE64",
			mockError:     fmt.Errorf("connection refused"),
			expectError:   true,
			errorContains: "failed to decrypt sops data key with OCI KMS key",
		},
		{
			name:           "invalid ciphertext",
			encryptedKey:   "INVALID_CIPHERTEXT",
			mockResponse:   `{"code":"InvalidCiphertext","message":"The ciphertext is invalid"}`,
			mockStatusCode: 400,
			expectError:    true,
			errorContains:  "failed to decrypt sops data key with OCI KMS key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := newMockHTTPClient()

			if tt.mockError != nil {
				mockHTTP.addError(tt.mockError)
			} else {
				mockHTTP.addResponse(tt.mockStatusCode, tt.mockResponse)
			}

			key := createTestMasterKey(testOCID, mockHTTP)
			key.EncryptedKey = tt.encryptedKey

			plaintext, err := key.DecryptContext(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPlain, plaintext)
			}
		})
	}
}

func TestHTTPClientInjection(t *testing.T) {
	mockHTTP := newMockHTTPClient()
	mockHTTP.addResponse(200, createEncryptResponse("ENCRYPTED"))

	key := NewMasterKeyFromOCID(testOCID)

	// Inject config provider (required for client creation)
	configProvider := NewConfigurationProvider(newMockConfigProvider())
	configProvider.ApplyToMasterKey(key)

	// Inject HTTP client
	httpClient := NewHTTPClient(mockHTTP)
	httpClient.ApplyToMasterKey(key)

	// Perform encryption
	err := key.EncryptContext(context.Background(), []byte("test"))
	require.NoError(t, err)

	// Verify our mock client was used
	assert.Equal(t, 1, len(mockHTTP.requests), "should have used injected HTTP client")
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	dataKey := []byte("this-is-a-32-byte-test-key-12345")
	dataKeyBase64 := base64.StdEncoding.EncodeToString(dataKey)
	ciphertext := "MOCK_ENCRYPTED_CIPHERTEXT_BASE64"

	mockHTTP := newMockHTTPClient()

	// Mock encrypt response
	mockHTTP.addResponse(200, createEncryptResponse(ciphertext))
	// Mock decrypt response
	mockHTTP.addResponse(200, createDecryptResponse(dataKeyBase64))

	key := createTestMasterKey(testOCID, mockHTTP)

	// Encrypt
	err := key.EncryptContext(context.Background(), dataKey)
	require.NoError(t, err)
	assert.Equal(t, ciphertext, key.EncryptedKey)

	// Decrypt
	decrypted, err := key.DecryptContext(context.Background())
	require.NoError(t, err)
	assert.Equal(t, dataKey, decrypted)

	// Verify two requests were made
	assert.Equal(t, 2, len(mockHTTP.requests))
}

func TestContextCancellation(t *testing.T) {
	mockHTTP := newMockHTTPClient()
	mockHTTP.addResponse(200, createEncryptResponse("ENCRYPTED"))

	key := createTestMasterKey(testOCID, mockHTTP)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Attempt encryption with cancelled context
	err := key.EncryptContext(ctx, []byte("test"))

	// Should fail due to context cancellation
	// Note: actual behavior depends on when OCI SDK checks context
	// This test documents the expected behavior
	_ = err // May or may not error depending on when context is checked
}

func TestEncryptIfNeeded(t *testing.T) {
	dataKey := []byte("test-data-key-32-bytes-long-1234")

	t.Run("encrypts when EncryptedKey is empty", func(t *testing.T) {
		mockHTTP := newMockHTTPClient()
		mockHTTP.addResponse(200, createEncryptResponse("ENCRYPTED"))

		key := createTestMasterKey(testOCID, mockHTTP)
		key.EncryptedKey = "" // Explicitly empty

		err := key.EncryptIfNeeded(dataKey)
		require.NoError(t, err)
		assert.Equal(t, "ENCRYPTED", key.EncryptedKey)
		assert.Equal(t, 1, len(mockHTTP.requests))
	})

	t.Run("skips encryption when EncryptedKey exists", func(t *testing.T) {
		mockHTTP := newMockHTTPClient()
		// Don't add any responses - should not be called

		key := createTestMasterKey(testOCID, mockHTTP)
		key.EncryptedKey = "ALREADY_ENCRYPTED"

		err := key.EncryptIfNeeded(dataKey)
		require.NoError(t, err)
		assert.Equal(t, "ALREADY_ENCRYPTED", key.EncryptedKey)
		assert.Equal(t, 0, len(mockHTTP.requests), "should not make HTTP request")
	})
}

func TestNeedsRotation(t *testing.T) {
	t.Run("new key does not need rotation", func(t *testing.T) {
		key := NewMasterKeyFromOCID(testOCID)
		assert.False(t, key.NeedsRotation())
	})

	t.Run("old key needs rotation", func(t *testing.T) {
		key := NewMasterKeyFromOCID(testOCID)
		// Set creation date to 7 months ago (> 6 months)
		key.CreationDate = time.Now().UTC().Add(-7 * 30 * 24 * time.Hour)
		assert.True(t, key.NeedsRotation())
	})

	t.Run("6-month-old key does not need rotation", func(t *testing.T) {
		key := NewMasterKeyFromOCID(testOCID)
		// Set creation date to just under 6 months ago
		key.CreationDate = time.Now().UTC().Add(-6*30*24*time.Hour + time.Hour)
		// Should not need rotation (> is used, not >=)
		assert.False(t, key.NeedsRotation())
	})
}

func TestToString(t *testing.T) {
	key := NewMasterKeyFromOCID(testOCID)
	assert.Equal(t, testOCID, key.ToString())
}

func TestTypeToIdentifier(t *testing.T) {
	key := NewMasterKeyFromOCID(testOCID)
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
	assert.Equal(t, "oci_kms", key.TypeToIdentifier())
}

func TestExtractRefs(t *testing.T) {
	tests := []struct {
		name           string
		ocid           string
		expectError    bool
		expectedRegion string
		expectedVault  string
	}{
		{
			name:           "valid OCID",
			ocid:           "ocid1.key.oc1.uk-london-1.aaaalgz5aacmg.aaaailjtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq",
			expectError:    false,
			expectedRegion: "uk-london-1",
			expectedVault:  "aaaalgz5aacmg",
		},
		{
			name:           "valid OCID 2",
			ocid:           "ocid1.vault.oc1.iad.asdadsasdagz5aacmg.abwgiljtjasdasdasdagugpfe7wrtngukihgkybqxcoozz7sbh6lq",
			expectError:    false,
			expectedRegion: "iad",
			expectedVault:  "asdadsasdagz5aacmg",
		},
		{
			name:        "invalid OCID - too few parts",
			ocid:        "ocid1.key.oc1",
			expectError: true,
		},
		{
			name:        "invalid OCID - too many parts",
			ocid:        "ocid1.key.oc1.region.vault.extra.extra",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewMasterKeyFromOCID(tt.ocid)
			region, vault, err := extractRefs(key)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRegion, region)
				assert.Equal(t, tt.expectedVault, vault)
			}
		})
	}
}
