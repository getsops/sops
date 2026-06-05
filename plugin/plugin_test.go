package plugin

import (
	"os/exec"
    "os"
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
        fmt.Println("{\"ciphertext\": \"cifra-secreta\"}")
    } else if action == "decrypt" {
        fmt.Println("{\"plaintext\": \"dGV4dG8tY2xhcm8=\"}") // base64 de "texto-claro"
    }
}
`

func TestPluginIPC(t *testing.T) {
    tmpDir := t.TempDir()
    pluginPath := filepath.Join(tmpDir, "sops-plugin-dummy")

    srcFile := filepath.Join(tmpDir, "main.go")
    os.WriteFile(srcFile, []byte(dummyPluginCode), 0644)

    err := exec.Command("go", "build", "-o", pluginPath, srcFile).Run()
    assert.NoError(t, err, "failed to compile dummy plugin")

    t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"))

    key := NewMasterKey("dummy", map[string]any{"minha_config": "valor"}, "10s")

    err = key.Encrypt([]byte("texto-claro"))
    assert.NoError(t, err)
    assert.Equal(t, "cifra-secreta", string(key.EncryptedDataKey()))

    plaintext, err := key.Decrypt()
    assert.NoError(t, err)
    assert.Equal(t, "texto-claro", string(plaintext))
}
