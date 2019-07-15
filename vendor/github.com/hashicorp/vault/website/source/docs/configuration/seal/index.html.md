---
layout: "docs"
page_title: "Seals - Configuration"
sidebar_title: "<code>seal</code>"
sidebar_current: "docs-configuration-seal"
description: |-
  The seal stanza configures the seal type to use for additional data protection.
---

# `seal` Stanza

The `seal` stanza configures the seal type to use for additional data
protection, such as using HSM or Cloud KMS solutions to encrypt and decrypt the
master key. This stanza is optional, and in the case of the master key, Vault
will use the Shamir algorithm to cryptographically split the master key if this
is not configured.

As of Vault 0.9.0, the seal can also be used for [seal wrapping][sealwrap] to
add an extra layer of protection and satisfy compliance and regulatory requirements. 
This feature is only available in Vault Enterprise.

For more examples, please choose a specific auto unsealing technology from the
sidebar.

## Configuration

Seal configuration can be done through the Vault configuration file using the
`seal` stanza:

```hcl
seal [NAME] {
  # ...
}
```

For example:

```hcl
seal "pkcs11" {
  # ...
}
```

For configuration options which also read an environment variable, the
environment variable will take precedence over values in the configuration file.

[sealwrap]: /docs/enterprise/sealwrap/index.html
