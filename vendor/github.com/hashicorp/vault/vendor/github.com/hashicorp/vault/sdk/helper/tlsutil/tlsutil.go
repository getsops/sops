package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/strutil"
)

var ErrInvalidCertParams = errors.New("invalid certificate parameters")

// TLSLookup maps the tls_min_version configuration to the internal value
var TLSLookup = map[string]uint16{
	"tls10": tls.VersionTLS10,
	"tls11": tls.VersionTLS11,
	"tls12": tls.VersionTLS12,
}

// cipherMap maps the cipher suite names to the internal cipher suite code.
var cipherMap = map[string]uint16{
	"TLS_RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
}

// ParseCiphers parse ciphersuites from the comma-separated string into recognized slice
func ParseCiphers(cipherStr string) ([]uint16, error) {
	suites := []uint16{}
	ciphers := strutil.ParseStringSlice(cipherStr, ",")
	for _, cipher := range ciphers {
		if v, ok := cipherMap[cipher]; ok {
			suites = append(suites, v)
		} else {
			return suites, fmt.Errorf("unsupported cipher %q", cipher)
		}
	}

	return suites, nil
}

// GetCipherName returns the name of a given cipher suite code or an error if the
// given cipher is unsupported.
func GetCipherName(cipher uint16) (string, error) {
	for cipherStr, cipherCode := range cipherMap {
		if cipherCode == cipher {
			return cipherStr, nil
		}
	}
	return "", fmt.Errorf("unsupported cipher %d", cipher)
}

func ClientTLSConfig(caCert []byte, clientCert []byte, clientKey []byte) (*tls.Config, error) {
	var tlsConfig *tls.Config
	var pool *x509.CertPool

	switch {
	case len(caCert) != 0:
		// Valid
	case len(clientCert) != 0 && len(clientKey) != 0:
		// Valid
	default:
		return nil, ErrInvalidCertParams
	}

	if len(caCert) != 0 {
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
	}

	tlsConfig = &tls.Config{
		RootCAs:    pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}

	var cert tls.Certificate
	var err error
	if len(clientCert) != 0 && len(clientKey) != 0 {
		cert, err = tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
