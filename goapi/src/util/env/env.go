package env

import (
	"fmt"
	"os"
)

const (
	DB_URL    = "DB_URL"
	CACHE_URL = "CACHE_URL"
	SQS_URL   = "SQS_URL"
)

type Env struct {
	DB_URL    string
	CACHE_URL string
	SQS_URL   string
}

func NewEnv() Env {
	return Env{
		DB_URL:    os.Getenv(DB_URL),
		CACHE_URL: os.Getenv(CACHE_URL),
		SQS_URL:   os.Getenv(SQS_URL),
	}
}

func (e Env) Validate() bool {
	if e.DB_URL == "" || e.CACHE_URL == "" || e.SQS_URL == "" {
		return false
	}
	return true
}

func (e *Env) OutPut() string {
	e.DB_URL = "********"
	return fmt.Sprintf("%v", e)
}
