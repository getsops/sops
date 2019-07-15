---
layout: "api"
page_title: "Oracle - Database - Secrets Engines - HTTP API"
sidebar_title: "Oracle"
sidebar_current: "api-http-secret-databases-oracle-maria"
description: |-
  The Oracle plugin for Vault's database secrets engine generates database credentials to access Oracle servers.
---

# Oracle Database Plugin HTTP API

The Oracle database plugin is one of the supported plugins for the database
secrets engine. This plugin generates database credentials dynamically based on
configured roles for the Oracle database.

## Configure Connection

In addition to the parameters defined by the [Database
Backend](/api/secret/databases/index.html#configure-connection), this plugin
has a number of parameters to further configure a connection.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/database/config/:name`     |

### Parameters
- `connection_url` `(string: <required>)` - Specifies the Oracle DSN.

- `max_open_connections` `(int: 2)` - Specifies the maximum number of open
  connections to the database.

- `max_idle_connections` `(int: 0)` - Specifies the maximum number of idle
  connections to the database. A zero uses the value of `max_open_connections`
  and a negative value disables idle connections. If larger than
  `max_open_connections` it will be reduced to be equal.

- `max_connection_lifetime` `(string: "0s")` - Specifies the maximum amount of
  time a connection may be reused. If <= 0s connections are reused forever.

### Sample Payload

```json
{
  "plugin_name": "oracle-database-plugin",
  "allowed_roles": "readonly",
  "connection_url": "system/Oracle@localhost:1521/OraDoc.localhost",
  "max_open_connections": 5,
  "max_connection_lifetime": "5s"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/database/config/oracle
```

## Statements

Statements are configured during role creation and are used by the plugin to
determine what is sent to the database on user creation, renewing, and
revocation. For more information on configuring roles see the [Role
API](/api/secret/databases/index.html#create-role) in the database secrets engine docs.

### Parameters

The following are the statements used by this plugin. If not mentioned in this
list the plugin does not support that statement type.

- `creation_statements` `(string: <required>)` – Specifies the database
  statements executed to create and configure a user. Must be a
  semicolon-separated string, a base64-encoded semicolon-separated string, a
  serialized JSON string array, or a base64-encoded serialized JSON string
  array. The '{{name}}' and '{{password}}' values will be substituted.

- `revocation_statements` `(string: "")` – Specifies the database statements to
  be executed to revoke a user. Must be a semicolon-separated string, a
  base64-encoded semicolon-separated string, a serialized JSON string array, or
  a base64-encoded serialized JSON string array. The '{{name}}' value will be
  substituted. If not provided defaults to a generic drop user statement.
