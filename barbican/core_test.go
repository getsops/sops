package barbican

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCoreEncryptionDecryption tests the basic encryption/decryption functionality
func TestCoreEncryptionDecryption(t *testing.T) {
	// Create mock servers
	barbicanServer := NewMockBarbicanServer()
	defer barbicanServer.Close()
	
	keystoneServer := NewMockKeystoneServer()
	defer keystoneServer.Close()
	
	// Test data
	testData := []byte("test-secret-data-12345")
	
	// Create master key with mock configuration
	key := NewMasterKey("550e8400-e29b-41d4-a716-446655440000")
	
	// Configure authentication
	config := &AuthConfig{
		AuthURL:   keystoneServer.URL(),
		Username:  "test-user",
		Password:  "test-password",
		ProjectID: "test-project",
	}
	
	authManager, err := NewAuthManager(config)
	require.NoError(t, err)
	
	key.AuthConfig = config
	key.authManager = authManager
	key.baseEndpoint = barbicanServer.URL()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Test encryption
	err = key.EncryptContext(ctx, testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, key.EncryptedKey)
	
	// Verify secret was stored in Barbican
	assert.Equal(t, 1, barbicanServer.GetSecretCount())
	
	// Test decryption
	decryptedData, err := key.DecryptContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, testData, decryptedData)
}

// TestBasicAuthenticationMethods tests authentication methods individually
func TestBasicAuthenticationMethods(t *testing.T) {
	barbicanServer := NewMockBarbicanServer()
	defer barbicanServer.Close()
	
	keystoneServer := NewMockKeystoneServer()
	defer keystoneServer.Close()
	
	testData := []byte("auth-test-data")
	
	tests := []struct {
		name   string
		config *AuthConfig
	}{
		{
			name: "Password Authentication",
			config: &AuthConfig{
				AuthURL:   keystoneServer.URL(),
				Username:  "test-user",
				Password:  "test-password",
				ProjectID: "test-project",
			},
		},
		{
			name: "Application Credential Authentication",
			config: &AuthConfig{
				AuthURL:                     keystoneServer.URL(),
				ApplicationCredentialID:     "app-cred-123",
				ApplicationCredentialSecret: "app-secret-456",
			},
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create master key
			secretRef := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", i)
			key := NewMasterKey(secretRef)
			
			// Configure authentication
			authManager, err := NewAuthManager(tt.config)
			require.NoError(t, err)
			
			key.AuthConfig = tt.config
			key.authManager = authManager
			key.baseEndpoint = barbicanServer.URL()
			
			// Test encryption
			err = key.EncryptContext(ctx, testData)
			assert.NoError(t, err)
			assert.NotEmpty(t, key.EncryptedKey)
			
			// Test decryption
			decryptedData, err := key.DecryptContext(ctx)
			assert.NoError(t, err)
			assert.Equal(t, testData, decryptedData)
		})
	}
}

// TestErrorHandling tests basic error scenarios
func TestErrorHandling(t *testing.T) {
	keystoneServer := NewMockKeystoneServer()
	defer keystoneServer.Close()
	
	testData := []byte("error-test-data")
	
	// Test with invalid credentials
	key := NewMasterKey("550e8400-e29b-41d4-a716-446655440000")
	
	config := &AuthConfig{
		AuthURL:   keystoneServer.URL(),
		Username:  "invalid-user",
		Password:  "invalid-password",
		ProjectID: "test-project",
	}
	
	authManager, err := NewAuthManager(config)
	require.NoError(t, err)
	
	key.AuthConfig = config
	key.authManager = authManager
	key.baseEndpoint = "http://localhost:99999" // Non-existent service
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test encryption - should fail
	err = key.EncryptContext(ctx, testData)
	assert.Error(t, err)
	assert.Empty(t, key.EncryptedKey)
}