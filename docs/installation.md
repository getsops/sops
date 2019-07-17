Installation
============

Stable release
--------------

Binaries and packages of the latest stable release are available
[here](https://github.com/mozilla/sops/releases).

Development branch
------------------

For the adventurous, unstable features are available in the develop branch, which you can install from source:

```bash
$ go get -u go.mozilla.org/sops/cmd/sops
$ cd $GOPATH/src/go.mozilla.org/sops/
$ git checkout develop
$ make install
```

Requires Go 1.12 or newer. If you don't have Go installed, refer to the
[official Go documentation](https://golang.org/doc/install).
