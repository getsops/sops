# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

PROJECT		:= go.mozilla.org/sops
GO 		:= GO15VENDOREXPERIMENT=1 go
GOLINT 		:= golint

all: test vet generate install functional-tests

install:
	$(GO) install go.mozilla.org/sops/cmd/sops

tag: all
	git tag -s $(TAGVER) -a -m "$(TAGMSG)"

lint:
	$(GOLINT) $(PROJECT)

vendor:
	govend -u

vet:
	$(GO) vet $(PROJECT)

test:
	touch coverage.txt
	$(GO) test -coverprofile=coverage_tmp.txt -covermode=atomic $(PROJECT) && cat coverage_tmp.txt >> coverage.txt
	$(GO) test $(PROJECT)/aes -coverprofile=coverage_tmp.txt -covermode=atomic && cat coverage_tmp.txt >> coverage.txt
	$(GO) test $(PROJECT)/cmd/sops -coverprofile=coverage_tmp.txt -covermode=atomic && cat coverage_tmp.txt >> coverage.txt
	$(GO) test $(PROJECT)/json -coverprofile=coverage_tmp.txt -covermode=atomic && cat coverage_tmp.txt >> coverage.txt
	$(GO) test $(PROJECT)/yaml -coverprofile=coverage_tmp.txt -covermode=atomic && cat coverage_tmp.txt >> coverage.txt
	gpg --import pgp/sops_functional_tests_key.asc 2>&1 1>/dev/null || exit 0
	$(GO) test $(PROJECT)/pgp -coverprofile=coverage_tmp.txt -covermode=atomic && cat coverage_tmp.txt >> coverage.txt
	$(GO) test $(PROJECT)/kms -coverprofile=coverage_tmp.txt -covermode=atomic && cat coverage_tmp.txt >> coverage.txt

showcoverage: test
	$(GO) tool cover -html=coverage.out

generate:
	$(GO) generate

functional-tests:
	$(GO) build -o functional-tests/sops go.mozilla.org/sops/cmd/sops
	cd functional-tests && cargo test

.PHONY: all test generate clean vendor functional-tests
