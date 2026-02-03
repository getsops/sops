// Package barbican provides OpenStack Barbican authentication and token management.
package barbican

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// AuthManager handles OpenStack authentication using various methods including
// password authentication, application credentials, and token authentication.
// It manages token caching and automatic renewal.
type AuthManager struct {
	config     *AuthConfig
	httpClient *http.Client
	tokenCache *TokenCache
}

// TokenCache manages authentication token lifecycle with thread-safe access.
// It automatically handles token expiration and renewal.
type TokenCache struct {
	token     string
	expiry    time.Time
	projectID string
	mutex     sync.RWMutex
}

// AuthResponse represents the response from Keystone authentication
type AuthResponse struct {
	Token struct {
		ExpiresAt string `json:"expires_at"`
		Project   struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"token"`
}

// AuthRequest represents different authentication request formats
type AuthRequest struct {
	Auth AuthMethod `json:"auth"`
}

// AuthMethod represents the authentication method interface
type AuthMethod struct {
	Identity Identity `json:"identity"`
	Scope    *Scope   `json:"scope,omitempty"`
}

// Identity contains the authentication credentials
type Identity struct {
	Methods             []string             `json:"methods"`
	Password            *PasswordAuth        `json:"password,omitempty"`
	ApplicationCredential *ApplicationCredAuth `json:"application_credential,omitempty"`
	Token               *TokenAuth           `json:"token,omitempty"`
}

// PasswordAuth represents password-based authentication
type PasswordAuth struct {
	User User `json:"user"`
}

// User represents user credentials
type User struct {
	Name     string  `json:"name,omitempty"`
	ID       string  `json:"id,omitempty"`
	Password string  `json:"password"`
	Domain   *Domain `json:"domain,omitempty"`
}

// ApplicationCredAuth represents application credential authentication
type ApplicationCredAuth struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// TokenAuth represents token-based authentication
type TokenAuth struct {
	ID string `json:"id"`
}

// Scope represents the authentication scope
type Scope struct {
	Project *Project `json:"project,omitempty"`
}

// Project represents project scope
type Project struct {
	ID     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Domain *Domain `json:"domain,omitempty"`
}

// Domain represents domain information
type Domain struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(config *AuthConfig) (*AuthManager, error) {
	if config == nil {
		return nil, NewConfigError("Authentication configuration is required")
	}

	// Validate security configuration
	if err := ValidateSecurityConfiguration(config); err != nil {
		return nil, err
	}

	// Create HTTP client with TLS configuration
	httpClient, err := createHTTPClient(config)
	if err != nil {
		return nil, NewTLSError("Failed to create HTTP client", err)
	}

	return &AuthManager{
		config:     config,
		httpClient: httpClient,
		tokenCache: &TokenCache{},
	}, nil
}

// createHTTPClient creates an HTTP client with proper TLS configuration
func createHTTPClient(config *AuthConfig) (*http.Client, error) {
	// Create security validator and TLS config
	securityConfig := SecurityConfigFromAuthConfig(config)
	validator := NewSecurityValidator(securityConfig)
	
	tlsConfig, err := validator.ValidateAndCreateTLSConfig()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		MaxIdleConnsPerHost: 5,
		TLSClientConfig:     tlsConfig,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}

// GetToken retrieves a valid authentication token, using cache if available
func (am *AuthManager) GetToken(ctx context.Context) (string, string, error) {
	// Check if we have a valid cached token
	am.tokenCache.mutex.RLock()
	if am.tokenCache.token != "" && time.Now().Before(am.tokenCache.expiry) {
		token := am.tokenCache.token
		projectID := am.tokenCache.projectID
		am.tokenCache.mutex.RUnlock()
		log.Debug("Using cached authentication token")
		return token, projectID, nil
	}
	am.tokenCache.mutex.RUnlock()

	// Need to authenticate
	log.Debug("Authenticating with OpenStack Keystone")
	return am.authenticate(ctx)
}

// authenticate performs the actual authentication with Keystone
func (am *AuthManager) authenticate(ctx context.Context) (string, string, error) {
	// Determine authentication method and build request
	authReq, err := am.buildAuthRequest()
	if err != nil {
		return "", "", NewAuthenticationError("Failed to build authentication request").WithCause(err)
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return "", "", NewAuthenticationError("Failed to marshal authentication request").WithCause(err)
	}

	// Build authentication URL
	authURL := strings.TrimSuffix(am.config.AuthURL, "/") + "/auth/tokens"

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", "", NewNetworkError("Failed to create authentication request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Log request details (without sensitive data)
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	sanitizedData := securityValidator.SanitizeForLogging(map[string]interface{}{
		"auth_url": authURL,
		"method":   "POST",
	})
	log.WithFields(sanitizedData).Debug("Sending authentication request")

	// Send request
	resp, err := am.httpClient.Do(req)
	if err != nil {
		return "", "", NewNetworkError("authentication request failed", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", NewNetworkError("Failed to read authentication response", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return "", "", am.handleAuthenticationError(resp.StatusCode, body)
	}

	// Extract token from header
	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return "", "", NewAuthenticationError("No authentication token received in response")
	}

	// Parse response to get project ID and expiration
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", "", NewAuthenticationError("failed to parse authentication response").WithCause(err)
	}

	// Parse expiration time
	expiry, err := time.Parse(time.RFC3339, authResp.Token.ExpiresAt)
	if err != nil {
		log.WithField("expires_at", authResp.Token.ExpiresAt).Warn("Failed to parse token expiration, using default")
		expiry = time.Now().Add(1 * time.Hour) // Default 1 hour expiration
	}

	projectID := authResp.Token.Project.ID

	// Cache the token
	am.tokenCache.mutex.Lock()
	am.tokenCache.token = token
	am.tokenCache.expiry = expiry.Add(-5 * time.Minute) // Refresh 5 minutes before expiry
	am.tokenCache.projectID = projectID
	am.tokenCache.mutex.Unlock()

	log.WithField("expires_at", expiry).Debug("Authentication successful, token cached")

	return token, projectID, nil
}

// buildAuthRequest builds the appropriate authentication request based on configuration
func (am *AuthManager) buildAuthRequest() (*AuthRequest, error) {
	config := am.config

	// Determine authentication method priority:
	// 1. Application Credentials (most secure)
	// 2. Token (for re-authentication)
	// 3. Password (fallback)

	var identity Identity
	var scope *Scope

	if config.ApplicationCredentialID != "" && config.ApplicationCredentialSecret != "" {
		// Application Credential authentication
		identity = Identity{
			Methods: []string{"application_credential"},
			ApplicationCredential: &ApplicationCredAuth{
				ID:     config.ApplicationCredentialID,
				Secret: config.ApplicationCredentialSecret,
			},
		}
		log.Debug("Using application credential authentication")
	} else if config.Token != "" {
		// Token authentication
		identity = Identity{
			Methods: []string{"token"},
			Token: &TokenAuth{
				ID: config.Token,
			},
		}
		log.Debug("Using token authentication")
	} else if config.Username != "" && config.Password != "" {
		// Password authentication
		user := User{
			Password: config.Password,
		}

		// Set user identifier (prefer ID over name)
		if config.Username != "" {
			user.Name = config.Username
		}

		// Set domain for user
		if config.DomainID != "" {
			user.Domain = &Domain{ID: config.DomainID}
		} else if config.DomainName != "" {
			user.Domain = &Domain{Name: config.DomainName}
		} else {
			user.Domain = &Domain{Name: "default"}
		}

		identity = Identity{
			Methods:  []string{"password"},
			Password: &PasswordAuth{User: user},
		}
		log.Debug("Using password authentication")
	} else {
		return nil, NewAuthenticationError("No valid authentication method found").
			WithSuggestions(
				"Provide application credentials (OS_APPLICATION_CREDENTIAL_ID/OS_APPLICATION_CREDENTIAL_SECRET)",
				"Provide token (OS_TOKEN)",
				"Provide username/password (OS_USERNAME/OS_PASSWORD)",
			)
	}

	// Set project scope if not using application credentials
	if config.ApplicationCredentialID == "" {
		project := &Project{}
		
		if config.ProjectID != "" {
			project.ID = config.ProjectID
		} else if config.ProjectName != "" {
			project.Name = config.ProjectName
			// Set domain for project
			if config.DomainID != "" {
				project.Domain = &Domain{ID: config.DomainID}
			} else if config.DomainName != "" {
				project.Domain = &Domain{Name: config.DomainName}
			} else {
				project.Domain = &Domain{Name: "default"}
			}
		} else {
			return nil, NewAuthenticationError("Project scope is required when not using application credentials").
				WithSuggestions(
					"Set project ID (OS_PROJECT_ID)",
					"Set project name (OS_PROJECT_NAME)",
				)
		}

		scope = &Scope{Project: project}
	}

	return &AuthRequest{
		Auth: AuthMethod{
			Identity: identity,
			Scope:    scope,
		},
	}, nil
}

// handleAuthenticationError creates appropriate error for authentication failures
func (am *AuthManager) handleAuthenticationError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return NewAuthenticationError(fmt.Sprintf("authentication failed with status %d", statusCode)).
			WithCode(statusCode).
			WithDetails(string(body)).
			WithSuggestions(
				"Verify your credentials are correct",
				"Check that your user account is not disabled",
				"Ensure you have access to the specified project",
			)
	case http.StatusForbidden:
		return NewAuthorizationError("Access denied").
			WithCode(statusCode).
			WithDetails(string(body)).
			WithSuggestions(
				"Verify your user has the required permissions",
				"Check that you're accessing the correct project/tenant",
				"Contact your OpenStack administrator",
			)
	case http.StatusBadRequest:
		return NewValidationError("Invalid authentication request").
			WithCode(statusCode).
			WithDetails(string(body)).
			WithSuggestions(
				"Check your authentication parameters",
				"Verify the authentication URL is correct",
				"Ensure all required fields are provided",
			)
	case http.StatusServiceUnavailable:
		return NewServiceUnavailableError("Authentication service is unavailable").
			WithCode(statusCode).
			WithDetails(string(body))
	default:
		return NewAPIError("Authentication failed", statusCode).
			WithDetails(string(body))
	}
}

// InvalidateToken clears the cached token, forcing re-authentication on next request
func (am *AuthManager) InvalidateToken() {
	am.tokenCache.mutex.Lock()
	defer am.tokenCache.mutex.Unlock()
	
	am.tokenCache.token = ""
	am.tokenCache.expiry = time.Time{}
	am.tokenCache.projectID = ""
	
	log.Debug("Authentication token invalidated")
}

// LoadConfigFromEnvironment loads authentication configuration from standard OpenStack environment variables
func LoadConfigFromEnvironment() *AuthConfig {
	config := &AuthConfig{
		AuthURL:                     os.Getenv("OS_AUTH_URL"),
		Region:                      os.Getenv("OS_REGION_NAME"),
		ProjectID:                   os.Getenv("OS_PROJECT_ID"),
		ProjectName:                 os.Getenv("OS_PROJECT_NAME"),
		DomainID:                    os.Getenv("OS_DOMAIN_ID"),
		DomainName:                  os.Getenv("OS_DOMAIN_NAME"),
		Username:                    os.Getenv("OS_USERNAME"),
		Password:                    os.Getenv("OS_PASSWORD"),
		ApplicationCredentialID:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
		Token:                       os.Getenv("OS_TOKEN"),
	}

	// Handle TLS configuration
	if os.Getenv("OS_INSECURE") == "true" || os.Getenv("OS_INSECURE") == "1" {
		config.Insecure = true
	}

	if caCert := os.Getenv("OS_CACERT"); caCert != "" {
		config.CACert = caCert
	}

	// Set default region if not specified
	if config.Region == "" {
		config.Region = "RegionOne"
	}

	// Set default domain if not specified and using password auth
	if config.DomainName == "" && config.DomainID == "" && config.Username != "" {
		config.DomainName = "default"
	}

	return config
}

// ValidateConfig validates the authentication configuration
func ValidateConfig(config *AuthConfig) error {
	if config == nil {
		return NewConfigError("Authentication configuration is required")
	}

	if config.AuthURL == "" {
		return NewConfigError("Authentication URL is required").
			WithSuggestions("Set OS_AUTH_URL environment variable")
	}

	// Validate authentication method
	hasAppCred := config.ApplicationCredentialID != "" && config.ApplicationCredentialSecret != ""
	hasToken := config.Token != ""
	hasPassword := config.Username != "" && config.Password != ""

	if !hasAppCred && !hasToken && !hasPassword {
		return NewAuthenticationError("No valid authentication method provided").
			WithSuggestions(
				"Set application credentials (OS_APPLICATION_CREDENTIAL_ID/OS_APPLICATION_CREDENTIAL_SECRET)",
				"Set token (OS_TOKEN)",
				"Set username/password (OS_USERNAME/OS_PASSWORD)",
			)
	}

	// Validate project scope for non-application credential auth
	if !hasAppCred {
		if config.ProjectID == "" && config.ProjectName == "" {
			return NewAuthenticationError("Project scope is required when not using application credentials").
				WithSuggestions(
					"Set OS_PROJECT_ID environment variable",
					"Set OS_PROJECT_NAME environment variable",
				)
		}
	}

	return nil
}

// AuthManagerStats represents authentication manager performance statistics
type AuthManagerStats struct {
	TokenCacheHits     int64         `json:"token_cache_hits"`
	TokenCacheMisses   int64         `json:"token_cache_misses"`
	AuthRequests       int64         `json:"auth_requests"`
	AuthFailures       int64         `json:"auth_failures"`
	AvgAuthTime        time.Duration `json:"avg_auth_time"`
	TokenRefreshCount  int64         `json:"token_refresh_count"`
	CachedTokenExpiry  time.Time     `json:"cached_token_expiry"`
}

// GetAuthManagerStats returns current authentication manager statistics
func (am *AuthManager) GetAuthManagerStats() *AuthManagerStats {
	am.tokenCache.mutex.RLock()
	defer am.tokenCache.mutex.RUnlock()
	
	// In a real implementation, these would be tracked over time
	return &AuthManagerStats{
		TokenCacheHits:     0, // Would be tracked in actual implementation
		TokenCacheMisses:   0, // Would be tracked in actual implementation
		AuthRequests:       0, // Would be tracked in actual implementation
		AuthFailures:       0, // Would be tracked in actual implementation
		AvgAuthTime:        0, // Would be calculated from timing data
		TokenRefreshCount:  0, // Would be tracked in actual implementation
		CachedTokenExpiry:  am.tokenCache.expiry,
	}
}

// OptimizeAuthConfig analyzes authentication patterns and suggests optimizations
func OptimizeAuthConfig(stats *AuthManagerStats, currentConfig *AuthConfig) *AuthConfig {
	optimized := *currentConfig // Copy current config
	
	// Analyze cache hit rate
	if stats.TokenCacheHits > 0 || stats.TokenCacheMisses > 0 {
		totalRequests := stats.TokenCacheHits + stats.TokenCacheMisses
		hitRate := float64(stats.TokenCacheHits) / float64(totalRequests)
		
		if hitRate < 0.8 { // Less than 80% cache hit rate
			log.WithField("cache_hit_rate", hitRate).Info("Low token cache hit rate detected")
			// In a real implementation, might suggest token refresh strategies
		}
	}
	
	// Analyze authentication failure rate
	if stats.AuthRequests > 0 {
		failureRate := float64(stats.AuthFailures) / float64(stats.AuthRequests)
		
		if failureRate > 0.1 { // More than 10% failure rate
			log.WithField("auth_failure_rate", failureRate).Warn("High authentication failure rate detected")
			// In a real implementation, might suggest credential validation or retry strategies
		}
	}
	
	return &optimized
}

// PrewarmTokenCache proactively authenticates to warm up the token cache
func (am *AuthManager) PrewarmTokenCache(ctx context.Context) error {
	log.Debug("Prewarming authentication token cache")
	_, _, err := am.GetToken(ctx)
	if err != nil {
		return WrapError(err, ErrorTypeAuthentication, "Failed to prewarm token cache")
	}
	
	log.Debug("Token cache prewarmed successfully")
	return nil
}

// RefreshTokenIfNeeded checks if the token needs refresh and refreshes it proactively
func (am *AuthManager) RefreshTokenIfNeeded(ctx context.Context, refreshThreshold time.Duration) error {
	am.tokenCache.mutex.RLock()
	needsRefresh := am.tokenCache.token != "" && time.Until(am.tokenCache.expiry) < refreshThreshold
	am.tokenCache.mutex.RUnlock()
	
	if needsRefresh {
		log.WithField("threshold", refreshThreshold).Debug("Proactively refreshing authentication token")
		_, _, err := am.authenticate(ctx)
		return err
	}
	
	return nil
}