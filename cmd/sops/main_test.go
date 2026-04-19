package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3/config"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

const nonMatchingCreationRuleConfig = `creation_rules:
  - path_regex: something-else/.*\.(json|yaml|yml|env|txt)$
    age: age15sq7kls08hzq8djpn26dda0fna3ccnw038568gcul9amjjjdaedq4xg2rr
`

const matchingCreationRuleConfig = `creation_rules:
  - path_regex: ""
    age: age15sq7kls08hzq8djpn26dda0fna3ccnw038568gcul9amjjjdaedq4xg2rr
`

func newTestCLIContext(t *testing.T, configPath string, inlineFlags map[string]string) *cli.Context {
	t.Helper()

	app := cli.NewApp()

	globalSet := flag.NewFlagSet("global", flag.ContinueOnError)
	globalSet.String("config", "", "")
	require.NoError(t, globalSet.Set("config", configPath))
	globalCtx := cli.NewContext(app, globalSet, nil)

	localSet := flag.NewFlagSet("local", flag.ContinueOnError)
	for _, name := range []string{"kms", "pgp", "gcp-kms", "hckms", "azure-kv", "hc-vault-transit", "age"} {
		localSet.String(name, "", "")
	}
	for name, value := range inlineFlags {
		require.NoError(t, localSet.Set(name, value))
	}

	return cli.NewContext(app, localSet, globalCtx)
}

func writeConfigFile(t *testing.T, dir string, contents string) string {
	t.Helper()

	configPath := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(contents), 0o600))
	return configPath
}

func TestLoadConfigIgnoresNonMatchingCreationRulesWhenInlineKeysAreProvided(t *testing.T) {
	dir := t.TempDir()
	configPath := writeConfigFile(t, dir, nonMatchingCreationRuleConfig)
	ctx := newTestCLIContext(t, configPath, map[string]string{
		"age": "age1xxfdafu5j4e5z7y5l6my6x07vjuh6unxersnwne4etpvykheq9gsj003fv",
	})

	conf, err := loadConfig(ctx, filepath.Join(dir, "secret.json"), nil)
	require.NoError(t, err)
	require.Nil(t, conf)
}

func TestLoadConfigReturnsNonMatchingCreationRuleErrorWithoutInlineKeys(t *testing.T) {
	dir := t.TempDir()
	configPath := writeConfigFile(t, dir, nonMatchingCreationRuleConfig)
	ctx := newTestCLIContext(t, configPath, nil)

	conf, err := loadConfig(ctx, filepath.Join(dir, "secret.json"), nil)
	require.Nil(t, conf)
	require.ErrorIs(t, err, config.ErrNoMatchingCreationRules)
}

func TestLoadConfigStillLoadsMatchingCreationRulesWithInlineKeys(t *testing.T) {
	dir := t.TempDir()
	configPath := writeConfigFile(t, dir, matchingCreationRuleConfig)
	ctx := newTestCLIContext(t, configPath, map[string]string{
		"age": "age1xxfdafu5j4e5z7y5l6my6x07vjuh6unxersnwne4etpvykheq9gsj003fv",
	})

	conf, err := loadConfig(ctx, filepath.Join(dir, "secret.json"), nil)
	require.NoError(t, err)
	require.NotNil(t, conf)
	require.Len(t, conf.KeyGroups, 1)
	require.Len(t, conf.KeyGroups[0], 1)
}
