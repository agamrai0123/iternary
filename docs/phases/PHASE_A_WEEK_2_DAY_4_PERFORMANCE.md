# Phase A Week 2 - Day 4: Performance Baseline & Optimization

**Date:** March 28, 2026 (Thursday)  
**Duration:** 2-3 hours  
**Goal:** Establish performance baselines for all endpoints

---

## Performance Testing Strategy

### Test Approach
1. **Baseline Metrics**: Measure response times for all endpoints
2. **Load Testing**: Test with realistic data volumes
3. **Database Queries**: Identify slow queries
4. **Memory Profiling**: Check for leaks
5. **Optimization**: Record recommendations

### Tools
- **Load Testing**: Apache JMeter, Locust, or Artillery
- **Profiling**: Go pprof, CPU/memory profilers
- **Monitoring**: Gin middleware metrics

---

## Test Environment Setup

### Requirements
- **Database**: Populated with test data
  - 50 group trips
  - 500 group members (10 per trip avg)
  - 300 expenses
  - 150 polls
- **Server**: Running on localhost:8080
- **Network**: Local (no network latency)

### Prepare Test Data

```bash
# Script: prepare_load_test_data.sql

-- Create 50 test trips
INSERT INTO group_trips (id, title, destination_id, owner_id, budget, duration, start_date, status, created_at)
SELECT 
  'trip-' || rownum as id,
  'Trip ' || rownum as title,
  'dest-001' as destination_id,
  'user-' || (mod(rownum-1, 5) + 1) as owner_id,
  round(dbms_random.value(10000, 100000)) as budget,
  round(dbms_random.value(3, 14)) as duration,
  trunc(sysdate) + round(dbms_random.value(1, 180)) as start_date,
  'planning' as status,
  sysdate as created_at
FROM dual
CONNECT BY rownum <= 50;

-- Create members for each trip
INSERT INTO group_members (id, trip_id, user_id, role, status, joined_at)
SELECT 
  'member-' || rownum as id,
  'trip-' || ceil(rownum / 10) as trip_id,
  'user-' || (mod(rownum-1, 5) + 1) as user_id,
  case when rownum mod 10 = 1 then 'owner' else 'editor' end as role,
  'active' as status,
  sysdate as joined_at
FROM dual
CONNECT BY rownum <= 500;

-- Create expenses for trips
INSERT INTO expenses (id, trip_id, description, amount, paid_by_id, category, created_at)
SELECT 
  'exp-' || rownum as id,
  'trip-' || ceil(rownum / 6) as trip_id,
  'Expense ' || rownum as description,
  round(dbms_random.value(10, 5000)) as amount,
  'user-' || (mod(rownum-1, 5) + 1) as paid_by_id,
  case round(dbms_random.value(1,5))
    when 1 then 'transportation'
    when 2 then 'accommodation'
    when 3 then 'food'
    when 4 then 'activities'
    else 'other'
  end as category,
  sysdate as created_at
FROM dual
CONNECT BY rownum <= 300;

-- Create polls
INSERT INTO polls (id, trip_id, title, poll_type, created_by_id, created_at)
SELECT 
  'poll-' || rownum as id,
  'trip-' || ceil(rownum / 3) as trip_id,
  'Poll ' || rownum as title,
  'multiple_choice' as poll_type,
  'user-1' as created_by_id,
  sysdate as created_at
FROM dual
CONNECT BY rownum <= 150;

COMMIT;
```

---

## Performance Test 1: Individual Endpoint Response Times

### Test Matrix

```
Endpoint                    | Expected Time | Status Code
GET /group-trips            | < 200ms       | 200
POST /group-trips           | < 500ms       | 201
GET /group-trips/{id}       | < 100ms       | 200
PUT /group-trips/{id}       | < 200ms       | 200
DELETE /group-trips/{id}    | < 200ms       | 204
GET /group-trips/{id}/members    | < 150ms  | 200
POST /group-trips/{id}/members/invite | < 300ms | 201
POST /group-trips/{id}/expenses      | < 400ms | 201
GET /group-trips/{id}/expenses       | < 200ms | 200
GET /group-trips/{id}/report         | < 300ms | 200
POST /group-trips/{id}/polls         | < 300ms | 201
GET /group-trips/{id}/polls          | < 200ms | 200
POST /group-trips/{id}/polls/{id}/vote | < 200ms | 201
```

### Measurement Script

**Using curl with timing:**

```bash
#!/bin/bash

echo "=== Endpoint Performance Baseline ==="
echo ""

BASE_URL="http://localhost:8080/api/v1"
AUTH_TOKEN="${JWT_TOKEN}"
TRIP_ID="trip-001"

# Test 1: GET all trips
echo "[TEST 1] GET /group-trips"
curl -w "Time: %{time_total}s, HTTP: %{http_code}\n" \
  -X GET "$BASE_URL/group-trips" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -s -o /dev/null

# Test 2: GET single trip
echo "[TEST 2] GET /group-trips/{id}"
curl -w "Time: %{time_total}s, HTTP: %{http_code}\n" \
  -X GET "$BASE_URL/group-trips/$TRIP_ID" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -s -o /dev/null

# Test 3: GET trip members
echo "[TEST 3] GET /group-trips/{id}/members"
curl -w "Time: %{time_total}s, HTTP: %{http_code}\n" \
  -X GET "$BASE_URL/group-trips/$TRIP_ID/members" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -s -o /dev/null

# Test 4: GET expenses
echo "[TEST 4] GET /group-trips/{id}/expenses"
curl -w "Time: %{time_total}s, HTTP: %{http_code}\n" \
  -X GET "$BASE_URL/group-trips/$TRIP_ID/expenses" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -s -o /dev/null

# Test 5: GET report (includes settlement calc)
echo "[TEST 5] GET /group-trips/{id}/report"
curl -w "Time: %{time_total}s, HTTP: %{http_code}\n" \
  -X GET "$BASE_URL/group-trips/$TRIP_ID/report" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -s -o /dev/null

# Test 6: GET polls
echo "[TEST 6] GET /group-trips/{id}/polls"
curl -w "Time: %{time_total}s, HTTP: %{http_code}\n" \
  -X GET "$BASE_URL/group-trips/$TRIP_ID/polls" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -s -o /dev/null
```

**Expected Output:**
```
[TEST 1] GET /group-trips
Time: 0.145s, HTTP: 200
[TEST 2] GET /group-trips/{id}
Time: 0.087s, HTTP: 200
[TEST 3] GET /group-trips/{id}/members
Time: 0.125s, HTTP: 200
...
```

### Record Results

| Endpoint | Expected | Actual | Pass? | Notes |
|----------|----------|--------|-------|-------|
| GET /group-trips | <200ms | ___ms | ✓/✗ | _____ |
| POST /group-trips | <500ms | ___ms | ✓/✗ | _____ |
| GET /group-trips/{id} | <100ms | ___ms | ✓/✗ | _____ |
| PUT /group-trips/{id} | <200ms | ___ms | ✓/✗ | _____ |
| DELETE /group-trips/{id} | <200ms | ___ms | ✓/✗ | _____ |
| GET /group-trips/{id}/members | <150ms | ___ms | ✓/✗ | _____ |
| POST /group-trips/{id}/members/invite | <300ms | ___ms | ✓/✗ | _____ |
| POST /group-trips/{id}/expenses | <400ms | ___ms | ✓/✗ | _____ |
| GET /group-trips/{id}/expenses | <200ms | ___ms | ✓/✗ | _____ |
| GET /group-trips/{id}/report | <300ms | ___ms | ✓/✗ | _____ |
| POST /group-trips/{id}/polls | <300ms | ___ms | ✓/✗ | _____ |
| GET /group-trips/{id}/polls | <200ms | ___ms | ✓/✗ | _____ |
| POST .../polls/{id}/vote | <200ms | ___ms | ✓/✗ | _____ |

---

## Performance Test 2: Load Testing

### Scenario: Concurrent Requests

**Scenario A: List Operations Under Load**
```
- 50 concurrent GET /group-trips requests
- Measure: Response time, error rate
- Expected: < 500ms p95, 0% errors
```

**Using Apache ab (Apache Bench):**

```bash
# Test 1: 50 concurrent requests, 1000 total
ab -n 1000 -c 50 \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  http://localhost:8080/api/v1/group-trips

# Test 2: Repeated create operations
ab -n 100 -c 10 \
  -p create_trip.json \
  -T application/json \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  http://localhost:8080/api/v1/group-trips
```

**Expected Output:**
```
This is ApacheBench, Version 2.3
...
Completed 1000 requests
Finished 1000 requests

Benchmarking localhost (be patient)
...

Requests per second:    123.45
Time per request:       405.67 [ms]
...
Failed requests:        0

Percentage of requests served within a certain time (ms)
  50%    350
  75%    450
  90%    500
  95%    520
  99%    600
```

### Record Results

| Test | Requests | Concurrency | RPS | p95 Time | Errors |
|------|----------|-------------|-----|----------|--------|
| List trips | 1000 | 50 | _____ | ___ms | ___ |
| Create trip | 100 | 10 | _____ | ___ms | ___ |

---

## Performance Test 3: Database Query Analysis

### Identify Slow Queries

**Enable Query Logging** (Temporary for testing):

Oracle:
```sql
CREATE TABLE sql_trace AS
SELECT * FROM v$sql WHERE sql_text LIKE '%group%';

-- Or enable SQL tracing in code:
ALTER SESSION SET sql_trace=TRUE;
```

PostgreSQL:
```sql
-- Enable query logging
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = on;
SELECT pg_reload_conf();
```

### Run Queries and Check Performance

```bash
# Run load test
ab -n 100 -c 10 \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  http://localhost:8080/api/v1/group-trips/{id}/report

# Check slow query log
tail -f /var/log/postgresql/postgresql.log | grep "duration"
```

### Analyze Execution Plans

```sql
-- Oracle
EXPLAIN PLAN FOR
SELECT * FROM group_trips WHERE id = 'trip-001';
SELECT * FROM table(dbms_xplan.display);

-- PostgreSQL
EXPLAIN ANALYZE
SELECT * FROM group_trips WHERE id = 'trip-001';
```

### Expected Execution Plans

**Good Plan:**
```
Index Scan using group_trips_pkey on group_trips
  Index Cond: (id = 'trip-001')
  Total Cost: 1.2
```

**Bad Plan (Full scan):**
```
Seq Scan on group_trips
  Filter: (id = 'trip-001')
  Total Cost: 145.3
```

### Record Query Performance

| Query | Type | Rows | Time | Plan |
|-------|------|------|------|------|
| GET trip | Index | 1 | <10ms | ✓ GOOD |
| List trips | Index | 50 | <100ms | ✓ GOOD |
| List expenses | Index | 300 | <150ms | ✓ GOOD |
| Settlement calc | Complex | N/A | <300ms | Review |

---

## Performance Test 4: Memory Profiling

### Go Profiling

**Enable pprof in main.go:**

```go
import (
  _ "net/http/pprof"
)

func init() {
  go http.ListenAndServe("localhost:6060", nil)
}
```

**Collect Profile:**

```bash
# Base profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Memory allocation
go tool pprof http://localhost:6060/debug/pprof/allocs

# CPU profile (30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

**Expected Output:**
```
Showing nodes accounting for 50MB, 100% of 50MB total

flat  flat%   sum%        cum   cum%
20MB 40.0% 40.0%        30MB 60.0%  database.QueryTx
15MB 30.0% 70.0%        15MB 30.0%  encoding/json.Marshal
10MB 20.0% 90.0%        10MB 20.0%  strings.Join
5MB  10.0% 100.0%        5MB 10.0%  runtime.mallocgc
```

### Performance Targets

- **Heap Memory**: < 100MB at rest
- **Alloc Rate**: < 10MB per 1000 requests
- **GC Pause**: < 50ms

---

## Performance Test 5: Stress Testing

### Sustained Load Test

**Scenario: Sustained high load for 5 minutes**

```bash
# Using Locust (Python)
# File: locustfile.py

from locust import HttpUser, task, between

class TripsUser(HttpUser):
    wait_time = between(1, 3)
    
    @task(3)
    def get_trips(self):
        self.client.get(
            "/api/v1/group-trips",
            headers={"Authorization": "Bearer {token}"}
        )
    
    @task(1)
    def get_trip_detail(self):
        self.client.get(
            "/api/v1/group-trips/trip-001",
            headers={"Authorization": "Bearer {token}"}
        )
    
    @task(2)
    def get_trip_report(self):
        self.client.get(
            "/api/v1/group-trips/trip-001/report",
            headers={"Authorization": "Bearer {token}"}
        )

# Run:
# locust -f locustfile.py -u 100 -r 10 --host=http://localhost:8080
```

### Monitor During Stress Test

```bash
# Terminal 1: Run load test
locust -f locustfile.py -u 100 -r 10 --host=http://localhost:8080

# Terminal 2: Monitor server
# Check CPU, memory, active connections
watch -n 1 'ps aux | grep itinerary'
netstat -an | grep 8080 | wc -l

# Terminal 3: Monitor database
# Oracle:
SELECT count(*) FROM v$session WHERE status='ACTIVE';

# PostgreSQL:
SELECT count(*) FROM pg_stat_activity;
```

### Expected Results

- **Success rate**: 99-100%
- **Response time p95**: < 1000ms
- **Server CPU**: < 80%
- **Server Memory**: < 200MB
- **Database connections**: < 50

---

## Optimization Recommendations

### If Performance Is Poor

#### Slow List Endpoints
```
Issue: GET /group-trips takes > 200ms
Solution:
1. Add pagination if not present
2. Create index on status column
3. Cache results for 5 minutes
4. Implement lazy loading of relationships
```

#### Slow Settlement Calculation
```
Issue: GET /group-trips/{id}/report takes > 500ms
Solution:
1. Cache settlement results
2. Recalculate only when expenses change
3. Use database-side settlement view
4. Limit historical data loaded
```

#### High Memory Usage
```
Issue: Memory grows > 100MB
Solution:
1. Reduce connection pool size
2. Clear request caches
3. Use streaming for large result sets
4. Check for goroutine leaks
```

### Optimization Checklist

- [ ] All endpoints under time targets
- [ ] under 100MB memory at rest
- [ ] Database indexes properly used
- [ ] No full scans on large tables
- [ ] Connection pooling configured
- [ ] Caching implemented where appropriate
- [ ] Query results paginated
- [ ] Large results streamed

---

## Performance Baseline Report

**Test Date:** _______________  
**Tested By:** _______________

### Summary

**Endpoints Performance:**
- Endpoints meeting targets: ___/13
- Average response time: ___ms
- Slowest endpoint: __________ (___ms)
- Fastest endpoint: __________ (___ms)

**Load Testing:**
- Max RPS achieved: _____
- Error rate under load: ____%
- p95 response time: ___ms

**Memory:**
- Heap size at rest: ___MB
- Peak heap size: ___MB
- No memory leaks: YES / NO

**Database:**
- All queries using indexes: YES / NO
- Slowest query: __________ (___ms)
- Connection pool: ___/___

### Bottlenecks Identified

1. ____________________
2. ____________________
3. ____________________

### Optimizations Recommended

1. ____________________
2. ____________________

### Action Items

- [ ] Fix slow endpoints (if any)
- [ ] Optimize database queries
- [ ] Review memory usage
- [ ] Document findings

---

## Thursday Summary

- Time spent: _____ hours
- Baseline established: YES / NO
- Performance acceptable: YES / NO
- Optimizations needed: YES / NO
- Ready for Day 5: YES / NO

---

## Appendix: Performance Tools

### Command Reference

```bash
# curl timing
curl -w "Response time: %{time_total}s\n" -O /dev/null http://localhost:8080/api/v1/group-trips

# Apache Bench
ab -n 1000 -c 50 -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/group-trips

# Go pprof
go test -cpuprofile=cpu.prof -memprofile=mem.prof -run TestXyz

# Analyze profile
go tool pprof -http=:8081 cpu.prof
```

### References

- Go Performance: https://golang.org/doc/diagnostics
- PostgreSQL Optimization: https://www.postgresql.org/docs/current/runtime-config.html
- Oracle Performance: https://docs.oracle.com/
