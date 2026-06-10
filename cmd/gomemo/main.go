package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var logger *slog.Logger

func main() {
	var err error
	var args config.Args

	args, err = config.ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing args: %v\n", err.Error())
		os.Exit(1)
	}

	if args.Help {
		config.PrintHelp()
		os.Exit(0)
	}

	if args.Version {
		fmt.Printf("Version: %s\n", version.Version)
		os.Exit(0)
	}

	initLogger(args.Debug)

	var cfg config.Config
	cfg, err = config.Load()
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading configuration: %v", err.Error()))
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

func initLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)

	if debug {
		logger.Debug("Debug logging enabled.")
	}

}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("Could not encode JSON response.", "status", status, "value", v)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
