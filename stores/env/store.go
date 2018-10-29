package env //import "go.mozilla.org/sops/stores/env"
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

// Store handles storage of env data
type Store struct {
}

func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branch, err := store.LoadPlainFile(in)
	if err != nil { return sops.Tree{}, err }

	storeMetadata := stores.Metadata{}
	index := -1
	for i, item := range branch {
		if item.Key == "_metadata" {
			storeMetadata, err = FromGOB64(fmt.Sprint(item.Value))
			if err != nil { return sops.Tree{}, err }
			index = i
			break
		}
	}
	if index == -1 { return sops.Tree{}, sops.MetadataNotFound }

	internalMetadata, err := storeMetadata.ToInternal()
	if err != nil { return sops.Tree{}, err }

	return sops.Tree{
		Branch:   append(branch[:index], branch[index+1:]...),
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
			// FIXME: Print line number
			return nil, fmt.Errorf("could not parse line: %s", line)
		}
		// FIXME: Trim key and value? Remove quotation marks?
		branch = append(branch, sops.TreeItem{
			Key:   line[:pos],
			Value: line[pos+1:],
		})
	}
	return branch, nil
}

func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	plain, err := store.EmitPlainFile(in.Branch)
	if err != nil { return nil, err }
	buffer := bytes.NewBuffer(plain)
	metadata := stores.MetadataFromInternal(in.Metadata)
	str, err := ToGOB64(metadata)
	if err != nil { return nil, err }
	line := fmt.Sprintf("_metadata=%s\n", str)
	buffer.WriteString(line)
	return buffer.Bytes(), nil
}

func (store *Store) EmitPlainFile(in sops.TreeBranch) ([]byte, error) {
	buffer := bytes.Buffer{}
	for _, item := range in {
		// FIXME: Check that item.Value is a scalar.
		// FIXME: Does Go know how to print the OS-specific EOL string? Do we care?
		line := fmt.Sprintf("%s=%s\n", item.Key, item.Value)
		buffer.WriteString(line)
	}
	return buffer.Bytes(), nil
}

func (Store) EmitValue(v interface{}) ([]byte, error) {
	// FIXME: Whot should this function do?
	panic("implement me")
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
