package storage

import (
	"fmt"
	"log/slog"
)

type Memory[T any] struct {
	logger *slog.Logger
	store  map[int64]T
}

func NewMemory[T any](logger *slog.Logger) *Memory[T] {
	return &Memory[T]{
		logger: logger,
		store:  make(map[int64]T),
	}
}

func (m *Memory[T]) Create(t T) (T, error) {
	// TODO right way is to have interface requiring each storage element to have an ID?
	// e := t
	// e.ID = 5
	// m.store[e.ID] = t // TODO helper to get nextID from last/next ID, mutex
	// return e, nil

	m.store[5] = t
	return t, nil
}

func (m *Memory[T]) FindByID(ID int64) (T, error) {
	if e, ok := m.store[ID]; ok {
		return e, nil
	} else {
		return *new(T), fmt.Errorf("could not find element with ID `%d`", ID)
	}
}

func (m *Memory[T]) DeleteByID(ID int64) (bool, error) {
	if _, ok := m.store[ID]; ok {
		delete(m.store, ID)
		return true, nil
	} else {
		return false, nil
	}
}

func (m *Memory[T]) All() ([]T, error) {
	var elems = make([]T, 0, len(m.store))

	for _, e := range m.store {
		elems = append(elems, e)
	}
	return elems, nil
}
