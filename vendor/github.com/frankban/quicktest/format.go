// Licensed under the MIT license, see LICENCE file for details.

package quicktest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kr/pretty"
)

// Format formats the given value as a string. It is used to print values in
// test failures unless that's changed by calling C.SetFormat.
func Format(v interface{}) string {
	switch v := v.(type) {
	case error:
		return formatErr(v)
	case fmt.Stringer:
		return "s" + quoteString(v.String())
	case string:
		return quoteString(v)
	}
	// The pretty.Sprint equivalent does not quote string values.
	return fmt.Sprintf("%# v", pretty.Formatter(v))
}

func formatErr(err error) string {
	s := fmt.Sprintf("%+v", err)
	if s != err.Error() {
		// The error has formatted itself with additional information.
		// Leave that as is.
		return s
	}
	return "e" + quoteString(s)
}

func quoteString(s string) string {
	// TODO think more about what to do about multi-line strings.
	if strings.Contains(s, `"`) && !strings.Contains(s, "\n") && strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

type formatFunc func(interface{}) string
