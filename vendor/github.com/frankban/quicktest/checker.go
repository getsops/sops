// Licensed under the MIT license, see LICENCE file for details.

package quicktest

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kr/pretty"
)

// Checker is implemented by types used as part of Check/Assert invocations.
type Checker interface {
	// Check checks that the obtained value (got) is correct with respect to
	// the checker's arguments (args). On failure, the returned error is
	// printed along with the checker arguments and any key-value pairs added
	// by calling the note function. Values are pretty-printed unless they are
	// of type Unquoted.
	//
	// When the check arguments are invalid, Check may return a BadCheck error,
	// which suppresses printing of the checker arguments. Values added with
	// note are still printed.
	//
	// If Check returns ErrSilent, neither the checker arguments nor the error
	// are printed. Again, values added with note are still printed.
	Check(got interface{}, args []interface{}, note func(key string, value interface{})) error

	// ArgNames returns the names of all required arguments, including the
	// mandatory got argument and any additional args.
	ArgNames() []string
}

// Equals is a Checker checking equality of two comparable values.
// For instance:
//
//     c.Assert(answer, qt.Equals, 42)
//
// Note that the following will fail:
//
//     c.Assert((*sometype)(nil), qt.Equals, nil)
//
// Use the IsNil checker below for this kind of nil check.
var Equals Checker = &equalsChecker{
	argNames: []string{"got", "want"},
}

type equalsChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got == args[0].
func (c *equalsChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	defer func() {
		// A panic is raised when the provided values are not comparable.
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	if want := args[0]; got != want {
		if _, ok := got.(error); ok && want == nil {
			return errors.New("got non-nil error")
		}
		return errors.New("values are not equal")
	}
	return nil
}

// CmpEquals returns a Checker checking equality of two arbitrary values
// according to the provided compare options. See DeepEquals as an example of
// such a checker, commonly used when no compare options are required.
// For instance:
//
//     c.Assert(list, qt.CmpEquals(cmpopts.SortSlices), []int{42, 47})
//     c.Assert(got, qt.CmpEquals(), []int{42, 47}) // Same as qt.DeepEquals.
//
func CmpEquals(opts ...cmp.Option) Checker {
	return cmpEquals(testing.Verbose, opts...)
}

func cmpEquals(verbose func() bool, opts ...cmp.Option) Checker {
	return &cmpEqualsChecker{
		argNames: []string{"got", "want"},
		opts:     opts,
		verbose:  verbose,
	}
}

type cmpEqualsChecker struct {
	argNames
	opts    cmp.Options
	verbose func() bool
}

// Check implements Checker.Check by checking that got == args[0] according to
// the compare options stored in the checker.
func (c *cmpEqualsChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	defer func() {
		// A panic is raised in some cases, for instance when trying to compare
		// structs with unexported fields and neither AllowUnexported nor
		// cmpopts.IgnoreUnexported are provided.
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	want := args[0]
	if diff := cmp.Diff(got, want, c.opts...); diff != "" {
		// Only output values when the verbose flag is set.
		if c.verbose() {
			note("diff (-got +want)", Unquoted(diff))
			return errors.New("values are not deep equal")
		}
		note("error", Unquoted("values are not deep equal"))
		note("diff (-got +want)", Unquoted(diff))
		return ErrSilent
	}
	return nil
}

// DeepEquals is a Checker deeply checking equality of two arbitrary values.
// For instance:
//
//     c.Assert(got, qt.DeepEquals, []int{42, 47})
//
var DeepEquals = CmpEquals()

// ContentEquals is like DeepEquals but any slices in the compared values will
// be sorted before being compared.
var ContentEquals = CmpEquals(cmpopts.SortSlices(func(x, y interface{}) bool {
	// TODO frankban: implement a proper sort function.
	return pretty.Sprint(x) < pretty.Sprint(y)
}))

// Matches is a Checker checking that the provided string or fmt.Stringer
// matches the provided regular expression pattern.
// For instance:
//
//     c.Assert("these are the voyages", qt.Matches, "these are .*")
//     c.Assert(net.ParseIP("1.2.3.4"), qt.Matches, "1.*")
//
var Matches Checker = &matchesChecker{
	argNames: []string{"got value", "regexp"},
}

type matchesChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got is a string or a
// fmt.Stringer and that it matches args[0].
func (c *matchesChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) error {
	pattern := args[0]
	switch v := got.(type) {
	case string:
		return match(v, pattern, "value does not match regexp", note)
	case fmt.Stringer:
		return match(v.String(), pattern, "value.String() does not match regexp", note)
	}
	note("value", got)
	return BadCheckf("value is not a string or a fmt.Stringer")
}

// ErrorMatches is a Checker checking that the provided value is an error whose
// message matches the provided regular expression pattern.
// For instance:
//
//     c.Assert(err, qt.ErrorMatches, "bad wolf .*")
//
var ErrorMatches Checker = &errorMatchesChecker{
	argNames: []string{"got error", "regexp"},
}

type errorMatchesChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got is an error whose
// Error() matches args[0].
func (c *errorMatchesChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) error {
	if got == nil {
		return errors.New("got nil error but want non-nil")
	}
	err, ok := got.(error)
	if !ok {
		note("got", got)
		return BadCheckf("first argument is not an error")
	}
	return match(err.Error(), args[0], "error does not match regexp", note)
}

// PanicMatches is a Checker checking that the provided function panics with a
// message matching the provided regular expression pattern.
// For instance:
//
//     c.Assert(func() {panic("bad wolf ...")}, qt.PanicMatches, "bad wolf .*")
//
var PanicMatches Checker = &panicMatchesChecker{
	argNames: []string{"function", "regexp"},
}

type panicMatchesChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got is a func() that panics
// with a message matching args[0].
func (c *panicMatchesChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	f := reflect.ValueOf(got)
	if f.Kind() != reflect.Func {
		note("got", got)
		return BadCheckf("first argument is not a function")
	}
	ftype := f.Type()
	if ftype.NumIn() != 0 {
		note("function", got)
		return BadCheckf("cannot use a function receiving arguments")
	}

	defer func() {
		r := recover()
		if r == nil {
			err = errors.New("function did not panic")
			return
		}
		msg := fmt.Sprint(r)
		note("panic value", msg)
		err = match(msg, args[0], "panic value does not match regexp", note)
	}()

	f.Call(nil)
	return nil
}

// IsNil is a Checker checking that the provided value is nil.
// For instance:
//
//     c.Assert(got, qt.IsNil)
//
var IsNil Checker = &isNilChecker{
	argNames: []string{"got"},
}

type isNilChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got is nil.
func (c *isNilChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	if got == nil {
		return nil
	}
	value := reflect.ValueOf(got)
	if canBeNil(value.Kind()) && value.IsNil() {
		return nil
	}
	return fmt.Errorf("%#v is not nil", got)
}

// HasLen is a Checker checking that the provided value has the given length.
// For instance:
//
//     c.Assert([]int{42, 47}, qt.HasLen, 2)
//     c.Assert(myMap, qt.HasLen, 42)
//
var HasLen Checker = &hasLenChecker{
	argNames: []string{"got", "want length"},
}

type hasLenChecker struct {
	argNames
}

// Check implements Checker.Check by checking that len(got) == args[0].
func (c *hasLenChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	v := reflect.ValueOf(got)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
	default:
		note("got", got)
		return BadCheckf("first argument has no length")
	}
	want, ok := args[0].(int)
	if !ok {
		note("length", args[0])
		return BadCheckf("length is not an int")
	}
	length := v.Len()
	note("len(got)", length)
	if length != want {
		return fmt.Errorf("unexpected length")
	}
	return nil
}

// Satisfies is a Checker checking that the provided value, when used as
// argument of the provided predicate function, causes the function to return
// true. The function must be of type func(T) bool, having got assignable to T.
// For instance:
//
//     // Check that an error from os.Open satisfies os.IsNotExist.
//     c.Assert(err, qt.Satisfies, os.IsNotExist)
//
//     // Check that a floating point number is a not-a-number.
//     c.Assert(f, qt.Satisfies, math.IsNaN)
//
var Satisfies Checker = &satisfiesChecker{
	argNames: []string{"arg", "predicate function"},
}

type satisfiesChecker struct {
	argNames
}

// Check implements Checker.Check by checking that args[0](got) == true.
func (c *satisfiesChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	// Original code at
	// <https://github.com/juju/testing/blob/master/checkers/bool.go>.
	// Copyright 2011 Canonical Ltd.
	// Licensed under the LGPLv3, see LICENCE file for details.
	predicate := args[0]
	f := reflect.ValueOf(predicate)
	ftype := f.Type()
	if ftype.Kind() != reflect.Func || ftype.NumIn() != 1 || ftype.NumOut() != 1 || ftype.Out(0).Kind() != reflect.Bool {
		note("predicate function", predicate)
		return BadCheckf("predicate function is not a func(T) bool")
	}
	v, t := reflect.ValueOf(got), ftype.In(0)
	if !v.IsValid() {
		if !canBeNil(t.Kind()) {
			note("predicate function", predicate)
			return BadCheckf("cannot use nil as type %v in argument to predicate function", t)
		}
		v = reflect.Zero(t)
	} else if !v.Type().AssignableTo(t) {
		note("arg", got)
		note("predicate function", predicate)
		return BadCheckf("cannot use value of type %v as type %v in argument to predicate function", v.Type(), t)
	}
	if f.Call([]reflect.Value{v})[0].Interface().(bool) {
		return nil
	}
	return fmt.Errorf("value does not satisfy predicate function")
}

// Not returns a Checker negating the given Checker.
// For instance:
//
//     c.Assert(got, qt.Not(qt.IsNil))
//     c.Assert(answer, qt.Not(qt.Equals), 42)
//
func Not(checker Checker) Checker {
	return &notChecker{
		Checker: checker,
	}
}

type notChecker struct {
	Checker
}

// Check implements Checker.Check by checking that the stored checker fails.
func (c *notChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	if nc, ok := c.Checker.(*notChecker); ok {
		return nc.Checker.Check(got, args, note)
	}
	err = c.Checker.Check(got, args, note)
	if IsBadCheck(err) {
		return err
	}
	if err != nil {
		return nil
	}
	return errors.New("unexpected success")
}

// Contains is a checker that checks that a map, slice, array
// or string contains a value. It's the same as using
// Any(Equals), except that it has a special case
// for strings - if the first argument is a string,
// the second argument must also be a string
// and strings.Contains will be used.
//
// For example:
//
//	c.Assert("hello world", qt.Contains, "world")
//	c.Assert([]int{3,5,7,99}, qt.Contains, 7)
var Contains Checker = &containsChecker{
	argNames: []string{"got", "want"},
}

type containsChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got contains args[0].
func (c *containsChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) error {
	if got, ok := got.(string); ok {
		want, ok := args[0].(string)
		if !ok {
			return BadCheckf("strings can only contain strings, not %T", args[0])
		}
		if strings.Contains(got, want) {
			return nil
		}
		return errors.New("no substring match found")
	}
	return Any(Equals).Check(got, args, note)
}

// Any returns a Checker that uses the given checker to check elements
// of a slice or array or the values from a map. It succeeds if any element
// passes the check.
//
// For example:
//
//	c.Assert([]int{3,5,7,99}, qt.Any(qt.Equals), 7)
//	c.Assert([][]string{{"a", "b"}, {"c", "d"}}, qt.Any(qt.DeepEquals), []string{"c", "d"})
//
// See also All and Contains.
func Any(c Checker) Checker {
	return &anyChecker{
		argNames:    append([]string{"container"}, c.ArgNames()[1:]...),
		elemChecker: c,
	}
}

type anyChecker struct {
	argNames
	elemChecker Checker
}

// Check implements Checker.Check by checking that one of the elements of
// got passes the c.elemChecker check.
func (c *anyChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) error {
	iter, err := newIter(got)
	if err != nil {
		return BadCheckf("%v", err)
	}
	for iter.next() {
		// For the time being, discard the notes added by the sub-checker,
		// because it's not clear what a good behaviour would be.
		// Should we print all the failed check for all elements? If there's only
		// one element in the container, the answer is probably yes,
		// but let's leave it for now.
		err := c.elemChecker.Check(
			iter.value().Interface(),
			args,
			func(key string, value interface{}) {},
		)
		if err == nil {
			return nil
		}
		if IsBadCheck(err) {
			return BadCheckf("at %s: %v", iter.key(), err)
		}
	}
	return errors.New("no matching element found")
}

// All returns a Checker that uses the given checker to check elements
// of slice or array or the values of a map. It succeeds if all elements
// pass the check.
// On failure it prints the error from the first index that failed.
//
// For example:
//
//	c.Assert([]int{3, 5, 8}, qt.All(qt.Not(qt.Equals)), 0)
//	c.Assert([][]string{{"a", "b"}, {"a", "b"}}, qt.All(qt.DeepEquals), []string{"c", "d"})
//
// See also Any and Contains.
func All(c Checker) Checker {
	return &allChecker{
		argNames:    append([]string{"container"}, c.ArgNames()[1:]...),
		elemChecker: c,
	}
}

type allChecker struct {
	argNames
	elemChecker Checker
}

// Check implement Checker.Check by checking that all the elements of got
// pass the c.elemChecker check.
func (c *allChecker) Check(got interface{}, args []interface{}, notef func(key string, value interface{})) error {
	iter, err := newIter(got)
	if err != nil {
		return BadCheckf("%v", err)
	}
	for iter.next() {
		// Store any notes added by the checker so
		// we can add our own note at the start
		// to say which element failed.
		var notes []note
		err := c.elemChecker.Check(
			iter.value().Interface(),
			args,
			func(key string, val interface{}) {
				notes = append(notes, note{key, val})
			},
		)
		if err == nil {
			continue
		}
		if IsBadCheck(err) {
			return BadCheckf("at %s: %v", iter.key(), err)
		}
		notef("error", Unquoted("mismatch at "+iter.key()))
		// TODO should we print the whole container value in
		// verbose mode?
		if err != ErrSilent {
			// If the error's not silent, the checker is expecting
			// the caller to print the error and the value that failed.
			notef("error", Unquoted(err.Error()))
			notef("first mismatched element", iter.value().Interface())
		}
		for _, n := range notes {
			notef(n.key, n.value)
		}
		return ErrSilent
	}
	return nil
}

// argNames helps implementing Checker.ArgNames.
type argNames []string

// ArgNames implements Checker.ArgNames by returning the argument names.
func (a argNames) ArgNames() []string {
	return a
}

// match checks that the given error message matches the given pattern.
func match(got string, pattern interface{}, msg string, note func(key string, value interface{})) error {
	regex, ok := pattern.(string)
	if !ok {
		note("regexp", pattern)
		return BadCheckf("regexp is not a string")
	}
	matches, err := regexp.MatchString("^("+regex+")$", got)
	if err != nil {
		note("regexp", regex)
		return BadCheckf("cannot compile regexp: %s", err)
	}
	if matches {
		return nil
	}
	return errors.New(msg)
}

// canBeNil reports whether a value or type of the given kind can be nil.
func canBeNil(k reflect.Kind) bool {
	switch k {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}
