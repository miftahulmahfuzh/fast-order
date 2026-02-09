#!/bin/bash
set -e

echo "Installing Playwright..."
cd frontend
npm install -D @playwright/test
npx playwright install --with-deps chromium

echo "Running integration tests..."
npm run test

echo "Integration tests complete!"
