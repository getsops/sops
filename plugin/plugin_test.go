package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dummyPluginCode = `package main
import (
    "encoding/json"
    "fmt"
    "os"
)
func main() {
    var req map[string]any
    json.NewDecoder(os.Stdin).Decode(&req)

    action := req["action"].(string)
    if action == "encrypt" {
        fmt.Println("{\"ciphertext\": \"secret_cypher\"}")
    } else if action == "decrypt" {
        fmt.Println("{\"plaintext\": \"cGxhaW4tYXMtZGF5Cg==\"}")
    }
}
`

func TestPluginIPC(t *testing.T) {
	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, "sops-plugin-dummy")

	srcFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(srcFile, []byte(dummyPluginCode), 0o644)

	err := exec.Command("go", "build", "-o", pluginPath, srcFile).Run()
	assert.NoError(t, err, "failed to compile dummy plugin")

	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	key := NewMasterKey("dummy", map[string]any{"my_config": "value"}, "10s", "dummy")

	err = key.Encrypt([]byte("plain-as-day\n"))
	assert.NoError(t, err)
	assert.Equal(t, "secret_cypher", string(key.EncryptedDataKey()))

	plaintext, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, "plain-as-day\n", string(plaintext))
}

const flexiblePluginCode = `package main
import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)
func main() {
    if os.Getenv("MOCK_HANG") == "true" {
        time.Sleep(2 * time.Second)
        return
    }
    if os.Getenv("MOCK_INVALID_JSON") == "true" {
        fmt.Println("invalid-json-response")
        return
    }
    if os.Getenv("MOCK_ERROR_RESPONSE") == "true" {
        fmt.Println("{\"error\": \"plugin failed custom error\"}")
        return
    }
    if os.Getenv("MOCK_EMPTY_RESPONSE") == "true" {
        fmt.Println("{}")
        return
    }
    if os.Getenv("MOCK_STDERR_RESPONSE") == "true" {
        fmt.Fprintln(os.Stderr, "custom plugin stderr message")
        os.Exit(1)
        return
    }
    if os.Getenv("MOCK_STDERR_WITH_SUCCESS") == "true" {
        fmt.Fprintln(os.Stderr, "non-fatal warning message")
    }
    var req map[string]any
    if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
        fmt.Printf("{\"error\": \"decode error: %v\"}\n", err)
        return
    }
    action := req["action"].(string)
    if action == "encrypt" {
        ciphertext := req["plaintext"].(string)
        fmt.Printf("{\"ciphertext\": \"%s\"}\n", ciphertext)
    } else if action == "decrypt" {
        ciphertext := req["ciphertext"].(string)
        fmt.Printf("{\"plaintext\": \"%s\"}\n", ciphertext)
    }
}
`

func setupFlexiblePlugin(t *testing.T) *MasterKey {
	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, "sops-plugin-flexible")

	srcFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(srcFile, []byte(flexiblePluginCode), 0o644)
	require.NoError(t, err)

	err = exec.Command("go", "build", "-o", pluginPath, srcFile).Run()
	require.NoError(t, err, "failed to compile flexible plugin")

	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	return NewMasterKey("flexible", map[string]any{"conf": "val"}, "10s", "flexible")
}

func TestMasterKeyBasicMethods(t *testing.T) {
	key := NewMasterKey("my-binary", map[string]any{"foo": "bar"}, "5s", "my-vault")

	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
	assert.Equal(t, "plugin:my-binary", key.ToString())
	assert.Equal(t, "SOPS_PLUGIN_MY_VAULT_", key.GetEnvPrefix())

	// Test default instance ID
	keyDefaultID := NewMasterKey("my-binary", nil, "", "")
	assert.Equal(t, "my-binary", keyDefaultID.InstanceID)
	assert.Equal(t, "SOPS_PLUGIN_MY_BINARY_", keyDefaultID.GetEnvPrefix())

	// Test ToMap
	key.SetEncryptedDataKey([]byte("encrypted-payload"))
	assert.Equal(t, []byte("encrypted-payload"), key.EncryptedDataKey())

	m := key.ToMap()
	assert.Equal(t, "my-binary", m["binary_name"])
	assert.Equal(t, map[string]any{"foo": "bar"}, m["config"])
	assert.Equal(t, "encrypted-payload", m["enc"])
	assert.NotEmpty(t, m["created_at"])

	// Test NeedsRotation
	assert.False(t, key.NeedsRotation())
}

func TestMasterKeyValidation(t *testing.T) {
	tests := []struct {
		name       string
		binaryName string
		wantErr    string
	}{
		{
			name:       "valid name",
			binaryName: "valid-plugin_123",
			wantErr:    "",
		},
		{
			name:       "invalid char slash",
			binaryName: "plugin/path",
			wantErr:    "invalid binary name: only alphanumeric, dashes, and underscores allowed",
		},
		{
			name:       "invalid char space",
			binaryName: "plugin name",
			wantErr:    "invalid binary name: only alphanumeric, dashes, and underscores allowed",
		},
		{
			name:       "empty name",
			binaryName: "",
			wantErr:    "invalid binary name: length must be between 1 and 128 characters",
		},
		{
			name:       "too long name",
			binaryName: strings.Repeat("a", 129),
			wantErr:    "invalid binary name: length must be between 1 and 128 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewMasterKey(tt.binaryName, nil, "", "")
			err := key.Encrypt([]byte("data"))
			if tt.wantErr == "" {
				if err != nil {
					assert.NotContains(t, err.Error(), "invalid binary name")
				}
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			_, err = key.Decrypt()
			if tt.wantErr == "" {
				if err != nil {
					assert.NotContains(t, err.Error(), "invalid binary name")
				}
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestMasterKeyEncryptIfNeeded(t *testing.T) {
	// If EncryptedKey is empty, EncryptIfNeeded should call Encrypt.
	// Since no plugin is present, it will fail.
	key := NewMasterKey("non-existent-plugin", nil, "", "")
	err := key.EncryptIfNeeded([]byte("data"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin executable not found")

	// If EncryptedKey is not empty, EncryptIfNeeded should do nothing and return nil.
	key.EncryptedKey = "already-encrypted"
	err = key.EncryptIfNeeded([]byte("data"))
	assert.NoError(t, err)
	assert.Equal(t, "already-encrypted", key.EncryptedKey)
}

func TestPluginHappyPath(t *testing.T) {
	key := setupFlexiblePlugin(t)

	// Encrypt
	dataKey := []byte("hello-world")
	err := key.Encrypt(dataKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, key.EncryptedKey)

	// Decrypt
	plaintext, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, dataKey, plaintext)
}

func TestPluginInvalidJSON(t *testing.T) {
	key := setupFlexiblePlugin(t)
	t.Setenv("MOCK_INVALID_JSON", "true")

	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violated IPC contract (invalid JSON)")

	_, err = key.Decrypt()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violated IPC contract (invalid JSON)")
}

func TestPluginErrorResponse(t *testing.T) {
	key := setupFlexiblePlugin(t)
	t.Setenv("MOCK_ERROR_RESPONSE", "true")

	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin sops-plugin-flexible error: plugin failed custom error")

	_, err = key.Decrypt()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin sops-plugin-flexible error: plugin failed custom error")
}

func TestPluginEmptyResponse(t *testing.T) {
	key := setupFlexiblePlugin(t)
	t.Setenv("MOCK_EMPTY_RESPONSE", "true")

	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin did not return ciphertext")

	_, err = key.Decrypt()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin did not return plaintext")
}

func TestPluginTimeout(t *testing.T) {
	key := setupFlexiblePlugin(t)
	key.Timeout = "100ms"
	t.Setenv("MOCK_HANG", "true")

	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin execution timed out")

	// Decrypt uses default timeout or context deadline. Let's pass a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err = key.DecryptContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin execution timed out")
}

func TestPluginNotFound(t *testing.T) {
	key := NewMasterKey("non-existent-plugin-12345", nil, "", "")
	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin executable not found")
}

func TestPluginRequestMarshalError(t *testing.T) {
	// Put an unmarshallable channel in the plugin config
	key := NewMasterKey("flexible", map[string]any{"bad": make(chan int)}, "", "")
	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal plugin request")
}

func TestPluginStderrResponse(t *testing.T) {
	key := setupFlexiblePlugin(t)
	t.Setenv("MOCK_STDERR_RESPONSE", "true")

	err := key.Encrypt([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin execution failed")
	assert.Contains(t, err.Error(), "custom plugin stderr message")

	_, err = key.Decrypt()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin execution failed")
	assert.Contains(t, err.Error(), "custom plugin stderr message")
}

func TestPluginStderrWithSuccess(t *testing.T) {
	key := setupFlexiblePlugin(t)
	t.Setenv("MOCK_STDERR_WITH_SUCCESS", "true")

	dataKey := []byte("hello-world")
	err := key.Encrypt(dataKey)
	assert.NoError(t, err)

	plaintext, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, dataKey, plaintext)
}

func TestPluginConcurrency(t *testing.T) {
	// Setup the flexible plugin binary compiled once
	_ = setupFlexiblePlugin(t)

	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Create a distinct MasterKey instance for this goroutine referencing the same binary name.
			// This simulates two different processes/calls accessing the same plugin concurrently.
			goroutineKey := NewMasterKey("flexible", map[string]any{"conf": "val"}, "10s", "flexible")

			dataKey := []byte(fmt.Sprintf("secret-payload-%d", id))

			// Test encrypt
			err := goroutineKey.Encrypt(dataKey)
			assert.NoError(t, err)
			assert.NotEmpty(t, goroutineKey.EncryptedKey)

			// Test decrypt
			plaintext, err := goroutineKey.Decrypt()
			assert.NoError(t, err)
			assert.Equal(t, dataKey, plaintext)
		}(i)
	}

	wg.Wait()
}

