package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/jjulito/devagent-cli/internal/provider"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the LLM with project context",
	Long:  `Send a question to the configured language model and show the response in the terminal.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppConfig.Validate(); err != nil {
			return err
		}

		llm, err := provider.New(AppConfig)
		if err != nil {
			return err
		}

		question := strings.Join(args, " ")

		output.Divider()
		output.ModelTag(llm.Name(), AppConfig.Model)
		output.Info("Thinking...")
		output.Divider()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		messages := []provider.Message{
			{
				Role:    provider.RoleSystem,
				Content: "You are an expert assistant for developers. Respond clearly, concisely, and usefully. If asked about code, include examples when relevant.",
			},
			{
				Role:    provider.RoleUser,
				Content: question,
			},
		}

		textCh, errCh := llm.ChatStream(ctx, messages)

		for {
			select {
			case text, ok := <-textCh:
				if !ok {
					output.StreamDone()
					output.Divider()
					output.Success("Completed")
					return nil
				}
				output.StreamChunk(text)
			case err, ok := <-errCh:
				if ok && err != nil {
					return err
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
