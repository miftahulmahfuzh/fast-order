package config

import "testing"

func TestResolveLLMConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want LLMConfig
	}{
		{
			name: "default falls back to DeepSeek",
			env: map[string]string{
				"LLM_API_KEY":  "ds-key",
				"LLM_BASE_URL": "https://deepseek.example/v1",
				"LLM_MODEL":    "deepseek-v3",
			},
			want: LLMConfig{
				Type:    "DEEPSEEK",
				APIKey:  "ds-key",
				BaseURL: "https://deepseek.example/v1",
				Model:   "deepseek-v3",
			},
		},
		{
			name: "explicit DeepSeek",
			env: map[string]string{
				"LLM_TYPE":     "DEEPSEEK",
				"LLM_API_KEY":  "ds-key",
				"LLM_BASE_URL": "https://deepseek.example/v1",
				"LLM_MODEL":    "deepseek-v3",
			},
			want: LLMConfig{
				Type:    "DEEPSEEK",
				APIKey:  "ds-key",
				BaseURL: "https://deepseek.example/v1",
				Model:   "deepseek-v3",
			},
		},
		{
			name: "Gemini uses OpenAI-compatible endpoint and Gemini vars",
			env: map[string]string{
				"LLM_TYPE":           "GEMINI",
				"GEMINI_LLM_API_KEY": "gm-key",
				"GEMINI_MODEL":       "gemini-3.1-flash-lite",
				"LLM_API_KEY":        "ds-key-should-be-ignored",
			},
			want: LLMConfig{
				Type:    "GEMINI",
				APIKey:  "gm-key",
				BaseURL: geminiOpenAIBaseURL,
				Model:   "gemini-3.1-flash-lite",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearLLMEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := resolveLLMConfig()
			if got != tt.want {
				t.Errorf("resolveLLMConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// clearLLMEnv resets all LLM-related env vars so each case starts clean.
func clearLLMEnv(t *testing.T) {
	for _, k := range []string{
		"LLM_TYPE", "LLM_API_KEY", "LLM_BASE_URL", "LLM_MODEL",
		"GEMINI_LLM_API_KEY", "GEMINI_MODEL",
	} {
		t.Setenv(k, "")
	}
}
