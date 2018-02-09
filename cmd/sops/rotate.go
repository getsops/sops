package main

import (
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
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
}

func rotate(opts rotateOpts) ([]byte, error) {
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	dataKey, err := common.DecryptTree(common.DecryptTreeOpts{
		Cipher: opts.Cipher, IgnoreMac: opts.IgnoreMAC, Tree: tree,
		KeyServices: opts.KeyServices,
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

	encryptedFile, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return encryptedFile, nil
}
