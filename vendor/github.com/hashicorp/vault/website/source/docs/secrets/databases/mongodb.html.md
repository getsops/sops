---
layout: "docs"
page_title: "MongoDB - Database - Secrets Engines"
sidebar_title: "MongoDB"
sidebar_current: "docs-secrets-databases-mongodb"
description: |-
  MongoDB is one of the supported plugins for the database secrets engine. This
  plugin generates database credentials dynamically based on configured roles
  for the MongoDB database.
---

# MongoDB Database Secrets Engine

MongoDB is one of the supported plugins for the database secrets engine. This
plugin generates database credentials dynamically based on configured roles for
the MongoDB database.

See the [database secrets engine](/docs/secrets/databases/index.html) docs for
more information about setting up the database secrets engine.

## Setup

1. Enable the database secrets engine if it is not already enabled:

    ```text
    $ vault secrets enable database
    Success! Enabled the database secrets engine at: database/
    ```

    By default, the secrets engine will enable at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure Vault with the proper plugin and connection information:

    ```text
    $ vault write database/config/my-mongodb-database \
        plugin_name=mongodb-database-plugin \
        allowed_roles="my-role" \
        connection_url="mongodb://{{username}}:{{password}}@mongodb.acme.com:27017/admin?ssl=true" \
        username="admin" \
        password="Password!"
    ```

1. Configure a role that maps a name in Vault to a MongoDB command that executes and
creates the database credential:

    ```text
    $ vault write database/roles/my-role \
        db_name=my-mongodb-database \
        creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }' \
        default_ttl="1h" \
        max_ttl="24h"
    Success! Data written to: database/roles/my-role
    ```

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

1. Generate a new credential by reading from the `/creds` endpoint with the name
of the role:

    ```text
    $ vault read database/creds/my-role
    Key                Value
    ---                -----
    lease_id           database/creds/my-role/2f6a614c-4aa2-7b19-24b9-ad944a8d4de6
    lease_duration     1h
    lease_renewable    true
    password           8cab931c-d62e-a73d-60d3-5ee85139cd66
    username           v-root-e2978cd0-
    ```

## API

The full list of configurable options can be seen in the [MongoDB database
plugin API](/api/secret/databases/mongodb.html) page.

For more information on the database secrets engine's HTTP API please see the
[Database secrets engine API](/api/secret/databases/index.html) page.
