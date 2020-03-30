/*
Package decrypt is the external API other Go programs can use to decrypt SOPS files. It is the only package in SOPS with
a stable API.
*/
package decrypt // import "go.mozilla.org/sops/v3/decrypt"

import (
	"fmt"
	"io/ioutil"
	"time"

	"go.mozilla.org/sops/v3/aes"
	"go.mozilla.org/sops/v3/cmd/sops/common"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/keyservice"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

// YamlFile opens a sops encrypted yaml file, parses it into a struct
// and returns an error if it fails. If a socket exists at /tmp/sops.sock,
// it will be used as a keyservice.
func YamlFile(path string, out interface{}) (err error) {
	// Read the file into an []byte
	encryptedData, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", path, err)
	}

	var svcs []keyservice.KeyServiceClient
	svcs = append(svcs, keyservice.NewLocalClient())
	// try connecting to unix:///tmp/sops.sock
	conn, err := grpc.Dial("unix:///tmp/sops.sock", []grpc.DialOption{grpc.WithInsecure()}...)
	if err == nil {
		// ignore errors but only add the keyservice if the dial call succeded
		svcs = append(svcs, keyservice.NewKeyServiceClient(conn))
	}

	store := common.StoreForFormat(formats.Yaml)

	// Load SOPS file and access the data key
	tree, err := store.LoadEncryptedFile(encryptedData)
	if err != nil {
		return err
	}
	key, err := tree.Metadata.GetDataKeyWithKeyServices(svcs)
	if err != nil {
		return err
	}

	// Decrypt the tree
	cipher := aes.NewCipher()
	mac, err := tree.Decrypt(key, cipher)
	if err != nil {
		return err
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
		return fmt.Errorf("Failed to verify data integrity. expected mac %q, got %q", originalMac, mac)
	}

	cleartext, err := store.EmitPlainFile(tree.Branches)
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}
	err = yaml.Unmarshal(cleartext, &out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cleartext into yaml: %w", err)
	}
	return nil
}

// File is a wrapper around Data that reads a local encrypted
// file and returns its cleartext data in an []byte
func File(path, format string) (cleartext []byte, err error) {
	// Read the file into an []byte
	encryptedData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read %q: %v", path, err)
	}

	// uses same logic as cli.
	formatFmt := formats.FormatForPathOrString(path, format)
	return DataWithFormat(encryptedData, formatFmt)
}

// DataWithFormat is a helper that takes encrypted data, and a format enum value,
// decrypts the data and returns its cleartext in an []byte.
func DataWithFormat(data []byte, format formats.Format) (cleartext []byte, err error) {

	store := common.StoreForFormat(format)

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

// Data is a helper that takes encrypted data and a format string,
// decrypts the data and returns its cleartext in an []byte.
// The format string can be `json`, `yaml`, `ini`, `dotenv` or `binary`.
// If the format string is empty, binary format is assumed.
func Data(data []byte, format string) (cleartext []byte, err error) {
	formatFmt := formats.FormatFromString(format)
	return DataWithFormat(data, formatFmt)
}
