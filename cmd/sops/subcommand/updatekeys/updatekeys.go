package updatekeys

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/keyservice"
)

// Opts represents key operation options and config
type Opts struct {
	InputPath       string
	ShamirThreshold int
	KeyServices     []keyservice.KeyServiceClient
	DecryptionOrder []string
	Interactive     bool
	ConfigPath      string
	InputType       string
}

// UpdateKeys update the keys for a given file
func UpdateKeys(opts Opts) error {
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
	opts.InputPath = path
	return updateFile(opts)
}

func updateFile(opts Opts) error {
	sc, err := config.LoadStoresConfig(opts.ConfigPath)
	if err != nil {
		return err
	}
	store := common.DefaultStoreForPathOrFormat(sc, opts.InputPath, opts.InputType)
	log.Printf("Syncing keys for file %s", opts.InputPath)
	tree, err := common.LoadEncryptedFile(store, opts.InputPath)
	if err != nil {
		return err
	}
	conf, err := config.LoadCreationRuleForFile(opts.ConfigPath, opts.InputPath, make(map[string]*string))
	if err != nil {
		return err
	}
	if conf == nil {
		return fmt.Errorf("The config file %s does not contain any creation rule", opts.ConfigPath)
	}

	diffs := common.DiffKeyGroups(tree.Metadata.KeyGroups, conf.KeyGroups)
	keysWillChange := false
	for _, diff := range diffs {
		if len(diff.Added) > 0 || len(diff.Removed) > 0 {
			keysWillChange = true
		}
	}

	// TODO: use conf.ShamirThreshold instead of tree.Metadata.ShamirThreshold in the next line?
	//       Or make this configurable?
	var shamirThreshold = tree.Metadata.ShamirThreshold
	if opts.ShamirThreshold != 0 {
		shamirThreshold = opts.ShamirThreshold
	}
	shamirThreshold = min(shamirThreshold, len(conf.KeyGroups))
	var shamirThresholdWillChange = tree.Metadata.ShamirThreshold != shamirThreshold

	if !keysWillChange && !shamirThresholdWillChange {
		log.Printf("File %s already up to date", opts.InputPath)
		return nil
	}
	fmt.Printf("The following changes will be made to the file's groups:\n")
	common.PrettyPrintShamirDiff(tree.Metadata.ShamirThreshold, shamirThreshold)
	common.PrettyPrintDiffs(diffs)

	if opts.Interactive {
		var response string
		for response != "y" && response != "n" {
			fmt.Printf("Is this okay? (y/n):")
			_, err = fmt.Scanln(&response)
			if err != nil {
				return err
			}
		}
		if response == "n" {
			log.Printf("File %s left unchanged", opts.InputPath)
			return nil
		}
	}
	key, err := tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices, opts.DecryptionOrder)
	if err != nil {
		return common.NewExitError(err, codes.CouldNotRetrieveKey)
	}
	tree.Metadata.KeyGroups = conf.KeyGroups
	tree.Metadata.ShamirThreshold = shamirThreshold
	errs := tree.Metadata.UpdateMasterKeysWithKeyServices(key, opts.KeyServices)
	if len(errs) > 0 {
		return fmt.Errorf("error updating one or more master keys: %s", errs)
	}
	output, err := store.EmitEncryptedFile(*tree)
	if err != nil {
		return common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	outputFile, err := os.Create(opts.InputPath)
	if err != nil {
		return fmt.Errorf("could not open file for writing: %s", err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(output)
	if err != nil {
		return fmt.Errorf("error writing to file: %s", err)
	}
	log.Printf("File %s synced with new keys", opts.InputPath)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
