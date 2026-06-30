package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erikw/gomemo/internal/notes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteStore struct {
	pool *pgxpool.Pool
}

func NewNoteStore(pool *pgxpool.Pool) *NoteStore {
	return &NoteStore{pool: pool}
}

func (s *NoteStore) Create(ctx context.Context, note *notes.Note) (*notes.Note, error) {
	row := s.pool.QueryRow(
		ctx,
		`INSERT INTO notes (title, content, created_at, modified_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, title, content, created_at, modified_at`,
		note.Title,
		note.Content,
		note.CreatedAt,
		note.ModifiedAt,
	)

	created := &notes.Note{}
	if err := scanNote(row, created); err != nil {
		return nil, fmt.Errorf("could not create note: %w", err)
	}
	return created, nil
}

func (s *NoteStore) All(ctx context.Context) ([]*notes.Note, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT id, title, content, created_at, modified_at
		 FROM notes
		 ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query notes: %w", err)
	}
	defer rows.Close()

	notesList := make([]*notes.Note, 0)
	for rows.Next() {
		n := &notes.Note{}
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.ModifiedAt); err != nil {
			return nil, fmt.Errorf("could not scan note row: %w", err)
		}
		notesList = append(notesList, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate note rows: %w", err)
	}

	return notesList, nil
}

func (s *NoteStore) FindByID(ctx context.Context, ID int64) (*notes.Note, error) {
	row := s.pool.QueryRow(
		ctx,
		`SELECT id, title, content, created_at, modified_at
		 FROM notes
		 WHERE id = $1`,
		ID,
	)

	n := &notes.Note{}
	if err := scanNote(row, n); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("could not find element with ID `%d`", ID)
		}
		return nil, fmt.Errorf("could not find note by ID: %w", err)
	}
	return n, nil
}

func (s *NoteStore) InsertWithID(ctx context.Context, ID int64, note *notes.Note) (*notes.Note, error) {
	row := s.pool.QueryRow(
		ctx,
		`INSERT INTO notes (id, title, content, created_at, modified_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO UPDATE
		 SET title = EXCLUDED.title,
		     content = EXCLUDED.content,
		     created_at = EXCLUDED.created_at,
		     modified_at = EXCLUDED.modified_at
		 RETURNING id, title, content, created_at, modified_at`,
		ID,
		note.Title,
		note.Content,
		note.CreatedAt,
		note.ModifiedAt,
	)

	updated := &notes.Note{}
	if err := scanNote(row, updated); err != nil {
		return nil, fmt.Errorf("could not upsert note with ID `%d`: %w", ID, err)
	}
	return updated, nil
}

func (s *NoteStore) DeleteByID(ctx context.Context, ID int64) (bool, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM notes WHERE id = $1`, ID)
	if err != nil {
		return false, fmt.Errorf("could not delete note with ID `%d`: %w", ID, err)
	}
	return tag.RowsAffected() > 0, nil
}

func (s *NoteStore) Clear(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `TRUNCATE TABLE notes RESTART IDENTITY`); err != nil {
		return fmt.Errorf("could not clear notes table: %w", err)
	}
	return nil
}

func (s *NoteStore) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.pool.Ping(ctx)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanNote(s scanner, note *notes.Note) error {
	return s.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.ModifiedAt)
}
