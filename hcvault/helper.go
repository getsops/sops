package hcvault

// TokenHelper is an interface that contains basic operations that must be
// implemented by a token helper
type TokenHelper interface {
	// Path displays a method-specific path; for the internal helper this
	// is the location of the token stored on disk; for the external helper
	// this is the location of the binary being invoked
	Path() string

	Erase() error
	Get() (string, error)
	Store(string) error
}

// DefaultTokenHelper returns the token helper that is configured for Vault.
func DefaultTokenHelper() (TokenHelper, error) {
	config, err := LoadConfig("")
	if err != nil {
		return nil, err
	}

	path := config.TokenHelper
	if path == "" {
		return NewInternalTokenHelper()
	}

	path, err = ExternalTokenHelperPath(path)
	if err != nil {
		return nil, err
	}
	return &ExternalTokenHelper{BinaryPath: path}, nil
}
