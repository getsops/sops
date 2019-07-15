---
layout: "api"
page_title: "LDAP - Auth Methods - HTTP API"
sidebar_title: "LDAP"
sidebar_current: "api-http-auth-ldap"
description: |-
  This is the API documentation for the Vault LDAP auth method.
---

# LDAP Auth Method (API)

This is the API documentation for the Vault LDAP auth method. For
general information about the usage and operation of the LDAP method, please
see the [Vault LDAP method documentation](/docs/auth/ldap.html).

This documentation assumes the LDAP method is mounted at the `/auth/ldap`
path in Vault. Since it is possible to enable auth methods at any location,
please update your API calls accordingly.

## Configure LDAP

This endpoint configures the LDAP auth method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/auth/ldap/config`         |

### Parameters

- `url` `(string: <required>)` – The LDAP server to connect to. Examples:
  `ldap://ldap.myorg.com`, `ldaps://ldap.myorg.com:636`. Multiple URLs can be
  specified with commas, e.g. `ldap://ldap.myorg.com,ldap://ldap2.myorg.com`;
  these will be tried in-order.
- `case_sensitive_names` `(bool: false)` – If set, user and group names
  assigned to policies within the backend will be case sensitive. Otherwise,
  names will be normalized to lower case. Case will still be preserved when
  sending the username to the LDAP server at login time; this is only for
  matching local user/group definitions.
- `starttls` `(bool: false)` – If true, issues a `StartTLS` command after
  establishing an unencrypted connection.
- `tls_min_version` `(string: tls12)` – Minimum TLS version to use. Accepted
  values are `tls10`, `tls11` or `tls12`.
- `tls_max_version` `(string: tls12)` – Maximum TLS version to use. Accepted
  values are `tls10`, `tls11` or `tls12`.
- `insecure_tls` `(bool: false)` – If true, skips LDAP server SSL certificate
  verification - insecure, use with caution!
- `certificate` `(string: "")` – CA certificate to use when verifying LDAP server
  certificate, must be x509 PEM encoded.
- `binddn` `(string: "")` – Distinguished name of object to bind when performing
  user search.  Example: `cn=vault,ou=Users,dc=example,dc=com`
- `bindpass` `(string: "")` – Password to use along with `binddn` when performing
  user search.
- `userdn` `(string: "")` – Base DN under which to perform user search. Example:
  `ou=Users,dc=example,dc=com`
- `userattr` `(string: "")` – Attribute on user attribute object matching the
  username passed when authenticating. Examples: `sAMAccountName`, `cn`, `uid`
- `discoverdn` `(bool: false)` – Use anonymous bind to discover the bind DN of a
  user.
- `deny_null_bind` `(bool: true)` – This option prevents users from bypassing
  authentication when providing an empty password.
- `upndomain` `(string: "")` – The userPrincipalDomain used to construct the UPN
  string for the authenticating user. The constructed UPN will appear as
  `[username]@UPNDomain`. Example: `example.com`, which will cause vault to bind
  as `username@example.com`.
- `groupfilter` `(string: "")` – Go template used when constructing the group
  membership query. The template can access the following context variables:
  \[`UserDN`, `Username`\]. The default is
  `(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))`,
  which is compatible with several common directory schemas. To support
  nested group resolution for Active Directory, instead use the following
  query: `(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))`.
- `groupdn` `(string: "")` – LDAP search base to use for group membership
  search. This can be the root containing either groups or users.  Example:
  `ou=Groups,dc=example,dc=com`
- `groupattr` `(string: "")` – LDAP attribute to follow on objects returned by
  `groupfilter` in order to enumerate user group membership. Examples: for
  groupfilter queries returning _group_ objects, use: `cn`. For queries
  returning _user_ objects, use: `memberOf`. The default is `cn`.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/ldap/config
```

### Sample Payload

```json
{
  "binddn": "cn=vault,ou=Users,dc=example,dc=com",
  "deny_null_bind": true,
  "discoverdn": false,
  "groupattr": "cn",
  "groupdn": "ou=Groups,dc=example,dc=com",
  "groupfilter": "(\u0026(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))",
  "insecure_tls": false,
  "starttls": false,
  "tls_max_version": "tls12",
  "tls_min_version": "tls12",
  "url": "ldaps://ldap.myorg.com:636",
  "userattr": "samaccountname",
  "userdn": "ou=Users,dc=example,dc=com"
}
```

## Read LDAP Configuration

This endpoint retrieves the LDAP configuration for the auth method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/ldap/config`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/ldap/config
```

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "binddn": "cn=vault,ou=Users,dc=example,dc=com",
    "bindpass": "",
    "certificate": "",
    "deny_null_bind": true,
    "discoverdn": false,
    "groupattr": "cn",
    "groupdn": "ou=Groups,dc=example,dc=com",
    "groupfilter": "(\u0026(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))",
    "insecure_tls": false,
    "starttls": false,
    "tls_max_version": "tls12",
    "tls_min_version": "tls12",
    "upndomain": "",
    "url": "ldaps://ldap.myorg.com:636",
    "userattr": "samaccountname",
    "userdn": "ou=Users,dc=example,dc=com"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## List LDAP Groups

This endpoint returns a list of existing groups in the method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/ldap/groups`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/auth/ldap/groups
```

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "scientists",
      "engineers"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## Read LDAP Group

This endpoint returns the policies associated with a LDAP group.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/ldap/groups/:name`     |

### Parameters

- `name` `(string: <required>)` – The name of the LDAP group

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/ldap/groups/admins
```

### Sample Response

```json
{
  "data": {
    "policies": [
      "admin",
      "default"
    ]
  },
  "renewable": false,
  "lease_id": ""
  "lease_duration": 0,
  "warnings": null
}
```

## Create/Update LDAP Group

This endpoint creates or updates LDAP group policies.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/auth/ldap/groups/:name`   |

### Parameters

- `name` `(string: <required>)` – The name of the LDAP group
- `policies` `(string: "")` – Comma-separated list of policies associated to the
  group.

### Sample Payload

```json
{
  "policies": "admin,default"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/ldap/groups/admins
```

## Delete LDAP Group

This endpoint deletes the LDAP group and policy association.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/auth/ldap/groups/:name`    |

### Parameters

- `name` `(string: <required>)` – The name of the LDAP group

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/ldap/groups/admins
```

## List LDAP Users

This endpoint returns a list of existing users in the method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/ldap/users`           |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/auth/ldap/users
```

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "mitchellh",
      "armon"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## Read LDAP User

This endpoint returns the policies associated with a LDAP user.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/ldap/users/:username` |

### Parameters

- `username` `(string: <required>)` – The username of the LDAP user

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/ldap/users/mitchellh
```

### Sample Response

```json
{
  "data": {
    "policies": [
      "admin",
      "default"
    ],
    "groups": ""
  },
  "renewable": false,
  "lease_id": ""
  "lease_duration": 0,
  "warnings": null
}
```

## Create/Update LDAP User

This endpoint creates or updates LDAP users policies and group associations.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/auth/ldap/users/:username`   |

### Parameters

- `username` `(string: <required>)` – The username of the LDAP user
- `policies` `(string: "")` – Comma-separated list of policies associated to the
  user.
- `groups` `(string: "")` – Comma-separated list of groups associated to the
  user.

### Sample Payload

```json
{
  "policies": "admin,default"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/ldap/users/mitchellh
```

## Delete LDAP User

This endpoint deletes the LDAP user and policy association.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/auth/ldap/users/:username` |

### Parameters

- `username` `(string: <required>)` – The username of the LDAP user

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/ldap/users/mitchellh
```

## Login with LDAP User

This endpoint allows you to log in with LDAP credentials

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/ldap/login/:username` |

### Parameters

- `username` `(string: <required>)` – The username of the LDAP user
- `password` `(string: <required>)` – The password for the LDAP user

### Sample Payload

```json
{
  "password": "MyPassword1"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/ldap/login/mitchellh
```

### Sample Response

```json
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "auth": {
    "client_token": "c4f280f6-fdb2-18eb-89d3-589e2e834cdb",
    "policies": [
      "admins",
      "default"
    ],
    "metadata": {
      "username": "mitchellh"
    },
    "lease_duration": 0,
    "renewable": false
  }
}
```
