// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestMetadata(t *testing.T) {
	md := &Metadata{EventID: "test event ID"}
	ctx := NewContext(context.Background(), md)
	got, err := FromContext(ctx)
	if err != nil {
		t.Fatalf("FromContext error: %v", err)
	}
	if !reflect.DeepEqual(got, md) {
		t.Fatalf("FromContext\nGot %v\nWant %v", got, md)
	}
}

func TestMetadataError(t *testing.T) {
	if _, err := FromContext(nil); err == nil {
		t.Errorf("FromContext got no error, wanted an error")
	}
	if _, err := FromContext(context.Background()); err == nil {
		t.Errorf("FromContext got no error, wanted an error")
	}
	if _, err := FromContext(NewContext(context.Background(), nil)); err == nil {
		t.Errorf("FromContext got no error, wanted an error")
	}
}
