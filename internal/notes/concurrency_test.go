package notes

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/testutil"
)

func TestServiceConcurrentListAndUpdate(t *testing.T) {
	ctx := context.Background()
	service := NewService(testutil.Logger(), storage.NewMemory[*Note](testutil.Logger()))

	created, err := service.Create(ctx, "initial", "content")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	const (
		listers    = 4
		updaters   = 4
		iterations = 1000
	)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < listers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				notes, err := service.GetAll(ctx)
				if err != nil {
					t.Errorf("GetAll() error = %v", err)
					return
				}
				for _, note := range notes {
					_ = note.Title
					_ = note.Content
				}
			}
		}()
	}

	for i := 0; i < updaters; i++ {
		wg.Add(1)
		go func(updater int) {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				title := fmt.Sprintf("title-%d-%d", updater, j)
				if _, err := service.Update(ctx, created.ID, &title, nil); err != nil {
					t.Errorf("Update() error = %v", err)
					return
				}
			}
		}(i)
	}

	close(start)
	wg.Wait()
}
