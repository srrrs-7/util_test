package main

import (
	"context"
	"errors"
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

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &helloServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestHelloServer_SayHello(t *testing.T) {
	tests := []struct {
		name        string
		request     *pb.HelloRequest
		expected    string
		expectError bool
	}{
		{
			name:        "Valid request with name",
			request:     &pb.HelloRequest{Name: "Alice"},
			expected:    "Hello, Alice!",
			expectError: false,
		},
		{
			name:        "Valid request with empty name",
			request:     &pb.HelloRequest{Name: ""},
			expected:    "Hello, !",
			expectError: false,
		},
		{
			name:        "Valid request with special characters",
			request:     &pb.HelloRequest{Name: "Bob@123"},
			expected:    "Hello, Bob@123!",
			expectError: false,
		},
		{
			name:        "Valid request with unicode",
			request:     &pb.HelloRequest{Name: "世界"},
			expected:    "Hello, 世界!",
			expectError: false,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.SayHello(ctx, tt.request)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
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

func TestHelloServer_SayHello_Unit(t *testing.T) {
	server := &helloServer{}

	tests := []struct {
		name     string
		request  *pb.HelloRequest
		expected string
	}{
		{
			name:     "Standard greeting",
			request:  &pb.HelloRequest{Name: "Test"},
			expected: "Hello, Test!",
		},
		{
			name:     "Empty name",
			request:  &pb.HelloRequest{Name: ""},
			expected: "Hello, !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := server.SayHello(ctx, tt.request)

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

func TestHelloServer_SayHello_Context(t *testing.T) {
	server := &helloServer{}

	t.Run("Context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := server.SayHello(ctx, &pb.HelloRequest{Name: "Test"})
		
		// The method should still work even with cancelled context
		// as it doesn't explicitly check context status
		if err != nil {
			t.Errorf("Unexpected error with cancelled context: %v", err)
		}
	})

	t.Run("Context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout

		_, err := server.SayHello(ctx, &pb.HelloRequest{Name: "Test"})
		
		// The method should still work as it doesn't check context deadline
		if err != nil {
			t.Errorf("Unexpected error with timed out context: %v", err)
		}
	})
}

func TestHelloServer_SayHello_ErrorCases(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	t.Run("Nil request handling", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// This will actually work in gRPC as it creates an empty request
		resp, err := client.SayHello(ctx, nil)
		if err == nil {
			t.Logf("Nil request handled gracefully, response: %s", resp.GetMessage())
		}
	})

	t.Run("Context deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(1 * time.Millisecond) // Ensure timeout

		_, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Test"})
		if err == nil {
			t.Error("Expected deadline exceeded error")
		} else if !strings.Contains(err.Error(), "deadline exceeded") && !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Errorf("Expected deadline exceeded error, got: %v", err)
		}
	})
}

func TestHelloServer_SayHelloStream(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	tests := []struct {
		name           string
		request        *pb.HelloRequest
		expectedCount  int
		expectedPrefix string
	}{
		{
			name:           "Standard streaming request",
			request:        &pb.HelloRequest{Name: "StreamTest"},
			expectedCount:  5,
			expectedPrefix: "Hello StreamTest! Message #",
		},
		{
			name:           "Empty name streaming",
			request:        &pb.HelloRequest{Name: ""},
			expectedCount:  5,
			expectedPrefix: "Hello ! Message #",
		},
		{
			name:           "Unicode name streaming",
			request:        &pb.HelloRequest{Name: "テスト"},
			expectedCount:  5,
			expectedPrefix: "Hello テスト! Message #",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.SayHelloStream(context.Background(), tt.request)
			if err != nil {
				t.Fatalf("SayHelloStream failed: %v", err)
			}

			var responses []string
			for i := 0; i < tt.expectedCount; i++ {
				resp, err := stream.Recv()
				if err != nil {
					t.Fatalf("Stream receive failed at message %d: %v", i+1, err)
				}
				responses = append(responses, resp.GetMessage())
			}

			// Check that we get EOF after expected messages
			_, err = stream.Recv()
			if err == nil {
				t.Error("Expected EOF after receiving all messages")
			}

			// Verify response count
			if len(responses) != tt.expectedCount {
				t.Errorf("Expected %d responses, got %d", tt.expectedCount, len(responses))
			}

			// Verify response content
			for i, response := range responses {
				expected := tt.expectedPrefix + string(rune('1'+i))
				if response != expected {
					t.Errorf("Response %d: expected '%s', got '%s'", i+1, expected, response)
				}
			}
		})
	}
}

func TestHelloServer_SayHelloStream_Timeout(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	t.Run("Stream with short timeout", func(t *testing.T) {
		// Create context with timeout shorter than expected stream duration
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		stream, err := client.SayHelloStream(ctx, &pb.HelloRequest{Name: "TimeoutTest"})
		if err != nil {
			t.Fatalf("SayHelloStream failed: %v", err)
		}

		// Try to receive messages until timeout
		messageCount := 0
		for {
			_, err := stream.Recv()
			if err != nil {
				if strings.Contains(err.Error(), "deadline exceeded") || strings.Contains(err.Error(), "context deadline exceeded") {
					t.Logf("Received expected timeout after %d messages", messageCount)
					break
				}
				if strings.Contains(err.Error(), "EOF") {
					t.Logf("Stream completed normally with %d messages", messageCount)
					break
				}
				t.Fatalf("Unexpected error: %v", err)
			}
			messageCount++
		}

		if messageCount == 0 {
			t.Error("Expected to receive at least one message before timeout")
		}
	})
}

func TestHelloServer_SayHelloStream_Cancellation(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	t.Run("Cancel stream midway", func(t *testing.T) {
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

		// Try to receive next message - should fail due to cancellation
		_, err = stream.Recv()
		if err == nil {
			t.Error("Expected error after context cancellation")
		} else if !strings.Contains(err.Error(), "context canceled") && !strings.Contains(err.Error(), "canceled") {
			t.Logf("Got error after cancellation (expected): %v", err)
		}
	})
}

// Mock stream for unit testing
type mockHelloServiceSayHelloStreamServer struct {
	grpc.ServerStream
	messages []*pb.HelloResponse
	sendErr  error
}

func (m *mockHelloServiceSayHelloStreamServer) Send(response *pb.HelloResponse) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.messages = append(m.messages, response)
	return nil
}

func TestHelloServer_SayHelloStream_Unit(t *testing.T) {
	server := &helloServer{}

	t.Run("Successful stream", func(t *testing.T) {
		mockStream := &mockHelloServiceSayHelloStreamServer{}
		
		err := server.SayHelloStream(&pb.HelloRequest{Name: "UnitTest"}, mockStream)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(mockStream.messages) != 5 {
			t.Errorf("Expected 5 messages, got %d", len(mockStream.messages))
		}

		for i, msg := range mockStream.messages {
			expected := "Hello UnitTest! Message #" + string(rune('1'+i))
			if msg.GetMessage() != expected {
				t.Errorf("Message %d: expected '%s', got '%s'", i+1, expected, msg.GetMessage())
			}
		}
	})

	t.Run("Stream send error", func(t *testing.T) {
		mockStream := &mockHelloServiceSayHelloStreamServer{
			sendErr: errors.New("mock send error"),
		}
		
		err := server.SayHelloStream(&pb.HelloRequest{Name: "ErrorTest"}, mockStream)
		if err == nil {
			t.Error("Expected error from stream.Send failure")
		}
		
		if err.Error() != "mock send error" {
			t.Errorf("Expected 'mock send error', got: %v", err)
		}
	})
}