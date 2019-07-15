---
layout: "api"
page_title: "Userpass - Auth Methods - HTTP API"
sidebar_title: "Username & Password"
sidebar_current: "api-http-auth-userpass"
description: |-
  This is the API documentation for the Vault username and password
  auth method.
---

# Userpass Auth Method (HTTP API)

This is the API documentation for the Vault Username & Password auth method. For
general information about the usage and operation of the Username and Password method, please
see the [Vault Userpass method documentation](/docs/auth/userpass.html).

This documentation assumes the Username & Password method is mounted at the `/auth/userpass`
path in Vault. Since it is possible to enable auth methods at any location,
please update your API calls accordingly.

## Create/Update User

Create a new user or update an existing user. This path honors the distinction between the `create` and `update` capabilities inside ACL policies.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/auth/userpass/users/:username`   |

### Parameters

- `username` `(string: <required>)` – The username for the user.
- `password` `(string: <required>)` - The password for the user. Only required
  when creating the user.
- `policies` `(string: "")` – Comma-separated list of policies. If set to empty
  string, only the `default` policy will be applicable to the user.
- `ttl` `(string: "")` - The lease duration which decides login expiration.
- `max_ttl` `(string: "")` - Maximum duration after which login should expire.
- `bound_cidrs` `(string: "", or list: [])` – If set, restricts usage of the
  login and token to client IPs falling within the range of the specified
  CIDR(s).

### Sample Payload

```json
{
  "password": "superSecretPassword",
  "policies": "admin,default",
  "bound_cidrs": ["127.0.0.1/32", "128.252.0.0/16"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/userpass/users/mitchellh
```

## Read User

Reads the properties of an existing username.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/userpass/users/:username`   |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/userpass/users/mitchellh
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "max_ttl": 0,
    "policies": ["default", "dev"],
    "ttl": 0
  },
  "warnings": null
}
```

## Delete User

This endpoint deletes the user from the method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/auth/userpass/users/:username` |

### Parameters

- `username` `(string: <required>)` - The username for the user.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/userpass/users/mitchellh
```

## Update Password on User

Update password for an existing user.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST` | `/auth/userpass/users/:username/password` |

### Parameters

- `username` `(string: <required>)` – The username for the user.
- `password` `(string: <required>)` - The password for the user.

### Sample Payload

```json
{
  "password": "superSecretPassword2",
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/userpass/users/mitchellh/password
```

## Update Policies on User

Update policies for an existing user.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST` | `/auth/userpass/users/:username/policies` |

### Parameters

- `username` `(string: <required>)` – The username for the user.
- `policies` `(string: "")` – Comma-separated list of policies. If set to empty

### Sample Payload

```json
{
  "policies": ["policy1", "policy2"],
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/userpass/users/mitchellh/policies
```

## List Users

List available userpass users.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/userpass/users`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST
    http://127.0.0.1:8200/v1/auth/userpass/users
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "mitchellh",
      "armon"
    ]
  }
}
```

## Login

Login with the username and password.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST` | `/auth/userpass/login/:username` |

### Parameters

- `username` `(string: <required>)` – The username for the user.
- `password` `(string: <required>)` - The password for the user.

### Sample Payload

```json
{
  "password": "superSecretPassword2",
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/userpass/login/mitchellh
```

### Sample Response

```json
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "warnings": null,
  "auth": {
    "client_token": "64d2a8f2-2a2f-5688-102b-e6088b76e344",
    "accessor": "18bb8f89-826a-56ee-c65b-1736dc5ea27d",
    "policies": ["default"],
    "metadata": {
      "username": "mitchellh"
    },
    "lease_duration": 7200,
    "renewable": true
  }
}
```
