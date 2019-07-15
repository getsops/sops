---
layout: "api"
page_title: "Active Directory - Secrets Engines - HTTP API"
sidebar_title: "Active Directory"
sidebar_current: "api-http-secret-active-directory"
description: |-
  This is the API documentation for the Vault Active Directory secrets engine.
---

# Active Directory Secrets Engine (API)

This is the API documentation for the Vault AD secrets engine. For general
information about the usage and operation of the AD secrets engine, please see
the [Vault Active Directory documentation](/docs/secrets/ad/index.html).

This documentation assumes the AD secrets engine is enabled at the `/ad` path
in Vault. Since it is possible to enable secrets engines at any location, please
update your API calls accordingly.

## Configuration

The `config` endpoint configures the LDAP connection and binding parameters, as well as the password rotation configuration.

### Password parameters

* `ttl` (string, optional) - The default password time-to-live in seconds. Once the ttl has passed, a password will be rotated the next time it's requested.
* `max_ttl` (string, optional) - The maximum password time-to-live in seconds. No role will be allowed to set a custom ttl greater than the `max_ttl`.
* `length` (string, optional) - The desired password length. Defaults to 64. Minimum is 14.
* `formatter` (string, optional) - Text into which the base64 password should be inserted, formatted like so: `mycustom{{PASSWORD}}`.

To meet Microsoft's password complexity requirements, all passwords begin with "?@09AZ" unless a `formatter` is provided. 
The `formatter` is for organizations with different, custom password requirements. It allows an organization to supply
text that fulfills those requirements. `{{PASSWORD}}` must appear exactly once and can be anywhere in the text.

### Connection parameters

* `url` (string, required) - The LDAP server to connect to. Examples: `ldaps://ldap.myorg.com`, `ldaps://ldap.myorg.com:636`. This can also be a comma-delineated list of URLs, e.g. `ldaps://ldap.myorg.com,ldaps://ldap.myorg.com:636`, in which case the servers will be tried in-order if there are errors during the connection process.
* `starttls` (bool, optional) - If true, issues a `StartTLS` command after establishing an unencrypted connection.
* `insecure_tls` - (bool, optional) - If true, skips LDAP server SSL certificate verification - insecure, use with caution!
* `certificate` - (string, optional) - CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded.

### Binding parameters

* `binddn` (string, required) - Distinguished name of object to bind when performing user and group search. Example: `cn=vault,ou=Users,dc=example,dc=com`
* `bindpass` (string, required) - Password to use along with `binddn` when performing user search.
* `userdn` (string, optional) - Base DN under which to perform user search. Example: `ou=Users,dc=example,dc=com`
* `upndomain` (string, optional) - userPrincipalDomain used to construct the UPN string for the authenticating user. The constructed UPN will appear as `[username]@UPNDomain`. Example: `example.com`, which will cause vault to bind as `username@example.com`.

### Other parameters

* `last_rotation_tolerance` (string, optional) - Tolerance duration to use when checking the last rotation time. 
Active Directory often shows a "pwdLastSet" time after Vault's because it takes
a while for password updates to be propagated across a large cluster. By default, if Active Directory's last rotation time is
within 5 seconds of Vault's, Vault considers itself to have been the last entity that rotated the password. However, if it's been 
more than 5 seconds, Vault thinks that something rotated the password out-of-band, and re-rotates it so it will "know" it and be
able to continue returning it. This may be too low for larger Active Directory clusters, and too high for smaller ones.

## Config management

At present, this endpoint does not confirm that the provided AD credentials are
valid AD credentials with proper permissions.

| Method   | Path                   |
| :--------------------- | :--------------------- |
| `POST`   | `/ad/config`           |
| `GET`    | `/ad/config`           |
| `DELETE` | `/ad/config`           |

### Sample Post Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/ad/config
```

### Sample Post Payload

```json
{
  "binddn": "domain-admin",
  "bindpass": "pa$$w0rd",
  "url": "ldaps://127.0.0.11",
  "userdn": "dc=example,dc=com"
}
```

### Sample Get Response Data

```json
{
  "binddn": "domain-admin",
  "certificate": "",
  "insecure_tls": false,
  "length": 64,
  "max_ttl": 2764800,
  "starttls": false,
  "tls_max_version": "tls12",
  "tls_min_version": "tls12",
  "ttl": 2764800,
  "upndomain": "",
  "url": "ldaps://127.0.0.11",
  "userdn": "dc=example,dc=com"
}

```

## Role management

The `roles` endpoint configures how Vault will manage the passwords for individual service accounts.

### Parameters

* `service_account_name` (string, required) - The name of a pre-existing service account in Active Directory that maps to this role.
* `ttl` (string, optional) - The password time-to-live in seconds. Defaults to the configuration `ttl` if not provided.

When adding a role, Vault verifies its associated service account exists.

| Method   | Path                   |
| :--------------------- | :--------------------- |
| `GET`    | `/ad/roles`            |
| `POST`   | `/ad/roles/:role_name` |
| `GET`    | `/ad/roles/:role_name` |
| `DELETE` | `/ad/roles/:role_name` |

### Sample Post Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/ad/roles/my-application
```

### Sample Post Payload

```json
{
  "service_account_name": "my-application@example.com",
  "ttl": 100
}
```

### Sample Get Role Response

```json
{
  "last_vault_rotation": "2018-05-24T17:14:38.677370855Z",
  "password_last_set": "2018-05-24T17:14:38.6038495Z",
  "service_account_name": "my-application@example.com",
  "ttl": 100
}
```

### Sample List Roles Response

Performing a `LIST` on the `/ad/roles` endpoint will list the names of all the roles Vault contains.

```json
[
  "my-application"
]
```

## Retrieving passwords

The `creds` endpoint offers the credential information for a given role.

| Method   | Path                   |
| :--------------------- | :--------------------- |
| `GET`    | `/ad/creds/:role_name` |

### Sample Get Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    http://127.0.0.1:8200/v1/ad/creds/my-application
```

### Sample Get Response

```json
{
  "current_password": "?@09AZnh4Q5N4O5zdLk/4F8aIMgsnpDM6tSQEZCge3Mz1wXcZEgZhOa6OR748F96",
  "last_password": "?@09AZSen9TzUwK7ZhafS7B0GuWGraQjfWEna5SwnmF/tVaKFqjXhhGV/Z0v/pBJ",
  "username": "my-application"
}
```
## Rotate Root Credentials

Rotate the `bindpass` to a new one known only to Vault.

### Risks

1. When the `bindpass` is rotated, it successfully gets rotated in Active Directory but Vault can't store it so it becomes unknown.
2. If the `binddn` in use applies to more than one entity in Active Directory, root credential rotation will fail because it's unclear which entity to perform the operation for.

### Mitigating Risks

1. Always have another account that can provision a new `binddn` and `bindpass` to replace one whose password becomes unknown.
2. Ensure the `binddn` in use only applies to one entity by including all distinguished name parameters possible. For example, use `"CN=vault-ad-test,CN=Users,DC=example,DC=com"` instead of `"CN=vault-ad-test"`.

### Endpoints

| Method   | Path                   |
| :--------------------- | :----------------------------------------- |
| `GET`    | `/ad/rotate-root`      | `204 (empty body) or 200 with warning`     |

Generally, `rotate-root` returns a 204. However, if `rotate-root` is already in progress, it may return a 200 with a warning that root credential rotation is already in progress.

### Sample Get Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    --data @payload.json \
    http://127.0.0.1:8200/v1/ad/rotate-root
```