package cmd

import (
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the LLM with project context",
	Long:  `Send a question to the configured language model and show the response in the terminal.`,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(askCmd)
}
