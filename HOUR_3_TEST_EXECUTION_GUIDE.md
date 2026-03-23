# Monday Hour 3: Test Execution - Complete Guide

**Time:** 11:00 AM - 12:00 PM  
**Goal:** Run 79 tests with >85% coverage  
**Status:** Ready to execute

---

## Test Files Available

Located in `d:\Learn\iternary\itinerary-backend\itinerary\`:

1. **auth_service_test.go** - Authentication tests
2. **config_test.go** - Configuration tests  
3. **error_test.go** - Error handling tests
4. **logger_test.go** - Logger tests
5. **metrics_test.go** - Metrics tests
6. **models_test.go** - Models tests (25+ tests)
7. **service_test.go** - Service tests
8. **template_helpers_test.go** - Template tests
9. **auth_service_test.go** - Auth service tests
10. **group_models_test.go** - Group models tests (25+ tests)
11. **group_service_test.go** - Group service tests (32+ tests)
12. **group_integration_test.go** - Integration tests (22+ tests)

**Total Tests:** 79 tests across 12 files

---

## Step-by-Step Execution

### Step 1: Navigate to Project Directory

```bash
cd /d/Learn/iternary/itinerary-backend
```

### Step 2: Run All Tests with Coverage

```bash
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s
```

**What this does:**
- `-v` = Verbose output (shows each test)
- `-cover` = Calculate coverage percentage
- `-coverprofile=coverage.out` = Save coverage data to file
- `-timeout 60s` = Allow 60 seconds for all tests to complete

**Expected Output:**
```
=== RUN   TestConfigLoading
--- PASS: TestConfigLoading (0.00s)
=== RUN   TestConfigProperties
--- PASS: TestConfigProperties (0.00s)
...
[more tests...]
...
ok  	github.com/yourusername/itinerary-backend/itinerary	5.234s	coverage: 85.2% of statements
```

**Expected Result:** 
- ✅ All 79 tests PASS
- ✅ Coverage: ≥85%

---

### Step 3: Generate HTML Coverage Report

After tests pass, generate an HTML report:

```bash
go tool cover -html=coverage.out -o coverage.html
```

This creates a visual representation showing:
- Green: Code that's tested
- Red: Code that's not tested
- Line-by-line coverage details

**Expected:** `coverage.html` file created (~5-10KB)

---

### Step 4: View Coverage Summary

Get a text summary of coverage:

```bash
go tool cover -func=coverage.out
```

**Expected Output:**
```
github.com/yourusername/itinerary-backend/itinerary/models.go:XX:	FunctionName	100.0%
github.com/yourusername/itinerary-backend/itinerary/service.go:XX:	FunctionName	85.0%
...
total:	(statements)	85.2%
```

---

### Step 5: Verify Test Results

Check test summary:

```bash
# Count passing tests
grep "PASS" test_results.log | wc -l
# Expected: 79

# Check for failures
grep "FAIL" test_results.log
# Expected: (no output = no failures)

# Get coverage percentage
go tool cover -func=coverage.out | grep total
# Expected: total: (statements) 85.0% or higher
```

---

## Individual Test Suites (Optional)

If you want to run tests by module:

### Model Tests Only
```bash
go test ./itinerary -v -run "TestGroup|TestExpense|TestPoll|TestSettlement" -timeout 30s
# Expected: 25 tests passing
```

### Service Tests Only
```bash
go test ./itinerary -v -run "TestCreate|TestSettlement|TestPermission|TestVoting" -timeout 30s
# Expected: 32+ tests passing
```

### Integration Tests Only
```bash
go test ./itinerary -v -run "Integration|Handler|Middleware" -timeout 30s
# Expected: 22+ tests passing
```

### Configuration Tests Only
```bash
go test ./itinerary -v -run "TestConfig" -timeout 10s
# Expected: All config tests passing
```

---

## Expected Test Coverage Breakdown

| Module | Coverage | Status |
|--------|----------|--------|
| Models | 85%+ | ✓ Expected |
| Service | 80%+ | ✓ Expected |
| Handlers | 75%+ | ✓ Expected |
| Database | 85%+ | ✓ Expected |
| Auth | 90%+ | ✓ Expected |
| Logger | 95%+ | ✓ Expected |
| **Total** | **85%+** | **✓ Target** |

---

## Quick Commands Cheat Sheet

```bash
# Navigate to project
cd /d/Learn/iternary/itinerary-backend

# Run all tests (simple)
go test ./itinerary -v

# Run all tests with coverage (recommended)
go test ./itinerary -v -cover -coverprofile=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Show coverage summary
go tool cover -func=coverage.out | tail -5

# Run tests matching pattern
go test ./itinerary -v -run "pattern"

# Run tests with timeout
go test ./itinerary -v -timeout 60s

# Run specific test file
go test ./itinerary -v -run TestConfigLoading

# Clean test cache (if needed)
go clean -testcache && go test ./itinerary -v
```

---

## Troubleshooting

### Issue: Tests fail with "connection refused"
**Cause:** Server not running or database not initialized

**Solution:**
```bash
# Make sure database exists
cd /d/Learn/iternary/itinerary-backend
sqlite3 itinerary.db ".tables"

# If no tables, apply schema first
sqlite3 itinerary.db < multicurrency_schema.sql
```

### Issue: "Cannot find package"
**Cause:** Go modules not initialized

**Solution:**
```bash
cd /d/Learn/iternary/itinerary-backend
go mod tidy
go mod download
go test ./itinerary -v
```

### Issue: Tests timeout
**Cause:** Tests taking > 60 seconds

**Solution:**
```bash
# Increase timeout
go test ./itinerary -v -timeout 120s
```

### Issue: Coverage file not created
**Cause:** Tests failed

**Solution:**
```bash
# Run tests verbosely to see errors
go test ./itinerary -v 2>&1 | tail -50
```

---

## Success Checklist

- [ ] Navigate to itinerary-backend directory
- [ ] Run: `go test ./itinerary -v -cover -coverprofile=coverage.out`
- [ ] All 79 tests show **PASS**
- [ ] Coverage shows **85%+**
- [ ] Run: `go tool cover -html=coverage.out -o coverage.html`
- [ ] coverage.html file created
- [ ] Review coverage report (optional)

---

## What Happens Next

✅ **If all tests pass:**
1. Coverage report generated
2. Move to Hour 4: Build and verification
3. Build binary
4. Start server
5. Run smoke tests

❌ **If tests fail:**
1. Check error messages in output
2. Use `-run "TestName" -v` to run specific failing test
3. Read test file to understand what it's testing
4. Fix the issue
5. Re-run tests

---

## Preview: What Tests Check

### Model Tests (25+ tests)
- ✓ GroupTrip validation (title, budget, duration)
- ✓ Expense validation (amounts, splits)
- ✓ Poll validation (questions, options)
- ✓ Settlement calculations
- ✓ GroupMember roles and permissions
- ✓ Status transitions

### Service Tests (32+ tests)
- ✓ Create operations with validation
- ✓ Settlement algorithm accuracy
- ✓ Permission checks
- ✓ Poll voting logic
- ✓ Expense splitting
- ✓ Invitation lifecycles

### Integration Tests (22+ tests)
- ✓ End-to-end workflows
- ✓ Error handling
- ✓ Middleware behavior
- ✓ HTTP status codes
- ✓ Data consistency

---

**Ready to execute Hour 3?** Follow the steps above and report back when tests complete!

Next: Hour 4 Build and Verification
