package env

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

const (
	MODE         = "MODE"
	API_PORT     = "API_PORT"
	DB_URL       = "DB_URL"
	CACHE_URL    = "CACHE_URL"
	CACHE_TTL    = "CACHE_TTL"
	CACHE_PREFIX = "CACHE_PREFIX"
	SQS_URL      = "SQS_URL"
)

type Env struct {
	MODE         string
	API_PORT     string
	DB_URL       string
	CACHE_URL    string
	CACHE_TTL    time.Duration
	CACHE_PREFIX string
	SQS_URL      string
}

type env struct {
	MODE         string
	API_PORT     string
	DB_URL       string
	CACHE_URL    string
	CACHE_TTL    string
	CACHE_PREFIX string
	SQS_URL      string
}

func NewEnv() *env {
	return &env{
		MODE:         os.Getenv(MODE),
		API_PORT:     os.Getenv(API_PORT),
		DB_URL:       os.Getenv(DB_URL),
		CACHE_URL:    os.Getenv(CACHE_URL),
		CACHE_TTL:    os.Getenv(CACHE_TTL),
		CACHE_PREFIX: os.Getenv(CACHE_PREFIX),
		SQS_URL:      os.Getenv(SQS_URL),
	}
}

func (e *env) Validate() (*Env, bool) {
	if e.MODE == "" || e.MODE == "debug" || e.MODE == "release" {
		slog.Error("ENV MODE is invalid debug or release", "mode", e.MODE)
		return nil, false
	}
	if e.API_PORT == "" {
		slog.Error("ENV API_PORT is empty")
		return nil, false
	}
	if e.DB_URL == "" {
		slog.Error("ENV DB_URL is empty")
		return nil, false
	}
	if e.CACHE_URL == "" {
		slog.Error("ENV CACHE_URL is empty")
		return nil, false
	}
	if e.CACHE_TTL == "" {
		slog.Error("ENV CACHE_TTL is empty")
		return nil, false
	}
	if e.CACHE_PREFIX == "" {
		slog.Error("ENV CACHE_PREFIX is empty")
		return nil, false
	}
	if e.SQS_URL == "" {
		slog.Error("ENV SQS_URL is empty")
		return nil, false
	}

	ttl, err := strconv.Atoi(e.CACHE_TTL)
	if err != nil {
		slog.Error("ENV CACHE_TTL is invalid", "error", err.Error())
		return nil, false
	}

	var environment Env
	environment.MODE = e.MODE
	environment.API_PORT = e.API_PORT
	environment.DB_URL = e.DB_URL
	environment.CACHE_URL = e.CACHE_URL
	environment.CACHE_TTL = time.Duration(ttl) * time.Second
	environment.CACHE_PREFIX = e.CACHE_PREFIX
	environment.SQS_URL = e.SQS_URL

	return &environment, true
}

func (e *Env) OutPut() string {
	e.DB_URL = "********"
	return fmt.Sprintf("%v", e)
}
