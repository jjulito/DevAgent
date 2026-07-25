package provider

import (
	"fmt"

	"github.com/jjulito/devagent-cli/internal/config"
)

func New(cfg *config.Config) (LLMProvider, error) {
	switch cfg.Provider {
	case "openrouter":
		return NewOpenRouter(cfg.OpenRouterAPIKey, cfg.Model), nil
	case "openai":
		return NewOpenAI(cfg.OpenAIAPIKey, cfg.Model), nil
	case "gemini":
		return NewGemini(cfg.GeminiAPIKey, cfg.Model), nil
	case "ollama":
		return NewOllama(cfg.OllamaHost, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}
