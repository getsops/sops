// +build go1.9

// Copyright 2018 Microsoft Corporation
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

package model

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// DeleteChildDirs deletes all child directories in the specified directory.
func DeleteChildDirs(dir string) error {
	children, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, child := range children {
		if child.IsDir() {
			childPath := filepath.Join(dir, child.Name())
			err = os.RemoveAll(childPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
