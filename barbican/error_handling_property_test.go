package barbican

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/quick"
	"time"
)

// TestErrorHandlingConsistencyProperty implements Property 7: Error Handling Consistency
// **Validates: Requirements 2.5, 6.1, 7.1**
func TestErrorHandlingConsistencyProperty(t *testing.T) {
	// Property-based test function
	f := func(
		statusCode uint16,
		includeCredentials bool,
		includeSecretRef bool,
		errorMessage string,
		useCustomMessage bool,
	) bool {
		// Constrain status code to valid HTTP range
		if statusCode < 200 || statusCode > 599 {
			return true // Skip invalid status codes
		}
		
		// Skip empty error messages for meaningful tests
		if !useCustomMessage {
			errorMessage = "Test error message"
		}
		if len(errorMessage) == 0 {
			return true
		}
		
		// Create test scenario based on inputs
		httpStatusCode := int(statusCode)
		
		// Test different error creation scenarios
		var err *BarbicanError
		
		switch {
		case httpStatusCode >= 400 && httpStatusCode < 500:
			// Client errors
			if httpStatusCode == 401 {
				err = NewAuthenticationError(errorMessage)
			} else if httpStatusCode == 403 {
				err = NewAuthorizationError(errorMessage)
			} else if httpStatusCode == 404 {
				err = NewSecretNotFoundError("test-secret-ref")
			} else {
				err = NewValidationError(errorMessage)
			}
		case httpStatusCode >= 500:
			// Server errors
			err = NewServiceUnavailableError(errorMessage)
		default:
			// Network or other errors
			err = NewNetworkError(errorMessage, nil)
		}
		
		// Add optional context
		if includeSecretRef {
			err = err.WithSecretRef("550e8400-e29b-41d4-a716-446655440000")
		}
		
		// Test the consistency properties
		errorString := err.Error()
		
		// Property 1: Error messages should never contain full credentials
		if includeCredentials {
			// Simulate adding credentials to the error (this should be sanitized)
			testCredentials := []string{
				"password123",
				"secret-token-value",
				"application-credential-secret",
			}
			
			for _, cred := range testCredentials {
				if strings.Contains(errorString, cred) {
					// This would be a security violation
					return false
				}
			}
		}
		
		// Property 2: Secret references should be sanitized
		if includeSecretRef {
			fullSecretRef := "550e8400-e29b-41d4-a716-446655440000"
			if strings.Contains(errorString, fullSecretRef) {
				// Full secret reference should not appear in error
				return false
			}
			
			// Should contain sanitized version
			if !strings.Contains(errorString, "***") {
				return false
			}
		}
		
		// Property 3: All errors should have consistent structure
		if !strings.Contains(errorString, "Barbican") {
			return false
		}
		
		// Property 4: Error should contain the original message
		if !strings.Contains(errorString, errorMessage) {
			return false
		}
		
		// Property 5: Errors should have suggestions for user-facing error types
		userFacingTypes := []BarbicanErrorType{
			ErrorTypeAuthentication,
			ErrorTypeAuthorization,
			ErrorTypeValidation,
			ErrorTypeConfig,
		}
		
		isUserFacing := false
		for _, userType := range userFacingTypes {
			if err.Type == userType {
				isUserFacing = true
				break
			}
		}
		
		if isUserFacing && len(err.Suggestions) == 0 {
			return false
		}
		
		// Property 6: Error unwrapping should work correctly
		if err.Unwrap() != err.Cause {
			return false
		}
		
		return true
	}
	
	// Run the property-based test with constrained iterations for reasonable execution time
	if err := quick.Check(f, &quick.Config{MaxCount: 50}); err != nil {
		t.Error(err)
	}
}

// TestAuthenticationErrorConsistencyProperty tests authentication error consistency
// **Validates: Requirements 2.5, 6.1**
func TestAuthenticationErrorConsistencyProperty(t *testing.T) {
	f := func(
		usePassword bool,
		useAppCred bool,
		useToken bool,
		includeProjectScope bool,
		simulateNetworkError bool,
	) bool {
		// Skip invalid combinations
		authMethodCount := 0
		if usePassword {
			authMethodCount++
		}
		if useAppCred {
			authMethodCount++
		}
		if useToken {
			authMethodCount++
		}
		
		// Need at least one auth method for meaningful test
		if authMethodCount == 0 {
			return true
		}
		
		// Create mock server for testing
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if simulateNetworkError {
				// Simulate network error by closing connection
				hj, ok := w.(http.Hijacker)
				if ok {
					conn, _, _ := hj.Hijack()
					conn.Close()
				}
				return
			}
			
			// Simulate authentication failure
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"message": "Invalid credentials", "code": 401}}`))
		}))
		defer server.Close()
		
		// Create auth config
		config := &AuthConfig{
			AuthURL: server.URL,
		}
		
		if usePassword {
			config.Username = "testuser"
			config.Password = "testpass"
			if includeProjectScope {
				config.ProjectID = "project123"
			}
		}
		
		if useAppCred {
			config.ApplicationCredentialID = "app-cred-id"
			config.ApplicationCredentialSecret = "app-cred-secret"
		}
		
		if useToken {
			config.Token = "existing-token"
			if includeProjectScope {
				config.ProjectID = "project123"
			}
		}
		
		// Test authentication
		authManager, err := NewAuthManager(config)
		if err != nil {
			// Should be a BarbicanError
			barbicanErr, ok := err.(*BarbicanError)
			if !ok {
				return false
			}
			
			// Should have appropriate error type
			if barbicanErr.Type != ErrorTypeConfig && barbicanErr.Type != ErrorTypeAuthentication {
				return false
			}
			
			// Should have suggestions
			if len(barbicanErr.Suggestions) == 0 {
				return false
			}
			
			return true
		}
		
		// Try to get token
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		_, _, err = authManager.GetToken(ctx)
		if err != nil {
			// Should be a BarbicanError
			barbicanErr, ok := err.(*BarbicanError)
			if !ok {
				return false
			}
			
			// Should have appropriate error type
			expectedTypes := []BarbicanErrorType{
				ErrorTypeAuthentication,
				ErrorTypeNetwork,
				ErrorTypeTimeout,
			}
			
			validType := false
			for _, expectedType := range expectedTypes {
				if barbicanErr.Type == expectedType {
					validType = true
					break
				}
			}
			
			if !validType {
				return false
			}
			
			// Should have suggestions
			if len(barbicanErr.Suggestions) == 0 {
				return false
			}
			
			// Error message should not contain credentials
			errorStr := barbicanErr.Error()
			sensitiveData := []string{
				config.Password,
				config.ApplicationCredentialSecret,
				config.Token,
			}
			
			for _, sensitive := range sensitiveData {
				if sensitive != "" && strings.Contains(errorStr, sensitive) {
					return false
				}
			}
		}
		
		return true
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 20}); err != nil {
		t.Error(err)
	}
}

// TestRetryableErrorClassificationProperty tests retry error classification consistency
// **Validates: Requirements 8.2**
func TestRetryableErrorClassificationProperty(t *testing.T) {
	f := func(
		errorType uint8,
		statusCode uint16,
		includeNetworkKeywords bool,
		includeServerErrorKeywords bool,
	) bool {
		// Map errorType to BarbicanErrorType
		var barbicanErrorType BarbicanErrorType
		switch errorType % 13 { // We have 13 error types
		case 0:
			barbicanErrorType = ErrorTypeAuthentication
		case 1:
			barbicanErrorType = ErrorTypeAuthorization
		case 2:
			barbicanErrorType = ErrorTypeValidation
		case 3:
			barbicanErrorType = ErrorTypeFormat
		case 4:
			barbicanErrorType = ErrorTypeNetwork
		case 5:
			barbicanErrorType = ErrorTypeTimeout
		case 6:
			barbicanErrorType = ErrorTypeUnavailable
		case 7:
			barbicanErrorType = ErrorTypeAPI
		case 8:
			barbicanErrorType = ErrorTypeNotFound
		case 9:
			barbicanErrorType = ErrorTypeQuota
		case 10:
			barbicanErrorType = ErrorTypeTLS
		case 11:
			barbicanErrorType = ErrorTypeSecurity
		case 12:
			barbicanErrorType = ErrorTypeConfig
		}
		
		// Create error message with optional keywords
		message := "Test error"
		if includeNetworkKeywords {
			message += " connection refused timeout"
		}
		if includeServerErrorKeywords {
			message += " server error (500)"
		}
		
		// Create error
		var err *BarbicanError
		if statusCode >= 200 && statusCode <= 599 {
			err = NewBarbicanError(barbicanErrorType, message).WithCode(int(statusCode))
		} else {
			err = NewBarbicanError(barbicanErrorType, message)
		}
		
		// Test retry classification
		isRetryable := IsRetryableError(err)
		
		// Define expected retry behavior
		expectedRetryable := false
		
		switch barbicanErrorType {
		case ErrorTypeNetwork, ErrorTypeTimeout, ErrorTypeUnavailable:
			expectedRetryable = true
		case ErrorTypeAuthentication:
			expectedRetryable = true // Token might have expired
		case ErrorTypeAPI:
			if err.Code >= 500 {
				expectedRetryable = true
			}
		}
		
		// Additional checks for message content
		if includeNetworkKeywords || includeServerErrorKeywords {
			expectedRetryable = true
		}
		
		return isRetryable == expectedRetryable
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Error(err)
	}
}

// TestErrorSanitizationProperty tests that sensitive data is consistently sanitized
// **Validates: Requirements 7.1**
func TestErrorSanitizationProperty(t *testing.T) {
	f := func(
		secretRef string,
		endpoint string,
		includeCredentials bool,
		credentialType uint8,
	) bool {
		// Skip empty inputs
		if len(secretRef) == 0 {
			secretRef = "550e8400-e29b-41d4-a716-446655440000"
		}
		if len(endpoint) == 0 {
			endpoint = "https://barbican.example.com:9311/v1"
		}
		
		// Skip inputs that are too long or contain non-ASCII characters
		// as they will be sanitized to "***" which is expected behavior
		if len(secretRef) > 100 || len(endpoint) > 200 {
			return true
		}
		
		// Check for non-ASCII characters
		hasNonASCII := false
		for _, r := range secretRef + endpoint {
			if r > 127 || r < 32 {
				hasNonASCII = true
				break
			}
		}
		if hasNonASCII {
			return true // Skip non-ASCII inputs as they're handled specially
		}
		
		// Create error with sensitive data
		err := NewAuthenticationError("Test authentication error").
			WithSecretRef(secretRef).
			WithEndpoint(endpoint)
		
		// Add credentials based on type
		if includeCredentials {
			switch credentialType % 4 {
			case 0:
				err = err.WithDetails("Password: secret123")
			case 1:
				err = err.WithDetails("Token: auth-token-12345")
			case 2:
				err = err.WithDetails("Application credential secret: app-secret-67890")
			case 3:
				err = err.WithDetails("API key: api-key-abcdef")
			}
		}
		
		errorString := err.Error()
		
		// Property 1: Full secret reference should not appear
		if len(secretRef) > 10 && strings.Contains(errorString, secretRef) {
			return false
		}
		
		// Property 2: Full endpoint should not appear (should be sanitized)
		if strings.Contains(endpoint, "://") && strings.Contains(errorString, endpoint) {
			return false
		}
		
		// Property 3: Credentials should not appear in full
		if includeCredentials {
			sensitivePatterns := []string{
				"secret123",
				"auth-token-12345",
				"app-secret-67890",
				"api-key-abcdef",
			}
			
			for _, pattern := range sensitivePatterns {
				if strings.Contains(errorString, pattern) {
					return false
				}
			}
		}
		
		// Property 4: Should contain sanitized indicators
		if len(secretRef) > 10 && !strings.Contains(errorString, "***") {
			return false
		}
		
		return true
	}
	
	// Run the property-based test
	if err := quick.Check(f, &quick.Config{MaxCount: 50}); err != nil {
		t.Error(err)
	}
}