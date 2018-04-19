package scan

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type TestCasePositions map[string]string

func LocateTestCases(filename string) TestCasePositions {
	return gatherTestCaseLineNumbers(parseFixtures(filename))
}

func gatherTestCaseLineNumbers(fixtures []*fixtureInfo) TestCasePositions {
	positions := make(TestCasePositions)
	for _, fixture := range fixtures {
		for _, test := range fixture.TestCases {
			key := fmt.Sprintf("Test%s/%s", fixture.StructName, test.Name)
			value := fmt.Sprintf("%s:%d", fixture.Filename, test.LineNumber)
			positions[key] = value
		}
	}
	return positions
}

func parseFixtures(filename string) (fixtures []*fixtureInfo) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	batch, err := scanForFixtures(string(source))
	if err != nil {
		return nil
	}
	for _, fixture := range batch {
		fixture.Filename = filename
		fixtures = append(fixtures, fixture)
		for _, test := range fixture.TestCases {
			test.LineNumber = lineNumber(string(source), test.CharacterPosition)
		}
	}
	return fixtures
}

func lineNumber(source string, position int) int {
	source = source[:position+1]
	return strings.Count(source, "\n") + 1
}
