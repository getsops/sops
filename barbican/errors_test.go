package barbican

import (
	"errors"
	"testing"
)

func TestBarbicanError_Creation(t *testing.T) {
	tests := []struct {
		name          string
		errorType     BarbicanErrorType
		message       string
		expectedType  BarbicanErrorType
		expectedMsg   string
	}{
		{
			name:         "authentication error",
			errorType:    ErrorTypeAuthentication,
			message:      "Invalid credentials",
			expectedType: ErrorTypeAuthentication,
			expectedMsg:  "Invalid credentials",
		},
		{
			name:         "validation error",
			errorType:    ErrorTypeValidation,
			message:      "Invalid input format",
			expectedType: ErrorTypeValidation,
			expectedMsg:  "Invalid input format",
		},
		{
			name:         "network error",
			errorType:    ErrorTypeNetwork,
			message:      "Connection failed",
			expectedType: ErrorTypeNetwork,
			expectedMsg:  "Connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewBarbicanError(tt.errorType, tt.message)
			
			if err.Type != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, err.Type)
			}
			
			if err.Message != tt.expectedMsg {
				t.Errorf("Expected message %q, got %q", tt.expectedMsg, err.Message)
			}
		})
	}
}

func TestBarbicanError_WithMethods(t *testing.T) {
	baseErr := NewBarbicanError(ErrorTypeAuthentication, "Base error")
	
	// Test method chaining
	err := baseErr.
		WithDetails("Additional details").
		WithSuggestions("Try this", "Or this").
		WithCode(401).
		WithCause(errors.New("underlying error")).
		WithSecretRef("550e8400-e29b-41d4-a716-446655440000").
		WithRegion("sjc3")
	
	if err.Details != "Additional details" {
		t.Errorf("Expected details to be set")
	}
	
	if len(err.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(err.Suggestions))
	}
	
	if err.Code != 401 {
		t.Errorf("Expected code 401, got %d", err.Code)
	}
	
	if err.Cause == nil {
		t.Errorf("Expected cause to be set")
	}
	
	if err.SecretRef != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected secret ref to be set")
	}
	
	if err.Region != "sjc3" {
		t.Errorf("Expected region to be set")
	}
}

func TestBarbicanError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewBarbicanError(ErrorTypeNetwork, "Network error").WithCause(cause)
	
	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be the cause")
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name        string
		constructor func() *BarbicanError
		expectedType BarbicanErrorType
	}{
		{
			name:        "authentication error",
			constructor: func() *BarbicanError { return NewAuthenticationError("Auth failed") },
			expectedType: ErrorTypeAuthentication,
		},
		{
			name:        "authorization error",
			constructor: func() *BarbicanError { return NewAuthorizationError("Access denied") },
			expectedType: ErrorTypeAuthorization,
		},
		{
			name:        "validation error",
			constructor: func() *BarbicanError { return NewValidationError("Invalid input") },
			expectedType: ErrorTypeValidation,
		},
		{
			name:        "secret ref format error",
			constructor: func() *BarbicanError { return NewSecretRefFormatError("invalid-ref") },
			expectedType: ErrorTypeFormat,
		},
		{
			name:        "network error",
			constructor: func() *BarbicanError { return NewNetworkError("Network failed", nil) },
			expectedType: ErrorTypeNetwork,
		},
		{
			name:        "timeout error",
			constructor: func() *BarbicanError { return NewTimeoutError("Timeout", nil) },
			expectedType: ErrorTypeTimeout,
		},
		{
			name:        "service unavailable error",
			constructor: func() *BarbicanError { return NewServiceUnavailableError("Service down") },
			expectedType: ErrorTypeUnavailable,
		},
		{
			name:        "API error",
			constructor: func() *BarbicanError { return NewAPIError("API failed", 500) },
			expectedType: ErrorTypeAPI,
		},
		{
			name:        "secret not found error",
			constructor: func() *BarbicanError { return NewSecretNotFoundError("secret-ref") },
			expectedType: ErrorTypeNotFound,
		},
		{
			name:        "quota exceeded error",
			constructor: func() *BarbicanError { return NewQuotaExceededError("Quota exceeded") },
			expectedType: ErrorTypeQuota,
		},
		{
			name:        "TLS error",
			constructor: func() *BarbicanError { return NewTLSError("TLS failed", nil) },
			expectedType: ErrorTypeTLS,
		},
		{
			name:        "security error",
			constructor: func() *BarbicanError { return NewSecurityError("Security issue") },
			expectedType: ErrorTypeSecurity,
		},
		{
			name:        "config error",
			constructor: func() *BarbicanError { return NewConfigError("Config invalid") },
			expectedType: ErrorTypeConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor()
			
			if err.Type != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, err.Type)
			}
			
			// Check that suggestions are provided for user-facing errors
			if len(err.Suggestions) == 0 && (tt.expectedType == ErrorTypeAuthentication || 
				tt.expectedType == ErrorTypeValidation || tt.expectedType == ErrorTypeConfig) {
				t.Errorf("Expected suggestions for user-facing error type %v", tt.expectedType)
			}
		})
	}
}

func TestSanitizeSecretRef(t *testing.T) {
	tests := []struct {
		name      string
		secretRef string
		expected  string
	}{
		{
			name:      "empty string",
			secretRef: "",
			expected:  "",
		},
		{
			name:      "UUID format",
			secretRef: "550e8400-e29b-41d4-a716-446655440000",
			expected:  "***40000",
		},
		{
			name:      "regional format",
			secretRef: "region:sjc3:550e8400-e29b-41d4-a716-446655440000",
			expected:  "***40000",
		},
		{
			name:      "URI format",
			secretRef: "https://barbican.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000",
			expected:  "***40000",
		},
		{
			name:      "invalid format",
			secretRef: "invalid-format",
			expected:  "***",
		},
		{
			name:      "short UUID",
			secretRef: "123",
			expected:  "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeSecretRef(tt.secretRef)
			
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSanitizeEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		expected string
	}{
		{
			name:     "empty string",
			endpoint: "",
			expected: "",
		},
		{
			name:     "HTTPS URL",
			endpoint: "https://barbican.example.com:9311/v1",
			expected: "https://barbican.example.com:9311/***",
		},
		{
			name:     "HTTP URL",
			endpoint: "http://barbican.example.com:9311/v1/secrets",
			expected: "http://barbican.example.com:9311/***",
		},
		{
			name:     "non-URL string",
			endpoint: "not-a-url",
			expected: "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeEndpoint(tt.endpoint)
			
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		errorType    BarbicanErrorType
		message      string
		expectNil    bool
		expectType   BarbicanErrorType
	}{
		{
			name:      "nil error",
			err:       nil,
			errorType: ErrorTypeNetwork,
			message:   "Network failed",
			expectNil: true,
		},
		{
			name:       "existing BarbicanError",
			err:        NewAuthenticationError("Auth failed"),
			errorType:  ErrorTypeNetwork,
			message:    "Network failed",
			expectNil:  false,
			expectType: ErrorTypeAuthentication, // Should preserve original type
		},
		{
			name:       "standard error",
			err:        errors.New("standard error"),
			errorType:  ErrorTypeNetwork,
			message:    "Network failed",
			expectNil:  false,
			expectType: ErrorTypeNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.errorType, tt.message)
			
			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil result")
				}
				return
			}
			
			if result == nil {
				t.Errorf("Expected non-nil result")
				return
			}
			
			if result.Type != tt.expectType {
				t.Errorf("Expected type %v, got %v", tt.expectType, result.Type)
			}
		})
	}
}

func TestErrorStringFormatting(t *testing.T) {
	err := NewBarbicanError(ErrorTypeAuthentication, "Invalid credentials").
		WithDetails("Username not found").
		WithSecretRef("550e8400-e29b-41d4-a716-446655440000").
		WithRegion("sjc3").
		WithSuggestions("Check username", "Verify password")
	
	errorStr := err.Error()
	
	// Verify error string contains expected components
	expectedComponents := []string{
		"Barbican authentication error",
		"Invalid credentials",
		"Details: Username not found",
		"Secret: ***40000", // Sanitized secret ref
		"Region: sjc3",
		"Suggestions: Check username; Verify password",
	}
	
	for _, component := range expectedComponents {
		if !contains(errorStr, component) {
			t.Errorf("Error string should contain %q, got: %s", component, errorStr)
		}
	}
	
	// Verify full secret reference is not exposed
	if contains(errorStr, "550e8400-e29b-41d4-a716-446655440000") {
		t.Errorf("Error string should not contain full secret reference")
	}
}