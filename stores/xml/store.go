package xml //import "github.com/getsops/sops/v3/stores/xml"

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/stores"
)

const (
	// attributePrefix is prepended to the name of an XML attribute, to tell
	// attributes apart from child elements inside a tree branch.
	attributePrefix = "@"
	// textKey is the key holding the text of an element which also has
	// attributes or child elements.
	textKey = "#text"
)

// Store handles storage of XML data.
type Store struct {
	config config.XMLStoreConfig
}

func NewStore(c *config.XMLStoreConfig) *Store {
	return &Store{config: *c}
}

func (store *Store) Name() string {
	return "xml"
}

// registerPrefixes records the namespace prefixes declared on an element. The
// decoder resolves names to namespace URIs, so the declarations are needed to
// write the names back the way they were written in the document.
func registerPrefixes(start xml.StartElement, prefixes map[string]string) {
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			prefixes[attr.Value] = attr.Name.Local
		}
	}
}

// qualifiedName renders an element or attribute name the way it appeared in
// the document, by mapping its namespace URI back to the declared prefix.
func qualifiedName(name xml.Name, prefixes map[string]string) string {
	if name.Space == "xmlns" {
		return "xmlns:" + name.Local
	}
	if prefix, ok := prefixes[name.Space]; ok && prefix != "" {
		return prefix + ":" + name.Local
	}
	return name.Local
}

// addItem appends a child to a branch. Children sharing their name with an
// earlier sibling are collected into a list, so that repeated elements behave
// like the arrays of the other stores.
func addItem(branch sops.TreeBranch, key string, value interface{}) sops.TreeBranch {
	for i, item := range branch {
		if item.Key != key {
			continue
		}
		if list, ok := item.Value.([]interface{}); ok {
			branch[i].Value = append(list, value)
		} else {
			branch[i].Value = []interface{}{item.Value, value}
		}
		return branch
	}
	return append(branch, sops.TreeItem{Key: key, Value: value})
}

// elementFromXMLDecoder reads the element which has just been opened on the
// decoder and returns its name and its value. Elements without attributes and
// without children become plain strings, all others become tree branches.
func (store Store) elementFromXMLDecoder(dec *xml.Decoder, start xml.StartElement, prefixes map[string]string) (string, interface{}, error) {
	registerPrefixes(start, prefixes)
	name := qualifiedName(start.Name, prefixes)

	var branch sops.TreeBranch
	for _, attr := range start.Attr {
		branch = append(branch, sops.TreeItem{
			Key:   attributePrefix + qualifiedName(attr.Name, prefixes),
			Value: attr.Value,
		})
	}

	var text string
	for {
		token, err := dec.Token()
		if err != nil {
			return name, nil, err
		}
		switch token := token.(type) {
		case xml.StartElement:
			childName, childValue, err := store.elementFromXMLDecoder(dec, token, prefixes)
			if err != nil {
				return name, nil, err
			}
			branch = addItem(branch, childName, childValue)
		case xml.CharData:
			text += string(token)
		case xml.Comment:
			branch = append(branch, sops.TreeItem{
				Key:   sops.Comment{Value: string(token)},
				Value: nil,
			})
		case xml.EndElement:
			if len(branch) == 0 {
				// The text of a leaf element is kept as it is written.
				return name, text, nil
			}
			// The text of an element with attributes or children is only kept
			// when it is more than the indentation of the document.
			if text = strings.TrimSpace(text); text != "" {
				branch = append(branch, sops.TreeItem{Key: textKey, Value: text})
			}
			return name, branch, nil
		}
	}
}

// treeBranchFromXML parses an XML document into a sops.TreeBranch holding the
// root element of the document.
func (store Store) treeBranchFromXML(in []byte) (sops.TreeBranch, error) {
	dec := xml.NewDecoder(bytes.NewReader(in))
	prefixes := make(map[string]string)
	var branch sops.TreeBranch
	for {
		token, err := dec.Token()
		if err == io.EOF {
			return branch, nil
		}
		if err != nil {
			return nil, err
		}
		switch token := token.(type) {
		case xml.StartElement:
			name, value, err := store.elementFromXMLDecoder(dec, token, prefixes)
			if err != nil {
				return nil, err
			}
			branch = addItem(branch, name, value)
		case xml.Comment:
			branch = append(branch, sops.TreeItem{
				Key:   sops.Comment{Value: string(token)},
				Value: nil,
			})
		}
	}
}

// rootIndex returns the position of the root element within the first branch
// of an XML tree. An XML document has exactly one root element.
func rootIndex(branches sops.TreeBranches) (int, error) {
	index := -1
	if len(branches) > 0 {
		for i, item := range branches[0] {
			if _, ok := item.Key.(sops.Comment); ok {
				continue
			}
			if index != -1 {
				return 0, fmt.Errorf("XML files must have exactly one root element, found more than one")
			}
			index = i
		}
	}
	if index == -1 {
		return 0, fmt.Errorf("XML files must have exactly one root element, found none")
	}
	return index, nil
}

// withRootValue returns a copy of the tree branches in which the value of the
// root element is replaced by the given branch.
func withRootValue(branches sops.TreeBranches, index int, value sops.TreeBranch) sops.TreeBranches {
	branch := make(sops.TreeBranch, len(branches[0]))
	copy(branch, branches[0])
	branch[index].Value = value
	out := make(sops.TreeBranches, len(branches))
	copy(out, branches)
	out[0] = branch
	return out
}

// escape returns the XML escaped representation of a value. Newlines and tabs
// are escaped as character references as well, so that they survive parsing.
func escape(v interface{}) string {
	var buffer bytes.Buffer
	xml.EscapeText(&buffer, []byte(stores.ValToString(v)))
	return buffer.String()
}

// indentString returns the string used for a single level of indentation.
func (store Store) indentString() (string, error) {
	if store.config.Indent < -1 {
		return "", fmt.Errorf("XML Indentation parameter smaller than -1 is not accepted")
	}
	if store.config.Indent == -1 {
		return "\t", nil
	}
	return strings.Repeat(" ", store.config.Indent), nil
}

// encodeBranch writes the items of a branch as XML elements and comments.
func (store Store) encodeBranch(buffer *bytes.Buffer, branch sops.TreeBranch, indent string, depth int) error {
	for _, item := range branch {
		if comment, ok := item.Key.(sops.Comment); ok {
			buffer.WriteString(strings.Repeat(indent, depth) + "<!--" + comment.Value + "-->\n")
			continue
		}
		key, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("Error encoding element: unexpected key of type %T", item.Key)
		}
		if err := store.encodeElement(buffer, key, item.Value, indent, depth); err != nil {
			return err
		}
	}
	return nil
}

// encodeElement writes a single value as an XML element with the given name.
func (store Store) encodeElement(buffer *bytes.Buffer, name string, value interface{}, indent string, depth int) error {
	switch value := value.(type) {
	case []interface{}:
		// A list is written as a repetition of elements with the same name.
		for _, item := range value {
			if _, ok := item.(sops.Comment); ok {
				continue
			}
			if err := store.encodeElement(buffer, name, item, indent, depth); err != nil {
				return err
			}
		}
		return nil
	case sops.TreeBranch:
		return store.encodeBranchElement(buffer, name, value, indent, depth)
	default:
		buffer.WriteString(strings.Repeat(indent, depth) + "<" + name + ">" + escape(value) + "</" + name + ">\n")
		return nil
	}
}

// encodeBranchElement writes a branch as an XML element, turning the items
// prefixed with "@" back into attributes and the "#text" item back into the
// text of the element.
func (store Store) encodeBranchElement(buffer *bytes.Buffer, name string, branch sops.TreeBranch, indent string, depth int) error {
	var attributes, text string
	var children sops.TreeBranch
	for _, item := range branch {
		key, ok := item.Key.(string)
		if !ok {
			children = append(children, item)
			continue
		}
		switch {
		case strings.HasPrefix(key, attributePrefix):
			attributes += " " + strings.TrimPrefix(key, attributePrefix) + `="` + escape(item.Value) + `"`
		case key == textKey:
			text = escape(item.Value)
		default:
			children = append(children, item)
		}
	}

	prefix := strings.Repeat(indent, depth)
	buffer.WriteString(prefix + "<" + name + attributes + ">")
	if len(children) == 0 {
		buffer.WriteString(text + "</" + name + ">\n")
		return nil
	}
	buffer.WriteString("\n")
	if text != "" {
		buffer.WriteString(prefix + indent + text + "\n")
	}
	if err := store.encodeBranch(buffer, children, indent, depth+1); err != nil {
		return err
	}
	buffer.WriteString(prefix + "</" + name + ">\n")
	return nil
}

// xmlFromTreeBranch returns the XML document corresponding to a branch.
func (store Store) xmlFromTreeBranch(branch sops.TreeBranch) ([]byte, error) {
	if _, err := rootIndex(sops.TreeBranches{branch}); err != nil {
		return nil, err
	}
	indent, err := store.indentString()
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString(xml.Header)
	if err := store.encodeBranch(&buffer, branch, indent, 0); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// encodeValue encodes a single value. Branches and lists are encoded as XML
// elements, simple values as their plain string representation.
func (store Store) encodeValue(v interface{}) ([]byte, error) {
	indent, err := store.indentString()
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	switch v := v.(type) {
	case sops.TreeBranch:
		err = store.encodeBranch(&buffer, v, indent, 0)
	case []interface{}:
		err = store.encodeElement(&buffer, "value", v, indent, 0)
	default:
		buffer.WriteString(stores.ValToString(v))
	}
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// LoadEncryptedFile loads an encrypted XML file's bytes onto a sops.Tree
// runtime object. The metadata is stored inside the root element, because an
// XML document cannot have more than one root element.
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}
	index, err := rootIndex(branches)
	if err != nil {
		return sops.Tree{}, err
	}
	rootBranch, ok := branches[0][index].Value.(sops.TreeBranch)
	if !ok {
		return sops.Tree{}, sops.MetadataNotFound
	}
	rootBranches, metadata, err := stores.ExtractMetadata(sops.TreeBranches{rootBranch}, stores.MetadataOpts{
		Flatten: stores.MetadataFlattenNone,
	})
	if err != nil {
		return sops.Tree{}, err
	}
	return sops.Tree{
		Branches: withRootValue(branches, index, rootBranches[0]),
		Metadata: metadata,
	}, nil
}

// LoadPlainFile loads a plaintext XML file's bytes onto a sops.TreeBranches
// runtime object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	branch, err := store.treeBranchFromXML(in)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	if _, err := rootIndex(sops.TreeBranches{branch}); err != nil {
		return nil, err
	}
	return sops.TreeBranches{
		branch,
	}, nil
}

// EmitEncryptedFile returns encrypted XML file bytes corresponding to a
// sops.Tree runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	index, err := rootIndex(in.Branches)
	if err != nil {
		return nil, err
	}
	rootBranch, ok := in.Branches[0][index].Value.(sops.TreeBranch)
	if !ok {
		// A root element without children is turned into a branch, so that
		// the metadata can be stored inside of it.
		rootBranch = sops.TreeBranch{
			sops.TreeItem{Key: textKey, Value: in.Branches[0][index].Value},
		}
	}
	rootBranches, err := stores.SerializeMetadata(sops.Tree{
		Branches: sops.TreeBranches{rootBranch},
		Metadata: in.Metadata,
	}, stores.MetadataOpts{
		Flatten: stores.MetadataFlattenNone,
	})
	if err != nil {
		return nil, fmt.Errorf("Error marshaling metadata: %s", err)
	}
	return store.EmitPlainFile(withRootValue(in.Branches, index, rootBranches[0]))
}

// EmitPlainFile returns the plaintext XML file bytes corresponding to a
// sops.TreeBranches object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	if len(in) != 1 {
		return nil, fmt.Errorf("Error marshaling to XML: there must be exactly one tree branch")
	}
	out, err := store.xmlFromTreeBranch(in[0])
	if err != nil {
		return nil, fmt.Errorf("Error marshaling to XML: %s", err)
	}
	return out, nil
}

// EmitValue returns a single value encoded in a generic interface{} as bytes
func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	return store.encodeValue(v)
}

// EmitExample returns the plaintext XML file bytes corresponding to an example
// document
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   "example",
				Value: stores.ExampleComplexTree.Branches[0],
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return bytes
}

// HasSopsTopLevelKey checks whether a "sops" key exists inside the root
// element of the document.
func (store *Store) HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	index, err := rootIndex(sops.TreeBranches{branch})
	if err != nil {
		return false
	}
	rootBranch, ok := branch[index].Value.(sops.TreeBranch)
	if !ok {
		return false
	}
	return stores.HasSopsTopLevelKey(rootBranch)
}
