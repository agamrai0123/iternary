# Phase A Week 2 - Day 2: API Endpoint Testing

**Date:** March 26, 2026 (Tuesday)  
**Duration:** 3-4 hours  
**Goal:** Test all 16 API endpoints with proper HTTP status codes

---

## Pre-Test Setup

### Start Application Server

```bash
cd /d/Learn/iternary/itinerary-backend
go run . &
# OR
./itinerary-backend.exe &
```

**Expected Output:**
```
[INFO] Itinerary service starting on port 8080
[INFO] Database connected to Oracle/PostgreSQL
[INFO] Routes registered: 16 group endpoints
[INFO] Server ready to accept requests
```

### Test Tools

Use any of the following:
- **curl** (command-line)
- **Postman** (GUI - recommended for documentation)
- **Thunder Client** (VS Code extension)
- **REST Client** extension in VS Code

---

## Test Scenario 1: Group Trip Management

### 1.1 CREATE Group Trip

**Endpoint:** `POST /api/v1/group-trips`

**Headers:**
```
Authorization: Bearer {valid_jwt_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Bali Adventure 2026",
  "destination_id": "dest-001",
  "budget": 50000,
  "duration": 7,
  "start_date": "2026-05-01",
  "description": "Amazing beach trip with friends"
}
```

**Expected Response:** HTTP 201 Created
```json
{
  "id": "trip-{uuid}",
  "title": "Bali Adventure 2026",
  "destination_id": "dest-001",
  "owner_id": "user-001",
  "budget": 50000,
  "duration": 7,
  "start_date": "2026-05-01",
  "status": "planning",
  "member_count": 1,
  "created_at": "2026-03-26T...",
  "updated_at": "2026-03-26T..."
}
```

**Test Cases:**
- [x] Valid creation returns 201
- [x] Response contains ID
- [x] Status defaults to "planning"
- [x] Member count is 1 (owner)
- [ ] Without title returns 400
- [ ] Without destination returns 400
- [ ] Without JWT token returns 401

**Results:**
| Test Case | Expected | Result | Status |
|-----------|----------|--------|--------|
| Valid creation | 201 | _____ | ✓/✗ |
| Missing title | 400 | _____ | ✓/✗ |
| No auth token | 401 | _____ | ✓/✗ |
| Invalid budget | 400 | _____ | ✓/✗ |

### 1.2 GET Group Trip

**Endpoint:** `GET /api/v1/group-trips/{id}`

**Expected Response:** HTTP 200 OK
```json
{
  "id": "trip-{uuid}",
  "title": "Bali Adventure 2026",
  "owner_id": "user-001",
  "member_count": 1,
  "members": [
    {
      "id": "member-001",
      "user_id": "user-001",
      "name": "Alice",
      "role": "owner",
      "status": "active",
      "joined_at": "2026-03-26T..."
    }
  ],
  "expenses_total": 0,
  "settlements": [],
  "created_at": "2026-03-26T...",
  "updated_at": "2026-03-26T..."
}
```

**Test Cases:**
- [x] Valid ID returns 200
- [x] Response contains member list
- [x] Expenses total field present
- [ ] Invalid ID returns 404
- [ ] Non-existent ID returns 404

### 1.3 LIST Group Trips

**Endpoint:** `GET /api/v1/group-trips?page=1&limit=10`

**Expected Response:** HTTP 200 OK
```json
{
  "trips": [
    {
      "id": "trip-001",
      "title": "Bali Adventure 2026",
      "member_count": 3,
      "expenses_total": 15000,
      "created_at": "2026-03-26T..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "pages": 1
  }
}
```

**Test Cases:**
- [x] Returns list of trips
- [x] Pagination info included
- [x] Default pagination works (no params)
- [ ] Page 999 returns empty list with 404 or 200?

### 1.4 UPDATE Group Trip

**Endpoint:** `PUT /api/v1/group-trips/{id}`

**Request Body:**
```json
{
  "title": "Bali Adventure 2026 - Updated",
  "budget": 60000,
  "status": "active"
}
```

**Expected Response:** HTTP 200 OK
```json
{
  "id": "trip-001",
  "title": "Bali Adventure 2026 - Updated",
  "budget": 60000,
  "status": "active",
  "updated_at": "2026-03-26T..."
}
```

**Test Cases:**
- [x] Valid update returns 200
- [x] Only owner can update (401 if not owner)
- [x] Invalid status returns 400
- [ ] Partial update works
- [ ] Cannot change owner_id

### 1.5 DELETE Group Trip

**Endpoint:** `DELETE /api/v1/group-trips/{id}`

**Expected Response:** HTTP 204 No Content (or 200 OK with message)

**Test Cases:**
- [x] Valid delete returns 204 or 200
- [x] Only owner can delete
- [x] Trip not found returns 404
- [ ] Cannot delete if expenses exist (409 Conflict)

---

## Test Scenario 2: Group Member Management

### 2.1 INVITE Group Member

**Endpoint:** `POST /api/v1/group-trips/{id}/members/invite`

**Request Body:**
```json
{
  "email": "bob@example.com",
  "role": "editor"
}
```

**Expected Response:** HTTP 201 Created
```json
{
  "id": "member-002",
  "trip_id": "trip-001",
  "user_email": "bob@example.com",
  "role": "editor",
  "status": "pending",
  "invited_at": "2026-03-26T...",
  "invitation_token": "token-{uuid}"
}
```

**Test Cases:**
- [x] Valid invite returns 201
- [x] Invitation has pending status
- [x] Only owner/editor can invite
- [ ] Invalid email returns 400
- [ ] Already member returns 409 Conflict
- [ ] Non-existent user returns 404

### 2.2 LIST Group Members

**Endpoint:** `GET /api/v1/group-trips/{id}/members`

**Expected Response:** HTTP 200 OK
```json
{
  "members": [
    {
      "id": "member-001",
      "user_id": "user-001",
      "name": "Alice",
      "role": "owner",
      "status": "active",
      "joined_at": "2026-03-26T..."
    },
    {
      "id": "member-002",
      "user_id": "user-002",
      "name": "Bob",
      "role": "editor",
      "status": "pending",
      "invited_at": "2026-03-26T..."
    }
  ]
}
```

### 2.3 RESPOND TO Invite

**Endpoint:** `POST /api/v1/group-trips/{id}/members/respond`

**Request Body:**
```json
{
  "invitation_token": "token-{uuid}",
  "accept": true
}
```

**Expected Response:** HTTP 200 OK
```json
{
  "id": "member-002",
  "status": "active",
  "message": "Invitation accepted"
}
```

**Test Cases:**
- [x] Accept invitation returns 200, status → active
- [x] Decline invitation returns 200, status → declined
- [ ] Invalid token returns 400
- [ ] Expired token returns 410 Gone
- [ ] Already responded returns 409 Conflict

### 2.4 REMOVE Group Member

**Endpoint:** `DELETE /api/v1/group-trips/{id}/members/{member_id}`

**Expected Response:** HTTP 204 No Content

**Test Cases:**
- [x] Valid remove returns 204
- [x] Only owner/editor can remove
- [x] Cannot remove owner
- [ ] Member not found returns 404

### 2.5 LEAVE Group Trip

**Endpoint:** `POST /api/v1/group-trips/{id}/members/leave`

**Expected Response:** HTTP 200 OK
```json
{
  "message": "Left the trip successfully",
  "trip_id": "trip-001"
}
```

**Test Cases:**
- [x] Member can leave
- [x] Owner cannot leave (must transfer ownership or delete)
- [ ] Not a member returns 403

---

## Test Scenario 3: Expense Management

### 3.1 ADD Expense

**Endpoint:** `POST /api/v1/group-trips/{id}/expenses`

**Request Body (Equal Split):**
```json
{
  "description": "Flight tickets",
  "amount": 3000,
  "paid_by_id": "user-001",
  "category": "transportation",
  "split_type": "equal",
  "split_among": ["user-001", "user-002", "user-003"]
}
```

**Expected Response:** HTTP 201 Created
```json
{
  "id": "expense-001",
  "trip_id": "trip-001",
  "description": "Flight tickets",
  "amount": 3000,
  "paid_by_id": "user-001",
  "category": "transportation",
  "split_type": "equal",
  "splits": [
    {"user_id": "user-001", "amount": 1000, "paid": true},
    {"user_id": "user-002", "amount": 1000, "paid": false},
    {"user_id": "user-003", "amount": 1000, "paid": false}
  ],
  "created_at": "2026-03-26T..."
}
```

**Test Cases (Equal Split):**
- [x] Valid equal split returns 201
- [x] Splits equally across members
- [x] Only member can add expense
- [ ] Amount exceeds budget returns 400
- [ ] Invalid category returns 400

**Request Body (Custom Split):**
```json
{
  "description": "Hotel",
  "amount": 1500,
  "paid_by_id": "user-001",
  "split_type": "custom",
  "split_details": [
    {"user_id": "user-001", "amount": 500},
    {"user_id": "user-002", "amount": 600},
    {"user_id": "user-003", "amount": 400}
  ]
}
```

**Test Cases (Custom Split):**
- [x] Custom split amounts match total
- [ ] Amounts don't match total returns 400
- [ ] User not in trip returns 403

### 3.2 LIST Expenses

**Endpoint:** `GET /api/v1/group-trips/{id}/expenses?category=transportation`

**Expected Response:** HTTP 200 OK
```json
{
  "expenses": [
    {
      "id": "expense-001",
      "description": "Flight tickets",
      "amount": 3000,
      "category": "transportation",
      "paid_by_name": "Alice",
      "created_at": "2026-03-26T..."
    }
  ],
  "summary": {
    "total": 3000,
    "by_category": {
      "transportation": 3000
    }
  }
}
```

### 3.3 GET Expense Report

**Endpoint:** `GET /api/v1/group-trips/{id}/report`

**Expected Response:** HTTP 200 OK
```json
{
  "trip_id": "trip-001",
  "trip_title": "Bali Adventure 2026",
  "total_expenses": 4500,
  "member_expenses": [
    {
      "user_id": "user-001",
      "name": "Alice",
      "paid": 3000,
      "owes": 500,
      "balance": 2500
    },
    {
      "user_id": "user-002",
      "name": "Bob",
      "paid": 1500,
      "owes": 1600,
      "balance": -100
    }
  ],
  "summary": {
    "total_expenses": 4500,
    "average_per_person": 1500,
    "largest_expense": 3000,
    "expense_count": 3
  }
}
```

---

## Test Scenario 4: Poll Management

### 4.1 CREATE Poll

**Endpoint:** `POST /api/v1/group-trips/{id}/polls`

**Request Body:**
```json
{
  "title": "Best restaurant for dinner?",
  "poll_type": "multiple_choice",
  "options": [
    "Italian",
    "Thai",
    "Mexican",
    "Japanese"
  ],
  "allow_multiple": false,
  "created_by_id": "user-001"
}
```

**Expected Response:** HTTP 201 Created
```json
{
  "id": "poll-001",
  "trip_id": "trip-001",
  "title": "Best restaurant for dinner?",
  "poll_type": "multiple_choice",
  "options": [
    {"id": "opt-001", "text": "Italian", "vote_count": 0},
    {"id": "opt-002", "text": "Thai", "vote_count": 0},
    {"id": "opt-003", "text": "Mexican", "vote_count": 0},
    {"id": "opt-004", "text": "Japanese", "vote_count": 0}
  ],
  "allow_multiple": false,
  "created_at": "2026-03-26T..."
}
```

### 4.2 GET Poll

**Endpoint:** `GET /api/v1/group-trips/{id}/polls/{poll_id}`

**Expected Response:** HTTP 200 OK with full poll details

### 4.3 LIST Polls

**Endpoint:** `GET /api/v1/group-trips/{id}/polls`

**Expected Response:** HTTP 200 OK
```json
{
  "polls": [
    {
      "id": "poll-001",
      "title": "Best restaurant for dinner?",
      "poll_type": "multiple_choice",
      "option_count": 4,
      "vote_count": 2,
      "created_at": "2026-03-26T..."
    }
  ]
}
```

### 4.4 VOTE on Poll

**Endpoint:** `POST /api/v1/group-trips/{id}/polls/{poll_id}/vote`

**Request Body:**
```json
{
  "option_id": "opt-002",
  "voted_by_id": "user-002"
}
```

**Expected Response:** HTTP 201 Created or 200 OK
```json
{
  "vote_id": "vote-001",
  "poll_id": "poll-001",
  "option_id": "opt-002",
  "option_text": "Thai",
  "voted_by_id": "user-002",
  "voted_at": "2026-03-26T..."
}
```

**Test Cases:**
- [x] Valid vote returns 201
- [x] Vote count increments on option
- [ ] Duplicate vote returns 409 Conflict
- [ ] Invalid option returns 400
- [ ] User not in trip returns 403

---

## Complete Workflow Test

### Scenario: Full Trip Planning Cycle

```
1. CREATE trip: POST /api/v1/group-trips
   Expected: 201, trip_id = T1

2. INVITE member: POST /api/v1/group-trips/{T1}/members/invite
   Expected: 201, invitation status = pending

3. RESPOND to invite: POST /api/v1/group-trips/{T1}/members/respond
   Expected: 200, status = active

4. ADD expense: POST /api/v1/group-trips/{T1}/expenses
   Expected: 201, expense_id = E1, splits created

5. ADD second expense: POST /api/v1/group-trips/{T1}/expenses
   Expected: 201, expense_id = E2

6. CREATE poll: POST /api/v1/group-trips/{T1}/polls
   Expected: 201, poll_id = P1

7. VOTE on poll: POST /api/v1/group-trips/{T1}/polls/{P1}/vote
   Expected: 201, vote increments count

8. GET report: GET /api/v1/group-trips/{T1}/report
   Expected: 200, report shows settlements

9. Verify all 8 steps succeed with proper HTTP codes
   Expected: All 2xx responses
```

---

## HTTP Status Code Verification

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | GET, successful operation result |
| 201 | Created | POST, resource successfully created |
| 204 | No Content | DELETE successful, no response body |
| 400 | Bad Request | Invalid input, validation failed |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Authenticated but no permission |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource conflict (duplicate, can't delete) |
| 500 | Server Error | Unexpected server error |

**Test Matrix:**

| Endpoint | Method | Success Code | Error Codes |
|----------|--------|--------------|------------|
| /group-trips | POST | 201 | 400, 401 |
| /group-trips | GET | 200 | 401 |
| /group-trips/{id} | GET | 200 | 401, 404 |
| /group-trips/{id} | PUT | 200 | 400, 401, 403, 404 |
| /group-trips/{id} | DELETE | 204 | 401, 403, 404, 409 |
| /group-trips/{id}/members/invite | POST | 201 | 400, 401, 403, 404, 409 |
| /group-trips/{id}/members | GET | 200 | 401, 404 |
| /group-trips/{id}/members/respond | POST | 200 | 400, 401, 404, 409 |
| /group-trips/{id}/members/{member_id} | DELETE | 204 | 401, 403, 404 |
| /group-trips/{id}/members/leave | POST | 200 | 401, 403, 404 |
| /group-trips/{id}/expenses | POST | 201 | 400, 401, 403, 404 |
| /group-trips/{id}/expenses | GET | 200 | 401, 404 |
| /group-trips/{id}/report | GET | 200 | 401, 404 |
| /group-trips/{id}/polls | POST | 201 | 400, 401, 403, 404 |
| /group-trips/{id}/polls | GET | 200 | 401, 404 |
| /group-trips/{id}/polls/{poll_id}/vote | POST | 201 | 400, 401, 403, 404, 409 |

---

## Postman Collection

Create a Postman collection with all 16 endpoints:

```json
{
  "info": {
    "name": "Itinerary API - Phase A Week 2",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080/api/v1",
      "type": "string"
    },
    {
      "key": "auth_token",
      "value": "{{JWT_TOKEN}}",
      "type": "string"
    },
    {
      "key": "trip_id",
      "value": "",
      "type": "string"
    }
  ],
  "item": [
    {
      "name": "Group Trips",
      "item": [
        {
          "name": "Create Trip",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{auth_token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{...}"
            },
            "url": {
              "raw": "{{base_url}}/group-trips",
              "host": ["{{base_url}}"],
              "path": ["group-trips"]
            }
          }
        }
      ]
    }
  ]
}
```

---

## Test Results Summary

**Tested by:** _________________  
**Date:** _________________

| Test Category | Passed | Failed | Total | Status |
|---------------|--------|--------|-------|--------|
| Trip Management (5 tests) | ___ | ___ | 5 | ✓/✗ |
| Member Management (5 tests) | ___ | ___ | 5 | ✓/✗ |
| Expense Management (3 tests) | ___ | ___ | 3 | ✓/✗ |
| Poll Management (4 tests) | ___ | ___ | 4 | ✓/✗ |
| **Total** | ___ | ___ | **16** | **✓/✗** |

**HTTP Status Codes Verified:**
- [ ] 200 OK
- [ ] 201 Created
- [ ] 204 No Content
- [ ] 400 Bad Request
- [ ] 401 Unauthorized
- [ ] 403 Forbidden
- [ ] 404 Not Found
- [ ] 409 Conflict

---

## Issues Found

No. | Issue | Severity | Resolution | Status |
----|-------|----------|-----------|--------|
| 1 | _____ | HIGH/MED/LOW | _____ | OPEN/FIXED |
| 2 | _____ | HIGH/MED/LOW | _____ | OPEN/FIXED |

**Total Issues:** _____

---

## Next Steps

- [ ] All tests passing? YES → Proceed to Day 3
- [ ] Some tests failing? → Fix issues and retest
- [ ] Need more time? → Extend to Wednesday

**Tuesday Summary:**
- Time spent: _____ hours
- Tests passed: _____/16
- Blocking issues: YES/NO
- Ready for Day 3: YES/NO
