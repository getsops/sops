package main

import (
	"fmt"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keyservice"
)

type setOpts struct {
	Cipher          sops.Cipher
	InputStore      sops.Store
	OutputStore     sops.Store
	InputPath       string
	IgnoreMAC       bool
	TreePath        []interface{}
	Value           interface{}
	KeyServices     []keyservice.KeyServiceClient
	DecryptionOrder []string
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
		Cipher:          opts.Cipher,
		IgnoreMac:       opts.IgnoreMAC,
		Tree:            tree,
		KeyServices:     opts.KeyServices,
		DecryptionOrder: opts.DecryptionOrder,
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
