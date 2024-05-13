#!/bin/bash

# Stop execution if any command fails
set -e

# Export the environment variable
export VITE_API_BASE_URL="http://api/v1/"

# Navigate to the frontend directory
cd frontend

# Build the frontend
echo "Building the frontend..."
pnpm run build

# Navigate back to the root directory where docker-compose.yml is located
cd ..

# Start docker-compose services
echo "Starting Docker Compose services..."
docker-compose up -d --build

# Optional: if you want docker-compose to run in detached mode, use `docker-compose up -d`
