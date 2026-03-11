package barbican

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthManager(t *testing.T) {
	tests := []struct {
		name        string
		config      *AuthConfig
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: &AuthConfig{
				AuthURL:  "https://keystone.example.com:5000/v3",
				Username: "test-user",
				Password: "test-password",
				ProjectID: "test-project",
			},
			expectError: false,
		},
		{
			name:        "Nil configuration",
			config:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewAuthManager(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, manager)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, manager)
				assert.Equal(t, tt.config, manager.config)
				assert.NotNil(t, manager.httpClient)
				assert.NotNil(t, manager.tokenCache)
			}
		})
	}
}

func TestCreateHTTPClient(t *testing.T) {
	tests := []struct {
		name   string
		config *AuthConfig
	}{
		{
			name: "Default configuration",
			config: &AuthConfig{
				Insecure: false,
			},
		},
		{
			name: "Insecure configuration",
			config: &AuthConfig{
				Insecure: true,
			},
		},
		{
			name: "With CA certificate content",
			config: &AuthConfig{
				CACert: "-----BEGIN CERTIFICATE-----\ntest-cert-content\n-----END CERTIFICATE-----",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := createHTTPClient(tt.config)
			
			// Note: We expect error for invalid cert content in test, but function should not panic
			if tt.name == "With CA certificate content" {
				// This will fail because it's not a valid certificate, but that's expected in test
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, 30*time.Second, client.Timeout)
			}
		})
	}
}

func TestBuildAuthRequest(t *testing.T) {
	tests := []struct {
		name        string
		config      *AuthConfig
		expectError bool
		expectMethod string
	}{
		{
			name: "Application credential authentication",
			config: &AuthConfig{
				ApplicationCredentialID:     "app-cred-id",
				ApplicationCredentialSecret: "app-cred-secret",
			},
			expectError:  false,
			expectMethod: "application_credential",
		},
		{
			name: "Token authentication",
			config: &AuthConfig{
				Token:     "existing-token",
				ProjectID: "test-project",
			},
			expectError:  false,
			expectMethod: "token",
		},
		{
			name: "Password authentication with project ID",
			config: &AuthConfig{
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project",
			},
			expectError:  false,
			expectMethod: "password",
		},
		{
			name: "Password authentication with project name",
			config: &AuthConfig{
				Username:    "test-user",
				Password:    "test-password",
				ProjectName: "test-project",
				DomainName:  "default",
			},
			expectError:  false,
			expectMethod: "password",
		},
		{
			name: "No authentication method",
			config: &AuthConfig{
				AuthURL: "https://keystone.example.com:5000/v3",
			},
			expectError: true,
		},
		{
			name: "Password auth without project scope",
			config: &AuthConfig{
				Username: "test-user",
				Password: "test-password",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &AuthManager{config: tt.config}
			
			authReq, err := manager.buildAuthRequest()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, authReq)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authReq)
				assert.Contains(t, authReq.Auth.Identity.Methods, tt.expectMethod)

				// Verify scope is set correctly for non-app-cred auth
				if tt.expectMethod != "application_credential" {
					assert.NotNil(t, authReq.Auth.Scope)
					assert.NotNil(t, authReq.Auth.Scope.Project)
				}
			}
		})
	}
}

func TestAuthManagerAuthenticate(t *testing.T) {
	// Create a mock Keystone server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/auth/tokens", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Mock successful response
		w.Header().Set("X-Subject-Token", "test-token-12345")
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
	}))
	defer server.Close()

	config := &AuthConfig{
		AuthURL:   server.URL,
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}

	manager, err := NewAuthManager(config)
	require.NoError(t, err)

	ctx := context.Background()
	token, projectID, err := manager.authenticate(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "test-token-12345", token)
	assert.Equal(t, "test-project-id", projectID)

	// Verify token is cached
	manager.tokenCache.mutex.RLock()
	assert.Equal(t, "test-token-12345", manager.tokenCache.token)
	assert.Equal(t, "test-project-id", manager.tokenCache.projectID)
	assert.True(t, time.Now().Before(manager.tokenCache.expiry))
	manager.tokenCache.mutex.RUnlock()
}

func TestAuthManagerGetToken(t *testing.T) {
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}

	manager, err := NewAuthManager(config)
	require.NoError(t, err)

	// Test with cached token
	manager.tokenCache.mutex.Lock()
	manager.tokenCache.token = "cached-token"
	manager.tokenCache.projectID = "cached-project"
	manager.tokenCache.expiry = time.Now().Add(1 * time.Hour)
	manager.tokenCache.mutex.Unlock()

	ctx := context.Background()
	token, projectID, err := manager.GetToken(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "cached-token", token)
	assert.Equal(t, "cached-project", projectID)
}

func TestAuthManagerInvalidateToken(t *testing.T) {
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}

	manager, err := NewAuthManager(config)
	require.NoError(t, err)

	// Set a cached token
	manager.tokenCache.mutex.Lock()
	manager.tokenCache.token = "test-token"
	manager.tokenCache.projectID = "test-project"
	manager.tokenCache.expiry = time.Now().Add(1 * time.Hour)
	manager.tokenCache.mutex.Unlock()

	// Invalidate the token
	manager.InvalidateToken()

	// Verify token is cleared
	manager.tokenCache.mutex.RLock()
	assert.Empty(t, manager.tokenCache.token)
	assert.Empty(t, manager.tokenCache.projectID)
	assert.True(t, manager.tokenCache.expiry.IsZero())
	manager.tokenCache.mutex.RUnlock()
}

func TestLoadConfigFromEnvironment(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"OS_AUTH_URL", "OS_REGION_NAME", "OS_PROJECT_ID", "OS_PROJECT_NAME",
		"OS_DOMAIN_ID", "OS_DOMAIN_NAME", "OS_USERNAME", "OS_PASSWORD",
		"OS_APPLICATION_CREDENTIAL_ID", "OS_APPLICATION_CREDENTIAL_SECRET",
		"OS_TOKEN", "OS_INSECURE", "OS_CACERT",
	}

	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
		os.Unsetenv(env)
	}

	// Restore environment after test
	defer func() {
		for _, env := range envVars {
			if val, exists := originalEnv[env]; exists {
				os.Setenv(env, val)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	// Set test environment variables
	testEnv := map[string]string{
		"OS_AUTH_URL":                        "https://keystone.example.com:5000/v3",
		"OS_REGION_NAME":                     "sjc3",
		"OS_PROJECT_ID":                      "test-project-id",
		"OS_USERNAME":                        "test-user",
		"OS_PASSWORD":                        "test-password",
		"OS_APPLICATION_CREDENTIAL_ID":       "app-cred-id",
		"OS_APPLICATION_CREDENTIAL_SECRET":   "app-cred-secret",
		"OS_INSECURE":                        "true",
		"OS_CACERT":                          "/path/to/ca.pem",
	}

	for key, value := range testEnv {
		os.Setenv(key, value)
	}

	config := LoadConfigFromEnvironment()

	assert.Equal(t, "https://keystone.example.com:5000/v3", config.AuthURL)
	assert.Equal(t, "sjc3", config.Region)
	assert.Equal(t, "test-project-id", config.ProjectID)
	assert.Equal(t, "test-user", config.Username)
	assert.Equal(t, "test-password", config.Password)
	assert.Equal(t, "app-cred-id", config.ApplicationCredentialID)
	assert.Equal(t, "app-cred-secret", config.ApplicationCredentialSecret)
	assert.True(t, config.Insecure)
	assert.Equal(t, "/path/to/ca.pem", config.CACert)
}

func TestLoadConfigFromEnvironmentDefaults(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"OS_AUTH_URL", "OS_REGION_NAME", "OS_PROJECT_ID", "OS_PROJECT_NAME",
		"OS_DOMAIN_ID", "OS_DOMAIN_NAME", "OS_USERNAME", "OS_PASSWORD",
		"OS_APPLICATION_CREDENTIAL_ID", "OS_APPLICATION_CREDENTIAL_SECRET",
		"OS_TOKEN", "OS_INSECURE", "OS_CACERT",
	}

	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
		os.Unsetenv(env)
	}

	// Restore environment after test
	defer func() {
		for _, env := range envVars {
			if val, exists := originalEnv[env]; exists {
				os.Setenv(env, val)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	// Set minimal environment
	os.Setenv("OS_USERNAME", "test-user")

	config := LoadConfigFromEnvironment()

	// Check defaults
	assert.Equal(t, "RegionOne", config.Region)
	assert.Equal(t, "default", config.DomainName)
	assert.False(t, config.Insecure)
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *AuthConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Nil configuration",
			config:      nil,
			expectError: true,
			errorMsg:    "Authentication configuration is required",
		},
		{
			name: "Missing auth URL",
			config: &AuthConfig{
				Username: "test-user",
				Password: "test-password",
			},
			expectError: true,
			errorMsg:    "Authentication URL is required",
		},
		{
			name: "No authentication method",
			config: &AuthConfig{
				AuthURL: "https://keystone.example.com:5000/v3",
			},
			expectError: true,
			errorMsg:    "No valid authentication method provided",
		},
		{
			name: "Password auth without project scope",
			config: &AuthConfig{
				AuthURL:  "https://keystone.example.com:5000/v3",
				Username: "test-user",
				Password: "test-password",
			},
			expectError: true,
			errorMsg:    "Project scope is required",
		},
		{
			name: "Valid password authentication",
			config: &AuthConfig{
				AuthURL:   "https://keystone.example.com:5000/v3",
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project",
			},
			expectError: false,
		},
		{
			name: "Valid application credential authentication",
			config: &AuthConfig{
				AuthURL:                     "https://keystone.example.com:5000/v3",
				ApplicationCredentialID:     "app-cred-id",
				ApplicationCredentialSecret: "app-cred-secret",
			},
			expectError: false,
		},
		{
			name: "Valid token authentication",
			config: &AuthConfig{
				AuthURL:   "https://keystone.example.com:5000/v3",
				Token:     "existing-token",
				ProjectID: "test-project",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthManagerAuthenticationFailure(t *testing.T) {
	// Create a mock server that returns authentication failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": {"message": "Invalid credentials"}}`))
	}))
	defer server.Close()

	config := &AuthConfig{
		AuthURL:   server.URL,
		Username:  "invalid-user",
		Password:  "invalid-password",
		ProjectID: "test-project",
	}

	manager, err := NewAuthManager(config)
	require.NoError(t, err)

	ctx := context.Background()
	token, projectID, err := manager.authenticate(ctx)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Empty(t, projectID)
	assert.Contains(t, err.Error(), "authentication failed with status 401")
}

func TestAuthManagerMissingToken(t *testing.T) {
	// Create a mock server that doesn't return X-Subject-Token header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	config := &AuthConfig{
		AuthURL:   server.URL,
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}

	manager, err := NewAuthManager(config)
	require.NoError(t, err)

	ctx := context.Background()
	token, projectID, err := manager.authenticate(ctx)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Empty(t, projectID)
	assert.Contains(t, err.Error(), "No authentication token received")
}

func TestTokenCacheExpiration(t *testing.T) {
	config := &AuthConfig{
		AuthURL:   "https://keystone.example.com:5000/v3",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}

	manager, err := NewAuthManager(config)
	require.NoError(t, err)

	// Set an expired token
	manager.tokenCache.mutex.Lock()
	manager.tokenCache.token = "expired-token"
	manager.tokenCache.projectID = "expired-project"
	manager.tokenCache.expiry = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	manager.tokenCache.mutex.Unlock()

	// GetToken should detect expired token and try to re-authenticate
	// This will fail because we don't have a real server, but it should not use the cached token
	ctx := context.Background()
	_, _, err = manager.GetToken(ctx)

	// Should get an error because authentication will fail, but importantly,
	// it should not return the expired cached token
	assert.Error(t, err)
}

// Unit tests for specific authentication methods and error handling
// Requirements: 2.1, 2.2, 2.3, 2.5

func TestPasswordAuthenticationFlow(t *testing.T) {
	tests := []struct {
		name           string
		config         *AuthConfig
		expectError    bool
		expectedMethod string
	}{
		{
			name: "Password auth with project ID",
			config: &AuthConfig{
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project-id",
			},
			expectError:    false,
			expectedMethod: "password",
		},
		{
			name: "Password auth with project name and domain ID",
			config: &AuthConfig{
				Username:    "test-user",
				Password:    "test-password",
				ProjectName: "test-project",
				DomainID:    "test-domain-id",
			},
			expectError:    false,
			expectedMethod: "password",
		},
		{
			name: "Password auth with project name and domain name",
			config: &AuthConfig{
				Username:    "test-user",
				Password:    "test-password",
				ProjectName: "test-project",
				DomainName:  "test-domain",
			},
			expectError:    false,
			expectedMethod: "password",
		},
		{
			name: "Password auth with default domain",
			config: &AuthConfig{
				Username:    "test-user",
				Password:    "test-password",
				ProjectName: "test-project",
			},
			expectError:    false,
			expectedMethod: "password",
		},
		{
			name: "Password auth missing project scope",
			config: &AuthConfig{
				Username: "test-user",
				Password: "test-password",
			},
			expectError: true,
		},
		{
			name: "Password auth missing username",
			config: &AuthConfig{
				Password:  "test-password",
				ProjectID: "test-project",
			},
			expectError: true,
		},
		{
			name: "Password auth missing password",
			config: &AuthConfig{
				Username:  "test-user",
				ProjectID: "test-project",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &AuthManager{config: tt.config}
			
			authReq, err := manager.buildAuthRequest()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, authReq)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authReq)
				assert.Contains(t, authReq.Auth.Identity.Methods, tt.expectedMethod)
				assert.NotNil(t, authReq.Auth.Identity.Password)
				assert.Equal(t, tt.config.Username, authReq.Auth.Identity.Password.User.Name)
				assert.Equal(t, tt.config.Password, authReq.Auth.Identity.Password.User.Password)
				
				// Verify scope is set
				assert.NotNil(t, authReq.Auth.Scope)
				assert.NotNil(t, authReq.Auth.Scope.Project)
				
				// Verify domain is set correctly
				if tt.config.DomainID != "" {
					assert.Equal(t, tt.config.DomainID, authReq.Auth.Identity.Password.User.Domain.ID)
				} else if tt.config.DomainName != "" {
					assert.Equal(t, tt.config.DomainName, authReq.Auth.Identity.Password.User.Domain.Name)
				} else {
					assert.Equal(t, "default", authReq.Auth.Identity.Password.User.Domain.Name)
				}
			}
		})
	}
}

func TestApplicationCredentialAuthenticationFlow(t *testing.T) {
	tests := []struct {
		name           string
		config         *AuthConfig
		expectError    bool
		expectedMethod string
	}{
		{
			name: "Valid application credential auth",
			config: &AuthConfig{
				ApplicationCredentialID:     "app-cred-id-123",
				ApplicationCredentialSecret: "app-cred-secret-456",
			},
			expectError:    false,
			expectedMethod: "application_credential",
		},
		{
			name: "Application credential auth missing ID",
			config: &AuthConfig{
				ApplicationCredentialSecret: "app-cred-secret-456",
			},
			expectError: true,
		},
		{
			name: "Application credential auth missing secret",
			config: &AuthConfig{
				ApplicationCredentialID: "app-cred-id-123",
			},
			expectError: true,
		},
		{
			name: "Application credential auth with empty ID",
			config: &AuthConfig{
				ApplicationCredentialID:     "",
				ApplicationCredentialSecret: "app-cred-secret-456",
			},
			expectError: true,
		},
		{
			name: "Application credential auth with empty secret",
			config: &AuthConfig{
				ApplicationCredentialID:     "app-cred-id-123",
				ApplicationCredentialSecret: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &AuthManager{config: tt.config}
			
			authReq, err := manager.buildAuthRequest()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, authReq)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authReq)
				assert.Contains(t, authReq.Auth.Identity.Methods, tt.expectedMethod)
				assert.NotNil(t, authReq.Auth.Identity.ApplicationCredential)
				assert.Equal(t, tt.config.ApplicationCredentialID, authReq.Auth.Identity.ApplicationCredential.ID)
				assert.Equal(t, tt.config.ApplicationCredentialSecret, authReq.Auth.Identity.ApplicationCredential.Secret)
				
				// Application credentials don't require explicit scope
				assert.Nil(t, authReq.Auth.Scope)
			}
		})
	}
}

func TestTokenAuthenticationFlow(t *testing.T) {
	tests := []struct {
		name           string
		config         *AuthConfig
		expectError    bool
		expectedMethod string
	}{
		{
			name: "Valid token auth with project ID",
			config: &AuthConfig{
				Token:     "existing-token-123",
				ProjectID: "test-project-id",
			},
			expectError:    false,
			expectedMethod: "token",
		},
		{
			name: "Valid token auth with project name",
			config: &AuthConfig{
				Token:       "existing-token-123",
				ProjectName: "test-project",
				DomainName:  "test-domain",
			},
			expectError:    false,
			expectedMethod: "token",
		},
		{
			name: "Token auth missing project scope",
			config: &AuthConfig{
				Token: "existing-token-123",
			},
			expectError: true,
		},
		{
			name: "Token auth with empty token",
			config: &AuthConfig{
				Token:     "",
				ProjectID: "test-project",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &AuthManager{config: tt.config}
			
			authReq, err := manager.buildAuthRequest()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, authReq)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authReq)
				assert.Contains(t, authReq.Auth.Identity.Methods, tt.expectedMethod)
				assert.NotNil(t, authReq.Auth.Identity.Token)
				assert.Equal(t, tt.config.Token, authReq.Auth.Identity.Token.ID)
				
				// Token auth requires explicit scope
				assert.NotNil(t, authReq.Auth.Scope)
				assert.NotNil(t, authReq.Auth.Scope.Project)
			}
		})
	}
}

func TestAuthenticationMethodPriority(t *testing.T) {
	// Test that application credentials take priority over other methods
	config := &AuthConfig{
		ApplicationCredentialID:     "app-cred-id",
		ApplicationCredentialSecret: "app-cred-secret",
		Token:                       "existing-token",
		Username:                    "test-user",
		Password:                    "test-password",
		ProjectID:                   "test-project",
	}

	manager := &AuthManager{config: config}
	authReq, err := manager.buildAuthRequest()

	assert.NoError(t, err)
	assert.Contains(t, authReq.Auth.Identity.Methods, "application_credential")
	assert.NotNil(t, authReq.Auth.Identity.ApplicationCredential)
	assert.Nil(t, authReq.Auth.Identity.Token)
	assert.Nil(t, authReq.Auth.Identity.Password)

	// Test that token takes priority over password when app creds not available
	config2 := &AuthConfig{
		Token:     "existing-token",
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}

	manager2 := &AuthManager{config: config2}
	authReq2, err := manager2.buildAuthRequest()

	assert.NoError(t, err)
	assert.Contains(t, authReq2.Auth.Identity.Methods, "token")
	assert.NotNil(t, authReq2.Auth.Identity.Token)
	assert.Nil(t, authReq2.Auth.Identity.Password)
}

func TestPasswordAuthenticationWithInvalidCredentials(t *testing.T) {
	// Create a mock server that returns 401 Unauthorized for invalid credentials
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request to check credentials
		var authReq AuthRequest
		json.NewDecoder(r.Body).Decode(&authReq)
		
		// Check if credentials are "invalid"
		if authReq.Auth.Identity.Password != nil && 
		   (authReq.Auth.Identity.Password.User.Name == "invalid-user" || 
		    authReq.Auth.Identity.Password.User.Password == "invalid-password") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"message": "The request you have made requires authentication.", "code": 401, "title": "Unauthorized"}}`))
			return
		}
		
		// Valid credentials - return success
		w.Header().Set("X-Subject-Token", "valid-token-123")
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
	}))
	defer server.Close()

	tests := []struct {
		name        string
		username    string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid username",
			username:    "invalid-user",
			password:    "valid-password",
			expectError: true,
			errorMsg:    "authentication failed with status 401",
		},
		{
			name:        "Invalid password",
			username:    "valid-user",
			password:    "invalid-password",
			expectError: true,
			errorMsg:    "authentication failed with status 401",
		},
		{
			name:        "Valid credentials",
			username:    "valid-user",
			password:    "valid-password",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				AuthURL:   server.URL,
				Username:  tt.username,
				Password:  tt.password,
				ProjectID: "test-project",
			}

			manager, err := NewAuthManager(config)
			require.NoError(t, err)

			ctx := context.Background()
			token, projectID, err := manager.authenticate(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Empty(t, projectID)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, projectID)
			}
		})
	}
}

func TestApplicationCredentialAuthenticationWithInvalidCredentials(t *testing.T) {
	// Create a mock server that returns 401 for invalid app credentials
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authReq AuthRequest
		json.NewDecoder(r.Body).Decode(&authReq)
		
		// Check if app credentials are "invalid"
		if authReq.Auth.Identity.ApplicationCredential != nil && 
		   (authReq.Auth.Identity.ApplicationCredential.ID == "invalid-id" || 
		    authReq.Auth.Identity.ApplicationCredential.Secret == "invalid-secret") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"message": "Invalid application credential", "code": 401, "title": "Unauthorized"}}`))
			return
		}
		
		// Valid credentials
		w.Header().Set("X-Subject-Token", "valid-token-123")
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
	}))
	defer server.Close()

	tests := []struct {
		name        string
		credID      string
		credSecret  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid credential ID",
			credID:      "invalid-id",
			credSecret:  "valid-secret",
			expectError: true,
			errorMsg:    "authentication failed with status 401",
		},
		{
			name:        "Invalid credential secret",
			credID:      "valid-id",
			credSecret:  "invalid-secret",
			expectError: true,
			errorMsg:    "authentication failed with status 401",
		},
		{
			name:        "Valid application credentials",
			credID:      "valid-id",
			credSecret:  "valid-secret",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				AuthURL:                     server.URL,
				ApplicationCredentialID:     tt.credID,
				ApplicationCredentialSecret: tt.credSecret,
			}

			manager, err := NewAuthManager(config)
			require.NoError(t, err)

			ctx := context.Background()
			token, projectID, err := manager.authenticate(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Empty(t, projectID)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, projectID)
			}
		})
	}
}

func TestTokenAuthenticationWithInvalidCredentials(t *testing.T) {
	// Create a mock server that returns 401 for invalid tokens
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authReq AuthRequest
		json.NewDecoder(r.Body).Decode(&authReq)
		
		// Check if token is "invalid"
		if authReq.Auth.Identity.Token != nil && 
		   strings.Contains(authReq.Auth.Identity.Token.ID, "invalid") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"message": "Invalid token", "code": 401, "title": "Unauthorized"}}`))
			return
		}
		
		// Valid token
		w.Header().Set("X-Subject-Token", "new-valid-token-123")
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
	}))
	defer server.Close()

	tests := []struct {
		name        string
		token       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid token",
			token:       "invalid-token-123",
			expectError: true,
			errorMsg:    "authentication failed with status 401",
		},
		{
			name:        "Expired token",
			token:       "invalid-expired-token",
			expectError: true,
			errorMsg:    "authentication failed with status 401",
		},
		{
			name:        "Valid token",
			token:       "valid-token-123",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				AuthURL:   server.URL,
				Token:     tt.token,
				ProjectID: "test-project",
			}

			manager, err := NewAuthManager(config)
			require.NoError(t, err)

			ctx := context.Background()
			token, projectID, err := manager.authenticate(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Empty(t, projectID)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, projectID)
			}
		})
	}
}

func TestAuthenticationNetworkErrors(t *testing.T) {
	tests := []struct {
		name        string
		serverFunc  func() *httptest.Server
		expectError bool
		errorMsg    string
	}{
		{
			name: "Connection refused",
			serverFunc: func() *httptest.Server {
				// Return a server that's immediately closed to simulate connection refused
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
				server.Close()
				return server
			},
			expectError: true,
			errorMsg:    "authentication request failed",
		},
		{
			name: "Server timeout",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Simulate timeout by sleeping longer than client timeout
					time.Sleep(35 * time.Second) // Client timeout is 30s
				}))
			},
			expectError: true,
			errorMsg:    "authentication request failed",
		},
		{
			name: "Invalid JSON response",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Subject-Token", "test-token")
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`invalid json response`))
				}))
			},
			expectError: true,
			errorMsg:    "failed to parse authentication response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverFunc()
			if tt.name != "Connection refused" {
				defer server.Close()
			}

			config := &AuthConfig{
				AuthURL:   server.URL,
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project",
			}

			manager, err := NewAuthManager(config)
			require.NoError(t, err)

			ctx := context.Background()
			
			// For timeout test, use a shorter context timeout
			if tt.name == "Server timeout" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
				defer cancel()
			}

			token, projectID, err := manager.authenticate(ctx)

			assert.Error(t, err)
			assert.Empty(t, token)
			assert.Empty(t, projectID)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

// TestAuthenticationTokenCachingProperty implements Property 3: Authentication Token Caching
// **Validates: Requirements 2.6, 8.3**
func TestAuthenticationTokenCachingProperty(t *testing.T) {
	// Property-based test function
	f := func(tokenLifetimeMinutes uint8, cacheBufferMinutes uint8) bool {
		// Constrain inputs to reasonable ranges
		if tokenLifetimeMinutes == 0 || tokenLifetimeMinutes > 120 {
			return true // Skip invalid inputs
		}
		if cacheBufferMinutes > tokenLifetimeMinutes {
			return true // Skip invalid inputs where buffer is larger than lifetime
		}

		// Create a mock Keystone server that tracks authentication calls
		authCallCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/auth/tokens") {
				authCallCount++
				
				// Generate unique token for each call
				token := fmt.Sprintf("test-token-%d", authCallCount)
				w.Header().Set("X-Subject-Token", token)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				// Set expiration based on test parameters
				expiry := time.Now().Add(time.Duration(tokenLifetimeMinutes) * time.Minute)
				response := AuthResponse{
					Token: struct {
						ExpiresAt string `json:"expires_at"`
						Project   struct {
							ID string `json:"id"`
						} `json:"project"`
					}{
						ExpiresAt: expiry.Format(time.RFC3339),
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
		defer server.Close()

		config := &AuthConfig{
			AuthURL:   server.URL,
			Username:  "test-user",
			Password:  "test-password",
			ProjectID: "test-project",
		}

		manager, err := NewAuthManager(config)
		if err != nil {
			t.Logf("Failed to create auth manager: %v", err)
			return false
		}

		ctx := context.Background()

		// First call should authenticate and cache token
		token1, projectID1, err := manager.GetToken(ctx)
		if err != nil {
			t.Logf("First GetToken call failed: %v", err)
			return false
		}

		if authCallCount != 1 {
			t.Logf("Expected 1 auth call after first GetToken, got %d", authCallCount)
			return false
		}

		// Verify token is cached
		manager.tokenCache.mutex.RLock()
		cachedToken := manager.tokenCache.token
		cachedProjectID := manager.tokenCache.projectID
		cachedExpiry := manager.tokenCache.expiry
		manager.tokenCache.mutex.RUnlock()

		if cachedToken != token1 {
			t.Logf("Token not cached correctly: expected %s, got %s", token1, cachedToken)
			return false
		}

		if cachedProjectID != projectID1 {
			t.Logf("ProjectID not cached correctly: expected %s, got %s", projectID1, cachedProjectID)
			return false
		}

		// Second call should use cached token (no new auth call)
		token2, projectID2, err := manager.GetToken(ctx)
		if err != nil {
			t.Logf("Second GetToken call failed: %v", err)
			return false
		}

		if authCallCount != 1 {
			t.Logf("Expected 1 auth call after second GetToken (should use cache), got %d", authCallCount)
			return false
		}

		if token1 != token2 {
			t.Logf("Second call should return same cached token: expected %s, got %s", token1, token2)
			return false
		}

		if projectID1 != projectID2 {
			t.Logf("Second call should return same cached projectID: expected %s, got %s", projectID1, projectID2)
			return false
		}

		// Verify that the cache expiry is set correctly (with buffer)
		expectedExpiry := time.Now().Add(time.Duration(tokenLifetimeMinutes)*time.Minute - 5*time.Minute)
		timeDiff := cachedExpiry.Sub(expectedExpiry).Abs()
		if timeDiff > 10*time.Second { // Allow 10 second tolerance for test execution time
			t.Logf("Cache expiry not set correctly: expected around %v, got %v (diff: %v)", expectedExpiry, cachedExpiry, timeDiff)
			return false
		}

		// Simulate token expiration by manually setting expiry to past
		manager.tokenCache.mutex.Lock()
		manager.tokenCache.expiry = time.Now().Add(-1 * time.Minute)
		manager.tokenCache.mutex.Unlock()

		// Third call should detect expired token and re-authenticate
		token3, projectID3, err := manager.GetToken(ctx)
		if err != nil {
			t.Logf("Third GetToken call (after expiry) failed: %v", err)
			return false
		}

		if authCallCount != 2 {
			t.Logf("Expected 2 auth calls after token expiry, got %d", authCallCount)
			return false
		}

		if token1 == token3 {
			t.Logf("Third call should return new token after expiry: got same token %s", token3)
			return false
		}

		// Verify new token is cached
		manager.tokenCache.mutex.RLock()
		newCachedToken := manager.tokenCache.token
		newCachedProjectID := manager.tokenCache.projectID
		manager.tokenCache.mutex.RUnlock()

		if newCachedToken != token3 {
			t.Logf("New token not cached correctly: expected %s, got %s", token3, newCachedToken)
			return false
		}

		if newCachedProjectID != projectID3 {
			t.Logf("New projectID not cached correctly: expected %s, got %s", projectID3, newCachedProjectID)
			return false
		}

		return true
	}

	// Run the property-based test with constrained iterations for reasonable execution time
	if err := quick.Check(f, &quick.Config{MaxCount: 10}); err != nil {
		t.Error(err)
	}
}