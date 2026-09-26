package projects

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int       `json:"owner_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(name string, ownerID int) (*Project, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO projects (name, owner_id, api_key) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	
	p := &Project{
		Name:    name,
		OwnerID: ownerID,
		APIKey:  apiKey,
	}

	err = r.db.QueryRow(query, name, ownerID, apiKey).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repository) GetByID(id int) (*Project, error) {
	query := `SELECT id, name, owner_id, api_key, created_at FROM projects WHERE id = $1`
	
	p := &Project{}
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.OwnerID, &p.APIKey, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return p, nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
