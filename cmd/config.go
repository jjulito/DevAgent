package cmd

import (
	"fmt"

	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show active configuration",
	Long:  `Displays the current configuration including the active provider, model, and connection details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := AppConfig

		output.Banner()
		fmt.Println()
		output.Info("Active Configuration")
		output.Divider()

		fmt.Printf("  Provider:    %s\n", cfg.Provider)
		fmt.Printf("  Model:       %s\n", cfg.Model)
		fmt.Printf("  Verbose:     %v\n", cfg.Verbose)
		fmt.Println()

		output.Dim("  API Keys:")
		if cfg.OpenRouterAPIKey != "" {
			fmt.Println("    OpenRouter: ✔ configured")
		} else {
			fmt.Println("    OpenRouter: ✖ not set")
		}
		if cfg.OpenAIAPIKey != "" {
			fmt.Println("    OpenAI:     ✔ configured")
		} else {
			fmt.Println("    OpenAI:     ✖ not set")
		}
		if cfg.GeminiAPIKey != "" {
			fmt.Println("    Gemini:     ✔ configured")
		} else {
			fmt.Println("    Gemini:     ✖ not set")
		}
		fmt.Println()

		output.Dim("  Connections:")
		fmt.Printf("    Ollama:  %s\n", cfg.OllamaHost)
		fmt.Printf("    Qdrant:  %s:%d\n", cfg.QdrantHost, cfg.QdrantPort)
		fmt.Printf("    MCP:     port %d\n", cfg.MCPPort)

		output.Divider()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
