SOPS: Secrets OPerationS
========================

**sop** is an editor of encrypted files that supports YAML, JSON and BINARY formats and encrypts with AWS KMS and PGP (via GnuPG). Watch `the demo <https://www.youtube.com/watch?v=YTEVyLXFiq0>`_.

.. image:: http://i.imgur.com/IL6dlhm.gif

.. image:: https://travis-ci.org/mozilla/sops.svg?branch=master
	:target: https://travis-ci.org/mozilla/sops

**Questions?** ping "ulfr" in `#security` on `irc.mozilla.org <https://wiki.mozilla.org/IRC>`_
(use a web client like `mibbit <https://chat.mibbit.com>`_ ).

.. sectnum::
.. contents:: Table of Contents

Installation
------------

* RHEL family::

	sudo yum install gcc git libffi-devel libyaml-devel make openssl openssl-devel python-devel python-pip
	sudo pip install --upgrade sops

* Debian family::

	sudo apt-get install gcc git libffi-dev libssl-dev libyaml-dev make openssl python-dev python-pip
	sudo pip install --upgrade sops

* MacOS Brew Install::

	brew install sops

* MacOS Manual Install::

	brew install libffi libyaml python [1]
	pip install sops

1. http://docs.python-guide.org/en/latest/starting/install/osx/#doing-it-right

In a virtualenv
~~~~~~~~~~~~~~~

Assuming you already have libffi and libyaml installed, the following commands will install sops in a virtualenv:

.. code:: bash

    $ sudo pip install virtualenv --upgrade
    $ virtualenv ~/sopsvenv
    $ source ~/sopsvenv/bin/activate
    $ pip install -U sops
    $ sops -v
    sops 1.9

Test with the dev PGP key
~~~~~~~~~~~~~~~~~~~~~~~~~
Clone the repository, load the test PGP key and open the test files::

	$ git clone https://github.com/mozilla/sops.git
	$ cd sops
	$ gpg --import tests/sops_functional_tests_key.asc
	$ sops example.yaml

This last step will decrypt `example.yaml` using the test private key. To create
your own secrets files using keys under your control, keep reading.

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
to access your data.

To decrypt a file in a `cat` fashion, use the `-d` flag:

.. code:: bash

	$ sops -d mynewtestfile.yaml

`sops` encrypted files contain the necessary information to decrypt their content.
All a user of `sops` needs is valid AWS credentials and the necessary
permissions on KMS keys.

Given that, the only command a `sops` user needs is:

.. code:: bash

	$ sops <file>

`<file>` will be opened, decrypted, passed to a text editor (vim by default),
encrypted if modified, and saved back to its original location. All of these
steps, apart from the actual editing, are transparent to the user.

Adding and removing keys
~~~~~~~~~~~~~~~~~~~~~~~~

When creating new files, `sops` uses the PGP and KMS defined in the command
line arguments `--kms` and `--pgp`, or from the environment variables
`SOPS_KMS_ARN` and `SOPS_PGP_FP`. That information is stored in the file under
the `sops` section, such that decrypting files does not require providing those
parameters again.

Master PGP and KMS keys can be added and removed from a `sops` file in one of
two ways: by using command line flag, or by editing the file directly.

Command line flag `--add-kms`, `--add-pgp`, `--rm-kms` and `--rm-pgp` can be
used to add and remove keys from a file. These flags use the comma separated
syntax as the `--kms` and `--pgp` arguments when creating new files.

.. code:: bash

	# add a new pgp key to the file and rotate the data key
	$ sops -r --add-pgp 85D77543B3D624B63CEA9E6DBC17301B491B3F21 example.yaml

	# remove a pgp key from the file and rotate the data key
	$ sops -r --rm-pgp 85D77543B3D624B63CEA9E6DBC17301B491B3F21 example.yaml

Alternatively, invoking `sops` with the flag **-s** will display the master keys
while editing. This method can be used to add or remove kms or pgp keys under the
sops section.

For example, to add a KMS master key to a file, add the following entry while
editing:

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

When removing keys, it is recommended to rotate the data key using `-r`,
otherwise owners of the removed key may have add access to the data key in the
past.

Assuming roles and using KMS in various AWS accounts
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

SOPS has the ability to use KMS in multiple AWS accounts by assuming roles in
each account. Being able to assume roles is a nice feature of AWS that allows
administrators to establish trust relationships between accounts, typically from
the most secure account to the least secure one. In our use-case, we use roles
to indicate that a user of the Master AWS account is allowed to make use of KMS
master keys in development and staging AWS accounts. Using roles, a single file
can be encrypted with KMS keys in multiple accounts, thus increasing reliability
and ease of use.

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

It is recommended to renew the data key on a regular basis. `sops` supports key
rotation via the `-r` flag. Invoking it on an existing file causes sops to
reencrypt the file with a new data key, which is then encrypted with the various
KMS and PGP master keys defined in the file.

.. code:: bash

	sops -r example.yaml

Using .sops.yaml conf to select KMS/PGP for new files
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

It is often tedious to specify the `--kms` and `--pgp` parameters for creation
of all new files. If your secrets are stored under a specific directory, like a
`git` repository, you can create a `.sops.yaml` configuration file at the root
directory to define which keys are used for which filename.

Let's take an example:

* file named **something.dev.yaml** should use one set of KMS A
* file named **something.prod.yaml** should use another set of KMS B
* other files use a third set of KMS C
* all live under **mysecretrepo/something.{dev,prod}.yaml**

Under those circumstances, a file placed at **mysecretrepo/.sops.yaml**
can manage the three sets of configurations for the three types of files:

.. code:: yaml

	# creation rules are evaluated sequentially, the first match wins
	creation_rules:
		# upon creation of a file that matches the pattern *.dev.yaml,
		# KMS set A is used
		- filename_regex: \.dev\.yaml$
		  kms: 'arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500,arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e+arn:aws:iam::361527076523:role/hiera-sops-prod'
		  pgp: '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'

		# prod files use KMS set B in the PROD IAM
		- filename_regex: \.prod\.yaml$
		  kms: 'arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e+arn:aws:iam::361527076523:role/hiera-sops-prod,arn:aws:kms:eu-central-1:361527076523:key/cb1fab90-8d17-42a1-a9d8-334968904f94+arn:aws:iam::361527076523:role/hiera-sops-prod'
		  pgp: '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'

		# Finally, if the rules above have not matched, this one is a
		# catchall that will encrypt the file using KMS set C
		# The absence of a filename_regex means it will match everything
		- kms: 'arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500,arn:aws:kms:us-west-2:142069644989:key/846cfb17-373d-49b9-8baf-f36b04512e47,arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e'
		  pgp: '1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A'

When creating any file under **mysecretrepo**, whether at the root or under
a subdirectory, sops will recursively look for a `.sops.yaml` file. If one is
found, the filename of the file being created is compared with the filename
regexes of the configuration file. The first regex that matches is selected,
and its KMS and PGP keys are used to encrypt the file.

Creating a new file with the right keys is now as simple as

.. code:: bash

	$ sops <newfile>.prod.yaml

Note that the configuration file is ignored when KMS or PGP parameters are
passed on the sops command line or in environment variables.

Important information on types
------------------------------

YAML and JSON type extensions
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

`sops` uses the file extension to decide which encryption method to use on the file
content. `YAML` and `JSON` files are treated as trees of data, and key/values are
extracted from the files to only encrypt the leaf values. The tree structure is also
used to check the integrity of the file.

Therefore, if a file is encrypted using a specific format, it need to be decrypted
in the same format. The easiest way to achieve this is to conserve the original file
extension after encrypting a file. For example::

	$ sops -e -i myfile.json

	$ sops -d myfile.json

If you want to change the extension of the file once encrypted, you need to provide
sops with the `--input-type` flag upon decryption. For example::

	$ sops -e myfile.json > myfile.json.enc

	$ sops -d --input-type json myfile.json.enc

YAML anchors
~~~~~~~~~~~~
`sops` only supports a subset of `YAML`'s many types. Encrypting YAML files that
contain strings, numbers and booleans will work fine, but files that contain anchors
will not work, because the anchors redefine the structure of the file at load time.

This file will not work in `sops`:

.. code:: yaml

	bill-to:  &id001
	    street: |
	        123 Tornado Alley
	        Suite 16
	    city:   East Centerville
	    state:  KS

	ship-to:  *id001

`sops` uses the path to a value as additional data in the AEAD encryption, and thus
dynamic paths generated by anchors break the authentication step.

JSON and TEXT file types do not support anchors and thus have no such limitation.

Top-level arrays
~~~~~~~~~~~~~~~~
`YAML` and `JSON` top-level arrays are not supported, because `sops` needs a top-level
`sops` key to store its metadata.
This file will not work in sops:

.. code:: yaml

	---
	  - some
	  - array
	  - elements

But this one will because because the `sops` key can be added at the same level as the
`data` key.

.. code:: yaml

	data:
	  - some
	  - array
	  - elements

Similarly, with `JSON` arrays, this document will not work:

.. code:: json

	[
	  "some",
	  "array",
	  "elements"
	]


But this one will work just fine:

.. code:: json

	{
	  "data": [
	    "some",
	    "array",
	    "elements"
	  ]
	}


Examples
--------

Take a look into the `examples <https://github.com/mozilla/sops/tree/master/examples>`_ folder for detailed use cases of sops in a CI environment. The section below describes specific tips for common use cases.

Creating a new file
~~~~~~~~~~~~~~~~~~~

The command below creates a new file with a data key encrypted by KMS and PGP.

.. code:: bash

	$ sops --kms "arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500" --pgp C9CAB0AF1165060DB58D6D6B2653B624D620786D /path/to/new/file.yaml

Encrypting an existing file
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Similar to the previous command, we tell sops to use one KMS and one PGP key.
The path points to an existing cleartext file, so we give sops flag `-e` to
encrypt the file, and redirect the output to a destination file.

.. code:: bash

	$ export SOPS_KMS_ARN="arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500"
	$ export SOPS_PGP_FP="C9CAB0AF1165060DB58D6D6B2653B624D620786D"
	$ sops -e /path/to/existing/file.yaml > /path/to/new/encrypted/file.yaml

Decrypt the file with `-d`.

.. code:: bash

	$ sops -d /path/to/new/encrypted/file.yaml

Encrypt or decrypt a file in place
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Rather than redirecting the output of `-e` or `-d`, sops can replace the
original file after encrypting or decrypting it.

.. code:: bash

	# file.yaml is in cleartext
	$ sops -e -i /path/to/existing/file.yaml
	# file.yaml is now encrypted
	$ sops -d -i /path/to/existing/file.yaml
	# file.yaml is back in cleartext

Encrypting binary files
~~~~~~~~~~~~~~~~~~~~~~~

`sops` primary use case is encrypting YAML and JSON configuration files, but it
also has the ability to manage binary files. When encrypting a binary, sops will
read the data as bytes, encrypt it, store the encrypted base64 under
`tree['data']` and write the result as JSON.

Note that the base64 encoding of encrypted data can actually make the encrypted
file larger than the cleartext one.

In-place encryption/decryption also works on binary files.

.. code::

	$ dd if=/dev/urandom of=/tmp/somerandom bs=1024
	count=512
	512+0 records in
	512+0 records out
	524288 bytes (524 kB) copied, 0.0466158 s, 11.2 MB/s

	$ sha512sum /tmp/somerandom
	9589bb20280e9d381f7a192000498c994e921b3cdb11d2ef5a986578dc2239a340b25ef30691bac72bdb14028270828dad7e8bd31e274af9828c40d216e60cbe /tmp/somerandom

	$ sops -e -i /tmp/somerandom
	please wait while a data encryption key is being generated and stored securely

	$ sops -d -i /tmp/somerandom

	$ sha512sum /tmp/somerandom
	9589bb20280e9d381f7a192000498c994e921b3cdb11d2ef5a986578dc2239a340b25ef30691bac72bdb14028270828dad7e8bd31e274af9828c40d216e60cbe /tmp/somerandom

Extract a sub-part of a document tree
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

`sops` can extract a specific part of a YAML or JSON document, by provided the
path in the `--extract` command line flag. This is useful to extract specific
values, like keys, without needing an extra parser.

.. code:: bash

	$ sops -d ~/git/svc/sops/example.yaml -t '["app2"]["key"]'
	-----BEGIN RSA PRIVATE KEY-----
	MIIBPAIBAAJBAPTMNIyHuZtpLYc7VsHQtwOkWYobkUblmHWRmbXzlAX6K8tMf3Wf
	ImcbNkqAKnELzFAPSBeEMhrBN0PyOC9lYlMCAwEAAQJBALXD4sjuBn1E7Y9aGiMz
	bJEBuZJ4wbhYxomVoQKfaCu+kH80uLFZKoSz85/ySauWE8LgZcMLIBoiXNhDKfQL
	vHECIQD6tCG9NMFWor69kgbX8vK5Y+QL+kRq+9HK6yZ9a+hsLQIhAPn4Ie6HGTjw
	fHSTXWZpGSan7NwTkIu4U5q2SlLjcZh/AiEA78NYRRBwGwAYNUqzutGBqyXKUl4u
	Erb0xAEyVV7e8J0CIQC8VBY8f8yg+Y7Kxbw4zDYGyb3KkXL10YorpeuZR4LuQQIg
	bKGPkMM4w5blyE1tqGN0T7sJwEx+EUOgacRNqM2ljVA=
	-----END RSA PRIVATE KEY-----

The tree path syntax uses regular python dictionary syntax, without the
variable name. Extract keys by naming them, and array elements by numbering
them.

.. code:: bash

	$ sops -d ~/git/svc/sops/example.yaml -t '["an_array"][1]'
	secretuser2

Using sops as a library in a python script
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can import sops as a module and use it in your python program.

.. code:: python

	import sops

	pathtype = sops.detect_filetype(path)
	tree = sops.load_file_into_tree(path, pathtype)
	sops_key, tree = sops.get_key(tree)
	tree = sops.walk_and_decrypt(tree, sops_key)
	sops.write_file(tree, path=path, filetype=pathtype)

Showing diffs in cleartext in git
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You most likely want to store encrypted files in a version controlled repository.
Sops can be used with git to decrypt files when showing diffs between versions.
This is very handy for reviewing changes or visualizing history.

To configure sops to decrypt files during diff, create a `.gitattributes` file
at the root of your repository that contains a filter and a command.

... code::

	*.yaml diff=sopsdiffer

Here we only care about YAML files. `sopsdiffer` is an arbitrary name that we map
to a sops command in the git configuration file of the repository.

.. code:: bash

	$ git config diff.sopsdiffer.textconv "sops -d"

	$ grep -A 1 sopsdiffer .git/config
	[diff "sopsdiffer"]
		textconv = "sops -d"

With this in place, calls to `git diff` will decrypt both previous and current
versions of the target file prior to displaying the diff. And it even works with
git client interfaces, because they call git diff under the hood!

Encrypting only parts of a file
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Note: this only works on YAML and JSON files, not on BINARY files.

By default, `sops` encrypts all the values of a YAML or JSON file and leaves the
keys in cleartext. In some instances, you may want to exclude some values from
being encrypted. This can be accomplished by adding the suffix **_unencrypted**
to any key of a file. When set, all values underneath the key that set the
**_unencrypted** prefix will be left in cleartext.

Note that, while in cleartext, unencrypted content is still added to the
checksum of the file, and thus cannot be modified outside of sops without
breaking the file integrity check.

The unencrypted suffix can be set to a different value using the
`--unencrypted-suffix` option.

Encryption Protocol
-------------------

When sops creates a file, it generates a random 256 bit data key and asks each
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

Upon save, sops browses the entire file as a key/value tree. Every time sops
encounters a leaf value (a value that does not have children), it encrypts the
value with AES256_GCM using the data key and a 256 bit random initialization
vector.

Each file uses a single data key to encrypt all values of a document, but each
value receives a unique initialization vector and has unique authentication data.

Additional data is used to guarantee the integrity of the encrypted data
and of the tree structure: when encrypting the tree, key names are concatenated
into a byte string that is used as AEAD additional data (aad) when encrypting
values. We expect that keys do not carry sensitive information, and
keeping them in cleartext allows for better diff and overall readability.

Any valid KMS or PGP master key can later decrypt the data key and access the
data.

Multiple master keys allow for sharing encrypted files without sharing master
keys, and provide a disaster recovery solution. The recommended way to use sops
is to have two KMS master keys in different regions and one PGP public key with
the private key stored offline. If, by any chance, both KMS master keys are
lost, you can always recover the encrypted data using the PGP private key.

Message Authentication Code
~~~~~~~~~~~~~~~~~~~~~~~~~~~

In addition to authenticating branches of the tree using keys as additional
data, sops computes a MAC on all the values to ensure that no value has been
added or removed fraudulently. The MAC is stored encrypted with AES_GCM and
the data key under tree->`sops`->`mac`.

Motivation
----------

Automating the distribution of secrets and credentials to components of an
infrastructure is a hard problem. We know how to encrypt secrets and share them
between humans, but extending that trust to systems is difficult. Particularly
when these systems follow devops principles and are created and destroyed
without human intervention. The issue boils down to establishing the initial
trust of a system that just joined the infrastructure, and providing it access
to the secrets it needs to configure itself.

The initial trust
~~~~~~~~~~~~~~~~~

In many infrastructures, even highly dynamic ones, the initial trust is
established by a human. An example is seen in Puppet by the way certificates are
issued: when a new system attempts to join a Puppetmaster, an administrator
must, by default, manually approve the issuance of the certificate the system
needs. This is cumbersome, and many puppetmasters are configured to auto-sign
new certificates to work around that issue. This is obviously not recommended
and far from ideal.

AWS provides a more flexible approach to trusting new systems. It uses a
powerful mechanism of roles and identities. In AWS, it is possible to verify
that a new system has been granted a specific role at creation, and it is
possible to map that role to specific resources. Instead of trusting new systems
directly, the administrator trusts the AWS permission model and its automation
infrastructure. As long as AWS keys are safe, and the AWS API is secure, we can
assume that trust is maintained and systems are who they say they are.

KMS, Trust and secrets distribution
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Using the AWS trust model, we can create fine grained access controls to
Amazon's Key Management Service (KMS). KMS is a service that encrypts and
decrypts data with AES_GCM, using keys that are never visible to users of the
service. Each KMS master key has a set of role-based access controls, and
individual roles are permitted to encrypt or decrypt using the master key. KMS
helps solve the problem of distributing keys, by shifting it into an access
control problem that can be solved using AWS's trust model.

Operational requirements
~~~~~~~~~~~~~~~~~~~~~~~~

When Mozilla's Services Operations team started revisiting the issue of
distributing secrets to EC2 instances, we set a goal to store these secrets
encrypted until the very last moment, when they need to be decrypted on target
systems. Not unlike many other organizations that operate sufficiently complex
automation, we found this to be a hard problem with a number of prerequisites:

1. Secrets must be stored in YAML files for easy integration into hiera

2. Secrets must be stored in GIT, and when a new CloudFormation stack is
   built, the current HEAD is pinned to the stack. (This allows secrets to
   be changed in GIT without impacting the current stack that may
   autoscale).

3. Entries must be encrypted separately. Encrypting entire files as blobs makes
   git conflict resolution almost impossible. Encrypting each entry
   separately is much easier to manage.

4. Secrets must always be encrypted on disk (admin laptop, upstream
   git repo, jenkins and S3) and only be decrypted on the target
   systems

SOPS can be used to encrypt YAML, JSON and BINARY files. In BINARY mode, the
content of the file is treated as a blob, the same way PGP would encrypt an
entire file. In YAML and JSON modes, however, the content of the file is
manipulated as a tree where keys are stored in cleartext, and values are
encrypted. hiera-eyaml does something similar, and over the years we learned
to appreciate its benefits, namely:

* diffs are meaningful. If a single value of a file is modified, only that
  value will show up in the diff. The diff is still limited to only showing
  encrypted data, but that information is already more granular that
  indicating that an entire file has changed.

* conflicts are easier to resolve. If multiple users are working on the
  same encrypted files, as long as they don't modify the same values,
  changes are easy to merge. This is an improvement over the PGP
  encryption approach where unsolvable conflicts often happen when
  multiple users work on the same file.

OpenPGP integration
~~~~~~~~~~~~~~~~~~~

OpenPGP gets a lot of bad press for being an outdated crypto protocol, and while
true, what really made us look for alternatives is the difficulty of managing and
distributing keys to systems. With KMS, we manage permissions to an API, not keys,
and that's a lot easier to do.

But PGP is not dead yet, and we still rely on it heavily as a backup solution:
all our files are encrypted with KMS and with one PGP public key, with its
private key stored securely for emergency decryption in the event that we lose
all our KMS master keys.

SOPS can be used without KMS entirely, the same way you would use an encrypted
PGP file: by referencing the pubkeys of each individual who has access to the file.
It can easily be done by providing sops with a comma-separated list of public keys
when creating a new file:

.. code:: bash

	$ sops --pgp "E60892BB9BD89A69F759A1A0A3D652173B763E8F,84050F1D61AF7C230A12217687DF65059EF093D3,85D77543B3D624B63CEA9E6DBC17301B491B3F21" mynewfile.yaml

Threat Model
------------

The security of the data stored using sops is as strong as the weakest
cryptographic mechanism. Values are encrypted using AES256_GCM which is the
strongest symetric encryption algorithm known today. Data keys are encrypted
in either KMS, which also uses AES256_GCM, or PGP which uses either RSA or
ECDSA keys.

Going from the most likely to the least likely, the threats are as follows:

Compromised AWS credentials grant access to KMS master key
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

An attacker with access to an AWS console can grant itself access to one of
the KMS master keys used to encrypt a sops data key. This threat should be
mitigated by protecting AWS accesses with strong controls, such as multi-factor
authentication, and also by performing regular audits of permissions granted
to AWS users.

Compromised PGP key
~~~~~~~~~~~~~~~~~~~

PGP keys are routinely mishandled, either because owners copy them from
machine to machine, or because the key is left forgotten on an unused machine
an attacker gains access to. When using PGP encryption, sops users should take
special care of PGP private keys, and store them on smart cards or offline
as often as possible.

Factorized RSA key
~~~~~~~~~~~~~~~~~~

sops doesn't apply any restriction on the size or type of PGP keys. A weak PGP
keys, for example 512 bits RSA, could be factorized by an attacker to gain
access to the private key and decrypt the data key. Users of sops should rely
on strong keys, such as 2048+ bits RSA keys, or 256+ bits ECDSA keys.

Weak AES cryptography
~~~~~~~~~~~~~~~~~~~~~

A vulnerability in AES256_GCM could potentially leak the data key or the KMS
master key used by a sops encrypted file. While no such vulnerability exists
today, we recommend that users keep their encrypted files reasonably private.

Backward compatibility
----------------------

`sops` will remain backward compatible on the major version, meaning that all
improvements brought to the 1.X branch (current) will maintain the file format
introduced in **1.0**.

License
-------
Mozilla Public License Version 2.0

Authors
-------
* Julien Vehent <jvehent@mozilla.com> (lead & maintainer)

* Daniel Thornton <dthornton@mozilla.com>
* Alexis Metaireau <alexis@mozilla.com>
* Rémy Hubscher <natim@mozilla.com>
* Todd Wolfson <todd@twolfson.com>
* Brian Hourigan <bhourigan@mozilla.com>

Credits
-------

`sops` is inspired by `hiera-eyaml <https://github.com/TomPoulton/hiera-eyaml>`_,
`credstash <https://github.com/LuminalOSS/credstash>`_ ,
`sneaker <https://github.com/codahale/sneaker>`_,
`password store <http://www.passwordstore.org/>`_ and too many years managing
PGP encrypted files by hand...
