---
layout: "docs"
page_title: "Elasticsearch - Database - Secrets Engines"
sidebar_title: "Elasticsearch"
sidebar_current: "docs-secrets-databases-elasticdb"
description: |-
  Elasticsearch is one of the supported plugins for the database secrets engine. This
  plugin generates database credentials dynamically based on configured roles
  for Elasticsearch.
---

# Elasticsearch Database Secrets Engine

Elasticsearch is one of the supported plugins for the database secrets engine. This
plugin generates database credentials dynamically based on configured roles for
Elasticsearch.

See the [database secrets engine](/docs/secrets/databases/index.html) docs for
more information about setting up the database secrets engine.

## Getting Started

To take advantage of this plugin, you must first enable Elasticsearch's native realm of security by activating X-Pack. These
instructions will walk you through doing this using Elasticsearch 7.1.1.

### Enable X-Pack Security in Elasticsearch

Read [Securing the Elastic Stack](https://www.elastic.co/guide/en/elastic-stack-overview/7.1/elasticsearch-security.html) and 
follow [its instructions for enabling X-Pack Security](https://www.elastic.co/guide/en/elasticsearch/reference/7.1/setup-xpack.html). 

### Enable Encrypted Communications

This plugin communicates with Elasticsearch's security API. ES requires TLS for these communications so they can be
encrypted.

To set up TLS in Elasticsearch, first read [encrypted communications](https://www.elastic.co/guide/en/elastic-stack-overview/7.1/encrypting-communications.html)
and go through its instructions on [encrypting HTTP client communications](https://www.elastic.co/guide/en/elasticsearch/reference/7.1/configuring-tls.html#tls-http). 

After enabling TLS on the Elasticsearch side, you'll need to convert the .p12 certificates you generated to other formats so they can be 
used by Vault. [Here is an example using OpenSSL](https://stackoverflow.com/questions/15144046/converting-pkcs12-certificate-into-pem-using-openssl) 
to convert our .p12 certs to the pem format.

Also, on the instance running Elasticsearch, we needed to install our newly generated CA certificate that was originally in the .p12 format.
We did this by converting the .p12 CA cert to a pem, and then further converting that 
[pem to a crt](https://stackoverflow.com/questions/13732826/convert-pem-to-crt-and-key), adding that crt to `/usr/share/ca-certificates/extra`, 
and using `sudo dpkg-reconfigure ca-certificates`.

The above instructions may vary if you are not using an Ubuntu machine. Please ensure you're using the methods specific to your operating
environment. Describing every operating environment is outside the scope of these instructions.

### Set Up Passwords

When done, verify that you've enabled X-Pack by running `$ $ES_HOME/bin/elasticsearch-setup-passwords interactive`. You'll
know it's been set up successfully if it takes you through a number of password-inputting steps.

### Create a Role for Vault

Next, in Elasticsearch, we recommend that you create a user just for Vault to use in managing secrets.

To do this, first create a role that will allow Vault the minimum privileges needed to administer users and passwords by performing a
POST to Elasticsearch. To do this, we used the `elastic` superuser whose password we created in the
`$ $ES_HOME/bin/elasticsearch-setup-passwords interactive` step.

```
$ curl \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"cluster": ["manage_security"]}' \
    http://elastic:$PASSWORD@localhost:9200/_xpack/security/role/vault
```

Next, create a user for Vault associated with that role.

```
$ curl \
    -X POST \
    -H "Content-Type: application/json" \
    -d @data.json \
    http://elastic:$PASSWORD@localhost:9200/_xpack/security/user/vault
```

The contents of `data.json` in this example are:
```
{
 "password" : "myPa55word",
 "roles" : [ "vault" ],
 "full_name" : "Hashicorp Vault",
 "metadata" : {
   "plugin_name": "Vault Plugin Database Elasticsearch",
   "plugin_url": "https://github.com/hashicorp/vault-plugin-database-elasticsearch"
 }
}
```

Now, Elasticsearch is configured and ready to be used with Vault.

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
    $ vault write database/config/my-elasticsearch-database \
        plugin_name="elasticsearch-database-plugin" \
        allowed_roles="internally-defined-role,externally-defined-role" \
        username=vault \
        password=myPa55word \
        url=http://localhost:9200 \
        ca_cert=/usr/share/ca-certificates/extra/elastic-stack-ca.crt.pem \
        client_cert=$ES_HOME/config/certs/elastic-certificates.crt.pem \
        client_key=$ES_HOME/config/certs/elastic-certificates.key.pem   
    ```

1. Configure a role that maps a name in Vault to a role definition in Elasticsearch.
This is considered the most secure type of role because nobody can perform
a privilege escalation by editing a role's privileges out-of-band in 
Elasticsearch:

    ```text
    $ vault write database/roles/internally-defined-role \
          db_name=my-elasticsearch-database \
          creation_statements='{"elasticsearch_role_definition": {"indices": [{"names":["*"], "privileges":["read"]}]}}' \
          default_ttl="1h" \
          max_ttl="24h"   
    Success! Data written to: database/roles/internally-defined-role
    ```

1. Alternatively, configure a role that maps a name in Vault to a pre-existing 
role definition in Elasticsearch:

    ```text
    $ vault write database/roles/externally-defined-role \
         db_name=my-elasticsearch-database \
         creation_statements='{"elasticsearch_roles": ["pre-existing-role-in-elasticsearch"]}' \
         default_ttl="1h" \
         max_ttl="24h"
    Success! Data written to: database/roles/externally-defined-role
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

The full list of configurable options can be seen in the [Elasticsearch database
plugin API](/api/secret/databases/elasticdb.html) page.

For more information on the database secrets engine's HTTP API please see the
[Database secrets engine API](/api/secret/databases/index.html) page.
