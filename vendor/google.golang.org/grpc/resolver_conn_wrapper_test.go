/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpc

import (
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc/resolver"
)

func TestParseTarget(t *testing.T) {
	for _, test := range []resolver.Target{
		{"dns", "", "google.com"},
		{"dns", "a.server.com", "google.com"},
		{"dns", "a.server.com", "google.com/?a=b"},
		{"passthrough", "", "/unix/socket/address"},
	} {
		str := test.Scheme + "://" + test.Authority + "/" + test.Endpoint
		got := parseTarget(str)
		if got != test {
			t.Errorf("parseTarget(%q) = %+v, want %+v", str, got, test)
		}
	}
}

func TestParseTargetString(t *testing.T) {
	for _, test := range []struct {
		targetStr string
		want      resolver.Target
	}{
		{"", resolver.Target{"", "", ""}},
		{":///", resolver.Target{"", "", ""}},
		{"a:///", resolver.Target{"a", "", ""}},
		{"://a/", resolver.Target{"", "a", ""}},
		{":///a", resolver.Target{"", "", "a"}},
		{"a://b/", resolver.Target{"a", "b", ""}},
		{"a:///b", resolver.Target{"a", "", "b"}},
		{"://a/b", resolver.Target{"", "a", "b"}},
		{"a://b/c", resolver.Target{"a", "b", "c"}},
		{"dns:///google.com", resolver.Target{"dns", "", "google.com"}},
		{"dns://a.server.com/google.com", resolver.Target{"dns", "a.server.com", "google.com"}},
		{"dns://a.server.com/google.com/?a=b", resolver.Target{"dns", "a.server.com", "google.com/?a=b"}},

		{"/", resolver.Target{"", "", "/"}},
		{"google.com", resolver.Target{"", "", "google.com"}},
		{"google.com/?a=b", resolver.Target{"", "", "google.com/?a=b"}},
		{"/unix/socket/address", resolver.Target{"", "", "/unix/socket/address"}},

		// If we can only parse part of the target.
		{"://", resolver.Target{"", "", "://"}},
		{"unix://domain", resolver.Target{"", "", "unix://domain"}},
		{"a:b", resolver.Target{"", "", "a:b"}},
		{"a/b", resolver.Target{"", "", "a/b"}},
		{"a:/b", resolver.Target{"", "", "a:/b"}},
		{"a//b", resolver.Target{"", "", "a//b"}},
		{"a://b", resolver.Target{"", "", "a://b"}},
	} {
		got := parseTarget(test.targetStr)
		if got != test.want {
			t.Errorf("parseTarget(%q) = %+v, want %+v", test.targetStr, got, test.want)
		}
	}
}

// The target string with unknown scheme should be kept unchanged and passed to
// the dialer.
func TestDialParseTargetUnknownScheme(t *testing.T) {
	for _, test := range []struct {
		targetStr string
		want      string
	}{
		{"/unix/socket/address", "/unix/socket/address"},

		// Special test for "unix:///".
		{"unix:///unix/socket/address", "unix:///unix/socket/address"},

		// For known scheme.
		{"passthrough://a.server.com/google.com", "google.com"},
	} {
		dialStrCh := make(chan string, 1)
		cc, err := Dial(test.targetStr, WithInsecure(), WithDialer(func(t string, _ time.Duration) (net.Conn, error) {
			dialStrCh <- t
			return nil, fmt.Errorf("test dialer, always error")
		}))
		if err != nil {
			t.Fatalf("Failed to create ClientConn: %v", err)
		}
		got := <-dialStrCh
		cc.Close()
		if got != test.want {
			t.Errorf("Dial(%q), dialer got %q, want %q", test.targetStr, got, test.want)
		}
	}
}
