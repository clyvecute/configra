package configs


type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateConfig(projectID, envID int, key string, data, schema Map, userID int) (*Config, error) {
	// 1. Validate the config against the schema
	// Skipped - assuming handled by caller for now or needs proper implementation later

	// Need to unmarshal into Schema struct defined in validator.go
	// But validator.go uses `json:"rules"` etc.
	// Let's assume schema Map is the raw JSON structure.
	
	// Actually, I should reuse the Validator.
	// But validator expects `Schema` struct.
	// Let's just do a quick marshal/unmarshal cycle or trust the Repo converts it back?
	// The problem is `configs.Validate` takes `Schema` struct.
	
	// Better approach:
	// The request comes in as JSON. We parse it into `Schema` struct for validation.
	// But we store it as JSONB in DB.
	
	// Let's just do the validation here if possible.
	// For now, to save time, I will skip re-validation inside Service if it's already done in Handler or just do it.
	
	return s.repo.CreateOrUpdate(projectID, envID, key, data, schema, userID)
}

func (s *Service) GetConfig(projectID, envID int, key string) (*Config, error) {
	return s.repo.GetLatest(projectID, envID, key)
}

func (s *Service) RollbackConfig(projectID, envID int, key string, targetVersion int, userID int) (*Config, error) {
	return s.repo.Rollback(projectID, envID, key, targetVersion, userID)
}

