package plugin

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "time"
)

const (
    // KeyTypeIdentifier is the string used to identify a Plugin MasterKey
    // in the metadata of an encrypted file.
    KeyTypeIdentifier = "plugin"
)

// MasterKey is a generic plugin wrapper that satisfies the SOPS MasterKey interface.
// It bridges the SOPS Go Core with external external executables via stdin/stdout.
type MasterKey struct {
    BinaryName   string
    PluginConfig map[string]any
    EncryptedKey string
    CreationDate time.Time
}

func NewMasterKey(binaryName string, config map[string]any) *MasterKey {
    return &MasterKey{
        BinaryName:   binaryName,
        PluginConfig: config,
        CreationDate: time.Now().UTC(),
    }
}

func (key *MasterKey) TypeToIdentifier() string {
    return KeyTypeIdentifier
}

func (key *MasterKey) ToString() string {
    return fmt.Sprintf("plugin:%s", key.BinaryName)
}

func (key MasterKey) ToMap() map[string]any {
    out := make(map[string]any)
    out["binary_name"] = key.BinaryName

    out["config"] = key.PluginConfig 

    out["enc"] = key.EncryptedKey
    out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
    return out
}

func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
    key.EncryptedKey = string(enc)
}

func (key *MasterKey) EncryptedDataKey() []byte {
    return []byte(key.EncryptedKey)
}

func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
    if key.EncryptedKey == "" {
        return key.Encrypt(dataKey)
    }
    return nil
}

func (key *MasterKey) NeedsRotation() bool {
    return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// Encrypt takes a SOPS data key, encrypts it via the external plugin, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
    return key.EncryptContext(context.Background(), dataKey)
}

func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
    req := map[string]any{
        "action":    "encrypt",
        "config":    key.PluginConfig,
        "plaintext": dataKey,
    }

    resp, err := executePlugin(ctx, key.BinaryName, req)
    if err != nil {
        return err
    }

    key.EncryptedKey = resp.Ciphertext
    return nil
}

// Decrypt decrypts the EncryptedKey field via the external plugin and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
    return key.DecryptContext(context.Background())
}

func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
    req := map[string]any{
        "action":     "decrypt",
        "config":     key.PluginConfig,
        "ciphertext": key.EncryptedKey,
    }

    resp, err := executePlugin(ctx, key.BinaryName, req)
    if err != nil {
        return nil, err
    }

    return resp.Plaintext, nil
}

// PluginResponse is the contract we expect from the plugin's stdout.
type PluginResponse struct {
    Plaintext  []byte `json:"plaintext,omitempty"`
    Ciphertext string `json:"ciphertext,omitempty"`
    Error      string `json:"error,omitempty"`
}

// executePlugin is the IPC Sandbox
func executePlugin(ctx context.Context, binaryName string, req map[string]any) (*PluginResponse, error) {
	// Binary naming convention: sops-plugin-<binaryName>
    executableName := fmt.Sprintf("sops-plugin-%s", binaryName)
    cmd := exec.CommandContext(ctx, executableName)

    reqBytes, _ := json.Marshal(req)
    cmd.Stdin = bytes.NewReader(reqBytes)

    var stdout bytes.Buffer
    cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("plugin execution failed (%s): %v. Stderr: %s", executableName, err, stderr.String())
    }

    var resp PluginResponse
    if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
        return nil, fmt.Errorf("plugin %s violated IPC contract (invalid JSON): %v", executableName, err)
    }

    if resp.Error != "" {
        return nil, fmt.Errorf("plugin %s error: %s", executableName, resp.Error)
    }

    return &resp, nil
}
