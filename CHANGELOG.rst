Changelog
=========

3.4.0
-----
Features:

    * `sops publish`, a new command for publishing sops encrypted secrets to S3, GCS, or Hashicorp Vault
    * Support for multiple Azure authentication mechanisms
    * Azure Keyvault support to the sops config file
    * `encrypted_regex` option to the sops config file

Bug fixes:

    * Return non-zero exit code for invalid CLI flags
    * Broken path handling for sops editing on Windows
    * `go lint/fmt` violations
    * Check for pgp fingerprint before slicing it

Project changes:

    * Build container using golang 1.12
    * Switch to using go modules
    * Hashicorp Vault server in Travis CI build
    * Mozilla Publice License file to repo
    * Replaced expiring test gpg keys

3.3.1
-----

Bug fixes:

* Make sure the pgp key fingerprint is longer than 16 characters before
  slicing it. (#463)
* Allow for `--set` value to be a string. (#461)

Project changes:

* Using `develop` as a staging branch to create releases off of. What
  is in `master` is now the current stable release.
* Upgrade to using Go 1.12 to build sops
* Updated all vendored packages

3.3.0
-----

New features:

* Multi-document support for YAML files
* Support referencing AWS KMS keys by their alias
* Support for INI files
* Support for AWS CLI profiles
* Comment support in .env files
* Added vi to the list of known editors
* Added a way to specify the GPG key server to use through the
  SOPS_GPG_KEYSERVER environment variable

Bug fixes:

* Now uses $HOME instead of ~ (which didn't work) to find the GPG home
* Fix panic when vim was not available as an editor, but other
  alternative editors were
* Fix issue with AWS KMS Encryption Contexts (#445) with more than one
  context value failing to decrypt intermittently. Includes an
  automatic fix for old files affected by this issue.

Project infrastructure changes:

* Added integration tests for AWS KMS
* Added Code of Conduct


3.2.0
-----

* Added --output flag to write output a file directly instead of
  through stdout
* Added support for dotenv files

3.1.1
-----

* Fix incorrect version number from previous release

3.1.0
-----

* Add support for Azure Key Service

* Fix bug that prevented JSON escapes in input files from working

3.0.5
-----

* Prevent files from being encrypted twice

* Fix empty comments not being decrypted correctly

* If keyservicecmd returns an error, log it.

* Initial sops workspace auditing support (still wip)

* Refactor Store interface to reflect operations SOPS performs

3.0.3
----

* --set now works with nested data structures and not just simple
  values

* Changed default log level to warn instead of info

* Avoid creating empty files when using the editor mode to create new
  files and not making any changes to the example files

* Output unformatted strings when using --extract instead of encoding
  them to yaml

* Allow forcing binary input and output types from command line flags

* Deprecate filename_regex in favor of path_regex. filename_regex had
  a bug and matched on the whole file path, when it should have only
  matched on the file name. path_regex on the other hand is documented
  to match on the whole file path.

* Add an encrypted-suffix option, the exact opposite of
  unencrypted-suffix

* Allow specifying unencrypted_suffix and encrypted_suffix rules in
  the .sops.yaml configuration file

* Introduce key service flag optionally prompting users on
  encryption/decryption

3.0.1
-----

* Don't consider io.EOF returned by Decoder.Token as error

* add IsBinary: true to FileHints when encoding with crypto/openpgp 

* some improvements to error messages

3.0.0
-----

* Shamir secret sharing scheme support allows SOPS to require multiple master
  keys to access a data key and decrypt a file. See `sops groups -help` and the
  documentation in README.

* Keyservice to forward access to a local master key on a socket, similar to
  gpg-agent. See `sops keyservice --help` and the documentation in README.

* Encrypt comments by default

* Support for Google Compute Platform KMS

* Refactor of the store logic to separate the internal representation SOPS
  has of files from the external representation used in JSON and YAML files

* Reencoding of versions as string on sops 1.X files.
  **WARNING** this change breaks backward compatibility.
  SOPS shows an error message with instructions on how to solve
  this if it happens.
  
* Added command to reconfigure the keys used to encrypt/decrypt a file based on the .sops.yaml config file

* Retrieve missing PGP keys from gpg.mozilla.org

* Improved error messages for errors when decrypting files


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
