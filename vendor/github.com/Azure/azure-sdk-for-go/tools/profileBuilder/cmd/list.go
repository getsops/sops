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

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/profileBuilder/model"
	"github.com/spf13/cobra"
)

const (
	inputLongName    = "input"
	inputShortName   = "i"
	inputDescription = "Specify the input JSON file to read for the list of packages."
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Creates a profile from a set of packages.",
	Long: `Reads a list of packages from stdin, where each line is treated as a Go package
identifier. These packages are then used to create a profile.

Often, the easiest way of invoking this command will be using a pipe operator
to specify the packages to include.

Example:
$> ../model/testdata/smallProfile.txt > profileBuilder list --name small_profile
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logWriter := ioutil.Discard
		if verboseFlag {
			logWriter = os.Stdout
		}

		outputLog := log.New(logWriter, "[STATUS] ", 0)
		errLog := log.New(os.Stderr, "[ERROR] ", 0)

		if !filepath.IsAbs(outputRootDir) {
			abs, err := filepath.Abs(outputRootDir)
			if err != nil {
				errLog.Fatalf("failed to convert to absolute path: %v", err)
			}
			outputRootDir = abs
		}
		outputLog.Printf("Output-Location set to: %s", outputRootDir)

		inputFile, err := cmd.Flags().GetString(inputLongName)
		if err != nil {
			errLog.Fatalf("failed to get %s: %v", inputLongName, err)
		}

		data, err := ioutil.ReadFile(inputFile)
		if err != nil {
			errLog.Fatalf("failed to read list: %v", err)
		}

		var listDef model.ListDefinition
		err = json.Unmarshal(data, &listDef)
		if err != nil {
			errLog.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if clearOutputFlag {
			if err := model.DeleteChildDirs(outputRootDir); err != nil {
				errLog.Fatalf("Unable to clear output-folder: %v", err)
			}
		}

		model.BuildProfile(listDef, profileName, outputRootDir, outputLog, errLog)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP(inputLongName, inputShortName, "", inputDescription)
	listCmd.MarkFlagRequired(inputLongName)
}
