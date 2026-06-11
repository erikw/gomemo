package notes

import (
	"fmt"
	"time"
)

type Note struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func (n Note) String() string {
	return fmt.Sprintf("Note{id=%d, title=%q}", n.ID, n.Title)
}
