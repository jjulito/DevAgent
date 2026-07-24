package provider

import (
	"fmt"

	"github.com/jjulito/devagent-cli/internal/config"
)

func New(cfg *config.Config) (LLMProvider, error) {
	switch cfg.Provider {
	case "openrouter":
		return NewOpenRouter(cfg.OpenRouterAPIKey, cfg.Model), nil
	case "ollama":
		return nil, fmt.Errorf("provider 'ollama' will be implemented in Phase 2")
	case "openai":
		return nil, fmt.Errorf("provider 'openai' will be implemented in Phase 2")
	case "gemini":
		return nil, fmt.Errorf("provider 'gemini' will be implemented in Phase 2")
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}
