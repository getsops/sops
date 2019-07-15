---
layout: "docs"
page_title: "Google Cloud - Auth Methods"
sidebar_title: "Google Cloud"
sidebar_current: "docs-auth-gcp"
description: |-
  The "gcp" auth method allows users and machines to authenticate to Vault using
  Google Cloud service accounts.
---

# Google Cloud Auth Method

The `gcp` auth method allows Google Cloud Platform entities to authenticate to
Vault. Vault treats Google Cloud as a trusted third party and verifies
authenticating entities against the Google Cloud APIs. This backend allows for
authentication of:

- Google Cloud IAM service accounts
- Google Compute Engine (GCE) instances

This backend focuses on identities specific to Google _Cloud_ and does not
support authenticating arbitrary Google or G Suite users or generic OAuth
against Google.

This plugin is developed in a separate GitHub repository at
[hashicorp/vault-plugin-auth-gcp][repo],
but is automatically bundled in Vault releases. Please file all feature
requests, bugs, and pull requests specific to the GCP plugin under that
repository.


## Authenticate

### Via the CLI Helper

Vault includes a CLI helper that obtains a signed JWT locally and sends the
request to Vault. This helper is only available for IAM-type roles.

```text
$ vault login -method=gcp \
    role="my-role" \
    service_account="authenticating-account@my-project.iam.gserviceaccount.com" \
    project="my-project" \
    jwt_exp="15m" \
    credentials=@path/to/signer/credentials.json
```

For more usage information, run `vault auth help gcp`.

### Via the CLI

```text
$ vault write -field=token auth/gcp/login \
    role="my-role" \
    jwt="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

See [Generating JWTs](#generating-jwts) for ways to obtain the JWT token.

### Via the API

```text
$ curl \
    --request POST \
    --data '{"role":"my-role", "jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}' \
    http://127.0.0.1:8200/v1/auth/gcp/login
```

See [API docs][api-docs] for expected response.

## Configuration

Auth methods must be configured in advance before users or machines can
authenticate. These steps are usually completed by an operator or configuration
management tool.

1. Enable the Google Cloud auth method:

    ```text
    $ vault auth enable gcp
    ```

1. Configure the auth method credentials:

    ```text
    $ vault write auth/gcp/config \
        credentials=@/path/to/credentials.json
    ```

    If you are using instance credentials or want to specify credentials via
    an environment variable, you can skip this step. To learn more, see the
    [Google Cloud Authentication](#google-cloud-authentication) section below.

1. Create a named role:

    For an `iam`-type role:

    ```text
    $ vault write auth/gcp/role/my-iam-role \
        type="iam" \
        policies="dev,prod" \
        bound_service_accounts="my-service@my-project.iam.gserviceaccount.com"
    ```

    For a `gce`-type role:

    ```text
    $ vault write auth/gcp/role/my-gce-role \
        type="gce" \
        policies="dev,prod" \
        bound_projects="my-project1,my-project2" \
        bound_zones="us-east1-b" \
        bound_labels="foo:bar,zip:zap"
    ```

    For the complete list of configuration options for each type, please see the
    [API documentation][api-docs].


## Authentication

The Google Cloud Vault auth method uses the official Google Cloud Golang SDK.
This means it supports the common ways of [providing credentials to Google
Cloud][cloud-creds].

1. The environment variable `GOOGLE_APPLICATION_CREDENTIALS`. This is specified
as the **path** to a Google Cloud credentials file, typically for a service
account. If this environment variable is present, the resulting credentials are
used. If the credentials are invalid, an error is returned.

1. Default instance credentials. When no environment variable is present, the
default service account credentials are used.

For more information on service accounts, please see the [Google Cloud Service
Accounts documentation][service-accounts].

To use this auth method, the service account must have the following minimum
scope:

```text
https://www.googleapis.com/auth/cloud-platform
```

### Required GCP Permissions

#### Vault Server Permissions

**For `iam`-type Vault roles**, Vault can be given the following roles:

```text
roles/iam.serviceAccountKeyAdmin
```

**For `gce`-type Vault roles**, Vault can be given the following roles:

```text
roles/compute.viewer
```

If you instead wish to create a custom role with only the exact GCP permissions
required, use the following list of permissions:

```text
iam.serviceAccounts.get
iam.serviceAccountKeys.get
compute.instances.get
compute.instanceGroups.list
compute.instanceGroups.listInstances
```

These allow Vault to:

* verify the service account, either directly authenticating or associated with 
  authenticating GCE instance, exists
* get the corresponding public keys for verifying JWTs signed by service account
  private keys.
* verify authenticating GCE instances exist
* compare bound fields for GCE roles (zone/region, labels, or membership
  in given instance groups)

#### Permissions For Authenticating Against Vault

Note that the previously mentioned permissions are given to the _Vault servers_. 
The IAM service account or GCE instance that is **authenticating against Vault**
must have the following role:

```text
roles/iam.serviceAccountTokenCreator
```

## Group Aliases

As of Vault 1.0, roles can specify an `add_group_aliases` boolean parameter
that adds [group aliases][identity-group-aliases] to the auth response. These
aliases can aid in building reusable policies since they are available as
interpolated values in Vault's policy engine. Once enabled, the auth response
will include the following aliases:

```json
[
  "project-$PROJECT_ID",
  "folder-$SUBFOLDER_ID",
  "folder-$FOLDER_ID",
  "organization-$ORG_ID"
]
```


## Implementation Details

This section describes the implementation details for how Vault communicates
with Google Cloud to authenticate and authorize JWT tokens. This information is
provided for those who are curious, but these details are not
required knowledge for using the auth method.

### IAM Login

IAM login applies only to roles of type `iam`. The Vault authentication workflow
for IAM service accounts looks like this:

[![Vault Google Cloud IAM Login Workflow](/img/vault-gcp-iam-auth-workflow.svg)](/img/vault-gcp-iam-auth-workflow.svg)

  1. The client generates a signed JWT using the IAM
  [`projects.serviceAccounts.signJwt`][signjwt-method] method. For examples of
  how to do this, see the [Obtaining JWT Tokens](#obtaining-jwt-tokens) section.

  2. The client sends this signed JWT to Vault along with a role name.

  3. Vault extracts the `kid` header value, which contains the ID of the
  key-pair used to generate the JWT, and the `sub` ID/email to find the service
  account key. If the service account does not exist or the key is not linked to
  the service account, Vault denies authentication.

  4. Vault authorizes the confirmed service account against the given role. If
  that is successful, a Vault token with the proper policies is returned.

### GCE Login

GCE login only applies to roles of type `gce` and **must be completed on an
instance running in GCE**. These steps will not work from your local laptop or
another cloud provider.

[![Vault Google Cloud GCE Login Workflow](/img/vault-gcp-gce-auth-workflow.svg)](/img/vault-gcp-gce-auth-workflow.svg)

  1. The client obtains an [instance identity metadata token][instance-identity]
  on a GCE instance.

  2. The client sends this JWT to Vault along with a role name.

  3. Vault extracts the `kid` header value, which contains the ID of the
  key-pair used to generate the JWT, to find the OAuth2 public cert to verify
  this JWT.

  4. Vault authorizes the confirmed instance against the given role, ensuring
  the instance matches the bound zones, regions, or instance groups. If that is
  successful, a Vault token with the proper policies is returned.


## Generating JWTs

This section details the various methods and examples for obtaining JWT
tokens.

### IAM

This describes how to use the GCP IAM [API method][signjwt-method] directly
to generate the signed JWT with the claims that Vault expects. Note the CLI
does this process for you and is much easier, and that there is very little
reason to do this yourself.

#### curl Example

Vault requires the following minimum claim set:

```json
{
  "sub": "$SERVICE_ACCOUNT_EMAIL_OR_ID",
  "aud": "vault/$ROLE",
  "exp": "$EXPIRATION"
}
```

For the API method, expiration is optional and will default to an hour.
If specified, expiration must be a
[NumericDate](https://tools.ietf.org/html/rfc7519#section-2) value (seconds from
Epoch). This value must be before the max JWT expiration allowed for a role.
This defaults to 15 minutes and cannot be more than 1 hour.

One you have all this information, the JWT token can be signed using curl and
[oauth2l](https://github.com/google/oauth2l):

```text
ROLE="my-role"
PROJECT="my-project"
SERVICE_ACCOUNT="service-account@my-project.iam.gserviceaccount.com"
OAUTH_TOKEN="$(oauth2l header cloud-platform)"
JWT_CLAIM="{\\\"aud\\\":\\\"vault/${ROLE}\\\", \\\"sub\\\": \\\"${SERVICE_ACCOUNT}\\\"}"

curl \
  --header "${OAUTH_TOKEN}" \
  --header "Content-Type: application/json" \
  --request POST \
  --data "{\"payload\": \"${JWT_CLAIM}\"}" \
  "https://iam.googleapis.com/v1/projects/${PROJECT}/serviceAccounts/${SERVICE_ACCOUNT}:signJwt"
```

#### gcloud Example

You can also do this through the (currently beta) gcloud command.

```text
$ gcloud beta iam service-accounts sign-jwt $INPUT_JWT_CLAIMS $OUTPUT_JWT_FILE \
    --iam-account=service-account@my-project.iam.gserviceaccount.com \
    --project=my-project
```

#### Golang Example

Read more on the
[Google Open Source blog](https://opensource.googleblog.com/2017/08/hashicorp-vault-and-google-cloud-iam.html).

### GCE

GCE tokens **can only be generated from a GCE instance**. The JWT token can be
obtained from the `service-accounts/default/identity` endpoint for a
instance's metadata server.

#### curl Example

```text
ROLE="my-gce-role"

curl \
  --header "Metadata-Flavor: Google" \
  --get \
  --data-urlencode "audience=http://vault/${ROLE}" \
  --data-urlencode "format=full" \
  "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity"
```

## API

The GCP Auth Plugin has a full HTTP API. Please see the
[API docs][api-docs] for more details.

[jwt]: https://tools.ietf.org/html/rfc7519
[signjwt-method]: https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt
[cloud-creds]: https://cloud.google.com/docs/authentication/production#providing_credentials_to_your_application
[service-accounts]: https://cloud.google.com/compute/docs/access/service-accounts
[api-docs]: /api/auth/gcp/index.html
[identity-group-aliases]: /api/secret/identity/group-alias.html
[instance-identity]: https://cloud.google.com/compute/docs/instances/verifying-instance-identity
[repo]: https://github.com/hashicorp/vault-plugin-auth-gcp
