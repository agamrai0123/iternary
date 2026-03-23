# Phase A Week 2 - Day 3: Settlement Algorithm Verification

**Date:** March 27, 2026 (Wednesday)  
**Duration:** 2-3 hours  
**Goal:** Verify settlement calculations and algorithms work correctly

---

## Overview

We'll verify three key algorithms:
1. **Expense Splitting**: Equal and custom splits
2. **Settlement Calculation**: Minimized transaction algorithm
3. **Poll Voting**: Duplicate prevention and vote counting

---

## Algorithm 1: Equal Expense Splitting

### Test Case 1.1: Simple Equal Split (3 people)

**Setup:**
- Trip: "Bali Adventure"
- Members: Alice (paid), Bob, Charlie
- Expense: $300 for flight (paid by Alice)

**Expected Split:**
```
Alice:   -$100 (paid $300, owes $100)
Bob:     +$100 (owes $100)
Charlie: +$100 (owes $100)
```

**Test:**
```bash
POST /api/v1/group-trips/{trip_id}/expenses
{
  "description": "Flight tickets",
  "amount": 300,
  "paid_by_id": "alice",
  "category": "transportation",
  "split_type": "equal",
  "split_among": ["alice", "bob", "charlie"]
}
```

**Verification:**
- [ ] Response has 3 splits
- [ ] Each split = $100
- [ ] Alice marked as "paid": true
- [ ] Bob, Charlie marked as "paid": false

**Result:** ✓ PASS / ✗ FAIL

---

### Test Case 1.2: Uneven Split (4 people, odd amount)

**Setup:**
- Expense: $100 for dinner split among 4 people
- Amount per person = $25

**SQL Verification:**
```sql
SELECT * FROM expense_splits 
WHERE expense_id = 'exp-{id}'
ORDER BY user_id;
-- Expected: 4 rows, each with amount = 25.00
```

**Verification:**
- [ ] No rounding errors
- [ ] Total = $100
- [ ] 4 splits of $25 each

**Result:** ✓ PASS / ✗ FAIL

---

## Algorithm 2: Custom Expense Splitting

### Test Case 2.1: Custom Allocation

**Setup:**
- Expense: $1000 for hotel (paid by Alice)
- Custom split requested:
  - Alice: $300
  - Bob: $400
  - Charlie: $300

**Test:**
```bash
POST /api/v1/group-trips/{trip_id}/expenses
{
  "description": "Hotel",
  "amount": 1000,
  "paid_by_id": "alice",
  "category": "accommodation",
  "split_type": "custom",
  "split_details": [
    {"user_id": "alice", "amount": 300},
    {"user_id": "bob", "amount": 400},
    {"user_id": "charlie", "amount": 300}
  ]
}
```

**Verification:**
- [ ] Response shows all 3 splits
- [ ] Alice owes: $300
- [ ] Bob owes: $400
- [ ] Charlie owes: $300
- [ ] Total sum = $1000

**Result:** ✓ PASS / ✗ FAIL

---

### Test Case 2.2: Reject Invalid Custom Split

**Setup:**
- Total expense: $1000
- Custom split amounts: $400 + $300 + $100 = $800 (NOT $1000)

**Test:**
```bash
POST /api/v1/group-trips/{trip_id}/expenses
{
  "description": "Hotel",
  "amount": 1000,
  "paid_by_id": "alice",
  "split_type": "custom",
  "split_details": [
    {"user_id": "alice", "amount": 400},
    {"user_id": "bob", "amount": 300},
    {"user_id": "charlie", "amount": 100}
  ]
}
```

**Expected Response:** HTTP 400 Bad Request
```json
{
  "error": "Custom split amounts must sum to total expense amount",
  "expected": 1000,
  "received": 800
}
```

**Verification:**
- [ ] Returns 400 (not 201)
- [ ] Error message explains mismatch
- [ ] Expense NOT created

**Result:** ✓ PASS / ✗ FAIL

---

## Algorithm 3: Settlement Calculation

This is the MOST IMPORTANT algorithm. It minimizes the number of transactions needed to settle debts.

### Background

**Problem:**
```
Alice: receives $1500 from Bob and Charlie
Bob:   receives $500 from Alice
Charlie: receives $1000 from Alice

Simple approach: 3 transactions
- Bob pays Alice $500
- Charlie pays Alice $1000

Settlement after expense accounting:
- Alice is owed $1500 total
- Bob owes Alice $500
- Charlie owes Alice $1000
- Solution: 2 transactions (minimal)
```

### Test Case 3.1: Simple Settlement (3 People, Linear)

**Setup:**
```
Expense 1: Alice pays $300 (split equally: A=$100 owed, B=$100 owed, C=$100 owed)
Expense 2: Bob pays $200 (split: A=$50 owed, B=$50, C=$100 owed)
```

**Calculate Balances:**
```
Alice: Paid $300, Owes $150 → Net: +$150 (creditor)
Bob:   Paid $200, Owes $150 → Net: +$50 (creditor)
Charlie: Paid $0, Owes $200 → Net: -$200 (debtor)
```

**Minimal Settlement:**
```
1. Charlie pays Alice $150
2. Charlie pays Bob $50
Total: 2 transactions (minimal)
```

**Test:**
```bash
GET /api/v1/group-trips/{trip_id}/report
```

**Verification:**
- [ ] Settlement shows 2 transactions
- [ ] Transactions zero out all balances
- [ ] No circular payments

**Result:** ✓ PASS / ✗ FAIL

---

### Test Case 3.2: Complex Settlement (5 People, Multiple Transactions)

**Setup:** Create multiple expenses with different payers
```
Expense 1: Alice pays $500 (split among 5)
Expense 2: Bob pays $600 (split among 5)
Expense 3: Charlie pays $400 (split among 5)
Expense 4: Dave pays $300 (split among 5)
```

**Calculate Shares:**
```
Each person's share:
- Share 1: $500/5 = $100 per person
- Share 2: $600/5 = $120 per person
- Share 3: $400/5 = $80 per person
- Share 4: $300/5 = $60 per person
Total per person: $360

Individual Balances:
- Alice: Paid $500, Owes $360 → Net: +$140
- Bob:   Paid $600, Owes $360 → Net: +$240
- Charlie: Paid $400, Owes $360 → Net: +$40
- Dave:   Paid $300, Owes $360 → Net: -$60
- Eve:    Paid $0, Owes $360 → Net: -$360
```

**Minimal Settlement (should be 3-4 transactions max):**
```
Algorithm result should show minimal transactions
Example (may vary):
1. Eve pays Alice $140
2. Eve pays Bob $160
3. Eve pays Charlie $40
4. Dave pays Bob $60

Wait, that's suboptimal. Correct settlement:
1. Eve pays Alice $140
2. Eve pays Bob $200
3. Eve pays Charlie $20
4. Dave pays Bob $40
5. Charlie pays Dave $20 (no - wrong direction)

Actually, we need to verify the algorithm produces minimal.
```

**Test:**
```bash
GET /api/v1/group-trips/{trip_id}/report
```

**Query Database:**
```sql
SELECT * FROM settlements 
WHERE trip_id = '{trip_id}'
ORDER BY settlement_id;

-- Count transactions
SELECT COUNT(*) as transaction_count 
FROM settlements 
WHERE trip_id = '{trip_id}';
-- Expected: 3-4 transactions maximum (not 5-6)
```

**Verification:**
- [ ] Settlement calculated without errors
- [ ] Transaction count ≤ 4 (minimal)
- [ ] All balances zero out

**Result** ✓ PASS / ✗ FAIL

---

### Test Case 3.3: Settlement Algorithm - Edge Case (Multiple Creditors/Debtors)

**Setup:** Create scenario where multiple people have credits/debits
```
Alice: Net +$200 (creditor)
Bob:   Net +$100 (creditor)
Charlie: Net -$150 (debtor)
Dave:   Net -$150 (debtor)
```

**Expected Settlement:**
```
1. Charlie pays Alice $150
2. Dave pays Alice $50
3. Dave pays Bob $100

Total: 3 transactions (optimal)
```

**Verification:**
- [ ] Algorithm finds ≤3 transactions
- [ ] No excess transactions
- [ ] All balances settle to zero

**Debug Settlement Table:**
```sql
SELECT * FROM settlements 
WHERE trip_id = '{trip_id}'
ORDER BY created_at;
```

**Result:** ✓ PASS / ✗ FAIL

---

## Algorithm 4: Poll Voting

### Test Case 4.1: Single Vote Per User (No Duplicates)

**Setup:**
- Poll: "Best restaurant?"
- Options: Italian, Thai, Mexican, Japanese
- Each user votes once

**Test Sequence:**
```
1. Alice votes for "Thai"
   Expected: vote_count for Thai = 1

2. Bob votes for "Italian"
   Expected: vote_count for Italian = 1

3. Charlie votes for "Thai"
   Expected: vote_count for Thai = 2
```

**Verification:**
```bash
GET /api/v1/group-trips/{trip_id}/polls/{poll_id}
```

**Response Should Show:**
```json
{
  "options": [
    {"id": "opt-1", "text": "Italian", "vote_count": 1},
    {"id": "opt-2", "text": "Thai", "vote_count": 2},
    {"id": "opt-3", "text": "Mexican", "vote_count": 0},
    {"id": "opt-4", "text": "Japanese", "vote_count": 0}
  ]
}
```

**Verification:**
- [ ] Vote counts are correct
- [ ] Thai has 2 votes (not 1)
- [ ] Italian has 1 vote
- [ ] Total votes = 3

**Result:** ✓ PASS / ✗ FAIL

---

### Test Case 4.2: Prevent Duplicate Votes

**Setup:**
- Alice already voted for "Thai"

**Test:** Alice tries to vote again for "Italian"
```bash
POST /api/v1/group-trips/{trip_id}/polls/{poll_id}/vote
{
  "option_id": "opt-1",
  "voted_by_id": "alice"
}
```

**Expected Response:** HTTP 409 Conflict
```json
{
  "error": "You have already voted on this poll",
  "previous_vote_option": "Thai"
}
```

**Verification:**
- [ ] Returns 409 (not 201)
- [ ] Vote count unchanged (still Thai=2, Italian=1)
- [ ] Error message explains conflict

**Result:** ✓ PASS / ✗ FAIL

---

### Test Case 4.3: Multiple Choice Voting (If Enabled)

**Setup:**
- Poll: "Which days work for trip?" (assume allow_multiple=true)
- Options: Friday, Saturday, Sunday

**Test:**
```
1. Alice votes for "Friday" and "Saturday"
   Expected: Both options increment by 1

2. Bob votes for "Saturday"
   Expected: Saturday count further increments
```

**Verification:**
- [ ] Alice can vote multiple times? (depends on design)
- [ ] Vote counts correct

**Result:** ✓ PASS / ✗ FAIL

---

## Comprehensive Algorithm Test Script

Create a shell script to test all algorithms:

**File: `test_algorithms.sh`**

```bash
#!/bin/bash

echo "=== Phase A Week 2 - Algorithm Verification ==="
echo ""

# Configuration
BASE_URL="http://localhost:8080/api/v1"
AUTH_TOKEN="${JWT_TOKEN}"
TRIP_ID=""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test 1: Equal Split
echo -e "${YELLOW}[TEST 1] Equal Expense Splitting${NC}"
curl -X POST "$BASE_URL/group-trips/$TRIP_ID/expenses" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d @- << 'EOF'
{
  "description": "Test: Equal Split $300",
  "amount": 300,
  "paid_by_id": "user-001",
  "category": "food",
  "split_type": "equal",
  "split_among": ["user-001", "user-002", "user-003"]
}
EOF
echo ""

# Test 2: Custom Split - Valid
echo -e "${YELLOW}[TEST 2] Custom Expense Splitting (Valid)${NC}"
curl -X POST "$BASE_URL/group-trips/$TRIP_ID/expenses" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d @- << 'EOF'
{
  "description": "Test: Custom Split $1000",
  "amount": 1000,
  "paid_by_id": "user-001",
  "category": "accommodation",
  "split_type": "custom",
  "split_details": [
    {"user_id": "user-001", "amount": 300},
    {"user_id": "user-002", "amount": 400},
    {"user_id": "user-003", "amount": 300}
  ]
}
EOF
echo ""

# Test 3: Custom Split - Invalid (amounts don't match)
echo -e "${YELLOW}[TEST 3] Custom Expense Splitting (Invalid - should fail with 400)${NC}"
curl -X POST "$BASE_URL/group-trips/$TRIP_ID/expenses" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d @- << 'EOF'
{
  "description": "Test: Invalid Custom Split",
  "amount": 1000,
  "paid_by_id": "user-001",
  "category": "accommodation",
  "split_type": "custom",
  "split_details": [
    {"user_id": "user-001", "amount": 300},
    {"user_id": "user-002", "amount": 300},
    {"user_id": "user-003", "amount": 200}
  ]
}
EOF
echo ""

# Test 4: Settlement Report
echo -e "${YELLOW}[TEST 4] Settlement Calculation${NC}"
curl -X GET "$BASE_URL/group-trips/$TRIP_ID/report" \
  -H "Authorization: Bearer $AUTH_TOKEN"
echo ""

# Test 5: Poll Voting
echo -e "${YELLOW}[TEST 5] Poll Voting${NC}"
POLL_ID="poll-001"
curl -X POST "$BASE_URL/group-trips/$TRIP_ID/polls/$POLL_ID/vote" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"option_id": "opt-001", "voted_by_id": "user-001"}'
echo ""

# Test 6: Prevent Duplicate Vote
echo -e "${YELLOW}[TEST 6] Prevent Duplicate Vote (should fail with 409)${NC}"
curl -X POST "$BASE_URL/group-trips/$TRIP_ID/polls/$POLL_ID/vote" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"option_id": "opt-002", "voted_by_id": "user-001"}'
echo ""

echo -e "${GREEN}=== Algorithm Tests Complete ===${NC}"
```

---

## Verification Checklist

**Equal Expense Splitting:**
- [ ] 3-person split: Each gets $100 from $300
- [ ] 4-person split: Each gets $25 from $100 (odd amounts handled)
- [ ] Floating point precision: No rounding errors

**Custom Expense Splitting:**
- [ ] Valid splits accepted: $400 + $300 + $300 = $1000
- [ ] Invalid splits rejected: Returns HTTP 400
- [ ] Error message clear and helpful

**Settlement Algorithm:**
- [ ] Test 1 (3 people): Max 2 transactions
- [ ] Test 2 (5 people): Max 4 transactions
- [ ] Test 3 (multiple creditors/debtors): Calculation correct
- [ ] No circular payments
- [ ] All balances zero out

**Poll Voting:**
- [ ] Vote count increments correctly
- [ ] Duplicate votes prevented (409)
- [ ] Vote marked with correct user

---

## SQL Verification Queries

### Verify Expense Splits

```sql
-- Check splits for expense
SELECT 
  es.expense_id,
  es.user_id,
  es.amount,
  es.paid
FROM expense_splits es
WHERE es.expense_id = 'exp-{id}'
ORDER BY es.user_id;
```

### Verify Settlements

```sql
-- Check settlements for trip
SELECT 
  s.id,
  s.creditor_id,
  s.debtor_id,
  s.amount,
  s.status
FROM settlements s
WHERE s.trip_id = '{trip_id}'
ORDER BY s.creditor_id, s.debtor_id;

-- Count transactions
SELECT COUNT(*) as settlement_count
FROM settlements
WHERE trip_id = '{trip_id}';
```

### Verify Poll Votes

```sql
-- Check vote counts per option
SELECT 
  po.poll_option_id,
  po.option_text,
  COUNT(pv.poll_vote_id) as vote_count
FROM poll_options po
LEFT JOIN poll_votes pv ON po.poll_option_id = pv.poll_option_id
WHERE po.poll_id = '{poll_id}'
GROUP BY po.poll_option_id, po.option_text;
```

---

## Performance Notes

**Expected Response Times:**
- Equal split calculation: < 10ms
- Custom split validation: < 5ms
- Settlement algorithm: < 100ms (even with 100 expenses)
- Poll vote recording: < 10ms

**If Slower:**
- Check database indexes
- Review query execution plans
- Profile Python/Go code

---

## Test Results

| Test | Result | Time | Notes |
|------|--------|------|-------|
| Equal split $300 ÷ 3 | ✓/✗ | ___ms | _____ |
| Custom split validation | ✓/✗ | ___ms | _____ |
| Invalid split rejection | ✓/✗ | ___ms | _____ |
| Settlement calc (3 ppl) | ✓/✗ | ___ms | _____ |
| Settlement calc (5 ppl) | ✓/✗ | ___ms | _____ |
| Poll vote increment | ✓/✗ | ___ms | _____ |
| Duplicate vote prevent | ✓/✗ | ___ms | _____ |

**Overall:** ___/7 tests passed

---

## Issues Found

No. | Issue | Severity | Impact | Resolution |
----|-------|----------|--------|-----------|
| 1 | _____ | HIGH/MED/LOW | _____ | _____ |

---

## Wednesday Summary

- Time spent: _____ hours
- All tests passed: YES / NO
- Algorithms verified: YES / NO
- Ready for Day 4: YES / NO
