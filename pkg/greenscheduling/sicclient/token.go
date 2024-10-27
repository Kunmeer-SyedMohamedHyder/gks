package sicclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Constants for token management
const (
	grantType         = "client_credentials"
	contentTypeHeader = "application/x-www-form-urlencoded"
)

// TokenResponse represents the structure of the token response.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"` // Include token type (e.g., "Bearer")
	ExpiresIn   int    `json:"expires_in"` // The token's lifetime in seconds
}

// TokenConfig holds the configuration for the token generation.
type TokenConfig struct {
	URL          string // The URL for the token endpoint
	ClientID     string // The client ID
	ClientSecret string // The client secret
}

// TokenInfo holds the access token and its expiration details.
type TokenInfo struct {
	Token  string    // Current access token
	Expiry time.Time // Expiration time of the current token
}

// TokenManager is responsible for managing the access token.
type TokenManager struct {
	config     TokenConfig  // Configuration for token management
	httpClient *http.Client // HTTP client for making requests
	tokenInfo  TokenInfo    // Token information including current token and expiry
}

// NewTokenManager creates a new instance of TokenManager.
func NewTokenManager(config TokenConfig) *TokenManager {
	return &TokenManager{
		config:     config,
		httpClient: &http.Client{},
	}
}

// GetToken returns the current access token and its type, generating a new one if necessary.
func (tm *TokenManager) GetToken() (string, error) {
	if time.Now().After(tm.tokenInfo.Expiry) { // Check if token is expired
		if err := tm.RefreshToken(); err != nil {
			return "", err
		}
	}
	return tm.tokenInfo.Token, nil // Return current access token
}

// RefreshToken refreshes the access token.
func (tm *TokenManager) RefreshToken() error {
	return tm.GenerateToken() // Reuse GenerateToken logic to get a new token
}

// GenerateToken generates a new access token.
func (tm *TokenManager) GenerateToken() error {
	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("client_id", tm.config.ClientID)
	data.Set("client_secret", tm.config.ClientSecret)

	dataEncoded := data.Encode() // Create the request body

	req, err := http.NewRequest("POST", tm.config.URL, bytes.NewBufferString(dataEncoded))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentTypeHeader) // Set content type header

	resp, err := tm.httpClient.Do(req) // Execute the request
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to generate token: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	// Update token and expiry time
	tm.tokenInfo.Token = tokenResp.AccessToken
	tm.tokenInfo.Expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}
