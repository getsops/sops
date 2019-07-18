# `sops` as a git differ

You most likely want to store your encrypted files in a version control
repository. `sops` can be used with git to decrypt files when showing diffs
between versions. This is very handy for reviewing changes or visualizing
history.

To configure sops to decrypt files during diff, create a `.gitattributes` file
at the root of your repository that contains a filter and a command:

```
*.yaml diff=sopsdiffer
```

Here we only care about YAML files. ``sopsdiffer`` is an arbitrary name that we
then map to a `sops` command in the git configuration file of the repository.

```bash
$ git config diff.sopsdiffer.textconv "sops -d"
```

With this in place, calls to ``git diff`` will decrypt both previous and
current versions of the target file prior to displaying the diff. And it even
works with git client interfaces, because they call git diff under the hood!
