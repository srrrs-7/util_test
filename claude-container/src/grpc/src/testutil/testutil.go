package testutil

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	pb "grpc-simple-api/proto/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const BufSize = 1024 * 1024

// TestServer provides utilities for testing gRPC servers
type TestServer struct {
	listener *bufconn.Listener
	server   *grpc.Server
	cleanup  func()
}

// NewTestServer creates a new test server with the given HelloService implementation
func NewTestServer(t *testing.T, service pb.HelloServiceServer) *TestServer {
	listener := bufconn.Listen(BufSize)
	server := grpc.NewServer()
	pb.RegisterHelloServiceServer(server, service)

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("Test server stopped: %v", err)
		}
	}()

	cleanup := func() {
		server.Stop()
		listener.Close()
	}

	return &TestServer{
		listener: listener,
		server:   server,
		cleanup:  cleanup,
	}
}

// Close stops the test server
func (ts *TestServer) Close() {
	if ts.cleanup != nil {
		ts.cleanup()
	}
}

// NewClient creates a new gRPC client connected to the test server
func (ts *TestServer) NewClient() (pb.HelloServiceClient, *grpc.ClientConn, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return ts.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewHelloServiceClient(conn)
	return client, conn, nil
}

// MockHelloServer provides a configurable mock server for testing
type MockHelloServer struct {
	pb.UnimplementedHelloServiceServer
	
	// SayHello configuration
	HelloResponse *pb.HelloResponse
	HelloError    error
	HelloDelay    time.Duration
	
	// SayHelloStream configuration
	StreamResponses []*pb.HelloResponse
	StreamError     error
	StreamDelay     time.Duration
	
	// Call tracking
	HelloCallCount  int
	StreamCallCount int
	LastHelloRequest  *pb.HelloRequest
	LastStreamRequest *pb.HelloRequest
}

// SayHello implements the HelloService.SayHello method
func (m *MockHelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	m.HelloCallCount++
	m.LastHelloRequest = req
	
	if m.HelloDelay > 0 {
		time.Sleep(m.HelloDelay)
	}
	
	if m.HelloError != nil {
		return nil, m.HelloError
	}
	
	if m.HelloResponse != nil {
		return m.HelloResponse, nil
	}
	
	// Default response
	return &pb.HelloResponse{
		Message: "Hello, " + req.GetName() + "!",
	}, nil
}

// SayHelloStream implements the HelloService.SayHelloStream method
func (m *MockHelloServer) SayHelloStream(req *pb.HelloRequest, stream pb.HelloService_SayHelloStreamServer) error {
	m.StreamCallCount++
	m.LastStreamRequest = req
	
	if m.StreamError != nil {
		return m.StreamError
	}
	
	responses := m.StreamResponses
	if responses == nil {
		// Default responses
		responses = []*pb.HelloResponse{
			{Message: "Hello " + req.GetName() + "! Message #1"},
			{Message: "Hello " + req.GetName() + "! Message #2"},
			{Message: "Hello " + req.GetName() + "! Message #3"},
			{Message: "Hello " + req.GetName() + "! Message #4"},
			{Message: "Hello " + req.GetName() + "! Message #5"},
		}
	}
	
	for _, response := range responses {
		if m.StreamDelay > 0 {
			time.Sleep(m.StreamDelay)
		}
		
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	
	return nil
}

// Reset clears all call tracking data
func (m *MockHelloServer) Reset() {
	m.HelloCallCount = 0
	m.StreamCallCount = 0
	m.LastHelloRequest = nil
	m.LastStreamRequest = nil
}

// AssertHelloCallCount verifies the number of SayHello calls
func (m *MockHelloServer) AssertHelloCallCount(t *testing.T, expected int) {
	if m.HelloCallCount != expected {
		t.Errorf("Expected %d SayHello calls, got %d", expected, m.HelloCallCount)
	}
}

// AssertStreamCallCount verifies the number of SayHelloStream calls
func (m *MockHelloServer) AssertStreamCallCount(t *testing.T, expected int) {
	if m.StreamCallCount != expected {
		t.Errorf("Expected %d SayHelloStream calls, got %d", expected, m.StreamCallCount)
	}
}

// TestStreamReceiver is a helper for testing streaming responses
type TestStreamReceiver struct {
	Messages []string
	Error    error
}

// ReceiveAll receives all messages from a stream until EOF or error
func (r *TestStreamReceiver) ReceiveAll(stream pb.HelloService_SayHelloStreamClient) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			r.Error = err
			return
		}
		r.Messages = append(r.Messages, resp.GetMessage())
	}
}

// AssertMessageCount verifies the expected number of messages received
func (r *TestStreamReceiver) AssertMessageCount(t *testing.T, expected int) {
	if len(r.Messages) != expected {
		t.Errorf("Expected %d messages, got %d", expected, len(r.Messages))
	}
}

// AssertNoError verifies that no error occurred during streaming
func (r *TestStreamReceiver) AssertNoError(t *testing.T) {
	if r.Error != nil {
		t.Errorf("Unexpected streaming error: %v", r.Error)
	}
}

// AssertMessages verifies the exact messages received
func (r *TestStreamReceiver) AssertMessages(t *testing.T, expected []string) {
	if len(r.Messages) != len(expected) {
		t.Errorf("Expected %d messages, got %d", len(expected), len(r.Messages))
		return
	}
	
	for i, expectedMsg := range expected {
		if i >= len(r.Messages) {
			t.Errorf("Missing message at index %d", i)
			continue
		}
		if r.Messages[i] != expectedMsg {
			t.Errorf("Message %d: expected '%s', got '%s'", i, expectedMsg, r.Messages[i])
		}
	}
}

// Common test error types
var (
	ErrMockServer     = errors.New("mock server error")
	ErrMockStream     = errors.New("mock stream error")
	ErrMockConnection = errors.New("mock connection error")
)

// WithTimeout creates a context with timeout for tests
func WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// DefaultTestTimeout is the default timeout for test operations
const DefaultTestTimeout = 5 * time.Second

// WithDefaultTimeout creates a context with the default test timeout
func WithDefaultTimeout() (context.Context, context.CancelFunc) {
	return WithTimeout(DefaultTestTimeout)
}