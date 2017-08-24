package main

import (
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
	cli "gopkg.in/urfave/cli.v1"
)

type RotateOpts struct {
	Cipher        sops.DataKeyCipher
	InputStore    sops.Store
	OutputStore   sops.Store
	InputPath     string
	IgnoreMAC     bool
	AddMasterKeys []struct {
		Key     keys.MasterKey
		ToGroup uint
	}
	RemoveMasterKeys []struct {
		Key       keys.MasterKey
		FromGroup uint
	}
	KeyServices []keyservice.KeyServiceClient
}

func Rotate(opts RotateOpts) ([]byte, error) {
	tree, err := loadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return nil, err
	}

	dataKey, err := decryptTree(decryptTreeOpts{
		Stash: make(map[string][]interface{}), Cipher: opts.Cipher, IgnoreMac: opts.IgnoreMAC, Tree: tree,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Add and remove master keys
	// Create a new data key
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}

	// Reencrypt the file with the new key
	err = encryptTree(encryptTreeOpts{
		Stash: make(map[string][]interface{}), DataKey: dataKey, Tree: tree, Cipher: opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	encryptedFile, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return encryptedFile, nil
}
