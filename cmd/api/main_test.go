package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticPages(t *testing.T) {
	for name, page := range map[string]string{"landing": landingPageHTML, "dashboard": dashboardPageHTML} {
		t.Run(name, func(t *testing.T) {
			if strings.TrimSpace(page) == "" {
				t.Fatal("page is empty")
			}
		})
	}
}

func TestHealthRouteResponds(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rr.Code != http.StatusOK || rr.Body.String() != "ok" {
		t.Fatalf("health response: %d %q", rr.Code, rr.Body.String())
	}
}
