package cmd

import (
	"os"

	"github.com/jjulito/devagent-cli/internal/config"
	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	flagProvider string
	flagModel    string
	verbose      bool

	AppConfig *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "devagent",
	Short: "CLI for AI-powered developer tools",
	Long:  `DevAgent is a modular CLI that allows you to analyze code, perform code reviews, semantic search, and more, using your own API key (BYOK).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if flagProvider != "" {
			cfg.Provider = flagProvider
		}
		if flagModel != "" {
			cfg.Model = flagModel
		}
		cfg.Verbose = verbose

		AppConfig = cfg
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "configuration file")
	rootCmd.PersistentFlags().StringVarP(&flagProvider, "provider", "p", "", "LLM provider (openrouter, ollama, openai, gemini)")
	rootCmd.PersistentFlags().StringVarP(&flagModel, "model", "m", "", "specific model to use")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output")
}
