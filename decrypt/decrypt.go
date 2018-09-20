/*
Package decrypt is the external API other Go programs can use to decrypt SOPS files. It is the only package in SOPS with
a stable API.
*/
package decrypt // import "go.mozilla.org/sops/decrypt"

import (
	"fmt"
	"io/ioutil"
	"time"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/aes"
	sopsdotenv "go.mozilla.org/sops/stores/dotenv"
	sopsjson "go.mozilla.org/sops/stores/json"
	sopsyaml "go.mozilla.org/sops/stores/yaml"
)

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
// The format string can be `json`, `yaml`, `dotenv` or `binary`.
// If the format string is empty, binary format is assumed.
func Data(data []byte, format string) (cleartext []byte, err error) {
	// Initialize a Sops JSON store
	var store sops.Store
	switch format {
	case "json":
		store = &sopsjson.Store{}
	case "yaml":
		store = &sopsyaml.Store{}
	case "dotenv":
		store = &sopsdotenv.Store{}
	default:
		store = &sopsjson.BinaryStore{}
	}
	// Load SOPS file and access the data key
	tree, err := store.LoadEncryptedFile(data)
	if err != nil {
		return nil, err
	}
	key, err := tree.Metadata.GetDataKey()
	if err != nil {
		return nil, err
	}

	// Decrypt the tree
	cipher := aes.NewCipher()
	mac, err := tree.Decrypt(key, cipher)
	if err != nil {
		return nil, err
	}

	// Compute the hash of the cleartext tree and compare it with
	// the one that was stored in the document. If they match,
	// integrity was preserved
	originalMac, err := cipher.Decrypt(
		tree.Metadata.MessageAuthenticationCode,
		key,
		tree.Metadata.LastModified.Format(time.RFC3339),
	)
	if originalMac != mac {
		return nil, fmt.Errorf("Failed to verify data integrity. expected mac %q, got %q", originalMac, mac)
	}

	return store.EmitPlainFile(tree.Branches)
}
