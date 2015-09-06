#!/usr/bin/env python
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

import os
from distutils.core import setup

def read(fname):
    return open(os.path.join(os.path.dirname(__file__), fname)).read()

setup(
    name            = "sops",
    packages        = ['sops'],
    version         = "0.2",
    author          = "Julien Vehent",
    author_email    = "jvehent@mozilla.com",
    description     = "Secrets OPerationS (sops) is an editor of encrypted files",
    license         = "MPL",
    keywords        = "mozilla secret credential encryption aws kms",
    url             = "https://github.com/mozilla-services/sops",
    long_description= read('README.rst'),
    install_requires= ['ruamel.yaml', 'json', 'boto3', 'cryptography'],
    classifiers     = [
        "Development Status :: 5 - Production/Stable",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "License :: OSI Approved :: Mozilla Public License 2.0 (MPL 2.0)",
    ],
)
