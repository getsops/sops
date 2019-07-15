---
layout: "docs"
page_title: "kv undelete - Command"
sidebar_title: "<code>undelete</code>"
sidebar_current: "docs-commands-kv-undelete"
description: |-
  The "kv undelete" command undeletes the data for the provided version and path
  in the key-value store. This restores the data, allowing it to be returned on
  get requests.
---

# kv undelete

~> **NOTE:** This is a [K/V Version 2](/docs/secrets/kv/kv-v2.html) secrets
engine command, and not available for Version 1.


The `kv undelete` command undoes the deletes of the data for the provided version
and path in the key-value store. This restores the data, allowing it to be
returned on get requests.

## Examples

Undelete version 3 of the key "creds":

```text
$ vault kv undelete -versions=3 secret/creds
Success! Data written to: secret/undelete/creds
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-versions` `([]int: <required>)` - Specifies the version number that should
be made current again.
