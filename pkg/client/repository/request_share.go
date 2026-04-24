package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// loadAccessTokenFromFiles: search token files in priority order (admin → app)
func loadAccessTokenFromFilesRequestHelper() string { // renamed to avoid duplicate
	dirs := []string{
		filepath.Join("etc", ".cmn", "client", "admin", "access_token"),
		filepath.Join("etc", ".cmn", "client", "app", "access_token"),
	}
	for _, f := range dirs {
		b, err := os.ReadFile(f)
		if err == nil && len(b) > 0 {
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}

// SendRequest is a helper function to make HTTP requests and handle responses.
// It abstracts away the boilerplate code for making requests, handling JSON, and decoding responses.
func SendRequest(method, endpoint string, requestBody interface{}, response interface{}) error {
	var req *http.Request
	var err error

	// Marshal the request body if it's not nil
	var jsonBody []byte
	if requestBody != nil {
		jsonBody, err = json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
	}

	// Create the HTTP request
	req, err = http.NewRequest(method, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header if token exists (for internal/private APIs)
	token := os.Getenv("CMN_ACCESS_TOKEN")
	if token == "" { // fallback to files
		token = loadAccessTokenFromFilesRequestHelper()
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status %s: %s", resp.Status, string(respBody))
	}

	// Decode the response body into the provided response struct
	if err := json.Unmarshal(respBody, &response); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}

	return nil
}
