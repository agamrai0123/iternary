# Phase A Week 2: Quick Start Checklist

**Start Date:** Monday  
**Duration:** 1 week  
**Objective:** Multi-currency support + Performance monitoring  
**Status:** ✅ READY TO EXECUTE

---

## 📋 Pre-Execution Checklist

Before Monday morning, verify:

- [ ] Go environment ready: `go version` (expect 1.21+)
- [ ] Database backup exists: `itinerary.db.phase_a_week_1_backup`
- [ ] Test data from Week 1 has 10+ trips
- [ ] All 79 tests passing: `go test ./itinerary -v`
- [ ] Documentation files downloaded
- [ ] Team notified of Week 2 start
- [ ] Monitoring setup ready (optional: Slack integration)

---

## 🚀 Monday Morning - 4-Hour Sprint

### Goal: Foundation Ready (Multi-currency DB + Tests Passing)

**Timeline:** 9 AM - 1 PM (4 hours)

#### Hour 1: Database Setup (0900-1000)

**Files to execute:**

1. **Apply Schema**
   ```bash
   cd itinerary-backend
   sqlite3 itinerary.db < ../PHASE_A_WEEK_2_MULTICURRENCY_SCHEMA.sql
   ```
   Expected: 10 new tables created, 8 currencies loaded, 8 languages loaded

2. **Verify Schema**
   ```bash
   sqlite3 itinerary.db ".schema user_preferences"
   sqlite3 itinerary.db "SELECT COUNT(*) FROM supported_currencies;"  # expect 8
   sqlite3 itinerary.db "SELECT COUNT(*) FROM supported_languages;"   # expect 8
   ```

3. **Check Time:** Should take ~5 minutes
   - [ ] Database backup exists
   - [ ] Schema applied successfully
   - [ ] Tables show in sqlite3 shell

---

#### Hour 2: Test Data Creation (1000-1100)

**Create 5 international test users:**

```bash
# User 1: Alice (USA, USD, English, New York)
curl -X POST http://localhost:8080/api/v1/auth/register \
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
curl -X POST http://localhost:8080/api/v1/auth/register \
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
curl -X POST http://localhost:8080/api/v1/auth/register \
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
curl -X POST http://localhost:8080/api/v1/auth/register \
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
curl -X POST http://localhost:8080/api/v1/auth/register \
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
```

**Create multi-currency trip:**

```bash
# Login Alice
alice_token=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice-us-001", "password": "Password123!"}' | jq -r .token)

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
trip_id=$(echo $trip | jq -r .id)

# Add members
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/members \
  -H "Authorization: Bearer $alice_token" \
  -d '{"username": "raj-in-001"}' > /dev/null
  
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/members \
  -H "Authorization: Bearer $alice_token" \
  -d '{"username": "bob-uk-001"}' > /dev/null
  
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/members \
  -H "Authorization: Bearer $alice_token" \
  -d '{"username": "yuki-jp-001"}' > /dev/null

# Add multi-currency expenses

# Expense 1: Alice pays $300 USD for flights
curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $alice_token" \
  -d '{
    "amount": 300,
    "currency": "USD",
    "description": "Flight tickets",
    "category": "transportation",
    "split_type": "equal",
    "paid_by": "alice-us-001"
  }' > /dev/null

# Expense 2: Raj pays ₹5000 INR for hotel (New Delhi leg)
raj_token=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -d '{"username": "raj-in-001", "password": "Password123!"}' | jq -r .token)

curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $raj_token" \
  -d '{
    "amount": 5000,
    "currency": "INR",
    "description": "Hotel in New Delhi (2 nights)",
    "category": "accommodation",
    "split_type": "equal",
    "paid_by": "raj-in-001"
  }' > /dev/null

# Expense 3: Bob pays £100 GBP for activities
bob_token=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -d '{"username": "bob-uk-001", "password": "Password123!"}' | jq -r .token)

curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $bob_token" \
  -d '{
    "amount": 100,
    "currency": "GBP",
    "description": "Temple tours and activities",
    "category": "activities",
    "split_type": "equal",
    "paid_by": "bob-uk-001"
  }' > /dev/null

# Expense 4: Yuki pays ¥10000 JPY for meals
yuki_token=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -d '{"username": "yuki-jp-001", "password": "Password123!"}' | jq -r .token)

curl -s -X POST http://localhost:8080/api/v1/group-trips/$trip_id/expenses \
  -H "Authorization: Bearer $yuki_token" \
  -d '{
    "amount": 10000,
    "currency": "JPY",
    "description": "Restaurant meals and street food",
    "category": "food",
    "split_type": "equal",
    "paid_by": "yuki-jp-001"
  }' > /dev/null
```

**Check Time:** Should take ~15 minutes
- [ ] 5 users created with different currencies
- [ ] 1 trip created with 4 members
- [ ] 4 expenses added in different currencies
- [ ] No API errors in curl responses

---

#### Hour 3: Test Execution (1100-1200)

**Run Full Test Suite:**

```bash
cd itinerary-backend

# Run all tests with verbose output and coverage
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s 2>&1 | tee test_output.log

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
go tool cover -func=coverage.out | grep total | awk '{print $3}'

# Expected output: coverage: XX.X% of statements
```

**Verify Test Results:**

```bash
# Check for passed tests
grep -c "PASS" test_output.log  # expect 79+

# Check for failed tests
grep "FAIL" test_output.log     # expect 0

# Open coverage report to see covered code
# Open coverage.html in browser
```

**Check Time:** Should take ~15 minutes
- [ ] All 79 tests pass (PASS in log)
- [ ] Coverage >85% of statements
- [ ] HTML coverage report generated
- [ ] No error logs

---

#### Hour 4: Build & Verification (1200-1300)

**Build Application:**

```bash
cd itinerary-backend

# Clean previous build
go clean
rm -f itinerary-backend.exe itinerary-backend

# Build with optimization
go build -o itinerary-backend.exe -ldflags="-s -w" .

# Verify build size
ls -lh itinerary-backend.exe  # expect <30MB

# Check for build warnings
go build -v . 2>&1 | grep -i warn
```

**Start Server:**

```bash
# Start in background
./itinerary-backend.exe &
server_pid=$!

# Wait for startup
sleep 2

# Test health endpoint
curl -s http://localhost:8080/api/v1/health | jq .

# Expected response:
# {
#   "status": "healthy",
#   "database": "connected",
#   "timestamp": "2026-03-24T14:30:00Z"
# }
```

**Smoke Test:**

```bash
# Test a few key endpoints

# 1. Get all trips
curl -s http://localhost:8080/api/v1/group-trips | jq . | head -20

# 2. Get performance dashboard
curl -s http://localhost:8080/api/performance-dashboard | jq .

# 3. Check for multi-currency trip
curl -s http://localhost:8080/api/v1/group-trips | jq '.[] | select(.name == "Asia Tour 2026")'
```

**Stop Server:**

```bash
kill $server_pid
```

**Check Time:** Should take ~15 minutes
- [ ] Build completes without errors
- [ ] Binary created: itinerary-backend.exe
- [ ] Server starts without errors
- [ ] Health endpoint responds
- [ ] Performance dashboard returns metrics
- [ ] API requests return valid JSON

---

### ✅ Monday Checklist

```
[ ] 9:00 AM - Database schema applied
    - 10 new tables created
    - 8 currencies loaded
    - 8 languages loaded

[ ] 10:00 AM - Test data created
    - 5 international users registered
    - 1 multi-currency trip created
    - 4 expenses in different currencies

[ ] 11:00 AM - Tests executed
    - 79 tests passed
    - Coverage >85%
    - HTML report generated

[ ] 12:00 PM - Build verified
    - Binary created successfully
    - Server starts cleanly
    - Endpoints responsive

[ ] 1:00 PM - Team sync
    - Report created: PHASE_A_WEEK_2_MONDAY_REPORT.md
    - Next steps assigned (Tuesday API testing)
    - Blockers identified (if any)
```

---

## 📊 Expected Monday Results

| Metric | Target | Status |
|--------|--------|--------|
| Database tables | 18 total (8 original + 10 new) | ✓ |
| Test users | 5 international | ✓ |
| Test trip | 1 multi-currency trip | ✓ |
| Test expenses | 4 in different currencies | ✓ |
| Tests passing | 79/79 (100%) | ? |
| Code coverage | >85% | ? |
| Build size | <30MB | ? |
| Server startup | <2 seconds | ? |
| API latency | <500ms avg | ? |

---

## 📅 Tuesday-Friday Quick Summary

### Tuesday: API Testing (3-4 hours)
- Test all 16 endpoints with multi-currency data
- Verify currency conversion in responses
- Test settlement calculations with mixed currencies
- **Expected:** All endpoints support currency parameter

### Wednesday: Algorithm Testing (2-3 hours)
- Run settlement algorithm with multi-currency expenses
- Test poll voting (should be unaffected)
- Verify calculation accuracy
- **Expected:** Settlements handle all currency combinations

### Thursday: Performance Monitoring (2-3 hours)
- Enable real-time metrics collection
- Trigger alerts intentionally
- Monitor P95/P99 response times
- **Expected:** Dashboard shows real-time metrics and alerts

### Friday: Release (2-3 hours)
- Documentation complete
- Performance baselines established
- Phase B kickoff prep
- **Expected:** Ready for Phase B

---

## 🔧 Troubleshooting - Monday

### Issue: "database locked" error

```bash
# SQLite is single-writer issue
# Solution: Close any other connections

# Kill any stray processes
pkill -f "sqlite3 itinerary.db"

# Remove lock file if exists
rm -f itinerary.db-wal itinerary.db-shm

# Try again
sqlite3 itinerary.db ".schema user_preferences"
```

### Issue: Tests fail with "table not found"

```bash
# Schema didn't apply correctly
# Solution: Reapply schema

sqlite3 itinerary.db < ../PHASE_A_WEEK_2_MULTICURRENCY_SCHEMA.sql

# Then re-run tests
go test ./itinerary -v
```

### Issue: Curl returns "401 Unauthorized"

```bash
# Token expired or missing
# Solution: Get new token

token=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice-us-001", "password": "Password123!"}' | jq -r .token)

# Use new token
curl -s http://localhost:8080/api/v1/group-trips \
  -H "Authorization: Bearer $token" | jq .
```

### Issue: Build fails with "module not found"

```bash
# Go modules issue
# Solution: Update modules

cd itinerary-backend
go mod tidy
go mod download
go build .
```

---

## 📞 Getting Help

**Document References:**

1. **PHASE_A_WEEK_2_ENHANCED_EXECUTION_PLAN.md** - Full design
2. **PHASE_A_WEEK_2_MULTICURRENCY_SCHEMA.sql** - Database schema
3. **PHASE_A_WEEK_2_MONDAY_KICKOFF.md** - Detailed procedures
4. **PHASE_A_WEEK_2_PERFORMANCE_MONITORING_GUIDE.md** - Monitoring setup

**Questions?**

- Database schema → See MULTICURRENCY_SCHEMA.sql comments
- API testing → See PHASE_A_GROUP_SCHEMA.sql examples
- Performance → See PERFORMANCE_MONITORING_GUIDE.md
- Tests → Run with `-v` flag for details

---

## 🎯 Success Criteria

✅ **Monday Must-Have:**
- Database schema applied with 0 errors
- 5 test users created in database
- Multi-currency trip with 4 expenses created
- All 79 tests passing
- Build successful, binary created
- Server starts and responds to endpoints

✅ **Week Must-Have:**
- Multi-currency support implemented end-to-end
- Performance monitoring collecting metrics
- Real-time alerts working
- All documentation updated
- Team trained on new features
- Ready for Phase B

---

**Phase A Week 2 is ready to execute!** 🚀

Start Monday at 9 AM. All commands provided. Expected completion: 1 week.

---

**Created:** 2026-03-23  
**Updated:** 2026-03-24  
**Status:** ✅ READY
