package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"go.mozilla.org/sops"
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
	fileBytes, err := ioutil.ReadFile(opts.InputPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error reading file: %s", err), exitCouldNotReadInputFile)
	}
	metadata, err := opts.InputStore.UnmarshalMetadata(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file metadata: %s", err), exitCouldNotReadInputFile)
	}
	branch, err := opts.InputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error loading file: %s", err), exitCouldNotReadInputFile)
	}
	tree := sops.Tree{
		Branch:   branch,
		Metadata: *metadata,
	}
	// Decrypt the file
	dataKey, err := tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return nil, cli.NewExitError(err.Error(), exitCouldNotRetrieveKey)
	}
	computedMac, err := tree.Decrypt(dataKey, opts.Cipher, make(map[string][]interface{}))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), exitErrorDecryptingTree)
	}
	fileMac, _, err := opts.Cipher.Decrypt(tree.Metadata.MessageAuthenticationCode, dataKey, tree.Metadata.LastModified.Format(time.RFC3339))
	if !opts.IgnoreMAC {
		if fileMac != computedMac {
			// If the file has an empty MAC, display "no MAC" instead of not displaying anything
			if fileMac == "" {
				fileMac = "no MAC"
			}
			return nil, cli.NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", fileMac, computedMac), exitMacMismatch)
		}
	}

	// Set the value
	key := opts.TreePath[len(opts.TreePath)-1]
	path := opts.TreePath[:len(opts.TreePath)-1]
	parent, err := tree.Branch.Truncate(path)
	if err != nil {
		return nil, cli.NewExitError("Could not truncate tree to the provided path", exitErrorInvalidSetFormat)
	}
	branch = parent.(sops.TreeBranch)
	tree.Branch = branch.InsertOrReplaceValue(key, opts.Value)

	// Encrypt the file
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}
	mac, err := tree.Encrypt(dataKey, opts.Cipher, make(map[string][]interface{}))
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), exitErrorEncryptingTree)
	}
	tree.Metadata.LastModified = time.Now().UTC()
	mac, err = opts.Cipher.Encrypt(mac, dataKey, tree.Metadata.LastModified.Format(time.RFC3339), nil)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), exitErrorEncryptingMac)
	}
	tree.Metadata.MessageAuthenticationCode = mac
	encryptedFile, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return encryptedFile, err
}
