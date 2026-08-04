package json //import "github.com/getsops/sops/v3/stores/json"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/stores"
)

// Store handles storage of JSON data.
type Store struct {
	config config.JSONStoreConfig
}

func NewStore(c *config.JSONStoreConfig) *Store {
	return &Store{config: *c}
}

func (store *Store) Name() string {
	return "json"
}

// BinaryStore handles storage of binary data in a JSON envelope.
type BinaryStore struct {
	store  Store
	config config.JSONBinaryStoreConfig
}

// The binary store uses a single key ("data") to store everything.
func (store *BinaryStore) IsSingleValueStore() bool {
	return true
}

func (store *BinaryStore) Name() string {
	return "binary"
}

func NewBinaryStore(c *config.JSONBinaryStoreConfig) *BinaryStore {
	return &BinaryStore{config: *c, store: *NewStore(&config.JSONStoreConfig{
		Indent: c.Indent,
	})}
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

var BinaryStoreEmitPlainError = errors.New("error emitting binary store")

// EmitPlainFile produces plaintext json file's bytes from its corresponding sops.TreeBranches object
func (store BinaryStore) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	if len(in) != 1 {
		return nil, fmt.Errorf("%w: there must be exactly one tree branch", BinaryStoreEmitPlainError)
	}
	// JSON stores a single object per file
	for _, item := range in[0] {
		if item.Key == "data" {
			if value, ok := item.Value.(string); ok {
				return []byte(value), nil
			} else {
				return nil, fmt.Errorf("%w: 'data' key in tree does not have a string value", BinaryStoreEmitPlainError)
			}
		}
	}
	return nil, fmt.Errorf("%w: no binary data found in tree", BinaryStoreEmitPlainError)
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
		} else if ok && delim.String() == "[" {
			item, err := store.sliceFromJSONDecoder(dec)
			if err != nil {
				return slice, err
			}
			slice = append(slice, item)
		} else {
			v, err := normalizeJSONNumber(t)
			if err != nil {
				return slice, err
			}
			slice = append(slice, v)
		}
	}
}

// normalizeJSONNumber converts a json.Number scalar (produced because the
// decoder runs with UseNumber) into an int for integers within the int64 range
// and a float64 otherwise. Non-number tokens are returned unchanged; a number
// representable as neither returns the json.Number.Float64 error.
func normalizeJSONNumber(t interface{}) (interface{}, error) {
	n, ok := t.(json.Number)
	if !ok {
		return t, nil
	}
	if i, err := n.Int64(); err == nil {
		return int(i), nil
	}
	f, err := n.Float64()
	if err != nil {
		return nil, err
	}
	return f, nil
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
		v, err := normalizeJSONNumber(value)
		if err != nil {
			return item, err
		}
		item.Value = v
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
	empty := true
	for _, item := range array {
		if _, ok := item.(sops.Comment); ok {
			continue
		}
		if !empty {
			out += ","
		}
		v, err := store.encodeValue(item)
		if err != nil {
			return nil, err
		}
		out += string(v)
		empty = false
	}
	out += "]"
	return []byte(out), nil
}

func (store Store) encodeTree(tree sops.TreeBranch) ([]byte, error) {
	out := "{"
	empty := true
	for _, item := range tree {
		if _, ok := item.Key.(sops.Comment); ok {
			continue
		}
		if !empty {
			out += ","
		}
		v, err := store.encodeValue(item.Value)
		if err != nil {
			return nil, fmt.Errorf("Error encoding value %s: %s", v, err)
		}
		k, err := json.Marshal(item.Key.(string))
		if err != nil {
			return nil, fmt.Errorf("Error encoding key %s: %s", k, err)
		}
		out += string(k) + `: ` + string(v)
		empty = false
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
	// Decode numbers as json.Number instead of the default float64, then
	// normalize each to int/float64 (see normalizeJSONNumber). The default
	// float64 silently loses precision for integers larger than 2^53.
	dec.UseNumber()
	value, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := value.(json.Delim); ok {
		if delim.String() != "{" {
			return nil, fmt.Errorf("SOPS only supports JSON files with a top-level object (starting with '{'), not arrays or other types. Got delimiter %s instead. To encrypt this file, wrap it in an object, e.g., {\"data\": [...]}", value)
		}
	} else {
		v, nerr := normalizeJSONNumber(value)
		if nerr != nil {
			v = value
		}
		return nil, fmt.Errorf("SOPS only supports JSON files with a top-level object (starting with '{'), not other JSON types. Got %#v of type %T instead", v, v)
	}
	return store.treeBranchFromJSONDecoder(dec)
}

func (store Store) reindentJSON(in []byte) ([]byte, error) {
	var out bytes.Buffer
	indent := "\t"
	if store.config.Indent > -1 {
		indent = strings.Repeat(" ", store.config.Indent)
	} else if store.config.Indent < -1 {
		return nil, errors.New("JSON Indentation parameter smaller than -1 is not accepted")
	}
	err := json.Indent(&out, in, "", indent)
	return out.Bytes(), err
}

// LoadEncryptedFile loads an encrypted secrets file onto a sops.Tree object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}
	branches, metadata, err := stores.ExtractMetadata(branches, stores.MetadataOpts{
		Flatten: stores.MetadataFlattenNone,
	})
	if err != nil {
		return sops.Tree{}, err
	}
	return sops.Tree{
		Branches: branches,
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads plaintext json file bytes onto a sops.TreeBranches object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	branch, err := store.treeBranchFromJSON(in)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	return sops.TreeBranches{
		branch,
	}, nil
}

// EmitEncryptedFile returns the encrypted bytes of the json file corresponding to a
// sops.Tree runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	branches, err := stores.SerializeMetadata(in, stores.MetadataOpts{
		Flatten: stores.MetadataFlattenNone,
	})
	if err != nil {
		return nil, fmt.Errorf("Error marshaling metadata: %s", err)
	}
	return store.EmitPlainFile(branches)
}

// EmitPlainFile returns the plaintext bytes of the json file corresponding to a
// sops.TreeBranches runtime object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	out, err := store.jsonFromTreeBranch(in[0])
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	out = append(out, '\n')
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

// HasSopsTopLevelKey checks whether a top-level "sops" key exists.
func (store *Store) HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	return stores.HasSopsTopLevelKey(branch)
}

// HasSopsTopLevelKey checks whether a top-level "sops" key exists.
func (store *BinaryStore) HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	return stores.HasSopsTopLevelKey(branch)
}
