package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	// Whether cgo is enabled or not; set at build time
	CgoEnabled bool

	Version           = "1.2.0"
	VersionPrerelease = "beta2"
	VersionMetadata   = ""
)
