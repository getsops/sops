package yaml //import "go.mozilla.org/sops/v3/stores/yaml"

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/stores"
)

// Store handles storage of YAML data
type Store struct {
}

func (store Store) appendCommentToList(comment string, list []interface{}) []interface{} {
	if comment != "" {
		for _, commentLine := range strings.Split(comment, "\n") {
			if commentLine != "" {
				list = append(list, sops.Comment{
					Value: commentLine[1:],
				})
			}
		}
	}
	return list
}

func (store Store) appendCommentToMap(comment string, branch sops.TreeBranch) sops.TreeBranch {
	if comment != "" {
		for _, commentLine := range strings.Split(comment, "\n") {
			if commentLine != "" {
				branch = append(branch, sops.TreeItem{
					Key: sops.Comment{
						Value: commentLine[1:],
					},
					Value: nil,
				})
			}
		}
	}
	return branch
}

func (store Store) nodeToTreeValue(node *yaml.Node, commentsWereHandled bool) (interface{}, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		panic("documents should never be passed here")
	case yaml.SequenceNode:
		var result []interface{}
		if !commentsWereHandled {
			result = store.appendCommentToList(node.HeadComment, result)
			result = store.appendCommentToList(node.LineComment, result)
		}
		for _, item := range node.Content {
			result = store.appendCommentToList(item.HeadComment, result)
			result = store.appendCommentToList(item.LineComment, result)
			val, err := store.nodeToTreeValue(item, true)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
			result = store.appendCommentToList(item.FootComment, result)
		}
		if !commentsWereHandled {
			result = store.appendCommentToList(node.FootComment, result)
		}
		return result, nil
	case yaml.MappingNode:
		branch := make(sops.TreeBranch, 0)
		return store.appendYamlNodeToTreeBranch(node, branch, false)
	case yaml.ScalarNode:
		var result interface{}
		node.Decode(&result)
		return result, nil
	case yaml.AliasNode:
		return store.nodeToTreeValue(node.Alias, false);
	}
	return nil, nil
}

func (store Store) appendYamlNodeToTreeBranch(node *yaml.Node, branch sops.TreeBranch, commentsWereHandled bool) (sops.TreeBranch, error) {
	var err error
	if !commentsWereHandled {
		branch = store.appendCommentToMap(node.HeadComment, branch)
		branch = store.appendCommentToMap(node.LineComment, branch)
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, item := range node.Content {
			branch, err = store.appendYamlNodeToTreeBranch(item, branch, false)
			if err != nil {
				return nil, err
			}
		}
	case yaml.SequenceNode:
		return nil, fmt.Errorf("YAML documents that are sequences are not supported")
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i + 1]
			branch = store.appendCommentToMap(key.HeadComment, branch)
			branch = store.appendCommentToMap(key.LineComment, branch)
			handleValueComments := value.Kind == yaml.ScalarNode || value.Kind == yaml.AliasNode
			if handleValueComments {
				branch = store.appendCommentToMap(value.HeadComment, branch)
				branch = store.appendCommentToMap(value.LineComment, branch)
			}
			var keyValue interface{}
			key.Decode(&keyValue)
			valueTV, err := store.nodeToTreeValue(value, handleValueComments)
			if err != nil {
				return nil, err
			}
			branch = append(branch, sops.TreeItem{
				Key:   keyValue,
				Value: valueTV,
			})
			if handleValueComments {
				branch = store.appendCommentToMap(value.FootComment, branch)
			}
			branch = store.appendCommentToMap(key.FootComment, branch)
		}
	case yaml.ScalarNode:
		// A empty document with a document start marker without comments results in null
		if node.ShortTag() == "!!null" {
			return branch, nil
		}
		return nil, fmt.Errorf("YAML documents that are values are not supported")
	case yaml.AliasNode:
		branch, err = store.appendYamlNodeToTreeBranch(node.Alias, branch, false)
	}
	if !commentsWereHandled {
		branch = store.appendCommentToMap(node.FootComment, branch)
	}
	return branch, nil
}

func (store Store) yamlDocumentNodeToTreeBranch(in yaml.Node) (sops.TreeBranch, error) {
	branch := make(sops.TreeBranch, 0)
	return store.appendYamlNodeToTreeBranch(&in, branch, false)
}

func (store *Store) addCommentsHead(node *yaml.Node, comments []string) []string {
	if len(comments) > 0 {
		comment := "#" + strings.Join(comments, "\n#")
		if len(node.HeadComment) > 0 {
			node.HeadComment = comment + "\n" + node.HeadComment
		} else {
			node.HeadComment = comment
		}
	}
	return nil
}

func (store *Store) addCommentsFoot(node *yaml.Node, comments []string) []string {
	if len(comments) > 0 {
		comment := "#" + strings.Join(comments, "\n#")
		if len(node.FootComment) > 0 {
			node.FootComment += "\n" + comment
		} else {
			node.FootComment = comment
		}
	}
	return nil
}

func (store *Store) treeValueToNode(in interface{}) *yaml.Node {
	switch in := in.(type) {
	case sops.TreeBranch:
		var mapping = &yaml.Node{}
		mapping.Kind = yaml.MappingNode
		store.appendTreeBranch(in, mapping)
		return mapping
	case []interface{}:
		var sequence = &yaml.Node{}
		sequence.Kind = yaml.SequenceNode
		store.appendSequence(in, sequence)
		return sequence
	default:
		var valueNode = &yaml.Node{}
		valueNode.Encode(in)
		return valueNode
	}
}

func (store *Store) appendSequence(in []interface{}, sequence *yaml.Node) {
	var comments []string
	var beginning bool = true
	for _, item := range in {
		if comment, ok := item.(sops.Comment); ok {
			comments = append(comments, comment.Value)
		} else {
			if beginning {
				comments = store.addCommentsHead(sequence, comments)
				beginning = false
			}
			itemNode := store.treeValueToNode(item)
			comments = store.addCommentsHead(itemNode, comments)
			sequence.Content = append(sequence.Content, itemNode)
		}
	}
	if len(comments) > 0 {
		if beginning {
			comments = store.addCommentsHead(sequence, comments)
		} else {
			comments = store.addCommentsFoot(sequence.Content[len(sequence.Content) - 1], comments)
		}
	}
}

func (store *Store) appendTreeBranch(branch sops.TreeBranch, mapping *yaml.Node) {
	var comments []string
	var beginning bool = true
	for _, item := range branch {
		if comment, ok := item.Key.(sops.Comment); ok {
			comments = append(comments, comment.Value)
		} else {
			if beginning {
				comments = store.addCommentsHead(mapping, comments)
				beginning = false
			}
			var keyNode = &yaml.Node{}
			keyNode.Encode(item.Key)
			comments = store.addCommentsHead(keyNode, comments)
			valueNode := store.treeValueToNode(item.Value)
			mapping.Content = append(mapping.Content, keyNode, valueNode)
		}
	}
	if len(comments) > 0 {
		if beginning {
			comments = store.addCommentsHead(mapping, comments)
		} else {
			comments = store.addCommentsFoot(mapping.Content[len(mapping.Content) - 1], comments)
		}
	}
}

// LoadEncryptedFile loads the contents of an encrypted yaml file onto a
// sops.Tree runtime object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	// Because we don't know what fields the input file will have, we have to
	// load the file in two steps.
	// First, we load the file's metadata, the structure of which is known.
	metadataHolder := stores.SopsFile{}
	err := yaml.Unmarshal(in, &metadataHolder)
	if err != nil {
		return sops.Tree{}, fmt.Errorf("Error unmarshalling input yaml: %s", err)
	}
	if metadataHolder.Metadata == nil {
		return sops.Tree{}, sops.MetadataNotFound
	}
	metadata, err := metadataHolder.Metadata.ToInternal()
	if err != nil {
		return sops.Tree{}, err
	}
	var data yaml.Node
	if err := yaml.Unmarshal(in, &data); err != nil {
		return sops.Tree{}, fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	var branches sops.TreeBranches
	d := yaml.NewDecoder(bytes.NewReader(in))
	for true {
		var data yaml.Node
		err := d.Decode(&data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return sops.Tree{}, fmt.Errorf("Error unmarshaling input YAML: %s", err)
		}

		branch, err := store.yamlDocumentNodeToTreeBranch(data)
		if err != nil {
			return sops.Tree{}, fmt.Errorf("Error unmarshaling input YAML: %s", err)
		}

		for i, elt := range branch {
			if elt.Key == "sops" { // Erase
				branch = append(branch[:i], branch[i+1:]...)
			}
		}
		branches = append(branches, branch)
	}
	return sops.Tree{
		Branches: branches,
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads the contents of a plaintext yaml file onto a
// sops.Tree runtime object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	var branches sops.TreeBranches
	d := yaml.NewDecoder(bytes.NewReader(in))
	for true {
		var data yaml.Node
		err := d.Decode(&data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling input YAML: %s", err)
		}

		branch, err := store.yamlDocumentNodeToTreeBranch(data)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling input YAML: %s", err)
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

// EmitEncryptedFile returns the encrypted bytes of the yaml file corresponding to a
// sops.Tree runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
    var b bytes.Buffer
	e := yaml.NewEncoder(io.Writer(&b))
	e.SetIndent(4)
	for _, branch := range in.Branches {
		// Document root
		var doc = yaml.Node{}
		doc.Kind = yaml.DocumentNode
		// Add global mapping
		var mapping = yaml.Node{}
		mapping.Kind = yaml.MappingNode
		doc.Content = append(doc.Content, &mapping)
		// Create copy of branch with metadata appended
		branch = append(sops.TreeBranch(nil), branch...)
		branch = append(branch, sops.TreeItem{
			Key: "sops",
			Value: stores.MetadataFromInternal(in.Metadata),
		})
		// Marshal branch to global mapping node
		store.appendTreeBranch(branch, &mapping)
		// Encode YAML
		err := e.Encode(&doc)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
		}
	}
	e.Close()
	return b.Bytes(), nil
}

// EmitPlainFile returns the plaintext bytes of the yaml file corresponding to a
// sops.TreeBranches runtime object
func (store *Store) EmitPlainFile(branches sops.TreeBranches) ([]byte, error) {
    var b bytes.Buffer
	e := yaml.NewEncoder(io.Writer(&b))
	e.SetIndent(4)
	for _, branch := range branches {
		// Document root
		var doc = yaml.Node{}
		doc.Kind = yaml.DocumentNode
		// Add global mapping
		var mapping = yaml.Node{}
		mapping.Kind = yaml.MappingNode
		// Marshal branch to global mapping node
		store.appendTreeBranch(branch, &mapping)
		if len(mapping.Content) == 0 {
			doc.HeadComment = mapping.HeadComment
		} else {
			doc.Content = append(doc.Content, &mapping)
		}
		// Encode YAML
		err := e.Encode(&doc)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
		}
	}
	e.Close()
	return b.Bytes(), nil
}

// EmitValue returns bytes corresponding to a single encoded value
// in a generic interface{} object
func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	n := store.treeValueToNode(v)
	return yaml.Marshal(n)
}

// EmitExample returns the bytes corresponding to an example complex tree
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleComplexTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}
