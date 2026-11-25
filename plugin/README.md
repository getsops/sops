# External key-provider plugins

SOPS can delegate wrapping and unwrapping its data key to an external Go
process. SOPS continues to own file encryption, key groups, and Shamir secret
sharing; a provider only implements the key operation it needs. This supports
providers that wrap offline, use remote KMS APIs, or combine secret retrieval
with local cryptography.

## Implement a provider

Implement these methods on a provider and start its process with
`plugin.Serve(provider)`:

```go
import (
    "context"

    sopsplugin "github.com/getsops/sops/v3/plugin"
)

type KeyProvider interface {
    Wrap(context.Context, sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error)
    Unwrap(context.Context, sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error)
}
```

`WrapKeyResponse.EncryptedDataKey` is opaque to SOPS. `Wrap` receives either
the 32-byte SOPS data key or a 33-byte Shamir share when SOPS uses multiple key
groups. `Unwrap` must return the same key material it was given. Optional `Metadata` is
stored alongside it and passed to later `Wrap` and `Unwrap` calls. The first
`Wrap` receives an empty metadata object; each successful `Wrap` replaces the
stored metadata. This lets a provider reuse public, per-key material when
rewrapping a data key.

The host's key-service runtime can reuse a plugin process across `Wrap` and
`Unwrap` calls. Providers must be concurrency-safe and must not rely on
process-local state for per-key data; use per-key metadata or an external state
store. The config and metadata values are JSON objects. SOPS passes plugin
config through unchanged and does not resolve paths inside it; relative paths
are interpreted by the plugin relative to its working directory. The plugin
handshake uses stdout, so send provider logs to stderr. Implementations should
honor request context cancellation.

The gRPC service is versioned as `sops.plugin.v1`. Evolve its protobuf messages
additively; breaking protocol changes should use a new service version.

Provider configuration and metadata are persisted in encrypted files and may
outlive several plugin releases. Keep old serialized forms readable when
changing your provider, and version provider-owned metadata when its shape or
meaning changes. A plugin upgrade must continue to unwrap data keys produced
by earlier compatible versions. `plugin.DecodeConfig` and
`plugin.DecodeMetadata` can decode the generic SDK maps into provider-defined
structs:

```go
type Config struct {
    KeyID string `json:"key_id"`
}

var config Config
if err := sopsplugin.DecodeConfig(request.Config, &config); err != nil {
    return sopsplugin.WrapKeyResponse{}, err
}
```

The generic config and metadata maps represent JSON numbers as `json.Number`;
typed decoding into Go integer or floating-point fields avoids lossy
`float64` conversions.

`sopsplugin.NewRuntime` is a host-side API used by SOPS's existing key-service
implementation; it is not a new server API for plugins to call. It reuses
plugin processes and their connections across requests until the owning
key-service shuts down. The normal CLI owns an in-process local key-service;
the optional key-service daemon can own one for its lifetime as well. Calls to
one plugin process may run concurrently, so provider implementations must be
concurrency-safe; do not rely on process memory for per-key state.

The plugin receives the plaintext SOPS data key and runs with the SOPS
process's user permissions; plugins are trusted code, not a security sandbox.

## Configure a plugin

Creation rules and key groups can include plugin keys:

```yaml
creation_rules:
  - path_regex: secrets\.yaml$
    plugins:
      - name: org-vault
        config:
          key_id: projects/example/keys/sops
```

The same object can be supplied as a JSON object or array to `--plugin`,
`--add-plugin`, or `--rm-plugin`.

SOPS only launches plugins listed in the user's local allowlist. By default,
the allowlist is `$XDG_CONFIG_HOME/sops/plugins.yaml` (or the platform's user
configuration directory):

```yaml
plugins:
  org-vault:
    path: /usr/local/lib/sops/sops-plugin-org-vault
    invocation_timeout: 45s
```

Executable paths must be absolute. Set `SOPS_PLUGIN_CONFIG` to use a different
allowlist file. The allowlist is local policy and is not read from the
repository's `.sops.yaml` or the encrypted file.

The legacy `org-vault: /absolute/path` form remains supported. The optional
`invocation_timeout` is a positive Go duration applied independently to each
`Wrap` and `Unwrap`; if omitted, there is no host-imposed invocation deadline.
An earlier caller deadline still takes precedence. Plugin startup is separately
bounded to 10 seconds. When an invocation deadline expires, SOPS returns the
timeout error and kills/evicts that plugin process so a stuck handler is not
reused. This can also interrupt other in-flight calls handled by that process;
SOPS does not retry timed-out operations because a provider may already have
performed side effects. An explicit caller cancellation does not evict a
healthy process.

Plugin config and provider metadata are stored in the SOPS file's plaintext
metadata. Do not put credentials or other secrets in either value; load secrets
from the plugin's normal runtime environment or credential sources instead.
