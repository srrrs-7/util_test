package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-simple-api/proto/api"

	"google.golang.org/grpc"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", req.GetName())
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *helloServer) SayHelloStream(req *pb.HelloRequest, stream pb.HelloService_SayHelloStreamServer) error {
	log.Printf("Stream request received: %v", req.GetName())

	for i := 0; i < 5; i++ {
		response := &pb.HelloResponse{
			Message: fmt.Sprintf("Hello %s! Message #%d", req.GetName(), i+1),
		}

		if err := stream.Send(response); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &helloServer{})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
