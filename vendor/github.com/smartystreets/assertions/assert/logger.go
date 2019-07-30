package assert

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// logger is meant be included as a pointer field on a struct. Leaving the
// instance as a nil reference will cause any calls on the *logger to forward
// to the corresponding functions from the standard log package. This is meant
// to be the behavior in production. In testing, set the field to a non-nil
// instance of a *logger to record log statements for later inspection.
type logger struct {
	*log.Logger

	Log   *bytes.Buffer
	Calls int
}

// capture creates a new *logger instance with an internal buffer. The prefix
// and flags default to the values of log.Prefix() and log.Flags(), respectively.
// This function is meant to be called from test code. See the godoc for the
// logger struct for details.
func capture() *logger {
	out := new(bytes.Buffer)
	inner := log.New(out, log.Prefix(), log.Flags())
	inner.SetPrefix("")
	return &logger{
		Log:    out,
		Logger: inner,
	}
}

// Fatal -> log.Fatal (except in testing it uses log.Print)
func (this *logger) Fatal(v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprint(v...))
		os.Exit(1)
	} else {
		this.Calls++
		this.Logger.Print(v...)
	}
}

// Panic -> log.Panic
func (this *logger) Panic(v ...interface{}) {
	if this == nil {
		s := fmt.Sprint(v...)
		this.Output(3, s)
		panic(s)
	} else {
		this.Calls++
		this.Logger.Panic(v...)
	}
}

// Print -> log.Print
func (this *logger) Print(v ...interface{}) {
	if this == nil {
		this.Output(3, fmt.Sprint(v...))
	} else {
		this.Calls++
		this.Logger.Print(v...)
	}
}

// Output -> log.Output
func (this *logger) Output(calldepth int, s string) error {
	if this == nil {
		return log.Output(calldepth, s)
	}
	this.Calls++
	return this.Logger.Output(calldepth, s)
}
