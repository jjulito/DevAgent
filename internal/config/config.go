package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

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

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	v.SetEnvPrefix("DEVAGENT")
	v.AutomaticEnv()

	v.SetDefault("provider", "openrouter")
	v.SetDefault("model", "")
	v.SetDefault("verbose", false)
	v.SetDefault("ollama_host", "http://localhost:11434")
	v.SetDefault("qdrant_host", "localhost")
	v.SetDefault("qdrant_port", 6334)
	v.SetDefault("mcp_port", 3000)

	bindEnvAliases(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	cfg.Provider = strings.ToLower(cfg.Provider)

	if cfg.Model == "" {
		if m, ok := defaultModels[cfg.Provider]; ok {
			cfg.Model = m
		}
	}

	return cfg, nil
}

func bindEnvAliases(v *viper.Viper) {
	v.BindEnv("openrouter_api_key", "OPENROUTER_API_KEY")
	v.BindEnv("openai_api_key", "OPENAI_API_KEY")
	v.BindEnv("gemini_api_key", "GEMINI_API_KEY")
	v.BindEnv("ollama_host", "OLLAMA_HOST")
	v.BindEnv("qdrant_host", "QDRANT_HOST")
	v.BindEnv("qdrant_port", "QDRANT_PORT")
	v.BindEnv("mcp_port", "MCP_PORT")
}

func (c *Config) APIKeyForProvider() string {
	switch c.Provider {
	case "openrouter":
		return c.OpenRouterAPIKey
	case "openai":
		return c.OpenAIAPIKey
	case "gemini":
		return c.GeminiAPIKey
	default:
		return ""
	}
}

func (c *Config) Validate() error {
	validProviders := []string{"openrouter", "openai", "gemini", "ollama"}
	valid := false
	for _, p := range validProviders {
		if c.Provider == p {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid provider %q, options: %s", c.Provider, strings.Join(validProviders, ", "))
	}

	if c.Provider != "ollama" && c.APIKeyForProvider() == "" {
		return fmt.Errorf("API key required for provider %q (configure in .env or environment variable)", c.Provider)
	}

	return nil
}
