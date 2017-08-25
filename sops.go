/*
Package Sops manages JSON, YAML and BINARY documents to be encrypted or decrypted.

This package should not be used directly. Instead, Sops users should install the
command line client via `go get -u go.mozilla.org/sops/cmd/sops`, or use the
decryption helper provided at `go.mozilla.org/sops/decrypt`.

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

	"go.mozilla.org/sops/shamir"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
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
	Encrypt(value interface{}, key []byte, path string, stash interface{}) (string, error)
	Decrypt(value string, key []byte, path string) (plaintext interface{}, stashValue interface{}, err error)
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

// TrimTreePathComponent trimps a tree path component so that it's a valid tree key
func TrimTreePathComponent(component string) (string, error) {
	if component[len(component) - 1] != ']' {
		return "", fmt.Errorf("Invalid component")
	}
	component = component[:len(component) - 1]
	component = strings.Replace(component, `"`, "", 2)
	component = strings.Replace(component, `'`, "", 2)
	return component, nil
}

// Truncate truncates the tree following Python dictionary access syntax, for example, ["foo"][2].
func (tree TreeBranch) Truncate(path string) (interface{}, error) {
	components := strings.Split(path, "[")
	var current interface{} = tree
	for _, component := range components {
		if component == "" {
			continue
		}
		component, err := TrimTreePathComponent(component)
		if err != nil {
			return nil, fmt.Errorf("Invalid tree path format string: %s", path)
		}
		i, err := strconv.Atoi(component)
		if err != nil {
			for _, item := range current.(TreeBranch) {
				if item.Key == component {
					current = item.Value
					break
				}
			}
		} else {
			v := reflect.ValueOf(current)
			current = v.Index(i).Interface()
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
			return nil, fmt.Errorf("Tree contains a non-string key (type %T): %s. Only string keys are" +
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

// Metadata holds information about a file encrypted by sops
type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	MessageAuthenticationCode string
	Version                   string
	KeySources                []KeySource
	// Shamir is true when the data key is split across multiple master keys
	// according to shamir's secret sharing algorithm
	Shamir bool
	// ShamirQuorum is the number of master keys required to recover the
	// original data key
	ShamirQuorum int
}

// KeySource is a collection of MasterKeys with a Name.
type KeySource struct {
	Name string
	Keys []MasterKey
}

// MasterKey provides a way of securing the key used to encrypt the Tree by encrypting and decrypting said key.
type MasterKey interface {
	Encrypt(dataKey []byte) error
	EncryptIfNeeded(dataKey []byte) error
	Decrypt() ([]byte, error)
	NeedsRotation() bool
	ToString() string
	ToMap() map[string]interface{}
}

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
	for _, ks := range m.KeySources {
		count += len(ks.Keys)
	}
	return count
}

// RemoveMasterKeys removes all of the provided keys from the metadata's KeySources, if they exist there.
func (m *Metadata) RemoveMasterKeys(keys []MasterKey) {
	for j, ks := range m.KeySources {
		var newKeys []MasterKey
		for _, k := range ks.Keys {
			matchFound := false
			for _, keyToRemove := range keys {
				if k.ToString() == keyToRemove.ToString() {
					matchFound = true
					break
				}
			}
			if !matchFound {
				newKeys = append(newKeys, k)
			}
		}
		m.KeySources[j].Keys = newKeys
	}
}

// UpdateMasterKeysIfNeeded encrypts the data key with all master keys if it's needed
func (m *Metadata) UpdateMasterKeysIfNeeded(dataKey []byte) (errs []error) {
	// If we're using Shamir and we've added or removed keys, we must
	// generate Shamir parts again and reencrypt with all keys
	if m.Shamir {
		return m.updateMasterKeysShamir(dataKey)
	}
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := k.EncryptIfNeeded(dataKey)
			if err != nil {
				errs = append(errs, fmt.Errorf("Failed to encrypt new data key with master key %q: %v\n", k.ToString(), err))
			}
		}
	}
	return
}

// updateMasterKeysShamir splits the data key into parts using Shamir's Secret
// Sharing algorithm and encrypts each part with a master key
func (m *Metadata) updateMasterKeysShamir(dataKey []byte) (errs []error) {
	keyCount := 0
	for _, ks := range m.KeySources {
		for range ks.Keys {
			keyCount++
		}
	}
	// If the quorum wasn't set, default to 2
	if m.ShamirQuorum == 0 {
		m.ShamirQuorum = 2
	}
	parts, err := shamir.Split(dataKey, keyCount, m.ShamirQuorum)
	if err != nil {
		errs = append(errs, fmt.Errorf("Could not split data key into parts for Shamir: %s", err))
		return
	}
	if len(parts) != keyCount {
		errs = append(errs, fmt.Errorf("Not enough parts obtained from Shamir. Need %d, got %d", keyCount, len(parts)))
		return
	}
	counter := 0
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := k.Encrypt(parts[counter])
			if err != nil {
				errs = append(errs, fmt.Errorf("Failed to encrypt Shamir part with master key %q: %v\n", k.ToString(), err))
			}
			counter++
		}
	}
	return
}

// UpdateMasterKeys encrypts the data key with all master keys
func (m *Metadata) UpdateMasterKeys(dataKey []byte) (errs []error) {
	if m.Shamir {
		return m.updateMasterKeysShamir(dataKey)
	}
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := k.Encrypt(dataKey)
			if err != nil {
				errs = append(errs, fmt.Errorf("Failed to encrypt new data key with master key %q: %v\n", k.ToString(), err))
			}
		}
	}
	return
}

// AddPGPMasterKeys parses the input comma separated string of GPG fingerprints, generates a PGP MasterKey for each fingerprint, and adds the keys to the PGP KeySource
func (m *Metadata) AddPGPMasterKeys(pgpFps string) {
	for i, ks := range m.KeySources {
		if ks.Name == "pgp" {
			var keys []MasterKey
			for _, k := range pgp.MasterKeysFromFingerprintString(pgpFps) {
				keys = append(keys, k)
				fmt.Printf("Adding new PGP master key: %X\n", k.Fingerprint)
			}
			ks.Keys = append(ks.Keys, keys...)
			m.KeySources[i] = ks
		}
	}
}

// AddKMSMasterKeys parses the input comma separated string of AWS KMS ARNs, generates a KMS MasterKey for each ARN, and then adds the keys to the KMS KeySource
func (m *Metadata) AddKMSMasterKeys(kmsArns string, context map[string]*string) {
	for i, ks := range m.KeySources {
		if ks.Name == "kms" {
			var keys []MasterKey
			for _, k := range kms.MasterKeysFromArnString(kmsArns, context) {
				keys = append(keys, k)
				fmt.Printf("Adding new KMS master key: %s\n", k.Arn)
			}
			ks.Keys = append(ks.Keys, keys...)
			m.KeySources[i] = ks
		}
	}
}

// RemovePGPMasterKeys takes a comma separated string of PGP fingerprints and removes the keys corresponding to those fingerprints from the metadata's KeySources
func (m *Metadata) RemovePGPMasterKeys(pgpFps string) {
	var keys []MasterKey
	for _, k := range pgp.MasterKeysFromFingerprintString(pgpFps) {
		keys = append(keys, k)
	}
	m.RemoveMasterKeys(keys)
}

// RemoveKMSMasterKeys takes a comma separated string of AWS KMS ARNs and removes the keys corresponding to those ARNs from the metadata's KeySources
func (m *Metadata) RemoveKMSMasterKeys(arns string) {
	var keys []MasterKey
	for _, k := range kms.MasterKeysFromArnString(arns, nil) {
		keys = append(keys, k)
	}
	m.RemoveMasterKeys(keys)
}

// ToMap converts the Metadata to a map for serialization purposes
func (m *Metadata) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["lastmodified"] = m.LastModified.Format(time.RFC3339)
	out["unencrypted_suffix"] = m.UnencryptedSuffix
	out["mac"] = m.MessageAuthenticationCode
	out["version"] = m.Version
	out["shamir"] = m.Shamir
	out["shamir_quorum"] = m.ShamirQuorum
	for _, ks := range m.KeySources {
		var keys []map[string]interface{}
		for _, k := range ks.Keys {
			keys = append(keys, k.ToMap())
		}
		out[ks.Name] = keys
	}
	return out
}

func (m Metadata) getDataKeyShamir() ([]byte, error) {
	var parts [][]byte
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			key, err := k.Decrypt()
			if err != nil {
				fmt.Printf("Key error: %s %s\n", k.ToString(), err)
			} else {
				parts = append(parts, key)
			}
		}
	}
	if len(parts) < m.ShamirQuorum {
		return nil, fmt.Errorf("Not enough parts to recover data key with Shamir. Need %d, have %d.", m.ShamirQuorum, len(parts))
	}
	dataKey, err := shamir.Combine(parts)
	if err != nil {
		return nil, fmt.Errorf("Could not get data key from shamir parts: %s", err)
	}
	return dataKey, nil
}

// getFirstDataKey retrieves the data key from the first MasterKey in the
// Metadata's KeySources that's able to return it.
func (m Metadata) getFirstDataKey() ([]byte, error) {
	errMsg := "Could not decrypt the data key with any of the master keys:\n"
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			key, err := k.Decrypt()
			if err == nil {
				return key, nil
			}
			keyType := "Unknown"
			if _, ok := k.(*pgp.MasterKey); ok {
				keyType = "GPG"
			} else if _, ok := k.(*kms.MasterKey); ok {
				keyType = "KMS"
			}
			errMsg += fmt.Sprintf("\t[%s]: %s:\t%s\n", keyType, k.ToString(), err)
		}
	}
	return nil, fmt.Errorf(errMsg)
}

// GetDataKey retrieves the data key.
func (m Metadata) GetDataKey() ([]byte, error) {
	if m.Shamir {
		return m.getDataKeyShamir()
	} else {
		return m.getFirstDataKey()
	}
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

// MapToMetadata tries to convert a map[string]interface{} obtained from an encrypted file into a Metadata struct.
func MapToMetadata(data map[string]interface{}) (Metadata, error) {
	var metadata Metadata
	mac, ok := data["mac"].(string)
	if !ok {
		fmt.Println("WARNING: no MAC was found on the input file. " +
			"Verification will fail. You can use --ignore-mac to skip verification.")
	}
	metadata.MessageAuthenticationCode = mac
	lastModified, err := time.Parse(time.RFC3339, data["lastmodified"].(string))
	if err != nil {
		return metadata, fmt.Errorf("Could not parse last modified date: %s", err)
	}
	metadata.LastModified = lastModified
	unencryptedSuffix, ok := data["unencrypted_suffix"].(string)
	if !ok {
		unencryptedSuffix = DefaultUnencryptedSuffix
	}
	metadata.UnencryptedSuffix = unencryptedSuffix
	if metadata.Version, ok = data["version"].(string); !ok {
		metadata.Version = strconv.FormatFloat(data["version"].(float64), 'f', -1, 64)
	}
	shamir, ok := data["shamir"].(bool)
	if ok {
		metadata.Shamir = shamir
	}
	if shamirQuorum, ok := data["shamir_quorum"].(float64); ok {
		metadata.ShamirQuorum = int(shamirQuorum)
	} else if shamirQuorum, ok := data["shamir_quorum"].(int); ok {
		metadata.ShamirQuorum = shamirQuorum
	}
	if k, ok := data["kms"].([]interface{}); ok {
		ks, err := mapKMSEntriesToKeySource(k)
		if err == nil {
			metadata.KeySources = append(metadata.KeySources, ks)
		}
	}

	if pgp, ok := data["pgp"].([]interface{}); ok {
		ks, err := mapPGPEntriesToKeySource(pgp)
		if err == nil {
			metadata.KeySources = append(metadata.KeySources, ks)
		}
	}
	return metadata, nil
}

func convertToMapStringInterface(in map[interface{}]interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for k, v := range in {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("Map contains non-string-key (Type %T): %s", k, k)
		}
		m[key] = v
	}
	return m, nil
}

func mapKMSEntriesToKeySource(in []interface{}) (KeySource, error) {
	var keys []MasterKey
	keysource := KeySource{Name: "kms", Keys: keys}
	for _, v := range in {
		entry, ok := v.(map[string]interface{})
		if !ok {
			m, ok := v.(map[interface{}]interface{})
			var err error
			entry, err = convertToMapStringInterface(m)
			if !ok || err != nil {
				fmt.Println("KMS entry has invalid format, skipping...")
				continue
			}
		}
		key := &kms.MasterKey{}
		key.Arn = entry["arn"].(string)
		key.EncryptedKey = entry["enc"].(string)
		role, ok := entry["role"].(string)
		if ok {
			key.Role = role
		}
		creationDate, err := time.Parse(time.RFC3339, entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		if _, ok := entry["context"]; ok {
			key.EncryptionContext = kms.ParseKMSContext(entry["context"])
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}

func mapPGPEntriesToKeySource(in []interface{}) (KeySource, error) {
	var keys []MasterKey
	keysource := KeySource{Name: "pgp", Keys: keys}
	for _, v := range in {
		entry, ok := v.(map[string]interface{})
		if !ok {
			m, ok := v.(map[interface{}]interface{})
			var err error
			entry, err = convertToMapStringInterface(m)
			if !ok || err != nil {
				fmt.Println("PGP entry has invalid format, skipping...")
				continue
			}
		}
		key := &pgp.MasterKey{}
		key.Fingerprint = entry["fp"].(string)
		key.EncryptedKey = entry["enc"].(string)
		creationDate, err := time.Parse(time.RFC3339, entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}
