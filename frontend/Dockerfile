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
# Use envsubst to substitute PORT and BACKEND_URL in nginx config
# Shell form (not JSON array) for proper variable expansion
CMD /bin/sh -c 'envsubst "$$PORT $$BACKEND_URL" < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && exec nginx -g "daemon off;"'
