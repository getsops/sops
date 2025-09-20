# AWS Secrets Manager Key/Value Format Guide

## Overview

AWS Secrets Manager supports two different formats for storing secrets:

1. **JSON Format** - Complex nested structures (default SOPS behavior)
2. **Key/Value Format** - Simple flat key/value pairs (enables console key/value editor)

## The Difference

### JSON Format (Complex Structure)

**SOPS File:**
```yaml
database:
  host: db.example.com
  port: 5432
  username: app_user
  password: secret123

api_keys:
  stripe: sk_test_123
  github: ghp_456
```

**Stored in AWS Secrets Manager as:**
```json
{
  "database": {
    "host": "db.example.com",
    "port": 5432,
    "username": "app_user", 
    "password": "secret123"
  },
  "api_keys": {
    "stripe": "sk_test_123",
    "github": "ghp_456"
  }
}
```

**AWS Console Experience:**
- ❌ Key/Value tab is **disabled**
- ⚠️ Only "Plaintext" tab available
- 📝 Must edit raw JSON manually
- 🔍 Harder to find specific values

### Key/Value Format (Flat Structure)

**SOPS File:**
```yaml
database_host: db.example.com
database_port: "5432"
database_username: app_user
database_password: secret123
api_key_stripe: sk_test_123
api_key_github: ghp_456
```

**Stored in AWS Secrets Manager as:**
```json
{
  "database_host": "db.example.com",
  "database_port": "5432", 
  "database_username": "app_user",
  "database_password": "secret123",
  "api_key_stripe": "sk_test_123",
  "api_key_github": "ghp_456"
}
```

**AWS Console Experience:**
- ✅ Key/Value tab is **enabled**
- 🎯 Easy individual key editing
- 🔍 Quick search and filtering
- 📋 Copy individual values easily
- 👥 Better for non-technical team members

## When to Use Each Format

### Use JSON Format When:
- You have complex hierarchical configuration
- You need nested objects or arrays
- You're migrating existing complex configs
- You primarily access secrets programmatically

### Use Key/Value Format When:
- You want easy AWS console management
- Team members need to update individual secrets
- You have simple configuration values
- You want better visibility and searchability

## Implementation

### Testing Both Formats

The integration tests demonstrate both approaches:

```bash
# Run integration tests (requires AWS credentials)
export SOPS_TEST_AWS_SECRET_NAME="sops-test-secret"
export SOPS_TEST_AWS_REGION="us-east-1"

go test -tags=integration ./publish -run TestAWSSecretsManagerDestination -v
```

**Test Coverage:**
- `TestAWSSecretsManagerDestination_Integration` - Complex JSON format
- `TestAWSSecretsManagerDestination_KeyValueFormat_Integration` - Simple key/value format

### Example Files

See the examples directory:
- `examples/example.yaml` - Complex nested structure (JSON format)
- `examples/keyvalue-secrets.yaml` - Flat key/value structure

## Best Practices

### For Key/Value Format:

1. **Use string values consistently:**
   ```yaml
   # Good - all strings
   port: "5432"
   debug: "true"
   timeout: "30"
   
   # Avoid - mixed types
   port: 5432      # number
   debug: true     # boolean
   timeout: "30"   # string
   ```

2. **Use descriptive flat keys:**
   ```yaml
   # Good - clear hierarchy in key names
   database_host: db.example.com
   database_port: "5432"
   api_key_stripe: sk_test_123
   
   # Avoid - unclear relationships
   host: db.example.com
   port: "5432"
   key1: sk_test_123
   ```

3. **Group related keys with prefixes:**
   ```yaml
   # Database settings
   db_host: localhost
   db_port: "5432"
   db_name: myapp
   
   # API keys
   api_stripe: sk_test_123
   api_github: ghp_456
   
   # Feature flags
   feature_new_ui: "true"
   feature_beta: "false"
   ```

### For JSON Format:

1. **Use when you need complex structures:**
   ```yaml
   database:
     primary:
       host: db1.example.com
       port: 5432
     replica:
       host: db2.example.com
       port: 5432
   
   api_endpoints:
     - name: stripe
       url: https://api.stripe.com
       timeout: 30
     - name: github
       url: https://api.github.com
       timeout: 10
   ```

## Migration Between Formats

### From JSON to Key/Value:

```bash
# Original complex structure
database:
  host: db.example.com
  port: 5432

# Flatten to key/value
database_host: db.example.com
database_port: "5432"
```

### From Key/Value to JSON:

```bash
# Original flat structure  
database_host: db.example.com
database_port: "5432"

# Group into nested structure
database:
  host: db.example.com
  port: 5432
```

## Console Screenshots Comparison

### JSON Format Console View:
```
┌─ AWS Secrets Manager Console ─────────────────────┐
│ Secret: myapp/config                              │
│                                                   │
│ Tabs: [Plaintext] Key/value(disabled)            │
│                                                   │
│ {                                                 │
│   "database": {                                   │
│     "host": "db.example.com",                     │
│     "port": 5432,                                 │
│     "password": "secret123"                       │
│   },                                              │
│   "api_keys": {                                   │
│     "stripe": "sk_test_123"                       │
│   }                                               │
│ }                                                 │
└───────────────────────────────────────────────────┘
```

### Key/Value Format Console View:
```
┌─ AWS Secrets Manager Console ─────────────────────┐
│ Secret: myapp/config-keyvalue                     │
│                                                   │
│ Tabs: Plaintext [Key/value]                      │
│                                                   │
│ ┌─ Key ──────────────┬─ Value ──────────────────┐ │
│ │ database_host      │ db.example.com           │ │
│ │ database_port      │ 5432                     │ │
│ │ database_password  │ ••••••••••               │ │
│ │ api_key_stripe     │ sk_test_123              │ │
│ └────────────────────┴──────────────────────────┘ │
│                                                   │
│ [Add] [Edit] [Delete] buttons for each key       │
└───────────────────────────────────────────────────┘
```

## Conclusion

Choose the format that best fits your team's workflow:

- **Key/Value Format**: Better for teams that need AWS console access
- **JSON Format**: Better for complex configurations accessed programmatically

Both formats work with SOPS encryption and publishing - the choice depends on your operational needs.
