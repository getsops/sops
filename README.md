<!-- THIS FILE HAS BEEN AUTOMATICALLY CONVERTED FROM README.rst.
     DO NOT MODIFY THIS FILE, MODIFY README.rst INSTEAD! -->

# SOPS\: Secrets OPerationS

<strong>SOPS</strong> is an editor of encrypted files that supports YAML\, JSON\, ENV\, INI and BINARY
formats and encrypts with AWS KMS\, GCP KMS\, Azure Key Vault\, age\, and PGP\.
\([demo](https\://www\.youtube\.com/watch\?v\=YTEVyLXFiq0)\)

![](https\://i\.imgur\.com/X0TM5NI\.gif)

---

[
![](https\://pkg\.go\.dev/badge/github\.com/getsops/sops/v3\.svg)
](https\://pkg\.go\.dev/github\.com/getsops/sops/v3)

<a id="download"></a>
## 1   Download

<a id="stable-release"></a>
### 1\.1   Stable release

Binaries and packages of the latest stable release are available at [https\://github\.com/getsops/sops/releases](https\://github\.com/getsops/sops/releases)\.

<a id="development-branch"></a>
### 1\.2   Development branch

For the adventurous\, unstable features are available in the <em class="title-reference">main</em> branch\, which you can install from source\:

```bash
$ mkdir -p $GOPATH/src/github.com/getsops/sops/
$ git clone https://github.com/getsops/sops.git $GOPATH/src/github.com/getsops/sops/
$ cd $GOPATH/src/github.com/getsops/sops/
$ make install
```

\(requires Go \>\= 1\.19\)

If you don\'t have Go installed\, set it up with\:

```bash
$ {apt,yum,brew} install golang
$ echo 'export GOPATH=~/go' >> ~/.bashrc
$ source ~/.bashrc
$ mkdir $GOPATH
```

Or whatever variation of the above fits your system and shell\.

To use <strong>SOPS</strong> as a library\, take a look at the [decrypt package](https\://pkg\.go\.dev/github\.com/getsops/sops/v3/decrypt)\.

#### Table of Contents

- [1   Download](\#download)

  - [1\.1   Stable release](\#stable\-release)
  - [1\.2   Development branch](\#development\-branch)
- [2   Usage](\#usage)

  - [2\.1   Test with the dev PGP key](\#test\-with\-the\-dev\-pgp\-key)
  - [2\.2   Encrypting using age](\#encrypting\-using\-age)
  - [2\.3   Encrypting using GCP KMS](\#encrypting\-using\-gcp\-kms)
  - [2\.4   Encrypting using Azure Key Vault](\#encrypting\-using\-azure\-key\-vault)
  - [2\.5   Encrypting and decrypting from other programs](\#encrypting\-and\-decrypting\-from\-other\-programs)
  - [2\.6   Encrypting using Hashicorp Vault](\#encrypting\-using\-hashicorp\-vault)
  - [2\.7   Adding and removing keys](\#adding\-and\-removing\-keys)

    - [2\.7\.1   <code>updatekeys</code> command](\#updatekeys\-command)
    - [2\.7\.2   <code>rotate</code> command](\#rotate\-command)
    - [2\.7\.3   Direct Editing](\#direct\-editing)
  - [2\.8   KMS AWS Profiles](\#kms\-aws\-profiles)
  - [2\.9   Assuming roles and using KMS in various AWS accounts](\#assuming\-roles\-and\-using\-kms\-in\-various\-aws\-accounts)
  - [2\.10   AWS KMS Encryption Context](\#aws\-kms\-encryption\-context)
  - [2\.11   Key Rotation](\#key\-rotation)
  - [2\.12   Using \.sops\.yaml conf to select KMS\, PGP and age for new files](\#using\-sops\-yaml\-conf\-to\-select\-kms\-pgp\-and\-age\-for\-new\-files)
  - [2\.13   Specify a different GPG executable](\#specify\-a\-different\-gpg\-executable)
  - [2\.14   Specify a different GPG key server](\#specify\-a\-different\-gpg\-key\-server)
  - [2\.15   Key groups](\#key\-groups)
  - [2\.16   Key service](\#key\-service)
  - [2\.17   Auditing](\#auditing)
  - [2\.18   Saving Output to a File](\#saving\-output\-to\-a\-file)
  - [2\.19   Passing Secrets to Other Processes](\#passing\-secrets\-to\-other\-processes)
  - [2\.20   Using the publish command](\#using\-the\-publish\-command)

    - [2\.20\.1   Publishing to Vault](\#publishing\-to\-vault)
- [3   Important information on types](\#important\-information\-on\-types)

  - [3\.1   YAML\, JSON\, ENV and INI type extensions](\#yaml\-json\-env\-and\-ini\-type\-extensions)
  - [3\.2   JSON and JSON\_binary indentation](\#json\-and\-json\-binary\-indentation)
  - [3\.3   YAML indentation](\#yaml\-indentation)
  - [3\.4   YAML anchors](\#yaml\-anchors)
  - [3\.5   YAML Streams](\#yaml\-streams)
  - [3\.6   Top\-level arrays](\#top\-level\-arrays)
- [4   Examples](\#examples)

  - [4\.1   Creating a new file](\#creating\-a\-new\-file)
  - [4\.2   Encrypting an existing file](\#encrypting\-an\-existing\-file)
  - [4\.3   Encrypt or decrypt a file in place](\#encrypt\-or\-decrypt\-a\-file\-in\-place)
  - [4\.4   Encrypting binary files](\#encrypting\-binary\-files)
  - [4\.5   Extract a sub\-part of a document tree](\#extract\-a\-sub\-part\-of\-a\-document\-tree)
  - [4\.6   Set a sub\-part in a document tree](\#set\-a\-sub\-part\-in\-a\-document\-tree)
  - [4\.7   Showing diffs in cleartext in git](\#showing\-diffs\-in\-cleartext\-in\-git)
  - [4\.8   Encrypting only parts of a file](\#encrypting\-only\-parts\-of\-a\-file)
- [5   Encryption Protocol](\#encryption\-protocol)

  - [5\.1   Message Authentication Code](\#message\-authentication\-code)
- [6   Motivation](\#motivation)

  - [6\.1   The initial trust](\#the\-initial\-trust)
  - [6\.2   KMS\, Trust and secrets distribution](\#kms\-trust\-and\-secrets\-distribution)
  - [6\.3   Operational requirements](\#operational\-requirements)
  - [6\.4   OpenPGP integration](\#openpgp\-integration)
- [7   Threat Model](\#threat\-model)

  - [7\.1   Compromised AWS credentials grant access to KMS master key](\#compromised\-aws\-credentials\-grant\-access\-to\-kms\-master\-key)
  - [7\.2   Compromised PGP key](\#compromised\-pgp\-key)
  - [7\.3   Factorized RSA key](\#factorized\-rsa\-key)
  - [7\.4   Weak AES cryptography](\#weak\-aes\-cryptography)
- [8   Backward compatibility](\#backward\-compatibility)
- [9   Security](\#security)
- [10   License](\#license)
- [11   Authors](\#authors)
- [12   Credits](\#credits)

<a id="usage"></a>
## 2   Usage

For a quick presentation of SOPS\, check out this Youtube tutorial\:

[
![](https\://img\.youtube\.com/vi/V2PRhxphH2w/0\.jpg)
](https\://www\.youtube\.com/watch\?v\=V2PRhxphH2w)

If you\'re using AWS KMS\, create one or multiple master keys in the IAM console
and export them\, comma separated\, in the <strong>SOPS\_KMS\_ARN</strong> env variable\. It is
recommended to use at least two master keys in different regions\.

```bash
export SOPS_KMS_ARN="arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e,arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
```

SOPS uses [aws\-sdk\-go\-v2](https\://github\.com/aws/aws\-sdk\-go\-v2) to communicate with AWS KMS\. It will automatically
read the credentials from the <code>\~/\.aws/credentials</code> file which can be created with the <code>aws configure</code> command\.

An example of the <code>\~/\.aws/credentials</code> file is shown below\:

```sh
$ cat ~/.aws/credentials
[default]
aws_access_key_id = AKI.....
aws_secret_access_key = mw......
```

In addition to the <code>\~/\.aws/credentials</code> file\, you can also use the <code>AWS\_ACCESS\_KEY\_ID</code> and <code>AWS\_SECRET\_ACCESS\_KEY</code>
environment variables to specify your credentials\:

```bash
export AWS_ACCESS_KEY_ID="AKI......"
export AWS_SECRET_ACCESS_KEY="mw......"
```

For more information and additional environment variables\, see
[specifying credentials](https\://aws\.github\.io/aws\-sdk\-go\-v2/docs/configuring\-sdk/\#specifying\-credentials)\.

If you want to use PGP\, export the fingerprints of the public keys\, comma
separated\, in the <strong>SOPS\_PGP\_FP</strong> env variable\.

```bash
export SOPS_PGP_FP="85D77543B3D624B63CEA9E6DBC17301B491B3F21,E60892BB9BD89A69F759A1A0A3D652173B763E8F"
```

Note\: you can use both PGP and KMS simultaneously\.

Then simply call <code>sops edit</code> with a file path as argument\. It will handle the
encryption/decryption transparently and open the cleartext file in an editor

```sh
$ sops edit mynewtestfile.yaml
mynewtestfile.yaml doesn't exist, creating it.
please wait while an encryption key is being generated and stored in a secure fashion
file written to mynewtestfile.yaml
```

Editing will happen in whatever <code>\$EDITOR</code> is set to\, or\, if it\'s not set\, in vim\.
Keep in mind that SOPS will wait for the editor to exit\, and then try to reencrypt
the file\. Some GUI editors \(atom\, sublime\) spawn a child process and then exit
immediately\. They usually have an option to wait for the main editor window to be
closed before exiting\. See [\#127](https\://github\.com/getsops/sops/issues/127) for
more information\.

The resulting encrypted file looks like this\:

```yaml
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
        - created_at: 1441570389.775376
          enc: CiC....Pm1Hm
          arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
        - created_at: 1441570391.925734
          enc: Ci...awNx
          arn: arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d
    pgp:
        - fp: 85D77543B3D624B63CEA9E6DBC17301B491B3F21
          created_at: 1441570391.930042
          enc: |
              -----BEGIN PGP MESSAGE-----
              hQIMA0t4uZHfl9qgAQ//UvGAwGePyHuf2/zayWcloGaDs0MzI+zw6CmXvMRNPUsA
              ...=oJgS
              -----END PGP MESSAGE-----
```

A copy of the encryption/decryption key is stored securely in each KMS and PGP
block\. As long as one of the KMS or PGP method is still usable\, you will be able
to access your data\.

To decrypt a file in a <code>cat</code> fashion\, use the <code>\-d</code> flag\:

```sh
$ sops decrypt mynewtestfile.yaml
```

SOPS encrypted files contain the necessary information to decrypt their content\.
All a user of SOPS needs is valid AWS credentials and the necessary
permissions on KMS keys\.

Given that\, the only command a SOPS user needs is\:

```sh
$ sops edit <file>
```

<em class="title-reference">\<file\></em> will be opened\, decrypted\, passed to a text editor \(vim by default\)\,
encrypted if modified\, and saved back to its original location\. All of these
steps\, apart from the actual editing\, are transparent to the user\.

The order in which available decryption methods are tried can be specified with
<code>\-\-decryption\-order</code> option or <strong>SOPS\_DECRYPTION\_ORDER</strong> environment variable
as a comma separated list\. The default order is <code>age\,pgp</code>\. Offline methods are
tried first and then the remaining ones\.

<a id="test-with-the-dev-pgp-key"></a>
### 2\.1   Test with the dev PGP key

If you want to test <strong>SOPS</strong> without having to do a bunch of setup\, you can use
the example files and pgp key provided with the repository\:

```
$ git clone https://github.com/getsops/sops.git
$ cd sops
$ gpg --import pgp/sops_functional_tests_key.asc
$ sops edit example.yaml
```

This last step will decrypt <code>example\.yaml</code> using the test private key\.

<a id="encrypting-using-age"></a>
### 2\.2   Encrypting using age

[age](https\://age\-encryption\.org/) is a simple\, modern\, and secure tool for
encrypting files\. It\'s recommended to use age over PGP\, if possible\.

You can encrypt a file for one or more age recipients \(comma separated\) using
the <code>\-\-age</code> option or the <strong>SOPS\_AGE\_RECIPIENTS</strong> environment variable\:

```sh
$ sops encrypt --age age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw test.yaml > test.enc.yaml
```

When decrypting a file with the corresponding identity\, SOPS will look for a
text file name <code>keys\.txt</code> located in a <code>sops</code> subdirectory of your user
configuration directory\. On Linux\, this would be <code>\$XDG\_CONFIG\_HOME/sops/age/keys\.txt</code>\.
On macOS\, this would be <code>\$HOME/Library/Application Support/sops/age/keys\.txt</code>\. On
Windows\, this would be <code>\%AppData\%\\sops\\age\\keys\.txt</code>\. You can specify the location
of this file manually by setting the environment variable <strong>SOPS\_AGE\_KEY\_FILE</strong>\.
Alternatively\, you can provide the key\(s\) directly by setting the <strong>SOPS\_AGE\_KEY</strong>
environment variable\.

The contents of this key file should be a list of age X25519 identities\, one
per line\. Lines beginning with <code>\#</code> are considered comments and ignored\. Each
identity will be tried in sequence until one is able to decrypt the data\.

Encrypting with SSH keys via age is not yet supported by SOPS\.

<a id="encrypting-using-gcp-kms"></a>
### 2\.3   Encrypting using GCP KMS

GCP KMS uses [Application Default Credentials](https\://developers\.google\.com/identity/protocols/application\-default\-credentials)\.
If you already logged in using

```sh
$ gcloud auth login
```

you can enable application default credentials using the sdk\:

```sh
$ gcloud auth application-default login
```

Encrypting/decrypting with GCP KMS requires a KMS ResourceID\. You can use the
cloud console the get the ResourceID or you can create one using the gcloud
sdk\:

```sh
$ gcloud kms keyrings create sops --location global
$ gcloud kms keys create sops-key --location global --keyring sops --purpose encryption
$ gcloud kms keys list --location global --keyring sops

# you should see
NAME                                                                   PURPOSE          PRIMARY_STATE
projects/my-project/locations/global/keyRings/sops/cryptoKeys/sops-key ENCRYPT_DECRYPT  ENABLED
```

Now you can encrypt a file using\:

```
$ sops encrypt --gcp-kms projects/my-project/locations/global/keyRings/sops/cryptoKeys/sops-key test.yaml > test.enc.yaml
```

And decrypt it using\:

```
$ sops decrypt test.enc.yaml
```

<a id="encrypting-using-azure-key-vault"></a>
### 2\.4   Encrypting using Azure Key Vault

The Azure Key Vault integration uses the
[default credential chain](https\://pkg\.go\.dev/github\.com/Azure/azure\-sdk\-for\-go/sdk/azidentity\#DefaultAzureCredential)
which tries several authentication methods\, in this order\:

1. [Environment credentials](https\://pkg\.go\.dev/github\.com/Azure/azure\-sdk\-for\-go/sdk/azidentity\#EnvironmentCredential)

   1. Service Principal with Client Secret
   1. Service Principal with Certificate
   1. User with username and password
   1. Configuration for multi\-tenant applications
1. [Workload Identity credentials](https\://pkg\.go\.dev/github\.com/Azure/azure\-sdk\-for\-go/sdk/azidentity\#WorkloadIdentityCredential)
1. [Managed Identity credentials](https\://pkg\.go\.dev/github\.com/Azure/azure\-sdk\-for\-go/sdk/azidentity\#ManagedIdentityCredential)
1. [Azure CLI credentials](https\://pkg\.go\.dev/github\.com/Azure/azure\-sdk\-for\-go/sdk/azidentity\#AzureCLICredential)

For example\, you can use a Service Principal with the following environment variables\:

```bash
AZURE_TENANT_ID
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
```

You can create a Service Principal using the CLI like this\:

```sh
$ az ad sp create-for-rbac -n my-keyvault-sp

{
    "appId": "<some-uuid>",
    "displayName": "my-keyvault-sp",
    "name": "http://my-keyvault-sp",
    "password": "<random-string>",
    "tenant": "<tenant-uuid>"
}
```

The <em class="title-reference">appId</em> is the client ID\, and the <em class="title-reference">password</em> is the client secret\.

Encrypting/decrypting with Azure Key Vault requires the resource identifier for
a key\. This has the following form\:

```
https://${VAULT_URL}/keys/${KEY_NAME}/${KEY_VERSION}
```

To create a Key Vault and assign your service principal permissions on it
from the commandline\:

```sh
# Create a resource group if you do not have one:
$ az group create --name sops-rg --location westeurope
# Key Vault names are globally unique, so generate one:
$ keyvault_name=sops-$(uuidgen | tr -d - | head -c 16)
# Create a Vault, a key, and give the service principal access:
$ az keyvault create --name $keyvault_name --resource-group sops-rg --location westeurope
$ az keyvault key create --name sops-key --vault-name $keyvault_name --protection software --ops encrypt decrypt
$ az keyvault set-policy --name $keyvault_name --resource-group sops-rg --spn $AZURE_CLIENT_ID \
    --key-permissions encrypt decrypt
# Read the key id:
$ az keyvault key show --name sops-key --vault-name $keyvault_name --query key.kid

https://sops.vault.azure.net/keys/sops-key/some-string
```

Now you can encrypt a file using\:

```
$ sops encrypt --azure-kv https://sops.vault.azure.net/keys/sops-key/some-string test.yaml > test.enc.yaml
```

And decrypt it using\:

```
$ sops decrypt test.enc.yaml
```

<a id="encrypting-and-decrypting-from-other-programs"></a>
### 2\.5   Encrypting and decrypting from other programs

When using <code>sops</code> in scripts or from other programs\, there are often situations where you do not want to write
encrypted or decrypted data to disk\. The best way to avoid this is to pass data to SOPS via stdin\, and to let
SOPS write data to stdout\. By default\, the encrypt and decrypt operations write data to stdout already\. To pass
data via stdin\, you need to pass <code>/dev/stdin</code> as the input filename\. Please note that this only works on
Unix\-like operating systems such as macOS and Linux\. On Windows\, you have to use named pipes\.

To decrypt data\, you can simply do\:

```sh
$ cat encrypted-data | sops decrypt /dev/stdin > decrypted-data
```

To control the input and output format\, pass <code>\-\-input\-type</code> and <code>\-\-output\-type</code> as appropriate\. By default\,
<code>sops</code> determines the input and output format from the provided filename\, which is <code>/dev/stdin</code> here\, and
thus will use the binary store which expects JSON input and outputs binary data on decryption\.

For example\, to decrypt YAML data and obtain the decrypted result as YAML\, use\:

```sh
$ cat encrypted-data | sops decrypt --input-type yaml --output-type yaml /dev/stdin > decrypted-data
```

To encrypt\, it is important to note that SOPS also uses the filename to look up the correct creation rule from
<code>\.sops\.yaml</code>\. Likely <code>/dev/stdin</code> will not match a creation rule\, or only match the fallback rule without
<code>path\_regex</code>\, which is usually not what you want\. For that\, <code>sops</code> provides the <code>\-\-filename\-override</code>
parameter which allows you to tell SOPS which filename to use to match creation rules\:

```sh
$ echo 'foo: bar' | sops encrypt --filename-override path/filename.sops.yaml /dev/stdin > encrypted-data
```

SOPS will find a matching creation rule for <code>path/filename\.sops\.yaml</code> in <code>\.sops\.yaml</code> and use that one to
encrypt the data from stdin\. This filename will also be used to determine the input and output store\. As always\,
the input store type can be adjusted by passing <code>\-\-input\-type</code>\, and the output store type by passing
<code>\-\-output\-type</code>\:

```sh
$ echo foo=bar | sops encrypt --filename-override path/filename.sops.yaml --input-type dotenv /dev/stdin > encrypted-data
```

<a id="encrypting-using-hashicorp-vault"></a>
### 2\.6   Encrypting using Hashicorp Vault

We assume you have an instance \(or more\) of Vault running and you have privileged access to it\. For instructions on how to deploy a secure instance of Vault\, refer to Hashicorp\'s official documentation\.

To easily deploy Vault locally\: \(DO NOT DO THIS FOR PRODUCTION\!\!\!\)

```sh
$ docker run -d -p8200:8200 vault:1.2.0 server -dev -dev-root-token-id=toor
```
```sh
$ # Substitute this with the address Vault is running on
$ export VAULT_ADDR=http://127.0.0.1:8200

$ # this may not be necessary in case you previously used `vault login` for production use
$ export VAULT_TOKEN=toor

$ # to check if Vault started and is configured correctly
$ vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.2.0
Cluster Name    vault-cluster-618cc902
Cluster ID      e532e461-e8f0-1352-8a41-fc7c11096908
HA Enabled      false

$ # It is required to enable a transit engine if not already done (It is suggested to create a transit engine specifically for SOPS, in which it is possible to have multiple keys with various permission levels)
$ vault secrets enable -path=sops transit
Success! Enabled the transit secrets engine at: sops/

$ # Then create one or more keys
$ vault write sops/keys/firstkey type=rsa-4096
Success! Data written to: sops/keys/firstkey

$ vault write sops/keys/secondkey type=rsa-2048
Success! Data written to: sops/keys/secondkey

$ vault write sops/keys/thirdkey type=chacha20-poly1305
Success! Data written to: sops/keys/thirdkey

$ sops encrypt --hc-vault-transit $VAULT_ADDR/v1/sops/keys/firstkey vault_example.yml

$ cat <<EOF > .sops.yaml
creation_rules:
    - path_regex: \.dev\.yaml$
      hc_vault_transit_uri: "$VAULT_ADDR/v1/sops/keys/secondkey"
    - path_regex: \.prod\.yaml$
      hc_vault_transit_uri: "$VAULT_ADDR/v1/sops/keys/thirdkey"
EOF

$ sops encrypt --verbose prod/raw.yaml > prod/encrypted.yaml
```

<a id="adding-and-removing-keys"></a>
### 2\.7   Adding and removing keys

When creating new files\, <code>sops</code> uses the PGP\, KMS and GCP KMS defined in the
command line arguments <code>\-\-kms</code>\, <code>\-\-pgp</code>\, <code>\-\-gcp\-kms</code> or <code>\-\-azure\-kv</code>\, or from
the environment variables <code>SOPS\_KMS\_ARN</code>\, <code>SOPS\_PGP\_FP</code>\, <code>SOPS\_GCP\_KMS\_IDS</code>\,
<code>SOPS\_AZURE\_KEYVAULT\_URLS</code>\. That information is stored in the file under the
<code>sops</code> section\, such that decrypting files does not require providing those
parameters again\.

Master PGP and KMS keys can be added and removed from a <code>sops</code> file in one of
three ways\:

1. By using a <code>\.sops\.yaml</code> file and the <code>updatekeys</code> command\.
1. By using command line flags\.
1. By editing the file directly\.

The SOPS team recommends the <code>updatekeys</code> approach\.

<a id="updatekeys-command"></a>
#### 2\.7\.1   updatekeys command

The <code>updatekeys</code> command uses the [\.sops\.yaml](\#using\-sops\-yaml\-conf\-to\-select\-kms\-pgp\-for\-new\-files)
configuration file to update \(add or remove\) the corresponding secrets in the
encrypted file\. Note that the example below uses the
[Block Scalar yaml construct](https\://yaml\-multiline\.info/) to build a space
separated list\.

```yaml
creation_rules:
    - pgp: >-
        85D77543B3D624B63CEA9E6DBC17301B491B3F21,
        FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4
```
```sh
$ sops updatekeys test.enc.yaml
```

SOPS will prompt you with the changes to be made\. This interactivity can be
disabled by supplying the <code>\-y</code> flag\.

<a id="rotate-command"></a>
#### 2\.7\.2   rotate command

The <code>rotate</code> command generates a new data encryption key and reencrypt all values
with the new key\. At te same time\, the command line flag <code>\-\-add\-kms</code>\, <code>\-\-add\-pgp</code>\,
<code>\-\-add\-gcp\-kms</code>\, <code>\-\-add\-azure\-kv</code>\, <code>\-\-rm\-kms</code>\, <code>\-\-rm\-pgp</code>\, <code>\-\-rm\-gcp\-kms</code>
and <code>\-\-rm\-azure\-kv</code> can be used to add and remove keys from a file\. These flags use
the comma separated syntax as the <code>\-\-kms</code>\, <code>\-\-pgp</code>\, <code>\-\-gcp\-kms</code> and <code>\-\-azure\-kv</code>
arguments when creating new files\.

Use <code>updatekeys</code> if you want to add a key without rotating the data key\.

```sh
# add a new pgp key to the file and rotate the data key
$ sops rotate -i --add-pgp 85D77543B3D624B63CEA9E6DBC17301B491B3F21 example.yaml

# remove a pgp key from the file and rotate the data key
$ sops rotate -i --rm-pgp 85D77543B3D624B63CEA9E6DBC17301B491B3F21 example.yaml
```

<a id="direct-editing"></a>
#### 2\.7\.3   Direct Editing

Alternatively\, invoking <code>sops edit</code> with the flag <strong>\-s</strong> will display the master keys
while editing\. This method can be used to add or remove <code>kms</code> or <code>pgp</code> keys under the
<code>sops</code> section\.

For example\, to add a KMS master key to a file\, add the following entry while
editing\:

```yaml
sops:
    kms:
        - arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
```

And\, similarly\, to add a PGP master key\, we add its fingerprint\:

```yaml
sops:
    pgp:
        - fp: 85D77543B3D624B63CEA9E6DBC17301B491B3F21
```

When the file is saved\, SOPS will update its metadata and encrypt the data key
with the freshly added master keys\. The removed entries are simply deleted from
the file\.

When removing keys\, it is recommended to rotate the data key using <code>\-r</code>\,
otherwise\, owners of the removed key may have add access to the data key in the
past\.

<a id="kms-aws-profiles"></a>
### 2\.8   KMS AWS Profiles

If you want to use a specific profile\, you can do so with <em class="title-reference">aws\_profile</em>\:

```yaml
sops:
    kms:
        - arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
          aws_profile: foo
```

If no AWS profile is set\, default credentials will be used\.

Similarly the <em class="title-reference">\-\-aws\-profile</em> flag can be set with the command line with any of the KMS commands\.

<a id="assuming-roles-and-using-kms-in-various-aws-accounts"></a>
### 2\.9   Assuming roles and using KMS in various AWS accounts

SOPS has the ability to use KMS in multiple AWS accounts by assuming roles in
each account\. Being able to assume roles is a nice feature of AWS that allows
administrators to establish trust relationships between accounts\, typically from
the most secure account to the least secure one\. In our use\-case\, we use roles
to indicate that a user of the Master AWS account is allowed to make use of KMS
master keys in development and staging AWS accounts\. Using roles\, a single file
can be encrypted with KMS keys in multiple accounts\, thus increasing reliability
and ease of use\.

You can use keys in various accounts by tying each KMS master key to a role that
the user is allowed to assume in each account\. The [IAM roles](http\://docs\.aws\.amazon\.com/IAM/latest/UserGuide/id\_roles\_use\.html)
documentation has full details on how this needs to be configured on AWS\'s side\.

From the point of view of SOPS\, you only need to specify the role a KMS key
must assume alongside its ARN\, as follows\:

```yaml
sops:
    kms:
        - arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
          role: arn:aws:iam::927034868273:role/sops-dev-xyz
```

The role must have permission to call Encrypt and Decrypt using KMS\. An example
policy is shown below\.

```json
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
```

You can specify a role in the <code>\-\-kms</code> flag and <code>SOPS\_KMS\_ARN</code> variable by
appending it to the ARN of the master key\, separated by a <strong>\+</strong> sign\:

```
<KMS ARN>+<ROLE ARN>
arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500+arn:aws:iam::927034868273:role/sops-dev-xyz
```

<a id="aws-kms-encryption-context"></a>
### 2\.10   AWS KMS Encryption Context

SOPS has the ability to use [AWS KMS key policy and encryption context](http\://docs\.aws\.amazon\.com/kms/latest/developerguide/encryption\-context\.html)
to refine the access control of a given KMS master key\.

When creating a new file\, you can specify the encryption context in the
<code>\-\-encryption\-context</code> flag by comma separated list of key\-value pairs\:

```sh
$ sops edit --encryption-context Environment:production,Role:web-server test.dev.yaml
```

The format of the Encrypt Context string is <code>\<EncryptionContext Key\>\:\<EncryptionContext Value\>\,\<EncryptionContext Key\>\:\<EncryptionContext Value\>\,\.\.\.</code>

The encryption context will be stored in the file metadata and does
not need to be provided at decryption\.

Encryption contexts can be used in conjunction with KMS Key Policies to define
roles that can only access a given context\. An example policy is shown below\:

```json
{
  "Effect": "Allow",
  "Principal": {
    "AWS": "arn:aws:iam::111122223333:role/RoleForExampleApp"
  },
  "Action": "kms:Decrypt",
  "Resource": "*",
  "Condition": {
    "StringEquals": {
      "kms:EncryptionContext:AppName": "ExampleApp",
      "kms:EncryptionContext:FilePath": "/var/opt/secrets/"
    }
  }
}
```

<a id="key-rotation"></a>
### 2\.11   Key Rotation

It is recommended to renew the data key on a regular basis\. <code>sops</code> supports key
rotation via the <code>rotate</code> command\. Invoking it on an existing file causes <code>sops</code>
to reencrypt the file with a new data key\, which is then encrypted with the various
KMS and PGP master keys defined in the file\.

Add the <code>\-i</code> option to write the rotated file back\, instead of printing it to
stdout\.

```sh
$ sops rotate example.yaml
```

<a id="using-sops-yaml-conf-to-select-kms-pgp-and-age-for-new-files"></a>
### 2\.12   Using \.sops\.yaml conf to select KMS\, PGP and age for new files

It is often tedious to specify the <code>\-\-kms</code> <code>\-\-gcp\-kms</code> <code>\-\-pgp</code> and <code>\-\-age</code> parameters for creation
of all new files\. If your secrets are stored under a specific directory\, like a
<code>git</code> repository\, you can create a <code>\.sops\.yaml</code> configuration file at the root
directory to define which keys are used for which filename\.

Let\'s take an example\:

- file named <strong>something\.dev\.yaml</strong> should use one set of KMS A\, PGP and age
- file named <strong>something\.prod\.yaml</strong> should use another set of KMS B\, PGP and age
- other files use a third set of KMS C and PGP
- all live under <strong>mysecretrepo/something\.\{dev\,prod\,gcp\}\.yaml</strong>

Under those circumstances\, a file placed at <strong>mysecretrepo/\.sops\.yaml</strong>
can manage the three sets of configurations for the three types of files\:

```yaml
# creation rules are evaluated sequentially, the first match wins
creation_rules:
    # upon creation of a file that matches the pattern *.dev.yaml,
    # KMS set A as well as PGP and age is used
    - path_regex: \.dev\.yaml$
      kms: 'arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500,arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e+arn:aws:iam::361527076523:role/hiera-sops-prod'
      pgp: 'FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4'
      age: 'age129h70qwx39k7h5x6l9hg566nwm53527zvamre8vep9e3plsm44uqgy8gla'

    # prod files use KMS set B in the PROD IAM, PGP and age
    - path_regex: \.prod\.yaml$
      kms: 'arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e+arn:aws:iam::361527076523:role/hiera-sops-prod,arn:aws:kms:eu-central-1:361527076523:key/cb1fab90-8d17-42a1-a9d8-334968904f94+arn:aws:iam::361527076523:role/hiera-sops-prod'
      pgp: 'FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4'
      age: 'age129h70qwx39k7h5x6l9hg566nwm53527zvamre8vep9e3plsm44uqgy8gla'
      hc_vault_uris: "http://localhost:8200/v1/sops/keys/thirdkey"

    # gcp files using GCP KMS
    - path_regex: \.gcp\.yaml$
      gcp_kms: projects/mygcproject/locations/global/keyRings/mykeyring/cryptoKeys/thekey

    # Finally, if the rules above have not matched, this one is a
    # catchall that will encrypt the file using KMS set C as well as PGP
    # The absence of a path_regex means it will match everything
    - kms: 'arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500,arn:aws:kms:us-west-2:142069644989:key/846cfb17-373d-49b9-8baf-f36b04512e47,arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e'
      pgp: 'FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4'
```

When creating any file under <strong>mysecretrepo</strong>\, whether at the root or under
a subdirectory\, SOPS will recursively look for a <code>\.sops\.yaml</code> file\. If one is
found\, the filename of the file being created is compared with the filename
regexes of the configuration file\. The first regex that matches is selected\,
and its KMS and PGP keys are used to encrypt the file\. It should be noted that
the looking up of <code>\.sops\.yaml</code> is from the working directory \(CWD\) instead of
the directory of the encrypting file \(see [Issue 242](https\://github\.com/getsops/sops/issues/242)\)\.

The <code>path\_regex</code> checks the path of the encrypting file relative to the <code>\.sops\.yaml</code> config file\. Here is another example\:

- files located under directory <strong>development</strong> should use one set of KMS A
- files located under directory <strong>production</strong> should use another set of KMS B
- other files use a third set of KMS C
```yaml
creation_rules:
    # upon creation of a file under development,
    # KMS set A is used
    - path_regex: .*/development/.*
      kms: 'arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500,arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e+arn:aws:iam::361527076523:role/hiera-sops-prod'
      pgp: 'FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4'

    # prod files use KMS set B in the PROD IAM
    - path_regex: .*/production/.*
      kms: 'arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e+arn:aws:iam::361527076523:role/hiera-sops-prod,arn:aws:kms:eu-central-1:361527076523:key/cb1fab90-8d17-42a1-a9d8-334968904f94+arn:aws:iam::361527076523:role/hiera-sops-prod'
      pgp: 'FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4'

    # other files use KMS set C
    - kms: 'arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500,arn:aws:kms:us-west-2:142069644989:key/846cfb17-373d-49b9-8baf-f36b04512e47,arn:aws:kms:us-west-2:361527076523:key/5052f06a-5d3f-489e-b86c-57201e06f31e'
      pgp: 'FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4'
```

Creating a new file with the right keys is now as simple as

```sh
$ sops edit <newfile>.prod.yaml
```

Note that the configuration file is ignored when KMS or PGP parameters are
passed on the SOPS command line or in environment variables\.

<a id="specify-a-different-gpg-executable"></a>
### 2\.13   Specify a different GPG executable

SOPS checks for the <code>SOPS\_GPG\_EXEC</code> environment variable\. If specified\,
it will attempt to use the executable set there instead of the default
of <code>gpg</code>\.

Example\: place the following in your <code>\~/\.bashrc</code>

```bash
SOPS_GPG_EXEC = 'your_gpg_client_wrapper'
```

<a id="specify-a-different-gpg-key-server"></a>
### 2\.14   Specify a different GPG key server

By default\, SOPS uses the key server <code>keys\.openpgp\.org</code> to retrieve the GPG
keys that are not present in the local keyring\.
This is no longer configurable\. You can learn more about why from this write\-up\: [SKS Keyserver Network Under Attack](https\://gist\.github\.com/rjhansen/67ab921ffb4084c865b3618d6955275f)\.

<a id="key-groups"></a>
### 2\.15   Key groups

By default\, SOPS encrypts the data key for a file with each of the master keys\,
such that if any of the master keys is available\, the file can be decrypted\.
However\, it is sometimes desirable to require access to multiple master keys
in order to decrypt files\. This can be achieved with key groups\.

When using key groups in SOPS\, data keys are split into parts such that keys from
multiple groups are required to decrypt a file\. SOPS uses Shamir\'s Secret Sharing
to split the data key such that each key group has a fragment\, each key in the
key group can decrypt that fragment\, and a configurable number of fragments \(threshold\)
are needed to decrypt and piece together the complete data key\. When decrypting a
file using multiple key groups\, SOPS goes through key groups in order\, and in
each group\, tries to recover the fragment of the data key using a master key from
that group\. Once the fragment is recovered\, SOPS moves on to the next group\,
until enough fragments have been recovered to obtain the complete data key\.

By default\, the threshold is set to the number of key groups\. For example\, if
you have three key groups configured in your SOPS file and you don\'t override
the default threshold\, then one master key from each of the three groups will
be required to decrypt the file\.

Management of key groups is done with the <code>sops groups</code> command\.

For example\, you can add a new key group with 3 PGP keys and 3 KMS keys to the
file <code>my\_file\.yaml</code>\:

```sh
$ sops groups add --file my_file.yaml --pgp fingerprint1 --pgp fingerprint2 --pgp fingerprint3 --kms arn1 --kms arn2 --kms arn3
```

Or you can delete the 1st group \(group number 0\, as groups are zero\-indexed\)
from <code>my\_file\.yaml</code>\:

```sh
$ sops groups delete --file my_file.yaml 0
```

Key groups can also be specified in the <code>\.sops\.yaml</code> config file\,
like so\:

```yaml
creation_rules:
    - path_regex: .*keygroups.*
      key_groups:
          # First key group
          - pgp:
                - fingerprint1
                - fingerprint2
            kms:
                - arn: arn1
                  role: role1
                  context:
                      foo: bar
                - arn: arn2
                  aws_profile: myprofile
          # Second key group
          - pgp:
                - fingerprint3
                - fingerprint4
            kms:
                - arn: arn3
                - arn: arn4
          # Third key group
          - pgp:
                - fingerprint5
```

Given this configuration\, we can create a new encrypted file like we normally
would\, and optionally provide the <code>\-\-shamir\-secret\-sharing\-threshold</code> command line
flag if we want to override the default threshold\. SOPS will then split the data
key into three parts \(from the number of key groups\) and encrypt each fragment with
the master keys found in each group\.

For example\:

```sh
$ sops edit --shamir-secret-sharing-threshold 2 example.json
```

Alternatively\, you can configure the Shamir threshold for each creation rule in the <code>\.sops\.yaml</code> config
with <code>shamir\_threshold</code>\:

```yaml
creation_rules:
    - path_regex: .*keygroups.*
      shamir_threshold: 2
      key_groups:
          # First key group
          - pgp:
                - fingerprint1
                - fingerprint2
            kms:
                - arn: arn1
                  role: role1
                  context:
                      foo: bar
                - arn: arn2
                  aws_profile: myprofile
          # Second key group
          - pgp:
                - fingerprint3
                - fingerprint4
            kms:
                - arn: arn3
                - arn: arn4
          # Third key group
          - pgp:
                - fingerprint5
```

And then run <code>sops edit example\.json</code>\.

The threshold \(<code>shamir\_threshold</code>\) is set to 2\, so this configuration will require
master keys from two of the three different key groups in order to decrypt the file\.
You can then decrypt the file the same way as with any other SOPS file\:

```sh
$ sops decrypt example.json
```

<a id="key-service"></a>
### 2\.16   Key service

There are situations where you might want to run SOPS on a machine that
doesn\'t have direct access to encryption keys such as PGP keys\. The <code>sops</code> key
service allows you to forward a socket so that SOPS can access encryption
keys stored on a remote machine\. This is similar to GPG Agent\, but more
portable\.

SOPS uses a client\-server approach to encrypting and decrypting the data
key\. By default\, SOPS runs a local key service in\-process\. SOPS uses a key
service client to send an encrypt or decrypt request to a key service\, which
then performs the operation\. The requests are sent using gRPC and Protocol
Buffers\. The requests contain an identifier for the key they should perform
the operation with\, and the plaintext or encrypted data key\. The requests do
not contain any cryptographic keys\, public or private\.

<strong>WARNING\: the key service connection currently does not use any sort of
authentication or encryption\. Therefore\, it is recommended that you make sure
the connection is authenticated and encrypted in some other way\, for example
through an SSH tunnel\.</strong>

Whenever we try to encrypt or decrypt a data key\, SOPS will try to do so first
with the local key service \(unless it\'s disabled\)\, and if that fails\, it will
try all other remote key services until one succeeds\.

You can start a key service server by running <code>sops keyservice</code>\.

You can specify the key services the <code>sops</code> binary uses with <code>\-\-keyservice</code>\.
This flag can be specified more than once\, so you can use multiple key
services\. The local key service can be disabled with
<code>enable\-local\-keyservice\=false</code>\.

For example\, to decrypt a file using both the local key service and the key
service exposed on the unix socket located in <code>/tmp/sops\.sock</code>\, you can run\:

```sh
$ sops decrypt --keyservice unix:///tmp/sops.sock file.yaml`
```

And if you only want to use the key service exposed on the unix socket located
in <code>/tmp/sops\.sock</code> and not the local key service\, you can run\:

```sh
$ sops decrypt --enable-local-keyservice=false --keyservice unix:///tmp/sops.sock file.yaml
```

<a id="auditing"></a>
### 2\.17   Auditing

Sometimes\, users want to be able to tell what files were accessed by whom in an
environment they control\. For this reason\, SOPS can generate audit logs to
record activity on encrypted files\. When enabled\, SOPS will write a log entry
into a pre\-configured PostgreSQL database when a file is decrypted\. The log
includes a timestamp\, the username SOPS is running as\, and the file that was
decrypted\.

In order to enable auditing\, you must first create the database and credentials
using the schema found in <code>audit/schema\.sql</code>\. This schema defines the
tables that store the audit events and a role named <code>sops</code> that only has
permission to add entries to the audit event tables\. The default password for
the role <code>sops</code> is <code>sops</code>\. You should change this password\.

Once you have created the database\, you have to tell SOPS how to connect to it\.
Because we don\'t want users of SOPS to be able to control auditing\, the audit
configuration file location is not configurable\, and must be at
<code>/etc/sops/audit\.yaml</code>\. This file should have strict permissions such
that only the root user can modify it\.

For example\, to enable auditing to a PostgreSQL database named <code>sops</code> running
on localhost\, using the user <code>sops</code> and the password <code>sops</code>\,
<code>/etc/sops/audit\.yaml</code> should have the following contents\:

```yaml
backends:
    postgres:
        - connection_string: "postgres://sops:sops@localhost/sops?sslmode=verify-full"
```

You can find more information on the <code>connection\_string</code> format in the
[PostgreSQL docs](https\://www\.postgresql\.org/docs/current/static/libpq\-connect\.html\#libpq\-connstring)\.

Under the <code>postgres</code> map entry in the above YAML is a list\, so one can
provide more than one backend\, and SOPS will log to all of them\:

```yaml
backends:
    postgres:
        - connection_string: "postgres://sops:sops@localhost/sops?sslmode=verify-full"
        - connection_string: "postgres://sops:sops@remotehost/sops?sslmode=verify-full"
```

<a id="saving-output-to-a-file"></a>
### 2\.18   Saving Output to a File

By default SOPS just dumps all the output to the standard output\. We can use the
<code>\-\-output</code> flag followed by a filename to save the output to the file specified\.
Beware using both <code>\-\-in\-place</code> and <code>\-\-output</code> flags will result in an error\.

<a id="passing-secrets-to-other-processes"></a>
### 2\.19   Passing Secrets to Other Processes

In addition to writing secrets to standard output and to files on disk\, SOPS
has two commands for passing decrypted secrets to a new process\: <code>exec\-env</code>
and <code>exec\-file</code>\. These commands will place all output into the environment of
a child process and into a temporary file\, respectively\. For example\, if a
program looks for credentials in its environment\, <code>exec\-env</code> can be used to
ensure that the decrypted contents are available only to this process and never
written to disk\.

```sh
# print secrets to stdout to confirm values
$ sops decrypt out.json
{
        "database_password": "jf48t9wfw094gf4nhdf023r",
        "AWS_ACCESS_KEY_ID": "AKIAIOSFODNN7EXAMPLE",
        "AWS_SECRET_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

# decrypt out.json and run a command
# the command prints the environment variable and runs a script that uses it
$ sops exec-env out.json 'echo secret: $database_password; ./database-import'
secret: jf48t9wfw094gf4nhdf023r

# launch a shell with the secrets available in its environment
$ sops exec-env out.json 'sh'
sh-3.2# echo $database_password
jf48t9wfw094gf4nhdf023r

# the secret is not accessible anywhere else
sh-3.2$ exit
$ echo your password: $database_password
your password:
```

If the command you want to run only operates on files\, you can use <code>exec\-file</code>
instead\. By default\, SOPS will use a FIFO to pass the contents of the
decrypted file to the new program\. Using a FIFO\, secrets are only passed in
memory which has two benefits\: the plaintext secrets never touch the disk\, and
the child process can only read the secrets once\. In contexts where this won\'t
work\, eg platforms like Windows where FIFOs unavailable or secret files that need
to be available to the child process longer term\, the <code>\-\-no\-fifo</code> flag can be
used to instruct SOPS to use a traditional temporary file that will get cleaned
up once the process is finished executing\. <code>exec\-file</code> behaves similar to
<code>find\(1\)</code> in that <code>\{\}</code> is used as a placeholder in the command which will be
substituted with the temporary file path \(whether a FIFO or an actual file\)\.

```sh
# operating on the same file as before, but as a file this time
$ sops exec-file out.json 'echo your temporary file: {}; cat {}'
your temporary file: /tmp/.sops894650499/tmp-file
{
        "database_password": "jf48t9wfw094gf4nhdf023r",
        "AWS_ACCESS_KEY_ID": "AKIAIOSFODNN7EXAMPLE",
        "AWS_SECRET_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

# launch a shell with a variable TMPFILE pointing to the temporary file
$ sops exec-file --no-fifo out.json 'TMPFILE={} sh'
sh-3.2$ echo $TMPFILE
/tmp/.sops506055069/tmp-file291138648
sh-3.2$ cat $TMPFILE
{
        "database_password": "jf48t9wfw094gf4nhdf023r",
        "AWS_ACCESS_KEY_ID": "AKIAIOSFODNN7EXAMPLE",
        "AWS_SECRET_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
sh-3.2$ ./program --config $TMPFILE
sh-3.2$ exit

# try to open the temporary file from earlier
$ cat /tmp/.sops506055069/tmp-file291138648
cat: /tmp/.sops506055069/tmp-file291138648: No such file or directory
```

Additionally\, on unix\-like platforms\, both <code>exec\-env</code> and <code>exec\-file</code>
support dropping privileges before executing the new program via the
<code>\-\-user \<username\></code> flag\. This is particularly useful in cases where the
encrypted file is only readable by root\, but the target program does not
need root privileges to function\. This flag should be used where possible
for added security\.

To overwrite the default file name \(<code>tmp\-file</code>\) in <code>exec\-file</code> use the
<code>\-\-filename \<filename\></code> parameter\.

```sh
# the encrypted file can't be read by the current user
$ cat out.json
cat: out.json: Permission denied

# execute sops as root, decrypt secrets, then drop privileges
$ sudo sops exec-env --user nobody out.json 'sh'
sh-3.2$ echo $database_password
jf48t9wfw094gf4nhdf023r

# dropped privileges, still can't load the original file
sh-3.2$ id
uid=4294967294(nobody) gid=4294967294(nobody) groups=4294967294(nobody)
sh-3.2$ cat out.json
cat: out.json: Permission denied
```

<a id="using-the-publish-command"></a>
### 2\.20   Using the publish command

<code>sops publish \$file</code> publishes a file to a pre\-configured destination \(this lives in the SOPS
config file\)\. Additionally\, support re\-encryption rules that work just like the creation rules\.

This command requires a <code>\.sops\.yaml</code> configuration file\. Below is an example\:

```yaml
destination_rules:
    - s3_bucket: "sops-secrets"
      path_regex: s3/*
      recreation_rule:
          pgp: F69E4901EDBAD2D1753F8C67A64535C4163FB307
    - gcs_bucket: "sops-secrets"
      path_regex: gcs/*
      recreation_rule:
          pgp: F69E4901EDBAD2D1753F8C67A64535C4163FB307
    - vault_path: "sops/"
      vault_kv_mount_name: "secret/" # default
      vault_kv_version: 2 # default
      path_regex: vault/*
      omit_extensions: true
```

The above configuration will place all files under <code>s3/\*</code> into the S3 bucket <code>sops\-secrets</code>\,
all files under <code>gcs/\*</code> into the GCS bucket <code>sops\-secrets</code>\, and the contents of all files under
<code>vault/\*</code> into Vault\'s KV store under the path <code>secrets/sops/</code>\. For the files that will be
published to S3 and GCS\, it will decrypt them and re\-encrypt them using the
<code>F69E4901EDBAD2D1753F8C67A64535C4163FB307</code> pgp key\.

You would deploy a file to S3 with a command like\: <code>sops publish s3/app\.yaml</code>

To publish all files in selected directory recursively\, you need to specify <code>\-\-recursive</code> flag\.

If you don\'t want file extension to appear in destination secret path\, use <code>\-\-omit\-extensions</code>
flag or <code>omit\_extensions\: true</code> in the destination rule in <code>\.sops\.yaml</code>\.

<a id="publishing-to-vault"></a>
#### 2\.20\.1   Publishing to Vault

There are a few settings for Vault that you can place in your destination rules\. The first
is <code>vault\_path</code>\, which is required\. The others are optional\, and they are
<code>vault\_address</code>\, <code>vault\_kv\_mount\_name</code>\, <code>vault\_kv\_version</code>\.

SOPS uses the official Vault API provided by Hashicorp\, which makes use of [environment
variables](https\://www\.vaultproject\.io/docs/commands/\#environment\-variables) for
configuring the client\.

<code>vault\_kv\_mount\_name</code> is used if your Vault KV is mounted somewhere other than <code>secret/</code>\.
<code>vault\_kv\_version</code> supports <code>1</code> and <code>2</code>\, with <code>2</code> being the default\.

If the destination secret path already exists in Vault and contains the same data as the source
file\, it will be skipped\.

Below is an example of publishing to Vault \(using token auth with a local dev instance of Vault\)\.

```sh
$ export VAULT_TOKEN=...
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ sops decrypt vault/test.yaml
example_string: bar
example_number: 42
example_map:
    key: value
$ sops publish vault/test.yaml
uploading /home/user/sops_directory/vault/test.yaml to http://127.0.0.1:8200/v1/secret/data/sops/test.yaml ? (y/n): y
$ vault kv get secret/sops/test.yaml
====== Metadata ======
Key              Value
---              -----
created_time     2019-07-11T03:32:17.074792017Z
deletion_time    n/a
destroyed        false
version          3

========= Data =========
Key               Value
---               -----
example_map       map[key:value]
example_number    42
example_string    bar
```

<a id="important-information-on-types"></a>
## 3   Important information on types

<a id="yaml-json-env-and-ini-type-extensions"></a>
### 3\.1   YAML\, JSON\, ENV and INI type extensions

SOPS uses the file extension to decide which encryption method to use on the file
content\. <code>YAML</code>\, <code>JSON</code>\, <code>ENV</code>\, and <code>INI</code> files are treated as trees of data\, and key/values are
extracted from the files to only encrypt the leaf values\. The tree structure is also
used to check the integrity of the file\.

Therefore\, if a file is encrypted using a specific format\, it needs to be decrypted
in the same format\. The easiest way to achieve this is to conserve the original file
extension after encrypting a file\. For example\:

```sh
$ sops encrypt -i myfile.json
$ sops decrypt myfile.json
```

If you want to change the extension of the file once encrypted\, you need to provide
<code>sops</code> with the <code>\-\-input\-type</code> flag upon decryption\. For example\:

```sh
$ sops encrypt myfile.json > myfile.json.enc

$ sops decrypt --input-type json myfile.json.enc
```

When operating on stdin\, use the <code>\-\-input\-type</code> and <code>\-\-output\-type</code> flags as follows\:

```sh
$ cat myfile.json | sops decrypt --input-type json --output-type json /dev/stdin
```

<a id="json-and-json-binary-indentation"></a>
### 3\.2   JSON and JSON\_binary indentation

SOPS indents <code>JSON</code> files by default using one <code>tab</code>\. However\, you can change
this default behaviour to use <code>spaces</code> by either using the additional <code>\-\-indent\=2</code> CLI option or
by configuring <code>\.sops\.yaml</code> with the code below\.

The special value <code>0</code> disables indentation\, and <code>\-1</code> uses a single tab\.

```yaml
stores:
    json:
        indent: 2
    json_binary:
        indent: 2
```

<a id="yaml-indentation"></a>
### 3\.3   YAML indentation

SOPS indents <code>YAML</code> files by default using 4 spaces\. However\, you can change
this default behaviour by either using the additional <code>\-\-indent\=2</code> CLI option or
by configuring <code>\.sops\.yaml</code> with\:

```yaml
stores:
    yaml:
        indent: 2
```
> [!NOTE]
> The YAML emitter used by sops only supports values between 2 and 9\. If you specify 1\,
> or 10 and larger\, the indent will be 2\.

<a id="yaml-anchors"></a>
### 3\.4   YAML anchors

SOPS only supports a subset of <code>YAML</code>\'s many types\. Encrypting YAML files that
contain strings\, numbers and booleans will work fine\, but files that contain anchors
will not work\, because the anchors redefine the structure of the file at load time\.

This file will not work in SOPS\:

```yaml
bill-to:  &id001
    street: |
        123 Tornado Alley
        Suite 16
    city:   East Centerville
    state:  KS

ship-to:  *id001
```

SOPS uses the path to a value as additional data in the AEAD encryption\, and thus
dynamic paths generated by anchors break the authentication step\.

JSON and TEXT file types do not support anchors and thus have no such limitation\.

<a id="yaml-streams"></a>
### 3\.5   YAML Streams

<code>YAML</code> supports having more than one \"document\" in a single file\, while
formats like <code>JSON</code> do not\. SOPS is able to handle both\. This means the
following multi\-document will be encrypted as expected\:

```yaml-stream
---
data: foo
---
data: bar
```

Note that the <code>sops</code> metadata\, i\.e\. the hash\, etc\, is computed for the physical
file rather than each internal \"document\"\.

<a id="top-level-arrays"></a>
### 3\.6   Top\-level arrays

<code>YAML</code> and <code>JSON</code> top\-level arrays are not supported\, because SOPS
needs a top\-level <code>sops</code> key to store its metadata\.

This file will not work in SOPS\:

```yaml
---
  - some
  - array
  - elements
```

But this one will work because the <code>sops</code> key can be added at the same level as the
<code>data</code> key\.

```yaml
data:
    - some
    - array
    - elements
```

Similarly\, with <code>JSON</code> arrays\, this document will not work\:

```json
[
  "some",
  "array",
  "elements"
]
```

But this one will work just fine\:

```json
{
  "data": [
    "some",
    "array",
    "elements"
  ]
}
```

<a id="examples"></a>
## 4   Examples

Take a look into the [examples folder](https\://github\.com/getsops/sops/tree/main/examples) for detailed use cases of SOPS in a CI environment\. The section below describes specific tips for common use cases\.

<a id="creating-a-new-file"></a>
### 4\.1   Creating a new file

The command below creates a new file with a data key encrypted by KMS and PGP\.

```sh
$ sops edit --kms "arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500" --pgp C9CAB0AF1165060DB58D6D6B2653B624D620786D /path/to/new/file.yaml
```

<a id="encrypting-an-existing-file"></a>
### 4\.2   Encrypting an existing file

Similar to the previous command\, we tell SOPS to use one KMS and one PGP key\.
The path points to an existing cleartext file\, so we give <code>sops</code> the flag <code>\-e</code> to
encrypt the file\, and redirect the output to a destination file\.

```sh
$ export SOPS_KMS_ARN="arn:aws:kms:us-west-2:927034868273:key/fe86dd69-4132-404c-ab86-4269956b4500"
$ export SOPS_PGP_FP="C9CAB0AF1165060DB58D6D6B2653B624D620786D"
$ sops encrypt /path/to/existing/file.yaml > /path/to/new/encrypted/file.yaml
```

Decrypt the file with <code>\-d</code>\.

```sh
$ sops decrypt /path/to/new/encrypted/file.yaml
```

<a id="encrypt-or-decrypt-a-file-in-place"></a>
### 4\.3   Encrypt or decrypt a file in place

Rather than redirecting the output of <code>\-e</code> or <code>\-d</code>\, <code>sops</code> can replace the
original file after encrypting or decrypting it\.

```sh
# file.yaml is in cleartext
$ sops encrypt -i /path/to/existing/file.yaml
# file.yaml is now encrypted
$ sops decrypt -i /path/to/existing/file.yaml
# file.yaml is back in cleartext
```

<a id="encrypting-binary-files"></a>
### 4\.4   Encrypting binary files

SOPS primary use case is encrypting YAML and JSON configuration files\, but it
also has the ability to manage binary files\. When encrypting a binary\, SOPS will
read the data as bytes\, encrypt it\, store the encrypted base64 under
<code>tree\[\'data\'\]</code> and write the result as JSON\.

Note that the base64 encoding of encrypted data can actually make the encrypted
file larger than the cleartext one\.

In\-place encryption/decryption also works on binary files\.

```sh
$ dd if=/dev/urandom of=/tmp/somerandom bs=1024
count=512
512+0 records in
512+0 records out
524288 bytes (524 kB) copied, 0.0466158 s, 11.2 MB/s

$ sha512sum /tmp/somerandom
9589bb20280e9d381f7a192000498c994e921b3cdb11d2ef5a986578dc2239a340b25ef30691bac72bdb14028270828dad7e8bd31e274af9828c40d216e60cbe /tmp/somerandom

$ sops encrypt -i /tmp/somerandom
please wait while a data encryption key is being generated and stored securely

$ sops decrypt -i /tmp/somerandom

$ sha512sum /tmp/somerandom
9589bb20280e9d381f7a192000498c994e921b3cdb11d2ef5a986578dc2239a340b25ef30691bac72bdb14028270828dad7e8bd31e274af9828c40d216e60cbe /tmp/somerandom
```

<a id="extract-a-sub-part-of-a-document-tree"></a>
### 4\.5   Extract a sub\-part of a document tree

SOPS can extract a specific part of a YAML or JSON document\, by provided the
path in the <code>\-\-extract</code> command line flag\. This is useful to extract specific
values\, like keys\, without needing an extra parser\.

```sh
$ sops decrypt --extract '["app2"]["key"]' ~/git/svc/sops/example.yaml
-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBAPTMNIyHuZtpLYc7VsHQtwOkWYobkUblmHWRmbXzlAX6K8tMf3Wf
ImcbNkqAKnELzFAPSBeEMhrBN0PyOC9lYlMCAwEAAQJBALXD4sjuBn1E7Y9aGiMz
bJEBuZJ4wbhYxomVoQKfaCu+kH80uLFZKoSz85/ySauWE8LgZcMLIBoiXNhDKfQL
vHECIQD6tCG9NMFWor69kgbX8vK5Y+QL+kRq+9HK6yZ9a+hsLQIhAPn4Ie6HGTjw
fHSTXWZpGSan7NwTkIu4U5q2SlLjcZh/AiEA78NYRRBwGwAYNUqzutGBqyXKUl4u
Erb0xAEyVV7e8J0CIQC8VBY8f8yg+Y7Kxbw4zDYGyb3KkXL10YorpeuZR4LuQQIg
bKGPkMM4w5blyE1tqGN0T7sJwEx+EUOgacRNqM2ljVA=
-----END RSA PRIVATE KEY-----
```

The tree path syntax uses regular python dictionary syntax\, without the
variable name\. Extract keys by naming them\, and array elements by numbering
them\.

```sh
$ sops decrypt --extract '["an_array"][1]' ~/git/svc/sops/example.yaml
secretuser2
```

<a id="set-a-sub-part-in-a-document-tree"></a>
### 4\.6   Set a sub\-part in a document tree

SOPS can set a specific part of a YAML or JSON document\, by providing
the path and value in the <code>set</code> command\. This is useful to set specific
values\, like keys\, without needing an editor\.

```sh
$ sops set ~/git/svc/sops/example.yaml '["app2"]["key"]' '"app2keystringvalue"'
```

The tree path syntax uses regular python dictionary syntax\, without the
variable name\. Set to keys by naming them\, and array elements by
numbering them\.

```sh
$ sops set ~/git/svc/sops/example.yaml '["an_array"][1]' '"secretuser2"'
```

The value must be formatted as json\.

```sh
$ sops set ~/git/svc/sops/example.yaml '["an_array"][1]' '{"uid1":null,"uid2":1000,"uid3":["bob"]}'
```

<a id="showing-diffs-in-cleartext-in-git"></a>
### 4\.7   Showing diffs in cleartext in git

You most likely want to store encrypted files in a version controlled repository\.
SOPS can be used with git to decrypt files when showing diffs between versions\.
This is very handy for reviewing changes or visualizing history\.

To configure SOPS to decrypt files during diff\, create a <code>\.gitattributes</code> file
at the root of your repository that contains a filter and a command\.

```text
*.yaml diff=sopsdiffer
```

Here we only care about YAML files\. <code>sopsdiffer</code> is an arbitrary name that we map
to a SOPS command in the git configuration file of the repository\.

```sh
$ git config diff.sopsdiffer.textconv "sops decrypt"

$ grep -A 1 sopsdiffer .git/config
[diff "sopsdiffer"]
    textconv = "sops decrypt"
```

With this in place\, calls to <code>git diff</code> will decrypt both previous and current
versions of the target file prior to displaying the diff\. And it even works with
git client interfaces\, because they call git diff under the hood\!

<a id="encrypting-only-parts-of-a-file"></a>
### 4\.8   Encrypting only parts of a file

Note\: this only works on YAML and JSON files\, not on BINARY files\.

By default\, SOPS encrypts all the values of a YAML or JSON file and leaves the
keys in cleartext\. In some instances\, you may want to exclude some values from
being encrypted\. This can be accomplished by adding the suffix <strong>\_unencrypted</strong>
to any key of a file\. When set\, all values underneath the key that set the
<strong>\_unencrypted</strong> suffix will be left in cleartext\.

Note that\, while in cleartext\, unencrypted content is still added to the
checksum of the file\, and thus cannot be modified outside of SOPS without
breaking the file integrity check\.
This behavior can be modified using <code>\-\-mac\-only\-encrypted</code> flag or <code>\.sops\.yaml</code>
config file which makes SOPS compute a MAC only over values it encrypted and
not all values\.

The unencrypted suffix can be set to a different value using the
<code>\-\-unencrypted\-suffix</code> option\.

Conversely\, you can opt in to only encrypt some values in a YAML or JSON file\,
by adding a chosen suffix to those keys and passing it to the <code>\-\-encrypted\-suffix</code> option\.

A third method is to use the <code>\-\-encrypted\-regex</code> which will only encrypt values under
keys that match the supplied regular expression\.  For example\, this command\:

```sh
$ sops encrypt --encrypted-regex '^(data|stringData)$' k8s-secrets.yaml
```

will encrypt the values under the <code>data</code> and <code>stringData</code> keys in a YAML file
containing kubernetes secrets\.  It will not encrypt other values that help you to
navigate the file\, like <code>metadata</code> which contains the secrets\' names\.

Conversely\, you can opt in to only leave certain keys without encrypting by using the
<code>\-\-unencrypted\-regex</code> option\, which will leave the values unencrypted of those keys
that match the supplied regular expression\. For example\, this command\:

```sh
$ sops encrypt --unencrypted-regex '^(description|metadata)$' k8s-secrets.yaml
```

will not encrypt the values under the <code>description</code> and <code>metadata</code> keys in a YAML file
containing kubernetes secrets\, while encrypting everything else\.

You can also specify these options in the <code>\.sops\.yaml</code> config file\.

Note\: these four options <code>\-\-unencrypted\-suffix</code>\, <code>\-\-encrypted\-suffix</code>\, <code>\-\-encrypted\-regex</code> and <code>\-\-unencrypted\-regex</code> are
mutually exclusive and cannot all be used in the same file\.

<a id="encryption-protocol"></a>
## 5   Encryption Protocol

When SOPS creates a file\, it generates a random 256 bit data key and asks each
KMS and PGP master key to encrypt the data key\. The encrypted version of the data
key is stored in the <code>sops</code> metadata under <code>sops\.kms</code> and <code>sops\.pgp</code>\.

For KMS\:

```yaml
sops:
    kms:
        - enc: CiC6yCOtzsnFhkfdIslYZ0bAf//gYLYCmIu87B3sy/5yYxKnAQEBAQB4usgjrc7JxYZH3SLJWGdGwH//4GC2ApiLvOwd7Mv+cmMAAAB+MHwGCSqGSIb3DQEHBqBvMG0CAQAwaAYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAyGdRODuYMHbA8Ozj8CARCAO7opMolPJUmBXd39Zlp0L2H9fzMKidHm1vvaF6nNFq0ClRY7FlIZmTm4JfnOebPseffiXFn9tG8cq7oi
          enc_ts: 1439568549.245995
          arn: arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e
```

For PGP\:

```yaml
sops:
    pgp:
        - fp: 85D77543B3D624B63CEA9E6DBC17301B491B3F21
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
```

SOPS then opens a text editor on the newly created file\. The user adds data to the
file and saves it when done\.

Upon save\, SOPS browses the entire file as a key/value tree\. Every time SOPS
encounters a leaf value \(a value that does not have children\)\, it encrypts the
value with AES256\_GCM using the data key and a 256 bit random initialization
vector\.

Each file uses a single data key to encrypt all values of a document\, but each
value receives a unique initialization vector and has unique authentication data\.

Additional data is used to guarantee the integrity of the encrypted data
and of the tree structure\: when encrypting the tree\, key names are concatenated
into a byte string that is used as AEAD additional data \(aad\) when encrypting
values\. We expect that keys do not carry sensitive information\, and
keeping them in cleartext allows for better diff and overall readability\.

Any valid KMS or PGP master key can later decrypt the data key and access the
data\.

Multiple master keys allow for sharing encrypted files without sharing master
keys\, and provide a disaster recovery solution\. The recommended way to use SOPS
is to have two KMS master keys in different regions and one PGP public key with
the private key stored offline\. If\, by any chance\, both KMS master keys are
lost\, you can always recover the encrypted data using the PGP private key\.

<a id="message-authentication-code"></a>
### 5\.1   Message Authentication Code

In addition to authenticating branches of the tree using keys as additional
data\, SOPS computes a MAC on all the values to ensure that no value has been
added or removed fraudulently\. The MAC is stored encrypted with AES\_GCM and
the data key under tree \-\> <code>sops</code> \-\> <code>mac</code>\.
This behavior can be modified using <code>\-\-mac\-only\-encrypted</code> flag or <code>\.sops\.yaml</code>
config file which makes SOPS compute a MAC only over values it encrypted and
not all values\.

<a id="motivation"></a>
## 6   Motivation

> 📝 <strong>A note from the maintainers</strong>
>
> This section was written by the original authors of SOPS while they were
> working at Mozilla\. It is kept here for historical reasons and to provide
> technical background on the project\. It is not necessarily representative
> of the views of the current maintainers\, nor are they currently affiliated
> with Mozilla\.

Automating the distribution of secrets and credentials to components of an
infrastructure is a hard problem\. We know how to encrypt secrets and share them
between humans\, but extending that trust to systems is difficult\. Particularly
when these systems follow devops principles and are created and destroyed
without human intervention\. The issue boils down to establishing the initial
trust of a system that just joined the infrastructure\, and providing it access
to the secrets it needs to configure itself\.

<a id="the-initial-trust"></a>
### 6\.1   The initial trust

In many infrastructures\, even highly dynamic ones\, the initial trust is
established by a human\. An example is seen in Puppet by the way certificates are
issued\: when a new system attempts to join a Puppetmaster\, an administrator
must\, by default\, manually approve the issuance of the certificate the system
needs\. This is cumbersome\, and many puppetmasters are configured to auto\-sign
new certificates to work around that issue\. This is obviously not recommended
and far from ideal\.

AWS provides a more flexible approach to trusting new systems\. It uses a
powerful mechanism of roles and identities\. In AWS\, it is possible to verify
that a new system has been granted a specific role at creation\, and it is
possible to map that role to specific resources\. Instead of trusting new systems
directly\, the administrator trusts the AWS permission model and its automation
infrastructure\. As long as AWS keys are safe\, and the AWS API is secure\, we can
assume that trust is maintained and systems are who they say they are\.

<a id="kms-trust-and-secrets-distribution"></a>
### 6\.2   KMS\, Trust and secrets distribution

Using the AWS trust model\, we can create fine grained access controls to
Amazon\'s Key Management Service \(KMS\)\. KMS is a service that encrypts and
decrypts data with AES\_GCM\, using keys that are never visible to users of the
service\. Each KMS master key has a set of role\-based access controls\, and
individual roles are permitted to encrypt or decrypt using the master key\. KMS
helps solve the problem of distributing keys\, by shifting it into an access
control problem that can be solved using AWS\'s trust model\.

<a id="operational-requirements"></a>
### 6\.3   Operational requirements

When Mozilla\'s Services Operations team started revisiting the issue of
distributing secrets to EC2 instances\, we set a goal to store these secrets
encrypted until the very last moment\, when they need to be decrypted on target
systems\. Not unlike many other organizations that operate sufficiently complex
automation\, we found this to be a hard problem with a number of prerequisites\:

1. Secrets must be stored in YAML files for easy integration into hiera
1. Secrets must be stored in GIT\, and when a new CloudFormation stack is
   built\, the current HEAD is pinned to the stack\. \(This allows secrets to
   be changed in GIT without impacting the current stack that may
   autoscale\)\.
1. Entries must be encrypted separately\. Encrypting entire files as blobs makes
   git conflict resolution almost impossible\. Encrypting each entry
   separately is much easier to manage\.
1. Secrets must always be encrypted on disk \(admin laptop\, upstream
   git repo\, jenkins and S3\) and only be decrypted on the target
   systems

SOPS can be used to encrypt YAML\, JSON and BINARY files\. In BINARY mode\, the
content of the file is treated as a blob\, the same way PGP would encrypt an
entire file\. In YAML and JSON modes\, however\, the content of the file is
manipulated as a tree where keys are stored in cleartext\, and values are
encrypted\. hiera\-eyaml does something similar\, and over the years we learned
to appreciate its benefits\, namely\:

- diffs are meaningful\. If a single value of a file is modified\, only that
  value will show up in the diff\. The diff is still limited to only showing
  encrypted data\, but that information is already more granular that
  indicating that an entire file has changed\.
- conflicts are easier to resolve\. If multiple users are working on the
  same encrypted files\, as long as they don\'t modify the same values\,
  changes are easy to merge\. This is an improvement over the PGP
  encryption approach where unsolvable conflicts often happen when
  multiple users work on the same file\.

<a id="openpgp-integration"></a>
### 6\.4   OpenPGP integration

OpenPGP gets a lot of bad press for being an outdated crypto protocol\, and while
true\, what really made us look for alternatives is the difficulty of managing and
distributing keys to systems\. With KMS\, we manage permissions to an API\, not keys\,
and that\'s a lot easier to do\.

But PGP is not dead yet\, and we still rely on it heavily as a backup solution\:
all our files are encrypted with KMS and with one PGP public key\, with its
private key stored securely for emergency decryption in the event that we lose
all our KMS master keys\.

SOPS can be used without KMS entirely\, the same way you would use an encrypted
PGP file\: by referencing the pubkeys of each individual who has access to the file\.
It can easily be done by providing SOPS with a comma\-separated list of public keys
when creating a new file\:

```sh
$ sops edit --pgp "E60892BB9BD89A69F759A1A0A3D652173B763E8F,84050F1D61AF7C230A12217687DF65059EF093D3,85D77543B3D624B63CEA9E6DBC17301B491B3F21" mynewfile.yaml
```

<a id="threat-model"></a>
## 7   Threat Model

The security of the data stored using SOPS is as strong as the weakest
cryptographic mechanism\. Values are encrypted using AES256\_GCM which is the
strongest symmetric encryption algorithm known today\. Data keys are encrypted
in either KMS\, which also uses AES256\_GCM\, or PGP which uses either RSA or
ECDSA keys\.

Going from the most likely to the least likely\, the threats are as follows\:

<a id="compromised-aws-credentials-grant-access-to-kms-master-key"></a>
### 7\.1   Compromised AWS credentials grant access to KMS master key

An attacker with access to an AWS console can grant itself access to one of
the KMS master keys used to encrypt a <code>sops</code> data key\. This threat should be
mitigated by protecting AWS accesses with strong controls\, such as multi\-factor
authentication\, and also by performing regular audits of permissions granted
to AWS users\.

<a id="compromised-pgp-key"></a>
### 7\.2   Compromised PGP key

PGP keys are routinely mishandled\, either because owners copy them from
machine to machine\, or because the key is left forgotten on an unused machine
an attacker gains access to\. When using PGP encryption\, SOPS users should take
special care of PGP private keys\, and store them on smart cards or offline
as often as possible\.

<a id="factorized-rsa-key"></a>
### 7\.3   Factorized RSA key

SOPS doesn\'t apply any restriction on the size or type of PGP keys\. A weak PGP
keys\, for example 512 bits RSA\, could be factorized by an attacker to gain
access to the private key and decrypt the data key\. Users of SOPS should rely
on strong keys\, such as 2048\+ bits RSA keys\, or 256\+ bits ECDSA keys\.

<a id="weak-aes-cryptography"></a>
### 7\.4   Weak AES cryptography

A vulnerability in AES256\_GCM could potentially leak the data key or the KMS
master key used by a SOPS encrypted file\. While no such vulnerability exists
today\, we recommend that users keep their encrypted files reasonably private\.

<a id="backward-compatibility"></a>
## 8   Backward compatibility

SOPS will remain backward compatible on the major version\, meaning that all
improvements brought to the 1\.X and 2\.X branches \(current\) will maintain the
file format introduced in <strong>1\.0</strong>\.

<a id="security"></a>
## 9   Security

Please report any security issues privately using [GitHub\'s advisory form](https\://github\.com/getsops/sops/security/advisories)\.

<a id="license"></a>
## 10   License

Mozilla Public License Version 2\.0

<a id="authors"></a>
## 11   Authors

SOPS was initially launched as a project at Mozilla in 2015 and has been
graciously donated to the CNCF as a Sandbox project in 2023\, now under the
stewardship of a [new group of maintainers](https\://github\.com/getsops/community/blob/main/MAINTAINERS\.md)\.

The original authors of the project were\:

- Adrian Utrilla \@autrilla
- Julien Vehent \@jvehent

Furthermore\, the project has been carried for a long time by AJ Bahnken \@ajvb\,
and had not been possible without the contributions of numerous [contributors](https\://github\.com/getsops/sops/graphs/contributors)\.

<a id="credits"></a>
## 12   Credits

SOPS was inspired by [hiera\-eyaml](https\://github\.com/TomPoulton/hiera\-eyaml)\,
[credstash](https\://github\.com/LuminalOSS/credstash)\,
[sneaker](https\://github\.com/codahale/sneaker)\,
[password store](http\://www\.passwordstore\.org/) and too many years managing
PGP encrypted files by hand\.\.\.

---

<img src="docs/images/cncf-color-bg.svg" alt="CNCF Sandbox Project" width="400">

<strong>We are a</strong> [Cloud Native Computing Foundation](https\://cncf\.io) <strong>sandbox project\.</strong>
