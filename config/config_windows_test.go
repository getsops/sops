//go:build windows

package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRegexMatchesWindowsPaths(t *testing.T) {
	conf := parseConfigFile([]byte(`
creation_rules:
  - path_regex: '^secrets/.*\.env$'
    pgp: test
destination_rules:
  - path_regex: '^secrets/.*\.env$'
    s3_bucket: example
    recreation_rule:
      pgp: test
`), t)

	filePath := filepath.Join("secrets", "prod.env")
	configPath := filepath.Join(t.TempDir(), ".sops.yaml")
	creationPath := filepath.Join(filepath.Dir(configPath), filePath)

	creationConfig, err := parseCreationRuleForFile(conf, configPath, creationPath, nil)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "test", creationConfig.KeyGroups[0][0].ToString())

	destinationConfig, err := parseDestinationRuleForFile(conf, filePath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, destinationConfig.Destination)
}
