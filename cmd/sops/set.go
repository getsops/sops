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
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
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
	key := opts.TreePath[len(opts.TreePath)-1]
	path := opts.TreePath[:len(opts.TreePath)-1]
	parent, err := tree.Branch.Truncate(path)
	if err != nil {
		return nil, common.NewExitError("Could not truncate tree to the provided path", codes.ErrorInvalidSetFormat)
	}
	branch := parent.(sops.TreeBranch)
	tree.Branch = branch.InsertOrReplaceValue(key, opts.Value)

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
	return encryptedFile, err
}
