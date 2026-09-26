package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Migrate(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
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

			if _, err := db.Exec(upSQL); err != nil {
				// In a real app we would check if migration already applied
				// For this demo, we ignore "relation already exists" errors loosely or just fail
				// But to confirm it works, let's just create a version table check later.
				// For now, let's just log and continue if it might be idempotent-ish or fail.
				// Actually, simpler: Wrap in transaction.
				return fmt.Errorf("failed to exec migration %s: %v", entry.Name(), err)
			}
		}
	}
	return nil
}
