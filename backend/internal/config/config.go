package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	LLM LLMConfig
	Port string
}

type LLMConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		LLM: LLMConfig{
			APIKey:  os.Getenv("LLM_API_KEY"),
			BaseURL: os.Getenv("LLM_BASE_URL"),
			Model:   os.Getenv("LLM_MODEL"),
		},
		Port: getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
