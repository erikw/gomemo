package notes

import (
	"testing"
	"time"
)

func TestNoteValidate(t *testing.T) {
	tests := []struct {
		name    string
		note    *Note
		wantErr bool
	}{
		{
			name: "valid note with title",
			note: &Note{
				Title:   "Test Title",
				Content: "Some content",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			note: &Note{
				Title:   "",
				Content: "Some content",
			},
			wantErr: true,
		},
		{
			name: "whitespace-only title",
			note: &Note{
				Title:   "   ",
				Content: "Some content",
			},
			wantErr: true,
		},
		{
			name: "valid title with content",
			note: &Note{
				Title:      "My Note",
				Content:    "",
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.note.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err != ErrTitleRequired {
				t.Errorf("Validate() error = %v, want ErrTitleRequired", err)
			}
		})
	}
}
