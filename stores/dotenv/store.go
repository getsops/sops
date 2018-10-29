package dotenv //import "go.mozilla.org/sops/stores/dotenv"

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/stores"
	"strings"
)

const SopsPrefix = "sops_"

// Store handles storage of dotenv data
type Store struct {
}

func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branch, err := store.LoadPlainFile(in)
	if err != nil { return sops.Tree{}, err }

	resultBranch := make(sops.TreeBranch, 0)
	metadata := stores.Metadata{}
	md_found := false
	for _, item := range branch {
		// FIXME: use sops_* items instead
		if strings.HasPrefix(item.Key.(string), SopsPrefix) {
			if item.Key == SopsPrefix + "metadata" {
				metadata, err = FromGOB64(fmt.Sprint(item.Value))
				if err != nil { return sops.Tree{}, err }
				md_found = true
				break
			}
		} else {
			resultBranch = append(resultBranch, item)
		}
	}
	if !md_found { return sops.Tree{}, sops.MetadataNotFound }

	internalMetadata, err := metadata.ToInternal()
	if err != nil { return sops.Tree{}, err }

	return sops.Tree{
		Branch:   resultBranch,
		Metadata: internalMetadata,
	}, nil
}

func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranch, error) {
	branch := make(sops.TreeBranch, 0)
	reader := bytes.NewReader(in)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { continue }
		pos := strings.Index(line, "=")
		if pos == -1 {
			return nil, fmt.Errorf("invalid dotenv input line: %s", line)
		}
		branch = append(branch, sops.TreeItem{
			Key:   line[:pos],
			Value: line[pos+1:],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("invalid dotenv input: %s", err)
	}
	return branch, nil
}

func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	metadata := stores.MetadataFromInternal(in.Metadata)
	mdItems, err := metadataToTreeItems(metadata)
	if err != nil { return nil, err }
	for key, value := range mdItems {
		if value == "" { continue }
		in.Branch = append(in.Branch, sops.TreeItem{Key: SopsPrefix + key, Value: value})
	}
	return store.EmitPlainFile(in.Branch)
}
func (store *Store) EmitPlainFile(in sops.TreeBranch) ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, item := range in {
		if isComplexValue(item.Value) {
			return nil, fmt.Errorf( "cannot use complex value in dotenv file: %s", item.Value)
		}
		line := fmt.Sprintf("%s=%s\n", item.Key, item.Value)
		buffer.WriteString(line)
	}
	return buffer.Bytes(), nil
}

func (Store) EmitValue(v interface{}) ([]byte, error) {
	// FIXME: What should this function do?
	panic("implement me")
}

func metadataToTreeItems(md stores.Metadata) (map[string] string, error) {
	// FIXME: encode all metadata in sops_* items
	mdGob, err := ToGOB64(md)
	if err != nil { return nil, err }
	return map[string] string {
		"version":                     md.Version,
		"last_modified":               md.LastModified,
		"unencrypted_suffix":          md.UnencryptedSuffix,
		"encrypted_suffix":            md.EncryptedSuffix,
		"message_authentication_code": md.MessageAuthenticationCode,
		"metadata":                    mdGob,
	}, nil
}

func isComplexValue(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	case sops.TreeBranch:
		return true
	}
	return false
}

func init() {
	gob.Register(stores.Metadata{})
}

func ToGOB64(m stores.Metadata) (string, error) {
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(m)
	if err != nil {
		return "", fmt.Errorf("could not base64-encode metadata: %s", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func FromGOB64(str string) (stores.Metadata, error) {
	metadata := stores.Metadata{}
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return metadata, fmt.Errorf("could not base64-decode metadata: %s", err)
	}
	buffer := bytes.Buffer{}
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&metadata)
	if err != nil {
		return stores.Metadata{}, fmt.Errorf("could not parse metadata: %s", err)
	}
	return metadata, nil
}
