package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	client *redis.Client
}

func NewCacheRepo(c *redis.Client) CacheRepo {
	return CacheRepo{
		client: c,
	}
}

func (r CacheRepo) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r CacheRepo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r CacheRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
