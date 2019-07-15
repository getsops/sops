---
layout: "api"
page_title: "Kubernetes - Auth Methods - HTTP API"
sidebar_title: "Kubernetes"
sidebar_current: "api-http-auth-kubernetes"
description: |-
  This is the API documentation for the Vault Kubernetes auth method plugin.
---

# Kubernetes Auth Method (API)

This is the API documentation for the Vault Kubernetes auth method plugin. To
learn more about the usage and operation, see the
[Vault Kubernetes auth method](/docs/auth/kubernetes.html).

This documentation assumes the Kubernetes method is mounted at the
`/auth/kubernetes` path in Vault. Since it is possible to enable auth methods at
any location, please update your API calls accordingly.

## Configure Method

The Kubernetes auth method validates service account JWTs and verifies their
existence with the Kubernetes TokenReview API. This endpoint configures the
public key used to validate the JWT signature and the necessary information to
access the Kubernetes API.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/kubernetes/config`    |

### Parameters
 - `kubernetes_host` `(string: <required>)` - Host must be a host string, a host:port pair, or a URL to the base of the Kubernetes API server.
 - `kubernetes_ca_cert` `(string: "")` - PEM encoded CA cert for use by the TLS client used to talk with the Kubernetes API. NOTE: Every line must end with a newline: \n
 - `token_reviewer_jwt` `(string: "")` - A service account JWT used to access the TokenReview
    API to validate other JWTs during login. If not set
    the JWT used for login will be used to access the API.
 - `pem_keys` `(array: [])` - Optional list of PEM-formatted public keys or certificates
    used to verify the signatures of Kubernetes service account
    JWTs. If a certificate is given, its public key will be
    extracted. Not every installation of Kubernetes exposes these
    keys.

### Sample Payload

```json
{
  "kubernetes_host": "https://192.168.99.100:8443",
  "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----\n.....\n-----END CERTIFICATE-----",
  "pem_keys": "-----BEGIN CERTIFICATE-----\n.....\n-----END CERTIFICATE-----"
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/kubernetes/config
```

## Read Config

Returns the previously configured config, including credentials.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/kubernetes/config`    |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/kubernetes/config
```

### Sample Response

```json
{
  "data":{
    "kubernetes_host": "https://192.168.99.100:8443",
    "kubernetes_ca_cert": "-----BEGIN CERTIFICATE-----.....-----END CERTIFICATE-----",
    "pem_keys": ["-----BEGIN CERTIFICATE-----.....", .....]
  }
}
```

## Create Role

Registers a role in the auth method. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/kubernetes/role/:name`|

### Parameters
- `name` `(string: <required>)` - Name of the role.
- `bound_service_account_names` `(array: <required>)` - List of service account
  names able to access this role. If set to "\*" all names are allowed, both this
  and bound_service_account_namespaces can not be "\*".
- `bound_service_account_namespaces` `(array: <required>)` - List of namespaces
  allowed to access this role. If set to "\*" all namespaces are allowed, both
  this and bound_service_account_names can not be set to "\*".
- `ttl` `(string: "")` - The TTL period of tokens issued using this role in
  seconds.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens
  issued in seconds using this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.
- `policies` `(array: [])` - Policies to be set on tokens issued using this
  role.

### Sample Payload

```json
{
  "bound_service_account_names": "vault-auth",
  "bound_service_account_namespaces": "default",
  "policies": [
    "dev",
    "prod"
  ],
  "max_ttl": 1800000,
}
```

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/kubernetes/role/dev-role
```
## Read Role

Returns the previously registered role configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`   | `/auth/kubernetes/role/:name` |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/kubernetes/role/dev-role
```

### Sample Response

```json
{
  "data":{
    "bound_service_account_names": "vault-auth",
    "bound_service_account_namespaces": "default",
    "max_ttl": 1800000,
    "ttl":0,
    "period": 0,
    "policies":[
      "dev",
      "prod"
    ]
  }
}
```

## List Roles

Lists all the roles that are registered with the auth method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/kubernetes/role`            |
| `GET`   | `/auth/kubernetes/role?list=true`   |

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/auth/kubernetes/role
```

### Sample Response

```json  
{
  "data": {
    "keys": [
      "dev-role",
      "prod-role"
    ]
  }
}
```

## Delete Role

Deletes the previously registered role.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/auth/kubernetes/role/:role`|

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```text
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/kubernetes/role/dev-role
```

## Login

Fetch a token. This endpoint takes a signed JSON Web Token (JWT) and
a role name for some entity. It verifies the JWT signature to authenticate that
entity and then authorizes the entity for the given role.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/kubernetes/login`            |

### Sample Payload

- `role` `(string: <required>)` - Name of the role against which the login is being
  attempted.
- `jwt` `(string: <required>)` - Signed [JSON Web
  Token](https://tools.ietf.org/html/rfc7519) (JWT) for authenticating a service
  account.

### Sample Payload

```json
{
  "role": "dev-role",
  "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Sample Request

```text
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/kubernetes/login
```

### Sample Response

```json
{
  "auth": {
    "client_token": "62b858f9-529c-6b26-e0b8-0457b6aacdb4",
    "accessor": "afa306d0-be3d-c8d2-b0d7-2676e1c0d9b4",
    "policies": [
      "default"
    ],
    "metadata": {
      "role": "test",
      "service_account_name": "vault-auth",
      "service_account_namespace": "default",
      "service_account_secret_name": "vault-auth-token-pd21c",
      "service_account_uid": "aa9aa8ff-98d0-11e7-9bb7-0800276d99bf"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
```
