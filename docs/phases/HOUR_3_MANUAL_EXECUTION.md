# 🎯 Hour 3 Test Execution - Manual Step-by-Step Guide

**Status:** Ready to Execute  
**Time:** 11:00 AM - 12:00 PM  
**Goal:** Run 79 tests across 12 test files, achieve >85% code coverage

---

## ⚡ Quick Start (Copy-Paste Commands)

### Option 1: Windows Command Prompt (CMD)

**Open Command Prompt and navigate:**
```cmd
cd D:\Learn\iternary\itinerary-backend
```

**Create/verify database:**
```cmd
sqlite3 itinerary.db < multicurrency_schema.sql
```

**Run all tests:**
```cmd
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s
```

**Generate coverage report:**
```cmd
go tool cover -html=coverage.out -o coverage.html
```

**View coverage summary:**
```cmd
go tool cover -func=coverage.out | tail -5
```

---

### Option 2: Git Bash (if you prefer)

```bash
cd /d/Learn/iternary/itinerary-backend
sqlite3 itinerary.db < multicurrency_schema.sql
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s 2>&1 | tee test_results.txt
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | tail -5
```

---

## 📋 What Tests Will Run (79 Total)

| Test File | Test Count | Categories |
|-----------|-----------|-----------|
| auth_service_test.go | ~4 | Token generation, validation, password hashing |
| config_test.go | ~5 | Config loading, properties |
| error_test.go | ~3 | Error handling, codes |
| logger_test.go | ~7 | Debug, info, error, warn, fields |
| metrics_test.go | ~4 | Metrics collection, reporting |
| models_test.go | ~3 | Model validation |
| service_test.go | ~5 | Service business logic |
| template_helpers_test.go | ~3 | Template rendering |
| group_models_test.go | **~25** | Group entities, validation, status |
| group_service_test.go | **~32** | Settlement, permissions, algorithms |
| group_integration_test.go | **~22** | End-to-end, middleware, workflows |
| **TOTAL** | **~79** | **Comprehensive coverage** |

---

## ✅ Expected Results

### When Tests Complete Successfully:

```
ok	itinerary	15.234s	coverage: 85.2% of statements
```

### What You Should See:

1. **Tests running output:**
   ```
   === RUN   TestConfigLoading
   --- PASS: TestConfigLoading (0.00s)
   === RUN   TestConfigProperties
   --- PASS: TestConfigProperties (0.01s)
   ...
   === RUN   TestGroupServiceSettlement
   --- PASS: TestGroupServiceSettlement (0.05s)
   ```

2. **Final summary line (at the very end):**
   ```
   ok	itinerary	Xs	coverage: X.X%
   ```

3. **Three files created:**
   - `coverage.out` (50-100KB, raw coverage data)
   - `coverage.html` (200-500KB, visual report - open in browser)
   - `test_summary.txt` (if using the batch script)

---

## 🔍 Key Milestones During Execution

### Phase 1: Test Discovery (10-15 seconds)
- Go compiler discovers all `_test.go` files
- Parses test functions (those starting with `Test`)
- Initializes test environment

### Phase 2: Setup Tests (1-2 seconds)
- Config tests: Load configuration
- Logger tests: Initialize logger
- Database isolation (if needed)

### Phase 3: Core Tests (5-8 seconds)
- Auth tests: ~0.5s
- Model tests: ~1s
- Service tests: ~2s
- Group tests: ~3s

### Phase 4: Integration Tests (3-5 seconds)
- Group integration: ~3s
- Handler tests: ~1s
- Full workflow tests: ~1s

### Phase 5: Coverage Analysis (1-2 seconds)
- Go measures code coverage
- Generates coverage.out file
- Reports: `coverage: 85+%`

**Total Expected Duration: 10-15 minutes**

---

## 📊 Coverage Target Breakdown

### Expected Coverage by Component:

| Component | Target | Type |
|-----------|--------|------|
| Auth service | 85%+ | Authentication, tokens, passwords |
| Configuration | 90%+ | Config loading and validation |
| Database | 80%+ | CRUD operations, queries |
| Models | 95%+ | Struct validation and methods |
| Service layer | 85%+ | Business logic |
| Group features | 80%+ | Complex settlement logic |
| Handlers | 70%+ | HTTP endpoint logic |
| Middleware | 75%+ | Request/response processing |
| **Overall** | **>85%** | **All code** |

---

## 🚀 Step-by-Step Execution

### Step 1: Open Terminal

**Windows:** Press `Win+R`, type `cmd`, press Enter
**Or:** Open Git Bash, PowerShell, or Terminal

### Step 2: Navigate to Backend

```cmd
cd D:\Learn\iternary\itinerary-backend
```

Expected: You should see the prompt: `D:\Learn\iternary\itinerary-backend>`

### Step 3: Verify Files Exist

```cmd
dir itinerary
```

Should show:
- ✅ auth_service_test.go
- ✅ group_service_test.go
- ✅ group_integration_test.go
- ✅ (and other test files)

### Step 4: Create/Update Database

```cmd
sqlite3 itinerary.db < multicurrency_schema.sql
```

Expected output:
- No errors
- File `itinerary.db` created or updated

Verify:
```cmd
sqlite3 itinerary.db "SELECT COUNT(*) FROM sqlite_master WHERE type='table';"
```

Should show: `25` (or similar - number of tables)

### Step 5: Run Tests (THE MAIN EVENT)

```cmd
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s
```

**What to watch for:**
- ✅ Each test shows as `--- PASS:` (should be green/successful)
- ⏱️ Tests run in ~10-15 seconds
- ✅ Final line shows: `ok itinerary ... coverage: XX.X%`
- ❌ NO `--- FAIL:` lines (if any tests fail, they'll show here)

### Step 6: Generate HTML Coverage Report

```cmd
go tool cover -html=coverage.out -o coverage.html
```

Expected:
- File `coverage.html` created (~300KB)

### Step 7: Check Coverage

```cmd
go tool cover -func=coverage.out | tail -5
```

Expected output:
```
total:                           (statements)  XX.X%
```

Where XX.X should be **>85%**

---

## 🎯 Success Criteria

### For Hour 3 to be COMPLETE:

- ✅ All 79 tests PASS (0 failures)
- ✅ Coverage >= 85%
- ✅ Files created: coverage.out, coverage.html
- ✅ No errors in output
- ✅ Duration approximately 10-15 minutes

### Example of SUCCESS:

```
ok	itinerary	14.523s	coverage: 85.2% of statements
```

### Example of PARTIAL SUCCESS:

```
ok	itinerary	14.523s	coverage: 82.1% of statements
```
(Coverage slightly below 85% - still acceptable)

### Example of FAILURE:

```
--- FAIL: TestGroupServiceSettlement (0.05s)
FAIL	itinerary	15.234s
```

---

## 🐛 Troubleshooting

### Problem: "command not found: go"
**Solution:**
- Ensure Go is installed: `go version`
- If not found, [install Go](https://golang.org/dl)
- Restart terminal after installation

### Problem: "sqlite3: command not found"
**Solution:**
- Install SQLite3: `choco install sqlite` (Windows with Chocolatey)
- Or download from [sqlite.org](https://www.sqlite.org/download.html)

### Problem: "database is locked" error
**Solution:**
```cmd
taskkill /IM itinerary-backend.exe /F
```
Then retry tests

### Problem: Tests fail with "connection refused"
**Solution:**
- Database might not exist or be initialized
- Run the schema application command first:
```cmd
sqlite3 itinerary.db < multicurrency_schema.sql
```

### Problem: Coverage is less than 85%
**Solution:**
- This might be normal - submit what you have
- Coverage of 80-85% is still very good
- Don't worry about reaching exactly 85%

### Problem: Takes longer than 15 minutes
**Solution:**
- Increase timeout (use `-timeout 120s` instead of `-timeout 60s`)
- Some machines run slower - this is fine

---

## 📈 After Tests Complete

### What to do next:

1. **Save the coverage report:**
   - Open `coverage.html` in your browser
   - Take a screenshot showing the coverage percentage

2. **Collect test results:**
   - Note the final line: `ok itinerary X.XXXs coverage: YY.Y%`
   - If there are any failures, note them

3. **Proceed to Hour 4:**
   - When tests pass, move to: [HOUR_4_VERIFICATION_GUIDE.md](HOUR_4_VERIFICATION_GUIDE.md)
   - Hour 4: Build binary and verify server startup

4. **Report results:**
   - Tests passed: YES/NO
   - Coverage percent: XX.X%
   - Any failures: (list if applicable)

---

## 💡 Understanding Test Output

### What each symbol means:

| Symbol | Meaning | Status |
|--------|---------|--------|
| `RUN` | Test started | 🟡 Running |
| `PASS` | Test succeeded | ✅ Good |
| `FAIL` | Test failed | ❌ Problem |
| `SKIP` | Test skipped | ⏭️ Ignored |

### Example output:
```
=== RUN   TestConfigLoading
--- PASS: TestConfigLoading (0.00s)
=== RUN   TestConfigProperties  
--- PASS: TestConfigProperties (0.01s)
=== RUN   TestGroupServiceSettlement
--- PASS: TestGroupServiceSettlement (0.05s)
ok  itinerary    14.523s coverage: 85.2% of statements
```

This means:
- 3 tests ran ✅
- All passed ✅
- Coverage is 85.2% ✅
- Total time: 14.5 seconds ✅

---

## 🎊 You're Ready!

Everything is prepared. Execute the commands from **Quick Start** section above to:

1. ✅ Set up database
2. ✅ Run 79 tests  
3. ✅ Generate coverage report
4. ✅ Complete Hour 3

**Expected time:** ~15 minutes  
**Expected result:** All tests passing, 85%+ coverage

---

## 📞 Quick Reference

| Need | Command |
|------|---------|
| Run all tests | `go test ./itinerary -v -cover -coverprofile=coverage.out` |
| See coverage | `go tool cover -func=coverage.out` |
| Open HTML report | `coverage.html` (double-click) |
| Run specific test | `go test ./itinerary -run TestConfigLoading -v` |
| See test output | `cat test_results.txt` (or open in editor) |
| Clean test cache | `go clean -testcache` |

---

## ✨ When You're Done

After successful execution:

1. ✅ Move to Hour 4: Build & Verification
2. ✅ Build binary: `go build -o itinerary-backend.exe .`
3. ✅ Start server: `./itinerary-backend.exe`
4. ✅ Verify endpoints working
5. ✅ Complete Monday sprint by 1 PM

**Let's go! Execute Hour 3 now! 🚀**
