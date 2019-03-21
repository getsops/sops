package main

import (
	"io/ioutil"
	"path/filepath"

	"fmt"

	wordwrap "github.com/mitchellh/go-wordwrap"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/version"
)

type encryptOpts struct {
	Cipher            sops.Cipher
	InputStore        sops.Store
	OutputStore       sops.Store
	InputPath         string
	KeyServices       []keyservice.KeyServiceClient
	UnencryptedSuffix string
	EncryptedSuffix   string
	KeyGroups         []sops.KeyGroup
	GroupThreshold    int
}

type fileAlreadyEncryptedError struct{}

func (err *fileAlreadyEncryptedError) Error() string {
	return "File already encrypted"
}

func (err *fileAlreadyEncryptedError) UserError() string {
	message := "The file you have provided contains a top-level entry called " +
		"'sops'. This is generally due to the file already being encrypted. " +
		"SOPS uses a top-level entry called 'sops' to store the metadata " +
		"required to decrypt the file. For this reason, SOPS can not " +
		"encrypt files that already contain such an entry.\n\n" +
		"If this is an unencrypted file, rename the 'sops' entry.\n\n" +
		"If this is an encrypted file and you want to edit it, use the " +
		"editor mode, for example: `sops my_file.yaml`"
	return wordwrap.WrapString(message, 75)
}

func ensureNoMetadata(opts encryptOpts, branch sops.TreeBranch) error {
	for _, b := range branch {
		if b.Key == "sops" {
			return &fileAlreadyEncryptedError{}
		}
	}
	return nil
}

func encrypt(opts encryptOpts) (encryptedFile []byte, err error) {
	// Load the file
	fileBytes, err := ioutil.ReadFile(opts.InputPath)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	branches, err := opts.InputStore.LoadPlainFile(fileBytes)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	if err := ensureNoMetadata(opts, branches[0]); err != nil {
		return nil, common.NewExitError(err, codes.FileAlreadyEncrypted)
	}
	path, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return nil, err
	}
	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			KeyGroups:         opts.KeyGroups,
			UnencryptedSuffix: opts.UnencryptedSuffix,
			EncryptedSuffix:   opts.EncryptedSuffix,
			Version:           version.Version,
			ShamirThreshold:   opts.GroupThreshold,
		},
		FilePath: path,
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

	encryptedFile, err = opts.OutputStore.EmitEncryptedFile(tree)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return
}
