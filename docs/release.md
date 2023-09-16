# Release procedure

This document describes the procedure for releasing a new version of SOPS. It
is intended for maintainers of the project, but may be useful for anyone
interested in the release process.

## Overview

The release is performed by creating a signed tag for the release, and pushing
it to GitHub. This will automatically trigger a GitHub Actions workflow that
builds the binaries, packages, SBOMs, and other artifacts for the release
using [GoReleaser](https://goreleaser.com), and uploads them to GitHub.

The configuration for GoReleaser is in the file
[`.goreleaser.yaml`](../.goreleaser.yaml). The configuration for the GitHub
Actions workflow is in the file
[`release.yml`](../.github/workflows/release.yml).

This configuration is quite sophisticated, and ensures at least the following:

- The release is built for multiple platforms and architectures, including
  Linux, macOS, and Windows, and for both AMD64 and ARM64.
- The release includes multiple packages in Debian and RPM formats.
- For every binary, a corresponding SBOM is generated and published.
- For all binaries, a checksum file is generated and signed using
  [Cosign](https://docs.sigstore.dev/cosign/overview/) with GitHub OIDC.
- Both Debian and Alpine Docker multi-arch images are built and pushed to GitHub
  Container Registry and Quay.io.
- The container images are signed using
  [Cosign](https://docs.sigstore.dev/cosign/overview/) with GitHub OIDC.
- [SLSA provenance](https://slsa.dev/provenance/v0.2) metadata is generated for
  release artifacts and container images.

## Preparation

- [ ] Ensure that all changes intended for the release are merged into the
  `main` branch. At present, this means that all pull requests attached to the
  milestone for the release are merged. If there are any pull requests that
  should not be included in the release, move them to a different milestone.
- [ ] Create a pull request to update the [`CHANGELOG.rst`](../CHANGELOG.rst)
  file. This should include a summary of all changes since the last release,
  including references to any relevant pull requests.
- [ ] In this same pull request, update the version number in `version/version.go`
  to the new version number.
- [ ] Get approval for the pull request from at least one other maintainer, and
  merge it into `main`.
- [ ] Ensure CI passes on the `main` branch.

## Release

- [ ] Ensure your local copy of the `main` branch is up-to-date:

  ```sh
  git checkout main
  git pull
  ```

- [ ] Create a **signed tag** for the release, using the following command:

  ```sh
  git tag -s -m <version> <version>
  ```

  where `<version>` is the version number of the release. The version number
  should be in the form `vX.Y.Z`, where `X`, `Y`, and `Z` are integers. The
  version number should be incremented according to
  [semantic versioning](https://semver.org/).
- [ ] Push the tag to GitHub:

  ```sh
  git push origin <version>
  ```

- [ ] Ensure the release is built successfully on GitHub Actions. This will
  automatically create a release on GitHub.
