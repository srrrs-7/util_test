# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a simple gRPC API implementation in Go featuring both unary and streaming RPC calls. The project demonstrates basic gRPC service patterns with a HelloService that provides greeting functionality.

## Commands

### Running the Application
```bash
# Run server (starts on port 50051)
go run server/main.go

# Run client (in separate terminal)
go run client/main.go

# Use compiled binaries
./bin/server
./bin/client
```

### Building
```bash
# Build both server and client
go build -o bin/server ./server
go build -o bin/client ./client

# Install/update dependencies
go mod tidy
```

### Protocol Buffer Code Generation
```bash
# Install required tools first
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# Generate Go code from proto files
protoc --go_out=proto --go_opt=paths=source_relative \
       --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
       api/hello.proto
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

# Run single test
go test -run TestHelloServer_SayHello ./server
```

## Architecture

The project follows standard gRPC patterns:

- **API Definition**: `api/hello.proto` defines the HelloService with SayHello (unary) and SayHelloStream (server streaming) RPCs
- **Generated Code**: `proto/api/` contains auto-generated Go protobuf and gRPC code
- **Server**: `server/main.go` implements the HelloService interface with both unary and streaming handlers
- **Client**: `client/main.go` demonstrates calling both RPC methods

The server implements `pb.UnimplementedHelloServiceServer` and registers with a gRPC server on port 50051. The streaming RPC sends 5 messages with 1-second delays between each.

## Go Module

Module name: `grpc-simple-api`
Go version: 1.24+
Key dependencies: google.golang.org/grpc, google.golang.org/protobuf