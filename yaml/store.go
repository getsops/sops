package yaml

import (
	"crypto/sha512"
	"fmt"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/decryptor"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"gopkg.in/yaml.v2"
	"strconv"
	"time"
)

type YAMLStore struct {
	Data     yaml.MapSlice
	metadata sops.Metadata
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
	case yaml.MapSlice:
		return store.WalkMapSlice(in, additionalAuthData, onLeaves)
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

func (store *YAMLStore) WalkMapSlice(in yaml.MapSlice, additionalAuthData string, onLeaves func(interface{}, string) (interface{}, error)) (yaml.MapSlice, error) {
	for i, item := range in {
		newV, err := store.WalkValue(item.Value, additionalAuthData+item.Key.(string)+":", onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
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

func (store *YAMLStore) LoadUnencrypted(data string) error {
	if err := yaml.Unmarshal([]byte(data), &store.Data); err != nil {
		return fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	return nil
}

func toBytes(in interface{}) ([]byte, error) {
	switch in := in.(type) {
	case string:
		return []byte(in), nil
	case int:
		return []byte(strconv.Itoa(in)), nil
	case float64:
		return []byte(strconv.FormatFloat(in, 'f', -1, 64)), nil
	case bool:
		return []byte(strconv.FormatBool(in)), nil
	case []byte:
		return in, nil
	default:
		return nil, fmt.Errorf("Could not convert unknown type %T to bytes", in)
	}
}

func (store *YAMLStore) Load(data, key string) error {
	if err := yaml.Unmarshal([]byte(data), &store.Data); err != nil {
		return fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	err := store.LoadMetadata(data)
	if err != nil {
		return fmt.Errorf("Could not get metadata from YAML file: %s", err)
	}
	for i, v := range store.Data {
		if v.Key == "sops" {
			store.Data = append(store.Data[:i], store.Data[i+1:]...)
			break
		}
	}
	hash := sha512.New()
	_, err = store.WalkValue(store.Data, "", func(in interface{}, additionalAuthData string) (interface{}, error) {
		v, err := decryptor.Decrypt(in.(string), key, []byte(additionalAuthData))
		if err != nil {
			return nil, err
		}
		bytes, err := toBytes(v)
		if err != nil {
			return nil, err
		}
		hash.Write(bytes)
		return v, err
	})
	if err != nil {
		return fmt.Errorf("Error walking tree: %s", err)
	}
	originalMac, err := decryptor.Decrypt(store.metadata.MessageAuthenticationCode, key, []byte(store.metadata.LastModified.Format("2006-01-02T15:04:05Z")))
	if err != nil {
		return fmt.Errorf("Error decrypting MAC: %s", err)
	}
	macHex := fmt.Sprintf("%X", hash.Sum(nil))
	if macHex != originalMac.(string) {
		return sops.MacMismatch
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

func (store *YAMLStore) DumpUnencrypted() (string, error) {
	out, err := yaml.Marshal(store.Data)
	if err != nil {
		return "", fmt.Errorf("Error marshaling to yaml: %s", err)
	}
	return string(out), nil
}

func (store *YAMLStore) LoadMetadata(in string) error {
	data := make(map[interface{}]interface{})
	encoded := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(in), &encoded); err != nil {
		return fmt.Errorf("Error unmarshalling input yaml: %s", err)
	}

	sopsYaml, err := yaml.Marshal(encoded["sops"])
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(sopsYaml, &data)
	if err != nil {
		return err
	}
	store.metadata.MessageAuthenticationCode = data["mac"].(string)
	lastModified, err := time.Parse("2006-01-02T15:04:05Z", data["lastmodified"].(string))
	if err != nil {
		return fmt.Errorf("Could not parse last modified date: %s", err)
	}
	store.metadata.LastModified = lastModified
	store.metadata.UnencryptedSuffix = data["unencrypted_suffix"].(string)
	store.metadata.Version = data["version"].(string)
	if k, ok := data["kms"].([]interface{}); ok {
		ks, err := store.kmsEntries(k)
		if err == nil {
			store.metadata.KeySources = append(store.metadata.KeySources, ks)
		}

	}

	if pgp, ok := data["pgp"].([]interface{}); ok {
		ks, err := store.pgpEntries(pgp)
		if err == nil {
			store.metadata.KeySources = append(store.metadata.KeySources, ks)
		}
	}
	return nil
}

func (store *YAMLStore) Metadata() sops.Metadata {
	return store.metadata
}

func (store *YAMLStore) kmsEntries(in []interface{}) (sops.KeySource, error) {
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
		creationDate, err := time.Parse("2006-01-02T15:04:05Z", entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}

func (store *YAMLStore) pgpEntries(in []interface{}) (sops.KeySource, error) {
	var keys []sops.MasterKey
	keysource := sops.KeySource{Name: "pgp", Keys: keys}
	for _, v := range in {
		entry := v.(map[interface{}]interface{})
		key := &pgp.GPGMasterKey{}
		key.Fingerprint = entry["fp"].(string)
		key.EncryptedKey = entry["enc"].(string)
		creationDate, err := time.Parse("2006-01-02T15:04:05Z", entry["created_at"].(string))
		if err != nil {
			return keysource, fmt.Errorf("Could not parse creation date: %s", err)
		}
		key.CreationDate = creationDate
		keysource.Keys = append(keysource.Keys, key)
	}
	return keysource, nil
}
