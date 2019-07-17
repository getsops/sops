// copyright 2016 google inc. all rights reserved.
//
// licensed under the apache license, version 2.0 (the "license");
// you may not use this file except in compliance with the license.
// you may obtain a copy of the license at
//
//     http://www.apache.org/licenses/license-2.0
//
// unless required by applicable law or agreed to in writing, software
// distributed under the license is distributed on an "as is" basis,
// without warranties or conditions of any kind, either express or implied.
// see the license for the specific language governing permissions and
// limitations under the license.

package martianurl

import "testing"

func TestMatchHost(t *testing.T) {
	tt := []struct {
		host, match string
		want        bool
	}{
		{"example.com", "example.com", true},
		{"example.com", "example.org", false},
		{"ample.com", "example.com", false},
		{"example.com", "ample.com", false},
		{"example.com", "example.*", true},
		{"www.example.com", "*.example.com", true},
		{"one.two.example.com", "*.example.com", false},
		{"one.two.example.com", "*.*.example.com", true},
		{"", "", false},
		{"", "foo", false},
	}

	for i, tc := range tt {
		if got := MatchHost(tc.host, tc.match); got != tc.want {
			t.Errorf("%d. MatchHost(%s, %s): got %t, want %t", i, tc.host, tc.match, got, tc.want)
		}
	}
}
