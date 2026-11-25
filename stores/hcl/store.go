/*
Package hcl implements a SOPS store for HCL (HashiCorp Configuration Language)
files such as Terraform `.tfvars` and Terragrunt `.hcl` files.

Mapping between HCL and the SOPS tree:

  - An attribute becomes a tree item whose key is the attribute name.
  - A block becomes a tree item whose key is the block header as it is written
    in the file, including labels and the opening brace, for example
    `resource "aws_instance" "web" {`. The trailing brace is what tells blocks
    apart from attributes holding an object. The value is the block body.
  - Literal strings, numbers and booleans become strings, ints/floats and bools,
    and are encrypted. `null` becomes nil.
  - Tuples and objects made only of literals become lists and mappings.
  - Every other expression (references, function calls, conditionals,
    templates with interpolation, ...) is opaque: its source text is wrapped as
    `${<source>}` and treated as a string. It is therefore encrypted like any
    other string, and written back verbatim on decryption. This is the same
    convention HCL's JSON syntax uses for expressions.
  - Comments become sops.Comment items.
*/
package hcl //import "github.com/getsops/sops/v3/stores/hcl"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/stores"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/text/unicode/norm"
)

const indent = "  "

// Store handles storage of HCL data.
type Store struct {
	config *config.HCLStoreConfig
}

// NewStore creates a new HCL store.
func NewStore(c *config.HCLStoreConfig) *Store {
	return &Store{config: c}
}

// Name returns the name of the store.
func (store *Store) Name() string {
	return "hcl"
}

// isExpr reports whether a tree string stands for an opaque expression.
func isExpr(s string) bool {
	return len(s) >= 3 && strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}")
}

func quote(s string) string {
	return string(hclwrite.TokensForValue(cty.StringVal(s)).Bytes())
}

// ---------------------------------------------------------------------------
// HCL -> tree

type span struct {
	start, end int
}

type comment struct {
	span
	text string
}

type loader struct {
	src      []byte
	comments []comment
}

func (l *loader) hasComment(s span) bool {
	for _, c := range l.comments {
		if c.start >= s.start && c.end <= s.end {
			return true
		}
	}
	return false
}

func commentText(raw string) string {
	switch {
	case strings.HasPrefix(raw, "/*"):
		raw = strings.TrimSuffix(strings.TrimPrefix(raw, "/*"), "*/")
		return strings.TrimSpace(raw)
	case strings.HasPrefix(raw, "//"):
		raw = raw[2:]
	default:
		raw = strings.TrimPrefix(raw, "#")
	}
	return strings.TrimPrefix(strings.TrimRight(raw, "\r\n"), " ")
}

func (l *loader) opaque(e hclsyntax.Expression) string {
	r := e.Range()
	return "${" + string(l.src[r.Start.Byte:r.End.Byte]) + "}"
}

func (l *loader) literal(s string) (interface{}, error) {
	if isExpr(s) {
		return nil, fmt.Errorf("string %q starts with ${ and ends with }, which cannot be told apart from an expression", s)
	}
	return s, nil
}

func number(v cty.Value) (interface{}, bool) {
	bf := v.AsBigFloat()
	if bf.IsInt() {
		if i, acc := bf.Int64(); acc == big.Exact && int64(int(i)) == i {
			return int(i), true
		}
	}
	f, _ := bf.Float64()
	if r, err := cty.ParseNumberVal(strconv.FormatFloat(f, 'f', -1, 64)); err == nil && r.RawEquals(v) {
		return f, true
	}
	return nil, false
}

// nfc reports whether the source of e is NFC-normal. HCL normalizes decoded
// strings to NFC, which would change the bytes of anything that is not.
func (l *loader) nfc(r hcl.Range) bool {
	src := l.src[r.Start.Byte:r.End.Byte]
	// \u and \U escapes are ASCII in the source but can decode to non-NFC text.
	return norm.NFC.IsNormal(src) && !bytes.Contains(src, []byte(`\u`)) && !bytes.Contains(src, []byte(`\U`))
}

func (l *loader) value(e hclsyntax.Expression) (interface{}, error) {
	r := e.Range()
	rs := span{r.Start.Byte, r.End.Byte}
	switch e := e.(type) {
	case *hclsyntax.UnaryOpExpr:
		if lv, ok := e.Val.(*hclsyntax.LiteralValueExpr); ok && e.Op == hclsyntax.OpNegate && lv.Val.Type() == cty.Number {
			if n, ok := number(lv.Val.Negate()); ok {
				return n, nil
			}
		}
	case *hclsyntax.LiteralValueExpr:
		switch {
		case e.Val.IsNull():
			return nil, nil
		case e.Val.Type() == cty.Bool:
			return e.Val.True(), nil
		case e.Val.Type() == cty.Number:
			if n, ok := number(e.Val); ok {
				return n, nil
			}
		}
	case *hclsyntax.TemplateExpr:
		var sb strings.Builder
		literal := true
		for _, p := range e.Parts {
			lv, ok := p.(*hclsyntax.LiteralValueExpr)
			if !ok || lv.Val.Type() != cty.String {
				literal = false
				break
			}
			sb.WriteString(lv.Val.AsString())
		}
		if literal && l.nfc(r) {
			return l.literal(sb.String())
		}
	case *hclsyntax.TupleConsExpr:
		if l.hasComment(rs) {
			break
		}
		out := make([]interface{}, 0, len(e.Exprs))
		for _, item := range e.Exprs {
			v, err := l.value(item)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case *hclsyntax.ObjectConsExpr:
		if l.hasComment(rs) {
			break
		}
		out := sops.TreeBranch{}
		ok := true
		for _, item := range e.Items {
			key := l.objectKey(item.KeyExpr)
			if key == nil {
				ok = false
				break
			}
			v, err := l.value(item.ValueExpr)
			if err != nil {
				return nil, err
			}
			out = append(out, sops.TreeItem{Key: *key, Value: v})
		}
		if ok {
			return out, nil
		}
	}
	return l.opaque(e), nil
}

// objectKey returns the key of an object entry if it is a plain identifier or
// a literal string, and nil otherwise.
func (l *loader) objectKey(e hclsyntax.Expression) *string {
	if kw := hcl.ExprAsKeyword(e); kw != "" {
		return &kw
	}
	if w, ok := e.(*hclsyntax.ObjectConsKeyExpr); ok {
		e = w.Wrapped
	}
	if t, ok := e.(*hclsyntax.TemplateExpr); ok && len(t.Parts) == 1 {
		if lv, ok := t.Parts[0].(*hclsyntax.LiteralValueExpr); ok && lv.Val.Type() == cty.String {
			if !l.nfc(t.Range()) {
				return nil
			}
			s := lv.Val.AsString()
			return &s
		}
	}
	return nil
}

type bodyItem struct {
	span
	attr  *hclsyntax.Attribute
	block *hclsyntax.Block
}

// body converts the contents of body, which spans the bytes [lo, hi) of the
// source, to a tree branch.
func (l *loader) body(b *hclsyntax.Body, lo, hi int) (sops.TreeBranch, error) {
	var items []bodyItem
	for _, a := range b.Attributes {
		items = append(items, bodyItem{span{a.SrcRange.Start.Byte, a.SrcRange.End.Byte}, a, nil})
	}
	for _, blk := range b.Blocks {
		items = append(items, bodyItem{span{blk.TypeRange.Start.Byte, blk.CloseBraceRange.End.Byte}, nil, blk})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].start < items[j].start })

	var branch sops.TreeBranch
	// addComments emits the comments located in [from, to).
	addComments := func(from, to int, prev *bodyItem) {
		for _, c := range l.comments {
			if c.start < from || c.end > to {
				continue
			}
			inline := prev != nil && prev.attr != nil &&
				!strings.Contains(string(l.src[prev.end:c.start]), "\n")
			branch = append(branch, sops.TreeItem{Key: sops.Comment{Value: c.text, Inline: inline}})
		}
	}
	pos := lo
	var prev *bodyItem
	for i := range items {
		it := &items[i]
		addComments(pos, it.start, prev)
		var ti sops.TreeItem
		if it.attr != nil {
			v, err := l.value(it.attr.Expr)
			if err != nil {
				return nil, fmt.Errorf("attribute %q: %w", it.attr.Name, err)
			}
			ti = sops.TreeItem{Key: it.attr.Name, Value: v}
		} else {
			blk := it.block
			header := blk.Type
			for i, label := range blk.Labels {
				if !l.nfc(blk.LabelRanges[i]) {
					return nil, fmt.Errorf("block label %q is not NFC-normalized, which HCL cannot represent", label)
				}
				header += " " + quote(label)
			}
			inner, err := l.body(blk.Body, blk.OpenBraceRange.End.Byte, blk.CloseBraceRange.Start.Byte)
			if err != nil {
				return nil, err
			}
			ti = sops.TreeItem{Key: header + " {", Value: inner}
		}
		branch = append(branch, ti)
		pos = it.end
		prev = it
	}
	addComments(pos, hi, prev)
	return branch, nil
}

// LoadPlainFile loads a plaintext HCL file's bytes onto a sops.TreeBranches runtime object
func (store *Store) LoadPlainFile(in []byte) (sops.TreeBranches, error) {
	file, diags := hclsyntax.ParseConfig(in, "", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, fmt.Errorf("Could not unmarshal input data: %s", diags.Error())
	}
	tokens, diags := hclsyntax.LexConfig(in, "", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, fmt.Errorf("Could not unmarshal input data: %s", diags.Error())
	}
	l := &loader{src: in}
	for _, t := range tokens {
		if t.Type == hclsyntax.TokenComment {
			l.comments = append(l.comments, comment{span{t.Range.Start.Byte, t.Range.End.Byte}, commentText(string(t.Bytes))})
		}
	}
	branch, err := l.body(file.Body.(*hclsyntax.Body), 0, len(in))
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal input data: %s", err)
	}
	return sops.TreeBranches{branch}, nil
}

// LoadEncryptedFile loads encrypted HCL file's bytes onto a sops.Tree runtime object
func (store *Store) LoadEncryptedFile(in []byte) (sops.Tree, error) {
	branches, err := store.LoadPlainFile(in)
	if err != nil {
		return sops.Tree{}, err
	}
	branches, metadata, err := stores.ExtractMetadata(branches, stores.MetadataOpts{})
	if err != nil {
		return sops.Tree{}, err
	}
	return sops.Tree{
		Branches: branches,
		Metadata: metadata,
	}, nil
}

// ---------------------------------------------------------------------------
// tree -> HCL

func isIdent(s string) bool {
	return hclsyntax.ValidIdentifier(s)
}

func scalar(v interface{}) (string, error) {
	switch v := v.(type) {
	case nil:
		return "null", nil
	case string:
		if isExpr(v) {
			return v[2 : len(v)-1], nil
		}
		return quote(v), nil
	case []byte:
		return quote(string(v)), nil
	case bool:
		return strconv.FormatBool(v), nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		if v != v || v > 1.7976931348623157e308 || v < -1.7976931348623157e308 {
			return "", fmt.Errorf("cannot represent %v in HCL", v)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	}
	return "", fmt.Errorf("unsupported value type %T", v)
}

func writeComments(sb *strings.Builder, c sops.Comment, pad string) {
	for _, line := range strings.Split(c.Value, "\n") {
		sb.WriteString(pad + strings.TrimRight("# "+line, " ") + "\n")
	}
}

func writeValue(sb *strings.Builder, v interface{}, pad string) error {
	switch v := v.(type) {
	case sops.TreeBranch:
		if len(v) == 0 {
			sb.WriteString("{}")
			return nil
		}
		sb.WriteString("{\n")
		if err := writeItems(sb, v, pad+indent, false); err != nil {
			return err
		}
		sb.WriteString(pad + "}")
	case []interface{}:
		if len(v) == 0 {
			sb.WriteString("[]")
			return nil
		}
		sb.WriteString("[\n")
		for _, e := range v {
			if c, ok := e.(sops.Comment); ok {
				writeComments(sb, c, pad+indent)
				continue
			}
			sb.WriteString(pad + indent)
			if err := writeValue(sb, e, pad+indent); err != nil {
				return err
			}
			sb.WriteString(",\n")
		}
		sb.WriteString(pad + "]")
	default:
		s, err := scalar(v)
		if err != nil {
			return err
		}
		sb.WriteString(s)
	}
	return nil
}

func isBlock(key string) bool {
	return strings.HasSuffix(key, "{")
}

// writeItems writes the items of a branch. Blocks are only allowed if blocks is true.
func writeItems(sb *strings.Builder, branch sops.TreeBranch, pad string, blocks bool) error {
	var pending []sops.Comment
	prevBlock, first := false, true
	// open reports that the last attribute line is still waiting for its newline,
	// so that an inline comment can be appended to it.
	open := false
	closeLine := func() {
		if open {
			sb.WriteString("\n")
			open = false
		}
	}
	flush := func() {
		for _, c := range pending {
			writeComments(sb, c, pad)
		}
		pending = nil
	}
	for _, item := range branch {
		if c, ok := item.Key.(sops.Comment); ok {
			if c.Inline && open {
				sb.WriteString(" # " + strings.ReplaceAll(c.Value, "\n", " "))
				closeLine()
			} else {
				closeLine()
				pending = append(pending, c)
			}
			continue
		}
		closeLine()
		key, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("unexpected key type %T", item.Key)
		}
		block := blocks && isBlock(key)
		if !first && (block || prevBlock) {
			sb.WriteString("\n")
		}
		first, prevBlock = false, block
		flush()
		if block {
			sb.WriteString(pad + key + "\n")
			inner, ok := item.Value.(sops.TreeBranch)
			if !ok {
				return fmt.Errorf("block %q: value must be a mapping", key)
			}
			if err := writeItems(sb, inner, pad+indent, true); err != nil {
				return err
			}
			sb.WriteString(pad + "}\n")
			continue
		}
		if !isIdent(key) {
			key = quote(key)
		}
		sb.WriteString(pad + key + " = ")
		if err := writeValue(sb, item.Value, pad); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		open = true
	}
	closeLine()
	flush()
	return nil
}

// EmitPlainFile returns the plaintext HCL file bytes corresponding to a sops.TreeBranches object
func (store *Store) EmitPlainFile(in sops.TreeBranches) ([]byte, error) {
	var sb strings.Builder
	for _, branch := range in {
		if err := writeItems(&sb, branch, "", true); err != nil {
			return nil, fmt.Errorf("Error marshaling to HCL: %s", err)
		}
	}
	out := hclwrite.Format([]byte(sb.String()))
	// Guard against values that would produce a syntactically invalid file.
	if _, diags := hclsyntax.ParseConfig(out, "", hcl.InitialPos); diags.HasErrors() {
		return nil, fmt.Errorf("Error marshaling to HCL: %s", diags.Error())
	}
	return out, nil
}

// EmitEncryptedFile returns encrypted HCL file bytes corresponding to a sops.Tree
// runtime object
func (store *Store) EmitEncryptedFile(in sops.Tree) ([]byte, error) {
	branches, err := stores.SerializeMetadata(in, stores.MetadataOpts{})
	if err != nil {
		return nil, fmt.Errorf("Error marshaling metadata: %s", err)
	}
	return store.EmitPlainFile(branches)
}

// EmitValue returns a single value encoded in a generic interface{} as bytes
func (store *Store) EmitValue(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case sops.TreeBranches:
		return store.EmitPlainFile(v)
	default:
		return json.Marshal(v)
	}
}

// EmitExample returns the plaintext HCL file bytes corresponding to the FlatTree example
func (store *Store) EmitExample() []byte {
	bytes, err := store.EmitPlainFile(stores.ExampleFlatTree.Branches)
	if err != nil {
		panic(err)
	}
	return bytes
}

// HasSopsTopLevelKey checks whether a top-level "sops" key exists.
func (store *Store) HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	return stores.HasSopsTopLevelKey(branch)
}
