package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct{}

func (s *MCPServer) Start() error {
	mcpSrv := server.NewMCPServer(
		"mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	mcpSrv.AddTool(mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	), s.helloHandler)

	if err := server.ServeStdio(mcpSrv); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	return nil
}

// TODO: implement Resource and Prompts

func (s *MCPServer) helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
