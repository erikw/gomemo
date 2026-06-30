package seed

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/storage"
	"gopkg.in/yaml.v3"
)

type fixtureNote struct {
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
}

type fixtures struct {
	Notes []fixtureNote `yaml:"notes"`
}

func Load(logger *slog.Logger, fixturesPath string, store storage.Storage[*notes.Note]) error {
	logger.Info("Loading fixtures", "path", fixturesPath)

	data, err := os.ReadFile(fixturesPath)
	if err != nil {
		return fmt.Errorf("failed to read fixtures file: %w", err)
	}

	var f fixtures
	if err := yaml.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("failed to parse fixtures YAML: %w", err)
	}

	ctx := context.Background()
	if err := store.Clear(ctx); err != nil {
		return fmt.Errorf("failed to clear existing data before seed: %w", err)
	}

	now := time.Now()
	for i, n := range f.Notes {
		time := now.Add(-time.Duration(len(f.Notes)-i) * time.Hour)
		note := &notes.Note{
			Title:      n.Title,
			Content:    n.Content,
			CreatedAt:  time,
			ModifiedAt: time,
		}

		_, err := store.Create(ctx, note)
		if err != nil {
			return fmt.Errorf("failed to seed note %d: %w", i, err)
		}
	}

	logger.Info("Fixtures loaded successfully", "count", len(f.Notes))
	return nil
}
