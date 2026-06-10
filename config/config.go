/*
Package config provides a way to find and load SOPS configuration files
*/
package config //import "github.com/getsops/sops/v3/config"

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keys"
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
	Merge     []keyGroup     `yaml:"merge"`
	Providers map[string]any `yaml:",inline"`
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
	Timeout                 string      `yaml:"timeout,omitempty"`
	KeyGroups               []keyGroup  `yaml:"key_groups"`
	ShamirThreshold         int         `yaml:"shamir_threshold"`
	UnencryptedSuffix       string      `yaml:"unencrypted_suffix"`
	EncryptedSuffix         string      `yaml:"encrypted_suffix"`
	UnencryptedRegex        string      `yaml:"unencrypted_regex"`
	EncryptedRegex          string      `yaml:"encrypted_regex"`
	UnencryptedCommentRegex string      `yaml:"unencrypted_comment_regex"`
	EncryptedCommentRegex   string      `yaml:"encrypted_comment_regex"`
	MACOnlyEncrypted        bool        `yaml:"mac_only_encrypted"`

	Providers map[string]any `yaml:",inline"`
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

func extractMasterKeys(group keyGroup, opts keys.CreationOptions) (sops.KeyGroup, error) {
	var kg sops.KeyGroup
	for _, k := range group.Merge {
		subKeyGroup, err := extractMasterKeys(k, opts)
		if err != nil {
			return nil, err
		}
		kg = append(kg, subKeyGroup...)
	}

	order := []string{"pgp", "kms", "gcp_kms", "hckms", "azure_keyvault", "hc_vault", "age", "plugins"}
	for _, providerName := range order {
		providerData, ok := group.Providers[providerName]
		if !ok {
			continue
		}
		lookupName := providerName
		if lookupName == "azure_keyvault" {
			lookupName = "azure_kv"
		}
		provider := keys.GetProvider(lookupName)
		if provider == nil {
			continue
		}
		masterKeys, err := provider.KeysFromConfig(providerData, opts)
		if err != nil {
			return nil, err
		}
		for _, mk := range masterKeys {
			kg = append(kg, mk)
		}
	}

	return deduplicateKeygroup(kg), nil
}

func getKeyGroupsFromCreationRule(cRule *creationRule, kmsEncryptionContext map[string]*string) ([]sops.KeyGroup, error) {
	var groups []sops.KeyGroup
	opts := keys.CreationOptions{
		KmsEncryptionContext: kmsEncryptionContext,
		GlobalConfig:         cRule.Providers,
	}

	if len(cRule.KeyGroups) > 0 {
		for _, group := range cRule.KeyGroups {
			keyGroup, err := extractMasterKeys(group, opts)
			if err != nil {
				return nil, err
			}
			groups = append(groups, keyGroup)
		}
	} else {
		var keyGroup sops.KeyGroup
		order := []string{"pgp", "kms", "gcp_kms", "hckms", "azure_keyvault", "hc_vault", "hc_vault_transit_uri", "hc_vault_uris", "age", "plugins"}
		
		for _, providerName := range order {
			providerData, ok := cRule.Providers[providerName]
			if !ok {
				continue
			}
			
			lookupName := providerName
			if lookupName == "azure_keyvault" {
				lookupName = "azure_kv"
			} else if lookupName == "hc_vault_transit_uri" || lookupName == "hc_vault_uris" {
				lookupName = "hc_vault"
			}
			
			provider := keys.GetProvider(lookupName)
			if provider == nil {
				continue
			}
			masterKeys, err := provider.KeysFromConfig(providerData, opts)
			if err != nil {
				return nil, err
			}
			for _, mk := range masterKeys {
				keyGroup = append(keyGroup, mk)
			}
		}
		
		if len(keyGroup) > 0 {
			groups = append(groups, deduplicateKeygroup(keyGroup))
		}
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
	if conf.CreationRules == nil {
		return nil, nil
	}

	configDir, err := filepath.Abs(filepath.Dir(confPath))
	if err != nil {
		return nil, err
	}

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

func LoadCreationRuleForFile(confPath string, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	conf, err := loadConfigFile(confPath)
	if err != nil {
		return nil, err
	}

	return parseCreationRuleForFile(conf, confPath, filePath, kmsEncryptionContext)
}

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
