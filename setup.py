import os
from setuptools import setup, find_packages


def read(fname):
    return open(os.path.join(os.path.dirname(__file__), fname)).read()


setup(
    name="sops",
    py_modules=['sops'],
    version="0.6.2",
    author="Julien Vehent",
    author_email="jvehent@mozilla.com",
    description="Secrets OPerationS (sops) is an editor of encrypted files",
    license="MPL",
    keywords="mozilla secret credential encryption aws kms",
    url="https://github.com/mozilla/sops",
    packages=find_packages(),
    zip_safe=True,
    long_description=read('README.rst'),
    install_requires=[
        'ruamel.yaml>=0.10.7', 'boto3>=1.1.3', 'cryptography>=0.9.3'],
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
