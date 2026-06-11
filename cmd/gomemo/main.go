package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/erikw/gomemo/internal/api"
	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/notes"
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

	notesService := notes.NewService(logger)
	notesHandler := notes.NewHandler(logger, notesService)

	router := api.NewRouter(logger, cfg, notesHandler)
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
