/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gopgagent /* import "go.mozilla.org/gopgagent" */

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPrompt(t *testing.T) {
	if os.Getenv("TEST_GPGAGENT_LIB") != "1" {
		t.Logf("skipping TestPrompt without $TEST_GPGAGENT_LIB == 1")
		return
	}
	conn, err := NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	req := &PassphraseRequest{
		Desc:     "Type 'foo' for testing",
		Error:    "seriously, or I'll be an error.",
		Prompt:   "foo",
		CacheKey: fmt.Sprintf("gpgagent_test-cachekey-%d", time.Now()),
	}
	s1, err := conn.GetPassphrase(req)
	if err != nil {
		t.Fatal(err)
	}
	t1 := time.Now()
	s2, err := conn.GetPassphrase(req)
	if err != nil {
		t.Fatal(err)
	}
	t2 := time.Now()
	if td := t2.Sub(t1); td > 1e9/5 {
		t.Errorf("cached passphrase took more than 1/5 second; took %d ns", td)
	}
	if s1 != s2 {
		t.Errorf("cached passphrase differed; got %q, want %q", s2, s1)
	}
	if s1 != "foo" {
		t.Errorf("got passphrase %q; want %q", s1, "foo")
	}
	err = conn.RemoveFromCache(req.CacheKey)
	if err != nil {
		t.Fatal(err)
	}

	req.NoAsk = true
	s3, err := conn.GetPassphrase(req)
	if err != ErrNoData {
		t.Errorf("after remove from cache, expected gpgagent.ErrNoData, got %q, %v", s3, err)
	}

	s4, err := conn.GetPassphrase(&PassphraseRequest{
		Desc:     "Press Cancel for testing",
		Error:    "seriously, or I'll be an error.",
		Prompt:   "cancel!",
		CacheKey: fmt.Sprintf("gpgagent_test-cachekey-%d", time.Now()),
	})
	if err != ErrCancel {
		t.Errorf("expected cancel, got %q, %v", s4, err)
	}
}
