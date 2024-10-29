// sicclient.go
package sicclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/sicclient/sicparams"
	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/sicclient/sicresponse"
)

// Config holds the configuration needed to initialize the SIC API client.
type Config struct {
	Hostname    string      // Hostname for the SIC API
	TokenConfig TokenConfig // Config for token generation and management
}

// Client represents the main client that interacts with SIC APIs.
type Client struct {
	hostname     string
	tokenManager *TokenManager
}

// New initializes a new SIC Client with the given hostname and token manager.
func New(config Config) *Client {
	tokenManager := NewTokenManager(config.TokenConfig)
	return &Client{
		hostname:     config.Hostname,
		tokenManager: tokenManager,
	}
}

// GetUsageByEntity fetches usage data by entity within the specified time range,
// applying optional filters, sorting, and pagination.
func (c *Client) GetUsageByEntity(startTime, endTime string, parameters *sicparams.Params) (*sicresponse.UsageByEntityResponse, error) {
	apiURL := fmt.Sprintf("https://%s/sustainability-insight-ctr/v1beta1/usage-by-entity?start-time=%s&end-time=%s",
		c.hostname, url.QueryEscape(startTime), url.QueryEscape(endTime))

	if parameters != nil {
		apiURL = fmt.Sprintf("%s&%s", apiURL, parameters.ToQueryParams().Encode())
	}

	var usageResponse sicresponse.UsageByEntityResponse
	if err := c.doGetRequest(apiURL, &usageResponse); err != nil {
		return nil, err
	}
	return &usageResponse, nil
}

// GetUsageSeries fetches usage data over a time series with specified intervals.
// It applies optional filters, sorting, and pagination.
func (c *Client) GetUsageSeries(startTime, endTime, interval string, parameters *sicparams.Params) (*sicresponse.UsageSeriesResponse, error) {
	apiURL := fmt.Sprintf("https://%s/sustainability-insight-ctr/v1beta1/usage-series?start-time=%s&end-time=%s&interval=%s",
		c.hostname, url.QueryEscape(startTime), url.QueryEscape(endTime), url.QueryEscape(interval))

	if parameters != nil {
		apiURL = fmt.Sprintf("%s&%s", apiURL, parameters.ToQueryParams().Encode())
	}

	var usageSeriesResponse sicresponse.UsageSeriesResponse
	if err := c.doGetRequest(apiURL, &usageSeriesResponse); err != nil {
		return nil, err
	}
	return &usageSeriesResponse, nil
}

// doGetRequest is a helper method that performs a GET request and decodes the response into the provided response struct.
func (c *Client) doGetRequest(apiURL string, response interface{}) error {
	// Get the token from the TokenManager
	token, err := c.tokenManager.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set Authorization header with the token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body) // Read body for better error handling
		if readErr != nil {
			return fmt.Errorf("API returned status %s, failed to read body: %v", resp.Status, readErr)
		}
		return fmt.Errorf("API returned status %s: %s", resp.Status, string(body))
	}

	// Unmarshal the response body into the provided response struct
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return nil
}
