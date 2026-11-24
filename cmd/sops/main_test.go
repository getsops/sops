package main

import (
	"testing"

	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
)

func TestParsePluginKeys(t *testing.T) {
	for _, value := range []string{
		`{"name":"example","config":{"key_id":"one"}}`,
		`[{"name":"example","config":{"key_id":"one"}},{"name":"other","config":{}}]`,
	} {
		keys, err := parsePluginKeys(value)
		require.NoError(t, err)
		require.NotEmpty(t, keys)
		require.IsType(t, &sopsplugin.MasterKey{}, keys[0])
	}
}

func TestParsePluginKeysRejectsInvalidInput(t *testing.T) {
	for _, value := range []string{"null", "[]", `"example"`, `{"config":{}}`} {
		_, err := parsePluginKeys(value)
		require.Error(t, err, value)
	}
}
