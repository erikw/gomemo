package storage

import "context"

type Storage[O Object] interface {
	Create(ctx context.Context, o O) (O, error)
	All(ctx context.Context) ([]O, error)
	FindByID(ctx context.Context, ID int64) (O, error)
	InsertWithID(ctx context.Context, ID int64, o O) (O, error)
	DeleteByID(ctx context.Context, ID int64) (bool, error)
}

type Identifiable interface {
	GetID() int64
	SetID(ID int64) error
}

type Object interface {
	Identifiable
}

// Cloneable returns an independent copy of an object.
type Cloneable[O any] interface {
	Clone() O
}

// MemoryObject is an object that Memory can own without exposing mutable state.
type MemoryObject[O any] interface {
	Object
	Cloneable[O]
}
