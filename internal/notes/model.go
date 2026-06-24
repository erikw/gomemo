package notes

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrTitleRequired = errors.New("the field Title is required")

type Note struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func (n *Note) String() string {
	return fmt.Sprintf("Note{id=%d, title=%q}", n.ID, n.Title)
}

func (n *Note) GetID() int64 {
	return n.ID
}

func (n *Note) SetID(ID int64) error {
	n.ID = ID
	return nil
}

func (n *Note) Validate() error {
	if strings.TrimSpace(n.Title) == "" {
		return ErrTitleRequired
	}
	return nil
}
