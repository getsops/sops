SOPS: Secrets OPerationS
========================

**sop** is an editor of encrypted files that supports YAML, JSON and TEXT formats and encrypts with AWS KMS and PGP (via GnuPG). Watch `the demo <https://www.youtube.com/watch?v=YTEVyLXFiq0>`_.

.. image:: http://i.imgur.com/IL6dlhm.gif

.. image:: https://travis-ci.org/mozilla/sops.svg?branch=master
	:target: https://travis-ci.org/mozilla/sops

.. sectnum::
.. contents:: Table of Contents

Requirements
------------
First install some libraries from your package manager:

* RHEL family::

	sudo yum install libyaml-devel python-devel libffi-devel pip

* Debian family::

	sudo apt-get install libyaml-dev python-dev libffi-dev python-pip

* MacOS::

	brew install libffi libyaml
	sudo easy_install pip

Then install `sops` from pip::

	sudo pip install --upgrade sops

note: on centos/sl, you may need to upgrade `botocore` after installing
sops to deal with a `requirement conflict
<https://github.com/boto/botocore/issues/660>`_.
Do so with `sudo pip install -U botocore`.

Usage
-----

If you're using AWS KMS, create one or multiple master keys in the IAM console
and export them, comma separated, in the **SOPS_KMS_ARN** env variable. It is
recommended to use at least two master keys in different regions.

.. code:: bash

	export SOPS_KMS_ARN="arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e,arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"

Your AWS credentials must be present in `~/.aws/credentials`. sops uses boto3.

.. code::

	$ cat ~/.aws/credentials 
	[default]
	aws_access_key_id = AKI.....
	aws_secret_access_key = mw......

If you want to use PGP, export the fingerprints of the public keys, comma
separated, in the **SOPS_PGP_FP** env variable.

.. code:: bash

	export SOPS_PGP_FP="85D77543B3D624B63CEA9E6DBC17301B491B3F21,E60892BB9BD89A69F759A1A0A3D652173B763E8F"

Note: you can use both PGP and KMS simultaneously.

Then simply call `sops` with a file path as argument. It will handle the
encryption/decryption transparently and open the cleartext file in an editor.

.. code:: bash

	$ sops mynewtestfile.yaml
	mynewtestfile.yaml doesn't exist, creating it.
	please wait while an encryption key is being generated and stored in a secure fashion
	[... editing happens in vim, or whatever $EDITOR is set to ...]
	file written to mynewtestfile.yaml

The resulting encrypted file looks like this:

.. code:: yaml

    myapp1: ENC[AES256_GCM,data:Tr7o=,iv:1=,aad:No=,tag:k=]
    app2:
        db:
            user: ENC[AES256_GCM,data:CwE4O1s=,iv:2k=,aad:o=,tag:w==]
            password: ENC[AES256_GCM,data:p673w==,iv:YY=,aad:UQ=,tag:A=]
        # private key for secret operations in app2
        key: |-
            ENC[AES256_GCM,data:Ea3kL5O5U8=,iv:DM=,aad:FKA=,tag:EA==]
    an_array:
    - ENC[AES256_GCM,data:v8jQ=,iv:HBE=,aad:21c=,tag:gA==]
    - ENC[AES256_GCM,data:X10=,iv:o8=,aad:CQ=,tag:Hw==]
    - ENC[AES256_GCM,data:KN=,iv:160=,aad:fI4=,tag:tNw==]
    sops:
        kms:
        -   created_at: 1441570389.775376
            enc: CiC....Pm1Hm
            arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
        -   created_at: 1441570391.925734
            enc: Ci...awNx
            arn: arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d
        pgp:
        -   fp: 85D77543B3D624B63CEA9E6DBC17301B491B3F21
            created_at: 1441570391.930042
            enc: |
                -----BEGIN PGP MESSAGE-----
                hQIMA0t4uZHfl9qgAQ//UvGAwGePyHuf2/zayWcloGaDs0MzI+zw6CmXvMRNPUsA
				...=oJgS
                -----END PGP MESSAGE-----

A copy of the encryption/decryption key is stored securely in each KMS and PGP
block. As long as one of the KMS or PGP method is still usable, you will be able
to access you data.

To decrypt a file in a `cat` fashion, use the `-d` flag:

.. code:: bash

	$ sops -d mynewtestfile.yaml

`sops` encrypted files contain the necessary information to decrypt their content.
All a user of `sops` need is valid AWS credentials and the necessary
permissions on KMS keys.

Given that, the only command a `sops` user need is:

.. code:: bash

	$ sops <file>

`<file>` will be opened, decrypted, passed to a text editor (vim by default),
encrypted if modified, and saved back to its original location. All of these
steps, apart from the actual editing, are transparent to the user.

Adding and removing keys
~~~~~~~~~~~~~~~~~~~~~~~~

When creating a new files, `sops` uses the PGP and KMS defined in the command
line arguments `--kms` and `--pgp`, or from the environment variables
`SOPS_KMS_ARN` and `SOPS_PGP_FP`. That information is stored in the file under
the `sops` section. When editing a file, it is trivial to add or remove keys:
invoke `sops` with the flag **-s** to display the master keys while editing, and
add or remove kms or pgp keys under the sops section.

For example, to add a KMS master key to a file, we would add the following
entry:

.. code:: yaml

	sops:
	    kms:
	    - arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e

And, similarly, to add a PGP master key, we add its fingerprint:

.. code:: yaml

	sops:
	    pgp:
	    - fp: 85D77543B3D624B63CEA9E6DBC17301B491B3F21

When the file is saved, `sops` will update its metadata and encrypt the data key
with the freshly added master keys. The removed entries are simply deleted from
the file.

Using KMS master keys in various AWS accounts
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can use keys in various accounts by tying each KMS master key to a role that
the user is allowed to assume in each account. The `IAM roles
<http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use.html>`_
documentation has full details on how this needs to be configured on AWS's side.
From the point of view of `sops`, you only need to specify the role a KMS key
must assume alongside its ARN, as follows:

.. code:: yaml

	sops:
	    kms:
	    -	arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
	        role: arn:aws:iam::927034868273:role/sops-dev-xyz

The role must have permission to call Encrypt and Decrypt using KMS. An example
policy is shown below.

.. code:: json

	{
	  "Sid": "Allow use of the key",
	  "Effect": "Allow",
	  "Action": [
		"kms:Encrypt",
		"kms:Decrypt",
		"kms:ReEncrypt*",
		"kms:GenerateDataKey*",
		"kms:DescribeKey"
	  ],
	  "Resource": "*",
	  "Principal": {
		"AWS": [
		  "arn:aws:iam::927034868273:role/sops-dev-xyz"
		]
	  }
	}

You can specify a role in the `--kms` flag and `SOPS_KMS_ARN` variable by
appending it to the ARN of the master key, separated by a **+** sign::

	<KMS ARN>+<ROLE ARN>
	arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500+arn:aws:iam::927034868273:role/sops-dev-xyz

Key Rotation
~~~~~~~~~~~~

It is recommend to renew the data key on a regular basis. `sops` supports key
rotation via the `-r` flag. A simple approach is to decrypt and reencrypt all
files in place with rotation enabled:

.. code:: bash

	for file in $(find . -type f -name "*.yaml"); do
		sops -d -i $file
		sops -e -i -r $file
	done

Note on YAML
------------

`sops` is designed to encrypt files that contain secrets, which are most likely
strings or numbers. It will not work on complex YAML files that use references
or anchors.

Cryptographic details
---------------------

When sops creates a file, it generates a random 256 bits data key and asks each
KMS and PGP master key to encrypt the data key. The encrypted version of the data
key is stored in the `sops` metadata under `sops.kms` and `sops.pgp`.

For KMS:

.. code:: yaml

    sops:
        kms:
        -   enc: CiC6yCOtzsnFhkfdIslYZ0bAf//gYLYCmIu87B3sy/5yYxKnAQEBAQB4usgjrc7JxYZH3SLJWGdGwH//4GC2ApiLvOwd7Mv+cmMAAAB+MHwGCSqGSIb3DQEHBqBvMG0CAQAwaAYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAyGdRODuYMHbA8Ozj8CARCAO7opMolPJUmBXd39Zlp0L2H9fzMKidHm1vvaF6nNFq0ClRY7FlIZmTm4JfnOebPseffiXFn9tG8cq7oi
            enc_ts: 1439568549.245995
            arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e

For PGP:

.. code:: yaml

    sops:
        pgp:
        -   fp: 85D77543B3D624B63CEA9E6DBC17301B491B3F21
            created_at: 1441570391.930042
            enc: |
                -----BEGIN PGP MESSAGE-----
                Version: GnuPG v1

                hQIMA0t4uZHfl9qgAQ//UvGAwGePyHuf2/zayWcloGaDs0MzI+zw6CmXvMRNPUsA
                pAgRKczJmDu4+XzN+cxX5Iq9xEWIbny9B5rOjwTXT3qcUYZ4Gkzbq4MWkjuPp/Iv
                qO4MJaYzoH5YxC4YORQ2LvzhA2YGsCzYnljmatGEUNg01yJ6r5mwFwDxl4Nc80Cn
                RwnHuGExK8j1jYJZu/juK1qRbuBOAuruIPPWVdFB845PA7waacG1IdUW3ZtBkOy3
                O0BIfG2ekRg0Nik6sTOhDUA+l2bewCcECI8FYCEjwHm9Sg5cxmP2V5m1mby+uKAm
                kewaoOyjbmV1Mh3iI1b/AQMr+/6ZE9MT2KnsoWosYamFyjxV5r1ZZM7cWKnOT+tu
                KOvGhTV1TeOfVpajNTNwtV/Oyh3mMLQ0F0HgCTqomQVqw5+sj7OWAASuD3CU/dyo
                pcmY5Qe0TNL1JsMNEH8LJDqSh+E0hsUxdY1ouVsg3ysf6mdM8ciWb3WRGxih1Vmf
                unfLy8Ly3V7ZIC8EHV8aLJqh32jIZV4i2zXIoO4ZBKrudKcECY1C2+zb/TziVAL8
                qyPe47q8gi1rIyEv5uirLZjgpP+JkDUgoMnzlX334FZ9pWtQMYW4Y67urAI4xUq6
                /q1zBAeHoeeeQK+YKDB7Ak/Y22YsiqQbNp2n4CKSKAE4erZLWVtDvSp+49SWmS/S
                XgGi+13MaXIp0ecPKyNTBjF+NOw/I3muyKr8EbDHrd2XgIT06QXqjYLsCb1TZ0zm
                xgXsOTY3b+ONQ2zjhcovanDp7/k77B+gFitLYKg4BLZsl7gJB12T8MQnpfSmRT4=
                =oJgS
                -----END PGP MESSAGE-----

sops then opens a text editor on the newly created file. The user adds data to the
file and saves it when done.

Upon save, sops browses the entire file as of a key/value tree. Every time sops
encounters a leaf value (a value that does not have children), it encrypts the
value with AES256_GCM using the data key and a 256 bits random initialization
vector.

Each file uses a single data key to encrypt all values of a document, but each
value receives a unique initialization vector and has unique authentication data.

Additional data is used to guarantee the integrity of the encrypted data
and of the tree structure: when encrypting the tree, key names are concatenated
into a byte string that is used as AEAD additional data (aad) when encrypting
the value. The `aad` field is not stored with the value but reconstructed from
the tree structure every time.

The result of AES256_GCM encryption is stored in the leaf of the tree using a
base64 encoded string format::

    ENC[AES256_GCM,
        data:CwE4O1s=,
        iv:S0fozGAOxNma/pWDUuk1iEaYw0wlba0VOLHjPxIok2k=,
        tag:XaGsYaL9LCkLWJI0uxnTYw==]

where:

* **data** is the encrypted value
* **iv** is the 256 bits initialization vector
* **tag** is the authentication tag

The encrypted file is written to disk with nested keys in cleartext and
values encrypted. We expect that keys do not carry sensitive information, and
keeping them in cleartext allows for better diff and overall readability.

Any valid KMS or PGP master key can later decrypt the data key and access the
data.

Multiple master keys allow for sharing encrypted files without sharing master
keys, and provide disaster recovery solution. The recommended way to use sops
is to have two KMS master keys in different region and one PGP public key with
the private key stored offline. If, by any chance, both KMS master keys are
lost, you can always recover the encrypted data using the PGP private key.

Message Authentication Code
~~~~~~~~~~~~~~~~~~~~~~~~~~~

In addition to authenticating branches of the tree using keys as additional
data, sops computes a MAC on all the values to ensure that no value has been
added or removed fraudulently. The MAC is stored encrypted with AES_GCM and
the data key under tree->`sops`->`mac`.

Threat Model
------------

The security of the data stored using sops is as strong as the weakest
cryptographic mechanism. Values are encrypted using AES256_GCM which is the
strongest symetric encryption algorithm known today. Data keys are encrypted
in either KMS, which also uses AES256_GCM, or PGP which uses either RSA or
ECDSA keys. 

Going from the most likely to the least likely, the threats are as follow:

1. Compromised AWS credentials grant access to KMS master key
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

An attacker with access to an AWS console can grant itself access to one of
the KMS master key used to encrypt a sops data key. This threat should be
mitigated by protecting AWS accesses with strong controls, such as multi-factor
authentication, and also by performing regular audits of permissions granted
to AWS users.

2. Compromised PGP key
~~~~~~~~~~~~~~~~~~~~~~

PGP keys are routinely mishandled, either because owners copy them from
machine to machine, or because the key is left forgotten on an unused machine
an attacker gains access to. When using PGP encryption, sops users should take
special care of PGP private keys, and store them on smart cards or offline
as often as possible.

3. Factorized RSA key
~~~~~~~~~~~~~~~~~~~~~

sops doesn't apply any restriction on the size or type of PGP keys. A weak PGP
keys, for example 512 bits RSA, could be factorized by an attacker to gain
access to the private key and decrypt the data key. Users of sops should rely
on strong keys, such as 2048+ bits RSA keys, or 256+ bits ECDSA keys.

4. Weak AES cryptography
~~~~~~~~~~~~~~~~~~~~~~~~

A vulnerability in AES256_GCM could potentially leak the data key or the KMS
master key used by a sops encrypted file. While no such vulnerability exists
today, we recommend that users keep their encrypted files reasonably private.

License
-------
Mozilla Public License Version 2.0

Authors
-------
* Julien Vehent <jvehent@mozilla.com>
* Daniel Thornton <dthornton@mozilla.com>
* Alexis Metaireau <alexis@mozilla.com>
* Rémy Hubscher <natim@mozilla.com>

Credits
-------

`sops` is inspired by `hiera-eyaml <https://github.com/TomPoulton/hiera-eyaml>`_,
`credstash <https://github.com/LuminalOSS/credstash>`_ ,
`sneaker <https://github.com/codahale/sneaker>`_,
`password store <http://www.passwordstore.org/>`_ and too many years managing
PGP encrypted files by hand...
