#!/usr/bin/env bash

# Display commands being run
set -x

# Only run apidiff checks on go1.12 (we only need it once).
# TODO(deklerk) We should pass an environment variable from kokoro to decide
# this logic instead.
if [[ `go version` != *"go1.12"* ]]; then
    exit 0
fi

if git log -1 | grep BREAKING_CHANGE_ACCEPTABLE; then
  exit 0
fi

export GO111MODULE=on

go install golang.org/x/exp/cmd/apidiff

# We compare against master@HEAD. This is unfortunate in some cases: if you're
# working on an out-of-date branch, and master gets some new feature (that has
# nothing to do with your work on your branch), you'll get an error message.
# Thankfully the fix is quite simple: rebase your branch.
git clone https://github.com/googleapis/gax-go /tmp/gax

for dir in "" "/v2"; do
  pkg="github.com/googleapis/gax-go$dir"
  echo "Testing $pkg"

  # cd to the exact directory that specifies the go module so that it doesn't
  # use the module cache. https://go-review.googlesource.com/c/exp/+/155058
  cd "/tmp/gax$dir"
  apidiff -w /tmp/pkg.master $pkg
  cd - > /dev/null

  # TODO(deklerk) there's probably a nicer way to do this that doesn't require
  # two invocations
  if ! apidiff -incompatible /tmp/pkg.master $pkg | (! read); then
    apidiff -incompatible /tmp/pkg.master $pkg
    exit 1
  fi
done
