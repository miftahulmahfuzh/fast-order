package llm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/miftah/fast-order/internal/config"
	"github.com/sony/gobreaker"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type ResilientLLM struct {
	llm     llms.Model
	breaker *gobreaker.CircuitBreaker
}

func NewResilientLLM(ctx context.Context, cfg *config.LLMConfig) (*ResilientLLM, error) {
	llmClient, err := openai.New(
		openai.WithToken(cfg.APIKey),
		openai.WithBaseURL(cfg.BaseURL),
		openai.WithModel(cfg.Model),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "llm",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 2
		},
	})

	return &ResilientLLM{
		llm:     llmClient,
		breaker: breaker,
	}, nil
}

// SanitizeOrderOutput cleans LLM output to ensure format compliance.
// Removes square brackets [] and normalizes separators to use ":"
func SanitizeOrderOutput(input string) string {
	// Remove square brackets
	result := input
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")

	// Normalize separator: replace " - " with " : "
	// This handles cases where LLM uses dash instead of colon
	result = strings.ReplaceAll(result, " - ", " : ")

	return result
}

func (r *ResilientLLM) GenerateFromSinglePrompt(ctx context.Context, prompt string) (string, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		return llms.GenerateFromSinglePrompt(ctx, r.llm, prompt)
	})

	if err != nil {
		log.Printf("[LLM] ERROR: %v", err)
		return "", fmt.Errorf("LLM API call failed: %w", err)
	}

	return result.(string), nil
}
