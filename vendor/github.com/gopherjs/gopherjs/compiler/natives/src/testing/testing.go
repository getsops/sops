// +build js

package testing

// The upstream callerName and frameSkip rely on runtime.Callers,
// and panic if there are zero callers found. However, runtime.Callers
// is not implemented for GopherJS at this time, so we can't use
// that implementation. Use these stubs instead.
func callerName(skip int) string {
	// Upstream callerName requires a functional runtime.Callers.
	// TODO: Implement if possible.
	return "<unknown>"
}
func (*common) frameSkip(skip int) int {
	// Upstream frameSkip requires a functional runtime.Callers.
	// TODO: Implement if possible.
	return skip
}
