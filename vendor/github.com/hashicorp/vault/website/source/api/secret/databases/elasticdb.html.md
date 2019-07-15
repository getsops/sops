---
layout: "api"
page_title: "Elasticsearch - Database - Secrets Engines - HTTP API"
sidebar_title: "Elasticsearch"
sidebar_current: "api-http-secret-databases-elasticsearch"
description: |-
  The Elasticsearch plugin for Vault's database secrets engine generates database credentials to access Elasticsearch.
---

# Elasticsearch Database Plugin HTTP API

The Elasticsearch database plugin is one of the supported plugins for the database
secrets engine. This plugin generates credentials dynamically based on
configured roles for Elasticsearch.

## Configure Connection

In addition to the parameters defined by the [Database
Backend](/api/secret/databases/index.html#configure-connection), this plugin
has a number of parameters to further configure a connection.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/database/config/:name`     |

### Parameters

- `url` `(string: <required>)` - The URL for Elasticsearch's API ("http://localhost:9200").
- `username` `(string: <required>)` - The username to be used in the connection URL ("vault"). 
- `password` `(string: <required>)` - The password to be used in the connection URL ("pa55w0rd").
- `ca_cert` `(string: "")` - The path to a PEM-encoded CA cert file to use to verify the Elasticsearch server's identity.
- `ca_path` `(string: "")` - The path to a directory of PEM-encoded CA cert files to use to verify the Elasticsearch server's identity.
- `client_cert` `(string: "")` - The path to the certificate for the Elasticsearch client to present for communication.
- `client_key` `(string: "")` - The path to the key for the Elasticsearch client to use for communication.
- `tls_server_name` `(string: "")` - This, if set, is used to set the SNI host when connecting via 1TLS.
- `insecure` `(bool: false)` - Not recommended. Default to false. Can be set to true to disable SSL verification.

### Sample Payload

```json
{
  "plugin_name": "elasticsearch-database-plugin",
  "allowed_roles": "internally-defined-role,externally-defined-role",
  "url": "http://localhost:9200",
  "username": "vault",
  "password": "myPa55word",
  "ca_cert": "/usr/share/ca-certificates/extra/elastic-stack-ca.crt.pem",
  "client_cert": "$ES_HOME/config/certs/elastic-certificates.crt.pem",
  "client_key": "$ES_HOME/config/certs/elastic-certificates.key.pem"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/database/config/my-elasticsearch-database
```

## Statements

Statements are configured during role creation and are used by the plugin to
determine what is sent to the database on user creation, renewing, and
revocation. For more information on configuring roles see the [Role
API](/api/secret/databases/index.html#create-role) in the database secrets engine docs.

### Parameters

The following are the statements used by this plugin. If not mentioned in this
list the plugin does not support that statement type.

- `creation_statements` `(string: <required>)` – Using JSON, either defines an
  `elasticsearch_role_definition` or a group of pre-existing `elasticsearch_roles`.
  The object specified by the `elasticsearch_role_definition` is the JSON directly
  passed through to the Elasticsearch API, so you can pass through anything shown
  [here](https://www.elastic.co/guide/en/elasticsearch/reference/6.6/security-api-put-role.html).
  For `elasticsearch_roles`, add the names of the roles only. They must pre-exist
  in Elasticsearch. Defining roles in Vault is more secure than using pre-existing
  roles because a privilege escalation could be performed by editing the roles used
  out-of-band in Elasticsearch.

### Sample Creation Statements
```json
{
  "elasticsearch_role_definition": {
    "indices": [{
      "names": ["*"],
      "privileges": ["read"]
    }]
  }
}
```
```json
{
  "elasticsearch_roles": ["pre-existing-role-in-elasticsearch"]
}
```