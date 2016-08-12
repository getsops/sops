package json

import (
	"encoding/json"
	"fmt"
	"go.mozilla.org/sops/decryptor"
)

type JSONStore struct {
	Data map[string]interface{}
}

func (store *JSONStore) WalkValue(in interface{}, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, additionalAuthData)
	case int:
		return onLeaves(in, additionalAuthData)
	case bool:
		return onLeaves(in, additionalAuthData)
	case map[string]interface{}:
		return store.WalkMap(in, additionalAuthData, onLeaves)
	case []interface{}:
		return store.WalkSlice(in, additionalAuthData, onLeaves)
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (store *JSONStore) WalkSlice(in []interface{}, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) ([]interface{}, error) {
	for i, v := range in {
		newV, err := store.WalkValue(v, additionalAuthData, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
	}
	return in, nil
}

func (store *JSONStore) WalkMap(in map[string]interface{}, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) (map[string]interface{}, error) {
	for k, v := range in {
		newV, err := store.WalkValue(v, additionalAuthData+k+":", onLeaves)
		if err != nil {
			return nil, err
		}
		in[k] = newV
	}
	return in, nil
}

func (store *JSONStore) Load(data, key string) error {
	err := json.Unmarshal([]byte(data), &store.Data)
	if err != nil {
		return fmt.Errorf("Could not unmarshal input data: %s", err)
	}

	_, err = store.WalkValue(store.Data, "", func(in interface{}, additionalAuthData string) (interface{}, error) {
		return decryptor.Decrypt(in.(string), key, []byte(additionalAuthData))
	})
	if err != nil {
		return fmt.Errorf("Error walking tree: %s", err)
	}
	return nil
}

func (store *JSONStore) Encrypt(in map[interface{}]interface{}) (string, error) {
	return "", nil
}

type SopsMetadata struct {
	Mac               string
	Version           string
	KMS               []map[string]string
	PGP               []map[string]string
	LastModifed       string `yaml:"lastmodified"`
	UnencryptedSuffix string `yaml:"unencrypted_suffix"`
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
