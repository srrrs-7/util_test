package main

import (
	"api/domain/usecase"
	"api/driver"
	"api/driver/model"
	"api/driver/repository"
	"api/util/env"
	"api/util/utillog"
	"context"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func init() {
	utillog.NewLogger()
}

func main() {
	env := env.NewEnv()
	if ok := env.Validate(); !ok {
		slog.Error("error init env", "env", env.OutPut())
		panic("error init env")
	}

	db, sqlDb := driver.NewDb(env.DB_URL)
	defer sqlDb.Close()
	cache := driver.NewCache(env.CACHE_URL)
	defer cache.Close()

	var queue *sqs.Client
	if env.MODE == "debug" {
		queue = driver.NewLocalQueue(env.SQS_URL)
	} else {
		queue = driver.NewQueue()
	}

	repository.NewDbRepo[model.User](db)

	usecase.NewWorkerUseCase(
		&sync.Mutex{},
		repository.NewQueueRepo[model.QueueModel](queue, env.SQS_URL),
		repository.NewCacheRepo[model.CacheModel](cache),
	).Work(context.Background())
	slog.Info("start worker")
}
