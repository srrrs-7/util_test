package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	pb "claude/driver/grpc"
	"claude/util/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type env struct {
	workerNum int
	queueHost string
}

func newEnv() (*env, error) {
	workerNum := os.Getenv(config.ENV_WORKER_NUM)
	num, err := strconv.Atoi(workerNum)
	if err != nil {
		return nil, fmt.Errorf("invalid WORKER_NUM: %w", err)
	}

	return &env{
		workerNum: num,
		queueHost: os.Getenv(config.ENV_QUEUE_HOST),
	}, nil
}

func (e *env) validate() error {
	if e.workerNum <= 0 {
		return fmt.Errorf("WORKER_NUM must be a positive integer")
	}
	if e.queueHost == "" {
		return fmt.Errorf("QUEUE_HOST must be set")
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

	conn, err := grpc.NewClient(
		e.queueHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to queue host %s: %w", e.queueHost, err))
	}
	defer conn.Close()

	client := pb.NewEnqueueServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prompt, err := readPrompt()
	if err != nil {
		panic(fmt.Errorf("failed to read prompt: %w", err))
	}
	fmt.Printf("Read prompt: \n%s\n", prompt)

	resp, err := client.Enqueue(ctx, &pb.EnqueueRequest{
		Prompt: prompt,
	})
	if err != nil {
		panic(fmt.Errorf("failed to enqueue prompt: %w", err))
	}

	fmt.Println("Enqueued prompt result: ", resp.Message)
}

func readPrompt() (string, error) {
	filepath := fmt.Sprintf("%s/%s", config.MASTER_PROMPT_FILE_PATH, config.INSTRUCTION_FILE_NAME)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("prompt file %s does not exist", filepath)
	}

	prompt, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", filepath, err)
	}

	return string(prompt), nil
}
