package main

import (
	"api/domain/usecase"
	"api/driver"
	"api/driver/model"
	"api/driver/repository"
	"api/util/env"
	"api/util/utillog"
	"context"
	"sync"
)

func init() {
	utillog.NewLogger()
}

func main() {
	env := env.NewEnv()
	if ok := env.Validate(); !ok {
		panic("error init env")
	}

	db, sqlDb := driver.NewDb(env.DB_URL)
	defer sqlDb.Close()
	cache := driver.NewCache(env.CACHE_URL)
	defer cache.Close()

	repository.NewDbRepo[model.User](db)

	usecase.NewWorkerUseCase(
		&sync.Mutex{},
		repository.NewQueueRepo[model.QueueModel](driver.NewQueue(), env.SQS_URL),
		repository.NewCacheRepo[model.CacheModel](cache),
	).Work(context.Background())
}
