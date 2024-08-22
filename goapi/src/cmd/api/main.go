package main

import (
	"api/domain/entity"
	"api/domain/usecase"
	"api/driver"
	"api/driver/model"
	"api/driver/repository"
	"api/handle"
	"api/handle/request"
	"api/util/env"
	"api/util/utillog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	cacheRepo := repository.NewCacheRepo[entity.CheckStatusEnt](cache, env.CACHE_TTL, env.CACHE_PREFIX)
	queueRepo := repository.NewQueueRepo[model.QueueModel[request.Params]](queue, env.SQS_URL)

	r := handle.NewServer(
		usecase.NewCheckUseCase(cacheRepo),
		usecase.NewCreateUseCase(queueRepo, cacheRepo),
	).Routing()

	go func() {
		slog.Info("Server Start")
		if err := http.ListenAndServe(env.API_PORT, r); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Shutdown Server")
}
