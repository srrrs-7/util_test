package main

import (
	"claude/config"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type env struct {
	workerID int
}

func newEnv() *env {
	workerID := os.Getenv(config.ENV_WORKER_ID)
	id, err := strconv.Atoi(workerID)
	if err != nil {
		panic(fmt.Errorf("invalid WORKER_ID: %w", err))
	}
	return &env{workerID: id}
}

func (e *env) validate() error {
	if e.workerID <= 0 {
		return fmt.Errorf("WORKER_ID must be a positive integer")
	}
	return nil
}

func main() {
	e := newEnv()
	if err := e.validate(); err != nil {
		panic(err)
	}

	workerFilepath := fmt.Sprintf(
		config.WORKER_PROMPT_FILE_PATH+config.WORKER_FILE_NAME,
		e.workerID,
	)

	prompt := fmt.Sprintf(`
		Read the file in %s every 10 seconds and complete the task written in the file.
	`, workerFilepath)

	if err := exec.Command("claude", "-p", prompt).Run(); err != nil {
		panic(err)
	}
}
