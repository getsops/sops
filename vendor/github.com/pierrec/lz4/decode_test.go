package lz4

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func unbase64(in string) []byte {
	p, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return p
}

func TestBlockDecode(t *testing.T) {
	appendLen := func(p []byte, size int) []byte {
		for size > 0xFF {
			p = append(p, 0xFF)
			size -= 0xFF
		}

		p = append(p, byte(size))
		return p
	}

	emitSeq := func(lit string, offset uint16, matchLen int) []byte {
		var b byte
		litLen := len(lit)
		if litLen < 15 {
			b = byte(litLen << 4)
			litLen = -1
		} else {
			b = 0xF0
			litLen -= 15
		}

		if matchLen < 4 || offset == 0 {
			out := []byte{b}
			if litLen >= 0 {
				out = appendLen(out, litLen)
			}
			return append(out, lit...)
		}

		matchLen -= 4
		if matchLen < 15 {
			b |= byte(matchLen)
			matchLen = -1
		} else {
			b |= 0x0F
			matchLen -= 15
		}

		out := []byte{b}
		if litLen >= 0 {
			out = appendLen(out, litLen)
		}

		if len(lit) > 0 {
			out = append(out, lit...)
		}

		out = append(out, byte(offset), byte(offset>>8))

		if matchLen >= 0 {
			out = appendLen(out, matchLen)
		}

		return out
	}
	concat := func(in ...[]byte) []byte {
		var p []byte
		for _, b := range in {
			p = append(p, b...)
		}
		return p
	}

	tests := []struct {
		name string
		src  []byte
		exp  []byte
	}{
		{
			"literal_only_short",
			emitSeq("hello", 0, 0),
			[]byte("hello"),
		},
		{
			"literal_only_long",
			emitSeq(strings.Repeat("A", 15+255+255+1), 0, 0),
			bytes.Repeat([]byte("A"), 15+255+255+1),
		},
		{
			"literal_only_long_1",
			emitSeq(strings.Repeat("A", 15), 0, 0),
			bytes.Repeat([]byte("A"), 15),
		},
		{
			"repeat_match_len",
			emitSeq("a", 1, 4),
			[]byte("aaaaa"),
		},
		{
			"repeat_match_len_2_seq",
			concat(emitSeq("a", 1, 4), emitSeq("B", 1, 4)),
			[]byte("aaaaaBBBBB"),
		},
		{
			"long_match",
			emitSeq("A", 1, 16),
			bytes.Repeat([]byte("A"), 17),
		},
		{
			"repeat_match_log_len_2_seq",
			concat(emitSeq("a", 1, 15), emitSeq("B", 1, 15), emitSeq("end", 0, 0)),
			[]byte(strings.Repeat("a", 16) + strings.Repeat("B", 16) + "end"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := make([]byte, len(test.exp))
			n := decodeBlock(buf, test.src)
			if n <= 0 {
				t.Log(-n)
			}

			if !bytes.Equal(buf, test.exp) {
				t.Fatalf("expected %q got %q", test.exp, buf)
			}
		})
	}
}
