package provider

import "context"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Content    string `json:"content"`
	TokensUsed int    `json:"tokens_used,omitempty"`
	Model      string `json:"model,omitempty"`
}

type LLMProvider interface {
	Chat(ctx context.Context, messages []Message) (*Response, error)
	ChatStream(ctx context.Context, messages []Message) (<-chan string, <-chan error)
	Name() string
}
