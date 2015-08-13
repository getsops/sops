SOPS: Secrets OPerationS
========================
`sops` is a cli that encrypt values of yaml, json or text files using AWS KMS.

Requirements
------------
* `boto3 <https://pypi.python.org/pypi/boto3/1.1.1>`_
* `ruamel.yaml <https://pypi.python.org/pypi/ruamel.yaml>`_; requires
  libyaml-devel and python-devel prior to `pip install`-ing it.

.. code::

	sudo yum install libyaml-devel python-devel
	sudo pip install ruamel.yaml

* `cryptography <https://pypi.python.org/pypi/cryptography>`_; requires
  libffi-devel prior to `pip install`-ing it.

.. code::

	sudo yum install libffi-devel
	sudo pip install cryptography

Usage
-----

Editing
~~~~~~~

`sops` encrypted file contain the necessary KMS information to decrypt their
content. All a user of `sops` need is valid AWS credentials and the necessary
permissions on KMS keys.

Given that, the only command a `sops` user need is:

.. code:: bash

	$ sops <file>

`<file>` will be opened, decrypted, passed to a text editor (vim by default),
encrypted if modified, and save back to its original location. All of these
steps, apart from the actual editing, are transparent to the user.

Creating
~~~~~~~~

In order to create a file, the KMS ARN must be provided to `sops`, either on the
command line in the `-k` flag, or in the environment variable **SOPS_KMS_ARN**.

`sops` automatically create a file if the given path doesn't exist (it will not
create folders, however).

.. code:: bash

	$ sops newfile.yaml -k arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
	newfile.yaml doesn't exist, creating it.
	new data key generated from kms: CiC6yCOtzsnFhkfdIs...
	file written to newfile.yaml

	$ ./sops -d newfile.yaml 2>/dev/null
	mysecretkey: value12345abcdef
	sops:
	  kms:
		enc: CiC6yCOtzsnFhkvfd...
		enc_ts: 1439502977.62264
		arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e

To create and encrypt a file without specifying the KMS ARN in `-k`:

.. code:: bash

	$ export SOPS_KMS_ARN="arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e"
	$ sops newfile.yaml

Reading
~~~~~~~

To read an encrypted file without opening an editor, use `-d` flag. The
content of the file is sent to **stdout**, and accompanying messages are
sent to **stderr** (can be ignored with `2>/dev/null`).

.. code:: bash

	$ sops -d newfile.yaml

License
-------
Mozilla Public License Version 2.0

Authors
-------
* Julien Vehent
