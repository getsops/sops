#!/usr/bin/env bash
#
# Regenerates the static CRL test fixtures in this directory:
#   static_ca.pem, static_leaf.pem, static_leaf.key,
#   static_revoked.crl.pem, static_empty.crl.pem
#
# static_leaf.key is the leaf's private key; the end-to-end tests use it to
# stand up a TLS server that presents static_leaf.pem. It is a throwaway test
# key, not a secret.
#
# Overwrites the fixture files in place and cleans up every intermediate it
# creates. Run it from anywhere; it operates on its own directory.
#
#   ./gen-fixtures.sh
#
# To use a different leaf serial, pass it as the first argument (decimal):
#
#   ./gen-fixtures.sh 255
#
# The serial is opaque to the end-to-end tests, so no Go changes are needed.

set -euo pipefail

# Work in the directory this script lives in, regardless of the caller's cwd.
cd "$(dirname "$0")"

# Leaf serial in decimal; OpenSSL's CRL database wants it in hex.
leaf_serial="${1:-42}"
leaf_serial_hex=$(printf '%X' "$leaf_serial")

cleanup() {
  rm -rf ca.key leaf.csr leaf.ext openssl.cnf demoCA
}
trap cleanup EXIT

# ---- CA ---------------------------------------------------------------------
openssl ecparam -name prime256v1 -genkey -noout -out ca.key
openssl req -x509 -new -key ca.key -sha256 -days 36500 \
  -subj "/CN=Test Vault Static CA" \
  -addext "keyUsage=critical,keyCertSign,cRLSign" \
  -addext "basicConstraints=critical,CA:TRUE" \
  -set_serial 1 -out static_ca.pem

# ---- Leaf (signed by the CA) ------------------------------------------------
openssl ecparam -name prime256v1 -genkey -noout -out static_leaf.key
openssl req -new -key static_leaf.key -subj "/CN=vault.example.com" -out leaf.csr
cat > leaf.ext <<'EOF'
keyUsage=critical,digitalSignature
extendedKeyUsage=serverAuth
subjectAltName=DNS:vault.example.com
EOF
openssl x509 -req -in leaf.csr -CA static_ca.pem -CAkey ca.key \
  -set_serial "$leaf_serial" -days 36500 -sha256 -extfile leaf.ext \
  -out static_leaf.pem

# ---- CRL scaffolding --------------------------------------------------------
mkdir -p demoCA/newcerts
echo 1000 > demoCA/crlnumber
cat > openssl.cnf <<'EOF'
[ ca ]
default_ca = CA_default
[ CA_default ]
database         = demoCA/index.txt
crlnumber        = demoCA/crlnumber
default_md       = sha256
default_crl_days = 36500
EOF

revdate=$(date -u +%y%m%d%H%M%SZ)
# macOS `date`; fall back to the GNU form on Linux.
expdate=$(date -u -v+100y +%y%m%d%H%M%SZ 2>/dev/null \
  || date -u -d "+36500 days" +%y%m%d%H%M%SZ)

# ---- Revoked CRL: leaf serial marked revoked --------------------------------
printf "R\t%s\t%s\t%s\tunknown\t/CN=vault.example.com\n" \
  "$expdate" "$revdate" "$leaf_serial_hex" > demoCA/index.txt
openssl ca -config openssl.cnf -gencrl -keyfile ca.key -cert static_ca.pem \
  -out static_revoked.crl.pem

# ---- Empty CRL: nothing revoked ---------------------------------------------
: > demoCA/index.txt
openssl ca -config openssl.cnf -gencrl -keyfile ca.key -cert static_ca.pem \
  -out static_empty.crl.pem

echo "Regenerated fixtures (leaf serial ${leaf_serial} / 0x${leaf_serial_hex})."
