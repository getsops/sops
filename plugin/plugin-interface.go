package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	// KeyTypeIdentifier is the string used to identify a Plugin MasterKey
	// in the metadata of an encrypted file.
	KeyTypeIdentifier = "plugin"
	// TimeoutFallback is the default timeout for plugin execution if not specified in the MasterKey config.
	TimeoutFallback = 10 * time.Second
	// MaxBinaryNameLength is a sanity check to prevent excessively long binary names that could cause DoSs.
	MaxBinaryNameLength = 128
	// MaxBinaryNameLength is another sanity check to prevent empty binary names (e.g., "sops-plugin-" is not a valid plugin binary name)
	MinBinaryNameLength = 1
)

// MasterKey is a generic plugin wrapper that satisfies the SOPS MasterKey interface.
// It bridges the SOPS Go Core with external external executables via stdin/stdout.
type MasterKey struct {
	BinaryName   string
	PluginConfig map[string]any
	InstanceID   string
	EncryptedKey string
	Timeout      string
	CreationDate time.Time
}

func NewMasterKey(
	binaryName string,
	config map[string]any,
	timeout string,
	instanceID string,
) *MasterKey {
	if instanceID == "" {
		instanceID = binaryName
	}
	return &MasterKey{
		BinaryName:   binaryName,
		InstanceID:   instanceID,
		PluginConfig: config,
		Timeout:      timeout,
		CreationDate: time.Now().UTC(),
	}
}

func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

func (key *MasterKey) ToString() string {
	return fmt.Sprintf("plugin:%s", key.BinaryName)
}

func (key *MasterKey) GetEnvPrefix() string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	normalized := reg.ReplaceAllString(key.InstanceID, "_")

	// we put a trailing underscore to make it easier for users to append their own suffixes in the plugin.
	// e.g, if instanceID is "my-vault", env prefix will be "SOPS_PLUGIN_MY_VAULT_",
	// and users can then use env vars like "SOPS_PLUGIN_MY_VAULT_TOKEN
	// or "SOPS_PLUGIN_MY_VAULT_KEY" in their plugin implementation.
	return "SOPS_PLUGIN_" + strings.ToUpper(normalized) + "_"
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
	// for Now, we have no criteria for rotation of plugin-based keys, so we'll just return false.
	// Maybe use a config inside of key??
	return false
}

// Encrypt takes a SOPS data key, encrypts it via the external plugin, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// validateBinaryName checks that the BinaryName field of the MasterKey is valid, not containing command injectable characters
// and do not exceed reasonable length limits to prevent DoSs.
//
// Valid nomenclature:
//   - Only alphanumeric characters, dashes, and underscores allowed.
//   - Length must be between 1 and 128 characters.
//   - The plugin binary path must follow "sops-plugin-<binaryName>" convention.
//   - Plugins are expected to be within the user's PATH for flexibility and security.
//   - Users are responsible for ensuring that their plugin binaries are secure and trustworthy.
func (key *MasterKey) validateBinaryName() error {
	if len(key.BinaryName) > MaxBinaryNameLength || len(key.BinaryName) < MinBinaryNameLength {
		return fmt.Errorf(
			"invalid binary name: length must be between %d and %d characters",
			MinBinaryNameLength,
			MaxBinaryNameLength,
		)
	}

	validBinaryName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validBinaryName.MatchString(key.BinaryName) {
		return fmt.Errorf("invalid binary name: only alphanumeric, dashes, and underscores allowed")
	}
	return nil
}

func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	req := map[string]any{
		"action":      "encrypt",
		"instance_id": key.InstanceID,
		"env_prefix":  key.GetEnvPrefix(),
		"config":      key.PluginConfig,
		"plaintext":   dataKey,
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, key.getTimeout())
		defer cancel()
	}

	err := key.validateBinaryName()
	if err != nil {
		return err
	}

	resp, err := executePlugin(ctx, key.BinaryName, req)
	if err != nil {
		return err
	}

	if resp.Ciphertext == "" {
		return fmt.Errorf("plugin did not return ciphertext")
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
		"action":      "decrypt",
		"instance_id": key.InstanceID,
		"env_prefix":  key.GetEnvPrefix(),
		"config":      key.PluginConfig,
		"ciphertext":  key.EncryptedKey,
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, key.getTimeout())
		defer cancel()
	}

	err := key.validateBinaryName()
	if err != nil {
		return nil, err
	}

	resp, err := executePlugin(ctx, key.BinaryName, req)
	if err != nil {
		return nil, err
	}

	if resp.Plaintext == nil {
		return nil, fmt.Errorf("plugin did not return plaintext")
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
func executePlugin(
	ctx context.Context,
	binaryName string,
	req map[string]any,
) (*PluginResponse, error) {
	// Binary naming convention: sops-plugin-<binaryName>
	executableName := fmt.Sprintf("sops-plugin-%s", binaryName)
	cmd := exec.CommandContext(ctx, executableName)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal plugin request: %v", err)
	}
	cmd.Stdin = bytes.NewReader(reqBytes)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("plugin execution timed out (%s)", executableName)
		}
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf(
				"plugin executable not found: %s. Please ensure that sops-plugin-%s is installed and in your PATH\nAvailable paths: %s",
				executableName,
				executableName,
				os.Getenv("PATH"),
			)
		}
		return nil, fmt.Errorf(
			"plugin execution failed (%s): %v. Stderr: %s",
			executableName,
			err,
			stderr.String(),
		)
	}

	var resp PluginResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf(
			"plugin %s violated IPC contract (invalid JSON): %v",
			executableName,
			err,
		)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("plugin %s error: %s", executableName, resp.Error)
	}

	return &resp, nil
}

// getTimeout is a helper function to parse the timeout string from the MasterKey config.
// falls back to 10s.
func (key *MasterKey) getTimeout() time.Duration {
	if key.Timeout != "" {
		if timeout, err := time.ParseDuration(key.Timeout); err == nil {
			return timeout
		}
	}
	return TimeoutFallback
}
