#!/bin/bash

# TODO(deklerk) Add integration tests when it's secure to do so. b/64723143

# Fail on any error
set -eo pipefail

# Display commands being run
set -x

# cd to project dir on Kokoro instance
cd github/gax-go

go version

# Set $GOPATH
export GOPATH="$HOME/go"
export GAX_HOME=$GOPATH/src/github.com/googleapis/gax-go
export PATH="$GOPATH/bin:$PATH"
mkdir -p $GAX_HOME

# Move code into $GOPATH and get dependencies
git clone . $GAX_HOME
cd $GAX_HOME

try3() { eval "$*" || eval "$*" || eval "$*"; }

download_deps() {
    if [[ `go version` == *"go1.11"* ]] || [[ `go version` == *"go1.12"* ]]; then
        export GO111MODULE=on
        # All packages, including +build tools, are fetched.
        try3 go mod download
    else
        # Because we don't provide -tags tools, the +build tools
        # dependencies aren't fetched.
        try3 go get -v -t ./...
    fi
}

download_deps
./internal/kokoro/check_incompat_changes.sh
./internal/kokoro/vet.sh
go test -race -v . 2>&1 | tee $KOKORO_ARTIFACTS_DIR/$KOKORO_GERRIT_CHANGE_NUMBER.txt

cd v2
download_deps
go test -race -v . 2>&1 | tee $KOKORO_ARTIFACTS_DIR/$KOKORO_GERRIT_CHANGE_NUMBER.txt
