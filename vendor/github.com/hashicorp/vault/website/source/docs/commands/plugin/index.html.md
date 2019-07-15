---
layout: "docs"
page_title: "plugin - Command"
sidebar_title: "<code>plugin</code>"
sidebar_current: "docs-commands-plugin"
description: |-
  The "plugin" command groups subcommands for interacting with
  Vault's plugins and the plugin catalog.
---

# plugin

The `plugin` command groups subcommands for interacting with Vault's plugins and
the plugin catalog

## Examples

List all available plugins in the catalog:

```text
$ vault plugin list

Plugins
-------
my-custom-plugin
# ...
```

Register a new plugin to the catalog:

```text
$ vault plugin register \
  -sha256=d3f0a8be02f6c074cf38c9c99d4d04c9c6466249 \
  my-custom-plugin
Success! Registered plugin: my-custom-plugin
```

Get information about a plugin in the catalog:

```text
$ vault plugin info my-custom-plugin
Key        Value
---        -----
command    my-custom-plugin
name       my-custom-plugin
sha256     d3f0a8be02f6c074cf38c9c99d4d04c9c6466249
```

## Usage

```text
Usage: vault plugin <subcommand> [options] [args]

  # ...

Subcommands:
    deregister    Deregister an existing plugin in the catalog
    list          Lists available plugins
    read          Read information about a plugin in the catalog
    register      Registers a new plugin in the catalog
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
