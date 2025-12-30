// Package barbican provides OpenStack Barbican API client functionality.
package barbican

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ClientInterface defines the operations needed from Barbican API.
// This interface allows for easy mocking and testing of Barbican operations.
type ClientInterface interface {
	// StoreSecret stores a secret payload in Barbican and returns the secret reference
	StoreSecret(ctx context.Context, payload []byte, metadata SecretMetadata) (string, error)
	// GetSecretPayload retrieves the payload of a secret from Barbican
	GetSecretPayload(ctx context.Context, secretRef string) ([]byte, error)
	// DeleteSecret removes a secret from Barbican
	DeleteSecret(ctx context.Context, secretRef string) error
	// ValidateSecretExists checks if a secret exists and is accessible
	ValidateSecretExists(ctx context.Context, secretRef string) error
}

// BarbicanClient provides an interface to OpenStack Barbican API with
// connection pooling, retry logic, and timeout handling.
type BarbicanClient struct {
	endpoint    string
	authManager *AuthManager
	httpClient  *http.Client
	config      *ClientConfig
}

// ClientConfig holds configuration for the Barbican client
type ClientConfig struct {
	// Timeout for HTTP requests
	Timeout time.Duration
	// MaxRetries for failed requests
	MaxRetries int
	// InitialRetryDelay for exponential backoff
	InitialRetryDelay time.Duration
	// MaxRetryDelay caps the retry delay
	MaxRetryDelay time.Duration
	// RetryMultiplier for exponential backoff
	RetryMultiplier float64
	// Insecure disables TLS certificate validation
	Insecure bool
	// CACert is the path to or content of a custom CA certificate
	CACert string
	
	// Connection pool configuration for performance optimization
	// MaxIdleConns controls the maximum number of idle connections across all hosts
	MaxIdleConns int
	// MaxIdleConnsPerHost controls the maximum idle connections per host
	MaxIdleConnsPerHost int
	// MaxConnsPerHost controls the maximum connections per host
	MaxConnsPerHost int
	// IdleConnTimeout is the maximum amount of time an idle connection will remain idle
	IdleConnTimeout time.Duration
	// DisableCompression disables compression for requests
	DisableCompression bool
}

// DefaultClientConfig returns a default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		InitialRetryDelay: 1 * time.Second,
		MaxRetryDelay:     30 * time.Second,
		RetryMultiplier:   2.0,
		Insecure:          false,
		
		// Connection pool defaults optimized for performance
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}
}

// HighPerformanceClientConfig returns a client configuration optimized for high-throughput scenarios
func HighPerformanceClientConfig() *ClientConfig {
	config := DefaultClientConfig()
	
	// Increase connection pool sizes for high throughput
	config.MaxIdleConns = 200
	config.MaxIdleConnsPerHost = 20
	config.MaxConnsPerHost = 100
	
	// Reduce timeouts for faster failure detection
	config.Timeout = 15 * time.Second
	config.MaxRetryDelay = 10 * time.Second
	
	// Keep connections alive longer for reuse
	config.IdleConnTimeout = 120 * time.Second
	
	return config
}

// LowLatencyClientConfig returns a client configuration optimized for low-latency scenarios
func LowLatencyClientConfig() *ClientConfig {
	config := DefaultClientConfig()
	
	// Reduce timeouts for faster responses
	config.Timeout = 10 * time.Second
	config.InitialRetryDelay = 500 * time.Millisecond
	config.MaxRetryDelay = 5 * time.Second
	
	// Fewer retries for faster failure
	config.MaxRetries = 2
	
	// Smaller connection pools to reduce overhead
	config.MaxIdleConns = 50
	config.MaxIdleConnsPerHost = 5
	config.MaxConnsPerHost = 25
	
	return config
}

// MultiRegionClientConfig returns a client configuration optimized for multi-region operations
func MultiRegionClientConfig() *ClientConfig {
	config := DefaultClientConfig()
	
	// Larger connection pools to handle multiple regions
	config.MaxIdleConns = 300
	config.MaxIdleConnsPerHost = 30
	config.MaxConnsPerHost = 150
	
	// Longer timeouts to handle cross-region latency
	config.Timeout = 60 * time.Second
	config.MaxRetryDelay = 60 * time.Second
	
	// More retries for network reliability across regions
	config.MaxRetries = 5
	config.RetryMultiplier = 1.5
	
	// Keep connections alive longer for cross-region reuse
	config.IdleConnTimeout = 300 * time.Second
	
	return config
}

// PerformanceMetrics tracks performance statistics for Barbican operations
type PerformanceMetrics struct {
	// Operation counters
	EncryptOperations int64
	DecryptOperations int64
	AuthOperations    int64
	
	// Timing statistics (in milliseconds)
	AvgEncryptTime    int64
	AvgDecryptTime    int64
	AvgAuthTime       int64
	
	// Error counters
	EncryptErrors     int64
	DecryptErrors     int64
	AuthErrors        int64
	
	// Connection pool statistics
	ActiveConnections int64
	IdleConnections   int64
	
	// Retry statistics
	TotalRetries      int64
	SuccessfulRetries int64
}

// GetPerformanceMetrics returns current performance metrics for a Barbican client
func (c *BarbicanClient) GetPerformanceMetrics() *PerformanceMetrics {
	// In a real implementation, this would collect actual metrics
	// For now, return a placeholder structure
	return &PerformanceMetrics{
		// These would be populated from actual monitoring
		EncryptOperations: 0,
		DecryptOperations: 0,
		AuthOperations:    0,
		AvgEncryptTime:    0,
		AvgDecryptTime:    0,
		AvgAuthTime:       0,
		EncryptErrors:     0,
		DecryptErrors:     0,
		AuthErrors:        0,
		ActiveConnections: 0,
		IdleConnections:   0,
		TotalRetries:      0,
		SuccessfulRetries: 0,
	}
}

// OptimizeClientConfig analyzes performance metrics and suggests configuration optimizations
func OptimizeClientConfig(metrics *PerformanceMetrics, currentConfig *ClientConfig) *ClientConfig {
	optimized := *currentConfig // Copy current config
	
	// Analyze error rates and adjust retry settings
	if metrics.EncryptErrors > 0 || metrics.DecryptErrors > 0 {
		totalOps := metrics.EncryptOperations + metrics.DecryptOperations
		errorRate := float64(metrics.EncryptErrors+metrics.DecryptErrors) / float64(totalOps)
		
		if errorRate > 0.1 { // More than 10% error rate
			// Increase retry attempts and delays
			optimized.MaxRetries = min(optimized.MaxRetries+1, 10)
			if optimized.MaxRetryDelay*2 < 120*time.Second {
				optimized.MaxRetryDelay = optimized.MaxRetryDelay * 2
			} else {
				optimized.MaxRetryDelay = 120 * time.Second
			}
			log.WithField("error_rate", errorRate).Info("Increased retry settings due to high error rate")
		}
	}
	
	// Analyze timing and adjust timeouts
	if metrics.AvgEncryptTime > 0 && metrics.AvgDecryptTime > 0 {
		avgOpTime := (metrics.AvgEncryptTime + metrics.AvgDecryptTime) / 2
		suggestedTimeout := time.Duration(avgOpTime*3) * time.Millisecond // 3x average time
		
		if suggestedTimeout > optimized.Timeout {
			optimized.Timeout = suggestedTimeout
			log.WithField("new_timeout", suggestedTimeout).Info("Increased timeout based on operation timing")
		}
	}
	
	// Analyze connection usage and adjust pool sizes
	if metrics.ActiveConnections > 0 {
		utilizationRate := float64(metrics.ActiveConnections) / float64(currentConfig.MaxConnsPerHost)
		
		if utilizationRate > 0.8 { // More than 80% utilization
			// Increase connection pool sizes
			optimized.MaxConnsPerHost = min(optimized.MaxConnsPerHost*2, 200)
			optimized.MaxIdleConnsPerHost = min(optimized.MaxIdleConnsPerHost*2, 50)
			log.WithField("utilization_rate", utilizationRate).Info("Increased connection pool sizes due to high utilization")
		}
	}
	
	return &optimized
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BatchOperation represents a batch operation request
type BatchOperation struct {
	Operation string      `json:"operation"` // "encrypt" or "decrypt"
	Payload   []byte      `json:"payload,omitempty"`
	SecretRef string      `json:"secret_ref,omitempty"`
	Metadata  SecretMetadata `json:"metadata,omitempty"`
}

// BatchResult represents the result of a batch operation
type BatchResult struct {
	Success   bool   `json:"success"`
	SecretRef string `json:"secret_ref,omitempty"`
	Payload   []byte `json:"payload,omitempty"`
	Error     string `json:"error,omitempty"`
}

// BatchOperationRequest represents a batch of operations
type BatchOperationRequest struct {
	Operations []BatchOperation `json:"operations"`
}

// BatchOperationResponse represents the response from a batch operation
type BatchOperationResponse struct {
	Results []BatchResult `json:"results"`
}

// SecretCreateRequest represents a request to create a secret in Barbican
type SecretCreateRequest struct {
	Name                 string            `json:"name,omitempty"`
	Algorithm            string            `json:"algorithm,omitempty"`
	BitLength            int               `json:"bit_length,omitempty"`
	Mode                 string            `json:"mode,omitempty"`
	SecretType           string            `json:"secret_type,omitempty"`
	PayloadContentType   string            `json:"payload_content_type,omitempty"`
	Payload              string            `json:"payload,omitempty"`
	PayloadContentEncoding string          `json:"payload_content_encoding,omitempty"`
	Expiration           *time.Time        `json:"expiration,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// SecretCreateResponse represents the response from creating a secret
type SecretCreateResponse struct {
	SecretRef string `json:"secret_ref"`
}

// SecretResponse represents a secret object from Barbican
type SecretResponse struct {
	SecretRef            string            `json:"secret_ref"`
	Name                 string            `json:"name"`
	Algorithm            string            `json:"algorithm"`
	BitLength            int               `json:"bit_length"`
	Mode                 string            `json:"mode"`
	SecretType           string            `json:"secret_type"`
	PayloadContentType   string            `json:"payload_content_type"`
	Status               string            `json:"status"`
	CreatedAt            string            `json:"created"`
	UpdatedAt            string            `json:"updated"`
	Expiration           *string           `json:"expiration"`
	ContentTypes         map[string]string `json:"content_types"`
	Metadata             map[string]string `json:"metadata"`
}

// ErrorResponse represents an error response from Barbican
type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Title   string `json:"title"`
	} `json:"error"`
}

// NewBarbicanClient creates a new Barbican client
func NewBarbicanClient(endpoint string, authManager *AuthManager, config *ClientConfig) (*BarbicanClient, error) {
	if endpoint == "" {
		return nil, NewConfigError("Barbican endpoint is required")
	}
	
	if authManager == nil {
		return nil, NewConfigError("Authentication manager is required")
	}

	if config == nil {
		config = DefaultClientConfig()
	}

	// Validate endpoint security
	securityValidator := NewSecurityValidator(SecurityConfigFromAuthConfig(authManager.config))
	if err := securityValidator.CheckEndpointSecurity(endpoint); err != nil {
		return nil, err
	}

	// Create HTTP client with connection pooling and timeout handling
	httpClient, err := createBarbicanHTTPClient(config)
	if err != nil {
		return nil, NewTLSError("Failed to create HTTP client", err)
	}

	// Ensure endpoint has proper format
	endpoint = strings.TrimSuffix(endpoint, "/")
	if !strings.HasSuffix(endpoint, "/v1") {
		endpoint = endpoint + "/v1"
	}

	return &BarbicanClient{
		endpoint:    endpoint,
		authManager: authManager,
		httpClient:  httpClient,
		config:      config,
	}, nil
}

// createBarbicanHTTPClient creates an HTTP client with proper configuration
func createBarbicanHTTPClient(config *ClientConfig) (*http.Client, error) {
	// Create security validator and TLS config
	securityConfig := &SecurityConfig{
		InsecureTLS:          config.Insecure,
		CACertPath:           config.CACert,
		MinTLSVersion:        tls.VersionTLS12,
		SanitizeLogs:         true,
		RedactCredentials:    true,
		ShowSecurityWarnings: true,
	}
	
	validator := NewSecurityValidator(securityConfig)
	tlsConfig, err := validator.ValidateAndCreateTLSConfig()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		IdleConnTimeout:     config.IdleConnTimeout,
		DisableCompression:  config.DisableCompression,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		TLSClientConfig:     tlsConfig,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}, nil
}

// StoreSecret stores a secret in Barbican and returns the secret reference
func (c *BarbicanClient) StoreSecret(ctx context.Context, payload []byte, metadata SecretMetadata) (string, error) {
	// Build the request
	request := SecretCreateRequest{
		Name:                   metadata.Name,
		Algorithm:              metadata.Algorithm,
		BitLength:              metadata.BitLength,
		Mode:                   metadata.Mode,
		SecretType:             metadata.SecretType,
		PayloadContentType:     metadata.ContentType,
		Payload:                base64.StdEncoding.EncodeToString(payload),
		PayloadContentEncoding: "base64",
		Expiration:             metadata.Expiration,
		Metadata:               metadata.Metadata,
	}

	// Set defaults if not provided
	if request.SecretType == "" {
		request.SecretType = "opaque"
	}
	if request.PayloadContentType == "" {
		request.PayloadContentType = "application/octet-stream"
	}
	if request.Name == "" {
		request.Name = "SOPS Data Key"
	}

	// Marshal request
	reqBody, err := json.Marshal(request)
	if err != nil {
		return "", NewAPIError("Failed to marshal secret create request", 0).WithCause(err)
	}

	// Make the API call with retry logic
	var response SecretCreateResponse
	err = c.doRequestWithRetry(ctx, "POST", "/secrets", reqBody, &response)
	if err != nil {
		return "", WrapError(err, ErrorTypeAPI, "Failed to store secret")
	}

	// Sanitize secret reference for logging
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	sanitizedRef := securityValidator.sanitizeValue("secret_ref", response.SecretRef)
	log.WithField("secret_ref", sanitizedRef).Debug("Secret stored successfully")
	
	return response.SecretRef, nil
}

// GetSecretPayload retrieves the payload of a secret from Barbican
func (c *BarbicanClient) GetSecretPayload(ctx context.Context, secretRef string) ([]byte, error) {
	// Validate secret reference
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	if err := securityValidator.ValidateSecretRef(secretRef); err != nil {
		return nil, err
	}

	// Extract UUID from secret reference
	uuid, err := extractUUIDFromSecretRef(secretRef)
	if err != nil {
		return nil, NewSecretRefFormatError(secretRef).WithCause(err)
	}

	// Build the path for payload retrieval
	path := fmt.Sprintf("/secrets/%s/payload", uuid)

	// Make the API call with retry logic
	var payload []byte
	err = c.doRequestWithRetry(ctx, "GET", path, nil, &payload)
	if err != nil {
		return nil, WrapError(err, ErrorTypeAPI, "Failed to retrieve secret payload").WithSecretRef(secretRef)
	}

	// Sanitize UUID for logging
	sanitizedUUID := securityValidator.sanitizeValue("secret_uuid", uuid)
	log.WithField("secret_uuid", sanitizedUUID).Debug("Secret payload retrieved successfully")
	
	return payload, nil
}

// DeleteSecret deletes a secret from Barbican
func (c *BarbicanClient) DeleteSecret(ctx context.Context, secretRef string) error {
	// Validate secret reference
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	if err := securityValidator.ValidateSecretRef(secretRef); err != nil {
		return err
	}

	// Extract UUID from secret reference
	uuid, err := extractUUIDFromSecretRef(secretRef)
	if err != nil {
		return NewSecretRefFormatError(secretRef).WithCause(err)
	}

	// Build the path for secret deletion
	path := fmt.Sprintf("/secrets/%s", uuid)

	// Make the API call with retry logic
	err = c.doRequestWithRetry(ctx, "DELETE", path, nil, nil)
	if err != nil {
		return WrapError(err, ErrorTypeAPI, "Failed to delete secret").WithSecretRef(secretRef)
	}

	// Sanitize UUID for logging
	sanitizedUUID := securityValidator.sanitizeValue("secret_uuid", uuid)
	log.WithField("secret_uuid", sanitizedUUID).Debug("Secret deleted successfully")
	
	return nil
}

// ValidateSecretExists checks if a secret exists in Barbican
func (c *BarbicanClient) ValidateSecretExists(ctx context.Context, secretRef string) error {
	// Validate secret reference
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	if err := securityValidator.ValidateSecretRef(secretRef); err != nil {
		return err
	}

	// Extract UUID from secret reference
	uuid, err := extractUUIDFromSecretRef(secretRef)
	if err != nil {
		return NewSecretRefFormatError(secretRef).WithCause(err)
	}

	// Build the path for secret metadata retrieval
	path := fmt.Sprintf("/secrets/%s", uuid)

	// Make the API call with retry logic
	var secret SecretResponse
	err = c.doRequestWithRetry(ctx, "GET", path, nil, &secret)
	if err != nil {
		return NewSecretNotFoundError(secretRef).WithCause(err)
	}

	// Sanitize UUID for logging
	sanitizedUUID := securityValidator.sanitizeValue("secret_uuid", uuid)
	log.WithField("secret_uuid", sanitizedUUID).Debug("Secret exists and is accessible")
	
	return nil
}

// StoreBatchSecrets stores multiple secrets in parallel for improved performance
func (c *BarbicanClient) StoreBatchSecrets(ctx context.Context, payloads [][]byte, metadatas []SecretMetadata) ([]string, error) {
	if len(payloads) != len(metadatas) {
		return nil, NewValidationError("Number of payloads must match number of metadata entries")
	}
	
	if len(payloads) == 0 {
		return []string{}, nil
	}
	
	// Use channels to collect results from parallel operations
	type storeResult struct {
		index     int
		secretRef string
		error     error
	}
	
	resultChan := make(chan storeResult, len(payloads))
	
	// Start storage operations in parallel
	for i, payload := range payloads {
		go func(idx int, data []byte, metadata SecretMetadata) {
			secretRef, err := c.StoreSecret(ctx, data, metadata)
			resultChan <- storeResult{
				index:     idx,
				secretRef: secretRef,
				error:     err,
			}
		}(i, payload, metadatas[i])
	}
	
	// Collect results in order
	results := make([]string, len(payloads))
	var errors []error
	
	for i := 0; i < len(payloads); i++ {
		result := <-resultChan
		if result.error != nil {
			errors = append(errors, fmt.Errorf("operation %d: %w", result.index, result.error))
		} else {
			results[result.index] = result.secretRef
		}
	}
	
	// Return partial results if some operations succeeded
	if len(errors) > 0 {
		log.WithField("failed_operations", len(errors)).WithField("total_operations", len(payloads)).Warn("Some batch store operations failed")
		
		// If all operations failed, return error
		if len(errors) == len(payloads) {
			return nil, fmt.Errorf("all batch store operations failed: %v", errors)
		}
	}
	
	log.WithField("successful_operations", len(payloads)-len(errors)).WithField("total_operations", len(payloads)).Debug("Batch store operations completed")
	return results, nil
}

// GetBatchSecretPayloads retrieves multiple secret payloads in parallel for improved performance
func (c *BarbicanClient) GetBatchSecretPayloads(ctx context.Context, secretRefs []string) ([][]byte, error) {
	if len(secretRefs) == 0 {
		return [][]byte{}, nil
	}
	
	// Use channels to collect results from parallel operations
	type retrieveResult struct {
		index   int
		payload []byte
		error   error
	}
	
	resultChan := make(chan retrieveResult, len(secretRefs))
	
	// Start retrieval operations in parallel
	for i, secretRef := range secretRefs {
		go func(idx int, ref string) {
			payload, err := c.GetSecretPayload(ctx, ref)
			resultChan <- retrieveResult{
				index:   idx,
				payload: payload,
				error:   err,
			}
		}(i, secretRef)
	}
	
	// Collect results in order
	results := make([][]byte, len(secretRefs))
	var errors []error
	
	for i := 0; i < len(secretRefs); i++ {
		result := <-resultChan
		if result.error != nil {
			errors = append(errors, fmt.Errorf("operation %d: %w", result.index, result.error))
		} else {
			results[result.index] = result.payload
		}
	}
	
	// Return partial results if some operations succeeded
	if len(errors) > 0 {
		log.WithField("failed_operations", len(errors)).WithField("total_operations", len(secretRefs)).Warn("Some batch retrieve operations failed")
		
		// If all operations failed, return error
		if len(errors) == len(secretRefs) {
			return nil, fmt.Errorf("all batch retrieve operations failed: %v", errors)
		}
	}
	
	log.WithField("successful_operations", len(secretRefs)-len(errors)).WithField("total_operations", len(secretRefs)).Debug("Batch retrieve operations completed")
	return results, nil
}

// doRequestWithRetry performs an HTTP request with exponential backoff retry logic
func (c *BarbicanClient) doRequestWithRetry(ctx context.Context, method, path string, body []byte, result interface{}) error {
	var lastErr error
	
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// Calculate delay for this attempt (exponential backoff)
		if attempt > 0 {
			delay := time.Duration(float64(c.config.InitialRetryDelay) * math.Pow(c.config.RetryMultiplier, float64(attempt-1)))
			if delay > c.config.MaxRetryDelay {
				delay = c.config.MaxRetryDelay
			}
			
			log.WithField("attempt", attempt).WithField("delay", delay).Debug("Retrying Barbican API request")
			
			select {
			case <-ctx.Done():
				return NewTimeoutError("Request cancelled by context", ctx.Err())
			case <-time.After(delay):
				// Continue with retry
			}
		}

		err := c.doRequest(ctx, method, path, body, result)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryableError(err) {
			log.WithError(err).Debug("Non-retryable error, not retrying")
			break
		}

		log.WithError(err).WithField("attempt", attempt+1).Debug("Retryable error occurred")
	}

	return NewAPIError("Request failed after maximum retry attempts", 0).
		WithCause(lastErr).
		WithDetails(fmt.Sprintf("Failed after %d attempts", c.config.MaxRetries+1))
}

// doRequest performs a single HTTP request to the Barbican API
func (c *BarbicanClient) doRequest(ctx context.Context, method, path string, body []byte, result interface{}) error {
	// Get authentication token
	token, projectID, err := c.authManager.GetToken(ctx)
	if err != nil {
		return WrapError(err, ErrorTypeAuthentication, "Failed to get authentication token")
	}

	// Build full URL
	fullURL := c.endpoint + path

	// Create request
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return NewNetworkError("Failed to create HTTP request", err)
	}

	// Set headers
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("X-Project-Id", projectID)
	req.Header.Set("Accept", "application/json")
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Log request details (without sensitive data)
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	sanitizedData := securityValidator.SanitizeForLogging(map[string]interface{}{
		"method": method,
		"url":    fullURL,
	})
	log.WithFields(sanitizedData).Debug("Making Barbican API request")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewNetworkError("HTTP request failed", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewNetworkError("Failed to read response body", err)
	}

	// Check for authentication errors and invalidate token if needed
	if resp.StatusCode == http.StatusUnauthorized {
		log.Debug("Authentication failed, invalidating token")
		c.authManager.InvalidateToken()
		return NewAuthenticationError("Authentication failed").WithCode(resp.StatusCode)
	}

	// Handle different response types based on status code
	return c.handleResponse(resp.StatusCode, respBody, result)
}

// handleResponse handles HTTP response status codes and creates appropriate errors
func (c *BarbicanClient) handleResponse(statusCode int, respBody []byte, result interface{}) error {
	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		// Success - parse response if result is provided
		if result != nil {
			// Handle different result types
			switch v := result.(type) {
			case *[]byte:
				// For payload retrieval, return raw bytes
				*v = respBody
			default:
				// For JSON responses, unmarshal
				if len(respBody) > 0 {
					if err := json.Unmarshal(respBody, result); err != nil {
						return NewAPIError("Failed to unmarshal response", statusCode).WithCause(err)
					}
				}
			}
		}
		return nil

	case http.StatusNoContent:
		// Success with no content (e.g., DELETE operations)
		return nil

	case http.StatusNotFound:
		return NewSecretNotFoundError("").WithCode(statusCode).WithDetails(string(respBody))

	case http.StatusForbidden:
		return NewAuthorizationError("Access forbidden").WithCode(statusCode).WithDetails(string(respBody))

	case http.StatusBadRequest:
		// Try to parse error response
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			return NewValidationError(errorResp.Error.Message).WithCode(statusCode)
		}
		return NewValidationError("Bad request").WithCode(statusCode).WithDetails(string(respBody))

	case http.StatusTooManyRequests:
		return NewAPIError("Rate limit exceeded", statusCode).
			WithDetails(string(respBody)).
			WithSuggestions("Wait before retrying", "Reduce request frequency")

	case http.StatusInsufficientStorage, http.StatusRequestEntityTooLarge:
		return NewQuotaExceededError("Storage quota exceeded").WithCode(statusCode).WithDetails(string(respBody))

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		// Server errors - these are retryable
		return NewServiceUnavailableError("Server error").WithCode(statusCode).WithDetails(string(respBody))

	default:
		// Other errors
		return NewAPIError("Unexpected response status", statusCode).WithDetails(string(respBody))
	}
}

// GetBarbicanEndpoint discovers the Barbican endpoint from the service catalog
func GetBarbicanEndpoint(authManager *AuthManager, region string) (string, error) {
	return GetBarbicanEndpointForRegion(authManager, region)
}

// GetBarbicanEndpointForRegion discovers the Barbican endpoint for a specific region
func GetBarbicanEndpointForRegion(authManager *AuthManager, region string) (string, error) {
	if authManager == nil || authManager.config == nil {
		return "", fmt.Errorf("authentication manager is required")
	}

	authURL := authManager.config.AuthURL
	if authURL == "" {
		return "", fmt.Errorf("auth URL is required to construct Barbican endpoint")
	}

	// Parse the auth URL to get the base
	u, err := url.Parse(authURL)
	if err != nil {
		return "", fmt.Errorf("invalid auth URL: %w", err)
	}

	// Use region if provided, otherwise use default
	if region == "" {
		region = "RegionOne"
	}

	// Construct region-specific Barbican endpoint
	// In a real implementation, this would query the Keystone service catalog
	// For now, we'll construct a standard endpoint URL with region-specific hostname
	var barbicanURL string
	
	// Check if the hostname already includes region information
	hostname := u.Hostname()
	if strings.Contains(hostname, region) {
		// Hostname already region-specific
		barbicanURL = fmt.Sprintf("%s://%s:9311", u.Scheme, hostname)
	} else {
		// Construct region-specific hostname
		// Common patterns: region.service.domain or service-region.domain
		if strings.Contains(hostname, ".") {
			parts := strings.Split(hostname, ".")
			if len(parts) >= 2 {
				// Insert region as subdomain: keystone.example.com -> barbican-region.example.com
				regionHostname := fmt.Sprintf("barbican-%s.%s", region, strings.Join(parts[1:], "."))
				barbicanURL = fmt.Sprintf("%s://%s:9311", u.Scheme, regionHostname)
			} else {
				// Fallback to simple region prefix
				barbicanURL = fmt.Sprintf("%s://%s-%s:9311", u.Scheme, hostname, region)
			}
		} else {
			// Simple hostname, add region suffix
			barbicanURL = fmt.Sprintf("%s://%s-%s:9311", u.Scheme, hostname, region)
		}
	}
	
	log.WithField("endpoint", barbicanURL).WithField("region", region).Debug("Constructed region-specific Barbican endpoint")
	return barbicanURL, nil
}

// GetMultiRegionEndpoints returns Barbican endpoints for multiple regions
func GetMultiRegionEndpoints(authManager *AuthManager, regions []string) (map[string]string, error) {
	if len(regions) == 0 {
		return nil, fmt.Errorf("no regions specified")
	}

	endpoints := make(map[string]string)
	var errors []error

	for _, region := range regions {
		endpoint, err := GetBarbicanEndpointForRegion(authManager, region)
		if err != nil {
			errors = append(errors, fmt.Errorf("region %s: %w", region, err))
			continue
		}
		endpoints[region] = endpoint
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("failed to get endpoints for any region: %v", errors)
	}

	if len(errors) > 0 {
		log.WithField("failed_regions", len(errors)).WithField("successful_regions", len(endpoints)).Warn("Some regions failed during endpoint discovery")
	}

	return endpoints, nil
}

// ConnectionPoolStats represents connection pool statistics
type ConnectionPoolStats struct {
	MaxIdleConns        int `json:"max_idle_conns"`
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host"`
	MaxConnsPerHost     int `json:"max_conns_per_host"`
	IdleConnTimeout     int `json:"idle_conn_timeout_seconds"`
	ActiveConnections   int `json:"active_connections"`
	IdleConnections     int `json:"idle_connections"`
}

// GetConnectionPoolStats returns current connection pool statistics
func (c *BarbicanClient) GetConnectionPoolStats() *ConnectionPoolStats {
	// In a real implementation, this would extract actual statistics from the HTTP transport
	// For now, return configuration values as a baseline
	return &ConnectionPoolStats{
		MaxIdleConns:        c.config.MaxIdleConns,
		MaxIdleConnsPerHost: c.config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     c.config.MaxConnsPerHost,
		IdleConnTimeout:     int(c.config.IdleConnTimeout.Seconds()),
		ActiveConnections:   0, // Would be populated from actual transport stats
		IdleConnections:     0, // Would be populated from actual transport stats
	}
}

// OptimizeConnectionPool analyzes usage patterns and optimizes connection pool settings
func (c *BarbicanClient) OptimizeConnectionPool(stats *ConnectionPoolStats) *ClientConfig {
	optimized := *c.config // Copy current config
	
	// Calculate utilization rates
	if stats.MaxConnsPerHost > 0 {
		utilizationRate := float64(stats.ActiveConnections) / float64(stats.MaxConnsPerHost)
		
		// If utilization is high, increase pool sizes
		if utilizationRate > 0.8 {
			optimized.MaxConnsPerHost = min(optimized.MaxConnsPerHost*2, 200)
			optimized.MaxIdleConnsPerHost = min(optimized.MaxIdleConnsPerHost*2, 50)
			optimized.MaxIdleConns = min(optimized.MaxIdleConns*2, 400)
			
			log.WithField("utilization_rate", utilizationRate).Info("Increased connection pool sizes due to high utilization")
		}
		
		// If utilization is very low, decrease pool sizes to save resources
		if utilizationRate < 0.2 && optimized.MaxConnsPerHost > 10 {
			optimized.MaxConnsPerHost = max(optimized.MaxConnsPerHost/2, 10)
			optimized.MaxIdleConnsPerHost = max(optimized.MaxIdleConnsPerHost/2, 2)
			optimized.MaxIdleConns = max(optimized.MaxIdleConns/2, 20)
			
			log.WithField("utilization_rate", utilizationRate).Info("Decreased connection pool sizes due to low utilization")
		}
	}
	
	return &optimized
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CreateOptimizedClient creates a new Barbican client with performance optimizations based on usage patterns
func CreateOptimizedClient(endpoint string, authManager *AuthManager, usagePattern string) (*BarbicanClient, error) {
	var config *ClientConfig
	
	switch usagePattern {
	case "high-throughput":
		config = HighPerformanceClientConfig()
	case "low-latency":
		config = LowLatencyClientConfig()
	case "multi-region":
		config = MultiRegionClientConfig()
	default:
		config = DefaultClientConfig()
	}
	
	return NewBarbicanClient(endpoint, authManager, config)
}