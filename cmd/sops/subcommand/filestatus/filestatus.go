package filestatus

import (
	"fmt"
	"io/ioutil"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/cmd/sops/common"
)

// Opts represent the input options for FileStatus
type Opts struct {
	InputStore sops.Store
	InputPath  string
}

// Status represents the status of a file
type Status struct {
	// Encrypted represents whether the file provided is encrypted by SOPS
	Encrypted bool
}

// FileStatus checks encryption status of a file
func FileStatus(opts Opts) (Status, error) {
	fileBytes, err := ioutil.ReadFile(opts.InputPath)
	if err != nil {
		return Status{}, common.NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	branches, err := opts.InputStore.LoadPlainFile(fileBytes)
	if err != nil {
		return Status{}, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	hasMetadata := checkMetadata((branches[0]))
	return Status{Encrypted: hasMetadata}, nil
}

func checkMetadata(branch sops.TreeBranch) bool {
	for _, b := range branch {
		if b.Key == "sops" {
			return true
		}
	}
	return false
}
