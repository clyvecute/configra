package main

import (
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
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		defer database.Close()

		// Auto-migrate database
		log.Println("Running database migrations...")
		if err := db.Migrate(database, "./internal/db/migrations"); err != nil {
			log.Fatalf("Database migration failed: %v", err)
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
	mux.HandleFunc("/v1/validate", configsHandler.Validate) // No auth needed for local check
	mux.HandleFunc("/v1/configs", authMiddleware.RequireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			configsHandler.Create(w, r)
		case http.MethodGet:
			configsHandler.Get(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})) // Protected
	mux.Handle("GET /v1/configs/{key}", authMiddleware.RequireAPIKey(configsHandler.Resource))
	mux.Handle("GET /v1/configs/{key}/{action}", authMiddleware.RequireAPIKey(configsHandler.Resource))
	mux.Handle("POST /v1/configs/{key}/{action}", authMiddleware.RequireAPIKey(configsHandler.Resource))
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!doctype html><html><head><title>Configra API</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>SwaggerUIBundle({url:'/openapi.yaml',dom_id:'#swagger-ui'});</script></body></html>`))
	})
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "openapi.yaml") })
	mux.HandleFunc("/v1/rollback", authMiddleware.RequireAPIKey(configsHandler.Rollback)) // Protected
	mux.HandleFunc("/fetch", authMiddleware.RequireAPIKey(configsHandler.FetchSource))    // Protected external fetch

	// Liveness check: process is running.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	// Readiness check: required backing services are reachable.
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := database.PingContext(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	// Dashboard handler — Space White minimalist control center
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(dashboardPageHTML))
	})

	// Root handler — HTML landing page for browser/portfolio visitors
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(landingPageHTML))
	})

	fmt.Printf("Starting Configra API on :%s\n", cfg.Port)
	// Apply CORS middleware to everything
	if err := http.ListenAndServe(":"+cfg.Port, middleware.CORS(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
