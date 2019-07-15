---
layout: "api"
page_title: "Nomad Secret Backend - HTTP API"
sidebar_title: "Nomad"
sidebar_current: "api-http-secret-nomad"
description: |-
  This is the API documentation for the Vault Nomad secret backend.
---

# Nomad Secret Backend HTTP API

This is the API documentation for the Vault Nomad secret backend. For general
information about the usage and operation of the Nomad backend, please see the
[Vault Nomad backend documentation](/docs/secrets/nomad/index.html).

This documentation assumes the Nomad backend is mounted at the `/nomad` path
in Vault. Since it is possible to mount secret backends at any location, please
update your API calls accordingly.

## Configure Access

This endpoint configures the access information for Nomad. This access
information is used so that Vault can communicate with Nomad and generate
Nomad tokens.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/nomad/config/access`       |

### Parameters

- `address` `(string: "")` – Specifies the address of the Nomad
  instance, provided as `"protocol://host:port"` like `"http://127.0.0.1:4646"`.
  This value can also be provided on individual calls with the NOMAD_ADDR 
  environment variable.

- `token` `(string: "")` – Specifies the Nomad Management token to use.
  This value can also be provided on individual calls with the NOMAD_TOKEN 
  environment variable.

- `max_token_name_length` `(int: <optional>)` – Specifies the maximum length to
  use for the name of the Nomad token generated with [Generate
  Credential](#generate-credential). If omitted, `0` is used and ignored,
  defaulting to the max value allowed by the Nomad version. For Nomad versions
  0.8.3 and earlier, the default is `64`. For Nomad version 0.8.4 and later, the default is
  `256`.

### Sample Payload

```json
{
  "address": "http://127.0.0.1:4646",
  "token": "adha...",
  "max_token_name_length": 256
}
```

### Sample Request

```
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/nomad/config/access
```

## Read Access Configuration

This endpoint queries for information about the Nomad connection.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/nomad/config/access`       |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/nomad/config/access
```

### Sample Response

```json
  "data": {
    "address": "http://localhost:4646/"
  }
```

## Configure Lease

This endpoint configures the lease settings for generated tokens.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/nomad/config/lease`        |

### Parameters

- `ttl` `(string: "")` – Specifies the ttl for the lease. This is provided
  as a string duration with a time suffix like `"30s"` or `"1h"` or as total 
  seconds.

- `max_ttl` `(string: "")` – Specifies the max ttl for the lease. This is 
  provided as a string duration with a time suffix like `"30s"` or `"1h"` or as 
  total seconds.

### Sample Payload

```json
{
  "ttl": 1800,
  "max_ttl": 3600
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/nomad/config/lease
```

## Read Lease Configuration

This endpoint queries for information about the Lease TTL for the specified mount.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/nomad/config/lease`        |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/nomad/config/lease
```

### Sample Response

```json
  "data": {
    "max_ttl": 86400,
    "ttl": 86400
  }
```

## Delete Lease Configuration

This endpoint deletes the lease configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/nomad/config/lease`         |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/nomad/config/lease
```

## Create/Update Role

This endpoint creates or updates the Nomad role definition in Vault. If the role does not exist, it will be created. If the role already exists, it will receive
updated attributes.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/nomad/role/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Nomad tokens. This is part of the request URL.

- `policies` `(string: "")` – Comma separated list of Nomad policies the token is going to be created against. These need to be created beforehand in Nomad.

- `global` `(bool: "false")` – Specifies if the token should be global, as defined in the [Nomad Documentation](https://www.nomadproject.io/guides/acl.html#acl-tokens).

- `type` `(string: "client")` - Specifies the type of token to create when
  using this role. Valid values are `"client"` or `"management"`.

### Sample Payload

To create a client token with a custom policy:

```json
{
  "policies": "readonly"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --header "X-Vault-Token: ..." \
    --data @payload.json \
    http://127.0.0.1:8200/v1/nomad/role/monitoring
```

## Read Role

This endpoint queries for information about a Nomad role with the given name.
If no role exists with that name, a 404 is returned.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/nomad/role/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to query. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/nomad/role/monitoring
```

### Sample Response

```json
{
  "data": {
    "lease": "0s",
    "policies": [
      "example"
    ],
    "token_type": "client"
  }
}
```

## List Roles

This endpoint lists all existing roles in the backend.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`    | `/nomad/role`              |
| `GET`     | `/nomad/role?list=true`    |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/nomad/role
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "example"
    ]
  }
}
```

## Delete Role

This endpoint deletes a Nomad role with the given name. Even if the role does
not exist, this endpoint will still return a successful response.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/nomad/role/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --request DELETE \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/nomad/role/example-role
```

## Generate Credential

This endpoint generates a dynamic Nomad token based on the given role
definition.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/nomad/creds/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of an existing role against
  which to create this Nomad token. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/nomad/creds/example
```

### Sample Response

```json
{
  "data": {
    "accessor_id": "c834ba40-8d84-b0c1-c084-3a31d3383c03",
    "secret_id": "65af6f07-7f57-bb24-cdae-a27f86a894ce"
  }
}
```
