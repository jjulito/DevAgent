package rag

import "context"

type Document struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
	Language string `json:"language"`
	Chunk    int    `json:"chunk"`
}

type SearchResult struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"`
}

type VectorEngine interface {
	Index(ctx context.Context, docs []Document) error
	Search(ctx context.Context, query string, topK int) ([]SearchResult, error)
	Delete(ctx context.Context, collection string) error
	CollectionExists(ctx context.Context, name string) (bool, error)
}
