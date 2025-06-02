package service

import (
	"concurrency/domain"
	"context"
	"time"
)

type QueueRepo interface {
	Enqueue(context.Context, string) (domain.QueueID, error)
	Dequeue(context.Context) error
}

type CacheRepo interface {
	Set(context.Context, string, string, time.Duration) error
	Get(context.Context, string) (string, error)
	Delete(context.Context, string) error
	MakeKey(context.Context, string, string) string
}
