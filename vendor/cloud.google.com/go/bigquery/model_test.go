// Copyright 2019 Google LLC
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

package bigquery

import (
	"testing"
	"time"

	"cloud.google.com/go/internal/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
	bq "google.golang.org/api/bigquery/v2"
)

func TestBQToModelMetadata(t *testing.T) {
	aTime := time.Date(2019, 3, 14, 0, 0, 0, 0, time.Local)
	aTimeMillis := aTime.UnixNano() / 1e6
	for _, test := range []struct {
		in   *bq.Model
		want *ModelMetadata
	}{
		{&bq.Model{}, &ModelMetadata{}},
		{
			&bq.Model{
				CreationTime:     aTimeMillis,
				Description:      "desc",
				Etag:             "etag",
				ExpirationTime:   aTimeMillis,
				FriendlyName:     "fname",
				LastModifiedTime: aTimeMillis,
				Location:         "loc",
				Labels:           map[string]string{"a": "b"},
			},
			&ModelMetadata{
				CreationTime:     aTime.Truncate(time.Millisecond),
				Description:      "desc",
				ETag:             "etag",
				ExpirationTime:   aTime.Truncate(time.Millisecond),
				Name:             "fname",
				LastModifiedTime: aTime.Truncate(time.Millisecond),
				Location:         "loc",
				Labels:           map[string]string{"a": "b"},
			},
		},
	} {
		got, err := bqToModelMetadata(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if diff := testutil.Diff(got, test.want, cmpopts.IgnoreUnexported(ModelMetadata{})); diff != "" {
			t.Errorf("%+v:\n, -got, +want:\n%s", test.in, diff)
		}
	}
}

func TestModelMetadataUpdateToBQ(t *testing.T) {
	aTime := time.Date(2019, 3, 14, 0, 0, 0, 0, time.Local)
	aTimeMillis := aTime.UnixNano() / 1e6

	for _, test := range []struct {
		in   ModelMetadataToUpdate
		want *bq.Model
	}{
		{
			ModelMetadataToUpdate{},
			&bq.Model{},
		},
		{
			ModelMetadataToUpdate{
				Description: "d",
				Name:        "n",
			},
			&bq.Model{
				Description:     "d",
				FriendlyName:    "n",
				ForceSendFields: []string{"Description", "FriendlyName"},
			},
		},
		{
			ModelMetadataToUpdate{
				ExpirationTime: aTime,
			},
			&bq.Model{
				ExpirationTime:  aTimeMillis,
				ForceSendFields: []string{"ExpirationTime"},
			},
		},
		{
			ModelMetadataToUpdate{
				labelUpdater: labelUpdater{
					setLabels:    map[string]string{"L": "V"},
					deleteLabels: map[string]bool{"D": true},
				},
			},
			&bq.Model{
				Labels:     map[string]string{"L": "V"},
				NullFields: []string{"Labels.D"},
			},
		},
	} {
		got, err := test.in.toBQ()
		if err != nil {
			t.Fatalf("%+v: %v", test.in, err)
		}
		if diff := testutil.Diff(got, test.want); diff != "" {
			t.Errorf("%+v:\n-got, +want:\n%s", test.in, diff)
		}
	}
}
