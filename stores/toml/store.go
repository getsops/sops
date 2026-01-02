package toml //import "github.com/getsops/sops/v3/stores/toml"

import (
	"errors"
	"fmt"
	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/stores"
	"sort"

	"github.com/pelletier/go-toml"
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

// positionKey is necessary to keep the order when we append in branches.
type positionKey struct {
	position int
	key      string
}

type byPosition []positionKey

func (b byPosition) Less(i, j int) bool { return b[i].position < b[j].position }
func (b byPosition) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byPosition) Len() int           { return len(b) }

var errUnexpectedValue = errors.New("unexpected value")

func tomlTreeToTreeBranch(tr *toml.Tree) (sops.TreeBranch, error) {
	treeItems := make(map[string]sops.TreeItem)
	pks := []positionKey{}

	for k, v := range tr.Values() {
		switch node := v.(type) {
		case []*toml.Tree:
			maxPosition := 0

			var treeBranches []interface{}

			for _, item := range node {
				if item.Position().Line > maxPosition {
					maxPosition = item.Position().Line
				}

				branch, errT := tomlTreeToTreeBranch(item)
				if errT != nil {
					return nil, fmt.Errorf("tomlTreeToTreeBranch: %w - %v, %s", errUnexpectedValue, v, k)
				}

				treeBranches = append(treeBranches, branch)
			}

			pks = append(pks, positionKey{position: maxPosition, key: k})
			treeItems[k] = sops.TreeItem{
				Key:   k,
				Value: treeBranches,
			}
		case *toml.Tree:
			pks = append(pks, positionKey{position: node.Position().Line, key: k})

			branch, errT := tomlTreeToTreeBranch(node)
			if errT != nil {
				return nil, fmt.Errorf("tomlTreeToTreeBranch: %w - %v, %s", errUnexpectedValue, v, k)
			}

			treeItems[k] = sops.TreeItem{
				Key:   k,
				Value: branch,
			}
		case *toml.PubTOMLValue:
			pks = append(pks, positionKey{position: node.Position().Line, key: k})
			treeItems[k] = sops.TreeItem{
				Key:   k,
				Value: node.Value(),
			}
		default:
			return nil, fmt.Errorf("tomlTreeToTreeBranch: %w - %v, %s", errUnexpectedValue, v, k)
		}
	}

	sort.Sort(byPosition(pks))

	var br sops.TreeBranch
	for _, pk := range pks {
		br = append(br, treeItems[pk.key])
	}

	return br, nil
}

func treeBranchToTOMLTree(stree sops.TreeBranch) (*toml.Tree, error) {
	ttree := &toml.PubTree{}
	values := make(map[string]interface{})

	for _, treeItem := range stree {
		var errT error

		values[treeItem.Key.(string)], errT = treeItemValueToTOML(treeItem.Value)
		if errT != nil {
			return nil, errT
		}
	}

	ttree.SetValues(values)

	return ttree, nil
}

func treeItemValueToTOML(treeItemValue interface{}) (interface{}, error) {
	switch treeItemValueTyped := treeItemValue.(type) {
	case sops.TreeBranch:
		return treeBranchToTOMLTree(treeItemValueTyped)

	case []interface{}:
		switch treeItemValueTyped[0].(type) {
		case sops.TreeBranch:
			var array []*toml.Tree

			for _, itm := range treeItemValueTyped {
				tb := itm.(sops.TreeBranch)

				tr, errT := treeBranchToTOMLTree(tb)
				if errT != nil {
					return nil, errT
				}

				array = append(array, tr)
			}

			return array, nil
		default:
			val := &toml.PubTOMLValue{}
			val.SetValue(treeItemValueTyped)

			return val, nil
		}

	default:
		val := &toml.PubTOMLValue{}
		val.SetValue(treeItemValueTyped)

		return val, nil
	}
}

// LoadEncryptedFile loads the contents of an encrypted toml file onto a
// sops.Tree runtime object.
func (s *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	data, err := toml.LoadBytes(in)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	// Because we don't know what fields the input file will have, we have to
	// load the file in two steps.
	// First, we load the file's metadata, the structure of which is known.
	metadataHolder := stores.SopsFile{}
	if err := data.Unmarshal(&metadataHolder); err != nil {
		return sops.Tree{}, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	if metadataHolder.Metadata == nil {
		return sops.Tree{}, sops.MetadataNotFound
	}

	metadata, err := metadataHolder.Metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	branch, errT := tomlTreeToTreeBranch(data)
	if errT != nil {
		return sops.Tree{}, fmt.Errorf("error transforming toml Tree: %w", err)
	}

	return sops.Tree{
		Branches: sops.TreeBranches{branch},
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads the contents of a plaintext toml file onto a
// sops.Tree runtime object.
func (s *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	tomlTree, err := toml.LoadBytes(in)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling input toml: %w", err)
	}

	branch, errT := tomlTreeToTreeBranch(tomlTree)
	if errT != nil {
		return nil, fmt.Errorf("error transforming toml Tree: %w", err)
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

	tomlTree, err := treeBranchToTOMLTree(in.Branches[0])
	if err != nil {
		return nil, errTOMLUniqueDocument
	}

	values := tomlTree.Values()
	values["sops"] = stores.MetadataFromInternal(in.Metadata)
	tomlTree.SetValues(values)

	return []byte(tomlTree.String()), nil
}

// EmitPlainFile returns the plaintext bytes of the toml file corresponding to a
// sops.TreeBranches runtime object.
func (s *Store) EmitPlainFile(branches sops.TreeBranches) ([]byte, error) {
	if len(branches) != 1 {
		return nil, errTOMLUniqueDocument
	}

	tomlTree, err := treeBranchToTOMLTree(branches[0])
	if err != nil {
		return nil, fmt.Errorf("emit plain file: %w", err)
	}

	return tomlTree.Marshal()
}

// EmitValue returns bytes corresponding to a single encoded value
// in a generic interface{} object.
func (s *Store) EmitValue(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case sops.TreeBranch:
		tomlTree, err := treeBranchToTOMLTree(v)
		if err != nil {
			return nil, fmt.Errorf("emit plain file: %w", err)
		}

		return tomlTree.Marshal()
	default:
		str, err := toml.ValueStringRepresentation(v, "", "", toml.OrderPreserve, false)
		if err != nil {
			return nil, fmt.Errorf("emit plain file: %w", err)
		}
		return []byte(str), nil
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
