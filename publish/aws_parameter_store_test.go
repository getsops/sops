package publish

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAWSParameterStoreDestination(t *testing.T) {
	dest := NewAWSParameterStoreDestination("us-east-1", "/myapp/config")
	assert.NotNil(t, dest)
	assert.Equal(t, "us-east-1", dest.region)
	assert.Equal(t, "/myapp/config", dest.parameterPath)

	// Test path normalization (should add leading slash)
	dest = NewAWSParameterStoreDestination("us-east-1", "myapp/config")
	assert.Equal(t, "/myapp/config", dest.parameterPath)
}

func TestAWSParameterStoreDestination_Path(t *testing.T) {
	// Test with specific parameter path (no trailing slash)
	dest := NewAWSParameterStoreDestination("us-east-1", "/myapp/database")
	path := dest.Path("config.yaml")
	assert.Equal(t, "/myapp/database", path)

	// Test with parameter path ending with slash
	dest = NewAWSParameterStoreDestination("us-east-1", "/myapp/configs/")
	path = dest.Path("api.yaml")
	assert.Equal(t, "/myapp/configs/api.yaml", path)

	// Test with empty parameter path (uses filename)
	dest = NewAWSParameterStoreDestination("us-east-1", "")
	path = dest.Path("standalone.yaml")
	assert.Equal(t, "/standalone.yaml", path)

	// Test with filename that already has leading slash
	dest = NewAWSParameterStoreDestination("us-east-1", "")
	path = dest.Path("/already-prefixed.yaml")
	assert.Equal(t, "/already-prefixed.yaml", path)
}

func TestAWSParameterStoreDestination_Upload(t *testing.T) {
	dest := NewAWSParameterStoreDestination("us-east-1", "/test-parameter")
	err := dest.Upload([]byte("test content"), "test.yaml")

	assert.NotNil(t, err)
	assert.IsType(t, &NotImplementedError{}, err)
	assert.Contains(t, err.Error(), "AWS Parameter Store does not support uploading encrypted sops files directly")
}

func TestNewAWSParameterStoreDestination_EmptyRegion(t *testing.T) {
	// Test that empty region is allowed (will use SDK defaults)
	dest := NewAWSParameterStoreDestination("", "/myapp/config")
	assert.NotNil(t, dest)
	assert.Equal(t, "", dest.region)
	assert.Equal(t, "/myapp/config", dest.parameterPath)
}
