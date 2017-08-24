package main

import (
	"fmt"
	"io/ioutil"

	"time"

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
	// Load the file
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
	// TODO: Add and remove master keys
	// Create a new data key
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}
	// Reencrypt the file with the new key
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
	return encryptedFile, nil
}
