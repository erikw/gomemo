package notes

import (
	"fmt"
	"time"
)

var db = map[int64]Note{
	1: {1, "Title of note", "Some content string here", time.Now(), time.Now()},
}

// TODO return pointer to Note?
func GetByID(ID int64) (Note, error) {
	if note, ok := db[ID]; ok {
		return note, nil
	} else {
		return Note{}, fmt.Errorf("could not find note with ID `%d`", ID)
	}
}
