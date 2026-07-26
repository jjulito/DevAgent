package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/jjulito/devagent-cli/internal/provider"
)

type Retriever struct {
	engine VectorEngine
	llm    provider.LLMProvider
}

func NewRetriever(engine VectorEngine, llm provider.LLMProvider) *Retriever {
	return &Retriever{engine: engine, llm: llm}
}

func (r *Retriever) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	if topK == 0 {
		topK = 5
	}

	results, err := r.engine.Search(ctx, query, topK)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return results, nil
}

func (r *Retriever) Ask(ctx context.Context, query string, topK int) (*provider.Response, error) {
	results, err := r.Search(ctx, query, topK)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no relevant code found, try indexing first: devagent index .")
	}

	var contextParts []string
	for _, res := range results {
		contextParts = append(contextParts, fmt.Sprintf("--- %s (score: %.2f) ---\n%s",
			res.Document.FilePath, res.Score, res.Document.Content))
	}

	codeContext := strings.Join(contextParts, "\n\n")

	messages := []provider.Message{
		{
			Role:    provider.RoleSystem,
			Content: "You are a code assistant. Answer the user's question using ONLY the provided code context. Be specific and reference file paths when applicable.",
		},
		{
			Role:    provider.RoleUser,
			Content: fmt.Sprintf("Code context:\n\n%s\n\nQuestion: %s", codeContext, query),
		},
	}

	return r.llm.Chat(ctx, messages)
}
