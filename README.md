# Sops
A Go decrypter for [sops](https://github.com/mozilla/sops).

[![Codecov branch](https://img.shields.io/codecov/c/github/mozilla/sops/autrilla.svg?maxAge=2592000)]()

## Install (while on go-sops branch)
```
git clone -b go-sops git@github.com:mozilla/sops $GOPATH/src/github.com/mozilla/sops
go install github.com/mozilla/sops
```

## Decrypt
`sops -d <sops yaml file>`
