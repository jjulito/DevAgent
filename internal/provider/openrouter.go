package provider

import (
	"net/http"
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
