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
