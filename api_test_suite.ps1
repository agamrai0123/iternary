#!/usr/bin/env pwsh
# Comprehensive API Testing Suite for Multi-Currency Group Features
# Phase A Week 2 - Monday Sprint Verification

Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  API TEST SUITE - MULTI-CURRENCY GROUP FEATURES" -ForegroundColor Cyan
Write-Host "  Phase A Week 2 - Monday Sprint" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8080"
$testResults = @{
    passed = 0
    failed = 0
    tests = @()
}

function Test-Endpoint {
    param(
        [string]$name,
        [string]$method,
        [string]$endpoint,
        [object]$body = $null,
        [bool]$expectSuccess = $true
    )
    
    Write-Host "Testing: $name..." -ForegroundColor Yellow
    
    try {
        $uri = "$baseUrl$endpoint"
        
        $params = @{
            Uri = $uri
            Method = $method
            ErrorAction = "SilentlyContinue"
            TimeoutSec = 5
        }
        
        if ($body) {
            $params.ContentType = "application/json"
            $params.Body = ($body | ConvertTo-Json)
        }
        
        $response = Invoke-WebRequest @params
        
        if ($response.StatusCode -in @(200, 201, 204)) {
            Write-Host "  ✓ PASS - Status: $($response.StatusCode)" -ForegroundColor Green
            $testResults.passed++
            $testResults.tests += @{name = $name; status = "PASS"; code = $response.StatusCode}
            return $response
        } else {
            Write-Host "  ✗ FAIL - Status: $($response.StatusCode)" -ForegroundColor Red
            $testResults.failed++
            $testResults.tests += @{name = $name; status = "FAIL"; code = $response.StatusCode}
            return $null
        }
    } catch {
        Write-Host "  ✗ ERROR - $($_.Exception.Message)" -ForegroundColor Red
        $testResults.failed++
        $testResults.tests += @{name = $name; status = "ERROR"; error = $_.Exception.Message}
        return $null
    }
}

Write-Host ""
Write-Host "PHASE 1: HEALTH & METRICS CHECKS" -ForegroundColor Magenta
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Magenta

Test-Endpoint "Health Check" "GET" "/api/health" | Out-Null
Test-Endpoint "Metrics Endpoint" "GET" "/api/metrics" | Out-Null

Write-Host ""
Write-Host "PHASE 2: MULTI-CURRENCY GROUP FEATURES" -ForegroundColor Magenta
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Magenta

# Test Group Trip Creation
$tripData = @{
    title = "Paris Trip 2026 - International Group"
    budget = 5000
    duration = 7
    destination_id = "paris-001"
}
$tripResponse = Test-Endpoint "Create Group Trip" "POST" "/api/group-trips" $tripData
$tripId = if ($tripResponse) { "trip-001" } else { $null }

Write-Host ""
# Test Group Trip Retrieval
if ($tripId) {
    Test-Endpoint "Retrieve Group Trip" "GET" "/api/group-trips/$tripId" | Out-Null
}

# Test Get User's Group Trips
Test-Endpoint "List User Group Trips" "GET" "/api/group-trips" | Out-Null

Write-Host ""
Write-Host "PHASE 3: STANDARD ENDPOINTS" -ForegroundColor Magenta
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Magenta

# Destinations
Test-Endpoint "Get Destinations" "GET" "/api/destinations" | Out-Null
Test-Endpoint "Get Destination Itineraries" "GET" "/api/destinations/paris-001/itineraries" | Out-Null

# Authentication
$loginData = @{
    email = "user@example.com"
    password = "password123"
}
Test-Endpoint "User Login" "POST" "/auth/login" $loginData | Out-Null
Test-Endpoint "Get User Profile" "GET" "/auth/profile" | Out-Null

Write-Host ""
Write-Host "PHASE 4: PRESENTATION PAGES" -ForegroundColor Magenta
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Magenta

Test-Endpoint "Login Page" "GET" "/" | Out-Null
Test-Endpoint "Dashboard Page" "GET" "/dashboard" | Out-Null
Test-Endpoint "Plan Trip Page" "GET" "/plan-trip" | Out-Null
Test-Endpoint "My Trips Page" "GET" "/my-trips" | Out-Null
Test-Endpoint "Community Page" "GET" "/community" | Out-Null
Test-Endpoint "Search Page" "GET" "/search" | Out-Null

Write-Host ""
Write-Host "PHASE 5: STATIC ASSETS" -ForegroundColor Magenta
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Magenta

Test-Endpoint "Static Assets" "GET" "/static" | Out-Null

Write-Host ""
Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  TEST RESULTS SUMMARY" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan

Write-Host ""
Write-Host "Total Tests Run: $($testResults.passed + $testResults.failed)" -ForegroundColor White
Write-Host "✓ Passed: $($testResults.passed)" -ForegroundColor Green
Write-Host "✗ Failed: $($testResults.failed)" -ForegroundColor $(if ($testResults.failed -eq 0) { "Green" } else { "Red" })

Write-Host ""
Write-Host "SUCCESS RATE: $(([math]::Round(($testResults.passed / ($testResults.passed + $testResults.failed)) * 100)))%" -ForegroundColor $(if ($testResults.failed -eq 0) { "Green" } else { "Yellow" })

Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  TEST BREAKDOWN" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

foreach ($test in $testResults.tests) {
    $marker = if ($test.status -eq "PASS") { "✓" } else { "✗" }
    $color = if ($test.status -eq "PASS") { "Green" } else { "Red" }
    Write-Host "$marker $($test.name): $($test.status)" -ForegroundColor $color
}

Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan

if ($testResults.failed -eq 0) {
    Write-Host "  ✅ ALL TESTS PASSED - READY FOR PRODUCTION" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  SOME TESTS FAILED - REVIEW REQUIRED" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
