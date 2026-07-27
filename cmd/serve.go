package cmd

import (
	"fmt"

	"github.com/jjulito/devagent-cli/internal/mcp"
	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start MCP server exposing project tools",
	Long:  `Starts a Model Context Protocol (MCP) server that exposes filesystem tools (read_file, list_dir, search_files, grep) to AI agents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := servePort
		if port == 0 {
			port = AppConfig.MCPPort
		}

		output.Banner()
		fmt.Println()
		output.Info(fmt.Sprintf("Starting MCP server on port %d...", port))
		output.Divider()
		output.Dim("Tools: read_file, list_dir, search_files, grep")
		output.Dim("Transport: SSE (Server-Sent Events)")
		output.Divider()

		return mcp.StartServer(mcp.ServerOptions{Port: port})
	},
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 0, "port to listen on (default from config)")
	rootCmd.AddCommand(serveCmd)
}
