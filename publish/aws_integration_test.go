//go:build integration

package publish

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests for AWS Secrets Manager and Parameter Store publishing.
// These tests require real AWS credentials and resources.
//
// To run these tests:
// 1. Set up AWS credentials (AWS_PROFILE, AWS_ACCESS_KEY_ID, etc.)
// 2. Set required environment variables:
//    - SOPS_TEST_AWS_REGION (default: us-east-1)
//    - SOPS_TEST_AWS_SECRET_NAME (secret for testing)
//    - SOPS_TEST_AWS_PARAMETER_NAME (parameter for testing)
// 3. Run with: go test -tags=integration ./publish -run TestAWS -v
//
// Prerequisites:
// - AWS credentials with Secrets Manager and Parameter Store permissions
// - Test secret and parameter resources should already exist or be creatable

var (
	testAWSRegion     = getEnvOrDefault("SOPS_TEST_AWS_REGION", "us-east-1")
	testSecretName    = os.Getenv("SOPS_TEST_AWS_SECRET_NAME")    // e.g., "sops-test-secret"
	testParameterName = os.Getenv("SOPS_TEST_AWS_PARAMETER_NAME") // e.g., "/sops-test/parameter"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestAWSSecretsManagerDestination_PlainText_Integration(t *testing.T) {
	if testSecretName == "" {
		t.Skip("Skipping integration test: SOPS_TEST_AWS_SECRET_NAME not set")
	}

	ctx := context.Background()
	dest := NewAWSSecretsManagerDestination(testAWSRegion, testSecretName)

	// Test data with complex nested structure
	// Note: This format stores as Plain Text JSON and does NOT enable key/value editor in AWS console
	// For key/value format, see TestAWSSecretsManagerDestination_KeyValue_Integration
	testData := map[string]interface{}{
		"database": map[string]interface{}{
			"host":     "localhost",
			"port":     float64(5432),
			"username": "testuser",
			"password": "supersecret",
		},
		"api_keys": map[string]interface{}{
			"stripe": "sk_test_123456",
			"github": "ghp_987654321",
		},
	}

	// Upload test data
	err := dest.UploadUnencrypted(testData, "test-secret")
	require.NoError(t, err, "Failed to upload secret to Secrets Manager")

	// Verify the secret was stored correctly
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(testAWSRegion))
	require.NoError(t, err, "Failed to load AWS config")

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(testSecretName),
	})
	require.NoError(t, err, "Failed to retrieve secret from Secrets Manager")

	// Parse and verify the stored data
	var storedData map[string]interface{}
	err = json.Unmarshal([]byte(*result.SecretString), &storedData)
	require.NoError(t, err, "Failed to parse stored secret JSON")

	assert.Equal(t, testData, storedData, "Stored data doesn't match original")

	// Test no-op behavior (upload same data again)
	err = dest.UploadUnencrypted(testData, "test-secret")
	assert.NoError(t, err, "No-op upload should succeed")
}

func TestAWSSecretsManagerDestination_KeyValue_Integration(t *testing.T) {
	if testSecretName == "" {
		t.Skip("Skipping integration test: SOPS_TEST_AWS_SECRET_NAME not set")
	}

	ctx := context.Background()
	// Use a different secret name for key/value test to avoid conflicts
	keyValueSecretName := testSecretName + "-keyvalue"
	dest := NewAWSSecretsManagerDestination(testAWSRegion, keyValueSecretName)

	// Test data with simple key/value pairs (no nested objects)
	// This format enables the key/value editor in AWS Secrets Manager console
	testData := map[string]interface{}{
		"database_host":     "db.example.com",
		"database_port":     "5432",
		"database_username": "app_user",
		"database_password": "super_secret_password",
		"api_key_stripe":    "sk_live_abcdef123456",
		"api_key_github":    "ghp_xyz789012345",
		"debug_mode":        "false",
		"log_level":         "info",
		"max_connections":   "100",
		"timeout_seconds":   "30",
	}

	// Upload test data
	err := dest.UploadUnencrypted(testData, "test-keyvalue-secret")
	require.NoError(t, err, "Failed to upload key/value secret to Secrets Manager")

	// Verify the secret was stored correctly
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(testAWSRegion))
	require.NoError(t, err, "Failed to load AWS config")

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(keyValueSecretName),
	})
	require.NoError(t, err, "Failed to retrieve key/value secret from Secrets Manager")

	// Parse and verify the stored data
	var storedData map[string]interface{}
	err = json.Unmarshal([]byte(*result.SecretString), &storedData)
	require.NoError(t, err, "Failed to parse stored key/value secret JSON")

	assert.Equal(t, testData, storedData, "Stored key/value data doesn't match original")

	// Verify that all values are stored as strings (important for key/value format)
	for key, value := range storedData {
		assert.IsType(t, "", value, "Value for key %s should be a string for key/value format", key)
	}

	// Test no-op behavior (upload same data again)
	err = dest.UploadUnencrypted(testData, "test-keyvalue-secret")
	assert.NoError(t, err, "No-op upload should succeed for key/value format")
}

func TestAWSParameterStoreDestination_Integration(t *testing.T) {
	if testParameterName == "" {
		t.Skip("Skipping integration test: SOPS_TEST_AWS_PARAMETER_NAME not set")
	}

	ctx := context.Background()
	dest := NewAWSParameterStoreDestination(testAWSRegion, testParameterName)

	// Test data
	testData := map[string]interface{}{
		"app_config": map[string]interface{}{
			"debug":       false,
			"log_level":   "info",
			"max_workers": float64(10), 
			"features": map[string]interface{}{
				"new_ui":       true,
				"beta_feature": false,
			},
		},
	}

	// Upload test data (this is the method used by the publish command)
	err := dest.UploadUnencrypted(testData, "test-config")
	require.NoError(t, err, "Failed to upload parameter to Parameter Store")

	// Verify the parameter was stored correctly
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(testAWSRegion))
	require.NoError(t, err, "Failed to load AWS config")

	client := ssm.NewFromConfig(cfg)
	result, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(testParameterName),
		WithDecryption: aws.Bool(true),
	})
	require.NoError(t, err, "Failed to retrieve parameter from Parameter Store")

	// Parse and verify the stored data
	var storedData map[string]interface{}
	err = json.Unmarshal([]byte(*result.Parameter.Value), &storedData)
	require.NoError(t, err, "Failed to parse stored parameter JSON")

	assert.Equal(t, testData, storedData, "Stored data doesn't match original")

	// Verify parameter type is always SecureString
	assert.Equal(t, "SecureString", string(result.Parameter.Type), "Parameter type should always be SecureString")

	// Test no-op behavior (upload same data again)
	err = dest.UploadUnencrypted(testData, "test-config")
	assert.NoError(t, err, "No-op upload should succeed")
}

func TestAWSParameterStoreDestination_EncryptedFile_Integration(t *testing.T) {
	if testParameterName == "" {
		t.Skip("Skipping integration test: SOPS_TEST_AWS_PARAMETER_NAME not set")
	}

	dest := NewAWSParameterStoreDestination(testAWSRegion, testParameterName+"-file")

	encryptedContent := []byte(`# SOPS encrypted file
database:
    host: ENC[AES256_GCM,data:xyz123,type:str]
    password: ENC[AES256_GCM,data:abc456,type:str]
sops:
    kms:
    - arn: arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012
    version: 3.8.1`)

	// Upload should return NotImplementedError
	err := dest.Upload(encryptedContent, "encrypted-test")
	require.NotNil(t, err, "Upload should return an error")
	assert.IsType(t, &NotImplementedError{}, err, "Should return NotImplementedError")
	assert.Contains(t, err.Error(), "AWS Parameter Store does not support uploading encrypted sops files directly")
}
