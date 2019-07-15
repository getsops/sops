---
layout: "docs"
page_title: "Identity - Secrets Engines"
sidebar_title: "Identity"
sidebar_current: "docs-secrets-identity"
description: |-
  The Identity secrets engine for Vault manages client identities.
---

# Identity Secrets Engine

Name: `identity`

The Identity secrets engine is the identity management solution for Vault. It
internally maintains the clients who are recognized by Vault. Each client is
internally termed as an `Entity`. An entity can have multiple `Aliases`. For
example, a single user who has accounts in both GitHub and LDAP, can be mapped
to a single entity in Vault that has 2 aliases, one of type GitHub and one of
type LDAP. When a client authenticates via any of the credential backend
(except the Token backend), Vault creates a new entity and attaches a new
alias to it, if a corresponding entity doesn't already exist. The entity identifier will
be tied to the authenticated token. When such tokens are put to use, their
entity identifiers are audit logged, marking a trail of actions performed by
specific users.

Identity store allows operators to **manage** the entities in Vault. Entities
can be created and aliases can be tied to entities, via the ACL'd API. There
can be policies set on the entities which adds capabilities to the tokens that
are tied to entity identifiers. The capabilities granted to tokens via the
entities are **an addition** to the existing capabilities of the token and
**not** a replacement. The capabilities of the token that get inherited from
entities are computed dynamically at request time. This provides flexibility in
controlling the access of tokens that are already issued.

This secrets engine will be mounted by default. This secrets engine cannot be
disabled or moved.

## Concepts

### Entities and Aliases

Each user will have multiple accounts with various identity providers. Users
can now be mapped as `Entities` and their corresponding accounts with
authentication providers can be mapped as `Aliases`. In essence, each entity is
made up of zero or more aliases.

### Entity Management

Entities in Vault **do not** automatically pull identity information from
anywhere. It needs to be explicitly managed by operators. This way, it is
flexible in terms of administratively controlling the number of entities to be
synced against Vault. In some sense, Vault will serve as a _cache_ of
identities and not as a _source_ of identities.

### Entity Policies

Vault policies can be assigned to entities which will grant _additional_
permissions to the token on top of the existing policies on the token. If the
token presented on the API request contains an identifier for the entity and if
that entity has a set of policies on it, then the token will be capable of
performing actions allowed by the policies on the entity as well.

This is a paradigm shift in terms of _when_ the policies of the token get
evaluated. Before identity, the policy names on the token were immutable (not
the contents of those policies though). But with entity policies, along with
the immutable set of policy names on the token, the evaluation of policies
applicable to the token through its identity will happen at request time. This
also adds enormous flexibility to control the behavior of already issued
tokens.

Its important to note that the policies on the entity are only a means to grant
_additional_ capabilities and not a replacement for the policies on the token.
To know the full set of capabilities of the token with an associated entity
identifier, the policies on the token should be taken into account.

### Mount Bound Aliases

Vault supports multiple authentication backends and also allows enabling the
same type of authentication backend on different mount paths. The alias name of
the user will be unique within the backend's mount. But identity store needs to
uniquely distinguish between conflicting alias names across different mounts of
these identity providers. Hence, the alias name in combination with the
authentication backend mount's accessor, serve as the unique identifier of an
alias.

### Implicit Entities

Operators can create entities for all the users of an auth mount beforehand and
assign policies to them, so that when users login, the desired capabilities to
the tokens via entities are already assigned. But if that's not done, upon a
successful user login from any of the authentication backends, Vault will
create a new entity and assign an alias against the login that was successful.

Note that the tokens created using the token authentication backend will not
have an associated identity information. Logging in using the authentication
backends is the only way to create tokens that have a valid entity identifiers.

### Identity Auditing

If the token used to make API calls have an associated entity identifier, it
will be audit logged as well. This leaves a trail of actions performed by
specific users.

### Identity Groups

In version 0.9, Vault identity has support for groups. A group can contain
multiple entities as its members. A group can also have subgroups. Policies set
on the group is granted to all members of the group. During request time, when
the token's entity ID is being evaluated for the policies that it has access
to; along with the policies on the entity itself, policies that are inherited
due to group memberships are also granted.

### Group Hierarchical Permissions

Entities can be direct members of groups, in which case they inherit the
policies of the groups they belong to. Entities can also be indirect members of
groups. For example, if a GroupA has GroupB as subgroup, then members of GroupB
are indirect members of GroupA. Hence, the members of GroupB will have access
to policies on both GroupA and GroupB.

### External vs Internal Groups

By default, the groups created in identity store are called the internal
groups. The membership management of these groups should be carried out
manually. A group can also be created as an external group. In this case, the
entity membership in the group is managed semi-automatically. External group
serves as a mapping to a group that is outside of the identity store. External
groups can have one (and only one) alias. This alias should map to a notion of
group that is outside of the identity store. For example, groups in LDAP, and
teams in GitHub. A username in LDAP, belonging to a group in LDAP, can get its
entity ID added as a member of a group in Vault automatically during *logins*
and *token renewals*. This works only if the group in Vault is an external
group and has an alias that maps to the group in LDAP. If the user is removed
from the group in LDAP, that change gets reflected in Vault only upon the
subsequent login or renewal operation.

## Identity Tokens

Identity information is used throughout Vault, but it can also be exported for
use by other applications. An authorized user/application can request a token
that encapsulates identity information for their associated entity. These
tokens are signed JWTs following the [OIDC ID
token](https://openid.net/specs/openid-connect-core-1_0.html#IDToken) structure.
The public keys used to authenticate the tokens are published by Vault on an
unauthenticated endpoint following OIDC discovery and JWKS conventions, which
should be a directly usable by JWT/OIDC libraries. An introspection endpoint is
also provided by Vault for token verification.

### Roles and Keys

OIDC-compliant ID tokens are generated against a role which allows configuration
of token claims via a templating system, token ttl, and a way to specify which
"key" will be used to sign the token. The role template is an optional parameter
to customize the token contents and is described in the next section. Token TTL
controls the expiration time of the token, after which verification library will
consider the token invalid. All roles have a Vault-generated `client_id`
attribute that is returned when the role is read. This value cannot be changed
and will be added to the token's `aud` parameter. JWT/OIDC libraries will often
require this value.

A role's `key` parameter links a role to an existing named key (multiple roles
may refer to the same key). It is not possible to generate an unsigned ID token.

A named key is a public/private key pair generated by Vault. The private key is
used to sign the identity tokens, and the public key is used by clients to
verify the signature. Key are regularly rotated, whereby a new key pair is
generated and the previous _public_ key is retained for a limited time for
verification purposes.

A named key's configuration specifies a rotation period, a verification ttl, and
signing algorithm. Rotation period specifies the frequency at which a new
signing key is generated and the private portion of the previous signing key is
deleted. Verification ttl is the time a public key is retained for verification,
after being rotated. By default, keys are rotated every 24 hours, and continue
to be available for verification for 24 hours after their rotation.


### Token Contents and Templates

Identity tokens will always contain, at a minimum, the claims required by OIDC:

* `iss` - Issuer URL
* `sub` - Requester's entity ID
* `aud` - `client_id` for the role
* `iss` - Time of issue
* `exp` - Expiration time for the token

In addition, the operator may configure per-role templates that allow a variety
of other entity information to be added to the token. The templates are
structured as JSON with replaceable parameters. The parameter syntax is the same
as that used for [ACL Path Templating](/docs/concepts/policies.html).

For example:

```json
{
  "color": {{identity.entity.metadata.color}},
  "userinfo": {
     "username": {{identity.entity.aliases.usermap_123.metadata.username}},
     "groups": {{identity.entity.group_names}}
  
  "nbf": {{time.now}},
}
```

When a token is requested, the resulting template might be populated as:

```json
{
  "color": "green",
  "userinfo": {
     "username": "bob",
     "groups": ["web", "engr", "default]
  
  "nbf": 1561411915,
}
```

which would be merged with the base OIDC claims into the final token:

```json
{
  "iss": "https://10.1.1.45:8200/v1/identity/oidc",
  "sub": "a2cd63d3-5364-406f-980e-8d71bb0692f5",
  "aud": "SxSouteCYPBoaTFy94hFghmekos",
  "iss": 1561411915,
  "exp": 1561412215,
  "color": "green",
  "userinfo": {
     "username": "bob",
     "groups": ["web", "engr", "default]
  },
  "nbf": 1561411915,
}
```

Note how the template is merged, with top level template keys becoming top level
token keys. For this reason, templates may not contain top level keys that
overwrite the standard OIDC claims.

Template parameters that are not present for an entity, such as a metadata that
isn't present, or an alias accessor which doesn't exist, are simply empty
strings or objects, depending on the data type.

Templates are configured on the role and may be optionally encoded as base64.

The full list of template parameters is shown below:

|                                    Name                                |                                    Description                                          |
| :--------------------------------------------------------------------- | :-------------------------------------------------------------------------------------- |
| `identity.entity.id`                                                   | The entity's ID                                                                         |
| `identity.entity.name`                                                 | The entity's name                                                                       |
| `identity.entity.group_ids`                                            | The IDs of the groups the entity is a member of                                         |
| `identity.entity.group_names`                                          | The names of the groups the entity is a member of                                       |
| `identity.entity.metadata`                                             | Metadata associated with the entity                                                     |
| `identity.entity.metadata.<<metadata key>>`                            | Metadata associated with the entity for the given key                                   |
| `identity.entity.aliases.<<mount accessor>>.id`                        | Entity alias ID for the given mount                                                     |
| `identity.entity.aliases.<<mount accessor>>.name`                      | Entity alias name for the given mount                                                   |
| `identity.entity.aliases.<<mount accessor>>.metadata`                  | Metadata associated with the alias for the given mount                                  |
| `identity.entity.aliases.<<mount accessor>>.metadata.<<metadata key>>` | Metadata associated with the alias for the given mount and metadata key                 |
| `time.now`                                                             | Current time as integral seconds since the Epoch                                        |
| `time.now.plus.<<duration>>`                                           | Current time plus a Go-parsable [duration](https://golang.org/pkg/time/#ParseDuration)  |                             |
| `time.now.minus.<<duration>>`                                          | Current time minus a Go-parsable [duration](https://golang.org/pkg/time/#ParseDuration) |                             |

### Token Generation
    
An authenticated client may request a token using the [token generation
endpoint](api/secret/identity/tokens.html#generate-a-signed-id-token). The token
will be generated per the requested role's specifications, for the requester's
entity. It is not possible to generate tokens for an arbitrary entity.

### Verifying Authenticity of ID Tokens Generated by Vault

An identity token may be verified by the client party using the public keys
published by Vault, or via a Vault-provided introspection endpoint.

Vault will serve standard "[.well-known](https://tools.ietf.org/html/rfc5785)"
endpoints that allow easy integration with OIDC verification libraries.
Configuring the libraries will typically involve providing an issuer URL and
client ID. The library will then handle key requests and can validate the
signature and claims requirements on tokens. This approach has the advantage of
only requiring _access_ to Vault, not _authorization_, as the .well-known
endpoints are unauthenticated.

Alternatively, the token may be sent to Vault for verification via an
[introspection endpoint](api/secret/identity/tokens.html#introspect-a-signed-id-token).
The response will indicate whether the token is "active" or not, as well as any
errors that occurred during validation. Beyond simply allowing the client to
delegate verification to Vault, using this endpoint incorporates the additional
check of whether the entity is still active or not, which is something that
cannot be determined from the token alone. Unlike the .well-known endpoint, accessing the
introspection endpoint does require a valid Vault token and sufficient
authorization.


### Issuer Considerations

The identity token system has one configurable parameter: issuer. The issuer
`iss` claim is particularly important for proper validation of the token by
clients, and special consideration should be given when using Identity Tokens
with [performance replication](docs/enterprise/replication/index.html).
Consumers of the token will request public keys from Vault using the issuer URL,
so it must be network reachable. Furthermore, the returned set of keys will include
an issuer that must match the request.

By default Vault will set the issuer to the Vault instance's
[`api_addr`](docs/configuration/index.html#api_addr). This means that tokens
issued in a given cluster should be validated within that same cluster.
Alternatively, the [`issuer`](api/secret/identity/tokens.html#issuer) parameter
may be configured explicitly. This address must point to the identity/oidc path
for the Vault instance (e.g.
`https://vault-1.example.com:8200/v1/identity/oidc`) and should be
reachable by any client trying to validate identity tokens.


## API

The Identity secrets engine has a full HTTP API. Please see the
[Identity secrets engine API](/api/secret/identity/index.html) for more
details.
