package barbican

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1*time.Second, config.InitialRetryDelay)
	assert.Equal(t, 30*time.Second, config.MaxRetryDelay)
	assert.Equal(t, 2.0, config.RetryMultiplier)
	assert.False(t, config.Insecure)
}

func TestNewBarbicanClient(t *testing.T) {
	// Create a test auth manager
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	authManager, err := NewAuthManager(config)
	require.NoError(t, err)
	
	tests := []struct {
		name        string
		endpoint    string
		authManager *AuthManager
		config      *ClientConfig
		expectError bool
	}{
		{
			name:        "Valid configuration",
			endpoint:    "https://barbican.example.com:9311",
			authManager: authManager,
			config:      DefaultClientConfig(),
			expectError: false,
		},
		{
			name:        "Valid configuration with nil config (uses default)",
			endpoint:    "https://barbican.example.com:9311",
			authManager: authManager,
			config:      nil,
			expectError: false,
		},
		{
			name:        "Empty endpoint",
			endpoint:    "",
			authManager: authManager,
			config:      DefaultClientConfig(),
			expectError: true,
		},
		{
			name:        "Nil auth manager",
			endpoint:    "https://barbican.example.com:9311",
			authManager: nil,
			config:      DefaultClientConfig(),
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewBarbicanClient(tt.endpoint, tt.authManager, tt.config)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.httpClient)
				assert.Equal(t, tt.authManager, client.authManager)
				
				// Check endpoint formatting
				expectedEndpoint := tt.endpoint
				if expectedEndpoint != "" && !strings.HasSuffix(expectedEndpoint, "/v1") {
					expectedEndpoint = expectedEndpoint + "/v1"
				}
				assert.Equal(t, expectedEndpoint, client.endpoint)
			}
		})
	}
}

func TestGetBarbicanEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		authManager *AuthManager
		region      string
		expectError bool
	}{
		{
			name: "Valid auth manager",
			authManager: &AuthManager{
				config: &AuthConfig{
					AuthURL: "https://keystone.example.com:5000/v3",
				},
			},
			region:      "sjc3",
			expectError: false,
		},
		{
			name:        "Nil auth manager",
			authManager: nil,
			region:      "sjc3",
			expectError: true,
		},
		{
			name: "Auth manager with nil config",
			authManager: &AuthManager{
				config: nil,
			},
			region:      "sjc3",
			expectError: true,
		},
		{
			name: "Auth manager with empty auth URL",
			authManager: &AuthManager{
				config: &AuthConfig{
					AuthURL: "",
				},
			},
			region:      "sjc3",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint, err := GetBarbicanEndpoint(tt.authManager, tt.region)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, endpoint)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, endpoint)
				assert.Contains(t, endpoint, "9311") // Barbican default port
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "Connection refused error",
			err:      fmt.Errorf("connection refused"),
			expected: true,
		},
		{
			name:     "Timeout error",
			err:      fmt.Errorf("request timeout"),
			expected: true,
		},
		{
			name:     "Server error 500",
			err:      fmt.Errorf("server error (500)"),
			expected: true,
		},
		{
			name:     "Server error 502",
			err:      fmt.Errorf("bad gateway (502)"),
			expected: true,
		},
		{
			name:     "Server error 503",
			err:      fmt.Errorf("service unavailable (503)"),
			expected: true,
		},
		{
			name:     "Server error 504",
			err:      fmt.Errorf("gateway timeout (504)"),
			expected: true,
		},
		{
			name:     "Authentication error",
			err:      fmt.Errorf("authentication failed (401)"),
			expected: true,
		},
		{
			name:     "Client error 400",
			err:      fmt.Errorf("bad request (400)"),
			expected: false,
		},
		{
			name:     "Not found error",
			err:      fmt.Errorf("not found (404)"),
			expected: false,
		},
		{
			name:     "Generic error",
			err:      fmt.Errorf("some other error"),
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMasterKeyGetBarbicanClient(t *testing.T) {
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
		SecretRef:    "550e8400-e29b-41d4-a716-446655440000",
		AuthConfig:   config,
		authManager:  authManager,
		baseEndpoint: "https://barbican.example.com:9311",
	}
	
	ctx := context.Background()
	client, err := key.getBarbicanClient(ctx)
	
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, client, key.client) // Should be cached
	
	// Test that subsequent calls return the same client
	client2, err := key.getBarbicanClient(ctx)
	assert.NoError(t, err)
	assert.Equal(t, client, client2)
}

// TestRetryLogicProperty tests Property 12: Retry Logic
// **Validates: Requirements 8.2**
func TestRetryLogicProperty(t *testing.T) {
	// Property-based test function
	f := func(maxRetries uint8, initialDelay uint8, multiplier uint8) bool {
		// Constrain inputs to very small ranges for fast execution
		if maxRetries > 2 {
			maxRetries = 2
		}
		if maxRetries == 0 {
			maxRetries = 1
		}
		if initialDelay == 0 {
			initialDelay = 10 // 10ms minimum
		}
		if initialDelay > 50 {
			initialDelay = 50 // 50ms maximum
		}
		if multiplier < 2 {
			multiplier = 2
		}
		if multiplier > 2 {
			multiplier = 2 // Fixed at 2x for predictability
		}

		// Track retry attempts and delays
		var attempts []time.Time
		var delays []time.Duration
		
		// Create a mock server that always returns retryable errors
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts = append(attempts, time.Now())
			w.WriteHeader(http.StatusInternalServerError) // Retryable error
			w.Write([]byte(`{"error": {"message": "Internal server error"}}`))
		}))
		defer server.Close()

		// Create a mock Keystone server for authentication
		keystoneServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/auth/tokens" {
				w.Header().Set("X-Subject-Token", "test-token")
				w.WriteHeader(http.StatusCreated)
				response := map[string]interface{}{
					"token": map[string]interface{}{
						"expires_at": time.Now().Add(time.Hour).Format(time.RFC3339),
						"project": map[string]interface{}{
							"id": "test-project",
						},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))
		defer keystoneServer.Close()

		// Create client configuration with fast test parameters
		config := &ClientConfig{
			Timeout:           1 * time.Second,  // Much shorter timeout
			MaxRetries:        int(maxRetries),
			InitialRetryDelay: time.Duration(initialDelay) * time.Millisecond,
			MaxRetryDelay:     200 * time.Millisecond, // Much shorter max delay
			RetryMultiplier:   float64(multiplier),
			Insecure:          true,
		}

		// Create auth manager
		authConfig := &AuthConfig{
			AuthURL:   keystoneServer.URL,
			Username:  "test-user",
			Password:  "test-password",
			ProjectID: "test-project",
		}
		
		authManager, err := NewAuthManager(authConfig)
		if err != nil {
			t.Logf("Failed to create auth manager: %v", err)
			return false
		}

		// Create Barbican client
		client, err := NewBarbicanClient(server.URL, authManager, config)
		if err != nil {
			t.Logf("Failed to create Barbican client: %v", err)
			return false
		}

		// Attempt an operation that will fail and trigger retries
		ctx := context.Background()
		metadata := SecretMetadata{
			Name:        "test-secret",
			SecretType:  "opaque",
			ContentType: "application/octet-stream",
		}
		
		_, err = client.StoreSecret(ctx, []byte("test-data"), metadata)
		
		// Should fail after all retries
		if err == nil {
			t.Logf("Expected error but got success")
			return false
		}

		// Verify the number of attempts (initial + retries)
		expectedAttempts := int(maxRetries) + 1
		if len(attempts) != expectedAttempts {
			t.Logf("Expected %d attempts, got %d. Error: %v", expectedAttempts, len(attempts), err)
			return false
		}

		// Calculate actual delays between attempts
		for i := 1; i < len(attempts); i++ {
			delay := attempts[i].Sub(attempts[i-1])
			delays = append(delays, delay)
		}

		// Verify exponential backoff pattern (with relaxed tolerance for fast execution)
		for i, delay := range delays {
			expectedDelay := time.Duration(float64(config.InitialRetryDelay) * 
				pow(config.RetryMultiplier, float64(i)))
			
			// Cap at MaxRetryDelay
			if expectedDelay > config.MaxRetryDelay {
				expectedDelay = config.MaxRetryDelay
			}

			// Allow generous tolerance for timing variations in fast tests
			tolerance := 50 * time.Millisecond
			if delay < expectedDelay-tolerance || delay > expectedDelay+tolerance*4 {
				t.Logf("Delay %d: expected ~%v, got %v (tolerance allows %v to %v)", 
					i+1, expectedDelay, delay, expectedDelay-tolerance, expectedDelay+tolerance*4)
				return false
			}
		}

		return true
	}

	// Run the property-based test with minimal iterations for fast execution
	if err := quick.Check(f, &quick.Config{MaxCount: 3}); err != nil {
		t.Error(err)
	}
}

// pow is a simple integer power function for calculating exponential backoff
func pow(base float64, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}