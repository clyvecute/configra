package configs_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"

	"github.com/clyvecute/configra/internal/configs"
	"github.com/clyvecute/configra/internal/db"
	"github.com/clyvecute/configra/internal/middleware"
	_ "github.com/lib/pq"
)

func TestConfigLifecycleAPI(t *testing.T) {
	dsn := os.Getenv("CONFIGRA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set CONFIGRA_TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err = database.Ping(); err != nil {
		t.Fatal(err)
	}
	if err = db.Migrate(database, filepath.Join("..", "db", "migrations")); err != nil {
		t.Fatal(err)
	}
	// Each job gets a fresh PostgreSQL service; clean any prior run data defensively.
	if _, err = database.Exec(`TRUNCATE config_versions,configs,environments,projects,users RESTART IDENTITY CASCADE`); err != nil {
		t.Fatal(err)
	}
	var userID, projectID int
	if err = database.QueryRow(`INSERT INTO users(email,password_hash) VALUES('integration@example.test','test') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if err = database.QueryRow(`INSERT INTO projects(name,owner_id,api_key) VALUES('integration', $1, 'integration-api-key') RETURNING id`, userID).Scan(&projectID); err != nil {
		t.Fatal(err)
	}
	var envID int
	if err = database.QueryRow(`INSERT INTO environments(project_id,name,slug) VALUES($1,'Test','test') RETURNING id`, projectID).Scan(&envID); err != nil {
		t.Fatal(err)
	}
	service := configs.NewService(configs.NewRepository(database), nil)
	handler := configs.NewHandler(service)
	auth := middleware.NewAuthMiddleware(database)
	mux := http.NewServeMux()
	mux.Handle("POST /v1/configs", auth.RequireAPIKey(func(w http.ResponseWriter, r *http.Request) { handler.Create(w, r) }))
	mux.Handle("GET /v1/configs/{key}", auth.RequireAPIKey(handler.Resource))
	mux.Handle("GET /v1/configs/{key}/{action}", auth.RequireAPIKey(handler.Resource))
	mux.Handle("POST /v1/configs/{key}/{action}", auth.RequireAPIKey(handler.Resource))
	server := httptest.NewServer(mux)
	defer server.Close()
	client := server.Client()
	schema := map[string]interface{}{"rules": map[string]interface{}{"enabled": map[string]interface{}{"type": "bool", "required": true}, "limit": map[string]interface{}{"type": "int", "required": true}}}
	push := func(enabled bool, limit int) int {
		t.Helper()
		body := map[string]interface{}{"env_id": envID, "key": "runtime", "data": map[string]interface{}{"enabled": enabled, "limit": limit}, "schema": schema}
		resp := doJSON(t, client, http.MethodPost, server.URL+"/v1/configs", body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("push status=%d", resp.StatusCode)
		}
		var cfg configs.Config
		if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
			t.Fatal(err)
		}
		return cfg.Version
	}
	if v := push(true, 10); v != 1 {
		t.Fatalf("first version=%d", v)
	}
	if v := push(false, 10); v != 2 {
		t.Fatalf("second version=%d", v)
	}
	get := func(url string) map[string]interface{} {
		t.Helper()
		resp := doJSON(t, client, http.MethodGet, server.URL+url, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET %s status=%d", url, resp.StatusCode)
		}
		var value map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		return value
	}
	latest := get(fmt.Sprintf("/v1/configs/runtime?env_id=%d", envID))
	if latest["version"] != float64(2) {
		t.Fatalf("latest version: %v", latest["version"])
	}
	old := get(fmt.Sprintf("/v1/configs/runtime?env_id=%d&version=1", envID))
	if old["version"] != float64(1) {
		t.Fatalf("historical version: %v", old["version"])
	}
	diff := get(fmt.Sprintf("/v1/configs/runtime/diff?env_id=%d&from=1&to=2", envID))
	changed := diff["changed"].(map[string]interface{})
	if len(changed) != 1 {
		t.Fatalf("changed keys = %v", changed)
	}
	if _, ok := changed["enabled"]; !ok {
		t.Fatalf("expected enabled-only diff: %v", changed)
	}
	invalid := doJSON(t, client, http.MethodPost, server.URL+"/v1/configs", map[string]interface{}{"env_id": envID, "key": "runtime", "data": map[string]interface{}{"enabled": "wrong", "limit": 10}, "schema": schema})
	invalid.Body.Close()
	if invalid.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid push status=%d", invalid.StatusCode)
	}
	versionsURL := fmt.Sprintf("/v1/configs/runtime/versions?env_id=%d", envID)
	before := getArray(t, client, server.URL+versionsURL)
	if len(before) != 2 {
		t.Fatalf("invalid push created version: count=%d", len(before))
	}
	rollback := doJSON(t, client, http.MethodPost, server.URL+"/v1/configs/runtime/rollback", map[string]int{"env_id": envID, "target_version": 1})
	defer rollback.Body.Close()
	if rollback.StatusCode != http.StatusOK {
		t.Fatalf("rollback status=%d", rollback.StatusCode)
	}
	var rolled configs.Config
	if err = json.NewDecoder(rollback.Body).Decode(&rolled); err != nil {
		t.Fatal(err)
	}
	if rolled.Version != 3 || rolled.Data["enabled"] != true || rolled.Data["limit"] != float64(10) {
		t.Fatalf("rollback result: %+v", rolled)
	}
	after := getArray(t, client, server.URL+versionsURL)
	if len(after) != 3 {
		t.Fatalf("rollback must append version 3, count=%d", len(after))
	}
	if after[0]["version"] != float64(3) || after[1]["version"] != float64(2) || after[2]["version"] != float64(1) {
		t.Fatalf("history not append-only: %+v", after)
	}
	if after[2]["author"] != "integration@example.test" {
		t.Fatalf("author attribution=%v", after[2]["author"])
	}

	// Exercise concurrent first pushes for a new key and ensure all versions
	// are allocated once, without duplicates or gaps.
	const writers = 20
	var wg sync.WaitGroup
	versions := make(chan int, writers)
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			body := map[string]interface{}{"env_id": envID, "key": "concurrent", "data": map[string]interface{}{"enabled": true, "limit": i}, "schema": schema}
			resp, err := client.Do(mustRequest(t, server.URL+"/v1/configs", body))
			if err != nil {
				errs <- err
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				errs <- fmt.Errorf("push status=%d", resp.StatusCode)
				return
			}
			var c configs.Config
			if err = json.NewDecoder(resp.Body).Decode(&c); err != nil {
				errs <- err
				return
			}
			versions <- c.Version
		}(i)
	}
	wg.Wait()
	close(versions)
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	got := make([]int, 0, writers)
	for v := range versions {
		got = append(got, v)
	}
	if len(got) != writers {
		t.Fatalf("successful concurrent pushes=%d, want %d", len(got), writers)
	}
	sort.Ints(got)
	for i, v := range got {
		if v != i+1 {
			t.Fatalf("concurrent version sequence=%v", got)
		}
	}
}

func mustRequest(t *testing.T, url string, body interface{}) *http.Request {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-Key", "integration-api-key")
	req.Header.Set("Content-Type", "application/json")
	return req
}

func doJSON(t *testing.T, client *http.Client, method, url string, body interface{}) *http.Response {
	t.Helper()
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-Key", "integration-api-key")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
func getArray(t *testing.T, client *http.Client, url string) []map[string]interface{} {
	t.Helper()
	resp := doJSON(t, client, http.MethodGet, url, nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("history status=%d", resp.StatusCode)
	}
	var items []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatal(err)
	}
	return items
}
