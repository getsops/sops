Per-file example
================
This directory is an example configuration for SOPS inside of a project. We will cover the files used and relevant scripts for developers.

This example is optimized for storing sensitive information next to related non-sensitive information (e.g. password next to username).

The downsides include:

- Slowing down developers by requiring usage of SOPS for non-sensitive information
- Losing dynamic configurations that rely on reusing variables (e.g. ``test = {'foo': {'bar': common['foo']['bar'], 'baz': false}}``)

  - There might be work arounds via YAML

Getting started
---------------
To use this example, run the following

.. code:: bash

    # From the `sops` root directory
    # Import the test key
    gpg --import tests/sops_functional_tests_key.asc

    # Navigate to our example directory
    cd examples/per_file

    # Decrypt our secrets
    bin/decrypt-config.sh

    # Optionally edit a secret
    # bin/edit-secret.sh config.enc/static_github.json

    # Run our script
    python main.py

Storage
-------
In both development and production, we will be storing the secrets file unencrypted on disk. This is for a few reasons:

- Can't store file in an encrypted manner because we would need to know the secret to decode it
- Loading it into memory at boot is impractical

  - Requires reimplementing SOPS' decryption logic to multiple languages which increases chance of human error which is bad for security
  - If someone uses an automatic process reloader during development, then it could get expensive with AWS

    - We could cache the results from AWS but those secrets would wind up being stored on disk

As peace of mind, think about this:

- Unencrypted on disk is fine because if the attacker ever gains access to the server, then they can run ``sops --decrypt`` as well.

Files
-----
- ``bin/decrypt-config.sh`` - Script to decrypt secret file
- ``bin/edit-config-file.sh`` - Script to edit a secret file and then decrypt it
- ``config`` - Directory containing decrypted secrets
- ``config.bak`` - Backup of ``config`` to prevent accidental data loss
- ``config.enc`` - Directory containing encrypted secrets

  - ``static.py`` - Python script to merge together secrets
  - ``static_github.json`` - File containing secrets

- ``.gitignore`` - Ignore file for ``config`` and ``config.bak``
- ``main.py`` - Example script

Usage
-----
Development
~~~~~~~~~~~
For development, each developer must have access to the PGP/KMS keys. This means:

- If we are using PGP, then each developer must have the private key installed on their local machine
- If we are using KMS, then each developer must have AWS access to the appropriate key

Testing
~~~~~~~
For testing in a public CI, we can copy ``config.enc`` to ``config``. The secret files will have structure with an additional ``sops`` key but not reveal any secret information.

..

    For convenience, we can run ``CONFIG_COPY_ONLY=TRUE bin/decrypt-config.sh`` which will use ``ln -s`` rather than ``sops --decrypt``.

For testing in a private CI where we need private information, see the `Production instructions <#production>`_.

Production
~~~~~~~~~~
For production, we have a few options:

- Build an archive (e.g. ``.tar.gz``) in a private CI which contains the secrets and deploy our service via the archive
- Install PGP private key/KMS credentials on production machine, decrypt secrets during deployment process on production machine
