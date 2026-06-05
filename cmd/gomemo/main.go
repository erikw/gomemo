package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var logger *slog.Logger

func init() {
	initLogger(true) // TODO move to main() to allow --debug cli flag?
}

func initLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove timestamp from output.
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		}),
	)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Error loading configuration")
		os.Exit(1)
	}

	logger.Info("Starting Gomemo.", "config", cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(20 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		// fmt.Println(middleware.GetReqID(r.Context()))

		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
	err = http.ListenAndServe(cfg.AddrString(), r)
	if err != nil {
		logger.Error("Error serving HTTP.", "error", err.Error)
	}
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("Could not encode JSON response.", "status", status, "value", v)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
