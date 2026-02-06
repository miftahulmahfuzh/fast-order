package llm

import (
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name     string
		params   GenerateOrderParams
		contains string
	}{
		{
			name: "Normal mode with menu",
			params: GenerateOrderParams{
				ListMenu:     "Cah buncis\nFillet ayam",
				CurrentOrders: "1. farid: nasi 1",
			},
			contains: "AVAILABLE MENU",
		},
		{
			name: "Nitro mode without menu",
			params: GenerateOrderParams{
				ListMenu:     "",
				CurrentOrders: "1. farid: nasi 1",
			},
			contains: "No menu provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := BuildPrompt(tt.params)
			if !contains(prompt, tt.contains) {
				t.Errorf("Prompt should contain %q", tt.contains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
