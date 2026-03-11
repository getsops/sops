package dotenv //import "github.com/getsops/sops/v3/stores/dotenv"

import (
	"bytes"
	"fmt"
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

func (store *Store) Name() string {
	return "dotenv"
}

// LoadEncryptedFile loads an encrypted file's bytes onto a sops.Tree runtime object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}
	branches, metadata, err := stores.ExtractMetadata(branches, stores.MetadataOpts{
		Flatten: stores.MetadataFlattenFull,
	})
	if err != nil {
		return sops.Tree{}, err
	}
	return sops.Tree{
		Branches: branches,
		Metadata: metadata,
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
	branches, err := stores.SerializeMetadata(in, stores.MetadataOpts{
		Flatten: stores.MetadataFlattenFull,
	})
	if err != nil {
		return nil, fmt.Errorf("Error marshaling metadata: %s", err)
	}
	return store.EmitPlainFile(branches)
}

// EmitPlainFile returns the plaintext file's bytes corresponding to a sops
// runtime object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, item := range in[0] {
		if stores.IsComplexValue(item.Value) {
			return nil, fmt.Errorf("cannot use complex value in dotenv file; offending key %s", item.Key)
		}
		var line string
		if comment, ok := item.Key.(sops.Comment); ok {
			line = fmt.Sprintf("#%s\n", comment.Value)
		} else {
			value, ok := item.Value.(string)
			if !ok {
				value = stores.ValToString(item.Value)
			} else {
				value = strings.ReplaceAll(value, "\n", "\\n")
			}

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

// Deprecated: use stores.IsComplexValue() instead!
func IsComplexValue(v interface{}) bool {
	return stores.IsComplexValue(v)
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
