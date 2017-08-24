package main

import (
	"io/ioutil"

	"fmt"

	"time"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/keyservice"
	"gopkg.in/urfave/cli.v1"
)

type EncryptOpts struct {
	Cipher            sops.DataKeyCipher
	InputStore        sops.Store
	OutputStore       sops.Store
	InputPath         string
	KeyServices       []keyservice.KeyServiceClient
	UnencryptedSuffix string
	KeyGroups         []sops.KeyGroup
	GroupQuorum       uint
}

func Encrypt(opts EncryptOpts) (decryptedFile []byte, err error) {
	fileBytes, err := ioutil.ReadFile(opts.InputPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error reading file: %s", err), exitCouldNotReadInputFile)
	}
	var tree sops.Tree
	branch, err := opts.InputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), exitCouldNotReadInputFile)
	}
	tree.Branch = branch
	tree.Metadata = sops.Metadata{
		KeyGroups:         opts.KeyGroups,
		UnencryptedSuffix: opts.UnencryptedSuffix,
		Version:           version,
		ShamirQuorum:      int(opts.GroupQuorum),
	}
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
	decryptedFile, err = opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), exitErrorDumpingTree)
	}
	return
}
