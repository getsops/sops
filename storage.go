package sops

import (
	"encoding/json"
	"fmt"

	"go.mozilla.org/sops/decryptor"
	"time"
)

type SopsMetadata struct {
	Mac               string
	Version           string
	KMS               []map[string]string
	PGP               []map[string]string
	LastModifed       time.Time `yaml:"lastmodified"`
	UnencryptedSuffix string    `yaml:"unencrypted_suffix"`
}

type JSONStore struct {
}

func (store JSONStore) Metadata(in string) (SopsMetadata, error) {
	var metadata SopsMetadata
	var encoded map[string]interface{}
	err := json.Unmarshal([]byte(in), &encoded)
	if err != nil {
		return metadata, err
	}

	sopsJson, err := json.Marshal(encoded["sops"])
	if err != nil {
		return metadata, err
	}

	err = json.Unmarshal(sopsJson, &metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}

func (store JSONStore) DecryptValue(in interface{}, key string) (interface{}, error) {
	switch in := in.(type) {
	case map[string]interface{}:
		return store.DecryptMap(in, key)
	case string:
		k, err := decryptor.Decrypt(in, key)
		if err != nil {
			return nil, fmt.Errorf("Could not decrypt \"%s\"", in)
		}
		return k, nil
	default:
	}
	return nil, fmt.Errorf("Value %s is of unknown type", in)
}

func (store JSONStore) DecryptMap(in map[string]interface{}, key string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for i, v := range in {
		v, err := store.DecryptValue(v, key)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}

	return out, nil
}

func (store JSONStore) Decrypt(in, key string) (string, error) {
	var encoded map[string]interface{}
	err := json.Unmarshal([]byte(in), &encoded)
	if err != nil {
		return "", err
	}

	v, err := store.DecryptMap(encoded, key)
	if err != nil {
		return "", err
	}
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (store JSONStore) Encrypt(in map[interface{}]interface{}) (string, error) {
	return "", nil
}
