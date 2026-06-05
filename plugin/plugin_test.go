package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
