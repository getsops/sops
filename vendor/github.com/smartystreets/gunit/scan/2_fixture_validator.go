package scan

import "go/ast"

type fixtureValidator struct {
	Parent      *fixtureCollector
	FixtureName string
}

func (this *fixtureValidator) Visit(node ast.Node) ast.Visitor {
	// We start at a TypeSpec and look for an embedded pointer field: `*gunit.Fixture`.
	field, isField := node.(*ast.Field)
	if !isField {
		return this
	}
	pointer, isPointer := field.Type.(*ast.StarExpr)
	if !isPointer {
		return this
	}

	selector, isSelector := pointer.X.(*ast.SelectorExpr)
	if !isSelector {
		return this
	}
	gunit, isGunit := selector.X.(*ast.Ident)
	if selector.Sel.Name != "Fixture" || !isGunit || gunit.Name != "gunit" {
		return this
	}
	this.Parent.Validate(this.FixtureName)
	return nil
}
