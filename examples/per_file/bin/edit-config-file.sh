#!/usr/bin/env bash
# Exit on first error
set -e

# Localize our filepath
filepath="$1"
if test "$filepath" = ""; then
  echo "Expected \`filepath\` but received nothing" 1>&2
  echo "Usage: $0 <filepath>" 1>&2
  exit 1
fi

# Load our file into SOPS and run our sync script
sops "$filepath"
bin/decrypt-config.sh
