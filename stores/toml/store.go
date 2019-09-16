package toml //import "go.mozilla.org/sops/stores/toml"

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/stores"
)

// Store handles storage of TOML data.
type Store struct {
}

// LoadEncryptedFile loads an encrypted TOML secrets file onto a sops.Tree object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	// Load metadata.
	metadataHolder := stores.SopsFile{}
	err := toml.Unmarshal(in, &metadataHolder)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("Error unmarshalling metadata: %s", err)
	}
	if metadataHolder.Metadata == nil {
		return sops.Tree{}, sops.MetadataNotFound
	}
	metadata, err := metadataHolder.Metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, fmt.Errorf("Error parsing TOML metadata: %s", err)
	}

	// Load data.
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}

	return sops.Tree{Metadata: metadata, Branches: branches}, nil
}

// LoadPlainFile loads a plaintext TOML file's bytes onto a sops.TreeBranches object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	// Load TOML data into unordered map of nested interfaces and
	// use the TOML metadata.Keys() to get back the original key order.
	var unordered map[string]interface{}
	md, err := toml.Decode(string(in), &unordered)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling TOML data: %s", err)
	}
	delete(unordered, "sops")

	r := tomlReader(md.Keys())
	branch := r.readToTreeBranch(unordered)
	return sops.TreeBranches{branch}, err
}

// EmitEncryptedFile produces an encrypted TOML file's bytes from its corresponding sops.Tree object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	var b bytes.Buffer
	for _, branch := range in.Branches {
		if err := printTreeBranchInTOML(&b, branch); err != nil {
			return nil, fmt.Errorf("Error marshalling data to TOML: %s", err)
		}
	}
	b.WriteByte('\n')

	metadata := stores.MetadataFromInternal(in.Metadata)
	var customMetadataHolder struct { // custom stores.SopsFile
		Metadata struct {
			*stores.Metadata

			// Omitempty: github.com/BurntSushi/toml doesn't support ",omitempty" tag for zero values
			ShamirThreshold interface{} `toml:"shamir_threshold"`
		} `toml:"sops"`
	}
	customMetadataHolder.Metadata.Metadata = &metadata
	if metadata.ShamirThreshold != 0 {
		customMetadataHolder.Metadata.ShamirThreshold = metadata.ShamirThreshold
	}

	enc := toml.NewEncoder(&b)
	if err := enc.Encode(customMetadataHolder); err != nil {
		return nil, fmt.Errorf("Error marshalling metadata to TOML: %s", err)
	}
	b.WriteByte('\n')

	return b.Bytes(), nil
}

// EmitPlainFile produces plaintext TOML file's bytes from its corresponding sops.TreeBranches object
func (store *Store) EmitPlainFile(branches sops.TreeBranches) ([]byte, error) {
	var b bytes.Buffer
	for _, branch := range branches {
		if err := printTreeBranchInTOML(&b, branch); err != nil {
			return nil, fmt.Errorf("Error encoding data: %s", err)
		}
	}
	return b.Bytes(), nil
}

// EmitValue returns a single value encoded in a generic interface{} as bytes
func (store *Store) EmitValue(value interface{}) ([]byte, error) {
	var b bytes.Buffer
	switch v := value.(type) {
	case sops.TreeBranch:
		if err := printTreeBranchInTOML(&b, v); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case sops.TreeItem:
		if err := printTreeItemInTOML(&b, v); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case []interface{}:
		return nil, fmt.Errorf("Error extracting array of %v items. Please, access an individual item.", len(v))
	default:
		if err := printTreeBranch(&b, sops.TreeBranch{sops.TreeItem{Key: "_delete", Value: v}}); err != nil {
			return nil, err
		}
		return b.Bytes()[len("_delete = "):], nil
	}
}

// EmitExample returns the plaintext TOML file bytes corresponding to the SimpleTree example
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleComplexTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}
