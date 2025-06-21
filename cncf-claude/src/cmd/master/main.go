package main

import (
	"claude/config"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type env struct {
	workerNum int
}

func newEnv() (*env, error) {
	workerNum := os.Getenv(config.ENV_WORKER_NUM)
	num, err := strconv.Atoi(workerNum)
	if err != nil {
		return nil, fmt.Errorf("invalid WORKER_NUM: %w", err)
	}

	return &env{
		workerNum: num,
	}, nil
}

func (e *env) validate() error {
	if e.workerNum <= 0 {
		return fmt.Errorf("WORKER_NUM must be a positive integer")
	}
	return nil
}

func main() {
	e, err := newEnv()
	if err != nil {
		panic(err)
	}
	if err := e.validate(); err != nil {
		panic(err)
	}

	var workerFiles string
	for i := range e.workerNum {
		workerFile := fmt.Sprintf(
			config.WORKER_PROMPT_FILE_PATH+config.WORKER_FILE_NAME,
			i+1,
		)
		workerFiles += fmt.Sprintf("%s ", workerFile)
		if _, err := os.Stat(workerFile); os.IsNotExist(err) {
			if err = os.MkdirAll(config.WORKER_PROMPT_FILE_PATH, os.ModePerm); err != nil {
				panic(fmt.Errorf("failed to create directory %s: %w", config.WORKER_PROMPT_FILE_PATH, err))
			}
			if _, err := os.Create(workerFile); err != nil {
				panic(fmt.Errorf("failed to create file %s: %w", workerFile, err))
			}
			fmt.Printf("Created worker file: %s\n", workerFile)
		}
	}

	instructFile := fmt.Sprintf("%s/%s", config.MASTER_PROMPT_FILE_PATH, config.INSTRUCTION_FILE_NAME)
	prompt := fmt.Sprintf(`
		While reading the files in %s, create a task execution plan for %s.
		Since each file is associated with a worker on a 1:1 basis, run the tasks in parallel to create the optimal execution plan.
	`, instructFile, workerFiles)

	if err := exec.Command("claude", "-p", prompt).Run(); err != nil {
		panic(err)
	}
}
