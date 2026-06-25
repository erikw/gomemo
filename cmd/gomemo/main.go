// TODO add JSON curl requests
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/erikw/gomemo/internal/api"
	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/seed"
	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/version"
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
	router := api.NewRouter(logger, cfg)

	// TODO storage should be confgurable from env/file.
	notesStore := storage.NewMemory[*notes.Note](logger)

	// Seed database in dev mode.
	// TODO extract this to a cli command `$ gomemo seed`?
	if cfg.Env == "dev" {
		fixturesPath := "data/dev.yaml"
		if err := seed.Load(logger, fixturesPath, notesStore); err != nil {
			logger.Error(fmt.Sprintf("Error seeding database: %v", err.Error()))
			os.Exit(1)
		}
	}

	notesService := notes.NewService(logger, notesStore)
	notesHandler := notes.NewHandler(logger, notesService)

	// Register all handlers that implement RouteRegistrar
	handlers := []api.RouteRegistrar{
		notesHandler,
		// Future handlers (auth, users, etc.) can be added here
	}
	for _, h := range handlers {
		h.RegisterRoutes(router.ChiRouter())
	}

	router.RunServer()
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
