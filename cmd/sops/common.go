package main

import (
	"fmt"
	"time"

	"io/ioutil"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/keyservice"
	"gopkg.in/urfave/cli.v1"
)

type decryptTreeOpts struct {
	Tree        *sops.Tree
	Stash       map[string][]interface{}
	KeyServices []keyservice.KeyServiceClient
	IgnoreMac   bool
	Cipher      sops.DataKeyCipher
}

func decryptTree(opts decryptTreeOpts) ([]byte, error) {
	dataKey, err := opts.Tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return nil, cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	computedMac, err := opts.Tree.Decrypt(dataKey, opts.Cipher, opts.Stash)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), exitErrorDecryptingTree)
	}
	fileMac, _, err := opts.Cipher.Decrypt(opts.Tree.Metadata.MessageAuthenticationCode, dataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339))
	if !opts.IgnoreMac {
		if fileMac != computedMac {
			// If the file has an empty MAC, display "no MAC" instead of not displaying anything
			if fileMac == "" {
				fileMac = "no MAC"
			}
			return nil, cli.NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", fileMac, computedMac), exitMacMismatch)
		}
	}
	return dataKey, nil
}

type encryptTreeOpts struct {
	Tree    *sops.Tree
	Stash   map[string][]interface{}
	Cipher  sops.DataKeyCipher
	DataKey []byte
}

func encryptTree(opts encryptTreeOpts) error {
	mac, err := opts.Tree.Encrypt(opts.DataKey, opts.Cipher, opts.Stash)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	opts.Tree.Metadata.LastModified = time.Now().UTC()
	mac, err = opts.Cipher.Encrypt(mac, opts.DataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339), nil)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), exitErrorEncryptingMac)
	}
	opts.Tree.Metadata.MessageAuthenticationCode = mac
	return nil
}

func loadEncryptedFile(inputStore sops.Store, inputPath string) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error reading file: %s", err), exitCouldNotReadInputFile)
	}
	metadata, err := inputStore.UnmarshalMetadata(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file metadata: %s", err), exitCouldNotReadInputFile)
	}
	branch, err := inputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree := sops.Tree{
		Branch:   branch,
		Metadata: *metadata,
	}
	return &tree, nil
}
