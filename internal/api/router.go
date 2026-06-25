package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router wraps the chi router and manages HTTP routes and middleware.
// Handlers are registered by implementing the RouteRegistrar interface.
type Router struct {
	logger      *slog.Logger
	config      config.Config
	chiRouter   chi.Router
	v1Router    chi.Router
}

func NewRouter(logger *slog.Logger, cfg config.Config) *Router {
	r := Router{
		logger:    logger,
		config:    cfg,
		chiRouter: chi.NewRouter(),
	}
	r.setupMiddlewares()
	r.setupRoutes()
	return &r
}

func (r *Router) ChiRouter() chi.Router { return r.chiRouter }

// RegisterV1Routes registers handlers under the /api/v1 prefix
func (r *Router) RegisterV1Routes(registrars ...RouteRegistrar) {
	for _, registrar := range registrars {
		registrar.RegisterRoutes(r.v1Router)
	}
}

func (r *Router) RunServer() {
	err := http.ListenAndServe(r.config.AddrString(), r.chiRouter)
	if err != nil {
		r.logger.Error("Error serving HTTP.", "error", err.Error())
	}
}

func (r *Router) setupMiddlewares() {
	r.chiRouter.Use(middleware.Logger)
	r.chiRouter.Use(middleware.RequestID)
	r.chiRouter.Use(middleware.Timeout(20 * time.Second))

}

func (r *Router) setupRoutes() {
	// Health endpoints - unversioned (best practice for monitoring)
	r.chiRouter.Get("/", r.getHealth)
	r.chiRouter.Get("/health", r.getHealth)

	// API v1 routes
	r.chiRouter.Route("/api/v1", func(v1Router chi.Router) {
		r.v1Router = v1Router
		// Handlers will register their routes here via RegisterV1Routes
	})
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
