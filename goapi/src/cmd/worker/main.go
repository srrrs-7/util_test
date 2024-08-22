package main

import (
	"api/domain/entity"
	"api/domain/usecase"
	"api/driver"
	"api/driver/model"
	"api/driver/repository"
	"api/handle/request"
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
	env, ok := env.NewEnv().Validate()
	if !ok {
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
		repository.NewQueueRepo[model.QueueModel[request.Params]](queue, env.SQS_URL),
		repository.NewCacheRepo[entity.CheckStatusEnt](cache, env.CACHE_TTL, env.CACHE_PREFIX),
	).Work(context.Background())

	slog.Info("start worker")
}
