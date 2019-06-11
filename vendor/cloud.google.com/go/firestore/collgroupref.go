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

package firestore

// A CollectionGroupRef is a reference to a group of collections sharing the
// same ID.
type CollectionGroupRef struct {
	c *Client

	// Use the methods of Query on a CollectionGroupRef to create and run queries.
	Query
}

func newCollectionGroupRef(c *Client, dbPath, collectionID string) *CollectionGroupRef {
	return &CollectionGroupRef{
		c: c,

		Query: Query{
			c:              c,
			collectionID:   collectionID,
			path:           dbPath,
			parentPath:     dbPath + "/documents",
			allDescendants: true,
		},
	}
}
