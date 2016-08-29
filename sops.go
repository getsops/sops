package sops

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DefaultUnencryptedSuffix is the default suffix a TreeItem key has to end with for sops to leave its Value unencrypted
const DefaultUnencryptedSuffix = "_unencrypted"

type sopsError string

func (e sopsError) Error() string { return string(e) }

// MacMismatch occurs when the computed MAC does not match the expected ones
const MacMismatch = sopsError("MAC mismatch")

// MetadataNotFound occurs when the input file is malformed and doesn't have sops metadata in it
const MetadataNotFound = sopsError("sops metadata not found")

// DataKeyCipher provides a way to encrypt and decrypt the data key used to encrypt and decrypt sops files, so that the data key can be stored alongside the encrypted content. A DataKeyCipher must be able to decrypt the values it encrypts.
type DataKeyCipher interface {
	Encrypt(value interface{}, key []byte, additionalAuthData []byte) (string, error)
	Decrypt(value string, key []byte, additionalAuthData []byte) (interface{}, error)
}

// TreeItem is an item inside sops's tree
type TreeItem struct {
	Key   string
	Value interface{}
}

// TreeBranch is a branch inside sops's tree. It is a slice of TreeItems and is therefore ordered
type TreeBranch []TreeItem

// Tree is the data structure used by sops to represent documents internally
type Tree struct {
	Branch   TreeBranch
	Metadata Metadata
}

// Truncate truncates the tree following Python dictionary access syntax, for example, ["foo"][2].
func (tree TreeBranch) Truncate(path string) (interface{}, error) {
	components := strings.Split(path, "[")
	var current interface{} = tree
	for _, component := range components {
		if component == "" {
			continue
		}
		if component[len(component)-1] != ']' {
			return nil, fmt.Errorf("Invalid tree path format string: %s", path)
		}
		component = component[:len(component)-1]
		component = strings.Replace(component, `"`, "", 2)
		component = strings.Replace(component, `'`, "", 2)
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
		newV, err := tree.walkValue(item.Value, append(path, item.Key), onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
	}
	return in, nil
}

// Encrypt walks over the tree and encrypts all values with the provided cipher, except those whose key ends with the UnencryptedSuffix specified on the Metadata struct. If encryption is successful, it returns the MAC for the encrypted tree.
func (tree Tree) Encrypt(key []byte, cipher DataKeyCipher) (string, error) {
	hash := sha512.New()
	_, err := tree.Branch.walkBranch(tree.Branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
		bytes, err := ToBytes(in)
		if !strings.HasSuffix(path[len(path)-1], tree.Metadata.UnencryptedSuffix) {
			var err error
			in, err = cipher.Encrypt(in, key, []byte(strings.Join(path, ":")+":"))
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
func (tree Tree) Decrypt(key []byte, cipher DataKeyCipher) (string, error) {
	hash := sha512.New()
	_, err := tree.Branch.walkBranch(tree.Branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
		var v interface{}
		if !strings.HasSuffix(path[len(path)-1], tree.Metadata.UnencryptedSuffix) {
			var err error
			v, err = cipher.Decrypt(in.(string), key, []byte(strings.Join(path, ":")+":"))
			if err != nil {
				return nil, fmt.Errorf("Could not decrypt value: %s", err)
			}
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
func (tree Tree) GenerateDataKey() ([]byte, error) {
	newKey := make([]byte, 32)
	_, err := rand.Read(newKey)
	if err != nil {
		return nil, fmt.Errorf("Could not generate random key: %s", err)
	}
	for _, ks := range tree.Metadata.KeySources {
		for _, k := range ks.Keys {
			k.Encrypt(newKey)
		}
	}
	return newKey, nil
}

// Metadata holds information about a file encrypted by sops
type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	MessageAuthenticationCode string
	Version                   string
	KeySources                []KeySource
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
	ToMap() map[string]string
}

// Store provides a way to load and save the sops tree along with metadata
type Store interface {
	Unmarshal(in []byte) (TreeBranch, error)
	UnmarshalMetadata(in []byte) (Metadata, error)
	Marshal(TreeBranch) ([]byte, error)
	MarshalWithMetadata(TreeBranch, Metadata) ([]byte, error)
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
		for i, k := range ks.Keys {
			for _, k2 := range keys {
				if k.ToString() == k2.ToString() {
					ks.Keys = append(ks.Keys[:i], ks.Keys[i+1:]...)
				}
			}
		}
		m.KeySources[j] = ks
	}
}

// UpdateMasterKeys encrypts the data key with all master keys if it's needed
func (m *Metadata) UpdateMasterKeys(dataKey []byte) {
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := k.EncryptIfNeeded(dataKey)
			if err != nil {
				fmt.Println("[WARNING]: could not encrypt data key with master key ", k.ToString())
			}
		}
	}
}

// AddPGPMasterKeys parses the input comma separated string of GPG fingerprints, generates a PGP MasterKey for each fingerprint, and adds the keys to the PGP KeySource
func (m *Metadata) AddPGPMasterKeys(pgpFps string) {
	for i, ks := range m.KeySources {
		if ks.Name == "pgp" {
			var keys []MasterKey
			for _, k := range pgp.MasterKeysFromFingerprintString(pgpFps) {
				keys = append(keys, &k)
				fmt.Println("Keys to add:", keys)
			}
			ks.Keys = append(ks.Keys, keys...)
			m.KeySources[i] = ks
		}
	}
}

// AddKMSMasterKeys parses the input comma separated string of AWS KMS ARNs, generates a KMS MasterKey for each ARN, and then adds the keys to the KMS KeySource
func (m *Metadata) AddKMSMasterKeys(kmsArns string) {
	for i, ks := range m.KeySources {
		if ks.Name == "kms" {
			var keys []MasterKey
			for _, k := range kms.MasterKeysFromArnString(kmsArns) {
				keys = append(keys, &k)
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
		keys = append(keys, &k)
	}
	m.RemoveMasterKeys(keys)
}

// RemoveKMSMasterKeys takes a comma separated string of AWS KMS ARNs and removes the keys corresponding to those ARNs from the metadata's KeySources
func (m *Metadata) RemoveKMSMasterKeys(arns string) {
	var keys []MasterKey
	for _, k := range kms.MasterKeysFromArnString(arns) {
		keys = append(keys, &k)
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
	for _, ks := range m.KeySources {
		var keys []map[string]string
		for _, k := range ks.Keys {
			keys = append(keys, k.ToMap())
		}
		out[ks.Name] = keys
	}
	return out
}

// GetDataKey retrieves the data key from the first MasterKey in the Metadata's KeySources that's able to return it.
func (m Metadata) GetDataKey() ([]byte, error) {
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			key, err := k.Decrypt()
			if err == nil {
				return key, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not get master key")
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
