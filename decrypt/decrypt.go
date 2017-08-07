package decrypt // import "go.mozilla.org/sops/decrypt"

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/aes"
	sopsjson "go.mozilla.org/sops/json"
	sopsyaml "go.mozilla.org/sops/yaml"
)

// Constant indicating a MAC tag mismatch.
var errMacMismatch error = errors.New("Failed to verify data integrity.")

// File is a wrapper around Data that reads a local encrypted
// file and returns its cleartext data in an []byte
func File(path, format string) (cleartext []byte, err error) {
	// Read the file into an []byte
	encryptedData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read %q: %v", path, err)
	}
	return Data(encryptedData, format)
}

// Data is a helper that takes encrypted data and a format string,
// decrypts the data and returns its cleartext in an []byte.
// The format string can be `json`, `yaml` or `binary`.
// If the format string is empty, binary format is assumed.
func Data(data []byte, format string) (cleartext []byte, err error) {
	// Initialize a Sops JSON store
	var store sops.Store
	switch format {
	case "json":
		store = &sopsjson.Store{}
	case "yaml":
		store = &sopsyaml.Store{}
	default:
		store = &sopsjson.BinaryStore{}
	}
	// Load Sops metadata from the document and access the data key
	metadata, err := store.UnmarshalMetadata(data)
	if err != nil {
		return nil, err
	}
	key, err := metadata.GetDataKey()
	if err != nil {
		return nil, err
	}

	// Load the encrypted document and create a tree structure
	// with the encrypted content and metadata
	branch, err := store.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	tree := sops.Tree{Branch: branch, Metadata: metadata}

	// Decrypt the tree
	cipher := aes.Cipher{}
	stash := make(map[string][]interface{})
	mac, err := tree.Decrypt(key, cipher, stash)
	if err != nil {
		return nil, err
	}

	// Compute the hash of the cleartext tree and compare it with
	// the one that was stored in the document. If they match,
	// integrity was preserved
	originalMac, _, err := cipher.Decrypt(
		metadata.MessageAuthenticationCode,
		key,
		metadata.LastModified.Format(time.RFC3339),
	)

	computedMacBytes, err := hex.DecodeString(mac)
	if err != nil {
		return nil, err
	}
	originalMacBytes, err := hex.DecodeString(originalMac.(string))
	if err != nil {
		return nil, err
	}

	// Use a constant-time MAC tag comparison to avoid timing attacks: https://codahale.com/a-lesson-in-timing-attacks/
	if subtle.ConstantTimeCompare(computedMacBytes, originalMacBytes) != 1 {
		return nil, errMacMismatch
	}

	return store.Marshal(tree.Branch)
}
