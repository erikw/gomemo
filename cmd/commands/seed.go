package commands

import (
	"fmt"

	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/seed"
	"github.com/erikw/gomemo/internal/storage"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with fixtures",
	Long:  `Load fixture data from a YAML file into the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := GetLogger()
		logger.Info("Starting seed command")

		// Initialize storage
		store := storage.NewMemory[*notes.Note](logger)

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
