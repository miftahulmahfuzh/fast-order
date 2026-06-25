package config

import (
	"os"

	"github.com/joho/godotenv"
)

// geminiOpenAIBaseURL is Google's OpenAI-compatible endpoint for Gemini models.
const geminiOpenAIBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"

type Config struct {
	LLM  LLMConfig
	Port string
}

type LLMConfig struct {
	Type    string
	APIKey  string
	BaseURL string
	Model   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		LLM:  resolveLLMConfig(),
		Port: getEnv("PORT", "8080"),
	}, nil
}

// resolveLLMConfig selects the active LLM provider based on LLM_TYPE.
// GEMINI uses Google's OpenAI-compatible endpoint so the same client path
// (and circuit breaker) is reused. Any other value falls back to DEEPSEEK,
// which is the default and reads the generic LLM_* variables.
func resolveLLMConfig() LLMConfig {
	llmType := getEnv("LLM_TYPE", "DEEPSEEK")

	if llmType == "GEMINI" {
		return LLMConfig{
			Type:    llmType,
			APIKey:  os.Getenv("GEMINI_LLM_API_KEY"),
			BaseURL: geminiOpenAIBaseURL,
			Model:   os.Getenv("GEMINI_MODEL"),
		}
	}

	return LLMConfig{
		Type:    "DEEPSEEK",
		APIKey:  os.Getenv("LLM_API_KEY"),
		BaseURL: os.Getenv("LLM_BASE_URL"),
		Model:   os.Getenv("LLM_MODEL"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
