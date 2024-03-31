package main

import (
	"fmt"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keyservice"
)

type unsetOpts struct {
	Cipher          sops.Cipher
	InputStore      sops.Store
	OutputStore     sops.Store
	InputPath       string
	IgnoreMAC       bool
	TreePath        []interface{}
	KeyServices     []keyservice.KeyServiceClient
	DecryptionOrder []string
}

func unset(opts unsetOpts) ([]byte, error) {
	// Load the file
	tree, err := common.LoadEncryptedFileWithBugFixes(common.GenericDecryptOpts{
		Cipher:      opts.Cipher,
		InputStore:  opts.InputStore,
		InputPath:   opts.InputPath,
		IgnoreMAC:   opts.IgnoreMAC,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	// Decrypt the file
	dataKey, err := common.DecryptTree(common.DecryptTreeOpts{
		Cipher:          opts.Cipher,
		IgnoreMac:       opts.IgnoreMAC,
		Tree:            tree,
		KeyServices:     opts.KeyServices,
		DecryptionOrder: opts.DecryptionOrder,
	})
	if err != nil {
		return nil, err
	}

	// Unset the value
	newBranch, err := tree.Branches[0].Unset(opts.TreePath)
	if err != nil {
		return nil, err
	}
	tree.Branches[0] = newBranch

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
	return encryptedFile, err
}
