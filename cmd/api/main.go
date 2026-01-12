package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/clyvecute/configra/internal/config"
	"github.com/clyvecute/configra/internal/configs"
	"github.com/clyvecute/configra/internal/db"
	"github.com/clyvecute/configra/internal/middleware"
)

func main() {
	cfg := config.Load()

	// Connect to DB
	database, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()

	// Auto-migrate database (simplifies cloud deployment)
	log.Println("Running database migrations...")
	// We need to find the migrations directory. In Docker, we'll COPY it to a known location.
	// Or we can embed it using `embed`. embedding is much better for single-binary deploys.
	// For now, let's assume relative path works if we setup Dockerfile correctly, 
	// but strictly speaking, embedding is the "Pro" way.
	// Let's stick to the existing Migrator which takes a path.
	// We will assume "./internal/db/migrations" is available relative to CWD.
	
	if err := db.Migrate(database, "./internal/db/migrations"); err != nil {
		log.Printf("Warning: Migration failed (might be already up to date or path issue): %v", err)
		// Don't fatal here in case it's just a path issue in dev, but in prod we want to know.
	} else {
		log.Println("Migrations applied successfully!")
	}
	
	mux := http.NewServeMux()

	// Initialize dependencies
	configsRepo := configs.NewRepository(database)
	configsService := configs.NewService(configsRepo)
	configsHandler := configs.NewHandler(configsService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(database)

	// Register routes
	mux.HandleFunc("/v1/validate", configsHandler.Validate) // No auth needed for local check check
	mux.HandleFunc("/v1/configs", authMiddleware.RequireAPIKey(configsHandler.Create)) // Protected
	mux.HandleFunc("/v1/rollback", authMiddleware.RequireAPIKey(configsHandler.Rollback)) // Protected


	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Root handler (Landing Page) - Friendly message for browser/portfolio visitors
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"service": "Configra API",
			"status":  "running",
			"docs":    "This is a JSON-only API. Use the CLI or API endpoints.",
			"endpoints": "/health, /v1/validate, /v1/configs, /v1/rollback",
		}
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Starting Configra API on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
