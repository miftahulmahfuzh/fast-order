#!/bin/bash
set -e

echo "=== Fast Order Integration Test ==="
echo ""

# Check if containers are running
if ! docker compose ps | grep -q "Up"; then
  echo "Containers not running. Starting with deploy-compose.sh..."
  ./deploy-compose.sh
  echo "Waiting for services to be ready..."
  sleep 5
fi

echo "Installing Playwright..."
cd frontend
npm install -D @playwright/test
npx playwright install --with-deps chromium

echo ""
echo "Running integration tests against Docker deployment..."
npx playwright test docker-integration.spec.ts --config=playwright.config.ts --reporter=list

echo ""
echo "=== Integration tests complete! ==="
