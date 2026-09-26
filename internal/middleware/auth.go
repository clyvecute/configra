package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/clyvecute/configra/pkg/utils"
)

type AuthMiddleware struct {
	db *sql.DB
}

func NewAuthMiddleware(db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

type contextKey string
const ProjectIDKey contextKey = "projectID"

func (m *AuthMiddleware) RequireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing api key"})
			return
		}

		if m.db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database connection unavailable"})
			return
		}

		var projectID int
		// Simple query to validate key and get ID
		err := m.db.QueryRow("SELECT id FROM projects WHERE api_key = $1", apiKey).Scan(&projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid api key"})
				return
			}
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "auth error"})
			return
		}

		// Store projectID in context
		ctx := context.WithValue(r.Context(), ProjectIDKey, projectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
