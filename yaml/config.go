package yaml

import (
	"fmt"
	"github.com/autrilla/yaml"
	"io/ioutil"
	"os"
	"path"
	"regexp"
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

type creationRule struct {
	FilenameRegex string `yaml:"filename_regex"`
	KMS           string
	PGP           string
}

// Load loads a sops config file into a temporary struct
func (f *configFile) load(bytes []byte) error {
	err := yaml.Unmarshal(bytes, f)
	if err != nil {
		return fmt.Errorf("Could not unmarshal config file: %s", err)
	}
	return nil
}

// MasterKeyStringsForFile returns a comma separated string of KMS ARNs and a comma separated list of PGP fingerprints. If the config bytes are left empty, the function will look for the config file by itself.
func MasterKeyStringsForFile(filepath string, confBytes []byte) (kms, pgp string, err error) {
	if confBytes == nil {
		confPath, err := FindConfigFile(filepath)
		if err != nil {
			return "", "", err
		}
		confBytes, err = ioutil.ReadFile(confPath)
	}
	if err != nil {
		return "", "", fmt.Errorf("Could not read config file: %s", err)
	}
	conf := configFile{}
	err = conf.load(confBytes)
	if err != nil {
		return "", "", fmt.Errorf("Error loading config: %s", err)
	}
	for _, rule := range conf.CreationRules {
		if match, _ := regexp.MatchString(rule.FilenameRegex, filepath); match {
			return rule.KMS, rule.PGP, nil
		}
	}
	return "", "", nil
}
