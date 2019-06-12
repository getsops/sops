package publish

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/config"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/logging"
	"go.mozilla.org/sops/version"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("PUBLISH")
}

type Opts struct {
	Interactive bool
	Cipher      sops.Cipher
	ConfigPath  string
	InputPath   string
	KeyServices []keyservice.KeyServiceClient
	InputStore  sops.Store
}

func Run(opts Opts) error {
	var fileContents []byte
	path, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("can't operate on a directory")
	}
	_, fileName := filepath.Split(path)

	conf, err := config.LoadDestinationRuleForFile(opts.ConfigPath, opts.InputPath, make(map[string]*string))
	if err != nil {
		return err
	}
	if conf.Destination == nil {
		return errors.New("no destination configured for this file")
	}

	// Check that this is a sops-encrypted file
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return err
	}

	// Re-encrypt if settings exist to do so
	if len(conf.KeyGroups[0]) != 0 {
		log.Debug("Re-encrypting tree before publishing")
		_, err = common.DecryptTree(common.DecryptTreeOpts{
			Cipher:      opts.Cipher,
			IgnoreMac:   false,
			Tree:        tree,
			KeyServices: opts.KeyServices,
		})
		if err != nil {
			return err
		}

		tree.Metadata = sops.Metadata{
			KeyGroups:         conf.KeyGroups,
			UnencryptedSuffix: conf.UnencryptedSuffix,
			EncryptedSuffix:   conf.EncryptedSuffix,
			Version:           version.Version,
			ShamirThreshold:   conf.ShamirThreshold,
		}

		dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
		if len(errs) > 0 {
			err = fmt.Errorf("Could not generate data key: %s", errs)
			return err
		}

		err = common.EncryptTree(common.EncryptTreeOpts{
			DataKey: dataKey,
			Tree:    tree,
			Cipher:  opts.Cipher,
		})
		if err != nil {
			return err
		}

		fileContents, err = opts.InputStore.EmitEncryptedFile(*tree)
		if err != nil {
			return common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
		}
	} else {
		fileContents, err = ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read file: %s", err)
		}
	}

	if opts.Interactive {
		var response string
		for response != "y" && response != "n" {
			fmt.Printf("\nuploading %s to %s ? (y/n): ", path, conf.Destination.Path(fileName))
			_, err := fmt.Scanln(&response)
			if err != nil {
				return err
			}
		}
		if response == "n" {
			return errors.New("Publish canceled")
		}
	}

	err = conf.Destination.Upload(fileContents, fileName)
	if err != nil {
		return err
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}