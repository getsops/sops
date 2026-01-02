package toml //import "github.com/getsops/sops/v3/stores/toml"

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/stores"

	"github.com/pelletier/go-toml/v2"
)

// Store handles storage of TOML data.
type Store struct {
	config *config.TOMLStoreConfig
}

func NewStore(c *config.TOMLStoreConfig) *Store {
	return &Store{config: c}
}

func (store *Store) Name() string {
	return "toml"
}

var errUnexpectedValue = errors.New("unexpected value")

// mapToTreeBranch converts a map[string]any to a sops.TreeBranch.
func mapToTreeBranch(m map[string]any) (sops.TreeBranch, error) {
	// Separate keys by type: simple values first, then complex types
	// (tables/arrays).
	var simpleKeys, complexKeys []string
	for k, v := range m {
		switch v.(type) {
		case map[string]any:
			complexKeys = append(complexKeys, k)
		case []any:
			// Check if it's an array of tables
			if arr, ok := v.([]any); ok && len(arr) > 0 {
				if _, isMap := arr[0].(map[string]any); isMap {
					complexKeys = append(complexKeys, k)
				} else {
					simpleKeys = append(simpleKeys, k)
				}
			} else {
				simpleKeys = append(simpleKeys, k)
			}
		default:
			simpleKeys = append(simpleKeys, k)
		}
	}

	// Sort each group independently.
	sortKeysNaturally(simpleKeys)
	sortKeysNaturally(complexKeys)

	// Combine: simple values first, then complex types.
	keys := append(simpleKeys, complexKeys...)

	var branch sops.TreeBranch
	for _, k := range keys {
		v := m[k]
		value, err := anyToTreeItemValue(v)
		if err != nil {
			return nil, fmt.Errorf("mapToTreeBranch: %w - %v, %s", errUnexpectedValue, v, k)
		}
		branch = append(branch, sops.TreeItem{
			Key:   k,
			Value: value,
		})
	}
	return branch, nil
}

// sortKeysNaturally sorts keys lexicographically for deterministic output.
func sortKeysNaturally(keys []string) {
	sort.Strings(keys)
}

// anyToTreeItemValue converts an any value from TOML unmarshaling
// to a sops TreeItem value.
func anyToTreeItemValue(v any) (any, error) {
	switch val := v.(type) {
	case map[string]any:
		return mapToTreeBranch(val)
	case []any:
		// Check if it's an array of maps (array of tables in TOML).
		if len(val) > 0 {
			if _, ok := val[0].(map[string]any); ok {
				// Yes, it's an array of tables.
				var branches []any
				for _, item := range val {
					if m, ok := item.(map[string]any); ok {
						branch, err := mapToTreeBranch(m)
						if err != nil {
							return nil, err
						}
						branches = append(branches, branch)
					} else {
						return nil, fmt.Errorf("anyToTreeItemValue: expected map in array, got %T", item)
					}
				}
				return branches, nil
			}
		}
		return val, nil
	default:
		return val, nil
	}
}

// treeBranchToMap converts a sops.TreeBranch to a map[string]any.
func treeBranchToMap(branch sops.TreeBranch) (map[string]any, error) {
	m := make(map[string]any)
	for _, item := range branch {
		key, ok := item.Key.(string)
		if !ok {
			// Skip non-string keys (like comments).
			continue
		}
		value, err := treeItemValueToInterface(item.Value)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}

// treeItemValueToInterface converts a sops TreeItem value to an any
// suitable for TOML marshaling.
func treeItemValueToInterface(value any) (any, error) {
	switch val := value.(type) {
	case sops.TreeBranch:
		return treeBranchToMap(val)
	case []any:
		// Check if it's an array of TreeBranches.
		if len(val) > 0 {
			if _, ok := val[0].(sops.TreeBranch); ok {
				var result []any
				for _, item := range val {
					if branch, ok := item.(sops.TreeBranch); ok {
						m, err := treeBranchToMap(branch)
						if err != nil {
							return nil, err
						}
						result = append(result, m)
					} else {
						return nil, fmt.Errorf("treeItemValueToInterface: expected TreeBranch in array, got %T", item)
					}
				}
				return result, nil
			}
		}
		return val, nil
	default:
		return val, nil
	}
}

// LoadEncryptedFile loads the contents of an encrypted toml file onto a
// sops.Tree runtime object.
func (s *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	var data map[string]any
	if err := toml.Unmarshal(in, &data); err != nil {
		return sops.Tree{}, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	// Because we don't know what fields the input file will have, we have to
	// load the file in two steps.
	//
	// First, we load the file's metadata, the structure of which is known.
	metadataHolder := stores.SopsFile{}
	if err := toml.Unmarshal(in, &metadataHolder); err != nil {
		return sops.Tree{}, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	if metadataHolder.Metadata == nil {
		return sops.Tree{}, sops.MetadataNotFound
	}

	metadata, err := metadataHolder.Metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	// Second, we load the rest of the file's contents into a generic tree
	// structure.
	branch, err := mapToTreeBranch(data)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("error transforming toml data: %w", err)
	}

	return sops.Tree{
		Branches: sops.TreeBranches{branch},
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads the contents of a plaintext toml file onto a
// sops.Tree runtime object.
func (s *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	var data map[string]any
	if err := toml.Unmarshal(in, &data); err != nil {
		return nil, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	branch, err := mapToTreeBranch(data)
	if err != nil {
		return nil, fmt.Errorf("error transforming toml data: %w", err)
	}

	return sops.TreeBranches{branch}, nil
}

var errTOMLUniqueDocument = errors.New("toml can only contain 1 document")

// EmitEncryptedFile returns the encrypted bytes of the toml file corresponding to a
// sops.Tree runtime object.
func (s *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	if len(in.Branches) != 1 {
		return nil, errTOMLUniqueDocument
	}

	data, err := treeBranchToMap(in.Branches[0])
	if err != nil {
		return nil, fmt.Errorf("error converting tree branch: %w", err)
	}

	data["sops"] = stores.MetadataFromInternal(in.Metadata)

	return s.marshalTOML(data)
}

// EmitPlainFile returns the plaintext bytes of the toml file corresponding to a
// sops.TreeBranches runtime object.
func (s *Store) EmitPlainFile(branches sops.TreeBranches) ([]byte, error) {
	if len(branches) != 1 {
		return nil, errTOMLUniqueDocument
	}

	data, err := treeBranchToMap(branches[0])
	if err != nil {
		return nil, fmt.Errorf("emit plain file: %w", err)
	}

	return s.marshalTOML(data)
}

// marshalTOML marshals data to TOML with custom formatting
func (s *Store) marshalTOML(data any) ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := toml.NewEncoder(&buf)
	encoder.SetIndentTables(true)
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("error marshalling toml: %w", err)
	}

	// Replace single quotes with double quotes for string values.
	result := bytes.ReplaceAll(buf.Bytes(), []byte("'"), []byte("\""))
	return result, nil
}

// EmitValue returns bytes corresponding to a single encoded value
// in a generic any object.
func (s *Store) EmitValue(v any) ([]byte, error) {
	switch val := v.(type) {
	case sops.TreeBranch:
		data, err := treeBranchToMap(val)
		if err != nil {
			return nil, fmt.Errorf("emit value: %w", err)
		}

		return s.marshalTOML(data)
	case string:
		// For strings, return quoted value.
		return []byte(fmt.Sprintf("%q", val)), nil
	default:
		// For simple values, format them appropriately.
		return []byte(fmt.Sprintf("%v", val)), nil
	}
}

// EmitExample returns the bytes corresponding to an example complex tree.
func (s *Store) EmitExample() []byte {
	bytes, err := s.EmitPlainFile(stores.ExampleComplexTree.Branches)
	if err != nil {
		panic(err)
	}

	return bytes
}

// HasSopsTopLevelKey checks whether a top-level "sops" key exists.
func (store *Store) HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	return stores.HasSopsTopLevelKey(branch)
}
