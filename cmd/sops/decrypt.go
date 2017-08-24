package main

import (
	"fmt"
	"io/ioutil"

	"time"

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
	if len(opts.Extract) > 0 {
		return Extract(tree, opts.Extract, opts.OutputStore)
	}
	decryptedFile, err = opts.OutputStore.Marshal(tree.Branch)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error dumping file: %s", err), exitErrorDumpingTree)
	}
	return decryptedFile, err
}

func Extract(tree sops.Tree, path []interface{}, outputStore sops.Store) (output []byte, err error) {
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
