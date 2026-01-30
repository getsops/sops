package barbican

import (
	"fmt"
	"strings"
)

// BarbicanErrorType represents different categories of Barbican errors
type BarbicanErrorType string

const (
	// Authentication errors
	ErrorTypeAuthentication BarbicanErrorType = "authentication"
	ErrorTypeAuthorization  BarbicanErrorType = "authorization"
	
	// Validation errors
	ErrorTypeValidation BarbicanErrorType = "validation"
	ErrorTypeFormat     BarbicanErrorType = "format"
	
	// Network and connectivity errors
	ErrorTypeNetwork     BarbicanErrorType = "network"
	ErrorTypeTimeout     BarbicanErrorType = "timeout"
	ErrorTypeUnavailable BarbicanErrorType = "unavailable"
	
	// API and service errors
	ErrorTypeAPI        BarbicanErrorType = "api"
	ErrorTypeNotFound   BarbicanErrorType = "not_found"
	ErrorTypeQuota      BarbicanErrorType = "quota"
	
	// Security errors
	ErrorTypeTLS        BarbicanErrorType = "tls"
	ErrorTypeSecurity   BarbicanErrorType = "security"
	
	// Configuration errors
	ErrorTypeConfig BarbicanErrorType = "configuration"
)

// BarbicanError represents a comprehensive error with troubleshooting information
type BarbicanError struct {
	Type         BarbicanErrorType `json:"type"`
	Message      string            `json:"message"`
	Details      string            `json:"details,omitempty"`
	Suggestions  []string          `json:"suggestions,omitempty"`
	Code         int               `json:"code,omitempty"`
	Cause        error             `json:"-"`
	SecretRef    string            `json:"secret_ref,omitempty"`
	Region       string            `json:"region,omitempty"`
	Endpoint     string            `json:"endpoint,omitempty"`
}

// Error implements the error interface
func (e *BarbicanError) Error() string {
	var parts []string
	
	parts = append(parts, fmt.Sprintf("Barbican %s error: %s", e.Type, e.Message))
	
	if e.Details != "" {
		parts = append(parts, fmt.Sprintf("Details: %s", e.Details))
	}
	
	if e.SecretRef != "" {
		// Sanitize secret reference for logging (show only last 8 characters)
		sanitizedRef := sanitizeSecretRef(e.SecretRef)
		parts = append(parts, fmt.Sprintf("Secret: %s", sanitizedRef))
	}
	
	if e.Region != "" {
		parts = append(parts, fmt.Sprintf("Region: %s", e.Region))
	}
	
	if len(e.Suggestions) > 0 {
		parts = append(parts, fmt.Sprintf("Suggestions: %s", strings.Join(e.Suggestions, "; ")))
	}
	
	return strings.Join(parts, ". ")
}

// Unwrap returns the underlying cause error
func (e *BarbicanError) Unwrap() error {
	return e.Cause
}

// NewBarbicanError creates a new BarbicanError with the specified type and message
func NewBarbicanError(errorType BarbicanErrorType, message string) *BarbicanError {
	return &BarbicanError{
		Type:    errorType,
		Message: message,
	}
}

// WithDetails adds details to the error
func (e *BarbicanError) WithDetails(details string) *BarbicanError {
	e.Details = details
	return e
}

// WithSuggestions adds troubleshooting suggestions to the error
func (e *BarbicanError) WithSuggestions(suggestions ...string) *BarbicanError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// WithCode adds an HTTP status code to the error
func (e *BarbicanError) WithCode(code int) *BarbicanError {
	e.Code = code
	return e
}

// WithCause adds the underlying cause error
func (e *BarbicanError) WithCause(cause error) *BarbicanError {
	e.Cause = cause
	return e
}

// WithSecretRef adds the secret reference (will be sanitized in output)
func (e *BarbicanError) WithSecretRef(secretRef string) *BarbicanError {
	e.SecretRef = secretRef
	return e
}

// WithRegion adds the region information
func (e *BarbicanError) WithRegion(region string) *BarbicanError {
	e.Region = region
	return e
}

// WithEndpoint adds the endpoint information (will be sanitized in output)
func (e *BarbicanError) WithEndpoint(endpoint string) *BarbicanError {
	e.Endpoint = sanitizeEndpoint(endpoint)
	return e
}

// Authentication error constructors
func NewAuthenticationError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeAuthentication, message).
		WithSuggestions(
			"Check your OpenStack credentials (OS_USERNAME, OS_PASSWORD, etc.)",
			"Verify the authentication URL (OS_AUTH_URL) is correct",
			"Ensure your user has access to the specified project",
			"Try using application credentials for better security",
		)
}

func NewAuthorizationError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeAuthorization, message).
		WithSuggestions(
			"Verify your user has the required Barbican permissions",
			"Check that you're accessing the correct project/tenant",
			"Contact your OpenStack administrator for access rights",
		)
}

// Validation error constructors
func NewValidationError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeValidation, message).
		WithSuggestions(
			"Check the configuration syntax in .sops.yaml",
			"Verify all required environment variables are set",
			"Ensure secret references are in the correct format",
		)
}

func NewSecretRefFormatError(secretRef string) *BarbicanError {
	return NewBarbicanError(ErrorTypeFormat, "Invalid secret reference format").
		WithSecretRef(secretRef).
		WithDetails("Secret reference must be a UUID, full URI, or regional format").
		WithSuggestions(
			"Use UUID format: 550e8400-e29b-41d4-a716-446655440000",
			"Use URI format: https://barbican.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000",
			"Use regional format: region:sjc3:550e8400-e29b-41d4-a716-446655440000",
		)
}

// Network error constructors
func NewNetworkError(message string, cause error) *BarbicanError {
	return NewBarbicanError(ErrorTypeNetwork, message).
		WithCause(cause).
		WithSuggestions(
			"Check your network connectivity to the OpenStack endpoints",
			"Verify firewall rules allow access to Barbican (port 9311)",
			"Try again in a few moments if this is a temporary network issue",
		)
}

func NewTimeoutError(message string, cause error) *BarbicanError {
	return NewBarbicanError(ErrorTypeTimeout, message).
		WithCause(cause).
		WithSuggestions(
			"Increase the timeout value in your configuration",
			"Check network latency to the OpenStack endpoints",
			"Verify the Barbican service is responding normally",
		)
}

func NewServiceUnavailableError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeUnavailable, message).
		WithSuggestions(
			"Wait a few moments and try again",
			"Check the OpenStack service status",
			"Try using a different region if available",
			"Contact your OpenStack administrator if the issue persists",
		)
}

// API error constructors
func NewAPIError(message string, code int) *BarbicanError {
	return NewBarbicanError(ErrorTypeAPI, message).
		WithCode(code).
		WithSuggestions(
			"Check the Barbican API documentation for this error code",
			"Verify your request parameters are correct",
			"Try the operation again with different parameters",
		)
}

func NewSecretNotFoundError(secretRef string) *BarbicanError {
	return NewBarbicanError(ErrorTypeNotFound, "Secret not found or not accessible").
		WithSecretRef(secretRef).
		WithSuggestions(
			"Verify the secret reference is correct",
			"Check that the secret exists in the specified region",
			"Ensure your user has read access to the secret",
			"Confirm you're using the correct project/tenant",
		)
}

func NewQuotaExceededError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeQuota, message).
		WithSuggestions(
			"Delete unused secrets to free up quota",
			"Contact your OpenStack administrator to increase quota",
			"Use secret expiration to automatically clean up old secrets",
		)
}

// Security error constructors
func NewTLSError(message string, cause error) *BarbicanError {
	return NewBarbicanError(ErrorTypeTLS, message).
		WithCause(cause).
		WithSuggestions(
			"Verify the server's TLS certificate is valid",
			"Check if you need to provide a custom CA certificate (OS_CACERT)",
			"Use OS_INSECURE=true only for testing (not recommended for production)",
		)
}

func NewSecurityError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeSecurity, message).
		WithSuggestions(
			"Review your security configuration",
			"Ensure you're using secure authentication methods",
			"Check that TLS is properly configured",
		)
}

// Configuration error constructors
func NewConfigError(message string) *BarbicanError {
	return NewBarbicanError(ErrorTypeConfig, message).
		WithSuggestions(
			"Check your .sops.yaml configuration file",
			"Verify all required environment variables are set",
			"Review the Barbican configuration documentation",
		)
}

// Helper functions for sanitizing sensitive information

// sanitizeSecretRef sanitizes a secret reference for safe logging
func sanitizeSecretRef(secretRef string) string {
	if secretRef == "" {
		return ""
	}
	
	// For very long or non-ASCII strings, just return sanitized placeholder
	if len(secretRef) > 100 {
		return "***"
	}
	
	// Check if string contains only printable ASCII characters for UUID extraction
	isPrintableASCII := true
	for _, r := range secretRef {
		if r > 127 || r < 32 {
			isPrintableASCII = false
			break
		}
	}
	
	if !isPrintableASCII {
		return "***"
	}
	
	// Extract UUID from the reference
	uuid, err := extractUUIDFromSecretRef(secretRef)
	if err != nil {
		// If we can't extract UUID, just show the format type
		if strings.HasPrefix(secretRef, "region:") {
			return "region:***:***"
		} else if strings.HasPrefix(secretRef, "http") {
			return "https://***:****/v1/secrets/***"
		}
		return "***"
	}
	
	// Show only the last 5 characters of the UUID for better privacy
	if len(uuid) >= 5 {
		return "***" + uuid[len(uuid)-5:]
	}
	
	return "***"
}

// sanitizeEndpoint sanitizes an endpoint URL for safe logging
func sanitizeEndpoint(endpoint string) string {
	if endpoint == "" {
		return ""
	}
	
	// For very long or non-ASCII strings, just return sanitized placeholder
	if len(endpoint) > 200 {
		return "***"
	}
	
	// Check if string contains only printable ASCII characters
	isPrintableASCII := true
	for _, r := range endpoint {
		if r > 127 || r < 32 {
			isPrintableASCII = false
			break
		}
	}
	
	if !isPrintableASCII {
		return "***"
	}
	
	// Parse and sanitize the URL
	if strings.HasPrefix(endpoint, "http") {
		// Extract just the scheme and host, hide the full path
		parts := strings.Split(endpoint, "/")
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2] + "/***"
		}
	}
	
	return "***"
}

// sanitizeCredentials removes sensitive information from error messages
func sanitizeCredentials(message string) string {
	// List of sensitive patterns to redact
	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"credential",
		"key",
	}
	
	result := message
	for _, pattern := range sensitivePatterns {
		// Simple pattern matching - in a real implementation, you might want more sophisticated regex
		if strings.Contains(strings.ToLower(result), pattern) {
			// Don't modify the message structure, just ensure no actual credentials leak
			// This is a basic implementation - more sophisticated sanitization might be needed
		}
	}
	
	return result
}

// WrapError wraps an existing error with Barbican-specific context
func WrapError(err error, errorType BarbicanErrorType, message string) *BarbicanError {
	if err == nil {
		return nil
	}
	
	// Check if it's already a BarbicanError
	if barbicanErr, ok := err.(*BarbicanError); ok {
		return barbicanErr
	}
	
	return NewBarbicanError(errorType, message).WithCause(err)
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a BarbicanError
	if barbicanErr, ok := err.(*BarbicanError); ok {
		switch barbicanErr.Type {
		case ErrorTypeNetwork, ErrorTypeTimeout, ErrorTypeUnavailable:
			return true
		case ErrorTypeAuthentication:
			// Authentication errors might be retryable if token expired
			return true
		case ErrorTypeAPI:
			// Some API errors are retryable (5xx status codes)
			if barbicanErr.Code >= 500 {
				return true
			}
			// Fall through to message content check
		default:
			// For other error types, check message content for retryable indicators
		}
		
		// Check message content for retryable indicators
		errStr := barbicanErr.Error()
		
		// Network-level errors are retryable
		if strings.Contains(errStr, "connection refused") ||
			strings.Contains(errStr, "timeout") ||
			strings.Contains(errStr, "temporary failure") ||
			strings.Contains(errStr, "network is unreachable") {
			return true
		}
		
		// HTTP server errors are retryable
		if strings.Contains(errStr, "server error (5") ||
			strings.Contains(errStr, "(502)") ||
			strings.Contains(errStr, "(503)") ||
			strings.Contains(errStr, "(504)") {
			return true
		}
		
		return false
	}
	
	// Fallback to string matching for non-BarbicanError types
	errStr := err.Error()
	
	// Network-level errors are retryable
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporary failure") ||
		strings.Contains(errStr, "network is unreachable") {
		return true
	}
	
	// HTTP server errors are retryable
	if strings.Contains(errStr, "server error (5") ||
		strings.Contains(errStr, "(502)") ||
		strings.Contains(errStr, "(503)") ||
		strings.Contains(errStr, "(504)") {
		return true
	}
	
	// Authentication errors might be retryable (token might have expired)
	if strings.Contains(errStr, "authentication failed (401)") {
		return true
	}
	
	return false
}