package main

import (
	"log"
	"net"
	"os"

	pb "chat/driver/grpc"

	"google.golang.org/grpc"
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

func main() {
	e := NewEnv()
	if err := e.validate(); err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", e.serverAddr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	defer listener.Close()
	log.Printf("Server running at %s", e.serverAddr)

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, pb.New())

	grpcServer.Serve(listener)
}
