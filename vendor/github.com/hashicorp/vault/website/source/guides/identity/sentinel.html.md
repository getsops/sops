---
layout: "guides"
page_title: "Sentinel - Guides"
sidebar_title: "Sentinel Policies"
sidebar_current: "guides-identity-sentinel"
description: |-
  Vault Enterprise supports Sentinel to provide a rich set of access control
  functionality. This guide walks through the creation and usage of role
  governing policies (RGPs) and endpoint governing policies (EGPs).
---

# Sentinel Policies

~> **Enterprise Only:** Sentinel is a part of _Vault Enterprise Premium_.

_Sentinel_ is a language framework for policy build to be embedded in Vault
Enterprise to enable fine-grained, logic-based policy decisions which cannot be
fully handled by the ACL policies.

**Role Governing Policies (RGPs)** and **Endpoint Governing
Policies (EGPs)** can be defined using Sentinel:

- RGPs are tied to particular tokens, identity entities, or identity groups
- EGPs are tied to particular paths (e.g. `aws/creds/`)

> This guide walks you through the authoring of Sentinel policies in Vault.  For
ACL policy authoring, refer to the [Policies](/guides/identity/policies.html)
guide.


## Reference Material

- [Sentinel Getting Started Guide](https://docs.hashicorp.com/sentinel/intro/getting-started/first-policy)
- [Sentinel](https://docs.hashicorp.com/sentinel/) documentation
- [Vault Sentinel](/docs/enterprise/sentinel/index.html) documentation
- [Security and Fundamentals at Scale with Vault](https://www.youtube.com/watch?time_continue=121&v=yiPbKICFkvQ)
- [Identity - Entities and Groups](/guides/identity/identity.html) guide


## Estimated Time to Complete

5 - 10 minutes


## Challenge

ACL policies are ***path-based*** that it has the following challenges:

- Cannot grant permissions based on logics other than paths
- Paths are merged in ACL policies which could potentially cause a conflict as
the number of policies grows

What if the policy requirement was to grant read permission on `secret/orders`
path ***only if*** the request came from an IP address within a certain CIDR?


## Solution

Use Sentinel policies (RGPs and/or EGPs) to fulfill more complex policy
requirements.

Sentinel can access properties of the incoming requests and make a decision
based on a certain set of conditions. Available properties include:

- **request** - Information about the request itself (path, operation type,
 parameters, etc.)
- **token** - Information about the token being used (creation time, attached
  policies, etc.)
- **identity** - Identity entities and all related data
- **mfa** - Information about successful MFA validations


## Prerequisites

To perform the tasks described in this guide, you need to have a ***Vault
Enterprise*** environment.  

### Policy requirements

Since this guide demonstrates the creation of policies, log in with highly
privileged token such as **`root`**. Required permissions are:

```shell
# To list policies
path "sys/policies/*"
{
  capabilities = ["list"]
}

# Create and manage EGPs
path "sys/policies/egp/*"
{
  capabilities = ["create", "read", "update", "delete", "list"]
}
```


## Steps

This guide demonstrates basic Sentinel policy authoring and management tasks.

1. [Write Sentinel Policies](#step1)
1. [Test the Sentinel Policies](#step2)
1. [Deploy your EGP policies](#step3)
1. [Delete Sentinel Policies](#step4)


### <a name="step1"></a>Step 1: Write Sentinel Policies

#### Anatomy of Sentinel Policies

```hcl
import "<library>"

<variable> = <value>

main = rule {
     <conditions_to_evaluate>
}
```

- **`import`** - Enables your policy to access reusable libraries. There are a
set of built-in [imports](https://docs.hashicorp.com/sentinel/imports/)
available to help define your policy rules.

- **`main`** (required) - Every Sentinel policy must have a **`main`** rule
which is evaluated to determine the result of a policy.

- **`rule`** - A first-class construct in Sentinel. It describes a set of
conditions resulting in either true or false. (NOTE: Refer to the [Boolean
Expressions](https://docs.hashicorp.com/sentinel/language/boolexpr) for the full
list of available operators in writing rules.)

- **`<variable>`** - Variables are dynamically typed in Sentinel. You can define
its value explicitly or implicitly by the host system or [function](https://docs.hashicorp.com/sentinel/language/functions).

~> **NOTE:** The Sentinel language supports many features such as functions,
loops, slices, etc. You can learn about all of this in the [complete language
guide](https://docs.hashicorp.com/sentinel/language/).

#### Policy requirements

In this guide, you are going to write Sentinel policies that fulfill the
following requirements:

1. Any incoming request against the "`secret/accounting/*`" to be performed
during the business hours (7:00 am to 6:00 pm during the work days).

1. Any `create`, `update` and `delete` operations against Key/Value secret
engine (mounted at "`secret`") **must** come from an internal IP of
`122.22.3.4/32` CIDR.


#### Sentinel Policies

Requirement #1: **`business-hrs.sentinel`**

```shell
import "time"

# Expect requests to only happen during work days (Monday through Friday)
# 0 for Sunday and 6 for Saturday
workdays = rule {
    time.now.weekday > 0 and time.now.weekday < 6
}

# Expect requests to only happen during work hours (7:00 am - 6:00 pm)
workhours = rule {
    time.now.hour > 7 and time.now.hour < 18
}

main = rule {
    workdays and workhours
}
```

Requirement #2: **`cidr-check.sentinel`**

```shell
import "sockaddr"
import "strings"

# Only care about create, update, and delete operations against secret path
precond = rule {
    request.operation in ["create", "update", "delete"] and
    strings.has_prefix(request.path, "secret/")
}

# Requests to come only from our private IP range
cidrcheck = rule {
    sockaddr.is_contained(request.connection.remote_addr, "122.22.3.4/32")
}

# Check the precondition before execute the cidrcheck
main = rule when precond {
    cidrcheck
}
```

> **NOTE:** The **`main`** has conditional rule (`when precond`) to ensure that
the rule gets evaluated only if the request is relevant.

~> Refer to the [Sentinel Properties](/docs/enterprise/sentinel/properties.html)
documentation for available properties which Vault injects to Sentinel to allow
fine-grained controls.



### <a name="step2"></a>Step 2: Test the Sentinel Policies

You can test the Sentinel policies prior to deployment in orders to validate
syntax and to document expected behavior.

1. First, you need to download the
[Sentinel simulator](https://docs.hashicorp.com/sentinel/downloads.html).

    **Example:**

    ```plaintext
    $ wget https://releases.hashicorp.com/sentinel/0.3.1/sentinel_0.3.1_darwin_amd64.zip
    $ unzip sentinel_0.3.1_darwin_amd64.zip -d /usr/local/bin
    ```

1. Create a sub-folder named, **`test`** where `cidr-check.sentinel` and
`business-hrs.sentinel` policies are located. Under the `test` folder, you want
to create a sub-folder for each policy: **`cidr-check`** and **`business-hrs`**.

    ```plaintext
    $ mkdir -p test/business-hrs
    $ mkdir -p test/cidr-check
    ```

    > **NOTE:** The test should be created under `/test/<policy_name>` folder.

1. Write a passing test case in a file named, **`success.json`** under
`test/business-hrs` directory.

    ```plaintext
    {
        "global": {
            "timespace": {
                "weekday": 1,
                "hour": 12
            }
        }
    }
    ```

    Under **`global`**, you specify the mock test data. In this example, the
    `weekday` is set to `1` which is **`Monday`** and `hour` is set to `12`
    which is **`noon`**. Therefore, the `main` should return `true`.

1. Write a failing test in a file named, **`fail.json`** under
`test/business-hrs`.

    ```plaintext
    {
        "global": {
            "timespace": {
                "weekday": 0,
                "hour": 12
            }
        }
    }
    ```

    The mock data is set to **`Sunday`** at **`noon`**; therefore, Therefore,
    the `main` should return `false`.


1. Similarly, write a passing test case for `cidr-check` policy,
**`test/cidr-check/success.json`**:


    ```plaintext
    {
      "global": {
        "request": {
          "connection": {
              "remote_addr": "122.22.3.4"
          },
          "operation": "create",
          "path": "secret/orders"        
        }
      }
    ```

    In this example, the `global` specifies the `create` operation is invoked on
    `secret/orders` endpoint which initiated from an IP address `122.22.3.4`.
    Therefore, the `main` should return `true`.

1. Write a failing test for `cidr-check` policy, **`test/cidr-check/fail.json`**.

    ```plaintext
    {
      "global": {
        "request": {
          "connection": {
            "remote_addr": "122.22.3.10"
          },
          "operation": "create",
          "path": "secret/orders"
        }
      },
      "test": {
        "precond": true,
        "main": false
      }
    }
    ```

    This test will fail because of the IP address mismatch. However, the
    `precond` should pass since the requested operation is `create` and the
    targeted endpoint is `secret/orders`.

    > The optional **`test`** definition adds more context to why the test
    should fail.  The expected behavior is that the test fails because `main`
    returns `false` but `precond` should return `true`.

1. Now, you have written both success and failure tests:

    ```plaintext
    ├── business-hrs.sentinel
    ├── cidr-check.sentinel
    └── test
        ├── business-hrs
        │   ├── fail.json
        │   └── success.json
        └── cidr-check
            ├── fail.json
            └── success.json
    ```

1. Execute the test:

    ```plaintext
    $ sentinel test

    PASS - business-hrs.sentinel
      PASS - test/business-hrs/success.json  PASS - test/business-hrs/fail.json
    PASS - cidr-check.sentinel
      PASS - test/cidr-check/success.json  PASS - test/cidr-check/fail.json
    ```

    > **NOTE:** If you want to see the tracing and log output for those tests,
    run the command with `-verbose` flag.

### <a name="step3"></a>Step 3: Deploy your EGP policies

Sentinel policies has three **enforcement levels**:

| Level          | Description                                                   |
|----------------|---------------------------------------------------------------|
| advisory       | The policy is allowed to fail. Can be used as a tool to educate new users. |
| soft-mandatory | The policy must pass unless an override is specified.         |
| hard-mandatory | The policy must pass no matter what!                          |

<br>

Since both policies are tied to specific paths, the policy type that you are
going to create is Endpoint Governing Policies (EGPs).

#### CLI command

1. Store the Base64 encoded `cidr-check.sentinel` policy in an environment
variable named `POLICY`.

    ```plaintext
    $ POLICY=$(base64 cidr-check.sentinel)
    ```

1. Create a policy `cidr-check` with enforcement level of **hard-mandatory** to
reject all requests coming from IP addressed that are not internal.

    ```plaintext
    $ vault write sys/policies/egp/cidr-check \
            policy="${POLICY}" \
            paths="secret/*" \
            enforcement_level="hard-mandatory"
    ```

1. You can read the policy by executing the following command:

    ```plaintext
    $ vault read sys/policies/egp/cidr-check    
    ```

1. Repeat the steps to create a policy named `business-hrs`.

    ```shell
    # Encode the business-hrs policy
    $ POLICY2=$(base64 business-hrs.sentinel)

    # Create a policy with soft-mandatory enforcement-level
    $ vault write sys/policies/egp/business-hrs \
            policy="${POLICY2}" \
            paths="secret/accounting/*" \
            enforcement_level="soft-mandatory"

    # To read the policy you just created
    $ vault read sys/policies/egp/business-hrs        
    ```


#### API call using cURL

To create EGP policies, use the `/sys/policies/egp` endpoint:

```plaintext
$ curl --header "X-Vault-Token: <TOKEN>" \
       --request PUT \
       --data <PAYLOAD> \
       <VAULT_ADDRESS>/v1/sys/policies/egp/<POLICY_NAME>
```

Where `<TOKEN>` is your valid token, and `<PAYLOAD>` includes the Base64 encoded
policy, endpoint paths, and enforcement level.


1. Store the Base64 encoded `cidr-check.sentinel` policy in an environment
variable named `POLICY`.

    ```plaintext
    $ POLICY=$(base64 cidr-check.sentinel)
    ```

1. Create API request payload.

    ```plaintext
    $ tee cidr-payload.json <<EOF
    {
      "policy": "${POLICY}",
      "paths": ["secret/*"],
      "enforcement_level": "hard-mandatory"
    }
    EOF
    ```

1. Create a policy `cidr-check` with enforcement level of **hard-mandatory** to
reject all requests coming from IP addressed that are not internal.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @cidr-payload.json \
           http://127.0.0.1:8200/v1/sys/policies/egp/cidr-check
    ```

1. Repeat the steps to create a policy named `business-hrs` with enforcement
level of soft-mandatory.

    ```shell
    # Encode the business-hrs policy
    $ POLICY2=$(base64 business-hrs.sentinel)

    # Create the request payload
    $ tee buz-hrs-payload.json <<EOF
    {
      "policy": "${POLICY2}",
      "paths": ["secret/accounting/*"],
      "enforcement_level": "soft-mandatory"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request PUT \
           --data @buz-hrs-payload.json \
           http://127.0.0.1:8200/v1/sys/policies/egp/business-hrs
    ```

1. You can list the EGPs that were created.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request LIST \
           http://127.0.0.1:8200/v1/sys/policies/egp | jq
    ```

#### Web UI

Open a web browser and launch the Vault UI (e.g. http://127.0.0.1:8200/ui) and
then login.

1. Select **Policies** and select the **Endpoint Governing Policies** tab.

1. Select **Create EGP policy**.

1. Enter **`business-hrs`** in the **Name** field.  

1. Enter the [**`business-hrs.sentinel`** policy](#sentinel-policies-1) in the
**Policy** editor.

1. Select **soft-mandatory** from the **Enforcement level** drop-down list.

1. Enter **`secret/accounting/*`** in the **Paths** field, and then click
**Create Policy**.

  ![EGP](/img/vault-sentinel-1.png)

1. Select **Endpoint Governing Policies** again, and then **Create EGP policy**.

1. Enter **`cidr-check`** in the **Name** field.  

1. Enter the [**`cidr-check.sentinel`** policy](#sentinel-policies-1) in the
**Policy** editor.

1. Leave the **Enforcement level** as hard-mandatory, and enter **`secret/*`**
in the **Paths** field.

1. Click **Create Policy**.

<br>

~> **NOTE:** Unlike ACL policies, EGPs are a _prefix walk_ which allows policies
to be applied at various points at Vault API.  If you have EGPs tied to
"**`secret/orders`**", "**`secret/*`**" and "**`*`**", all EGPs will be
evaluated for a request on "**`secret/orders`**".


#### Verification

Once the policies were deployed, `create`, `update` and `delete` operations
coming from an IP address other than `122.22.3.4` will be denied.

```plaintext
$ vault kv put secret/accounting/test acct_no="293472309423"

Error writing data to secret/accounting/test: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/secret/accounting/test
Code: 400. Errors:

* 1 error occurred:

* egp standard policy "cidr-check" evaluation resulted in denial.

The specific error was:
<nil>

A trace of the execution for policy "cidr-check" is available:

Result: false

Description: Check the precondition before execute the cidrcheck

Rule "main" (byte offset 442) = false
  false (offset 314): sockaddr.is_contained

Rule "cidrcheck" (byte offset 291) = false

Rule "precond" (byte offset 113) = true
  true (offset 134): request.operation in ["create", "update", "delete"]
  true (offset 194): strings.has_prefix
```

Similarly, you will get an error if any request is made outside of the business
hours defined by the `business-hrs` policy.

!> **NOTE:** Like with ACL policies, **`root`** tokens are ***NOT*** subject to
Sentinel policy checks.


### <a name="step4"></a>Step 4: Delete Sentinel Policies

#### CLI Command

To delete EGPs:

```shell
# Delete the business-hrs EGP
$ vault delete sys/policies/egp/business-hrs

# Delete the cidr-check EGP
$ vault delete sys/policies/egp/cidr-check
```


#### API call using cURL

To delete EGPs:

```shell
# Delete the business-hrs EGP
$ curl --header "X-Vault-Token: ..." \
       --request DELETE \
       http://127.0.0.1:8200/v1/sys/policies/egp/business-hrs

# Delete the cidr-check EGP
$ curl --header "X-Vault-Token: ..." \
      --request DELETE \
      http://127.0.0.1:8200/v1/sys/policies/egp/cidr-check
```

#### Web UI

1. Select **Policies** and select the **Endpoint Governing Policies** tab.

1. Select **Delete** from the policy menu for `business-hrs`.

    ![Delete EGP](/img/vault-sentinel-2.png)

1. When prompted, click **Delete** again to confirm.

1. Repeat the steps to delete `cidr-check` policy.



## Next steps

Refer to the [Sentinel Properties](/docs/enterprise/sentinel/properties.html)
documentation for the full list of properties available in Vault to write
fine-grained policies to meet your organizational policy requirements.
