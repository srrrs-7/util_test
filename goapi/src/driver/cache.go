package driver

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewCache(dsn string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(dsn),
		Password: "",
		DB:       0,
	})
}
