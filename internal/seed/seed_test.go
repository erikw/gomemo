package seed

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/testutil"
)

func TestLoadFixtures(t *testing.T) {
	dir := t.TempDir()
	fixturesPath := filepath.Join(dir, "fixtures.yaml")
	content := []byte("notes:\n  - title: first\n    content: first-content\n  - title: second\n    content: second-content\n")
	if err := os.WriteFile(fixturesPath, content, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := storage.NewMemory[*notes.Note](testutil.Logger())
	if err := Load(testutil.Logger(), fixturesPath, store); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	all, err := store.All(context.Background())
	if err != nil {
		t.Fatalf("All() error = %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("len(All()) = %d, want 2", len(all))
	}
}

func TestLoadMissingFile(t *testing.T) {
	store := storage.NewMemory[*notes.Note](testutil.Logger())
	err := Load(testutil.Logger(), filepath.Join(t.TempDir(), "does-not-exist.yaml"), store)
	if err == nil {
		t.Fatalf("Load() error = nil, want non-nil")
	}
}
