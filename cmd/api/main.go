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
		log.Printf("Warning: Failed to connect to DB: %v. Database-backed features will be disabled.", err)
	} else {
		defer database.Close()
		
		// Auto-migrate database
		log.Println("Running database migrations...")
		if err := db.Migrate(database, "./internal/db/migrations"); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
		} else {
			log.Println("Migrations applied successfully!")
		}
	}
	
	mux := http.NewServeMux()

	// Initialize dependencies
	sentinelClient := configs.NewSentinelClient(cfg.SentinelURL)
	configsRepo := configs.NewRepository(database)
	configsService := configs.NewService(configsRepo, sentinelClient)
	configsHandler := configs.NewHandler(configsService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(database)

	// Register routes
	mux.HandleFunc("/v1/validate", configsHandler.Validate) // No auth needed for local check check
	mux.HandleFunc("/v1/configs", authMiddleware.RequireAPIKey(configsHandler.Create)) // Protected
	mux.HandleFunc("/v1/rollback", authMiddleware.RequireAPIKey(configsHandler.Rollback)) // Protected
	mux.HandleFunc("/fetch", configsHandler.FetchSource) // Internal/External fetch for UI


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

	fmt.Printf("Starting Configra API on :%s\n", cfg.Port)
	// Apply CORS middleware to everything
	if err := http.ListenAndServe(":"+cfg.Port, middleware.CORS(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
