package storage

// TODO all funcs should take context as first arg
// TODO cant use `any`, we need interface that they have an ID.
type Storage[T any] interface {
	Create(t T) (T, error)
	All() ([]T, error)
	FindByID(id int64) (T, error)
	DeleteByID(id int64) (bool, error)
}
