package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/stores"
	"github.com/getsops/sops/v3/version"
	"github.com/mitchellh/go-wordwrap"
)

type encryptConfig struct {
	UnencryptedSuffix string
	EncryptedSuffix   string
	UnencryptedRegex  string
	EncryptedRegex    string
	MACOnlyEncrypted  bool
	KeyGroups         []sops.KeyGroup
	GroupThreshold    int
}

type encryptOpts struct {
	Cipher      sops.Cipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	KeyServices []keyservice.KeyServiceClient
	encryptConfig
}

type fileAlreadyEncryptedError struct{}

func (err *fileAlreadyEncryptedError) Error() string {
	return "File already encrypted"
}

func (err *fileAlreadyEncryptedError) UserError() string {
	message := "The file you have provided contains a top-level entry called " +
		"'" + stores.SopsMetadataKey + "', or for flat file formats top-level entries starting with " +
		"'" + stores.SopsMetadataKey + "_'. This is generally due to the file already being encrypted. " +
		"SOPS uses a top-level entry called '" + stores.SopsMetadataKey + "' to store the metadata " +
		"required to decrypt the file. For this reason, SOPS can not " +
		"encrypt files that already contain such an entry.\n\n" +
		"If this is an unencrypted file, rename the '" + stores.SopsMetadataKey + "' entry.\n\n" +
		"If this is an encrypted file and you want to edit it, use the " +
		"editor mode, for example: `sops my_file.yaml`"
	return wordwrap.WrapString(message, 75)
}

func ensureNoMetadata(opts encryptOpts, branch sops.TreeBranch) error {
	if opts.OutputStore.HasSopsTopLevelKey(branch) {
		return &fileAlreadyEncryptedError{}
	}
	return nil
}

func metadataFromEncryptionConfig(config encryptConfig) sops.Metadata {
	return sops.Metadata{
		KeyGroups:         config.KeyGroups,
		UnencryptedSuffix: config.UnencryptedSuffix,
		EncryptedSuffix:   config.EncryptedSuffix,
		UnencryptedRegex:  config.UnencryptedRegex,
		EncryptedRegex:    config.EncryptedRegex,
		MACOnlyEncrypted:  config.MACOnlyEncrypted,
		Version:           version.Version,
		ShamirThreshold:   config.GroupThreshold,
	}
}

func encrypt(opts encryptOpts) (encryptedFile []byte, err error) {
	// Load the file
	fileBytes, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	branches, err := opts.InputStore.LoadPlainFile(fileBytes)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	if len(branches) < 1 {
		return nil, common.NewExitError("File cannot be completely empty, it must contain at least one document", codes.NeedAtLeastOneDocument)
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
		Metadata: metadataFromEncryptionConfig(opts.encryptConfig),
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
