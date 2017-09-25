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
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/gcpkms"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
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
	KMS    []kmsKey
	GCPKMS []gcpKmsKey `yaml:"gcp_kms"`
	PGP    []string
}

type gcpKmsKey struct {
	ResourceID string `yaml:"resource_id"`
}

type kmsKey struct {
	Arn     string             `yaml:"arn"`
	Role    string             `yaml:"role,omitempty"`
	Context map[string]*string `yaml:"context"`
}

type creationRule struct {
	FilenameRegex   string `yaml:"filename_regex"`
	KMS             string
	PGP             string
	GCPKMS          string     `yaml:"gcp_kms"`
	KeyGroups       []keyGroup `yaml:"key_groups"`
	ShamirThreshold int        `yaml:"shamir_threshold"`
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
	KeyGroups       []sops.KeyGroup
	ShamirThreshold int
}

func loadForFileFromBytes(confBytes []byte, filePath string, kmsEncryptionContext map[string]*string) (*Config, error) {
	conf := configFile{}
	err := conf.load(confBytes)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err)
	}
	var rule *creationRule

	for _, r := range conf.CreationRules {
		if match, _ := regexp.MatchString(r.FilenameRegex, filePath); match {
			rule = &r
			break
		}
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
		for _, k := range kms.MasterKeysFromArnString(rule.KMS, kmsEncryptionContext) {
			keyGroup = append(keyGroup, k)
		}
		for _, k := range gcpkms.MasterKeysFromResourceIDString(rule.GCPKMS) {
			keyGroup = append(keyGroup, k)
		}
		groups = append(groups, keyGroup)
	}
	return &Config{
		KeyGroups:       groups,
		ShamirThreshold: rule.ShamirThreshold,
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
