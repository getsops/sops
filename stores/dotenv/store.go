package dotenv //import "github.com/getsops/sops/v3/stores/dotenv"

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/stores"
)

// SopsPrefix is the prefix for all metadatada entry keys
const SopsPrefix = stores.SopsMetadataKey + "_"

// Store handles storage of dotenv data
type Store struct {
	config config.DotenvStoreConfig
}

func NewStore(c *config.DotenvStoreConfig) *Store {
	return &Store{config: *c}
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

	stores.DecodeNewLines(mdMap)
	err = stores.DecodeNonStrings(mdMap)
	if err != nil {
		return sops.Tree{}, err
	}
	metadata, err := stores.UnflattenMetadata(mdMap)
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
	var branch sops.TreeBranch

	for _, line := range bytes.Split(in, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			branch = append(branch, sops.TreeItem{
				Key:   sops.Comment{Value: string(line[1:])},
				Value: nil,
			})
		} else {
			pos := bytes.Index(line, []byte("="))
			if pos == -1 {
				return nil, fmt.Errorf("invalid dotenv input line: %s", line)
			}
			branch = append(branch, sops.TreeItem{
				Key:   string(line[:pos]),
				Value: strings.Replace(string(line[pos+1:]), "\\n", "\n", -1),
			})
		}
	}

	branches = append(branches, branch)
	return branches, nil
}

// EmitEncryptedFile returns the encrypted file's bytes corresponding to a sops
// runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	metadata := stores.MetadataFromInternal(in.Metadata)
	mdItems, err := stores.FlattenMetadata(metadata)
	if err != nil {
		return nil, err
	}

	stores.EncodeNonStrings(mdItems)
	stores.EncodeNewLines(mdItems)

	var keys []string
	for k := range mdItems {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		var value = mdItems[key]
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
			line = fmt.Sprintf("#%s\n", comment.Value)
		} else {
			value := strings.Replace(item.Value.(string), "\n", "\\n", -1)
			line = fmt.Sprintf("%s=%s\n", item.Key, value)
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

func isComplexValue(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	case sops.TreeBranch:
		return true
	}
	return false
}

// HasSopsTopLevelKey checks whether a top-level "sops" key exists.
func (store *Store) HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	for _, b := range branch {
		if key, ok := b.Key.(string); ok {
			if strings.HasPrefix(key, SopsPrefix) {
				return true
			}
		}
	}
	return false
}
