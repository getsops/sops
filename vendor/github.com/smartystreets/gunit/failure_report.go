package gunit

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type failureReport struct {
	Stack   []string
	Method  string
	Fixture string
	Package string
	Failure string
}

func newFailureReport(failure string) string {
	report := &failureReport{Failure: failure}
	report.ScanStack()
	return report.String()
}

func (this *failureReport) ScanStack() {
	for x := maxStackDepth; x >= 0; x-- {
		pc, file, line, ok := runtime.Caller(x)
		if !ok { // stack frame still too high
			continue
		}
		if !strings.HasSuffix(file, "_test.go") {
			continue
		}
		name := runtime.FuncForPC(pc).Name() // example: bitbucket.org/smartystreets/project/package.(*SomeFixture).TestSomething
		this.ParseTestName(name)
		this.Stack = append(this.Stack, fmt.Sprintf("%s:%d", file, line))
	}
}

func (this *failureReport) ParseTestName(name string) {
	if len(this.Method) > 0 {
		return
	}
	parts := strings.Split(name, ".")
	partCount := len(parts)
	last := partCount - 1
	if partCount < 3 {
		return
	}
	if !strings.HasPrefix(parts[last], "Test") {
		return
	}
	this.Method = parts[last]
	this.Fixture = parts[last-1]
	this.Package = strings.Join(parts[0:last-1], ".")
}

func (this failureReport) String() string {
	buffer := new(bytes.Buffer)
	fmt.Fprintf(buffer, "Test:     %s.%s()\n", this.Fixture, this.Method)
	for i, stack := range this.Stack {
		fmt.Fprintf(buffer, "(%d):      %s\n", len(this.Stack)-i-1, stack)
	}
	fmt.Fprintf(buffer, this.Failure)
	return buffer.String() + "\n\n"
}

const maxStackDepth = 24
