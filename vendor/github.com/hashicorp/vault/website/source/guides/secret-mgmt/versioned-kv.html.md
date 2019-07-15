---
layout: "guides"
page_title: "Versioned KV Secret Engine - Guides"
sidebar_title: "Versioned KV Secret Engine"
sidebar_current: "guides-secret-mgmt-versioned-kv"
description: |-
  Vault 0.10.0 introduced version 2 of key-value secret engine which supports
  versioning of your secrets so that you can undo the accidental deletion of
  secrets, or compare the different versions of the secret.
---

# Versioned Key/Value Secret Engine

The [Static Secrets](/guides/secret-mgmt/static-secrets.html) guide introduced
the basics of working with key-value secret engine. **Vault 0.10** introduced [_K/V
Secrets Engine v2 with Secret
Versioning_](https://www.hashicorp.com/blog/vault-0-10).  This guide
demonstrates the new features introduced by the key-value secret engine v2.


## Reference Material

- [Static Secrets guide](/guides/secret-mgmt/static-secrets.html)
- [KV Secrets Engine - Version 2](/docs/secrets/kv/kv-v2.html)
- [KV Secrets Engine - Version 2 (API)](/api/secret/kv/kv-v2.html)

~> **NOTE:** An [interactive
tutorial](https://www.katacoda.com/hashicorp/scenarios/vault-static-secrets) is
also available if you do not have a Vault environment to perform the steps
described in this guide.


## Estimated Time to Complete

10 minutes


## Challenge

The KV secret engine v1 does not provide a way to version or roll back secrets.
This made it difficult to recover from unintentional data loss or overwrite when
more than one user is writing at the same path.


## Solution

Run the **version 2** of KV secret engine which can retain a configurable
number of secret versions. This enables older versions' data to be retrievable
in case of unwanted deletion or updates of the data.  In addition, its
_Check-and-Set_ operations can be used to protect the data from being overwritten
unintentionally.

![Versioned KV](/img/vault-versioned-kv-1.png)

## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault. Make sure
that your Vault server has been [initialized and
unsealed](/intro/getting-started/deploy.html).

### Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for
initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# To view in Web UI
path "sys/mounts" {
  capabilities = [ "read", "update" ]
}

# Write and manage secrets in key-value secret engine
path "secret*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# To enable secret engines
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

This guide demonstrates the basic commands for working with KV secret engine v2.  

You will perform the following:

1. [Check the KV secret engine version](#step1)
2. [Write secrets](#step2)
3. [Retrieve a specific version of secret](#step3)
4. [Specify the number of versions to keep](#step4)
5. [Delete versions of secret](#step5)
6. [Permanently delete data](#step6)


### <a name="step1"></a>Step 1: Check the KV secret engine version
(**Persona:** devops)

Before beginning, verify that you are using the v2 of the KV secret engine.

#### CLI command

To check the KV secret engine version:

```plaintext
$ vault secrets list -format=json
...
"secret/": {
  "type": "kv",
  "description": "key/value secret storage",
  "accessor": "kv_f05b8b9c",
  "config": {
    "default_lease_ttl": 0,
    "max_lease_ttl": 0,
    "force_no_cache": false
  },
  "options": {
    "version": "2"
  },
  ...
```

The indicated **`version`** should be **`2`**. If the version is **`1`**,
upgrade it to v2.

```plaintext
$ vault kv enable-versioning secret/
```

#### API call using cURL

To check the KV secret engine version:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       <VAULT_ADDRESS>/v1/sys/mounts
```

Where `<TOKEN>` is your valid token, and `<VAULT_ADDRESS>` is where your vault
server is running.


**Example:**

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/sys/mounts | jq
...
  "secret/": {
    "accessor": "kv_f05b8b9c",
    "config": {
      "default_lease_ttl": 0,
      "force_no_cache": false,
      "max_lease_ttl": 0,
      "plugin_name": ""
    },
    "description": "key/value secret storage",
    "local": false,
    "options": {
      "version": "2"
    },
    "seal_wrap": false,
    "type": "kv"
  },
...
```

The indicated **`version`** should be **`2`**. If the version is **`1`**,
upgrade it to v2.

```plaintext
$ cat payload.json
{
  "options": {
      "version": "2"
  }
}

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/sys/mounts/secret/tune
```

#### Web UI

Open a web browser and launch the Vault UI (e.g. `http://127.0.0.1:8200/ui`) and
then login.

![Web UI](/img/vault-versioned-kv-2.png)

If `secret/` does not indicates **`v2`**, you can upgrade it from `v1` to `v2`
by executing the following CLI command:

```plaintext
$ vault kv enable-versioning secret/
```

Alternatively, you can enable KV secret engine v2 at a different path by
clicking **Enable new engine**. Select **KV** from the **Secret engine type**
drop-down list.  Be sure that the **Version** is set to be **Version 2**.

![Enabling kv-v2](/img/vault-versioned-kv-3.png)

Click **Enable Engine** to complete.


### <a name="step2"></a>Step 2: Write Secrets

To understand how the versioning works, let's write some test data.

#### CLI commands

To write secrets, run `vault kv put` command instead of `vault write`:

```plaintext
$ vault kv put secret/customer/acme name="ACME Inc." contact_email="jsmith@acme.com"
Key              Value
---              -----
created_time     2018-04-14T00:05:47.115378933Z
deletion_time    n/a
destroyed        false
version          1
```

To update the existing secret, run the `vault kv put` command again:

```plaintext
$ vault kv put secret/customer/acme name="ACME Inc." contact_email="john.smith@acme.com"
Key              Value
---              -----
created_time     2018-04-14T00:13:35.296018431Z
deletion_time    n/a
destroyed        false
version          2
```

Now you have two versions of the `secret/customer/acme` data. Run `vault kv get`
to read the data.

```plaintext
$ vault kv get secret/customer/acme
====== Metadata ======
Key              Value
---              -----
created_time     2018-04-14T00:13:35.296018431Z
deletion_time    n/a
destroyed        false
version          2

======== Data ========
Key              Value
---              -----
contact_email    john.smith@acme.com
name             ACME Inc.
```

#### API call using cURL

Write some data at `secret/customer/acme`:

```plaintext
$ tee payload.json <<EOF
{
  "data": {
    "name": "ACME Inc.",
    "contact_email": "jsmith@acme.com"
  }
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret/data/customer/acme
```

Notice that the endpoint for KV v2 is **`/secret/data/<path>`**; therefore, to
write secrets at `secret/customer/acme`, the API endpoint becomes
`/secret/data/customer/acme`.

Update the secret to create another version:

```plaintext
$ tee payload.json <<EOF
{
  "data": {
    "name": "ACME Inc.",
    "contact_email": "john.smith@acme.com"
  }
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret/data/customer/acme
```

Now you have two versions of the `secret/customer/acme` data. Read back the secret.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/data/customer/acme
{
   "request_id": "7233b69d-35d9-6c1b-ae81-9a679a03082d",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 0,
   "data": {
     "data": {
       "contact_email": "john.smith@acme.com",
       "name": "ACME Inc."
     },
     "metadata": {
       "created_time": "2018-04-14T00:59:11.27903511Z",
       "deletion_time": "",
       "destroyed": false,
       "version": 2
     }
   },
   "wrap_info": null,
   "warnings": null,
   "auth": null
}
```


#### Web UI

In the Web UI, select `secret/` and then click **Create secret**.

![Write Secret](/img/vault-versioned-kv-5.png)

Click **Save**.

To update the existing secret, select **Edit**, change the `contact_email`
value, and then click **Save**.

![Write Secret](/img/vault-versioned-kv-6.png)


### <a name="step3"></a>Step 3: Retrieve a Specific Version of Secret

You may run into a situation where you need to view the secret before an update.

#### CLI commands

To retrieve the version 1 of the secret written at `secret/customer/acme`:

```plaintext
$ vault kv get -version=1 secret/customer/acme
====== Metadata ======
Key              Value
---              -----
created_time     2018-04-14T00:05:47.115378933Z
deletion_time    n/a
destroyed        false
version          1

======== Data ========
Key              Value
---              -----
contact_email    jsmith@acme.com
name             ACME Inc.
```

To read the **metadata** of `secret/customer/acme`:

```plaintext
$ vault kv metadata get secret/customer/acme
======= Metadata =======
Key                Value
---                -----
created_time       2018-04-14T00:05:47.115378933Z
current_version    2
max_versions       0
oldest_version     0
updated_time       2018-04-14T00:13:35.296018431Z

====== Version 1 ======
Key              Value
---              -----
created_time     2018-04-14T00:05:47.115378933Z
deletion_time    n/a
destroyed        false

====== Version 2 ======
Key              Value
---              -----
created_time     2018-04-14T00:13:35.296018431Z
deletion_time    n/a
destroyed        false
```


#### API call using cURL

To retrieve the version 1 of the secret written at `secret/customer/acme`:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/data/customer/acme?version=1 | jq
{
 "request_id": "3bf5a2c1-d89b-9dd5-9bb5-0bc61a4a6d83",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
   "data": {
     "contact_email": "jsmith@acme.com",
     "name": "ACME Inc."
   },
   "metadata": {
     "created_time": "2018-04-14T00:05:47.115378933Z",
     "deletion_time": "",
     "destroyed": false,
     "version": 1
   }
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```

To read the **metadata** of `secret/customer/acme`:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/metadata/customer/acme | jq
{
 "request_id": "34708262-59cd-9a94-247f-3b1db0909050",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
   "created_time": "2018-04-14T00:05:47.115378933Z",
   "current_version": 2,
   "max_versions": 0,
   "oldest_version": 0,
   "updated_time": "2018-04-14T00:13:35.296018431Z",
   "versions": {
     "1": {
       "created_time": "2018-04-14T00:05:47.115378933Z",
       "deletion_time": "",
       "destroyed": false
     },
     "2": {
       "created_time": "2018-04-14T00:13:35.296018431Z",
       "deletion_time": "",
       "destroyed": false
     }
   }
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```


### <a name="step4"></a>Step 4: Specify the number of versions to keep

By default, the `kv-v2` secret engine keeps up to 10 versions.  Let's limit the
maximum number of versions to keep to be 4.

#### CLI command

To set the `secret/` to keep up to 4 versions:

```shell
$ vault write secret/config max_versions=4
Success! Data written to: secret/config

# View the configuration settings
$ vault read secret/config
Key             Value
---             -----
cas_required    false
max_versions    4
```

Alternatively, to limit the number of versions only on the
**`secret/customer/acme`** path rather than the entire `secret/` engine:

```plaintext
$ vault kv metadata put -max-versions=4 secret/customer/acme
```

Overwrite the data a few more times to see what happens to the data.

```plaintext
$ vault kv metadata get secret/customer/acme
======= Metadata =======
Key                Value
---                -----
created_time       2018-04-14T00:42:25.677078177Z
current_version    6
max_versions       0
oldest_version     3
updated_time       2018-04-16T00:17:23.930473344Z

====== Version 3 ======
Key              Value
---              -----
created_time     2018-04-16T00:15:59.880368849Z
deletion_time    n/a
destroyed        false

====== Version 4 ======
Key              Value
---              -----
created_time     2018-04-16T00:16:18.941331243Z
deletion_time    n/a
destroyed        false

====== Version 5 ======
Key              Value
---              -----
created_time     2018-04-16T00:16:34.407951572Z
deletion_time    n/a
destroyed        false

====== Version 6 ======
Key              Value
---              -----
created_time     2018-04-16T00:17:23.930473344Z
deletion_time    n/a
destroyed        false
```

In this example, the current version is 6. Notice that version 1 and 2 do not
show up in the metadata. Because the kv secret engine is configured to keep only
4 versions, the oldest two versions are permanently deleted and you won't be
able to read them.

```plaintext
$ vault kv get -version=1 secret/customer/acme
No value found at secret/data/customer/data
```

#### API call using cURL

To set the `secret/` to keep up to 4 versions:

```plaintext
$ tee payload.json<<EOF
{
  "max_versions": 4,
  "cas_required": false
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json
       http://127.0.0.1:8200/v1/secret/config
```

To view the configuration:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/config | jq
{
 "request_id": "8addfed1-41eb-6a19-8342-93f493c51538",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
   "cas_required": false,
   "max_versions": 4
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```

Alternatively, to limit the number of versions only on the
**`secret/customer/acme`** path rather than the entire `secret/` engine:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json
       http://127.0.0.1:8200/v1/secret/metadata/customer/acme
```

Invoke the `secret/metadata/customer/acme` endpoint instead.


Overwrite the data a few more times to see what happens to the data.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/metadata/customer/acme | jq
{
 "request_id": "f2dd7f69-294c-e5c3-d582-f723005ea243",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
   "created_time": "2018-04-14T00:42:25.677078177Z",
   "current_version": 6,
   "max_versions": 0,
   "oldest_version": 3,
   "updated_time": "2018-04-16T00:17:23.930473344Z",
   "versions": {
     "3": {
       "created_time": "2018-04-16T00:15:59.880368849Z",
       "deletion_time": "",
       "destroyed": false
     },
     "4": {
       "created_time": "2018-04-16T00:16:18.941331243Z",
       "deletion_time": "",
       "destroyed": false
     },
     "5": {
       "created_time": "2018-04-16T00:16:34.407951572Z",
       "deletion_time": "",
       "destroyed": false
     },
     "6": {
       "created_time": "2018-04-16T00:17:23.930473344Z",
       "deletion_time": "",
       "destroyed": false
     }
   }
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```

In this example, the current version is 6. Notice that version 1 and 2 do not
show up in the metadata. Because the kv secret engine is configured to keep only
4 versions, the oldest two versions are permanently deleted and you won't be
able to read them.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       http://127.0.0.1:8200/v1/secret/data/customer/acme?version=1 | jq
{
 "errors": []
}
```


### <a name="step5"></a>Step 5: Delete versions of secret


#### CLI command

Let's delete versions 4 and 5:

```shell
$ vault kv delete -versions="4,5" secret/customer/acme
Success! Data deleted (if it existed) at: secret/customer/acme

# Check the metadata
$ vault kv metadata get secret/customer/acme
...
====== Version 4 ======
Key              Value
---              -----
created_time     2018-04-16T00:12:25.404198622Z
deletion_time    2018-04-16T01:04:01.160426888Z
destroyed        false

====== Version 5 ======
Key              Value
---              -----
created_time     2018-04-16T00:12:47.527981267Z
deletion_time    2018-04-16T01:04:01.160427742Z
destroyed        false
...
```

The metadata on versions 4 and 5 reports its deletion timestamp
(`deletion_time`); however, the `destroyed` parameter is set to `false`.

If version 5 was deleted by mistake and you wish to recover, invoke the `vault
kv undelete` command:

```plaintext
$ vault kv undelete -versions=5 secret/customer/acme
Success! Data written to: secret/undelete/customer/acme
```


#### API call using cURL

Let's delete versions 4 and 5:

```shell
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{ "versions":[4,5] }'
       http://127.0.0.1:8200/v1/secret/delete/customer/acme

# Check the metadata
$ curl --header "X-Vault-Token: ..." \
      http://127.0.0.1:8200/v1/secret/metadata/customer/acme | jq
...
"4": {
   "created_time": "2018-04-16T00:16:18.941331243Z",
   "deletion_time": "2018-04-16T01:17:42.003111567Z",
   "destroyed": false
 },
 "5": {
   "created_time": "2018-04-16T00:16:34.407951572Z",
   "deletion_time": "2018-04-16T01:17:42.003111978Z",
   "destroyed": false
 },
...
```

The metadata on versions 4 and 5 reports its deletion timestamp
(`deletion_time`); however, the `destroyed` parameter is set to `false`.

If version 5 was deleted by mistake and you wish to recover, invoke the
`/secret/undelete` endpoint:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{ "versions":[5] }'
       http://127.0.0.1:8200/v1/secret/undelete/customer/acme
```


### <a name="step6"></a>Step 6: Permanently delete data

#### CLI command

To permanently delete a version of secret:

```shell
$ vault kv destroy -versions=4 secret/customer/acme
Success! Data written to: secret/destroy/customer/acme

# Check the metadata
$ vault kv metadata get secret/customer/acme
...
====== Version 4 ======
Key              Value
---              -----
created_time     2018-04-16T00:12:25.404198622Z
deletion_time    2018-04-16T01:04:01.160426888Z
destroyed        true
...
```

The metadata indicates that Version 4 is destroyed.

If you wish to destroy all the keys and versions at `secret/customer/acme`,
invoke the `vault kv metadata delete` command:

```plaintext
$ vault kv metadata delete secret/customer/acme
Success! Data deleted (if it existed) at: secret/metadata/customer/acme
```


#### API call using cURL

To permanently delete a version of secret:

```shell
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{ "versions":[4] }'
       http://127.0.0.1:8200/v1/secret/destroy/customer/acme

# Check the metadata
$ curl --header "X-Vault-Token: ..." \
      http://127.0.0.1:8200/v1/secret/metadata/customer/acme | jq
...
  "4": {
    "created_time": "2018-04-16T00:16:18.941331243Z",
    "deletion_time": "2018-04-16T01:17:42.003111567Z",
    "destroyed": true
  },
...
```

The metadata indicates that Version 4 is destroyed.

If you wish to destroy all the keys and versions at `secret/customer/acme`,
invoke the `secret/metadata` endpoint:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request DELETE
       http://127.0.0.1:8200/v1/secret/metadata/customer/acme
```


### Additional Discussion

The v2 of KV secret engine supports a **_Check-And-Set_** operation to prevent
unintentional secret overwrite. When you pass the `cas` flag to Vault, it first
checks if the key already exists.

By default, _Check-And-Set_ operation is not enabled on the KV secret engine;
therefore, write is always allowed (no checking is performed).  

```plaintext
$ vault read secret/config
Key             Value
---             -----
cas_required    false
max_versions    0
```

#### CLI command

To enable the **_Check-And-Set_** operation:

```shell
# Enable cas_requied on the secret engine mounted at secret/
$ vault write secret/config cas-required=true

# Enable cas_requied only on the secret/partner path
$ vault kv metadata put -cas-required=true secret/partner
```

Once check-and-set is enabled, every write operation requires `cas` value to be
passed. If you are sure that you want to overwrite the existing key-value, set
`cas` to match the current version. Set `cas` to `0` if you want to write the
secret _only if_ the key does not exists.

**Example:**

```shell
# To write if the key does not already exists
$ vault kv put -cas=0 secret/partner name="Example Co." partner_id="123456789"
Key              Value
---              -----
created_time     2018-04-16T22:58:15.798753323Z
deletion_time    n/a
destroyed        false
version          1

# To overwrite the secret, you must specify the current version with -cas flag
$ vault kv put -cas=1 secret/partner name="Example Co." partner_id="ABCDEFGHIJKLMN"
Key              Value
---              -----
created_time     2018-04-16T23:00:28.66552289Z
deletion_time    n/a
destroyed        false
version          2
```

#### API call using cURL

To enable the **_Check-And-Set_** operation:

```shell
$ tee payload.json<<EOF
{
  "max_versions": 10,
  "cas_required": true
}
EOF

# Enable cas_requied on the secret engine mounted at secret/
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json
       http://127.0.0.1:8200/v1/secret/config

# Enable cas_requied only on the secret/partner path
$ curl --header "X-Vault-Token: ..." \
      --request POST \
      --data @payload.json
      http://127.0.0.1:8200/v1/secret/metadata/partner
```

Once check-and-set is enabled, every write operation requires `cas` value to be
passed. If you are sure that you want to overwrite the existing key-value, set
`cas` to match the current version. Set `cas` to `0` if you want to write the
secret _only if_ the key does not exists.

**Example:**

```shell
# Write if the key does not already exists
$ tee payload.json <<EOF
{
  "options": {
    "cas": 0
  },
  "data": {
    "name": "Example Co.",
    "partner_id": "123456789"
  }
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret/data/partner

# To overwrite the secret, you must pass the current version
$ tee payload.json <<EOF
{
  "options": {
    "cas": 1
  },
  "data": {
    "name": "Example Co.",
    "partner_id": "ABCDEFGHIJKLMN"
  }
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       http://127.0.0.1:8200/v1/secret/data/partner
```

<br>

~> If the **`cas`** value is missing in your write request, the
"`check-and-set parameter required for this call`" error will be returned.  If
the `cas` does not match the current version number, you will receive the
"`check-and-set parameter did not match the current version`" message.



## Next steps

This guide introduced the CLI commands and API endpoints to read and write
static secrets in the key-value secret engine. Read [Secret as a Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html) guide to learn about the
usage of database secret engine.
