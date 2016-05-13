import os
import codecs
from setuptools import setup, find_packages

here = os.path.abspath(os.path.dirname(__file__))
with codecs.open(os.path.join(here, 'README.rst'), encoding='utf-8') as f:
        README = f.read()

setup(
    name="sops",
    py_modules=['sops'],
    version="1.13",
    author="Julien Vehent",
    author_email="jvehent@mozilla.com",
    description="Secrets OPerationS (sops) is an editor of encrypted files",
    license="MPL",
    keywords="mozilla secret credential encryption aws kms",
    url="https://github.com/mozilla/sops",
    packages=find_packages(),
    zip_safe=True,
    long_description=README,
    install_requires=[
        'ruamel.yaml>=0.10.7',
        'boto3>=1.1.3',
        'cryptography>=0.9.3',
        'setuptools>=11.3'],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "License :: OSI Approved :: Mozilla Public License 2.0 (MPL 2.0)",
    ],
    entry_points={
        'console_scripts': [
            'sops = sops:main'
        ]
    }
)
