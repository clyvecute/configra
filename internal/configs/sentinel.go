package configs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SentinelClient handles communication with the Sentinel Config Linter service.
type SentinelClient struct {
	BaseURL string
	Timeout time.Duration
}

// SentinelResponse matches the expected JSON response from the Sentinel server.
type SentinelResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
	Score  int      `json:"score,omitempty"` // For the "premium" system health visualizer
}

func NewSentinelClient(url string) *SentinelClient {
	return &SentinelClient{
		BaseURL: url,
		Timeout: 5 * time.Second,
	}
}

// Lint sends the configuration and schema to the Sentinel service for deep linting.
func (c *SentinelClient) Lint(schema Schema, config map[string]interface{}) (bool, []string, error) {
	if c.BaseURL == "" {
		return true, nil, nil // Skip if not configured
	}

	payload := map[string]interface{}{
		"schema": schema,
		"config": config,
	}
	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Post(fmt.Sprintf("%s/v1/lint", c.BaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false, nil, fmt.Errorf("sentinel unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("sentinel returned status: %d", resp.StatusCode)
	}

	var result SentinelResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, nil, fmt.Errorf("failed to decode sentinel response: %w", err)
	}

	return result.Valid, result.Errors, nil
}
