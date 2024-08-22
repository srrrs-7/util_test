package repository

import (
	"api/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo[T any] struct {
	client *redis.Client
	ttl    time.Duration
	prefix string
}

func NewCacheRepo[T any](c *redis.Client, t time.Duration, p string) CacheRepo[T] {
	return CacheRepo[T]{
		client: c,
		ttl:    t,
		prefix: p,
	}
}

func (r CacheRepo[T]) Set(ctx context.Context, key entity.QueueId, value T) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}

	k := fmt.Sprintf("%s:%s", r.prefix, key.String())
	return r.client.Set(ctx, k, j, r.ttl).Err()
}

func (r CacheRepo[T]) Get(ctx context.Context, key entity.QueueId) (*T, error) {
	k := fmt.Sprintf("%s:%s", r.prefix, key.String())
	s, err := r.client.Get(ctx, k).Result()
	if err != nil {
		return nil, err
	}

	var j T
	if err = json.Unmarshal([]byte(s), &j); err != nil {
		return nil, err
	}

	return &j, nil
}

func (r CacheRepo[T]) Delete(ctx context.Context, key entity.QueueId) error {
	k := fmt.Sprintf("%s:%s", r.prefix, key.String())
	return r.client.Del(ctx, k).Err()
}
