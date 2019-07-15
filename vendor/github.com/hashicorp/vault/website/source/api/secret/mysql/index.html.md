---
layout: "api"
page_title: "MySQL - Secrets Engines - HTTP API"
sidebar_title: "MySQL <sup>DEPRECATED</sup>"
sidebar_current: "api-http-secret-mysql"
description: |-
  This is the API documentation for the Vault MySQL secrets engine.
---

# MySQL Secrets Engine (API)

~> **Deprecation Note:** This secrets engine is deprecated in favor of the
combined databases secrets engine added in v0.7.1. See the API documentation for
the new implementation of this secrets engine at
[MySQL/MariaDB database plugin HTTP API](/api/secret/databases/mysql-maria.html).

This is the API documentation for the Vault MySQL secrets engine. For general
information about the usage and operation of the MySQL secrets engine, please
see the [Vault MySQL documentation](/docs/secrets/mysql/index.html).

This documentation assumes the MySQL secrets engine is enabled at the `/mysql`
path in Vault. Since it is possible to enable secrets engines at any location,
please update your API calls accordingly.

## Configure Connection

This endpoint configures the connection DSN used to communicate with MySQL.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/mysql/config/connection`   |

### Parameters

- `connection_url` `(string: <required>)` – Specifies the MySQL DSN.

- `max_open_connections` `(int: 2)` – Specifies the maximum number of open
  connections to the database.

- `max_idle_connections` `(int: 0)` – Specifies the maximum number of idle
  connections to the database. A zero uses the value of `max_open_connections`
  and a negative value disables idle connections. If larger than
  `max_open_connections` it will be reduced to be equal.

- `verify_connection` `(bool: true)` – Specifies if the connection is verified
  during initial configuration.

### Sample Payload

```json
{
  "connection_url": "mysql:host=localhost;dbname=testdb"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/mysql/config/connection
```

## Configure Lease

This endpoint configures the lease settings for generated credentials. If not
configured, leases default to 1 hour. This is a root protected endpoint.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/mysql/config/lease`        |

### Parameters

- `lease` `(string: <required>)` – Specifies the lease value provided as a
  string duration with time suffix. "h" (hour) is the largest suffix.

- `lease_max` `(string: <required>)` – Specifies the maximum lease value
  provided as a string duration with time suffix. "h" (hour) is the largest
  suffix.

### Sample Payload

```json
{
  "lease": "12h",
  "lease_max": "24h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/mysql/config/lease
```

## Create Role

This endpoint creates or updates the role definition.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/mysql/roles/:name`         |

### Parameters

- `sql` `(string: <required>)` – Specifies the SQL statements executed to create
  and configure a user. Must be a semicolon-separated string, a base64-encoded
  semicolon-separated string, a serialized JSON string array, or a
  base64-encoded serialized JSON string array. The '{{name}}' and
  '{{password}}' values will be substituted.

- `revocation_sql` `(string: "")` – Specifies the SQL statements executed to
  revoke a user. Must be a semicolon-separated string, a base64-encoded
  semicolon-separated string, a serialized JSON string array, or a
  base64-encoded serialized JSON string array. The '{{name}}' value will be
  substituted.

- `rolename_length` `(int: 4)` – Specifies how many characters from the role
  name will be used to form the mysql username interpolated into the '{{name}}'
  field of the sql parameter.  

- `displayname_length` `(int: 4)` – Specifies how many characters from the token
  display name will be used to form the mysql username interpolated into the
  '{{name}}' field of the sql parameter.  

- `username_length` `(int: 16)` – Specifies the maximum total length in
  characters of the mysql username interpolated into the '{{name}}' field of the
  sql parameter.

### Sample Payload

```json
{
  "sql": "CREATE USER ..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/mysql/roles/my-role
```

## Read Role

This endpoint queries the role definition.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/mysql/roles/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/mysql/roles/my-role
```

### Sample Response

```json
{
  "data": {
    "sql": "CREATE USER..."
  }
}
```

## List Roles

This endpoint returns a list of available roles. Only the role names are
returned, not any values.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/mysql/roles`               |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/mysql/roles
```

### Sample Response

```json
{
  "auth": null,
  "data": {
    "keys": ["dev", "prod"]
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
```

## Delete Role

This endpoint deletes the role definition.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/mysql/roles/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/mysql/roles/my-role
```

## Generate Credentials

This endpoint generates a new set of dynamic credentials based on the named
role.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/mysql/creds/:name`         |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create
  credentials against. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/mysql/creds/my-role
```

### Sample Response

```json
{
  "data": {
    "username": "user-role-aefa63",
    "password": "132ae3ef-5a64-7499-351e-bfe59f3a2a21"
  }
}
```
