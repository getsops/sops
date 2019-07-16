package yaml //import "go.mozilla.org/sops/stores/yaml"

import (
	"fmt"

	"github.com/mozilla-services/yaml"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/stores"
)

// Store handles storage of YAML data
type Store struct {
}

func (store Store) mapSliceToTreeBranch(in yaml.MapSlice) sops.TreeBranch {
	branch := make(sops.TreeBranch, 0)
	for _, item := range in {
		if comment, ok := item.Key.(yaml.Comment); ok {
			// Convert the yaml comment to a generic sops comment
			branch = append(branch, sops.TreeItem{
				Key: sops.Comment{
					Value: comment.Value,
				},
				Value: nil,
			})
		} else {
			branch = append(branch, sops.TreeItem{
				Key:   item.Key,
				Value: store.yamlValueToTreeValue(item.Value),
			})
		}
	}
	return branch
}

func (store Store) yamlValueToTreeValue(in interface{}) interface{} {
	switch in := in.(type) {
	case map[interface{}]interface{}:
		return store.yamlMapToTreeBranch(in)
	case yaml.MapSlice:
		return store.mapSliceToTreeBranch(in)
	case []interface{}:
		return store.yamlSliceToTreeValue(in)
	case yaml.Comment:
		return sops.Comment{Value: in.Value}
	default:
		return in
	}
}

func (store *Store) yamlSliceToTreeValue(in []interface{}) []interface{} {
	for i, v := range in {
		in[i] = store.yamlValueToTreeValue(v)
	}
	return in
}

func (store *Store) yamlMapToTreeBranch(in map[interface{}]interface{}) sops.TreeBranch {
	branch := make(sops.TreeBranch, 0)
	for k, v := range in {
		branch = append(branch, sops.TreeItem{
			Key:   k.(string),
			Value: store.yamlValueToTreeValue(v),
		})
	}
	return branch
}

func (store Store) treeValueToYamlValue(in interface{}) interface{} {
	switch in := in.(type) {
	case sops.TreeBranch:
		return store.treeBranchToYamlMap(in)
	case sops.Comment:
		return yaml.Comment{in.Value}
	case []interface{}:
		var out []interface{}
		for _, v := range in {
			out = append(out, store.treeValueToYamlValue(v))
		}
		return out
	default:
		return in
	}
}

func (store Store) treeBranchToYamlMap(in sops.TreeBranch) yaml.MapSlice {
	branch := make(yaml.MapSlice, 0)
	for _, item := range in {
		if comment, ok := item.Key.(sops.Comment); ok {
			branch = append(branch, yaml.MapItem{
				Key:   store.treeValueToYamlValue(comment),
				Value: nil,
			})
		} else {
			branch = append(branch, yaml.MapItem{
				Key:   item.Key,
				Value: store.treeValueToYamlValue(item.Value),
			})
		}
	}
	return branch
}

// LoadEncryptedFile loads the contents of an encrypted yaml file onto a
// sops.Tree runtime object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	var data []yaml.MapSlice
	if err := (yaml.CommentUnmarshaler{}).UnmarshalDocuments(in, &data); err != nil {
		return sops.Tree{}, fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	// Because we don't know what fields the input file will have, we have to
	// load the file in two steps.
	// First, we load the file's metadata, the structure of which is known.
	metadataHolder := stores.SopsFile{}
	err := yaml.Unmarshal(in, &metadataHolder)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("Error unmarshalling input yaml: %s", err)
	}
	if metadataHolder.Metadata == nil {
		return sops.Tree{}, sops.MetadataNotFound
	}
	metadata, err := metadataHolder.Metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, err
	}
	var branches sops.TreeBranches
	for _, doc := range data {
		for i, item := range doc {
			if item.Key == "sops" { // Erase
				doc = append(doc[:i], doc[i+1:]...)
			}
		}
		branches = append(branches, store.mapSliceToTreeBranch(doc))
	}
	return sops.Tree{
		Branches: branches,
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads the contents of a plaintext yaml file onto a
// sops.Tree runtime obejct
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	var data []yaml.MapSlice
	if err := (yaml.CommentUnmarshaler{}).UnmarshalDocuments(in, &data); err != nil {
		return nil, fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}

	var branches sops.TreeBranches
	for _, doc := range data {
		branches = append(branches, store.mapSliceToTreeBranch(doc))
	}
	return branches, nil
}

// EmitEncryptedFile returns the encrypted bytes of the yaml file corresponding to a
// sops.Tree runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	out := []byte{}
	for i, branch := range in.Branches {
		if i > 0 {
			out = append(out, "---\n"...)
		}
		yamlMap := store.treeBranchToYamlMap(branch)
		yamlMap = append(yamlMap, yaml.MapItem{Key: "sops", Value: stores.MetadataFromInternal(in.Metadata)})
		tout, err := (&yaml.YAMLMarshaler{Indent: 4}).Marshal(yamlMap)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
		}
		out = append(out, tout...)
	}
	return out, nil
}

// EmitPlainFile returns the plaintext bytes of the yaml file corresponding to a
// sops.TreeBranches runtime object
func (store *Store) EmitPlainFile(branches sops.TreeBranches) ([]byte, error) {
	var out []byte
	for i, branch := range branches {
		if i > 0 {
			out = append(out, "---\n"...)
		}
		yamlMap := store.treeBranchToYamlMap(branch)
		tmpout, err := (&yaml.YAMLMarshaler{Indent: 4}).Marshal(yamlMap)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
		}
		out = append(out[:], tmpout[:]...)
	}
	return out, nil
}

// EmitValue returns bytes corresponding to a single encoded value
// in a generic interface{} object
func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	v = store.treeValueToYamlValue(v)
	return (&yaml.YAMLMarshaler{Indent: 4}).Marshal(v)
}

// EmitExample returns the bytes corresponding to an example complex tree
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleComplexTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}
