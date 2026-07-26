package rag

import (
	"context"
	"fmt"

	"github.com/jjulito/devagent-cli/internal/provider"
)

type Indexer struct {
	engine VectorEngine
	llm    provider.LLMProvider
}

func NewIndexer(engine VectorEngine, llm provider.LLMProvider) *Indexer {
	return &Indexer{engine: engine, llm: llm}
}

func (idx *Indexer) IndexDirectory(ctx context.Context, path string) (int, error) {
	docs, err := ChunkDirectory(ChunkOptions{
		RootPath: path,
		MaxLines: 60,
		Overlap:  10,
	})
	if err != nil {
		return 0, fmt.Errorf("chunking failed: %w", err)
	}

	if len(docs) == 0 {
		return 0, fmt.Errorf("no indexable files found in %s", path)
	}

	if err := idx.engine.Index(ctx, docs); err != nil {
		return 0, fmt.Errorf("indexing failed: %w", err)
	}

	return len(docs), nil
}
