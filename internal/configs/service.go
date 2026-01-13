package configs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Service struct {
	repo     *Repository
	sentinel *SentinelClient
}

func NewService(repo *Repository, sentinel *SentinelClient) *Service {
	return &Service{repo: repo, sentinel: sentinel}
}

func (s *Service) CreateConfig(projectID, envID int, key string, data, schema Map, userID int) (*Config, error) {
	// 1. Convert Map to Schema struct for internal validation
	schemaBytes, _ := json.Marshal(schema)
	var schemaStruct Schema
	if err := json.Unmarshal(schemaBytes, &schemaStruct); err != nil {
		return nil, fmt.Errorf("invalid schema format: %w", err)
	}

	// 2. Perform internal validation
	if err := Validate(schemaStruct, data); err != nil {
		return nil, fmt.Errorf("local validation failed: %w", err)
	}

	// 3. Optional: Deep Linting with Sentinel
	if s.sentinel != nil && s.sentinel.BaseURL != "" {
		valid, errs, err := s.sentinel.Lint(schemaStruct, data)
		if err != nil {
			// We log but don't necessarily block if Sentinel is down, 
			// depending on how strict we want to be. 
			// For "premium" feel, we might want to block or at least flag it.
			fmt.Printf("Sentinel linting error (skipped): %v\n", err)
		} else if !valid {
			return nil, &ValidationError{Errors: append([]string{"Sentinel deep linting failed"}, errs...)}
		}
	}

	return s.repo.CreateOrUpdate(projectID, envID, key, data, schema, userID)
}

func (s *Service) GetConfig(projectID, envID int, key string) (*Config, error) {
	return s.repo.GetLatest(projectID, envID, key)
}

func (s *Service) RollbackConfig(projectID, envID int, key string, targetVersion int, userID int) (*Config, error) {
	return s.repo.Rollback(projectID, envID, key, targetVersion, userID)
}

func (s *Service) FetchExternal(url string) (map[string]interface{}, error) {
	// Auto-transform Gist URLs to raw versions
	if strings.Contains(url, "gist.github.com") && !strings.Contains(url, "/raw") {
		url = url + "/raw"
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch from source: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode source data (expected JSON): %w", err)
	}

	return data, nil
}

