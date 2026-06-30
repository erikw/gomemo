package notes

import (
	"context"
	"testing"
	"time"

	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/testutil"
)

func newServiceForTest() (*Service, *storage.Memory[*Note]) {
	store := storage.NewMemory[*Note](testutil.Logger())
	return NewService(testutil.Logger(), store), store
}

func TestServiceCreate(t *testing.T) {
	service, _ := newServiceForTest()

	note, err := service.Create(context.Background(), "title", "content")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if note.ID != 0 {
		t.Fatalf("Create() ID = %d, want 0", note.ID)
	}
	if note.Title != "title" {
		t.Fatalf("Create() Title = %q, want %q", note.Title, "title")
	}
}

func TestServiceCreateInvalidTitle(t *testing.T) {
	service, _ := newServiceForTest()

	_, err := service.Create(context.Background(), "   ", "content")
	if err == nil {
		t.Fatalf("Create() error = nil, want non-nil")
	}
	if err != ErrTitleRequired {
		t.Fatalf("Create() error = %v, want %v", err, ErrTitleRequired)
	}
}

func TestServiceGetAllPaginated(t *testing.T) {
	service, _ := newServiceForTest()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := service.Create(ctx, "title", "content")
		if err != nil {
			t.Fatalf("Create() #%d error = %v", i, err)
		}
	}

	notes, total, err := service.GetAllPaginated(ctx, 2, 1)
	if err != nil {
		t.Fatalf("GetAllPaginated() error = %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("len(notes) = %d, want 2", len(notes))
	}
	if total != 3 {
		t.Fatalf("total = %d, want 3", total)
	}

	notes, total, err = service.GetAllPaginated(ctx, 2, 10)
	if err != nil {
		t.Fatalf("GetAllPaginated() with high offset error = %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("len(notes) with high offset = %d, want 0", len(notes))
	}
	if total != 3 {
		t.Fatalf("total with high offset = %d, want 3", total)
	}
}

func TestServiceUpdate(t *testing.T) {
	service, _ := newServiceForTest()
	ctx := context.Background()

	created, err := service.Create(ctx, "old", "content")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	before := created.ModifiedAt
	time.Sleep(2 * time.Millisecond)
	newTitle := "new"
	updated, err := service.Update(ctx, created.ID, &newTitle, nil)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated.Title != "new" {
		t.Fatalf("Update() Title = %q, want %q", updated.Title, "new")
	}
	if !updated.ModifiedAt.After(before) {
		t.Fatalf("Update() ModifiedAt = %v, want after %v", updated.ModifiedAt, before)
	}
}

func TestServiceDeleteByID(t *testing.T) {
	service, _ := newServiceForTest()
	ctx := context.Background()

	created, err := service.Create(ctx, "title", "content")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	deleted, err := service.DeleteByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeleteByID() error = %v", err)
	}
	if !deleted {
		t.Fatalf("DeleteByID() = false, want true")
	}
}
