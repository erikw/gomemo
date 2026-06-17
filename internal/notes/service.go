package notes

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/erikw/gomemo/internal/storage"
)

var ErrTitleRequired = errors.New("the field Title is required")

type Service struct {
	logger *slog.Logger
	store  storage.Storage[Note]
}

func NewService(logger *slog.Logger, store storage.Storage[Note]) *Service {
	store.Create(Note{
		1, "Title of note",
		"Some content string here", time.Now(), time.Now()},
	)

	return &Service{
		logger: logger,
		store:  store,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Note, error) {
	// TODO pass ctx to DB. Set custom timeout?

	var notes []Note
	var err error
	if notes, err = s.store.All(); err != nil {
		s.logger.Error("could not retrieve all Notes")
		return make([]Note, 0, 0), err
	}

	return notes, nil
}

func (s *Service) GetByID(ctx context.Context, ID int64) (Note, error) {
	// TODO pass ctx to DB. Set custom timeout?
	var note Note
	var err error
	if note, err = s.store.FindByID(ID); err != nil {
		s.logger.Error("could not create a new Note in storage")
		return Note{}, err
	}

	return note, nil
}

func (s *Service) Create(ctx context.Context, title string, content string) (Note, error) {
	// TODO pass ctx to DB. Set custom timeout?

	if strings.TrimSpace(title) == "" {
		return Note{}, ErrTitleRequired
	}

	note := Note{
		Title:      title,
		Content:    content,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	var createdNote Note
	var err error
	if createdNote, err = s.store.Create(note); err != nil {
		s.logger.Error("could not create a new Note in storage")
		return Note{}, err
	}

	return createdNote, nil
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
