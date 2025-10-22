package main

import (
	"fmt"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/audit"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
)

type rotateOpts struct {
	Cipher           sops.Cipher
	InputStore       sops.Store
	OutputStore      sops.Store
	InputPath        string
	IgnoreMAC        bool
	AddMasterKeys    []keys.MasterKey
	RemoveMasterKeys []keys.MasterKey
	KeyServices      []keyservice.KeyServiceClient
	DecryptionOrder  []string
}

func rotate(opts rotateOpts) ([]byte, error) {
	tree, err := common.LoadEncryptedFileWithBugFixes(common.GenericDecryptOpts{
		Cipher:          opts.Cipher,
		InputStore:      opts.InputStore,
		InputPath:       opts.InputPath,
		IgnoreMAC:       opts.IgnoreMAC,
		KeyServices:     opts.KeyServices,
		DecryptionOrder: opts.DecryptionOrder,
	})
	if err != nil {
		return nil, err
	}

	audit.SubmitEvent(audit.RotateEvent{
		File: tree.FilePath,
	})

	_, err = common.DecryptTree(common.DecryptTreeOpts{
		Cipher:          opts.Cipher,
		IgnoreMac:       opts.IgnoreMAC,
		Tree:            tree,
		KeyServices:     opts.KeyServices,
		DecryptionOrder: opts.DecryptionOrder,
	})
	if err != nil {
		return nil, err
	}

	// Add new master keys
	for _, key := range opts.AddMasterKeys {
		tree.Metadata.KeyGroups[0] = append(tree.Metadata.KeyGroups[0], key)
	}
	// Remove master keys
	for _, rmKey := range opts.RemoveMasterKeys {
		for i := range tree.Metadata.KeyGroups {
			for j, groupKey := range tree.Metadata.KeyGroups[i] {
				if rmKey.ToString() == groupKey.ToString() {
					tree.Metadata.KeyGroups[i] = append(tree.Metadata.KeyGroups[i][:j], tree.Metadata.KeyGroups[i][j+1:]...)
				}
			}
		}
	}

	// Create a new data key
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}

	// Reencrypt the file with the new key
	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey, Tree: tree, Cipher: opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	encryptedFile, err := opts.OutputStore.EmitEncryptedFile(*tree)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return encryptedFile, nil
}
