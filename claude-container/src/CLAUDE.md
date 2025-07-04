# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a multi-component repository containing:
1. **gRPC Service** - A simple Go-based gRPC API with unary and streaming RPC calls
2. **Orchestrator** - A tmux-based orchestration system for managing multiple Claude Code instances in separate panes

## Commands

### Main Orchestrator
```bash
# Start the tmux orchestration environment
make orchestra
# or directly:
./orchestrator/scripts/claude.sh
```

### gRPC Service (from /grpc/src/)
```bash
# Navigate to gRPC directory first
cd grpc/src

# Run server (starts on port 50051)
go run server/main.go

# Run client (in separate terminal)
go run client/main.go

# Build both components
go build -o bin/server ./server
go build -o bin/client ./client

# Install/update dependencies
go mod tidy

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test suites
go test ./server -v          # Server unit tests
go test ./client -v          # Client unit tests  
go test . -v                 # Integration tests

# Run single test
go test -run TestHelloServer_SayHello ./server

# Generate protobuf code (if proto files change)
# First install required tools:
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# Then generate:
protoc --go_out=proto --go_opt=paths=source_relative \
       --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
       api/hello.proto
```

## Architecture

### Repository Structure
- **Root**: Contains Makefile for orchestration entry point
- **grpc/src/**: Complete gRPC service implementation in Go
- **orchestrator/**: tmux-based multi-pane management system

### gRPC Service Architecture
- **API Definition**: `grpc/src/api/hello.proto` defines HelloService with SayHello (unary) and SayHelloStream (server streaming) RPCs
- **Generated Code**: `grpc/src/proto/api/` contains auto-generated protobuf and gRPC code
- **Server**: `grpc/src/server/main.go` implements HelloService interface
- **Client**: `grpc/src/client/main.go` demonstrates RPC method calls
- **Testing**: Unit tests for both server and client components, plus integration tests
- **Test Utilities**: `grpc/src/testutil/testutil.go` provides TestServer (in-memory gRPC server), MockHelloServer, and TestStreamReceiver for comprehensive testing

### Orchestrator System
The orchestrator creates a tmux session with multiple panes and launches Claude Code instances in each pane. It uses:
- Custom tmux configuration with vim-style keybindings (prefix: C-x)
- Fish shell as default
- Manager-worker pattern for task distribution across panes
- Custom status bar with system information

## Go Module Details

**Module**: `grpc-simple-api` (Go 1.24+)
**Key Dependencies**: 
- google.golang.org/grpc v1.58.3
- google.golang.org/protobuf v1.31.0

## Working Directory Context

When working on gRPC components, change to `/workspace/grpc/src/` directory first as it contains the Go module and all related files. This is where you'll find:
- `go.mod` and `go.sum` for dependency management
- All Go source files (server/, client/, testutil/)
- The `api/` directory with protobuf definitions
- Generated code in `proto/api/`

The server runs on port 50051 by default for local development.