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
	branch, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}

	var resultBranch sops.TreeBranch
	mdMap := make(map[string]interface{})
	for _, item := range branch {
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
		Branch:   resultBranch,
		Metadata: internalMetadata,
	}, nil
}

func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranch, error) {
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
	return branch, nil
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
		in.Branch = append(in.Branch, sops.TreeItem{Key: SopsPrefix + key, Value: value})
	}
	return store.EmitPlainFile(in.Branch)
}

func (store *Store) EmitPlainFile(in sops.TreeBranch) ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, item := range in {
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
	return flatten(mdMap), nil
}

func mapToMetadata(m map[string]interface{}) (stores.Metadata, error) {
	m = unflatten(m)
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

func flatten(m map[string]interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	flattenRecursive(m, []string{}, func(ks []string, v interface{}) {
		if s, ok := v.(string); ok {
			v = strings.Replace(s, "\n", "\\n", -1)
			v = fmt.Sprintf(`"%s"`, v)
		}
		r[strings.Join(ks, "_")] = v
	})
	return r
}

func flattenRecursive(v interface{}, ks []string, cb func([]string, interface{})) {
	if m, ok := v.(map[string]interface{}); ok {
		for k, v := range m {
			newks := append(ks, k)
			flattenRecursive(v, newks, cb)
		}
	} else if s, ok := v.([]interface{}); ok {
		for i, e := range s {
			newks := append(ks, fmt.Sprint(i))
			flattenRecursive(e, newks, cb)
		}
	} else {
		cb(ks, v)
	}
}

type TokenizerFunc func(string) []string

func unflatten(m map[string]interface{}) map[string]interface{} {
	var tree = make(map[string]interface{})
	for k, v := range m {
		ks := strings.Split(k, "_")
		tr := tree
		for _, tk := range ks[:len(ks)-1] {
			trnew, ok := tr[tk]
			if !ok {
				trnew = make(map[string]interface{})
				tr[tk] = trnew
			}
			tr = trnew.(map[string]interface{})
		}
		if s, ok := v.(string); ok {
			if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
				s = s[1 : len(s)-1]
				v = strings.Replace(s, "\\n", "\n", -1)
			}
		}
		tr[ks[len(ks)-1]] = v
	}
	return tree
}
