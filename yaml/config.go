package yaml //import "go.mozilla.org/sops/yaml"

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/mozilla-services/yaml"
	"go.mozilla.org/sops"
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
	KMS string
	PGP string
}

type creationRule struct {
	FilenameRegex string `yaml:"filename_regex"`
	KMS           string
	PGP           string
	KeyGroups     []keyGroup `yaml:"key_groups"`
}

// Load loads a sops config file into a temporary struct
func (f *configFile) load(bytes []byte) error {
	err := yaml.Unmarshal(bytes, f)
	if err != nil {
		return fmt.Errorf("Could not unmarshal config file: %s", err)
	}
	return nil
}

func KeyGroupsForFile(filepath string, confBytes []byte, kmsEncryptionContext map[string]*string) ([]sops.KeyGroup, error) {
	var err error
	if confBytes == nil {
		var confPath string
		confPath, err = FindConfigFile(".")
		if err != nil {
			return nil, err
		}
		confBytes, err = ioutil.ReadFile(confPath)
	}
	if err != nil {
		return nil, fmt.Errorf("Could not read config file: %s", err)
	}
	conf := configFile{}
	err = conf.load(confBytes)
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %s", err)
	}
	var groups []sops.KeyGroup
	for _, rule := range conf.CreationRules {
		if match, _ := regexp.MatchString(rule.FilenameRegex, filepath); match {
			if len(rule.KeyGroups) > 0 {
				for _, group := range rule.KeyGroups {
					var keyGroup sops.KeyGroup
					for _, k := range pgp.MasterKeysFromFingerprintString(group.PGP) {
						keyGroup = append(keyGroup, k)
					}
					for _, k := range kms.MasterKeysFromArnString(group.KMS, kmsEncryptionContext) {
						keyGroup = append(keyGroup, k)
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
				groups = append(groups, keyGroup)
			}
			return groups, nil
		}
	}
	return nil, nil
}
