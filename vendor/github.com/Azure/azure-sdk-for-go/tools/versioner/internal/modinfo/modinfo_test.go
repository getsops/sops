package modinfo

import (
	"regexp"
	"testing"
)

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

func Test_ScenarioA(t *testing.T) {
	// scenario A has no breaking changes, additive only
	mod, err := GetModuleInfo("../../testdata/scenarioa/foo", "../../testdata/scenarioa/foo/stage")
	if err != nil {
		t.Fatalf("failed to get module info: %v", err)
	}
	if mod.BreakingChanges() {
		t.Fatal("no breaking changes in scenario A")
	}
	if !mod.NewExports() {
		t.Fatal("expected new exports in scenario A")
	}
	if mod.VersionSuffix() {
		t.Fatalf("unexpected version suffix in scenario A")
	}
	regex := regexp.MustCompile(`testdata/scenarioa/foo$`)
	if !regex.MatchString(mod.DestDir()) {
		t.Fatalf("bad destination dir: %s", mod.DestDir())
	}
}

func Test_ScenarioB(t *testing.T) {
	// scenario B has a breaking change
	mod, err := GetModuleInfo("../../testdata/scenariob/foo", "../../testdata/scenariob/foo/stage")
	if err != nil {
		t.Fatalf("failed to get module info: %v", err)
	}
	if !mod.BreakingChanges() {
		t.Fatal("expected breaking changes in scenario B")
	}
	if !mod.NewExports() {
		t.Fatal("expected new exports in scenario B")
	}
	if !mod.VersionSuffix() {
		t.Fatalf("expected version suffix in scenario B")
	}
	regex := regexp.MustCompile(`testdata/scenariob/foo/v2$`)
	if !regex.MatchString(mod.DestDir()) {
		t.Fatalf("bad destination dir: %s", mod.DestDir())
	}
}

func Test_ScenarioC(t *testing.T) {
	// scenario C has no new exports or breaking changes (function body/doc changes only)
	mod, err := GetModuleInfo("../../testdata/scenarioc/foo", "../../testdata/scenarioc/foo/stage")
	if err != nil {
		t.Fatalf("failed to get module info: %v", err)
	}
	if mod.BreakingChanges() {
		t.Fatal("unexpected breaking changes in scenario C")
	}
	if mod.NewExports() {
		t.Fatal("unexpected new exports in scenario C")
	}
	if mod.VersionSuffix() {
		t.Fatalf("unexpected version suffix in scenario C")
	}
	regex := regexp.MustCompile(`testdata/scenarioc/foo$`)
	if !regex.MatchString(mod.DestDir()) {
		t.Fatalf("bad destination dir: %s", mod.DestDir())
	}
}

func Test_ScenarioD(t *testing.T) {
	// scenario D has a breaking change on top of a v2 release
	mod, err := GetModuleInfo("../../testdata/scenariod/foo/v2", "../../testdata/scenariod/foo/stage")
	if err != nil {
		t.Fatalf("failed to get module info: %v", err)
	}
	if !mod.BreakingChanges() {
		t.Fatal("expected breaking changes in scenario D")
	}
	if mod.NewExports() {
		t.Fatal("unexpected new exports in scenario D")
	}
	if !mod.VersionSuffix() {
		t.Fatalf("expected version suffix in scenario D")
	}
	regex := regexp.MustCompile(`testdata/scenariod/foo/v3$`)
	if !regex.MatchString(mod.DestDir()) {
		t.Fatalf("bad destination dir: %s", mod.DestDir())
	}
}

func Test_ScenarioE(t *testing.T) {
	// scenario E has a new export on top of a v2 release
	mod, err := GetModuleInfo("../../testdata/scenarioe/foo/v2", "../../testdata/scenarioe/foo/stage")
	if err != nil {
		t.Fatalf("failed to get module info: %v", err)
	}
	if mod.BreakingChanges() {
		t.Fatal("unexpected breaking changes in scenario E")
	}
	if !mod.NewExports() {
		t.Fatal("expected new exports in scenario E")
	}
	if !mod.VersionSuffix() {
		t.Fatalf("expected version suffix in scenario E")
	}
	regex := regexp.MustCompile(`testdata/scenarioe/foo/v2$`)
	if !regex.MatchString(mod.DestDir()) {
		t.Fatalf("bad destination dir: %s", mod.DestDir())
	}
}

func Test_ScenarioF(t *testing.T) {
	// scenario F is a new module
	mod, err := GetModuleInfo("../../testdata/scenariof/foo", "../../testdata/scenariof/foo/stage")
	if err != nil {
		t.Fatalf("failed to get module info: %v", err)
	}
	if mod.BreakingChanges() {
		t.Fatal("unexpected breaking changes in scenario F")
	}
	if !mod.NewExports() {
		t.Fatal("expected new exports in scenario F")
	}
	if mod.VersionSuffix() {
		t.Fatalf("unexpected version suffix in scenario F")
	}
	if !mod.NewModule() {
		t.Fatal("expected new module in scenario F")
	}
	regex := regexp.MustCompile(`testdata/scenariof/foo$`)
	if !regex.MatchString(mod.DestDir()) {
		t.Fatalf("bad destination dir: %s", mod.DestDir())
	}
}
