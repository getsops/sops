package publish

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAWSSecretsManagerDestination(t *testing.T) {
	dest := NewAWSSecretsManagerDestination("us-east-1", "myapp/database")
	assert.NotNil(t, dest)
	assert.Equal(t, "us-east-1", dest.region)
	assert.Equal(t, "myapp/database", dest.secretName)
}

func TestAWSSecretsManagerDestination_Path(t *testing.T) {
	// Test with specified secret name
	dest := NewAWSSecretsManagerDestination("us-east-1", "myapp/database")
	path := dest.Path("config.yaml")
	expected := "arn:aws:secretsmanager:us-east-1:*:secret:myapp/database"
	assert.Equal(t, expected, path)

	// Test without specified secret name (uses filename)
	dest = NewAWSSecretsManagerDestination("us-west-2", "")
	path = dest.Path("api-keys.yaml")
	expected = "arn:aws:secretsmanager:us-west-2:*:secret:api-keys.yaml"
	assert.Equal(t, expected, path)
}

func TestAWSSecretsManagerDestination_Upload(t *testing.T) {
	dest := NewAWSSecretsManagerDestination("us-east-1", "test-secret")
	err := dest.Upload([]byte("test content"), "test.yaml")

	// Should return NotImplementedError
	assert.NotNil(t, err)
	assert.IsType(t, &NotImplementedError{}, err)
	assert.Contains(t, err.Error(), "AWS Secrets Manager does not support uploading encrypted sops files directly")
}

func TestNewAWSSecretsManagerDestination_EmptyRegion(t *testing.T) {
	// Test that empty region is allowed (will use SDK defaults)
	dest := NewAWSSecretsManagerDestination("", "myapp/database")
	assert.NotNil(t, dest)
	assert.Equal(t, "", dest.region)
	assert.Equal(t, "myapp/database", dest.secretName)
}
