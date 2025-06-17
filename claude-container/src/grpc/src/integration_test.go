package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	pb "grpc-simple-api/proto/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// helloServer implementation for integration tests
type integrationHelloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *integrationHelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", req.GetName())
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *integrationHelloServer) SayHelloStream(req *pb.HelloRequest, stream pb.HelloService_SayHelloStreamServer) error {
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

// TestFullGRPCIntegration tests the complete gRPC server and client integration
func TestFullGRPCIntegration(t *testing.T) {
	// Start the actual server
	lis, err := net.Listen("tcp", ":0") // Use random available port
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &integrationHelloServer{})

	// Start server in background
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()
	defer s.Stop()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Connect client to the server
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	t.Run("Integration SayHello", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Integration"})
		if err != nil {
			t.Fatalf("SayHello failed: %v", err)
		}

		expected := "Hello, Integration!"
		if response.GetMessage() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, response.GetMessage())
		}
	})

	t.Run("Integration SayHelloStream", func(t *testing.T) {
		stream, err := client.SayHelloStream(context.Background(), &pb.HelloRequest{Name: "StreamIntegration"})
		if err != nil {
			t.Fatalf("SayHelloStream failed: %v", err)
		}

		var responses []string
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				t.Fatalf("Stream receive failed: %v", err)
			}
			responses = append(responses, resp.GetMessage())
		}

		expectedCount := 5
		if len(responses) != expectedCount {
			t.Errorf("Expected %d responses, got %d", expectedCount, len(responses))
		}

		for i, response := range responses {
			expected := "Hello StreamIntegration! Message #" + string(rune('1'+i))
			if response != expected {
				t.Errorf("Response %d: expected '%s', got '%s'", i+1, expected, response)
			}
		}
	})
}

// TestConcurrentClients tests multiple clients connecting simultaneously
func TestConcurrentClients(t *testing.T) {
	// Start the actual server
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &integrationHelloServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()
	defer s.Stop()

	time.Sleep(100 * time.Millisecond)

	const numClients = 10
	done := make(chan bool, numClients)
	errors := make(chan error, numClients)

	// Start multiple clients concurrently
	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				errors <- err
				return
			}
			defer conn.Close()

			client := pb.NewHelloServiceClient(conn)
			
			// Test unary call
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Client" + string(rune('0'+clientID))})
			if err != nil {
				errors <- err
				return
			}

			expected := "Hello, Client" + string(rune('0'+clientID)) + "!"
			if response.GetMessage() != expected {
				errors <- err
				return
			}

			// Test streaming call
			stream, err := client.SayHelloStream(context.Background(), &pb.HelloRequest{Name: "StreamClient" + string(rune('0'+clientID))})
			if err != nil {
				errors <- err
				return
			}

			messageCount := 0
			for {
				_, err := stream.Recv()
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					errors <- err
					return
				}
				messageCount++
			}

			if messageCount != 5 {
				errors <- err
				return
			}

			done <- true
		}(i)
	}

	// Wait for all clients to complete
	completedClients := 0
	timeout := time.After(30 * time.Second)

	for completedClients < numClients {
		select {
		case <-done:
			completedClients++
		case err := <-errors:
			t.Fatalf("Client error: %v", err)
		case <-timeout:
			t.Fatalf("Test timeout: only %d/%d clients completed", completedClients, numClients)
		}
	}

	t.Logf("Successfully completed %d concurrent client connections", numClients)
}

// TestServerShutdown tests graceful server shutdown
func TestServerShutdown(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &integrationHelloServer{})

	serverStopped := make(chan bool)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server stopped: %v", err)
		}
		serverStopped <- true
	}()

	time.Sleep(100 * time.Millisecond)

	// Connect a client
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	// Make a successful request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "ShutdownTest"})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}

	if response.GetMessage() != "Hello, ShutdownTest!" {
		t.Errorf("Unexpected response: %s", response.GetMessage())
	}

	// Stop the server
	s.Stop()

	// Wait for server to stop
	select {
	case <-serverStopped:
		t.Log("Server stopped gracefully")
	case <-time.After(5 * time.Second):
		t.Error("Server did not stop within timeout")
	}

	// Try to make another request - should fail
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	_, err = client.SayHello(ctx2, &pb.HelloRequest{Name: "AfterShutdown"})
	if err == nil {
		t.Error("Expected error after server shutdown")
	} else {
		t.Logf("Got expected error after shutdown: %v", err)
	}
}

// TestLongRunningStream tests streaming with longer duration
func TestLongRunningStream(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &integrationHelloServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()
	defer s.Stop()

	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	t.Run("Stream with timing verification", func(t *testing.T) {
		start := time.Now()
		
		stream, err := client.SayHelloStream(context.Background(), &pb.HelloRequest{Name: "TimingTest"})
		if err != nil {
			t.Fatalf("SayHelloStream failed: %v", err)
		}

		messageCount := 0
		for {
			_, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				t.Fatalf("Stream receive failed: %v", err)
			}
			messageCount++
		}

		duration := time.Since(start)
		
		if messageCount != 5 {
			t.Errorf("Expected 5 messages, got %d", messageCount)
		}

		// The stream should take at least 4 seconds (5 messages with 1 second delay between)
		// But allow some tolerance for test execution
		expectedMinDuration := 4 * time.Second
		if duration < expectedMinDuration {
			t.Errorf("Stream completed too quickly: %v (expected at least %v)", duration, expectedMinDuration)
		}

		t.Logf("Stream completed in %v with %d messages", duration, messageCount)
	})
}