package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	pb "claude/driver/grpc"
	"claude/util/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	WAIT_TIME = 3 * time.Second
)

type env struct {
	workerID  int
	queueHost string
}

func newEnv() *env {
	workerID := os.Getenv(config.ENV_WORKER_ID)
	id, err := strconv.Atoi(workerID)
	if err != nil {
		panic(fmt.Errorf("invalid WORKER_ID: %w", err))
	}
	return &env{
		workerID:  id,
		queueHost: os.Getenv(config.ENV_QUEUE_HOST)}
}

func (e *env) validate() error {
	if e.workerID <= 0 {
		return fmt.Errorf("WORKER_ID must be a positive integer")
	}
	if e.queueHost == "" {
		return fmt.Errorf("QUEUE_HOST must be set")
	}
	return nil
}

func main() {
	e := newEnv()
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

	client := pb.NewDequeueServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		resp, err := client.Dequeue(ctx, &pb.DequeueRequest{})
		if err != nil {
			fmt.Printf("Error calling Dequeue: %v\n", err)
		}
		if resp == nil || resp.Prompt == "" {
			fmt.Println("No tasks available, waiting...")
			time.Sleep(WAIT_TIME)
			continue
		}
		fmt.Printf("Received prompt: \n%s\n", resp.Prompt)

		out, err := exec.Command("claude", "-p", resp.Prompt).Output()
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
		fmt.Printf("Command output: \n%s\n", string(out))

		time.Sleep(WAIT_TIME)
	}
}
