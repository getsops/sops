---
layout: "api"
page_title: "Transit - Secrets Engines - HTTP API"
sidebar_title: "Transit"
sidebar_current: "api-http-secret-transit"
description: |-
  This is the API documentation for the Vault Transit secrets engine.
---

# Transit Secrets Engine (API)

This is the API documentation for the Vault Transit secrets engine. For general
information about the usage and operation of the Transit secrets engine, please
see the [transit documentation](/docs/secrets/transit/index.html).

This documentation assumes the transit secrets engine is enabled at the
`/transit` path in Vault. Since it is possible to enable secrets engines at any
location, please update your API calls accordingly.

## Create Key

This endpoint creates a new named encryption key of the specified type. The
values set here cannot be changed after key creation.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name`        |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  create. This is specified as part of the URL.

- `convergent_encryption` `(bool: false)` – If enabled, the key will support
  convergent encryption, where the same plaintext creates the same ciphertext.
  This requires _derived_ to be set to `true`. When enabled, each
  encryption(/decryption/rewrap/datakey) operation will derive a `nonce` value
  rather than randomly generate it.

- `derived` `(bool: false)` – Specifies if key derivation is to be used. If
  enabled, all encrypt/decrypt requests to this named key must provide a context
  which is used for key derivation.

- `exportable` `(bool: false)` -  Enables keys to be exportable. This
  allows for all the valid keys in the key ring to be exported. Once set, this
  cannot be disabled.

- `allow_plaintext_backup` `(bool: false)` - If set, enables taking backup of
  named key in the plaintext format. Once set, this cannot be disabled.

- `type` `(string: "aes256-gcm96")` – Specifies the type of key to create. The
  currently-supported types are:

    - `aes256-gcm96` – AES-256 wrapped with GCM using a 96-bit nonce size AEAD
      (symmetric, supports derivation and convergent encryption)
    - `chacha20-poly1305` – ChaCha20-Poly1305 AEAD (symmetric, supports
      derivation and convergent encryption)
    - `ed25519` – ED25519 (asymmetric, supports derivation). When using
      derivation, a sign operation with the same context will derive the same
      key and signature; this is a signing analogue to `convergent_encryption`.
    - `ecdsa-p256` – ECDSA using the P-256 elliptic curve (asymmetric)
    - `rsa-2048` - RSA with bit size of 2048 (asymmetric)
    - `rsa-4096` - RSA with bit size of 4096 (asymmetric)

### Sample Payload

```json
{
  "type": "ecdsa-p256",
  "derived": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/keys/my-key
```

## Read Key

This endpoint returns information about a named encryption key. The `keys`
object shows the creation time of each key version; the values are not the keys
themselves. Depending on the type of key, different information may be returned,
e.g. an asymmetric key will return its public key in a standard format for the
type.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/transit/keys/:name`        |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  read. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/transit/keys/my-key
```

### Sample Response

```json
{
  "data": {
    "type": "aes256-gcm96",
    "deletion_allowed": false,
    "derived": false,
    "exportable": false,
    "allow_plaintext_backup": false,
    "keys": {
      "1": 1442851412
    },
    "min_decryption_version": 1,
    "min_encryption_version": 0,
    "name": "foo",
    "supports_encryption": true,
    "supports_decryption": true,
    "supports_derivation": true,
    "supports_signing": false
  }
}
```

## List Keys

This endpoint returns a list of keys. Only the key names are returned (not the
actual keys themselves).

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/transit/keys`              |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/transit/keys
```

### Sample Response

```json
{
  "data": {
    "keys": ["foo", "bar"]
  },
  "lease_duration": 0,
  "lease_id": "",
  "renewable": false
}
```

## Delete Key

This endpoint deletes a named encryption key. It will no longer be possible to
decrypt any data encrypted with the named key. Because this is a potentially
catastrophic operation, the `deletion_allowed` tunable must be set in the key's
`/config` endpoint.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/transit/keys/:name`        |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  delete. This is specified as part of the URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/transit/keys/my-key
```

## Update Key Configuration

This endpoint allows tuning configuration values for a given key. (These values
are returned during a read operation on the named key.)

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name/config` |

### Parameters

- `min_decryption_version` `(int: 0)` – Specifies the minimum version of
  ciphertext allowed to be decrypted. Adjusting this as part of a key rotation
  policy can prevent old copies of ciphertext from being decrypted, should they
  fall into the wrong hands. For signatures, this value controls the minimum
  version of signature that can be verified against. For HMACs, this controls
  the minimum version of a key allowed to be used as the key for verification.

- `min_encryption_version` `(int: 0)` – Specifies the minimum version of the
  key that can be used to encrypt plaintext, sign payloads, or generate HMACs.
  Must be `0` (which will use the latest version) or a value greater or equal
  to `min_decryption_version`.

- `deletion_allowed` `(bool: false)` - Specifies if the key is allowed to be
  deleted.

- `exportable` `(bool: false)` -  Enables keys to be exportable. This
  allows for all the valid keys in the key ring to be exported. Once set, this
  cannot be disabled.

- `allow_plaintext_backup` `(bool: false)` - If set, enables taking backup of
  named key in the plaintext format. Once set, this cannot be disabled.

### Sample Payload

```json
{
  "deletion_allowed": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/keys/my-key/config
```

## Rotate Key

This endpoint rotates the version of the named key. After rotation, new
plaintext requests will be encrypted with the new version of the key. To upgrade
ciphertext to be encrypted with the latest version of the key, use the `rewrap`
endpoint. This is only supported with keys that support encryption and
decryption operations.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name/rotate` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    http://127.0.0.1:8200/v1/transit/keys/my-key/rotate
```

## Export Key

This endpoint returns the named key. The `keys` object shows the value of the
key for each version. If `version` is specified, the specific version will be
returned. If `latest` is provided as the version, the current key will be
provided. Depending on the type of key, different information may be returned.
The key must be exportable to support this operation and the version must still
be valid.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/transit/export/:key_type/:name(/:version)` |

### Parameters

- `key_type` `(string: <required>)` – Specifies the type of the key to export.
  This is specified as part of the URL. Valid values are:

    - `encryption-key`
    - `signing-key`
    - `hmac-key`

- `name` `(string: <required>)` – Specifies the name of the key to read
  information about. This is specified as part of the URL.

- `version` `(string: "")` – Specifies the version of the key to read. If omitted,
  all versions of the key will be returned. This is specified as part of the
  URL. If the version is set to `latest`, the current key will be returned.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/transit/export/encryption-key/my-key/1
```

### Sample Response

```json
{
  "data": {
    "name": "foo",
    "keys": {
      "1": "eyXYGHbTmugUJn6EtYD/yVEoF6pCxm4R/cMEutUm3MY=",
      "2": "Euzymqx6iXjS3/NuGKDCiM2Ev6wdhnU+rBiKnJ7YpHE="
    }
  }
}
```

## Encrypt Data

This endpoint encrypts the provided plaintext using the named key. This path
supports the `create` and `update` policy capabilities as follows: if the user
has the `create` capability for this endpoint in their policies, and the key
does not exist, it will be upserted with default values (whether the key
requires derivation depends on whether the context parameter is empty or not).
If the user only has `update` capability and the key does not exist, an error
will be returned.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/encrypt/:name`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  encrypt against. This is specified as part of the URL.

- `plaintext` `(string: <required>)` – Specifies **base64 encoded** plaintext to
  be encoded.

- `context` `(string: "")` – Specifies the **base64 encoded** context for key
  derivation. This is required if key derivation is enabled for this key.

- `key_version` `(int: 0)` – Specifies the version of the key to use for
  encryption. If not set, uses the latest version. Must be greater than or
  equal to the key's `min_encryption_version`, if set.

- `nonce` `(string: "")` – Specifies the **base64 encoded** nonce value. This
  must be provided if convergent encryption is enabled for this key and the key
  was generated with Vault 0.6.1. Not required for keys created in 0.6.2+. The
  value must be exactly 96 bits (12 bytes) long and the user must ensure that
  for any given context (and thus, any given encryption key) this nonce value is
  **never reused**.

- `batch_input` `(array<object>: nil)` – Specifies a list of items to be
  encrypted in a single batch. When this parameter is set, if the parameters
  'plaintext', 'context' and 'nonce' are also set, they will be ignored. The
  format for the input is:

    ```json
    [
      {
        "context": "c2FtcGxlY29udGV4dA==",
        "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
      },
      {
        "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
        "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
      },
    ]
    ```

- `type` `(string: "aes256-gcm96")` –This parameter is required when encryption
  key is expected to be created. When performing an upsert operation, the type
  of key to create.

- `convergent_encryption` `(string: "")` – This parameter will only be used when
  a key is expected to be created.  Whether to support convergent encryption.
  This is only supported when using a key with key derivation enabled and will
  require all requests to carry both a context and 96-bit (12-byte) nonce. The
  given nonce will be used in place of a randomly generated nonce. As a result,
  when the same context and nonce are supplied, the same ciphertext is
  generated. It is _very important_ when using this mode that you ensure that
  all nonces are unique for a given context.  Failing to do so will severely
  impact the ciphertext's security.

### Sample Payload

```json
{
  "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
}
```

!> Vault HTTP API imposes a maximum request size of 32MB to prevent a denial
of service attack. This can be tuned per [`listener`
block](/docs/configuration/listener/tcp.html) in the Vault server
configuration.


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/encrypt/my-key
```

### Sample Response

```json
{
  "data": {
    "ciphertext": "vault:v1:abcdefgh"
  }
}
```

## Decrypt Data

This endpoint decrypts the provided ciphertext using the named key.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/decrypt/:name`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  decrypt against. This is specified as part of the URL.

- `ciphertext` `(string: <required>)` – Specifies the ciphertext to decrypt.

- `context` `(string: "")` – Specifies the **base64 encoded** context for key
  derivation. This is required if key derivation is enabled.

- `nonce` `(string: "")` – Specifies a base64 encoded nonce value used during
  encryption. Must be provided if convergent encryption is enabled for this key
  and the key was generated with Vault 0.6.1. Not required for keys created in
  0.6.2+.

- `batch_input` `(array<object>: nil)` – Specifies a list of items to be
  decrypted in a single batch. When this parameter is set, if the parameters
  'ciphertext', 'context' and 'nonce' are also set, they will be ignored. Format
  for the input goes like this:

    ```json
    [
      {
        "context": "c2FtcGxlY29udGV4dA==",
        "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
      },
      {
        "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
        "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
      },
    ]
    ```

### Sample Payload

```json
{
  "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/decrypt/my-key
```

### Sample Response

```json
{
  "data": {
    "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Rewrap Data

This endpoint rewraps the provided ciphertext using the latest version of the
named key. Because this never returns plaintext, it is possible to delegate this
functionality to untrusted users or scripts.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/rewrap/:name`      |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  re-encrypt against. This is specified as part of the URL.

- `ciphertext` `(string: <required>)` – Specifies the ciphertext to re-encrypt.

- `context` `(string: "")` – Specifies the **base64 encoded** context for key
  derivation. This is required if key derivation is enabled.

- `key_version` `(int: 0)` – Specifies the version of the key to use for the
  operation. If not set, uses the latest version. Must be greater than or equal
  to the key's `min_encryption_version`, if set.

- `nonce` `(string: "")` – Specifies a base64 encoded nonce value used during
  encryption. Must be provided if convergent encryption is enabled for this key
  and the key was generated with Vault 0.6.1. Not required for keys created in
  0.6.2+.

- `batch_input` `(array<object>: nil)` – Specifies a list of items to be
  decrypted in a single batch. When this parameter is set, if the parameters
  'ciphertext', 'context' and 'nonce' are also set, they will be ignored. Format
  for the input goes like this:

    ```json
    [
      {
        "context": "c2FtcGxlY29udGV4dA==",
        "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
      },
      {
        "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
        "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
      },
    ]
    ```

### Sample Payload

```json
{
  "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/rewrap/my-key
```

### Sample Response

```json
{
  "data": {
    "ciphertext": "vault:v2:abcdefgh"
  }
}
```

## Generate Data Key

This endpoint generates a new high-entropy key and the value encrypted with the
named key. Optionally return the plaintext of the key as well. Whether plaintext
is returned depends on the path; as a result, you can use Vault ACL policies to
control whether a user is allowed to retrieve the plaintext value of a key. This
is useful if you want an untrusted user or operation to generate keys that are
then made available to trusted users.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/datakey/:type/:name` |

### Parameters

- `type` `(string: <required>)` – Specifies the type of key to generate. If
  `plaintext`, the plaintext key will be returned along with the ciphertext. If
  `wrapped`, only the ciphertext value will be returned. This is specified as
  part of the URL.

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  use to encrypt the datakey. This is specified as part of the URL.

- `context` `(string: "")` – Specifies the key derivation context, provided as a
  base64-encoded string. This must be provided if derivation is enabled.

- `nonce` `(string: "")` – Specifies a nonce value, provided as base64 encoded.
  Must be provided if convergent encryption is enabled for this key and the key
  was generated with Vault 0.6.1. Not required for keys created in 0.6.2+. The
  value must be exactly 96 bits (12 bytes) long and the user must ensure that
  for any given context (and thus, any given encryption key) this nonce value is
  **never reused**.

- `bits` `(int: 256)` – Specifies the number of bits in the desired key. Can be
  128, 256, or 512.

### Sample Payload

```json
{
  "context": "Ab3=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/datakey/plaintext/my-key
```

### Sample Response

```json
{
  "data": {
    "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveAo=",
    "ciphertext": "vault:v1:abcdefgh"
  }
}
```

## Generate Random Bytes

This endpoint returns high-quality random bytes of the specified length.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/random(/:bytes)`   |

### Parameters

- `bytes` `(int: 32)` – Specifies the number of bytes to return. This value can
  be specified either in the request body, or as a part of the URL.

- `format` `(string: "base64")` – Specifies the output encoding. Valid options
  are `hex` or `base64`.

### Sample Payload

```json
{
  "format": "hex"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/random/164
```

### Sample Response

```json
{
  "data": {
    "random_bytes": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Hash Data

This endpoint returns the cryptographic hash of given data using the specified
algorithm.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/hash(/:algorithm)` |

### Parameters

- `algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use. This
  can also be specified as part of the URL. Currently-supported algorithms are:

    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: <required>)` – Specifies the **base64 encoded** input data.

- `format` `(string: "hex")` – Specifies the output encoding. This can be either
  `hex` or `base64`.

### Sample Payload

```json
{
  "input": "adba32=="
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/hash/sha2-512
```

### Sample Response

```json
{
  "data": {
    "sum": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

## Generate HMAC

This endpoint returns the digest of given data using the specified hash
algorithm and the named key. The key can be of any type supported by `transit`;
the raw key will be marshaled into bytes to be used for the HMAC function. If
the key is of a type that supports rotation, the latest (current) version will
be used.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/hmac/:name(/:algorithm)` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  generate hmac against. This is specified as part of the URL.

- `key_version` `(int: 0)` – Specifies the version of the key to use for the
  operation. If not set, uses the latest version. Must be greater than or equal
  to the key's `min_encryption_version`, if set.

- `algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use. This
  can also be specified as part of the URL. Currently-supported algorithms are:

    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: "")` – Specifies the **base64 encoded** input data. One of 
  `input` or `batch_input` must be supplied.

- `batch_input` `(array<object>: nil)` – Specifies a list of items for processing.
  When this parameter is set, if the parameter 'input' is also set, it will be 
  ignored.  Responses are returned in the 'batch_results' array component of the 
  'data' element of the response. If the input data value of an item is invalid, the 
  corresponding item in the 'batch_results' will have the key 'error' with a value 
  describing the error. The format for batch_input is:

    ```json
    {
      "batch_input": [
        {
          "input": "adba32=="
        },
        {
          "input": "aGVsbG8gd29ybGQuCg=="
        }
      ]
    }
    ```


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/hmac/my-key/sha2-512
```

### Sample Payload

```json
{
  "input": "adba32=="
}
```

### Sample Response

```json
{
  "data": {
    "hmac": "dGhlIHF1aWNrIGJyb3duIGZveAo="
  }
}
```

### Sample Payload with batch_input

```json
{
  "batch_input": [
    {
      "input": "adba32=="
    },
    {
      "input": "adba32=="
    },
    {},
    {
      "input": ""
    }
  ]
}
```

### Sample Response for batch_input

```json
{
  "data": {
    "batch_results": [
      {
        "hmac": "vault:v1:1jFhRYWHiddSKgEFyVRpX8ieX7UU+748NBwHKecXE3hnGBoAxrfgoD5U0yAvji7b5X6V1fP"
      },
      {
        "hmac": "vault:v1:1jFhRYWHiddSKgEFyVRpX8ieX7UU+748NBwHKecXE3hnGBoAxrfgoD5U0yAvji7b5X6V1fP"
      },
      {
        "error": "missing input for HMAC"
      },
      {
        "hmac": "vault:v1:/wsSP6iQ9ECO9RRkefKLXey9sDntzSjoiW0vBrWfUsYB0ISroyC6plUt/jN7gcOv9O+Ecow"
      }
    ]
  }
}
```


## Sign Data

This endpoint returns the cryptographic signature of the given data using the
named key and the specified hash algorithm. The key must be of a type that
supports signing.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/sign/:name(/:hash_algorithm)` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key to
  use for signing. This is specified as part of the URL.

- `key_version` `(int: 0)` – Specifies the version of the key to use for
  signing. If not set, uses the latest version. Must be greater than or equal
  to the key's `min_encryption_version`, if set.

- `hash_algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use for
  supporting key types (notably, not including `ed25519` which specifies its
  own hash algorithm). This can also be specified as part of the URL.
  Currently-supported algorithms are:

    - `sha1`
    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: "")` – Specifies the **base64 encoded** input data. One of 
  `input` or `batch_input` must be supplied.

- `batch_input` `(array<object>: nil)` – Specifies a list of items for processing.
  When this parameter is set, any supplied 'input' or 'context' parameters will be 
  ignored.  Responses are returned in the 'batch_results' array component of the 
  'data' element of the response. If the input data value of an item is invalid, the 
  corresponding item in the 'batch_results' will have the key 'error' with a value 
  describing the error. The format for batch_input is:

    ```json
    {
      "batch_input": [
        {
          "input": "adba32==",
          "context": "abcd"
        },
        {
          "input": "aGVsbG8gd29ybGQuCg==",
          "context": "efgh"
        }
      ]
    }
    ```

- `context` `(string: "")` - Base64 encoded context for key derivation.
   Required if key derivation is enabled; currently only available with ed25519
   keys.

- `prehashed` `(bool: false)` - Set to `true` when the input is already hashed.
  If the key type is `rsa-2048` or `rsa-4096`, then the algorithm used to hash
  the input should be indicated by the `hash_algorithm` parameter.  Just as the
  value to sign should be the base64-encoded representation of the exact binary
  data you want signed, when set, `input` is expected to be base64-encoded
  binary hashed data, not hex-formatted. (As an example, on the command line,
  you could generate a suitable input via `openssl dgst -sha256 -binary |
  base64`.)

- `signature_algorithm` `(string: "pss")` – When using a RSA key, specifies the RSA
  signature algorithm to use for signing. Supported signature types are:

    - `pss`
    - `pkcs1v15`

- `marshaling_algorithm` `(string: "asn1")` – Specifies the way in which the signature should be marshaled. This currently only applies to ECDSA keys. Supported types are:

    - `asn1`: The default, used by OpenSSL and X.509
    - `jws`: The version used by JWS (and thus for JWTs). Selecting this will
      also change the output encoding to URL-safe Base64 encoding instead of
      standard Base64-encoding.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/sign/my-key/sha2-512
```

### Sample Payload

```json
{
  "input": "adba32=="
}
```

### Sample Response

```json
{
  "data": {
    "signature": "vault:v1:MEUCIQCyb869d7KWuA0hBM9b5NJrmWzMW3/pT+0XYCM9VmGR+QIgWWF6ufi4OS2xo1eS2V5IeJQfsi59qeMWtgX0LipxEHI="
  }
}
```

### Sample Payload with batch_input

 Given an ed25519 key with derived keys set, the context parameter is expected for each batch_input item, and 
 the response will include the derived public key for each item.
```
{
  "batch_input": [
    {
      "input": "adba32==",
      "context": "efgh"
    },
    {
      "input": "adba32==",
      "context": "abcd"
    },
    {}
  ]
}
```

### Sample Response for batch_input
```
{
  "data": {
    "batch_results": [
      {
        "signature": "vault:v1:+R3cxAy6j4KriYzAyExU6p1glnyT/eLDSaUZO7gr8a8kgi/zSynNbOBSDJcGaAfLD1OF2hGupYBYTjmZMNoVAA==",
        "publickey": "2fQIaaem7+EhSGs3jUebAS/8qP2+sUrmxOmgqZIZc0c="
      },
      {
        "signature": "vault:v1:3hBwA88lnuAVJqb5rCCEstzKYaBTeSdejk356BTCE/nKwySOhzQH3mWCvJZwbRptNGa7ia5ykosYYdJz+aIKDA==",
        "publickey": "goDXuePo7L9z6HOw+a54O4HeV189BLECK9nAUudwp4Y="
      },
      {
        "error": "missing input"
      }
    ]
  },
}
```

## Verify Signed Data

This endpoint returns whether the provided signature is valid for the given
data.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/transit/verify/:name(/:hash_algorithm)` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the encryption key that
  was used to generate the signature or HMAC.

- `hash_algorithm` `(string: "sha2-256")` – Specifies the hash algorithm to use. This
  can also be specified as part of the URL. Currently-supported algorithms are:

    - `sha1`
    - `sha2-224`
    - `sha2-256`
    - `sha2-384`
    - `sha2-512`

- `input` `(string: "")` – Specifies the **base64 encoded** input data. One of 
  `input` or `batch_input` must be supplied.
  
- `signature` `(string: "")` – Specifies the signature output from the
  `/transit/sign` function. Either this must be supplied or `hmac` must be
  supplied.

- `hmac` `(string: "")` – Specifies the signature output from the
  `/transit/hmac` function. Either this must be supplied or `signature` must be
  supplied.

- `batch_input` `(array<object>: nil)` – Specifies a list of items for processing.
  When this parameter is set, any supplied 'input', 'hmac' or 'signature' parameters 
  will be ignored.  'batch_input' items should contain an 'input' parameter and
  either an 'hmac' or 'signature' parameter. All items in the batch must consistently
  supply either 'hmac' or 'signature' parameters.  It is an error for some items to
  supply 'hmac' while others supply 'signature'. Responses are returned in the 
  'batch_results' array component of the 'data' element of the response. If the 
  input data value of an item is invalid, the corresponding item in the 'batch_results' 
  will have the key 'error' with a value describing the error. The format for batch_input is:

    ```json
    {
      "batch_input": [
        {
          "input": "adba32==",
          "hmac": "vault:v1:1jFhRYWHiddSKgEFyVRpX8ieX7UU+748NBwHKecXE3hnGBoAxrfgoD5U0yAvji7b5X6V1fP"
        },
        {
          "input": "aGVsbG8gd29ybGQuCg==",
          "hmac": "vault:v1:/wsSP6iQ9ECO9RRkefKLXey9sDntzSjoiW0vBrWfUsYB0ISroyC6plUt/jN7gcOv9O+Ecow"
        }
      ]
    }
    ```

- `context` `(string: "")` - Base64 encoded context for key derivation.
   Required if key derivation is enabled; currently only available with ed25519
   keys.

- `prehashed` `(bool: false)` - Set to `true` when the input is already
   hashed. If the key type is `rsa-2048` or `rsa-4096`, then the algorithm used
   to hash the input should be indicated by the `hash_algorithm` parameter.

- `signature_algorithm` `(string: "pss")` – When using a RSA key, specifies the RSA
  signature algorithm to use for signature verification. Supported signature types
  are:

    - `pss`
    - `pkcs1v15`

- `marshaling_algorithm` `(string: "asn1")` – Specifies the way in which the signature was originally marshaled. This currently only applies to ECDSA keys. Supported types are:

    - `asn1`: The default, used by OpenSSL and X.509
    - `jws`: The version used by JWS (and thus for JWTs). Selecting this will
      also expect the input encoding to URL-safe Base64 encoding instead of
      standard Base64-encoding.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/verify/my-key/sha2-512
```

### Sample Payload

```json
{
  "input": "abcd13==",
  "signature": "vault:v1:MEUCIQCyb869d7KWuA..."
}
```

### Sample Response

```json
{
  "data": {
    "valid": true
  }
}
```

### Sample Payload with batch_input

```
{
  "batch_input": [
    {
      "input": "adba32==",
      "context": "abcd",
      "signature": "vault:v1:3hBwA88lnuAVJqb5rCCEstzKYaBTeSdejk356BTCE/nKwySOhzQH3mWCvJZwbRptNGa7ia5ykosYYdJz+aIKDA=="
    },
    {
      "input": "adba32==",
      "context": "efgh",
      "signature": "vault:v1:3hBwA88lnuAVJqb5rCCEstzKYaBTeSdejk356BTCE/nKwySOhzQH3mWCvJZwbRptNGa7ia5ykosYYdJz+aIKDA=="
    },
    {
      "input": "",
      "context": "abcd",
      "signature": "vault:v1:C/pxm5V1RI6kqudLdbLdj5Bpm2P38FKgvxoV69oNXphvJukRcQIqjZO793jCa2JPYPG21Y7vquDWy/Ff4Ma4AQ=="
    }
  ]
}
```

### Sample Response for batch_input
```
{
  "data": {
    "batch_results": [
      {
        "valid": true
      },
      {
        "valid": false
      },
      {
        "valid": true
      }
    ]
  },
}
```

## Backup Key

This endpoint returns a plaintext backup of a named key. The backup contains all
the configuration data and keys of all the versions along with the HMAC key.
The response from this endpoint can be used with the `/restore` endpoint to
restore the key.

| Method  | Path                    |
| :---------------------- | :--------------------- |
| `GET`   | `/transit/backup/:name` |

### Parameters

 - `name` `(string: <required>)` - Name of the key.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/transit/backup/aes
```

### Sample Response

```json
{
  "data": {
    "backup": "eyJwb2xpY3kiOnsibmFtZSI6ImFlcyIsImtleXMiOnsiMSI6eyJrZXkiOiJXK3k4Z0dOMHdiTDJLOU95NXFPN1laMGtjdzMvR0ZiNWM4STBzdlNMMnFNPSIsImhtYWNfa2V5IjoiUDBTcjh1YTJaZERNUTdPd2h4RGp1Z0U5d0JSR3Q2QXl6K0t4TzN5Z2M5ST0iLCJ0aW1lIjoiMjAxNy0xMi0wOFQxMTo1MDowOC42MTM4MzctMDU6MDAiLCJlY194IjpudWxsLCJlY195IjpudWxsLCJlY19kIjpudWxsLCJyc2Ffa2V5IjpudWxsLCJwdWJsaWNfa2V5IjoiIiwiY3JlYXRpb25fdGltZSI6MTUxMjc1MTgwOH19LCJkZXJpdmVkIjpmYWxzZSwia2RmIjowLCJjb252ZXJnZW50X2VuY3J5cHRpb24iOmZhbHNlLCJleHBvcnRhYmxlIjpmYWxzZSwibWluX2RlY3J5cHRpb25fdmVyc2lvbiI6MSwibWluX2VuY3J5cHRpb25fdmVyc2lvbiI6MCwibGF0ZXN0X3ZlcnNpb24iOjEsImFyY2hpdmVfdmVyc2lvbiI6MSwiZGVsZXRpb25fYWxsb3dlZCI6ZmFsc2UsImNvbnZlcmdlbnRfdmVyc2lvbiI6MCwidHlwZSI6MCwiYmFja3VwX2luZm8iOnsidGltZSI6IjIwMTctMTItMDhUMTE6NTA6MjkuMjI4MTU3LTA1OjAwIiwidmVyc2lvbiI6MX0sInJlc3RvcmVfaW5mbyI6bnVsbH0sImFyY2hpdmVkX2tleXMiOnsia2V5cyI6W3sia2V5IjpudWxsLCJobWFjX2tleSI6bnVsbCwidGltZSI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwiZWNfeCI6bnVsbCwiZWNfeSI6bnVsbCwiZWNfZCI6bnVsbCwicnNhX2tleSI6bnVsbCwicHVibGljX2tleSI6IiIsImNyZWF0aW9uX3RpbWUiOjB9LHsia2V5IjoiVyt5OGdHTjB3YkwySzlPeTVxTzdZWjBrY3czL0dGYjVjOEkwc3ZTTDJxTT0iLCJobWFjX2tleSI6IlAwU3I4dWEyWmRETVE3T3doeERqdWdFOXdCUkd0NkF5eitLeE8zeWdjOUk9IiwidGltZSI6IjIwMTctMTItMDhUMTE6NTA6MDguNjEzODM3LTA1OjAwIiwiZWNfeCI6bnVsbCwiZWNfeSI6bnVsbCwiZWNfZCI6bnVsbCwicnNhX2tleSI6bnVsbCwicHVibGljX2tleSI6IiIsImNyZWF0aW9uX3RpbWUiOjE1MTI3NTE4MDh9XX19Cg=="
  }
}
```

## Restore Key

This endpoint restores the backup as a named key. This will restore the key
configurations and all the versions of the named key along with HMAC keys. The
input to this endpoint should be the output of `/backup` endpoint.

 ~> For safety, by default the backend will refuse to restore to an existing
 key. If you want to reuse a key name, it is recommended you delete the key
 before restoring. It is a good idea to attempt restoring to a different key
 name first to verify that the operation successfully completes.

| Method   | Path                        |
| :-------------------------- | :--------------------- |
| `POST`   | `/transit/restore(/:name)`  |

### Parameters

 - `backup` `(string: <required>)` - Backed up key data to be restored. This
   should be the output from the `/backup` endpoint.

 - `name` `(string: <optional>)` - If set, this will be the name of the
   restored key.

 - `force` `(bool: false)` - If set, force the restore to proceed even if a key
   by this name already exists.

### Sample Payload

```json
  "backup": "eyJwb2xpY3kiOnsibmFtZSI6ImFlcyIsImtleXMiOnsiMSI6eyJrZXkiOiJXK3k4Z0dOMHdiTDJLOU95NXFPN1laMGtjdzMvR0ZiNWM4STBzdlNMMnFNPSIsImhtYWNfa2V5IjoiUDBTcjh1YTJaZERNUTdPd2h4RGp1Z0U5d0JSR3Q2QXl6K0t4TzN5Z2M5ST0iLCJ0aW1lIjoiMjAxNy0xMi0wOFQxMTo1MDowOC42MTM4MzctMDU6MDAiLCJlY194IjpudWxsLCJlY195IjpudWxsLCJlY19kIjpudWxsLCJyc2Ffa2V5IjpudWxsLCJwdWJsaWNfa2V5IjoiIiwiY3JlYXRpb25fdGltZSI6MTUxMjc1MTgwOH19LCJkZXJpdmVkIjpmYWxzZSwia2RmIjowLCJjb252ZXJnZW50X2VuY3J5cHRpb24iOmZhbHNlLCJleHBvcnRhYmxlIjpmYWxzZSwibWluX2RlY3J5cHRpb25fdmVyc2lvbiI6MSwibWluX2VuY3J5cHRpb25fdmVyc2lvbiI6MCwibGF0ZXN0X3ZlcnNpb24iOjEsImFyY2hpdmVfdmVyc2lvbiI6MSwiZGVsZXRpb25fYWxsb3dlZCI6ZmFsc2UsImNvbnZlcmdlbnRfdmVyc2lvbiI6MCwidHlwZSI6MCwiYmFja3VwX2luZm8iOnsidGltZSI6IjIwMTctMTItMDhUMTE6NTA6MjkuMjI4MTU3LTA1OjAwIiwidmVyc2lvbiI6MX0sInJlc3RvcmVfaW5mbyI6bnVsbH0sImFyY2hpdmVkX2tleXMiOnsia2V5cyI6W3sia2V5IjpudWxsLCJobWFjX2tleSI6bnVsbCwidGltZSI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwiZWNfeCI6bnVsbCwiZWNfeSI6bnVsbCwiZWNfZCI6bnVsbCwicnNhX2tleSI6bnVsbCwicHVibGljX2tleSI6IiIsImNyZWF0aW9uX3RpbWUiOjB9LHsia2V5IjoiVyt5OGdHTjB3YkwySzlPeTVxTzdZWjBrY3czL0dGYjVjOEkwc3ZTTDJxTT0iLCJobWFjX2tleSI6IlAwU3I4dWEyWmRETVE3T3doeERqdWdFOXdCUkd0NkF5eitLeE8zeWdjOUk9IiwidGltZSI6IjIwMTctMTItMDhUMTE6NTA6MDguNjEzODM3LTA1OjAwIiwiZWNfeCI6bnVsbCwiZWNfeSI6bnVsbCwiZWNfZCI6bnVsbCwicnNhX2tleSI6bnVsbCwicHVibGljX2tleSI6IiIsImNyZWF0aW9uX3RpbWUiOjE1MTI3NTE4MDh9XX19Cg=="
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/restore
```

## Trim Key

This endpoint trims older key versions setting a minimum version for the
keyring. Once trimmed, previous versions of the key cannot be recovered.

| Method   | Path                       |
| :------------------------- | :--------------------- |
| `POST`   | `/transit/keys/:name/trim` |

### Parameters

- `min_version` `(int: <required>)` - The minimum version for the key ring. All
  versions before this version will be permanently deleted. This value can at
  most be equal to the lesser of `min_decryption_version` and
  `min_encryption_version`. This is not allowed to be set when either
  `min_encryption_version` or `min_decryption_version` is set to zero.

### Sample Payload

```json
{
    "min_version": 2
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/keys/my-key/trim
```

## Configure Cache

This endpoint is used to configure the transit engine's cache. Note that configuration
changes will not be applied until the transit plugin is reloaded which can be achieved
 using the [`/sys/plugins/reload/backend`][sys-plugin-reload-backend] endpoint.

| Method   | Path                       |
| :------------------------- | :--------------------- |
| `POST`   | `/transit/cache-config` |

### Parameters

- `size` `(int: 0)` - Specifies the size in terms of number of entries. A size of
  `0` means unlimited. A _Least Recently Used_ (LRU) caching strategy is used for a
  non-zero cache size.

### Sample Payload

```json
{
  "size": 456
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..."
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/transit/cache-config
```

## Read Transit Cache Configuration

This endpoint retrieves configurations for the transit engine's cache.

| Method   | Path                       |
| :------------------------- | :--------------------- |
| `GET`   | `/transit/cache-config` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..."
    --request GET \
    http://127.0.0.1:8200/v1/transit/cache-config
```

### Sample Response

```json
  "data": {
    "size": 0
  },
```

[sys-plugin-reload-backend]: /api/system/plugins-reload-backend.html#reload-plugins