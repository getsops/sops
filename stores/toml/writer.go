package toml //import "go.mozilla.org/sops/stores/toml"

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"go.mozilla.org/sops"
)

func printTreeBranchInTOML(w io.Writer, branch sops.TreeBranch) error {
	return printTreeBranch(newIndenter(w, ""), branch)
}

func printTreeItemInTOML(w io.Writer, item sops.TreeItem) error {
	return printTreeItem(newIndenter(w, ""), item)
}

// TOML indenter.
type indenter struct {
	w           io.Writer
	path        string
	indentation []byte
	written     bool
}

func newIndenter(w io.Writer, key string) *indenter {
	var indentation []byte
	if indenter, ok := w.(*indenter); ok {
		// Has parent indenter?
		if len(indenter.path) > 0 {
			key = indenter.path + "." + key
		}
		if indenter.written {
			w.Write([]byte("\n"))
		}
		indentation = []byte{' ', ' '}
	}
	return &indenter{w: w, path: key, indentation: indentation}
}

func (w *indenter) PrintTableName() {
	fmt.Fprintf(w.w, "[%v]\n", w.path)
}

func (w *indenter) PrintArrayName() {
	fmt.Fprintf(w.w, "[[%v]]\n", w.path)
}

func (w *indenter) Write(p []byte) (n int, err error) {
	w.written = true
	return w.w.Write(append(w.indentation, p...))
}

func printTreeBranch(w io.Writer, branch sops.TreeBranch) error {
	for _, item := range branch {
		printTreeItem(w, item)
	}
	return nil
}

func printTreeItem(w io.Writer, item sops.TreeItem) error {
	if _, ok := item.Key.(sops.Comment); ok {
		return nil
	}
	key, ok := toString(item.Key)
	if !ok {
		return fmt.Errorf("Error encoding item.Key of type=%T, value=%v", item.Key, item.Key)
	}
	switch v := item.Value.(type) {
	case sops.TreeBranch:
		indenter := newIndenter(w, key)
		indenter.PrintTableName()
		if err := printTreeBranch(indenter, v); err != nil {
			return err
		}
	case []interface{}:
		// Look up type of the values without reflection.
		var isTreeBranch bool
		for _, item := range v {
			_, isTreeBranch = item.(sops.TreeBranch)
			break
		}
		if isTreeBranch {
			for _, item := range v {
				indenter := newIndenter(w, key)
				indenter.PrintArrayName()
				branch, ok := item.(sops.TreeBranch)
				if !ok {
					return fmt.Errorf("Error encoding array: unexpected item of type %T", item)
				}
				if err := printTreeBranch(indenter, branch); err != nil {
					return err
				}
			}
		} else {
			enc := toml.NewEncoder(w)
			if err := enc.Encode(map[string]interface{}{key: v}); err != nil {
				return fmt.Errorf("Error encoding %v: %s", v, err)
			}
		}
	default:
		enc := toml.NewEncoder(w)
		if err := enc.Encode(map[string]interface{}{key: v}); err != nil {
			return fmt.Errorf("Error encoding %v: %s", item.Key, err)
		}
	}
	return nil
}

func toString(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case toml.Key:
		return v[0], true
	default:
		return "", false
	}
}
