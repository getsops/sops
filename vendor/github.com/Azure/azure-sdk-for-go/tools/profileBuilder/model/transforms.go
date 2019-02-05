// +build go1.9

// Copyright 2018 Microsoft Corporation and contributors
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

// Package model holds the business logic for the operations made available by
// profileBuilder.
//
// This package is not governed by the SemVer associated with the rest of the
// Azure-SDK-for-Go.
package model

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/tools/imports"
)

// ListDefinition represents a JSON file that contains a list of packages to include
type ListDefinition struct {
	Include      []string          `json:"include"`
	PathOverride map[string]string `json:"pathOverride"`
}

const (
	armPathModifier = "mgmt"
	aliasFileName   = "models.go"
)

var packageName = regexp.MustCompile(`services[/\\](?P<provider>[\w\-\.\d_\\/]+)[/\\](?:(?P<arm>` + armPathModifier + `)[/\\])?(?P<version>v?\d{4}-\d{2}-\d{2}[\w\d\.\-]*|v?\d+[\.\d+\.\d\w\-]*)[/\\](?P<group>[/\\\w\d\-\._]+)`)

// BuildProfile takes a list of packages and creates a profile
func BuildProfile(packageList ListDefinition, name, outputLocation string, outputLog, errLog *log.Logger) {
	wg := &sync.WaitGroup{}
	wg.Add(len(packageList.Include))
	for _, pkgDir := range packageList.Include {
		if !filepath.IsAbs(pkgDir) {
			abs, err := filepath.Abs(pkgDir)
			if err != nil {
				errLog.Fatalf("failed to convert to absolute path: %v", err)
			}
			pkgDir = abs
		}
		go func(pd string) {
			fs := token.NewFileSet()
			packages, err := parser.ParseDir(fs, pd, func(f os.FileInfo) bool {
				// exclude test files
				return !strings.HasSuffix(f.Name(), "_test.go")
			}, 0)
			if err != nil {
				errLog.Fatalf("failed to parse '%s': %v", pd, err)
			}
			if len(packages) < 1 {
				errLog.Fatalf("didn't find any packages in '%s'", pd)
			}
			if len(packages) > 1 {
				errLog.Fatalf("found more than one package in '%s'", pd)
			}
			for pn := range packages {
				p := packages[pn]
				// trim any non-exported nodes
				if exp := ast.PackageExports(p); !exp {
					errLog.Fatalf("package '%s' doesn't contain any exports", pn)
				}
				// construct the import path from the outputLocation
				// e.g. D:\work\src\github.com\Azure\azure-sdk-for-go\profiles\2017-03-09\compute\mgmt\compute
				// becomes github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/compute/mgmt/compute
				i := strings.Index(pd, "github.com")
				if i == -1 {
					errLog.Fatalf("didn't find 'github.com' in '%s'", pd)
				}
				importPath := strings.Replace(pd[i:], "\\", "/", -1)
				ap, err := NewAliasPackage(p, importPath)
				if err != nil {
					errLog.Fatalf("failed to create alias package: %v", err)
				}
				updateAliasPackageUserAgent(ap, name)
				// build the profile output directory, if there's an override path use that
				var aliasPath string
				var ok bool
				if aliasPath, ok = packageList.PathOverride[importPath]; !ok {
					var err error
					aliasPath, err = getAliasPath(pd)
					if err != nil {
						errLog.Fatalf("failed to calculate alias directory: %v", err)
					}
				}
				aliasPath = filepath.Join(outputLocation, aliasPath)
				if _, err := os.Stat(aliasPath); os.IsNotExist(err) {
					err = os.MkdirAll(aliasPath, os.ModeDir)
					if err != nil {
						errLog.Fatalf("failed to create alias directory: %v", err)
					}
				}
				writeAliasPackage(ap, aliasPath, outputLog, errLog)
			}
			wg.Done()
		}(pkgDir)
	}
	wg.Wait()
	outputLog.Print(len(packageList.Include), " packages generated.")
}

// getAliasPath takes an existing API Version path and converts the path to a path which uses the new profile layout.
func getAliasPath(packageDir string) (string, error) {
	// we want to transform this:
	//  .../services/compute/mgmt/2016-03-30/compute
	// into this:
	//  compute/mgmt/compute
	// i.e. remove everything to the left of /services along with the API version
	packageDir = strings.TrimSuffix(packageDir, string(filepath.Separator))
	matches := packageName.FindAllStringSubmatch(packageDir, -1)
	if matches == nil {
		return "", fmt.Errorf("path '%s' does not resemble a known package path", packageDir)
	}

	output := []string{
		matches[0][1],
	}

	if matches[0][2] == armPathModifier {
		output = append(output, armPathModifier)
	}
	output = append(output, matches[0][4])

	return filepath.Join(output...), nil
}

// updateAliasPackageUserAgent updates the "UserAgent" function in the generated profile, if it is present.
func updateAliasPackageUserAgent(ap *AliasPackage, profileName string) {
	var userAgent *ast.FuncDecl
	for _, decl := range ap.Files[aliasFileName].Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == "UserAgent" {
			userAgent = fd
			break
		}
	}
	if userAgent == nil {
		return
	}

	// Grab the expression being returned.
	retResults := &userAgent.Body.List[0].(*ast.ReturnStmt).Results[0]

	// Append a string literal to the result
	updated := &ast.BinaryExpr{
		Op: token.ADD,
		X:  *retResults,
		Y: &ast.BasicLit{
			Value: fmt.Sprintf(`" profiles/%s"`, profileName),
		},
	}
	*retResults = updated
}

// writeAliasPackage adds the MSFT Copyright Header, then writes the alias package to disk.
func writeAliasPackage(ap *AliasPackage, outputPath string, outputLog, errLog *log.Logger) {
	files := token.NewFileSet()

	err := os.MkdirAll(path.Dir(outputPath), os.ModePerm|os.ModeDir)
	if err != nil {
		errLog.Fatalf("error creating directory: %v", err)
	}

	aliasFile := filepath.Join(outputPath, aliasFileName)
	outputFile, err := os.Create(aliasFile)
	if err != nil {
		errLog.Fatalf("error creating file: %v", err)
	}

	// TODO: This should really be added by the `goalias` package itself. Doing it here is a work around
	fmt.Fprintln(outputFile, "// +build go1.9")
	fmt.Fprintln(outputFile)

	generatorStampBuilder := new(bytes.Buffer)

	fmt.Fprintf(generatorStampBuilder, "// Copyright %4d Microsoft Corporation\n", time.Now().Year())
	fmt.Fprintln(generatorStampBuilder, `//
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
// limitations under the License.`)

	fmt.Fprintln(outputFile, generatorStampBuilder.String())

	generatorStampBuilder.Reset()

	fmt.Fprintln(generatorStampBuilder, "// This code was auto-generated by:")
	fmt.Fprintln(generatorStampBuilder, "// github.com/Azure/azure-sdk-for-go/tools/profileBuilder")

	fmt.Fprintln(generatorStampBuilder)
	fmt.Fprint(outputFile, generatorStampBuilder.String())

	outputLog.Printf("Writing File: %s", aliasFile)

	file := ap.ModelFile()

	var b bytes.Buffer
	printer.Fprint(&b, files, file)
	res, err := imports.Process(aliasFile, b.Bytes(), nil)
	if err != nil {
		errLog.Fatalf("failed to process imports: %v", err)
	}
	fmt.Fprintf(outputFile, "%s", res)
	outputFile.Close()

	// be sure to specify the file for formatting not the directory; this is to
	// avoid race conditions when formatting parent/child directories (foo and foo/fooapi)
	if err := exec.Command("gofmt", "-w", aliasFile).Run(); err != nil {
		errLog.Fatalf("error formatting profile '%s': %v", aliasFile, err)
	}
}
