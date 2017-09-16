Changelog
=========

3.0.0
-----

* Shamir secret sharing scheme support allows SOPS to require multiple master
  keys to access a data key and decrypt a file. See `sops groups -help`.

* Keyservice to forward access to a local master key on a socket, similar to
  gpg-agent. See `sops keyservice --help`.

* Encrypt comments by default

* Support for Google Compute Platform KMS

* Refactor of the store logic to separate the internal representation SOPS
  has of files from the external representation used in JSON and YAML files

* Reencoding of versions as string on sops 1.X files, may break backward
  compatibility but will be handled automatically.

2.0.0
-----

* [major] rewrite in Go

1.14
----

* [medium] Support AWS KMS Encryption Contexts
* [minor] Support insertion in encrypted documents via --set
* [minor] Read location of gpg binary from SOPS_GPG_EXEC env variables

1.13
----

* [minor] handle $EDITOR variable with parameters

1.12
----

* [minor] make sure filename_regex gets applied to file names, not paths
* [minor] move check of latest version under the -V flag
* [medium] fix handling of binary data to preserve file integrity
* [minor] try to use configuration when encrypting existing files
