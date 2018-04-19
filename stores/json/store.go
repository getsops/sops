package json //import "go.mozilla.org/sops/stores/json"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/stores"
)

// Store handles storage of JSON data.
type Store struct {
}

// BinaryStore handles storage of binary data in a JSON envelope.
type BinaryStore struct {
	store Store
}

func (store BinaryStore) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	return store.store.LoadEncryptedFile(in)
}

func (store BinaryStore) LoadPlainFile(in []byte) (sops.TreeBranch, error) {
	return sops.TreeBranch{
		sops.TreeItem{
			Key:   "data",
			Value: string(in),
		},
	}, nil
}

func (store BinaryStore) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	return store.store.EmitEncryptedFile(in)
}

func (store BinaryStore) EmitPlainFile(in sops.TreeBranch) ([]byte, error) {
	for _, item := range in {
		if item.Key == "data" {
			return []byte(item.Value.(string)), nil
		}
	}
	return nil, fmt.Errorf("No binary data found in tree")
}

func (store BinaryStore) EmitValue(v interface{}) ([]byte, error) {
	return nil, fmt.Errorf("Binary files are not structured and extracting a single value is not possible")
}

func (store Store) sliceFromJSONDecoder(dec *json.Decoder) ([]interface{}, error) {
	var slice []interface{}
	for {
		t, err := dec.Token()
		if err != nil {
			return slice, err
		}
		if delim, ok := t.(json.Delim); ok && delim.String() == "]" {
			return slice, nil
		} else if ok && delim.String() == "{" {
			item, err := store.treeBranchFromJSONDecoder(dec)
			if err != nil {
				return slice, err
			}
			slice = append(slice, item)
		} else {
			slice = append(slice, t)
		}
	}
}

var errEndOfObject = fmt.Errorf("End of object")

func (store Store) treeItemFromJSONDecoder(dec *json.Decoder) (sops.TreeItem, error) {
	var item sops.TreeItem
	key, err := dec.Token()
	if err != nil && err != io.EOF {
		return item, err
	}
	if k, ok := key.(string); ok {
		item.Key = k
	} else if d, ok := key.(json.Delim); ok && d.String() == "}" {
		return item, errEndOfObject
	} else {
		return item, fmt.Errorf("Expected JSON object key, got %s of type %T instead", key, key)
	}
	value, err := dec.Token()
	if err != nil {
		return item, err
	}
	if delim, ok := value.(json.Delim); ok {
		if delim.String() == "[" {
			v, err := store.sliceFromJSONDecoder(dec)
			if err != nil {
				return item, err
			}
			item.Value = v
		}
		if delim.String() == "{" {
			v, err := store.treeBranchFromJSONDecoder(dec)
			if err != nil {
				return item, err
			}
			item.Value = v
		}
	} else {
		item.Value = value
	}
	return item, nil

}

func (store Store) treeBranchFromJSONDecoder(dec *json.Decoder) (sops.TreeBranch, error) {
	var tree sops.TreeBranch
	for {
		item, err := store.treeItemFromJSONDecoder(dec)
		if err == io.EOF {
			return tree, nil
		}
		if err == errEndOfObject {
			return tree, nil
		}
		if err != nil {
			return tree, err
		}
		tree = append(tree, item)
	}
}

func (store Store) encodeValue(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case sops.TreeBranch:
		return store.encodeTree(v)
	case []interface{}:
		return store.encodeArray(v)
	default:
		return json.Marshal(v)
	}
}

func (store Store) encodeArray(array []interface{}) ([]byte, error) {
	out := "["
	for i, item := range array {
		if _, ok := item.(sops.Comment); ok {
			continue
		}
		v, err := store.encodeValue(item)
		if err != nil {
			return nil, err
		}
		out += string(v)
		if i != len(array)-1 {
			out += ","
		}
	}
	out += "]"
	return []byte(out), nil
}

func (store Store) encodeTree(tree sops.TreeBranch) ([]byte, error) {
	out := "{"
	for i, item := range tree {
		if _, ok := item.Key.(sops.Comment); ok {
			continue
		}
		v, err := store.encodeValue(item.Value)
		if err != nil {
			return nil, fmt.Errorf("Error encoding value %s: %s", v, err)
		}
		out += `"` + item.Key.(string) + `": ` + string(v)
		if i != len(tree)-1 {
			out += ","
		}
	}
	return []byte(out + "}"), nil
}

func (store Store) jsonFromTreeBranch(branch sops.TreeBranch) ([]byte, error) {
	out, err := store.encodeTree(branch)
	if err != nil {
		return nil, err
	}
	return store.reindentJSON(out)
}

func (store Store) treeBranchFromJSON(in []byte) (sops.TreeBranch, error) {
	dec := json.NewDecoder(bytes.NewReader(in))
	dec.Token()
	return store.treeBranchFromJSONDecoder(dec)
}

func (store Store) reindentJSON(in []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "\t")
	return out.Bytes(), err
}

func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	// Because we don't know what fields the input file will have, we have to
	// load the file in two steps.
	// First, we load the file's metadata, the structure of which is known.
	metadataHolder := stores.SopsFile{}
	err := json.Unmarshal(in, &metadataHolder)
	if err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			if err.Value == "number" && err.Struct == "Metadata" && err.Field == "version" {
				return sops.Tree{},
					fmt.Errorf("SOPS versions higher than 2.0.10 can not automatically decrypt JSON files " +
						"created with SOPS 1.x. In order to be able to decrypt this file, you can either edit it " +
						"manually and make sure the JSON value under `sops -> version` is a string and not a " +
						"number, or you can rotate the file's key with any version of SOPS between 2.0 and 2.0.10 " +
						"using `sops -r your_file.json`")
			}
		}
		return sops.Tree{}, fmt.Errorf("Error unmarshalling input json: %s", err)
	}
	if metadataHolder.Metadata == nil {
		return sops.Tree{}, sops.MetadataNotFound
	}
	metadata, err := metadataHolder.Metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, err
	}
	// After that, we load the whole file into a map.
	branch, err := store.treeBranchFromJSON(in)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	// Discard metadata, as we already loaded it.
	for i, item := range branch {
		if item.Key == "sops" {
			branch = append(branch[:i], branch[i+1:]...)
		}
	}
	return sops.Tree{
		Branch:   branch,
		Metadata: metadata,
	}, nil
}

func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranch, error) {
	branch, err := store.treeBranchFromJSON(in)
	if err != nil {
		return branch, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	return branch, nil
}

func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	tree := append(in.Branch, sops.TreeItem{Key: "sops", Value: stores.MetadataFromInternal(in.Metadata)})
	out, err := store.jsonFromTreeBranch(tree)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	return out, nil
}

func (store *Store) EmitPlainFile(in sops.TreeBranch) ([]byte, error) {
	out, err := store.jsonFromTreeBranch(in)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	return out, nil
}

func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	s, err := store.encodeValue(v)
	if err != nil {
		return nil, err
	}
	return store.reindentJSON(s)
}
