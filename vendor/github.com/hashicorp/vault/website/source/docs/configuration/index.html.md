---
layout: "docs"
page_title: "Server Configuration"
sidebar_title: "Configuration"
sidebar_current: "docs-configuration"
description: |-
  Vault server configuration reference.
---

# Vault Configuration

Outside of development mode, Vault servers are configured using a file.
The format of this file is [HCL](https://github.com/hashicorp/hcl) or JSON.
An example configuration is shown below:

```javascript
storage "consul" {
  address = "127.0.0.1:8500"
  path    = "vault"
}

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = 1
}

telemetry {
  statsite_address = "127.0.0.1:8125"
  disable_hostname = true
}
```

After the configuration is written, use the `-config` flag with `vault server`
to specify where the configuration is.

## Parameters

- `storage` <tt>([StorageBackend][storage-backend]: \<required\>)</tt> –
  Configures the storage backend where Vault data is stored. Please see the
  [storage backends documentation][storage-backend] for the full list of
  available storage backends. Running Vault in HA mode would require
  coordination semantics to be supported by the backend. If the storage backend
  supports HA coordination, HA backend options can also be specified in this
  parameter block. If not, a separate `ha_storage` parameter should be
  configured with a backend that supports HA, along with corresponding HA
  options.

- `ha_storage` <tt>([StorageBackend][storage-backend]: nil)</tt> – Configures
  the storage backend where Vault HA coordination will take place. This must be
  an HA-supporting backend. If not set, HA will be attempted on the backend
  given in the `storage` parameter. This parameter is not required if the
  storage backend supports HA coordination and if HA specific options are
  already specified with `storage` parameter.

- `listener` <tt>([Listener][listener]: \<required\>)</tt> – Configures how
  Vault is listening for API requests.

- `seal` <tt>([Seal][seal]: nil)</tt> – Configures the seal type to use for
  auto-unsealing, as well as for
  [seal wrapping][sealwrap] as an additional layer of data protection.

- `cluster_name` `(string: <generated>)` – Specifies the identifier for the
  Vault cluster. If omitted, Vault will generate a value. When connecting to
  Vault Enterprise, this value will be used in the interface.

- `cache_size` `(string: "32000")` – Specifies the size of the read cache used
  by the physical storage subsystem. The value is in number of entries, so the
  total cache size depends on the size of stored entries.

- `disable_cache` `(bool: false)` – Disables all caches within Vault, including
  the read cache used by the physical storage subsystem. This will very
  significantly impact performance.

- `disable_mlock` `(bool: false)` – Disables the server from executing the
  `mlock` syscall. `mlock` prevents memory from being swapped to disk. Disabling
  `mlock` is not recommended in production, but is fine for local development
  and testing.

    Disabling `mlock` is not recommended unless the systems running Vault only
    use encrypted swap or do not use swap at all. Vault only supports memory
    locking on UNIX-like systems that support the mlock() syscall (Linux, FreeBSD, etc).
    Non UNIX-like systems (e.g. Windows, NaCL, Android) lack the primitives to keep a
    process's entire memory address space from spilling to disk and is therefore
    automatically disabled on unsupported platforms.

    On Linux, to give the Vault executable the ability to use the `mlock`
    syscall without running the process as root, run:

    ```shell
    sudo setcap cap_ipc_lock=+ep $(readlink -f $(which vault))
    ```

    Note: Since each plugin runs as a separate process, you need to do the same for each plugin in your [plugins directory](https://www.vaultproject.io/docs/internals/plugins.html#plugin-directory).

    If you use a Linux distribution with a modern version of systemd, you can add
    the following directive to the "[Service]" configuration section:

    ```ini
    LimitMEMLOCK=infinity
    ```

- `plugin_directory` `(string: "")` – A directory from which plugins are
  allowed to be loaded. Vault must have permission to read files in this
  directory to successfully load plugins.

- `telemetry` <tt>([Telemetry][telemetry]: &lt;none&gt;)</tt> – Specifies the telemetry
  reporting system.

- `log_level` `(string: "")` – Specifies the log level to use; overridden by
  CLI and env var parameters. On SIGHUP, Vault will update the log level to the
  current value specified here (including overriding the CLI/env var
  parameters). Not all parts of Vault's logging can have its level be changed
  dynamically this way; in particular, secrets/auth plugins are currently not
  updated dynamically. Supported log levels: Trace, Debug, Error, Warn, Info.

- `default_lease_ttl` `(string: "768h")` – Specifies the default lease duration
  for tokens and secrets. This is specified using a label suffix like `"30s"` or
  `"1h"`. This value cannot be larger than `max_lease_ttl`.

- `max_lease_ttl` `(string: "768h")` – Specifies the maximum possible lease
  duration for tokens and secrets. This is specified using a label
  suffix like `"30s"` or `"1h"`.

- `default_max_request_duration` `(string: "90s")` – Specifies the default
  maximum request duration allowed before Vault cancels the request. This can
  be overridden per listener via the `max_request_duration` value.

- `raw_storage_endpoint` `(bool: false)` – Enables the `sys/raw` endpoint which
  allows the decryption/encryption of raw data into and out of the security
  barrier. This is a highly privileged endpoint.

- `ui` `(bool: false)` – Enables the built-in web UI, which is available on all
  listeners (address + port) at the `/ui` path. Browsers accessing the standard
  Vault API address will automatically redirect there. This can also be provided
  via the environment variable `VAULT_UI`. For more information, please see the
  [ui configuration documentation](/docs/configuration/ui/index.html).

- `pid_file` `(string: "")` - Path to the file in which the Vault server's
  Process ID (PID) should be stored.

### High Availability Parameters

The following parameters are used on backends that support [high availability][high-availability].

- `api_addr` `(string: "")` - Specifies the address (full URL) to advertise to
  other Vault servers in the cluster for client redirection. This value is also
  used for [plugin backends][plugins]. This can also be provided via the
  environment variable `VAULT_API_ADDR`. In general this should be set as a full
  URL that points to the value of the [`listener`](#listener) address.

- `cluster_addr` `(string: "")` -  – Specifies the address to advertise to other
  Vault servers in the cluster for request forwarding. This can also be provided
  via the environment variable `VAULT_CLUSTER_ADDR`. This is a full URL, like
  `api_addr`, but Vault will ignore the scheme (all cluster members always
  use TLS with a private key/certificate).

- `disable_clustering` `(bool: false)` – Specifies whether clustering features
  such as request forwarding are enabled. Setting this to true on one Vault node
  will disable these features _only when that node is the active node_.

### Vault Enterprise Parameters

The following parameters are only used with Vault Enterprise

- `disable_sealwrap` `(bool: false)` – Disables using [seal wrapping][sealwrap]
  for any value except the master key. If this value is toggled, the new
  behavior will happen lazily (as values are read or written).

- `disable_performance_standby` `(bool: false)` – Specifies whether performance
  standbys should be disabled on this node. Setting this to true on one Vault
  node will disable this feature when this node is Active or Standby. It's
  recommended to sync this setting across all nodes in the cluster.

[storage-backend]: /docs/configuration/storage/index.html
[listener]: /docs/configuration/listener/index.html
[seal]: /docs/configuration/seal/index.html
[sealwrap]: /docs/enterprise/sealwrap/index.html
[telemetry]: /docs/configuration/telemetry.html
[high-availability]: /docs/concepts/ha.html
[plugins]: /docs/plugin/index.html
