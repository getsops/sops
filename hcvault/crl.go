package hcvault

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// SopsHCVaultCRLFile can be set as an environment variable with a
	// comma-separated list of CRL file paths as value. Checks every
	// non-root certificate in the Vault server's chain against a
	// currently valid CRL.
	SopsHCVaultCRLFile = "SOPS_HC_VAULT_CRL_FILE"
)

// CRLKind classifies a CRL (Certificate Revocation List) verification failure
// by the remediation it calls for.
type CRLKind int

const (
	// CRLKindConfig indicates a CRL configuration or input problem: an empty
	// SOPS_HC_VAULT_CRL_FILE entry, a non-https Vault address, an unreadable,
	// unparseable, or empty CRL file, or an unsupported HTTP transport.
	// Remediation: fix the SOPS CRL configuration.
	CRLKindConfig CRLKind = iota
	// CRLKindNoValidCRL indicates that no currently valid CRL could be applied
	// to the Vault server's certificate chain (no verified chain with an
	// issuer, or no in-effect CRL signed by the issuer).
	// Remediation: supply a current CRL for the issuer.
	CRLKindNoValidCRL
	// CRLKindRevoked indicates that a certificate in the Vault server's chain
	// was found on a valid CRL. This is a security-relevant event.
	// Remediation: rotate/replace the revoked certificate.
	CRLKindRevoked
)

// CRLError is returned when CRL verification of the Vault server's certificate
// chain fails. It is a typed error so that callers can distinguish CRL failures
// from other transport errors, and distinguish CRL failure kinds from each other,
// via errors.As — without the hcvault package needing to know about CLI exit codes.
type CRLError struct {
	// Kind classifies the failure for remediation and exit-code mapping.
	Kind CRLKind
	// Err is the underlying error carrying the human-readable detail.
	Err error
}

// newCRLError wraps a formatted message in a *CRLError of the given kind.
func newCRLError(kind CRLKind, format string, args ...interface{}) *CRLError {
	return &CRLError{
		Kind: kind,
		Err:  fmt.Errorf(format, args...),
	}
}

// Error implements the error interface.
func (e *CRLError) Error() string {
	if e.Err == nil {
		return "CRL verification failed"
	}
	return e.Err.Error()
}

// Unwrap allows errors.Is/errors.As to traverse to the underlying error.
func (e *CRLError) Unwrap() error {
	return e.Err
}

// logCRLError logs a CRL verification failure once, at the Vault request
// boundary where the error surfaces.
func logCRLError(err error) {
	var crlErr *CRLError
	if !errors.As(err, &crlErr) {
		return
	}
	log.WithField("error", crlErr.Error()).Warn("Vault server certificate chain failed CRL check")
}

// configureCRL enables CRL checking when SopsHCVaultCRLFile is set.
// Requires an https Vault address. When enabled it layers CRL verification
// onto client, preserving its existing TLS settings.
func configureCRL(address string, client *http.Client) (*http.Client, error) {
	paths, err := getCRLFiles()
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return client, nil
	}

	log.WithField("files", len(paths)).Debug("CRL checking enabled")

	if !strings.HasPrefix(strings.ToLower(address), "https://") {
		log.WithField("address", address).Warn("CRL checking requires an https:// Vault address")
		return nil, newCRLError(CRLKindConfig,
			"%s is set but the Vault address %q is not an https:// endpoint; "+
				"CRL checking requires TLS", SopsHCVaultCRLFile, address)
	}

	crls, err := loadCRLs(paths)
	if err != nil {
		return nil, err
	}

	return applyCRLToClient(client, crls)
}

// getCRLFiles returns the list of CRL file paths configured through the
// SopsHCVaultCRLFile environment variable. Empty comma separated values are
// treated as a configuration error.
func getCRLFiles() ([]string, error) {
	value, ok := os.LookupEnv(SopsHCVaultCRLFile)
	if !ok || strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	paths := make([]string, 0, len(parts))
	for idx, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			log.WithField("entry", idx+1).Warn("CRL file list has an empty entry")
			return nil, newCRLError(CRLKindConfig, "%s's entry %d is empty", SopsHCVaultCRLFile, idx+1)
		}
		paths = append(paths, p)
	}
	return paths, nil
}

// loadCRLs reads and parses the CRLs from the provided file paths. Each file
// may contain a PEM-encoded ("X509 CRL") or DER-encoded CRL.
func loadCRLs(paths []string) ([]*x509.RevocationList, error) {
	var crls []*x509.RevocationList
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			log.WithFields(logrus.Fields{"file": p, "error": err}).Warn("Cannot read CRL file")
			return nil, newCRLError(CRLKindConfig, "cannot read CRL file %q: %w", p, err)
		}
		parsed, err := parseCRLs(data)
		if err != nil {
			log.WithFields(logrus.Fields{"file": p, "error": err}).Warn("CRL file is not a valid CRL")
			return nil, newCRLError(CRLKindConfig, "cannot parse CRL file %q: %w", p, err)
		}
		if len(parsed) == 0 {
			log.WithField("file", p).Warn("CRL file is valid PEM but contains no X509 CRL block")
			return nil, newCRLError(CRLKindConfig, "CRL file %q does not contain a certificate revocation list", p)
		}
		if log.IsLevelEnabled(logrus.DebugLevel) {
			for _, crl := range parsed {
				log.WithFields(logrus.Fields{
					"file":       p,
					"issuer":     crl.Issuer.CommonName,
					"thisUpdate": crl.ThisUpdate,
					"nextUpdate": crl.NextUpdate,
					"revoked":    len(crl.RevokedCertificateEntries),
				}).Debug("Loaded CRL")
			}
		}
		crls = append(crls, parsed...)
	}
	return crls, nil
}

// parseCRLs parses one or more CRLs from the given bytes. It supports
// PEM-encoded input (possibly containing multiple "X509 CRL" blocks) as well
// as a single DER-encoded CRL.
func parseCRLs(data []byte) ([]*x509.RevocationList, error) {
	var crls []*x509.RevocationList

	// Attempt PEM decoding first; fall back to raw DER if no PEM blocks are
	// present.
	rest := data
	foundPEM := false
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		foundPEM = true
		if block.Type != "X509 CRL" {
			continue
		}
		crl, err := x509.ParseRevocationList(block.Bytes)
		if err != nil {
			return nil, err
		}
		crls = append(crls, crl)
	}
	if foundPEM {
		return crls, nil
	}

	crl, err := x509.ParseRevocationList(data)
	if err != nil {
		return nil, err
	}
	return []*x509.RevocationList{crl}, nil
}

// applyCRLToClient returns an *http.Client, derived from client, whose TLS
// configuration performs CRL checking of the Vault server's certificate chain.
// The provided client is not mutated.
func applyCRLToClient(client *http.Client, crls []*x509.RevocationList) (*http.Client, error) {
	c := *client

	var transport *http.Transport
	switch rt := client.Transport.(type) {
	case nil:
		transport = http.DefaultTransport.(*http.Transport).Clone()
	case *http.Transport:
		transport = rt.Clone()
	default:
		log.WithField("transport", fmt.Sprintf("%T", client.Transport)).Warn("Cannot enable CRL checking: unsupported HTTP transport")
		return nil, newCRLError(CRLKindConfig, "cannot enable CRL checking: HTTP client uses a non-standard transport (%T)", client.Transport)
	}

	var tlsConfig *tls.Config
	if transport.TLSClientConfig != nil {
		tlsConfig = transport.TLSClientConfig.Clone()
	} else {
		tlsConfig = &tls.Config{}
	}

	crlVerify := newCRLVerifier(crls)
	if existing := tlsConfig.VerifyPeerCertificate; existing != nil {
		tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if err := existing(rawCerts, verifiedChains); err != nil {
				return err
			}
			return crlVerify(rawCerts, verifiedChains)
		}
	} else {
		tlsConfig.VerifyPeerCertificate = crlVerify
	}

	transport.TLSClientConfig = tlsConfig
	c.Transport = transport
	return &c, nil
}

// newCRLVerifier checks every non-root certificate in the chain against a
// currently valid CRL signed by its direct issuer.
func newCRLVerifier(crls []*x509.RevocationList) func([][]byte, [][]*x509.Certificate) error {
	return func(_ [][]byte, verifiedChains [][]*x509.Certificate) error {
		chain, err := verifiedChainWithIssuer(verifiedChains)
		if err != nil {
			return err
		}

		now := time.Now()
		for i := 0; i < len(chain)-1; i++ {
			if err := verifyCertificateCRL(chain[i], chain[i+1], crls, now); err != nil {
				return err
			}
		}
		return nil
	}
}

func verifiedChainWithIssuer(chains [][]*x509.Certificate) ([]*x509.Certificate, error) {
	if len(chains) == 0 || len(chains[0]) < 2 {
		return nil, newCRLError(CRLKindNoValidCRL,
			"cannot perform CRL check: no verified certificate chain with an issuer",
		)
	}
	return chains[0], nil
}

func verifyCertificateCRL(
	cert, issuer *x509.Certificate,
	crls []*x509.RevocationList,
	now time.Time,
) error {
	var applicable []*x509.RevocationList
	for _, crl := range crls {
		if crl.CheckSignatureFrom(issuer) != nil {
			continue
		}
		if !isCRLCurrentlyValid(crl, now) {
			continue
		}
		applicable = append(applicable, crl)
	}

	if len(applicable) == 0 {
		return newCRLError(CRLKindNoValidCRL,
			"no valid CRL found for issuer %q of vault server certificate chain (certificate %q)",
			issuer.Subject.CommonName,
			cert.Subject.CommonName,
		)
	}

	for _, crl := range applicable {
		if isCertificateRevoked(cert, crl) {
			return newCRLError(CRLKindRevoked,
				"vault server certificate chain entry %q (serial %s) has been revoked",
				cert.Subject.CommonName,
				cert.SerialNumber,
			)
		}
	}
	return nil
}

func isCRLCurrentlyValid(crl *x509.RevocationList, now time.Time) bool {
	if !crl.NextUpdate.IsZero() && now.After(crl.NextUpdate) {
		return false
	}
	return !now.Before(crl.ThisUpdate)
}

func isCertificateRevoked(cert *x509.Certificate, crl *x509.RevocationList) bool {
	for _, revoked := range crl.RevokedCertificateEntries {
		if revoked.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			return true
		}
	}
	return false
}
