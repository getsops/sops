package sops

import (
	"errors"
	"github.com/autrilla/sops/decryptor"
	"gopkg.in/yaml.v2"
	"time"
)

type Store interface {
	Metadata(in string) (SopsMetadata, error)
	Decrypt(in, key string) (string, error)
	Encrypt(in []map[interface{}]interface{}) (string, error)
}

type YAMLStore struct {
}
type JSONStore struct {
}

type SopsMetadata struct {
	Mac               string
	Version           string
	KMS               []map[string]string
	PGP               []map[string]string
	LastModifed       time.Time `yaml:"lastmodified"`
	UnencryptedSuffix string    `yaml:"unencrypted_suffix"`
}

func (store YAMLStore) DecryptValue(in interface{}, decryptionKey string) interface{} {
	switch in := in.(type) {
	case string:
		v, err := decryptor.Decrypt(in, decryptionKey)
		if err != nil {
			return nil
		}
		return v
	case map[interface{}]interface{}:
		return store.DecryptMap(in, decryptionKey)
	case yaml.MapSlice:
		return store.DecryptMapSlice(in, decryptionKey)
	case []interface{}:
		return store.DecryptSlice(in, decryptionKey)
	default:
	}
	return nil
}

func (store YAMLStore) DecryptMap(in map[interface{}]interface{}, decryptionKey string) map[interface{}]interface{} {
	branch := make(map[interface{}]interface{})
	for k, v := range in {
		branch[k] = store.DecryptValue(v, decryptionKey)
	}
	return branch
}

func (store YAMLStore) DecryptSlice(in []interface{}, decryptionKey string) []interface{} {
	list := make([]interface{}, len(in))
	for i, v := range in {
		list[i] = store.DecryptValue(v, decryptionKey)
	}
	return list
}

func (store YAMLStore) DecryptMapSlice(in yaml.MapSlice, decryptionKey string) yaml.MapSlice {
	out := make(yaml.MapSlice, len(in))
	for i, v := range in {
		plaintext := store.DecryptValue(v.Value, decryptionKey)
		out[i] = yaml.MapItem{v.Key, plaintext}
	}
	return out
}

func (store YAMLStore) Decrypt(in, key string) (string, error) {
	encoded := make(yaml.MapSlice, 0)
	if err := yaml.Unmarshal([]byte(in), &encoded); err != nil {
		return "", errors.New("Error unmarshalling input yaml")
	}

	decoded := store.DecryptMapSlice(encoded, key)
	out, err := yaml.Marshal(decoded)
	return string(out), err
}

func (store YAMLStore) Metadata(in string) (SopsMetadata, error) {
	sops := SopsMetadata{}
	encoded := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(in), &encoded); err != nil {
		return sops, errors.New("Error unmarshalling input yaml")
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
