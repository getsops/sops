---
layout: "api"
page_title: "AliCloud - Secrets Engines - HTTP API"
sidebar_title: "AliCloud"
sidebar_current: "docs-http-secret-alicloud"
description: |-
  This is the API documentation for the Vault AliCloud secrets engine.
---

# AliCloud Secrets Engine (API)

This is the API documentation for the Vault AliCloud secrets engine. For general
information about the usage and operation of the AliCloud secrets engine, please see
the [Vault AliCloud documentation](/docs/secrets/alicloud/index.html).

This documentation assumes the AliCloud secrets engine is enabled at the `/alicloud` path
in Vault. Since it is possible to enable secrets engines at any location, please
update your API calls accordingly.

## Config management

This endpoint configures the root RAM credentials to communicate with AliCloud. AliCloud
will use credentials in the following order:

1. [Environment variables](https://github.com/aliyun/alibaba-cloud-sdk-go/blob/master/sdk/auth/credentials/providers/env.go)
2. A static credential configuration set at this endpoint
3. Instance metadata (recommended)

To use instance metadata, leave the static credential configuration unset.

At present, this endpoint does not confirm that the provided AliCloud credentials are
valid AliCloud credentials with proper permissions.

Please see the [Vault AliCloud documentation](/docs/secrets/alicloud/index.html) for
the policies that should be attached to the access key you provide.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/alicloud/config`           |
| `GET`    | `/alicloud/config`           |

### Parameters

* `access_key` (string, required) - The ID of an access key with appropriate policies.
* `secret_key` (string, required) - The secret for that key.

### Sample Post Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/alicloud/config
```

### Sample Post Payload

```json
{
  "access_key": "0wNEpMMlzy7szvai",
  "secret_key": "PupkTg8jdmau1cXxYacgE736PJj4cA"
}
```

### Sample Get Response Data

```json
{
    "access_key": "0wNEpMMlzy7szvai"
}
```

## Role management

The `role` endpoint configures how Vault will generate credentials for users of each role.

### Parameters

* `name` (string, required) – Specifies the name of the role to generate credentials against. This is part of the request URL.
* `remote_policies` (string, optional) - The names and types of a pre-existing policies to be applied to the generate access token. Example: "name:AliyunOSSReadOnlyAccess,type:System".
* `inline_policies` (string, optional) - The policy document JSON to be generated and attached to the access token.
* `role_arn` (string, optional) - The ARN of a role that will be assumed to obtain STS credentials. See [Vault AliCloud documentation](/docs/secrets/alicloud/index.html) regarding trusted actors.
* `ttl` (int, optional) - The duration in seconds after which the issued token should expire. Defaults to 0, in which case the value will fallback to the system/mount defaults.
* `max_ttl` (int, optional) - The maximum allowed lifetime of tokens issued using this role.

| Method   | Path                        |
| :---------------------------| :--------------------- |
| `GET`    | `/alicloud/role`            |
| `POST`   | `/alicloud/role/:role_name` |
| `GET`    | `/alicloud/role/:role_name` |
| `DELETE` | `/alicloud/role/:role_name` |

### Sample Post Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/alicloud/role/my-application
```

### Sample Post Payload Using Policies

```json
{
  "remote_policies": [
    "name:AliyunOSSReadOnlyAccess,type:System",
    "name:AliyunRDSReadOnlyAccess,type:System"
  ],
  "inline_policies": "[{\"Statement\": [{\"Action\": [\"ram:Get*\",\"ram:List*\"],\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}]"
}
```

### Sample Get Role Response Using Policies

```json
{
	"inline_policies": [{
		"hash": "49796debb24d39b7a61485f9b0c97e04",
		"policy_document": {
			"Statement": [{
				"Action": ["ram:Get*", "ram:List*"],
				"Effect": "Allow",
				"Resource": "*"
			}],
			"Version": "1"
		}
	}],
	"max_ttl": 0,
	"remote_policies": [{
		"name": "AliyunOSSReadOnlyAccess",
		"type": "System"
	}, {
		"name": "AliyunRDSReadOnlyAccess",
		"type": "System"
	}],
	"role_arn": "",
	"ttl": 0
}
```

### Sample Post Payload Using Assume-Role

```json
{
  "role_arn": "acs:ram::5138828231865461:role/hastrustedactors"
}
```

### Sample Get Role Response Using Assume-Role

```json
{
	"inline_policies": null,
	"max_ttl": 0,
	"remote_policies": null,
	"role_arn": "acs:ram::5138828231865461:role/hastrustedactors",
	"ttl": 0
}
```

### Sample List Roles Response

Performing a `LIST` on the `/alicloud/roles` endpoint will list the names of all the roles Vault contains.

```json
[
  "policy-based",
  "role-based"
]
```

## Generate RAM Credentials

This endpoint generates dynamic RAM credentials based on the named role. This
role must be created before queried.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/alicloud/creds/:name`      |

### Parameters

* `name` (string, required) – Specifies the name of the role to generate credentials against. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/alicloud/creds/example-role
```

### Sample Response for Roles Using Policies

```json
{
  "access_key": "0wNEpMMlzy7szvai",
  "secret_key": "PupkTg8jdmau1cXxYacgE736PJj4cA"
}

```

### Sample Response for Roles Using Assume-Role

```json
{
	"access_key": "STS.L4aBSCSJVMuKg5U1vFDw",
	"expiration": "2018-08-15T22:04:07Z",
	"secret_key": "wyLTSmsyPGP1ohvvw8xYgB29dlGI8KMiH2pKCNZ9",
	"security_token": "CAESrAIIARKAAShQquMnLIlbvEcIxO6wCoqJufs8sWwieUxu45hS9AvKNEte8KRUWiJWJ6Y+YHAPgNwi7yfRecMFydL2uPOgBI7LDio0RkbYLmJfIxHM2nGBPdml7kYEOXmJp2aDhbvvwVYIyt/8iES/R6N208wQh0Pk2bu+/9dvalp6wOHF4gkFGhhTVFMuTDRhQlNDU0pWTXVLZzVVMXZGRHciBTQzMjc0KgVhbGljZTCpnJjwySk6BlJzYU1ENUJuCgExGmkKBUFsbG93Eh8KDEFjdGlvbkVxdWFscxIGQWN0aW9uGgcKBW9zczoqEj8KDlJlc291cmNlRXF1YWxzEghSZXNvdXJjZRojCiFhY3M6b3NzOio6NDMyNzQ6c2FtcGxlYm94L2FsaWNlLyo="
}
```
