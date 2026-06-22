package storage

// TODO all funcs should take context as first arg
// TODO cant use `any`, we need interface that they have an ID.
type Storage[O Object] interface {
	Create(o O) (O, error)
	All() ([]O, error)
	FindByID(ID int64) (O, error)
	InsertWithID(ID int64, o O) (O, error)
	DeleteByID(ID int64) (bool, error)
}

type Identifiable interface {
	GetID() int64
	SetID(ID int64) error
}

type Object interface {
	Identifiable
}
