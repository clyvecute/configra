package configs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clyvecute/configra/internal/middleware"
	"github.com/clyvecute/configra/pkg/utils"
)

type ValidateRequest struct {
	Schema Schema                 `json:"schema"`
	Config map[string]interface{} `json:"config"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := Validate(req.Schema, req.Config); err != nil {
		// If validation fails, return 400 with the error details
		// Ideally we cast err to ValidationError to get structured data, but string is fine for now
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "valid"})
}

type CreateRequest struct {
	ProjectID int                    `json:"project_id"`
	EnvID     int                    `json:"env_id"`
	Key       string                 `json:"key"`
	Data      map[string]interface{} `json:"data"`
	Schema    map[string]interface{} `json:"schema"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// Basic validation
	if req.EnvID == 0 || req.Key == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}

	// Security: Get ProjectID from context (set by middleware)
	// We ignore req.ProjectID to prevent spoofing
	projectID, ok := r.Context().Value(middleware.ProjectIDKey).(int)
	if !ok || projectID == 0 {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized project scope"})
		return
	}

	// Call Service
	cfg, err := h.service.CreateConfig(projectID, req.EnvID, req.Key, req.Data, req.Schema, 1)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, cfg)
}

type RollbackRequest struct {
	ProjectID     int    `json:"project_id"`
	EnvID         int    `json:"env_id"`
	Key           string `json:"key"`
	TargetVersion int    `json:"target_version"`
}

func (h *Handler) Rollback(w http.ResponseWriter, r *http.Request) {
	var req RollbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.EnvID == 0 || req.Key == "" || req.TargetVersion == 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}

	// Security: Get ProjectID from context
	projectID, ok := r.Context().Value(middleware.ProjectIDKey).(int)
	if !ok || projectID == 0 {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized project scope"})
		return
	}

	cfg, err := h.service.RollbackConfig(projectID, req.EnvID, req.Key, req.TargetVersion, 1) // default admin ID
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, cfg)
}

// Get retrieves the latest config version for a given project/env/key.
// Query params: key (required), env_id (required).
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	key := r.URL.Query().Get("key")
	envIDStr := r.URL.Query().Get("env_id")
	if key == "" || envIDStr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required query params: key, env_id"})
		return
	}

	var envID int
	if _, err := fmt.Sscanf(envIDStr, "%d", &envID); err != nil || envID <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "env_id must be a positive integer"})
		return
	}

	projectID, ok := r.Context().Value(middleware.ProjectIDKey).(int)
	if !ok || projectID == 0 {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized project scope"})
		return
	}

	cfg, err := h.service.GetConfig(projectID, envID, key)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if cfg == nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "config not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, cfg)
}

type FetchRequest struct {
	URL string `json:"url"`
}

func (h *Handler) FetchSource(w http.ResponseWriter, r *http.Request) {
	var req FetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	data, err := h.service.FetchExternal(req.URL)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, data)
}
