# Phase A Week 2 - Day 1: Database & Test Verification

**Date:** March 24-25, 2026 (Monday)  
**Duration:** 3-4 hours  
**Goal:** Get database running and all tests passing

---

## Task 1: Database Setup (60 minutes)

### Oracle Database Setup

```sql
-- Step 1: Connect to Oracle
sqlplus system/password@XEPDB1

-- Step 2: Execute Phase A schema
@docs/PHASE_A_GROUP_SCHEMA.sql

-- Step 3: Verify tables created
SELECT table_name FROM user_tables 
WHERE table_name IN ('GROUP_TRIPS', 'GROUP_MEMBERS', 'EXPENSES', 
                      'EXPENSE_SPLITS', 'POLLS', 'POLL_OPTIONS', 
                      'POLL_VOTES', 'SETTLEMENTS')
ORDER BY table_name;

-- Expected: 8 rows

-- Step 4: Verify indexes
SELECT index_name, table_name, column_name
FROM user_ind_columns
ORDER BY table_name, index_name;

-- Expected: 12 indexes

-- Step 5: Verify views
SELECT view_name FROM user_views
WHERE view_name LIKE '%GROUP%' OR view_name LIKE '%SETTLEMENT%';

-- Expected: 2 views
```

### PostgreSQL Database Setup

```sql
-- Step 1: Connect to PostgreSQL
psql -U postgres -d itinerary

-- Step 2: Execute Phase A schema
\i docs/PHASE_A_GROUP_SCHEMA.sql

-- Step 3: Verify tables created
\dt group_*
\dt *expense*
\dt poll*
\dt *settlement*

-- Step 4: Verify indexes
\di group_*
\di *expense*

-- Step 5: Verify views
\dv *group*
\dv *settlement*
```

### Test Data Insertion

```sql
-- Insert test data for verification
INSERT INTO users (id, name, email, password_hash, created_at)
VALUES 
  ('user-001', 'Alice', 'alice@example.com', 'hash1', NOW()),
  ('user-002', 'Bob', 'bob@example.com', 'hash2', NOW()),
  ('user-003', 'Charlie', 'charlie@example.com', 'hash3', NOW());

INSERT INTO destinations (id, name, country, created_at)
VALUES ('dest-001', 'Bali', 'Indonesia', NOW());

INSERT INTO group_trips (id, title, destination_id, owner_id, budget, duration, start_date, status, created_at, updated_at)
VALUES ('trip-001', 'Bali Adventure', 'dest-001', 'user-001', 50000, 7, '2026-05-01', 'planning', NOW(), NOW());

-- Verify data
SELECT * FROM group_trips;
```

### Success Criteria ✓

- [x] Database connection established
- [x] 8 tables created
- [x] 12 indexes created
- [x] 2 views created
- [x] Test data inserted
- [x] No errors in schema

---

## Task 2: Test Execution (90 minutes)

### 2.1 Run Model Tests

```bash
cd /d/Learn/iternary/itinerary-backend
go test ./itinerary -v -run "TestGroupTrip|TestExpense|TestPoll|TestGroupMember" -timeout 30s
```

**Expected Output:** 25 test functions passing

**Tests Covered:**
- TestGroupTripValidation
- TestExpenseValidation
- TestExpenseSplitValidation
- TestPollValidation
- TestSettlementValidation
- TestGroupMemberRoles
- TestGroupMemberStatuses
- TestExpenseCategories
- TestPollTypes
- (+17 more tests)

### 2.2 Run Service Tests

```bash
go test ./itinerary -v -run "TestCreate|TestSettlement|TestPoll|TestPermission" -timeout 30s
```

**Expected Output:** 32+ test functions passing

**Tests Covered:**
- TestCreateGroupTripValidation
- TestEqualExpenseSplit
- TestSettlementCalculation
- TestGroupMemberRolePermissions
- TestPollVoting
- TestExpenseCategoryOrganization
- TestGroupTripStatusTransitions
- TestGroupMemberInvitationLifecycle
- (+24+ more tests)

### 2.3 Run Integration Tests

```bash
go test ./itinerary -v -run "Integration|Handler|Middleware|Error" -timeout 30s
```

**Expected Output:** 22+ test functions passing

### 2.4 Run All Tests with Coverage

```bash
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s
go tool cover -html=coverage.out -o coverage.html
```

**Coverage Targets:**
- Total: > 80%
- Models package: > 85%
- Service package: > 80%
- Handlers package: > 75%

**View Coverage Report:**
- Open `coverage.html` in browser to see line-by-line coverage

---

## Task 3: Application Build (30 minutes)

### 3.1 Build Project

```bash
cd /d/Learn/iternary/itinerary-backend
go build -o itinerary-backend.exe .
```

**Expected:** Binary created with size ~15-25 MB

### 3.2 Verify Build

```bash
ls -lh itinerary-backend.exe
./itinerary-backend.exe --version  # If version flag implemented
```

### 3.3 Check for Warnings

```bash
go build -v .
```

**Expected:** All imports resolve without warnings

---

## Verification Checklist

**Database Setup:**
- [ ] Connected to Oracle/PostgreSQL
- [ ] PHASE_A_GROUP_SCHEMA.sql executed
- [ ] 8 tables created and verified
- [ ] 12 indexes created and verified
- [ ] 2 views created and verified
- [ ] Test data inserted
- [ ] Schema execution time: _____ seconds

**Test Execution:**
- [ ] go test ./itinerary - all pass
- [ ] Model tests: 25/25 passing
- [ ] Service tests: 32+/32+ passing
- [ ] Integration tests: 22+/22+ passing
- [ ] Total: 79+ tests passing
- [ ] Code coverage: _____ %
- [ ] Coverage > 80%? YES / NO

**Application Build:**
- [ ] Build succeeds with no errors
- [ ] Binary created: itinerary-backend.exe
- [ ] Binary size: _____ MB
- [ ] No compilation warnings
- [ ] Ready to start

---

## Monday Summary

**Time Spent:** _____ hours (est. 3-4 hours)

**Status:** ✓ COMPLETE / ✗ NEEDS WORK

**Issues Encountered:**
```
(List any issues and resolutions)
```

**Next Steps:**
- Proceed to Tuesday: API Endpoint Testing
- OR resolve any failing tests before continuing

---

## Appendix: Troubleshooting

### Issue: Database Connection Failed

**Symptom:** Cannot connect to Oracle/PostgreSQL

**Resolution:**
1. Verify database is running
2. Check connection string in config.json
3. Verify credentials are correct
4. Check firewall rules

### Issue: Tests Failing

**Symptom:** Some test cases fail

**Resolution:**
1. Read test error message carefully
2. Check test logs: `test_models.log`, `test_service.log`
3. Verify all dependencies installed: `go mod tidy`
4. Clear test cache: `go clean -testcache`

### Issue: Code Coverage Low

**Symptom:** Coverage < 80%

**Resolution:**
1. Identify untested functions in coverage.html
2. Add test cases for those functions
3. Run coverage again: `go test -cover ./itinerary`
4. Target: > 85% for critical paths

---

## Files Modified/Created This Day

- Database schema verified: `docs/PHASE_A_GROUP_SCHEMA.sql`
- Test logs created: `test_models.log`, `test_service.log`
- Coverage report: `coverage.out`, `coverage.html`
- Binary created: `itinerary-backend.exe`
