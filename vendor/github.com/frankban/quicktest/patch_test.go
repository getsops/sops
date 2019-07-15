// Licensed under the MIT license, see LICENCE file for details.

package quicktest_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPatchSetInt(t *testing.T) {
	c := qt.New(t)
	i := 99
	c.Patch(&i, 88)
	c.Assert(i, qt.Equals, 88)
	c.Done()
	c.Assert(i, qt.Equals, 99)
}

func TestPatchSetError(t *testing.T) {
	c := qt.New(t)
	oldErr := errors.New("foo")
	newErr := errors.New("bar")
	err := oldErr
	c.Patch(&err, newErr)
	c.Assert(err, qt.Equals, newErr)
	c.Done()
	c.Assert(err, qt.Equals, oldErr)
}

func TestPatchSetErrorToNil(t *testing.T) {
	c := qt.New(t)
	oldErr := errors.New("foo")
	err := oldErr
	c.Patch(&err, nil)
	c.Assert(err, qt.Equals, nil)
	c.Done()
	c.Assert(err, qt.Equals, oldErr)
}

func TestPatchSetMapToNil(t *testing.T) {
	c := qt.New(t)
	oldMap := map[string]int{"foo": 1234}
	m := oldMap
	c.Patch(&m, nil)
	c.Assert(m, qt.IsNil)
	c.Done()
	c.Assert(m, qt.DeepEquals, oldMap)
}

func TestSetPatchPanicsWhenNotAssignable(t *testing.T) {
	c := qt.New(t)
	i := 99
	type otherInt int
	c.Assert(func() { c.Patch(&i, otherInt(88)) }, qt.PanicMatches, `reflect\.Set: value of type quicktest_test\.otherInt is not assignable to type int`)
}

func TestSetenv(t *testing.T) {
	c := qt.New(t)
	const envName = "SOME_VAR"
	os.Setenv(envName, "initial")
	c.Setenv(envName, "new value")
	c.Check(os.Getenv(envName), qt.Equals, "new value")
	c.Done()
	c.Check(os.Getenv(envName), qt.Equals, "initial")
}

func TestMkdir(t *testing.T) {
	c := qt.New(t)
	dir := c.Mkdir()
	c.Assert(c, qt.Not(qt.Equals), "")
	info, err := os.Stat(dir)
	c.Assert(err, qt.Equals, nil)
	c.Assert(info.IsDir(), qt.Equals, true)
	f, err := os.Create(filepath.Join(dir, "hello"))
	c.Assert(err, qt.Equals, nil)
	f.Close()
	c.Done()
	_, err = os.Stat(dir)
	c.Assert(err, qt.Not(qt.IsNil))
}
