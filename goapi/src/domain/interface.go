package domain

import "context"

type Queuer[T any] interface {
	EnQueue(ctx context.Context, msg string) (string, error)
	DeQueue(ctx context.Context) (*T, error)
	DelQueue(ctx context.Context, id string) error
}

type Cacher[T any] interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (*T, error)
	Delete(ctx context.Context, key string) error
}

type Querier[T any] interface {
	Select(ctx context.Context, id string) ([]T, error)
	Insert(ctx context.Context, value T) error
	Update(ctx context.Context, value T) error
	Delete(ctx context.Context, value T) error
}
