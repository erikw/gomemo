package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	logger *slog.Logger
	config config.Config
}

func NewRouter(logger *slog.Logger, cfg config.Config) *Router {
	return &Router{
		logger: logger,
		config: cfg,
	}
}

func (r *Router) RunServer() {
	cr := chi.NewRouter()
	cr.Use(middleware.Logger)
	cr.Use(middleware.RequestID)
	cr.Use(middleware.Timeout(20 * time.Second))

	cr.Get("/", func(w http.ResponseWriter, req *http.Request) {

		// fmt.Println(middleware.GetReqID(r.Context()))

		r.respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
	err := http.ListenAndServe(r.config.AddrString(), cr)
	if err != nil {
		r.logger.Error("Error serving HTTP.", "error", err.Error)
	}
}

func (r *Router) respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		r.logger.Error("Could not encode JSON response.", "status", status, "value", v)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
