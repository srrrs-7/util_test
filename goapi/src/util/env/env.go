package env

import (
	"fmt"
	"os"
)

const (
	MODE      = "MODE"
	API_PORT  = "API_PORT"
	DB_URL    = "DB_URL"
	CACHE_URL = "CACHE_URL"
	SQS_URL   = "SQS_URL"
)

type Env struct {
	MODE      string
	API_PORT  string
	DB_URL    string
	CACHE_URL string
	SQS_URL   string
}

func NewEnv() Env {
	return Env{
		MODE:      os.Getenv(MODE),
		API_PORT:  os.Getenv(API_PORT),
		DB_URL:    os.Getenv(DB_URL),
		CACHE_URL: os.Getenv(CACHE_URL),
		SQS_URL:   os.Getenv(SQS_URL),
	}
}

func (e Env) Validate() bool {
	if e.MODE == "" || e.API_PORT == "" || e.DB_URL == "" || e.CACHE_URL == "" || e.SQS_URL == "" {
		return false
	}
	return true
}

func (e *Env) OutPut() string {
	e.DB_URL = "********"
	return fmt.Sprintf("%v", e)
}
