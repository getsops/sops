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

package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	storage_v1_tests "cloud.google.com/go/storage/internal/test/conformance"
	"github.com/golang/protobuf/jsonpb"
)

func TestSigningV4Conformance(t *testing.T) {
	oldUTCNow := utcNow
	defer func() {
		utcNow = oldUTCNow
	}()

	dir := "internal/test/conformance"

	inBytes, err := ioutil.ReadFile(dir + "/service-account")
	if err != nil {
		t.Fatal(err)
	}
	serviceAccount := map[string]string{}
	if err := json.Unmarshal(inBytes, &serviceAccount); err != nil {
		t.Fatal(err)
	}
	googleAccessID := serviceAccount["client_email"]
	privateKey := serviceAccount["private_key"]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if !strings.Contains(f.Name(), ".json") {
			continue
		}

		inBytes, err := ioutil.ReadFile(dir + "/" + f.Name())
		if err != nil {
			t.Fatalf("%s: %v", f.Name(), err)
		}

		var testfile storage_v1_tests.TestFile
		if err := jsonpb.Unmarshal(bytes.NewReader(inBytes), &testfile); err != nil {
			t.Fatalf("unmarshalling %s: %v", f.Name(), err)
		}

		for _, tc := range testfile.SigningV4Tests {
			t.Run(tc.Description, func(t *testing.T) {
				utcNow = func() time.Time {
					return time.Unix(tc.Timestamp.Seconds, 0).UTC()
				}

				gotURL, err := SignedURL(tc.Bucket, tc.Object, &SignedURLOptions{
					GoogleAccessID: googleAccessID,
					PrivateKey:     []byte(privateKey),
					Method:         tc.Method,
					Expires:        utcNow().Add(time.Duration(tc.Expiration) * time.Second),
					Scheme:         SigningSchemeV4,
					Headers:        headersAsSlice(tc.Headers),
				})
				if err != nil {
					t.Fatal(err)
				}
				wantURL, err := url.Parse(tc.ExpectedUrl)
				if err != nil {
					t.Fatal(err)
				}
				// Sort the headers.
				wantURL.RawQuery = wantURL.Query().Encode()

				if gotURL != wantURL.String() {
					t.Fatalf("\nwant:\t%s\ngot:\t%s", wantURL.String(), gotURL)
				}
			})
		}
	}
}

func headersAsSlice(m map[string]string) []string {
	var s []string
	for k, v := range m {
		s = append(s, fmt.Sprintf("%s:%s", k, v))
	}
	return s
}
