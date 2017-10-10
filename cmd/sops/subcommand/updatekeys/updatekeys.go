package updatekeys

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/config"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
)

type Opts struct {
	InputPath   string
	GroupQuorum int
	KeyServices []keyservice.KeyServiceClient
	Interactive bool
	ConfigPath  string
}

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
	return updateFile(opts)
}

func updateFile(opts Opts) error {
	store := common.DefaultStoreForPath(opts.InputPath)
	log.Printf("Syncing keys for file %s", opts.InputPath)
	tree, err := common.LoadEncryptedFile(store, opts.InputPath)
	if err != nil {
		return err
	}
	conf, err := config.LoadForFile(opts.ConfigPath, opts.InputPath, make(map[string]*string))
	if err != nil {
		return err
	}
	diffs := diffKeyGroups(tree.Metadata.KeyGroups, conf.KeyGroups)
	keysWillChange := false
	for _, diff := range diffs {
		if len(diff.added) > 0 || len(diff.removed) > 0 {
			keysWillChange = true
		}
	}
	if !keysWillChange {
		log.Printf("File %s already up to date", opts.InputPath)
		return nil
	}
	fmt.Printf("The following changes will be made to the file's groups:\n")
	for i, diff := range diffs {
		color.New(color.Underline).Printf("Group %d\n", i+1)
		for _, c := range diff.common {
			fmt.Printf("    %s\n", c.ToString())
		}
		for _, c := range diff.added {
			color.New(color.FgGreen).Printf("+++ %s\n", c.ToString())
		}
		for _, c := range diff.removed {
			color.New(color.FgRed).Printf("--- %s\n", c.ToString())
		}
	}
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
	key, err := tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return fmt.Errorf("error getting data key: %s", err)
	}
	tree.Metadata.KeyGroups = conf.KeyGroups
	if opts.GroupQuorum != 0 {
		tree.Metadata.ShamirThreshold = opts.GroupQuorum
	}
	tree.Metadata.ShamirThreshold = min(tree.Metadata.ShamirThreshold, len(tree.Metadata.KeyGroups))
	errs := tree.Metadata.UpdateMasterKeysWithKeyServices(key, opts.KeyServices)
	if len(errs) > 0 {
		return fmt.Errorf("error updating one or more master keys: %s", errs)
	}
	output, err := store.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return fmt.Errorf("error marshaling tree: %s", err)
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

type diff struct {
	common  []keys.MasterKey
	added   []keys.MasterKey
	removed []keys.MasterKey
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func diffKeyGroups(ours, theirs []sops.KeyGroup) []diff {
	var diffs []diff
	for i := 0; i < max(len(ours), len(theirs)); i++ {
		var diff diff
		var ourGroup, theirGroup sops.KeyGroup
		if len(ours) > i {
			ourGroup = ours[i]
		}
		if len(theirs) > i {
			theirGroup = theirs[i]
		}
		ourKeys := make(map[string]struct{})
		theirKeys := make(map[string]struct{})
		for _, key := range ourGroup {
			ourKeys[key.ToString()] = struct{}{}
		}
		for _, key := range theirGroup {
			if _, ok := ourKeys[key.ToString()]; ok {
				diff.common = append(diff.common, key)
			} else {
				diff.added = append(diff.added, key)
			}
			theirKeys[key.ToString()] = struct{}{}
		}
		for _, key := range ourGroup {
			if _, ok := theirKeys[key.ToString()]; !ok {
				diff.removed = append(diff.removed, key)
			}
		}
		diffs = append(diffs, diff)
	}
	return diffs
}
