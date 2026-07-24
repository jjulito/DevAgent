package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

type OpenRouterProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewOpenRouter(apiKey, model string) *OpenRouterProvider {
	return &OpenRouterProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

func (p *OpenRouterProvider) Chat(ctx context.Context, messages []Message) (*Response, error) {
	body := orRequest{
		Model:    p.model,
		Messages: toORMessages(messages),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var orResp orResponse
	if err := json.Unmarshal(respBody, &orResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if orResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", orResp.Error.Message)
	}

	if len(orResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &Response{
		Content:    orResp.Choices[0].Message.Content,
		TokensUsed: orResp.Usage.TotalTokens,
		Model:      orResp.Model,
	}, nil
}

func (p *OpenRouterProvider) ChatStream(ctx context.Context, messages []Message) (<-chan string, <-chan error) {
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterURL, bytes.NewReader(jsonBody))
		if err != nil {
			errCh <- fmt.Errorf("create request: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.apiKey)

		resp, err := p.client.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("request failed: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			errCh <- fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
			return
		}

		buf := make([]byte, 4096)
		reader := resp.Body
		var remainder string

		for {
			n, readErr := reader.Read(buf)
			if n > 0 {
				lines := strings.Split(remainder+string(buf[:n]), "\n")
				remainder = lines[len(lines)-1]

				for _, line := range lines[:len(lines)-1] {
					line = strings.TrimSpace(line)
					if !strings.HasPrefix(line, "data: ") {
						continue
					}
					data := strings.TrimPrefix(line, "data: ")
					if data == "[DONE]" {
						return
					}

					var chunk orStreamChunk
					if err := json.Unmarshal([]byte(data), &chunk); err != nil {
						continue
					}
					if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
						textCh <- chunk.Choices[0].Delta.Content
					}
				}
			}
			if readErr != nil {
				if readErr != io.EOF {
					errCh <- readErr
				}
				return
			}
		}
	}()

	return textCh, errCh
}

type orRequest struct {
	Model    string      `json:"model"`
	Messages []orMessage `json:"messages"`
	Stream   bool        `json:"stream,omitempty"`
}

type orMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type orResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Model string `json:"model"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type orStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func toORMessages(msgs []Message) []orMessage {
	out := make([]orMessage, len(msgs))
	for i, m := range msgs {
		out[i] = orMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}
	return out
}
