package repository

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type CacheRepo[T any] struct {
	client *redis.Client
}

func NewCacheRepo[T any](c *redis.Client) CacheRepo[T] {
	return CacheRepo[T]{
		client: c,
	}
}

func (r CacheRepo[T]) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r CacheRepo[T]) Get(ctx context.Context, key string) (*T, error) {
	s, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var j T
	if err = json.Unmarshal([]byte(s), &j); err != nil {
		return nil, err
	}

	return &j, nil
}

func (r CacheRepo[T]) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
