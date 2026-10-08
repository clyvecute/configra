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
const ActorIDKey contextKey = "actorID"
const ActorNameKey contextKey = "actorName"

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

		var projectID, actorID int
		var actorName string
		// Project credentials are attributed to their owning user. Ownerless legacy
		// projects retain a stable project-level label in audit history.
		err := m.db.QueryRow(`SELECT p.id, COALESCE(p.owner_id, 0), COALESCE(u.email, 'project:' || p.id::text)
			FROM projects p LEFT JOIN users u ON u.id=p.owner_id WHERE p.api_key = $1`, apiKey).Scan(&projectID, &actorID, &actorName)
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
		ctx = context.WithValue(ctx, ActorIDKey, actorID)
		ctx = context.WithValue(ctx, ActorNameKey, actorName)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
