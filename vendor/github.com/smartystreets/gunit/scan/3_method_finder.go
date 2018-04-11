package scan

import (
	"go/ast"
	"strings"
)

type fixtureMethodFinder struct {
	fixtures map[string]*fixtureInfo
}

func newFixtureMethodFinder(fixtures map[string]*fixtureInfo) *fixtureMethodFinder {
	return &fixtureMethodFinder{fixtures: fixtures}
}

func (this *fixtureMethodFinder) Find(file *ast.File) map[string]*fixtureInfo {
	ast.Walk(this, file) // Calls this.Visit(...) recursively.
	return this.fixtures
}

func (this *fixtureMethodFinder) Visit(node ast.Node) ast.Visitor {
	function, isFunction := node.(*ast.FuncDecl)
	if !isFunction {
		return this
	}

	if function.Recv.NumFields() == 0 {
		return nil
	}

	receiver, isPointer := function.Recv.List[0].Type.(*ast.StarExpr)
	if !isPointer {
		return this
	}

	fixtureName := receiver.X.(*ast.Ident).Name
	fixture, functionMatchesFixture := this.fixtures[fixtureName]
	if !functionMatchesFixture {
		return nil
	}

	if !isExportedAndVoidAndNiladic(function) {
		return this
	}

	this.attach(function, fixture)
	return nil
}

func isExportedAndVoidAndNiladic(function *ast.FuncDecl) bool {
	if isExported := function.Name.IsExported(); !isExported {
		return false
	}
	if isNiladic := function.Type.Params.NumFields() == 0; !isNiladic {
		return false
	}
	isVoid := function.Type.Results.NumFields() == 0
	return isVoid
}

func (this *fixtureMethodFinder) attach(function *ast.FuncDecl, fixture *fixtureInfo) {
	if IsTestCase(function.Name.Name) {
		fixture.TestCases = append(fixture.TestCases, &testCaseInfo{
			CharacterPosition: int(function.Pos()),
			Name:              function.Name.Name,
		})
	}
}

func IsTestCase(name string) bool {
	return strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "LongTest") ||
		strings.HasPrefix(name, "FocusTest") ||
		strings.HasPrefix(name, "FocusLongTest") ||
		strings.HasPrefix(name, "SkipTest") ||
		strings.HasPrefix(name, "SkipLongTest")
}
