#!/bin/bash
# Quick Start: Run Itinerary Platform Locally
# Linux/macOS Bash Script

echo "🚀 Itinerary Platform - Local Setup"
echo "==================================="
echo ""

# Check if Docker is installed
echo "📦 Checking Docker installation..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found! Download from https://www.docker.com/products/docker-desktop"
    exit 1
fi
echo "✅ Docker found"

# Check if docker-compose is available
echo ""
echo "📦 Checking Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose not found!"
    exit 1
fi
echo "✅ Docker Compose found"

# Start Docker Compose
echo ""
echo "🐳 Starting Docker Compose services..."
echo "   - PostgreSQL Database (port 5432)"
echo "   - Redis Cache (port 6379)"
echo "   - Go Backend (port 8080)"
echo ""

cd itinerary-backend
docker-compose up --build

echo ""
echo "🎉 Application started!"
echo ""
echo "📍 Access URLs:"
echo "   • Frontend: http://localhost:8080"
echo "   • API: http://localhost:8080/api"
echo "   • Database: localhost:5432"
echo "   • Cache: localhost:6379"
