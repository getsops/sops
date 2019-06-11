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

package modinfo

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/delta"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/report"
)

var (
	verSuffixRegex = regexp.MustCompile(`v\d+$`)
)

// HasVersionSuffix returns true if the specified path has a version suffix in the form vN.
func HasVersionSuffix(path string) bool {
	return verSuffixRegex.MatchString(path)
}

// FindVersionSuffix returns the version suffix or the empty string.
func FindVersionSuffix(path string) string {
	return verSuffixRegex.FindString(path)
}

// Provider provides information about a module staged for release.
type Provider interface {
	DestDir() string
	NewExports() bool
	BreakingChanges() bool
	VersionSuffix() bool
	NewModule() bool
	GenerateReport() report.Package
}

type module struct {
	lhs  exports.Content
	rhs  exports.Content
	dest string
}

// GetModuleInfo collects information about a module staged for release.
// baseline is the directory for the current module
// staged is the directory for the module staged for release
func GetModuleInfo(baseline, staged string) (Provider, error) {
	// TODO: verify staged is a child of baseline
	lhs, err := exports.Get(baseline)
	if err != nil {
		// if baseline has no content then this is a v1 package
		if ei, ok := err.(exports.LoadPackageErrorInfo); !ok || ei.PkgCount() != 0 {
			return nil, fmt.Errorf("failed to get exports for package '%s': %s", baseline, err)
		}
	}
	rhs, err := exports.Get(staged)
	if err != nil {
		return nil, fmt.Errorf("failed to get exports for package '%s': %s", staged, err)
	}
	mod := module{
		lhs:  lhs,
		rhs:  rhs,
		dest: baseline,
	}
	// calculate the destination directory
	// if there are breaking changes calculate the new directory
	if mod.BreakingChanges() {
		dest := filepath.Dir(staged)
		v := 2
		if verSuffixRegex.MatchString(baseline) {
			// baseline has a version, get the number and increment it
			s := string(baseline[len(baseline)-1])
			v, err = strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("failed to convert '%s' to int: %v", s, err)
			}
			v++
		}
		mod.dest = filepath.Join(dest, fmt.Sprintf("v%d", v))
	}
	return mod, nil
}

// DestDir returns the fully qualified module destination directory.
func (m module) DestDir() string {
	return m.dest
}

// NewExports returns true if the module contains any additive changes.
func (m module) NewExports() bool {
	if m.lhs.IsEmpty() {
		return true
	}
	adds := delta.GetExports(m.lhs, m.rhs)
	return !adds.IsEmpty()
}

// BreakingChanges returns true if the module contains breaking changes.
func (m module) BreakingChanges() bool {
	if m.lhs.IsEmpty() {
		return false
	}
	// check for changed content
	if len(delta.GetConstTypeChanges(m.lhs, m.rhs)) > 0 ||
		len(delta.GetFuncSigChanges(m.lhs, m.rhs)) > 0 ||
		len(delta.GetInterfaceMethodSigChanges(m.lhs, m.rhs)) > 0 ||
		len(delta.GetStructFieldChanges(m.lhs, m.rhs)) > 0 {
		return true
	}
	// check for removed content
	if removed := delta.GetExports(m.rhs, m.lhs); !removed.IsEmpty() {
		return true
	}
	return false
}

// VersionSuffix returns true if the module path contains a version suffix.
func (m module) VersionSuffix() bool {
	return verSuffixRegex.MatchString(m.dest)
}

// NewModule returns true if the module is new, i.e. v1.0.0.
func (m module) NewModule() bool {
	return m.lhs.IsEmpty()
}

// GenerateReport generates a package report for the module.
func (m module) GenerateReport() report.Package {
	return report.Generate(m.lhs, m.rhs, false, false)
}
