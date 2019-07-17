// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"os"
	"strings"
	"testing"

	stdlog "log"
)

func TestLog(t *testing.T) {
	buf := new(bytes.Buffer)

	stdlog.SetOutput(buf)
	defer stdlog.SetOutput(os.Stdout)

	// Reset log level after tests.
	defer func(l int) { level = l }(level)
	level = Debug

	Infof("log: %s test", "info")
	if got, want := buf.String(), "INFO: log: info test\n"; !strings.HasSuffix(got, want) {
		t.Errorf("Infof(): got %q, want to contain %q", got, want)
	}

	Debugf("log: %s test", "debug")
	if got, want := buf.String(), "DEBUG: log: debug test\n"; !strings.HasSuffix(got, want) {
		t.Errorf("Debugf(): got %q, want to contain %q", got, want)
	}

	Errorf("log: %s test", "error")
	if got, want := buf.String(), "ERROR: log: error test\n"; !strings.HasSuffix(got, want) {
		t.Errorf("Errorf(): got %q, want to contain %q", got, want)
	}
}
