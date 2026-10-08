package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Migrate(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %v", err)
	}
	var hasSchema bool
	if err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema=current_schema() AND table_name='configs')`).Scan(&hasSchema); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`); err != nil {
		return fmt.Errorf("failed to initialize migration tracking: %w", err)
	}
	// Existing installations predate migration tracking. Baseline the current
	// schema rather than replaying CREATE statements over live data.
	if hasSchema {
		if _, err := db.Exec(`INSERT INTO schema_migrations(name) VALUES('001_initial_schema.sql') ON CONFLICT DO NOTHING`); err != nil {
			return fmt.Errorf("failed to baseline existing schema: %w", err)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			var alreadyApplied bool
			if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name=$1)`, entry.Name()).Scan(&alreadyApplied); err != nil {
				return err
			}
			if alreadyApplied {
				continue
			}
			path := filepath.Join(migrationsDir, entry.Name())
			fmt.Printf("Applying migration: %s\n", entry.Name())

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", path, err)
			}

			// Simple migration: split by "-- Down" and take the first part
			sqlContent := string(content)
			parts := strings.Split(sqlContent, "-- Down")
			upSQL := parts[0]

			tx, err := db.Begin()
			if err != nil {
				return err
			}
			if _, err := tx.Exec(upSQL); err != nil {
				tx.Rollback()
				// In a real app we would check if migration already applied
				// For this demo, we ignore "relation already exists" errors loosely or just fail
				// But to confirm it works, let's just create a version table check later.
				// For now, let's just log and continue if it might be idempotent-ish or fail.
				// Actually, simpler: Wrap in transaction.
				return fmt.Errorf("failed to exec migration %s: %v", entry.Name(), err)
			}
			if _, err = tx.Exec(`INSERT INTO schema_migrations(name) VALUES($1)`, entry.Name()); err != nil {
				tx.Rollback()
				return err
			}
			if err = tx.Commit(); err != nil {
				return err
			}
		}
	}
	return nil
}
