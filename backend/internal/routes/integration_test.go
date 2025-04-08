package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Integration tests for the HTTP endpoints
// Note: These tests require a running server to test against

// getBaseURL returns the base URL for the API server
// This allows testing against a locally running server or a deployed one
func getBaseURL() string {
	// Check if TEST_API_URL is set in environment variables
	url := os.Getenv("TEST_API_URL")
	if url == "" {
		// Default to localhost if not specified
		url = "http://localhost:8080"
	}
	return url
}

func TestGetEndpointsIntegration(t *testing.T) {
	// Skip if the SKIP_INTEGRATION_TESTS environment variable is set
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	baseURL := getBaseURL()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	testCases := []struct {
		name           string
		endpoint       string
		expectedStatus int
		isArray        bool // Flag to indicate if the response is expected to be an array
	}{
		{
			name:           "Get Forecast by ID",
			endpoint:       "/forecasts/4",
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:           "Get Forecast Points by ID",
			endpoint:       "/forecast-points/4",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get All Forecast Points",
			endpoint:       "/forecast-points",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get Latest Forecast Points",
			endpoint:       "/forecast-points/latest",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get Latest Forecast Points By User",
			endpoint:       "/forecast-points/latest_by_user?user_id=2",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get Ordered Forecast Points",
			endpoint:       "/forecast-points/ordered/4",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get All Scores",
			endpoint:       "/scores/all",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get Average Scores",
			endpoint:       "/scores/average",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get Average Score By Forecast ID",
			endpoint:       "/scores/average/4",
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:           "List Users",
			endpoint:       "/users",
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s%s", baseURL, tc.endpoint)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error making request to %s: %v", tc.endpoint, err)
			}
			defer resp.Body.Close()

			t.Logf("Response status code: %d for endpoint %s", resp.StatusCode, tc.endpoint)

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d for endpoint %s",
					tc.expectedStatus, resp.StatusCode, tc.endpoint)
			}

			if resp.StatusCode != http.StatusOK {
				t.Logf("Response body: %s", resp.Body)
			}

			// Verify we got valid JSON back (if status is OK)
			if resp.StatusCode == http.StatusOK {
				// Read the entire body
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Failed to read response body: %v", err)
					return
				}

				// Debugging: Print first few bytes to see what we're dealing with
				if len(bodyBytes) > 0 {

					// Try both approaches - direct decoding as array or object
					var arrayResult []interface{}
					arrayErr := json.Unmarshal(bodyBytes, &arrayResult)

					var mapResult map[string]interface{}
					mapErr := json.Unmarshal(bodyBytes, &mapResult)

					// See which one worked
					if arrayErr == nil && mapErr != nil {
						t.Logf("Successfully decoded as array with %d elements", len(arrayResult))
						if !tc.isArray {
							t.Errorf("Expected object response but got array for endpoint %s", tc.endpoint)
						}
					} else if mapErr == nil && arrayErr != nil {
						t.Logf("Successfully decoded as object with %d fields", len(mapResult))
						if tc.isArray {
							t.Errorf("Expected array response but got object for endpoint %s", tc.endpoint)
						}
					} else if arrayErr == nil && mapErr == nil {
						// This is unlikely but could happen with "{}" or "[]"
						t.Logf("Response could be decoded as both array and object")
					} else {
						t.Errorf("Failed to decode as either array or object. Array error: %v, Object error: %v",
							arrayErr, mapErr)
					}
				} else {
					t.Logf("Empty response body from endpoint %s", tc.endpoint)
				}
			}
		})
	}
}

func TestPostEndpointsIntegration(t *testing.T) {
	// Skip if the SKIP_INTEGRATION_TESTS environment variable is set
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	baseURL := getBaseURL()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	testCases := []struct {
		name           string
		endpoint       string
		body           string
		expectedStatus int
		isArray        bool // Flag to indicate if the response is expected to be an array
	}{
		// Forecast endpoints
		{
			name:     "List forecasts",
			endpoint: "/forecasts",
			body: `{
				"list_type": "open",
				"category": "politics"
			}`,
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:     "List forecasts fail",
			endpoint: "/forecasts",
			body: `{
				"list_type": "xddd",
				"category": "none"
			}`,
			expectedStatus: http.StatusBadRequest,
			isArray:        false,
		},
		// Single-score endpoints
		{
			name:     "Get scores (user_id and forecast_id)",
			endpoint: "/scores",
			body: `{
				"user_id": 2,
				"forecast_id": 4
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:     "Get scores (user_id)",
			endpoint: "/scores",
			body: `{
				"user_id": 2
			}`,
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:     "Get scores (forecast_id)",
			endpoint: "/scores",
			body: `{
				"forecast_id": 4
			}`,
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get scores (no body)",
			endpoint:       "/scores",
			body:           `{}`,
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		// Aggregate score endpoints
		{
			name:     "Get aggregate scores",
			endpoint: "/scores/aggregate",
			body: `{
				"category": "politics"
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:     "Get aggregate scores (user_id)",
			endpoint: "/scores/aggregate",
			body: `{
				"category": "politics",
				"user_id": 2
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:           "Get aggregate scores (no body)",
			endpoint:       "/scores/aggregate",
			body:           `{}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:     "Get aggregate scores fail (user_id not existing)",
			endpoint: "/scores/aggregate",
			body: `{
				"category": "politics",
				"user_id": 999
			}`,
			expectedStatus: http.StatusInternalServerError,
			isArray:        false,
		},
		{
			name:     "Get aggregate scores by user",
			endpoint: "/scores/aggregate",
			body: `{
				"category": "economy",
				"by_user": true
			}`,
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		{
			name:           "Get aggregate scores by user (no body)",
			endpoint:       "/scores/aggregate",
			body:           `{"by_user": true}`,
			expectedStatus: http.StatusOK,
			isArray:        true,
		},
		// User endpoints (non-protected)
		{
			name:     "Create user",
			endpoint: "/users",
			body: `{
				"username": "testuser",
				"password": "testpassword"
			}`,
			expectedStatus: http.StatusCreated,
			isArray:        false,
		},
		{
			name:           "Create user fail",
			endpoint:       "/users",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
			isArray:        false,
		},
		{
			name:     "Login user",
			endpoint: "/users/login",
			body: `{
				"username": "testuser",
				"password": "testpassword"
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:     "Reset password",
			endpoint: "/users/reset-password",
			body: `{
				"user_id": 1,
				"new_password": "testing123"
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:     "Login after reset",
			endpoint: "/users/login",
			body: `{
				"username": "testuser",
				"password": "testing123"
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:           "Login user fail",
			endpoint:       "/users/login",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
			isArray:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s%s", baseURL, tc.endpoint)

			req, err := http.NewRequest("POST", url, strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error making request to %s: %v", tc.endpoint, err)
			}
			defer resp.Body.Close()

			t.Logf("Response status code: %d for endpoint %s", resp.StatusCode, tc.endpoint)

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d for endpoint %s",
					tc.expectedStatus, resp.StatusCode, tc.endpoint)
			}

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				t.Logf("Response body: %s", string(bodyBytes))
			}

			// Verify we got valid JSON back (if status is OK)
			if resp.StatusCode == http.StatusOK {
				// Read the entire body
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Failed to read response body: %v", err)
					return
				}

				// Debugging: Print first few bytes to see what we're dealing with
				if len(bodyBytes) > 0 {

					// Try both approaches - direct decoding as array or object
					var arrayResult []interface{}
					arrayErr := json.Unmarshal(bodyBytes, &arrayResult)

					var mapResult map[string]interface{}
					mapErr := json.Unmarshal(bodyBytes, &mapResult)

					// See which one worked
					if arrayErr == nil && mapErr != nil {
						t.Logf("Successfully decoded as array with %d elements", len(arrayResult))
						if !tc.isArray {
							t.Errorf("Expected object response but got array for endpoint %s", tc.endpoint)
						}
					} else if mapErr == nil && arrayErr != nil {
						t.Logf("Successfully decoded as object with %d fields", len(mapResult))
						if tc.isArray {
							t.Errorf("Expected array response but got object for endpoint %s", tc.endpoint)
						}
					} else if arrayErr == nil && mapErr == nil {
						// This is unlikely but could happen with "{}" or "[]"
						t.Logf("Response could be decoded as both array and object")
					} else {
						t.Errorf("Failed to decode as either array or object. Array error: %v, Object error: %v",
							arrayErr, mapErr)
					}
				} else {
					t.Logf("Empty response body from endpoint %s", tc.endpoint)
				}
			}
		})
	}
}

// Add this helper function to get a JWT token
func getAuthToken(t *testing.T) string {
	baseURL := getBaseURL()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	test_username := os.Getenv("TEST_USERNAME")
	test_password := os.Getenv("TEST_PASSWORD")
	if test_username == "" || test_password == "" {
		t.Fatalf("TEST_USERNAME and TEST_PASSWORD must be set in the environment")
	}

	loginBody := `{
		"username": "` + test_username + `",
		"password": "` + test_password + `"
	}`

	// Make login request
	req, err := http.NewRequest("POST", baseURL+"/users/login",
		strings.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Error creating login request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error logging in: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to login, status code: %d", resp.StatusCode)
	}

	// Parse response to get token
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("Error decoding login response: %v", err)
	}

	return loginResp.Token
}

func TestProtectedEndpointsIntegration(t *testing.T) {
	// Skip if the SKIP_INTEGRATION_TESTS environment variable is set
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	baseURL := getBaseURL()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Get token for authentication
	token := getAuthToken(t)

	testCases := []struct {
		name           string
		method         string
		endpoint       string
		body           string
		expectedStatus int
		isArray        bool
	}{
		// Forecast endpoints
		{
			name:     "Create forecast",
			method:   "POST",
			endpoint: "/api/forecasts/create",
			body: `{
				"question": "Will this test pass?",
				"category": "testing",
				"resolution_criteria": "the test will pass"
			}`,
			expectedStatus: http.StatusCreated,
			isArray:        false,
		},
		{
			name:           "Create forecast fail (no body)",
			method:         "POST",
			endpoint:       "/api/forecasts/create",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
			isArray:        false,
		},
		// Forecast point endpoints
		{
			name:     "Create forecast point",
			method:   "POST",
			endpoint: "/api/forecast-points",
			body: `{
				"forecast_id": 30,
				"point_forecast": 0.65,
				"reason": "the test will pass"
			}`,
			expectedStatus: http.StatusCreated,
			isArray:        false,
		},
		{
			name:     "Create forecast point fail (forecast not found)",
			method:   "POST",
			endpoint: "/api/forecast-points",
			body: `{
				"forecast_id": 500,
				"point_forecast": 0.65,
				"user_id": 2,
				"reason": "the test will pass"
			}`,
			expectedStatus: http.StatusBadRequest,
			isArray:        false,
		},
		{
			name:     "Resolve forecast",
			method:   "PUT",
			endpoint: "/api/resolve",
			body: `{
				"id": 2,
				"resolution": "1",
				"comment": "the test will pass",
				"user_id": 2
			}`,
			expectedStatus: http.StatusOK,
			isArray:        false,
		},
		{
			name:     "Resolve forecast fail (forecast already resolved)",
			method:   "PUT",
			endpoint: "/api/resolve",
			body: `{
				"id": 188,
				"resolution": "1",
				"comment": "the test will pass",
				"user_id": 2
			}`,
			expectedStatus: http.StatusBadRequest,
			isArray:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s%s", baseURL, tc.endpoint)

			req, err := http.NewRequest(tc.method, url, strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			// Add JWT token to the request
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error making request to %s: %v", tc.endpoint, err)
			}
			defer resp.Body.Close()

			t.Logf("Response status code: %d for endpoint %s", resp.StatusCode, tc.endpoint)

			if resp.StatusCode != tc.expectedStatus {
				bodyBytes, _ := io.ReadAll(resp.Body)
				t.Errorf("Expected status %d, got %d for endpoint %s. Response: %s",
					tc.expectedStatus, resp.StatusCode, tc.endpoint, string(bodyBytes))
				return
			}

			// Verify we got valid JSON back (if status is OK)
			if resp.StatusCode == http.StatusOK {
				// Read the entire body
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Failed to read response body: %v", err)
					return
				}

				// Skip empty responses
				if len(bodyBytes) == 0 {
					t.Logf("Empty response body from endpoint %s", tc.endpoint)
					return
				}

				// Try both approaches - direct decoding as array or object
				var arrayResult []interface{}
				arrayErr := json.Unmarshal(bodyBytes, &arrayResult)

				var mapResult map[string]interface{}
				mapErr := json.Unmarshal(bodyBytes, &mapResult)

				// See which one worked
				if arrayErr == nil && mapErr != nil {
					t.Logf("Successfully decoded as array with %d elements", len(arrayResult))
					if !tc.isArray {
						t.Errorf("Expected object response but got array for endpoint %s", tc.endpoint)
					}
				} else if mapErr == nil && arrayErr != nil {
					t.Logf("Successfully decoded as object with %d fields", len(mapResult))
					if tc.isArray {
						t.Errorf("Expected array response but got object for endpoint %s", tc.endpoint)
					}
				} else if arrayErr == nil && mapErr == nil {
					// This is unlikely but could happen with "{}" or "[]"
					t.Logf("Response could be decoded as both array and object")
				} else {
					t.Errorf("Failed to decode as either array or object. Array error: %v, Object error: %v",
						arrayErr, mapErr)
				}
			}
		})
	}
}
