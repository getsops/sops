package common

import (
	"fmt"
	"time"

	"io/ioutil"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/keyservice"
	"gopkg.in/urfave/cli.v1"
)

type DecryptTreeOpts struct {
	Tree        *sops.Tree
	Stash       map[string][]interface{}
	KeyServices []keyservice.KeyServiceClient
	IgnoreMac   bool
	Cipher      sops.DataKeyCipher
}

func DecryptTree(opts DecryptTreeOpts) ([]byte, error) {
	dataKey, err := opts.Tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return nil, cli.NewExitError(err.Error(), codes.CouldNotRetrieveKey)
	}
	computedMac, err := opts.Tree.Decrypt(dataKey, opts.Cipher, opts.Stash)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), codes.ErrorDecryptingTree)
	}
	fileMac, _, err := opts.Cipher.Decrypt(opts.Tree.Metadata.MessageAuthenticationCode, dataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339))
	if !opts.IgnoreMac {
		if fileMac != computedMac {
			// If the file has an empty MAC, display "no MAC" instead of not displaying anything
			if fileMac == "" {
				fileMac = "no MAC"
			}
			return nil, cli.NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", fileMac, computedMac), codes.MacMismatch)
		}
	}
	return dataKey, nil
}

type EncryptTreeOpts struct {
	Tree    *sops.Tree
	Stash   map[string][]interface{}
	Cipher  sops.DataKeyCipher
	DataKey []byte
}

func EncryptTree(opts EncryptTreeOpts) error {
	mac, err := opts.Tree.Encrypt(opts.DataKey, opts.Cipher, opts.Stash)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), codes.ErrorEncryptingTree)
	}
	opts.Tree.Metadata.LastModified = time.Now().UTC()
	mac, err = opts.Cipher.Encrypt(mac, opts.DataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339), opts.Stash)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), codes.ErrorEncryptingMac)
	}
	opts.Tree.Metadata.MessageAuthenticationCode = mac
	return nil
}

func LoadEncryptedFile(inputStore sops.Store, inputPath string) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	metadata, err := inputStore.UnmarshalMetadata(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file metadata: %s", err), codes.CouldNotReadInputFile)
	}
	branch, err := inputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), codes.CouldNotReadInputFile)
	}
	tree := sops.Tree{
		Branch:   branch,
		Metadata: metadata,
	}
	return &tree, nil
}
