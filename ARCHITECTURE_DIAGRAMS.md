# Triply Project - Architecture & Data Flow Diagrams

## 1. Complete System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                                │
│                   (Web Browser / API Client)                        │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                    HTTP Requests (JSON)
                    REST API Calls
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    HTTP/ROUTING LAYER                               │
│                         (routes.go)                                 │
│  GET /api/destinations                                              │
│  POST /api/itineraries                                              │
│  GET /api/user-trips         ← Protected (requires auth)            │
└────────────────────────────┬────────────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
    ┌─────────────────────┐      ┌──────────────────────┐
    │  MIDDLEWARE LAYER   │      │  MIDDLEWARE LAYER    │
    │                     │      │                      │
    │ • Auth validation   │      │ • Metrics tracking   │
    │ • Token checking    │      │ • Request logging    │
    │ • User context      │      │ • Error handling     │
    └─────────────────────┘      └──────────────────────┘
              │                             │
              └──────────────┬──────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    HANDLER LAYER                                   │
│                  (handlers.go)                                      │
│  • GetDestinations()    - Parse params, call service               │
│  • CreateItinerary()    - Validate input, call service             │
│  • GetUserTrips()       - Check auth, call service                 │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                  Call Service Methods
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    SERVICE LAYER                                    │
│              (service.go, auth_service.go, etc)                     │
│                                                                      │
│  • Business Logic    - Validate inputs, apply rules                │
│  • Orchestration     - Call multiple DB methods                    │
│  • Data Transform    - Enrich/format data                          │
│  • Logging           - Track operations                            │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                 Execute DB Queries
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                  DATABASE/DATA LAYER                                │
│                    (database.go)                                    │
│  • Execute SQL queries                                              │
│  • Parse results into models                                        │
│  • Handle DB errors                                                │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      DATABASE                                       │
│                   (SQLite / PostgreSQL)                             │
│  • Users table                                                      │
│  • Destinations table                                               │
│  • Itineraries table                                                │
│  • Group Trips table                                                │
│  • etc.                                                             │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│               INFRASTRUCTURE & UTILITIES                            │
│  ┌─────────┬──────────┬──────────┬──────────┬──────────────────┐   │
│  │ Logger  │ Config   │ Metrics  │ Errors   │ Template Helpers │   │
│  └─────────┴──────────┴──────────┴──────────┴──────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Request-Response Flow Example: Getting Destinations

```
USER CLICKS: "Browse Destinations"
                    │
                    ▼
BROWSER SENDS: GET /api/destinations?page=1&page_size=10
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 1. ROUTE MATCHING (routes.go)                            │
│    router.GET("/api/destinations", handlers.GetDestination...)
│    ✓ Match found                                          │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 2. MIDDLEWARE EXECUTION                                  │
│    • MetricsMiddleware - Start tracking                  │
│    • LoggerMiddleware - Log request                      │
│    • ErrorHandlerMiddleware                              │
│    All pass, continue to handler                         │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 3. HANDLER (handlers.go - GetDestinations)              │
│                                                           │
│    page := c.Query("page")           → "1"              │
│    pageSize := c.Query("page_size")  → "10"             │
│                                                           │
│    Call Service:                                         │
│    destinations, total, err :=                           │
│        h.service.GetDestinations(1, 10)                 │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 4. SERVICE (service.go - GetDestinations)              │
│                                                           │
│    // Validate                                           │
│    if page < 1 { page = 1 }                             │
│    if pageSize > 100 { pageSize = 100 }                 │
│                                                           │
│    // Call Database                                      │
│    return s.db.GetDestinations(page, pageSize)         │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 5. DATABASE (database.go - GetDestinations)            │
│                                                           │
│    SQL Query Executed:                                   │
│    SELECT id, name, country FROM destinations           │
│    LIMIT 10 OFFSET 0                                    │
│                                                           │
│    Results parsed into Go structs:                       │
│    []Destination{                                        │
│        {ID: "dest_1", Name: "Paris", Country: "France"}│
│        {ID: "dest_2", Name: "Tokyo", Country: "Japan"} │
│        ...                                              │
│    }                                                     │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼ (Return values back up the stack)
┌──────────────────────────────────────────────────────────┐
│ 6. SERVICE RETURNS                                        │
│    destinations []Destination (10 items)                 │
│    total 50                                               │
│    err nil                                                │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 7. HANDLER FORMAT RESPONSE                              │
│                                                           │
│    if err != nil {                                       │
│        return error response                             │
│    }                                                      │
│                                                           │
│    c.JSON(200, gin.H{                                   │
│        "data": destinations,                             │
│        "total": 50,                                       │
│        "page": 1,                                         │
│        "page_size": 10                                    │
│    })                                                     │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ 8. MIDDLEWARE (Response Phase)                           │
│    • MetricsMiddleware - Record success                  │
│    • LoggerMiddleware - Log response                     │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────────────────────┐
│ HTTP RESPONSE SENT TO BROWSER:                           │
│ Status: 200 OK                                            │
│ Content-Type: application/json                            │
│                                                           │
│ {                                                         │
│   "data": [                                               │
│     {"id":"d1","name":"Paris","country":"France"},      │
│     {"id":"d2","name":"Tokyo","country":"Japan"},       │
│     ...                                                  │
│   ],                                                      │
│   "total": 50,                                            │
│   "page": 1,                                              │
│   "page_size": 10                                         │
│ }                                                         │
└──────────────────────────────────────────────────────────┘
                    │
                    ▼
BROWSER RECEIVES DATA & DISPLAYS TO USER
```

---

## 3. Authentication Flow

```
USER ENTERS CREDENTIALS IN LOGIN FORM
│
├─ username: "john_doe"
├─ password: "mypassword123"
│
▼
POST /api/auth/login
│
▼
┌────────────────────────────────────────┐
│ Handler: Login()                        │
│ • Receives credentials                  │
│ • Calls service.Login()                 │
└────────────┬───────────────────────────┘
             │
             ▼
┌────────────────────────────────────────┐
│ AuthService: Login()                    │
│ 1. Get user from database              │
│ 2. Hash provided password              │
│ 3. Compare with stored hash            │
│ 4. If match: GenerateToken()           │
│ 5. If no match: return InvalidPassword  │
└────────────┬───────────────────────────┘
             │
             ▼
┌────────────────────────────────────────┐
│ GenerateToken() - cryptographic random │
│ • Create 32 random bytes               │
│ • Base64 encode                        │
│ • Result: "a7kD9x2mZ1..." (unique)    │
└────────────┬───────────────────────────┘
             │
             ▼
┌────────────────────────────────────────┐
│ Save Session in Database               │
│ INSERT INTO sessions (...)             │
│   id: "session_abc123"                 │
│   user_id: "user_john_123"            │
│   token: "a7kD9x2mZ1..."              │
│   expires_at: 2026-04-19 10:00:00     │
│   created_at: 2026-04-12 10:00:00     │
└────────────┬───────────────────────────┘
             │
             ▼
┌────────────────────────────────────────┐
│ Response to Client                     │
│ {                                      │
│   "token": "a7kD9x2mZ1...",           │
│   "expires_in": 604800,                │
│   "user": {                            │
│     "id": "user_123",                  │
│     "username": "john_doe"             │
│   }                                    │
│ }                                      │
└────────────┬───────────────────────────┘
             │
             ▼
CLIENT STORES TOKEN (localStorage or cookie)

═══════════════════════════════════════════════

USER MAKES PROTECTED REQUEST:
GET /api/user-trips

WITH HEADER:
Authorization: Bearer a7kD9x2mZ1...

             │
             ▼
┌────────────────────────────────────────┐
│ AuthMiddleware.RequireAuth()           │
│ 1. Extract token from header           │
│ 2. Validate token format               │
│ 3. Query database for session          │
│ 4. Check if token matches              │
│ 5. Check if not expired                │
└────────────┬───────────────────────────┘
             │
      ┌──────┴──────┐
      │             │
      ▼ Valid       ▼ Invalid
   ┌───────┐    ┌────────────┐
   │ ✓ OK  │    │ ✗ 401 Error│
   │ Pass  │    │ Unauthorized
   │ to    │    │ Reject req │
   │Handler│    └────────────┘
   └───┬───┘
       │
       ▼
┌────────────────────────────────────────┐
│ Handler Executes                       │
│ c.Set("user_id", "user_123")           │
│ c.Set("token", "a7kD9x2mZ1...")        │
│ → Can access user context              │
└─────────────────────────────────────────┘
```

---

## 4. Group Trip Workflow

```
┌─────────────────────────────────────────────────┐
│ STEP 1: Create Group Trip                       │
│ POST /api/group-trips                           │
│                                                  │
│ Request Body:                                   │
│ {                                               │
│   "title": "Paris Summer 2026",                 │
│   "destination_id": "dest_paris",               │
│   "budget": 20000,                              │
│   "initial_members": ["user_2", "user_3"]      │
│ }                                               │
│                                                  │
│ Created by: user_1 (owner)                      │
└──────────────────┬────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│ STEP 2: Group Trip Created                      │
│ Database:                                        │
│ group_trips table:                              │
│   id: "group_abc123"                            │
│   title: "Paris Summer 2026"                    │
│   owner_id: "user_1"                            │
│   status: "draft"                               │
│                                                  │
│ group_members table:                            │
│   member_1: user_1 (owner)                      │
│   member_2: user_2 (pending)                    │
│   member_3: user_3 (pending)                    │
└──────────────────┬────────────────────────────┘
                   │
       ┌───────────┼───────────┐
       ▼           ▼           ▼
   user_1     user_2      user_3
   (owner)   (invited)   (invited)
       │         │           │
       │    Accept?      Accept?
       │    YES/NO       YES/NO
       │         │           │
       ▼         ▼           ▼
    ACTIVE    ACTIVE      DECLINED
    
    
┌─────────────────────────────────────────────────┐
│ STEP 3: Add Expenses                            │
│ POST /api/group-trips/group_abc123/expenses     │
│                                                  │
│ user_1 pays for accommodation: 3000             │
│ user_2 pays for flight: 1500                    │
│ user_3 pays for food: 800                       │
│                                                  │
│ Total: 5300 (split among 3 people)             │
└──────────────────┬────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│ STEP 4: Calculate Settlements                   │
│                                                  │
│ Per person share: 5300 / 3 = 1766.67           │
│                                                  │
│ Settlements:                                    │
│ • user_1 paid 3000, owes 1766.67              │
│   → user_1 should get back 1233.33             │
│                                                  │
│ • user_2 paid 1500, owes 1766.67              │
│   → user_2 owes 266.67                         │
│                                                  │
│ • user_3 paid 800, owes 1766.67               │
│   → user_3 owes 966.67                         │
│                                                  │
│ Summary:                                        │
│ user_2 pays user_1: 266.67                     │
│ user_3 pays user_1: 966.67                     │
└─────────────────────────────────────────────────┘
```

---

## 5. File Access Pattern

```
Request to GET /api/itineraries/:id
                    │
                    ▼
┌─────────────────────────────────────┐
│ routes.go                            │
│ Route definition                     │
│ Handler: handlers.GetItineraryDetail │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ handlers.go                          │
│ func (h *Handlers)                  │
│       GetItineraryDetail()          │
│ ├─ Parse itinerary ID               │
│ ├─ Call service.GetItinerary()      │
│ └─ Return JSON response             │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ service.go                           │
│ func (s *Service)                   │
│       GetItinerary()                │
│ ├─ Validate input                   │
│ ├─ Call db.GetItinerary()          │
│ ├─ Enrich with items                │
│ └─ Return data                      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ database.go                          │
│ func (db *Database)                 │
│       GetItinerary()                │
│ ├─ Execute SQL SELECT               │
│ ├─ Parse into struct                │
│ ├─ Fetch items separately           │
│ └─ Return query results             │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ models.go                            │
│ type Itinerary struct {              │
│   ID string                          │
│   Title string                       │
│   Items []ItineraryItem             │
│   ...                                │
│ }                                    │
└─────────────────────────────────────┘
```

---

## 6. Data Model Relationships

```
                         ┌──────────────────┐
                         │      USERS       │ ◄─────────┐
                         ├──────────────────┤           │
                         │ id (PK)          │           │
                         │ username         │           │
                         │ email            │           │
                         │ created_at       │           │
                         └────────┬─────────┘           │
                                  │                     │
                    ┌─────────────┼─────────────┐       │
                    │             │             │       │
                    ▼             ▼             ▼       │
        ┌────────────────┐  ┌──────────────┐  ┌──────────────┐
        │  ITINERARIES   │  │ GROUP_TRIPS  │  │ USER_PLANS   │
        ├────────────────┤  ├──────────────┤  ├──────────────┤
        │ id             │  │ id           │  │ id           │
        │ user_id ──┐    │  │ owner_id ────┼──►               │
        │ dest_id   │    │  │ dest_id      │  │ user_id ──────┼────┐
        │ title     │    │  │ created_at   │  │ original_id   │    │
        │ items ────┼────┼──┤              │  │               │    │
        │ created_at│    │  │              │  └──────────────┘    │
        └───┬───────┘    │  └──────────────┘                      │
            │            │        │                               │
            │            │        ▼                               │
            │            │  ┌──────────────────┐                 │
            │            └─►│  DESTINATIONS    │                 │
            │              ├──────────────────┤                 │
            │              │ id               │                 │
            │              │ name             │                 │
            │              │ country          │                 │
            │              │ image            │                 │
            │              └──────────────────┘                 │
            │                                                   │
            ▼                                                   │
   ┌──────────────────┐                                         │
   │ ITINERARY_ITEMS  │                                         │
   ├──────────────────┤                                         │
   │ id               │                                         │
   │ itinerary_id ────┼─────────────────────────────────────────┘
   │ day              │  (Foreign Key Reference)
   │ type             │
   │ name             │
   │ price            │
   │ created_at       │
   └──────────────────┘


         ┌──────────────────┐
         │ GROUP_MEMBERS    │
         ├──────────────────┤
         │ id               │
         │ group_trip_id ───┼───┐
         │ user_id          │   │
         │ role             │   │
         │ status           │   │
         │ joined_at        │   │
         └──────────────────┘   │
                                ▼
                        ┌──────────────────┐
                        │ GROUP_TRIPS      │
                        └──────────────────┘

         ┌──────────────────┐
         │ EXPENSES         │
         ├──────────────────┤
         │ id               │
         │ group_trip_id ───┼───┐
         │ paid_by_user_id  │   │
         │ amount           │   │
         │ description      │   │
         │ created_at       │   │
         └──────────────────┘   │
                                ▼
                        ┌──────────────────┐
                        │ GROUP_TRIPS      │
                        └──────────────────┘
```

---

## 7. HTTP Request Parts Breakdown

```
Example Request:
POST /api/itineraries?page=1 HTTP/1.1
Authorization: Bearer token_abc123
Content-Type: application/json

{
  "user_id": "user_123",
  "title": "Summer Paris"
}

Components:
├─ METHOD: POST                    (What verb/action)
├─ PATH: /api/itineraries          (Which endpoint)
├─ QUERY PARAMS: ?page=1           (URL parameters)
├─ HEADERS:                         (Metadata)
│  ├─ Authorization                (Who is making request)
│  └─ Content-Type: application/json (Body format)
└─ BODY:                            (Data sent)
   {
     "user_id": "user_123",
     "title": "Summer Paris"
   }

In Go Handler:
func (h *Handlers) CreateItinerary(c *gin.Context) {
    page := c.Query("page")              // From query params
    token := c.GetHeader("Authorization") // From headers
    
    var body map[string]interface{}
    c.BindJSON(&body)  // Parse body JSON
    
    user_id := body["user_id"]
    title := body["title"]
}
```

---

## 8. Error Handling Flow

```
Error Occurs in Code
        │
        ▼
┌──────────────────────────────────────┐
│ Error Type Determined                │
│                                      │
│ • NotFound        → 404              │
│ • Unauthorized    → 401              │
│ • BadRequest      → 400              │
│ • Internal        → 500              │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│ Create APIError                      │
│                                      │
│ NewNotFoundError("entity not found")│
│   └─ APIError{                       │
│       Code: NOT_FOUND                │
│       Message: "entity not found"   │
│       StatusCode: 404               │
│   }                                  │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│ Log Error                            │
│                                      │
│ logger.Error("database_error",       │
│     "error", err.Error())           │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│ Return to Client                     │
│                                      │
│ c.JSON(404, apiError.ToJSON())      │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│ Response Sent                        │
│                                      │
│ HTTP/1.1 404 Not Found              │
│ Content-Type: application/json       │
│                                      │
│ {                                    │
│   "code": "NOT_FOUND",              │
│   "message": "entity not found",    │
│   "status_code": 404                │
│ }                                    │
└──────────────────────────────────────┘
```

---

## 9. Middleware Stack Execution Order

```
REQUEST ARRIVES
        │
        ▼
┌──────────────────────────────────────┐
│ PanicRecoveryMiddleware              │
│ (Catch unexpected errors)            │
└──────────┬───────────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│ RequestLoggerMiddleware              │
│ (Log: method, path, query params)    │
└──────────┬───────────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│ MetricsMiddleware                    │
│ (Start timer, track requests)        │
└──────────┬───────────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│ ErrorHandlerMiddleware               │
│ (Setup error handling)               │
└──────────┬───────────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│ ErrorLoggerMiddleware                │
│ (Log errors)                         │
└──────────┬───────────────────────────┘
           │
           ▼
    ┌──────────────┐
    │ Route Match? │
    └──────┬───────┘
           │
    ┌──────┴──────┐
    │             │
    ▼ YES         ▼ NO
┌─────────┐  ┌─────────────┐
│ Check   │  │ 404 Error   │
│ Auth?   │  │ Return      │
└──┬───┬──┘  └─────────────┘
   │   │
   │   └─ Protected route?
   │
   ▼ NO/YES
┌──────────────────────────────────┐
│ Optional/Required Auth           │
│ Validate token if present        │
└──────────┬───────────────────────┘
           │
           ▼
    ┌──────────────┐
    │ Auth Valid?  │
    └──────┬───────┘
           │
    ┌──────┴──────────┐
    │                 │
    ▼ YES           ▼ NO
┌─────────┐     ┌────────────┐
│ Call    │     │ 401 Return │
│ Handler │     │ Abort      │
└────┬────┘     └────────────┘
     │
     ▼
┌──────────────────────────────────┐
│ HANDLER EXECUTES                │
│ (business logic here)           │
└──────────┬───────────────────────┘
           │
           ▼ (response ready)
┌──────────────────────────────────┐
│ MetricsMiddleware (exit)         │
│ Stop timer, record metrics       │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ RequestLogger (exit)             │
│ Log: status, response time       │
└──────────┬───────────────────────┘
           │
           ▼
    RESPONSE SENT TO CLIENT
```

---

## 10. Testing Pyramid

```
                             ┌─────────────────┐
                             │  END-TO-END     │
                             │  TESTS (E2E)   │
                             │                │
                             │ Full API flow  │
                             │ Through all    │
                             │ layers         │
                             └─────────────────┘
                            
                        ┌─────────────────────┐
                        │   INTEGRATION       │
                        │   TESTS             │
                        │                     │
                        │ Database + Service  │
                        │ Multiple components │
                        └─────────────────────┘

            ┌────────────────────────────────────┐
            │   UNIT TESTS                       │
            │                                    │
            │ Single function/method             │
            │ Mocked dependencies                │
            │ Fast execution                     │
            └────────────────────────────────────┘

TARGET: 85%+ Code Coverage
```

---

This visual guide should help you understand how all the components fit together!
