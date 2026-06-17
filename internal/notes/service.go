package notes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// TODO move to proper storage layer, first in-memory module
var db = map[int64]Note{
	1: {1, "Title of note", "Some content string here", time.Now(), time.Now()},
}

var ErrTitleRequired = errors.New("the field Title is required")

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Note, error) {
	// TODO pass ctx to DB. Set custom timeout?

	notes := make([]Note, 0, len(db))
	for _, note := range db {
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Service) GetByID(ctx context.Context, ID int64) (Note, error) {
	// TODO pass ctx to DB. Set custom timeout?
	if note, ok := db[ID]; ok {
		return note, nil
	} else {
		return Note{}, fmt.Errorf("could not find note with ID `%d`", ID)
	}
}

func (s *Service) Create(ctx context.Context, title string, content string) (Note, error) {
	// TODO pass ctx to DB. Set custom timeout?

	if strings.TrimSpace(title) == "" {
		return Note{}, ErrTitleRequired
	}

	note := Note{
		ID:         5, // TODO keep track of sequenced ID, with mutex
		Title:      title,
		Content:    content,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	db[note.ID] = note
	return note, nil
}

func (s *Service) DeleteByID(ctx context.Context, ID int64) (bool, error) {
	// TODO pass ctx to DB. Set custom timeout?
	if _, ok := db[ID]; ok {
		delete(db, ID)
		return true, nil
	} else {
		return false, nil
	}
}
