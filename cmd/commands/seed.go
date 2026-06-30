package commands

import (
	"fmt"

	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/seed"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with fixtures",
	Long:  `Load fixture data from a YAML file into the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := GetLogger()
		logger.Info("Starting seed command")

		cfg, err := config.Load()
		if err != nil {
			logger.Error(fmt.Sprintf("Error loading configuration: %v", err.Error()))
			return err
		}

		store, cleanup, err := initializeNotesStore(logger, cfg)
		if err != nil {
			logger.Error("Error initializing storage", "error", err.Error(), "storageType", cfg.StorageType)
			return err
		}
		defer cleanup()

		// Load fixtures from dev.yaml
		fixturesPath := "data/dev.yaml"
		if err := seed.Load(logger, fixturesPath, store); err != nil {
			logger.Error(fmt.Sprintf("Error seeding database: %v", err.Error()))
			return err
		}

		logger.Info("Database seeded successfully")
		return nil
	},
}
