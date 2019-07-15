---
layout: "api"
page_title: "/sys/namespaces - HTTP API"
sidebar_title: "<code>/sys/namespaces</code>"
sidebar_current: "api-http-system-namespaces"
description: |-
  The `/sys/namespaces` endpoint is used manage namespaces in Vault.
---

# `/sys/namespaces`

The `/sys/namespaces` endpoint is used manage namespaces in Vault.

## List Namespaces

This endpoints lists all the namespaces.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/sys/namespaces`            |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    -X LIST \
    http://127.0.0.1:8200/v1/sys/namespaces
```

### Sample Response

```json
[
    "ns1/",
    "ns2/"
]

```

## Create Namespace

This endpoint creates a namespace at the givent path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/sys/namespaces/:path`      |

### Parameters

- `path` `(string: <required>)` – Specifies the path where the namespace
  will be namespace. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/sys/namespaces/ns1
```

## Delete Namespace

This endpoint deletes a namespace at the specified path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/namespaces/:path`      | `204 (empty body)    ` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/namespaces/ns1
```

## Read Namespace Information

This endpoint get the metadata for the given namespace path.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/namespaces/:path`      |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/namespaces/ns1
```

### Sample Response

```json
{
  "id": "gsudj",
  "path": "ns1/"
}
```
