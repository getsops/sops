#!/usr/bin/env bash
# Exit on first error
set -e

# Define our secret files
secret_files="secret.enc.json"

# For each of our files in our encrypted config
for file in $secret_files; do
  # Determine src and target for our file
  src_file="config/$file"
  target_file="$(echo "config/$file" | sed -E "s/.enc.json/.json/")"

  # If we only want to copy, then perform a copy
  # DEV: We allow `CONFIG_COPY_ONLY` to handle tests in Travis CI
  if test "$CONFIG_COPY_ONLY" = "TRUE"; then
    cp "$src_file" "$target_file"
  # Otherwise, decrypt it
  else
    sops --decrypt "$src_file" > "$target_file"
  fi
done
