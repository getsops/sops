# Key groups

By default, `sops` encrypts its [data key](internals/encryption_protocol.md)
with each of the [master keys](encryption_providers/README.md), such that if
any of the master keys is available, the file can be decrypted. However, it is
sometimes desirable to require access to multiple master keys in order to
decrypt files. This can be achieved with key groups.

When using key groups in sops, the data key is split into fragments such that
master keys from multiple groups are required to decrypt a file. `sops` uses
[Shamir's Secret Sharing](TODO) to split the data key such that each key group
has a fragment, each key in the key group can decrypt that fragment, and a
configurable number of fragments (threshold) are needed to decrypt and piece
together the complete data key. When decrypting a file using multiple key
groups, `sops` goes through the key groups in order, and for each group, tries
to recover its corresponding fragment of the data key using a master key from
that group. Once the fragment is recovered, `sops` moves on to the next group,
until enough fragments have been recovered to obtain the complete data key.

By default, the threshold is set to the number of key groups. For example, if
you have three key groups configured in your SOPS file and you don't override
the default threshold, then one master key from each of the three groups will
be required to decrypt the file.

## The `sops groups` subcommand

Management of key groups is performed through the `sops groups` subcommand.

### Adding a key group to a file

You can add groups to a file through the `sops groups add` command. For example, to add a new key group to an existing file with 3 PGP keys and 3 AWS KMS keys, one could do as follows:

```bash
$ sops groups add --file my_file.enc --pgp fingerprint1 --pgp fingerprint2 --pgp fingerprint3 --kms arn1 --kms arn2 --kms arn3
```

### Deleting key groups from a file

Similarly, key groups can be deleted from encrypted files by their index with
the `sops groups delete` command. For example, to delete the first key group
(group number 0, as groups are zero-indexed) in a file:

```bash
$ sops groups delete --file my_file.enc 0
```

?> Currently, you can see the order of the key groups in a file by inspecting
the encrypted file directly without using `sops` (i.e., with `cat`, `less` or
your favorite editor)

## `.sops.yaml`

The [`.sops.yaml` configuration file](sops_yaml_config_file.md) also contains options related to
key groups, and as such it can be used as an alternative to the `sops groups`
subcommand.

When creating a new file using the rules on `.sops.yaml`, the Shamir threshold
can also be specified through the `--shamir-secret-sharing-threshold` flag.

## The `sops updatekeys` subcommand

In some cases, it might be more convenient to use rules like those in
`.sops.yaml` to manage key groups. However, these only apply to new files. The
`sops updatekeys` subcommand updates the key groups in existing files based on
the matching rules in the `.sops.yaml` configuration file. Essentially, this
updates the key groups in the file such that they end up as if the file was
newly created.

```bash
$ sops updatekeys example.yaml
2019/07/18 16:03:43 Syncing keys for file example.yaml
The following changes will be made to the file's groups:
Group 1
    1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A
--- 85D77543B3D624B63CEA9E6DBC17301B491B3F21
Is this okay? (y/n):y
2019/07/18 16:03:53 File example.yaml synced with new keys
```

This can be more convenient than the `sops groups` subcommand when, for
example, we want to delete a single key from a key group.
