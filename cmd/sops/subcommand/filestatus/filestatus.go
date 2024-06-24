package filestatus

import (
	"fmt"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/common"
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
	encrypted, err := cfs(opts.InputStore, opts.InputPath)
	if err != nil {
		return Status{}, fmt.Errorf("cannot check file status: %w", err)
	}
	return Status{Encrypted: encrypted}, nil
}

// cfs checks and reports on file encryption status.
//
// It tries to decrypt the input file with the provided store.
// It returns true if the file contains sops metadata, false
// if it doesn't or Version or MessageAuthenticationCode are
// not found.
// It reports any error encountered different from
// sops.MetadataNotFound, as that is used to detect a sops
// encrypted file.
func cfs(s sops.Store, inputpath string) (bool, error) {
	tree, err := common.LoadEncryptedFile(s, inputpath)
	if err != nil && err == sops.MetadataNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cannot load encrypted file: %w", err)
	}

	// NOTE: even if it's a file that sops recognize as containing
	// valid metadata, we want to ensure some metadata are present
	// to report the file as encrypted.
	if tree.Metadata.Version == "" {
		return false, nil
	}
	if tree.Metadata.MessageAuthenticationCode == "" {
		return false, nil
	}

	return true, nil
}
