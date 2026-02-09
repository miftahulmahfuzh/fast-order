# Postmortem: Frontend Deployment to Render

**Date:** 2026-02-09
**Service:** Fast Order Frontend (React + Vite + Nginx)
**Platform:** Render.com

## Executive Summary

Deploying the frontend to Render required 6 iterative fixes to properly configure nginx environment variable substitution. The core issue was that `envsubst` was replacing nginx's built-in runtime variables (`$host`, `$remote_addr`, `$scheme`, etc.) with empty strings, causing nginx configuration errors.

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
    location /api/ {
        proxy_pass $BACKEND_URL;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
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

## Environment Variables Required

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port nginx listens on | `10000` |
| `BACKEND_URL` | Backend API URL for proxy | (required) |

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
