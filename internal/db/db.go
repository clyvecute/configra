package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

func Connect(cfg Config) (*sql.DB, error) {
	var dsn string
	if cfg.URL != "" {
		dsn = cfg.URL
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Retry connection
	for i := 0; i < 5; i++ {
		if err = db.Ping(); err == nil {
			log.Println("Successfully connected to database")
			return db, nil
		}
		log.Printf("Failed to connect to DB, retrying in 2s... (%v)", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %v", err)
}
