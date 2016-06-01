# Sops
A Go decrypter for [sops](https://github.com/mozilla/sops).

## Install (while on go-sops branch)
```
git clone -b go-sops git@github.com:mozilla/sops $GOPATH/src/github.com/mozilla/sops
go install github.com/mozilla/sops
```

## Decrypt
`sops -d <sops yaml file>`
