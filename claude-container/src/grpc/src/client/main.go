package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "grpc-simple-api/proto/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	// Test regular RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("SayHello failed: %v", err)
	}
	log.Printf("Response: %s", response.GetMessage())

	// Test streaming RPC
	stream, err := client.SayHelloStream(context.Background(), &pb.HelloRequest{Name: "Stream"})
	if err != nil {
		log.Fatalf("SayHelloStream failed: %v", err)
	}

	log.Println("Streaming responses:")
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream receive failed: %v", err)
		}
		log.Printf("Stream response: %s", response.GetMessage())
	}
}
