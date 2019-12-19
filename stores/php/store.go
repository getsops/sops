package php //import "go.mozilla.org/sops/v3/stores/php"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/stores"
)

// Store handles storage of JSON data.
type Store struct {
}

// BinaryStore handles storage of binary data in a JSON envelope.
type BinaryStore struct {
	store Store
}

// LoadEncryptedFile loads an encrypted json file onto a sops.Tree object
func (store BinaryStore) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	return store.store.LoadEncryptedFile(in)
}

// LoadPlainFile loads a plaintext json file onto a sops.Tree encapsulated
// within a sops.TreeBranches object
func (store BinaryStore) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	return sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   "data",
				Value: string(in),
			},
		},
	}, nil
}

// EmitEncryptedFile produces an encrypted json file's bytes from its corresponding sops.Tree object
func (store BinaryStore) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	return store.store.EmitEncryptedFile(in)
}

// EmitPlainFile produces plaintext json file's bytes from its corresponding sops.TreeBranches object
func (store BinaryStore) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	// JSON stores a single object per file
	for _, item := range in[0] {
		if item.Key == "data" {
			return []byte(item.Value.(string)), nil
		}
	}
	return nil, fmt.Errorf("No binary data found in tree")
}

// EmitValue extracts a value from a generic interface{} object representing a structured set
// of binary files
func (store BinaryStore) EmitValue(v interface{}) ([]byte, error) {
	return nil, fmt.Errorf("Binary files are not structured and extracting a single value is not possible")
}

// EmitExample returns the example's plaintext json file bytes
func (store BinaryStore) EmitExample() []byte {
	return []byte("Welcome to SOPS! Edit this file as you please!")
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
		if v == nil {
			return []byte{}, nil
		}
		str, ok := v.(string)
		if !ok {
			return json.Marshal(v)
		}
		a := []byte(str)
		return a, nil
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
	out := ""
	for i, item := range tree {
		v, err := store.encodeValue(item.Value)
		if err != nil {
			return nil, fmt.Errorf("Error encoding value %s: %s", v, err)
		}
		k := []byte(item.Key.(string))

		if string(v) == "null" || v == nil || len(v) == 0 {
			out += string(k)
		} else if string(k) == "" {
			out += string(v)
		} else {
			out += string(k) + `=` + string(v)
		}
		if i != len(tree)-1 {
			out += "\n"
		}
	}
	return []byte(out), nil
}

func (store Store) jsonFromTreeBranch(branch sops.TreeBranch) ([]byte, error) {
	out, err := store.encodeTree(branch)
	if err != nil {
		return nil, err
	}
	return out, nil
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

// LoadEncryptedFile loads an encrypted secrets file onto a sops.Tree object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	// Because we don't know what fields the input file will have, we have to
	// load the file in two steps.
	// First, we load the file's metadata, the structure of which is known.
	var resultBranch sops.TreeBranch
	var metadata sops.Metadata
	branches, err := store.LoadPlainFile(in)

	for _, branch := range branches {
		for _, item := range branch {
			if item.Key == "sops" {
				v := []byte(item.Value.(string))

				m := make(map[string]interface{})
				err = json.Unmarshal(v, &m)
				if err != nil {
					log.Fatal("Tampered sops information", err)
				}
				for k, v := range m {
					if s, ok := v.(string); ok {
						m[k] = strings.Replace(s, "\\n", "\n", -1)
					}
				}
				m = stores.Unflatten(m)
				var md stores.Metadata
				inrec, err := json.Marshal(m)
				if err != nil {
					log.Fatal("Tampered sops information", err)
				}
				err = json.Unmarshal(inrec, &md)
				if err != nil {
					log.Fatal("Tampered sops information", err)
				}
				metadata, err = md.ToInternal()
				if err != nil {
					log.Fatal("Tampered sops information", err)
				}
			} else {
				resultBranch = append(resultBranch, item)
			}
		}
	}
	return sops.Tree{
		Branches: sops.TreeBranches{
			resultBranch,
		},
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads plaintext json file bytes onto a sops.TreeBranches object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	var branches sops.TreeBranches
	var branch sops.TreeBranch
	var isComments, isMetadata bool
	var metadata []byte
	lines := bytes.Split(in, []byte("\n"))
	for i, line := range lines {
		if len(line) < 2 {
			branch = append(branch, sops.TreeItem{
				Key:   string(line),
				Value: nil,
			})
			continue
		}

		if isMetadata || (len(line) > 6 && string(line[:6]) == "sops={") {
			metadata = append(metadata, line[5:]...)
			if i == len(lines)-1 {
				branch = append(branch, sops.TreeItem{
					Key:   "sops",
					Value: string(metadata),
				})
			}
			isMetadata = true
		}

		if isComments || string(line[:2]) == "/*" {
			if string(line[len(line)-2:]) == "*/" {
				branch = append(branch, sops.TreeItem{
					Key:   string(line),
					Value: nil,
				})
				isComments = false
				continue
			}
			branch = append(branch, sops.TreeItem{
				Key:   string(line),
				Value: nil,
			})
			isComments = true
		}

		if !isComments {
			if line[0] == '#' || string(line[:2]) == "//" {
				branch = append(branch, sops.TreeItem{
					Key:   string(line),
					Value: nil,
				})
			} else {
				pos := bytes.Index(line, []byte("="))
				if pos == -1 {
					branch = append(branch, sops.TreeItem{
						Key:   string(line[pos+1:]),
						Value: nil,
					})
					continue
				}
				branch = append(branch, sops.TreeItem{
					Key:   string(line[:pos]),
					Value: string(line[pos+1:]),
				})
			}
		}
	}

	branches = append(branches, branch)
	return branches, nil
}

// EmitEncryptedFile returns the encrypted bytes of the json file corresponding to a
// sops.Tree runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	tree := append(in.Branches[0], sops.TreeItem{Key: "sops", Value: stores.MetadataFromInternal(in.Metadata)})
	out, err := store.jsonFromTreeBranch(tree)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	return out, nil
}

// EmitPlainFile returns the plaintext bytes of the json file corresponding to a
// sops.TreeBranches runtime object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	out, err := store.jsonFromTreeBranch(in[0])
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	return out, nil
}

// EmitValue returns bytes corresponding to a single encoded value
// in a generic interface{} object
func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	s, err := store.encodeValue(v)
	if err != nil {
		return nil, err
	}
	return store.reindentJSON(s)
}

// EmitExample returns the bytes corresponding to an example complex tree
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleComplexTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}
