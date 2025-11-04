#!/bin/bash

# Script to fix and deploy the application on the server
# Usage: ./fix_and_deploy.sh

set -e  # Exit on error

echo "=== Excel Template Engine - Fix and Deploy ==="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed!"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose is not installed!"
    exit 1
fi

# Display versions
print_info "Docker version: $(docker --version)"
print_info "Docker Compose version: $(docker-compose --version)"
echo ""

# Step 1: Stop existing containers
print_info "Stopping existing containers..."
docker-compose down || true
echo ""

# Step 2: Remove old images (optional but recommended)
print_warn "Removing old images..."
docker rmi acts-service 2>/dev/null || true
docker rmi excel-template-engine-app 2>/dev/null || true
docker rmi $(docker images -f "dangling=true" -q) 2>/dev/null || true
echo ""

# Step 3: Clean build cache
print_info "Cleaning build cache..."
docker builder prune -f || true
echo ""

# Step 4: Build new images
print_info "Building new images (this may take a few minutes)..."
docker-compose build --no-cache
echo ""

# Step 5: Start services
print_info "Starting services..."
docker-compose up -d
echo ""

# Step 6: Wait for services to be ready
print_info "Waiting for services to start..."
sleep 10

# Step 7: Check status
print_info "Checking container status..."
docker-compose ps
echo ""

# Step 8: Display logs
print_info "Recent logs:"
docker-compose logs --tail=50
echo ""

# Step 9: Test connectivity
print_info "Testing service health..."
echo ""

# Wait a bit more for the app to fully start
sleep 5

# Try to connect to the app
if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    print_info "✓ Application is responding on http://localhost:8080"
else
    print_warn "Application might still be starting or health endpoint not available"
    print_info "Check logs with: docker-compose logs -f app"
fi

# Check MongoDB
if docker exec acts-mongodb mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    print_info "✓ MongoDB is running"
else
    print_warn "MongoDB might still be starting"
fi

echo ""
print_info "=== Deployment Complete ==="
echo ""
print_info "Useful commands:"
echo "  View logs:          docker-compose logs -f"
echo "  View app logs:      docker-compose logs -f app"
echo "  View mongo logs:    docker-compose logs -f mongodb"
echo "  Restart services:   docker-compose restart"
echo "  Stop services:      docker-compose down"
echo "  Check status:       docker-compose ps"
echo ""

