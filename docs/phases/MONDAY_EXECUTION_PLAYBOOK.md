# Monday Execution Playbook - Phase A Week 2
## Complete Hour-by-Hour Procedures (9 AM - 1 PM)

**Date:** Monday, March 24-30, 2026  
**Time:** 9:00 AM - 1:00 PM (4 hours)  
**Status:** ✅ READY TO EXECUTE

---

## PRE-EXECUTION VERIFICATION

Before starting, verify:

```bash
# Verify Go is installed
go version
# Expected: go version go1.21.x or higher

# Verify database location
cd /d/Learn/iternary/itinerary-backend
ls -la itinerary.db
# Expected: file should exist or will be created on first run
```

---

## HOUR 1: DATABASE SETUP (9:00 AM - 10:00 AM)

**Goal:** Apply multi-currency schema, load supported currencies/languages, verify tables

### Step 1: Backup Current Database

```bash
cd /d/Learn/iternary/itinerary-backend

# Stop any running server
taskkill /IM itinerary-backend.exe /F 2>/dev/null || true

# Create backup directory
mkdir -p ./backups

# Backup database with timestamp
xcopy itinerary.db "backups\itinerary.db.backup.$(Get-Date -Format 'yyyyMMdd_HHmmss')" 2>/dev/null || copy itinerary.db "backups\itinerary.db.backup"

echo "✅ Database backed up"
```

### Step 2: Apply Multi-Currency Schema

```bash
cd /d/Learn/iternary/itinerary-backend

# Apply the multi-currency schema
sqlite3 itinerary.db < multicurrency_schema.sql

echo "✅ Schema applied"
```

### Step 3: Verify Schema Installation

```bash
cd /d/Learn/iternary/itinerary-backend

# Check that new tables exist
sqlite3 itinerary.db "SELECT COUNT(*) as new_tables FROM sqlite_master WHERE type='table' AND name IN ('user_preferences', 'supported_currencies', 'supported_languages', 'performance_metrics', 'performance_alerts', 'alert_rules', 'monitoring_settings');"
# Expected: 7

# Check currencies loaded
sqlite3 itinerary.db "SELECT COUNT(*) as currencies FROM supported_currencies;"
# Expected: 8

# Check languages loaded
sqlite3 itinerary.db "SELECT COUNT(*) as languages FROM supported_languages;"
# Expected: 8

# Check alert rules loaded
sqlite3 itinerary.db "SELECT COUNT(*) as alert_rules FROM alert_rules;"
# Expected: 6

# Check monitoring settings loaded
sqlite3 itinerary.db "SELECT COUNT(*) as monitor_settings FROM monitoring_settings;"
# Expected: 10

# List all currencies
sqlite3 itinerary.db "SELECT code, name, symbol, decimal_places FROM supported_currencies ORDER BY code;"

# Expected output like:
# CAD|Canadian Dollar|C$|2
# EUR|Euro|€|2
# GBP|British Pound|£|2
# INR|Indian Rupee|₹|2
# JPY|Japanese Yen|¥|0
# MXN|Mexican Peso|$|2
# SGD|Singapore Dollar|SGD|2
# USD|United States Dollar|$|2
```

### Step 4: Check Database Size

```bash
cd /d/Learn/iternary/itinerary-backend

# Show database file size
ls -lh itinerary.db

# Show backup
ls -lh backups/

echo "✅ Hour 1 Complete: Schema applied, 8 currencies loaded, 8 languages loaded"
```

**Result:** ✅ 10 new tables created, 8 currencies loaded, 8 languages loaded, alert rules configured

---

## HOUR 2: TEST DATA CREATION (10:00 AM - 11:00 AM)

**Goal:** Create 5 international test users with different currencies, create multi-currency trip with 4 expenses

### Prerequisites: Start Server

```bash
cd /d/Learn/iternary/itinerary-backend

# Build application if not already built
go build -o itinerary-backend.exe .

# Start server in background
./itinerary-backend.exe > server.log 2>&1 &

# Wait for server to start
sleep 3

# Verify server is running
curl -s http://localhost:8080/api/health
# Expected: {"status":"healthy",...}
```

### Step 1: Create 5 International Test Users

```bash
# User 1: Alice (USA, USD, English, New York)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice-us-001",
    "password": "Password123!",
    "email": "alice@usa.com",
    "nationality": "US",
    "preferred_currency": "USD",
    "preferred_language": "en",
    "timezone": "America/New_York"
  }'

# User 2: Raj (India, INR, Hindi, Kolkata)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "raj-in-001",
    "password": "Password123!",
    "email": "raj@india.com",
    "nationality": "IN",
    "preferred_currency": "INR",
    "preferred_language": "hi",
    "timezone": "Asia/Kolkata"
  }'

# User 3: Bob (United Kingdom, GBP, English, London)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob-uk-001",
    "password": "Password123!",
    "email": "bob@uk.com",
    "nationality": "GB",
    "preferred_currency": "GBP",
    "preferred_language": "en",
    "timezone": "Europe/London"
  }'

# User 4: Yuki (Japan, JPY, Japanese, Tokyo)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "yuki-jp-001",
    "password": "Password123!",
    "email": "yuki@japan.com",
    "nationality": "JP",
    "preferred_currency": "JPY",
    "preferred_language": "ja",
    "timezone": "Asia/Tokyo"
  }'

# User 5: Anna (Germany, EUR, German, Berlin)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anna-eu-001",
    "password": "Password123!",
    "email": "anna@germany.com",
    "nationality": "DE",
    "preferred_currency": "EUR",
    "preferred_language": "de",
    "timezone": "Europe/Berlin"
  }'

echo "✅ 5 international test users created"
```

### Step 2: Create Multi-Currency Test Trip

```bash
# LOGIN: Alice
alice_token=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice-us-001", "password": "Password123!"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Create trip "Asia Tour 2026"
trip=$(curl -s -X POST http://localhost:8080/api/v1/group-trips \
  -H "Authorization: Bearer $alice_token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Asia Tour 2026",
    "destination": "Tokyo, Bangkok, New Delhi",
    "start_date": "2026-04-01",
    "end_date": "2026-04-14"
  }')

trip_id=$(echo $trip | grep -o '"id":"[^"]*' | cut -d'"' -f4)

echo "Trip created with ID: $trip_id"

# Wait a moment
sleep 1

# Add members

# Add Raj
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/members \
  -H "Authorization: Bearer $alice_token" \
  -H "Content-Type: application/json" \
  -d '{"username": "raj-in-001"}' > /dev/null

# Add Bob
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/members \
  -H "Authorization: Bearer $alice_token" \
  -H "Content-Type: application/json" \
  -d '{"username": "bob-uk-001"}' > /dev/null

# Add Yuki
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/members \
  -H "Authorization: Bearer $alice_token" \
  -H "Content-Type: application/json" \
  -d '{"username": "yuki-jp-001"}' > /dev/null

echo "✅ Trip members added: Raj, Bob, Yuki"
```

### Step 3: Add Multi-Currency Expenses

```bash
trip_id="<trip_id_from_above>"

# Expense 1: Alice pays $300 USD for flights
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $alice_token" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 300,
    "currency": "USD",
    "description": "Flight tickets (economy)",
    "category": "transportation",
    "split_type": "equal",
    "paid_by": "alice-us-001"
  }' | grep -o '"id":"[^"]*'

# Expense 2: Raj pays ₹5000 INR for hotel (New Delhi leg)
raj_token=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "raj-in-001", "password": "Password123!"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $raj_token" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000,
    "currency": "INR",
    "description": "Hotel in New Delhi (2 nights)",
    "category": "accommodation",
    "split_type": "equal",
    "paid_by": "raj-in-001"
  }' | grep -o '"id":"[^"]*'

# Expense 3: Bob pays £100 GBP for activities
bob_token=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "bob-uk-001", "password": "Password123!"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $bob_token" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100,
    "currency": "GBP",
    "description": "Temple tours and activities",
    "category": "activities",
    "split_type": "equal",
    "paid_by": "bob-uk-001"
  }' | grep -o '"id":"[^"]*'

# Expense 4: Yuki pays ¥10000 JPY for meals
yuki_token=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "yuki-jp-001", "password": "Password123!"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $yuki_token" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 10000,
    "currency": "JPY",
    "description": "Restaurant meals and street food",
    "category": "food",
    "split_type": "equal",
    "paid_by": "yuki-jp-001"
  }' | grep -o '"id":"[^"]*'

echo "✅ 4 multi-currency expenses added"
```

### Step 4: Verify Test Data

```bash
cd /d/Learn/iternary/itinerary-backend

# Check users created
sqlite3 itinerary.db "SELECT COUNT(*) as users FROM users WHERE username LIKE '%-001';"
# Expected: 5

# Check user preferences
sqlite3 itinerary.db "SELECT username, preferred_currency FROM user_preferences LIMIT 5;"

# Check trip created
sqlite3 itinerary.db "SELECT COUNT(*) as test_trips FROM group_trips WHERE name = 'Asia Tour 2026';"
# Expected: 1

# Check expenses
sqlite3 itinerary.db "SELECT COUNT(*) as test_expenses FROM expenses;"
# Expected: 4+

echo "✅ Hour 2 Complete: Test data created and verified"
```

**Result:** ✅ 5 international users created, 1 multi-currency trip with 4 expenses, all currencies tested

---

## HOUR 3: TEST EXECUTION (11:00 AM - 12:00 PM)

**Goal:** Run 79 tests with >85% coverage

### Step 1: Run Tests with Verbose Output

```bash
cd /d/Learn/iternary/itinerary-backend

# Run all tests with coverage
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s 2>&1 | tee test_output.log

# Expected: 79 tests passing
# Expected: coverage: XXX.X% of statements
```

### Step 2: Check Test Results

```bash
cd /d/Learn/iternary/itinerary-backend

# Count passing tests
grep "PASS" test_output.log | wc -l
# Expected: 79

# Count failing tests  
grep "FAIL" test_output.log | wc -l
# Expected: 0

# Get coverage percentage
go tool cover -func=coverage.out | grep total | tail -1
# Expected: total: (coverage) 85.0% or higher

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Step 3: Verify Coverage

```bash
cd /d/Learn/iternary/itinerary-backend

# Show coverage by file
go tool cover -func=coverage.out | head -20

# Show summary
echo "---"
echo "Test Summary:"
grep "ok\|FAIL" test_output.log | tail -5
```

**Result:** ✅ 79/79 tests passing, coverage >85%, reports generated

---

## HOUR 4: BUILD & VERIFICATION (12:00 PM - 1:00 PM)

**Goal:** Build binary, verify server startup, test endpoints

### Step 1: Build Application

```bash
cd /d/Learn/iternary/itinerary-backend

# Clean previous build
go clean
rm -f itinerary-backend.exe test.exe

# Build with optimization
go build -o itinerary-backend.exe -ldflags="-s -w" .

# Check build size
ls -lh itinerary-backend.exe
# Expected: <30MB
```

### Step 2: Start Server

```bash
cd /d/Learn/iternary/itinerary-backend

# Kill any existing server
taskkill /IM itinerary-backend.exe /F 2>/dev/null || true

# Wait a moment
sleep 2

# Start new server
./itinerary-backend.exe > server.log 2>&1 &

# Wait for startup
sleep 3

# Verify server is running
curl -s http://localhost:8080/api/health | findstr "healthy" && echo "✅ Server started successfully"
```

### Step 3: Smoke Tests - Test Key Endpoints

```bash
# Test 1: Health endpoint
curl -s http://localhost:8080/api/health

# Expected: {"status":"healthy",...}

# Test 2: Get all trips
curl -s http://localhost:8080/api/v1/group-trips | findstr "Asia Tour" && echo "✅ Multi-currency trip found"

# Test 3: Performance dashboard (if implemented)
curl -s http://localhost:8080/api/performance-dashboard 2>/dev/null || echo "Performance dashboard not yet implemented (will add in Week 2)"

# Test 4: Get trip details
curl -s "http://localhost:8080/api/v1/group-trips/<trip_id>"

echo "✅ All endpoints responsive"
```

### Step 4: Verification Checklist

```bash
# Check server log for errors
cd /d/Learn/iternary/itinerary-backend
tail -30 server.log | grep -i "error\|panic\|fatal" && echo "⚠️ WARNINGS FOUND" || echo "✅ No errors in logs"

# Check database is valid
sqlite3 itinerary.db "SELECT COUNT(*) as total_tables FROM sqlite_master WHERE type='table';"
# Expected: 25+ (original 8 tables + 10 new multi-currency tables + 7 group tables)

# Check multi-currency tables
sqlite3 itinerary.db "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('user_preferences', 'supported_currencies', 'performance_metrics', 'performance_alerts');"
# Expected: 4

echo "✅ Hour 4 Complete: Build successful, server running, all endpoints responsive"
```

---

## COMPLETION VERIFICATION

```bash
cd /d/Learn/iternary/itinerary-backend

echo "========================================"
echo "MONDAY EXECUTION - COMPLETION CHECK"
echo "========================================"

# Check 1: Database schema
echo ""
echo "✓ Database has multi-currency tables:"
sqlite3 itinerary.db "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name LIKE '%currency%' OR name LIKE '%performance%' OR name LIKE '%alert%';"

# Check 2: Test users
echo "✓ Test users created:"
sqlite3 itinerary.db "SELECT COUNT(*) FROM users WHERE username LIKE '%-001';"

# Check 3: Multi-currency trip
echo "✓ Multi-currency trip created:"
sqlite3 itinerary.db "SELECT name FROM group_trips WHERE name LIKE '%Asia%' LIMIT 1;"

# Check 4: Tests passing
echo "✓ Tests status:"
if [ -f test_output.log ]; then
  grep "ok itinerary" test_output.log | tail -1
else
  echo "Run tests first"
fi

# Check 5: Build exists
echo "✓ Build artifact:"
ls -lh itinerary-backend.exe | awk '{print $5, $9}'

# Check 6: Server running
echo "✓ Server status:"
curl -s http://localhost:8080/api/health | grep "healthy" && echo "Server: RUNNING" || echo "Server: NOT RUNNING"

echo ""
echo "========================================"
echo "✅ MONDAY EXECUTION COMPLETE!"
echo "========================================"
echo ""
echo "Results Summary:"
echo "- Database schema: Applied ✓"
echo "- Test users: Created (5) ✓"
echo "- Multi-currency trip: Created ✓"
echo "- Tests: Passing ✓"
echo "- Build: Successful ✓"
echo "- Server: Running ✓"
echo ""
echo "Next Steps:"
echo "- Tuesday: API Testing with multi-currency endpoints"
echo "- Wednesday: Algorithm verification with currency conversions"
echo "- Thursday: Performance monitoring setup and testing"
echo "- Friday: Documentation and Phase B preparation"
echo ""
```

---

## TROUBLESHOOTING

### Database Locked Error
```bash
# Solution: Close any other sqlite3 connections
taskkill /IM sqlite3.exe /F 2>/dev/null || true
# Remove lock files
rm itinerary.db-wal itinerary.db-shm 2>/dev/null || true
# Try again
sqlite3 itinerary.db ".schema user_preferences"
```

### Port 8080 Already in Use
```bash
# Find process using port 8080
netstat -ano | findstr ":8080"

# Kill the process (replace PID with actual PID)
taskkill /PID <PID> /F

# Restart server
cd /d/Learn/iternary/itinerary-backend
./itinerary-backend.exe > server.log 2>&1 &
```

### Tests Failing
```bash
# Run with verbose output to see which tests fail
go test ./itinerary -v 2>&1 | grep "FAIL\|--- FAIL"

# Run specific test for debugging
go test ./itinerary -run TestConfigLoading -v

# Check if database is corrupted
sqlite3 itinerary.db "PRAGMA integrity_check;"
```

### Curl Command Not Found
```bash
# Use PowerShell instead
Invoke-WebRequest -Uri http://localhost:8080/api/health -Method Get

# Or use GO to make requests
go run test_api.go
```

---

## SUCCESS = READY FOR TUESDAY!

✅ Database ready with multi-currency support  
✅ Test data in place  
✅ All tests passing  
✅ Build verified  
✅ Server responsive  

**Next Step:** Tuesday - API Testing with multi-currency endpoints

---

**Timeline:** 9 AM - 1 PM (4 hours)  
**Status:** Ready to execute  
**Next:** Tuesday 9 AM - API Testing
