---
layout: "api"
page_title: "/sys/rekey - HTTP API"
sidebar_title: "<code>/sys/rekey</code>"
sidebar_current: "api-http-system-rekey"
description: |-
  The `/sys/rekey` endpoints are used to rekey the unseal keys for Vault.
---

# `/sys/rekey`

The `/sys/rekey` endpoints are used to rekey the unseal keys for Vault.

On seals that support stored keys (e.g. HSM PKCS11), the recovery key share(s)
can be provided to rekey the master key since no unseal keys are available. The
secret shares, secret threshold, and stored shares parameters must be set to 1.
Upon successful rekey, no split unseal key shares are returned.

## Read Rekey Progress

This endpoint reads the configuration and progress of the current rekey attempt.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/rekey/init`            |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/rekey/init
```

### Sample Response

```json
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "t": 3,
  "n": 5,
  "progress": 1,
  "required": 3,
  "pgp_fingerprints": ["abcd1234"],
  "backup": true,
  "verification_required": false
}
```

If a rekey is started, then `n` is the new shares to generate and `t` is the
threshold required for the new shares. `progress` is how many unseal keys have
been provided for this rekey, where `required` must be reached to complete. The
`nonce` for the current rekey operation is also displayed. If PGP keys are being
used to encrypt the final shares, the key fingerprints and whether the final
keys will be backed up to physical storage will also be displayed.
`verification_required` indicates whether verification was enabled for this
operation.

## Start Rekey

This endpoint initializes a new rekey attempt. Only a single rekey attempt can
take place at a time, and changing the parameters of a rekey requires canceling
and starting a new rekey, which will also provide a new nonce.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/rekey/init`            |

### Parameters

- `secret_shares` `(int: <required>)` – Specifies the number of shares to split
  the master key into.

- `secret_threshold` `(int: <required>)` – Specifies the number of shares
  required to reconstruct the master key. This must be less than or equal to
  `secret_shares`.

- `pgp_keys` `(array<string>: nil)` – Specifies an array of PGP public keys used
  to encrypt the output unseal keys. Ordering is preserved. The keys must be
  base64-encoded from their original binary representation. The size of this
  array must be the same as `secret_shares`.

- `backup` `(bool: false)` – Specifies if using PGP-encrypted keys, whether
  Vault should also store a plaintext backup of the PGP-encrypted keys at
  `core/unseal-keys-backup` in the physical storage backend. These can then
  be retrieved and removed via the `sys/rekey/backup` endpoint.

- `require_verification` `(bool: false)` – This turns on verification
  functionality. When verification is turned on, after successful authorization
  with the current unseal keys, the new unseal keys are returned but the master
  key is not actually rotated. The new keys must be provided to authorize the
  actual rotation of the master key. This ensures that the new keys have been
  successfully saved and protects against a risk of the keys being lost after
  rotation but before they can be persisted. This can be used with without
  `pgp_keys`, and when used with it, it allows ensuring that the returned keys
  can be successfully decrypted before committing to the new shares, which the
  backup functionality does not provide.

### Sample Payload

```json
{
  "secret_shares": 10,
  "secret_threshold": 5
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/rekey/init
```

## Cancel Rekey

This endpoint cancels any in-progress rekey. This clears the rekey settings as
well as any progress made. This must be called to change the parameters of the
rekey. Note: verification is still a part of a rekey. If rekeying is canceled
during the verification flow, the current unseal keys remain valid.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/rekey/init`            |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/rekey/init
```

## Read Backup Key

This endpoint returns the backup copy of PGP-encrypted unseal keys. The returned
value is the nonce of the rekey operation and a map of PGP key fingerprint to
hex-encoded PGP-encrypted key.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/sys/rekey/backup`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/rekey/backup
```

### Sample Response

```json
{
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "keys": {
    "abcd1234": "..."
  }
}
```

## Delete Backup Key

This endpoint deletes the backup copy of PGP-encrypted unseal keys.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE` | `/sys/rekey/backup`          |

### Sample Request

```
$ curl \
    --header "X-Vault-Token" \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/rekey/backup
```

## Submit Key

This endpoint is used to enter a single master key share to progress the rekey
of the Vault. If the threshold number of master key shares is reached, Vault
will complete the rekey. Otherwise, this API must be called multiple times until
that threshold is met. The rekey nonce operation must be provided with each
call.

When the operation is complete, this will return a response like the example
below; otherwise the response will be the same as the `GET` method against
`sys/rekey/init`, providing status on the operation itself.

If verification was requested, successfully completing this flow will
immediately put the operation into a verification state, and provide the nonce
for the verification operation.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/rekey/update`          |

### Parameters

- `key` `(string: <required>)` – Specifies a single master share key.

- `nonce` `(string: <required>)` – Specifies the nonce of the rekey operation.

### Sample Payload

```json
{
  "key": "AB32...",
  "nonce": "abcd1234..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token" \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/rekey/update
```

### Sample Response

```json
{
  "complete": true,
  "keys": ["one", "two", "three"],
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "pgp_fingerprints": ["abcd1234"],
  "keys_base64": ["base64keyvalue"],
  "backup": true,
  "verification_required": true,
  "verification_nonce": "8b112c9e-2738-929d-bcc2-19aff249ff10"
}
```

If the keys are PGP-encrypted, an array of key fingerprints will also be
provided (with the order in which the keys were used for encryption) along with
whether or not the keys were backed up to physical storage.

## Read Rekey Verification Progress

This endpoint reads the configuration and progress of the current rekey
verification attempt.

| Method   | Path                           |
| :----------------------------- | :--------------------- |
| `GET`    | `/sys/rekey/verify`            |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/sys/rekey/verify
```

### Sample Response

```json
{
  "nonce": "8b112c9e-2738-929d-bcc2-19aff249ff10",
  "t": 3,
  "n": 5,
  "progress": 1
}
```

`n` is the total number of new shares that were generated and `t` is the
threshold required for the new shares to pass verification. `progress` is how
many of the new unseal keys have been provided for this verification operation.
The `nonce` for the current rekey operation is also displayed.

## Cancel Rekey Verification

This endpoint cancels any in-progress rekey verification operation. This clears
any progress made and resets the nonce. Unlike a `DELETE` against
`sys/rekey/init`, this only resets the current verification operation, not the
entire rekey atttempt. The return value is the same as `GET` along with the new
nonce.

| Method   | Path                           |
| :----------------------------- | :--------------------- |
| `DELETE` | `/sys/rekey/verify`            | `200 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token" \
    --request DELETE \
    http://127.0.0.1:8200/v1/sys/rekey/verify
```

### Sample Response

```json
{
  "nonce": "5827bbc1-0110-5725-cc21-beddc129d942",
  "t": 3,
  "n": 5,
  "progress": 0
}
```

## Submit Verification Key

This endpoint is used to enter a single new key share to progress the rekey
verification operation.  If the threshold number of new key shares is reached,
Vault will complete the rekey by performing the actual rotation of the master
key. Otherwise, this API must be called multiple times until that threshold is
met. The nonce must be provided with each call.

When the operation is complete, this will return a response like the example
below; otherwise the response will be the same as the `GET` method against
`sys/rekey/verify`, providing status on the operation itself.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/sys/rekey/verify`          |

### Parameters

- `key` `(string: <required>)` – Specifies a single master share key from the
  new set of shares.

- `nonce` `(string: <required>)` – Specifies the nonce of the rekey
  verification operation.

### Sample Payload

```json
{
  "key": "A58d...",
  "nonce": "5a27bbc1..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token" \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/rekey/verify
```

### Sample Response

```json
{
  "nonce": "5827bbc1-0110-5725-cc21-beddc129d942",
  "complete": true
}
```
