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

// Unmarshal takes a YAML document as input and unmarshals it into a sops tree, returning the tree
func (store Store) Unmarshal(in []byte) (sops.TreeBranch, error) {
	var data yaml.MapSlice
	if err := (yaml.CommentUnmarshaler{}).Unmarshal(in, &data); err != nil {
		return nil, fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	for i, item := range data {
		if item.Key == "sops" {
			data = append(data[:i], data[i+1:]...)
		}
	}
	return store.mapSliceToTreeBranch(data), nil
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

// Marshal takes a sops tree branch and marshals it into a yaml document
func (store Store) Marshal(tree sops.TreeBranch) ([]byte, error) {
	yamlMap := store.treeBranchToYamlMap(tree)
	out, err := (&yaml.YAMLMarshaler{Indent: 4}).Marshal(yamlMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
	}
	return out, nil
}

// MarshalWithMetadata takes a sops tree branch and metadata and marshals them into a yaml document
func (store Store) MarshalWithMetadata(tree sops.TreeBranch, metadata sops.Metadata) ([]byte, error) {
	yamlMap := store.treeBranchToYamlMap(tree)
	yamlMap = append(yamlMap, yaml.MapItem{Key: "sops", Value: stores.MetadataFromInternal(metadata)})
	out, err := (&yaml.YAMLMarshaler{Indent: 4}).Marshal(yamlMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
	}
	return out, nil
}

// MarshalValue takes any value and marshals it into a yaml document
func (store Store) MarshalValue(v interface{}) ([]byte, error) {
	v = store.treeValueToYamlValue(v)
	return (&yaml.YAMLMarshaler{Indent: 4}).Marshal(v)
}

// UnmarshalMetadata takes a yaml document as a string and extracts sops' metadata from it
func (store *Store) UnmarshalMetadata(in []byte) (sops.Metadata, error) {
	file := stores.SopsFile{}
	err := yaml.Unmarshal(in, &file)
	if err != nil {
		return sops.Metadata{}, fmt.Errorf("Error unmarshalling input yaml: %s", err)
	}
	if file.Metadata == nil {
		return sops.Metadata{}, sops.MetadataNotFound
	}
	return file.Metadata.ToInternal()
}
