package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
	"go.yaml.in/yaml/v3"
)

// loadJSONSchema loads the JSON schema from the schema directory
func loadJSONSchema(t *testing.T) *gojsonschema.Schema {
	schemaPath := filepath.Join("..", "schema", "sops.json")
	schemaBytes, err := os.ReadFile(schemaPath)
	require.NoError(t, err, "Failed to read JSON schema file")

	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	require.NoError(t, err, "Failed to parse JSON schema")

	return schema
}

// validateYAMLAgainstSchema validates a YAML file against the JSON schema
func validateYAMLAgainstSchema(t *testing.T, schema *gojsonschema.Schema, yamlPath string) *gojsonschema.Result {
	yamlBytes, err := os.ReadFile(yamlPath)
	require.NoError(t, err, "Failed to read YAML file: %s", yamlPath)

	// Parse YAML to Go object
	var config interface{}
	err = yaml.Unmarshal(yamlBytes, &config)
	require.NoError(t, err, "Failed to parse YAML: %s", yamlPath)

	// Convert to JSON for schema validation
	jsonBytes, err := json.Marshal(config)
	require.NoError(t, err, "Failed to convert to JSON: %s", yamlPath)

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	require.NoError(t, err, "Schema validation failed with error: %s", yamlPath)

	return result
}

// TestSchemaValidTestCases tests that all valid test cases pass schema validation
func TestSchemaValidTestCases(t *testing.T) {
	schema := loadJSONSchema(t)

	validTestCases := []string{
		"valid-basic.yaml",
		"valid-complete.yaml",
		"valid-keygroups.yaml",
		"valid-stores.yaml",
		"valid-destination.yaml",
		"valid-azure.yaml",
		"valid-merge.yaml",
	}

	for _, testCase := range validTestCases {
		t.Run(testCase, func(t *testing.T) {
			testPath := filepath.Join("..", "schema", "test-cases", testCase)
			result := validateYAMLAgainstSchema(t, schema, testPath)

			if !result.Valid() {
				t.Errorf("Valid test case %s failed schema validation:", testCase)
				for _, err := range result.Errors() {
					t.Errorf("  - %s", err)
				}
			}
			assert.True(t, result.Valid(), "Valid test case should pass schema validation")
		})
	}
}

// TestSchemaInvalidTestCases tests that all invalid test cases fail schema validation
func TestSchemaInvalidTestCases(t *testing.T) {
	schema := loadJSONSchema(t)

	invalidTestCases := []string{
		"invalid-unknown-field.yaml",
		"invalid-shamir-threshold.yaml",
		"invalid-kms-missing-arn.yaml",
		"invalid-azure-missing-key.yaml",
		"invalid-vault-version.yaml",
		"invalid-stores-unknown.yaml",
	}

	for _, testCase := range invalidTestCases {
		t.Run(testCase, func(t *testing.T) {
			testPath := filepath.Join("..", "schema", "test-cases", testCase)
			result := validateYAMLAgainstSchema(t, schema, testPath)

			if result.Valid() {
				t.Errorf("Invalid test case %s passed schema validation but should have failed", testCase)
			}
			assert.False(t, result.Valid(), "Invalid test case should fail schema validation")

			// Log validation errors for debugging
			t.Logf("Expected validation errors for %s:", testCase)
			for _, err := range result.Errors() {
				t.Logf("  - %s", err)
			}
		})
	}
}

// TestSchemaAgainstRootSopsYaml tests the schema against the root .sops.yaml file
func TestSchemaAgainstRootSopsYaml(t *testing.T) {
	schema := loadJSONSchema(t)
	sopsYamlPath := filepath.Join("..", ".sops.yaml")

	// Check if the file exists
	if _, err := os.Stat(sopsYamlPath); os.IsNotExist(err) {
		t.Skip("Root .sops.yaml file does not exist")
		return
	}

	result := validateYAMLAgainstSchema(t, schema, sopsYamlPath)
	if !result.Valid() {
		t.Errorf("Root .sops.yaml failed schema validation:")
		for _, err := range result.Errors() {
			t.Errorf("  - %s", err)
		}
	}
	assert.True(t, result.Valid(), "Root .sops.yaml should pass schema validation")
}

// TestSchemaStructureMatchesConfig tests that schema structure aligns with config structs
func TestSchemaStructureMatchesConfig(t *testing.T) {
	schema := loadJSONSchema(t)

	// Test that basic creation_rule fields are accepted
	basicConfig := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"path_regex": "\\.yaml$",
				"pgp":        "ABC123",
				"age":        "age1xxx",
				"kms":        "arn:aws:kms:us-east-1:123456789012:key/xxx",
			},
		},
	}

	jsonBytes, err := json.Marshal(basicConfig)
	require.NoError(t, err)

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	require.NoError(t, err)
	assert.True(t, result.Valid(), "Basic config should be valid")
}

// TestSchemaKeyGroupsMergeField tests that the merge field in key_groups is supported
func TestSchemaKeyGroupsMergeField(t *testing.T) {
	schema := loadJSONSchema(t)

	// Test key_groups with merge field
	configWithMerge := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"key_groups": []map[string]interface{}{
					{
						"merge": []map[string]interface{}{
							{
								"pgp": []string{"ABC123"},
							},
							{
								"age": []string{"age1xxx"},
							},
						},
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(configWithMerge)
	require.NoError(t, err)

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	require.NoError(t, err)

	if !result.Valid() {
		for _, err := range result.Errors() {
			t.Logf("Validation error: %s", err)
		}
	}
	assert.True(t, result.Valid(), "Config with merge field should be valid")
}

// TestSchemaHCVaultFieldVariants tests both hc_vault and hc_vault_transit_uri
func TestSchemaHCVaultFieldVariants(t *testing.T) {
	schema := loadJSONSchema(t)

	// Test with hc_vault (short form)
	configWithHCVault := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"key_groups": []map[string]interface{}{
					{
						"hc_vault": []string{"https://vault.example.com/v1/transit/keys/my-key"},
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(configWithHCVault)
	require.NoError(t, err)

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	require.NoError(t, err)
	assert.True(t, result.Valid(), "Config with hc_vault should be valid")

	// Test with hc_vault_transit_uri (long form)
	configWithHCVaultTransit := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"hc_vault_transit_uri": "https://vault.example.com/v1/transit/keys/my-key",
			},
		},
	}

	jsonBytes, err = json.Marshal(configWithHCVaultTransit)
	require.NoError(t, err)

	documentLoader = gojsonschema.NewBytesLoader(jsonBytes)
	result, err = schema.Validate(documentLoader)
	require.NoError(t, err)
	assert.True(t, result.Valid(), "Config with hc_vault_transit_uri should be valid")
}

// TestSchemaArrayAndStringFormats tests that both string and array formats are accepted
func TestSchemaArrayAndStringFormats(t *testing.T) {
	schema := loadJSONSchema(t)

	// Test with string format (comma-separated)
	configWithStrings := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"pgp": "ABC123,DEF456",
				"age": "age1xxx,age2yyy",
			},
		},
	}

	jsonBytes, err := json.Marshal(configWithStrings)
	require.NoError(t, err)

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	require.NoError(t, err)
	assert.True(t, result.Valid(), "Config with string format should be valid")

	// Test with array format
	configWithArrays := map[string]interface{}{
		"creation_rules": []map[string]interface{}{
			{
				"pgp": []string{"ABC123", "DEF456"},
				"age": []string{"age1xxx", "age2yyy"},
			},
		},
	}

	jsonBytes, err = json.Marshal(configWithArrays)
	require.NoError(t, err)

	documentLoader = gojsonschema.NewBytesLoader(jsonBytes)
	result, err = schema.Validate(documentLoader)
	require.NoError(t, err)
	assert.True(t, result.Valid(), "Config with array format should be valid")
}

// TestSchemaRecreationRuleCompleteness tests that recreation_rule supports all creation_rule fields
func TestSchemaRecreationRuleCompleteness(t *testing.T) {
	schema := loadJSONSchema(t)

	// Test recreation_rule with various fields
	configWithRecreation := map[string]interface{}{
		"destination_rules": []map[string]interface{}{
			{
				"s3_bucket": "my-bucket",
				"recreation_rule": map[string]interface{}{
					"kms":                       "arn:aws:kms:us-east-1:123456789012:key/xxx",
					"pgp":                       "ABC123",
					"encrypted_regex":           "^(password|secret)",
					"shamir_threshold":          2,
					"mac_only_encrypted":        true,
					"unencrypted_suffix":        "_public",
					"encrypted_comment_regex":   "^encrypted:",
					"unencrypted_comment_regex": "^public:",
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(configWithRecreation)
	require.NoError(t, err)

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := schema.Validate(documentLoader)
	require.NoError(t, err)

	if !result.Valid() {
		for _, err := range result.Errors() {
			t.Logf("Validation error: %s", err)
		}
	}
	assert.True(t, result.Valid(), "Recreation rule with all fields should be valid")
}
