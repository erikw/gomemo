package commands

import (
	"fmt"

	"github.com/erikw/gomemo/internal/api"
	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/seed"
	"github.com/erikw/gomemo/internal/storage"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  `Start the Gomemo web server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := GetLogger()
		logger.Info("Starting serve command")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			logger.Error(fmt.Sprintf("Error loading configuration: %v", err.Error()))
			return err
		}

		logger.Info("Starting Gomemo.", "config", cfg)
		router := api.NewRouter(logger, cfg)

		// Initialize storage
		notesStore := storage.NewMemory[*notes.Note](logger)

		// For in-memory storage, auto-seed fixtures on every run
		if cfg.IsMemoryStorage() {
			fixturesPath := "data/dev.yaml"
			if err := seed.Load(logger, fixturesPath, notesStore); err != nil {
				logger.Error(fmt.Sprintf("Error seeding database: %v", err.Error()))
				return err
			}
		}

		// Initialize services
		notesService := notes.NewService(logger, notesStore)
		notesHandler := notes.NewHandler(logger, notesService)

		// Register all handlers under /api/v1
		router.RegisterV1Routes(notesHandler)

		// Start server
		router.RunServer()
		return nil
	},
}
