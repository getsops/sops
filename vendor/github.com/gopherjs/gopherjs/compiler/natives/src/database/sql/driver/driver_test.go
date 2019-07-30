// +build js

package driver

var valueConverterTests = []valueConverterTest{
	{Bool, "true", true, ""},
	{Bool, "True", true, ""},
	{Bool, []byte("t"), true, ""},
	{Bool, true, true, ""},
	{Bool, "1", true, ""},
	{Bool, 1, true, ""},
	{Bool, int64(1), true, ""},
	{Bool, uint16(1), true, ""},
	{Bool, "false", false, ""},
	{Bool, false, false, ""},
	{Bool, "0", false, ""},
	{Bool, 0, false, ""},
	{Bool, int64(0), false, ""},
	{Bool, uint16(0), false, ""},
	{c: Bool, in: "foo", err: "sql/driver: couldn't convert \"foo\" into type bool"},
	{c: Bool, in: 2, err: "sql/driver: couldn't convert 2 into type bool"},
	{DefaultParameterConverter, now, now, ""},
	{DefaultParameterConverter, (*int64)(nil), nil, ""},
	{DefaultParameterConverter, &answer, answer, ""},
	{DefaultParameterConverter, &now, now, ""},
	//{DefaultParameterConverter, i(9), int64(9), ""}, // TODO: Fix.
	{DefaultParameterConverter, f(0.1), float64(0.1), ""},
	{DefaultParameterConverter, b(true), true, ""},
	//{DefaultParameterConverter, bs{1}, []byte{1}, ""}, // TODO: Fix.
	{DefaultParameterConverter, s("a"), "a", ""},
	{DefaultParameterConverter, is{1}, nil, "unsupported type driver.is, a slice of int"},
}
