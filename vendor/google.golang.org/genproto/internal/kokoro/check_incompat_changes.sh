#!/usr/bin/env bash

# Display commands being run
set -x

# Only run apidiff checks on go1.11 (we only need it once).
# TODO(deklerk) We should pass an environment variable from kokoro to decide
# this logic instead.
if [[ `go version` != *"go1.11"* ]]; then
    exit 0
fi

try3() { eval "$*" || eval "$*" || eval "$*"; }

try3 go get -u golang.org/x/exp/cmd/apidiff

# We compare against master@HEAD. This is unfortunate in some cases: if you're
# working on an out-of-date branch, and master gets some new feature (that has
# nothing to do with your work on your branch), you'll get an error message.
# Thankfully the fix is quite simple: rebase your branch.
git clone https://github.com/google/go-genproto /tmp/genproto

V1_DIRS=`find . -type d -regex '.*v1$'`
V1_SUBDIRS=`find . -type d -regex '.*v1\/.*'`
for dir in $V1_DIRS $V1_SUBDIRS; do
  # turns things like ./foo/bar into foo/bar
  dir_without_junk=`echo $dir | sed -n "s#\(\.\/\)\(.*\)#\2#p"`
  pkg="google.golang.org/genproto/$dir_without_junk"
  echo "Testing $pkg"

  cd /tmp/genproto
  apidiff -w /tmp/pkg.master $pkg
  cd - > /dev/null

  # TODO(deklerk) there's probably a nicer way to do this that doesn't require
  # two invocations
  if ! apidiff -incompatible /tmp/pkg.master $pkg | (! read); then
    apidiff -incompatible /tmp/pkg.master $pkg
    exit 1
  fi
done
