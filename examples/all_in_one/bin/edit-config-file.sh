#!/usr/bin/env bash
# Exit on first error
set -e

# Define our secret files
secret_files="secret.enc.json"

# Look up our file
filepath="$1"
if test "$filepath" = ""; then
  echo "Expected \`filepath\` but received nothing" 1>&2
  echo "Usage: $0 <filepath>" 1>&2
  exit 1
fi

# If our file is a secret
filename="$(basename "$filepath")"
if echo "$secret_files" | grep "$filename"; then
  # Load it into SOPS and run our sync script
  sops "$filepath"
  bin/decrypt-config.sh
# Otherwise (it's a normal file)
else
  # Resolve our editor via `sops` logic
  editor="$EDITOR"
  if test "$editor" = ""; then
    editor="$(which vim nano | head -n 1)"
  fi
  if test "$editor" = ""; then
    echo "Expected \`EDITOR\` environment variable to be defined but it was not" 1>&2
    exit 1
  fi

  # Edit our file
  "$editor" "$filepath"
fi
