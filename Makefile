# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

PROJECT := github.com/getsops/sops/v3
GO      := GOPROXY=https://proxy.golang.org go
GOLINT  := golint

all: test vet generate install functional-tests
origin-build: test vet generate install functional-tests-all

install:
	$(GO) install github.com/getsops/sops/v3/cmd/sops

tag: all
	git tag -s $(TAGVER) -a -m "$(TAGMSG)"

lint:
	$(GOLINT) $(PROJECT)

vendor:
	$(GO) mod tidy
	$(GO) mod vendor

vet:
	$(GO) vet $(PROJECT)

test: vendor
	gpg --import pgp/sops_functional_tests_key.asc 2>&1 1>/dev/null || exit 0
	./test.sh

showcoverage: test
	$(GO) tool cover -html=coverage.out

generate: keyservice/keyservice.pb.go
	$(GO) generate

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<

functional-tests:
	$(GO) build -o functional-tests/sops github.com/getsops/sops/v3/cmd/sops
	cd functional-tests && cargo test

# Ignored tests are ones that require external services (e.g. AWS KMS)
# 	TODO: Once `--include-ignored` lands in rust stable, switch to that.
functional-tests-all:
	$(GO) build -o functional-tests/sops github.com/getsops/sops/v3/cmd/sops
	cd functional-tests && cargo test && cargo test -- --ignored

.PHONY: all test generate clean vendor functional-tests
