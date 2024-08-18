package domain

import "context"

type Queuer interface {
	Send(msg string) (string, error)
	Receive() (string, error)
	Delete(id string) error
}

type Cacher interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Querier[T any] interface {
	Select(ctx context.Context, id string) ([]T, error)
	Insert(ctx context.Context, value T) error
	Update(ctx context.Context, value T) error
	Delete(ctx context.Context, value T) error
}
