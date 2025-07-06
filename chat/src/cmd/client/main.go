package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "chat/driver/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	EnvServerAddr = "SERVER_ADDR"
)

type env struct {
	serverAddr string
}

func NewEnv() env {
	return env{
		serverAddr: os.Getenv(EnvServerAddr),
	}
}

func (e env) validate() error {
	return nil
}

func getFrom() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your name: ")

	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read: %v", err)
		return "(anonymous)"
	}
	trimmed := strings.TrimSpace(name)

	if len(trimmed) == 0 {
		return "(anonymous)"
	}
	return trimmed
}

func newClient(cc grpc.ClientConnInterface, from string) *pb.Client {
	return &pb.Client{
		ChatServiceClient: pb.NewChatServiceClient(cc),
		From:              from,
	}
}

func main() {
	e := NewEnv()
	if err := e.validate(); err != nil {
		panic(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(e.serverAddr, opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		client := newClient(conn, getFrom())
		client.NewClient()
	}()

	<-quit
}
