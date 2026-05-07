package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBarbicanConfigurationIntegration tests the complete Barbican configuration workflow
func TestBarbicanConfigurationIntegration(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "sops-barbican-config-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a .sops.yaml configuration file
	configContent := `creation_rules:
  - path_regex: \.prod\.yaml$
    barbican:
      - "550e8400-e29b-41d4-a716-446655440000"
      - "region:us-west-1:660e8400-e29b-41d4-a716-446655440001"
    barbican_auth_url: "https://keystone.example.com:5000/v3"
    barbican_region: "us-east-1"
  - path_regex: \.dev\.yaml$
    barbican: "770e8400-e29b-41d4-a716-446655440002"
    barbican_auth_url: "https://keystone-dev.example.com:5000/v3"
    barbican_region: "us-west-2"
  - path_regex: ""
    key_groups:
    - barbican:
      - secret_ref: "880e8400-e29b-41d4-a716-446655440003"
        region: "eu-central-1"
      kms:
      - arn: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
      pgp:
      - "85D77543B3D624B63CEA9E6DBC17301B491B3F21"
`

	configPath := filepath.Join(tempDir, ".sops.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.Nil(t, err)

	// Test 1: Production file should match first rule
	prodFilePath := filepath.Join(tempDir, "secrets.prod.yaml")
	conf, err := LoadCreationRuleForFile(configPath, prodFilePath, nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 2, len(conf.KeyGroups[0]))

	// Verify both keys are Barbican keys
	for _, key := range conf.KeyGroups[0] {
		assert.Equal(t, "barbican", key.TypeToIdentifier())
	}

	// Test 2: Development file should match second rule
	devFilePath := filepath.Join(tempDir, "config.dev.yaml")
	conf, err = LoadCreationRuleForFile(configPath, devFilePath, nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 1, len(conf.KeyGroups[0]))
	assert.Equal(t, "barbican", conf.KeyGroups[0][0].TypeToIdentifier())

	// Test 3: Other files should match default rule with mixed key types
	otherFilePath := filepath.Join(tempDir, "other.yaml")
	conf, err = LoadCreationRuleForFile(configPath, otherFilePath, nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 3, len(conf.KeyGroups[0])) // 1 Barbican + 1 KMS + 1 PGP

	// Count key types to verify mixed configuration
	keyTypeCounts := make(map[string]int)
	for _, key := range conf.KeyGroups[0] {
		keyTypeCounts[key.TypeToIdentifier()]++
	}

	assert.Equal(t, 1, keyTypeCounts["barbican"])
	assert.Equal(t, 1, keyTypeCounts["kms"])
	assert.Equal(t, 1, keyTypeCounts["pgp"])

	// Test 4: Verify configuration validation works
	invalidConfigContent := `creation_rules:
  - path_regex: ""
    barbican: "invalid-secret-ref"
`
	invalidConfigPath := filepath.Join(tempDir, ".sops-invalid.yaml")
	err = os.WriteFile(invalidConfigPath, []byte(invalidConfigContent), 0644)
	assert.Nil(t, err)

	_, err = LoadCreationRuleForFile(invalidConfigPath, otherFilePath, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid Barbican secret reference")
}

// TestBarbicanConfigurationEdgeCases tests edge cases and error conditions
func TestBarbicanConfigurationEdgeCases(t *testing.T) {
	// Test empty Barbican key list
	emptyConfig := []byte(`
creation_rules:
  - path_regex: ""
    barbican: []
`)
	conf, err := parseCreationRuleForFile(parseConfigFile(emptyConfig, t), "/conf/path", "test", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 0, len(conf.KeyGroups[0])) // No keys should be created

	// Test mixed string and array format
	mixedConfig := []byte(`
creation_rules:
  - path_regex: "mixed*"
    barbican: "550e8400-e29b-41d4-a716-446655440000,660e8400-e29b-41d4-a716-446655440001"
  - path_regex: ""
    barbican:
      - "770e8400-e29b-41d4-a716-446655440002"
      - "880e8400-e29b-41d4-a716-446655440003"
`)
	
	// Test string format (comma-separated)
	conf, err = parseCreationRuleForFile(parseConfigFile(mixedConfig, t), "/conf/path", "mixed_test", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 2, len(conf.KeyGroups[0]))

	// Test array format
	conf, err = parseCreationRuleForFile(parseConfigFile(mixedConfig, t), "/conf/path", "other_test", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 2, len(conf.KeyGroups[0]))

	// Test configuration with only auth settings (no keys)
	authOnlyConfig := []byte(`
creation_rules:
  - path_regex: ""
    barbican_auth_url: "https://keystone.example.com:5000/v3"
    barbican_region: "us-east-1"
`)
	conf, err = parseCreationRuleForFile(parseConfigFile(authOnlyConfig, t), "/conf/path", "test", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.KeyGroups))
	assert.Equal(t, 0, len(conf.KeyGroups[0])) // No keys should be created
}