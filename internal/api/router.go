package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/httpx"
	"github.com/erikw/gomemo/internal/notes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	logger       *slog.Logger
	config       config.Config
	notesHandler *notes.Handler
}

func NewRouter(logger *slog.Logger, cfg config.Config, notesHandler *notes.Handler) *Router {
	return &Router{
		logger:       logger,
		config:       cfg,
		notesHandler: notesHandler,
	}
}

func (r *Router) RunServer() {
	cr := chi.NewRouter()
	cr.Use(middleware.Logger)
	cr.Use(middleware.RequestID)
	cr.Use(middleware.Timeout(20 * time.Second))

	r.setupRoutes(cr)

	err := http.ListenAndServe(r.config.AddrString(), cr)
	if err != nil {
		r.logger.Error("Error serving HTTP.", "error", err.Error)
	}
}

func (r *Router) setupRoutes(cr chi.Router) {
	cr.Get("/", r.getHealth)
	cr.Get("/health", r.getHealth)

	cr.Get("/notes/{noteID}", r.notesHandler.HandleGetByID)
}

func (r *Router) getHealth(w http.ResponseWriter, req *http.Request) {
	// fmt.Println(middleware.GetReqID(r.Context()))

	resp := map[string]string{
		"status": "ok",
	}
	status := http.StatusOK
	if err := httpx.RespondJSON(w, status, resp); err != nil {
		r.logger.Error("Could not encode JSON response.", "status", status, "value", resp)
	}
}
