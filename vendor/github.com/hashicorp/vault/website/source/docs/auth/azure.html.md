---
layout: "docs"
page_title: "Azure - Auth Methods"
sidebar_title: "Azure"
sidebar_current: "docs-auth-azure"
description: |-
  The azure auth method plugin allows automated authentication of Azure Active
  Directory.
---

# Azure Auth Method

The `azure` auth method allows authentication against Vault using
Azure Active Directory credentials. It treats Azure as a Trusted Third Party 
and expects a [JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519) 
signed by Azure Active Directory for the configured tenant.

Currently supports authentication for:

  * [Azure Managed Service Identity (MSI)](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview)

## Prerequisites:

The following documentation assumes that the method has been
[mounted](/docs/plugin/index.html) at `auth/azure`.

* A configured [Azure AD application](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-integrating-applications) which is used as the resource for generating MSI access tokens.
* Client credentials (shared secret) for accessing the Azure Resource Manager with read access to compute endpoints. See [Azure AD Service to Service Client Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-protocols-oauth-service-to-service)

Required Azure API permissions to be granted to Vault user:

* `Microsoft.Compute/virtualMachines/*/read`
* `Microsoft.Compute/virtualMachineScaleSets/*/read`

If Vault is hosted on Azure, Vault can use MSI to access Azure instead of a shared secret.  MSI must be [enabled](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/qs-configure-portal-windows-vm) on the VMs hosting Vault. 

The next sections review how the authN/Z workflows work. If you
have already reviewed these sections, here are some quick links to:

  * [Usage](#usage)
  * [API documentation](/api/auth/azure/index.html) docs.

## Authentication

### Via the CLI

The default path is `/auth/azure`. If this auth method was enabled at a different
path, specify `auth/my-path/login` instead.

```text
$ vault write auth/azure/login \
    role="dev-role" \
    jwt="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
    subscription_id="12345-..." \
    resource_group_name="test-group" \
    vm_name="test-vm"
```

The `role` and `jwt` parameters are required. When using bound_service_principal_ids and bound_groups in the token roles, all the information is required in the JWT.  When using other bound_* parameters, calls to Azure APIs will be made and subscription id, resource group name, and vm name are all required and can be obtained through instance metadata.

For example:

```text
$ vault write auth/azure/login role="dev-role" \
     jwt="$(curl -s 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.hashicorp.com%2F' -H Metadata:true | jq -r '.access_token')" \
     subscription_id=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq -r '.compute | .subscriptionId')  \
     resource_group_name=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq -r '.compute | .resourceGroupName') \
     vm_name=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance?api-version=2017-08-01" | jq -r '.compute | .name')
```

### Via the API

The default endpoint is `auth/azure/login`. If this auth method was enabled
at a different path, use that value instead of `azure`.

```sh
$ curl \
    --request POST \
    --data '{"role": "dev-role", "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}' \
    https://127.0.0.1:8200/v1/auth/azure/login
```

The response will contain the token at `auth.client_token`:

```json
{
  "auth": {
    "client_token": "f33f8c72-924e-11f8-cb43-ac59d697597c",
    "accessor": "0e9e354a-520f-df04-6867-ee81cae3d42d",
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

## Configuration

Auth methods must be configured in advance before machines can authenticate.
These steps are usually completed by an operator or configuration management 
tool.

### Via the CLI

1. Enable Azure authentication in Vault:

    ```text
    $ vault auth enable azure
    ```

1. Configure the Azure auth method:

    ```text
    $ vault write auth/azure/config \
        tenant_id= 7cd1f227-ca67-4fc6-a1a4-9888ea7f388c \
        resource=https://vault.hashicorp.com \
        client_id=dd794de4-4c6c-40b3-a930-d84cd32e9699 \
        client_secret=IT3B2XfZvWnfB98s1cie8EMe7zWg483Xy8zY004=
    ```

    For the complete list of configuration options, please see the API
    documentation.

1. Create a role:

    ```text
    $ vault write auth/azure/role/dev-role \
        policies="prod,dev" \
        bound_subscription_ids=6a1d5988-5917-4221-b224-904cd7e24a25 \
        bound_resource_groups=vault
    ```

    Roles are associated with an authentication type/entity and a set of Vault
    policies. Roles are configured with constraints specific to the
    authentication type, as well as overall constraints and configuration for
    the generated auth tokens.

    For the complete list of role options, please see the [API documentation](/api/auth/azure/index.html).

### Via the API

1. Enable Azure authentication in Vault:

    ```sh
    $ curl \
        --header "X-Vault-Token: ..." \
        --request POST \
        --data '{"type": "azure"}' \
        https://127.0.0.1:8200/v1/sys/auth/azure
    ```

1. Configure the Azure auth method:

    ```sh
    $ curl \
        --header "X-Vault-Token: ..." \
        --request POST \
        --data '{"tenant_id": "...", "resource": "..."}' \
        https://127.0.0.1:8200/v1/auth/azure/config
    ```

1. Create a role:

    ```sh
    $ curl \
        --header "X-Vault-Token: ..." \
        --request POST \
        --data '{"policies": ["dev", "prod"], ...}' \
        https://127.0.0.1:8200/v1/auth/azure/role/dev-role
    ```

### Plugin Setup

~> The following section is only relevant if you decide to enable the azure auth
method as an external plugin. The azure plugin method is integrated into Vault as
a builtin method by default.

Assuming you have saved the binary `vault-plugin-auth-azure` to some folder and
configured the [plugin directory](/docs/internals/plugins.html#plugin-directory)
for your server at `path/to/plugins`:


1. Enable the plugin in the catalog:

    ```text
    $ vault write sys/plugins/catalog/auth/azure-auth \
        command="vault-plugin-auth-azure" \
        sha256="..."
    ```

1. Enable the azure auth method as a plugin:

    ```text
    $ vault auth enable -path=azure azure-auth
    ```

## API

The Azure Auth Plugin has a full HTTP API. Please see the [API documentation]
(/api/auth/azure/index.html) for more details.

