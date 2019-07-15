---
layout: "docs"
page_title: "Plugin System"
sidebar_title: "Plugins"
sidebar_current: "docs-internals-plugins"
description: |-
  Learn about Vault's plugin system.
---

# Plugin System
All Vault auth and secret backends are considered plugins. This simple concept
allows both built-in and external plugins to be treated like Legos. Any plugin
can exist at multiple different locations. Different versions of a plugin may
be at each one, with each version differing from Vault's version.

## Built-In Plugins
Built-in plugins are shipped with Vault, often for commonly used implementations,
and require no additional operator intervention to run. Built-in plugins are
just like any other backend code inside Vault.

To use a different or edited version of a built-in plugin, you would first edit
the plugin's code or navigate to the Vault version holding the version of the
plugin you desire. Then, you'd `$ cd` into the `cmd/:plugin-name` directory
contained alongside that plugin's code. For instance, for AppRole, you would:
`$ cd vault/builtin/credential/approle/cmd/approle`. Once in that directory,
you would run `$ go build` to obtain a new binary for the AppRole plugin. Then
you would add it to the plugin catalog as per normal, and enable it.

# Plugin Architecture
Vault's plugins are completely separate, standalone applications that Vault
executes and communicates with over RPC. This means the plugin process does not
share the same memory space as Vault and therefore can only access the
interfaces and arguments given to it. This also means a crash in a plugin can not
crash the entirety of Vault.

It is possible to enable a custom plugin with a name that's identical to a 
built-in plugin. In such a situation, Vault will always choose the custom plugin
when enabling it. 

## Plugin Communication
Vault creates a mutually authenticated TLS connection for communication with the
plugin's RPC server. While invoking the plugin process, Vault passes a [wrapping
token](https://www.vaultproject.io/docs/concepts/response-wrapping.html) to the
plugin process' environment. This token is single use and has a short TTL. Once
unwrapped, it provides the plugin with a uniquely generated TLS certificate and
private key for it to use to talk to the original Vault process.

The [`api_addr`][api_addr] must be set in order for the plugin process establish
communication with the Vault server during mount time. If the storage backend
has HA enabled and supports automatic host address detection (e.g. Consul),
Vault will automatically attempt to determine the `api_addr` as well.

~> Note: Reading the original connection's TLS connection state is not supported
in plugins.

## Plugin Registration
An important consideration of Vault's plugin system is to ensure the plugin
invoked by Vault is authentic and maintains integrity. There are two components
that a Vault operator needs to configure before external plugins can be run, the
plugin directory and the plugin catalog entry.

### Plugin Directory
The plugin directory is a configuration option of Vault, and can be specified in
the [configuration file](https://www.vaultproject.io/docs/configuration/index.html).
This setting specifies a directory that all plugin binaries must live. A plugin
can not be added to Vault unless it exists in the plugin directory. There is no
default for this configuration option, and if it is not set plugins can not be
added to Vault.

~> Warning: A Vault operator should take care to lock down the permissions on
this directory to ensure a plugin can not be modified by an unauthorized user
between the time of the SHA check and the time of plugin execution.

### Plugin Catalog
The plugin catalog is Vault's list of approved plugins. The catalog is stored in
Vault's barrier and can only be updated by a Vault user with sudo permissions.
Upon adding a new plugin, the plugin name, SHA256 sum of the executable, and the
command that should be used to run the plugin must be provided. The catalog will
make sure the executable referenced in the command exists in the plugin
directory. When added to the catalog the plugin is not automatically executed,
it instead becomes visible to backends and can be executed by them. For more
information on the plugin catalog please see the [Plugin Catalog API
docs](/api/system/plugins-catalog.html).

An example plugin submission looks like:

```
$ vault write sys/plugins/catalog/database/myplugin-database-plugin \
    sha256=<expected SHA256 Hex value of the plugin binary> \
    command="myplugin"
Success! Data written to: sys/plugins/catalog/database/myplugin-database-plugin
```

### Plugin Execution
When a backend wants to run a plugin, it first looks up the plugin, by name, in
the catalog. It then checks the executable's SHA256 sum against the one
configured in the plugin catalog. Finally Vault runs the command configured in
the catalog, sending along the JWT formatted response wrapping token and mlock
settings (like Vault, plugins support [the use of mlock when available](https://www.vaultproject.io/docs/configuration/index.html#disable_mlock)).

# Plugin Development

~> Advanced topic! Plugin development is a highly advanced topic in Vault, and
is not required knowledge for day-to-day usage. If you don't plan on writing any
plugins, we recommend not reading this section of the documentation.

Because Vault communicates to plugins over a RPC interface, you can build and
distribute a plugin for Vault without having to rebuild Vault itself. This makes
it easy for you to build a Vault plugin for your organization's internal use,
for a proprietary API that you don't want to open source, or to prototype
something before contributing it back to the main project.

In theory, because the plugin interface is HTTP, you could even develop a plugin
using a completely different programming language! (Disclaimer, you would also
have to re-implement the plugin API which is not a trivial amount of work.)

Developing a plugin is simple. The only knowledge necessary to write
a plugin is basic command-line skills and basic knowledge of the
[Go programming language](http://golang.org).

Your plugin implementation needs to satisfy the interface for the plugin
type you want to build. You can find these definitions in the docs for the
backend running the plugin.

```go
package main

import (
	"os"

	myPlugin "your/plugin/import/path"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: myPlugin.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
```

And that's basically it! You would just need to change `myPlugin` to your actual
plugin. For more information on how to register and enable your plugin, check out the [Building Plugin Backends](https://learn.hashicorp.com/vault/developer/plugin-backends) tutorial.

[api_addr]: /docs/configuration/index.html#api_addr
