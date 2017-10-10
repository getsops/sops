package main

import (
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
)

type decryptOpts struct {
	Cipher      sops.Cipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	IgnoreMAC   bool
	Extract     []interface{}
	KeyServices []keyservice.KeyServiceClient
}

func decrypt(opts decryptOpts) (decryptedFile []byte, err error) {
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	_, err = common.DecryptTree(common.DecryptTreeOpts{
		Cipher:      opts.Cipher,
		IgnoreMac:   opts.IgnoreMAC,
		Tree:        tree,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	if len(opts.Extract) > 0 {
		return extract(tree, opts.Extract, opts.OutputStore)
	}
	decryptedFile, err = opts.OutputStore.Marshal(tree.Branch)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error dumping file: %s", err), codes.ErrorDumpingTree)
	}
	return decryptedFile, err
}

func extract(tree *sops.Tree, path []interface{}, outputStore sops.Store) (output []byte, err error) {
	v, err := tree.Branch.Truncate(path)
	if err != nil {
		return nil, fmt.Errorf("error truncating tree: %s", err)
	}
	if newBranch, ok := v.(sops.TreeBranch); ok {
		tree.Branch = newBranch
		decrypted, err := outputStore.Marshal(tree.Branch)
		if err != nil {
			return nil, common.NewExitError(fmt.Sprintf("Error dumping file: %s", err), codes.ErrorDumpingTree)
		}
		return decrypted, err
	}
	bytes, err := outputStore.MarshalValue(v)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error dumping tree: %s", err), codes.ErrorDumpingTree)
	}
	return bytes, nil
}
