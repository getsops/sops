---
layout: "api"
page_title: "KV - Secrets Engines - HTTP API"
sidebar_title: "K/V Version 1"
sidebar_current: "api-http-secret-kv-v1"
description: |-
  This is the API documentation for the Vault KV secrets engine.
---

# KV Secrets Engine - Version 1 (API)

This is the API documentation for the Vault KV secrets engine. For general
information about the usage and operation of the kv secrets engine, please
see the [Vault kv documentation](/docs/secrets/kv/index.html).

This documentation assumes the kv secrets engine is enabled at the
`/secret` path in Vault. Since it is possible to enable secrets engines at any
location, please update your API calls accordingly.

## Read Secret

This endpoint retrieves the secret at the specified location.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/secret/:path`              |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to read.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://127.0.0.1:8200/v1/secret/my-secret
```

### Sample Response

```json
{
  "auth": null,
  "data": {
    "foo": "bar",
    "ttl": "1h"
  },
  "lease_duration": 3600,
  "lease_id": "",
  "renewable": false
}
```

_Note_: the `lease_duration` field, which will be populated if a "ttl" field
was included in the data, is advisory. No lease is created. This is a way for
writers to indicate how often a given value should be re-read by the client.
See the [Vault KV secrets engine documentation](/docs/secrets/kv/index.html)
for more details.

## List Secrets

This endpoint returns a list of key names at the specified location. Folders are
suffixed with `/`. The input must be a folder; list on a file will not return a
value. Note that no policy-based filtering is performed on keys; do not encode
sensitive information in key names. The values themselves are not accessible via
this command.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/secret/:path`              |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to list.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/secret/my-secret
```

### Sample Response

The example below shows output for a query path of `secret/` when there are
secrets at `secret/foo` and `secret/foo/bar`; note the difference in the two
entries.

```json
{
  "auth": null,
  "data": {
    "keys": ["foo", "foo/"]
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
```

## Create/Update Secret

This endpoint stores a secret at the specified location. If the value does not
yet exist, the calling token must have an ACL policy granting the `create`
capability. If the value already exists, the calling token must have an ACL
policy granting the `update` capability.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/secret/:path`              |
| `PUT`    | `/secret/:path`              |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secrets to
  create/update. This is specified as part of the URL.

- `:key` `(string: "")` – Specifies a key, paired with an associated value, to
  be held at the given location. Multiple key/value pairs can be specified, and
  all will be returned on a read operation. A key called `ttl` will trigger
  some special behavior. See the [Vault KV secrets engine
  documentation](/docs/secrets/kv/index.html) for details.

### Sample Payload

```json
{
  "foo": "bar",
  "zip": "zap"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/secret/my-secret
```

## Delete Secret

This endpoint deletes the secret at the specified location.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/secret/:path`              |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the secret to delete.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/secret/my-secret
```
