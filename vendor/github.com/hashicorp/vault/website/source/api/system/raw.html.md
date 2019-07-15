---
layout: "api"
page_title: "/sys/raw - HTTP API"
sidebar_title: "<code>/sys/raw</code>"
sidebar_current: "api-http-system-raw"
description: |-
  The `/sys/raw` endpoint is used to access the raw underlying store in Vault.
---

# `/sys/raw`

The `/sys/raw` endpoint is used to access the raw underlying store in Vault.

This endpoint is off by default.  See the 
[Vault configuration documentation](/docs/configuration/index.html) to
enable.

## Read Raw

This endpoint reads the value of the key at the given path. This is the raw path
in the storage backend and not the logical path that is exposed via the mount
system.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/raw/:path`             |

### Parameters

- `path` `(string: <required>)` – Specifies the raw path in the storage backend.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    ---header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/raw/secret/foo
```

### Sample Response

```json
{
  "value": "{'foo':'bar'}"
}
```

## Create/Update Raw

This endpoint updates the value of the key at the given path. This is the raw
path in the storage backend and not the logical path that is exposed via the
mount system.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/raw/:path`             |

### Parameters

- `path` `(string: <required>)` – Specifies the raw path in the storage backend.
  This is specified as part of the URL.

- `value` `(string: <required>)` – Specifies the value of the key.

### Sample Payload

```json
{
  "value": "{\"foo\": \"bar\"}"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/raw/secret/foo
```

## List Raw

This endpoint returns a list keys for a given path prefix.

**This endpoint requires 'sudo' capability.**

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/sys/raw/:prefix` |
| `GET`   | `/sys/raw/:prefix?list=true` |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/sys/raw/logical
```

### Sample Response

```json
{
  "data":{
    "keys":[
      "abcd-1234...",
      "efgh-1234...",
      "ijkl-1234..."
    ]
  }
}
```

## Delete Raw

This endpoint deletes the key with given path. This is the raw path in the
storage backend and not the logical path that is exposed via the mount system.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/raw/:path`             |

### Parameters

- `path` `(string: <required>)` – Specifies the raw path in the storage backend.
  This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/raw/secret/foo
```
