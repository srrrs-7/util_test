package main

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	pb "grpc-simple-api/proto/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// Mock server for testing client
type mockHelloServer struct {
	pb.UnimplementedHelloServiceServer
	responses []string
	streamResponses [][]string
	errors map[string]error
}

func (s *mockHelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	if err, exists := s.errors["SayHello"]; exists {
		return nil, err
	}
	
	message := "Hello, " + req.GetName() + "!"
	if len(s.responses) > 0 {
		message = s.responses[0]
		s.responses = s.responses[1:]
	}
	
	return &pb.HelloResponse{Message: message}, nil
}

func (s *mockHelloServer) SayHelloStream(req *pb.HelloRequest, stream pb.HelloService_SayHelloStreamServer) error {
	if err, exists := s.errors["SayHelloStream"]; exists {
		return err
	}

	responses := []string{
		"Hello " + req.GetName() + "! Message #1",
		"Hello " + req.GetName() + "! Message #2",
		"Hello " + req.GetName() + "! Message #3",
		"Hello " + req.GetName() + "! Message #4",
		"Hello " + req.GetName() + "! Message #5",
	}
	
	if len(s.streamResponses) > 0 {
		responses = s.streamResponses[0]
		s.streamResponses = s.streamResponses[1:]
	}

	for _, response := range responses {
		if err := stream.Send(&pb.HelloResponse{Message: response}); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond) // Short delay to simulate streaming
	}
	
	return nil
}

func setupTestServer(mockServer pb.HelloServiceServer) func() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, mockServer)
	
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	
	return func() {
		s.Stop()
		lis.Close()
	}
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func getTestClient() (pb.HelloServiceClient, *grpc.ClientConn, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", 
		grpc.WithContextDialer(bufDialer), 
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	
	client := pb.NewHelloServiceClient(conn)
	return client, conn, nil
}

func TestClientSayHello(t *testing.T) {
	mockServer := &mockHelloServer{}
	cleanup := setupTestServer(mockServer)
	defer cleanup()

	tests := []struct {
		name        string
		request     *pb.HelloRequest
		expected    string
		expectError bool
	}{
		{
			name:        "Successful hello request",
			request:     &pb.HelloRequest{Name: "ClientTest"},
			expected:    "Hello, ClientTest!",
			expectError: false,
		},
		{
			name:        "Empty name request",
			request:     &pb.HelloRequest{Name: ""},
			expected:    "Hello, !",
			expectError: false,
		},
		{
			name:        "Unicode name request",
			request:     &pb.HelloRequest{Name: "世界"},
			expected:    "Hello, 世界!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, conn, err := getTestClient()
			if err != nil {
				t.Fatalf("Failed to create test client: %v", err)
			}
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.SayHello(ctx, tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if resp.GetMessage() != tt.expected {
				t.Errorf("Expected message '%s', got '%s'", tt.expected, resp.GetMessage())
			}
		})
	}
}

func TestClientSayHello_WithTimeout(t *testing.T) {
	mockServer := &mockHelloServer{}
	cleanup := setupTestServer(mockServer)
	defer cleanup()

	client, conn, err := getTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer conn.Close()

	t.Run("Request with very short timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(1 * time.Millisecond) // Ensure timeout

		_, err := client.SayHello(ctx, &pb.HelloRequest{Name: "TimeoutTest"})
		if err == nil {
			t.Error("Expected timeout error")
		} else if !strings.Contains(err.Error(), "deadline exceeded") && !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Errorf("Expected deadline exceeded error, got: %v", err)
		}
	})
}

func TestClientSayHelloStream(t *testing.T) {
	mockServer := &mockHelloServer{}
	cleanup := setupTestServer(mockServer)
	defer cleanup()

	client, conn, err := getTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer conn.Close()

	tests := []struct {
		name          string
		request       *pb.HelloRequest
		expectedCount int
		expectedPrefix string
	}{
		{
			name:          "Successful stream request",
			request:       &pb.HelloRequest{Name: "StreamTest"},
			expectedCount: 5,
			expectedPrefix: "Hello StreamTest! Message #",
		},
		{
			name:          "Empty name stream",
			request:       &pb.HelloRequest{Name: ""},
			expectedCount: 5,
			expectedPrefix: "Hello ! Message #",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.SayHelloStream(context.Background(), tt.request)
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

			if len(responses) != tt.expectedCount {
				t.Errorf("Expected %d responses, got %d", tt.expectedCount, len(responses))
			}

			for i, response := range responses {
				expected := tt.expectedPrefix + string(rune('1'+i))
				if response != expected {
					t.Errorf("Response %d: expected '%s', got '%s'", i+1, expected, response)
				}
			}
		})
	}
}

func TestClientSayHelloStream_WithCancellation(t *testing.T) {
	mockServer := &mockHelloServer{}
	cleanup := setupTestServer(mockServer)
	defer cleanup()

	client, conn, err := getTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer conn.Close()

	t.Run("Cancel stream after first message", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		stream, err := client.SayHelloStream(ctx, &pb.HelloRequest{Name: "CancelTest"})
		if err != nil {
			t.Fatalf("SayHelloStream failed: %v", err)
		}

		// Receive first message
		_, err = stream.Recv()
		if err != nil {
			t.Fatalf("Failed to receive first message: %v", err)
		}

		// Cancel the context
		cancel()

		// Try to receive more messages - should eventually fail
		messageCount := 1
		for messageCount < 10 { // Prevent infinite loop
			_, err = stream.Recv()
			if err != nil {
				if strings.Contains(err.Error(), "context canceled") || strings.Contains(err.Error(), "canceled") {
					t.Logf("Stream properly cancelled after %d messages", messageCount)
					return
				}
				if err.Error() == "EOF" {
					t.Log("Stream completed normally before cancellation took effect")
					return
				}
				t.Fatalf("Unexpected error: %v", err)
			}
			messageCount++
		}
		
		t.Error("Expected stream to be cancelled or complete")
	})
}

func TestClientConnectionFailure(t *testing.T) {
	t.Run("Connection to non-existent server", func(t *testing.T) {
		conn, err := grpc.Dial("localhost:99999", 
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithTimeout(100*time.Millisecond))
		
		if err == nil {
			conn.Close()
			t.Error("Expected connection error to non-existent server")
		} else {
			t.Logf("Got expected connection error: %v", err)
		}
	})
}

func TestClientWithMockServerErrors(t *testing.T) {
	mockServer := &mockHelloServer{
		errors: map[string]error{
			"SayHello": context.DeadlineExceeded,
		},
	}
	cleanup := setupTestServer(mockServer)
	defer cleanup()

	client, conn, err := getTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer conn.Close()

	t.Run("Server returns error for SayHello", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.SayHello(ctx, &pb.HelloRequest{Name: "ErrorTest"})
		if err == nil {
			t.Error("Expected error from server")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}

func TestClientWithCustomResponses(t *testing.T) {
	mockServer := &mockHelloServer{
		responses: []string{"Custom response 1", "Custom response 2"},
	}
	cleanup := setupTestServer(mockServer)
	defer cleanup()

	client, conn, err := getTestClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer conn.Close()

	t.Run("Custom server responses", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// First request
		resp1, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Test1"})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if resp1.GetMessage() != "Custom response 1" {
			t.Errorf("Expected 'Custom response 1', got '%s'", resp1.GetMessage())
		}

		// Second request
		resp2, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Test2"})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if resp2.GetMessage() != "Custom response 2" {
			t.Errorf("Expected 'Custom response 2', got '%s'", resp2.GetMessage())
		}
	})
}