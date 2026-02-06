#!/bin/bash
set -e

echo "Stopping containers..."
docker compose down

echo "Building images..."
docker compose build --no-cache

echo "Starting containers..."
docker compose up -d

echo "Deployment complete!"
docker compose ps
