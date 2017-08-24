package main

import (
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/keyservice"
	"gopkg.in/urfave/cli.v1"
)

type DecryptOpts struct {
	Cipher      sops.DataKeyCipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	IgnoreMAC   bool
	Extract     []interface{}
	KeyServices []keyservice.KeyServiceClient
}

func Decrypt(opts DecryptOpts) (decryptedFile []byte, err error) {
	tree, err := loadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	_, err = decryptTree(decryptTreeOpts{
		Stash: make(map[string][]interface{}), Cipher: opts.Cipher, IgnoreMac: opts.IgnoreMAC, Tree: tree,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	if len(opts.Extract) > 0 {
		return Extract(tree, opts.Extract, opts.OutputStore)
	}
	decryptedFile, err = opts.OutputStore.Marshal(tree.Branch)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error dumping file: %s", err), exitErrorDumpingTree)
	}
	return decryptedFile, err
}

func Extract(tree *sops.Tree, path []interface{}, outputStore sops.Store) (output []byte, err error) {
	v, err := tree.Branch.Truncate(path)
	if err != nil {
		return nil, fmt.Errorf("error truncating tree: %s", err)
	}
	if newBranch, ok := v.(sops.TreeBranch); ok {
		tree.Branch = newBranch
		decrypted, err := outputStore.Marshal(tree.Branch)
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Error dumping file: %s", err), exitErrorDumpingTree)
		}
		return decrypted, err
	} else {
		bytes, err := outputStore.MarshalValue(v)
		if err != nil {
			return nil, cli.NewExitError(fmt.Sprintf("Error dumping tree: %s", err), exitErrorDumpingTree)
		}
		return bytes, nil
	}
}
