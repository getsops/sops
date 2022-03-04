# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

PROJECT		:= go.mozilla.org/sops/v3
GO 		:= GOPROXY=https://proxy.golang.org go
GOLINT 		:= golint

all: test vet generate install functional-tests
origin-build: test vet generate install functional-tests-all

install:
	$(GO) install go.mozilla.org/sops/v3/cmd/sops

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
	$(GO) build -o functional-tests/sops go.mozilla.org/sops/v3/cmd/sops
	cd functional-tests && cargo test

# Ignored tests are ones that require external services (e.g. AWS KMS)
# 	TODO: Once `--include-ignored` lands in rust stable, switch to that.
functional-tests-all:
	$(GO) build -o functional-tests/sops go.mozilla.org/sops/v3/cmd/sops
	cd functional-tests && cargo test && cargo test -- --ignored

# Creates variables during target re-definition. Basically this block allows the particular variables to be used in the final target
build-deb-%: OS = $(word 1,$(subst -, ,$*))
build-deb-%: ARCH = $(word 2,$(subst -, ,$*))
build-deb-%: FPM_ARCH = $(word 3,$(subst -, ,$*))
# Poor-mans function with parameters being split out from the variable part of it's name
build-deb-%:
	rm -rf tmppkg
	mkdir -p tmppkg/usr/local/bin
	GOOS=$(OS) GOARCH="$(ARCH)" CGO_ENABLED=0 go build -mod vendor -o tmppkg/usr/local/bin/sops go.mozilla.org/sops/v3/cmd/sops
	fpm -C tmppkg -n sops --license MPL2.0 --vendor mozilla \
		--description "Sops is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP." \
		-m "AJ Bahnken <ajvb+sops@mozilla.com>" \
		--url https://go.mozilla.org/sops \
		--architecture $(FPM_ARCH) \
		-v "$$(grep '^const Version' version/version.go |cut -d \" -f 2)" \
		-s dir -t deb .

# Create .deb packages for multiple architectures
deb-pkg: vendor build-deb-linux-amd64-x86_64 build-deb-linux-arm64-arm64

# Creates variables during target re-definition. Basically this block allows the particular variables to be used in the final target
build-rpm-%: OS = $(word 1,$(subst -, ,$*))
build-rpm-%: ARCH = $(word 2,$(subst -, ,$*))
build-rpm-%: FPM_ARCH = $(word 3,$(subst -, ,$*))
# Poor-mans function with parameters being split out from the variable part of it's name
build-rpm-%:
	rm -rf tmppkg
	mkdir -p tmppkg/usr/local/bin
	GOOS=$(OS) GOARCH="$(ARCH)" CGO_ENABLED=0 go build -mod vendor -o tmppkg/usr/local/bin/sops go.mozilla.org/sops/v3/cmd/sops
	fpm -C tmppkg -n sops --license MPL2.0 --vendor mozilla \
		--description "Sops is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP." \
		-m "AJ Bahnken <ajvb+sops@mozilla.com>" \
		--url https://go.mozilla.org/sops \
		--architecture $(FPM_ARCH) \
		--rpm-os $(OS) \
		-v "$$(grep '^const Version' version/version.go |cut -d \" -f 2)" \
		-s dir -t rpm .

# Create .rpm packages for multiple architectures
rpm-pkg: vendor build-rpm-linux-amd64-x86_64 build-rpm-linux-arm64-arm64

dmg-pkg: install
ifneq ($(OS),darwin)
		echo 'you must be on MacOS and set OS=darwin on the make command line to build an OSX package'
else
	rm -rf tmppkg
	mkdir -p tmppkg/usr/local/bin
	cp $$GOPATH/bin/sops tmppkg/usr/local/bin/
	fpm -C tmppkg -n sops --license MPL2.0 --vendor mozilla \
		--description "Sops is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP." \
		-m "Mozilla Security <security@mozilla.org>" \
		--url https://go.mozilla.org/sops \
		--architecture x86_64 \
		-v "$$(grep '^const Version' version/version.go |cut -d \" -f 2)" \
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
