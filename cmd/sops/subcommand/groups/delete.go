package groups

import (
	"os"

	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/common"
	"go.mozilla.org/sops/keyservice"
)

// DeleteOpts are the options for deleting a key group from a SOPS file
type DeleteOpts struct {
	InputPath      string
	InputStore     sops.Store
	OutputStore    sops.Store
	Group          uint
	GroupThreshold int
	InPlace        bool
	KeyServices    []keyservice.KeyServiceClient
}

// Delete deletes a key group from a SOPS file
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

	if opts.GroupThreshold != 0 {
		tree.Metadata.ShamirThreshold = opts.GroupThreshold
	}

	if len(tree.Metadata.KeyGroups) < tree.Metadata.ShamirThreshold {
		return fmt.Errorf("removing this key group will make the Shamir threshold impossible to satisfy: "+
			"Shamir threshold is %d, but we only have %d key groups", tree.Metadata.ShamirThreshold,
			len(tree.Metadata.KeyGroups))
	}

	tree.Metadata.UpdateMasterKeysWithKeyServices(dataKey, opts.KeyServices)
	output, err := opts.OutputStore.MarshalWithMetadata(tree.Branch, tree.Metadata)
	if err != nil {
		return err
	}
	var outputFile = os.Stdout
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
