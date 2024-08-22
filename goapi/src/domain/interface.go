package domain

import (
	"api/domain/entity"
	"api/driver/model"
	"context"
)

type Queuer[T any] interface {
	EnQueue(ctx context.Context, msg string) (entity.QueueId, error)
	DeQueue(ctx context.Context) (*model.QueueModel[T], error)
	DelQueue(ctx context.Context, id string) error
}

type Cacher[T any] interface {
	Set(ctx context.Context, key entity.QueueId, value T) error
	Get(ctx context.Context, key entity.QueueId) (*T, error)
	Delete(ctx context.Context, key entity.QueueId) error
}

type Querier[T any] interface {
	Select(ctx context.Context, id string) ([]T, error)
	Insert(ctx context.Context, value T) error
	Update(ctx context.Context, value T) error
	Delete(ctx context.Context, value T) error
}
