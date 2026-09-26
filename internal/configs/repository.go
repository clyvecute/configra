package configs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Config struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	EnvID     int       `json:"env_id"`
	Key       string    `json:"key"`
	Version   int       `json:"version"` // Current version
	Data      Map       `json:"data"`    // Hydrated from version
	Schema    Map       `json:"schema"`  // Hydrated from version
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Map map[string]interface{}

func (m Map) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Map) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}



type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateOrUpdate handles the logic of creating a config key if it doesn't exist,
// and then appending a new version to it.
func (r *Repository) CreateOrUpdate(projectID, envID int, key string, data, schema Map, userID int) (*Config, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Get or Create Config Parent
	var configID int
	err = tx.QueryRow(`
		INSERT INTO configs (project_id, environment_id, key)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, environment_id, key) DO UPDATE 
			SET updated_at = NOW()
		RETURNING id`, projectID, envID, key).Scan(&configID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to upsert config parent: %v", err)
	}

	// 2. Get latest version number
	var currentVersion int
	err = tx.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM config_versions WHERE config_id = $1`, configID).Scan(&currentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get max version: %v", err)
	}

	newVersion := currentVersion + 1

	// 3. Insert new version
	schemaJSON, _ := json.Marshal(schema)
	dataJSON, _ := json.Marshal(data)

	_, err = tx.Exec(`
		INSERT INTO config_versions (config_id, version, data, schema, created_by)
		VALUES ($1, $2, $3, $4, $5)`,
		configID, newVersion, dataJSON, schemaJSON, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert version: %v", err)
	}

	// 4. Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Config{
		ID:        configID,
		ProjectID: projectID,
		EnvID:     envID,
		Key:       key,
		Version:   newVersion,
		Data:      data,
		Schema:    schema,
	}, nil
}

func (r *Repository) GetLatest(projectID, envID int, key string) (*Config, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}
	// Join configs and config_versions to get the latest data
	query := `
		SELECT c.id, c.updated_at, v.version, v.data, v.schema
		FROM configs c
		JOIN config_versions v ON c.id = v.config_id
		WHERE c.project_id = $1 AND c.environment_id = $2 AND c.key = $3
		ORDER BY v.version DESC
		LIMIT 1`
	
	row := r.db.QueryRow(query, projectID, envID, key)

	var c Config
	c.ProjectID = projectID
	c.EnvID = envID
	c.Key = key
	
	var dataBytes, schemaBytes []byte

	if err := row.Scan(&c.ID, &c.UpdatedAt, &c.Version, &dataBytes, &schemaBytes); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	json.Unmarshal(dataBytes, &c.Data)
	json.Unmarshal(schemaBytes, &c.Schema)

	return &c, nil
}
// Rollback finds a specific version of a config and creates a NEW version (latest + 1)
// with that old content. This preserves history (immutable).
func (r *Repository) Rollback(projectID, envID int, key string, targetVersion int, userID int) (*Config, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Find the Config ID
	var configID int
	err = tx.QueryRow(`
		SELECT id FROM configs 
		WHERE project_id = $1 AND environment_id = $2 AND key = $3`,
		projectID, envID, key).Scan(&configID)
	if err != nil {
		return nil, fmt.Errorf("config not found: %v", err)
	}

	// 2. Fetch the data from the TARGET version
	var oldData, oldSchema []byte
	err = tx.QueryRow(`
		SELECT data, schema FROM config_versions 
		WHERE config_id = $1 AND version = $2`,
		configID, targetVersion).Scan(&oldData, &oldSchema)
	if err != nil {
		return nil, fmt.Errorf("target version %d not found: %v", targetVersion, err)
	}

	// 3. Get current max version
	var currentVersion int
	err = tx.QueryRow(`SELECT MAX(version) FROM config_versions WHERE config_id = $1`, configID).Scan(&currentVersion)
	if err != nil {
		return nil, err
	}
	
	newVersion := currentVersion + 1

	// 4. Insert new version as a copy of the old one
	_, err = tx.Exec(`
		INSERT INTO config_versions (config_id, version, data, schema, created_by)
		VALUES ($1, $2, $3, $4, $5)`,
		configID, newVersion, oldData, oldSchema, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create rollback version: %v", err)
	}

	// 5. Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the new config state
	var c Config
	c.ID = configID
	c.ProjectID = projectID
	c.EnvID = envID
	c.Key = key
	c.Version = newVersion
	json.Unmarshal(oldData, &c.Data)
	json.Unmarshal(oldSchema, &c.Schema)

	return &c, nil
}
