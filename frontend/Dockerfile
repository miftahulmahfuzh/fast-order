# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Runtime stage
FROM nginx:alpine
# Install envsubst for variable substitution
RUN apk add --no-cache gettext
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/default.conf.template
# Use sed to substitute __PORT__ and __BACKEND_URL__ placeholders
# Default PORT to 10000 if not set (Render's default)
CMD /bin/sh -c 'export PORT=${PORT:-10000} && sed -e "s|__PORT__|$$PORT|g" -e "s|__BACKEND_URL__|$$BACKEND_URL|g" /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && exec nginx -g "daemon off;"'
