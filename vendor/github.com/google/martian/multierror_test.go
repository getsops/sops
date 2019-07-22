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

package martian

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMultiError(t *testing.T) {
	merr := NewMultiError()

	if !merr.Empty() {
		t.Fatal("Empty(): got false, want true")
	}

	var errs []error
	for i := 0; i < 3; i++ {
		err := fmt.Errorf("%d. error", i)
		errs = append(errs, err)
		merr.Add(err)
	}

	if merr.Empty() {
		t.Fatal("Empty(): got true, want false")
	}
	if got, want := merr.Errors(), errs; !reflect.DeepEqual(got, want) {
		t.Errorf("Errors(): got %v, want %v", got, want)
	}

	want := "0. error\n1. error\n2. error"
	if got := merr.Error(); got != want {
		t.Errorf("Error(): got %q, want %q", got, want)
	}
}
