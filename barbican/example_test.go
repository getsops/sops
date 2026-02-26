package barbican_test

import (
	"fmt"
	"log"
	"os"

	"github.com/getsops/sops/v3/barbican"
)

// This example demonstrates basic usage of Barbican MasterKey.
// Note: This example requires a real OpenStack environment and will not run in tests.
func ExampleMasterKey_usage() {
	// Set up OpenStack authentication environment variables
	os.Setenv("OS_AUTH_URL", "https://keystone.example.com:5000/v3")
	os.Setenv("OS_USERNAME", "sops-user")
	os.Setenv("OS_PASSWORD", "secret")
	os.Setenv("OS_PROJECT_ID", "abc123")
	os.Setenv("OS_DOMAIN_NAME", "default")

	// Create a Barbican MasterKey with a secret reference
	masterKey := &barbican.MasterKey{
		SecretRef: "550e8400-e29b-41d4-a716-446655440000",
		Region:    "sjc3",
	}

	// Example data key to encrypt
	dataKey := []byte("my-secret-data-key")

	// Encrypt the data key
	err := masterKey.Encrypt(dataKey)
	if err != nil {
		log.Fatalf("Failed to encrypt: %v", err)
	}

	fmt.Printf("Encrypted key stored in Barbican secret: %s\n", masterKey.EncryptedKey)

	// Decrypt the data key
	decryptedKey, err := masterKey.Decrypt()
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}

	fmt.Printf("Decrypted data key length: %d bytes\n", len(decryptedKey))
}

// This example demonstrates parsing secret references from a string.
func ExampleMasterKeysFromSecretRefString() {
	// Parse multiple secret references
	secretRefs := "550e8400-e29b-41d4-a716-446655440000,region:dfw3:660e8400-e29b-41d4-a716-446655440001"
	
	masterKeys, err := barbican.MasterKeysFromSecretRefString(secretRefs)
	if err != nil {
		log.Fatalf("Failed to parse secret references: %v", err)
	}

	fmt.Printf("Parsed %d master keys\n", len(masterKeys))
	for i, key := range masterKeys {
		fmt.Printf("Key %d: SecretRef=%s, Region=%s\n", i+1, key.SecretRef, key.Region)
	}
	// Output: Parsed 2 master keys
	// Key 1: SecretRef=550e8400-e29b-41d4-a716-446655440000, Region=
	// Key 2: SecretRef=region:dfw3:660e8400-e29b-41d4-a716-446655440001, Region=dfw3
}

// This example demonstrates different authentication methods for OpenStack.
func ExampleAuthConfig_methods() {
	// Password authentication
	passwordAuth := &barbican.AuthConfig{
		AuthURL:     "https://keystone.example.com:5000/v3",
		Username:    "sops-user",
		Password:    "secret",
		ProjectID:   "abc123",
		DomainName:  "default",
	}

	// Application credential authentication (recommended)
	appCredAuth := &barbican.AuthConfig{
		AuthURL:                     "https://keystone.example.com:5000/v3",
		ApplicationCredentialID:     "app-cred-id",
		ApplicationCredentialSecret: "app-cred-secret",
	}

	// Token authentication
	tokenAuth := &barbican.AuthConfig{
		AuthURL: "https://keystone.example.com:5000/v3",
		Token:   "existing-token",
	}

	fmt.Printf("Password auth configured for user: %s\n", passwordAuth.Username)
	fmt.Printf("App credential auth configured with ID: %s\n", appCredAuth.ApplicationCredentialID)
	fmt.Printf("Token auth configured with token length: %d\n", len(tokenAuth.Token))
	// Output: Password auth configured for user: sops-user
	// App credential auth configured with ID: app-cred-id
	// Token auth configured with token length: 14
}