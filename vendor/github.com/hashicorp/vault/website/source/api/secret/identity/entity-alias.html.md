---
layout: "api"
page_title: "Identity Secret Backend: Entity Alias - HTTP API"
sidebar_title: "Entity Alias"
sidebar_current: "api-http-secret-identity-entity-alias"
description: |-
  This is the API documentation for managing entity aliases in the identity store.
---

## Create an Entity Alias

This endpoint creates a new alias for an entity.

| Method   | Path                       |
| :------------------------- | :----------------------|
| `POST`   | `/identity/entity-alias`   |

### Parameters

- `name` `(string: <required>)` - Name of the alias. Name should be the identifier
  of the client in the authentication source. For example, if the alias belongs
  to userpass backend, the name should be a valid username within userpass
  backend. If alias belongs to GitHub, it should be the GitHub username.

- `id` `(string: <optional>)` - ID of the entity alias. If set, updates the
  corresponding entity alias.

- `canonical_id` `(string: <required>)` - Entity ID to which this alias belongs to.

- `mount_accessor` `(string: <required>)` - Accessor of the mount to which the
  alias should belong to.

### Sample Payload

```json
{
  "name": "testuser",
  "canonical_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
  "mount_accessor": "auth_userpass_e50b1a44"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/entity-alias
```

### Sample Response

```json
{
  "data": {
    "canonical_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
    "id": "34982d3d-e3ce-5d8b-6e5f-b9bb34246c31"
  }
}
```

## Read Entity Alias by ID

This endpoint queries the entity alias by its identifier.

| Method   | Path                             |
| :------------------------------- | :--------------------- |
| `GET`    | `/identity/entity-alias/id/:id`  |

### Parameters

- `id` `(string: <required>)` – Identifier of entity alias.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/identity/entity-alias/id/34982d3d-e3ce-5d8b-6e5f-b9bb34246c31
```

### Sample Response

```json
{
  "data": {
    "creation_time": "2017-07-25T21:41:09.820717636Z",
    "canonical_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
    "id": "34982d3d-e3ce-5d8b-6e5f-b9bb34246c31",
    "last_update_time": "2017-07-25T21:41:09.820717636Z",
    "metadata": {
      "group": "san_francisco",
      "region": "west"
    },
    "mount_accessor": "auth_userpass_e50b1a44",
    "mount_path": "userpass/",
    "mount_type": "userpass",
    "name": "testuser"
  }
}
```

## Update Entity Alias by ID

This endpoint is used to update an existing entity alias.

| Method   | Path                              |
| :-------------------------------- | :--------------------- |
| `POST`    | `/identity/entity-alias/id/:id`  |

### Parameters

- `id` `(string: <required>)` – Identifier of the entity alias.

- `name` `(string: <required>)` - Name of the alias. Name should be the identifier
  of the client in the authentication source. For example, if the alias belongs
  to userpass backend, the name should be a valid username within userpass
  backend. If alias belongs to GitHub, it should be the GitHub username.

- `canonical_id` `(string: <required>)` - Entity ID to which this alias belongs to.

- `mount_accessor` `(string: <required>)` - Accessor of the mount to which the
  alias should belong to.

### Sample Payload

```json
{
  "name": "testuser",
  "canonical_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
  "mount_accessor": "auth_userpass_e50b1a44"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/identity/entity-alias/id/34982d3d-e3ce-5d8b-6e5f-b9bb34246c31
```

### Sample Response

```json
{
  "data": {
    "canonical_id": "404e57bc-a0b1-a80f-0a73-b6e92e8a52d3",
    "id": "34982d3d-e3ce-5d8b-6e5f-b9bb34246c31"
  }
}
```

### Delete Entity Alias by ID

This endpoint deletes an alias from its corresponding entity.

| Method     | Path                             |
| :------------------------------- | :----------------------|
| `DELETE`   | `/identity/entity-alias/id/:id`  |

### Parameters

- `id` `(string: <required>)` – Identifier of the entity alias.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/identity/entity-alias/id/34982d3d-e3ce-5d8b-6e5f-b9bb34246c31
```

### List Entity Aliases by ID

This endpoint returns a list of available entity aliases by their identifiers.

| Method   | Path                                  |
| :------------------------------------ | :--------------------- |
| `LIST`   | `/identity/entity-alias/id`           |
| `GET`    | `/identity/entity-alias/id?list=true` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/identity/entity-alias/id
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "2e8217fa-8cb6-8aec-9e22-3196d74ca2ba",
      "91ebe973-ec86-84db-3c7c-f760415326de",
      "92308b08-4139-3ec6-7af2-8e98166b4e0c",
      "a3b042e6-5cc1-d5a9-8874-d53a51954de2",
      "d5844921-017f-e496-2a9a-23d4a2f3e8a3"
    ]
  }
}
```

