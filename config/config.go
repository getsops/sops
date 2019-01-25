/*
Package config provides a way to find and load SOPS configuration files
*/
package config //import "go.mozilla.org/sops/config"

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/mozilla-services/yaml"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/azkv"
	"go.mozilla.org/sops/gcpkms"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/logging"
	"go.mozilla.org/sops/pgp"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("CONFIG")
}

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
	maxDepth       = 100
	configFileName = ".sops.yaml"
)

// FindConfigFile looks for a sops config file in the current working directory and on parent directories, up to the limit defined by the maxDepth constant.
func FindConfigFile(start string) (string, error) {
	filepath := path.Dir(start)
	for i := 0; i < maxDepth; i++ {
		_, err := fs.Stat(path.Join(filepath, configFileName))
		if err != nil {
			filepath = path.Join(filepath, "..")
		} else {
			return path.Join(filepath, configFileName), nil
		}
	}
	return "", fmt.Errorf("Config file not found")
}

type configFile struct {
	CreationRules []creationRule `yaml:"creation_rules"`
}

type keyGroup struct {
	KMS     []kmsKey
	GCPKMS  []gcpKmsKey  `yaml:"gcp_kms"`
	AzureKV []azureKVKey `yaml:"azure_keyvault"`
	PGP     []string
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

type creationRule struct {
	FilenameRegex     string `yaml:"filename_regex"`
	PathRegex         string `yaml:"path_regex"`
	KMS               string
	AwsProfile        string `yaml:"aws_profile"`
	PGP               string
	GCPKMS            string     `yaml:"gcp_kms"`
	AzureKeyVault     string     `yaml:"azure_keyvault"`
	KeyGroups         []keyGroup `yaml:"key_groups"`
	ShamirThreshold   int        `yaml:"shamir_threshold"`
	UnencryptedSuffix string     `yaml:"unencrypted_suffix"`
	EncryptedSuffix   string     `yaml:"encrypted_suffix"`
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
	KeyGroups         []sops.KeyGroup
	ShamirThreshold   int
	UnencryptedSuffix string
	EncryptedSuffix   string
}

func loadForFileFromBytes(confBytes []byte, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	conf := configFile{}
	err := conf.load(confBytes)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err)
	}
	var rule *creationRule

	for _, r := range conf.CreationRules {
		if r.PathRegex == "" && r.FilenameRegex == "" {
			rule = &r
			break
		}
		if r.PathRegex != "" && r.FilenameRegex != "" {
			return nil, fmt.Errorf("error loading config: both filename_regex and path_regex were found, use only path_regex")
		}
		if r.FilenameRegex != "" {
			if match, _ := regexp.MatchString(r.FilenameRegex, filePath); match {
				log.Warn("The key: filename_regex will be removed in a future release. Instead use key: path_regex in your .sops.yaml file")
				rule = &r
				break
			}
		}
		if r.PathRegex != "" {
			if match, _ := regexp.MatchString(r.PathRegex, filePath); match {
				rule = &r
				break
			}
		}
	}

	if rule == nil {
		return nil, fmt.Errorf("error loading config: no matching creation rules found")
	}

	if rule.UnencryptedSuffix != "" && rule.EncryptedSuffix != "" {
		return nil, fmt.Errorf("error loading config: cannot use both encrypted_suffix and unencrypted_suffix for the same rule")
	}

	var groups []sops.KeyGroup
	if len(rule.KeyGroups) > 0 {
		for _, group := range rule.KeyGroups {
			var keyGroup sops.KeyGroup
			for _, k := range group.PGP {
				keyGroup = append(keyGroup, pgp.NewMasterKeyFromFingerprint(k))
			}
			for _, k := range group.KMS {
				keyGroup = append(keyGroup, kms.NewMasterKey(k.Arn, k.Role, k.Context))
			}
			for _, k := range group.GCPKMS {
				keyGroup = append(keyGroup, gcpkms.NewMasterKeyFromResourceID(k.ResourceID))
			}
			groups = append(groups, keyGroup)
		}
	} else {
		var keyGroup sops.KeyGroup
		for _, k := range pgp.MasterKeysFromFingerprintString(rule.PGP) {
			keyGroup = append(keyGroup, k)
		}
		for _, k := range kms.MasterKeysFromArnString(rule.KMS, kmsEncryptionContext, rule.AwsProfile) {
			keyGroup = append(keyGroup, k)
		}
		for _, k := range gcpkms.MasterKeysFromResourceIDString(rule.GCPKMS) {
			keyGroup = append(keyGroup, k)
		}
		azureKeys, err := azkv.MasterKeysFromURLs(rule.AzureKeyVault)
		if err != nil {
			return nil, err
		}
		for _, k := range azureKeys {
			keyGroup = append(keyGroup, k)
		}
		groups = append(groups, keyGroup)
	}
	return &Config{
		KeyGroups:         groups,
		ShamirThreshold:   rule.ShamirThreshold,
		UnencryptedSuffix: rule.UnencryptedSuffix,
		EncryptedSuffix:   rule.EncryptedSuffix,
	}, nil
}

// LoadForFile load the configuration for a given SOPS file from the config file at confPath. A kmsEncryptionContext
// should be provided for configurations that do not contain key groups, as there's no way to specify context inside
// a SOPS config file outside of key groups.
func LoadForFile(confPath string, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %s", err)
	}
	return loadForFileFromBytes(confBytes, filePath, kmsEncryptionContext)
}
