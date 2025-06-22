package main

import (
	"fmt"
	"net"
	"os"

	pb "claude/driver/grpc"
	"claude/util/config"
	"claude/util/queue"

	"google.golang.org/grpc"
)

type env struct {
	queueHost string
}

func newEnv() *env {
	queueHost := os.Getenv(config.ENV_QUEUE_HOST)
	if queueHost == "" {
		panic(fmt.Errorf("QUEUE_HOST environment variable is not set"))
	}

	return &env{queueHost: queueHost}
}

func (e *env) validate() error {
	if e.queueHost == "" {
		return fmt.Errorf("QUEUE_HOST must be set")
	}
	return nil
}

func init() {
	// Initialize the global scope queue
	queue.NewQueue()
}

func main() {
	e := newEnv()
	if err := e.validate(); err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", e.queueHost)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen: %v", err))
	}
	defer l.Close()

	grpcSrv := grpc.NewServer()

	pb.RegisterDequeueServiceServer(grpcSrv, &pb.Dequeue{})
	pb.RegisterEnqueueServiceServer(grpcSrv, &pb.Enqueue{})
	pb.RegisterQueueStatusServiceServer(grpcSrv, &pb.QueueStatus{})

	fmt.Printf("Queue server is running on %s\n", e.queueHost)

	// graceful shutdown
	defer grpcSrv.GracefulStop()

	if err := grpcSrv.Serve(l); err != nil {
		panic(fmt.Sprintf("Failed to serve: %v", err))
	}
}
