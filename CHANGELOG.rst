Changelog
=========

3.9.3
-----

Improvements:

* Dependency updates (#1699, #1703, #1710, #1714, #1715, #1723).
* Add ``persist-credentials: false`` to checkouts in GitHub workflows (#1704).
* Tests: use container images from https://github.com/getsops/ci-container-images (#1722).

Bugfixes:

* GnuPG: do not incorrectly trim fingerprint in presence of exclamation marks for specfic subkey selection (#1720).
* ``updatekeys`` subcommand: fix ``--input-type`` CLI flag being ignored (#1721).

Project changes:

* CI dependency updates (#1698, #1708, #1717).
* Rust dependency updates (#1707, #1716, #1725).


3.9.2
-----

Improvements:

* Dependency updates (#1645, #1649, #1653, #1662, #1686, #1693).
* Update compiled Protobuf definitions (#1688).
* Remove unused variables and simplify conditional (##1687).

Bugfixes:

* Handle whitespace in Azure Key Vault URLs (#1652).
* Correctly handle comments during JSON serialization (#1647).

Project changes:

* CI dependency updates (#1644, #1648, #1654, #1664, #1673, #1677, #1685).
* Rust dependency updates (#1655, #1663, #1670, #1676, #1689).
* Update and improve Protobuf code generation (#1688).


3.9.1
-----

Improvements:

* Dependency updates (#1550, #1554, #1558, #1562, #1565, #1568, #1575, #1581, #1589, #1593, #1602, #1603, #1618, #1629, #1635, #1639, #1640).
* Clarify naming of the configuration file in the documentation (#1569).
* Build with Go 1.22 (#1589).
* Specify filename of missing file in error messages (#1625).
* ``updatekeys`` subcommand: show changes in ``shamir_threshold`` (#1609).

Bugfixes:

* Fix the URL used for determining the latest SOPS version (#1553).
* ``updatekeys`` subcommand: actually use option ``--shamir-secret-sharing-threshold`` (#1608).
* Fix ``--config`` being ignored in subcommands by ``loadConfig`` (#1613).
* Allow ``edit`` subcommand to create files (#1596).
* Do not encrypt if a key group is empty, or there are no key groups (#1600).
* Do not ignore config errors when trying to parse a config file (#1614).

Project changes:

* CI dependency updates (#1551, #1555, #1559, #1564, #1566, #1574, #1584, #1586, #1590, #1592, #1619, #1628, #1634).
* Improve CI workflows (#1548, #1630).
* Ignore user-set environment variable ``SOPS_AGE_KEY_FILE`` in tests (#1595).
* Add example of using Age recipients in ``.sops.yaml`` (#1607).
* Add linting check for Rust code formatting (#1604).
* Set Rust version globally via ``rust-toolchain.toml`` for functional tests (#1612).
* Improve test coverage (#1617).
* Improve tests (#1622, #1624).
* Simplify branch rules to check DCO and ``check`` task instead of an explicit list of tasks in the CLI workflow (#1621).
* Build with Go 1.22 and 1.23 in CI and update Vault to 1.14 (#1531).
* Build release with Go 1.22 (#1615).
* Fix Dependabot config for Docker; add Dependabot config for Rust (#1632).
* Lock Rust package versions for functional tests for improved reproducibility (#1637).
* Rust dependency updates (#1638).

3.9.0
-----
Features:

* Add ``--mac-only-encrypted`` to compute MAC only over values which end up encrypted (#973)
* Allow configuration of indentation for YAML and JSON stores (#1273, #1372)
* Introduce a ``--pristine`` flag to ``sops exec-env`` (#912)
* Allow to pass multiple paths to ``sops updatekeys`` (#1274)
* Allow to override ``fileName`` with different value (#1332)
* Sort masterkeys according to ``--decryption-order`` (#1345)
* Add separate subcommands for encryption, decryption, rotating, editing, and setting values (#1391)
* Add ``filestatus`` command (#545)
* Add command ``unset`` (#1475)
* Merge key for key groups and make keys unique (#1493)
* Support using comments to select parts to encrypt (#974, #1392)

Deprecations:

* Deprecate the ``--background`` option to ``exec-env`` and ``exec-file`` (#1379)

Improvements:

* Warn/fail if the wrong number of arguments is provided (#1342)
* Warn if more than one command is used (#1388)
* Dependency updates (#1327, #1328, #1330, #1336, #1334, #1344, #1348, #1354, #1357, #1360, #1373, #1381, #1383, #1385, #1408, #1428, #1429, #1427, #1439, #1454, #1460, #1466, #1489, #1519, #1525, #1528, #1540, #1543, #1545)
* Build with Go 1.21 (#1427)
* Improve README.rst (#1339, #1399, #1350)
* Fix typos (#1337, #1477, #1484)
* Polish the ``sops help`` output a bit (#1341, #1544)
* Improve and fix tests (#1346, #1349, #1370, #1390, #1396, #1492)
* Create a constant for the ``sops`` metadata key (#1398)
* Refactoring: move extraction of encryption and rotation options to separate functions (#1389)

Bug fixes:

* Respect ``aws_profile`` from keygroup config (#1049)
* Fix a bug where not having a config results in a panic (#1371)
* Consolidate Flatten/Unflatten pre/post processing (#1356)
* INI and DotEnv stores: ``shamir_threshold`` is an integer (#1394)
* Make check whether file contains invalid keys for encryption dependent on output store (#1393)
* Do not panic if ``updatekeys`` is used with a config that has no creation rules defined (#1506)
* ``exec-file``: if ``--filename`` is used, use the provided filename without random suffix (#1474)
* Do not use DotEnv store for ``exec-env``, but specialized environment serializing code (#1436)
* Decryption: do not fail if no matching ``creation_rule`` is present in config file (#1434)

Project changes:

* CI dependency updates (#1347, #1359, #1376, #1382, #1386, #1425, #1432, #1498, #1503, #1508, #1510, #1516, #1521, #1492, #1534)
* Adjust Makefile to new goreleaser 6.0.0 release (#1526)

3.8.1
-----
Improvements:

* Improve handling of errors when binary store handles bad data (#1289)
* On macOS, prefer ``XDG_CONFIG_HOME`` over os.UserConfigDir() (#1291)
* Dependency updates (#1306, #1319, #1325)
* pgp: better error reporting for missing GPG binary during import of keys (#1286)
* Fix descriptions of unencrypted-regex and encrypted-regex flags, and ensure unencrypted_regex is considered in config validation (#1300)
* stores/json: improve error messages when parsing invalid JSON (#1307)

Bug fixes:

* pgp: improve handling of GnuPG home dir (#1298)
* Do not crash if an empty YAML file is encrypted (#1290)
* Handling of various ignored errors (#1304, #1311)
* pgp: do not require abs path for ``SOPS_GPG_EXEC`` (#1309)
* Report key rotation errors (#1317)
* Ensure wrapping of errors in main package (#1318)

Project changes:

* Enrich AWS authentication documentation (#1272)
* Add linting for RST and MD files (#1287)
* Delete SOPS encrypted file we don't have keys for (#1288)
* CI dependency updates (#1295, #1301)
* pgp: make error the last return value (#1310)
* Improve documentation files (#1320)

3.8.0
-----
Features:

* Support ``--version`` without network requests using ``--disable-version-check`` (#1115)
* Support ``--input-type`` for updatekeys command (#1116)

Improvements:

* pgp: modernize and improve, and add tests (#1054, #1282)
* azkv: update SDK to latest, add tests, tidy (#1067, #1092, #1256)
* age: improve identity loading, add tests, tidy (#1064)
* kms: AWS SDK V2, allow creds config, add tests (#1065, #1257)
* gcpkms: update SDK to latest, add tests, tidy (#1072, #1255)
* hcvault: update API, add tests, tidy (#1085)
* Do not report version when upstream ``--version`` check fails (#1124)
* Use GitHub endpoints in ``--version`` command (#1261)
* Close temporary file before invoking editor to widen support on Windows (#1265)
* Update dependencies (#1063, #1091, #1147, #1242, #1260, #1264, #1275, #1280, #1283)
* Deal with various deprecations of dependencies (#1113, #1262)

Bug fixes:

* Ensure YAML comments are not displaced (#1069)
* Ensure default Google credentials can be used again after introduction of ``GOOGLE_CREDENTIALS`` (#1249)
* Avoid duplicate logging of errors in some key sources (#1146, #1281)
* Using ``--set`` on a root level key does no longer truncate existing values (#899)
* Ensure stable order of SOPS parameters in dotenv file (#1101)

Project changes:

* Update Go to 1.20 (#1148)
* Update rustc functional tests to v1.70.0 (#1234)
* Remove remaining CircleCI workflow (#1237)
* Run CLI workflow on main (#1243)
* Delete obsolete ``validation/`` artifact (#1248)
* Rename Go module to ``github.com/getsops/sops/v3`` (#1247)
* Revamp release automation, including (Cosign) signed container images and checksums file, SLSA3 provenance and SBOMs (#1250)
* Update various bits of documentation (#1244)
* Add missing ``--encrypt`` flag from Vault example (#1060)
* Add documentation on how to use age in ``.sops.yaml`` (#1192)
* Improve Make targets and address various issues (#1258)
* Ensure clean working tree in CI (#1267)
* Fix CHANGELOG.rst formatting (#1269)
* Pin GitHub Actions to full length commit SHA and add CodeQL (#1276)
* Enable Dependabot for Docker, GitHub Actions and Go Mod (#1277)
* Generate versioned ``.intoto.jsonl`` (#1278)
* Update CI dependencies (#1279)

3.7.3
-----
Changes:

* Upgrade dependencies (#1024, #1045)
* Build alpine container in CI (#1018, #1032, #1025)
* keyservice: accept KeyServiceServer in LocalClient (#1035)
* Add support for GCP Service Account within ``GOOGLE_CREDENTIALS`` (#953)

Bug fixes:

* Upload the correct binary for the linux amd64 build (#1026)
* Fix bug when specifying multiple age recipients (#966)
* Allow for empty yaml maps (#908)
* Limit AWS role names to 64 characters (#1037)

3.7.2
-----
Changes:

* README updates (#861, #860)
* Various test fixes (#909, #906, #1008)
* Added Linux and Darwin arm64 releases (#911, #891)
* Upgrade to go v1.17 (#1012)
* Support SOPS_AGE_KEY environment variable (#1006)

Bug fixes:

* Make sure comments in yaml files are not duplicated (#866)
* Make sure configuration file paths work correctly relative to the config file in us (#853)

3.7.1
-----
Changes:

* Security fix
* Add release workflow (#843)
* Fix issue where CI wouldn't run against master (#848)
* Trim extra whitespace around age keys (#846)

3.7.0
-----
Features:

* Add support for age (#688)
* Add filename to exec-file (#761)

Changes:

* On failed decryption with GPG, return the error returned by GPG to the sops user (#762)
* Use yaml.v3 instead of modified yaml.v2 for handling YAML files (#791)
* Update aws-sdk-go to version v1.37.18 (#823)

Project Changes:

* Switch from TravisCI to Github Actions (#792)

3.6.1
-----
Features:

* Add support for --unencrypted-regex (#715)

Changes:

* Use keys.openpgp.org instead of gpg.mozilla.org (#732)
* Upgrade AWS SDK version (#714)
* Support --input-type for exec-file (#699)

Bug fixes:

* Fixes broken Vault tests (#731)
* Revert "Add standard newline/quoting behavior to dotenv store" (#706)


3.6.0
-----
Features:

* Support for encrypting data through the use of Hashicorp Vault (#655)
* ``sops publish`` now supports ``--recursive`` flag for publishing all files in a directory (#602)
* ``sops publish`` now supports ``--omit-extensions`` flag for omitting the extension in the destination path (#602)
* sops now supports JSON arrays of arrays (#642)

Improvements:

* Updates and standardization for the dotenv store (#612, #622)
* Close temp files after using them for edit command (#685)

Bug fixes:

* AWS SDK usage now correctly resolves the ``~/.aws/config`` file (#680)
* ``sops updatekeys`` now correctly matches config rules (#682)
* ``sops updatekeys`` now correctly uses the config path cli flag (#672)
* Partially empty sops config files don't break the use of sops anymore (#662)
* Fix possible infinite loop in PGP's passphrase prompt call (#690)

Project changes:

* Dockerfile now based off of golang version 1.14 (#649)
* Push alpine version of docker image to Dockerhub (#609)
* Push major, major.minor, and major.minor.patch tagged docker images to Dockerhub (#607)
* Removed out of date contact information (#668)
* Update authors in the cli help text (#645)


3.5.0
-----
Features:

* ``sops exec-env`` and ``sops exec-file``, two new commands for utilizing sops secrets within a temporary file or env vars

Bug fixes:

* Sanitize AWS STS session name, as sops creates it based off of the machines hostname
* Fix for ``decrypt.Data`` to support ``.ini`` files
* Various package fixes related to switching to Go Modules
* Fixes for Vault-related tests running locally and in CI.

Project changes:

* Change to proper use of go modules, changing to primary module name to ``go.mozilla.org/sops/v3``
* Change tags to requiring a ``v`` prefix.
* Add documentation for ``sops updatekeys`` command

3.4.0
-----
Features:

* ``sops publish``, a new command for publishing sops encrypted secrets to S3, GCS, or Hashicorp Vault
* Support for multiple Azure authentication mechanisms
* Azure Keyvault support to the sops config file
* ``encrypted_regex`` option to the sops config file

Bug fixes:

* Return non-zero exit code for invalid CLI flags
* Broken path handling for sops editing on Windows
* ``go lint/fmt`` violations
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
* Allow for ``--set`` value to be a string. (#461)

Project changes:

* Using ``develop`` as a staging branch to create releases off of. What
  is in ``master`` is now the current stable release.
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
-----

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
  keys to access a data key and decrypt a file. See ``sops groups -help`` and the
  documentation in README.

* Keyservice to forward access to a local master key on a socket, similar to
  gpg-agent. See ``sops keyservice --help`` and the documentation in README.

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
