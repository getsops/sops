package init

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"go.yaml.in/yaml/v3"
)

type CreationRule struct {
	PathRegex      string   `yaml:"path_regex"`
	Age            []string `yaml:"age"`
	PGP            []string `yaml:"pgp"`
	AWSKMS         []string `yaml:"kms"`
	GCPKMS         []string `yaml:"gcp_kms"`
	AzureKeyVault  []string `yaml:"azure_keyvault"`
	HuaweiCloud    []string `yaml:"hckms"`
	HashicorpVault []string `yaml:"hc_vault_transit_uri"`
}

type ConfigFile struct {
	CreationRules []CreationRule
}

type InitCommandArgs struct {
	ConfigFilePath string
	IsVerbose      bool
	CreationRule   CreationRule
}

func Init(args InitCommandArgs) error {

	fileInfo, err := os.Stat(args.ConfigFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return common.NewExitError(
			fmt.Errorf("%s does not exist, please ensure this is an existing directory", args.ConfigFilePath),
			codes.CouldNotReadInputFile,
		)
	}

	if !fileInfo.IsDir() {
		return common.NewExitError(
			fmt.Errorf("argument is not a directory, please provide a directory path"),
			codes.ErrorGeneric,
		)
	}

	if err != nil {
		return common.NewExitError(
			fmt.Errorf("failed to retrieve file stats: %e", err),
			codes.CouldNotReadInputFile,
		)
	}

	args.ConfigFilePath = filepath.Join(args.ConfigFilePath, ".sops.yaml")

	_, err = os.Stat(args.ConfigFilePath)
	if !errors.Is(err, os.ErrNotExist) {
		return common.NewExitError(
			fmt.Errorf("%s already exists", args.ConfigFilePath),
			codes.ErrorGeneric,
		)
	}

	if args.IsVerbose {
		fmt.Printf("generating .sops.yaml -> %s\n", args.ConfigFilePath)
	}

	defaultConfigFile := ConfigFile{
		CreationRules: []CreationRule{
			CreationRule{
				PathRegex: "secrets/[^/]+\\.(yaml|json|env|ini)$",
				PGP: []string{
					"2504791468b153b8a3963cc97ba53d1919c5dfd4!",
					"CHANGE-ME",
				},
				Age: []string{
					"age12zlz6lvcdk6eqaewfylg35w0syh58sm7gh53q5vvn7hd7c6nngyseftjxl",
					"CHANGE-ME",
				},
			},
		},
	}

	result, err := yaml.Marshal(&defaultConfigFile)
	if err != nil {
		return common.NewExitError(
			fmt.Errorf("error: %e", err),
			codes.ErrorReadingConfig,
		)
	}

	var mutatedResult = fmt.Sprintf(
		"# Example .sops.yaml\n# Please update the values according to your setup\n# https://getsops.io/docs/\n\n%s\n",
		string(result),
	)

	file, err := os.Create(args.ConfigFilePath)
	if err != nil {
		return common.NewExitError(
			fmt.Errorf("failed to open file: %e", err),
			codes.ErrorGeneric,
		)
	}

	_, err = file.Write([]byte(mutatedResult))
	if err != nil {
		return common.NewExitError(
			fmt.Errorf("failed to write to file: %e", err),
			codes.CouldNotWriteOutputFile,
		)
	}

	if err = file.Close(); err != nil {
		return common.NewExitError(
			fmt.Errorf("failed to close file descriptor: %e", err),
			codes.ErrorGeneric,
		)
	}

	if args.IsVerbose {
		fmt.Printf("generated yaml\noutput:\n%s\n", mutatedResult)
	}

	return nil

}
