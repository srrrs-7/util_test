package main

import (
	"api/driver"
	"api/driver/model"
	"api/driver/repository"
	"api/util/env"
	"api/util/utillog"
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
	repository.NewCacheRepo[model.CacheModel](cache)
	repository.NewQueueRepo[model.QueueModel](driver.NewQueue(), env.SQS_URL)

}
