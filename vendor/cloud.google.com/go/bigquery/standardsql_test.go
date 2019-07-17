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

	"cloud.google.com/go/internal/testutil"
	bq "google.golang.org/api/bigquery/v2"
)

func TestBQToStandardSQLDataType(t *testing.T) {
	for _, test := range []struct {
		in   *bq.StandardSqlDataType
		want *StandardSQLDataType
	}{
		{&bq.StandardSqlDataType{}, &StandardSQLDataType{}},
		{
			&bq.StandardSqlDataType{
				TypeKind: "INT64",
			},
			&StandardSQLDataType{
				TypeKind: "INT64",
			},
		},
		{
			&bq.StandardSqlDataType{
				TypeKind: "ARRAY",
				ArrayElementType: &bq.StandardSqlDataType{
					TypeKind: "INT64",
				},
			},
			&StandardSQLDataType{
				TypeKind: "ARRAY",
				ArrayElementType: &StandardSQLDataType{
					TypeKind: "INT64",
				},
			},
		},
	} {
		got, err := bqToStandardSQLDataType(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if diff := testutil.Diff(got, test.want); diff != "" {
			t.Errorf("%+v: -got, +want:\n%s", test.in, diff)
		}
	}
}

func TestBQToStandardSQLField(t *testing.T) {
	for _, test := range []struct {
		in   *bq.StandardSqlField
		want *StandardSQLField
	}{
		{&bq.StandardSqlField{}, &StandardSQLField{}},
		{
			&bq.StandardSqlField{
				Name: "foo",
				Type: &bq.StandardSqlDataType{
					TypeKind: "NUMERIC",
				},
			},
			&StandardSQLField{
				Name: "foo",
				Type: &StandardSQLDataType{
					TypeKind: "NUMERIC",
				},
			},
		},
	} {
		got, err := bqToStandardSQLField(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if diff := testutil.Diff(got, test.want); diff != "" {
			t.Errorf("%+v: -got, +want:\n%s", test.in, diff)
		}
	}
}

func TestBQToStandardSQLStructType(t *testing.T) {
	for _, test := range []struct {
		in   *bq.StandardSqlStructType
		want *StandardSQLStructType
	}{
		{&bq.StandardSqlStructType{}, &StandardSQLStructType{}},
		{
			&bq.StandardSqlStructType{
				Fields: []*bq.StandardSqlField{
					{
						Name: "foo",
						Type: &bq.StandardSqlDataType{
							TypeKind: "STRING",
						},
					},
				},
			},
			&StandardSQLStructType{
				Fields: []*StandardSQLField{
					{
						Name: "foo",
						Type: &StandardSQLDataType{
							TypeKind: "STRING",
						},
					},
				},
			},
		},
	} {
		got, err := bqToStandardSQLStructType(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if diff := testutil.Diff(got, test.want); diff != "" {
			t.Errorf("%+v: -got, +want:\n%s", test.in, diff)
		}
	}
}
