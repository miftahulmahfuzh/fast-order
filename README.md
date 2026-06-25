# Fast Order ⚡

Stateless webapp for generating WhatsApp lunch orders using LLM. Paste menu + current orders, get your order formatted and copied to clipboard instantly.

## Features

- 🚀 **Lightning Fast** - Paste, generate, copy in seconds
- 🤖 **LLM Powered** - Smart order generation with context awareness
- 📋 **Auto-Copy** - Result copied to clipboard automatically
- ⌨️ **Keyboard Shortcuts** - ENTER to generate, ESC to clear
- 🔥 **Three Modes** - Normal, Nitro, and First-Touch mode for any ordering scenario

## Order Modes

| Mode | When to Use | Input | Shortcut |
|------|-------------|-------|----------|
| **First-Touch** | You're the first to order | Menu only | `Shift+Enter` in List Menu |
| **Normal** | Others have already ordered | Menu + Current Orders | `ENTER` in Current Orders |
| **Nitro** | No menu available, choose from existing | Current Orders only | Leave List Menu empty + `ENTER` |

### First-Touch Mode 🆕
Generate order #1 when you're the first person ordering. Just paste the menu and press `Shift+Enter`.

### Normal Mode
Append your order to an existing list. Paste both the menu and current orders, then press `ENTER`.

### Nitro Mode
Generate an order based on dishes others have already ordered, without needing the menu. Leave the List Menu empty and press `ENTER`.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# 1. Copy env file
cp backend/.env.example .env

# 2. Edit .env with your LLM API credentials
# LLM_API_KEY=your_api_key_here
# LLM_BASE_URL=https://api.deepseek.com/v1
# LLM_MODEL=deepseek-chat

# 3. Start services
docker compose up -d

# 4. Open http://localhost:5173
```

### Manual Setup

**Backend:**
```bash
cd backend
cp .env.example .env
# Edit .env with your API credentials
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# Open http://localhost:5173
```

## Deployment

### Option 1: Docker Compose (Production)

```bash
# Build and start
docker compose up -d --build

# View logs
docker compose logs -f

# Stop
docker compose down
```

### Option 2: Single Server Deployment

Build frontend and serve via backend:

```bash
# 1. Build frontend
cd frontend
npm run build

# 2. Copy dist to backend
cp -r dist ../backend/static

# 3. Update backend main.go to serve static files
# Add: r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

# 4. Run backend
cd ../backend
go run main.go
```

### Option 3: Cloud Deployment

**Render/Railway/etc:**

1. **Backend**: Deploy as Go service
   - Set env vars: `LLM_API_KEY`, `LLM_BASE_URL`, `LLM_MODEL`
   - Port: 8080

2. **Frontend**: Deploy as static site
   - Build: `npm run build`
   - Update `API_BASE` in `frontend/src/lib/api.ts` to backend URL
   - Deploy `dist/` folder

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Shift+Enter` (in List Menu) | First-Touch Mode - Generate order #1 |
| `Enter` (in Current Orders) | Normal/Nitro Mode - Generate & copy |
| `Ctrl+Shift+C` | Generate & copy (respects current mode) |
| `ESC` | Clear all fields |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_TYPE` | Active provider: `DEEPSEEK` (default) or `GEMINI` | `DEEPSEEK` |
| `LLM_API_KEY` | DeepSeek/OpenAI-compatible API key | - |
| `LLM_BASE_URL` | DeepSeek/OpenAI-compatible base URL | - |
| `LLM_MODEL` | DeepSeek/OpenAI model name | - |
| `GEMINI_LLM_API_KEY` | Google AI API key (used when `LLM_TYPE=GEMINI`) | - |
| `GEMINI_MODEL` | Gemini model name, e.g. `gemini-3.1-flash-lite` | - |
| `PORT` | Backend port | 8080 |

### Choosing a provider

Set `LLM_TYPE=GEMINI` to route generation through Google's Gemini (Flash Lite is
significantly faster). Gemini is called via Google's OpenAI-compatible endpoint, so
it reuses the same client path and circuit breaker. Set `LLM_TYPE=DEEPSEEK` (or omit
it) to use the DeepSeek/OpenAI-compatible `LLM_*` settings.

> **Performance note:** The backend parses the order numbering itself and only asks
> the LLM for the dishes (a single short line), then assembles the numbered result in
> Go. This keeps the existing list verbatim and makes Normal and Nitro modes
> noticeably faster. The API response includes a `durationMs` field reporting LLM time.

## Tech Stack

- **Backend**: Go 1.24, langchaingo, gorilla/mux
- **Frontend**: React, Vite, TypeScript
- **LLM**: OpenAI-compatible (DeepSeek, OpenAI, etc.)

## Testing

The project includes two integration test scripts for different scenarios:

### `./test-integration.sh`
Runs **all** tests against the local development environment.
```bash
./test-integration.sh
```
- Use during active development
- No Docker required
- Faster feedback loop
- Runs the full test suite via `npm run test`

### `./test-docker-integration.sh`
Runs Docker-specific integration tests against the deployed stack.
```bash
./test-docker-integration.sh
```
- Use before deploying or after infrastructure changes
- Automatically ensures Docker containers are running
- Tests against the actual production-like environment
- Validates full deployment health

### Why Two Scripts?

| Script | Environment | Scope | Use Case |
|--------|-------------|-------|----------|
| `test-integration.sh` | Local dev | All tests | Feature development, quick iteration |
| `test-docker-integration.sh` | Docker deployed | Integration only | Pre-deploy validation, CI/CD |

The separation keeps local development fast while ensuring deployment-specific concerns (container health, networking, etc.) are validated separately.

### Backend Unit Tests

```bash
cd backend
go test ./...
```

### End-to-End Test (against a running backend)

The Go e2e test hits the live `/api/generate-order` endpoint in **nitro mode** and
reports how long generation took (both server-reported `durationMs` and full
client round-trip). It is gated behind the `e2e` build tag so it never runs during
normal `go test ./...`.

```bash
# 1. Start the stack (reads provider config from .env, e.g. LLM_TYPE=GEMINI)
docker compose down && docker compose up -d --build

# 2. Run the e2e test against the container (backend on :8089)
cd backend
BASE_URL=http://localhost:8089 go test -tags e2e -v -count=1 ./e2e/...
```

`-count=1` disables Go's test cache so you get a fresh timing measurement each run.
`BASE_URL` defaults to `http://localhost:8089` if unset.

## License

MIT
