package scan

import "go/ast"

type fixtureCollector struct {
	candidates map[string]*fixtureInfo
	fixtures   map[string]*fixtureInfo
}

func newFixtureCollector() *fixtureCollector {
	return &fixtureCollector{
		candidates: make(map[string]*fixtureInfo),
		fixtures:   make(map[string]*fixtureInfo),
	}
}

func (this *fixtureCollector) Collect(file *ast.File) map[string]*fixtureInfo {
	ast.Walk(this, file) // Calls this.Visit(...) recursively which populates this.fixtures
	return this.fixtures
}

func (this *fixtureCollector) Visit(node ast.Node) ast.Visitor {
	switch t := node.(type) {
	case *ast.TypeSpec:
		name := t.Name.Name
		this.candidates[name] = &fixtureInfo{StructName: name}
		return &fixtureValidator{Parent: this, FixtureName: name}
	default:
		return this
	}
}

func (this *fixtureCollector) Validate(fixture string) {
	this.fixtures[fixture] = this.candidates[fixture]
	delete(this.candidates, fixture)
}

func (this *fixtureCollector) Invalidate(fixture string) {
	this.Validate(fixture)
}
