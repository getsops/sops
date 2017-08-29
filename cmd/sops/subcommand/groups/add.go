package groups

import (
	"os"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
)

type AddOpts struct {
	InputPath   string
	InputStore  sops.Store
	OutputStore sops.Store
	Group       sops.KeyGroup
	GroupQuorum int
	InPlace     bool
	KeyServices []keyservice.KeyServiceClient
}

func Add(opts AddOpts) error {
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return err
	}
	dataKey, err := tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return err
	}
	tree.Metadata.KeyGroups = append(tree.Metadata.KeyGroups, opts.Group)

	if opts.GroupQuorum != 0 {
		tree.Metadata.ShamirQuorum = opts.GroupQuorum
	}
	tree.Metadata.UpdateMasterKeysWithKeyServices(dataKey, opts.KeyServices)
	output, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return err
	}
	var outputFile *os.File = os.Stdout
	if opts.InPlace {
		var err error
		outputFile, err = os.Create(opts.InputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()
	}
	outputFile.Write(output)
	return nil
}
