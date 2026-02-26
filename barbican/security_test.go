package barbican

import (
	"crypto/tls"
	"testing"
)

func TestSecurityValidator_ValidateAndCreateTLSConfig(t *testing.T) {
	tests := []struct {
		name           string
		config         *SecurityConfig
		expectError    bool
		expectInsecure bool
		expectMinTLS   uint16
	}{
		{
			name:           "default secure config",
			config:         DefaultSecurityConfig(),
			expectError:    false,
			expectInsecure: false,
			expectMinTLS:   tls.VersionTLS12,
		},
		{
			name: "insecure TLS config",
			config: &SecurityConfig{
				InsecureTLS:          true,
				MinTLSVersion:        tls.VersionTLS12,
				ShowSecurityWarnings: false, // Disable warnings for test
			},
			expectError:    false,
			expectInsecure: true,
			expectMinTLS:   tls.VersionTLS12,
		},
		{
			name: "custom CA cert content",
			config: &SecurityConfig{
				CACertContent: `-----BEGIN CERTIFICATE-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7d7Qj8...
-----END CERTIFICATE-----`,
				MinTLSVersion:        tls.VersionTLS12,
				ShowSecurityWarnings: false,
			},
			expectError:    true, // Invalid cert content will cause error
			expectInsecure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewSecurityValidator(tt.config)
			
			tlsConfig, err := validator.ValidateAndCreateTLSConfig()
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if tlsConfig.InsecureSkipVerify != tt.expectInsecure {
				t.Errorf("Expected InsecureSkipVerify=%v, got %v", tt.expectInsecure, tlsConfig.InsecureSkipVerify)
			}
			
			if tlsConfig.MinVersion != tt.expectMinTLS {
				t.Errorf("Expected MinVersion=%v, got %v", tt.expectMinTLS, tlsConfig.MinVersion)
			}
		})
	}
}

func TestSecurityValidator_ValidateAuthConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *AuthConfig
		expectError bool
		errorType   BarbicanErrorType
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorType:   ErrorTypeConfig,
		},
		{
			name: "missing auth URL",
			config: &AuthConfig{
				Username: "user",
				Password: "pass",
			},
			expectError: true,
			errorType:   ErrorTypeConfig,
		},
		{
			name: "valid password auth",
			config: &AuthConfig{
				AuthURL:     "https://keystone.example.com:5000/v3",
				Username:    "user",
				Password:    "pass",
				ProjectID:   "project123",
			},
			expectError: false,
		},
		{
			name: "valid app credential auth",
			config: &AuthConfig{
				AuthURL:                     "https://keystone.example.com:5000/v3",
				ApplicationCredentialID:     "app-cred-id",
				ApplicationCredentialSecret: "app-cred-secret",
			},
			expectError: false,
		},
		{
			name: "no auth method",
			config: &AuthConfig{
				AuthURL: "https://keystone.example.com:5000/v3",
			},
			expectError: true,
			errorType:   ErrorTypeConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewSecurityValidator(DefaultSecurityConfig())
			
			err := validator.ValidateAuthConfig(tt.config)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				
				if barbicanErr, ok := err.(*BarbicanError); ok {
					if barbicanErr.Type != tt.errorType {
						t.Errorf("Expected error type %v, got %v", tt.errorType, barbicanErr.Type)
					}
				} else {
					t.Errorf("Expected BarbicanError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecurityValidator_SanitizeForLogging(t *testing.T) {
	validator := NewSecurityValidator(DefaultSecurityConfig())
	
	input := map[string]interface{}{
		"username":     "testuser",
		"password":     "secret123",
		"token":        "auth-token-value",
		"url":          "https://user:pass@example.com/path",
		"normal_field": "normal_value",
	}
	
	sanitized := validator.SanitizeForLogging(input)
	
	// Check that sensitive fields are redacted
	if sanitized["password"] == "secret123" {
		t.Errorf("Password was not sanitized")
	}
	
	if sanitized["token"] == "auth-token-value" {
		t.Errorf("Token was not sanitized")
	}
	
	// Check that URL credentials are sanitized
	if url, ok := sanitized["url"].(string); ok {
		if url == "https://user:pass@example.com/path" {
			t.Errorf("URL credentials were not sanitized")
		}
	}
	
	// Check that normal fields are preserved
	if sanitized["normal_field"] != "normal_value" {
		t.Errorf("Normal field was incorrectly modified")
	}
}

func TestSecurityValidator_ValidateSecretRef(t *testing.T) {
	validator := NewSecurityValidator(DefaultSecurityConfig())
	
	tests := []struct {
		name        string
		secretRef   string
		expectError bool
	}{
		{
			name:        "empty secret ref",
			secretRef:   "",
			expectError: true,
		},
		{
			name:        "valid UUID",
			secretRef:   "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:        "valid regional format",
			secretRef:   "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:        "secret ref with spaces",
			secretRef:   "550e8400 e29b-41d4-a716-446655440000",
			expectError: true,
		},
		{
			name:        "secret ref with newlines",
			secretRef:   "550e8400-e29b-41d4-a716-446655440000\n",
			expectError: true,
		},
		{
			name:        "invalid format",
			secretRef:   "not-a-valid-uuid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSecretRef(tt.secretRef)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecureCleanup(t *testing.T) {
	validator := NewSecurityValidator(DefaultSecurityConfig())
	cleanup := NewSecureCleanup(validator)
	
	// Test adding temporary secrets
	cleanup.AddTemporarySecret("550e8400-e29b-41d4-a716-446655440000")
	cleanup.AddTemporarySecret("660e8400-e29b-41d4-a716-446655440001")
	
	if cleanup.GetTemporarySecretCount() != 2 {
		t.Errorf("Expected 2 temporary secrets, got %d", cleanup.GetTemporarySecretCount())
	}
	
	// Note: We can't test actual cleanup without a mock client
	// This would require implementing a mock ClientInterface
}

func TestBarbicanError_Error(t *testing.T) {
	err := NewAuthenticationError("Invalid credentials").
		WithDetails("Username or password is incorrect").
		WithSecretRef("550e8400-e29b-41d4-a716-446655440000").
		WithRegion("sjc3").
		WithSuggestions("Check your credentials", "Verify the auth URL")
	
	errorStr := err.Error()
	
	// Check that error contains expected components
	if !contains(errorStr, "authentication error") {
		t.Errorf("Error string should contain error type")
	}
	
	if !contains(errorStr, "Invalid credentials") {
		t.Errorf("Error string should contain message")
	}
	
	if !contains(errorStr, "Username or password is incorrect") {
		t.Errorf("Error string should contain details")
	}
	
	if !contains(errorStr, "sjc3") {
		t.Errorf("Error string should contain region")
	}
	
	// Check that secret reference is sanitized (should not contain full UUID)
	if contains(errorStr, "550e8400-e29b-41d4-a716-446655440000") {
		t.Errorf("Error string should not contain full secret reference")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}