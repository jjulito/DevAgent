package config

type Config struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
	Verbose  bool   `mapstructure:"verbose"`

	OpenRouterAPIKey string `mapstructure:"openrouter_api_key"`
	OpenAIAPIKey     string `mapstructure:"openai_api_key"`
	GeminiAPIKey     string `mapstructure:"gemini_api_key"`
	OllamaHost       string `mapstructure:"ollama_host"`

	QdrantHost string `mapstructure:"qdrant_host"`
	QdrantPort int    `mapstructure:"qdrant_port"`
	MCPPort    int    `mapstructure:"mcp_port"`
}

var defaultModels = map[string]string{
	"openrouter": "meta-llama/llama-4-scout:free",
	"openai":     "gpt-4o-mini",
	"gemini":     "gemini-2.5-flash",
	"ollama":     "llama3.2",
}
