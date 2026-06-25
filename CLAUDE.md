# Fast Order - Claude Context

## Project Overview
Stateless webapp for generating WhatsApp lunch orders using LLM. Auto-copies to clipboard for maximum speed.

## Tech Stack
- **Backend**: Go 1.24, langchaingo, gorilla/mux
- **Frontend**: React, Vite, TypeScript, Lucide React
- **LLM**: OpenAI-compatible API (DeepSeek/OpenAI)

## Architecture
```
┌─────────────┐         ┌─────────────┐
│   React     │ ──API──▶│    Go       │ ──▶ LLM API
│  (Vite)     │         │  Backend    │
└─────────────┘         └─────────────┘
     Port 5173             Port 8080
```

## Key Files
- `backend/main.go` - Server entry with CORS
- `backend/internal/llm/llm.go` - Resilient LLM with circuit breaker
- `backend/internal/llm/prompt.go` - Order prompt templates
- `backend/internal/handler/order.go` - API handlers
- `frontend/src/App.tsx` - Main UI with keyboard shortcuts
- `frontend/src/lib/api.ts` - API client

## Development
```bash
# Backend
cd backend && go run main.go

# Frontend
cd frontend && npm run dev
```

## Environment Variables
```bash
LLM_API_KEY=your_key
LLM_BASE_URL=https://api.deepseek.com/v1
LLM_MODEL=deepseek-chat
PORT=8080
```

## Features
- **Nitro Mode**: Works without menu (chooses from existing orders)
- **Keyboard Shortcuts**: ENTER (generate), ESC (clear), Ctrl+Shift+C (generate)
- **Auto-copy**: Generated order copied to clipboard automatically

> [!IMPORTANT]

> when user type in "l" : you must run "golangci-lint run --output.text.path=lint_output.txt 2>&1; cat lint_output.txt" . then fix 10 issues in NEW parallel subagents. USE opus model for this task.

> when user type in "p" : you must git add . commit push in a new subagent. USE haiku model for this task.

> when user type in "s" : you must summarize what you just did in the current session into a concise task item in the most appropriate todos.md file in this codebase. put it in Completed Tasks section . do it in a new subagent. USE haiku model for this task.

> when user type in "j": you must summarize everything you just did in the current session into a one liner task description that user can put in daily task report in Jira. do it in a new subagent. USE haiku model for this task.

> when user type in "q": you must locate the most relevant package_readme.md (representing a single package) that got influenced the most by the current session, and update that package_readme.md to reflect the current state of that package. USE sonnet model for this task.
