package storage_test

import (
	"context"
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
