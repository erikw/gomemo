package storage_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/testutil"
)

type testObject struct {
	id    int64
	Title string
}

func (o *testObject) GetID() int64         { return o.id }
func (o *testObject) SetID(ID int64) error { o.id = ID; return nil }
func (o *testObject) Clone() *testObject   { clone := *o; return &clone }

func TestMemoryCRUD(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemory[*testObject](testutil.Logger())

	n1, err := store.Create(ctx, &testObject{Title: "first"})
	if err != nil {
		t.Fatalf("Create(first) error = %v", err)
	}
	if n1.id != 0 {
		t.Fatalf("first ID = %d, want 0", n1.id)
	}

	n2, err := store.Create(ctx, &testObject{Title: "second"})
	if err != nil {
		t.Fatalf("Create(second) error = %v", err)
	}
	if n2.id != 1 {
		t.Fatalf("second ID = %d, want 1", n2.id)
	}

	got, err := store.FindByID(ctx, n1.id)
	if err != nil {
		t.Fatalf("FindByID(%d) error = %v", n1.id, err)
	}
	if got.Title != "first" {
		t.Fatalf("FindByID(%d).Title = %q, want %q", n1.id, got.Title, "first")
	}

	all, err := store.All(ctx)
	if err != nil {
		t.Fatalf("All() error = %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("len(All()) = %d, want 2", len(all))
	}

	updated := &testObject{
		id:    n1.id,
		Title: "first-updated",
	}
	if _, err := store.InsertWithID(ctx, n1.id, updated); err != nil {
		t.Fatalf("InsertWithID(%d) error = %v", n1.id, err)
	}

	got, err = store.FindByID(ctx, n1.id)
	if err != nil {
		t.Fatalf("FindByID(%d) after update error = %v", n1.id, err)
	}
	if got.Title != "first-updated" {
		t.Fatalf("updated title = %q, want %q", got.Title, "first-updated")
	}

	deleted, err := store.DeleteByID(ctx, n2.id)
	if err != nil {
		t.Fatalf("DeleteByID(%d) error = %v", n2.id, err)
	}
	if !deleted {
		t.Fatalf("DeleteByID(%d) = false, want true", n2.id)
	}

	deleted, err = store.DeleteByID(ctx, n2.id)
	if err != nil {
		t.Fatalf("DeleteByID(%d) second call error = %v", n2.id, err)
	}
	if deleted {
		t.Fatalf("DeleteByID(%d) second call = true, want false", n2.id)
	}
}

func TestMemoryDoesNotExposeStoredObjects(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemory[*testObject](testutil.Logger())

	created, err := store.Create(ctx, &testObject{Title: "original"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	created.Title = "changed-created-result"
	found, err := store.FindByID(ctx, created.id)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.Title != "original" {
		t.Fatalf("stored title after mutating Create result = %q, want %q", found.Title, "original")
	}

	found.Title = "changed-find-result"
	all, err := store.All(ctx)
	if err != nil {
		t.Fatalf("All() error = %v", err)
	}
	all[0].Title = "changed-all-result"

	found, err = store.FindByID(ctx, created.id)
	if err != nil {
		t.Fatalf("FindByID() after mutations error = %v", err)
	}
	if found.Title != "original" {
		t.Fatalf("stored title after mutating read results = %q, want %q", found.Title, "original")
	}
}

func TestMemoryConcurrentReadsAndWrites(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemory[*testObject](testutil.Logger())

	created, err := store.Create(ctx, &testObject{Title: "initial"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	const (
		readers    = 4
		writers    = 4
		iterations = 1000
	)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				if _, err := store.FindByID(ctx, created.id); err != nil {
					t.Errorf("FindByID() error = %v", err)
					return
				}
				if _, err := store.All(ctx); err != nil {
					t.Errorf("All() error = %v", err)
					return
				}
			}
		}()
	}

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(writer int) {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				if _, err := store.InsertWithID(ctx, created.id, &testObject{
					id:    created.id,
					Title: fmt.Sprintf("writer-%d-%d", writer, j),
				}); err != nil {
					t.Errorf("InsertWithID() error = %v", err)
					return
				}
			}
		}(i)
	}

	close(start)
	wg.Wait()
}
