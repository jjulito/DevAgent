package mcp

import (
	"context"
	"fmt"
	"log"
	"net/http"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerOptions struct {
	Port int
}

func StartServer(opts ServerOptions) error {
	server := gomcp.NewServer(
		&gomcp.Implementation{
			Name:    "devagent-mcp",
			Version: "1.0.0",
		},
		nil,
	)

	gomcp.AddTool(server,
		&gomcp.Tool{
			Name:        "read_file",
			Description: "Read the contents of a file",
		},
		func(ctx context.Context, req *gomcp.CallToolRequest, args ReadFileArgs) (*gomcp.CallToolResult, any, error) {
			content, err := HandleReadFile(ctx, args)
			if err != nil {
				return nil, nil, err
			}
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: content}},
			}, nil, nil
		},
	)

	gomcp.AddTool(server,
		&gomcp.Tool{
			Name:        "list_dir",
			Description: "List the contents of a directory",
		},
		func(ctx context.Context, req *gomcp.CallToolRequest, args ListDirArgs) (*gomcp.CallToolResult, any, error) {
			content, err := HandleListDir(ctx, args)
			if err != nil {
				return nil, nil, err
			}
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: content}},
			}, nil, nil
		},
	)

	gomcp.AddTool(server,
		&gomcp.Tool{
			Name:        "search_files",
			Description: "Search for files matching a glob pattern",
		},
		func(ctx context.Context, req *gomcp.CallToolRequest, args SearchFilesArgs) (*gomcp.CallToolResult, any, error) {
			content, err := HandleSearchFiles(ctx, args)
			if err != nil {
				return nil, nil, err
			}
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: content}},
			}, nil, nil
		},
	)

	gomcp.AddTool(server,
		&gomcp.Tool{
			Name:        "grep",
			Description: "Search for a text pattern in files",
		},
		func(ctx context.Context, req *gomcp.CallToolRequest, args GrepArgs) (*gomcp.CallToolResult, any, error) {
			content, err := HandleGrep(ctx, args)
			if err != nil {
				return nil, nil, err
			}
			return &gomcp.CallToolResult{
				Content: []gomcp.Content{&gomcp.TextContent{Text: content}},
			}, nil, nil
		},
	)

	addr := fmt.Sprintf(":%d", opts.Port)
	handler := gomcp.NewSSEHandler(func(*http.Request) *gomcp.Server { return server }, nil)

	log.Printf("MCP Server listening on http://localhost%s", addr)
	return http.ListenAndServe(addr, handler)
}
