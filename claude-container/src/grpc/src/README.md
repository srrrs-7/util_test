# Simple gRPC API with Go

This project implements a simple gRPC API with Go, featuring both unary and streaming RPC calls.

## Project Structure

```
grpc-simple-api/
├── api/
│   └── hello.proto          # Protocol buffer definition
├── proto/
│   └── api/
│       ├── hello.pb.go      # Generated Go protobuf code
│       └── hello_grpc.pb.go # Generated Go gRPC code
├── server/
│   └── main.go             # gRPC server implementation
├── client/
│   └── main.go             # gRPC client implementation
├── bin/
│   ├── server              # Compiled server binary
│   └── client              # Compiled client binary
├── go.mod                  # Go module file
└── README.md               # This file
```

## Features

The API provides two RPC methods:

1. **SayHello** - Unary RPC that returns a greeting message
2. **SayHelloStream** - Server streaming RPC that returns multiple greeting messages

## Running the Application

### Prerequisites

- Go 1.24 or later
- Protocol Buffer compiler (protoc)

### Running the Server

```bash
# Build and run the server
go run server/main.go

# Or use the compiled binary
./bin/server
```

The server will start listening on port 50051.

### Running the Client

In another terminal:

```bash
# Build and run the client
go run client/main.go

# Or use the compiled binary
./bin/client
```

The client will:
1. Make a unary RPC call to SayHello
2. Make a streaming RPC call to SayHelloStream and receive multiple responses

## Development

### Regenerating Protocol Buffer Code

If you modify the `.proto` file, regenerate the Go code:

```bash
# Make sure protoc-gen-go and protoc-gen-go-grpc are installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# Generate Go code
protoc --go_out=proto --go_opt=paths=source_relative \
       --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
       api/hello.proto
```

### Building

```bash
# Build both server and client
go build -o bin/server ./server
go build -o bin/client ./client

# Install dependencies
go mod tidy
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test suites
go test ./server -v          # Server unit tests
go test ./client -v          # Client unit tests  
go test . -v                 # Integration tests

# Run tests with coverage
go test -cover ./...
```

## Example Output

When running the client, you should see output like:

```
2024/01/01 12:00:00 Response: Hello, World!
2024/01/01 12:00:00 Streaming responses:
2024/01/01 12:00:00 Stream response: Hello Stream! Message #1
2024/01/01 12:00:01 Stream response: Hello Stream! Message #2
2024/01/01 12:00:02 Stream response: Hello Stream! Message #3
2024/01/01 12:00:03 Stream response: Hello Stream! Message #4
2024/01/01 12:00:04 Stream response: Hello Stream! Message #5
```