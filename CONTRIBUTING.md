# Contributing to SOPS

Mozilla welcomes contributions from everyone. Here are a few guidelines and instructions if you're thinking of helping with the development of SOPS.

# Getting started

* Make sure you have Go 1.6 or greater installed. You can find information on how to install Go [here](https://golang.org/dl/)
* After following the [Go installation guide](https://golang.org/doc/install), run `go get go.mozilla.org/sops`. This will automatically clone this repository.
* Switch into sops's directory, which will be in `$GOPATH/src/go.mozilla.org/sops`.
* Run the tests with `make test`. They should all pass.
* Fork the project on GitHub.
* Add your fork to git's remotes:
  * If you use SSH authentication: `git remote add <your username> git@github.com:<your username>/sops.git`.
  * Otherwise: `git remote add <your username> https://github.com/<your username>/sops.git`.
* Make any changes you want to sops, commit them, and push them to your fork.
* Create a pull request, and a contributor will come by and review your code. They may ask for some changes, and hopefully your contribution will be merged to the `master` branch!

# Guidelines

* Unless it's particularly hard, changes that fix a bug should have a regression test to make sure that the bug is not introduced again.
* New features and changes to existing features should be documented, and, if possible, tested.

# Communication

If you need any help contributing to sops, several contributors are on the `#go` channel on [Mozilla's IRC server](https://wiki.mozilla.org/IRC).
