# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This project is a Kubernetes development environment using Kind (Kubernetes in Docker) with Docker Compose orchestration. The setup creates containerized Kubernetes clusters for local development and testing.

Key components:
- **Docker Compose services**: Two identical Kind cluster services (`cluster-1`, `cluster-2`) that run in privileged containers
- **Kind configuration**: Multi-node cluster with 1 control plane and 1 worker node
- **Kubernetes manifests**: Sample nginx deployment for testing cluster functionality
- **Custom Docker image**: Based on Docker-in-Docker with kubectl, kind, and helm pre-installed

## Common Commands

### Cluster Management
```bash
# Start the cluster environment
make cluster-up

# Build the Docker images
make cluster-build

# Stop the cluster environment
make cluster-down

# View cluster logs
make cluster-logs
```

### Kubernetes Operations
```bash
# Create a Kind cluster
make cluster

# Delete the Kind cluster
make delete

# Access the cluster container
make exec

# View all Kubernetes resources
make get-all
```

### Node Management
```bash
# Get nodes as JSON
make get-nodes

# Get detailed node information
make get-nodes-detail
```

### Pod Management
```bash
# Deploy nginx test pod
make apply-pod

# Get pods as JSON
make get-pods

# Get detailed pod information
make get-pods-detail

# Describe nginx deployment
make describe-deploy

# Describe nginx pods
make describe-pod

# View pod resource usage
make top-pods
```

## Development Workflow

1. Use `make cluster-up` to start the Docker Compose services
2. Use `make cluster` to create the Kind cluster inside the container
3. Use `make apply-pod` to deploy test workloads
4. Use `make exec` to access the cluster container for manual kubectl operations
5. Use `make cluster-down` when finished

The `/src` directory is mounted into both cluster containers, allowing you to modify Kubernetes manifests and cluster configuration without rebuilding images.