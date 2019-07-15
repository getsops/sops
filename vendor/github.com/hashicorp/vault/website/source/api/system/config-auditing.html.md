---
layout: "api"
page_title: "/sys/config/auditing - HTTP API"
sidebar_title: "<code>/sys/config/auditing</code>"
sidebar_current: "api-http-system-config-auditing"
description: |-
  The `/sys/config/auditing` endpoint is used to configure auditing settings.
---

# `/sys/config/auditing/request-headers`

The `/sys/config/auditing` endpoint is used to configure auditing settings.

## Read All Audited Request Headers

This endpoint lists the request headers that are configured to be audited.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/auditing/request-headers` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/config/auditing/request-headers
```

### Sample Response

```json
{
  "headers": {
    "X-Forwarded-For": {
      "hmac": true
    }
  }
}
```

## Read Single Audit Request Header

This endpoint lists the information for the given request header.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/config/auditing/request-headers/:name` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the request header to
  query. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/config/auditing/request-headers/my-header
```

### Sample Response

```json
{
  "X-Forwarded-For": {
    "hmac": true
  }
}
```

## Create/Update Audit Request Header

This endpoint enables auditing of a header.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/config/auditing/request-headers/:name` |

### Parameters

- `hmac` `(bool: false)` – Specifies if this header's value should be HMAC'ed in
  the audit logs.

### Sample Payload

```json
{
  "hmac": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/config/auditing/request-headers/my-header
```

## Delete Audit Request Header

This endpoint disables auditing of the given request header.

- **`sudo` required** – This endpoint requires `sudo` capability in addition to
  any path-specific capabilities.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/config/auditing/request-headers/:name` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/config/auditing/request-headers/my-header
```
