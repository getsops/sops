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

	"log"

	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
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
	if component[len(component)-1] != ']' {
		return "", fmt.Errorf("Invalid component")
	}
	component = component[:len(component)-1]
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
func (tree Tree) GenerateDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient) ([]byte, []error) {
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

// RemoveMasterKeys removes all of the provided keys from the metadata's KeySources, if they exist there.
func (m *Metadata) RemoveMasterKeys(masterKeys []keys.MasterKey) {
	// TODO: Reimplement this with KeyGroups. It's unclear how it should behave.
	panic("Unimplemented")
}

// AddPGPMasterKeys parses the input comma separated string of GPG fingerprints, generates a PGP MasterKey for each fingerprint, and adds the keys to the PGP KeySource
func (m *Metadata) AddPGPMasterKeys(pgpFps string) {
	// TODO: Reimplement this with KeyGroups. It's unclear how it should behave.
	panic("Unimplemented")
}

// AddKMSMasterKeys parses the input comma separated string of AWS KMS ARNs, generates a KMS MasterKey for each ARN, and then adds the keys to the KMS KeySource
func (m *Metadata) AddKMSMasterKeys(kmsArns string, context map[string]*string) {
	// TODO: Reimplement this with KeyGroups. It's unclear how it should behave.
	panic("Unimplemented")
}

// RemovePGPMasterKeys takes a comma separated string of PGP fingerprints and removes the keys corresponding to those fingerprints from the metadata's KeySources
func (m *Metadata) RemovePGPMasterKeys(pgpFps string) {
	// TODO: Reimplement this with KeyGroups. It's unclear how it should behave.
	panic("Unimplemented")
}

// RemoveKMSMasterKeys takes a comma separated string of AWS KMS ARNs and removes the keys corresponding to those ARNs from the metadata's KeySources
func (m *Metadata) RemoveKMSMasterKeys(arns string) {
	// TODO: Reimplement this with KeyGroups. It's unclear how it should behave.
	panic("Unimplemented")
}

func (m *Metadata) UpdateMasterKeysWithKeyServices(dataKey []byte, svcs []keyservice.KeyServiceClient) (errs []error) {
	if len(svcs) == 0 {
		return []error{
			fmt.Errorf("No key services provided, can not update master keys."),
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
			for _, svc := range svcs {
				rsp, err := svc.Encrypt(context.Background(), &keyservice.EncryptRequest{
					Key:       &svcKey,
					Plaintext: part,
				})
				if err != nil {
					errs = append(errs, fmt.Errorf("Failed to encrypt new data key with master key %q: %v\n", key.ToString(), err))
					continue
				}
				key.SetEncryptedDataKey(rsp.Ciphertext)
				// Only need to encrypt the key successfully with one service
				break
			}
		}
	}
	return
}

// UpdateMasterKeys encrypts the data key with all master keys
func (m *Metadata) UpdateMasterKeys(dataKey []byte) (errs []error) {
	return m.UpdateMasterKeysWithKeyServices(dataKey, []keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	})
}

// ToMap converts the Metadata to a map for serialization purposes
func (m *Metadata) ToMap() map[string]interface{} {
	// TODO: This doesn't belong here. This is serialization logic.
	// It should probably be rewritten so that sops.Metadata gets mapped to some sort of stores.Metadata struct,
	// which then gets serialized directly
	out := make(map[string]interface{})
	out["lastmodified"] = m.LastModified.Format(time.RFC3339)
	out["unencrypted_suffix"] = m.UnencryptedSuffix
	out["mac"] = m.MessageAuthenticationCode
	out["version"] = m.Version
	out["shamir_quorum"] = m.ShamirQuorum
	if len(m.KeyGroups) == 1 {
		for k, v := range m.keyGroupToMap(m.KeyGroups[0]) {
			out[k] = v
		}
	} else {
		// This is very bad and I should feel bad
		var groups []map[string][]map[string]interface{}
		for _, group := range m.KeyGroups {
			groups = append(groups, m.keyGroupToMap(group))
		}
		out["key_groups"] = groups
	}
	return out
}

func (m *Metadata) keyGroupToMap(group KeyGroup) (keys map[string][]map[string]interface{}) {
	keys = make(map[string][]map[string]interface{})
	for _, k := range group {
		switch k := k.(type) {
		case *pgp.MasterKey:
			keys["pgp"] = append(keys["pgp"], k.ToMap())
		case *kms.MasterKey:
			keys["kms"] = append(keys["kms"], k.ToMap())
		}
	}
	return
}

// GetDataKeyWithKeyServices retrieves the data key, asking KeyServices to decrypt it with each
// MasterKey in the Metadata's KeySources until one of them succeeds.
func (m Metadata) GetDataKeyWithKeyServices(svcs []keyservice.KeyServiceClient) ([]byte, error) {
	errMsg := "Could not decrypt the data key with any of the master keys:\n"
	var parts [][]byte
	for _, group := range m.KeyGroups {
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
				break
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

// MapToMetadata tries to convert a map[string]interface{} obtained from an encrypted file into a Metadata struct.
func MapToMetadata(data map[string]interface{}) (Metadata, error) {
	// TODO: This doesn't belong here. This is serialization logic.
	// It should probably be rewritten so that the YAML/JSON gets mapped to a struct which then gets mapped to a
	// sops.Metadata
	// TODO: Use KeyGroups
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
	if shamirQuorum, ok := data["shamir_quorum"].(float64); ok {
		metadata.ShamirQuorum = int(shamirQuorum)
	} else if shamirQuorum, ok := data["shamir_quorum"].(int); ok {
		metadata.ShamirQuorum = shamirQuorum
	}
	if keyGroups, ok := data["key_groups"]; ok {
		var kgs []KeyGroup
		if gs, ok := keyGroups.([]interface{}); ok {
			for _, g := range gs {
				g := g.(map[interface{}]interface{})
				var group KeyGroup
				if k, ok := g["kms"].([]interface{}); ok {
					ks, err := mapKMSEntriesToKeySlice(k)
					if err == nil {
						group = append(group, ks...)
					}
				}
				if pgp, ok := g["pgp"].([]interface{}); ok {
					ks, err := mapPGPEntriesToKeySlice(pgp)
					if err == nil {
						group = append(group, ks...)
					}
				}
				kgs = append(kgs, group)
			}
		}
		metadata.KeyGroups = kgs
	} else {
		// Old data format, just one KeyGroup
		var group KeyGroup
		if k, ok := data["kms"].([]interface{}); ok {
			ks, err := mapKMSEntriesToKeySlice(k)
			if err == nil {
				group = append(group, ks...)
			}
		}
		if pgp, ok := data["pgp"].([]interface{}); ok {
			ks, err := mapPGPEntriesToKeySlice(pgp)
			if err == nil {
				group = append(group, ks...)
			}
		}
		metadata.KeyGroups = []KeyGroup{group}
	}
	return metadata, nil
}

func convertToMapStringInterface(in map[interface{}]interface{}) (map[string]interface{}, error) {
	// TODO: This doesn't belong here. This is serialization logic.
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

func mapKMSEntriesToKeySlice(in []interface{}) ([]keys.MasterKey, error) {
	// TODO: This doesn't belong here. This is serialization logic.
	var keys []keys.MasterKey
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
			return keys, fmt.Errorf("Could not parse creation date: %s", err)
		}
		if _, ok := entry["context"]; ok {
			key.EncryptionContext = kms.ParseKMSContext(entry["context"])
		}
		key.CreationDate = creationDate
		keys = append(keys, key)
	}
	return keys, nil
}

func mapPGPEntriesToKeySlice(in []interface{}) ([]keys.MasterKey, error) {
	// TODO: This doesn't belong here. This is serialization logic.
	var keys []keys.MasterKey
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
			return keys, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keys = append(keys, key)
	}
	return keys, nil
}
