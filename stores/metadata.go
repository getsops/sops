package stores

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/go-viper/mapstructure/v2"

	"github.com/getsops/sops/v3"
)

// MetadataFlatten is an enum type
type MetadataFlatten int

const (
	MetadataFlattenNone MetadataFlatten = iota
	MetadataFlattenBelowTop
	MetadataFlattenFull
)

type MetadataOpts struct {
	Flatten MetadataFlatten
}

// SopsPrefix is the prefix for all metadatada entry keys
const SopsPrefix = SopsMetadataKey + "_"

func sopsToGoMap(mapping sops.TreeBranch) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, item := range mapping {
		if _, ok := item.Key.(sops.Comment); ok {
			continue
		}
		key, ok := item.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Unexpected key type %T", item.Key)
		}
		value, err := sopsToGo(item.Value)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func sopsToGoSlice(slice []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(slice))
	for idx, item := range slice {
		if _, ok := item.(sops.Comment); ok {
			continue
		}
		value, err := sopsToGo(item)
		if err != nil {
			return nil, err
		}
		result[idx] = value
	}
	return result, nil
}

func sopsToGo(value interface{}) (interface{}, error) {
	switch value := value.(type) {
	case sops.TreeBranch:
		return sopsToGoMap(value)
	case []interface{}:
		return sopsToGoSlice(value)
	default:
		return value, nil
	}
}

func treeBranchToMetadata(meta sops.TreeBranch) (metadata, error) {
	var md metadata
	m, err := sopsToGoMap(meta)
	if err != nil {
		return md, err
	}
	config := mapstructure.DecoderConfig{
		Result:           &md,
		WeaklyTypedInput: true,
	}
	d, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return md, err
	}
	err = d.Decode(m)
	return md, err
}

// Extract SOPS metadata from tree branches.
func ExtractMetadata(branches sops.TreeBranches, opts MetadataOpts) (sops.TreeBranches, sops.Metadata, error) {
	var metadataTree sops.TreeBranch
	if opts.Flatten != MetadataFlattenFull {
		first := true
		for bi, branch := range branches {
			i := 0
			for i < len(branch) {
				if branch[i].Key == SopsMetadataKey {
					if bi == 0 {
						if !first {
							return nil, sops.Metadata{}, fmt.Errorf("Found duplicate %v entry", SopsMetadataKey)
						}
						first = false
						if tree, ok := branch[i].Value.(sops.TreeBranch); ok {
							metadataTree = tree
						} else {
							return nil, sops.Metadata{}, fmt.Errorf("Found %v entry that is not a mapping", SopsMetadataKey)
						}
					}
					branch = append(branch[:i], branch[i+1:]...)
				} else {
					i += 1
				}
			}
			branches[bi] = branch
		}
	} else {
		if len(branches) >= 1 {
			branch := branches[0]
			for i := 0; i < len(branch); i += 1 {
				if key, ok := branch[i].Key.(string); ok {
					if strings.HasPrefix(key, SopsPrefix) {
						entry := branch[i]
						entry.Key = key[len(SopsPrefix):]
						metadataTree = append(metadataTree, entry)
						branch = append(branch[:i], branch[i+1:]...)
						i -= 1
					}
				}
			}
			branches[0] = branch
		}
	}
	if metadataTree == nil {
		return nil, sops.Metadata{}, sops.MetadataNotFound
	}
	if opts.Flatten != MetadataFlattenNone {
		var err error
		metadataTree, err = unflattenTreeBranch(metadataTree)
		if err != nil {
			return nil, sops.Metadata{}, err
		}
	}
	md, err := treeBranchToMetadata(metadataTree)
	if err != nil {
		return nil, sops.Metadata{}, err
	}
	metadata, err := md.ToInternal()
	if err != nil {
		return nil, sops.Metadata{}, err
	}
	return branches, metadata, nil
}

type mapKey struct {
	Name string
	Key  reflect.Value
}

// byName implements sort.Interface for []mapKey
type byName []mapKey

func (mapKeys byName) Len() int {
	return len(mapKeys)
}

func (mapKeys byName) Swap(i, j int) {
	mapKeys[i], mapKeys[j] = mapKeys[j], mapKeys[i]
}

func (mapKeys byName) Less(i, j int) bool {
	return mapKeys[i].Name < mapKeys[j].Name
}

func goToSops(value interface{}) (interface{}, error) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		result := make([]interface{}, val.Len())
		for j := 0; j < val.Len(); j++ {
			v, err := goToSops(val.Index(j).Interface())
			if err != nil {
				return nil, err
			}
			result[j] = v
		}
		return result, nil
	case reflect.Map:
		keys := val.MapKeys()
		sortedKeys := make([]mapKey, len(keys))
		for idx, key := range keys {
			sortedKeys[idx] = mapKey{
				Name: key.Interface().(string),
				Key:  key,
			}
		}
		sort.Sort(byName(sortedKeys))
		result := make(sops.TreeBranch, len(sortedKeys))
		for idx, key := range sortedKeys {
			v, err := goToSops(val.MapIndex(key.Key).Interface())
			if err != nil {
				return nil, err
			}
			result[idx] = sops.TreeItem{
				Key:   key.Name,
				Value: v,
			}
		}
		return result, nil
	default:
		return value, nil
	}
}

func metadataToTreeBranch(md metadata) (sops.TreeBranch, error) {
	var mdMap map[string]interface{}
	config := mapstructure.DecoderConfig{
		Result: &mdMap,
	}
	d, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return nil, err
	}
	err = d.Decode(md)
	if err != nil {
		return nil, err
	}
	metadata, err := goToSops(mdMap)
	if err != nil {
		return nil, err
	}
	if tb, ok := metadata.(sops.TreeBranch); ok {
		return tb, nil
	}
	return nil, fmt.Errorf("Internal error: unexpected metadata conversion result %T", metadata)
}

func SerializeMetadata(data sops.Tree, opts MetadataOpts) (sops.TreeBranches, error) {
	md, err := metadataToTreeBranch(metadataFromInternal(data.Metadata))
	if err != nil {
		return nil, fmt.Errorf("Error while serializing metadata: %e", err)
	}
	if opts.Flatten != MetadataFlattenNone {
		var prefix string
		if opts.Flatten == MetadataFlattenFull {
			prefix = SopsPrefix
		}
		md, err = flattenTreeBranch(md, prefix)
		if err != nil {
			return nil, fmt.Errorf("Error while flatting metadata: %e", err)
		}
	}
	if opts.Flatten != MetadataFlattenFull {
		md = sops.TreeBranch{
			sops.TreeItem{
				Key:   SopsMetadataKey,
				Value: md,
			},
		}
	}
	var result sops.TreeBranches
	for _, branch := range data.Branches {
		newBranch := make(sops.TreeBranch, 0, len(branch)+len(md))
		for _, item := range branch {
			if key, ok := item.Key.(string); ok {
				if opts.Flatten == MetadataFlattenFull {
					if strings.HasPrefix(key, SopsPrefix) {
						return nil, fmt.Errorf("Found key %q in encrypted data, which starts with the reserved key prefix %q for SOPS metadata", key, SopsPrefix)
					}
				} else {
					if key == SopsMetadataKey {
						return nil, fmt.Errorf("Found key %q in encrypted data, which is a reserved key used for SOPS metadata", key)
					}
				}
			}
			newBranch = append(newBranch, item)
		}
		for _, item := range md {
			newBranch = append(newBranch, item)
		}
		result = append(result, newBranch)
	}
	return result, nil
}
