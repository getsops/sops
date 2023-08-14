# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

PROJECT             := github.com/getsops/sops/v3
PROJECT_DIR         := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
GO                  := GOPROXY=https://proxy.golang.org go

GITHUB_REPOSITORY   ?= github.com/getsops/sops

STATICCHECK         := $(PROJECT_DIR)/bin/staticcheck
STATICCHECK_VERSION := latest

GORELEASER          := $(PROJECT_DIR)/bin/goreleaser
GORELEASER_VERSION  ?= v1.20.0

.PHONY: all
all: test vet generate install functional-tests

.PHONY: origin-build
origin-build: test vet generate install functional-tests-all

install:
	$(GO) install github.com/getsops/sops/v3/cmd/sops

tag: all
	git tag -s $(TAGVER) -a -m "$(TAGMSG)"

.PHONY: staticcheck
staticcheck: install-staticcheck
	$(STATICCHECK) ./...

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

vet:
	$(GO) vet ./...

.PHONY: test
test: vendor
	gpg --import pgp/sops_functional_tests_key.asc 2>&1 1>/dev/null || exit 0
	./test.sh

showcoverage: test
	$(GO) tool cover -html=coverage.out

.PHONY: generate
generate: keyservice/keyservice.pb.go
	$(GO) generate

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. $<

.PHONY: functional-tests
functional-tests:
	$(GO) build -o functional-tests/sops github.com/getsops/sops/v3/cmd/sops
	cd functional-tests && cargo test

.PHONY: functional-tests-all
functional-tests-all:
	$(GO) build -o functional-tests/sops github.com/getsops/sops/v3/cmd/sops
	# Ignored tests are ones that require external services (e.g. AWS KMS)
	# 	TODO: Once `--include-ignored` lands in rust stable, switch to that.
	cd functional-tests && cargo test && cargo test -- --ignored

.PHONY: release-snapshot
release-snapshot: install-goreleaser
	GITHUB_REPOSITORY=$(GITHUB_REPOSITORY) $(GORELEASER) release --clean --snapshot --skip-sign

.PHONY: install-staticcheck
install-staticcheck:
	$(call go-install-tool,$(STATICCHECK),honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION),$(STATICCHECK_VERSION))

.PHONY: install-goreleaser
install-goreleaser:
	$(call go-install-tool,$(GORELEASER),github.com/goreleaser/goreleaser@$(GORELEASER_VERSION),$(GORELEASER_VERSION))

# go-install-tool will 'go install' any package $2 and install it to $1.
define go-install-tool
@[ -f $(1)-$(3) ] || { \
set -e ;\
GOBIN=$$(dirname $(1)) go install $(2) ;\
touch $(1)-$(3) ;\
}
endef
