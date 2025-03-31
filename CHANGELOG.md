# Changelog

## 3.10.1

This is a re-release of 3.10.0 with no code changes.

Due to a failure during the 3.10.0 release, the
[commit cached by the Go infrastructure for 3.10.0](https://github.com/getsops/sops/commit/200bb6d8ab4063330bc99697255b3583501b3877)
is different from
[the commit tagged in the repository](https://github.com/getsops/sops/commit/4ed7060298fbcd00cafa359121ca62091b85bb6f).
To avoid confusion, we decided to push another release where the tag in the repository
will coincide with the commit cached by Go.

Project changes:

* CI dependency updates ([#1826](https://github.com/getsops/sops/pull/1826)).

## 3.10.0

Security fixes:

* Cherry-pick a fix for a timing vulnerability in the Shamir Secret Sharing code.
  The code was vendored from HashiCorp's Vault project, and the issue was fixed
  there two years ago; see [GHSA-vq4h-9ghm-qmrr](https://github.com/advisories/GHSA-vq4h-9ghm-qmrr)
  for details ([#1813](https://github.com/getsops/sops/pull/1813)).

Features:

* Add `--input-type` option for `sops filestatus` subcommand ([#1601](https://github.com/getsops/sops/pull/1601)).
* Allow to set the editor `sops` should use with the `SOPS_EDITOR` environment variable.
  If not set, `sops` falls back to `EDITOR` as before ([#1611](https://github.com/getsops/sops/pull/1611)).
* Allow users to disable the latest version check with the environment variable `SOPS_DISABLE_VERSION_CHECK`.
  Setting it to `1`, `t`, `T`, `TRUE`, `true`, or `True` explicitly
  disables the check ([#1684](https://github.com/getsops/sops/pull/1684)).
* Allow users to explicitly enable the latest version check with the `--check-for-updates`
  option ([#1816](https://github.com/getsops/sops/pull/1816)).
* Add duplicate section support for INI store ([#1452](https://github.com/getsops/sops/pull/1452)).
* Add check to prevent duplicate keys in YAML files ([#1203](https://github.com/getsops/sops/pull/1203)).
* Add `--same-process` option for the `sops exec-env` to use the `execve` syscall
  instead of starting the command in a child process ([#880](https://github.com/getsops/sops/pull/880)).
* Add `--idempotent` option for the `sops set` subcommand that will only
  write the file if a change happened ([#1754](https://github.com/getsops/sops/pull/1754)).
* Encrypt and decrypt `time.Time` objects that can appear in YAML files
  when using dates and timestamps ([#1759](https://github.com/getsops/sops/pull/1759)).
* Allow to encrypt and decrypt from `stdin` without having to provide
  platform-specific device names. This only works when using the
  `sops encrypt` and `sops decrypt` subcommands ([#1690](https://github.com/getsops/sops/pull/1690)).
* Allow to set the SOPS config location with the environment variable
  `SOPS_CONFIG` ([#1701](https://github.com/getsops/sops/pull/1701)).
* Support the `--config` option in the `sops publish` subcommand ([#1779](https://github.com/getsops/sops/pull/1779)).
* Omit empty master key metadata from encrypted files ([#1571](https://github.com/getsops/sops/pull/1571)).
* Add SSH support for Age ([#1692](https://github.com/getsops/sops/pull/1692)).
* Support Age identities with passphrases ([#1400](https://github.com/getsops/sops/pull/1400)).
* Add Age plugin support ([#1641](https://github.com/getsops/sops/pull/1641)).
* Allow to set the `SOPS_AGE_KEY_CMD` environment variable to an executable that
  returns Age keys ([#1811](https://github.com/getsops/sops/pull/1811)).
* Add support for `oauth2.TokenSource` injection from key service clients in
  GCP KMS ([#1794](https://github.com/getsops/sops/pull/1794)).
* Support `GOOGLE_OAUTH_ACCESS_TOKEN` for GCP KMS ([#1578](https://github.com/getsops/sops/pull/1578)).

Improvements:

* Dependency updates ([#1743](https://github.com/getsops/sops/pull/1743), [#1745](https://github.com/getsops/sops/pull/1745),
  [#1751](https://github.com/getsops/sops/pull/1751), [#1763](https://github.com/getsops/sops/pull/1763),
  [#1769](https://github.com/getsops/sops/pull/1769), [#1773](https://github.com/getsops/sops/pull/1773),
  [#1784](https://github.com/getsops/sops/pull/1784), [#1797](https://github.com/getsops/sops/pull/1797),
  [#1802](https://github.com/getsops/sops/pull/1802), [#1806](https://github.com/getsops/sops/pull/1806),
  [#1809](https://github.com/getsops/sops/pull/1809), [#1814](https://github.com/getsops/sops/pull/1814)).
* Fix typos ([#1765](https://github.com/getsops/sops/pull/1765)).
* Make sure that tests do not pick up `keys.txt` from user's `$HOME` dir ([#1766](https://github.com/getsops/sops/pull/1766)).
* Consolidate passphrase reading functionality in Age code ([#1775](https://github.com/getsops/sops/pull/1775)).
* Fix some problems reported by the `staticcheck` linter ([#1780](https://github.com/getsops/sops/pull/1780)).
* Improve documentation of Shamir Secret Sharing code to ease maintenance ([#1813](https://github.com/getsops/sops/pull/1813)).
* Make sure all files are properly formatted ([#1817](https://github.com/getsops/sops/pull/1817)).
* `sops` now warns if it finds a `.sops.yml` file while searching for a
  `.sops.yaml` config file ([#1820](https://github.com/getsops/sops/pull/1820)).

Bugfixes:

* Add trailing newline at the end of JSON files ([#1476](https://github.com/getsops/sops/pull/1476)).
* Check GnuPG decryption result for non-empty size. Certain older versions return
  an empty result with a successful return code when a AEAD cipher from a newer
  version was used ([#1776](https://github.com/getsops/sops/pull/1776)).
* Fix caching of `Metadata.DataKey` ([#1781](https://github.com/getsops/sops/pull/1781)).
* If `--filename-override` is specified, convert it to an absolute path same as regular
  filenames ([#1793](https://github.com/getsops/sops/pull/1793)).

Deprecations:

* The current behavior that `sops --version` always checks whether the current
  version is the latest is deprecated and will no longer be the default eventually.
  It is best to right now always specify `--disable-version-check` or `--check-for-updates`
  to `sops --version`, or alternatively set the environment variable `SOPS_DISABLE_VERSION_CHECK=true`
  to already get the planned default behavior today. ([#1816](https://github.com/getsops/sops/pull/1816)).

Project changes:

* Go 1.22 is no longer support; CI now also builds with Go 1.24 ([#1819](https://github.com/getsops/sops/pull/1819)).
* CI dependency updates ([#1746](https://github.com/getsops/sops/pull/1746),
  [#1750](https://github.com/getsops/sops/pull/1750), [#1770](https://github.com/getsops/sops/pull/1770),
  [#1782](https://github.com/getsops/sops/pull/1782), [#1795](https://github.com/getsops/sops/pull/1795),
  [#1801](https://github.com/getsops/sops/pull/1801), [#1808](https://github.com/getsops/sops/pull/1808)).
* Rust dependency updates for functional tests ([#1744](https://github.com/getsops/sops/pull/1744),
  [#1762](https://github.com/getsops/sops/pull/1762), [#1768](https://github.com/getsops/sops/pull/1768),
  [#1783](https://github.com/getsops/sops/pull/1783), [#1796](https://github.com/getsops/sops/pull/1796),
  [#1800](https://github.com/getsops/sops/pull/1800), [#1807](https://github.com/getsops/sops/pull/1807)).
* Bump Rust version for functional tests to 1.85 ([#1783](https://github.com/getsops/sops/pull/1783)).
* Release environment updates ([#1700](https://github.com/getsops/sops/pull/1700),
  [#1761](https://github.com/getsops/sops/pull/1761)).
* The changelog is now a MarkDown document ([#1741](https://github.com/getsops/sops/pull/1741)).
* We now also build a Windows ARM64 binary ([#1791](https://github.com/getsops/sops/pull/1791)).
* In the `updatekey.Opts` structure, `GroupQuorum` was renamed to `ShamirThreshold`
  ([#1631](https://github.com/getsops/sops/pull/1631)).
* Produce multiple Windows binaries ([#1823](https://github.com/getsops/sops/pull/1823)).

## 3.9.4

Improvements:

* Dependency updates ([#1727](https://github.com/getsops/sops/pull/1727), [#1732](https://github.com/getsops/sops/pull/1732),
  [#1734](https://github.com/getsops/sops/pull/1734), [#1739](https://github.com/getsops/sops/pull/1739)).

Bugfixes:

* Prevent key deduplication to identify different AWS KMS keys that only differ by
  role, context, or profile ([#1733](https://github.com/getsops/sops/pull/1733)).
* Update part of Azure SDK which prevented decryption in some cases ([#1695](https://github.com/getsops/sops/issue/1695),
  [#1734](https://github.com/getsops/sops/pull/1734)).

Project changes:

* CI dependency updates ([#1730](https://github.com/getsops/sops/pull/1730), [#1738](https://github.com/getsops/sops/pull/1738)).
* Rust dependency updates ([#1728](https://github.com/getsops/sops/pull/1728), [#1731](https://github.com/getsops/sops/pull/1731),
  [#1735](https://github.com/getsops/sops/pull/1735)).

## 3.9.3

Improvements:

* Dependency updates ([#1699](https://github.com/getsops/sops/pull/1699), [#1703](https://github.com/getsops/sops/pull/1703),
  [#1710](https://github.com/getsops/sops/pull/1710), [#1714](https://github.com/getsops/sops/pull/1714),
  [#1715](https://github.com/getsops/sops/pull/1715), [#1723](https://github.com/getsops/sops/pull/1723)).
* Add `persist-credentials: false` to checkouts in GitHub workflows ([#1704](https://github.com/getsops/sops/pull/1704)).
* Tests: use container images from
  [https://github.com/getsops/ci-container-images](https://github.com/getsops/ci-container-images)
  ([#1722](https://github.com/getsops/sops/pull/1722)).

Bugfixes:

* GnuPG: do not incorrectly trim fingerprint in presence of exclamation
  marks for specfic subkey selection ([#1720](https://github.com/getsops/sops/pull/1720)).
* `updatekeys` subcommand: fix `--input-type` CLI flag being ignored ([#1721](https://github.com/getsops/sops/pull/1721)).

Project changes:

* CI dependency updates ([#1698](https://github.com/getsops/sops/pull/1698), [#1708](https://github.com/getsops/sops/pull/1708),
  [#1717](https://github.com/getsops/sops/pull/1717)).
* Rust dependency updates ([#1707](https://github.com/getsops/sops/pull/1707), [#1716](https://github.com/getsops/sops/pull/1716),
  [#1725](https://github.com/getsops/sops/pull/1725)).

## 3.9.2

Improvements:

* Dependency updates ([#1645](https://github.com/getsops/sops/pull/1645), [#1649](https://github.com/getsops/sops/pull/1649),
  [#1653](https://github.com/getsops/sops/pull/1653), [#1662](https://github.com/getsops/sops/pull/1662),
  [#1686](https://github.com/getsops/sops/pull/1686), [#1693](https://github.com/getsops/sops/pull/1693)).
* Update compiled Protobuf definitions ([#1688](https://github.com/getsops/sops/pull/1688)).
* Remove unused variables and simplify conditional (##1687).

Bugfixes:

* Handle whitespace in Azure Key Vault URLs ([#1652](https://github.com/getsops/sops/pull/1652)).
* Correctly handle comments during JSON serialization ([#1647](https://github.com/getsops/sops/pull/1647)).

Project changes:

* CI dependency updates ([#1644](https://github.com/getsops/sops/pull/1644), [#1648](https://github.com/getsops/sops/pull/1648),
  [#1654](https://github.com/getsops/sops/pull/1654), [#1664](https://github.com/getsops/sops/pull/1664),
  [#1673](https://github.com/getsops/sops/pull/1673), [#1677](https://github.com/getsops/sops/pull/1677),
  [#1685](https://github.com/getsops/sops/pull/1685)).
* Rust dependency updates ([#1655](https://github.com/getsops/sops/pull/1655), [#1663](https://github.com/getsops/sops/pull/1663),
  [#1670](https://github.com/getsops/sops/pull/1670), [#1676](https://github.com/getsops/sops/pull/1676),
  [#1689](https://github.com/getsops/sops/pull/1689)).
* Update and improve Protobuf code generation ([#1688](https://github.com/getsops/sops/pull/1688)).

## 3.9.1

Improvements:

* Dependency updates ([#1550](https://github.com/getsops/sops/pull/1550), [#1554](https://github.com/getsops/sops/pull/1554),
  [#1558](https://github.com/getsops/sops/pull/1558), [#1562](https://github.com/getsops/sops/pull/1562),
  [#1565](https://github.com/getsops/sops/pull/1565), [#1568](https://github.com/getsops/sops/pull/1568),
  [#1575](https://github.com/getsops/sops/pull/1575), [#1581](https://github.com/getsops/sops/pull/1581),
  [#1589](https://github.com/getsops/sops/pull/1589), [#1593](https://github.com/getsops/sops/pull/1593),
  [#1602](https://github.com/getsops/sops/pull/1602), [#1603](https://github.com/getsops/sops/pull/1603),
  [#1618](https://github.com/getsops/sops/pull/1618), [#1629](https://github.com/getsops/sops/pull/1629),
  [#1635](https://github.com/getsops/sops/pull/1635), [#1639](https://github.com/getsops/sops/pull/1639),
  [#1640](https://github.com/getsops/sops/pull/1640)).
* Clarify naming of the configuration file in the documentation ([#1569](https://github.com/getsops/sops/pull/1569)).
* Build with Go 1.22 ([#1589](https://github.com/getsops/sops/pull/1589)).
* Specify filename of missing file in error messages ([#1625](https://github.com/getsops/sops/pull/1625)).
* `updatekeys` subcommand: show changes in `shamir_threshold` ([#1609](https://github.com/getsops/sops/pull/1609)).

Bugfixes:

* Fix the URL used for determining the latest SOPS version ([#1553](https://github.com/getsops/sops/pull/1553)).
* `updatekeys` subcommand: actually use option
  `--shamir-secret-sharing-threshold` ([#1608](https://github.com/getsops/sops/pull/1608)).
* Fix `--config` being ignored in subcommands by `loadConfig` ([#1613](https://github.com/getsops/sops/pull/1613)).
* Allow `edit` subcommand to create files ([#1596](https://github.com/getsops/sops/pull/1596)).
* Do not encrypt if a key group is empty, or there are no key groups ([#1600](https://github.com/getsops/sops/pull/1600)).
* Do not ignore config errors when trying to parse a config file ([#1614](https://github.com/getsops/sops/pull/1614)).

Project changes:

* CI dependency updates ([#1551](https://github.com/getsops/sops/pull/1551), [#1555](https://github.com/getsops/sops/pull/1555),
  [#1559](https://github.com/getsops/sops/pull/1559), [#1564](https://github.com/getsops/sops/pull/1564),
  [#1566](https://github.com/getsops/sops/pull/1566), [#1574](https://github.com/getsops/sops/pull/1574),
  [#1584](https://github.com/getsops/sops/pull/1584), [#1586](https://github.com/getsops/sops/pull/1586),
  [#1590](https://github.com/getsops/sops/pull/1590), [#1592](https://github.com/getsops/sops/pull/1592),
  [#1619](https://github.com/getsops/sops/pull/1619), [#1628](https://github.com/getsops/sops/pull/1628),
  [#1634](https://github.com/getsops/sops/pull/1634)).
* Improve CI workflows ([#1548](https://github.com/getsops/sops/pull/1548), [#1630](https://github.com/getsops/sops/pull/1630)).
* Ignore user-set environment variable `SOPS_AGE_KEY_FILE` in tests ([#1595](https://github.com/getsops/sops/pull/1595)).
* Add example of using Age recipients in `.sops.yaml` ([#1607](https://github.com/getsops/sops/pull/1607)).
* Add linting check for Rust code formatting ([#1604](https://github.com/getsops/sops/pull/1604)).
* Set Rust version globally via `rust-toolchain.toml` for functional tests ([#1612](https://github.com/getsops/sops/pull/1612)).
* Improve test coverage ([#1617](https://github.com/getsops/sops/pull/1617)).
* Improve tests ([#1622](https://github.com/getsops/sops/pull/1622), [#1624](https://github.com/getsops/sops/pull/1624)).
* Simplify branch rules to check DCO and `check` task instead of an explicit
  list of tasks in the CLI workflow ([#1621](https://github.com/getsops/sops/pull/1621)).
* Build with Go 1.22 and 1.23 in CI and update Vault to 1.14 ([#1531](https://github.com/getsops/sops/pull/1531)).
* Build release with Go 1.22 ([#1615](https://github.com/getsops/sops/pull/1615)).
* Fix Dependabot config for Docker; add Dependabot config for Rust ([#1632](https://github.com/getsops/sops/pull/1632)).
* Lock Rust package versions for functional tests for improved
  reproducibility ([#1637](https://github.com/getsops/sops/pull/1637)).
* Rust dependency updates ([#1638](https://github.com/getsops/sops/pull/1638)).

## 3.9.0

Features:

* Add `--mac-only-encrypted` to compute MAC only over values which
  end up encrypted ([#973](https://github.com/getsops/sops/pull/973))
* Allow configuration of indentation for YAML and JSON stores ([#1273](https://github.com/getsops/sops/pull/1273),
  [#1372](https://github.com/getsops/sops/pull/1372))
* Introduce a `--pristine` flag to `sops exec-env` ([#912](https://github.com/getsops/sops/pull/912))
* Allow to pass multiple paths to `sops updatekeys` ([#1274](https://github.com/getsops/sops/pull/1274))
* Allow to override `fileName` with different value ([#1332](https://github.com/getsops/sops/pull/1332))
* Sort masterkeys according to `--decryption-order` ([#1345](https://github.com/getsops/sops/pull/1345))
* Add separate subcommands for encryption, decryption, rotating, editing,
  and setting values ([#1391](https://github.com/getsops/sops/pull/1391))
* Add `filestatus` command ([#545](https://github.com/getsops/sops/pull/545))
* Add command `unset` ([#1475](https://github.com/getsops/sops/pull/1475))
* Merge key for key groups and make keys unique ([#1493](https://github.com/getsops/sops/pull/1493))
* Support using comments to select parts to encrypt ([#974](https://github.com/getsops/sops/pull/974),
  [#1392](https://github.com/getsops/sops/pull/1392))

Deprecations:

* Deprecate the `--background` option to `exec-env` and `exec-file` ([#1379](https://github.com/getsops/sops/pull/1379))

Improvements:

* Warn/fail if the wrong number of arguments is provided ([#1342](https://github.com/getsops/sops/pull/1342))
* Warn if more than one command is used ([#1388](https://github.com/getsops/sops/pull/1388))
* Dependency updates ([#1327](https://github.com/getsops/sops/pull/1327),
  [#1328](https://github.com/getsops/sops/pull/1328), [#1330](https://github.com/getsops/sops/pull/1330),
  [#1336](https://github.com/getsops/sops/pull/1336), [#1334](https://github.com/getsops/sops/pull/1334),
  [#1344](https://github.com/getsops/sops/pull/1344), [#1348](https://github.com/getsops/sops/pull/1348),
  [#1354](https://github.com/getsops/sops/pull/1354), [#1357](https://github.com/getsops/sops/pull/1357),
  [#1360](https://github.com/getsops/sops/pull/1360), [#1373](https://github.com/getsops/sops/pull/1373),
  [#1381](https://github.com/getsops/sops/pull/1381), [#1383](https://github.com/getsops/sops/pull/1383),
  [#1385](https://github.com/getsops/sops/pull/1385), [#1408](https://github.com/getsops/sops/pull/1408),
  [#1428](https://github.com/getsops/sops/pull/1428), [#1429](https://github.com/getsops/sops/pull/1429),
  [#1427](https://github.com/getsops/sops/pull/1427), [#1439](https://github.com/getsops/sops/pull/1439),
  [#1454](https://github.com/getsops/sops/pull/1454), [#1460](https://github.com/getsops/sops/pull/1460),
  [#1466](https://github.com/getsops/sops/pull/1466), [#1489](https://github.com/getsops/sops/pull/1489),
  [#1519](https://github.com/getsops/sops/pull/1519), [#1525](https://github.com/getsops/sops/pull/1525),
  [#1528](https://github.com/getsops/sops/pull/1528), [#1540](https://github.com/getsops/sops/pull/1540),
  [#1543](https://github.com/getsops/sops/pull/1543), [#1545](https://github.com/getsops/sops/pull/1545))
* Build with Go 1.21 ([#1427](https://github.com/getsops/sops/pull/1427))
* Improve README.rst ([#1339](https://github.com/getsops/sops/pull/1339),
  [#1399](https://github.com/getsops/sops/pull/1399), [#1350](https://github.com/getsops/sops/pull/1350))
* Fix typos ([#1337](https://github.com/getsops/sops/pull/1337), [#1477](https://github.com/getsops/sops/pull/1477),
  [#1484](https://github.com/getsops/sops/pull/1484))
* Polish the `sops help` output a bit ([#1341](https://github.com/getsops/sops/pull/1341),
  [#1544](https://github.com/getsops/sops/pull/1544))
* Improve and fix tests ([#1346](https://github.com/getsops/sops/pull/1346),
  [#1349](https://github.com/getsops/sops/pull/1349), [#1370](https://github.com/getsops/sops/pull/1370),
  [#1390](https://github.com/getsops/sops/pull/1390), [#1396](https://github.com/getsops/sops/pull/1396),
  [#1492](https://github.com/getsops/sops/pull/1492))
* Create a constant for the `sops` metadata key ([#1398](https://github.com/getsops/sops/pull/1398))
* Refactoring: move extraction of encryption and rotation options to
  separate functions ([#1389](https://github.com/getsops/sops/pull/1389))

Bug fixes:

* Respect `aws_profile` from keygroup config ([#1049](https://github.com/getsops/sops/pull/1049))
* Fix a bug where not having a config results in a panic ([#1371](https://github.com/getsops/sops/pull/1371))
* Consolidate Flatten/Unflatten pre/post processing ([#1356](https://github.com/getsops/sops/pull/1356))
* INI and DotEnv stores: `shamir_threshold` is an integer ([#1394](https://github.com/getsops/sops/pull/1394))
* Make check whether file contains invalid keys for encryption dependent
  on output store ([#1393](https://github.com/getsops/sops/pull/1393))
* Do not panic if `updatekeys` is used with a config that has no creation
  rules defined ([#1506](https://github.com/getsops/sops/pull/1506))
* `exec-file`: if `--filename` is used, use the provided filename without
  random suffix ([#1474](https://github.com/getsops/sops/pull/1474))
* Do not use DotEnv store for `exec-env`, but specialized environment
  serializing code ([#1436](https://github.com/getsops/sops/pull/1436))
* Decryption: do not fail if no matching `creation_rule` is present in
  config file ([#1434](https://github.com/getsops/sops/pull/1434))

Project changes:

* CI dependency updates ([#1347](https://github.com/getsops/sops/pull/1347),
  [#1359](https://github.com/getsops/sops/pull/1359), [#1376](https://github.com/getsops/sops/pull/1376),
  [#1382](https://github.com/getsops/sops/pull/1382), [#1386](https://github.com/getsops/sops/pull/1386),
  [#1425](https://github.com/getsops/sops/pull/1425), [#1432](https://github.com/getsops/sops/pull/1432),
  [#1498](https://github.com/getsops/sops/pull/1498), [#1503](https://github.com/getsops/sops/pull/1503),
  [#1508](https://github.com/getsops/sops/pull/1508), [#1510](https://github.com/getsops/sops/pull/1510),
  [#1516](https://github.com/getsops/sops/pull/1516), [#1521](https://github.com/getsops/sops/pull/1521),
  [#1492](https://github.com/getsops/sops/pull/1492), [#1534](https://github.com/getsops/sops/pull/1534))
* Adjust Makefile to new goreleaser 6.0.0 release ([#1526](https://github.com/getsops/sops/pull/1526))

## 3.8.1

Improvements:

* Improve handling of errors when binary store handles bad data ([#1289](https://github.com/getsops/sops/pull/1289))
* On macOS, prefer `XDG_CONFIG_HOME` over os.UserConfigDir() ([#1291](https://github.com/getsops/sops/pull/1291))
* Dependency updates ([#1306](https://github.com/getsops/sops/pull/1306),
  [#1319](https://github.com/getsops/sops/pull/1319), [#1325](https://github.com/getsops/sops/pull/1325))
* pgp: better error reporting for missing GPG binary during import of keys ([#1286](https://github.com/getsops/sops/pull/1286))
* Fix descriptions of `unencrypted-regex` and `encrypted-regex` flags, and
  ensure `unencrypted_regex` is considered in config validation ([#1300](https://github.com/getsops/sops/pull/1300))
* stores/json: improve error messages when parsing invalid JSON ([#1307](https://github.com/getsops/sops/pull/1307))

Bug fixes:

* pgp: improve handling of GnuPG home dir ([#1298](https://github.com/getsops/sops/pull/1298))
* Do not crash if an empty YAML file is encrypted ([#1290](https://github.com/getsops/sops/pull/1290))
* Handling of various ignored errors ([#1304](https://github.com/getsops/sops/pull/1304),
  [#1311](https://github.com/getsops/sops/pull/1311))
* pgp: do not require abs path for `SOPS_GPG_EXEC` ([#1309](https://github.com/getsops/sops/pull/1309))
* Report key rotation errors ([#1317](https://github.com/getsops/sops/pull/1317))
* Ensure wrapping of errors in main package ([#1318](https://github.com/getsops/sops/pull/1318))

Project changes:

* Enrich AWS authentication documentation ([#1272](https://github.com/getsops/sops/pull/1272))
* Add linting for RST and MD files ([#1287](https://github.com/getsops/sops/pull/1287))
* Delete SOPS encrypted file we don't have keys for ([#1288](https://github.com/getsops/sops/pull/1288))
* CI dependency updates ([#1295](https://github.com/getsops/sops/pull/1295), [#1301](https://github.com/getsops/sops/pull/1301))
* pgp: make error the last return value ([#1310](https://github.com/getsops/sops/pull/1310))
* Improve documentation files ([#1320](https://github.com/getsops/sops/pull/1320))

## 3.8.0

Features:

* Support `--version` without network requests using `--disable-version-check` ([#1115](https://github.com/getsops/sops/pull/1115))
* Support `--input-type` for updatekeys command ([#1116](https://github.com/getsops/sops/pull/1116))

Improvements:

* pgp: modernize and improve, and add tests ([#1054](https://github.com/getsops/sops/pull/1054),
  [#1282](https://github.com/getsops/sops/pull/1282))
* azkv: update SDK to latest, add tests, tidy ([#1067](https://github.com/getsops/sops/pull/1067),
  [#1092](https://github.com/getsops/sops/pull/1092), [#1256](https://github.com/getsops/sops/pull/1256))
* age: improve identity loading, add tests, tidy ([#1064](https://github.com/getsops/sops/pull/1064))
* kms: AWS SDK V2, allow creds config, add tests ([#1065](https://github.com/getsops/sops/pull/1065),
  [#1257](https://github.com/getsops/sops/pull/1257))
* gcpkms: update SDK to latest, add tests, tidy ([#1072](https://github.com/getsops/sops/pull/1072),
  [#1255](https://github.com/getsops/sops/pull/1255))
* hcvault: update API, add tests, tidy ([#1085](https://github.com/getsops/sops/pull/1085))
* Do not report version when upstream `--version` check fails ([#1124](https://github.com/getsops/sops/pull/1124))
* Use GitHub endpoints in `--version` command ([#1261](https://github.com/getsops/sops/pull/1261))
* Close temporary file before invoking editor to widen support on Windows ([#1265](https://github.com/getsops/sops/pull/1265))
* Update dependencies ([#1063](https://github.com/getsops/sops/pull/1063),
  [#1091](https://github.com/getsops/sops/pull/1091), [#1147](https://github.com/getsops/sops/pull/1147),
  [#1242](https://github.com/getsops/sops/pull/1242), [#1260](https://github.com/getsops/sops/pull/1260),
  [#1264](https://github.com/getsops/sops/pull/1264), [#1275](https://github.com/getsops/sops/pull/1275),
  [#1280](https://github.com/getsops/sops/pull/1280), [#1283](https://github.com/getsops/sops/pull/1283))
* Deal with various deprecations of dependencies ([#1113](https://github.com/getsops/sops/pull/1113),
  [#1262](https://github.com/getsops/sops/pull/1262))

Bug fixes:

* Ensure YAML comments are not displaced ([#1069](https://github.com/getsops/sops/pull/1069))
* Ensure default Google credentials can be used again after introduction
  of `GOOGLE_CREDENTIALS` ([#1249](https://github.com/getsops/sops/pull/1249))
* Avoid duplicate logging of errors in some key sources ([#1146](https://github.com/getsops/sops/pull/1146),
  [#1281](https://github.com/getsops/sops/pull/1281))
* Using `--set` on a root level key does no longer truncate existing values ([#899](https://github.com/getsops/sops/pull/899))
* Ensure stable order of SOPS parameters in dotenv file ([#1101](https://github.com/getsops/sops/pull/1101))

Project changes:

* Update Go to 1.20 ([#1148](https://github.com/getsops/sops/pull/1148))
* Update rustc functional tests to v1.70.0 ([#1234](https://github.com/getsops/sops/pull/1234))
* Remove remaining CircleCI workflow ([#1237](https://github.com/getsops/sops/pull/1237))
* Run CLI workflow on main ([#1243](https://github.com/getsops/sops/pull/1243))
* Delete obsolete `validation/` artifact ([#1248](https://github.com/getsops/sops/pull/1248))
* Rename Go module to `github.com/getsops/sops/v3` ([#1247](https://github.com/getsops/sops/pull/1247))
* Revamp release automation, including (Cosign) signed container images
  and checksums file, SLSA3 provenance and SBOMs ([#1250](https://github.com/getsops/sops/pull/1250))
* Update various bits of documentation ([#1244](https://github.com/getsops/sops/pull/1244))
* Add missing `--encrypt` flag from Vault example ([#1060](https://github.com/getsops/sops/pull/1060))
* Add documentation on how to use age in `.sops.yaml` ([#1192](https://github.com/getsops/sops/pull/1192))
* Improve Make targets and address various issues ([#1258](https://github.com/getsops/sops/pull/1258))
* Ensure clean working tree in CI ([#1267](https://github.com/getsops/sops/pull/1267))
* Fix CHANGELOG.rst formatting ([#1269](https://github.com/getsops/sops/pull/1269))
* Pin GitHub Actions to full length commit SHA and add CodeQL ([#1276](https://github.com/getsops/sops/pull/1276))
* Enable Dependabot for Docker, GitHub Actions and Go Mod ([#1277](https://github.com/getsops/sops/pull/1277))
* Generate versioned `.intoto.jsonl` ([#1278](https://github.com/getsops/sops/pull/1278))
* Update CI dependencies ([#1279](https://github.com/getsops/sops/pull/1279))

## 3.7.3

Changes:

* Upgrade dependencies ([#1024](https://github.com/getsops/sops/pull/1024), [#1045](https://github.com/getsops/sops/pull/1045))
* Build alpine container in CI ([#1018](https://github.com/getsops/sops/pull/1018),
  [#1032](https://github.com/getsops/sops/pull/1032), [#1025](https://github.com/getsops/sops/pull/1025))
* keyservice: accept KeyServiceServer in LocalClient ([#1035](https://github.com/getsops/sops/pull/1035))
* Add support for GCP Service Account within `GOOGLE_CREDENTIALS` ([#953](https://github.com/getsops/sops/pull/953))

Bug fixes:

* Upload the correct binary for the linux amd64 build ([#1026](https://github.com/getsops/sops/pull/1026))
* Fix bug when specifying multiple age recipients ([#966](https://github.com/getsops/sops/pull/966))
* Allow for empty yaml maps ([#908](https://github.com/getsops/sops/pull/908))
* Limit AWS role names to 64 characters ([#1037](https://github.com/getsops/sops/pull/1037))

## 3.7.2

Changes:

* README updates ([#861](https://github.com/getsops/sops/pull/861), [#860](https://github.com/getsops/sops/pull/860))
* Various test fixes ([#909](https://github.com/getsops/sops/pull/909),
  [#906](https://github.com/getsops/sops/pull/906), [#1008](https://github.com/getsops/sops/pull/1008))
* Added Linux and Darwin arm64 releases ([#911](https://github.com/getsops/sops/pull/911),
  [#891](https://github.com/getsops/sops/pull/891))
* Upgrade to go v1.17 ([#1012](https://github.com/getsops/sops/pull/1012))
* Support SOPS_AGE_KEY environment variable ([#1006](https://github.com/getsops/sops/pull/1006))

Bug fixes:

* Make sure comments in yaml files are not duplicated ([#866](https://github.com/getsops/sops/pull/866))
* Make sure configuration file paths work correctly relative to the
  config file in us ([#853](https://github.com/getsops/sops/pull/853))

## 3.7.1

Changes:

* Security fix
* Add release workflow ([#843](https://github.com/getsops/sops/pull/843))
* Fix issue where CI wouldn't run against master ([#848](https://github.com/getsops/sops/pull/848))
* Trim extra whitespace around age keys ([#846](https://github.com/getsops/sops/pull/846))

## 3.7.0

Features:

* Add support for age ([#688](https://github.com/getsops/sops/pull/688))
* Add filename to exec-file ([#761](https://github.com/getsops/sops/pull/761))

Changes:

* On failed decryption with GPG, return the error returned by GPG to the
  sops user ([#762](https://github.com/getsops/sops/pull/762))
* Use yaml.v3 instead of modified yaml.v2 for handling YAML files ([#791](https://github.com/getsops/sops/pull/791))
* Update aws-sdk-go to version v1.37.18 ([#823](https://github.com/getsops/sops/pull/823))

Project Changes:

* Switch from TravisCI to Github Actions ([#792](https://github.com/getsops/sops/pull/792))

## 3.6.1

Features:

* Add support for --unencrypted-regex ([#715](https://github.com/getsops/sops/pull/715))

Changes:

* Use keys.openpgp.org instead of gpg.mozilla.org ([#732](https://github.com/getsops/sops/pull/732))
* Upgrade AWS SDK version ([#714](https://github.com/getsops/sops/pull/714))
* Support --input-type for exec-file ([#699](https://github.com/getsops/sops/pull/699))

Bug fixes:

* Fixes broken Vault tests ([#731](https://github.com/getsops/sops/pull/731))
* Revert "Add standard newline/quoting behavior to dotenv store" ([#706](https://github.com/getsops/sops/pull/706))

## 3.6.0

Features:

* Support for encrypting data through the use of Hashicorp Vault ([#655](https://github.com/getsops/sops/pull/655))
* `sops publish` now supports `--recursive` flag for publishing all files
  in a directory ([#602](https://github.com/getsops/sops/pull/602))
* `sops publish` now supports `--omit-extensions` flag for omitting the
  extension in the destination path ([#602](https://github.com/getsops/sops/pull/602))
* sops now supports JSON arrays of arrays ([#642](https://github.com/getsops/sops/pull/642))

Improvements:

* Updates and standardization for the dotenv store ([#612](https://github.com/getsops/sops/pull/612),
  [#622](https://github.com/getsops/sops/pull/622))
* Close temp files after using them for edit command ([#685](https://github.com/getsops/sops/pull/685))

Bug fixes:

* AWS SDK usage now correctly resolves the `~/.aws/config` file ([#680](https://github.com/getsops/sops/pull/680))
* `sops updatekeys` now correctly matches config rules ([#682](https://github.com/getsops/sops/pull/682))
* `sops updatekeys` now correctly uses the config path cli flag ([#672](https://github.com/getsops/sops/pull/672))
* Partially empty sops config files don't break the use of sops anymore ([#662](https://github.com/getsops/sops/pull/662))
* Fix possible infinite loop in PGP's passphrase prompt call ([#690](https://github.com/getsops/sops/pull/690))

Project changes:

* Dockerfile now based off of golang version 1.14 ([#649](https://github.com/getsops/sops/pull/649))
* Push alpine version of docker image to Dockerhub ([#609](https://github.com/getsops/sops/pull/609))
* Push major, major.minor, and major.minor.patch tagged docker images to
  Dockerhub ([#607](https://github.com/getsops/sops/pull/607))
* Removed out of date contact information ([#668](https://github.com/getsops/sops/pull/668))
* Update authors in the cli help text ([#645](https://github.com/getsops/sops/pull/645))

## 3.5.0

Features:

* `sops exec-env` and `sops exec-file`, two new commands for utilizing sops
  secrets within a temporary file or env vars

Bug fixes:

* Sanitize AWS STS session name, as sops creates it based off of the machines hostname
* Fix for `decrypt.Data` to support `.ini` files
* Various package fixes related to switching to Go Modules
* Fixes for Vault-related tests running locally and in CI.

Project changes:

* Change to proper use of go modules, changing to primary module name to
  `go.mozilla.org/sops/v3`
* Change tags to requiring a `v` prefix.
* Add documentation for `sops updatekeys` command

## 3.4.0

Features:

* `sops publish`, a new command for publishing sops encrypted secrets to
  S3, GCS, or Hashicorp Vault
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

## 3.3.1

Bug fixes:

* Make sure the pgp key fingerprint is longer than 16 characters before
  slicing it. ([#463](https://github.com/getsops/sops/pull/463))
* Allow for `--set` value to be a string. ([#461](https://github.com/getsops/sops/pull/461))

Project changes:

* Using `develop` as a staging branch to create releases off of. What
  is in `master` is now the current stable release.
* Upgrade to using Go 1.12 to build sops
* Updated all vendored packages

## 3.3.0

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
* Fix issue with AWS KMS Encryption Contexts ([#445](https://github.com/getsops/sops/pull/445))
  with more than one context value failing to decrypt intermittently.
  Includes an automatic fix for old files affected by this issue.

Project infrastructure changes:

* Added integration tests for AWS KMS
* Added Code of Conduct

## 3.2.0

* Added --output flag to write output a file directly instead of
  through stdout
* Added support for dotenv files

## 3.1.1

* Fix incorrect version number from previous release

## 3.1.0

* Add support for Azure Key Service

* Fix bug that prevented JSON escapes in input files from working

## 3.0.5

* Prevent files from being encrypted twice

* Fix empty comments not being decrypted correctly

* If keyservicecmd returns an error, log it.

* Initial sops workspace auditing support (still wip)

* Refactor Store interface to reflect operations SOPS performs

## 3.0.3

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

## 3.0.1

* Don't consider io.EOF returned by Decoder.Token as error

* add IsBinary: true to FileHints when encoding with crypto/openpgp

* some improvements to error messages

## 3.0.0

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

* Added command to reconfigure the keys used to encrypt/decrypt a file based on
  the `.sops.yaml` config file

* Retrieve missing PGP keys from gpg.mozilla.org

* Improved error messages for errors when decrypting files

## 2.0.0

* [major] rewrite in Go

## 1.14

* [medium] Support AWS KMS Encryption Contexts
* [minor] Support insertion in encrypted documents via --set
* [minor] Read location of gpg binary from SOPS_GPG_EXEC env variables

## 1.13

* [minor] handle $EDITOR variable with parameters

## 1.12

* [minor] make sure filename_regex gets applied to file names, not paths
* [minor] move check of latest version under the -V flag
* [medium] fix handling of binary data to preserve file integrity
* [minor] try to use configuration when encrypting existing files
