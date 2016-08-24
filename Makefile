# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

PROJECT		:= go.mozilla.org/sops
GO 		:= GO15VENDOREXPERIMENT=1 go
GOGETTER	:= GOPATH=$(shell pwd)/.tmpdeps go get -d
GOLINT 		:= golint

all: test vet generate install

install:
	$(GO) install github.com/mozilla-services/autograph

tag: all
	git tag -s $(TAGVER) -a -m "$(TAGMSG)"

lint:
	$(GOLINT) $(PROJECT)

vet:
	$(GO) vet $(PROJECT)

test:
	$(GO) test .
	$(GO) test ./aes
	$(GO) test ./json
	$(GO) test ./kms
	$(GO) test ./yaml
	$(GO) test ./pgp

showcoverage: test
	$(GO) tool cover -html=coverage.out

generate:
	$(GO) generate

.PHONY: all test generate clean autograph
