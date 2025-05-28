package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo[T any] struct {
	client *redis.Client
}

func NewCacheRepo[T any](dns string) CacheRepo[T] {
	return CacheRepo[T]{
		client: redis.NewClient(&redis.Options{
			Addr:     dns,
			Password: "",
			DB:       0,
		}),
	}
}

func (r CacheRepo[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, j, ttl).Err()
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

func (r CacheRepo[T]) MakeKey(ctx context.Context, prefix, suffix string) string {
	return prefix + ":" + suffix
}
