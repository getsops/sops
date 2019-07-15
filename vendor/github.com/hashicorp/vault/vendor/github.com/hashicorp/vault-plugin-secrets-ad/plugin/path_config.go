package plugin

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/client"
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/util"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configPath       = "config"
	configStorageKey = "config"

	// This length is arbitrarily chosen but should work for
	// most Active Directory minimum and maximum length settings.
	// A bit tongue-in-cheek since programmers love their base-2 exponents.
	defaultPasswordLength = 64

	defaultTLSVersion = "tls12"
)

func (b *backend) readConfig(ctx context.Context, storage logical.Storage) (*configuration, error) {
	entry, err := storage.Get(ctx, configStorageKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	config := &configuration{&passwordConf{}, &client.ADConf{}, 0}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: configPath,
		Fields:  b.configFields(),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.configUpdateOperation,
			logical.ReadOperation:   b.configReadOperation,
			logical.DeleteOperation: b.configDeleteOperation,
		},
		HelpSynopsis:    configHelpSynopsis,
		HelpDescription: configHelpDescription,
	}
}

func (b *backend) configFields() map[string]*framework.FieldSchema {
	fields := ldaputil.ConfigFields()
	fields["ttl"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "In seconds, the default password time-to-live.",
	}
	fields["max_ttl"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "In seconds, the maximum password time-to-live.",
	}
	fields["length"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Default:     defaultPasswordLength,
		Description: "The desired length of passwords that Vault generates.",
	}
	fields["formatter"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `Text to insert the password into, ex. "customPrefix{{PASSWORD}}customSuffix".`,
	}
	fields["last_rotation_tolerance"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "The number of seconds after a Vault rotation where, if Active Directory shows a later rotation, it should be considered out-of-band.",
		Default:     5,
	}
	return fields
}

func (b *backend) configUpdateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	// Build and validate the ldap conf.
	activeDirectoryConf, err := ldaputil.NewConfigEntry(nil, fieldData)
	if err != nil {
		return nil, err
	}
	if err := activeDirectoryConf.Validate(); err != nil {
		return nil, err
	}

	// Build the password conf.
	ttl := fieldData.Get("ttl").(int)
	maxTTL := fieldData.Get("max_ttl").(int)
	length := fieldData.Get("length").(int)
	formatter := fieldData.Get("formatter").(string)
	lastRotationTolerance := fieldData.Get("last_rotation_tolerance").(int)

	if ttl == 0 {
		ttl = int(b.System().DefaultLeaseTTL().Seconds())
	}
	if maxTTL == 0 {
		maxTTL = int(b.System().MaxLeaseTTL().Seconds())
	}
	if ttl > maxTTL {
		return nil, errors.New("ttl must be smaller than or equal to max_ttl")
	}
	if ttl < 1 {
		return nil, errors.New("ttl must be positive")
	}
	if maxTTL < 1 {
		return nil, errors.New("max_ttl must be positive")
	}
	if err := util.ValidatePwdSettings(formatter, length); err != nil {
		return nil, err
	}

	passwordConf := &passwordConf{
		TTL:       ttl,
		MaxTTL:    maxTTL,
		Length:    length,
		Formatter: formatter,
	}

	config := &configuration{passwordConf, &client.ADConf{ConfigEntry: activeDirectoryConf}, lastRotationTolerance}
	entry, err := logical.StorageEntryJSON(configStorageKey, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

func (b *backend) configReadOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	config, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	// NOTE:
	// "password" is intentionally not returned by this endpoint,
	// as we lean away from returning sensitive information unless it's absolutely necessary.
	// Also, we don't return the full ADConf here because not all parameters are used by this engine.
	configMap := map[string]interface{}{
		"url":                     config.ADConf.Url,
		"starttls":                config.ADConf.StartTLS,
		"insecure_tls":            config.ADConf.InsecureTLS,
		"certificate":             config.ADConf.Certificate,
		"binddn":                  config.ADConf.BindDN,
		"userdn":                  config.ADConf.UserDN,
		"upndomain":               config.ADConf.UPNDomain,
		"tls_min_version":         config.ADConf.TLSMinVersion,
		"tls_max_version":         config.ADConf.TLSMaxVersion,
		"last_rotation_tolerance": config.LastRotationTolerance,
	}
	if !config.ADConf.LastBindPasswordRotation.Equal(time.Time{}) {
		configMap["last_bind_password_rotation"] = config.ADConf.LastBindPasswordRotation
	}
	for k, v := range config.PasswordConf.Map() {
		configMap[k] = v
	}

	resp := &logical.Response{
		Data: configMap,
	}
	return resp, nil
}

func (b *backend) configDeleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configStorageKey); err != nil {
		return nil, err
	}
	return nil, nil
}

const (
	configHelpSynopsis = `
Configure the AD server to connect to, along with password options.
`
	configHelpDescription = `
This endpoint allows you to configure the AD server to connect to and its
configuration options. When you add, update, or delete a config, it takes
immediate effect on all subsequent actions. It does not apply itself to roles
or creds added in the past.

The AD URL can use either the "ldap://" or "ldaps://" schema. In the former
case, an unencrypted connection will be made with a default port of 389, unless
the "starttls" parameter is set to true, in which case TLS will be used. In the
latter case, a SSL connection will be established with a default port of 636.

## A NOTE ON ESCAPING

It is up to the administrator to provide properly escaped DNs. This includes
the user DN, bind DN for search, and so on.

The only DN escaping performed by this backend is on usernames given at login
time when they are inserted into the final bind DN, and uses escaping rules
defined in RFC 4514.

Additionally, Active Directory has escaping rules that differ slightly from the
RFC; in particular it requires escaping of '#' regardless of position in the DN
(the RFC only requires it to be escaped when it is the first character), and
'=', which the RFC indicates can be escaped with a backslash, but does not
contain in its set of required escapes. If you are using Active Directory and
these appear in your usernames, please ensure that they are escaped, in
addition to being properly escaped in your configured DNs.

For reference, see https://www.ietf.org/rfc/rfc4514.txt and
http://social.technet.microsoft.com/wiki/contents/articles/5312.active-directory-characters-to-escape.aspx
`
)
