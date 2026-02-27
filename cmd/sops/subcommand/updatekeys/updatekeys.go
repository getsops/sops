package updatekeys

import (
	"bytes"
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
	Global          bool // apply updatekey to all managed files
	DryRun          bool // do not modify files in global mode, only show intended changes
}

// UpdateKeys update the keys for a given file
func UpdateKeys(opts Opts) error {
	if opts.Global {
		return updateAll(opts)
	}
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

func updateAll(opts Opts) error {
    // Root scoped to config file directory or current working directory
    root := "."
    if opts.ConfigPath != "" {
        root = filepath.Dir(opts.ConfigPath)
    }
    absRoot, err := filepath.Abs(root)
    if err != nil {
        return err
    }

    log.Printf("Global updatekeys: scanning %s", absRoot)

    var updated, skipped int
    var errs []error
    var filesToUpdate []string

    err = filepath.Walk(absRoot, func(p string, info os.FileInfo, walkErr error) error {
        if walkErr != nil {
            errs = append(errs, walkErr)
            return nil
        }
        if info.IsDir() {
            // skip common large/irrelevant dirs
            base := filepath.Base(p)
            if base == ".git" || base == "vendor" || base == ".idea" || base == "node_modules" {
                return filepath.SkipDir
            }
            return nil
        }

        // Skip the config file itself
        if filepath.Base(p) == ".sops.yaml" || filepath.Base(p) == ".sops.yml" {
            skipped++
            return nil
        }

        // Determine if this file is a SOPS-managed file (contains SOPS metadata); if not, skip.
        data, rerr := os.ReadFile(p)
        if rerr != nil {
            errs = append(errs, fmt.Errorf("read failed for %s: %w", p, rerr))
            return nil
        }

        // Heuristic: look for common SOPS metadata markers, this could be better?
        hasMeta := bytes.Contains(data, []byte("sops:")) || bytes.Contains(data, []byte(`"sops"`))
        if !hasMeta {
            skipped++
            return nil
        }

        // Determine if this file has a creation rule; if not, skip
        conf, cerr := config.LoadCreationRuleForFile(opts.ConfigPath, p, make(map[string]*string))
        if cerr != nil || conf == nil {
            log.Printf("Ignoring file %s: no matching creation rule", p)
            skipped++
            return nil
        }
        fileOpts := opts
        fileOpts.InputPath = p
        if opts.DryRun {
            would, werr := wouldUpdate(fileOpts)
            if werr != nil {
                errs = append(errs, fmt.Errorf("check failed for %s: %w", p, werr))
                return nil
            }
            if would {
                filesToUpdate = append(filesToUpdate, p)
            }
        } else {
            if uErr := updateFile(fileOpts); uErr != nil {
                errs = append(errs, fmt.Errorf("update failed for %s: %w", p, uErr))
            } else {
                updated++
            }
        }
        return nil
    })
    if err != nil {
        errs = append(errs, err)
    }

    if opts.DryRun {
        log.Printf("Global dry-run updatekeys complete: would update %d files, skipped %d, errors %d", len(filesToUpdate), skipped, len(errs))
        if len(filesToUpdate) > 0 {
            fmt.Printf("Files that would be updated:\n")
            for _, f := range filesToUpdate {
                fmt.Printf("  %s\n", f)
            }
        }
    } else {
        log.Printf("Global updatekeys complete: updated=%d skipped=%d errors=%d", updated, skipped, len(errs))
    }
    if len(errs) > 0 {
        return fmt.Errorf("global updatekeys finished with errors: first=%v (total %d)", errs[0], len(errs))
    }
    return nil
}

func wouldUpdate(opts Opts) (bool, error) {
    sc, err := config.LoadStoresConfig(opts.ConfigPath)
    if err != nil {
        return false, err
    }
    store := common.DefaultStoreForPathOrFormat(sc, opts.InputPath, opts.InputType)
    tree, err := common.LoadEncryptedFile(store, opts.InputPath)
    if err != nil {
        return false, err
    }
    conf, err := config.LoadCreationRuleForFile(opts.ConfigPath, opts.InputPath, make(map[string]*string))
    if err != nil {
        return false, err
    }
    if conf == nil {
        return false, fmt.Errorf("The config file %s does not contain any creation rule", opts.ConfigPath)
    }

    diffs := common.DiffKeyGroups(tree.Metadata.KeyGroups, conf.KeyGroups)
    keysWillChange := false
    for _, diff := range diffs {
        if len(diff.Added) > 0 || len(diff.Removed) > 0 {
            keysWillChange = true
        }
    }

    var shamirThreshold = tree.Metadata.ShamirThreshold
    if opts.ShamirThreshold != 0 {
        shamirThreshold = opts.ShamirThreshold
    }
    shamirThreshold = min(shamirThreshold, len(conf.KeyGroups))
    shamirThresholdWillChange := tree.Metadata.ShamirThreshold != shamirThreshold

    return keysWillChange || shamirThresholdWillChange, nil
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
		if opts.DryRun {
			log.Printf("[dry-run] File %s already up to date", opts.InputPath)
			return nil
		}
		log.Printf("File %s already up to date", opts.InputPath)
		return nil
	}
	fmt.Printf("The following changes will be made to the file's groups:\n")
	common.PrettyPrintShamirDiff(tree.Metadata.ShamirThreshold, shamirThreshold)
	common.PrettyPrintDiffs(diffs)

	if opts.DryRun {
		log.Printf("[dry-run] Would update file %s (no changes written)", opts.InputPath)
		return nil
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
