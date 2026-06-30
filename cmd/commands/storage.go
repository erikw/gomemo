package commands

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/storage"
	pgstorage "github.com/erikw/gomemo/internal/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type storeCleanupFunc func()

func initializeNotesStore(logger *slog.Logger, cfg config.Config) (storage.Storage[*notes.Note], storeCleanupFunc, error) {
	if cfg.IsMemoryStorage() {
		return storage.NewMemory[*notes.Note](logger), func() {}, nil
	}

	if !cfg.IsPostgresStorage() {
		return nil, nil, fmt.Errorf("unsupported storage type: %q", cfg.StorageType)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("could not ping postgres: %w", err)
	}

	return pgstorage.NewNoteStore(pool), func() {
		pool.Close()
	}, nil
}
