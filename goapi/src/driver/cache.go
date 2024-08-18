package driver

import "github.com/redis/go-redis/v9"

func NewCache(dsn string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})
}
