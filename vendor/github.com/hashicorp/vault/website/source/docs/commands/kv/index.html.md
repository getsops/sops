---
layout: "docs"
page_title: "kv - Command"
sidebar_title: "<code>kv</code>"
sidebar_current: "docs-commands-kv"
description: |-
  The "kv" command groups subcommands for interacting with Vault's key/value
  secret engine.
---

# kv

The `kv` command groups subcommands for interacting with Vault's key/value
secrets engine (both [K/V Version 1](/docs/secrets/kv/kv-v1.html) and [K/V
Version 2](/docs/secrets/kv/kv-v2.html).


## Examples

Create or update the key named "creds" in the K/V Version 2 enabled at "secret"
with the value "passcode=my-long-passcode":

```text
$ vault kv put secret/creds passcode=my-long-passcode
Key              Value
---              -----
created_time     2019-06-28T15:53:30.395814Z
deletion_time    n/a
destroyed        false
version          1
```

Read this value back:

```text
$ vault kv get secret/creds
====== Metadata ======
Key              Value
---              -----
created_time     2019-06-28T15:53:30.395814Z
deletion_time    n/a
destroyed        false
version          1

====== Data ======
Key         Value
---         -----
passcode    my-long-passcode
```

Get metadata for the key named "creds":

```text
$ vault kv metadata get secret/creds
========== Metadata ==========
Key                     Value
---                     -----
cas_required            false
created_time            2019-06-28T15:53:30.395814Z
current_version         1
delete_version_after    0s
max_versions            0
oldest_version          0
updated_time            2019-06-28T15:53:30.395814Z

====== Version 1 ======
Key              Value
---              -----
created_time     2019-06-28T15:53:30.395814Z
deletion_time    n/a
destroyed        false
```


Get a specific version of the key named "creds":

```text
$ vault kv get -version=1 secret/creds
====== Metadata ======
Key              Value
---              -----
created_time     2019-06-28T15:53:30.395814Z
deletion_time    n/a
destroyed        false
version          1

====== Data ======
Key         Value
---         -----
passcode    my-long-passcode
```


## Usage

```text
Usage: vault kv <subcommand> [options] [args]

  # ...

Subcommands:
    delete               Deletes versions in the KV store
    destroy              Permanently removes one or more versions in the KV store
    enable-versioning    Turns on versioning for a KV store
    get                  Retrieves data from the KV store
    list                 List data or secrets
    metadata             Interact with Vault's Key-Value storage
    patch                Sets or updates data in the KV store without overwriting
    put                  Sets or updates data in the KV store
    rollback             Rolls back to a previous version of data
    undelete             Undeletes versions in the KV store
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
