package storage

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Memory[O Object] struct {
	logger *slog.Logger
	store  map[int64]O
	mu     sync.Mutex
	nextID int64
}

func NewMemory[O Object](logger *slog.Logger) *Memory[O] {
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
	m.store[m.nextID] = o

	m.nextID++

	return o, nil
}

func (m *Memory[O]) FindByID(ctx context.Context, ID int64) (O, error) {
	if o, ok := m.store[ID]; ok {
		return o, nil
	} else {
		return *new(O), fmt.Errorf("could not find element with ID `%d`", ID)
	}
}

func (m *Memory[O]) InsertWithID(ctx context.Context, ID int64, o O) (O, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// m.logger.Debug(fmt.Sprintf("Inserting with ID %d obj %v", ID, o))
	m.store[ID] = o

	return o, nil
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
	var objs = make([]O, 0, len(m.store))

	for _, o := range m.store {
		objs = append(objs, o)
	}
	return objs, nil
}

func (m *Memory[O]) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = make(map[int64]O)
	m.nextID = 0
	return nil
}
