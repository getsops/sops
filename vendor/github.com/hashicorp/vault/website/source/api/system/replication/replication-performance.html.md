---
layout: "api"
page_title: "/sys/replication - HTTP API"
sidebar_title: "<code>/sys/replication/performance</code>"
sidebar_current: "api-http-system-replication-performance"
description: |-
  The '/sys/replication/performance' endpoint focuses on managing general operations in Vault Enterprise Performance Replication
---

# `/sys/replication/performance`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

## Check Performance Status

This endpoint prints information about the status of replication (mode,
sync progress, etc).

This is an authenticated endpoint.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/performance/status`    |

### Sample Request

```
$ curl \
    http://127.0.0.1:8200/v1/sys/replication/performance/status
```

### Sample Response from Primary

The printed status of the replication environment. As an example, for a
primary, it will look something like:

```json
{
  "data": {
    "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
    "known_secondaries": [
      "2"
    ],
    "last_wal": 87,
    "merkle_root": "c31e40f6ff02f32c37b70e6a4d58732ac812abf0",
    "mode": "primary"
  },
}
```

### Sample Response from Secondary

The printed status of the replication environment. As an example, for a
secondary, it will look something like:

```json
{
  "data": {
    "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
    "known_primary_cluster_addrs": [
      "https://127.0.0.1:8201"
    ],
    "last_remote_wal": 87,
    "merkle_root": "c31e40f6ff02f32c37b70e6a4d58732ac812abf0",
    "mode": "secondary",
    "primary_cluster_addr": "https://127.0.0.1:8201",
    "secondary_id": "2",
    "state": "stream-wals"
  },
}
```

## Enable Performance Primary Replication

This endpoint enables replication in primary mode. This is used when replication
is currently disabled on the cluster (if the cluster is already a secondary, it
must be promoted).

!> Only one primary should be active at a given time. Multiple primaries may
result in data loss!

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/primary/enable` |

### Parameters

- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address,
  such as through a TCP-based load balancer. If not set, uses vault's configured
  cluster address.

### Sample Payload

```json
{}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/enable
```

## Demote Performance Primary

This endpoint demotes a performance primary cluster to a performance secondary.
This secondary cluster will not attempt to connect to a primary (see the update-primary call),
but will maintain knowledge of its cluster ID and can be reconnected to the same
replication set without wiping local storage.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/primary/demote` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/demote
```

## Disable Performance Primary

This endpoint disables performance replication entirely on the cluster. Any
performance secondaries will no longer be able to connect. Caution: re-enabling
this node as a primary or secondary will change its cluster ID; in the secondary
case this means a wipe of the underlying storage when connected to a primary,
and in the primary case, secondaries connecting back to the cluster (even if
they have connected before) will require a wipe of the underlying storage.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/primary/disable` |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/disable
```

## Generate Performance Secondary Token

This endpoint generates a performance secondary activation token for the
cluster with the given opaque identifier, which must be unique. This
identifier can later be used to revoke a secondary's access.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/sys/replication/performance/primary/secondary-token` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

- `ttl` `(string: "30m")` – Specifies the TTL for the secondary activation
  token.

### Sample Payload

```json
{
  "id": "us-east-1"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/secondary-token
```

### Sample Response

```json
{
  "request_id": "",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": null,
  "warnings": null,
  "wrap_info": {
    "token": "fb79b9d3-d94e-9eb6-4919-c559311133d6",
    "ttl": 300,
    "creation_time": "2016-09-28T14:41:00.56961496-04:00",
    "wrapped_accessor": ""
  }
}
```

## Revoke Performance Secondary Token

This endpoint revokes a performance secondary's ability to connect to the
performance primary cluster; the secondary will immediately be disconnected and
will not be allowed to connect again unless given a new activation token.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/primary/revoke-secondary` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

### Sample Payload

```json
{
  "id": "us-east"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/revoke-secondary
```

## Create Mounts Filter

This endpoint is used to modify the mounts that are filtered to a secondary.
Filtering can be specified in whitelist mode or blacklist mode.  In whitelist
mode the secret and auth mounts that are specified are included to the
selected secondary.  In blacklist mode, the mount paths are excluded.

| Method   | Path                                                     |
| :------------------------------------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/primary/mount-filter/:id` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

- `mode` `(string: "whitelist")` – Specifies the filtering mode.  Available values
  are "whitelist" and blacklist".

- `paths` `(array: [])` – The list of mount paths that are filtered.

### Sample Payload

```json
{
  "mode": "whitelist",
  "paths": ["secret/"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/mount-filter/us-east-1
```

## Read Mounts Filter

This endpoint is used to read the mode and the mount paths that are filtered
for a secondary.

| Method   | Path                                                     |
| :------------------------------------------------------- | :--------------------- |
| `GET`    | `/sys/replication/performance/primary/mount-filter/:id`  | `200 (empty body)` |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/mount-filter/us-east-1
```

### Sample Response

```json
{
  "mode": "whitelist",
  "paths": ["secret/"]
}
```

## Delete Mounts Filter

This endpoint is used to delete the mount filters for a secondary.

| Method   | Path                                                     |
| :------------------------------------------------------- | :--------------------- |
| `DELETE` | `/sys/replication/performance/primary/mount-filter/:id`  |

### Parameters

- `id` `(string: <required>)` – Specifies an opaque identifier, e.g. 'us-east'

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/replication/performance/primary/mount-filter/us-east-1
```

## Enable Performance Secondary

This endpoint enables performance replication on a secondary using a secondary activation
token.

!> This will immediately clear all data in the secondary cluster!

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/secondary/enable` |

### Parameters

- `token` `(string: <required>)` – Specifies the secondary activation token fetched from the primary.

- `primary_api_addr` `(string: "")` – Set this to the API address (normal Vault
  address) to override the value embedded in the token. This can be useful if
  the primary's redirect address is not accessible directly from this cluster
  (e.g. through a load balancer).

- `ca_file` `(string: "")` – Specifies the path to a CA root file (PEM format)
  that the secondary can use when unwrapping the token from the primary. If this
  and ca_path are not given, defaults to system CA roots.

- `ca_path` `(string: "")` – Specifies  the path to a CA root directory
  containing PEM-format files that the secondary can use when unwrapping the
  token from the primary. If this and ca_file are not given, defaults to system
  CA roots.

### Sample Payload

```json
{
  "token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/secondary/enable
```

## Promote Performance Secondary

This endpoint promotes the performance secondary cluster to performance primary.
For data safety and security reasons, new secondary tokens will need to be
issued to other secondaries, and there should never be more than one performance
primary at a time.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/secondary/promote` |

### Parameters

- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address
  (e.g. through a load balancer).
- `force` `(bool: false)` - If true the cluster will be promoted even if it fails 
  certain safety checks. Caution: Forcing promotion could result in data loss if 
  data isn't fully replicated.

### Sample Payload

```json
{}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/secondary/promote
```

## Disable Performance Secondary

This endpoint disables performance replication entirely on the cluster. The cluster will no
longer be able to connect to the performance primary.

!> Re-enabling this node as a performance primary or secondary will change its cluster ID;
in the secondary case this means a wipe of the underlying storage when connected
to a primary, and in the primary case, secondaries connecting back to the
cluster (even if they have connected before) will require a wipe of the
underlying storage.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/secondary/disable` |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/performance/secondary/disable
```

## Update Performance Secondary's Primary

This endpoint changes a performance secondary cluster's assigned primary cluster using a
secondary activation token. This does not wipe all data in the cluster.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/performance/secondary/update-primary` |

### Parameters

- `token` `(string: <required>)` – Specifies the secondary activation token
  fetched from the primary. If you set this to a blank string, the cluster will
  stay a secondary but clear its knowledge of any past primary (and thus not
  attempt to connect to the previous primary). This can be useful if the primary
  is down to stop the secondary from trying to reconnect to it.

- `primary_api_addr` `(string: )` – Specifies the API address (normal Vault
  address) to override the value embedded in the token. This can be useful if
  the primary's redirect address is not accessible directly from this cluster.

- `ca_file` `(string: "")` – Specifies the path to a CA root file (PEM format)
  that the secondary can use when unwrapping the token from the primary. If this
  and ca_path are not given, defaults to system CA roots.

- `ca_path` `string: ()` – Specifies the path to a CA root directory containing
  PEM-format files that the secondary can use when unwrapping the token from the
  primary. If this and ca_file are not given, defaults to system CA roots.

### Sample Payload

```json
{
  "token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/performance/secondary/update-primary
```
