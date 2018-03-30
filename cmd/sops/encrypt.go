package main

import (
	"io/ioutil"

	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
)

type encryptOpts struct {
	Cipher            sops.Cipher
	InputStore        sops.Store
	OutputStore       sops.Store
	InputPath         string
	KeyServices       []keyservice.KeyServiceClient
	UnencryptedSuffix string
	KeyGroups         []sops.KeyGroup
	GroupThreshold    int
}

func encrypt(opts encryptOpts) (encryptedFile []byte, err error) {
	// Load the file
	fileBytes, err := ioutil.ReadFile(opts.InputPath)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	var tree sops.Tree
	branch, err := opts.InputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	tree.Branch = branch
	tree.Metadata = sops.Metadata{
		KeyGroups:         opts.KeyGroups,
		UnencryptedSuffix: opts.UnencryptedSuffix,
		Version:           version,
		ShamirThreshold:   opts.GroupThreshold,
	}
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	encryptedFile, err = opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return
}
