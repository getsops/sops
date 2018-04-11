package scan

type fixtureInfo struct {
	Filename   string
	StructName string
	TestCases  []*testCaseInfo
}

type testCaseInfo struct {
	CharacterPosition int
	LineNumber        int
	Name              string
}
