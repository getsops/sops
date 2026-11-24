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

`WrapKeyResponse.EncryptedDataKey` is opaque to SOPS. Optional `Metadata` is
stored alongside it and passed back to `Unwrap`. The config and metadata values
are JSON objects. The plugin receives the plaintext SOPS data key and runs with
the SOPS process's user permissions; plugins are trusted code, not a security
sandbox.

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
  org-vault: /usr/local/lib/sops/sops-plugin-org-vault
```

Executable paths must be absolute. Set `SOPS_PLUGIN_CONFIG` to use a different
allowlist file. The allowlist is local policy and is not read from the
repository's `.sops.yaml` or the encrypted file.

Plugin config and provider metadata are stored in the SOPS file's plaintext
metadata. Do not put credentials or other secrets in either value; load secrets
from the plugin's normal runtime environment or credential sources instead.
