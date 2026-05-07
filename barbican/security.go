package barbican

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
)

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	// TLS configuration
	InsecureTLS       bool
	CACertPath        string
	CACertContent     string
	SkipHostVerify    bool
	MinTLSVersion     uint16
	
	// Credential security
	SanitizeLogs      bool
	RedactCredentials bool
	
	// Warnings
	ShowSecurityWarnings bool
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		InsecureTLS:          false,
		SkipHostVerify:       false,
		MinTLSVersion:        tls.VersionTLS12,
		SanitizeLogs:         true,
		RedactCredentials:    true,
		ShowSecurityWarnings: true,
	}
}

// SecurityValidator validates and enforces security policies
type SecurityValidator struct {
	config *SecurityConfig
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator(config *SecurityConfig) *SecurityValidator {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	return &SecurityValidator{config: config}
}

// ValidateAndCreateTLSConfig creates a TLS configuration with security validation
func (sv *SecurityValidator) ValidateAndCreateTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:         sv.config.MinTLSVersion,
		InsecureSkipVerify: sv.config.InsecureTLS || sv.config.SkipHostVerify,
	}
	
	// Show security warnings for insecure configurations
	if sv.config.ShowSecurityWarnings {
		if sv.config.InsecureTLS {
			sv.logSecurityWarning("TLS certificate validation is disabled. This is insecure and should only be used for testing.")
		}
		
		if sv.config.SkipHostVerify {
			sv.logSecurityWarning("TLS hostname verification is disabled. This reduces security.")
		}
		
		if sv.config.MinTLSVersion < tls.VersionTLS12 {
			sv.logSecurityWarning("TLS version is set below TLS 1.2. This may be insecure.")
		}
	}
	
	// Load custom CA certificate if provided
	if sv.config.CACertPath != "" || sv.config.CACertContent != "" {
		caCertPool, err := sv.loadCACertificates()
		if err != nil {
			return nil, NewTLSError("Failed to load CA certificates", err)
		}
		tlsConfig.RootCAs = caCertPool
		
		log.Debug("Custom CA certificates loaded for TLS validation")
	}
	
	return tlsConfig, nil
}

// loadCACertificates loads CA certificates from file or content
func (sv *SecurityValidator) loadCACertificates() (*x509.CertPool, error) {
	caCertPool := x509.NewCertPool()
	
	var caCertData []byte
	var err error
	
	// Try to load from file path first
	if sv.config.CACertPath != "" {
		if _, statErr := os.Stat(sv.config.CACertPath); statErr == nil {
			caCertData, err = os.ReadFile(sv.config.CACertPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA certificate file %s: %w", sv.config.CACertPath, err)
			}
			log.WithField("ca_cert_path", sv.config.CACertPath).Debug("Loaded CA certificate from file")
		} else {
			return nil, fmt.Errorf("CA certificate file not found: %s", sv.config.CACertPath)
		}
	} else if sv.config.CACertContent != "" {
		// Use certificate content directly
		caCertData = []byte(sv.config.CACertContent)
		log.Debug("Using CA certificate from content")
	}
	
	if len(caCertData) == 0 {
		return nil, fmt.Errorf("no CA certificate data provided")
	}
	
	// Parse and add certificates to pool
	if !caCertPool.AppendCertsFromPEM(caCertData) {
		return nil, fmt.Errorf("failed to parse CA certificate data")
	}
	
	return caCertPool, nil
}

// ValidateAuthConfig validates authentication configuration for security issues
func (sv *SecurityValidator) ValidateAuthConfig(config *AuthConfig) error {
	if config == nil {
		return NewConfigError("Authentication configuration is required")
	}
	
	var warnings []string
	var errors []string
	
	// Validate authentication URL
	if config.AuthURL == "" {
		errors = append(errors, "Authentication URL is required")
	} else {
		if !strings.HasPrefix(config.AuthURL, "https://") {
			if sv.config.ShowSecurityWarnings {
				warnings = append(warnings, "Authentication URL is not using HTTPS. This is insecure.")
			}
		}
	}
	
	// Validate authentication method security
	hasAppCred := config.ApplicationCredentialID != "" && config.ApplicationCredentialSecret != ""
	hasToken := config.Token != ""
	hasPassword := config.Username != "" && config.Password != ""
	
	if !hasAppCred && !hasToken && !hasPassword {
		errors = append(errors, "No valid authentication method provided")
	}
	
	// Security recommendations
	if sv.config.ShowSecurityWarnings {
		if hasPassword && !hasAppCred {
			warnings = append(warnings, "Consider using application credentials instead of username/password for better security")
		}
		
		if config.Insecure {
			warnings = append(warnings, "TLS certificate validation is disabled. This should only be used for testing")
		}
	}
	
	// Log warnings
	for _, warning := range warnings {
		sv.logSecurityWarning(warning)
	}
	
	// Return errors
	if len(errors) > 0 {
		return NewConfigError(strings.Join(errors, "; "))
	}
	
	return nil
}

// SanitizeForLogging sanitizes sensitive information for safe logging
func (sv *SecurityValidator) SanitizeForLogging(data map[string]interface{}) map[string]interface{} {
	if !sv.config.SanitizeLogs {
		return data
	}
	
	sanitized := make(map[string]interface{})
	
	for key, value := range data {
		sanitized[key] = sv.sanitizeValue(key, value)
	}
	
	return sanitized
}

// sanitizeValue sanitizes a single value based on its key
func (sv *SecurityValidator) sanitizeValue(key string, value interface{}) interface{} {
	if !sv.config.RedactCredentials {
		return value
	}
	
	keyLower := strings.ToLower(key)
	
	// List of sensitive field patterns
	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"credential",
		"key",
		"auth",
		"x-auth-token",
		"x-subject-token",
	}
	
	for _, pattern := range sensitivePatterns {
		if strings.Contains(keyLower, pattern) {
			return sv.redactValue(value)
		}
	}
	
	// Special handling for URLs that might contain credentials
	if keyLower == "url" || keyLower == "endpoint" || keyLower == "auth_url" {
		if str, ok := value.(string); ok {
			return sv.sanitizeURL(str)
		}
	}
	
	return value
}

// redactValue redacts a sensitive value
func (sv *SecurityValidator) redactValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			return ""
		}
		if len(v) <= 4 {
			return "***"
		}
		// Show first and last character with *** in between
		return string(v[0]) + "***" + string(v[len(v)-1])
	default:
		return "***"
	}
}

// sanitizeURL sanitizes URLs that might contain credentials
func (sv *SecurityValidator) sanitizeURL(url string) string {
	if url == "" {
		return ""
	}
	
	// Check if URL contains credentials (user:pass@host format)
	if strings.Contains(url, "@") && (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		// Find the @ symbol and replace everything before it (after the scheme)
		parts := strings.Split(url, "://")
		if len(parts) == 2 {
			hostPart := parts[1]
			if atIndex := strings.Index(hostPart, "@"); atIndex != -1 {
				// Replace credentials with ***
				return parts[0] + "://***@" + hostPart[atIndex+1:]
			}
		}
	}
	
	return url
}

// ValidateSecretRef validates a secret reference for security issues
func (sv *SecurityValidator) ValidateSecretRef(secretRef string) error {
	if secretRef == "" {
		return NewValidationError("Secret reference cannot be empty")
	}
	
	// Check for obviously invalid or suspicious patterns
	if strings.Contains(secretRef, " ") {
		return NewSecretRefFormatError(secretRef).
			WithDetails("Secret reference contains spaces")
	}
	
	if strings.Contains(secretRef, "\n") || strings.Contains(secretRef, "\r") {
		return NewSecretRefFormatError(secretRef).
			WithDetails("Secret reference contains newline characters")
	}
	
	// Validate format using existing validation
	if !isValidSecretRef(secretRef) {
		return NewSecretRefFormatError(secretRef)
	}
	
	return nil
}

// CheckEndpointSecurity validates endpoint security
func (sv *SecurityValidator) CheckEndpointSecurity(endpoint string) error {
	if endpoint == "" {
		return NewConfigError("Endpoint cannot be empty")
	}
	
	// Warn about non-HTTPS endpoints
	if sv.config.ShowSecurityWarnings && !strings.HasPrefix(endpoint, "https://") {
		sv.logSecurityWarning(fmt.Sprintf("Endpoint %s is not using HTTPS. This is insecure.", sv.sanitizeURL(endpoint)))
	}
	
	// Check for localhost/private IPs in production-like environments
	if sv.config.ShowSecurityWarnings {
		if strings.Contains(endpoint, "localhost") || strings.Contains(endpoint, "127.0.0.1") {
			sv.logSecurityWarning("Using localhost endpoint. Ensure this is intended for local development only.")
		}
	}
	
	return nil
}

// logSecurityWarning logs a security warning
func (sv *SecurityValidator) logSecurityWarning(message string) {
	log.WithField("type", "security_warning").Warn(message)
}

// SecureCleanup provides secure cleanup utilities
type SecureCleanup struct {
	temporarySecrets []string
	validator        *SecurityValidator
}

// NewSecureCleanup creates a new secure cleanup manager
func NewSecureCleanup(validator *SecurityValidator) *SecureCleanup {
	return &SecureCleanup{
		temporarySecrets: make([]string, 0),
		validator:        validator,
	}
}

// AddTemporarySecret adds a secret to the cleanup list
func (sc *SecureCleanup) AddTemporarySecret(secretRef string) {
	sc.temporarySecrets = append(sc.temporarySecrets, secretRef)
	log.WithField("secret_ref", sc.validator.sanitizeValue("secret_ref", secretRef)).Debug("Added temporary secret for cleanup")
}

// CleanupTemporarySecrets cleans up all temporary secrets
func (sc *SecureCleanup) CleanupTemporarySecrets(client ClientInterface) error {
	if len(sc.temporarySecrets) == 0 {
		return nil
	}
	
	var errors []error
	cleaned := 0
	
	for _, secretRef := range sc.temporarySecrets {
		err := client.DeleteSecret(nil, secretRef)
		if err != nil {
			sanitizedRef := sc.validator.sanitizeValue("secret_ref", secretRef)
			log.WithError(err).WithField("secret_ref", sanitizedRef).Warn("Failed to cleanup temporary secret")
			errors = append(errors, fmt.Errorf("failed to cleanup secret %s: %w", sanitizedRef, err))
		} else {
			cleaned++
		}
	}
	
	log.WithField("cleaned", cleaned).WithField("total", len(sc.temporarySecrets)).Debug("Temporary secret cleanup completed")
	
	// Clear the list
	sc.temporarySecrets = sc.temporarySecrets[:0]
	
	if len(errors) > 0 {
		return fmt.Errorf("cleanup completed with %d errors: %v", len(errors), errors)
	}
	
	return nil
}

// GetTemporarySecretCount returns the number of temporary secrets pending cleanup
func (sc *SecureCleanup) GetTemporarySecretCount() int {
	return len(sc.temporarySecrets)
}

// SecurityConfigFromAuthConfig creates a SecurityConfig from AuthConfig
func SecurityConfigFromAuthConfig(authConfig *AuthConfig) *SecurityConfig {
	config := DefaultSecurityConfig()
	
	if authConfig != nil {
		config.InsecureTLS = authConfig.Insecure
		config.CACertPath = authConfig.CACert
		// If CACert looks like content (contains newlines), treat it as content
		if strings.Contains(authConfig.CACert, "\n") {
			config.CACertContent = authConfig.CACert
			config.CACertPath = ""
		}
	}
	
	return config
}

// ValidateSecurityConfiguration performs comprehensive security validation
func ValidateSecurityConfiguration(authConfig *AuthConfig) error {
	securityConfig := SecurityConfigFromAuthConfig(authConfig)
	validator := NewSecurityValidator(securityConfig)
	
	// Validate auth configuration
	if err := validator.ValidateAuthConfig(authConfig); err != nil {
		return err
	}
	
	// Validate TLS configuration
	_, err := validator.ValidateAndCreateTLSConfig()
	if err != nil {
		return err
	}
	
	return nil
}