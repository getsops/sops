package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleConfigWithAWSSecretsManagerDestinationRules = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
destination_rules:
  - aws_region: "us-east-1"
    aws_secrets_manager_secret_name: "myapp/database"
    path_regex: "^secrets/.*"
  - aws_region: "us-west-2"
    aws_secrets_manager_secret_name: "api"
    path_regex: "^west-secrets/.*"
`)

var sampleConfigWithAWSParameterStoreDestinationRules = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
destination_rules:
  - aws_region: "us-east-1"
    aws_parameter_store_path: "/myapp/config"
    path_regex: "^parameters/.*"
  - aws_region: "us-west-2"
    aws_parameter_store_path: "/myapp/west/"
    path_regex: "^west-parameters/.*"
`)

var sampleConfigWithMixedAWSDestinationRules = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
destination_rules:
  - aws_region: "us-east-1"
    aws_secrets_manager_secret_name: "myapp/database"
    path_regex: "^secrets/.*"
  - aws_region: "us-east-1"
    aws_parameter_store_path: "/myapp/config"
    path_regex: "^parameters/.*"
  - s3_bucket: "mybucket"
    path_regex: "^s3/.*"
`)

func TestLoadConfigFileWithAWSSecretsManagerDestinationRules(t *testing.T) {
	conf, err := parseDestinationRuleForFile(parseConfigFile(sampleConfigWithAWSSecretsManagerDestinationRules, t), "secrets/database.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	path := conf.Destination.Path("database.yaml")
	assert.Contains(t, path, "arn:aws:secretsmanager:us-east-1:*:secret:myapp/database")

	// Test second rule with different region
	conf, err = parseDestinationRuleForFile(parseConfigFile(sampleConfigWithAWSSecretsManagerDestinationRules, t), "west-secrets/api.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	path = conf.Destination.Path("api.yaml")
	assert.Contains(t, path, "arn:aws:secretsmanager:us-west-2:*:secret:api")
}

func TestLoadConfigFileWithAWSParameterStoreDestinationRules(t *testing.T) {
	conf, err := parseDestinationRuleForFile(parseConfigFile(sampleConfigWithAWSParameterStoreDestinationRules, t), "parameters/app.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Equal(t, "/myapp/config", conf.Destination.Path("app.yaml"))

	// Test with path ending with slash
	conf, err = parseDestinationRuleForFile(parseConfigFile(sampleConfigWithAWSParameterStoreDestinationRules, t), "west-parameters/config.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Equal(t, "/myapp/west/config.yaml", conf.Destination.Path("config.yaml"))
}

func TestLoadConfigFileWithMixedAWSDestinationRules(t *testing.T) {
	// Test AWS Secrets Manager
	conf, err := parseDestinationRuleForFile(parseConfigFile(sampleConfigWithMixedAWSDestinationRules, t), "secrets/database.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Contains(t, conf.Destination.Path("database.yaml"), "arn:aws:secretsmanager:us-east-1:*:secret:myapp/database")

	// Test AWS Parameter Store
	conf, err = parseDestinationRuleForFile(parseConfigFile(sampleConfigWithMixedAWSDestinationRules, t), "parameters/config.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Equal(t, "/myapp/config", conf.Destination.Path("config.yaml"))

	// Test S3
	conf, err = parseDestinationRuleForFile(parseConfigFile(sampleConfigWithMixedAWSDestinationRules, t), "s3/backup.yaml", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Contains(t, conf.Destination.Path("backup.yaml"), "s3://mybucket/backup.yaml")
}

func TestValidateMultipleDestinationsInRule(t *testing.T) {
	invalidConfig := []byte(`
destination_rules:
  - aws_secrets_manager_secret_name: "my-secret"
    aws_parameter_store_path: "/my/path"
    path_regex: "^invalid/.*"
`)

	_, err := parseDestinationRuleForFile(parseConfigFile(invalidConfig, t), "invalid/test.yaml", nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "more than one destinations were found")
}

func TestValidateConflictingAWSDestinations(t *testing.T) {
	invalidConfig := []byte(`
destination_rules:
  - aws_secrets_manager_secret_name: "my-secret"
    s3_bucket: "mybucket"
    path_regex: "^invalid/.*"
`)

	_, err := parseDestinationRuleForFile(parseConfigFile(invalidConfig, t), "invalid/test.yaml", nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "more than one destinations were found")
}
