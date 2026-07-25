package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaProvider struct {
	host   string
	model  string
	client *http.Client
}

func NewOllama(host, model string) *OllamaProvider {
	return &OllamaProvider{
		host:   host,
		model:  model,
		client: &http.Client{},
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) chatURL() string {
	return p.host + "/v1/chat/completions"
}

func (p *OllamaProvider) Chat(ctx context.Context, messages []Message) (*Response, error) {
	body := orRequest{
		Model:    p.model,
		Messages: toORMessages(messages),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.chatURL(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed (is Ollama running at %s?): %w", p.host, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var orResp orResponse
	if err := json.Unmarshal(respBody, &orResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(orResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &Response{
		Content:    orResp.Choices[0].Message.Content,
		TokensUsed: orResp.Usage.TotalTokens,
		Model:      p.model,
	}, nil
}

func (p *OllamaProvider) ChatStream(ctx context.Context, messages []Message) (<-chan string, <-chan error) {
	textCh := make(chan string, 64)
	errCh := make(chan error, 1)

	go func() {
		defer close(textCh)
		defer close(errCh)

		body := orRequest{
			Model:    p.model,
			Messages: toORMessages(messages),
			Stream:   true,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			errCh <- fmt.Errorf("marshal request: %w", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.chatURL(), bytes.NewReader(jsonBody))
		if err != nil {
			errCh <- fmt.Errorf("create request: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		parseSSEStream(p.client, req, textCh, errCh)
	}()

	return textCh, errCh
}
