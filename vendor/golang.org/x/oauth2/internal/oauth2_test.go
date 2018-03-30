// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseINI(t *testing.T) {
	tests := []struct {
		ini  string
		want map[string]map[string]string
	}{
		{
			`root = toor
[foo]  
bar = hop
ini = nin
`,
			map[string]map[string]string{
				"":    {"root": "toor"},
				"foo": {"bar": "hop", "ini": "nin"},
			},
		},
		{
			`[empty]
[section]
empty=
`,
			map[string]map[string]string{
				"":        {},
				"empty":   {},
				"section": {"empty": ""},
			},
		},
		{
			`ignore
[invalid
=stuff
;comment=true
`,
			map[string]map[string]string{
				"": {},
			},
		},
	}
	for _, tt := range tests {
		result, err := ParseINI(strings.NewReader(tt.ini))
		if err != nil {
			t.Errorf("ParseINI(%q) error %v, want: no error", tt.ini, err)
			continue
		}
		if !reflect.DeepEqual(result, tt.want) {
			t.Errorf("ParseINI(%q) = %#v, want: %#v", tt.ini, result, tt.want)
		}
	}
}
