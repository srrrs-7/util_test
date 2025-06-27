# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Docker Services
- `make up`: Start claude-master, claude-worker, and queue services with 1 worker instance
- `make build`: Build all Docker images and clean up dangling images
- `make down`: Stop all services
- `make gopher`: Run gopher development container with shell access
- `make claude-master`: Execute claude-master service
- `make claude-worker`: Access claude-worker container shell
- `make mcp`: Run MCP server

### Protocol Buffers
- `make grpc`: Generate Go protobuf files from proto definitions in `src/driver/grpc/proto/`

## Architecture Overview

This is a distributed AI task processing system built in Go with the following components:

### Core Services
1. **Queue Service** (`src/cmd/queue/main.go`): gRPC server managing task queues on port 8080
2. **Master Service** (`src/cmd/master/main.go`): Reads prompts from `/prompt/tasks/instruction.md` and enqueues them
3. **Worker Service** (`src/cmd/worker/main.go`): Continuously dequeues tasks and executes them using the `claude` CLI
4. **MCP Service** (`src/cmd/mcp/main.go`): Model Context Protocol server with tool capabilities

### Key Components
- **gRPC Layer** (`src/driver/grpc/`): Handles Enqueue, Dequeue, and QueueStatus operations
- **Queue Implementation** (`src/util/queue/`): In-memory task queue management
- **Configuration** (`src/util/config/static.go`): Environment variables and file paths
- **MCP Integration** (`src/driver/mcp/server.go`): Uses mark3labs/mcp-go library

### Environment Variables
- `WORKER_NUM`: Number of workers to spawn (master service)
- `WORKER_ID`: Unique worker identifier (worker service)  
- `QUEUE_HOST`: Queue service address (default: "queue:8080")

### File Structure
- `/prompt/tasks/instruction.md`: Master service reads prompts from here
- `/prompt/worker_tasks/`: Worker service prompt directory
- `src/driver/grpc/proto/queue.proto`: Protocol buffer definitions

### Dependencies
- Go 1.24.2
- gRPC for service communication
- mark3labs/mcp-go for MCP protocol support
- Google Protocol Buffers

## Development Notes

The system uses Docker Compose for orchestration. Workers continuously poll the queue and execute tasks using the external `claude` CLI command. The MCP server provides a "hello_world" tool as an example implementation.

Master service validates that `WORKER_NUM` is positive and `QUEUE_HOST` is set. Workers validate `WORKER_ID` is positive and `QUEUE_HOST` is configured.