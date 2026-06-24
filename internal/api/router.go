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
	logger    *slog.Logger
	config    config.Config
	chiRouter chi.Router
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
	r.chiRouter.Get("/", r.getHealth)
	r.chiRouter.Get("/health", r.getHealth)
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
