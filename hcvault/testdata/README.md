# Static CRL test fixtures

This directory holds fixed PEM/key material used by the CRL tests in
`../keysource_test.go` (`TestMasterKey_CRL_EndToEnd`). The Dockerized Vault
server booted in `TestMain` runs a TLS listener that presents `static_leaf.pem`,
and the tests drive `MasterKey.Encrypt` / `MasterKey.Decrypt` over HTTPS with
`SOPS_HC_VAULT_CRL_FILE` set, so CRL checking runs as part of the real TLS
handshake against a live Vault.

## Files

| File                    | Description                                                     |
| ----------------------- | --------------------------------------------------------------- |
| `static_ca.pem`         | Self-signed CA, `CN=Test Vault Static CA`, serial 1.            |
| `static_leaf.pem`       | Leaf `CN=vault.example.com`, serial 42, signed by the CA.       |
| `static_leaf.key`       | The leaf's private key; Vault's TLS listener serves this cert.  |
| `static_revoked.crl.pem`| Revokes the leaf (serial 42).                                  |
| `static_empty.crl.pem`  | Revokes nothing — the happy-path CRL.                           |

`static_leaf.pem` / `static_leaf.key` are bind-mounted read-only into the Vault
container (at `/fixtures`) by `TestMain`. The key is a throwaway test key, not
a secret.

All material uses P-256 / ECDSA and a ~100 year validity window, so the
certificates and CRLs do not expire for the practical lifetime of these tests.
The `NextUpdate` of each CRL is likewise far in the future, so
`isCRLCurrentlyValid` treats them as in effect. That window ends in 2126; if
you are reading this close to then, regenerate the fixtures with a fresh one.

## Regenerating / replacing the fixtures

Run [`gen-fixtures.sh`](gen-fixtures.sh); it drives OpenSSL to rebuild the
fixtures in place and cleans up its intermediates. You can run it from anywhere
— it operates on its own directory.

```sh
./gen-fixtures.sh
```

### Changing the leaf serial

The fixtures use leaf serial `42`. To use a different one, pass it (in decimal)
to the script. The serial is otherwise opaque to the tests, so no Go changes
are needed:

```sh
./gen-fixtures.sh 255
```

### Checking a regenerated fixture

Confirm the CA issued the leaf and the CRL revokes the serial you expect:

```sh
openssl verify -CAfile static_ca.pem static_leaf.pem
openssl crl -in static_revoked.crl.pem -noout -text
```

Then run the tests:

```sh
go test ../ -run TestMasterKey_CRL_EndToEnd -v
```
