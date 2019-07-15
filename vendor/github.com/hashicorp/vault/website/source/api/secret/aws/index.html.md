---
layout: "api"
page_title: "AWS - Secrets Engines - HTTP API"
sidebar_title: "AWS"
sidebar_current: "api-http-secret-aws"
description: |-
  This is the API documentation for the Vault AWS secrets engine.
---

# AWS Secrets Engine (API)

This is the API documentation for the Vault AWS secrets engine. For general
information about the usage and operation of the AWS secrets engine, please see
the [Vault AWS documentation](/docs/secrets/aws/index.html).

This documentation assumes the AWS secrets engine is enabled at the `/aws` path
in Vault. Since it is possible to enable secrets engines at any location, please
update your API calls accordingly.

## Configure Root IAM Credentials

This endpoint configures the root IAM credentials to communicate with AWS. There
are multiple ways to pass root IAM credentials to the Vault server, specified
below with the highest precedence first. If credentials already exist, this will
overwrite them.

The official AWS SDK is used for sourcing credentials from env vars, shared
files, or IAM/ECS instances.

- Static credentials provided to the API as a payload

- Credentials in the `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, and `AWS_REGION`
  environment variables **on the server**

- Shared credentials files

- Assigned IAM role or ECS task role credentials

At present, this endpoint does not confirm that the provided AWS credentials are
valid AWS credentials with proper permissions.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/aws/config/root`           |

### Parameters

- `max_retries` `(int: -1)` - Number of max retries the client should use for
  recoverable errors. The default (`-1`) falls back to the AWS SDK's default
  behavior.

- `access_key` `(string: <required>)` – Specifies the AWS access key ID.

- `secret_key` `(string: <required>)` – Specifies the AWS secret access key.

- `region` `(string: <optional>)` – Specifies the AWS region. If not set it
  will use the `AWS_REGION` env var, `AWS_DEFAULT_REGION` env var, or
  `us-east-1` in that order.

- `iam_endpoint` `(string: <optional>)` – Specifies a custom HTTP IAM endpoint to use.

- `sts_endpoint` `(string: <optional>)` – Specifies a custom HTTP STS endpoint to use.

### Sample Payload

```json
{
  "access_key": "AKIA...",
  "secret_key": "2J+...",
  "region": "us-east-1"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/aws/config/root
```

## Rotate Root IAM Credentials

When you have configured Vault with static credentials, you can use this
endpoint to have Vault rotate the access key it used. Note that, due to AWS
eventual consistency, after calling this endpoint, subsequent calls from Vault
to AWS may fail for a few seconds until AWS becomes consistent again.


In order to call this endpoint, Vault's AWS access key MUST be the only access
key on the IAM user; otherwise, generation of a new access key will fail. Once
this method is called, Vault will now be the only entity that knows the AWS
secret key is used to access AWS.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/aws/config/rotate-root`    |

### Parameters

There are no parameters to this operation.

### Sample Request

```$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8211/v1/aws/config/rotate-root
```

### Sample Response

```json
{
  "data": {
    "access_key": "AKIA..."
  }
}
```

The new access key Vault uses is returned by this operation.

## Configure Lease

This endpoint configures lease settings for the AWS secrets engine. It is
optional, as there are default values for `lease` and `lease_max`.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/aws/config/lease`          |

### Parameters

- `lease` `(string: <required>)` – Specifies the lease value provided as a
  string duration with time suffix. "h" (hour) is the largest suffix.

- `lease_max` `(string: <required>)` – Specifies the maximum lease value
  provided as a string duration with time suffix. "h" (hour) is the largest
  suffix.

### Sample Payload

```json
{
  "lease": "30m",
  "lease_max": "12h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/aws/config/lease
```

## Read Lease

This endpoint returns the current lease settings for the AWS secrets engine.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/aws/config/lease`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/aws/config/lease
```

### Sample Response

```json
{
  "data": {
    "lease": "30m0s",
    "lease_max": "12h0m0s"
  }
}
```

## Create/Update Role

This endpoint creates or updates the role with the given `name`. If a role with
the name does not exist, it will be created. If the role exists, it will be
updated with the new attributes.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/aws/roles/:name`           |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create. This
  is part of the request URL.

- `credential_type` `(string: <required>)` – Specifies the type of credential to be used when
  retrieving credentials from the role. Must be one of `iam_user`,
  `assumed_role`, or `federation_token`.

- `role_arns` `(list: [])` – Specifies the ARNs of the AWS roles this Vault role
  is allowed to assume. Required when `credential_type` is `assumed_role` and
  prohibited otherwise. This is a comma-separated string or JSON array.

- `policy_arns` `(list: [])` – Specifies the ARNs of the AWS managed policies to
  be attached to IAM users when they are requested. Valid only when
  `credential_type` is `iam_user`. When `credential_type` is `iam_user`, at
  least one of `policy_arns` or `policy_document` must be specified. This is a
  comma-separated string or JSON array.

- `policy_document` `(string)` – The IAM policy document for the role. The
  behavior depends on the credential type. With `iam_user`, the policy document
  will be attached to the IAM user generated and augment the permissions the IAM
  user has. With `assumed_role` and `federation_token`, the policy document will
  act as a filter on what the credentials can do.

- `default_sts_ttl` `(string)` - The default TTL for STS credentials. When a TTL is not
  specified when STS credentials are requested, and a default TTL is specified
  on the role, then this default TTL will be used. Valid only when
  `credential_type` is one of `assumed_role` or `federation_token`.

- `max_sts_ttl` `(string)` - The max allowed TTL for STS credentials (credentials
  TTL are capped to `max_sts_ttl`). Valid only when `credential_type` is one of 
  `assumed_role` or `federation_token`.

- `user_path` `(string)` - The path for the user name. Valid only when
  `credential_type` is `iam_user`. Default is `/`

Legacy parameters:

These parameters are supported for backwards compatibility only. They cannot be
mixed with the parameters listed above.

- `policy` `(string: <required unless arn provided>)` – Specifies the IAM policy
  in JSON format.

- `arn` `(string: <required unless policy provided>)` – Specifies the full ARN
  reference to the desired existing policy.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/aws/roles/example-role
```

### Sample Payloads

Using an inline IAM policy:

```json
{
  "credential_type": "federation_token",
  "policy_document": "{\"Version\": \"...\"}"
}
```

Using an ARN:

```json
{
  "credential_type": "assumed_role",
  "role_arns": "arn:aws:iam::123456789012:role/DeveloperRole"
}
```

## Read Role

This endpoint queries an existing role by the given name. If the role does not
exist, a 404 is returned.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/aws/roles/:name`           |

If invalid role data was supplied to the role from an earlier version of Vault,
then it will show up in the response as `invalid_data`.

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/aws/roles/example-role
```

### Sample Responses

For an inline IAM policy:

```json
{
  "data": {
    "policy_document": "{\"Version\": \"...\"}",
    "policy_arns": [],
    "credential_types": ["assumed_role"],
    "role_arns": [],
  }
}
```

For a role ARN:

```json
{
  "data": {
    "policy_document": "",
    "policy_arns": [],
    "credential_types": ["assumed_role"],
    "role_arns": ["arn:aws:iam::123456789012:role/example-role"]
  }
}
```

## List Roles

This endpoint lists all existing roles in the secrets engine.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/aws/roles`                 |

### Sample Request

```
$ curl
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/aws/roles
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "example-role"
    ]
  }
}
```

## Delete Role

This endpoint deletes an existing role by the given name. If the role does not
exist, a 404 is returned.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`  | `/aws/roles/:name`           |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/aws/roles/example-role
```

## Generate Credentials

This endpoint generates credentials based on the named role. This role must be
created before queried.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/aws/creds/:name`           |
| `GET`    | `/aws/sts/:name`             |

The `/aws/creds` and `/aws/sts` endpoints are almost identical. The exception is
when retrieving credentials for a role that was specified with the legacy `arn`
or `policy` parameter. In this case, credentials retrieved through `/aws/sts`
must be of either the `assumed_role` or `federation_token` types, and
credentials retrieved through `/aws/creds` must be of the `iam_user` type.

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to generate
  credentials against. This is part of the request URL.
- `role_arn` `(string)` – The ARN of the role to assume if `credential_type` on
  the Vault role is `assumed_role`. Must match one of the allowed role ARNs in
  the Vault role. Optional if the Vault role only allows a single AWS role ARN;
  required otherwise.
- `ttl` `(string: "3600s")` – Specifies the TTL for the use of the STS token.
  This is specified as a string with a duration suffix. Valid only when
  `credential_type` is `assumed_role` or `federation_token`. When not specified,
  the `default_sts_ttl` set for the role will be used. If that is also not set, then
  the default value of `3600s` will be used. AWS places limits
  on the maximum TTL allowed. See the AWS documentation on the `DurationSeconds`
  parameter for
  [AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)
  (for `assumed_role` credential types) and
  [GetFederationToken](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html)
  (for `federation_token` credential types) for more details.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/aws/creds/example-role
```

### Sample Response

```json
{
  "data": {
    "access_key": "AKIA...",
    "secret_key": "xlCs...",
    "security_token": null
  }
}
```
