---
layout: "guides"
page_title: "Vault Auto-unseal using AWS KMS - Guides"
sidebar_title: "Vault Auto-unseal with AWS KMS"
sidebar_current: "guides-operations-autounseal-aws-kms"
description: |-
  In this guide, we'll show an example of how to use Terraform to provision an
  instance that can utilize an encryption key from AWS Key Management Services
  to unseal Vault.
---


# Vault Auto-unseal using AWS Key Management Service

~> **Enterprise Only:** Vault auto-unseal feature is a part of _Vault Enterprise Pro_.

When a Vault server is started, it starts in a
[***sealed***](/docs/concepts/seal.html) state and it does not know how to
decrypt data. Before any operation can be performed on the Vault, it must be
unsealed. Unsealing is the process of constructing the master key necessary to
decrypt the data encryption key.

![Unseal with Shamir's Secret Sharing](/img/vault-autounseal.png)

This guide demonstrates an example of how to use Terraform to provision an
instance that can utilize an encryption key from [AWS Key Management Services
(KMS)](https://aws.amazon.com/kms/) to unseal Vault.

## Reference Material

- [Vault Auto Unseal](/docs/configuration/seal/index.html)
- [Configuration: `awskms` Seal](/docs/configuration/seal/awskms.html)


## Estimated Time to Complete

10 minutes

## Personas

The steps described in this guide are typically performed by **operations**
persona.

## Challenge

Vault unseal operation requires a quorum of existing unseal keys split by
Shamir's Secret sharing algorithm. This is done so that the "_keys to the
kingdom_" won't fall into one person's hand.  However, this process is manual
and can become painful when you have many Vault clusters as there are now
many different key holders with many different keys.

## Solution

Vault Enterprise supports opt-in automatic unsealing via cloud technologies:
Amazon KMS, Azure Key Vault or GCP Cloud KMS. This feature enables operators to
delegate the unsealing process to trusted cloud providers to ease operations in
the event of partial failure and to aid in the creation of new or ephemeral
clusters.

![Unseal with AWS KMS](/img/vault-autounseal-2.png)

## Prerequisites

This guide assumes the following:   

- Access to **Vault Enterprise 0.9.0 or later** 
- A URL to download Vault Enterprise from (an Amazon S3 bucket will suffice)
- AWS account for provisioning cloud resources
- [Terraform installed](https://www.terraform.io/intro/getting-started/install.html)
and basic understanding of its usage

### Download demo assets

Clone or download the demo assets from the
[hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/operations/aws-kms-unseal/terraform-aws)
GitHub repository to perform the steps described in this guide.


## Steps

This guide demonstrates how to implement and use the Auto-unseal feature using
AWS KMS. Included is a Terraform configuration that has the following:   

* Ubuntu 16.04 LTS with Vault Enterprise    
* An instance profile granting the Amazon EC2 instance to an AWS KMS key
* Vault configured with access to an AWS KMS key   


[![YouTube](/img/vault-autounseal-4.png)](https://youtu.be/iRyqOEDFIiY)


You are going to perform the following steps:

1. [Provision the Cloud Resources](#step-1-provision-the-cloud-resources)
1. [Test the Auto-unseal Feature](#step-2-test-the-auto-unseal-feature)
1. [Clean Up](#step-3-clean-up)


### Step 1: Provision the Cloud Resources

**Task 1:** Be sure to set your working directory to where the
[`/operations/aws-kms-unseal/terraform-aws`](#download-demo-assets) folder is
located.

The working directory should contain the provided Terraform files:

```bash
~/git/vault-guides/operations/aws-kms-unseal/terraform$ tree
.
├── README.md
├── instance-profile.tf
├── instance.tf
├── main.tf
├── ssh-key.tf
├── terraform.tfvars.example
├── userdata.tpl
└── variables.tf
```

**Task 2:** Set your AWS credentials as environment variables:

```plaintext
$ export AWS_ACCESS_KEY_ID = "<YOUR_AWS_ACCESS_KEY_ID>"

$ export AWS_SECRET_ACCESS_KEY = "<YOUR_AWS_SECRET_ACCESS_KEY>"
```

Create a file named **`terraform.tfvars`** and specify your Vault Enterprise
binary download URL.

**Example:**

```plaintext
vault_url = "https://s3-us-west-2.amazonaws.com/hc-enterprise-binaries/vault/ent/0.10.3/vault-enterprise_0.10.3%2Bent_linux_amd64.zip"
```

**Task 3:** Perform a **`terraform init`** to pull down the necessary provider
resources. Then **`terraform plan`** to verify your changes and the resources that
will be created. If all looks good, then perform a **`terraform apply`** to
provision the resources.

```plaintext
$ terraform init
Initializing provider plugins...
...
Terraform has been successfully initialized!


$ terraform plan
...
Plan: 15 to add, 0 to change, 0 to destroy.


$ terraform apply
...
Apply complete! Resources: 15 added, 0 changed, 0 destroyed.

Outputs:

connections = Connect to Vault via SSH   ssh ubuntu@192.0.2.1 -i private.key
Vault Enterprise web interface  http://192.0.2.1:8200/ui
```

**NOTE:** The Terraform output will display the public IP address to SSH into
your server as well as the Vault Enterprise web interface address.


### Step 2: Test the Auto-unseal Feature

SSH into the provisioned EC2 instance.

```plaintext
$ ssh ubuntu@192.0.2.1 -i private.key
...
Are you sure you want to continue connecting (yes/no)? yes
```
When you are prompted, enter "yes" to continue.

To verify that Vault has been installed, run `vault status` command which should
return "_server is not yet initialized_" message.

```plaintext
$ export VAULT_ADDR=http://127.0.0.1:8200

$ vault status
Error checking seal status: Error making API request.

URL: GET http://127.0.0.1:8200/v1/sys/seal-status
Code: 400. Errors:

* server is not yet initialized
```

Run the **`vault operator init`** command to initialize the Vault server by
setting its key share to be **`1`** as follow:

```plaintext
$ vault operator init -stored-shares=1 -recovery-shares=1 -recovery-threshold=1 -key-shares=1 -key-threshold=1
Recovery Key 1: oOxAQfxcZitjqZfF3984De8rUckPeahQDUvmJ1A4JrQ=
Initial Root Token: 54c4dbe3-d45b-79d9-18d0-602831a6a991

Vault initialized successfully.

Recovery key initialized with 1 keys and a key threshold of 1. Please
securely distribute the above keys.
```

Stop and start the Vault server:

```plaintext
$ sudo systemctl stop vault

$ vault status
Error checking seal status: Get http://127.0.0.1:8200/v1/sys/seal-status: dial tcp 127.0.0.1:8200: getsockopt: connection refused

$ sudo systemctl start vault
```

Check the Vault status to verify that it has been started and unsealed.

```plaintext
$ vault status
Type: shamir
Sealed: false
Key Shares: 1
Key Threshold: 1
Unseal Progress: 0
Unseal Nonce:
Version: 0.9.6+prem.hsm
Cluster Name: vault-cluster-01cf6f33
Cluster ID: fb787d8a-b882-fee8-b461-445320cde311

High-Availability Enabled: false
```

Log into Vault using the generated initial root token:

```plaintext
$ vault login 54c4dbe3-d45b-79d9-18d0-602831a6a991
Successfully authenticated! You are now logged in.
token: 54c4dbe3-d45b-79d9-18d0-602831a6a991
token_duration: 0
token_policies: [root]
```

Review the Vault configuration file (`/etc/vault.d/vault.hcl`).

```plaintext
$ cat /etc/vault.d/vault.hcl
storage "file" {
  path = "/opt/vault"
}
listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}
seal "awskms" {
  kms_key_id = "d7c1ffd9-8cce-45e7-be4a-bb38dd205966"
}
ui=true
```

Notice the Vault configuration file defines the [`awskms`
stanza](/docs/configuration/seal/awskms.html) which sets the AWS KMS key ID to
use for encryption and decryption.

At this point, you should be able to launch the Vault Enterprise UI by entering
the address provided in the `terraform apply` outputs (e.g. http://192.0.2.1:8200/ui)
and log in with your initial root token.

![Vault Enterprise UI Login](/img/vault-autounseal-3.png)


### Step 3: Clean Up

Once completed, execute the following commands to clean up:

```plaintext
$ terraform destroy -force

$ rm -rf .terraform terraform.tfstate* private.key
```


## Next steps

Once you have a Vault environment setup, the next step is to write policies.
Read [Policies](/guides/identity/policies.html) to learn how to write policies
to govern the behavior of clients.
