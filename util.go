// pkg/util.go

package simpleturso

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Basic JWT structure validation (no signature verification)
func isValidJWTStructure(token string) bool {
	if token == "" {
		return false
	}

	// JWT should have 3 parts separated by dots: header.payload.signature
	parts := strings.Split(token, ".")
	return len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != ""
}

// Basic Turso URL validation
func isValidTursoUrl(url string) bool {
	if url == "" {
		return false
	}

	parts := strings.Split(url, "://")
	if len(parts) != 2 {
		return false
	}

	prefix := parts[0]
	if prefix != "https" && prefix != "libsql" {
		return false
	}

	if !strings.HasSuffix(url, ".turso.io") {
		return false
	}

	return true
}

// executeTursoRequest sends a request to Turso and returns an error if it fails
func executeTursoRequestWithResponse(req TursoRequest, cfg *Config) (*TursoResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", cfg.DbUrl, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+cfg.DbKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the body once
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("turso returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode from the bytes we just read
	var tursoResp TursoResponse
	if err := json.Unmarshal(bodyBytes, &tursoResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w (Body: %s)", err, string(bodyBytes))
	}

	return &tursoResp, nil
}
