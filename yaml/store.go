package yaml

import (
	"fmt"
	"go.mozilla.org/sops/decryptor"
	"gopkg.in/yaml.v2"
)

type YAMLStore struct {
	Data map[interface{}]interface{}
}

func (store *YAMLStore) WalkValue(in interface{}, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, additionalAuthData)
	case int:
		return onLeaves(in, additionalAuthData)
	case bool:
		return onLeaves(in, additionalAuthData)
	case map[interface{}]interface{}:
		return store.WalkMap(in, additionalAuthData, onLeaves)
	case []interface{}:
		return store.WalkSlice(in, additionalAuthData, onLeaves)
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (store *YAMLStore) WalkSlice(in []interface{}, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) ([]interface{}, error) {
	for i, v := range in {
		newV, err := store.WalkValue(v, additionalAuthData, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
	}
	return in, nil
}

func (store *YAMLStore) WalkMap(in map[interface{}]interface{}, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) (map[interface{}]interface{}, error) {
	for k, v := range in {
		newV, err := store.WalkValue(v, additionalAuthData+k.(string)+":", onLeaves)
		if err != nil {
			return nil, err
		}
		in[k] = newV
	}
	return in, nil
}

func (store *YAMLStore) Load(data, key string) error {
	if err := yaml.Unmarshal([]byte(data), &store.Data); err != nil {
		return fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	delete(store.Data, "sops")
	_, err := store.WalkValue(store.Data, "", func(in interface{}, additionalAuthData string) (interface{}, error) {
		return decryptor.Decrypt(in.(string), key, []byte(additionalAuthData))
	})
	if err != nil {
		return fmt.Errorf("Error walking tree: %s", err)
	}
	return nil
}

func (store *YAMLStore) Dump(key string) (string, error) {
	_, err := store.WalkValue(store.Data, "", func(in interface{}, additionalAuthData string) (interface{}, error) {
		return decryptor.Encrypt(in, key, []byte(additionalAuthData))
	})
	if err != nil {
		return "", fmt.Errorf("Error walking tree: %s", err)
	}
	out, err := yaml.Marshal(store.Data)
	if err != nil {
		return "", fmt.Errorf("Error marshaling to yaml: %s", err)
	}
	return string(out), nil
}

type KMS struct {
	Arn        string `yaml:"arn"`
	Role       string `yaml:"role"`
	CreatedAt  string `yaml:"created_at"`
	EncodedKey string `yaml:"enc"`
}

type PGP struct {
	Fingerprint string `yaml:"fp"`
	CreatedAt   string `yaml:"created_at"`
	EncodedKey  string `yaml:"enc"`
}

type SopsMetadata struct {
	Mac               string
	Version           string
	KMS               []KMS
	PGP               []PGP
	LastModifed       string `yaml:"lastmodified"`
	UnencryptedSuffix string `yaml:"unencrypted_suffix"`
}

func (store YAMLStore) Metadata(in string) (SopsMetadata, error) {
	sops := SopsMetadata{}
	encoded := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(in), &encoded); err != nil {
		return sops, fmt.Errorf("Error unmarshalling input yaml: %s", err)
	}

	sopsYaml, err := yaml.Marshal(encoded["sops"])
	if err != nil {
		return sops, err
	}
	err = yaml.Unmarshal(sopsYaml, &sops)
	if err != nil {
		return sops, err
	}
	return sops, nil
}
