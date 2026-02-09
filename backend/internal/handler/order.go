package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/miftah/fast-order/internal/llm"
)

type GenerateOrderRequest struct {
	Mode         string `json:"mode"`
	ListMenu     string `json:"listMenu"`
	CurrentOrders string `json:"currentOrders"`
}

type GenerateOrderResponse struct {
	GeneratedMessage string `json:"generatedMessage"`
	Error           string `json:"error,omitempty"`
}

type OrderHandler struct {
	llm *llm.ResilientLLM
}

func NewOrderHandler(llmClient *llm.ResilientLLM) *OrderHandler {
	return &OrderHandler{llm: llmClient}
}

func (h *OrderHandler) GenerateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req GenerateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GenerateOrderResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Default mode to "normal" if not provided
	if req.Mode == "" {
		req.Mode = "normal"
	}

	// Validate based on mode
	switch req.Mode {
	case "first-touch":
		if req.ListMenu == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GenerateOrderResponse{
				Error: "List menu is required for first-touch mode",
			})
			return
		}
	case "nitro":
		if req.CurrentOrders == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GenerateOrderResponse{
				Error: "Current orders is required for nitro mode",
			})
			return
		}
	case "normal":
		if req.CurrentOrders == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GenerateOrderResponse{
				Error: "Current orders is required for normal mode",
			})
			return
		}
	default:
		// Invalid mode defaults to normal behavior
		if req.CurrentOrders == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GenerateOrderResponse{
				Error: "Current orders is required",
			})
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	prompt := llm.BuildPrompt(llm.GenerateOrderParams{
		Mode:          req.Mode,
		ListMenu:      req.ListMenu,
		CurrentOrders: req.CurrentOrders,
	})

	result, err := h.llm.GenerateFromSinglePrompt(ctx, prompt)
	if err != nil {
		log.Printf("[Handler] LLM error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(GenerateOrderResponse{
			Error: "Failed to generate order",
		})
		return
	}

	// Sanitize output to ensure format compliance (remove [], normalize separators)
	sanitizedResult := llm.SanitizeOrderOutput(result)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GenerateOrderResponse{
		GeneratedMessage: sanitizedResult,
	})
}
