SOPS: Secrets OPerationS
========================
`sops` is a secrets management tool that encrypts YAML, JSON and TEXT files
using AWS KMS and/or PGP (via GnuPG).

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

.. code:: yaml

    myapp1: ENC[AES256_GCM,data:Tr7oc19nc6t1m9OrUeo=,iv:1vzlPZLfy6wa14/x17P8Ix8wEGDeY0v2dIboZmmwpww=,aad:NpobRzMzpDOkqijzONm8KglltzG+aBV7BJAxtm77veo=,tag:kaYqRgGGBhXhODSSmIZwyA==]
    app2:
        db:
            user: ENC[AES256_GCM,data:CwE4O1s=,iv:S0fozGAOxNma/pWDUuk1iEaYw0wlba0VOLHjPxIok2k=,aad:nEVizsMMyBXOxySnOHw/trTFBSW72nh+Q80YU7TPgIo=,tag:XaGsYaL9LCkLWJI0uxnTYw==]
            password: ENC[AES256_GCM,data:p673JCgHYw==,iv:EOOeivCp/Fd80xFdMYX0QeZn6orGTK8CeckmipjKqYY=,aad:UAhi/SHK0aCzptnFkFG4dW8Vv1ASg7TDHD6lui9mmKQ=,tag:QE6uuhRx+cGInwSVdmxXzA==]
        # private key for secret operations in app2
        key: |-
            ENC[AES256_GCM,data:Ea3zTFSOlg1PDZmBa1U2dtKl3pO4nTmaFswJx41fPfq3u8O2/Bq1UVfXn2SrO13obfr6xH4zuUceCDTvW2qvphlan5ir609EXt4dE2TEEcjVKhmAHf4LMwlZVAbvTJtlsnvo/aYJH95uctjsSX5h8pBlLaTGBGYwMrZuMyRU6vdcMWyha+piJckUc9sq7fevy1TSqIxf1Usbn/0NEklWm2VSNzQ2Urqtny6EXar+xU7NfYSRJ3mqmcJZ14oIeXPdpk962RwMEFWdYrbE7D59kWU2BgMjDxYJD5KXpWiw2YCrA/wsATxVCbZlwqC+TJFA5WAUZX756mFhV/t2Li3zQyDNUe6KkMXV9qwf/oV1j5sVRVFsKDYIBqhi3qWBVA+SO9RloQMjhru+IsdbQcS4LKq/1DrBENeZuJ0djUAxKLVfJzMGUf89ju3m9IEPovW8mfF0RbfAGRwFHMO9nEXCxrTLERf3owdR3u4j5/rNBpIvvy1z+2dy6sAx/eyNdS+cn5qO9BPAxsXpSwkaI96rlBagwH1Pfxus0x/D00j93OpE+M8MgQ/9LA68FlCFU4OAQlvw8f7MPoxnq+/+gFTS/qqjTR6EoUuX5NH2WY93YCC5TCbe4GOXyP0H05PbIWq55UMVLNcpAyac3gO4kL5O5U8=,iv:Dl61tsemKH0fdmNul/PmEEsRYFAh8GorR8GRupus/EM=,aad:Ft2aSYYukD1x8pMj1WvmodLjJV6waPy5FqdlImWyQKA=,tag:EPg4KpWqni/buCFjFL857A==]
    an_array:
    - ENC[AES256_GCM,data:v8dfh92oL8IcgjQ=,iv:HgNNPlQh9GNdE+YPvG4Ufpb2I0sIlEpCsOW3lJA1uBE=,aad:21GroP5gb9sCTxZIahN1NhMGqRPQZZksAr5Q7eCeHRc=,tag:gLsjVqot9+Pqck9LJC+bVA==]
    - ENC[AES256_GCM,data:X1LMy27AE9SI4h0=,iv:oA1kSg9esGxAvi3qhpcM6Ewrh+p0CFV5cgf6jSPpM08=,aad:CZ7FGJNko6367sd6PwbrIgN/V7Rly4TptbQ1gVsXT1Q=,tag:HerE4nTstX2QZhMn3CPZcw==]
    - ENC[AES256_GCM,data:KNkH9iI0bSyvcP3E+BRbqfcPUv3YBbCmtvbK1y+sHMI6Z1kXnkX4RoyYiZZXrM680Nh/p0TxNOdNsA==,iv:1h3KbThwTsRaVF+k+dnSwfocSEoyT00X279Dg1Wro60=,aad:foCwpM862VeAD2/7bHRJHAYISneTUJweoSRl2oAdsI4=,tag:tNuCjsNqIy5FVDRu39dQcw==]
    sops:
        kms:
        -   created_at: 1441570389.775376
            enc: CiC6yCOtzsnFhkfdIslYZ0bAf//gYLYCmIu87B3sy/5yYxKnAQEBAgB4usgjrc7JxYZH3SLJWGdGwH//4GC2ApiLvOwd7Mv+cmMAAAB+MHwGCSqGSIb3DQEHBqBvMG0CAQAwaAYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAxn6jfG4e44/phCddICARCAOzfGN/7WlU0MouQRXv22Pix46dSocMH1K7Xf47WqF1rCEcuN1aMVBj+IxwOgOVxVsr0Kze4lnMqPm1Hm
            arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
        -   created_at: 1441570391.925734
            enc: CiBdfsKZbRNf/Li8Tf2SjeSdP76DineB1sbPjV0TV+meTxKnAQEBAgB4XX7CmW0TX/y4vE39ko3knT++g4p3gdbGz41dE1fpnk8AAAB+MHwGCSqGSIb3DQEHBqBvMG0CAQAwaAYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAxGzsadorzSGbp73+ECARCAO0hc3cYxgNF2OU5TfTj8iyt/S6DTKDO+gwcHc3sy3ELQ/pUjSFJScYOQmqYpvsznhZ4YjHQWDdbRawNx
            arn: arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d
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
value with AES256_GCM using the data key, a 256 bits random initialization vector
and 256 bits of random additional data. While the same data key is used to
encrypt all values of a document, each value receives a unique initialization
vector and unique authentication data.

The result of AES256_GCM encryption is stored in the leaf of the tree using a
simple key/value format::

    ENC[AES256_GCM,
        data:CwE4O1s=,
        iv:S0fozGAOxNma/pWDUuk1iEaYw0wlba0VOLHjPxIok2k=,
        aad:nEVizsMMyBXOxySnOHw/trTFBSW72nh+Q80YU7TPgIo=,
        tag:XaGsYaL9LCkLWJI0uxnTYw==]

where:

* **data** is the encrypted value
* **iv** is the 256 bits initialization vector
* **aad** is the 256 bits additional data
* **tag** is the authentication tag

The encrypted file is written to disk with nested keys in cleartext and
encrypted values. We expect that keys do not carry sensitive information, and
keeping them in cleartext allows for better diff and overall readability.

Any valid KMS or PGP master key can later decrypt the data key and access the
data.

Multiple master keys allow for sharing encrypted files without sharing master
keys, and provide disaster recovery solution. The recommended way to use sops
is to have two KMS master keys in different region and one PGP public key with
the private key stored offline. If, by any chance, both KMS master keys are
lost, you can always recover the encrypted data using the PGP private key.

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

Credits
-------

`sops` is inspired by projects like `hiera-eyaml
<https://github.com/TomPoulton/hiera-eyaml>`_, `credstash
<https://github.com/LuminalOSS/credstash>`_ and `sneaker
<https://github.com/codahale/sneaker>`_. 
