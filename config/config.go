/*
Package config provides a way to find and load SOPS configuration files
*/
package config //import "github.com/getsops/sops/v3/config"

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/gcpkms"
	"github.com/getsops/sops/v3/hckms"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/pgp"
	"github.com/getsops/sops/v3/publish"
	"go.yaml.in/yaml/v3"
)

type fileSystem interface {
	Stat(name string) (os.FileInfo, error)
}

type osFS struct {
	stat func(string) (os.FileInfo, error)
}

func (fs osFS) Stat(name string) (os.FileInfo, error) {
	return fs.stat(name)
}

var fs fileSystem = osFS{stat: os.Stat}

const (
	maxDepth            = 100
	configFileName      = ".sops.yaml"
	alternateConfigName = ".sops.yml"
)

// ConfigFileResult contains the path to a config file and any warnings
type ConfigFileResult struct {
	Path    string
	Warning string
}

// LookupConfigFile looks for a sops config file in the current working directory
// and on parent directories, up to the maxDepth limit.
// It returns a result containing the file path and any warnings.
func LookupConfigFile(start string) (ConfigFileResult, error) {
	filepath := path.Dir(start)
	var foundAlternatePath string

	for i := 0; i < maxDepth; i++ {
		configPath := path.Join(filepath, configFileName)
		_, err := fs.Stat(configPath)
		if err == nil {
			result := ConfigFileResult{Path: configPath}

			if foundAlternatePath != "" {
				result.Warning = fmt.Sprintf(
					"ignoring %q when searching for config file; the config file must be called %q; using %q instead",
					foundAlternatePath, configFileName, configPath)
			}
			return result, nil
		}

		// Check for alternate filename if we haven't found one yet
		if foundAlternatePath == "" {
			alternatePath := path.Join(filepath, alternateConfigName)
			_, altErr := fs.Stat(alternatePath)
			if altErr == nil {
				foundAlternatePath = alternatePath
			}
		}

		filepath = path.Join(filepath, "..")
	}

	// No config file found
	result := ConfigFileResult{}
	if foundAlternatePath != "" {
		result.Warning = fmt.Sprintf(
			"ignoring %q when searching for config file; the config file must be called %q",
			foundAlternatePath, configFileName)
	}

	return result, fmt.Errorf("config file not found")
}

// FindConfigFile looks for a sops config file in the current working directory and on parent directories, up to the limit defined by the maxDepth constant.
func FindConfigFile(start string) (string, error) {
	result, err := LookupConfigFile(start)
	return result.Path, err
}

type DotenvStoreConfig struct{}

type INIStoreConfig struct{}

type JSONStoreConfig struct {
	Indent int `yaml:"indent"`
}

type JSONBinaryStoreConfig struct {
	Indent int `yaml:"indent"`
}

type YAMLStoreConfig struct {
	Indent int `yaml:"indent"`
}

type StoresConfig struct {
	Dotenv     DotenvStoreConfig     `yaml:"dotenv"`
	INI        INIStoreConfig        `yaml:"ini"`
	JSONBinary JSONBinaryStoreConfig `yaml:"json_binary"`
	JSON       JSONStoreConfig       `yaml:"json"`
	YAML       YAMLStoreConfig       `yaml:"yaml"`
}

type configFile struct {
	CreationRules    []creationRule    `yaml:"creation_rules"`
	DestinationRules []destinationRule `yaml:"destination_rules"`
	Stores           StoresConfig      `yaml:"stores"`
}

type keyGroup struct {
	Merge   []keyGroup   `yaml:"merge"`
	KMS     []kmsKey     `yaml:"kms"`
	GCPKMS  []gcpKmsKey  `yaml:"gcp_kms"`
	HCKms   []hckmsKey   `yaml:"hckms"`
	AzureKV []azureKVKey `yaml:"azure_keyvault"`
	Vault   []string     `yaml:"hc_vault"`
	Age     []string     `yaml:"age"`
	PGP     []string     `yaml:"pgp"`
}

type gcpKmsKey struct {
	ResourceID string `yaml:"resource_id"`
}

type kmsKey struct {
	Arn        string             `yaml:"arn"`
	Role       string             `yaml:"role,omitempty"`
	Context    map[string]*string `yaml:"context"`
	AwsProfile string             `yaml:"aws_profile"`
}

type azureKVKey struct {
	VaultURL string `yaml:"vaultUrl"`
	Key      string `yaml:"key"`
	Version  string `yaml:"version"`
}

type hckmsKey struct {
	KeyID string `yaml:"key_id"`
}

type destinationRule struct {
	PathRegex        string       `yaml:"path_regex"`
	S3Bucket         string       `yaml:"s3_bucket"`
	S3Prefix         string       `yaml:"s3_prefix"`
	GCSBucket        string       `yaml:"gcs_bucket"`
	GCSPrefix        string       `yaml:"gcs_prefix"`
	VaultPath        string       `yaml:"vault_path"`
	VaultAddress     string       `yaml:"vault_address"`
	VaultKVMountName string       `yaml:"vault_kv_mount_name"`
	VaultKVVersion   int          `yaml:"vault_kv_version"`
	RecreationRule   creationRule `yaml:"recreation_rule,omitempty"`
	OmitExtensions   bool         `yaml:"omit_extensions"`
}

type creationRule struct {
	PathRegex               string      `yaml:"path_regex"`
	KMS                     interface{} `yaml:"kms"` // string or []string
	AwsProfile              string      `yaml:"aws_profile"`
	Age                     interface{} `yaml:"age"`     // string or []string
	PGP                     interface{} `yaml:"pgp"`     // string or []string
	GCPKMS                  interface{} `yaml:"gcp_kms"` // string or []string
	HCKms                   []string    `yaml:"hckms"`
	AzureKeyVault           interface{} `yaml:"azure_keyvault"`       // string or []string
	VaultURI                interface{} `yaml:"hc_vault_transit_uri"` // string or []string
	KeyGroups               []keyGroup  `yaml:"key_groups"`
	ShamirThreshold         int         `yaml:"shamir_threshold"`
	UnencryptedSuffix       string      `yaml:"unencrypted_suffix"`
	EncryptedSuffix         string      `yaml:"encrypted_suffix"`
	UnencryptedRegex        string      `yaml:"unencrypted_regex"`
	EncryptedRegex          string      `yaml:"encrypted_regex"`
	UnencryptedCommentRegex string      `yaml:"unencrypted_comment_regex"`
	EncryptedCommentRegex   string      `yaml:"encrypted_comment_regex"`
	MACOnlyEncrypted        bool        `yaml:"mac_only_encrypted"`
}

// Helper methods to safely extract keys as []string
func (c *creationRule) GetKMSKeys() ([]string, error) {
	return parseKeyField(c.KMS, "kms")
}

func (c *creationRule) GetAgeKeys() ([]string, error) {
	return parseKeyField(c.Age, "age")
}

func (c *creationRule) GetPGPKeys() ([]string, error) {
	return parseKeyField(c.PGP, "pgp")
}

func (c *creationRule) GetGCPKMSKeys() ([]string, error) {
	return parseKeyField(c.GCPKMS, "gcp_kms")
}

func (c *creationRule) GetAzureKeyVaultKeys() ([]string, error) {
	return parseKeyField(c.AzureKeyVault, "azure_keyvault")
}

func (c *creationRule) GetVaultURIs() ([]string, error) {
	return parseKeyField(c.VaultURI, "hc_vault_transit_uri")
}

// Utility function to handle both string and []string
func parseKeyField(field interface{}, fieldName string) ([]string, error) {
	if field == nil {
		return []string{}, nil
	}

	switch v := field.(type) {
	case string:
		if v == "" {
			return []string{}, nil
		}
		// Existing CSV parsing logic
		keys := strings.Split(v, ",")
		result := make([]string, 0, len(keys))
		for _, key := range keys {
			trimmed := strings.TrimSpace(key)
			if trimmed != "" { // Skip empty strings (fixes trailing comma issue)
				result = append(result, trimmed)
			}
		}
		return result, nil
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				result[i] = str
			} else {
				return nil, fmt.Errorf("invalid %s key configuration: expected string in list, got %T", fieldName, item)
			}
		}
		return result, nil
	case []string:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid %s key configuration: expected string, []string, or nil, got %T", fieldName, field)
	}
}

func NewStoresConfig() *StoresConfig {
	storesConfig := &StoresConfig{}
	storesConfig.JSON.Indent = -1
	storesConfig.JSONBinary.Indent = -1
	return storesConfig
}

// Load loads a sops config file into a temporary struct
func (f *configFile) load(bytes []byte) error {
	err := yaml.Unmarshal(bytes, f)
	if err != nil {
		return fmt.Errorf("Could not unmarshal config file: %s", err)
	}
	return nil
}

// Config is the configuration for a given SOPS file
type Config struct {
	KeyGroups               []sops.KeyGroup
	ShamirThreshold         int
	UnencryptedSuffix       string
	EncryptedSuffix         string
	UnencryptedRegex        string
	EncryptedRegex          string
	UnencryptedCommentRegex string
	EncryptedCommentRegex   string
	MACOnlyEncrypted        bool
	Destination             publish.Destination
	OmitExtensions          bool
}

// ErrPathNotAbsolute is returned by MatchRulesForFile when either of its path
// arguments is empty or relative.
var ErrPathNotAbsolute = errors.New("path must be absolute")

// MatchResult holds the rules from .sops.yaml that match a given file path.
// At most one creation_rule and at most one destination_rule will be matched
// (first-match-wins, preserving today's semantics in parseCreationRuleForFile
// and parseDestinationRuleForFile).
//
// API stability: The JSON output of the `sops config` subcommand is the
// public, versioned contract (see "schema_version"). The Go types in this
// package — MatchResult, CreationRuleMatch, DestinationRuleMatch, and the
// associated accessor methods — are internal-stable-only and may change
// without a major version bump. External consumers should depend on the
// JSON output, not the Go API.
type MatchResult struct {
	ConfigPath      string                // absolute, as received
	FilePath        string                // absolute, as received
	CreationRule    *CreationRuleMatch    // nil if no creation_rule matched
	DestinationRule *DestinationRuleMatch // nil if no destination_rule matched
}

// CreationRuleMatch describes which creation rule from .sops.yaml matched the
// queried file path. Rule is the raw internal struct; cmd-side packages
// translate it to a public JSON view.
type CreationRuleMatch struct {
	RuleIndex int
	Rule      creationRule
}

// DestinationRuleMatch describes which destination rule matched.
type DestinationRuleMatch struct {
	RuleIndex int
	Rule      destinationRule
}

// PathRegex returns the rule's path_regex value (or "" for the catch-all rule).
func (m *CreationRuleMatch) PathRegex() string { return m.Rule.PathRegex }

// KMSEntry is a serializable view of a KMS key entry from .sops.yaml.
// Entries without role/context/profile come from the short string-form syntax.
type KMSEntry struct {
	Arn        string
	Role       string
	Context    map[string]*string
	AwsProfile string
}

// KMSEntries returns the rule's KMS entries normalized to a slice. Handles
// both the short string form ("arn1,arn2") and the object form (yaml map with
// arn/role/context/aws_profile fields).
func (m *CreationRuleMatch) KMSEntries() ([]KMSEntry, error) {
	switch v := m.Rule.KMS.(type) {
	case nil:
		return nil, nil
	case string:
		out := []KMSEntry{}
		for _, k := range splitKMSString(v) {
			out = append(out, KMSEntry{Arn: k})
		}
		return out, nil
	case []interface{}:
		out := []KMSEntry{}
		for _, item := range v {
			switch x := item.(type) {
			case string:
				out = append(out, KMSEntry{Arn: x})
			case map[string]interface{}:
				entry := KMSEntry{}
				if s, ok := x["arn"].(string); ok {
					entry.Arn = s
				}
				if s, ok := x["role"].(string); ok {
					entry.Role = s
				}
				if s, ok := x["aws_profile"].(string); ok {
					entry.AwsProfile = s
				}
				if ctx, ok := x["context"].(map[string]interface{}); ok {
					entry.Context = make(map[string]*string, len(ctx))
					for k, v := range ctx {
						if s, ok := v.(string); ok {
							s := s
							entry.Context[k] = &s
						}
					}
				}
				out = append(out, entry)
			default:
				return nil, fmt.Errorf("unsupported kms entry type %T", item)
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported kms field type %T", v)
	}
}

// splitKMSString splits a comma-separated string form like "arn1,arn2" into
// a slice. Mirrors how parseKeyField handles the short form.
func splitKMSString(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (m *CreationRuleMatch) AgeRecipients() ([]string, error)      { return m.Rule.GetAgeKeys() }
func (m *CreationRuleMatch) PGPFingerprints() ([]string, error)    { return m.Rule.GetPGPKeys() }
func (m *CreationRuleMatch) GCPKMSResourceIDs() ([]string, error)  { return m.Rule.GetGCPKMSKeys() }
func (m *CreationRuleMatch) AzureKeyVaults() ([]string, error)     { return m.Rule.GetAzureKeyVaultKeys() }
func (m *CreationRuleMatch) HCVaultTransitURIs() ([]string, error) { return m.Rule.GetVaultURIs() }

// HCKmsKeyIDs returns the rule's HC KMS key_id values (the rule struct field
// is already a typed []string, no interface{} unwrap needed).
func (m *CreationRuleMatch) HCKmsKeyIDs() []string {
	out := make([]string, 0, len(m.Rule.HCKms))
	out = append(out, m.Rule.HCKms...)
	return out
}

func (m *CreationRuleMatch) ShamirThreshold() int            { return m.Rule.ShamirThreshold }
func (m *CreationRuleMatch) UnencryptedSuffix() string       { return m.Rule.UnencryptedSuffix }
func (m *CreationRuleMatch) EncryptedSuffix() string         { return m.Rule.EncryptedSuffix }
func (m *CreationRuleMatch) UnencryptedRegex() string        { return m.Rule.UnencryptedRegex }
func (m *CreationRuleMatch) EncryptedRegex() string          { return m.Rule.EncryptedRegex }
func (m *CreationRuleMatch) UnencryptedCommentRegex() string { return m.Rule.UnencryptedCommentRegex }
func (m *CreationRuleMatch) EncryptedCommentRegex() string   { return m.Rule.EncryptedCommentRegex }
func (m *CreationRuleMatch) MACOnlyEncrypted() bool          { return m.Rule.MACOnlyEncrypted }

// KeyGroupEntry is a serializable view of a key_group entry from .sops.yaml.
// Recipient lists are normalized; nested merge groups recurse.
type KeyGroupEntry struct {
	Merge          []KeyGroupEntry
	KMS            []KMSEntry
	GCPKMS         []string // resource IDs
	HCKms          []string // key IDs
	AzureKeyVault  []AzureKeyVaultEntry
	HCVaultTransit []string // URIs
	Age            []string
	PGP            []string
}

// AzureKeyVaultEntry is the serializable form of an Azure Key Vault key from
// .sops.yaml (used inside KeyGroupEntry).
type AzureKeyVaultEntry struct {
	VaultURL string
	Key      string
	Version  string
}

// KeyGroups returns the rule's key_groups normalized to a serializable view.
// Nested merge directives recurse.
func (m *CreationRuleMatch) KeyGroups() []KeyGroupEntry {
	return convertKeyGroups(m.Rule.KeyGroups)
}

func (m *DestinationRuleMatch) PathRegex() string    { return m.Rule.PathRegex }
func (m *DestinationRuleMatch) OmitExtensions() bool { return m.Rule.OmitExtensions }

// S3 returns ("", "", false) if this rule does not target S3.
func (m *DestinationRuleMatch) S3() (bucket, prefix string, ok bool) {
	if m.Rule.S3Bucket == "" {
		return "", "", false
	}
	return m.Rule.S3Bucket, m.Rule.S3Prefix, true
}

// GCS returns ("", "", false) if this rule does not target GCS.
func (m *DestinationRuleMatch) GCS() (bucket, prefix string, ok bool) {
	if m.Rule.GCSBucket == "" {
		return "", "", false
	}
	return m.Rule.GCSBucket, m.Rule.GCSPrefix, true
}

// Vault returns address, path, kvMountName, kvVersion and ok=false if no
// Vault destination is configured for this rule.
func (m *DestinationRuleMatch) Vault() (address, path, kvMountName string, kvVersion int, ok bool) {
	if m.Rule.VaultPath == "" {
		return "", "", "", 0, false
	}
	return m.Rule.VaultAddress, m.Rule.VaultPath, m.Rule.VaultKVMountName, m.Rule.VaultKVVersion, true
}

// RecreationRule converts the nested creation_rule inside this destination
// rule. Returns nil if the user did not specify a recreation_rule. The
// returned match has RuleIndex == 0 because the nested rule has no top-level
// position; the view layer treats this field as meaningless for nested
// recreation rules.
func (m *DestinationRuleMatch) RecreationRule() *CreationRuleMatch {
	if reflect.DeepEqual(m.Rule.RecreationRule, creationRule{}) {
		return nil
	}
	return &CreationRuleMatch{RuleIndex: 0, Rule: m.Rule.RecreationRule}
}

func convertKeyGroups(in []keyGroup) []KeyGroupEntry {
	out := make([]KeyGroupEntry, 0, len(in))
	for _, g := range in {
		// Defensive copies on the string slices: callers must not be able to
		// mutate the internal keyGroup via the returned view.
		entry := KeyGroupEntry{
			Merge:          convertKeyGroups(g.Merge),
			HCVaultTransit: append([]string{}, g.Vault...),
			Age:            append([]string{}, g.Age...),
			PGP:            append([]string{}, g.PGP...),
		}
		for _, k := range g.KMS {
			entry.KMS = append(entry.KMS, KMSEntry{
				Arn:        k.Arn,
				Role:       k.Role,
				Context:    k.Context,
				AwsProfile: k.AwsProfile,
			})
		}
		for _, k := range g.GCPKMS {
			entry.GCPKMS = append(entry.GCPKMS, k.ResourceID)
		}
		for _, k := range g.HCKms {
			entry.HCKms = append(entry.HCKms, k.KeyID)
		}
		for _, k := range g.AzureKV {
			entry.AzureKeyVault = append(entry.AzureKeyVault, AzureKeyVaultEntry{
				VaultURL: k.VaultURL,
				Key:      k.Key,
				Version:  k.Version,
			})
		}
		out = append(out, entry)
	}
	return out
}

func deduplicateKeygroup(group sops.KeyGroup) sops.KeyGroup {
	var deduplicatedKeygroup sops.KeyGroup

	unique := make(map[string]bool)
	for _, v := range group {
		key := fmt.Sprintf("%T/%v", v, v.ToString())
		if _, ok := unique[key]; ok {
			// key already contained, therefore not unique
			continue
		}

		deduplicatedKeygroup = append(deduplicatedKeygroup, v)
		unique[key] = true
	}

	return deduplicatedKeygroup
}

func extractMasterKeys(group keyGroup) (sops.KeyGroup, error) {
	var keyGroup sops.KeyGroup
	for _, k := range group.Merge {
		subKeyGroup, err := extractMasterKeys(k)
		if err != nil {
			return nil, err
		}
		keyGroup = append(keyGroup, subKeyGroup...)
	}

	for _, k := range group.Age {
		keys, err := age.MasterKeysFromRecipients(k)
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			keyGroup = append(keyGroup, key)
		}
	}
	for _, k := range group.PGP {
		keyGroup = append(keyGroup, pgp.NewMasterKeyFromFingerprint(k))
	}
	for _, k := range group.KMS {
		keyGroup = append(keyGroup, kms.NewMasterKeyWithProfile(k.Arn, k.Role, k.Context, k.AwsProfile))
	}
	for _, k := range group.GCPKMS {
		keyGroup = append(keyGroup, gcpkms.NewMasterKeyFromResourceID(k.ResourceID))
	}
	for _, k := range group.HCKms {
		key, err := hckms.NewMasterKey(k.KeyID)
		if err != nil {
			return nil, err
		}
		keyGroup = append(keyGroup, key)
	}
	for _, k := range group.AzureKV {
		if key, err := azkv.NewMasterKeyWithOptionalVersion(k.VaultURL, k.Key, k.Version); err == nil {
			keyGroup = append(keyGroup, key)
		} else {
			return nil, err
		}
	}
	for _, k := range group.Vault {
		if masterKey, err := hcvault.NewMasterKeyFromURI(k); err == nil {
			keyGroup = append(keyGroup, masterKey)
		} else {
			return nil, err
		}
	}
	return deduplicateKeygroup(keyGroup), nil
}

func getKeysWithValidation(getKeysFunc func() ([]string, error), keyType string) ([]string, error) {
	keys, err := getKeysFunc()
	if err != nil {
		return nil, fmt.Errorf("invalid %s key configuration: %w", keyType, err)
	}
	return keys, nil
}

func getKeyGroupsFromCreationRule(cRule *creationRule, kmsEncryptionContext map[string]*string) ([]sops.KeyGroup, error) {
	var groups []sops.KeyGroup
	if len(cRule.KeyGroups) > 0 {
		for _, group := range cRule.KeyGroups {
			keyGroup, err := extractMasterKeys(group)
			if err != nil {
				return nil, err
			}
			groups = append(groups, keyGroup)
		}
	} else {
		var keyGroup sops.KeyGroup
		ageKeys, err := getKeysWithValidation(cRule.GetAgeKeys, "age")
		if err != nil {
			return nil, err
		}

		if len(ageKeys) > 0 {
			ageKeys, err := age.MasterKeysFromRecipients(strings.Join(ageKeys, ","))
			if err != nil {
				return nil, err
			} else {
				for _, ak := range ageKeys {
					keyGroup = append(keyGroup, ak)
				}
			}
		}
		pgpKeys, err := getKeysWithValidation(cRule.GetPGPKeys, "pgp")
		if err != nil {
			return nil, err
		}
		for _, k := range pgp.MasterKeysFromFingerprintString(strings.Join(pgpKeys, ",")) {
			keyGroup = append(keyGroup, k)
		}
		kmsKeys, err := getKeysWithValidation(cRule.GetKMSKeys, "kms")
		if err != nil {
			return nil, err
		}
		for _, k := range kms.MasterKeysFromArnString(strings.Join(kmsKeys, ","), kmsEncryptionContext, cRule.AwsProfile) {
			keyGroup = append(keyGroup, k)
		}
		gcpkmsKeys, err := getKeysWithValidation(cRule.GetGCPKMSKeys, "gcpkms")
		if err != nil {
			return nil, err
		}
		for _, k := range gcpkms.MasterKeysFromResourceIDString(strings.Join(gcpkmsKeys, ",")) {
			keyGroup = append(keyGroup, k)
		}
		hckmsMasterKeys, err := hckms.NewMasterKeyFromKeyIDString(strings.Join(cRule.HCKms, ","))
		if err != nil {
			return nil, err
		}
		for _, k := range hckmsMasterKeys {
			keyGroup = append(keyGroup, k)
		}
		azKeys, err := getKeysWithValidation(cRule.GetAzureKeyVaultKeys, "azure_keyvault")
		if err != nil {
			return nil, err
		}
		azureKeys, err := azkv.MasterKeysFromURLs(strings.Join(azKeys, ","))
		if err != nil {
			return nil, err
		}
		for _, k := range azureKeys {
			keyGroup = append(keyGroup, k)
		}
		vaultKeyUris, err := getKeysWithValidation(cRule.GetVaultURIs, "vault")
		if err != nil {
			return nil, err
		}
		vaultKeys, err := hcvault.NewMasterKeysFromURIs(strings.Join(vaultKeyUris, ","))
		if err != nil {
			return nil, err
		}
		for _, k := range vaultKeys {
			keyGroup = append(keyGroup, k)
		}
		groups = append(groups, keyGroup)
	}
	return groups, nil
}

func loadConfigFile(confPath string) (*configFile, error) {
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %s", err)
	}
	conf := &configFile{}
	conf.Stores = *NewStoresConfig()
	err = conf.load(confBytes)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err)
	}
	return conf, nil
}

func configFromRule(rule *creationRule, kmsEncryptionContext map[string]*string) (*Config, error) {
	cryptRuleCount := 0
	if rule.UnencryptedSuffix != "" {
		cryptRuleCount++
	}
	if rule.EncryptedSuffix != "" {
		cryptRuleCount++
	}
	if rule.UnencryptedRegex != "" {
		cryptRuleCount++
	}
	if rule.EncryptedRegex != "" {
		cryptRuleCount++
	}
	if rule.UnencryptedCommentRegex != "" {
		cryptRuleCount++
	}
	if rule.EncryptedCommentRegex != "" {
		cryptRuleCount++
	}

	if cryptRuleCount > 1 {
		return nil, fmt.Errorf("error loading config: cannot use more than one of encrypted_suffix, unencrypted_suffix, encrypted_regex, unencrypted_regex, encrypted_comment_regex, or unencrypted_comment_regex for the same rule")
	}

	groups, err := getKeyGroupsFromCreationRule(rule, kmsEncryptionContext)
	if err != nil {
		return nil, err
	}

	return &Config{
		KeyGroups:               groups,
		ShamirThreshold:         rule.ShamirThreshold,
		UnencryptedSuffix:       rule.UnencryptedSuffix,
		EncryptedSuffix:         rule.EncryptedSuffix,
		UnencryptedRegex:        rule.UnencryptedRegex,
		EncryptedRegex:          rule.EncryptedRegex,
		UnencryptedCommentRegex: rule.UnencryptedCommentRegex,
		EncryptedCommentRegex:   rule.EncryptedCommentRegex,
		MACOnlyEncrypted:        rule.MACOnlyEncrypted,
	}, nil
}

func parseDestinationRuleForFile(conf *configFile, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	var rule *creationRule
	var dRule *destinationRule

	if len(conf.DestinationRules) > 0 {
		for _, r := range conf.DestinationRules {
			if r.PathRegex == "" {
				dRule = &r
				rule = &dRule.RecreationRule
				break
			}
			if r.PathRegex != "" {
				if match, _ := regexp.MatchString(r.PathRegex, filePath); match {
					dRule = &r
					rule = &dRule.RecreationRule
					break
				}
			}
		}
	}

	if dRule == nil {
		return nil, fmt.Errorf("error loading config: no matching destination found in config")
	}

	var dest publish.Destination
	destinationCount := 0
	if dRule.S3Bucket != "" {
		destinationCount++
	}
	if dRule.GCSBucket != "" {
		destinationCount++
	}
	if dRule.VaultPath != "" {
		destinationCount++
	}

	if destinationCount > 1 {
		return nil, fmt.Errorf("error loading config: more than one destinations were found in a single destination rule, you can only use one per rule")
	}
	if dRule.S3Bucket != "" {
		dest = publish.NewS3Destination(dRule.S3Bucket, dRule.S3Prefix)
	}
	if dRule.GCSBucket != "" {
		dest = publish.NewGCSDestination(dRule.GCSBucket, dRule.GCSPrefix)
	}
	if dRule.VaultPath != "" {
		dest = publish.NewVaultDestination(dRule.VaultAddress, dRule.VaultPath, dRule.VaultKVMountName, dRule.VaultKVVersion)
	}

	config, err := configFromRule(rule, kmsEncryptionContext)
	if err != nil {
		return nil, err
	}
	config.Destination = dest
	config.OmitExtensions = dRule.OmitExtensions

	return config, nil
}

func parseCreationRuleForFile(conf *configFile, confPath, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	// If config file doesn't contain CreationRules (it's empty or only contains DestionationRules), assume it does not exist
	if conf.CreationRules == nil {
		return nil, nil
	}

	configDir, err := filepath.Abs(filepath.Dir(confPath))
	if err != nil {
		return nil, err
	}

	// compare file path relative to path of config file
	filePath = strings.TrimPrefix(filePath, configDir+string(filepath.Separator))

	var rule *creationRule

	for _, r := range conf.CreationRules {
		if r.PathRegex == "" {
			rule = &r
			break
		}
		reg, err := regexp.Compile(r.PathRegex)
		if err != nil {
			return nil, fmt.Errorf("can not compile regexp: %w", err)
		}
		if reg.MatchString(filePath) {
			rule = &r
			break
		}
	}

	if rule == nil {
		return nil, fmt.Errorf("error loading config: no matching creation rules found")
	}

	config, err := configFromRule(rule, kmsEncryptionContext)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// matchRulesForConfig matches a file path against an already-parsed config.
// Both absConfPath and absFilePath MUST be absolute (caller's responsibility).
// This is the file-IO-free core that production callers reach via MatchRulesForFile
// and that tests can exercise directly with parseConfigFile-built fixtures.
func matchRulesForConfig(conf *configFile, absConfPath, absFilePath string) (*MatchResult, error) {
	result := &MatchResult{
		ConfigPath: absConfPath,
		FilePath:   absFilePath,
	}

	matchPath := normalizeMatchPath(absConfPath, absFilePath)

	for i, r := range conf.CreationRules {
		if r.PathRegex == "" {
			rule := r
			result.CreationRule = &CreationRuleMatch{RuleIndex: i, Rule: rule}
			break
		}
		reg, err := regexp.Compile(r.PathRegex)
		if err != nil {
			return nil, fmt.Errorf("can not compile regexp: %w", err)
		}
		if reg.MatchString(matchPath) {
			rule := r
			result.CreationRule = &CreationRuleMatch{RuleIndex: i, Rule: rule}
			break
		}
	}

	for i, r := range conf.DestinationRules {
		if r.PathRegex == "" {
			rule := r
			result.DestinationRule = &DestinationRuleMatch{RuleIndex: i, Rule: rule}
			break
		}
		reg, err := regexp.Compile(r.PathRegex)
		if err != nil {
			return nil, fmt.Errorf("can not compile regexp: %w", err)
		}
		if reg.MatchString(matchPath) {
			rule := r
			result.DestinationRule = &DestinationRuleMatch{RuleIndex: i, Rule: rule}
			break
		}
	}

	return result, nil
}

// normalizeMatchPath returns the path that should be matched against rule
// path_regex values. When absFilePath is inside the config file's directory
// tree, the config-dir prefix is stripped so users' path_regex values can be
// written as repo-relative. When the file is outside the tree (or on a
// different Windows drive), the absolute path is used as-is.
//
// Uses filepath.Rel for platform-aware path arithmetic; correctly handles
// trailing separators, mixed separators on Windows, and case-insensitive
// drive letters.
func normalizeMatchPath(absConfPath, absFilePath string) string {
	configDir := filepath.Dir(absConfPath)
	rel, err := filepath.Rel(configDir, absFilePath)
	if err != nil {
		// Different Windows drives, or otherwise unrelatable. Use absolute.
		return absFilePath
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		// File is above the config dir. Use absolute.
		return absFilePath
	}
	return rel
}

// MatchRulesForFile loads the config at absConfPath and returns the
// creation_rule and destination_rule (each at most one) that apply to
// absFilePath. Both arguments MUST be absolute paths — MatchRulesForFile
// does NOT call filepath.Abs (which would implicitly depend on os.Getwd())
// and does NOT perform .sops.yaml auto-discovery. The CLI layer is
// responsible for FindConfigFile, the --config flag, and resolving any
// relative input to absolute before calling this function.
//
// File existence on absFilePath is NOT checked; only path matching occurs.
func MatchRulesForFile(absConfPath, absFilePath string) (*MatchResult, error) {
	if absConfPath == "" || !filepath.IsAbs(absConfPath) {
		return nil, ErrPathNotAbsolute
	}
	if absFilePath == "" || !filepath.IsAbs(absFilePath) {
		return nil, ErrPathNotAbsolute
	}
	conf, err := loadConfigFile(absConfPath)
	if err != nil {
		return nil, err
	}
	return matchRulesForConfig(conf, absConfPath, absFilePath)
}

// LoadCreationRuleForFile load the configuration for a given SOPS file from the config file at confPath. A kmsEncryptionContext
// should be provided for configurations that do not contain key groups, as there's no way to specify context inside
// a SOPS config file outside of key groups.
func LoadCreationRuleForFile(confPath string, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	conf, err := loadConfigFile(confPath)
	if err != nil {
		return nil, err
	}

	return parseCreationRuleForFile(conf, confPath, filePath, kmsEncryptionContext)
}

// LoadDestinationRuleForFile works the same as LoadCreationRuleForFile, but gets the "creation_rule" from the matching destination_rule's
// "recreation_rule".
func LoadDestinationRuleForFile(confPath string, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	conf, err := loadConfigFile(confPath)
	if err != nil {
		return nil, err
	}
	return parseDestinationRuleForFile(conf, filePath, kmsEncryptionContext)
}

func LoadStoresConfig(confPath string) (*StoresConfig, error) {
	conf, err := loadConfigFile(confPath)
	if err != nil {
		return nil, err
	}
	return &conf.Stores, nil
}
