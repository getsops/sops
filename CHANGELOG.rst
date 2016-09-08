Changelog
=========

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
