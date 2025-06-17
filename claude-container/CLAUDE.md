# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

This is a containerized development environment setup for Claude Code. The project uses Docker Compose to create a sandbox environment with Claude Code pre-installed.

### Container Structure
- Base image: `oven/bun` with additional development tools (vim, curl, tmux)
- Claude Code is installed globally via `bun install -g @anthropic-ai/claude-code`
- Working directory: `/src` (mounted from host `./src` directory)
- Port 8080 exposed for development services

## Development Commands

### Container Management
- `make up` - Start the Claude container in detached mode
- `make claude` - Execute bash shell inside the running Claude container

### Environment Setup
The tmux script at `src/scritps/tmux.sh` sets up a multi-pane tmux session for development workflow.

## Notes
- There's a typo in the Makefile: `cluade` should be `claude` in the `make up` command
- The project structure is minimal, focused on providing a containerized Claude Code environment