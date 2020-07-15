package dotenv

// The dotenv parser is designed around the following rules:
//
// Comments:
//
// * Comments may be written by starting a line with the `#` character.
//   End-of-line comments are not currently supported, as there is no way to
//   encode a comment's position in a `sops.TreeItem`.
//
// Newline handling:
//
// * If a value is unquoted or single-quoted and contains the character
//   sequence `\n` (`0x5c6e`), it IS NOT decoded to a line feed (`0x0a`).
//
// * If a value is double-quoted and contains the character sequence `\n`
//   (`0x5c6e`), it IS decoded to a line feed (`0x0a`).
//
// Whitespace trimming:
//
// * For comments, the whitespace immediately after the `#` character and any
//   trailing whitespace is trimmed.
//
// * If a value is unquoted and contains any leading or trailing whitespace, it
//   is trimmed.
//
// * If a value is either single- or double-quoted and contains any leading or
//   trailing whitespace, it is left untrimmed.
//
// Quotation handling:
//
// * If a value is surrounded by single- or double-quotes, the quotation marks
//   are interpreted and not included in the value.
//
// * Any number of single-quote characters may appear in a double-quoted
//   value, or within a single-quoted value if they are escaped (i.e.,
//   `'foo\'bar'`).
//
// * Any number of double-quote characters may appear in a single-quoted
//   value, or within a double-quoted value if they are escaped (i.e.,
//   `"foo\"bar"`).

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"go.mozilla.org/sops/v3"
)

var KeyRegexp = regexp.MustCompile(`^[A-Za-z_]+[A-Za-z0-9_]*$`)

func parse(data []byte) (items []sops.TreeItem, err error) {
	reader := bytes.NewReader(data)

	for {
		var b byte
		var item *sops.TreeItem

		b, err = reader.ReadByte()

		if err != nil {
			break
		}

		if isWhitespace(b) {
			continue
		}

		if b == '#' {
			item, err = parseComment(reader)
		} else {
			reader.UnreadByte()
			item, err = parseKeyValue(reader)
		}

		if err != nil {
			break
		}

		if item == nil {
			continue
		}

		items = append(items, *item)
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func parseComment(reader io.ByteScanner) (item *sops.TreeItem, err error) {
	var builder strings.Builder
	var whitespace bytes.Buffer

	for {
		var b byte
		b, err = reader.ReadByte()

		if err != nil {
			break
		}

		if b == '\n' {
			break
		}

		if isWhitespace(b) {
			whitespace.WriteByte(b)
			continue
		}

		if builder.Len() == 0 {
			whitespace.Reset()
		}

		_, err = io.Copy(&builder, &whitespace)

		if err != nil {
			break
		}

		builder.WriteByte(b)
	}

	if builder.Len() == 0 {
		return
	}

	item = &sops.TreeItem{Key: sops.Comment{builder.String()}, Value: nil}
	return
}

func parseKeyValue(reader io.ByteScanner) (item *sops.TreeItem, err error) {
	var key, value string

	key, err = parseKey(reader)
	if err != nil {
		return
	}

	value, err = parseValue(reader)
	if err != nil {
		return
	}

	item = &sops.TreeItem{Key: key, Value: value}
	return
}

func parseKey(reader io.ByteScanner) (key string, err error) {
	var builder strings.Builder

	for {
		var b byte
		b, err = reader.ReadByte()

		if err != nil {
			break
		}

		if b == '=' {
			break
		}

		builder.WriteByte(b)
	}

	key = builder.String()

	if !KeyRegexp.MatchString(key) {
		err = fmt.Errorf("invalid dotenv key: %q", key)
	}

	return
}

func parseValue(reader io.ByteScanner) (value string, err error) {
	var first byte
	first, err = reader.ReadByte()

	if err != nil {
		return
	}

	if first == '\'' {
		return parseSingleQuoted(reader)
	}

	if first == '"' {
		return parseDoubleQuoted(reader)
	}

	reader.UnreadByte()
	return parseUnquoted(reader)
}

func parseSingleQuoted(reader io.ByteScanner) (value string, err error) {
	var builder strings.Builder
	escaping := false

	for {
		var b byte
		b, err = reader.ReadByte()

		if err != nil {
			break
		}

		if !escaping && b == '\'' {
			break
		}

		if !escaping && b == '\\' {
			escaping = true
			continue
		}

		if escaping && b != '\'' {
			builder.WriteByte('\\')
		}

		escaping = false
		builder.WriteByte(b)
	}

	value = builder.String()
	return
}

func parseDoubleQuoted(reader io.ByteScanner) (value string, err error) {
	var builder strings.Builder
	escaping := false

	for {
		var b byte
		b, err = reader.ReadByte()

		if err != nil {
			break
		}

		if !escaping && b == '"' {
			break
		}

		if !escaping && b == '\\' {
			escaping = true
			continue
		}

		if escaping && b == 'n' {
			b = '\n'
		} else if escaping && b != '"' {
			builder.WriteByte('\\')
		}

		escaping = false
		builder.WriteByte(b)
	}

	value = builder.String()
	return
}

func parseUnquoted(reader io.ByteScanner) (value string, err error) {
	var builder strings.Builder
	var whitespace bytes.Buffer

	for {
		var b byte
		b, err = reader.ReadByte()

		if err != nil {
			break
		}

		if b == '\n' {
			break
		}

		if isWhitespace(b) {
			whitespace.WriteByte(b)
			continue
		}

		if builder.Len() == 0 {
			whitespace.Reset()
		}

		_, err = io.Copy(&builder, &whitespace)

		if err != nil {
			break
		}

		builder.WriteByte(b)
	}

	value = builder.String()
	return
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}
