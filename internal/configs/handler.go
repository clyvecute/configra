package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	actorID, _ := r.Context().Value(middleware.ActorIDKey).(int)
	cfg, err := h.service.CreateConfig(projectID, req.EnvID, req.Key, req.Data, req.Schema, actorID)
	if err != nil {
		var valErr *ValidationError
		if errors.As(err, &valErr) || strings.Contains(err.Error(), "validation failed") {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
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

	actorID, _ := r.Context().Value(middleware.ActorIDKey).(int)
	cfg, err := h.service.RollbackConfig(projectID, req.EnvID, req.Key, req.TargetVersion, actorID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, cfg)
}

func (h *Handler) RollbackREST(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EnvID         int `json:"env_id"`
		TargetVersion int `json:"target_version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.EnvID <= 0 || req.TargetVersion <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "env_id and target_version are required"})
		return
	}
	projectID, ok := r.Context().Value(middleware.ProjectIDKey).(int)
	if !ok || projectID == 0 {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized project scope"})
		return
	}
	actorID, _ := r.Context().Value(middleware.ActorIDKey).(int)
	cfg, err := h.service.RollbackConfig(projectID, req.EnvID, r.PathValue("key"), req.TargetVersion, actorID)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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

	if versionStr := r.URL.Query().Get("version"); versionStr != "" {
		version, err := strconv.Atoi(versionStr)
		if err != nil || version <= 0 {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "version must be a positive integer"})
			return
		}
		cfg, err := h.service.GetVersion(projectID, envID, key, version)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if cfg == nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "config version not found"})
			return
		}
		utils.WriteJSON(w, http.StatusOK, cfg)
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

func (h *Handler) Resource(w http.ResponseWriter, r *http.Request) {
	projectID, ok := r.Context().Value(middleware.ProjectIDKey).(int)
	if !ok || projectID == 0 {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized project scope"})
		return
	}
	key := r.PathValue("key")
	if key == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "key is required"})
		return
	}
	if r.Method == http.MethodPost && r.PathValue("action") == "rollback" {
		h.RollbackREST(w, r)
		return
	}
	envID, err := strconv.Atoi(r.URL.Query().Get("env_id"))
	if err != nil || envID <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "env_id must be a positive integer"})
		return
	}
	switch {
	case r.Method == http.MethodGet && r.PathValue("action") == "versions":
		versions, err := h.service.ListVersions(projectID, envID, key)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusOK, versions)
	case r.Method == http.MethodGet && r.PathValue("action") == "diff":
		from, e1 := strconv.Atoi(r.URL.Query().Get("from"))
		to, e2 := strconv.Atoi(r.URL.Query().Get("to"))
		if e1 != nil || e2 != nil || from <= 0 || to <= 0 {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "from and to must be positive versions"})
			return
		}
		a, e := h.service.GetVersion(projectID, envID, key, from)
		if e != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": e.Error()})
			return
		}
		b, e := h.service.GetVersion(projectID, envID, key, to)
		if e != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": e.Error()})
			return
		}
		if a == nil || b == nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "config version not found"})
			return
		}
		changed := map[string]interface{}{}
		for k, v := range b.Data {
			if old, ok := a.Data[k]; !ok || fmt.Sprint(old) != fmt.Sprint(v) {
				changed[k] = map[string]interface{}{"from": a.Data[k], "to": v}
			}
		}
		for k, v := range a.Data {
			if _, ok := b.Data[k]; !ok {
				changed[k] = map[string]interface{}{"from": v, "to": nil}
			}
		}
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"from": from, "to": to, "before": a.Data, "after": b.Data, "changed": changed})
	case r.Method == http.MethodGet:
		version := r.URL.Query().Get("version")
		cfg, err := h.service.GetConfig(projectID, envID, key)
		if version != "" {
			v, e := strconv.Atoi(version)
			if e != nil || v <= 0 {
				utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "version must be a positive integer"})
				return
			}
			cfg, err = h.service.GetVersion(projectID, envID, key, v)
		}
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if cfg == nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "config not found"})
			return
		}
		utils.WriteJSON(w, http.StatusOK, cfg)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
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
