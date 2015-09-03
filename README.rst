SOPS: Secrets OPerationS
========================
`sops` is a cli that encrypt values of yaml, json or text files using AWS KMS.

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

Input some cleartext yaml:

.. code:: yaml

    myapp1: t00m4nys3cr3tz
    app2:
        db:
            user: bob
            password: c4r1b0u
        # private key for secret operations in app2
        key: |
            -----BEGIN RSA PRIVATE KEY-----
            MIIBPAIBAAJBAPTMNIyHuZtpLYc7VsHQtwOkWYobkUblmHWRmbXzlAX6K8tMf3Wf
            Erb0xAEyVV7e8J0CIQC8VBY8f8yg+Y7Kxbw4zDYGyb3KkXL10YorpeuZR4LuQQIg
            bKGPkMM4w5blyE1tqGN0T7sJwEx+EUOgacRNqM2ljVA=
            -----END RSA PRIVATE KEY-----
    an_array:
        - secretuser1 # a super secret user
        - secretuser2
    sops:
        kms:
            enc: CiC6yCOtzsnFhkfdIslYZ0bAf//gYLYCmIu87B3sy/5yYxKnAQEBAQB4usg......
            enc_ts: 1439587921.752637
            arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e

After saving the file and exiting, it is automatically encrypted. Keys are
still in cleartext, but value are now unreadable.

.. code:: yaml

    myapp1: ENC[AES256_GCM,data=s4mlbkPqyk+GFDluAHY=,iv=7c9X8CwZyK5PsRRmUpzxL4CeQmp7+ry6mVemJtmpR7U=,aad=CFVNHUiz8xupOCMNYUlF4l+TcCjGaxayiknL9tQtolw=,tag=5ecBRedoXPJJ3uBjaj7J1w==]
    app2:
        db:
            user: ENC[AES256_GCM,data=CmkT,iv=xnUTxXU4g5lKEqetiZrM2s+m20idUUt9xGU6XitsIic=,aad=KidFJD6ioPXKqz+BYVYXtHk8Dd6e1yvhPx6kO5BOJTs=,tag=7WDZbBf2oqMuXi3YH4m2Ig==]
            password: ENC[AES256_GCM,data=zw2yh6Oz8Q==,iv=Apme9l8h+OwdwgbozsuXa1mVK+b821eoQNEBBSF6Ihs=,aad=SZFoaQDlNe0SkRaX65zB7E8SDyhkr9uVBI+3GWUBKsQ=,tag=We5dwW455S1M4ob1HzAu7Q==]
        # private key for secret operations in app2
        key: |-
            ENC[AES256_GCM,data=feo0o1qW8p4Nw2tN9/QAt0zoeGZHgolORWXH+7hk4Oc5nQcA/Ve3mYQ9TKSZAtzsYr+OEnEVUAg/RzvXy40F9dXsv2ugux+DLS1SIWddKRAdeL283vjnsDtydc3+AP+UuEuCyIVHqT8uKcqQnenzzu/yx07scIwcMQ6Vs2RnQ3WwrOrkBsbPQ4PpuPsrlck7EbcKkMnoIe09AMN3/J3A+mlmOGBxAio1ahFXpeBzzYeoRkjffojvigT2ULZy92Kx1afRSnWXNmUtMKqbJDIIvYulWHW5efAnulk/nHZ0Bhy+wxV0jqXAp7mKiIlGuydxRZ6DPon7jhABWV5d93EZdZJ+/33sUiOyQKIEukqae17C2Hrt0QoGg7OhG/O0oyTKiql0Nj6KC3bFbkdM7sSFbsIbv/of0P5Kb2zr3VYAJriJqUWMKzj3i7M6z9+wxrTVxuMQ4Lvzw3aHDjNOgJobkfjxxMYUvF5l2OWRFrdxtY9WxBYAcDzkJnBYtkPnUzlEc/8ieypefqOBlcphOvzl+EjM1I0N4OGG5ij5nNsHQ/MSoM3FJpjROQKclhz8ZN5CH41LUemP3AdddPpoUuwzHCxR8NskUhyHBlep0iZL9xGFLL7SwYEACKxk2BCwHMWeNmXfKo6co+wjCmn+un3FANE=,iv=NworRcR7VnLgW30c4W9OmVgBaY7tA1fd090JQpBM5ho=,aad=sbwFbTuEr9FbPd/ofR7BL9NORUpfmNd+X3Q+tJqmj8g=,tag=wc7RWWBArQrTMt3AAbSwZQ==]
    an_array:
    - ENC[AES256_GCM,data=L3Y0Bzn2M6yERcU=,iv=FslXY0z783MXhjCaz9ZZTqNaEwBWZkspNHAtHJaENH0=,aad=x0x9+PnDW81oLbYufq72RmaRZB29IPCALCL94KtmsvQ=,tag=qPyqJ3I9JM6wIJDOmgmJkQ==]
    - ENC[AES256_GCM,data=To5dwUDJi4Mh3hc=,iv=03vcf/AJaUKcHKEnGPq7ih8/xaKHewYiFkQcWOsh7So=,aad=nxUVG7rA+TjyK9BrzVtDGbCp7Iu7BCRLjYvZSnI5iCI=,tag=41ExX9KH+jRYvn51aaP6OA==]
    sops:
        kms:
            enc: CiC6yCOtzsnFhkfdIslYZ0bAf//gYLYCmIu87B3sy/5yYxKnAQEBAQB4usgjrc7JxYZH3SLJWGdGwH//4GC2ApiLvOwd7Mv+cmMAAAB+MHwGCSqGSIb3DQEHBqBvMG0CAQAwaAYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAwkRAZG5vQyIKvIKPwCARCAO9zQ43qeQ8loKu0HzXRnpqi6MK/+TpbO22sH0NkVXddXNTl7lfPjKc6gJynrEVdu6aCslUYIid+3FONY
            enc_ts: 1439587921.752637
            arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e

To decrypt, using flag `-d`.

.. code:: bash

	$ sops -d newfile.yaml
	myapp1: t00m4nys3cr3tz
	app2:
		db:
			user: bob
	[...]	

Set the env variable **SOPS_KMS_ARN** to your KMS ARN value to avoid
needing to set the `-k` flag every time you create a file.

.. code:: bash

	$ export SOPS_KMS_ARN="arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e"
	$ sops newfile.yaml

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

License
-------
Mozilla Public License Version 2.0

Authors
-------
* Julien Vehent

Credits
-------

`sops` is inspired by projects like `hiera-eyaml
<https://github.com/TomPoulton/hiera-eyaml>`_, `credstash
<https://github.com/LuminalOSS/credstash>`_ and `sneaker
<https://github.com/codahale/sneaker>`_. 
