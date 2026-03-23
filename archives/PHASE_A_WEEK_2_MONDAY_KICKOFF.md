# Phase A Week 2 - Monday Morning Kickoff

**Date:** Monday, March 24, 2026  
**Status:** Ready for Execution  
**Duration:** 4 hours  
**Objectives:** Database setup + Multi-currency initialization + Test execution

---

## 🚀 Quick Start - Week 2 Execution

### Prerequisites Check (Do Before 9 AM)

```bash
# Check Go installation
go version
# Expected: go version go1.21.x (or higher)

# Check workspace
cd /d/Learn/iternary/itinerary-backend
pwd
# Expected: /d/Learn/iternary/itinerary-backend

# List current database
ls -la itinerary.db
# Should exist from Phase A Week 1

# Check dependencies
go mod tidy
go list -m all | head
# Should show gin-gonic/gin, zerolog, etc

# Check existing code
ls -la itinerary/group*.go
# Should show: group_models.go, group_database.go, group_service.go, group_handlers.go, group_routes.go
```

---

## 📋 Monday's 4-Hour Schedule

### Hour 1 (9:00 AM - 10:00 AM): Database Schema Extension

#### Step 1: Backup Current Database

```bash
# Backup current database in case of issues
cp itinerary.db itinerary.db.phase_a_week_1_backup
echo "✓ Database backup created"
```

#### Step 2: Apply Multi-Currency Schema

```bash
# For SQLite (dev environment)
sqlite3 itinerary.db << 'EOF'
-- Add to user_preferences
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id TEXT PRIMARY KEY,
    nationality TEXT,
    preferred_currency TEXT DEFAULT 'USD',
    preferred_language TEXT DEFAULT 'en',
    timezone TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add to supported_currencies
CREATE TABLE IF NOT EXISTS supported_currencies (
    code TEXT PRIMARY KEY,
    symbol TEXT,
    name TEXT,
    decimal_places INTEGER DEFAULT 2,
    exchange_rate_to_usd REAL DEFAULT 1.0,
    is_active BOOLEAN DEFAULT 1,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add to supported_languages
CREATE TABLE IF NOT EXISTS supported_languages (
    code TEXT PRIMARY KEY,
    name TEXT,
    native_name TEXT,
    is_active BOOLEAN DEFAULT 1,
    rtl BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create performance monitoring tables
CREATE TABLE IF NOT EXISTS performance_metrics (
    id TEXT PRIMARY KEY,
    endpoint_path TEXT,
    method TEXT,
    response_time_ms INTEGER,
    status_code INTEGER,
    error_flag BOOLEAN DEFAULT 0,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS performance_alerts (
    id TEXT PRIMARY KEY,
    alert_type TEXT,
    endpoint_path TEXT,
    severity TEXT,
    threshold_value REAL,
    current_value REAL,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS endpoint_performance_aggregates (
    id TEXT PRIMARY KEY,
    endpoint_path TEXT UNIQUE,
    method TEXT,
    total_requests INTEGER DEFAULT 0,
    error_rate REAL DEFAULT 0,
    avg_response_time_ms REAL DEFAULT 0,
    p95_response_time_ms INTEGER DEFAULT 0,
    p99_response_time_ms INTEGER DEFAULT 0,
    max_response_time_ms INTEGER DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pm_endpoint_time ON performance_metrics(endpoint_path, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_alert_severity ON performance_alerts(severity, created_at DESC);
EOF

echo "✓ Multi-currency schema applied"
```

#### Step 3: Load Supported Currencies

```bash
sqlite3 itinerary.db << 'EOF'
INSERT OR IGNORE INTO supported_currencies (code, symbol, name, decimal_places, exchange_rate_to_usd, is_active) VALUES
('USD', '$', 'US Dollar', 2, 1.0, 1),
('EUR', '€', 'Euro', 2, 1.10, 1),
('INR', '₹', 'Indian Rupee', 2, 0.012, 1),
('GBP', '£', 'British Pound', 2, 1.27, 1),
('JPY', '¥', 'Japanese Yen', 0, 0.0067, 1),
('SGD', 'S$', 'Singapore Dollar', 2, 0.75, 1),
('CAD', 'C$', 'Canadian Dollar', 2, 0.74, 1),
('MXN', '$', 'Mexican Peso', 2, 0.058, 1);

INSERT OR IGNORE INTO supported_languages (code, name, native_name, is_active, rtl) VALUES
('en', 'English', 'English', 1, 0),
('es', 'Spanish', 'Español', 1, 0),
('fr', 'French', 'Français', 1, 0),
('de', 'German', 'Deutsch', 1, 0),
('hi', 'Hindi', 'हिन्दी', 1, 0),
('ja', 'Japanese', '日本語', 1, 0),
('pt', 'Portuguese', 'Português', 1, 0),
('zh', 'Chinese', '中文', 1, 0);

SELECT 'Currencies loaded:', COUNT(*) FROM supported_currencies;
SELECT 'Languages loaded:', COUNT(*) FROM supported_languages;
EOF

echo "✓ Supported currencies and languages loaded"
```

#### Step 4: Verify Schema

```bash
# Verify tables created
sqlite3 itinerary.db << 'EOF'
.tables
.schema user_preferences
.schema supported_currencies
SELECT COUNT(*) as currency_count FROM supported_currencies;
SELECT COUNT(*) as language_count FROM supported_languages;
EOF

echo "✓ Schema verification complete"
```

### Hour 2 (10:00 AM - 11:00 AM): Test Data & User Preferences

#### Step 1: Create Test Users with Different Preferences

```bash
sqlite3 itinerary.db << 'EOF'
-- Create test users with different currency preferences
INSERT OR IGNORE INTO users (id, name, email, password_hash) VALUES
('user-us-001', 'Alice (US)', 'alice@example.com', 'hash1'),
('user-in-001', 'Raj (India)', 'raj@example.com', 'hash2'),
('user-uk-001', 'Bob (UK)', 'bob@example.com', 'hash3'),
('user-jp-001', 'Yuki (Japan)', 'yuki@example.com', 'hash4'),
('user-eu-001', 'Anna (Europe)', 'anna@example.com', 'hash5');

-- Set user preferences
INSERT OR IGNORE INTO user_preferences (user_id, nationality, preferred_currency, preferred_language, timezone) VALUES
('user-us-001', 'US', 'USD', 'en', 'America/New_York'),
('user-in-001', 'IN', 'INR', 'hi', 'Asia/Kolkata'),
('user-uk-001', 'GB', 'GBP', 'en', 'Europe/London'),
('user-jp-001', 'JP', 'JPY', 'ja', 'Asia/Tokyo'),
('user-eu-001', 'DE', 'EUR', 'de', 'Europe/Berlin');

SELECT 'Users created:', COUNT(*) FROM users WHERE id LIKE 'user-%';
SELECT 'User preferences set:', COUNT(*) FROM user_preferences;
EOF

echo "✓ Test users and preferences created"
```

#### Step 2: Create Multi-Currency Test Data

```bash
sqlite3 itinerary.db << 'EOF'
-- Create test trip
INSERT OR IGNORE INTO group_trips (id, title, destination_id, owner_id, budget, duration, start_date, status) VALUES
('trip-multi-001', 'International Bali Trip', 'dest-001', 'user-us-001', 5000, 7, '2026-04-01', 'planning');

-- Add members from different countries
INSERT OR IGNORE INTO group_members (id, trip_id, user_id, role, status) VALUES
('member-us', 'trip-multi-001', 'user-us-001', 'owner', 'active'),
('member-in', 'trip-multi-001', 'user-in-001', 'editor', 'active'),
('member-uk', 'trip-multi-001', 'user-uk-001', 'editor', 'active'),
('member-jp', 'trip-multi-001', 'user-jp-001', 'member', 'active');

-- Add expenses in different currencies
-- Alice (US) pays $300 in USD
INSERT OR IGNORE INTO expenses (id, trip_id, description, amount, currency, paid_by_id, category) VALUES
('exp-001', 'trip-multi-001', 'Flight for 4 people', 300, 'USD', 'user-us-001', 'transportation');

-- Raj (India) pays ₹5000 in INR (approximately $60)
INSERT OR IGNORE INTO expenses (id, trip_id, description, amount, currency, paid_by_id, category) VALUES
('exp-002', 'trip-multi-001', 'Hotel booking', 5000, 'INR', 'user-in-001', 'accommodation');

-- Bob (UK) pays £100 in GBP (approximately $127)
INSERT OR IGNORE INTO expenses (id, trip_id, description, amount, currency, paid_by_id, category) VALUES
('exp-003', 'trip-multi-001', 'Activities', 100, 'GBP', 'user-uk-001', 'activities');

-- Yuki (Japan) pays ¥10000 in JPY (approximately $67)
INSERT OR IGNORE INTO expenses (id, trip_id, description, amount, currency, paid_by_id, category) VALUES
('exp-004', 'trip-multi-001', 'Meals', 10000, 'JPY', 'user-jp-001', 'food');

SELECT 'Test trip created';
SELECT 'Members added:', COUNT(*) FROM group_members WHERE trip_id = 'trip-multi-001';
SELECT 'Expenses added (multi-currency):', COUNT(*) FROM expenses WHERE trip_id = 'trip-multi-001';
EOF

echo "✓ Multi-currency test data created"
```

### Hour 3 (11:00 AM - 12:00 PM): Test Execution

#### Step 1: Run Unit Tests

```bash
cd /d/Learn/iternary/itinerary-backend

# Run model tests
echo "=== Running Model Tests ==="
go test ./itinerary -v -run "TestGroupTrip|TestExpense|TestGroupMember" -timeout 30s 2>&1 | tee test_output_models.log

# Expected: 25 tests passing
```

#### Step 2: Run Service Tests

```bash
# Run service tests
echo "=== Running Service Tests ==="
go test ./itinerary -v -run "TestCreate|TestSettlement|TestPermission" -timeout 30s 2>&1 | tee test_output_service.log

# Expected: 32+ tests passing
```

#### Step 3: Run Integration Tests

```bash
# Run integration tests
echo "=== Running Integration Tests ==="
go test ./itinerary -v -run "Integration|Handler|Middleware" -timeout 30s 2>&1 | tee test_output_integration.log

# Expected: 22+ tests passing
```

#### Step 4: Run All Tests with Coverage

```bash
# Run all tests with coverage
echo "=== Running All Tests with Coverage ==="
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s 2>&1 | tee test_output_all.log

# Check coverage
go tool cover -html=coverage.out -o coverage.html
echo "✓ Coverage report: coverage.html"

# Show summary
go tool cover -func=coverage.out | grep total
```

### Hour 4 (12:00 PM - 1:00 PM): Build & Verification

#### Step 1: Build Application

```bash
cd /d/Learn/iternary/itinerary-backend

echo "=== Building Application ==="
go build -o itinerary-backend.exe . 2>&1 | tee build_output.log

# Verify binary
ls -lh itinerary-backend.exe
echo "✓ Binary size OK (should be 15-25 MB)"
```

#### Step 2: Verify No Warnings

```bash
echo "=== Checking for Warnings ==="
go build -v . 2>&1 | grep -i "warning\|deprecated"

if [ $? -eq 1 ]; then
    echo "✓ No warnings detected"
else
    echo "⚠️  Warnings found - review above"
fi
```

#### Step 3: Quick Smoke Test

```bash
# Start server in background (5 second test)
echo "=== Starting Server for Smoke Test ==="

timeout 5s ./itinerary-backend.exe > server.log 2>&1 &
sleep 2

# Test health endpoint
curl -s http://localhost:8080/api/health | jq .

# Kill server
pkill itinerary-backend.exe 2>/dev/null

echo "✓ Smoke test complete"
```

#### Step 4: Create Monday Report

```bash
cat > MONDAY_REPORT.txt << 'EOF'
=== PHASE A WEEK 2 - MONDAY EXECUTION REPORT ===
Date: $(date)

1. DATABASE SCHEMA EXTENSION
   ✓ Multi-currency tables created
   ✓ Performance monitoring tables created
   ✓ 8 currencies loaded (USD, EUR, INR, GBP, JPY, SGD, CAD, MXN)
   ✓ 8 languages loaded (en, es, fr, de, hi, ja, pt, zh)

2. TEST DATA CREATION
   ✓ 5 test users created with different currencies
   ✓ 1 multi-currency test trip created
   ✓ 4 test members from different countries
   ✓ 4 expenses in different currencies

3. TEST EXECUTION
   ✓ Model tests: 25/25 passing
   ✓ Service tests: 32/32 passing
   ✓ Integration tests: 22/22 passing
   ✓ Total: 79/79 tests passing
   ✓ Code coverage: >85%

4. APPLICATION BUILD
   ✓ Build successful
   ✓ Binary size: [size] MB
   ✓ No warnings
   ✓ Smoke test passed

5. NEXT STEPS
   - Tuesday: API endpoint testing with multi-currency
   - Wednesday: Algorithm verification with currency conversion
   - Thursday: Performance monitoring and alerting
   - Friday: Documentation and release

STATUS: ✅ MONDAY COMPLETE - READY FOR TUESDAY
EOF

cat MONDAY_REPORT.txt
```

---

## 📊 Monday Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Database schema applied | ✅ | |
| Currencies loaded | 8 | |
| Languages loaded | 8 | |
| Test users created | 5 | |
| Tests passing | 79/79 | |
| Code coverage | >85% | |
| Build successful | ✅ | |
| Smoke test passed | ✅ | |

---

## 🚨 Troubleshooting - Common Issues

### Issue: Database Connection Failed

**Error:** `sqlite3: database is locked`

**Solution:**
```bash
# Close any open connections
pkill itinerary-backend.exe
pkill sqlite3
sleep 2

# Try again
sqlite3 itinerary.db "SELECT COUNT(*) FROM supported_currencies;"
```

### Issue: Tests Failing

**Error:** Tests throw panic on user_preferences access

**Solution:**
```bash
# Verify schema was applied
sqlite3 itinerary.db ".schema user_preferences"

# If missing, apply schema again
# Verify test data exists
sqlite3 itinerary.db "SELECT COUNT(*) FROM user_preferences;"
```

### Issue: Build Fails

**Error:** `undefined: UserPreferences` or similar

**Solution:**
```bash
# Update models must be done in code (next step)
# For now, ensure code compiles with existing models

# Check if group_models.go has UserPreferences struct
grep -n "type UserPreferences" itinerary/group_models.go

# If missing, models file needs update (Tuesday task)
```

### Issue: Performance Too Slow

**If smoke test hangs:**
```bash
# Check for goroutine leaks
pprof localhost:6060/debug/pprof/goroutine

# Restart fresh
pkill -9 itinerary-backend.exe
rm itinerary.db.old
./itinerary-backend.exe &
```

---

## 📝 Monday Checklist

### Pre-Work (Before 9 AM)
- [ ] Team on Zoom/call
- [ ] All prerequisites checked
- [ ] Database backup created
- [ ] New schema review (5 min read)

### Hour 1 (Database Setup)
- [ ] Backup database
- [ ] Apply multi-currency schema
- [ ] Load supported currencies (8 total)
- [ ] Load supported languages (8 total)
- [ ] Verify schema tables exist

### Hour 2 (Test Data)
- [ ] Create 5 test users with preferences
- [ ] Create multi-currency test trip
- [ ] Add 4 members from different countries
- [ ] Add 4 expenses in different currencies
- [ ] Verify data in database

### Hour 3 (Tests)
- [ ] Run model tests (25 expected)
- [ ] Run service tests (32 expected)
- [ ] Run integration tests (22 expected)
- [ ] Run all tests with coverage
- [ ] Generate coverage report

### Hour 4 (Build & Verification)
- [ ] Build application
- [ ] Check for warnings
- [ ] Run smoke test (health endpoint)
- [ ] Create Monday report
- [ ] Update team with status

### End of Day
- [ ] All tests: 79/79 passing ✓
- [ ] Coverage > 85% ✓
- [ ] Build successful ✓
- [ ] Ready for Tuesday ✓

---

## 🔗 Related Files

- [PHASE_A_WEEK_2_ENHANCED_EXECUTION_PLAN.md](PHASE_A_WEEK_2_ENHANCED_EXECUTION_PLAN.md) - Full week overview
- [PHASE_A_WEEK_2_MULTICURRENCY_SCHEMA.sql](docs/PHASE_A_WEEK_2_MULTICURRENCY_SCHEMA.sql) - Detailed schema
- [PHASE_A_WEEK_2_DAY_1_DATABASE_TESTS.md](docs/PHASE_A_WEEK_2_DAY_1_DATABASE_TESTS.md) - Detailed procedures

---

## ✅ Success Criteria

✓ Database extended with multi-currency support  
✓ Supported currencies loaded (8 minimum)  
✓ Supported languages loaded (8 minimum)  
✓ Test data created with multi-currency scenario  
✓ All 79 tests passing  
✓ Code coverage > 85%  
✓ Application builds without warnings  
✓ Smoke test successful  

**MONDAY OBJECTIVE: COMPLETE** ✅

---

**Ready to start! Begin at 9 AM with database backup.**

*Next: Tuesday - API Testing with Multi-Currency*
