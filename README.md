# Fast Order ⚡

Stateless webapp for generating WhatsApp lunch orders using LLM. Paste menu + current orders, get your order formatted and copied to clipboard instantly.

## Features

- 🚀 **Lightning Fast** - Paste, generate, copy in seconds
- 🤖 **LLM Powered** - Smart order generation with context awareness
- 📋 **Auto-Copy** - Result copied to clipboard automatically
- ⌨️ **Keyboard Shortcuts** - ENTER to generate, ESC to clear
- 🔥 **Nitro Mode** - Works without menu (chooses from existing orders)

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
| `ENTER` | Generate & copy order |
| `ESC` | Clear all fields |
| `Ctrl+Shift+C` | Generate & copy order |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_API_KEY` | Your LLM provider API key | - |
| `LLM_BASE_URL` | LLM API base URL | - |
| `LLM_MODEL` | Model name | - |
| `PORT` | Backend port | 8080 |

## Tech Stack

- **Backend**: Go 1.24, langchaingo, gorilla/mux
- **Frontend**: React, Vite, TypeScript
- **LLM**: OpenAI-compatible (DeepSeek, OpenAI, etc.)

## License

MIT
