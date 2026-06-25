package llm

import (
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name     string
		params   GenerateOrderParams
		contains string
	}{
		{
			name: "Normal mode with menu references the menu",
			params: GenerateOrderParams{
				Mode:          "normal",
				ListMenu:      "Cah buncis\nFillet ayam",
				CurrentOrders: "1. farid: nasi 1",
			},
			contains: "Cah buncis",
		},
		{
			name: "First-touch mode references the menu",
			params: GenerateOrderParams{
				Mode:     "first-touch",
				ListMenu: "Cah buncis\nFillet ayam",
			},
			contains: "Choose from this menu",
		},
		{
			name: "Nitro mode references existing orders",
			params: GenerateOrderParams{
				Mode:          "nitro",
				CurrentOrders: "1. farid: nasi 1, jamur crispy",
			},
			contains: "jamur crispy",
		},
		{
			name: "All modes instruct items-only output",
			params: GenerateOrderParams{
				Mode:          "nitro",
				CurrentOrders: "1. farid: nasi 1",
			},
			contains: "comma-separated line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := BuildPrompt(tt.params)
			if !strings.Contains(prompt, tt.contains) {
				t.Errorf("Prompt should contain %q\ngot:\n%s", tt.contains, prompt)
			}
		})
	}
}
