---
layout: "api"
page_title: "/sys/audit-hash - HTTP API"
sidebar_title: "<code>/sys/audit-hash</code>"
sidebar_current: "api-http-system-audit-hash"
description: |-
  The `/sys/audit-hash` endpoint is used to hash data using an audit device's
  hash function and salt.
---

# `/sys/audit-hash`

The `/sys/audit-hash` endpoint is used to calculate the hash of the data used by
an audit device's hash function and salt. This can be used to search audit logs
for a hashed value when the original value is known.

## Calculate Hash

This endpoint hashes the given input data with the specified audit device's
hash function and salt. This endpoint can be used to discover whether a given
plaintext string (the `input` parameter) appears in the audit log in obfuscated
form.

The audit log records requests and responses. Since the Vault API is JSON-based,
any binary data returned from an API call (such as a DER-format certificate) is
base64-encoded by the Vault server in the response. As a result such information
should also be base64-encoded to supply into the `input` parameter.

| Method | Path                    |
| :---------------------- | :----------------- |
| `POST` | `/sys/audit-hash/:path` |

### Parameters

- `path` `(string: <required>)` – Specifies the path of the audit device to
  generate hashes for. This is part of the request URL.

- `input` `(string: <required>)` – Specifies the input string to hash.

### Sample Payload

```json
{
  "input": "my-secret-vault"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/audit-hash/example-audit
```

### Sample Response

```json
{
  "hash": "hmac-sha256:08ba35..."
}
```
