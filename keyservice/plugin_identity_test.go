package keyservice

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluginPromptIdentityDistinguishesConfigurationWithoutShowingIt(t *testing.T) {
	first, err := json.Marshal(map[string]string{"key_id": "first-secret-looking-value"})
	require.NoError(t, err)
	second, err := json.Marshal(map[string]string{"key_id": "second-secret-looking-value"})
	require.NoError(t, err)

	firstLabel := keyToString(&Key{KeyType: &Key_PluginKey{PluginKey: &PluginKey{
		Name:       "kms",
		ConfigJson: first,
	}}})
	secondLabel := keyToString(&Key{KeyType: &Key_PluginKey{PluginKey: &PluginKey{
		Name:       "kms",
		ConfigJson: second,
	}}})

	require.NotEqual(t, firstLabel, secondLabel)
	require.Contains(t, firstLabel, `"kms"`)
	require.NotContains(t, firstLabel, "first-secret-looking-value")
	require.NotContains(t, secondLabel, "second-secret-looking-value")
	require.True(t, strings.HasPrefix(firstLabel, "SOPS plugin "))
}
