---
layout: "api"
page_title: "Azure - Secrets Engines - HTTP API"
sidebar_title: "Azure"
sidebar_current: "api-http-secret-azure"
description: |-
  This is the API documentation for the Vault Azure secrets engine.
---

# Azure Secrets Engine (API)

This is the API documentation for the Vault Azure
secrets engine. For general information about the usage and operation of
the Azure secrets engine, please see the main [Azure secrets documentation][docs].

This documentation assumes the Azure secrets engine is enabled at the `/azure` path
in Vault. Since it is possible to mount secrets engines at any path, please
update your API calls accordingly.

## Configure Access

Configures the credentials required for the plugin to perform API calls
to Azure. These credentials will be used to query roles and create/delete
service principals. Environment variables will override any parameters set in the config.

| Method   | Path                     |
| :------------------------| :------------------------ |
| `POST`   | `/azure/config`            |

- `subscription_id` (`string: <required>`) - The subscription id for the Azure Active Directory.
  This value can also be provided with the AZURE_SUBSCRIPTION_ID environment variable.
- `tenant_id` (`string: <required>`) - The tenant id for the Azure Active Directory.
  This value can also be provided with the AZURE_TENANT_ID environment variable.
- `client_id` (`string:""`) - The OAuth2 client id to connect to Azure. This value can also be provided
  with the AZURE_CLIENT_ID environment variable. See [authentication](#Authentication) for more details.
- `client_secret` (`string:""`) - The OAuth2 client secret to connect to Azure. This value can also be
  provided with the AZURE_CLIENT_ID environment variable. See [authentication](#Authentication) for more details.
- `environment` (`string:""`) - The Azure environment. This value can also be provided with the AZURE_ENVIRONMENT
  environment variable. If not specified, Vault will use Azure Public Cloud.

### Sample Payload

```json
{
  "subscription_id": "94ca80...",
  "tenant_id": "d0ac7e...",
  "client_id": "e607c4...",
  "client_secret": "9a6346...",
  "environment": "AzureGermanCloud"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/azure/config
```

## Read Config

Return the stored configuration, omitting `client_secret`.

| Method   | Path                     |
| :------------------------| :------------------------ |
| `GET`    | `/azure/config`            |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/azure/config
```

### Sample Response

```json
{
  "data": {
    "subscription_id": "94ca80...",
    "tenant_id": "d0ac7e...",
    "client_id": "e607c4...",
    "environment": "AzureGermanCloud"
  },
  ...
}
```

## Delete Config

Deletes the stored Azure configuration and credentials.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/auth/azure/config`         |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/auth/azure/config
```


## Create/Update Role

Create or update a Vault role. Either `application_object_id` or
`azure_roles` must be provided, and these resources must exist for this
call to succeed. See the Azure secrets [roles docs][roles] for more
information about roles.

| Method   | Path                     |
| :------------------------| :------------------------ |
| `POST`   | `/azure/roles/:name`     |


### Parameters

- `azure_roles` (`string: ""`) - List of Azure roles to be assigned to the generated service
   principal. The array must be in JSON format, properly escaped as a string. See [roles docs][roles]
   for details on role definition.
- `application_object_id` (`string: ""`) - Application Object ID for an existing service principal that will
   be used instead of creating dynamic service principals. If present, `azure_roles` will be ignored. See
   [roles docs][roles] for details on role definition.
- `ttl` (`string: ""`) – Specifies the default TTL for service principals generated using this role.
   Accepts time suffixed strings ("1h") or an integer number of seconds. Defaults to the system/engine default TTL time.
- `max_ttl` (`string: ""`) – Specifies the maximum TTL for service principals generated using this role. Accepts time
   suffixed strings ("1h") or an integer number of seconds. Defaults to the system/engine max TTL time.

### Sample Payload

```json
{
  "azure_roles": "[
    {
      \"role_name\": \"Contributor\",
      \"scope\":  \"/subscriptions/<uuid>/resourceGroup/Website\"
    },
    {
      \"role_id\": \"/subscriptions/<uuid>/providers/Microsoft.Authorization/roleDefinitions/<uuid>\",
      \"scope\":  \"/subscriptions/<uuid>\"
    }
  ]",
  "ttl": 3600,
  "max_ttl": "24h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/azure/roles/my-role
```


## List Roles

Lists all of the roles that are registered with the plugin.

| Method   | Path                     |
| :------------------------| :------------------------ |
| `LIST`   | `/azure/roles`           |


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/azure/roles
```

### Sample Response

```json
{
  "data": {
     "keys": [
       "my-role-one",
       "my-role-two"
     ]
   }
 }
```

## Generate Credentials

This endpoint generates a new service principal based on the named role.

| Method   | Path                     |
| :------------------------| :------------------------ |
| `GET`    | `/azure/creds/:name`     |

### Parameters

- `name` (`string: <required>`) - Specifies the name of the role to create credentials against.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/azure/creds/my-role
```

### Sample Response

```json
{
  "data": {
    "client_id": "408bf248-dd4e-4be5-919a-7f6207a307ab",
    "client_secret": "ad06228a-2db9-4e0a-8a5d-e047c7f32594",
    ...
  }
}
```

## Revoking/Renewing Secrets

See docs on how to [renew](/api/system/leases.html#renew-lease) and [revoke](/api/system/leases.html#revoke-lease) leases.


[docs]: /docs/secrets/azure/index.html
[roles]: /docs/secrets/azure/index.html#roles
