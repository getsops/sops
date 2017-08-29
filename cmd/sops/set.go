package main

import (
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
	"gopkg.in/urfave/cli.v1"
)

type SetOpts struct {
	Cipher      sops.DataKeyCipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	IgnoreMAC   bool
	TreePath    []interface{}
	Value       interface{}
	KeyServices []keyservice.KeyServiceClient
}

func Set(opts SetOpts) ([]byte, error) {
	// Load the file
	// TODO: Issue #173: if the file does not exist, create it with the contents passed in as opts.Value
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	// Decrypt the file
	dataKey, err := common.DecryptTree(common.DecryptTreeOpts{
		Stash:       make(map[string][]interface{}),
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
		return nil, cli.NewExitError("Could not truncate tree to the provided path", codes.ErrorInvalidSetFormat)
	}
	branch := parent.(sops.TreeBranch)
	tree.Branch = branch.InsertOrReplaceValue(key, opts.Value)

	err = common.EncryptTree(common.EncryptTreeOpts{
		Stash: make(map[string][]interface{}), DataKey: dataKey, Tree: tree, Cipher: opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	encryptedFile, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return encryptedFile, err
}
