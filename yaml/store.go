package yaml

import (
	"fmt"
	"github.com/autrilla/yaml"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"strconv"
	"time"
)

// Store handles storage of YAML data
type Store struct {
}

func (store Store) mapSliceToTreeBranch(in yaml.MapSlice) sops.TreeBranch {
	branch := make(sops.TreeBranch, 0)
	for _, item := range in {
		branch = append(branch, sops.TreeItem{
			Key:   item.Key.(string),
			Value: store.yamlValueToTreeValue(item.Value),
		})
	}
	return branch
}

// Unmarshal takes a YAML document as input and unmarshals it into a sops tree, returning the tree
func (store Store) Unmarshal(in []byte) (sops.TreeBranch, error) {
	var data yaml.MapSlice
	if err := yaml.Unmarshal(in, &data); err != nil {
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
		branch = append(branch, yaml.MapItem{
			Key:   item.Key,
			Value: store.treeValueToYamlValue(item.Value),
		})
	}
	return branch
}

// Marshal takes a sops tree branch and marshals it into a yaml document
func (store Store) Marshal(tree sops.TreeBranch) ([]byte, error) {
	yamlMap := store.treeBranchToYamlMap(tree)
	out, err := yaml.Marshal(yamlMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
	}
	return out, nil
}

// MarshalWithMetadata takes a sops tree branch and metadata and marshals them into a yaml document
func (store Store) MarshalWithMetadata(tree sops.TreeBranch, metadata sops.Metadata) ([]byte, error) {
	yamlMap := store.treeBranchToYamlMap(tree)
	yamlMap = append(yamlMap, yaml.MapItem{Key: "sops", Value: metadata.ToMap()})
	out, err := yaml.Marshal(yamlMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
	}
	return out, nil
}

// UnmarshalMetadata takes a yaml document as a string and extracts sops' metadata from it
func (store *Store) UnmarshalMetadata(in []byte) (sops.Metadata, error) {
	var metadata sops.Metadata
	var ok bool
	data := make(map[interface{}]interface{})
	err := yaml.Unmarshal(in, &data)
	if err != nil {
		return metadata, fmt.Errorf("Error unmarshalling input yaml: %s", err)
	}
	if data, ok = data["sops"].(map[interface{}]interface{}); !ok {
		return metadata, sops.MetadataNotFound
	}
	metadata.MessageAuthenticationCode = data["mac"].(string)
	lastModified, err := time.Parse(time.RFC3339, data["lastmodified"].(string))
	if err != nil {
		return metadata, fmt.Errorf("Could not parse last modified date: %s", err)
	}
	metadata.LastModified = lastModified
	metadata.UnencryptedSuffix = data["unencrypted_suffix"].(string)
	if metadata.Version, ok = data["version"].(string); !ok {
		metadata.Version = strconv.FormatFloat(data["version"].(float64), 'f', -1, 64)
	}
	if k, ok := data["kms"].([]interface{}); ok {
		ks, err := store.kmsEntries(k)
		if err == nil {
			metadata.KeySources = append(metadata.KeySources, ks)
		}

	}

	if pgp, ok := data["pgp"].([]interface{}); ok {
		ks, err := store.pgpEntries(pgp)
		if err == nil {
			metadata.KeySources = append(metadata.KeySources, ks)
		}
	}
	return metadata, nil
}

func (store *Store) kmsEntries(in []interface{}) (sops.KeySource, error) {
	var keys []sops.MasterKey
	keysource := sops.KeySource{Name: "kms", Keys: keys}
	for _, v := range in {
		entry := v.(map[interface{}]interface{})
		key := &kms.MasterKey{}
		key.Arn = entry["arn"].(string)
		key.EncryptedKey = entry["enc"].(string)
		role, ok := entry["role"].(string)
		if ok {
			key.Role = role
		}
		creationDate, err := time.Parse(time.RFC3339, entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}

func (store *Store) pgpEntries(in []interface{}) (sops.KeySource, error) {
	var keys []sops.MasterKey
	keysource := sops.KeySource{Name: "pgp", Keys: keys}
	for _, v := range in {
		entry := v.(map[interface{}]interface{})
		key := &pgp.MasterKey{}
		key.Fingerprint = entry["fp"].(string)
		key.EncryptedKey = entry["enc"].(string)
		creationDate, err := time.Parse(time.RFC3339, entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}
