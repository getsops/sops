package json

import (
	"encoding/json"
	"fmt"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"time"
)

type JSONStore struct {
}

func (store JSONStore) Load(in string) (sops.TreeBranch, error) {
	var branch sops.TreeBranch
	err := json.Unmarshal([]byte(in), branch)
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

func (store JSONStore) Dump(tree sops.TreeBranch) (string, error) {
	out, err := json.Marshal(tree)
	if err != nil {
		return "", fmt.Errorf("Error marshaling to json: %s", err)
	}
	return string(out), nil
}

func (store JSONStore) DumpWithMetadata(tree sops.TreeBranch, metadata sops.Metadata) (string, error) {
	tree = append(tree, sops.TreeItem{Key: "sops", Value: metadata.ToMap()})
	out, err := json.Marshal(tree)
	if err != nil {
		return "", fmt.Errorf("Error marshaling to json: %s", err)
	}
	return string(out), nil
}

func (store JSONStore) LoadMetadata(in string) (sops.Metadata, error) {
	var metadata sops.Metadata
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(in), &data); err != nil {
		return metadata, fmt.Errorf("Error unmarshalling input json: %s", err)
	}
	data = data["sops"].(map[string]interface{})
	metadata.MessageAuthenticationCode = data["mac"].(string)
	lastModified, err := time.Parse(sops.DateFormat, data["lastmodified"].(string))
	if err != nil {
		return metadata, fmt.Errorf("Could not parse last modified date: %s", err)
	}
	metadata.LastModified = lastModified
	metadata.UnencryptedSuffix = data["unencrypted_suffix"].(string)
	metadata.Version = data["version"].(string)
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

func (store JSONStore) kmsEntries(in []interface{}) (sops.KeySource, error) {
	var keys []sops.MasterKey
	keysource := sops.KeySource{Name: "kms", Keys: keys}
	for _, v := range in {
		entry := v.(map[interface{}]interface{})
		key := &kms.KMSMasterKey{}
		key.Arn = entry["arn"].(string)
		key.EncryptedKey = entry["enc"].(string)
		role, ok := entry["role"].(string)
		if ok {
			key.Role = role
		}
		creationDate, err := time.Parse(sops.DateFormat, entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}

func (store JSONStore) pgpEntries(in []interface{}) (sops.KeySource, error) {
	var keys []sops.MasterKey
	keysource := sops.KeySource{Name: "pgp", Keys: keys}
	for _, v := range in {
		entry := v.(map[interface{}]interface{})
		key := &pgp.GPGMasterKey{}
		key.Fingerprint = entry["fp"].(string)
		key.EncryptedKey = entry["enc"].(string)
		creationDate, err := time.Parse(sops.DateFormat, entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}
