package storage

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Memory is a concurrency-safe in-memory store that clones values at its boundaries.
//
// The store owns the instances in its map. Reads and write results return independent
// copies so callers cannot mutate map-owned values after an operation completes.
type Memory[O MemoryObject[O]] struct {
	logger *slog.Logger
	store  map[int64]O
	mu     sync.RWMutex
	nextID int64
}

func NewMemory[O MemoryObject[O]](logger *slog.Logger) *Memory[O] {
	return &Memory[O]{
		logger: logger,
		store:  make(map[int64]O),
		nextID: 0,
	}
}

func (m *Memory[O]) Create(ctx context.Context, o O) (O, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_ = o.SetID(m.nextID)
	m.store[m.nextID] = o.Clone()

	m.nextID++

	return o.Clone(), nil
}

func (m *Memory[O]) FindByID(ctx context.Context, ID int64) (O, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if o, ok := m.store[ID]; ok {
		return o.Clone(), nil
	} else {
		return *new(O), fmt.Errorf("could not find element with ID `%d`", ID)
	}
}

func (m *Memory[O]) InsertWithID(ctx context.Context, ID int64, o O) (O, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[ID] = o.Clone()

	return o.Clone(), nil
}

func (m *Memory[O]) DeleteByID(ctx context.Context, ID int64) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.store[ID]; ok {
		delete(m.store, ID)
		return true, nil
	} else {
		return false, nil
	}
}

func (m *Memory[O]) All(ctx context.Context) ([]O, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var objs = make([]O, 0, len(m.store))

	for _, o := range m.store {
		objs = append(objs, o.Clone())
	}
	return objs, nil
}
