package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	client *redis.Client
}

func NewCacheRepo(dns string) CacheRepo {
	return CacheRepo{
		client: redis.NewClient(&redis.Options{
			Addr:     dns,
			Password: "",
			DB:       0,
		}),
	}
}

func (r CacheRepo) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, j, ttl).Err()
}

func (r CacheRepo) Get(ctx context.Context, key string) (string, error) {
	s, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return s, nil
}

func (r CacheRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r CacheRepo) MakeKey(ctx context.Context, prefix, suffix string) string {
	return prefix + ":" + suffix
}
