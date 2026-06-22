package notes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/erikw/gomemo/internal/storage"
)

var ErrTitleRequired = errors.New("the field Title is required")

type Service struct {
	logger *slog.Logger
	store  storage.Storage[*Note]
}

func NewService(logger *slog.Logger, store storage.Storage[*Note]) *Service {
	// TODO hard code else where. Implement fixtures from YAML or such.
	_, _ = store.Create(&Note{
		Title:      "Title of note",
		Content:    "Some content string here",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	})

	return &Service{
		logger: logger,
		store:  store,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*Note, error) {
	// TODO pass ctx to DB. Set custom timeout?

	var notes []*Note
	var err error
	if notes, err = s.store.All(); err != nil {
		s.logger.Error("could not retrieve all Notes")
		return make([]*Note, 0), err
	}

	return notes, nil
}

func (s *Service) GetByID(ctx context.Context, ID int64) (*Note, error) {
	// TODO pass ctx to DB. Set custom timeout?
	var note *Note
	var err error
	if note, err = s.store.FindByID(ID); err != nil {
		s.logger.Error(fmt.Sprintf("could not find Note with ID %d", ID))
		return &Note{}, err
	}

	return note, nil
}

func (s *Service) Create(ctx context.Context, title string, content string) (*Note, error) {
	// TODO pass ctx to DB. Set custom timeout?

	note := &Note{
		Title:      title,
		Content:    content,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if err := s.validate(note); err != nil {
		s.logger.Warn(fmt.Sprintf("Note is not valid: %s", err.Error()))
		return nil, err
	}

	var createdNote *Note
	var err error
	if createdNote, err = s.store.Create(note); err != nil {
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

	if err := s.validate(note); err != nil {
		s.logger.Warn(fmt.Sprintf("Note is not valid: %s", err.Error()))
		return nil, err
	}

	note, err = s.store.InsertWithID(ID, note)
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Could not update Note by ID: %s", err.Error()))
		return nil, err
	}

	return note, nil
}

// TODO in Go design, keep this in service or Note itself?
func (s *Service) validate(note *Note) error {
	if strings.TrimSpace(note.Title) == "" {
		return ErrTitleRequired
	}
	return nil
}

func (s *Service) DeleteByID(ctx context.Context, ID int64) (bool, error) {
	// TODO pass ctx to DB. Set custom timeout?

	deleted, err := s.store.DeleteByID(ID)
	if err != nil {
		s.logger.Error("could not create a new Note in storage")
	}

	return deleted, err
}

func (s *Service) nextID() int64 {
	return 5
}
