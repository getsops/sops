# Publishing files

The `sops publish` command publishes a file to a pre-configured destination.
These destinations typically feed into production systems that consume the
secrets.

Destination configuration resides in the [`.sops.yaml` configuration
file](sops_yaml_config_file.md).

## Publication targets

A variety of publication targets are supported. We make a distinction between
publication targets that manage at-rest encryption themselves and those that do
not. 

// TODO this is confusing
// For instance, why is Vault secure enough to store plain text secrets, but S3 isn't?
// After all, S3 encrypts at rest as well, and also has access controls around it.
// I think this might have been designed with the needs of Mozilla in mind.

When the target does not manage encryption itself, SOPS will still be
responsible for keeping the file encrypted. As such, when publishing to these
targets, SOPS offers the option to reencrypt the files with a new set of keys.
Typically, you'd reencrypt the file with keys that only production systems have
access to. Recreation rules are supported in the [`.sops.yaml` configuration
file](sops_yaml_config_file.md).

For targets that manage encryption themselves, SOPS stores the plain-text,
unencrypted data on the target, and the target is responsible for ensuring the
data is stored securely, encrypted at rest and with appropriate access
controls.

### AWS S3

SOPS can publish files to S3 buckets.

?> S3 does *not* manage encryption itself

### Google Cloud Storage

SOPS can publish files to Google Cloud Storage buckets.

?> Google Cloud Storage does *not* manage encryption itself

### Hashicorp Vault
