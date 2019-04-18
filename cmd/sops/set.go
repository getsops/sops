package main

import (
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
)

type setOpts struct {
	Cipher      sops.Cipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	IgnoreMAC   bool
	TreePath    []interface{}
	Value       interface{}
	KeyServices []keyservice.KeyServiceClient
}

func set(opts setOpts) ([]byte, error) {
	// Load the file
	// TODO: Issue #173: if the file does not exist, create it with the contents passed in as opts.Value
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
		Cipher:      opts.Cipher,
		IgnoreMac:   opts.IgnoreMAC,
		Tree:        tree,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	// Set the value
	tree.Branches[0] = tree.Branches[0].Set(opts.TreePath, opts.Value)

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
