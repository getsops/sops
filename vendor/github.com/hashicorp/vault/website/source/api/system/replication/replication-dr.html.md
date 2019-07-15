---
layout: "api"
page_title: "/sys/replication - HTTP API"
sidebar_title: "<code>/sys/replication/dr</code>"
sidebar_current: "api-http-system-replication-dr"
description: |-
  The '/sys/replication/dr' endpoint focuses on managing general operations in Vault Enterprise Disaster Recovery replication
---

# `/sys/replication/dr`

~> **Enterprise Only** – These endpoints require Vault Enterprise.

## Check DR Status

This endpoint prints information about the status of replication (mode,
sync progress, etc).

This is an authenticated endpoint.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/dr/status`    |

### Sample Request

```
$ curl \
    http://127.0.0.1:8200/v1/sys/replication/dr/status
```

### Sample Response from Primary

The printed status of the replication environment. As an example, for a
primary, it will look something like:

```json
{
  "data": {
    "cluster_id": "d4095d41-3aee-8791-c421-9bc7f88f7c3e",
    "known_secondaries": [],
    "last_wal": 241,
    "merkle_root": "56794a98e52598f35974024fba6691f047e772e9",
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
    "last_remote_wal": 241,
    "merkle_root": "56794a98e52598f35974024fba6691f047e772e9",
    "mode": "secondary",
    "primary_cluster_addr": "https://127.0.0.1:8201",
    "secondary_id": "3",
    "state": "stream-wals"
  },
}
```

## Enable DR Primary Replication

This endpoint enables DR replication in primary mode. This is used when DR replication
is currently disabled on the cluster (if the cluster is already a secondary, it
must be promoted).

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/enable` |

### Parameters

- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address,
  such as through a TCP-based load balancer.

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
    http://127.0.0.1:8200/v1/sys/replication/dr/primary/enable
```

## Demote DR Primary

This endpoint demotes a DR primary cluster to a secondary. This DR secondary cluster
will not attempt to connect to a primary (see the update-primary call), but will
maintain knowledge of its cluster ID and can be reconnected to the same
DR replication set without wiping local storage.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/demote` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/dr/primary/demote
```

## Disable DR Primary

This endpoint disables DR replication entirely on the cluster. Any secondaries will
no longer be able to connect. Caution: re-enabling this node as a primary or
secondary will change its cluster ID; in the secondary case this means a wipe of
the underlying storage when connected to a primary, and in the primary case,
secondaries connecting back to the cluster (even if they have connected before)
will require a wipe of the underlying storage.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/disable` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/replication/dr/primary/disable
```

## Generate DR Secondary Token

This endpoint generates a DR secondary activation token for the
cluster with the given opaque identifier, which must be unique. This
identifier can later be used to revoke a DR secondary's access.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/sys/replication/dr/primary/secondary-token` |

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
    http://127.0.0.1:8200/v1/sys/replication/dr/primary/secondary-token
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

## Revoke DR Secondary Token

This endpoint revokes a DR secondary's ability to connect to the DR primary cluster;
the DR secondary will immediately be disconnected and will not be allowed to
connect again unless given a new activation token.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/primary/revoke-secondary` |

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
    http://127.0.0.1:8200/v1/sys/replication/dr/primary/revoke-secondary
```

## Enable DR Secondary

This endpoint enables replication on a DR secondary using a DR secondary activation
token.

!> This will immediately clear all data in the secondary cluster!

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/enable` |

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
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/enable
```

## Promote DR Secondary

This endpoint promotes the DR secondary cluster to DR primary. For data safety and
security reasons, new secondary tokens will need to be issued to other
secondaries, and there should never be more than one primary at a time.

If the DR secondary's primary cluster is also in a performance replication set,
the DR secondary will be promoted into that replication set. Care should be
taken when promoting to ensure multiple performance primary clusters are not
activate at the same time.

If the DR secondary's primary cluster is a performance secondary, the promoted
cluster will attempt to connect to the performance primary cluster using the
same secondary token.

This endpoint requires a DR Operation Token to be provided as means of
authorization. See the [DR Operation Token API
docs](#generate-disaster-recovery-operation-token) for more information.

!> Only one performance primary should be active at a given time. Multiple primaries may
result in data loss!

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/promote` |

### Parameters

- `dr_operation_token` `(string: <required>)` - DR operation token used to authorize this request.
- `primary_cluster_addr` `(string: "")` – Specifies the cluster address that the
  primary gives to secondary nodes. Useful if the primary's cluster address is
  not directly accessible and must be accessed via an alternate path/address
  (e.g. through a load balancer).
- `force` `(bool: false)` - If true the cluster will be promoted even if it fails
  certain safety checks. Caution: Forcing promotion could result in data loss if
  data isn't fully replicated.

### Sample Payload

```json
{
  "dr_operation_token": "ijH8tphEHaBtgx+IvPfxDsSi2LV4j9k+Lad6eqT5cJw="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/promote
```

### Sample Response

```json
{
  "progress": 0,
  "required": 1,
  "complete": false,
  "request_id": "ad8f9074-0e24-d30e-83cd-595c9652ff89",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "complete": false,
    "progress": 0,
    "required": 1
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

## Update DR Secondary's Primary

This endpoint changes a DR secondary cluster's assigned primary cluster using a
secondary activation token. This does not wipe all data in the cluster.

This endpoint requires a DR Operation Token to be provided as means of
authorization. See the [DR Operation Token API
docs](#generate-disaster-recovery-operation-token) for more information.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/update-primary` |

### Parameters

- `dr_operation_token` `(string: <required>)` - DR operation token used to authorize this request.

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
  "dr_operation_token": "...",
  "token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/update-primary
```

## Generate Disaster Recovery Operation Token

The `/sys/replication/dr/secondary/generate-operation-token` endpoint is used to create a new Disaster
Recovery operation token for a DR secondary. These tokens are used to authorize
certain DR Operation. They should be treated like traditional root tokens by
being generated when needed and deleted soon after.

## Read Generation Progress

This endpoint reads the configuration and process of the current generation
attempt.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/replication/dr/secondary/generate-operation-token/attempt` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/generate-operation-token/attempt
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "pgp_fingerprint": "",
  "complete": false
}
```

If a generation is started, `progress` is how many unseal keys have been
provided for this generation attempt, where `required` must be reached to
complete. The `nonce` for the current attempt and whether the attempt is
complete is also displayed. If a PGP key is being used to encrypt the final
token, its fingerprint will be returned. Note that if an OTP is being used to
encode the final token, it will never be returned.

## Start Token Generation

This endpoint initializes a new generation attempt. Only a single
generation attempt can take place at a time.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/replication/dr/secondary/generate-operation-token/attempt` |

### Parameters

- `pgp_key` `(string: <optional>)` – Specifies a base64-encoded PGP public key.
  The raw bytes of the token will be encrypted with this value before being
  returned to the final unseal key provider.

### Sample Request

```
$ curl \
    --request PUT \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/generate-operation-token/attempt
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "otp": "2vPFYG8gUSW9npwzyvxXMug0",
  "otp_length" :24,
  "complete": false
}
```

## Cancel Generation

This endpoint cancels any in-progress generation attempt. This clears any
progress made. This must be called to change the OTP or PGP key being used.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/replication/dr/secondary/generate-operation-token/attempt` |

### Sample Request

```
$ curl \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/generate-operation-token/attempt
```

## Provide Key Share to Generate Token

This endpoint is used to enter a single master key share to progress the
generation attempt. If the threshold number of master key shares is reached,
Vault will complete the generation and issue the new token.  Otherwise,
this API must be called multiple times until that threshold is met. The attempt
nonce must be provided with each call.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/replication/dr/secondary/generate-operation-token/update`  |

### Parameters

- `key` `(string: <required>)` – Specifies a single master key share.

- `nonce` `(string: <required>)` – Specifies the nonce of the attempt.

### Sample Payload

```json
{
  "key": "acbd1234",
  "nonce": "ad235"
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/generate-operation-token/update
```

### Sample Response

This returns a JSON-encoded object indicating the attempt nonce, and completion
status, and the encoded token, if the attempt is complete.

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 3,
  "required": 3,
  "pgp_fingerprint": "",
  "complete": true,
  "encoded_token": "FPzkNBvwNDeFh4SmGA8c+w=="
}
```

## Delete DR Operation Token

This endpoint revokes the DR Operation Token. This token does not have a TTL
and therefore should be deleted when it is no longer needed.


| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/replication/dr/secondary/operation-token/delete` |

### Parameters

- `dr_operation_token` `(string: <required>)` - DR operation token used to authorize this request.

### Sample Payload

```json
{
  "dr_operation_token": "..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/replication/dr/secondary/operation-token/delete
```
