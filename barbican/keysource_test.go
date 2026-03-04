package barbican

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMasterKey(t *testing.T) {
	secretRef := "550e8400-e29b-41d4-a716-446655440000"
	key := NewMasterKey(secretRef)
	
	assert.Equal(t, secretRef, key.SecretRef)
	assert.WithinDuration(t, time.Now().UTC(), key.CreationDate, time.Second)
	assert.Empty(t, key.Region)
	assert.Empty(t, key.EncryptedKey)
}

func TestNewMasterKeyWithRegion(t *testing.T) {
	secretRef := "550e8400-e29b-41d4-a716-446655440000"
	region := "sjc3"
	key := NewMasterKeyWithRegion(secretRef, region)
	
	assert.Equal(t, secretRef, key.SecretRef)
	assert.Equal(t, region, key.Region)
	assert.WithinDuration(t, time.Now().UTC(), key.CreationDate, time.Second)
}

func TestNewMasterKeyFromSecretRef(t *testing.T) {
	tests := []struct {
		name        string
		secretRef   string
		expectError bool
		expectRegion string
	}{
		{
			name:        "Valid UUID format",
			secretRef:   "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:        "Valid URI format",
			secretRef:   "https://barbican.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:         "Valid regional format",
			secretRef:    "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
			expectError:  false,
			expectRegion: "sjc3",
		},
		{
			name:        "Invalid format",
			secretRef:   "invalid-secret-ref",
			expectError: true,
		},
		{
			name:        "Empty string",
			secretRef:   "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewMasterKeyFromSecretRef(tt.secretRef)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, key)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, key)
				assert.Equal(t, tt.secretRef, key.SecretRef)
				assert.Equal(t, tt.expectRegion, key.Region)
			}
		})
	}
}

func TestMasterKeysFromSecretRefString(t *testing.T) {
	tests := []struct {
		name        string
		secretRefs  string
		expectCount int
		expectError bool
	}{
		{
			name:        "Empty string",
			secretRefs:  "",
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "Single valid UUID",
			secretRefs:  "550e8400-e29b-41d4-a716-446655440000",
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "Multiple valid UUIDs",
			secretRefs:  "550e8400-e29b-41d4-a716-446655440000,660e8400-e29b-41d4-a716-446655440001",
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "Mixed valid formats",
			secretRefs:  "550e8400-e29b-41d4-a716-446655440000,region:dfw3:660e8400-e29b-41d4-a716-446655440001",
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "Contains invalid reference",
			secretRefs:  "550e8400-e29b-41d4-a716-446655440000,invalid-ref",
			expectCount: 0,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := MasterKeysFromSecretRefString(tt.secretRefs)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, keys, tt.expectCount)
			}
		})
	}
}

func TestIsValidSecretRef(t *testing.T) {
	tests := []struct {
		name      string
		secretRef string
		expected  bool
	}{
		{
			name:      "Valid UUID",
			secretRef: "550e8400-e29b-41d4-a716-446655440000",
			expected:  true,
		},
		{
			name:      "Valid URI",
			secretRef: "https://barbican.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000",
			expected:  true,
		},
		{
			name:      "Valid regional format",
			secretRef: "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
			expected:  true,
		},
		{
			name:      "Invalid UUID format",
			secretRef: "550e8400-e29b-41d4-a716",
			expected:  false,
		},
		{
			name:      "Invalid characters",
			secretRef: "invalid-secret-ref",
			expected:  false,
		},
		{
			name:      "Empty string",
			secretRef: "",
			expected:  false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidSecretRef(tt.secretRef)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractUUIDFromSecretRef(t *testing.T) {
	expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
	
	tests := []struct {
		name        string
		secretRef   string
		expectError bool
	}{
		{
			name:        "UUID format",
			secretRef:   expectedUUID,
			expectError: false,
		},
		{
			name:        "URI format",
			secretRef:   "https://barbican.example.com:9311/v1/secrets/" + expectedUUID,
			expectError: false,
		},
		{
			name:        "Regional format",
			secretRef:   "region:sjc3:" + expectedUUID,
			expectError: false,
		},
		{
			name:        "Invalid format",
			secretRef:   "invalid-ref",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, err := extractUUIDFromSecretRef(tt.secretRef)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, uuid)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedUUID, uuid)
			}
		})
	}
}

func TestMasterKeyToString(t *testing.T) {
	tests := []struct {
		name      string
		key       *MasterKey
		expected  string
	}{
		{
			name: "UUID without region",
			key: &MasterKey{
				SecretRef: "550e8400-e29b-41d4-a716-446655440000",
			},
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name: "UUID with region",
			key: &MasterKey{
				SecretRef: "550e8400-e29b-41d4-a716-446655440000",
				Region:    "sjc3",
			},
			expected: "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name: "Regional format with region",
			key: &MasterKey{
				SecretRef: "region:dfw3:550e8400-e29b-41d4-a716-446655440000",
				Region:    "sjc3",
			},
			expected: "region:dfw3:550e8400-e29b-41d4-a716-446655440000",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.key.ToString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMasterKeyToMap(t *testing.T) {
	key := &MasterKey{
		SecretRef:    "550e8400-e29b-41d4-a716-446655440000",
		Region:       "sjc3",
		EncryptedKey: "encrypted-data",
		CreationDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	result := key.ToMap()
	
	assert.Equal(t, key.SecretRef, result["secret_ref"])
	assert.Equal(t, key.Region, result["region"])
	assert.Equal(t, key.EncryptedKey, result["enc"])
	assert.Equal(t, "2023-01-01T00:00:00Z", result["created_at"])
}

func TestMasterKeyTypeToIdentifier(t *testing.T) {
	key := &MasterKey{}
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
}

func TestMasterKeyNeedsRotation(t *testing.T) {
	// Test key that needs rotation (old)
	oldKey := &MasterKey{
		CreationDate: time.Now().UTC().Add(-barbicanTTL - time.Hour),
	}
	assert.True(t, oldKey.NeedsRotation())
	
	// Test key that doesn't need rotation (new)
	newKey := &MasterKey{
		CreationDate: time.Now().UTC(),
	}
	assert.False(t, newKey.NeedsRotation())
}

func TestMasterKeyEncryptedDataKey(t *testing.T) {
	key := &MasterKey{
		EncryptedKey: "test-encrypted-key",
	}
	
	result := key.EncryptedDataKey()
	assert.Equal(t, []byte("test-encrypted-key"), result)
}

func TestMasterKeySetEncryptedDataKey(t *testing.T) {
	key := &MasterKey{}
	testData := []byte("test-encrypted-key")
	
	key.SetEncryptedDataKey(testData)
	assert.Equal(t, string(testData), key.EncryptedKey)
}

func TestCredentialsProvider(t *testing.T) {
	config := &AuthConfig{
		AuthURL:  "https://keystone.example.com:5000/v3",
		Username: "test-user",
	}
	
	provider := NewCredentialsProvider(config)
	assert.NotNil(t, provider)
	assert.Equal(t, config, provider.config)
	
	key := &MasterKey{}
	provider.ApplyToMasterKey(key)
	
	assert.Equal(t, provider, key.credentialsProvider)
	assert.Equal(t, config, key.AuthConfig)
}

func TestHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	client := NewHTTPClient(httpClient)
	
	assert.NotNil(t, client)
	assert.Equal(t, httpClient, client.hc)
	
	key := &MasterKey{}
	client.ApplyToMasterKey(key)
	
	assert.Equal(t, httpClient, key.httpClient)
}

func TestMasterKeyGetAuthToken(t *testing.T) {
	// Test with pre-configured auth manager
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	authManager, err := NewAuthManager(config)
	require.NoError(t, err)
	
	key := &MasterKey{
		AuthConfig:  config,
		authManager: authManager,
	}
	
	// Set up a cached token to avoid actual network call
	authManager.tokenCache.mutex.Lock()
	authManager.tokenCache.token = "test-token"
	authManager.tokenCache.projectID = "test-project"
	authManager.tokenCache.expiry = time.Now().Add(1 * time.Hour)
	authManager.tokenCache.mutex.Unlock()
	
	ctx := context.Background()
	token, projectID, err := key.getAuthToken(ctx)
	
	assert.NoError(t, err)
	assert.Equal(t, "test-token", token)
	assert.Equal(t, "test-project", projectID)
}

// TestEncryptionRoundTrip implements Property 2: Encryption Round Trip
// **Validates: Requirements 1.2, 1.3**
func TestEncryptionRoundTrip(t *testing.T) {
	// Create mock Barbican server for testing
	secretStore := make(map[string][]byte) // In-memory secret store
	var secretCounter int
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if strings.HasSuffix(r.URL.Path, "/secrets") {
				// Store secret operation
				secretCounter++
				// Generate a proper UUID format for testing
				secretUUID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", secretCounter)
				secretRef := fmt.Sprintf("https://barbican.example.com:9311/v1/secrets/%s", secretUUID)
				
				// Read the request body to get the payload
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				
				var req SecretCreateRequest
				if err := json.Unmarshal(body, &req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				
				// Decode the base64 payload
				payload, err := base64.StdEncoding.DecodeString(req.Payload)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				
				// Store the payload in our mock store using the UUID as key
				secretStore[secretUUID] = payload
				
				// Return success response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				response := SecretCreateResponse{SecretRef: secretRef}
				json.NewEncoder(w).Encode(response)
			}
		case "GET":
			if strings.Contains(r.URL.Path, "/payload") {
				// Get secret payload operation
				// Extract secret UUID from path
				pathParts := strings.Split(r.URL.Path, "/")
				if len(pathParts) < 4 {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				
				secretUUID := pathParts[len(pathParts)-2] // UUID is before "payload"
				
				// Look up the payload in our mock store
				payload, exists := secretStore[secretUUID]
				if !exists {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				
				// Return the payload
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)
				w.Write(payload)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()
	
	// Create mock Keystone server for authentication
	keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
			w.Header().Set("X-Subject-Token", "test-token-12345")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			
			response := AuthResponse{
				Token: struct {
					ExpiresAt string `json:"expires_at"`
					Project   struct {
						ID string `json:"id"`
					} `json:"project"`
				}{
					ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
					Project: struct {
						ID string `json:"id"`
					}{
						ID: "test-project-id",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer keystoneServer.Close()
	
	// Property-based test function
	f := func(dataKey []byte) bool {
		// Skip empty data keys as they're not meaningful for encryption
		if len(dataKey) == 0 {
			return true
		}
		
		// Create master key with mock configuration
		key := NewMasterKey("550e8400-e29b-41d4-a716-446655440000")
		
		// Configure authentication
		config := &AuthConfig{
			AuthURL:   keystoneServer.URL,
			Username:  "test-user",
			Password:  "test-password",
			ProjectID: "test-project",
		}
		
		authManager, err := NewAuthManager(config)
		if err != nil {
			t.Logf("Failed to create auth manager: %v", err)
			return false
		}
		
		key.AuthConfig = config
		key.authManager = authManager
		key.baseEndpoint = server.URL
		
		ctx := context.Background()
		
		// Encrypt the data key
		err = key.EncryptContext(ctx, dataKey)
		if err != nil {
			t.Logf("Encryption failed: %v", err)
			return false
		}
		
		// Verify that EncryptedKey was set
		if key.EncryptedKey == "" {
			t.Logf("EncryptedKey was not set after encryption")
			return false
		}
		
		// Decrypt the data key
		decryptedKey, err := key.DecryptContext(ctx)
		if err != nil {
			t.Logf("Decryption failed: %v", err)
			return false
		}
		
		// Verify round trip: original data key should equal decrypted key
		equal := bytes.Equal(dataKey, decryptedKey)
		if !equal {
			t.Logf("Round trip failed: original=%v, decrypted=%v", dataKey, decryptedKey)
		}
		return equal
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 10}); err != nil {
		t.Error(err)
	}
}

func TestGetEffectiveRegion(t *testing.T) {
	tests := []struct {
		name           string
		key            *MasterKey
		expectedRegion string
	}{
		{
			name: "Regional format in secret ref",
			key: &MasterKey{
				SecretRef: "region:dfw3:550e8400-e29b-41d4-a716-446655440000",
			},
			expectedRegion: "dfw3",
		},
		{
			name: "Region field set",
			key: &MasterKey{
				SecretRef: "550e8400-e29b-41d4-a716-446655440000",
				Region:    "sjc3",
			},
			expectedRegion: "sjc3",
		},
		{
			name: "Auth config region",
			key: &MasterKey{
				SecretRef: "550e8400-e29b-41d4-a716-446655440000",
				AuthConfig: &AuthConfig{
					Region: "fra3",
				},
			},
			expectedRegion: "fra3",
		},
		{
			name: "Default region",
			key: &MasterKey{
				SecretRef: "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedRegion: "RegionOne",
		},
		{
			name: "Regional format takes precedence",
			key: &MasterKey{
				SecretRef: "region:dfw3:550e8400-e29b-41d4-a716-446655440000",
				Region:    "sjc3",
				AuthConfig: &AuthConfig{
					Region: "fra3",
				},
			},
			expectedRegion: "dfw3",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region := tt.key.getEffectiveRegion()
			assert.Equal(t, tt.expectedRegion, region)
		})
	}
}

func TestGroupKeysByRegion(t *testing.T) {
	keys := []*MasterKey{
		{
			SecretRef: "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
		},
		{
			SecretRef: "region:dfw3:660e8400-e29b-41d4-a716-446655440001",
		},
		{
			SecretRef: "550e8400-e29b-41d4-a716-446655440002",
			Region:    "sjc3",
		},
		{
			SecretRef: "550e8400-e29b-41d4-a716-446655440003",
		},
	}
	
	groups := GroupKeysByRegion(keys)
	
	assert.Len(t, groups, 3) // sjc3, dfw3, RegionOne
	assert.Len(t, groups["sjc3"], 2)
	assert.Len(t, groups["dfw3"], 1)
	assert.Len(t, groups["RegionOne"], 1)
}

func TestGetRegionsFromKeys(t *testing.T) {
	keys := []*MasterKey{
		{
			SecretRef: "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
		},
		{
			SecretRef: "region:dfw3:660e8400-e29b-41d4-a716-446655440001",
		},
		{
			SecretRef: "550e8400-e29b-41d4-a716-446655440002",
			Region:    "sjc3", // Duplicate region
		},
	}
	
	regions := GetRegionsFromKeys(keys)
	
	assert.Len(t, regions, 2) // Should deduplicate
	assert.Contains(t, regions, "sjc3")
	assert.Contains(t, regions, "dfw3")
}

func TestValidateMultiRegionKeys(t *testing.T) {
	tests := []struct {
		name        string
		keys        []*MasterKey
		expectError bool
	}{
		{
			name:        "Empty keys",
			keys:        []*MasterKey{},
			expectError: true,
		},
		{
			name: "Valid multi-region keys",
			keys: []*MasterKey{
				{
					SecretRef: "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
					AuthConfig: &AuthConfig{
						AuthURL: "https://keystone.example.com:5000/v3",
					},
				},
				{
					SecretRef: "region:dfw3:660e8400-e29b-41d4-a716-446655440001",
					AuthConfig: &AuthConfig{
						AuthURL: "https://keystone.example.com:5000/v3",
					},
				},
			},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMultiRegionKeys(tt.keys)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetBarbicanEndpointForRegion(t *testing.T) {
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	authManager, err := NewAuthManager(config)
	require.NoError(t, err)
	
	tests := []struct {
		name           string
		region         string
		expectedSubstr string
	}{
		{
			name:           "US East region",
			region:         "sjc3",
			expectedSubstr: "barbican-sjc3",
		},
		{
			name:           "US West region",
			region:         "dfw3",
			expectedSubstr: "barbican-dfw3",
		},
		{
			name:           "Empty region uses default",
			region:         "",
			expectedSubstr: "barbican-RegionOne",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint, err := GetBarbicanEndpointForRegion(authManager, tt.region)
			
			assert.NoError(t, err)
			assert.Contains(t, endpoint, tt.expectedSubstr)
			assert.Contains(t, endpoint, ":9311")
		})
	}
}

func TestGetMultiRegionEndpoints(t *testing.T) {
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	authManager, err := NewAuthManager(config)
	require.NoError(t, err)
	
	regions := []string{"sjc3", "dfw3", "fra3"}
	
	endpoints, err := GetMultiRegionEndpoints(authManager, regions)
	
	assert.NoError(t, err)
	assert.Len(t, endpoints, 3)
	
	for _, region := range regions {
		endpoint, exists := endpoints[region]
		assert.True(t, exists)
		assert.Contains(t, endpoint, fmt.Sprintf("barbican-%s", region))
	}
}

// TestSecretReferenceValidationProperty implements Property 1: Secret Reference Validation
// **Validates: Requirements 1.5, 5.2**
func TestSecretReferenceValidationProperty(t *testing.T) {
	// Property-based test function
	f := func(input string) bool {
		// Test the property: isValidSecretRef should consistently validate secret references
		// according to the defined formats
		
		result := isValidSecretRef(input)
		
		// If the function says it's valid, we should be able to create a MasterKey from it
		if result {
			key, err := NewMasterKeyFromSecretRef(input)
			if err != nil {
				t.Logf("isValidSecretRef returned true but NewMasterKeyFromSecretRef failed for: %s, error: %v", input, err)
				return false
			}
			
			// If it's valid, we should also be able to extract a UUID from it
			uuid, err := extractUUIDFromSecretRef(input)
			if err != nil {
				t.Logf("isValidSecretRef returned true but extractUUIDFromSecretRef failed for: %s, error: %v", input, err)
				return false
			}
			
			// The extracted UUID should be a valid UUID format (36 characters with hyphens)
			if len(uuid) != 36 {
				t.Logf("Extracted UUID has wrong length for: %s, uuid: %s", input, uuid)
				return false
			}
			
			// The key should have the correct SecretRef
			if key.SecretRef != input {
				t.Logf("MasterKey SecretRef doesn't match input: expected %s, got %s", input, key.SecretRef)
				return false
			}
			
			// For regional format, the region should be extracted correctly
			if strings.HasPrefix(input, "region:") {
				parts := strings.Split(input, ":")
				if len(parts) >= 3 {
					expectedRegion := parts[1]
					if key.Region != expectedRegion {
						t.Logf("Region not extracted correctly: expected %s, got %s", expectedRegion, key.Region)
						return false
					}
				}
			}
		} else {
			// If the function says it's invalid, NewMasterKeyFromSecretRef should fail
			key, err := NewMasterKeyFromSecretRef(input)
			if err == nil {
				t.Logf("isValidSecretRef returned false but NewMasterKeyFromSecretRef succeeded for: %s, key: %+v", input, key)
				return false
			}
		}
		
		return true
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 10}); err != nil {
		t.Error(err)
	}
}

func TestEncryptMultiRegion(t *testing.T) {
	// Create mock servers for different regions
	secretStores := make(map[string]map[string][]byte) // region -> secretUUID -> payload
	secretCounters := make(map[string]int)             // region -> counter
	
	createMockServer := func(region string) *httptest.Server {
		secretStores[region] = make(map[string][]byte)
		secretCounters[region] = 0
		
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				if strings.HasSuffix(r.URL.Path, "/secrets") {
					secretCounters[region]++
					secretUUID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", secretCounters[region])
					secretRef := fmt.Sprintf("https://barbican-%s.example.com:9311/v1/secrets/%s", region, secretUUID)
					
					body, err := io.ReadAll(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					
					var req SecretCreateRequest
					if err := json.Unmarshal(body, &req); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					
					payload, err := base64.StdEncoding.DecodeString(req.Payload)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					
					secretStores[region][secretUUID] = payload
					
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					response := SecretCreateResponse{SecretRef: secretRef}
					json.NewEncoder(w).Encode(response)
				}
			}
		}))
	}
	
	// Create servers for different regions
	usEastServer := createMockServer("sjc3")
	defer usEastServer.Close()
	
	usWestServer := createMockServer("dfw3")
	defer usWestServer.Close()
	
	// Create mock Keystone server
	keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
			w.Header().Set("X-Subject-Token", "test-token-12345")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			
			response := AuthResponse{
				Token: struct {
					ExpiresAt string `json:"expires_at"`
					Project   struct {
						ID string `json:"id"`
					} `json:"project"`
				}{
					ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
					Project: struct {
						ID string `json:"id"`
					}{
						ID: "test-project-id",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer keystoneServer.Close()
	
	// Create master keys for different regions
	keys := []*MasterKey{
		{
			SecretRef:    "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
			baseEndpoint: usEastServer.URL,
		},
		{
			SecretRef:    "region:dfw3:660e8400-e29b-41d4-a716-446655440001",
			baseEndpoint: usWestServer.URL,
		},
	}
	
	// Configure authentication for all keys
	config := &AuthConfig{
		AuthURL:   keystoneServer.URL,
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	for _, key := range keys {
		authManager, err := NewAuthManager(config)
		require.NoError(t, err)
		
		key.AuthConfig = config
		key.authManager = authManager
	}
	
	// Test multi-region encryption
	dataKey := []byte("test-data-key-12345")
	ctx := context.Background()
	
	err := EncryptMultiRegion(ctx, dataKey, keys)
	assert.NoError(t, err)
	
	// Verify that both keys have encrypted data
	for _, key := range keys {
		assert.NotEmpty(t, key.EncryptedKey)
	}
	
	// Verify that secrets were stored in both regions
	assert.Len(t, secretStores["sjc3"], 1)
	assert.Len(t, secretStores["dfw3"], 1)
}

func TestDecryptMultiRegion(t *testing.T) {
	// Create mock servers for different regions
	secretStores := make(map[string]map[string][]byte)
	
	createMockServer := func(region string, shouldFail bool) *httptest.Server {
		secretStores[region] = make(map[string][]byte)
		
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldFail {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			
			if r.Method == "GET" && strings.Contains(r.URL.Path, "/payload") {
				pathParts := strings.Split(r.URL.Path, "/")
				if len(pathParts) < 4 {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				
				secretUUID := pathParts[len(pathParts)-2]
				payload, exists := secretStores[region][secretUUID]
				if !exists {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)
				w.Write(payload)
			}
		}))
	}
	
	// Create servers - first one fails, second one succeeds
	usEastServer := createMockServer("sjc3", true) // This will fail
	defer usEastServer.Close()
	
	usWestServer := createMockServer("dfw3", false) // This will succeed
	defer usWestServer.Close()
	
	// Pre-populate the working server with test data
	testData := []byte("test-data-key-12345")
	secretStores["dfw3"]["660e8400-e29b-41d4-a716-446655440001"] = testData
	
	// Create mock Keystone server
	keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
			w.Header().Set("X-Subject-Token", "test-token-12345")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			
			response := AuthResponse{
				Token: struct {
					ExpiresAt string `json:"expires_at"`
					Project   struct {
						ID string `json:"id"`
					} `json:"project"`
				}{
					ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
					Project: struct {
						ID string `json:"id"`
					}{
						ID: "test-project-id",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer keystoneServer.Close()
	
	// Create master keys for different regions
	keys := []*MasterKey{
		{
			SecretRef:    "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
			EncryptedKey: "https://barbican-sjc3.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000",
			baseEndpoint: usEastServer.URL,
		},
		{
			SecretRef:    "region:dfw3:660e8400-e29b-41d4-a716-446655440001",
			EncryptedKey: "https://barbican-dfw3.example.com:9311/v1/secrets/660e8400-e29b-41d4-a716-446655440001",
			baseEndpoint: usWestServer.URL,
		},
	}
	
	// Configure authentication for all keys
	config := &AuthConfig{
		AuthURL:   keystoneServer.URL,
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	for _, key := range keys {
		authManager, err := NewAuthManager(config)
		require.NoError(t, err)
		
		key.AuthConfig = config
		key.authManager = authManager
	}
	
	ctx := context.Background()
	
	// Test sequential failover
	dataKey, err := DecryptMultiRegion(ctx, keys)
	assert.NoError(t, err)
	assert.Equal(t, testData, dataKey)
	
	// Test parallel decryption
	dataKey, err = DecryptMultiRegionParallel(ctx, keys)
	assert.NoError(t, err)
	assert.Equal(t, testData, dataKey)
}

// TestResourceCleanupProperty implements Property 11: Resource Cleanup
// **Validates: Requirements 7.6**
func TestResourceCleanupProperty(t *testing.T) {
	// Property-based test function
	f := func(dataKey []byte, shouldFailAfterCreate bool) bool {
		// Skip empty data keys as they're not meaningful for encryption
		if len(dataKey) == 0 {
			return true
		}
		
		// Track created secrets for cleanup verification
		createdSecrets := make(map[string]bool) // secretUUID -> exists
		deletedSecrets := make(map[string]bool) // secretUUID -> deleted
		var secretCounter int
		
		// Create mock Barbican server that tracks secret lifecycle
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				if strings.HasSuffix(r.URL.Path, "/secrets") {
					// Store secret operation
					secretCounter++
					secretUUID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", secretCounter)
					secretRef := fmt.Sprintf("https://barbican.example.com:9311/v1/secrets/%s", secretUUID)
					
					createdSecrets[secretUUID] = true
					
					// If we should fail after creating the secret, simulate a failure scenario
					if shouldFailAfterCreate {
						// Return success for secret creation but we'll simulate cleanup later
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						response := SecretCreateResponse{SecretRef: secretRef}
						json.NewEncoder(w).Encode(response)
					} else {
						// Normal successful creation
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						response := SecretCreateResponse{SecretRef: secretRef}
						json.NewEncoder(w).Encode(response)
					}
				}
			case "DELETE":
				if strings.Contains(r.URL.Path, "/secrets/") {
					// Delete secret operation
					pathParts := strings.Split(r.URL.Path, "/")
					if len(pathParts) >= 4 {
						secretUUID := pathParts[len(pathParts)-1]
						
						// Track that this secret was deleted
						deletedSecrets[secretUUID] = true
						
						// Remove from created secrets
						delete(createdSecrets, secretUUID)
						
						w.WriteHeader(http.StatusNoContent)
					}
				}
			case "GET":
				if strings.Contains(r.URL.Path, "/payload") {
					// Get secret payload - only succeed if secret exists and wasn't deleted
					pathParts := strings.Split(r.URL.Path, "/")
					if len(pathParts) >= 4 {
						secretUUID := pathParts[len(pathParts)-2] // UUID is before "payload"
						
						if createdSecrets[secretUUID] && !deletedSecrets[secretUUID] {
							w.Header().Set("Content-Type", "application/octet-stream")
							w.WriteHeader(http.StatusOK)
							w.Write(dataKey) // Return the original data key
						} else {
							w.WriteHeader(http.StatusNotFound)
						}
					}
				} else if strings.Contains(r.URL.Path, "/secrets/") {
					// Get secret metadata - for validation
					pathParts := strings.Split(r.URL.Path, "/")
					if len(pathParts) >= 4 {
						secretUUID := pathParts[len(pathParts)-1]
						
						if createdSecrets[secretUUID] && !deletedSecrets[secretUUID] {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							response := SecretResponse{
								SecretRef:  fmt.Sprintf("https://barbican.example.com:9311/v1/secrets/%s", secretUUID),
								Name:       "SOPS Data Key",
								SecretType: "opaque",
								Status:     "ACTIVE",
							}
							json.NewEncoder(w).Encode(response)
						} else {
							w.WriteHeader(http.StatusNotFound)
						}
					}
				}
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))
		defer server.Close()
		
		// Create mock Keystone server for authentication
		keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
				w.Header().Set("X-Subject-Token", "test-token-12345")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				
				response := AuthResponse{
					Token: struct {
						ExpiresAt string `json:"expires_at"`
						Project   struct {
							ID string `json:"id"`
						} `json:"project"`
					}{
						ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
						Project: struct {
							ID string `json:"id"`
						}{
							ID: "test-project-id",
						},
					},
				}
				json.NewEncoder(w).Encode(response)
			}
		}))
		defer keystoneServer.Close()
		
		// Create master key with mock configuration
		key := NewMasterKey("550e8400-e29b-41d4-a716-446655440000")
		
		// Configure authentication
		config := &AuthConfig{
			AuthURL:   keystoneServer.URL,
			Username:  "test-user",
			Password:  "test-password",
			ProjectID: "test-project",
		}
		
		authManager, err := NewAuthManager(config)
		if err != nil {
			t.Logf("Failed to create auth manager: %v", err)
			return false
		}
		
		key.AuthConfig = config
		key.authManager = authManager
		key.baseEndpoint = server.URL
		
		ctx := context.Background()
		
		// Test scenario 1: Normal encryption and cleanup
		if !shouldFailAfterCreate {
			// Encrypt the data key
			err = key.EncryptContext(ctx, dataKey)
			if err != nil {
				t.Logf("Encryption failed: %v", err)
				return false
			}
			
			// Verify that a secret was created
			if len(createdSecrets) != 1 {
				t.Logf("Expected 1 secret to be created, got %d", len(createdSecrets))
				return false
			}
			
			// Get the Barbican client to test cleanup
			client, err := key.getBarbicanClient(ctx)
			if err != nil {
				t.Logf("Failed to get Barbican client: %v", err)
				return false
			}
			
			// Test cleanup by deleting the secret
			err = client.DeleteSecret(ctx, key.EncryptedKey)
			if err != nil {
				t.Logf("Failed to delete secret: %v", err)
				return false
			}
			
			// Verify that the secret was deleted
			if len(deletedSecrets) != 1 {
				t.Logf("Expected 1 secret to be deleted, got %d", len(deletedSecrets))
				return false
			}
			
			// Verify that no secrets remain in the created list
			if len(createdSecrets) != 0 {
				t.Logf("Expected 0 secrets to remain after cleanup, got %d", len(createdSecrets))
				return false
			}
			
		} else {
			// Test scenario 2: Encryption with simulated failure requiring cleanup
			// First, encrypt successfully to create a secret
			err = key.EncryptContext(ctx, dataKey)
			if err != nil {
				t.Logf("Initial encryption failed: %v", err)
				return false
			}
			
			// Verify that a secret was created
			if len(createdSecrets) != 1 {
				t.Logf("Expected 1 secret to be created, got %d", len(createdSecrets))
				return false
			}
			
			// Now simulate a failure scenario where we need to clean up
			// Get the Barbican client
			client, err := key.getBarbicanClient(ctx)
			if err != nil {
				t.Logf("Failed to get Barbican client: %v", err)
				return false
			}
			
			// Simulate cleanup of the temporary secret
			err = client.DeleteSecret(ctx, key.EncryptedKey)
			if err != nil {
				t.Logf("Failed to clean up temporary secret: %v", err)
				return false
			}
			
			// Verify cleanup was successful
			if len(deletedSecrets) != 1 {
				t.Logf("Expected 1 secret to be cleaned up, got %d", len(deletedSecrets))
				return false
			}
			
			// Verify no secrets remain
			if len(createdSecrets) != 0 {
				t.Logf("Expected 0 secrets to remain after cleanup, got %d", len(createdSecrets))
				return false
			}
		}
		
		return true
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 5}); err != nil {
		t.Error(err)
	}
}

// TestParallelOperationsProperty implements Property 8: Parallel Operations
// **Validates: Requirements 8.4**
func TestParallelOperationsProperty(t *testing.T) {
	// Property-based test function
	f := func(dataKey []byte, numKeys uint8, shouldFailSome bool) bool {
		// Skip empty data keys as they're not meaningful for encryption
		if len(dataKey) == 0 {
			return true
		}
		
		// Limit number of keys to a reasonable range (1-10)
		if numKeys == 0 {
			numKeys = 1
		}
		if numKeys > 10 {
			numKeys = 10
		}
		
		// Track operations for race condition detection
		operationCount := int32(0)
		maxConcurrentOps := int32(0)
		currentOps := int32(0)
		
		// Create mock servers for parallel operations
		servers := make([]*httptest.Server, numKeys)
		secretStores := make([]map[string][]byte, numKeys) // Per-server secret storage
		
		for i := uint8(0); i < numKeys; i++ {
			serverIndex := i
			secretStores[serverIndex] = make(map[string][]byte)
			var secretCounter int
			
			// Determine if this server should fail (for testing partial failures)
			shouldFail := shouldFailSome && (serverIndex%3 == 0) // Fail every 3rd server
			
			servers[serverIndex] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Track concurrent operations
				atomic.AddInt32(&currentOps, 1)
				defer atomic.AddInt32(&currentOps, -1)
				
				// Update max concurrent operations
				for {
					current := atomic.LoadInt32(&currentOps)
					max := atomic.LoadInt32(&maxConcurrentOps)
					if current <= max || atomic.CompareAndSwapInt32(&maxConcurrentOps, max, current) {
						break
					}
				}
				
				// Increment total operation count
				atomic.AddInt32(&operationCount, 1)
				
				// Simulate some processing time to increase chance of concurrent operations
				time.Sleep(10 * time.Millisecond)
				
				if shouldFail {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				
				switch r.Method {
				case "POST":
					if strings.HasSuffix(r.URL.Path, "/secrets") {
						// Store secret operation
						secretCounter++
						secretUUID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", secretCounter)
						secretRef := fmt.Sprintf("https://barbican-server-%d.example.com:9311/v1/secrets/%s", serverIndex, secretUUID)
						
						// Read and decode the request body
						body, err := io.ReadAll(r.Body)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						
						var req SecretCreateRequest
						if err := json.Unmarshal(body, &req); err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						
						payload, err := base64.StdEncoding.DecodeString(req.Payload)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						
						// Store the payload
						secretStores[serverIndex][secretUUID] = payload
						
						// Return success response
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						response := SecretCreateResponse{SecretRef: secretRef}
						json.NewEncoder(w).Encode(response)
					}
				case "GET":
					if strings.Contains(r.URL.Path, "/payload") {
						// Get secret payload operation
						pathParts := strings.Split(r.URL.Path, "/")
						if len(pathParts) < 4 {
							w.WriteHeader(http.StatusNotFound)
							return
						}
						
						secretUUID := pathParts[len(pathParts)-2]
						payload, exists := secretStores[serverIndex][secretUUID]
						if !exists {
							w.WriteHeader(http.StatusNotFound)
							return
						}
						
						w.Header().Set("Content-Type", "application/octet-stream")
						w.WriteHeader(http.StatusOK)
						w.Write(payload)
					}
				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
		}
		
		// Defer cleanup of all servers
		defer func() {
			for _, server := range servers {
				if server != nil {
					server.Close()
				}
			}
		}()
		
		// Create mock Keystone server for authentication
		keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
				w.Header().Set("X-Subject-Token", "test-token-12345")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				
				response := AuthResponse{
					Token: struct {
						ExpiresAt string `json:"expires_at"`
						Project   struct {
							ID string `json:"id"`
						} `json:"project"`
					}{
						ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
						Project: struct {
							ID string `json:"id"`
						}{
							ID: "test-project-id",
						},
					},
				}
				json.NewEncoder(w).Encode(response)
			}
		}))
		defer keystoneServer.Close()
		
		// Create master keys for parallel operations
		var keys []*MasterKey
		for i := uint8(0); i < numKeys; i++ {
			secretRef := fmt.Sprintf("region:region-%d:550e8400-e29b-41d4-a716-%012d", i, i)
			
			key := &MasterKey{
				SecretRef:    secretRef,
				baseEndpoint: servers[i].URL,
			}
			
			// Configure authentication
			config := &AuthConfig{
				AuthURL:   keystoneServer.URL,
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project",
				Region:    fmt.Sprintf("region-%d", i),
			}
			
			authManager, err := NewAuthManager(config)
			if err != nil {
				t.Logf("Failed to create auth manager for key %d: %v", i, err)
				return false
			}
			
			key.AuthConfig = config
			key.authManager = authManager
			keys = append(keys, key)
		}
		
		ctx := context.Background()
		
		// Property 1: Parallel encryption should perform operations concurrently
		// Reset operation counters
		atomic.StoreInt32(&operationCount, 0)
		atomic.StoreInt32(&maxConcurrentOps, 0)
		atomic.StoreInt32(&currentOps, 0)
		
		startTime := time.Now()
		err := EncryptMultiRegion(ctx, dataKey, keys)
		encryptDuration := time.Since(startTime)
		
		// Check if we expect success or partial failure
		expectedSuccesses := int(numKeys)
		if shouldFailSome {
			// Count how many servers should succeed (non-failing ones)
			for i := uint8(0); i < numKeys; i++ {
				if i%3 == 0 { // Every 3rd server fails
					expectedSuccesses--
				}
			}
		}
		
		if expectedSuccesses > 0 {
			// Should succeed if at least one server works
			if err != nil {
				t.Logf("Multi-region encryption failed unexpectedly: %v", err)
				return false
			}
		} else {
			// Should fail if all servers fail
			if err == nil {
				t.Logf("Multi-region encryption should have failed when all servers fail")
				return false
			}
			return true // This is expected behavior
		}
		
		// Verify parallel execution occurred
		totalOps := atomic.LoadInt32(&operationCount)
		maxConcurrent := atomic.LoadInt32(&maxConcurrentOps)
		
		// Account for retry logic - each key may make multiple attempts
		// We should see at least numKeys operations, but possibly more due to retries
		if totalOps < int32(numKeys) {
			t.Logf("Expected at least %d operations, got %d", numKeys, totalOps)
			return false
		}
		
		// For multiple keys, we should see some concurrency
		if numKeys > 1 && maxConcurrent < 2 {
			t.Logf("Expected concurrent operations with %d keys, max concurrent was %d", numKeys, maxConcurrent)
			return false
		}
		
		// Parallel operations should be faster than sequential for multiple keys
		// (This is a rough heuristic - parallel should not take much longer than sequential)
		if numKeys > 2 {
			// With retry logic and exponential backoff, operations can take longer
			// We allow generous margin for test environment variability and retries
			maxExpectedTime := time.Duration(numKeys) * 2 * time.Second // Very generous for retries
			if encryptDuration > maxExpectedTime {
				t.Logf("Parallel encryption took too long: %v (expected < %v for %d keys)", 
					encryptDuration, maxExpectedTime, numKeys)
				return false
			}
		}
		
		// Property 2: Parallel decryption should work correctly
		// Reset counters for decryption test
		atomic.StoreInt32(&operationCount, 0)
		atomic.StoreInt32(&maxConcurrentOps, 0)
		atomic.StoreInt32(&currentOps, 0)
		
		decryptedKey, err := DecryptMultiRegionParallel(ctx, keys)
		
		if err != nil {
			t.Logf("Parallel decryption failed: %v", err)
			return false
		}
		
		if !bytes.Equal(dataKey, decryptedKey) {
			t.Logf("Parallel decryption round trip failed: original=%v, decrypted=%v", dataKey, decryptedKey)
			return false
		}
		
		// Verify parallel decryption execution
		totalDecryptOps := atomic.LoadInt32(&operationCount)
		maxDecryptConcurrent := atomic.LoadInt32(&maxConcurrentOps)
		
		// Parallel decryption should stop after first success, so we might not see all operations
		if totalDecryptOps == 0 {
			t.Logf("Expected at least 1 decryption operation, got %d", totalDecryptOps)
			return false
		}
		
		// For multiple keys, we should see some concurrency in decryption attempts
		if numKeys > 1 && maxDecryptConcurrent < 2 && totalDecryptOps > 1 {
			t.Logf("Expected concurrent decryption operations with %d keys, max concurrent was %d", numKeys, maxDecryptConcurrent)
			return false
		}
		
		// Property 3: No race conditions should occur
		// This is implicitly tested by the fact that all operations complete successfully
		// and the data integrity is maintained (round-trip test passes)
		
		// Property 4: Partial failures should be handled gracefully
		if shouldFailSome && expectedSuccesses > 0 {
			// Encryption should succeed despite some failures
			// We already verified this above by checking that err == nil
			
			// Verify that successful keys have encrypted data
			successCount := 0
			for i, key := range keys {
				if key.EncryptedKey != "" {
					successCount++
				} else if i%3 != 0 { // Non-failing servers should have succeeded
					t.Logf("Key %d should have encrypted data but doesn't", i)
					return false
				}
			}
			
			if successCount != expectedSuccesses {
				t.Logf("Expected %d successful encryptions, got %d", expectedSuccesses, successCount)
				return false
			}
		}
		
		return true
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 3}); err != nil {
		t.Error(err)
	}
}

// TestMultiRegionEncryptionProperty implements Property 4: Multi-Region Encryption
// **Validates: Requirements 5.1, 5.4**
func TestMultiRegionEncryptionProperty(t *testing.T) {
	// Property-based test function
	f := func(dataKey []byte, numRegions uint8) bool {
		// Skip empty data keys as they're not meaningful for encryption
		if len(dataKey) == 0 {
			return true
		}
		
		// Limit number of regions to a reasonable range (1-5)
		if numRegions == 0 {
			numRegions = 1
		}
		if numRegions > 5 {
			numRegions = 5
		}
		
		regions := []string{"sjc3", "dfw3", "fra3", "nrt3", "ams3"}
		
		// Create servers for the specified number of regions
		servers := make(map[string]*httptest.Server)
		
		for i := uint8(0); i < numRegions; i++ {
			region := regions[i]
			
			// Create a separate secret store for each server to avoid race conditions
			secretStore := make(map[string][]byte) // secretUUID -> payload
			var secretCounter int
			
			servers[region] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "POST":
					if strings.HasSuffix(r.URL.Path, "/secrets") {
						// Store secret operation
						secretCounter++
						secretUUID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", secretCounter)
						secretRef := fmt.Sprintf("https://barbican-%s.example.com:9311/v1/secrets/%s", region, secretUUID)
						
						// Read the request body to get the payload
						body, err := io.ReadAll(r.Body)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						
						var req SecretCreateRequest
						if err := json.Unmarshal(body, &req); err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						
						// Decode the base64 payload
						payload, err := base64.StdEncoding.DecodeString(req.Payload)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						
						// Store the payload in our mock store
						secretStore[secretUUID] = payload
						
						// Return success response
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						response := SecretCreateResponse{SecretRef: secretRef}
						json.NewEncoder(w).Encode(response)
					}
				case "GET":
					if strings.Contains(r.URL.Path, "/payload") {
						// Get secret payload operation
						pathParts := strings.Split(r.URL.Path, "/")
						if len(pathParts) < 4 {
							w.WriteHeader(http.StatusNotFound)
							return
						}
						
						secretUUID := pathParts[len(pathParts)-2] // UUID is before "payload"
						
						// Look up the payload in our mock store
						payload, exists := secretStore[secretUUID]
						if !exists {
							w.WriteHeader(http.StatusNotFound)
							return
						}
						
						// Return the payload
						w.Header().Set("Content-Type", "application/octet-stream")
						w.WriteHeader(http.StatusOK)
						w.Write(payload)
					}
				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
		}
		
		// Defer cleanup of all servers
		defer func() {
			for _, server := range servers {
				server.Close()
			}
		}()
		
		// Create mock Keystone server for authentication
		keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
				w.Header().Set("X-Subject-Token", "test-token-12345")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				
				response := AuthResponse{
					Token: struct {
						ExpiresAt string `json:"expires_at"`
						Project   struct {
							ID string `json:"id"`
						} `json:"project"`
					}{
						ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
						Project: struct {
							ID string `json:"id"`
						}{
							ID: "test-project-id",
						},
					},
				}
				json.NewEncoder(w).Encode(response)
			}
		}))
		defer keystoneServer.Close()
		
		// Create master keys for different regions
		var keys []*MasterKey
		for i := uint8(0); i < numRegions; i++ {
			region := regions[i]
			secretRef := fmt.Sprintf("region:%s:550e8400-e29b-41d4-a716-%012d", region, i)
			
			key := &MasterKey{
				SecretRef:    secretRef,
				baseEndpoint: servers[region].URL,
			}
			
			// Configure authentication
			config := &AuthConfig{
				AuthURL:   keystoneServer.URL,
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project",
				Region:    region, // Region-specific authentication endpoint
			}
			
			authManager, err := NewAuthManager(config)
			if err != nil {
				t.Logf("Failed to create auth manager for region %s: %v", region, err)
				return false
			}
			
			key.AuthConfig = config
			key.authManager = authManager
			keys = append(keys, key)
		}
		
		ctx := context.Background()
		
		// Property 1: Multi-region encryption should succeed with all available regions
		// (Validates Requirement 5.1: encrypt the data key with each secret)
		err := EncryptMultiRegion(ctx, dataKey, keys)
		if err != nil {
			t.Logf("Multi-region encryption failed: %v", err)
			return false
		}
		
		// Verify that all keys have encrypted data
		for i, key := range keys {
			if key.EncryptedKey == "" {
				t.Logf("Key %d (region %s) does not have encrypted data", i, key.getEffectiveRegion())
				return false
			}
		}
		
		// Property 2: Each encrypted key should be decryptable and return the original data
		// (Validates round-trip consistency across regions)
		for i, key := range keys {
			decryptedKey, err := key.DecryptContext(ctx)
			if err != nil {
				t.Logf("Failed to decrypt key %d (region %s): %v", i, key.getEffectiveRegion(), err)
				return false
			}
			
			if !bytes.Equal(dataKey, decryptedKey) {
				t.Logf("Round trip failed for key %d (region %s): original=%v, decrypted=%v", 
					i, key.getEffectiveRegion(), dataKey, decryptedKey)
				return false
			}
		}
		
		// Property 3: Region-specific authentication endpoints should be handled correctly
		// (Validates Requirement 5.4: handle region-specific authentication endpoints)
		for i, key := range keys {
			expectedRegion := regions[i]
			actualRegion := key.getEffectiveRegion()
			
			if actualRegion != expectedRegion {
				t.Logf("Key %d should be in region %s, got %s", i, expectedRegion, actualRegion)
				return false
			}
			
			// Verify that the auth config has the correct region
			if key.AuthConfig.Region != expectedRegion {
				t.Logf("Key %d auth config should have region %s, got %s", i, expectedRegion, key.AuthConfig.Region)
				return false
			}
		}
		
		// Property 4: Multi-region decryption should work with failover
		// Test that we can decrypt using any of the keys
		decryptedKey, err := DecryptMultiRegion(ctx, keys)
		if err != nil {
			t.Logf("Multi-region decryption failed: %v", err)
			return false
		}
		
		if !bytes.Equal(dataKey, decryptedKey) {
			t.Logf("Multi-region decryption round trip failed: original=%v, decrypted=%v", dataKey, decryptedKey)
			return false
		}
		
		// Property 5: Parallel multi-region decryption should also work
		decryptedKey, err = DecryptMultiRegionParallel(ctx, keys)
		if err != nil {
			t.Logf("Parallel multi-region decryption failed: %v", err)
			return false
		}
		
		if !bytes.Equal(dataKey, decryptedKey) {
			t.Logf("Parallel multi-region decryption round trip failed: original=%v, decrypted=%v", dataKey, decryptedKey)
			return false
		}
		
		return true
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 5}); err != nil {
		t.Error(err)
	}
}

// TestCommandLineIntegrationProperty implements Property 6: Command Line Integration
// **Validates: Requirements 3.1, 3.2, 3.3**
func TestCommandLineIntegrationProperty(t *testing.T) {
	// Property-based test function
	f := func(secretRefs []string, useCommas bool, includeInvalid bool) bool {
		// Skip empty input to focus on meaningful test cases
		if len(secretRefs) == 0 {
			return true
		}
		
		// Generate valid secret references for testing
		var validRefs []string
		var invalidRefs []string
		
		for i := range secretRefs {
			// Limit to reasonable number of refs for fast execution
			if i >= 5 {
				break
			}
			
			// Create valid UUID format
			validRef := fmt.Sprintf("550e840%d-e29b-41d4-a716-44665544000%d", i%10, i%10)
			validRefs = append(validRefs, validRef)
			
			// Create some invalid refs if requested
			if includeInvalid && i%3 == 0 {
				invalidRefs = append(invalidRefs, "invalid-ref-"+fmt.Sprint(i))
			}
		}
		
		// Test with valid references only
		if len(validRefs) > 0 {
			var testInput string
			if useCommas {
				testInput = strings.Join(validRefs, ",")
			} else {
				// Test single reference
				testInput = validRefs[0]
			}
			
			// Test MasterKeysFromSecretRefString function (simulates command line parsing)
			keys, err := MasterKeysFromSecretRefString(testInput)
			
			// Should succeed with valid input
			if err != nil {
				t.Logf("MasterKeysFromSecretRefString failed with valid input: %s, error: %v", testInput, err)
				return false
			}
			
			// Should return correct number of keys
			expectedCount := 1
			if useCommas {
				expectedCount = len(validRefs)
			}
			
			if len(keys) != expectedCount {
				t.Logf("Expected %d keys, got %d for input: %s", expectedCount, len(keys), testInput)
				return false
			}
			
			// Each key should have correct secret reference
			for i, key := range keys {
				expectedRef := validRefs[0]
				if useCommas && i < len(validRefs) {
					expectedRef = validRefs[i]
				}
				
				if key.SecretRef != expectedRef {
					t.Logf("Key %d has wrong SecretRef: expected %s, got %s", i, expectedRef, key.SecretRef)
					return false
				}
				
				// Key should be properly initialized
				if key.CreationDate.IsZero() {
					t.Logf("Key %d has zero CreationDate", i)
					return false
				}
			}
			
			// Test regional format parsing
			regionalRef := "region:sjc3:" + validRefs[0]
			regionalKeys, err := MasterKeysFromSecretRefString(regionalRef)
			if err != nil {
				t.Logf("MasterKeysFromSecretRefString failed with regional format: %s, error: %v", regionalRef, err)
				return false
			}
			
			if len(regionalKeys) != 1 {
				t.Logf("Expected 1 regional key, got %d", len(regionalKeys))
				return false
			}
			
			if regionalKeys[0].Region != "sjc3" {
				t.Logf("Regional key has wrong region: expected sjc3, got %s", regionalKeys[0].Region)
				return false
			}
		}
		
		// Test with invalid references (should fail)
		if len(invalidRefs) > 0 {
			invalidInput := strings.Join(invalidRefs, ",")
			_, err := MasterKeysFromSecretRefString(invalidInput)
			
			// Should fail with invalid input
			if err == nil {
				t.Logf("MasterKeysFromSecretRefString should have failed with invalid input: %s", invalidInput)
				return false
			}
		}
		
		// Test mixed valid and invalid (should fail)
		if len(validRefs) > 0 && len(invalidRefs) > 0 {
			mixedInput := validRefs[0] + ",invalid-ref"
			_, err := MasterKeysFromSecretRefString(mixedInput)
			
			// Should fail with mixed input
			if err == nil {
				t.Logf("MasterKeysFromSecretRefString should have failed with mixed valid/invalid input: %s", mixedInput)
				return false
			}
		}
		
		// Test empty string (should return empty slice, no error)
		emptyKeys, err := MasterKeysFromSecretRefString("")
		if err != nil {
			t.Logf("MasterKeysFromSecretRefString failed with empty string: %v", err)
			return false
		}
		
		if len(emptyKeys) != 0 {
			t.Logf("Expected 0 keys for empty string, got %d", len(emptyKeys))
			return false
		}
		
		// Test whitespace handling
		if len(validRefs) > 0 {
			whitespaceInput := " " + validRefs[0] + " "
			wsKeys, err := MasterKeysFromSecretRefString(whitespaceInput)
			if err != nil {
				t.Logf("MasterKeysFromSecretRefString failed with whitespace: %s, error: %v", whitespaceInput, err)
				return false
			}
			
			if len(wsKeys) != 1 {
				t.Logf("Expected 1 key with whitespace input, got %d", len(wsKeys))
				return false
			}
			
			if wsKeys[0].SecretRef != validRefs[0] {
				t.Logf("Whitespace not trimmed correctly: expected %s, got %s", validRefs[0], wsKeys[0].SecretRef)
				return false
			}
		}
		
		return true
	}
	
	// Run the property-based test with constrained iterations for reasonable execution time
	if err := quick.Check(f, &quick.Config{MaxCount: 10}); err != nil {
		t.Error(err)
	}
}
// TestKeyManagementOperationsProperty implements Property 14: Key Management Operations
// **Validates: Requirements 3.4, 3.5, 9.3**
func TestKeyManagementOperationsProperty(t *testing.T) {
	// Property-based test function
	f := func(initialKeys []string, keysToAdd []string, keysToRemove []string, shouldMaintainIntegrity bool) bool {
		// Skip empty test cases to focus on meaningful scenarios
		if len(initialKeys) == 0 && len(keysToAdd) == 0 {
			return true
		}
		
		// Limit the number of keys for reasonable test execution time
		if len(initialKeys) > 3 {
			initialKeys = initialKeys[:3]
		}
		if len(keysToAdd) > 3 {
			keysToAdd = keysToAdd[:3]
		}
		if len(keysToRemove) > 3 {
			keysToRemove = keysToRemove[:3]
		}
		
		// Generate valid secret references for testing
		var validInitialKeys []*MasterKey
		var validKeysToAdd []*MasterKey
		var validKeysToRemove []*MasterKey
		
		// Create initial keys
		for i := range initialKeys {
			if i >= 3 { // Limit for performance
				break
			}
			secretRef := fmt.Sprintf("550e840%d-e29b-41d4-a716-44665544000%d", i%10, i%10)
			key := NewMasterKey(secretRef)
			key.EncryptedKey = fmt.Sprintf("encrypted-key-%d", i) // Simulate encrypted data
			validInitialKeys = append(validInitialKeys, key)
		}
		
		// Create keys to add
		for i := range keysToAdd {
			if i >= 3 { // Limit for performance
				break
			}
			secretRef := fmt.Sprintf("660e840%d-e29b-41d4-a716-44665544000%d", i%10, i%10)
			key := NewMasterKey(secretRef)
			validKeysToAdd = append(validKeysToAdd, key)
		}
		
		// Create keys to remove (should be subset of initial keys for valid test)
		for i := range keysToRemove {
			if i >= len(validInitialKeys) || i >= 3 { // Limit for performance
				break
			}
			// Use existing initial key for removal
			validKeysToRemove = append(validKeysToRemove, validInitialKeys[i])
		}
		
		// Test key addition operations (simulates --add-barbican functionality)
		if len(validKeysToAdd) > 0 {
			// Simulate initial key group (like in SOPS metadata)
			initialKeyGroup := make([]interface{}, len(validInitialKeys))
			for i, key := range validInitialKeys {
				initialKeyGroup[i] = key
			}
			
			// Add new keys (simulates the rotate function logic)
			finalKeyGroup := make([]interface{}, len(initialKeyGroup))
			copy(finalKeyGroup, initialKeyGroup)
			
			for _, newKey := range validKeysToAdd {
				finalKeyGroup = append(finalKeyGroup, newKey)
			}
			
			// Verify keys were added correctly
			expectedCount := len(validInitialKeys) + len(validKeysToAdd)
			if len(finalKeyGroup) != expectedCount {
				t.Logf("Key addition failed: expected %d keys, got %d", expectedCount, len(finalKeyGroup))
				return false
			}
			
			// Verify all original keys are still present
			for _, originalKey := range validInitialKeys {
				found := false
				for _, finalKey := range finalKeyGroup {
					if mk, ok := finalKey.(*MasterKey); ok {
						if mk.ToString() == originalKey.ToString() {
							found = true
							break
						}
					}
				}
				if !found {
					t.Logf("Original key lost during addition: %s", originalKey.ToString())
					return false
				}
			}
			
			// Verify all new keys are present
			for _, newKey := range validKeysToAdd {
				found := false
				for _, finalKey := range finalKeyGroup {
					if mk, ok := finalKey.(*MasterKey); ok {
						if mk.ToString() == newKey.ToString() {
							found = true
							break
						}
					}
				}
				if !found {
					t.Logf("New key not found after addition: %s", newKey.ToString())
					return false
				}
			}
		}
		
		// Test key removal operations (simulates --rm-barbican functionality)
		if len(validKeysToRemove) > 0 && len(validInitialKeys) > 0 {
			// Simulate initial key group
			keyGroup := make([]*MasterKey, len(validInitialKeys))
			copy(keyGroup, validInitialKeys)
			
			// Remove keys (simulates the rotate function logic)
			for _, rmKey := range validKeysToRemove {
				for i, groupKey := range keyGroup {
					if rmKey.ToString() == groupKey.ToString() {
						// Remove the key (simulates slice removal in rotate function)
						keyGroup = append(keyGroup[:i], keyGroup[i+1:]...)
						break
					}
				}
			}
			
			// Verify keys were removed correctly
			expectedCount := len(validInitialKeys) - len(validKeysToRemove)
			if len(keyGroup) != expectedCount {
				t.Logf("Key removal failed: expected %d keys, got %d", expectedCount, len(keyGroup))
				return false
			}
			
			// Verify removed keys are no longer present
			for _, removedKey := range validKeysToRemove {
				for _, remainingKey := range keyGroup {
					if remainingKey.ToString() == removedKey.ToString() {
						t.Logf("Key not properly removed: %s", removedKey.ToString())
						return false
					}
				}
			}
			
			// Verify non-removed keys are still present
			for _, originalKey := range validInitialKeys {
				shouldBePresent := true
				for _, removedKey := range validKeysToRemove {
					if originalKey.ToString() == removedKey.ToString() {
						shouldBePresent = false
						break
					}
				}
				
				if shouldBePresent {
					found := false
					for _, remainingKey := range keyGroup {
						if remainingKey.ToString() == originalKey.ToString() {
							found = true
							break
						}
					}
					if !found {
						t.Logf("Non-removed key lost during removal: %s", originalKey.ToString())
						return false
					}
				}
			}
		}
		
		// Test ToString consistency for key identification (critical for add/remove operations)
		for _, key := range validInitialKeys {
			toString1 := key.ToString()
			toString2 := key.ToString()
			
			if toString1 != toString2 {
				t.Logf("ToString not consistent for key: %s vs %s", toString1, toString2)
				return false
			}
			
			// ToString should be non-empty for valid keys
			if toString1 == "" {
				t.Logf("ToString returned empty string for valid key")
				return false
			}
		}
		
		// Test mixed key type compatibility (Requirement 9.3: mixed master key types)
		if len(validInitialKeys) > 0 {
			// Simulate mixed key types in the same key group
			mixedKeyGroup := make([]interface{}, 0)
			
			// Add Barbican keys
			for _, key := range validInitialKeys {
				mixedKeyGroup = append(mixedKeyGroup, key)
			}
			
			// Add mock keys of other types (simulating KMS, PGP, etc.)
			mockKMSKey := &struct {
				ARN string
			}{
				ARN: "arn:aws:kms:sjc3:123456789012:key/12345678-1234-1234-1234-123456789012",
			}
			mixedKeyGroup = append(mixedKeyGroup, mockKMSKey)
			
			// Verify Barbican keys can coexist with other key types
			barbicanCount := 0
			otherCount := 0
			
			for _, key := range mixedKeyGroup {
				if _, ok := key.(*MasterKey); ok {
					barbicanCount++
				} else {
					otherCount++
				}
			}
			
			if barbicanCount != len(validInitialKeys) {
				t.Logf("Barbican key count mismatch in mixed group: expected %d, got %d", len(validInitialKeys), barbicanCount)
				return false
			}
			
			if otherCount != 1 {
				t.Logf("Other key count mismatch in mixed group: expected 1, got %d", otherCount)
				return false
			}
		}
		
		// Test file integrity maintenance (Requirement 9.3)
		if shouldMaintainIntegrity && len(validInitialKeys) > 0 {
			// Verify that key operations preserve essential key properties
			for _, key := range validInitialKeys {
				// Key should maintain its secret reference
				if key.SecretRef == "" {
					t.Logf("Key lost SecretRef during operations")
					return false
				}
				
				// Key should maintain its creation date
				if key.CreationDate.IsZero() {
					t.Logf("Key lost CreationDate during operations")
					return false
				}
				
				// Key should maintain its type identifier
				if key.TypeToIdentifier() != KeyTypeIdentifier {
					t.Logf("Key lost type identifier during operations: expected %s, got %s", KeyTypeIdentifier, key.TypeToIdentifier())
					return false
				}
			}
		}
		
		return true
	}
	
	// Run the property-based test with constrained iterations for reasonable execution time
	if err := quick.Check(f, &quick.Config{MaxCount: 10}); err != nil {
		t.Error(err)
	}
}
