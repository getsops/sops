package build

import (
	"fmt"
	gobuild "go/build"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/go/importgraphutil"
)

// Natives augment the standard library with GopherJS-specific changes.
// This test ensures that none of the standard library packages are modified
// in a way that adds imports which the original upstream standard library package
// does not already import. Doing that can increase generated output size or cause
// other unexpected issues (since the cmd/go tool does not know about these extra imports),
// so it's best to avoid it.
//
// It checks all standard library packages. Each package is considered as a normal
// package, as a test package, and as an external test package.
func TestNativesDontImportExtraPackages(t *testing.T) {
	// Calculate the forward import graph for all standard library packages.
	// It's needed for populateImportSet.
	stdOnly := gobuild.Default
	stdOnly.GOPATH = "" // We only care about standard library, so skip all GOPATH packages.
	forward, _, err := importgraphutil.BuildNoTests(&stdOnly)
	if err != nil {
		t.Fatalf("importgraphutil.BuildNoTests: %v", err)
	}

	// populateImportSet takes a slice of imports, and populates set with those
	// imports, as well as their transitive dependencies. That way, the set can
	// be quickly queried to check if a package is in the import graph of imports.
	//
	// Note, this does not include transitive imports of test/xtest packages,
	// which could cause some false positives. It currently doesn't, but if it does,
	// then support for that should be added here.
	populateImportSet := func(imports []string, set *stringSet) {
		for _, p := range imports {
			(*set)[p] = struct{}{}
			switch p {
			case "sync":
				(*set)["github.com/gopherjs/gopherjs/nosync"] = struct{}{}
			}
			transitiveImports := forward.Search(p)
			for p := range transitiveImports {
				(*set)[p] = struct{}{}
			}
		}
	}

	// Check all standard library packages.
	//
	// The general strategy is to first import each standard library package using the
	// normal build.Import, which returns a *build.Package. That contains Imports, TestImports,
	// and XTestImports values that are considered the "real imports".
	//
	// That list of direct imports is then expanded to the transitive closure by populateImportSet,
	// meaning all packages that are indirectly imported are also added to the set.
	//
	// Then, github.com/gopherjs/gopherjs/build.parseAndAugment(*build.Package) returns []*ast.File.
	// Those augmented parsed Go files of the package are checked, one file at at time, one import
	// at a time. Each import is verified to belong in the set of allowed real imports.
	for _, pkg := range gotool.ImportPaths([]string{"std"}) {
		// Normal package.
		{
			// Import the real normal package, and populate its real import set.
			bpkg, err := gobuild.Import(pkg, "", gobuild.ImportComment)
			if err != nil {
				t.Fatalf("gobuild.Import: %v", err)
			}
			realImports := make(stringSet)
			populateImportSet(bpkg.Imports, &realImports)

			// Use parseAndAugment to get a list of augmented AST files.
			fset := token.NewFileSet()
			files, err := parseAndAugment(bpkg, false, fset)
			if err != nil {
				t.Fatalf("github.com/gopherjs/gopherjs/build.parseAndAugment: %v", err)
			}

			// Verify imports of normal augmented AST files.
			for _, f := range files {
				fileName := fset.File(f.Pos()).Name()
				normalFile := !strings.HasSuffix(fileName, "_test.go")
				if !normalFile {
					continue
				}
				for _, imp := range f.Imports {
					importPath, err := strconv.Unquote(imp.Path.Value)
					if err != nil {
						t.Fatalf("strconv.Unquote(%v): %v", imp.Path.Value, err)
					}
					if importPath == "github.com/gopherjs/gopherjs/js" {
						continue
					}
					if _, ok := realImports[importPath]; !ok {
						t.Errorf("augmented normal package %q imports %q in file %v, but real %q doesn't:\nrealImports = %v", bpkg.ImportPath, importPath, fileName, bpkg.ImportPath, realImports)
					}
				}
			}
		}

		// Test package.
		{
			// Import the real test package, and populate its real import set.
			bpkg, err := gobuild.Import(pkg, "", gobuild.ImportComment)
			if err != nil {
				t.Fatalf("gobuild.Import: %v", err)
			}
			realTestImports := make(stringSet)
			populateImportSet(bpkg.TestImports, &realTestImports)

			// Use parseAndAugment to get a list of augmented AST files.
			fset := token.NewFileSet()
			files, err := parseAndAugment(bpkg, true, fset)
			if err != nil {
				t.Fatalf("github.com/gopherjs/gopherjs/build.parseAndAugment: %v", err)
			}

			// Verify imports of test augmented AST files.
			for _, f := range files {
				fileName, pkgName := fset.File(f.Pos()).Name(), f.Name.String()
				testFile := strings.HasSuffix(fileName, "_test.go") && !strings.HasSuffix(pkgName, "_test")
				if !testFile {
					continue
				}
				for _, imp := range f.Imports {
					importPath, err := strconv.Unquote(imp.Path.Value)
					if err != nil {
						t.Fatalf("strconv.Unquote(%v): %v", imp.Path.Value, err)
					}
					if importPath == "github.com/gopherjs/gopherjs/js" {
						continue
					}
					if _, ok := realTestImports[importPath]; !ok {
						t.Errorf("augmented test package %q imports %q in file %v, but real %q doesn't:\nrealTestImports = %v", bpkg.ImportPath, importPath, fileName, bpkg.ImportPath, realTestImports)
					}
				}
			}
		}

		// External test package.
		{
			// Import the real external test package, and populate its real import set.
			bpkg, err := gobuild.Import(pkg, "", gobuild.ImportComment)
			if err != nil {
				t.Fatalf("gobuild.Import: %v", err)
			}
			realXTestImports := make(stringSet)
			populateImportSet(bpkg.XTestImports, &realXTestImports)

			// Add _test suffix to import path to cause parseAndAugment to use external test mode.
			bpkg.ImportPath += "_test"

			// Use parseAndAugment to get a list of augmented AST files, then check only the external test files.
			fset := token.NewFileSet()
			files, err := parseAndAugment(bpkg, true, fset)
			if err != nil {
				t.Fatalf("github.com/gopherjs/gopherjs/build.parseAndAugment: %v", err)
			}

			// Verify imports of external test augmented AST files.
			for _, f := range files {
				fileName, pkgName := fset.File(f.Pos()).Name(), f.Name.String()
				xTestFile := strings.HasSuffix(fileName, "_test.go") && strings.HasSuffix(pkgName, "_test")
				if !xTestFile {
					continue
				}
				for _, imp := range f.Imports {
					importPath, err := strconv.Unquote(imp.Path.Value)
					if err != nil {
						t.Fatalf("strconv.Unquote(%v): %v", imp.Path.Value, err)
					}
					if importPath == "github.com/gopherjs/gopherjs/js" {
						continue
					}
					if _, ok := realXTestImports[importPath]; !ok {
						t.Errorf("augmented external test package %q imports %q in file %v, but real %q doesn't:\nrealXTestImports = %v", bpkg.ImportPath, importPath, fileName, bpkg.ImportPath, realXTestImports)
					}
				}
			}
		}
	}
}

// stringSet is used to print a set of strings in a more readable way.
type stringSet map[string]struct{}

func (m stringSet) String() string {
	s := make([]string, 0, len(m))
	for v := range m {
		s = append(s, v)
	}
	return fmt.Sprintf("%q", s)
}
