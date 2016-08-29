package json

import (
	"encoding/json"
	"fmt"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"strconv"
	"time"
)

// Store handles storage of JSON data. It's not finished yet, and therefore you should not use it.
type Store struct {
}

// Unmarshal takes an input json string and returns a sops tree branch
func (store Store) Unmarshal(in []byte) (sops.TreeBranch, error) {
	var branch sops.TreeBranch
	err := json.Unmarshal(in, branch)
	if err != nil {
		return branch, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	for i, item := range branch {
		if item.Key == "sops" {
			branch = append(branch[:i], branch[i+1:]...)
		}
	}
	return branch, nil
}

// Marshal takes a sops tree branch and returns a json formatted string
func (store Store) Marshal(tree sops.TreeBranch) ([]byte, error) {
	out, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	return out, nil
}

// MarshalWithMetadata takes a sops tree branch and sops metadata and marshals them to json.
func (store Store) MarshalWithMetadata(tree sops.TreeBranch, metadata sops.Metadata) ([]byte, error) {
	tree = append(tree, sops.TreeItem{Key: "sops", Value: metadata.ToMap()})
	out, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to json: %s", err)
	}
	return out, nil
}

// UnmarshalMetadata takes a json string and extracts sops' metadata from it
func (store Store) UnmarshalMetadata(in []byte) (sops.Metadata, error) {
	var ok bool
	var metadata sops.Metadata
	data := make(map[string]interface{})
	if err := json.Unmarshal(in, &data); err != nil {
		return metadata, fmt.Errorf("Error unmarshalling input json: %s", err)
	}
	if data, ok = data["sops"].(map[string]interface{}); !ok {
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

func (store Store) kmsEntries(in []interface{}) (sops.KeySource, error) {
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

func (store Store) pgpEntries(in []interface{}) (sops.KeySource, error) {
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
