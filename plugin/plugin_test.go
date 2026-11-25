package plugin_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keyservice"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestRuntimeReusesPluginProcessAndClosesIt(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	setAllowlist(t, binary, "lifecycle")

	runtime := sopsplugin.NewRuntime(context.Background())
	t.Cleanup(func() {
		require.NoError(t, runtime.Close())
	})
	dataKey := []byte("01234567890123456789012345678901")
	const concurrentWraps = 8
	type wrapResult struct {
		response sopsplugin.WrapKeyResponse
		err      error
	}
	results := make(chan wrapResult, concurrentWraps)
	var waitGroup sync.WaitGroup
	for range concurrentWraps {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			response, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
				DataKey: dataKey,
			})
			results <- wrapResult{response: response, err: err}
		}()
	}
	waitGroup.Wait()
	close(results)
	var wrapped sopsplugin.WrapKeyResponse
	for result := range results {
		require.NoError(t, result.err)
		if wrapped.EncryptedDataKey == nil {
			wrapped = result.response
		} else {
			require.Equal(t, wrapped.Metadata["process_id"], result.response.Metadata["process_id"])
		}
	}
	processID, ok := wrapped.Metadata["process_id"].(string)
	require.True(t, ok)

	unwrapped, err := runtime.UnwrapKey(context.Background(), "lifecycle", sopsplugin.UnwrapKeyRequest{
		Metadata:         wrapped.Metadata,
		EncryptedDataKey: wrapped.EncryptedDataKey,
	})
	require.NoError(t, err)
	require.Equal(t, dataKey, unwrapped.DataKey)

	again, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: dataKey,
	})
	require.NoError(t, err)
	require.Equal(t, processID, again.Metadata["process_id"])

	require.NoError(t, runtime.Close())
	_, err = runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: dataKey,
	})
	require.ErrorContains(t, err, "runtime is closed")
}

func TestRuntimeKeepsPluginAvailableAfterExplicitCancellation(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	setAllowlist(t, binary, "lifecycle")

	runtime := sopsplugin.NewRuntime(context.Background())
	t.Cleanup(func() {
		require.NoError(t, runtime.Close())
	})
	dataKey := []byte("01234567890123456789012345678901")
	first, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: dataKey,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	startedFile := filepath.Join(t.TempDir(), "request-started")
	result := make(chan error, 1)
	go func() {
		_, err := runtime.WrapKey(ctx, "lifecycle", sopsplugin.WrapKeyRequest{
			Config: map[string]interface{}{
				"delay_ms":     10_000,
				"started_file": startedFile,
			},
			DataKey: dataKey,
		})
		result <- err
	}()
	startupDeadline := time.After(5 * time.Second)
	poll := time.NewTicker(10 * time.Millisecond)
	defer poll.Stop()
	started := false
	for !started {
		select {
		case err := <-result:
			cancel()
			require.NoError(t, err)
			t.Fatal("plugin returned before the test canceled the request")
		case <-startupDeadline:
			cancel()
			t.Fatal("plugin did not start the delayed request")
		case <-poll.C:
			if _, err := os.Stat(startedFile); err == nil {
				started = true
			}
		}
	}

	cancel()
	select {
	case err := <-result:
		require.ErrorContains(t, err, "Canceled")
	case <-time.After(5 * time.Second):
		t.Fatal("plugin call did not stop after context cancellation")
	}

	afterCancellation, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: dataKey,
	})
	require.NoError(t, err)
	require.Equal(t, first.Metadata["process_id"], afterCancellation.Metadata["process_id"])
}

func TestRuntimeEvictsPluginAfterInvocationTimeout(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	for _, operation := range []string{"Wrap", "Unwrap"} {
		t.Run(operation, func(t *testing.T) {
			setAllowlistWithInvocationTimeout(t, binary, "lifecycle", 100*time.Millisecond)
			runtime := sopsplugin.NewRuntime(context.Background())
			t.Cleanup(func() {
				require.NoError(t, runtime.Close())
			})
			dataKey := []byte("01234567890123456789012345678901")
			first, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
				DataKey: dataKey,
			})
			require.NoError(t, err)

			started := time.Now()
			switch operation {
			case "Wrap":
				_, err = runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
					Config: map[string]interface{}{
						"delay_ms":            10_000,
						"ignore_cancellation": true,
					},
					DataKey: dataKey,
				})
			case "Unwrap":
				_, err = runtime.UnwrapKey(context.Background(), "lifecycle", sopsplugin.UnwrapKeyRequest{
					Config: map[string]interface{}{
						"delay_ms":            10_000,
						"ignore_cancellation": true,
					},
					Metadata:         first.Metadata,
					EncryptedDataKey: first.EncryptedDataKey,
				})
			}
			require.ErrorContains(t, err, "DeadlineExceeded")
			require.Less(t, time.Since(started), 5*time.Second)

			afterTimeout, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
				DataKey: dataKey,
			})
			require.NoError(t, err)
			require.NotEqual(t, first.Metadata["process_id"], afterTimeout.Metadata["process_id"])
		})
	}
}

func TestTimedOutPluginEvictsConcurrentCalls(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	setAllowlistWithInvocationTimeout(t, binary, "lifecycle", 10*time.Second)

	runtime := sopsplugin.NewRuntime(context.Background())
	t.Cleanup(func() {
		require.NoError(t, runtime.Close())
	})
	dataKey := []byte("01234567890123456789012345678901")
	first, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: dataKey,
	})
	require.NoError(t, err)

	firstContext, cancelFirst := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFirst()
	firstStarted := filepath.Join(t.TempDir(), "first-started")
	firstResult := make(chan error, 1)
	go func() {
		_, err := runtime.WrapKey(firstContext, "lifecycle", sopsplugin.WrapKeyRequest{
			Config: map[string]interface{}{
				"delay_ms":            10_000,
				"ignore_cancellation": true,
				"started_file":        firstStarted,
			},
			DataKey: dataKey,
		})
		firstResult <- err
	}()
	waitForTestFile(t, firstStarted)

	secondStarted := filepath.Join(t.TempDir(), "second-started")
	secondResult := make(chan error, 1)
	go func() {
		_, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
			Config: map[string]interface{}{
				"delay_ms":            10_000,
				"ignore_cancellation": true,
				"started_file":        secondStarted,
			},
			DataKey: dataKey,
		})
		secondResult <- err
	}()
	waitForTestFile(t, secondStarted)

	select {
	case err := <-firstResult:
		require.ErrorContains(t, err, "DeadlineExceeded")
	case <-time.After(5 * time.Second):
		t.Fatal("first plugin call did not time out")
	}
	select {
	case err := <-secondResult:
		require.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent plugin call was not interrupted by process eviction")
	}

	afterTimeout, err := runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: dataKey,
	})
	require.NoError(t, err)
	require.NotEqual(t, first.Metadata["process_id"], afterTimeout.Metadata["process_id"])
}

func TestOneShotPluginInvocationTimeout(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	setAllowlistWithInvocationTimeout(t, binary, "lifecycle", 100*time.Millisecond)

	started := time.Now()
	_, err := sopsplugin.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		Config: map[string]interface{}{
			"delay_ms":            10_000,
			"ignore_cancellation": true,
		},
		DataKey: []byte("01234567890123456789012345678901"),
	})
	require.ErrorContains(t, err, "DeadlineExceeded")
	require.Less(t, time.Since(started), 5*time.Second)
}

func TestPluginInvocationTimeoutThroughLocalKeyService(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	setAllowlistWithInvocationTimeout(t, binary, "lifecycle", 100*time.Millisecond)
	key, err := sopsplugin.NewMasterKey("lifecycle", map[string]interface{}{
		"delay_ms":            10_000,
		"ignore_cancellation": true,
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	started := time.Now()
	errs := metadata.UpdateMasterKeysWithKeyServices(dataKey, services)
	require.NotEmpty(t, errs)
	require.ErrorContains(t, errors.Join(errs...), "DeadlineExceeded")
	require.Less(t, time.Since(started), 5*time.Second)

	key.Config = map[string]interface{}{}
	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
}

func TestRuntimeCancelsPluginStartupWithRequestContext(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/lifecyclefixture")
	setAllowlist(t, binary, "lifecycle")
	t.Setenv("SOPS_TEST_PLUGIN_START_DELAY", "1")

	runtime := sopsplugin.NewRuntime(context.Background())
	t.Cleanup(func() {
		require.NoError(t, runtime.Close())
	})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err := runtime.WrapKey(ctx, "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: []byte("01234567890123456789012345678901"),
	})
	require.Error(t, err)
	require.Less(t, time.Since(started), 5*time.Second)

	t.Setenv("SOPS_TEST_PLUGIN_START_DELAY", "")
	_, err = runtime.WrapKey(context.Background(), "lifecycle", sopsplugin.WrapKeyRequest{
		DataKey: []byte("01234567890123456789012345678901"),
	})
	require.NoError(t, err)
}

func newPluginTestLocalClient(t testing.TB) keyservice.LocalClient {
	t.Helper()
	client := keyservice.NewLocalClient()
	t.Cleanup(func() {
		require.NoError(t, client.Close())
	})
	return client
}

func TestDecodeHelpersPreserveLargeIntegers(t *testing.T) {
	const largeInteger = "9007199254740993"
	config := struct {
		Value int64 `json:"value"`
	}{}
	require.NoError(t, sopsplugin.DecodeConfig(
		map[string]interface{}{"value": json.Number(largeInteger)},
		&config,
	))
	require.Equal(t, int64(9007199254740993), config.Value)

	metadata := struct {
		Version uint64 `json:"version"`
	}{}
	require.NoError(t, sopsplugin.DecodeMetadata(
		map[string]interface{}{"version": json.Number(largeInteger)},
		&metadata,
	))
	require.Equal(t, uint64(9007199254740993), metadata.Version)
}

func TestPluginViaSOPSKeyService(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/fixture")
	setAllowlist(t, binary, "fixture")

	key, err := sopsplugin.NewMasterKey("fixture", map[string]interface{}{"label": "test-key"})
	require.NoError(t, err)
	metadata := sops.Metadata{
		KeyGroups: []sops.KeyGroup{{key}},
	}
	dataKey := []byte("01234567890123456789012345678901")
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(dataKey, services))
	require.NotEmpty(t, key.EncryptedKey)
	require.Equal(t, map[string]interface{}{"label": "test-key"}, key.ProviderMetadata)

	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, dataKey, decoded)
}

func TestPluginRewrapUsesPersistedMetadataAfterPublicKeyFileIsRemoved(t *testing.T) {
	binary := buildFixturePlugin(t, "github.com/getsops/sops/v3/plugin/testdata/azkvfixture")
	setAllowlist(t, binary, "azkv-test")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	keyDir := t.TempDir()
	privateKeyPath := filepath.Join(keyDir, "private.pem")
	publicKeyPath := filepath.Join(keyDir, "public.pem")
	require.NoError(t, os.WriteFile(privateKeyPath, pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}), 0600))
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	require.NoError(t, os.WriteFile(publicKeyPath, publicKeyPEM, 0600))

	key, err := sopsplugin.NewMasterKey("azkv-test", map[string]interface{}{
		"vaultUrl":       "https://example.vault.azure.net",
		"key":            "sops",
		"version":        "test-version",
		"publicKeyFile":  publicKeyPath,
		"privateKeyFile": privateKeyPath,
	})
	require.NoError(t, err)
	metadata := sops.Metadata{KeyGroups: []sops.KeyGroup{{key}}}
	services := []keyservice.KeyServiceClient{newPluginTestLocalClient(t)}
	firstDataKey := []byte("01234567890123456789012345678901")
	secondDataKey := []byte("abcdefghijabcdefghijabcdefghij12")

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(firstDataKey, services))
	require.Equal(t, string(publicKeyPEM), key.ProviderMetadata["public_key"])
	require.NoError(t, os.Remove(publicKeyPath))

	require.Empty(t, metadata.UpdateMasterKeysWithKeyServices(secondDataKey, services))
	metadata.DataKey = nil
	decoded, err := metadata.GetDataKeyWithKeyServices(services, nil)
	require.NoError(t, err)
	require.Equal(t, secondDataKey, decoded)
}

func TestPluginRequiresLocalAllowlist(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "plugins.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte("plugins: {}\n"), 0600))
	t.Setenv("SOPS_PLUGIN_CONFIG", configPath)

	_, err := sopsplugin.WrapKey(context.Background(), "not-allowed", sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{},
		DataKey: []byte("01234567890123456789012345678901"),
	})
	require.ErrorContains(t, err, "not in the local allowlist")
}

func TestPluginNameCannotSelectAnExecutablePath(t *testing.T) {
	_, err := sopsplugin.WrapKey(context.Background(), "../not-allowed", sopsplugin.WrapKeyRequest{
		Config:  map[string]interface{}{},
		DataKey: []byte("01234567890123456789012345678901"),
	})
	require.ErrorContains(t, err, "invalid plugin name")
}

func TestPluginRejectsInvalidKeyMaterialLengths(t *testing.T) {
	for _, size := range []int{31, 34} {
		t.Run(fmt.Sprintf("%d bytes", size), func(t *testing.T) {
			_, err := sopsplugin.WrapKey(context.Background(), "unused", sopsplugin.WrapKeyRequest{
				DataKey: make([]byte, size),
			})
			require.ErrorContains(t, err, "must be 32 or 33 bytes")
		})
	}
}

func TestInvalidPluginInvocationTimeoutIsRejected(t *testing.T) {
	for _, timeout := range []string{"not-a-duration", "0s", "-1s"} {
		t.Run(timeout, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "plugins.yaml")
			content, err := yaml.Marshal(map[string]interface{}{
				"plugins": map[string]interface{}{
					"unused": map[string]interface{}{
						"path":               filepath.Join(t.TempDir(), "plugin"),
						"invocation_timeout": timeout,
					},
				},
			})
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(configPath, content, 0600))
			t.Setenv("SOPS_PLUGIN_CONFIG", configPath)

			_, err = sopsplugin.WrapKey(context.Background(), "unused", sopsplugin.WrapKeyRequest{
				DataKey: []byte("01234567890123456789012345678901"),
			})
			if timeout == "not-a-duration" {
				require.ErrorContains(t, err, "parse plugin")
			} else {
				require.ErrorContains(t, err, "invocation_timeout must be a positive Go duration")
			}
		})
	}
}

func buildFixturePlugin(t *testing.T, packagePath string) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	require.True(t, ok)
	projectRoot := filepath.Dir(filepath.Dir(source))
	name := "sops-plugin-fixture"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(t.TempDir(), name)
	command := exec.Command("go", "build", "-o", binary, packagePath)
	command.Dir = projectRoot
	output, err := command.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("build fixture plugin: %s", output))
	return binary
}

func setAllowlist(t *testing.T, binary, name string) {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "plugins.yaml")
	content, err := yaml.Marshal(map[string]interface{}{
		"plugins": map[string]string{name: binary},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configPath, content, 0600))
	t.Setenv("SOPS_PLUGIN_CONFIG", configPath)
}

func setAllowlistWithInvocationTimeout(t *testing.T, binary, name string, timeout time.Duration) {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "plugins.yaml")
	content, err := yaml.Marshal(map[string]interface{}{
		"plugins": map[string]interface{}{
			name: map[string]interface{}{
				"path":               binary,
				"invocation_timeout": timeout.String(),
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configPath, content, 0600))
	t.Setenv("SOPS_PLUGIN_CONFIG", configPath)
}

func waitForTestFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	poll := time.NewTicker(10 * time.Millisecond)
	defer poll.Stop()
	for {
		if _, err := os.Stat(path); err == nil {
			return
		} else if !errors.Is(err, os.ErrNotExist) {
			require.NoError(t, err)
		}
		select {
		case <-deadline.C:
			t.Fatalf("timed out waiting for %q", path)
		case <-poll.C:
		}
	}
}
