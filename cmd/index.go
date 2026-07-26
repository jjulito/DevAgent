package cmd

import (
	"context"
	"fmt"

	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/jjulito/devagent-cli/internal/provider"
	"github.com/jjulito/devagent-cli/internal/rag"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index [path]",
	Short: "Index source code for semantic search",
	Long:  `Walks the given directory, chunks source files, and indexes them in the vector store (Qdrant) for semantic search.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		if err := AppConfig.Validate(); err != nil {
			return err
		}

		llm, err := provider.New(AppConfig)
		if err != nil {
			return err
		}

		engine := rag.NewQdrant(AppConfig.QdrantHost, AppConfig.QdrantPort, "devagent")
		indexer := rag.NewIndexer(engine, llm)

		output.Info(fmt.Sprintf("Indexing %s ...", path))

		count, err := indexer.IndexDirectory(context.Background(), path)
		if err != nil {
			return err
		}

		output.Success(fmt.Sprintf("Indexed %d chunks", count))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
