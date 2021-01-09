package ini //import "go.mozilla.org/sops/v3/stores/ini"

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strconv"
	"strings"

	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/stores"
	"gopkg.in/ini.v1"
)

// Store handles storage of ini data.
type Store struct {
}

func (store Store) encodeTree(branches sops.TreeBranches) ([]byte, error) {
	iniFile := ini.Empty()
	for _, branch := range branches {
		for _, item := range branch {
			if _, ok := item.Key.(sops.Comment); ok {
				continue
			}
			section, err := iniFile.NewSection(item.Key.(string))
			if err != nil {
				return nil, fmt.Errorf("Error encoding section %s: %s", item.Key, err)
			}
			itemTree, ok := item.Value.(sops.TreeBranch)
			if !ok {
				return nil, fmt.Errorf("Error encoding section: Section values should always be TreeBranches")
			}

			first := 0
			if len(itemTree) > 0 {
				if sectionComment, ok := itemTree[0].Key.(sops.Comment); ok {
					section.Comment = sectionComment.Value
					first = 1
				}
			}

			var lastItem *ini.Key
			for i := first; i < len(itemTree); i++ {
				keyVal := itemTree[i]
				if comment, ok := keyVal.Key.(sops.Comment); ok {
					if lastItem != nil {
						lastItem.Comment = comment.Value
					}
				} else {
					lastItem, err = section.NewKey(keyVal.Key.(string), store.valToString(keyVal.Value))
					if err != nil {
						return nil, fmt.Errorf("Error encoding key: %s", err)
					}
				}
			}
		}
	}
	var buffer bytes.Buffer
	iniFile.WriteTo(&buffer)
	return buffer.Bytes(), nil
}

func (store Store) stripCommentChar(comment string) string {
	if strings.HasPrefix(comment, ";") {
		comment = strings.TrimLeft(comment, "; ")
	} else if strings.HasPrefix(comment, "#") {
		comment = strings.TrimLeft(comment, "# ")
	}
	return comment
}

func (store Store) valToString(v interface{}) string {
	switch v := v.(type) {
	case fmt.Stringer:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%s", v)
	}
}

func (store Store) iniFromTreeBranches(branches sops.TreeBranches) ([]byte, error) {
	return store.encodeTree(branches)
}

func (store Store) treeBranchesFromIni(in []byte) (sops.TreeBranches, error) {
	iniFile, err := ini.Load(in)
	if err != nil {
		return nil, err
	}
	var branch sops.TreeBranch
	for _, section := range iniFile.Sections() {

		item, err := store.treeItemFromSection(section)
		if err != nil {
			return sops.TreeBranches{branch}, err
		}
		branch = append(branch, item)
	}
	return sops.TreeBranches{branch}, nil
}

func (store Store) treeItemFromSection(section *ini.Section) (sops.TreeItem, error) {
	var sectionItem sops.TreeItem
	sectionItem.Key = section.Name()
	var items sops.TreeBranch

	if section.Comment != "" {
		items = append(items, sops.TreeItem{
			Key: sops.Comment{
				Value: store.stripCommentChar(section.Comment),
			},
			Value: nil,
		})
	}

	for _, key := range section.Keys() {
		item := sops.TreeItem{Key: key.Name(), Value: key.Value()}
		items = append(items, item)
		if key.Comment != "" {
			items = append(items, sops.TreeItem{
				Key: sops.Comment{
					Value: store.stripCommentChar(key.Comment),
				},
				Value: nil,
			})
		}
	}
	sectionItem.Value = items
	return sectionItem, nil
}

// LoadEncryptedFile loads encrypted INI file's bytes onto a sops.Tree runtime object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	iniFileOuter, err := ini.Load(in)
	if err != nil {
		return sops.Tree{}, err
	}

	sopsSection, err := iniFileOuter.GetSection("sops")
	if err != nil {
		return sops.Tree{}, sops.MetadataNotFound
	}

	metadataHolder, err := store.iniSectionToMetadata(sopsSection)
	if err != nil {
		return sops.Tree{}, err
	}

	metadata, err := metadataHolder.ToInternal()
	if err != nil {
		return sops.Tree{}, err
	}
	// After that, we load the whole file into a map.
	branches, err := store.treeBranchesFromIni(in)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	// Discard metadata, as we already loaded it.
	for bi, branch := range branches {
		for s, sectionBranch := range branch {
			if sectionBranch.Key == "sops" {
				branch = append(branch[:s], branch[s+1:]...)
				branches[bi] = branch
			}
		}
	}
	return sops.Tree{
		Branches: branches,
		Metadata: metadata,
	}, nil
}

func (store *Store) iniSectionToMetadata(sopsSection *ini.Section) (stores.Metadata, error) {

	metadataHash := make(map[string]interface{})
	for k, v := range sopsSection.KeysHash() {
		metadataHash[k] = strings.Replace(v, "\\n", "\n", -1)
	}
	m := stores.Unflatten(metadataHash)
	var md stores.Metadata
	inrec, err := json.Marshal(m)
	if err != nil {
		return md, err
	}
	err = json.Unmarshal(inrec, &md)
	return md, err
}

// LoadPlainFile loads a plaintext INI file's bytes onto a sops.TreeBranches runtime object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	branches, err := store.treeBranchesFromIni(in)
	if err != nil {
		return branches, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	return branches, nil
}

// EmitEncryptedFile returns encrypted INI file bytes corresponding to a sops.Tree
// runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {

	metadata := stores.MetadataFromInternal(in.Metadata)
	newBranch, err := store.encodeMetadataToIniBranch(metadata)
	if err != nil {
		return nil, err
	}
	sectionItem := sops.TreeItem{Key: "sops", Value: newBranch}
	branch := sops.TreeBranch{sectionItem}

	in.Branches = append(in.Branches, branch)

	out, err := store.iniFromTreeBranches(in.Branches)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to ini: %s", err)
	}
	return out, nil
}

func (store *Store) encodeMetadataToIniBranch(md stores.Metadata) (sops.TreeBranch, error) {
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
	if err != nil {
		return nil, err
	}
	branch := sops.TreeBranch{}
	for key, value := range flat {
		if value == nil {
			continue
		}
		branch = append(branch, sops.TreeItem{Key: key, Value: value})
	}
	return branch, nil
}

// EmitPlainFile returns the plaintext INI file bytes corresponding to a sops.TreeBranches object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	out, err := store.iniFromTreeBranches(in)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to ini: %s", err)
	}
	return out, nil
}

func (store Store) encodeValue(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case sops.TreeBranches:
		return store.encodeTree(v)
	default:
		return json.Marshal(v)
	}
}

// EmitValue returns a single value encoded in a generic interface{} as bytes
func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	return store.encodeValue(v)
}

// EmitExample returns the plaintext INI file bytes corresponding to the SimpleTree example
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleSimpleTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}
