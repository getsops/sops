package dotenv //import "go.mozilla.org/sops/stores/dotenv"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/stores"
)

const SopsPrefix = "sops_"

// Store handles storage of dotenv data
type Store struct {
}

func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}

	var resultBranch sops.TreeBranch
	mdMap := make(map[string]interface{})
	for _, item := range branches[0] {
		s := item.Key.(string)
		if strings.HasPrefix(s, SopsPrefix) {
			s = s[len(SopsPrefix):]
			mdMap[s] = item.Value
		} else {
			resultBranch = append(resultBranch, item)
		}
	}

	metadata, err := mapToMetadata(mdMap)
	if err != nil {
		return sops.Tree{}, err
	}
	internalMetadata, err := metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, err
	}

	return sops.Tree{
		Branches: sops.TreeBranches{
			resultBranch,
		},
		Metadata: internalMetadata,
	}, nil
}

func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	var branches sops.TreeBranches
	var branch sops.TreeBranch

	for _, line := range bytes.Split(in, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		pos := bytes.Index(line, []byte("="))
		if pos == -1 {
			return nil, fmt.Errorf("invalid dotenv input line: %s", line)
		}
		branch = append(branch, sops.TreeItem{
			Key:   string(line[:pos]),
			Value: string(line[pos+1:]),
		})
	}

	branches = append(branches, branch)
	return branches, nil
}

func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	metadata := stores.MetadataFromInternal(in.Metadata)
	mdItems, err := metadataToMap(metadata)
	if err != nil {
		return nil, err
	}
	for key, value := range mdItems {
		if value == nil {
			continue
		}
		in.Branches[0] = append(in.Branches[0], sops.TreeItem{Key: SopsPrefix + key, Value: value})
	}
	return store.EmitPlainFile(in.Branches)
}

func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, item := range in[0] {
		if isComplexValue(item.Value) {
			return nil, fmt.Errorf("cannot use complex value in dotenv file: %s", item.Value)
		}
		line := fmt.Sprintf("%s=%s\n", item.Key, item.Value)
		buffer.WriteString(line)
	}
	return buffer.Bytes(), nil
}

func (Store) EmitValue(v interface{}) ([]byte, error) {
	if s, ok := v.(string); ok {
		return []byte(s), nil
	}
	return nil, fmt.Errorf("the dotenv store only supports emitting strings, got %T", v)
}

func metadataToMap(md stores.Metadata) (map[string]interface{}, error) {
	var mdMap map[string]interface{}
	inrec, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &mdMap)
	if err != nil {
		return nil, err
	}
	flat := stores.Flatten(mdMap)
	for k, v := range flat {
		if s, ok := v.(string); ok {
			flat[k] = strings.Replace(s, "\n", "\\n", -1)
		}
	}
	return flat, nil
}

func mapToMetadata(m map[string]interface{}) (stores.Metadata, error) {
	for k, v := range m {
		if s, ok := v.(string); ok {
			m[k] = strings.Replace(s, "\\n", "\n", -1)
		}
	}
	m = stores.Unflatten(m)
	var md stores.Metadata
	inrec, err := json.Marshal(m)
	if err != nil {
		return md, err
	}
	err = json.Unmarshal(inrec, &md)
	return md, err
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
