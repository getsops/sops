package plugin_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keyservice"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestPluginViaSOPSKeyService(t *testing.T) {
	binary := buildFixturePlugin(t)
	setAllowlist(t, binary)

	key, err := sopsplugin.NewMasterKey("fixture", map[string]interface{}{"label": "test-key"})
	require.NoError(t, err)
	metadata := sops.Metadata{
		KeyGroups: []sops.KeyGroup{{key}},
	}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{keyservice.NewLocalClient()}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.NotEmpty(t, key.EncryptedKey)
	require.Equal(t, map[string]interface{}{"label": "test-key"}, key.ProviderMetadata)

	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
}

func TestPluginRequiresLocalAllowlist(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "plugins.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte("plugins: {}\n"), 0600))
	t.Setenv("SOPS_PLUGIN_CONFIG", configPath)

	_, err := sopsplugin.WrapKey(context.Background(), "not-allowed", sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{},
		DataKey: []byte("data key"),
	})
	require.ErrorContains(t, err, "not in the local allowlist")
}

func TestPluginNameCannotSelectAnExecutablePath(t *testing.T) {
	_, err := sopsplugin.WrapKey(context.Background(), "../not-allowed", sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{},
		DataKey: []byte("data key"),
	})
	require.ErrorContains(t, err, "invalid plugin name")
}

func buildFixturePlugin(t *testing.T) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	require.True(t, ok)
	projectRoot := filepath.Dir(filepath.Dir(source))
	name := "sops-plugin-fixture"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(t.TempDir(), name)
	command := exec.Command("go", "build", "-o", binary, "github.com/getsops/sops/v3/plugin/testdata/fixture")
	command.Dir = projectRoot
	output, err := command.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("build fixture plugin: %s", output))
	return binary
}

func setAllowlist(t *testing.T, binary string) {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "plugins.yaml")
	content, err := yaml.Marshal(map[string]interface{}{
		"plugins": map[string]string{"fixture": binary},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configPath, content, 0600))
	t.Setenv("SOPS_PLUGIN_CONFIG", configPath)
}
