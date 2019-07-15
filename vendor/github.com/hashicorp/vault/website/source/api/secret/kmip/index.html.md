---
layout: "api"
page_title: "KMIP - Secrets Engines - HTTP API"
sidebar_title: "KMIP <sup>ENTERPRISE</sup>"
sidebar_current: "api-http-secret-kmip"
description: |-
  This is the API documentation for the Vault KMIP secrets engine.
---

# KMIP Secrets Engine (API)

This is the API documentation for the Vault KMIP secrets engine. For general 
information about the usage and operation of
the KMIP secrets engine, please see [these docs](/docs/secrets/kmip/index.html).

This documentation assumes the KMIP secrets engine is enabled at the `/kmip` path
in Vault. Since it is possible to mount secrets engines at any path, please
update your API calls accordingly.

## Write Config

| Method | Path           |
|:-------|:---------------|
| `POST` | `/kmip/config` |

This endpoint configures shared information for the secrets engine. After writing
to it the KMIP engine will generate a CA and start listening for KMIP requests.
If the server was already running and any non-client settings are changed, the 
server will be restarted using the new settings.

### Parameters

- `listen_addrs` (`list: ["127.0.0.1:5696"] || string`) - Address and port the 
   KMIP server should listen on. Can be given as a JSON list or a 
   comma-separated string list. If multiple values are given, all will be 
   listened on.
   
- `connection_timeout` (`int: 1 || string:"1s"`) - Duration in either an integer 
   number of seconds (10) or an integer time unit (10s) within which connections
   must become ready.

- `server_hostnames` (`list: ["localhost"] || string`) - Hostnames to include in 
   the server's TLS certificate as SAN DNS names. The first will be used as the 
   common name (CN).

- `server_ips` (`list: [] || string`) - IPs to include in the server's TLS 
   certificate as SAN IP addresses. Localhost (IPv4 and IPv6) will be automatically
   included.
   
- `tls_ca_key_type` (`string: "ec"`) - CA key type, `rsa` or `ec`.

- `tls_ca_key_bits` (`int: 521`) - CA key bits, valid values depend on key type.

- `tls_min_version` (`string: "tls12"`) - Minimum TLS version to accept.

- `default_tls_client_key_type` (`string: "ec"`): - Client certificate key type, 
  `rsa` or `ec`.

- `default_tls_client_key_bits` (`int: 521`): - Client certificate key bits, valid 
  values depend on key type.
  
- `default_tls_client_ttl` (`int: 86400 || string:"24h"`) – Client certificate 
  TTL in either an integer number of seconds (10) or an integer time unit (10s).

### Sample Payload

```json
{
    "listen_addrs":                "127.0.0.1:5696,192.168.1.2:9000",
    "connection_timeout":          "1s",
    "server_hostnames":            "myhostname1,myhostname2",
    "server_ips":                  "192.168.1.2",
    "tls_ca_key_type":             "ec",
    "tls_ca_key_bits":             521,
    "tls_min_version":             "tls11",
    "default_tls_client_key_type": "ec",
    "default_tls_client_key_bits": 224,
    "default_tls_client_ttl":      86400,
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/kmip/config
```

## Read Config

| Method | Path           |
|:-------|:---------------|
| `GET`  | `/kmip/config` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/kmip/config
```

### Sample Response

```json
{
  "data": {
    "listen_addrs":                ["127.0.0.1:5696", "192.168.1.2:9000"],
    "connection_timeout":          "1s",
    "server_hostnames":            ["myhostname1", "myhostname2"],
    "server_ips":                  ["192.168.1.2"],
    "tls_ca_key_type":             "ec",
    "tls_ca_key_bits":             521,
    "tls_min_version":             "tls11",
    "default_tls_client_key_type": "ec",
    "default_tls_client_key_bits": 224,
    "default_tls_client_ttl":      86400,
  }
}
```

## Read CA

| Method | Path       |
|:-------|:-----------|
| `GET`  | `/kmip/ca` |

Returns the CA certificates in PEM format. Returns an error if config has never
been written.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/kmip/ca
```

### Sample Response

```json
{
  "data": {
    "ca_pem": "-----BEGIN CERTIFICATE-----\nMIICNzCCAZigAwIBAgIUApNsRil/dzQy3XT+yjZQEpcA49kwCgYIKoZIzj0EAwIw\nHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MB4XDTE5MDYyNDE4MzIzM1oX\nDTI5MDYyMTE4MzMwM1owKjEoMCYGA1UEAxMfdmF1bHQta21pcC1kZWZhdWx0LWlu\ndGVybWVkaWF0ZTCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEAGWJGwPjGGoXivBv\nLJwR+fIG3z6Ei06bhZgTaRW/U3eA5oivxubxOVZPe1BJGWCsIVNjxMZAN4Pswki7\nAHme9bdJAUbQw33tC1iAb0wjzIpoPv1+pdSk6wYZTCKzOYWCbsTb3SOIetpk7sQw\niM17agwIRK9qGvX3Q4PBfEKEpstAjoaJo2YwZDAOBgNVHQ8BAf8EBAMCAQYwEgYD\nVR0TAQH/BAgwBgEB/wIBCTAdBgNVHQ4EFgQUKMwPpRxU2Uzydv21bc8ePfUpGFEw\nHwYDVR0jBBgwFoAUwrPrJc9EsU6kTWJ5hXkJV4PEq9swCgYIKoZIzj0EAwIDgYwA\nMIGIAkIBRCarRMer42Ni/fKQBTi+uFk+2sPyCxCYDWTfMFAusC51dC2F91mUL77R\nkHxauSkh5gcZVAch/dg/L0ewP0AZUBUCQgE1VqoBN9klFky7LHfl62p6PgprH7d1\nYCvYVbWdBNnEdrL2P9aKsuCewdqycZVJLmM36cHnOAEGg1yea8soQL0Ylw==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIICKTCCAYugAwIBAgIUOBgW1GCH+n5gC6m8Ff5jq+5DmO8wCgYIKoZIzj0EAwIw\nHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MB4XDTE5MDYyNDE4MzIzM1oX\nDTI5MDYyMTE4MzMwM1owHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MIGb\nMBAGByqGSM49AgEGBSuBBAAjA4GGAAQA7vkbmKJR+SVBTJjAFnma0ynTIi64doZA\n5oOXIAExvOyyI2KBNfqXxgzt/51u9vvixQf3VX/1Jph+0fkIcIYUEmIBFAH7Th1X\n0EOOdmMHfN0YkXDEUUdKIZyQxgA7o3DF+JAVg1cdBV7S8jZyXik7pL+IFnlYdfvN\nUZcArUkMfKo1cZajZjBkMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/\nAgEKMB0GA1UdDgQWBBTCs+slz0SxTqRNYnmFeQlXg8Sr2zAfBgNVHSMEGDAWgBTC\ns+slz0SxTqRNYnmFeQlXg8Sr2zAKBggqhkjOPQQDAgOBiwAwgYcCQgGjKAC371/5\npxgYdLVBmVC6Aa+oOvwGfnich2YLSLbThySED7+fXl1BY43VU703ad6M34fStf6z\nwFZvVZVK188DCQJBJcSZ7YA3PjOre+epJHtAba+1CkAdbSAeGhBDgHdIEP1/FDvx\n+U2QYeVZ7kAVnkzPxa17V0yqjxDtQDTiOw/ZV5c=\n-----END CERTIFICATE-----"
  },
}
```

## Write scope

| Method | Path                 |
|:-------|:---------------------|
| `POST` | `/kmip/scope/:scope` |

Creates a new scope with the given name.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://127.0.0.1:8200/v1/kmip/scope/myscope
```

## List scopes

| Method | Path          |
|:-------|:--------------|
| `LIST` | `/kmip/scope` |

List existing scopes.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/kmip/scope
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "myscope"
    ]
  },
}
```

## Delete scope

| Method   | Path                 |
|:---------|:---------------------|
| `DELETE` | `/kmip/scope/:scope` |

Delete a scope by name.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `force` (`bool: false`) - Force scope deletion. If KMIP managed objects have
  been created within the scope this param must be provided or the deletion will
  fail.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/kmip/scope/myscope
```

## Write role

| Method | Path                            |
|:-------|:--------------------------------|
| `POST` | `/kmip/scope/:scope/role/:role` |

Creates or updates a role. 

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.
- `operation_none` (`bool: false`) - Remove all permissions
  from this role. May not be specified with any other 
  `operation_` params.
- `operation_all` (`bool: false`) - Grant all permissions
  to this role. May not be specified with any other 
  `operation_` params.
- `operation_activate` (`bool: false`) - Grant permission to use the KMIP 
  `Activate` operation.
- `operation_add_attribute` (`bool: false`) - Grant permission to use the KMIP 
  `Add Attribute` operation.
- `operation_create` (`bool: false`) - Grant permission to use the KMIP 
  `Create` operation.
- `operation_destroy` (`bool: false`) - Grant permission to use the KMIP 
  `Destroy` operation.
- `operation_discover_versions` (`bool: false`) - Grant permission to use the KMIP 
  `Discover Version` operation.
- `operation_get` (`bool: false`) - Grant permission to use the KMIP 
  `Get` operation.
- `operation_get_attributes` (`bool: false`) - Grant permission to use the KMIP 
  `Get Attributes` operation.
- `operation_locate` (`bool: false`) - Grant permission to use the KMIP 
  `Locate` operation.
- `operation_rekey` (`bool: false`) - Grant permission to use the KMIP 
  `Rekey` operation.
- `operation_revoke` (`bool: false`) - Grant permission to use the KMIP 
  `Revoke` operation.


### Sample Payload

```json
{
  "operation_activate": true,
  "operation_add_attribute": true,
  "operation_create": true,
  "operation_destroy": true,
  "operation_discover_versions": true,
  "operation_get": true,
  "operation_get_attributes": true,
  "operation_locate": true,
  "operation_rekey": true,
  "operation_revoke": true
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole
```

## Read role

| Method | Path                            |
|:-------|:--------------------------------|
| `GET`  | `/kmip/scope/:scope/role/:role` |

Read a role.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole
```

### Sample Response

```json
{
  "data": {
    "operation_activate": true,
    "operation_add_attribute": true,
    "operation_create": true,
    "operation_destroy": true,
    "operation_discover_versions": true,
    "operation_get": true,
    "operation_get_attributes": true,
    "operation_locate": true,
    "operation_rekey": true,
    "operation_revoke": true
  },
}
```

## List roles

| Method | Path                      |
|:-------|:--------------------------|
| `LIST` | `/kmip/scope/:scope/role` |

List roles with a scope.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "myrole"
    ]
  },
}
```

## Delete role

| Method   | Path                            |
|:---------|:--------------------------------|
| `DELETE` | `/kmip/scope/:scope/role/:role` |

Delete a role by name.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole
```


## Generate credential

| Method | Path                                                |
|:-------|:----------------------------------------------------|
| `POST` | `/kmip/scope/:scope/role/:role/credential/generate` |

Create a new client certificate tied to the given role and scope.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.
- `format` (`string: "pem"`) - Format to return the certificate, private key,
  and CA chain in.  One of `pem`, `pem_bundle`, or `der`.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole/credential/generate
```

### Sample Response

```json
{
  "data": {
    "ca_chain": [
      "-----BEGIN CERTIFICATE-----\nMIICNzCCAZigAwIBAgIUKOGtsdXdMjjGni52EsaMQ7ozhCEwCgYIKoZIzj0EAwIw\nHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MB4XDTE5MDYyNDE4NTgyMVoX\nDTI5MDYyMTE4NTg1MVowKjEoMCYGA1UEAxMfdmF1bHQta21pcC1kZWZhdWx0LWlu\ndGVybWVkaWF0ZTCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEATHNhNvU0GMtzl6A\nPbNaCoF0jV3z09RCfLKEqMl/MXv/AlPcfiqCQeOWBwWHv76epPWkCCo+IlNq8ldQ\neVe52p6mABMvRjE6BZ/eLea27zImI6waK7nZ2hqx0npb8ivdbwmrgp0NQnv0sJ+o\nPeLa2vh9wDK1NJebmOv0yRAbCw2CH7Rbo2YwZDAOBgNVHQ8BAf8EBAMCAQYwEgYD\nVR0TAQH/BAgwBgEB/wIBCTAdBgNVHQ4EFgQU2naFRym+xfFvZm2TNRBXNf3MJSsw\nHwYDVR0jBBgwFoAUFrA/R807R0BnIt395KzaXdP4n00wCgYIKoZIzj0EAwIDgYwA\nMIGIAkIAkb8EdHCXgPpQsKYedMz4X2j5CFSVdZTWsPVw1XuSXIsIsc6018V4z9Kp\nkPacsHZTBR636y2toqRPDG4y9MLqFFkCQgCV1jEkiNhhKc+ZWuDjerdqNvLnCbe+\n7t4fiG9zQgWwh6IxL11cNyGVz9gS9af32DtuYf0xwFLOwLgn1RadC9Pd7Q==\n-----END CERTIFICATE-----",
      "-----BEGIN CERTIFICATE-----\nMIICKTCCAYugAwIBAgIUOcs4pXlp+UgGiUKfKlcxIE/woPEwCgYIKoZIzj0EAwIw\nHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MB4XDTE5MDYyNDE4NTgyMVoX\nDTI5MDYyMTE4NTg1MVowHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MIGb\nMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAcst7uNwu77WtLDkbz4ILYDiQ3BgS++qU\nOoNKcKyvNe8YX6PtrdQWPTaxT4MZNHZvTv+BAQTQqGLKrstpkjXPh+sBn7V4trkT\nMCtxUjIGneURUXS4IC/KJEA60P7ep7MrGnJfG/N4m+Q/a6BuxKhdEavXtepniCMz\npHw4DCpW/9m2t16jZjBkMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/\nAgEKMB0GA1UdDgQWBBQWsD9HzTtHQGci3f3krNpd0/ifTTAfBgNVHSMEGDAWgBQW\nsD9HzTtHQGci3f3krNpd0/ifTTAKBggqhkjOPQQDAgOBiwAwgYcCQR7iNoA4nBV3\ndSn8nfafklFvHZxoKR1j3nn+56z4JHD6TNr//GNqQiqnM3P//Tce+E4KzEax4xRg\nhaLURgPLNBjOAkIAqW+1/+v9D0vXOU1WPc+/oFvhSjYnr5qqcTL7by5fsmMXzAIe\nLODXiODxdppXXnMZPCPZh6MGgUwEGYeCnaXopWc=\n-----END CERTIFICATE-----"
    ],
    "certificate": "-----BEGIN CERTIFICATE-----\nMIICOzCCAZygAwIBAgIUeOkn0HAdoh31nGkVKdafpCNuhFEwCgYIKoZIzj0EAwIw\nKjEoMCYGA1UEAxMfdmF1bHQta21pcC1kZWZhdWx0LWludGVybWVkaWF0ZTAeFw0x\nOTA2MjQxOTAwMDlaFw0xOTA2MjUxOTAwMzlaMCAxDjAMBgNVBAsTBWlsVjYzMQ4w\nDAYDVQQDEwUyRnlWTjCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEAA0rIy0h2DL3\nzmTXVj2v22Kz0N1EUUATlRgBj1XBsBA1Pdd7CSZoefmh/u6Z8TjtRX9Z1aj9Bb/d\nJxS3zB4mguULAF4k7bLH1gKXMVC6NYjjk3mfxH5jG4QY8S8n6uyqzNgI5KRJ2Hyj\nm8549Nvq3rvs8yOVXPSOGzkJ5KdUmSvXicMQo2cwZTAOBgNVHQ8BAf8EBAMCA6gw\nEwYDVR0lBAwwCgYIKwYBBQUHAwIwHQYDVR0OBBYEFEuzruLILCil5Fp32ZjE4AhD\nU268MB8GA1UdIwQYMBaAFNp2hUcpvsXxb2ZtkzUQVzX9zCUrMAoGCCqGSM49BAMC\nA4GMADCBiAJCAeeuaIsgO9ro7opzZ9y9hSHkKB5WA5Qc7ePoSiKHNNbVvIJMkjRQ\nC9YtUMQNnQ8wE6D/9xvR+9OBIi7t16iHGPGbAkIA6WIG6HHRNUXnHPIiW8iy/04O\nfVqZgJHJEeyGQbwdaehs+Z5xOz6TA4Z3uZOAMnPcb+KDwchnQ8CJnmT/KnnT5D8=\n-----END CERTIFICATE-----",
    "private_key": "-----BEGIN EC PRIVATE KEY-----\nMIHcAgEBBEIBB4xDj9SUtb6Z466lVQIf3ucy21q5S2Fp9bzTQ0Ch5Vg2+DhUZUa1\nDjKvDdICY6hLPBFAwcOUFdDXr4kH/i8wuRWgBwYFK4EEACOhgYkDgYYABAANKyMt\nIdgy985k11Y9r9tis9DdRFFAE5UYAY9VwbAQNT3XewkmaHn5of7umfE47UV/WdWo\n/QW/3ScUt8weJoLlCwBeJO2yx9YClzFQujWI45N5n8R+YxuEGPEvJ+rsqszYCOSk\nSdh8o5vOePTb6t677PMjlVz0jhs5CeSnVJkr14nDEA==\n-----END EC PRIVATE KEY-----",
    "serial_number": "728181095563584845125173905844944137943705466376"
  },
}
```

## Lookup credential

| Method | Path                                              |
|:-------|:--------------------------------------------------|
| `GET`  | `/kmip/scope/:scope/role/:role/credential/lookup` |

Read a certificate by serial number. The private key cannot be obtained except
at generation time.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.
- `serial_number` (`string: <required>`) - Serial number of certificate to revoke.
- `format` (`string: "pem"`) - Format to return the certificate, private key,
  and CA chain in.  One of `pem`, `pem_bundle`, or `der`.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole/credential/lookup?serial_number=728181095563584845125173905844944137943705466376
```

### Sample Response

```json
{
  "data": {
    "ca_chain": [
      "-----BEGIN CERTIFICATE-----\nMIICNzCCAZigAwIBAgIUGptwpwpVvxlx3sBniJ7TRGD9gCkwCgYIKoZIzj0EAwIw\nHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MB4XDTE5MDYyNDE5MDY0N1oX\nDTI5MDYyMTE5MDcxN1owKjEoMCYGA1UEAxMfdmF1bHQta21pcC1kZWZhdWx0LWlu\ndGVybWVkaWF0ZTCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEADO48mMu5V2PTbcg\nq0JPB5ReWwnUHhfFh/+XLP8ZM112JpOFutlcUYYZ23jAlvrlYZ+m1E0ASr0592ZM\n9CwIXy3zAJChPrV3tiofhINR5PPqCF42FcfNj4l7VN/XeYMN6dslX+O4dPn/DsbH\nZi7kWr5KSOR939ULFaRMYe3l2MxaYZ2do2YwZDAOBgNVHQ8BAf8EBAMCAQYwEgYD\nVR0TAQH/BAgwBgEB/wIBCTAdBgNVHQ4EFgQUPP7VJOGk3qR0qKqx3TLN1R8JDiQw\nHwYDVR0jBBgwFoAUBHr+hhaorPU2jIF35DTBDhL7uWowCgYIKoZIzj0EAwIDgYwA\nMIGIAkIA7G82rqLYb6bKrQZzhpNwvVIFOSocEJrUbP0E0D8dEeOmKs43C70P5e0s\nTrrpNAMEsK6vXWtM+QcrZZp+yyM6k3QCQgG8cxFIl8tgoMKWe0+cDeOoHtczopRy\nSk+Tt7DNNP9sfYK11g7w8xzbtW4ZuZKKoYRbxN+eQHn5c+8akMSt4h71Dg==\n-----END CERTIFICATE-----",
      "-----BEGIN CERTIFICATE-----\nMIICKDCCAYugAwIBAgIUWv6jrjNbsvdX43l4s10HaJkSxOMwCgYIKoZIzj0EAwIw\nHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MB4XDTE5MDYyNDE5MDY0N1oX\nDTI5MDYyMTE5MDcxN1owHTEbMBkGA1UEAxMSdmF1bHQta21pcC1kZWZhdWx0MIGb\nMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAP6C8d9ZUalKBM1NdALtEMlv+dwFnK88F\n8bp7i6hV55vER45FtKKciQwWoA91FjfWTrDYPHb1X4OPZvcjQGnIJ1AAj+BSzEWr\neJXNo46RxLLl+cndiVDqlbJlhE9qVn9ueLHhPIPNSFZneY9cTj5+EOPyKiBCo4xB\ndTtVr29lLu/JwM2jZjBkMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/\nAgEKMB0GA1UdDgQWBBQEev6GFqis9TaMgXfkNMEOEvu5ajAfBgNVHSMEGDAWgBQE\nev6GFqis9TaMgXfkNMEOEvu5ajAKBggqhkjOPQQDAgOBigAwgYYCQUlJqNoWCz4H\npjMNphxD4A8lfWtIrajGUhSxE9+JWRzoPpEJSwVobvryU2SO5u0sfqxtcmX/sBjY\n12N5QVFfqpB3AkErsjg8eMkh+OMalmWxRYtTuZt+i4DPm1CKEVIkUT8ZBXYTIl9V\nG3TG8lmby/8e+YUwJEKVvOy6tVI8ExEoVslwKw==\n-----END CERTIFICATE-----"
    ],
    "certificate": "-----BEGIN CERTIFICATE-----\nMIICOjCCAZygAwIBAgIUf4zFBobFJMkSIvM7CfceSVfYNggwCgYIKoZIzj0EAwIw\nKjEoMCYGA1UEAxMfdmF1bHQta21pcC1kZWZhdWx0LWludGVybWVkaWF0ZTAeFw0x\nOTA2MjQxOTA3MTBaFw0xOTA2MjUxOTA3NDBaMCAxDjAMBgNVBAsTBW5BcUswMQ4w\nDAYDVQQDEwU0Qjd2STCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEAdxHrbr/EXUz\nzWCd9HMUDus6r/3QF1Y3u9dPD2UwM76J3aICmykkm7xoYpoyg4chBEDxBWh2YkGT\na4WFMoXBa+k1AZhdvlj8tjOUlYZrTCLB9FBPCGz3JB4f5cmbG5JVsQ8qnBPiyV3e\nU21cWM6mWlhZKHWIdBU2pj+eXW78K5LMu2sWo2cwZTAOBgNVHQ8BAf8EBAMCA6gw\nEwYDVR0lBAwwCgYIKwYBBQUHAwIwHQYDVR0OBBYEFAT0QZOpZCTMCz7F8+BvF2xs\nZSfkMB8GA1UdIwQYMBaAFDz+1SThpN6kdKiqsd0yzdUfCQ4kMAoGCCqGSM49BAMC\nA4GLADCBhwJBPxBV4DgPi5zihRnxu7zTNeqe/xlvrEt1uTff8QtW3JsigbBDHV+A\nxBe7vc8mL8VQPG7BFKvvxuQvOAeeQ+AR8ZoCQgDtbaWgLtfbzKvwlY48e6dLeBpK\nDu1DaZq+79EON2lhWQ+ULHblJc5cK0F6Ff5OC89aDnV1TWQDHeR91mZdYiWZZQ==\n-----END CERTIFICATE-----",
    "serial_number": "728181095563584845125173905844944137943705466376"
  },
}
```

## List credential serial numbers

| Method | Path                                       |
|:-------|:-------------------------------------------|
| `LIST` | `/kmip/scope/:scope/role/:role/credential` |

List the serial numbers of all certificates within a role.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole/credential
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "728181095563584845125173905844944137943705466376"
    ]
  },
}
```

## Revoke credential

| Method | Path                                              |
|:-------|:--------------------------------------------------|
| `POST` | `/kmip/scope/:scope/role/:role/credential/revoke` |

Delete a certificate, thereby revoking it.

### Parameters

- `scope` (`string: <required>`) - Name of scope. This is part of the request URL.
- `role` (`string: <required>`) - Name of role. This is part of the request URL.
- `serial_number` (`string: ""`) - Serial number of certificate to revoke.
  Exactly one of `serial_number` or `certificate` must be provided.
- `certificate` (`string: """`) - Certificate to revoke, in PEM format.
  Exactly one of `serial_number` or `certificate` must be provided.

### Sample Payload

```json
{
    "serial_number": "728181095563584845125173905844944137943705466376"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://127.0.0.1:8200/v1/kmip/scope/myscope/role/myrole/credential/revoke
```
