# Fast Order Webapp Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a stateless local webapp that generates WhatsApp lunch orders using LLM, with auto-copy to clipboard for maximum ordering speed.

**Architecture:** Go backend (handles LLM calls) + React (Vite) frontend (brutalist UI). Backend uses existing llm package from code_examples. Frontend follows UI_UX_GUIDE.

**Tech Stack:**
- Backend: Go 1.24+, langchaingo, gorilla/mux
- Frontend: React, Vite, TypeScript, Lucide React
- LLM: OpenAI-compatible API (DeepSeek/OpenAI)

---

## Task 1: Backend - Setup Go Module and Config

**Files:**
- Create: `backend/go.mod`
- Create: `backend/main.go` (minimal server setup)
- Create: `backend/internal/config/config.go`
- Create: `backend/.env.example`

**Step 1: Initialize Go module**

```bash
cd backend
go mod init github.com/miftah/fast-order
```

**Step 2: Create go.mod with dependencies**

File: `backend/go.mod`
```go
module github.com/miftah/fast-order

go 1.24.4

require (
	github.com/joho/godotenv v1.5.1
	github.com/gorilla/mux v1.8.1
	github.com/tmc/langchaingo v0.1.14
	github.com/sony/gobreaker v1.0.0
)
```

Run: `go mod tidy`

**Step 3: Create config loader**

File: `backend/internal/config/config.go`
```go
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
```

**Step 4: Create minimal main.go**

File: `backend/main.go`
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/miftah/fast-order/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
```

**Step 5: Create .env.example**

File: `backend/.env.example`
```env
LLM_API_KEY=your_api_key_here
LLM_BASE_URL=https://api.deepseek.com/v1
LLM_MODEL=deepseek-chat
PORT=8080
```

**Step 6: Test server runs**

```bash
cp .env.example .env
# Edit .env with your API key
go run main.go
```

Expected output: `Server starting on :8080`

Test: `curl http://localhost:8080/health`
Expected: `OK`

**Step 7: Commit**

```bash
git add backend/
git commit -m "feat: initialize Go backend with config and health endpoint"
```

---

## Task 2: Backend - LLM Service Integration

**Files:**
- Create: `backend/internal/llm/llm.go`
- Create: `backend/internal/llm/prompt.go`
- Create: `backend/internal/llm/llm_test.go` (basic test)

**Step 1: Copy/adapt llm.go from code_examples**

File: `backend/internal/llm/llm.go`
```go
package llm

import (
	"context"
	"fmt"
	"log"
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
```

**Step 2: Create prompt builder**

File: `backend/internal/llm/prompt.go`
```go
package llm

const systemPrompt = `You are an order formatter for an Indonesian office lunch catering WhatsApp group.

Your task: Generate a WhatsApp order message that:
1. Preserves all existing orders before the user
2. Appends the user's order at the next number
3. Maintains the exact format of previous orders

USER: miftah
ALWAYS USE: "nasi 1" (never "nasi 1/2")
LAUK COUNT: Exactly 2-3 lauk (no more, no less)
PROTEIN REQUIREMENT: At least 1 protein dish (e.g., fillet ayam, ati ampela, dendeng sapi, udang, ikan, ceker)

OUTPUT FORMAT: Match the format of existing orders exactly.`

type GenerateOrderParams struct {
	ListMenu     string // Optional - full menu text
	CurrentOrders string // Required - current order list
}

func BuildPrompt(params GenerateOrderParams) string {
	if params.ListMenu == "" {
		// Nitro mode: no menu provided
		return systemPrompt + `

CURRENT ORDERS:
` + params.CurrentOrders + `

NOTE: No menu provided. Choose miftah's order from dishes that appear in existing orders above.

Generate the next order message with miftah's order appended.`
	}

	// Normal mode: with menu
	return systemPrompt + `

AVAILABLE MENU:
` + params.ListMenu + `

CURRENT ORDERS:
` + params.CurrentOrders + `

Generate the next order message with miftah's order appended.`
}
```

**Step 3: Create simple test**

File: `backend/internal/llm/llm_test.go`
```go
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
```

**Step 4: Run tests**

```bash
cd backend
go test ./internal/llm/... -v
```

Expected: `PASS`

**Step 5: Commit**

```bash
git add backend/internal/llm/
git commit -m "feat: add LLM service with prompt builder"
```

---

## Task 3: Backend - HTTP API Handler

**Files:**
- Create: `backend/internal/handler/order.go`
- Modify: `backend/main.go` (add routes and CORS)

**Step 1: Create request/response types and handler**

File: `backend/internal/handler/order.go`
```go
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

	if req.CurrentOrders == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GenerateOrderResponse{
			Error: "Current orders is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	prompt := llm.BuildPrompt(llm.GenerateOrderParams{
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GenerateOrderResponse{
		GeneratedMessage: result,
	})
}
```

**Step 2: Update main.go with routes and CORS**

File: `backend/main.go`
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	llmClient, err := llm.NewResilientLLM(os.Background(), &cfg.LLM)
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
	handler := corsMiddleware(r)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

**Step 3: Test the API manually**

```bash
# Start server
go run main.go

# In another terminal, test the endpoint
curl -X POST http://localhost:8080/api/generate-order \
  -H "Content-Type: application/json" \
  -d '{"listMenu": "Cah buncis\nFillet ayam", "currentOrders": "1. farid: nasi 1"}'
```

Expected: JSON response with `generatedMessage`

**Step 4: Commit**

```bash
git add backend/
git commit -m "feat: add /api/generate-order endpoint with CORS"
```

---

## Task 4: Frontend - Project Initialization

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/index.html`
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/vite-env.d.ts`

**Step 1: Initialize Vite + React + TypeScript**

```bash
cd frontend
npm create vite@latest . -- --template react-ts
```

**Step 2: Create package.json**

File: `frontend/package.json`
```json
{
  "name": "fast-order-frontend",
  "private": true,
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "lucide-react": "^0.468.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.12",
    "@types/react-dom": "^18.3.1",
    "@vitejs/plugin-react": "^4.3.4",
    "typescript": "^5.7.2",
    "vite": "^6.0.7"
  }
}
```

**Step 3: Create vite.config.ts**

File: `frontend/vite.config.ts`
```ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
  },
})
```

**Step 4: Create tsconfig.json**

File: `frontend/tsconfig.json`
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

**Step 5: Create tsconfig.node.json**

File: `frontend/tsconfig.node.json`
```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}
```

**Step 6: Create index.html**

File: `frontend/index.html`
```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Fast Order</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

**Step 7: Create minimal main.tsx**

File: `frontend/src/main.tsx`
```tsx
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './styles/globals.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
```

**Step 8: Create vite-env.d.ts**

File: `frontend/src/vite-env.d.ts`
```typescript
/// <reference types="vite/client" />
```

**Step 9: Install dependencies and test**

```bash
cd frontend
npm install
npm run dev
```

Expected: Dev server running on http://localhost:5173

**Step 10: Commit**

```bash
git add frontend/
git commit -m "feat: initialize React + Vite + TypeScript frontend"
```

---

## Task 5: Frontend - Brutalist Styles and Layout

**Files:**
- Create: `frontend/src/styles/globals.css`
- Create: `frontend/src/App.tsx`

**Step 1: Create globals.css with brutalist styles**

File: `frontend/src/styles/globals.css`
```css
:root {
  /* Colors */
  --primary: #FF4D00;
  --primary-dark: #E64500;
  --neutral-950: #0A0A0A;
  --neutral-500: #6B7280;
  --neutral-200: #EAEAEA;
  --neutral-100: #F3F3F3;
  --white: #FFFFFF;
  --success: #22C55E;
  --error: #EF4444;

  /* Spacing */
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-6: 24px;
  --space-8: 32px;
  --space-12: 48px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  background: var(--neutral-100);
  color: var(--neutral-950);
  line-height: 1.5;
}

/* Page Container */
.page-container {
  background: var(--white);
  border: 2px solid var(--neutral-950);
  min-height: 100vh;
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  padding: var(--space-12) var(--space-8) var(--space-8);
  border-bottom: 2px solid var(--neutral-950);
  text-align: center;
}

.page-title {
  font-family: 'Archivo', sans-serif;
  font-weight: 900;
  font-size: 2.5rem;
  color: var(--neutral-950);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0;
}

/* Page Content */
.page-content {
  padding: var(--space-8);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* Field Container */
.field-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.field-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.field-label {
  font-family: 'Inter', sans-serif;
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--neutral-950);
}

.field-required {
  color: var(--error);
  font-size: 0.75rem;
}

/* Textarea */
.brutalist-textarea {
  width: 100%;
  min-height: 150px;
  padding: var(--space-3) var(--space-4);
  background: var(--white);
  border: 2px solid var(--neutral-950);
  border-radius: 0;
  font-family: 'Inter', sans-serif;
  font-size: 1rem;
  color: var(--neutral-950);
  resize: vertical;
}

.brutalist-textarea:focus {
  outline: 2px solid var(--primary);
  outline-offset: -2px;
}

.brutalist-textarea::placeholder {
  color: var(--neutral-500);
}

/* Clear Button */
.btn-clear {
  padding: var(--space-2) var(--space-3);
  background: var(--neutral-200);
  color: var(--neutral-950);
  border: 2px solid var(--neutral-950);
  border-radius: 0;
  font-family: 'Inter', sans-serif;
  font-weight: 500;
  font-size: 0.75rem;
  cursor: pointer;
  transition: background 200ms;
}

.btn-clear:hover {
  background: var(--neutral-100);
}

/* Generate Button */
.btn-generate {
  width: 100%;
  padding: var(--space-4) var(--space-6);
  background: var(--primary);
  color: var(--white);
  border: 2px solid var(--neutral-950);
  border-radius: 0;
  font-family: 'Inter', sans-serif;
  font-weight: 700;
  font-size: 1rem;
  cursor: pointer;
  transition: background 200ms;
}

.btn-generate:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-generate:disabled {
  background: var(--neutral-200);
  cursor: not-allowed;
}

/* Status Message */
.status-message {
  padding: var(--space-3) var(--space-4);
  border: 2px solid var(--neutral-950);
  font-family: 'Inter', sans-serif;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.status-message.success {
  background: #F0FDF4;
  border-color: var(--success);
}

.status-message.error {
  background: #FEF2F2;
  border-color: var(--error);
}

.status-message.loading {
  background: var(--neutral-100);
}

/* Spinner */
.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--neutral-200);
  border-top: 2px solid var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
```

**Step 2: Create App.tsx with main layout**

File: `frontend/src/App.tsx`
```tsx
import { useState } from 'react'
import './styles/globals.css'

function App() {
  const [listMenu, setListMenu] = useState('')
  const [currentOrders, setCurrentOrders] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState<{ type: 'idle' | 'success' | 'error', message: string }>({
    type: 'idle',
    message: '',
  })

  const handleGenerate = async () => {
    if (!currentOrders.trim()) {
      setStatus({ type: 'error', message: 'Current orders is required' })
      return
    }

    setIsLoading(true)
    setStatus({ type: 'idle', message: '' })

    try {
      const response = await fetch('http://localhost:8080/api/generate-order', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ listMenu, currentOrders }),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Failed to generate order')
      }

      // Copy to clipboard
      await navigator.clipboard.writeText(data.generatedMessage)

      setStatus({ type: 'success', message: 'Order copied to clipboard!' })
    } catch (error) {
      setStatus({
        type: 'error',
        message: error instanceof Error ? error.message : 'An error occurred',
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleGenerate()
    }
  }

  return (
    <div className="page-container">
      <header className="page-header">
        <h1 className="page-title">FAST ORDER</h1>
      </header>

      <main className="page-content">
        {/* List Menu - Optional */}
        <div className="field-container">
          <div className="field-header">
            <span className="field-label">LIST MENU</span>
            {listMenu && (
              <button className="btn-clear" onClick={() => setListMenu('')}>
                CLEAR
              </button>
            )}
          </div>
          <textarea
            className="brutalist-textarea"
            placeholder="Paste menu here... (Ctrl+V to paste, then TAB)"
            value={listMenu}
            onChange={(e) => setListMenu(e.target.value)}
          />
        </div>

        {/* Current Orders - Required */}
        <div className="field-container">
          <div className="field-header">
            <span className="field-label">
              CURRENT ORDERS <span className="field-required">*</span>
            </span>
            {currentOrders && (
              <button className="btn-clear" onClick={() => setCurrentOrders('')}>
                CLEAR
              </button>
            )}
          </div>
          <textarea
            className="brutalist-textarea"
            placeholder="Paste current orders here... (Ctrl+V to paste, then ENTER to generate)"
            value={currentOrders}
            onChange={(e) => setCurrentOrders(e.target.value)}
            onKeyDown={handleKeyDown}
          />
        </div>

        {/* Generate Button */}
        <button
          className="btn-generate"
          onClick={handleGenerate}
          disabled={!currentOrders.trim() || isLoading}
        >
          {isLoading ? 'GENERATING...' : 'GENERATE & COPY'}
        </button>

        {/* Status Message */}
        {status.type !== 'idle' && (
          <div className={`status-message ${status.type}`}>
            {status.type === 'loading' && <div className="spinner" />}
            {status.message}
          </div>
        )}
      </main>
    </div>
  )
}

export default App
```

**Step 3: Test the UI**

```bash
npm run dev
```

Expected: Brutalist UI with two textareas and generate button

**Step 4: Commit**

```bash
git add frontend/src/
git commit -m "feat: add brutalist UI and main App layout"
```

---

## Task 6: Frontend - Components and API Integration

**Files:**
- Create: `frontend/src/components/TextAreaField.tsx`
- Create: `frontend/src/components/GenerateButton.tsx`
- Create: `frontend/src/components/StatusMessage.tsx`
- Create: `frontend/src/hooks/useClipboard.ts`
- Create: `frontend/src/lib/api.ts`
- Modify: `frontend/src/App.tsx` (refactor to use components)

**Step 1: Create TextAreaField component**

File: `frontend/src/components/TextAreaField.tsx`
```tsx
import { ReactNode } from 'react'

interface TextAreaFieldProps {
  label: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  onKeyDown?: (e: React.KeyboardEvent) => void
  hint?: ReactNode
}

export function TextAreaField({
  label,
  value,
  onChange,
  placeholder,
  required = false,
  onKeyDown,
  hint,
}: TextAreaFieldProps) {
  return (
    <div className="field-container">
      <div className="field-header">
        <span className="field-label">
          {label}
          {required && <span className="field-required"> *</span>}
        </span>
        {value && (
          <button
            className="btn-clear"
            onClick={() => onChange('')}
          >
            CLEAR
          </button>
        )}
      </div>
      <textarea
        className="brutalist-textarea"
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={onKeyDown}
      />
      {hint && <div className="field-hint">{hint}</div>}
    </div>
  )
}
```

**Step 2: Create GenerateButton component**

File: `frontend/src/components/GenerateButton.tsx`
```tsx
interface GenerateButtonProps {
  onClick: () => void
  disabled?: boolean
  loading?: boolean
}

export function GenerateButton({ onClick, disabled, loading }: GenerateButtonProps) {
  return (
    <button
      className="btn-generate"
      onClick={onClick}
      disabled={disabled}
    >
      {loading ? 'GENERATING...' : 'GENERATE & COPY'}
    </button>
  )
}
```

**Step 3: Create StatusMessage component**

File: `frontend/src/components/StatusMessage.tsx`
```tsx
interface StatusMessageProps {
  type: 'idle' | 'success' | 'error' | 'loading'
  message: string
}

export function StatusMessage({ type, message }: StatusMessageProps) {
  if (type === 'idle') return null

  return (
    <div className={`status-message ${type}`}>
      {type === 'loading' && <div className="spinner" />}
      {message}
    </div>
  )
}
```

**Step 4: Create useClipboard hook**

File: `frontend/src/hooks/useClipboard.ts`
```tsx
import { useState } from 'react'

export function useClipboard() {
  const [isCopied, setIsCopied] = useState(false)

  const copy = async (text: string): Promise<boolean> => {
    try {
      await navigator.clipboard.writeText(text)
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 2000)
      return true
    } catch {
      setIsCopied(false)
      return false
    }
  }

  return { copy, isCopied }
}
```

**Step 5: Create API client**

File: `frontend/src/lib/api.ts`
```typescript
const API_BASE = 'http://localhost:8080'

export interface GenerateOrderParams {
  listMenu: string
  currentOrders: string
}

export interface GenerateOrderResponse {
  generatedMessage: string
  error?: string
}

export async function generateOrder(
  params: GenerateOrderParams
): Promise<GenerateOrderResponse> {
  const response = await fetch(`${API_BASE}/api/generate-order`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params),
  })

  const data: GenerateOrderResponse = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'Failed to generate order')
  }

  return data
}
```

**Step 6: Refactor App.tsx to use components**

File: `frontend/src/App.tsx` (replace entire file)
```tsx
import { useState } from 'react'
import { TextAreaField } from './components/TextAreaField'
import { GenerateButton } from './components/GenerateButton'
import { StatusMessage } from './components/StatusMessage'
import { generateOrder } from './lib/api'
import './styles/globals.css'

function App() {
  const [listMenu, setListMenu] = useState('')
  const [currentOrders, setCurrentOrders] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState<{ type: 'idle' | 'success' | 'error', message: string }>({
    type: 'idle',
    message: '',
  })

  const handleGenerate = async () => {
    if (!currentOrders.trim()) {
      setStatus({ type: 'error', message: 'Current orders is required' })
      return
    }

    setIsLoading(true)
    setStatus({ type: 'idle', message: '' })

    try {
      const data = await generateOrder({ listMenu, currentOrders })

      // Copy to clipboard
      await navigator.clipboard.writeText(data.generatedMessage)

      setStatus({ type: 'success', message: 'Order copied to clipboard!' })
    } catch (error) {
      setStatus({
        type: 'error',
        message: error instanceof Error ? error.message : 'An error occurred',
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleGenerate()
    }
  }

  return (
    <div className="page-container">
      <header className="page-header">
        <h1 className="page-title">FAST ORDER</h1>
      </header>

      <main className="page-content">
        <TextAreaField
          label="LIST MENU"
          value={listMenu}
          onChange={setListMenu}
          placeholder="Paste menu here... (Ctrl+V to paste, then TAB)"
          hint="Optional - leave empty for Nitro Mode"
        />

        <TextAreaField
          label="CURRENT ORDERS"
          value={currentOrders}
          onChange={setCurrentOrders}
          placeholder="Paste current orders here... (Ctrl+V to paste, then ENTER to generate)"
          required
          onKeyDown={handleKeyDown}
        />

        <GenerateButton
          onClick={handleGenerate}
          disabled={!currentOrders.trim() || isLoading}
          loading={isLoading}
        />

        <StatusMessage type={status.type} message={status.message} />
      </main>
    </div>
  )
}

export default App
```

**Step 7: Test integration**

Start backend: `cd backend && go run main.go`
Start frontend: `cd frontend && npm run dev`

Test the full flow.

**Step 8: Commit**

```bash
git add frontend/src/
git commit -m "feat: refactor into components and add API client"
```

---

## Task 7: Frontend - Keyboard Shortcuts and Polish

**Files:**
- Modify: `frontend/src/App.tsx` (add keyboard shortcuts)
- Modify: `frontend/src/styles/globals.css` (add hint styles)

**Step 1: Add hint styles to globals.css**

Add to `frontend/src/styles/globals.css`:
```css
/* Field Hint */
.field-hint {
  font-family: 'Inter', sans-serif;
  font-size: 0.75rem;
  color: var(--neutral-500);
  margin-top: var(--space-1);
}
```

**Step 2: Add keyboard shortcuts and polish to App.tsx**

File: `frontend/src/App.tsx`
```tsx
import { useState, useEffect } from 'react'
import { TextAreaField } from './components/TextAreaField'
import { GenerateButton } from './components/GenerateButton'
import { StatusMessage } from './components/StatusMessage'
import { generateOrder } from './lib/api'
import './styles/globals.css'

function App() {
  const [listMenu, setListMenu] = useState('')
  const [currentOrders, setCurrentOrders] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState<{ type: 'idle' | 'success' | 'error', message: string }>({
    type: 'idle',
    message: '',
  })

  // Global keyboard shortcuts
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      // ESC to clear all
      if (e.key === 'Escape') {
        setListMenu('')
        setCurrentOrders('')
        setStatus({ type: 'idle', message: '' })
      }

      // Ctrl+Shift+C to generate (when not typing)
      if (e.ctrlKey && e.shiftKey && e.key === 'C') {
        e.preventDefault()
        handleGenerate()
      }
    }

    window.addEventListener('keydown', handleGlobalKeyDown)
    return () => window.removeEventListener('keydown', handleGlobalKeyDown)
  }, [listMenu, currentOrders, isLoading])

  const handleGenerate = async () => {
    if (!currentOrders.trim()) {
      setStatus({ type: 'error', message: 'Current orders is required' })
      return
    }

    if (isLoading) return

    setIsLoading(true)
    setStatus({ type: 'idle', message: '' })

    try {
      const data = await generateOrder({ listMenu, currentOrders })

      await navigator.clipboard.writeText(data.generatedMessage)

      setStatus({ type: 'success', message: 'Order copied to clipboard! Press Ctrl+V to paste in WhatsApp' })
    } catch (error) {
      setStatus({
        type: 'error',
        message: error instanceof Error ? error.message : 'An error occurred',
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleGenerate()
    }
  }

  return (
    <div className="page-container">
      <header className="page-header">
        <h1 className="page-title">FAST ORDER</h1>
      </header>

      <main className="page-content">
        <TextAreaField
          label="LIST MENU"
          value={listMenu}
          onChange={setListMenu}
          placeholder="Paste menu here... (Ctrl+V to paste, then TAB)"
          hint="Optional - leave empty for Nitro Mode"
        />

        <TextAreaField
          label="CURRENT ORDERS"
          value={currentOrders}
          onChange={setCurrentOrders}
          placeholder="Paste current orders here... (Ctrl+V to paste, then ENTER to generate)"
          required
          onKeyDown={handleKeyDown}
        />

        <GenerateButton
          onClick={handleGenerate}
          disabled={!currentOrders.trim() || isLoading}
          loading={isLoading}
        />

        <StatusMessage type={status.type} message={status.message} />

        {status.type === 'idle' && (
          <div className="field-hint" style={{ textAlign: 'center', marginTop: 'var(--space-2)' }}>
            Shortcuts: ENTER to generate • ESC to clear • Ctrl+Shift+C to generate
          </div>
        )}
      </main>
    </div>
  )
}

export default App
```

**Step 3: Test all shortcuts**

- ENTER in Current Orders → generates and copies
- ESC → clears all fields
- Ctrl+Shift+C → generates and copies

**Step 4: Final commit**

```bash
git add frontend/src/
git commit -m "feat: add keyboard shortcuts and polish"
```

---

## Final Verification

**Run the full app:**

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```

**Test the complete flow:**
1. Paste menu into List Menu field
2. Paste current orders into Current Orders field
3. Press ENTER
4. Verify message is copied to clipboard
5. Paste into a text editor to verify format

**Test Nitro Mode:**
1. Leave List Menu empty
2. Paste current orders (with existing orders)
3. Press ENTER
4. Verify LLM chooses from existing orders

---

## Task Dependencies Summary

```
Task 1 (Backend Setup) ──────┐
                              │
Task 2 (LLM Service) ─────────┼── Task 3 (API Handler)
                              │
Task 4 (Frontend Init) ───────┼── Task 5 (Styles & Layout) ── Task 6 (Components) ── Task 7 (Shortcuts)
```

---

## Notes

- Backend runs on port 8080
- Frontend dev server runs on port 5173
- Production: Build frontend (`npm run build`), serve via Go backend
- All commits should follow conventional commits format
- State is NOT persisted - refresh loses data (by design)
