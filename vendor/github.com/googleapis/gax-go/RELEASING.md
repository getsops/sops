# How to release v1

1. Determine the current release version with `git tag -l`. It should look
   something like `vX.Y.Z`. We'll call the current version `$CV` and the new
   version `$NV`.
1. On master, run `git log $CV..` to list all the changes since the last
   release.
   a. NOTE: Some commits may pertain to only v1 or v2. Manually introspect
   each commit to figure which occurred in v1.
1. Edit `CHANGES.md` to include a summary of the changes.
1. Mail the CL containing the `CHANGES.md` changes. When the CL is approved,
   submit it.
1. Without submitting any other CLs:
   a. Switch to master.
   b. `git pull`
   c. Tag the repo with the next version: `git tag $NV`. It should be of the
   form `v1.Y.Z`.
   d. Push the tag: `git push origin $NV`.
1. Update [the releases page](https://github.com/googleapis/google-cloud-go/releases)
   with the new release, copying the contents of the CHANGES.md.

# How to release v2

Same process as v1, once again noting that the commit list may include v1
commits (which should be pruned out). Note also whilst v1 tags are `v1.Y.Z`, v2
tags are `v2.Y.Z`.

# On releasing multiple major versions

Please see https://github.com/golang/go/wiki/Modules#releasing-modules-v2-or-higher.
