package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type QdrantEngine struct {
	host       string
	port       int
	collection string
	client     *http.Client
}

func NewQdrant(host string, port int, collection string) *QdrantEngine {
	return &QdrantEngine{
		host:       host,
		port:       port,
		collection: collection,
		client:     &http.Client{},
	}
}

func (q *QdrantEngine) baseURL() string {
	return fmt.Sprintf("http://%s:%d", q.host, q.port)
}

func (q *QdrantEngine) CollectionExists(ctx context.Context, name string) (bool, error) {
	url := fmt.Sprintf("%s/collections/%s", q.baseURL(), name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("qdrant connection failed: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (q *QdrantEngine) ensureCollection(ctx context.Context) error {
	exists, err := q.CollectionExists(ctx, q.collection)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	url := fmt.Sprintf("%s/collections/%s", q.baseURL(), q.collection)
	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     384,
			"distance": "Cosine",
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(req)
	if err != nil {
		return fmt.Errorf("create collection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create collection error (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (q *QdrantEngine) Index(ctx context.Context, docs []Document) error {
	if err := q.ensureCollection(ctx); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/collections/%s/points", q.baseURL(), q.collection)

	type point struct {
		ID      string                 `json:"id"`
		Vector  []float32              `json:"vector"`
		Payload map[string]interface{} `json:"payload"`
	}

	batchSize := 100
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}

		var points []point
		for _, doc := range docs[i:end] {
			vec := simpleHash(doc.Content, 384)
			points = append(points, point{
				ID:     doc.ID,
				Vector: vec,
				Payload: map[string]interface{}{
					"content":   doc.Content,
					"file_path": doc.FilePath,
					"language":  doc.Language,
					"chunk":     doc.Chunk,
				},
			})
		}

		body := map[string]interface{}{"points": points}
		jsonBody, _ := json.Marshal(body)

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonBody))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := q.client.Do(req)
		if err != nil {
			return fmt.Errorf("index batch failed: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("index error (status %d)", resp.StatusCode)
		}
	}

	return nil
}

func (q *QdrantEngine) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	url := fmt.Sprintf("%s/collections/%s/points/search", q.baseURL(), q.collection)

	vec := simpleHash(query, 384)

	body := map[string]interface{}{
		"vector":       vec,
		"limit":        topK,
		"with_payload": true,
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Result []struct {
			ID      string                 `json:"id"`
			Score   float64                `json:"score"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse search results: %w", err)
	}

	var results []SearchResult
	for _, r := range result.Result {
		doc := Document{
			ID: r.ID,
		}
		if v, ok := r.Payload["content"].(string); ok {
			doc.Content = v
		}
		if v, ok := r.Payload["file_path"].(string); ok {
			doc.FilePath = v
		}
		if v, ok := r.Payload["language"].(string); ok {
			doc.Language = v
		}

		results = append(results, SearchResult{
			Document: doc,
			Score:    r.Score,
		})
	}

	return results, nil
}

func (q *QdrantEngine) Delete(ctx context.Context, collection string) error {
	url := fmt.Sprintf("%s/collections/%s", q.baseURL(), collection)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return fmt.Errorf("delete collection failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func simpleHash(text string, dims int) []float32 {
	vec := make([]float32, dims)
	for i, ch := range text {
		vec[i%dims] += float32(ch) * 0.001
	}

	var norm float32
	for _, v := range vec {
		norm += v * v
	}
	if norm > 0 {
		norm = float32(1.0 / float64(norm))
		for i := range vec {
			vec[i] *= norm
		}
	}

	return vec
}
