// Package plugin provides the Go SDK and process runtime for SOPS key-provider
// plugins.
package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/shamir"
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

// WrapKeyRequest contains a SOPS data key or Shamir share, public
// configuration, and metadata from any previous wrapping operation.
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

// Runtime manages plugin processes and reuses them until Close is called.
// Provider methods may be invoked concurrently and must be safe for concurrent
// use. Runtime does not persist provider state between SOPS operations; use
// request metadata or an external store for per-key state.
type Runtime struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu      sync.Mutex
	clients map[string]runtimePlugin
	closed  bool
}

type runtimePlugin struct {
	client            *hplugin.Client
	provider          KeyProvider
	invocationTimeout time.Duration
}

type pluginConfig struct {
	path              string
	invocationTimeout time.Duration
}

// NewRuntime creates a plugin runtime whose child processes live until the
// runtime is closed or parent is canceled.
func NewRuntime(parent context.Context) *Runtime {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &Runtime{
		ctx:     ctx,
		cancel:  cancel,
		clients: make(map[string]runtimePlugin),
	}
}

// Close stops every plugin process owned by the runtime. It is safe to call
// Close more than once.
func (runtime *Runtime) Close() error {
	if runtime == nil {
		return nil
	}
	runtime.mu.Lock()
	if runtime.closed {
		runtime.mu.Unlock()
		return nil
	}
	runtime.closed = true
	clients := runtime.clients
	runtime.clients = nil
	runtime.mu.Unlock()

	runtime.cancel()
	for _, plugin := range clients {
		plugin.client.Kill()
	}
	return nil
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
	if ctx == nil {
		return WrapKeyResponse{}, fmt.Errorf("plugin request context is nil")
	}
	if err := validateWrapRequest(request); err != nil {
		return WrapKeyResponse{}, err
	}
	plugin, err := start(ctx, name)
	if err != nil {
		return WrapKeyResponse{}, err
	}
	defer plugin.client.Kill()

	callContext, cancel := withInvocationTimeout(ctx, plugin.invocationTimeout)
	defer cancel()
	return wrapWithProvider(callContext, name, plugin.provider, request)
}

// Runtime.WrapKey starts or reuses an allowlisted plugin and asks it to wrap a SOPS data key.
func (runtime *Runtime) WrapKey(ctx context.Context, name string, request WrapKeyRequest) (WrapKeyResponse, error) {
	if ctx == nil {
		return WrapKeyResponse{}, fmt.Errorf("plugin request context is nil")
	}
	if err := validateWrapRequest(request); err != nil {
		return WrapKeyResponse{}, err
	}
	plugin, err := runtime.provider(ctx, name)
	if err != nil {
		return WrapKeyResponse{}, err
	}
	callContext, cancel := withInvocationTimeout(ctx, plugin.invocationTimeout)
	defer cancel()
	response, err := wrapWithProvider(callContext, name, plugin.provider, request)
	if err != nil && invocationDeadlineExceeded(callContext) {
		runtime.evict(name, plugin.client)
	}
	return response, err
}

func validateWrapRequest(request WrapKeyRequest) error {
	if len(request.DataKey) == 0 {
		return fmt.Errorf("cannot wrap an empty data key")
	}
	if !validKeyMaterialSize(len(request.DataKey)) {
		return fmt.Errorf("SOPS data key or Shamir share must be %d or %d bytes, got %d", keys.DataKeySize, keys.DataKeySize+shamir.ShareOverhead, len(request.DataKey))
	}
	if _, err := encodeObject(request.Config); err != nil {
		return fmt.Errorf("invalid plugin configuration: %w", err)
	}
	if _, err := encodeObject(request.Metadata); err != nil {
		return fmt.Errorf("invalid plugin metadata: %w", err)
	}
	return nil
}

func wrapWithProvider(ctx context.Context, name string, provider KeyProvider, request WrapKeyRequest) (WrapKeyResponse, error) {
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
	if ctx == nil {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin request context is nil")
	}
	if err := validateUnwrapRequest(request); err != nil {
		return UnwrapKeyResponse{}, err
	}
	plugin, err := start(ctx, name)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}
	defer plugin.client.Kill()

	callContext, cancel := withInvocationTimeout(ctx, plugin.invocationTimeout)
	defer cancel()
	return unwrapWithProvider(callContext, name, plugin.provider, request)
}

// Runtime.UnwrapKey starts or reuses an allowlisted plugin and asks it to unwrap a SOPS data key.
func (runtime *Runtime) UnwrapKey(ctx context.Context, name string, request UnwrapKeyRequest) (UnwrapKeyResponse, error) {
	if ctx == nil {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin request context is nil")
	}
	if err := validateUnwrapRequest(request); err != nil {
		return UnwrapKeyResponse{}, err
	}
	plugin, err := runtime.provider(ctx, name)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}
	callContext, cancel := withInvocationTimeout(ctx, plugin.invocationTimeout)
	defer cancel()
	response, err := unwrapWithProvider(callContext, name, plugin.provider, request)
	if err != nil && invocationDeadlineExceeded(callContext) {
		runtime.evict(name, plugin.client)
	}
	return response, err
}

func validateUnwrapRequest(request UnwrapKeyRequest) error {
	if len(request.EncryptedDataKey) == 0 {
		return fmt.Errorf("cannot unwrap an empty data key")
	}
	if _, err := encodeObject(request.Config); err != nil {
		return fmt.Errorf("invalid plugin configuration: %w", err)
	}
	if _, err := encodeObject(request.Metadata); err != nil {
		return fmt.Errorf("invalid plugin metadata: %w", err)
	}
	return nil
}

func unwrapWithProvider(ctx context.Context, name string, provider KeyProvider, request UnwrapKeyRequest) (UnwrapKeyResponse, error) {
	response, err := provider.Unwrap(ctx, request)
	if err != nil {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin %q failed to unwrap the data key: %w", name, err)
	}
	if len(response.DataKey) == 0 {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin %q returned an empty data key", name)
	}
	if !validKeyMaterialSize(len(response.DataKey)) {
		return UnwrapKeyResponse{}, fmt.Errorf("plugin %q returned key material with invalid size: expected %d or %d bytes, got %d", name, keys.DataKeySize, keys.DataKeySize+shamir.ShareOverhead, len(response.DataKey))
	}
	return response, nil
}

func validKeyMaterialSize(size int) bool {
	return size == keys.DataKeySize || size == keys.DataKeySize+shamir.ShareOverhead
}

func (runtime *Runtime) provider(requestContext context.Context, name string) (runtimePlugin, error) {
	if runtime == nil {
		return runtimePlugin{}, fmt.Errorf("plugin runtime is nil")
	}
	if err := ValidatePluginName(name); err != nil {
		return runtimePlugin{}, err
	}
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	if runtime.closed {
		return runtimePlugin{}, fmt.Errorf("plugin runtime is closed")
	}
	if err := runtime.ctx.Err(); err != nil {
		return runtimePlugin{}, fmt.Errorf("plugin runtime is canceled: %w", err)
	}
	if existing, ok := runtime.clients[name]; ok {
		if !existing.client.Exited() {
			return existing, nil
		}
		existing.client.Kill()
		delete(runtime.clients, name)
	}
	plugin, err := startForRuntime(runtime.ctx, requestContext, name)
	if err != nil {
		return runtimePlugin{}, err
	}
	runtime.clients[name] = plugin
	return plugin, nil
}

func (runtime *Runtime) evict(name string, client *hplugin.Client) {
	runtime.mu.Lock()
	plugin, ok := runtime.clients[name]
	if ok && plugin.client == client {
		delete(runtime.clients, name)
	}
	runtime.mu.Unlock()
	if ok && plugin.client == client {
		client.Kill()
	}
}

func startForRuntime(lifetimeContext, requestContext context.Context, name string) (runtimePlugin, error) {
	startupContext, cancelStartup := context.WithCancel(lifetimeContext)
	callbackDone := make(chan struct{})
	stopRequestCancellation := context.AfterFunc(requestContext, func() {
		cancelStartup()
		close(callbackDone)
	})
	plugin, err := start(startupContext, name)
	if !stopRequestCancellation() {
		<-callbackDone
	}
	if err != nil {
		cancelStartup()
		return runtimePlugin{}, err
	}
	if err := requestContext.Err(); err != nil {
		plugin.client.Kill()
		cancelStartup()
		return runtimePlugin{}, err
	}
	if err := lifetimeContext.Err(); err != nil {
		plugin.client.Kill()
		cancelStartup()
		return runtimePlugin{}, fmt.Errorf("plugin runtime is canceled: %w", err)
	}
	return plugin, nil
}

func start(ctx context.Context, name string) (runtimePlugin, error) {
	if err := ValidatePluginName(name); err != nil {
		return runtimePlugin{}, err
	}
	config, err := resolvePlugin(name)
	if err != nil {
		return runtimePlugin{}, err
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
		Cmd:              exec.CommandContext(ctx, config.path),
		AllowedProtocols: []hplugin.Protocol{hplugin.ProtocolGRPC},
		StartTimeout:     10 * time.Second,
		Stderr:           os.Stderr,
		Logger:           hclog.NewNullLogger(),
	})
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return runtimePlugin{}, fmt.Errorf("start plugin %q: %w", name, err)
	}
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		client.Kill()
		return runtimePlugin{}, fmt.Errorf("connect to plugin %q: %w", name, err)
	}
	provider, ok := raw.(KeyProvider)
	if !ok {
		client.Kill()
		return runtimePlugin{}, fmt.Errorf("plugin %q does not implement the SOPS key-provider API", name)
	}
	return runtimePlugin{
		client:            client,
		provider:          provider,
		invocationTimeout: config.invocationTimeout,
	}, nil
}

// ValidatePluginName rejects names that could be interpreted as executable
// paths or command arguments.
func ValidatePluginName(name string) error {
	if !pluginNamePattern.MatchString(name) {
		return fmt.Errorf("invalid plugin name %q", name)
	}
	return nil
}

func resolvePlugin(name string) (pluginConfig, error) {
	configPath, err := pluginConfigPath()
	if err != nil {
		return pluginConfig{}, err
	}
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return pluginConfig{}, fmt.Errorf("read plugin allowlist %q: %w", configPath, err)
	}
	var config struct {
		Plugins map[string]yaml.Node `yaml:"plugins"`
	}
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return pluginConfig{}, fmt.Errorf("parse plugin allowlist %q: %w", configPath, err)
	}
	entry, ok := config.Plugins[name]
	if !ok {
		return pluginConfig{}, fmt.Errorf("plugin %q is not in the local allowlist %q", name, configPath)
	}
	var configured struct {
		Path              string `yaml:"path"`
		InvocationTimeout string `yaml:"invocation_timeout"`
	}
	if entry.Kind == yaml.ScalarNode {
		if err := entry.Decode(&configured.Path); err != nil {
			return pluginConfig{}, fmt.Errorf("parse plugin %q allowlist entry: %w", name, err)
		}
	} else if entry.Kind == yaml.MappingNode {
		if err := entry.Decode(&configured); err != nil {
			return pluginConfig{}, fmt.Errorf("parse plugin %q allowlist entry: %w", name, err)
		}
	} else {
		return pluginConfig{}, fmt.Errorf("plugin %q allowlist entry must be a path or mapping", name)
	}
	if configured.Path == "" {
		return pluginConfig{}, fmt.Errorf("plugin %q executable path is empty", name)
	}
	if !filepath.IsAbs(configured.Path) {
		return pluginConfig{}, fmt.Errorf("plugin %q executable path must be absolute", name)
	}
	invocationTimeout := time.Duration(0)
	if configured.InvocationTimeout != "" {
		invocationTimeout, err = time.ParseDuration(configured.InvocationTimeout)
		if err != nil {
			return pluginConfig{}, fmt.Errorf("parse plugin %q invocation_timeout %q: %w", name, configured.InvocationTimeout, err)
		}
		if invocationTimeout <= 0 {
			return pluginConfig{}, fmt.Errorf("plugin %q invocation_timeout must be a positive Go duration", name)
		}
	}
	info, err := os.Stat(configured.Path)
	if err != nil {
		return pluginConfig{}, fmt.Errorf("stat plugin %q executable %q: %w", name, configured.Path, err)
	}
	if !info.Mode().IsRegular() {
		return pluginConfig{}, fmt.Errorf("plugin %q executable %q is not a regular file", name, configured.Path)
	}
	return pluginConfig{
		path:              configured.Path,
		invocationTimeout: invocationTimeout,
	}, nil
}

func withInvocationTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}

func invocationDeadlineExceeded(ctx context.Context) bool {
	if err := ctx.Err(); err != nil {
		return err == context.DeadlineExceeded
	}
	deadline, ok := ctx.Deadline()
	return ok && !time.Now().Before(deadline)
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

// DecodeConfig decodes plugin configuration into a provider-defined value.
// The destination must be a pointer. Numeric fields should use concrete types
// in the destination struct to avoid interface-number conversions.
func DecodeConfig(config map[string]interface{}, destination interface{}) error {
	return decodeInto(config, destination, "plugin configuration")
}

// DecodeMetadata decodes per-key provider metadata into a provider-defined
// value. The destination must be a pointer.
func DecodeMetadata(metadata map[string]interface{}, destination interface{}) error {
	return decodeInto(metadata, destination, "plugin metadata")
}

func decodeInto(value map[string]interface{}, destination interface{}, description string) error {
	data, err := encodeObject(value)
	if err != nil {
		return fmt.Errorf("encode %s: %w", description, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode %s: %w", description, err)
	}
	return nil
}

func decodeObject(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return map[string]interface{}{}, nil
	}
	var object map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decodeJSONValue(decoder, &object); err != nil {
		return nil, err
	}
	if object == nil {
		return nil, fmt.Errorf("expected a JSON object")
	}
	return object, nil
}

func decodeJSONValue(decoder *json.Decoder, destination interface{}) error {
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing interface{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}
