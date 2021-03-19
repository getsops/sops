package filestatus

import (
	"fmt"
	"io/ioutil"

	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/cmd/sops/common"
)

// Opts represent the input options for FileStatus
type Opts struct {
	InputStore sops.Store
	InputPath  string
}

// Status represents the status of a file
type Status struct {
	// Encrypted represents whether the file provided is encrypted by SOPS
	Encrypted bool `json:"encrypted"`
}

// FileStatus checks encryption status of a file
func FileStatus(opts Opts) (Status, error) {
	encrypted, err := cfs(opts.InputPath)
	if err != nil {
		return Status{}, fmt.Errorf("cannot check file status: %w", err)
	}
	return Status{Encrypted: encrypted}, nil
}

func cfs(inputpath string) (bool, error) {
	fileBytes, err := ioutil.ReadFile(inputpath)
	if err != nil {
		return false, fmt.Errorf("cannot read input file: %w", err)
	}

	store := common.DefaultStoreForPath(inputpath)
	tree, err := store.LoadEncryptedFile(fileBytes)
	if err != nil && err == sops.MetadataNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cannot load encrypted file: %w", err)
	}

	if tree.Metadata.Version == "" {
		return false, nil
	}
	if tree.Metadata.MessageAuthenticationCode == "" {
		return false, nil
	}

	return true, nil
}
