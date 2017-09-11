/*
Package Sops manages JSON, YAML and BINARY documents to be encrypted or decrypted.

This package should not be used directly. Instead, Sops users should install the
command line client via `go get -u go.mozilla.org/sops/cmd/sops`, or use the
decryption helper provided at `go.mozilla.org/sops/decrypt`.

We do not guarantee API stability for any package other than `go.mozilla.org/sops/decrypt`.

A Sops document is a Tree composed of a data branch with arbitrary key/value pairs
and a metadata branch with encryption and integrity information.

In JSON and YAML formats, the structure of the cleartext tree is preserved, keys are
stored in cleartext and only values are encrypted. Keeping the values in cleartext
provides better readability when storing Sops documents in version controls, and allows
for merging competing changes on documents. This is a major difference between Sops
and other encryption tools that store documents as encrypted blobs.

In BINARY format, the cleartext data is treated as a single blob and the encrypted
document is in JSON format with a single `data` key and a single encrypted value.

Sops allows operators to encrypt their documents with multiple master keys. Each of
the master key defined in the document is able to decrypt it, allowing users to
share documents amongst themselves without sharing keys, or using a PGP key as a
backup for KMS.

In practice, this is achieved by generating a data key for each document that is used
to encrypt all values, and encrypting the data with each master key defined. Being
able to decrypt the data key gives access to the document.

The integrity of each document is guaranteed by calculating a Message Access Control
that is stored encrypted by the data key. When decrypting a document, the MAC should
be recalculated and compared with the MAC stored in the document to verify that no
fraudulent changes have been applied. The MAC covers keys and values as well as their
ordering.
*/
package sops //import "go.mozilla.org/sops"

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"log"

	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/shamir"
	"golang.org/x/net/context"
)

// DefaultUnencryptedSuffix is the default suffix a TreeItem key has to end with for sops to leave its Value unencrypted
const DefaultUnencryptedSuffix = "_unencrypted"

type sopsError string

func (e sopsError) Error() string {
	return string(e)
}

// MacMismatch occurs when the computed MAC does not match the expected ones
const MacMismatch = sopsError("MAC mismatch")

// MetadataNotFound occurs when the input file is malformed and doesn't have sops metadata in it
const MetadataNotFound = sopsError("sops metadata not found")

// DataKeyCipher provides a way to encrypt and decrypt the data key used to encrypt and decrypt sops files, so that the data key can be stored alongside the encrypted content. A DataKeyCipher must be able to decrypt the values it encrypts.
type DataKeyCipher interface {
	Encrypt(value interface{}, key []byte, additionalData string, stash interface{}) (string, error)
	Decrypt(value string, key []byte, additionalData string) (plaintext interface{}, stashValue interface{}, err error)
}

// Comment represents a comment in the sops tree for the file formats that actually support them.
type Comment struct {
	Value string
}

// TreeItem is an item inside sops's tree
type TreeItem struct {
	Key   interface{}
	Value interface{}
}

// TreeBranch is a branch inside sops's tree. It is a slice of TreeItems and is therefore ordered
type TreeBranch []TreeItem

// InsertOrReplaceValue replaces the value under the provided key with the newValue provided,
// or inserts a new key-value if it didn't exist already.
func (branch TreeBranch) InsertOrReplaceValue(key interface{}, newValue interface{}) TreeBranch {
	replaced := false
	for i, kv := range branch {
		if kv.Key == key {
			branch[i].Value = newValue
			replaced = true
			break
		}
	}
	if !replaced {
		return append(branch, TreeItem{Key: key, Value: newValue})
	}
	return branch
}

// Tree is the data structure used by sops to represent documents internally
type Tree struct {
	Branch   TreeBranch
	Metadata Metadata
}

// Truncate truncates the tree to the path specified
func (tree TreeBranch) Truncate(path []interface{}) (interface{}, error) {
	log.Printf("Truncating tree to %s", path)
	var current interface{} = tree
	for _, component := range path {
		switch component := component.(type) {
		case string:
			found := false
			for _, item := range current.(TreeBranch) {
				if item.Key == component {
					current = item.Value
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("component ['%s'] not found", component)
			}
		case int:
			if reflect.ValueOf(current).Kind() != reflect.Slice {
				return nil, fmt.Errorf("component [%d] is integer, but tree part is not a slice", component)
			}
			if reflect.ValueOf(current).Len() <= component {
				return nil, fmt.Errorf("component [%d] accesses out of bounds", component)
			}
			current = reflect.ValueOf(current).Index(component).Interface()
		}
	}
	return current, nil
}

func (tree TreeBranch) walkValue(in interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, path)
	case []byte:
		return onLeaves(string(in), path)
	case int:
		return onLeaves(in, path)
	case bool:
		return onLeaves(in, path)
	case float64:
		return onLeaves(in, path)
	case TreeBranch:
		return tree.walkBranch(in, path, onLeaves)
	case []interface{}:
		return tree.walkSlice(in, path, onLeaves)
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (tree TreeBranch) walkSlice(in []interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) ([]interface{}, error) {
	for i, v := range in {
		if _, ok := v.(Comment); ok {
			continue
		}
		newV, err := tree.walkValue(v, path, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
	}
	return in, nil
}

func (tree TreeBranch) walkBranch(in TreeBranch, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (TreeBranch, error) {
	for i, item := range in {
		if _, ok := item.Key.(Comment); ok {
			continue
		}
		key, ok := item.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Tree contains a non-string key (type %T): %s. Only string keys are"+
				"supported", item.Key, item.Key)
		}
		newV, err := tree.walkValue(item.Value, append(path, key), onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
	}
	return in, nil
}

// Encrypt walks over the tree and encrypts all values with the provided cipher, except those whose key ends with the UnencryptedSuffix specified on the Metadata struct. If encryption is successful, it returns the MAC for the encrypted tree.
func (tree Tree) Encrypt(key []byte, cipher DataKeyCipher, stash map[string][]interface{}) (string, error) {
	hash := sha512.New()
	_, err := tree.Branch.walkBranch(tree.Branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
		bytes, err := ToBytes(in)
		unencrypted := false
		for _, v := range path {
			if strings.HasSuffix(v, tree.Metadata.UnencryptedSuffix) {
				unencrypted = true
			}
		}
		if !unencrypted {
			var err error
			pathString := strings.Join(path, ":") + ":"
			// Pop from the left of the stash
			var stashValue interface{}
			if len(stash[pathString]) > 0 {
				var newStash []interface{}
				stashValue, newStash = stash[pathString][0], stash[pathString][1:len(stash[pathString])]
				stash[pathString] = newStash
			}
			in, err = cipher.Encrypt(in, key, pathString, stashValue)
			if err != nil {
				return nil, fmt.Errorf("Could not encrypt value: %s", err)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("Could not convert %s to bytes: %s", in, err)
		}
		hash.Write(bytes)
		return in, err
	})
	if err != nil {
		return "", fmt.Errorf("Error walking tree: %s", err)
	}
	return fmt.Sprintf("%X", hash.Sum(nil)), nil
}

// Decrypt walks over the tree and decrypts all values with the provided cipher, except those whose key ends with the UnencryptedSuffix specified on the Metadata struct. If decryption is successful, it returns the MAC for the decrypted tree.
func (tree Tree) Decrypt(key []byte, cipher DataKeyCipher, stash map[string][]interface{}) (string, error) {
	log.Print("Decrypting SOPS tree")
	hash := sha512.New()
	_, err := tree.Branch.walkBranch(tree.Branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
		var v interface{}
		unencrypted := false
		for _, v := range path {
			if strings.HasSuffix(v, tree.Metadata.UnencryptedSuffix) {
				unencrypted = true
			}
		}
		if !unencrypted {
			var err error
			var stashValue interface{}
			pathString := strings.Join(path, ":") + ":"
			v, stashValue, err = cipher.Decrypt(in.(string), key, pathString)
			if err != nil {
				return nil, fmt.Errorf("Could not decrypt value: %s", err)
			}
			stash[pathString] = append(stash[pathString], stashValue)
		} else {
			v = in
		}
		bytes, err := ToBytes(v)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %s to bytes: %s", in, err)
		}
		hash.Write(bytes)
		return v, err
	})
	if err != nil {
		return "", fmt.Errorf("Error walking tree: %s", err)
	}
	return fmt.Sprintf("%X", hash.Sum(nil)), nil

}

// GenerateDataKey generates a new random data key and encrypts it with all MasterKeys.
func (tree Tree) GenerateDataKey() ([]byte, []error) {
	newKey := make([]byte, 32)
	_, err := rand.Read(newKey)
	if err != nil {
		return nil, []error{fmt.Errorf("Could not generate random key: %s", err)}
	}
	return newKey, tree.Metadata.UpdateMasterKeys(newKey)
}

// GenerateDataKey generates a new random data key and encrypts it with all MasterKeys.
func (tree *Tree) GenerateDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient) ([]byte, []error) {
	newKey := make([]byte, 32)
	_, err := rand.Read(newKey)
	if err != nil {
		return nil, []error{fmt.Errorf("Could not generate random key: %s", err)}
	}
	return newKey, tree.Metadata.UpdateMasterKeysWithKeyServices(newKey, svcs)
}

// Metadata holds information about a file encrypted by sops
type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	MessageAuthenticationCode string
	Version                   string
	KeyGroups                 []KeyGroup
	// ShamirQuorum is the number of key groups required to recover the
	// original data key
	ShamirQuorum int
	// DataKey caches the decrypted data key so it doesn't have to be decrypted with a master key every time it's needed
	DataKey []byte
}

type KeyGroup []keys.MasterKey

// Store provides a way to load and save the sops tree along with metadata
type Store interface {
	Unmarshal(in []byte) (TreeBranch, error)
	UnmarshalMetadata(in []byte) (Metadata, error)
	Marshal(TreeBranch) ([]byte, error)
	MarshalWithMetadata(TreeBranch, Metadata) ([]byte, error)
	MarshalValue(interface{}) ([]byte, error)
}

// MasterKeyCount returns the number of master keys available
func (m *Metadata) MasterKeyCount() int {
	count := 0
	for _, group := range m.KeyGroups {
		count += len(group)
	}
	return count
}

func (m *Metadata) UpdateMasterKeysWithKeyServices(dataKey []byte, svcs []keyservice.KeyServiceClient) (errs []error) {
	if len(svcs) == 0 {
		return []error{
			fmt.Errorf("No key services provided, cansnot update master keys."),
		}
	}
	var parts [][]byte
	if len(m.KeyGroups) == 1 {
		// If there's only one key group, we can't do Shamir. All keys
		// in the group encrypt the whole data key.
		parts = append(parts, dataKey)
	} else {
		var err error
		if m.ShamirQuorum == 0 {
			m.ShamirQuorum = len(m.KeyGroups)
		}
		log.Printf("Multiple KeyGroups found, proceeding with Shamir with quorum %d", m.ShamirQuorum)
		parts, err = shamir.Split(dataKey, len(m.KeyGroups), m.ShamirQuorum)
		if err != nil {
			errs = append(errs, fmt.Errorf("Could not split data key into parts for Shamir: %s", err))
			return
		}
		if len(parts) != len(m.KeyGroups) {
			errs = append(errs, fmt.Errorf("Not enough parts obtained from Shamir. Need %d, got %d", len(m.KeyGroups), len(parts)))
			return
		}
	}
	for i, group := range m.KeyGroups {
		part := parts[i]
		for _, key := range group {
			svcKey := keyservice.KeyFromMasterKey(key)
			var keyErrs []error
			encrypted := false
			for _, svc := range svcs {
				rsp, err := svc.Encrypt(context.Background(), &keyservice.EncryptRequest{
					Key:       &svcKey,
					Plaintext: part,
				})
				if err != nil {
					keyErrs = append(keyErrs, fmt.Errorf("Failed to encrypt new data key with master key %q: %v\n", key.ToString(), err))
					continue
				}
				key.SetEncryptedDataKey(rsp.Ciphertext)
				encrypted = true
				// Only need to encrypt the key successfully with one service
				break
			}
			if !encrypted {
				errs = append(errs, keyErrs...)
			}
		}
	}
	m.DataKey = dataKey
	return
}

// UpdateMasterKeys encrypts the data key with all master keys
func (m *Metadata) UpdateMasterKeys(dataKey []byte) (errs []error) {
	return m.UpdateMasterKeysWithKeyServices(dataKey, []keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	})
}

// GetDataKeyWithKeyServices retrieves the data key, asking KeyServices to decrypt it with each
// MasterKey in the Metadata's KeySources until one of them succeeds.
func (m Metadata) GetDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient) ([]byte, error) {
	if m.DataKey != nil {
		return m.DataKey, nil
	}
	errMsg := "Could not decrypt the data key with any of the master keys:\n"
	var parts [][]byte
	for _, group := range m.KeyGroups {
	keysLoop:
		for _, key := range group {
			svcKey := keyservice.KeyFromMasterKey(key)
			for _, svc := range svcs {
				rsp, err := svc.Decrypt(
					context.Background(),
					&keyservice.DecryptRequest{
						Ciphertext: key.EncryptedDataKey(),
						Key:        &svcKey,
					})
				if err != nil {
					errMsg += fmt.Sprintf("\t%s: %s", key.ToString(), err)
					continue
				}
				parts = append(parts, rsp.Plaintext)
				// All keys in a key group encrypt the same part, so as soon
				// as we decrypt it successfully with one key, we need to
				// proceed with the next group
				break keysLoop
			}
		}
	}
	var dataKey []byte
	if len(m.KeyGroups) > 1 {
		if len(parts) < m.ShamirQuorum {
			return nil, fmt.Errorf("Not enough parts to recover data key with Shamir. Need %d, have %d.", m.ShamirQuorum, len(parts))
		}
		var err error
		dataKey, err = shamir.Combine(parts)
		if err != nil {
			return nil, fmt.Errorf("Could not get data key from shamir parts: %s", err)
		}
	} else {
		if len(parts) != 1 {
			return nil, fmt.Errorf("%s", errMsg)
		}
		dataKey = parts[0]
	}
	m.DataKey = dataKey
	return dataKey, nil
}

// GetDataKey retrieves the data key from the first MasterKey in the Metadata's KeySources that's able to return it,
// using the local KeyService
func (m Metadata) GetDataKey() ([]byte, error) {
	return m.GetDataKeyWithKeyServices([]keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	})
}

// ToBytes converts a string, int, float or bool to a byte representation.
func ToBytes(in interface{}) ([]byte, error) {
	switch in := in.(type) {
	case string:
		return []byte(in), nil
	case int:
		return []byte(strconv.Itoa(in)), nil
	case float64:
		return []byte(strconv.FormatFloat(in, 'f', -1, 64)), nil
	case bool:
		return []byte(strings.Title(strconv.FormatBool(in))), nil
	case []byte:
		return in, nil
	default:
		return nil, fmt.Errorf("Could not convert unknown type %T to bytes", in)
	}
}
