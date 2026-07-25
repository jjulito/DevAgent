package cmd

import (
	"context"
	"fmt"

	"github.com/jjulito/devagent-cli/internal/output"
	"github.com/jjulito/devagent-cli/internal/provider"
	"github.com/jjulito/devagent-cli/internal/review"
	"github.com/spf13/cobra"
)

var reviewStaged bool

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Automated code review from git diff",
	Long:  `Generates an AI-powered code review of the current git diff. Analyzes changes for bugs, style issues, and suggests improvements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppConfig.Validate(); err != nil {
			return err
		}

		llm, err := provider.New(AppConfig)
		if err != nil {
			return err
		}

		output.Divider()
		output.ModelTag(llm.Name(), AppConfig.Model)

		mode := "unstaged"
		if reviewStaged {
			mode = "staged"
		}
		output.Info(fmt.Sprintf("Reviewing %s changes...", mode))
		output.Divider()

		diff, textCh, errCh := review.RunStream(context.Background(), llm, review.Options{
			Staged: reviewStaged,
		})

		if diff != nil {
			output.Dim(fmt.Sprintf("Files changed: %d", diff.FileCount))
			for _, f := range diff.Files {
				output.Dim(fmt.Sprintf("  • %s", f))
			}
			output.Divider()
		}

		for {
			select {
			case text, ok := <-textCh:
				if !ok {
					output.StreamDone()
					output.Divider()
					output.Success("Review completed")
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
	reviewCmd.Flags().BoolVar(&reviewStaged, "staged", false, "review staged changes only")
	rootCmd.AddCommand(reviewCmd)
}
