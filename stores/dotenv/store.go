package dotenv //import "go.mozilla.org/sops/v3/stores/dotenv"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/stores"
)

// SopsPrefix is the prefix for all metadatada entry keys
const SopsPrefix = "sops_"

// Store handles storage of dotenv data
type Store struct {
}

// LoadEncryptedFile loads an encrypted file's bytes onto a sops.Tree runtime object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}

	var resultBranch sops.TreeBranch
	mdMap := make(map[string]interface{})
	for _, item := range branches[0] {
		switch key := item.Key.(type) {
		case string:
			if strings.HasPrefix(key, SopsPrefix) {
				key = key[len(SopsPrefix):]
				mdMap[key] = item.Value
			} else {
				resultBranch = append(resultBranch, item)
			}
		case sops.Comment:
			resultBranch = append(resultBranch, item)
		default:
			panic(fmt.Sprintf("Unexpected type: %T (value %#v)", key, key))
		}
	}

	metadata, err := mapToMetadata(mdMap)
	if err != nil {
		return sops.Tree{}, err
	}
	internalMetadata, err := metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, err
	}

	return sops.Tree{
		Branches: sops.TreeBranches{
			resultBranch,
		},
		Metadata: internalMetadata,
	}, nil
}

// LoadPlainFile returns the contents of a plaintext file loaded onto a
// sops runtime object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	var branches sops.TreeBranches
	items, err := parse(in)
	if err != nil {
		return nil, err
	}
	branches = append(branches, items)
	return branches, nil
}

// EmitEncryptedFile returns the encrypted file's bytes corresponding to a sops
// runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	metadata := stores.MetadataFromInternal(in.Metadata)
	mdItems, err := metadataToMap(metadata)
	if err != nil {
		return nil, err
	}
	for key, value := range mdItems {
		if value == nil {
			continue
		}
		in.Branches[0] = append(in.Branches[0], sops.TreeItem{Key: SopsPrefix + key, Value: value})
	}
	return store.EmitPlainFile(in.Branches)
}

// EmitPlainFile returns the plaintext file's bytes corresponding to a sops
// runtime object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, item := range in[0] {
		if isComplexValue(item.Value) {
			return nil, fmt.Errorf("cannot use complex value in dotenv file: %s", item.Value)
		}
		var line string
		if comment, ok := item.Key.(sops.Comment); ok {
			line = fmt.Sprintf("# %s\n", comment.Value)
		} else {
			value := strings.Replace(item.Value.(string), `'`, `\'`, -1)
			line = fmt.Sprintf("%s='%s'\n", item.Key, value)
		}
		buffer.WriteString(line)
	}
	return buffer.Bytes(), nil
}

// EmitValue returns a single value as bytes
func (Store) EmitValue(v interface{}) ([]byte, error) {
	if s, ok := v.(string); ok {
		return []byte(s), nil
	}
	return nil, fmt.Errorf("the dotenv store only supports emitting strings, got %T", v)
}

// EmitExample returns the bytes corresponding to an example Flat Tree runtime object
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleFlatTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}

func metadataToMap(md stores.Metadata) (map[string]interface{}, error) {
	var mdMap map[string]interface{}
	inrec, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &mdMap)
	if err != nil {
		return nil, err
	}
	flat := stores.Flatten(mdMap)
	for k, v := range flat {
		if s, ok := v.(string); ok {
			flat[k] = strings.Replace(s, "\n", "\\n", -1)
		}
	}
	return flat, nil
}

func mapToMetadata(m map[string]interface{}) (stores.Metadata, error) {
	for k, v := range m {
		if s, ok := v.(string); ok {
			m[k] = strings.Replace(s, "\\n", "\n", -1)
		}
	}
	m = stores.Unflatten(m)
	var md stores.Metadata
	inrec, err := json.Marshal(m)
	if err != nil {
		return md, err
	}
	err = json.Unmarshal(inrec, &md)
	return md, err
}

func isComplexValue(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	case sops.TreeBranch:
		return true
	}
	return false
}
