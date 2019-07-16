# JWK Key Generator

The `jwk-keygen` command line utility generates keypairs used for asymmetric
encryption and signing algorithms in JSON Web Key (JWK) format.

## Usage

The utility requires specification of both desired algorithm (`alg`) and key
usage (`use`) to remind that same keypair should never be used both for
encryption and signing.

Algorithms are selected via the `--alg` flag, which influence the `alg` header.
For JWE (`--use=enc`), `--alg` specifies the key management algorithm (e.g.
`RSA-OAEP`). For JWS (`--use=sig`), `--alg` specifies the signature algorithm
(e.g. `PS256`).

Output file is determined by specified usage, algorithm and Key ID, e.g.
`jwk-keygen --use=sig --alg=RS512 --kid=test` produces files
`jwk_sig_RS512_test` and `jwk_sig_RS512_test.pub`. Keys are sent to stdout when
no Key ID is specified: neither pre-defined nor random one.

## Examples

### RSA 2048

Generate RSA/2048 key for encryption and output to stdout.

    jwk-keygen --use enc --alg RSA-OAEP

### Custom key length

Generate RSA/4096 key for signing and store to files.

    jwk-keygen --use sig --alg RS256 --bits 4096 --kid test

