# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

PROJECT		:= go.mozilla.org/sops
GO 		:= GO15VENDOREXPERIMENT=1 GO111MODULE=on GOPROXY=https://proxy.golang.org go
GOLINT 		:= golint

all: test vet generate install functional-tests
origin-build: test vet generate install functional-tests-all

install:
	$(GO) install go.mozilla.org/sops/cmd/sops

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
	$(GO) build -o functional-tests/sops go.mozilla.org/sops/cmd/sops
	cd functional-tests && cargo test

# Ignored tests are ones that require external services (e.g. AWS KMS)
# 	TODO: Once `--include-ignored` lands in rust stable, switch to that.
functional-tests-all:
	$(GO) build -o functional-tests/sops go.mozilla.org/sops/cmd/sops
	cd functional-tests && cargo test && cargo test -- --ignored

deb-pkg: install
	rm -rf tmppkg
	mkdir -p tmppkg/usr/local/bin
	cp $$GOPATH/bin/sops tmppkg/usr/local/bin/
	fpm -C tmppkg -n sops --license MPL2.0 --vendor mozilla \
		--description "Sops is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP." \
		-m "Julien Vehent <jvehent+sops@mozilla.com>" \
		--url https://go.mozilla.org/sops \
		--architecture x86_64 \
		-v "$$(git describe --abbrev=0 --tags)" \
		-s dir -t deb .

rpm-pkg: install
	rm -rf tmppkg
	mkdir -p tmppkg/usr/local/bin
	cp $$GOPATH/bin/sops tmppkg/usr/local/bin/
	fpm -C tmppkg -n sops --license MPL2.0 --vendor mozilla \
		--description "Sops is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP." \
		-m "Julien Vehent <jvehent+sops@mozilla.com>" \
		--url https://go.mozilla.org/sops \
		--architecture x86_64 \
		-v "$$(git describe --abbrev=0 --tags)" \
		-s dir -t rpm .

dmg-pkg: install
ifneq ($(OS),darwin)
		echo 'you must be on MacOS and set OS=darwin on the make command line to build an OSX package'
else
	rm -rf tmppkg
	mkdir -p tmppkg/usr/local/bin
	cp $$GOPATH/bin/sops tmppkg/usr/local/bin/
	fpm -C tmppkg -n sops --license MPL2.0 --vendor mozilla \
		--description "Sops is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP." \
		-m "Julien Vehent <jvehent+sops@mozilla.com>" \
		--url https://go.mozilla.org/sops \
		--architecture x86_64 \
		-v "$$(git describe --abbrev=0 --tags)" \
		-s dir -t osxpkg \
		--osxpkg-identifier-prefix org.mozilla.sops \
		-p tmppkg/sops-$$(git describe --abbrev=0 --tags).pkg .
	hdiutil makehybrid -hfs -hfs-volume-name "Mozilla Sops" \
		-o tmppkg/sops-$$(git describe --abbrev=0 --tags).dmg tmpdmg
endif

download-index:
	bash make_download_page.sh

mock:
	go get github.com/vektra/mockery/.../
	mockery -dir vendor/github.com/aws/aws-sdk-go/service/kms/kmsiface/ -name KMSAPI -output kms/mocks

.PHONY: all test generate clean vendor functional-tests mock
