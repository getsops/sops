package sops

import (
	"crypto/sha512"
	"fmt"
	"go.mozilla.org/sops/aes"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "2006-01-02T15:04:05Z"

const DefaultUnencryptedSuffix = "_unencrypted"

type Error string

func (e Error) Error() string { return string(e) }

const MacMismatch = Error("MAC mismatch")

type TreeItem struct {
	Key   string
	Value interface{}
}

type TreeBranch []TreeItem

type Tree struct {
	Branch   TreeBranch
	Metadata Metadata
}

func (tree TreeBranch) WalkValue(in interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, path)
	case int:
		return onLeaves(in, path)
	case bool:
		return onLeaves(in, path)
	case TreeBranch:
		return tree.WalkBranch(in, path, onLeaves)
	case []interface{}:
		return tree.WalkSlice(in, path, onLeaves)
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (tree TreeBranch) WalkSlice(in []interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) ([]interface{}, error) {
	for i, v := range in {
		newV, err := tree.WalkValue(v, path, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
	}
	return in, nil
}

func (tree TreeBranch) WalkBranch(in TreeBranch, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (TreeBranch, error) {
	for i, item := range in {
		newV, err := tree.WalkValue(item.Value, append(path, item.Key), onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
	}
	return in, nil
}

func (tree Tree) Encrypt(key string) (string, error) {
	hash := sha512.New()
	_, err := tree.Branch.WalkBranch(tree.Branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
		bytes, err := toBytes(in)
		if !strings.HasSuffix(path[len(path)-1], tree.Metadata.UnencryptedSuffix) {
			var err error
			in, err = aes.Encrypt(in, key, []byte(strings.Join(path, ":")+":"))
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

func (tree Tree) Decrypt(key string) (string, error) {
	hash := sha512.New()
	_, err := tree.Branch.WalkBranch(tree.Branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
		var v interface{}
		if !strings.HasSuffix(path[len(path)-1], tree.Metadata.UnencryptedSuffix) {
			var err error
			v, err = aes.Decrypt(in.(string), key, []byte(strings.Join(path, ":")+":"))
			if err != nil {
				return nil, fmt.Errorf("Could not decrypt value: %s", err)
			}
		} else {
			v = in
		}
		bytes, err := toBytes(v)
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

type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	MessageAuthenticationCode string
	Version                   string
	KeySources                []KeySource
}

type KeySource struct {
	Name string
	Keys []MasterKey
}

type MasterKey interface {
	Encrypt(dataKey string) error
	EncryptIfNeeded(dataKey string) error
	Decrypt() (string, error)
	NeedsRotation() bool
	ToString() string
	ToMap() map[string]string
}

type Store interface {
	Load(in string) (TreeBranch, error)
	LoadMetadata(in string) (Metadata, error)
	Dump(TreeBranch) (string, error)
	DumpWithMetadata(TreeBranch, Metadata) (string, error)
}

func (m *Metadata) MasterKeyCount() int {
	count := 0
	for _, ks := range m.KeySources {
		count += len(ks.Keys)
	}
	return count
}

func (m *Metadata) RemoveMasterKeys(keys []MasterKey) {
	for _, ks := range m.KeySources {
		for i, k := range ks.Keys {
			for _, k2 := range keys {
				if k.ToString() == k2.ToString() {
					ks.Keys = append(ks.Keys[:i], ks.Keys[i+1:]...)
				}
			}
		}
	}
}

func (m *Metadata) UpdateMasterKeys(dataKey string) {
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := k.EncryptIfNeeded(dataKey)
			if err != nil {
				fmt.Println("[WARNING]: could not encrypt data key with master key ", k.ToString())
			}
		}
	}
}

func (m *Metadata) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["lastmodified"] = m.LastModified.Format("2006-01-02T15:04:05Z")
	out["unencrypted_suffix"] = m.UnencryptedSuffix
	out["mac"] = m.MessageAuthenticationCode
	out["version"] = m.Version
	for _, ks := range m.KeySources {
		keys := make([]map[string]string, 0)
		for _, k := range ks.Keys {
			keys = append(keys, k.ToMap())
		}
		out[ks.Name] = keys
	}
	return out
}

func toBytes(in interface{}) ([]byte, error) {
	switch in := in.(type) {
	case string:
		return []byte(in), nil
	case int:
		return []byte(strconv.Itoa(in)), nil
	case float64:
		return []byte(strconv.FormatFloat(in, 'f', -1, 64)), nil
	case bool:
		return []byte(strconv.FormatBool(in)), nil
	case []byte:
		return in, nil
	default:
		return nil, fmt.Errorf("Could not convert unknown type %T to bytes", in)
	}
}
