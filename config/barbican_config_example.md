# Barbican Configuration Examples

This document provides examples of how to configure OpenStack Barbican support in SOPS configuration files.

## Basic Configuration

### Single Barbican Key

```yaml
# .sops.yaml
creation_rules:
  - path_regex: \.prod\.yaml$
    barbican: "550e8400-e29b-41d4-a716-446655440000"
    barbican_auth_url: "https://keystone.example.com:5000/v3"
    barbican_region: "us-east-1"
```

### Multiple Barbican Keys

```yaml
# .sops.yaml
creation_rules:
  - path_regex: \.prod\.yaml$
    barbican:
      - "550e8400-e29b-41d4-a716-446655440000"
      - "region:us-west-1:660e8400-e29b-41d4-a716-446655440001"
    barbican_auth_url: "https://keystone.example.com:5000/v3"
    barbican_region: "us-east-1"
```

## Advanced Configuration

### Mixed Key Types with Key Groups

```yaml
# .sops.yaml
creation_rules:
  - path_regex: ""
    key_groups:
    - barbican:
      - secret_ref: "550e8400-e29b-41d4-a716-446655440000"
        region: "us-east-1"
      - secret_ref: "660e8400-e29b-41d4-a716-446655440001"
        region: "us-west-1"
      kms:
      - arn: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
      pgp:
      - "85D77543B3D624B63CEA9E6DBC17301B491B3F21"
```

### Multiple Rules for Different Environments

```yaml
# .sops.yaml
creation_rules:
  - path_regex: \.prod\.yaml$
    barbican:
      - "550e8400-e29b-41d4-a716-446655440000"
      - "region:us-west-1:660e8400-e29b-41d4-a716-446655440001"
    barbican_auth_url: "https://keystone.example.com:5000/v3"
    barbican_region: "us-east-1"
  - path_regex: \.dev\.yaml$
    barbican: "770e8400-e29b-41d4-a716-446655440002"
    barbican_auth_url: "https://keystone-dev.example.com:5000/v3"
    barbican_region: "us-west-2"
  - path_regex: ""
    kms: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
```

## Configuration Options

### Barbican-Specific Settings

- `barbican`: Barbican secret reference(s). Can be:
  - A single UUID: `"550e8400-e29b-41d4-a716-446655440000"`
  - A regional reference: `"region:us-east-1:550e8400-e29b-41d4-a716-446655440000"`
  - A full URI: `"https://barbican.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000"`
  - A comma-separated list: `"uuid1,uuid2,uuid3"`
  - An array of references: `["uuid1", "uuid2", "uuid3"]`

- `barbican_auth_url`: OpenStack Keystone authentication URL (e.g., `"https://keystone.example.com:5000/v3"`)

- `barbican_region`: Default OpenStack region for Barbican keys (e.g., `"us-east-1"`)

### Key Groups Format

When using key groups, Barbican keys can be specified with additional metadata:

```yaml
barbican:
  - secret_ref: "550e8400-e29b-41d4-a716-446655440000"
    region: "us-east-1"
  - secret_ref: "660e8400-e29b-41d4-a716-446655440001"
    region: "us-west-1"
```

## Validation

The configuration system validates:

1. **Secret Reference Format**: Must be a valid UUID, regional reference, or full URI
2. **Auth URL Format**: Must be a valid HTTP or HTTPS URL
3. **Region Format**: Cannot be empty or whitespace-only
4. **Key Accessibility**: References must point to accessible Barbican secrets (when possible)

## Error Messages

Common validation errors:

- `invalid Barbican secret reference 'invalid-ref'`: The secret reference format is invalid
- `barbican_auth_url must be a valid HTTP or HTTPS URL`: The auth URL is malformed
- `barbican_region cannot be empty or whitespace`: The region field is empty or contains only whitespace
- `no valid authentication method provided`: Missing authentication configuration in environment variables