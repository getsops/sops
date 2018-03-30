package scan

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

//////////////////////////////////////////////////////////////////////////////

func scanForFixtures(code string) ([]*fixtureInfo, error) {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", code, 0)
	if err != nil {
		return nil, err
	}
	// ast.Print(fileset, file) // helps with debugging...
	return findAndListFixtures(file)
}

func findAndListFixtures(file *ast.File) ([]*fixtureInfo, error) {
	collection := newFixtureCollector().Collect(file)
	collection = newFixtureMethodFinder(collection).Find(file)
	return listFixtures(collection)
}

func listFixtures(collection map[string]*fixtureInfo) ([]*fixtureInfo, error) {
	var fixtures []*fixtureInfo
	errorMessage := new(bytes.Buffer)

	for _, fixture := range collection {
		fixtures = append(fixtures, fixture)
	}
	if errorMessage.Len() > 0 {
		return nil, errors.New(errorMessage.String())
	}
	return fixtures, nil
}
