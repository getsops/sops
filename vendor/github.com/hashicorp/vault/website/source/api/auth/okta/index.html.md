---
layout: "api"
page_title: "Okta - Auth Methods - HTTP API"
sidebar_title: "Okta"
sidebar_current: "api-http-auth-okta"
description: |-
  This is the API documentation for the Vault Okta auth method.
---

# Okta Auth Method (API)

This is the API documentation for the Vault Okta auth method. For
general information about the usage and operation of the Okta method, please
see the [Vault Okta method documentation](/docs/auth/okta.html).

This documentation assumes the Okta method is mounted at the `/auth/okta`
path in Vault. Since it is possible to enable auth methods at any location,
please update your API calls accordingly.

## Create Configuration

Configures the connection parameters for Okta. This path honors the
distinction between the `create` and `update` capabilities inside ACL policies.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/okta/config`          |

### Parameters

- `org_name` `(string: <required>)` - Name of the organization to be used in the
  Okta API.
- `api_token` `(string: "")` - Okta API token. This is required to query Okta
  for user group membership. If this is not supplied only locally configured
  groups will be enabled.
- `base_url` `(string: "")` -  If set, will be used as the base domain
  for API requests.  Examples are okta.com, oktapreview.com, and okta-emea.com.
- `ttl` `(string: "")` - Duration after which authentication will be expired.
- `max_ttl` `(string: "")` - Maximum duration after which authentication will
  be expired.
- `bypass_okta_mfa` `(bool: false)` - Whether to bypass an Okta MFA request.
  Useful if using one of Vault's built-in MFA mechanisms, but this will also
  cause certain other statuses to be ignored, such as `PASSWORD_EXPIRED`.

### Sample Payload

```json
{
  "org_name": "example",
  "api_token": "abc123"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/okta/config
```

## Read Configuration

Reads the Okta configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/okta/config`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/okta/config
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "org_name": "example",
    "api_token": "abc123",
    "base_url": "okta.com",
    "ttl": "",
    "max_ttl": ""
  },
  "warnings": null
}
```

## List Users

List the users configured in the Okta method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/okta/users`           |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/auth/okta/users
```

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "fred",
	    "jane"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## Register User

Registers a new user and maps a set of policies to it.  

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/okta/users/:username` |

### Parameters

- `username` `(string: <required>)` - Name of the user.
- `groups` `(array: [])` - List or comma-separated string of groups associated with the user.
- `policies` `(array: [])` - List or comma-separated string of policies associated with the user.

```json
{
  "policies": [
    "dev",
    "prod"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/okta/users/fred
```

## Read User

Reads the properties of an existing username.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`   | `/auth/okta/users/:username` |

### Parameters

- `username` `(string: <required>)` - Username for this user.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/okta/users/test-user
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "policies": [
      "default",
      "dev",
    ],
    "groups": []
  },
  "warnings": null
}
```

## Delete User

Deletes an existing username from the method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`   | `/auth/okta/users/:username` |

### Parameters

- `username` `(string: <required>)` - Username for this user.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/okta/users/test-user
```

## List Groups

List the groups configured in the Okta method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/okta/groups`           |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/auth/okta/groups
```

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "admins",
      "dev-users"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## Register Group

Registers a new group and maps a set of policies to it.  

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/okta/groups/:name` |

### Parameters

- `name` `(string: <required>)` - The name of the group.
- `policies` `(array: [])` - The list or comma-separated string of policies associated with the group.

```json
{
  "policies": [
    "dev",
    "prod"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/okta/groups/admins
```

## Read Group

Reads the properties of an existing group.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`   | `/auth/okta/groups/:name`     |

### Parameters

- `name` `(string: <required>)` - The name for the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/okta/groups/admins
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "policies": [
      "default",
      "admin"
    ]
  },
  "warnings": null
}
```

## Delete Group

Deletes an existing group from the method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`   | `/auth/okta/groups/:name` |

### Parameters

- `name` `(string: <required>)` - The name for the group.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/okta/users/test-user
```

## Login

Login with the username and password.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/okta/login/:username` |

### Parameters

- `username` `(string: <required>)` - Username for this user.
- `password` `(string: <required>)` - Password for the authenticating user.

### Sample Payload

```json
{
  "password": "Password!"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/okta/login/fred
```

### Sample Response

```javascript
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
      "username": "fred",
      "policies": "default"
    },
  },
  "lease_duration": 7200,
  "renewable": true
}
 ```
