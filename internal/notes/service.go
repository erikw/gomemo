package notes

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/erikw/gomemo/internal/storage"
)

type Service struct {
	logger *slog.Logger
	store  storage.Storage[*Note]
}

func NewService(logger *slog.Logger, store storage.Storage[*Note]) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*Note, error) {
	var notes []*Note
	var err error
	if notes, err = s.store.All(ctx); err != nil {
		s.logger.Error("could not retrieve all Notes")
		return make([]*Note, 0), err
	}

	return notes, nil
}

func (s *Service) GetAllPaginated(ctx context.Context, limit int, offset int) ([]*Note, int64, error) {
	notes, err := s.GetAll(ctx)
	if err != nil {
		return make([]*Note, 0), 0, err
	}

	total := int64(len(notes))
	
	if offset >= len(notes) {
		return make([]*Note, 0), total, nil
	}

	end := offset + limit
	if end > len(notes) {
		end = len(notes)
	}

	return notes[offset:end], total, nil
}

func (s *Service) GetByID(ctx context.Context, ID int64) (*Note, error) {
	var note *Note
	var err error
	if note, err = s.store.FindByID(ctx, ID); err != nil {
		s.logger.Error(fmt.Sprintf("could not find Note with ID %d", ID))
		return &Note{}, err
	}

	return note, nil
}

func (s *Service) Create(ctx context.Context, title string, content string) (*Note, error) {
	note := &Note{
		Title:      title,
		Content:    content,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if err := note.Validate(); err != nil {
		s.logger.Warn(fmt.Sprintf("Note is not valid: %s", err.Error()))
		return nil, err
	}

	var createdNote *Note
	var err error
	if createdNote, err = s.store.Create(ctx, note); err != nil {
		s.logger.Error("could not create a new Note in storage")
		return &Note{}, err
	}

	return createdNote, nil
}

func (s *Service) Update(ctx context.Context, ID int64, title *string, content *string) (*Note, error) {

	note, err := s.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	modified := false
	if title != nil && *title != note.Title {
		note.Title = *title
		modified = true
	}
	if content != nil && *content != note.Content {
		note.Content = *content
		modified = true
	}
	if modified {
		note.ModifiedAt = time.Now()
	}

	if err := note.Validate(); err != nil {
		s.logger.Warn(fmt.Sprintf("Note is not valid: %s", err.Error()))
		return nil, err
	}

	note, err = s.store.InsertWithID(ctx, ID, note)
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Could not update Note by ID: %s", err.Error()))
		return nil, err
	}

	return note, nil
}

func (s *Service) DeleteByID(ctx context.Context, ID int64) (bool, error) {
	deleted, err := s.store.DeleteByID(ctx, ID)
	if err != nil {
		s.logger.Error("could not create a new Note in storage")
	}

	return deleted, err
}

func (s *Service) nextID() int64 {
	return 5
}
