package common

import (
	"fmt"
	"time"

	"io/ioutil"

	"strings"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/cmd/sops/codes"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/stores/json"
	"go.mozilla.org/sops/stores/yaml"
	"gopkg.in/urfave/cli.v1"
)

// DecryptTreeOpts are the options needed to decrypt a tree
type DecryptTreeOpts struct {
	// Tree is the tree to be decrypted
	Tree *sops.Tree
	// KeyServices are the key services to be used for decryption of the data key
	KeyServices []keyservice.KeyServiceClient
	// IgnoreMac is whether or not to ignore the Message Authentication Code included in the SOPS tree
	IgnoreMac bool
	// Cipher is the cryptographic cipher to use to decrypt the values inside the tree
	Cipher sops.Cipher
}

// DecryptTree decrypts the tree passed in through the DecryptTreeOpts and additionally returns the decrypted data key
func DecryptTree(opts DecryptTreeOpts) (dataKey []byte, err error) {
	dataKey, err = opts.Tree.Metadata.GetDataKeyWithKeyServices(opts.KeyServices)
	if err != nil {
		return nil, NewExitError(err, codes.CouldNotRetrieveKey)
	}
	computedMac, err := opts.Tree.Decrypt(dataKey, opts.Cipher)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Error decrypting tree: %s", err), codes.ErrorDecryptingTree)
	}
	fileMac, err := opts.Cipher.Decrypt(opts.Tree.Metadata.MessageAuthenticationCode, dataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339))
	if !opts.IgnoreMac {
		if fileMac != computedMac {
			// If the file has an empty MAC, display "no MAC" instead of not displaying anything
			if fileMac == "" {
				fileMac = "no MAC"
			}
			return nil, NewExitError(fmt.Sprintf("MAC mismatch. File has %s, computed %s", fileMac, computedMac), codes.MacMismatch)
		}
	}
	return dataKey, nil
}

// EncryptTreeOpts are the options needed to encrypt a tree
type EncryptTreeOpts struct {
	// Tree is the tree to be encrypted
	Tree *sops.Tree
	// Cipher is the cryptographic cipher to use to encrypt the values inside the tree
	Cipher sops.Cipher
	// DataKey is the key the cipher should use to encrypt the values inside the tree
	DataKey []byte
}

// EncryptTree encrypts the tree passed in through the EncryptTreeOpts
func EncryptTree(opts EncryptTreeOpts) error {
	unencryptedMac, err := opts.Tree.Encrypt(opts.DataKey, opts.Cipher)
	if err != nil {
		return NewExitError(fmt.Sprintf("Error encrypting tree: %s", err), codes.ErrorEncryptingTree)
	}
	opts.Tree.Metadata.LastModified = time.Now().UTC()
	opts.Tree.Metadata.MessageAuthenticationCode, err = opts.Cipher.Encrypt(unencryptedMac, opts.DataKey, opts.Tree.Metadata.LastModified.Format(time.RFC3339))
	if err != nil {
		return NewExitError(fmt.Sprintf("Could not encrypt MAC: %s", err), codes.ErrorEncryptingMac)
	}
	return nil
}

// LoadEncryptedFile loads an encrypted SOPS file, returning a SOPS tree
func LoadEncryptedFile(inputStore sops.Store, inputPath string) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	metadata, err := inputStore.UnmarshalMetadata(fileBytes)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Error loading file metadata: %s", err), codes.CouldNotReadInputFile)
	}
	branch, err := inputStore.Unmarshal(fileBytes)
	if err != nil {
		return nil, NewExitError(fmt.Sprintf("Error loading file: %s", err), codes.CouldNotReadInputFile)
	}
	tree := sops.Tree{
		Branch:   branch,
		Metadata: metadata,
	}
	return &tree, nil
}

func NewExitError(i interface{}, exitCode int) *cli.ExitError {
	if userErr, ok := i.(sops.UserError); ok {
		return NewExitError(userErr.UserError(), exitCode)
	}
	return cli.NewExitError(i, exitCode)
}

func IsYAMLFile(path string) bool {
	return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
}

func IsJSONFile(path string) bool {
	return strings.HasSuffix(path, ".json")
}

func DefaultStoreForPath(path string) sops.Store {
	if IsYAMLFile(path) {
		return &yaml.Store{}
	} else if IsJSONFile(path) {
		return &json.Store{}
	}
	return &json.BinaryStore{}
}
