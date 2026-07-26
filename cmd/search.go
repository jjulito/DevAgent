package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/jjulito/devagent-cli/internal/provider"
	"github.com/jjulito/devagent-cli/internal/rag"
	"github.com/spf13/cobra"
)

var searchTopK int

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Semantic search over indexed code",
	Long:  `Searches the indexed codebase using semantic similarity. Returns the most relevant code chunks matching your query.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

		if err := AppConfig.Validate(); err != nil {
			return err
		}

		llm, err := provider.New(AppConfig)
		if err != nil {
			return err
		}

		engine := rag.NewQdrant(AppConfig.QdrantHost, AppConfig.QdrantPort, "devagent")
		retriever := rag.NewRetriever(engine, llm)

		output.Info(fmt.Sprintf("Searching: %q", query))
		output.Divider()

		results, err := retriever.Search(context.Background(), query, searchTopK)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			output.Warn("No results found. Try indexing first: devagent index .")
			return nil
		}

		for i, r := range results {
			fmt.Printf("\n%s #%d — %s (score: %.3f)\n",
				"📄", i+1, r.Document.FilePath, r.Score)
			output.Dim(fmt.Sprintf("  Language: %s | Chunk: %d", r.Document.Language, r.Document.Chunk))

			lines := strings.Split(r.Document.Content, "\n")
			preview := lines
			if len(preview) > 10 {
				preview = preview[:10]
			}
			for _, line := range preview {
				fmt.Printf("  %s\n", line)
			}
			if len(lines) > 10 {
				output.Dim(fmt.Sprintf("  ... (%d more lines)", len(lines)-10))
			}
		}

		output.Divider()
		output.Success(fmt.Sprintf("Found %d results", len(results)))
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchTopK, "top", 5, "number of results to return")
	rootCmd.AddCommand(searchCmd)
}
