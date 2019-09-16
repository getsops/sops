package toml //import "go.mozilla.org/sops/stores/toml"

import (
	"github.com/BurntSushi/toml"
	"go.mozilla.org/sops"
)

// TOML reader that preserves original order via toml.Metadata.Keys(),
// see https://godoc.org/github.com/BurntSushi/toml#MetaData.Keys
// and converts the unoredered data (multi-level tree of
// map[string]interface{}) into sops.BranchTree.
type tomlReader []toml.Key

func (c tomlReader) keysInOrder() []string {
	var keys []string
	visited := map[string]struct{}{}
	for _, key := range c {
		_, ok := visited[key[0]]
		if !ok {
			keys = append(keys, key[0])
			visited[key[0]] = struct{}{}
		}
	}
	return keys
}

func (c tomlReader) table(key string) tomlReader {
	var keys []toml.Key
	for _, k := range c {
		if len(k) > 1 && k[0] == key {
			keys = append(keys, k[1:])
		}
	}
	return tomlReader(keys)
}

func (c tomlReader) arrayItem(key string, at int) tomlReader {
	var keys []toml.Key
	arrayAt := -1
	for _, k := range c {
		if len(k) == 1 && k[0] == key {
			arrayAt++
			continue
		}
		if len(k) > 1 && k[0] == key && at == arrayAt {
			keys = append(keys, k[1:])
		}
	}
	return tomlReader(keys)
}

func (c tomlReader) readToTreeBranch(unordered map[string]interface{}) sops.TreeBranch {
	var branch sops.TreeBranch
	for _, key := range c.keysInOrder() {
		value := unordered[key]
		switch v := value.(type) {
		case map[string]interface{}:
			branch = append(branch, sops.TreeItem{
				Key:   key,
				Value: c.table(key).readToTreeBranch(v),
			})
		case []map[string]interface{}:
			var array []interface{}
			for i, item := range v {
				array = append(array, c.arrayItem(key, i).readToTreeBranch(item))
			}
			branch = append(branch, sops.TreeItem{
				Key:   key,
				Value: array,
			})
		default:
			branch = append(branch, sops.TreeItem{
				Key:   key,
				Value: value,
			})
		}
	}
	return branch
}
