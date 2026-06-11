package notes

import (
	"fmt"
	"log/slog"
	"time"
)

var db = map[int64]Note{
	1: {1, "Title of note", "Some content string here", time.Now(), time.Now()},
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// TODO return pointer to Note?
func (s *Service) GetByID(ID int64) (Note, error) {
	if note, ok := db[ID]; ok {
		return note, nil
	} else {
		return Note{}, fmt.Errorf("could not find note with ID `%d`", ID)
	}
}
