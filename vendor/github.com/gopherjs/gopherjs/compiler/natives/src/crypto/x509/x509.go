// +build js

package x509

import "os"

func loadSystemRoots() (*CertPool, error) {
	// no system roots
	return NewCertPool(), nil
}

func execSecurityRoots() (*CertPool, error) {
	return nil, os.ErrNotExist
}
