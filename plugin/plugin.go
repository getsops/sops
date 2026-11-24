// Package plugin provides the Go SDK and process runtime for SOPS key-provider
// plugins.
package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"
	"go.yaml.in/yaml/v3"
)

const (
	pluginName       = "sops_key_provider"
	magicCookieKey   = "SOPS_PLUGIN"
	magicCookieValue = "sops-plugin"
)

var pluginNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// WrapKeyRequest contains the SOPS data key, public configuration, and
// metadata from any previous wrapping operation.
type WrapKeyRequest struct {
	Config map[string]interface{}
	// Metadata is empty on the first wrap and contains the previous response's
	// metadata when rewrapping an existing key.
	Metadata map[string]interface{}
	DataKey  []byte
}

// WrapKeyResponse contains the provider's opaque wrapped data key and
// metadata that replaces the previous metadata for this key.
type WrapKeyResponse struct {
	EncryptedDataKey []byte
	Metadata         map[string]interface{}
}

// UnwrapKeyRequest contains the wrapped data key, provider configuration, and
// per-encryption metadata.
type UnwrapKeyRequest struct {
	Config           map[string]interface{}
	Metadata         map[string]interface{}
	EncryptedDataKey []byte
}

// UnwrapKeyResponse contains the unwrapped SOPS data key.
type UnwrapKeyResponse struct {
	DataKey []byte
}

// KeyProvider is implemented by a SOPS plugin process.
type KeyProvider interface {
	// Wrap should honor context cancellation while wrapping the data key.
	Wrap(context.Context, WrapKeyRequest) (WrapKeyResponse, error)
	// Unwrap should honor context cancellation while unwrapping the data key.
	Unwrap(context.Context, UnwrapKeyRequest) (UnwrapKeyResponse, error)
}

// Serve starts a key-provider plugin process. It does not return.
func Serve(provider KeyProvider) {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: hplugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   magicCookieKey,
			MagicCookieValue: magicCookieValue,
		},
		Plugins: hplugin.PluginSet{
			pluginName: &grpcPlugin{Impl: provider},
		},
		GRPCServer: hplugin.DefaultGRPCServer,
		Logger:     hclog.NewNullLogger(),
	})
}

// WrapKey starts an allowlisted plugin and asks it to wrap a SOPS data key.
func WrapKey(ctx context.Context, name string, request WrapKeyRequest) (WrapKeyResponse, error) {
	if len(request.DataKey) == 0 {
		return WrapKeyResponse{}, fmt.Errorf("cannot wrap an empty data key")
	}
	if _, err := encodeObject(request.Config); err != nil {
		return WrapKeyResponse{}, fmt.Errorf("invalid plugin configuration: %w", err)
	}
	if _, err := encodeObject(request.Metadata); err != nil {
		return WrapKeyResponse{}, fmt.Errorf("invalid plugin metadata: %w", err)
	}
	provider, close, err := start(ctx, name)
	if err != nil {
		return WrapKeyResponse{}, err
	}
	defer close()

	response, err := provider.Wrap(ctx, request)
	if err != nil {
		return WrapKeyResponse{}, fmt.Errorf("plugin %q failed to wrap the data key: %w", name, err)
	}
	if len(response.EncryptedDataKey) == 0 {
		return WrapKeyResponse{}, fmt.Errorf("plugin %q returned an empty wrapped data key", name)
	}
	if _, err := encodeObject(response.Metadata); err != nil {
		return WrapKeyResponse{}, fmt.Errorf("plugin %q returned invalid metadata: %w", name, err)
	}
	return response, nil
}

// UnwrapKey starts an allowlisted plugin and asks it to unwrap a SOPS data key.
func UnwrapKey(ctx context.Context, name string, request UnwrapKeyRequest) (UnwrapKeyResponse, error) {
	if len(request.EncryptedDataKey) == 0 {
		return UnwrapKeyResponse{}, fmt.Errorf("cannot unwrap an empty data key")
	}
	if _, err := encodeObject(request.Config); err != nil {
		return UnwrapKeyResponse{}, fmt.Errorf("invalid plugin configuration: %w", err)
	}
	if _, err := encodeObject(request.Metadata); err != nil {
		return UnwrapKeyResponse{}, fmt.Errorf("invalid plugin metadata: %w", err)
	}
	provider, close, err := start(ctx, name)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}
	defer close()

	response, err := provider.Unwrap(ctx, request)
	if err != nil {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin %q failed to unwrap the data key: %w", name, err)
	}
	if len(response.DataKey) == 0 {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin %q returned an empty data key", name)
	}
	return response, nil
}

func start(ctx context.Context, name string) (KeyProvider, func(), error) {
	if err := ValidatePluginName(name); err != nil {
		return nil, nil, err
	}
	path, err := resolvePlugin(name)
	if err != nil {
		return nil, nil, err
	}
	client := hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig: hplugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   magicCookieKey,
			MagicCookieValue: magicCookieValue,
		},
		Plugins: hplugin.PluginSet{
			pluginName: &grpcPlugin{},
		},
		Cmd:              exec.CommandContext(ctx, path),
		AllowedProtocols: []hplugin.Protocol{hplugin.ProtocolGRPC},
		StartTimeout:     10 * time.Second,
		Stderr:           os.Stderr,
		Logger:           hclog.NewNullLogger(),
	})
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, fmt.Errorf("start plugin %q: %w", name, err)
	}
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		client.Kill()
		return nil, nil, fmt.Errorf("connect to plugin %q: %w", name, err)
	}
	provider, ok := raw.(KeyProvider)
	if !ok {
		client.Kill()
		return nil, nil, fmt.Errorf("plugin %q does not implement the SOPS key-provider API", name)
	}
	return provider, client.Kill, nil
}

// ValidatePluginName rejects names that could be interpreted as executable
// paths or command arguments.
func ValidatePluginName(name string) error {
	if !pluginNamePattern.MatchString(name) {
		return fmt.Errorf("invalid plugin name %q", name)
	}
	return nil
}

func resolvePlugin(name string) (string, error) {
	configPath, err := pluginConfigPath()
	if err != nil {
		return "", err
	}
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("read plugin allowlist %q: %w", configPath, err)
	}
	var config struct {
		Plugins map[string]string `yaml:"plugins"`
	}
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return "", fmt.Errorf("parse plugin allowlist %q: %w", configPath, err)
	}
	path, ok := config.Plugins[name]
	if !ok {
		return "", fmt.Errorf("plugin %q is not in the local allowlist %q", name, configPath)
	}
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("plugin %q executable path must be absolute", name)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat plugin %q executable %q: %w", name, path, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("plugin %q executable %q is not a regular file", name, path)
	}
	return path, nil
}

func pluginConfigPath() (string, error) {
	if path := os.Getenv("SOPS_PLUGIN_CONFIG"); path != "" {
		return path, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find user config directory for plugin allowlist: %w", err)
	}
	return filepath.Join(configDir, "sops", "plugins.yaml"), nil
}

func encodeObject(value map[string]interface{}) ([]byte, error) {
	if value == nil {
		value = map[string]interface{}{}
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var object map[string]interface{}
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	if object == nil {
		return nil, fmt.Errorf("expected a JSON object")
	}
	return data, nil
}

func decodeObject(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return map[string]interface{}{}, nil
	}
	var object map[string]interface{}
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	if object == nil {
		return nil, fmt.Errorf("expected a JSON object")
	}
	return object, nil
}
