package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miftah/fast-order/internal/config"
	"github.com/miftah/fast-order/internal/handler"
	"github.com/miftah/fast-order/internal/llm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize LLM
	llmClient, err := llm.NewResilientLLM(context.Background(), &cfg.LLM)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers
	orderHandler := handler.NewOrderHandler(llmClient)

	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/generate-order", orderHandler.GenerateOrder).Methods("POST")

	// CORS middleware
	corsHandler := corsMiddleware(r)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins for production (Render frontend -> Render backend)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
