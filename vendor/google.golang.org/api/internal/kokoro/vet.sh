# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

# Fail on error, and display commands being run.
set -ex

# Only run the linter on go1.11, since it needs type aliases (and we only care about its output once).
# TODO(deklerk) We should pass an environment variable from kokoro to decide this logic instead.
if [[ `go version` != *"go1.11"* ]]; then
    exit 0
fi

pwd

# Fail if a dependency was added without the necessary go.mod/go.sum change
# being part of the commit.
GO111MODULE=on go mod tidy
git diff go.mod | tee /dev/stderr | (! read)
git diff go.sum | tee /dev/stderr | (! read)

try3() { eval "$*" || eval "$*" || eval "$*"; }

try3 go get -u \
  golang.org/x/lint/golint \
  golang.org/x/tools/cmd/goimports \
  honnef.co/go/tools/cmd/staticcheck

# Look at all .go files (ignoring .pb.go files) and make sure they have a Copyright. Fail if any don't.
git ls-files "*[^.pb].go" | xargs grep -L "\(Copyright [0-9]\{4,\}\)" 2>&1 | tee /dev/stderr | (! read)
gofmt -s -d -l . 2>&1 | tee /dev/stderr | (! read)
goimports -l . 2>&1 | tee /dev/stderr | (! read)

# Runs the linter. Regrettably the linter is very simple and does not provide the ability to exclude rules or files,
# so we rely on inverse grepping to do this for us.
golint ./... 2>&1 | ( \
  grep -v "gen.go" | \
  grep -v "disco.go" | \
  grep -v "exported const DefaultDelayThreshold should have comment" | \
  grep -v "exported const DefaultBundleCountThreshold should have comment" | \
  grep -v "exported const DefaultBundleByteThreshold should have comment" | \
  grep -v "exported const DefaultBufferedByteLimit should have comment" | \
  grep -v "error var Done should have name of the form ErrFoo" | \
  grep -v "exported method APIKey.RoundTrip should have comment or be unexported" | \
  grep -v "exported method MarshalStyle.JSONReader should have comment or be unexported" | \
  grep -v "UnmarshalJSON should have comment or be unexported" | \
  grep -v "MarshalJSON should have comment or be unexported" | \
  grep -vE "\.pb\.go:" || true) | tee /dev/stderr | (! read)

staticcheck -go 1.9 ./... 2>&1 | ( \
  grep -v "SA1019" | \
  grep -v "S1007" | \
  grep -v "error var Done should have name of the form ErrFoo" | \
  grep -v "examples" | \
  grep -v "gen.go" || true) | tee /dev/stderr | (! read)
