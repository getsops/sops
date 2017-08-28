package groups

import (
	"os"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
)

type DeleteOpts struct {
	InputPath   string
	InputStore  sops.Store
	OutputStore sops.Store
	Group       uint
	GroupQuorum uint
	InPlace     bool
	KeyServices []keyservice.KeyServiceClient
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Delete(opts DeleteOpts) error {
	tree, err := common.LoadEncryptedFile(opts.InputStore, opts.InputPath)
	if err != nil {
		return err
	}
	dataKey, err := tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return err
	}
	tree.Metadata.KeyGroups = append(tree.Metadata.KeyGroups[:opts.Group], tree.Metadata.KeyGroups[opts.Group+1:]...)

	if opts.GroupQuorum != 0 {
		tree.Metadata.ShamirQuorum = int(opts.GroupQuorum)
	}
	// The quorum should always be smaller or equal to the number of key groups
	tree.Metadata.ShamirQuorum = min(tree.Metadata.ShamirQuorum, len(tree.Metadata.KeyGroups))

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
