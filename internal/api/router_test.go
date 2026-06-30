package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/testutil"
	"github.com/go-chi/chi/v5"
)

type testRegistrar struct {
	called bool
}

func (r *testRegistrar) RegisterRoutes(router chi.Router) {
	r.called = true
	router.Get("/test-route", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func TestHealthEndpoints(t *testing.T) {
	router := NewRouter(testutil.Logger(), config.Config{Host: "127.0.0.1", Port: "8080"})

	for _, path := range []string{"/", "/health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		router.ChiRouter().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d, want %d", path, rec.Code, http.StatusOK)
		}
	}
}

func TestRegisterV1Routes(t *testing.T) {
	router := NewRouter(testutil.Logger(), config.Config{Host: "127.0.0.1", Port: "8080"})
	registrar := &testRegistrar{}
	router.RegisterV1Routes(registrar)

	if !registrar.called {
		t.Fatalf("RegisterRoutes() was not called")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test-route", nil)
	rec := httptest.NewRecorder()
	router.ChiRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("GET /api/v1/test-route status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}
