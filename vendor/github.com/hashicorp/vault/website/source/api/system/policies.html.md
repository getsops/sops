---
layout: "api"
page_title: "/sys/policies/ - HTTP API"
sidebar_title: "<code>/sys/policies</code>"
sidebar_current: "api-http-system-policies"
description: |-
  The `/sys/policies/` endpoints are used to manage ACL, RGP, and EGP policies in Vault.
---

# `/sys/policies/`

The `/sys/policies` endpoints are used to manage ACL, RGP, and EGP policies in Vault.


~> **NOTE**: This endpoint is only available in Vault version 0.9+. Please also note that RGPs and EGPs are Vault Enterprise Premium features and the associated endpoints are not available in Vault Open Source or Vault Enterprise Pro.

## List ACL Policies

This endpoint lists all configured ACL policies.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/sys/policies/acl`          |

### Sample Request

```
$ curl \
    -X LIST --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policies/acl
```

### Sample Response

```json
{
  "keys": ["root", "my-policy"]
}
```

## Read ACL Policy

This endpoint retrieves information about the named ACL policy.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/policies/acl/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to retrieve.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policies/acl/my-policy
```

### Sample Response

```json
{
  "name": "deploy",
  "policy": "path \"secret/foo\" {..."
}
```

## Create/Update ACL Policy

This endpoint adds a new or updates an existing ACL policy. Once a policy is
updated, it takes effect immediately to all associated users.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/policies/acl/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to create.
  This is specified as part of the request URL.

- `policy` `(string: <required>)` - Specifies the policy document. This can be
  base64-encoded to avoid string escaping.

### Sample Payload

```json
{
  "policy": "path \"secret/foo\" {..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/policies/acl/my-policy
```

## Delete ACL Policy

This endpoint deletes the ACL policy with the given name. This will immediately
affect all users associated with this policy. (A deleted policy set on a token
acts as an empty policy.)

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/policies/acl/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to delete.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/policies/acl/my-policy
```

## List RGP Policies

This endpoint lists all configured RGP policies.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/sys/policies/rgp`          |

### Sample Request

```
$ curl \
    -X LIST --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policies/rgp
```

### Sample Response

```json
{
  "keys": ["webapp", "database"]
}
```

## Read RGP Policy

This endpoint retrieves information about the named RGP policy.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/policies/rgp/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to retrieve.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policies/rgp/webapp
```

### Sample Response

```json
{
  "name": "webapp",
  "policy": "rule main = {...",
  "enforcement_level": "soft-mandatory"
}
```

## Create/Update RGP Policy

This endpoint adds a new or updates an existing RGP policy. Once a policy is
updated, it takes effect immediately to all associated users.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/policies/rgp/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to create.
  This is specified as part of the request URL.

- `policy` `(string: <required>)` - Specifies the policy document. This can be
  base64-encoded to avoid string escaping.

- `enforcement_level` `(string: <required>)` - Specifies the enforcement level
  to use. This must be one of `advisory`, `soft-mandatory`, or
  `hard-mandatory`.

### Sample Payload

```json
{
  "policy": "rule main = {...",
  "enforcement_level": "soft-mandatory"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/policies/rgp/webapp
```

## Delete RGP Policy

This endpoint deletes the RGP policy with the given name. This will immediately
affect all users associated with this policy. (A deleted policy set on a token
acts as an empty policy.)

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/policies/rgp/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to delete.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/policies/rgp/webapp
```

## List EGP Policies

This endpoint lists all configured EGP policies. Since EGP policies act on a
path, this endpoint returns two identifiers:

 * `keys` contains a mapping of names to associated paths in a format that
   `vault list` understands
 * `name_path_map` contains an object mapping names to paths and glob status in
   a more machine-friendly format

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/sys/policies/egp`          |

### Sample Request

```
$ curl \
    -X LIST --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policies/egp
```

### Sample Response

```json
{
  "keys": [ "breakglass" ]
}
```

## Read EGP Policy

This endpoint retrieves information about the named EGP policy.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/policies/egp/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to retrieve.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/policies/egp/breakglass
```

### Sample Response

```json
{
  "enforcement_level": "soft-mandatory",
  "name": "breakglass",
  "paths": [ "*" ],
  "policy": "rule main = {..."
}
```

## Create/Update EGP Policy

This endpoint adds a new or updates an existing EGP policy. Once a policy is
updated, it takes effect immediately to all associated users.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/policies/egp/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to create.
  This is specified as part of the request URL.

- `policy` `(string: <required>)` - Specifies the policy document. This can be
  base64-encoded to avoid string escaping.

- `enforcement_level` `(string: <required>)` - Specifies the enforcement level
  to use. This must be one of `advisory`, `soft-mandatory`, or
  `hard-mandatory`.

- `paths` `(string or array: required)` - Specifies the paths on which this EGP
  should be applied, either as a comma-separated list or an array. Glob
  characters can denote suffixes, e.g. `secret/*`; a path of `*` will affect
  all authenticated and login requests.

### Sample Payload

```json
{
  "policy": "rule main = {...",
  "paths": [ "*", "secret/*", "transit/keys/*" ],
  "enforcement_level": "soft-mandatory"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/policies/egp/breakglass
```

## Delete EGP Policy

This endpoint deletes the EGP policy with the given name from all paths on which it was configured.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/policies/egp/:name`    |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the policy to delete.
  This is specified as part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/policies/egp/breakglass
```
