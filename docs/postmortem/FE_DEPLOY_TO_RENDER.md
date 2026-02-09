# Postmortem: Frontend Deployment to Render

**Date:** 2026-02-09

**Service:** Fast Order Frontend (React + Vite + Nginx)

**Platform:** Render.com

## Executive Summary

Deploying the frontend to Render required two phases of fixes:
1. **Phase 1 (6 commits):** Fixing nginx environment variable substitution issues
2. **Phase 2 (1 commit):** Fixing nginx proxy configuration for backend connectivity (307 redirect loop)

The final solution required: proper envsubst escaping, using HTTPS for backend connection, removing incorrect Host header override, and enabling SSL SNI.

## Timeline

| Commit | Time | Change |
|--------|------|--------|
| `a0e720a` | Initial | Add Render deployment configuration |
| `f02938f` | Fix #1 | Use `$PORT` env var for Render compatibility |
| `1928d3c` | Fix #2 | Use shell form CMD for proper envsubst |
| `2968535` | Fix #3 | Use `__PLACEHOLDER__` syntax to avoid nginx variable conflict |
| `079a068` | Fix #4 | Simplify envsubst usage |
| `7912fab` | Fix #5 | Use envsubst SHELL-FORMAT to preserve nginx variables |
| `5e20ed7` | Fix #6 | Properly escape envsubst variables in CMD |
| `30bca85` | Fix #7 | Fix nginx proxy for backend connectivity (307 redirect loop) |

## The Problem

### Root Cause

The nginx Docker image uses `envsubst` to substitute environment variables in configuration templates. However, when called without specifying which variables to substitute, `envsubst` replaces **ALL** `$variable` patterns, including nginx's built-in runtime variables:

```nginx
# Template file
proxy_set_header Host $host;           # $host is an nginx variable
proxy_set_header X-Real-IP $remote_addr;  # $remote_addr is an nginx variable
```

When `envsubst` runs without constraints:
```
$host         → "" (empty, because no env var named "host")
$remote_addr  → "" (empty, because no env var named "remote_addr")
```

Resulting in broken nginx configuration:
```nginx
proxy_set_header Host ;                  # Syntax error!
proxy_set_header X-Real-IP ;             # Syntax error!
```

### Error Messages Encountered

1. **Invalid number of arguments:**
   ```
   [emerg] invalid number of arguments in "proxy_set_header" directive
   ```

2. **Host not found:**
   ```
   [emerg] host not found in "$PORT" of the "listen" directive
   ```

## The Solution Journey

### Attempt 1: Initial Render Setup (`a0e720a`)
```dockerfile
CMD ["/bin/sh", "-c", "envsubst '$$BACKEND_URL' < ..."]
```
**Problem:** `$$BACKEND_URL` expands to process ID + "BACKEND_URL", not the env var.

### Attempt 2: Add PORT variable (`f02938f`)
Changed `listen 80` to `listen $PORT` for Render compatibility.

### Attempt 3: Shell form CMD (`1928d8c`)
Switched from JSON array to string form of CMD.

### Attempt 4: PLACEHOLDER syntax (`2968535`)
Tried using `__PLACEHOLDER__` as a delimiter that wouldn't conflict.

### Attempt 5: Simplify envsubst (`079a068`)
Removed placeholder approach, went back to simpler approach.

### Attempt 6: SHELL-FORMAT (`7912fab`)
```dockerfile
CMD /bin/sh -c '... envsubst "$$PORT $$BACKEND_URL" < ...'
```
**Problem:** `$$PORT` still expands to PID + "PORT", not `$PORT`.

### Attempt 7: Proper escaping (`5e20ed7`) ✅
```dockerfile
CMD ["/bin/sh", "-c", "... envsubst '\\$PORT \\$BACKEND_URL' < ..."]
```
**Success:** The `\\$` escaping ensures envsubst receives literal `$PORT` and `$BACKEND_URL`.

---

## Phase 2: Backend Proxy Configuration Issue

### The Problem

After successfully deploying the frontend, the app returned **HTTP 307 redirects** when trying to call the backend API, causing a redirect loop.

#### Symptom Flow
1. Browser sends: `POST /api/generate-order` to `https://fast-order-1.onrender.com`
2. Frontend nginx returns: **307 redirect**
3. Browser follows redirect to: `POST /generate-order` (without `/api`)
4. Frontend returns: **405 Method Not Allowed** (no route matches)

#### Render Frontend Logs
```
127.0.0.1 - - [09/Feb/2026:06:35:12 +0000] "POST /api/generate-order HTTP/1.1" 307 0
127.0.0.1 - - [09/Feb/2026:06:35:12 +0000] "POST /generate-order HTTP/1.1" 405 559
```

### Root Cause Analysis

The issue had **three contributing factors**:

#### 1. Host Header Mismatch (`changeOrigin` equivalent missing)
```nginx
# WRONG: This sends the Frontend's domain to the Backend
proxy_set_header Host $host;  # Host = fast-order-1.onrender.com (frontend)
```

When connecting to the backend but sending the frontend's Host header, Render's routing became confused. The load balancer saw a request for the frontend (which forces HTTPS), returning a 307 redirect.

**Vite proxy worked because** it has `changeOrigin: true`, which changes the Host header to match the target (backend).

#### 2. HTTP vs HTTPS Protocol
```bash
# WRONG: Using HTTP triggers Render's automatic HTTPS redirect
BACKEND_URL=http://fast-order-xvkq.onrender.com
```

Render applications are served over HTTPS. While port 80 (HTTP) is open, Render's load balancer redirects HTTP traffic to HTTPS, causing 307 responses.

#### 3. Missing SSL SNI Configuration
When connecting to HTTPS backends on shared infrastructure like Render, SNI (Server Name Indication) is required so the load balancer knows which app's SSL certificate to use.

### The Solution (`30bca85`) ✅

#### Step 1: Update Environment Variable
```bash
# Change from HTTP to HTTPS
BACKEND_URL=https://fast-order-xvkq.onrender.com
```

#### Step 2: Fix nginx.conf
```nginx
location /api/ {
    proxy_pass $BACKEND_URL;

    # REMOVE THIS LINE (was causing Host header mismatch)
    # proxy_set_header Host $host;

    # ADD THIS LINE (required for HTTPS on Render)
    proxy_ssl_server_name on;

    # Keep these for logging/IP tracking
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_redirect off;
}
```

### Why This Works

1. **`BACKEND_URL` with HTTPS:** Nginx talks to the backend securely, avoiding Render's HTTP→HTTPS redirect
2. **Removing `Host $host`:** Nginx now defaults to using the hostname from `proxy_pass` URL (backend's host), mimicking `changeOrigin: true` behavior
3. **`proxy_ssl_server_name on`:** Enables SNI so Render knows which SSL certificate and app to target during the TLS handshake

### Testing Checklist

- [ ] Local frontend → Local backend: Works
- [ ] Local frontend → Render backend (via Vite proxy): Works
- [ ] curl → Render backend directly: Works
- [ ] Render frontend → Render backend: **Now works** ✅

## Final Working Configuration

### Dockerfile
```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Runtime stage
FROM nginx:alpine
RUN apk add --no-cache gettext
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Use envsubst with SHELL-FORMAT to only substitute specific variables
# This prevents replacing nginx's built-in variables like $host, $remote_addr, etc.
CMD ["/bin/sh", "-c", "export PORT=${PORT:-10000} && envsubst '\\$PORT \\$BACKEND_URL' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && exec nginx -g 'daemon off;'"]
```

### nginx.conf
```nginx
server {
    listen $PORT;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy - BACKEND_URL will be substituted by envsubst
    # IMPORTANT: Use HTTPS (e.g., https://backend.onrender.com)
    # DO NOT set Host header - let nginx use backend's hostname
    location /api/ {
        proxy_pass $BACKEND_URL;

        # Enable SNI for HTTPS (required for Render)
        proxy_ssl_server_name on;

        # Don't set Host - let nginx use the hostname from BACKEND_URL
        # This mimics 'changeOrigin: true' behavior in Vite proxy

        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_redirect off;
    }
}
```

## Key Learnings

### 1. envsubst Substitutes Everything by Default
`envsubst < template > output` replaces ALL `$variable` patterns, not just environment variables.

### 2. Use SHELL-FORMAT to Constrain Substitution
```bash
# Only substitute specific variables
envsubst '$VAR1 $VAR2' < template > output
```

### 3. Escaping is Tricky in Docker CMD
| Syntax | Result |
|--------|--------|
| `"$$VAR"` | PID + "VAR" (wrong) |
| `'\$VAR'` | Literal `$VAR` (correct) |
| `"\\$VAR"` | Literal `$VAR` (correct in JSON array) |

### 4. JSON Array CMD is More Reliable
```dockerfile
# Preferred: JSON array (no shell processing)
CMD ["/bin/sh", "-c", "..."]

# Risky: String form (shell processes twice)
CMD /bin/sh -c '...'
```

### 5. Variable Distinguishability
Use `${VAR}` syntax for environment variables and `$var` for nginx variables to make templates clearer.

### 6. Nginx Proxy Headers Matter
When proxying to a different backend, **don't override the Host header** unless necessary. Let nginx use the hostname from the `proxy_pass` URL.
```nginx
# WRONG: Sends frontend's host to backend
proxy_set_header Host $host;

# CORRECT: Let nginx use backend's hostname from proxy_pass
# (don't set Host header at all)
```

### 7. Use HTTPS for Backend Communication
On platforms like Render that force HTTPS, always use HTTPS URLs for backend communication to avoid 307 redirects.

### 8. Enable SNI for HTTPS Proxies
When using HTTPS with `proxy_pass`, enable `proxy_ssl_server_name on` for proper SSL handshake on shared infrastructure.

## Environment Variables Required

| Variable | Description | Default | Notes |
|----------|-------------|---------|-------|
| `PORT` | Port nginx listens on | `10000` | Render's default |
| `BACKEND_URL` | Backend API URL for proxy | (required) | **Must use HTTPS** (e.g., `https://backend.onrender.com`) |

## References

- [nginx HTTP Proxy Module Documentation](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
- [Stack Overflow: envsubst with nginx in Docker](https://stackoverflow.com/questions/56649582/substitute-environment-variables-in-nginx-config-from-docker-compose)
- [GitHub Issue: envsubst conflicts with internal variables](https://github.com/nginxinc/docker-nginx/issues/529)
- [Unix StackExchange: Replacing only specific variables with envsubst](https://unix.stackexchange.com/questions/294378/replacing-only-specific-variables-with-envsubst)

## Preventive Measures

1. **Test envsubst locally** before deploying:
   ```bash
   docker run --rm -e PORT=8080 -e BACKEND_URL=http://localhost:3000 your-image nginx -T
   ```

2. **Use ${VAR} syntax** for environment variables in templates to visually distinguish from nginx variables.

3. **Always specify SHELL-FORMAT** when using envsubst with nginx configs.

4. **Consider using nginx-plus** with `env` directive for production environments needing complex variable handling.

5. **Test proxy configuration** by testing from local dev to deployed backend before deploying frontend:
   ```typescript
   // In vite.config.ts
   proxy: {
     '/api': {
       target: 'https://your-backend.onrender.com',
       changeOrigin: true,
     }
   }
   ```

6. **Always use HTTPS** for backend URLs on platforms that force HTTPS (Render, Heroku, etc.).

7. **Enable SNI** when proxying to HTTPS backends: `proxy_ssl_server_name on;`
