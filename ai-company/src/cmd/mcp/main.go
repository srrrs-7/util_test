package main

import (
	"claude/driver/mcp"
	"fmt"
)

func main() {
	mcpSrv := mcp.MCPServer{}

	fmt.Printf("Starting MCP server...\n")
	if err := mcpSrv.Start(); err != nil {
		fmt.Printf("Failed to start MCP server: %v\n", err)
	}
}
