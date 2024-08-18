package main

import (
	"api/domain/usecase"
	"api/driver"
	"api/driver/model"
	"api/driver/repository"
	"api/handle"
	"api/util/env"
	"api/util/utillog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	cacheRepo := repository.NewCacheRepo[model.CacheModel](cache)
	queueRepo := repository.NewQueueRepo[model.QueueModel](driver.NewQueue(), env.SQS_URL)

	r := handle.NewServer(
		usecase.NewCheckUseCase(cacheRepo),
		usecase.NewCreateUseCase(queueRepo, cacheRepo),
	).Routing()

	go func() {
		if err := http.ListenAndServe(env.API_PORT, r); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Shutdown Server ...")
}
