#!/usr/bin/env powershell
# Quick Start: Run Itinerary Platform Locally
# Windows PowerShell Script

Write-Host "🚀 Itinerary Platform - Local Setup" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is installed
Write-Host "📦 Checking Docker installation..." -ForegroundColor Yellow
if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Docker not found! Download from https://www.docker.com/products/docker-desktop" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Docker found" -ForegroundColor Green

# Check if docker-compose is available
Write-Host ""
Write-Host "📦 Checking Docker Compose..." -ForegroundColor Yellow
if (!(Get-Command docker-compose -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Docker Compose not found!" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Docker Compose found" -ForegroundColor Green

# Start Docker Compose
Write-Host ""
Write-Host "🐳 Starting Docker Compose services..." -ForegroundColor Yellow
Write-Host "   - PostgreSQL Database (port 5432)" -ForegroundColor Cyan
Write-Host "   - Redis Cache (port 6379)" -ForegroundColor Cyan
Write-Host "   - Go Backend (port 8080)" -ForegroundColor Cyan
Write-Host ""

Set-Location -Path "itinerary-backend"
docker-compose up --build

Write-Host ""
Write-Host "🎉 Application started!" -ForegroundColor Green
Write-Host ""
Write-Host "📍 Access URLs:" -ForegroundColor Cyan
Write-Host "   • Frontend: http://localhost:8080" -ForegroundColor Cyan
Write-Host "   • API: http://localhost:8080/api" -ForegroundColor Cyan
Write-Host "   • Database: localhost:5432" -ForegroundColor Cyan
Write-Host "   • Cache: localhost:6379" -ForegroundColor Cyan
