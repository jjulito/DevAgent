package review

import (
	"context"
	"fmt"
	"time"

	"github.com/jjulito/devagent-cli/internal/git"
	"github.com/jjulito/devagent-cli/internal/provider"
)

type Options struct {
	Staged bool
	Path   string
	Lang   string
}

type Result struct {
	Review    string
	FileCount int
	Files     []string
}

const systemPrompt = `You are a senior code reviewer. Analyze the following git diff and provide:

1. **Summary**: Brief overview of the changes
2. **Issues**: Bugs, security risks, or logic errors (if any)
3. **Suggestions**: Improvements for readability, performance, or best practices
4. **Rating**: Overall quality (Excellent / Good / Needs Work / Critical Issues)

Be concise and actionable. Focus on what matters.`

func Run(ctx context.Context, llm provider.LLMProvider, opts Options) (*Result, error) {
	diff, err := git.GetDiff(git.DiffOptions{
		Staged: opts.Staged,
		Path:   opts.Path,
	})
	if err != nil {
		return nil, err
	}

	maxDiff := diff.Raw
	if len(maxDiff) > 15000 {
		maxDiff = maxDiff[:15000] + "\n\n... [truncated, diff too large]"
	}

	userPrompt := fmt.Sprintf("Review this diff (%d files changed):\n\n```diff\n%s\n```", diff.FileCount, maxDiff)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: systemPrompt},
		{Role: provider.RoleUser, Content: userPrompt},
	}

	reviewCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	resp, err := llm.Chat(reviewCtx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM review failed: %w", err)
	}

	return &Result{
		Review:    resp.Content,
		FileCount: diff.FileCount,
		Files:     diff.Files,
	}, nil
}

func RunStream(ctx context.Context, llm provider.LLMProvider, opts Options) (*git.DiffResult, <-chan string, <-chan error) {
	diff, err := git.GetDiff(git.DiffOptions{
		Staged: opts.Staged,
		Path:   opts.Path,
	})
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return nil, nil, errCh
	}

	maxDiff := diff.Raw
	if len(maxDiff) > 15000 {
		maxDiff = maxDiff[:15000] + "\n\n... [truncated, diff too large]"
	}

	userPrompt := fmt.Sprintf("Review this diff (%d files changed):\n\n```diff\n%s\n```", diff.FileCount, maxDiff)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: systemPrompt},
		{Role: provider.RoleUser, Content: userPrompt},
	}

	reviewCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	_ = cancel

	textCh, errCh := llm.ChatStream(reviewCtx, messages)
	return diff, textCh, errCh
}
