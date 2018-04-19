A Gomega release is a tagged sha and a GitHub release.  To cut a release:

1. Ensure CHANGELOG.md is uptodate.
2. Update GOMEGA_VERSION in `gomega_dsl.go`
3. Push a commit with the version number as the commit message (e.g. `v1.3.0`)
4. Create a new [GitHub release](https://help.github.com/articles/creating-releases/) with the version number as the tag  (e.g. `v1.3.0`).  List the key changes in the release notes.