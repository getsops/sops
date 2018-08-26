package ini //import "go.mozilla.org/sops/stores/ini"

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/stores"
	"strconv"
	"reflect"
	"gopkg.in/ini.v1"
	"strings"
	"sort"
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
		item := sops.TreeItem{Key:key.Name(), Value:key.Value()}
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
			if sectionBranch.Key == "sops"{
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

	metadata := stores.Metadata{}
	m := reflect.ValueOf(&metadata).Elem()

	for _, key := range sopsSection.Keys() {

		if strings.Contains(key.Name(), ".") {
			parts := strings.SplitN(key.Name(), ".", 2)
			if len(parts) != 2 {
				return metadata, fmt.Errorf("Bad metadata format: key %s makes no sense", key.Name())
			}
			prefix := parts[0]
			remainder := parts[1]
			// Is a slice
			if strings.Contains(prefix, "[") {
				k, i, err := parseSliceKey(prefix)
				if err != nil {
					return metadata, fmt.Errorf("Bad metadata format: %s", err)
				}

				f := m.FieldByName(k)
				ensureReflectedSliceLength(f, i+1)
				sliceItem := f.Index(i)
				setMetadataField(sliceItem, remainder, key)
			} else {
				return metadata, fmt.Errorf("Bad metadata format: expected array but have %s", prefix)
			}
		} else {
			err := setMetadataField(m, key.Name(), key)
			if err != nil {
				return metadata, err
			}
		}
	}

	return metadata, nil
}

func ensureReflectedSliceLength(slice reflect.Value, length int) {
	if slice.Len() < length  {
		expanded := reflect.MakeSlice(slice.Type(), slice.Len()+1, slice.Cap()+1)
		reflect.Copy(expanded, slice)
		slice.Set(expanded)
	}
}

func setMetadataField(m reflect.Value, name string, iniKey *ini.Key) error {
	f := m.FieldByName(name)
	switch f.Kind() {
	case reflect.String:
		f.SetString(strings.Replace(iniKey.String(), "\\n", "\n", -1))
	case reflect.Int:
		val, err := iniKey.Int64()
		if err != nil {
			return err
		}
		f.SetInt(val)
	}
	return nil
}

func parseSliceKey(key string) (string, int, error) {
	openBracket := strings.IndexRune(key, '[')
	closeBracket := strings.IndexRune(key, ']')
	name := key[:openBracket]
	indexStr := key[openBracket+1:closeBracket]
	i, err := strconv.Atoi(indexStr)
	return name, i, err
}

func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	branches, err := store.treeBranchesFromIni(in)
	if err != nil {
		return branches, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	return branches, nil
}

func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {

	metadata := stores.MetadataFromInternal(in.Metadata)
	newBranch, err := store.encodeMetadataToIniBranch(metadata, "sops")
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

func (store *Store) encodeMetadataToIniBranch(metadata interface{}, prefix string) (sops.TreeBranch, error) {

	branch := sops.TreeBranch{}

	m := reflect.ValueOf(metadata)
	r, err := encodeMetadataItem("", m.Type().Kind(), m)

	// Keys are sorted so sops section is stable (for nice diffs)
	var keys []string
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		branch = append(branch, sops.TreeItem{Key: k, Value: r[k]})
	}

	return branch, err
}

func encodeMetadataItem(prefix string, kind reflect.Kind, field reflect.Value) (map[string]interface{}, error) {

	result := make(map[string]interface{}, 0)

	switch kind {
	case reflect.Slice:
		slf := field
		for j := 0; j < slf.Len(); j++ {
			item := slf.Index(j)
			p := fmt.Sprintf("%s[%d]", prefix, j)
			r, err := encodeMetadataItem(p, item.Type().Kind(), item)
			if err != nil {
				return result, err
			}
			for k, v := range r {
				result[k] = v
			}
		}
	case reflect.Struct:
		for i := 0; i < field.NumField(); i++ {
			sf := field.Type().Field(i)
			var name string
			if prefix == "" {
				name = sf.Name
			} else {
				name = fmt.Sprintf("%s.%s", prefix, sf.Name)
			}
			r, err := encodeMetadataItem(name, sf.Type.Kind(), field.Field(i))
			if err != nil {
				return result, err
			}
			for k, v := range r {
				result[k] = v
			}
		}
	case reflect.Int:
		if field.Int() != 0 {
			result[prefix] = string(field.Int())
		}
	case reflect.String:
		if field.String() != "" {
			result[prefix] = strings.Replace(field.String(), "\n", "\\n", -1)
		}
	default:
		return result, fmt.Errorf("Cannot encode %s, unexpected type %s", prefix, kind)
	}

	return result, nil
}


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

func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	return store.encodeValue(v)
}
