SOPS: Secrets OPerationS
========================

**SOPS** is an editor of encrypted files that supports YAML, JSON, ENV, INI and BINARY
formats and encrypts with AWS KMS, GCP KMS, Azure Key Vault, HuaweiCloud KMS, age, and PGP.
(`demo <https://www.youtube.com/watch?v=YTEVyLXFiq0>`_)

.. image:: https://i.imgur.com/X0TM5NI.gif

------------

.. image:: https://pkg.go.dev/badge/github.com/getsops/sops/v3.svg
    :target: https://pkg.go.dev/github.com/getsops/sops/v3

Documentation
-------------

You can find the SOPS documentation on `getsops.io <https://getsops.io/>`_ under `"Docs" <https://getsops.io/docs/>`_.

Security
--------

Please report any security issues privately using `GitHub's advisory form <https://github.com/getsops/sops/security/advisories>`_.

License
-------
Mozilla Public License Version 2.0

Authors
-------

SOPS was initially launched as a project at Mozilla in 2015 and has been
graciously donated to the CNCF as a Sandbox project in 2023, now under the
stewardship of a `new group of maintainers <https://github.com/getsops/community/blob/main/MAINTAINERS.md>`_.

The original authors of the project were:

* Adrian Utrilla @autrilla
* Julien Vehent @jvehent

Furthermore, the project has been carried for a long time by AJ Bahnken @ajvb,
and had not been possible without the contributions of numerous `contributors <https://github.com/getsops/sops/graphs/contributors>`_.

Credits
-------

SOPS was inspired by `hiera-eyaml <https://github.com/TomPoulton/hiera-eyaml>`_,
`credstash <https://github.com/LuminalOSS/credstash>`_,
`sneaker <https://github.com/codahale/sneaker>`_,
`password store <http://www.passwordstore.org/>`_ and too many years managing
PGP encrypted files by hand...

-----

.. image:: docs/images/cncf-color-bg.svg
   :width: 400
   :alt: CNCF Sandbox Project

**We are a** `Cloud Native Computing Foundation <https://cncf.io>`_ **sandbox project.**
