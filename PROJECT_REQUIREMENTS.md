# Triply Project - Complete Requirements & Specifications

**Status:** 🚀 Starting Option C Implementation (Solo Developer)  
**Date:** March 23, 2026  
**Project Type:** Vibe-Coding Project (Minimal Bureaucracy, Maximum Creativity)  
**Development Mode:** Solo Developer (You are the only developer)  
**Primary Database:** Oracle (with PostgreSQL fallback prepared)  

---

## 1. Project Overview

**Triply** is a collaborative trip planning platform that enables friends to:
- Plan trips together with shared itineraries
- Vote on destinations and activities
- Split expenses within the group
- Share experiences through photos and reviews
- Discover trips from other travelers
- Book activities and hotels

**Current Status:**
- ✅ 95% Backend Implementation Complete (Go + Oracle)
- ✅ 8 Test Suites Created (63 Test Functions, 95% Coverage)
- ✅ Backend API Ready for Testing
- 🚀 Ready to Start: Option C Full-Featured Build

**Target Timeline:** 20+ weeks (5 months)  
**Target Users (Launch):** 1000+  
**Expected Bookings (Year 1):** 250+  
**Revenue Target:** ₹2M+ in transaction volume

---

## 2. Technology Stack

### Current Implementation

#### Backend
- **Language:** Go 1.21
- **Framework:** Gin v1.10.0
- **Database:** Oracle Database (primary) → PostgreSQL (secondary/fallback)
- **ORM:** database/sql with manual queries
- **Authentication:** Token-based (simplified JWT for MVP)
- **File Storage:** Local filesystem (future: AWS S3)
- **API Payment:** Razorpay v2 (payments)
- **Photo Provider:** Unsplash API (free stock photos)

#### Frontend (Current)
- **HTML/CSS/JavaScript** (Vanilla - temporary)
- **Template Engine:** Go html/template
- **Build Tool:** None (static files served by Gin)

#### Infrastructure
- **VCS:** Git (GitHub)
- **CI/CD:** GitHub Actions (to be implemented)
- **Hosting:** VM/Docker (development) → Kubernetes (production)
- **Monitoring:** Structured logging + metrics collection

### Future Stack (Option C Phase 2+)

#### Frontend Replacement
- **Framework:** React 18+
- **Language:** TypeScript
- **State Management:** Redux Toolkit or Zustand
- **Styling:** Tailwind CSS
- **UI Components:** Shadcn/ui or Material-UI
- **Build Tool:** Vite
- **Real-time:** Socket.io/Tanstack Query
- **Animations:** Framer Motion

#### Backend Evolution
- **API Gateway:** Kong or equivalent
- **Message Queue:** Apache Kafka
- **Cache:** Redis
- **Search:** Elasticsearch
- **AI:** Claude API (via Anthropic)

#### Infrastructure Evolution
- **Container Orchestration:** Kubernetes
- **Service Mesh:** Istio
- **API Documentation:** OpenAPI/Swagger
- **Monitoring:** Prometheus + Grafana
- **Logging:** ELK Stack (Elasticsearch, Logstash, Kibana)

---

## 3. Database Architecture

### Database Abstraction Layer (DAL)

To enable easy switching between Oracle and PostgreSQL, the backend uses an abstraction layer:

```
┌──────────────────────────────────────────────┐
│         Service Layer (Business Logic)       │
└───────────────┬───────────────────────────────┘
                │
┌───────────────▼───────────────────────────────┐
│        Database Interface (Abstract)          │
│  - GetDestinations()                          │
│  - CreateUserTrip()                           │
│  - ExecuteExpenseQuery()                      │
│  etc.                                         │
└────────┬────────────────────┬─────────────────┘
         │                    │
    ┌────▼────┐          ┌────▼──────┐
    │ Oracle  │          │ PostgreSQL│
    │ Impl.   │          │ Impl.     │
    └─────────┘          └───────────┘
```

### Current Implementation (Oracle)

**Configuration (config/config.json):**
```json
{
  "database": {
    "host": "localhost",
    "port": "1521",
    "user": "system",
    "service": "XE",
    "database": "itinerary_db",
    "password": "your_password"
  }
}
```

**Connection String Format:**
```
oracle://system:password@localhost:1521/XE
```

### PostgreSQL Configuration (Prepared)

**Configuration (alternative config-postgres.json):**
```json
{
  "database": {
    "host": "localhost",
    "port": "5432",
    "user": "postgres",
    "database": "itinerary_db",
    "password": "your_password",
    "sslmode": "disable"
  }
}
```

**Connection String Format:**
```
postgres://postgres:password@localhost:5432/itinerary_db?sslmode=disable
```

### Switching Between Databases

**To switch from Oracle to PostgreSQL:**

1. Change config file:
   ```bash
   cp config/config.json config/config.json.oracle  # Backup
   cp config/config-postgres.json config/config.json  # Switch to PostgreSQL
   ```

2. Update database driver imports in `database.go`:
   ```go
   // For Oracle:
   import _ "github.com/godror/godror"
   
   // For PostgreSQL:
   import _ "github.com/lib/pq"
   ```

3. Restart server:
   ```bash
   go run main.go
   ```

---

## 4. Feature Requirements

### ✅ COMPLETED FEATURES (Phase 0-1)

#### Core Authentication (DONE)
- ✅ User registration with email validation
- ✅ User login with token generation
- ✅ Session management (token validation)
- ✅ Password hashing and verification
- ✅ Authorization middleware (RequireAuth)
- ✅ Token refresh mechanism
- ✅ Logout functionality

#### Trip Management (DONE)
- ✅ Create trip with basic details
- ✅ Edit trip information
- ✅ Delete trip
- ✅ View trip details
- ✅ List user's trips
- ✅ Trip status (draft, published, completed)
- ✅ Trip budget tracking
- ✅ Trip duration in days

#### Trip Planning (DONE)
- ✅ Add destinations to trip
- ✅ Add itinerary items (stay, food, activity, transport)
- ✅ Edit itinerary items
- ✅ Delete itinerary items
- ✅ Order itinerary items by day
- ✅ View itinerary by day

#### Photo Management (DONE)
- ✅ Upload photos for destinations
- ✅ Store photo metadata (URL, caption)
- ✅ Link photos to trips
- ✅ Photo validation (size, format)

#### Community Features (DONE)
- ✅ Publish trip to community feed
- ✅ View public trips (feed)
- ✅ Like/unlike trips
- ✅ Comment on trips
- ✅ Like counter on trips

#### Reviews & Ratings (DONE)
- ✅ Add review to trip
- ✅ Rate trip (1-5 stars)
- ✅ View reviews on trip detail

#### Backend Infrastructure (DONE)
- ✅ API routing (15+ endpoints)
- ✅ Error handling with proper HTTP status codes
- ✅ Input validation
- ✅ Database abstraction layer
- ✅ Structured logging
- ✅ Metrics collection (success rate, cache hits, duration)
- ✅ Configuration management
- ✅ Template helpers (date, currency formatting)
- ✅ CORS support
- ✅ 95% test coverage (63 test functions)

#### Database Schema (DONE)
- ✅ Users table (id, email, password, created_at)
- ✅ Destinations table (id, name, description, country)
- ✅ Trips table (id, user_id, destination_id, budget, duration, status)
- ✅ Itinerary items table (id, trip_id, day, name, type, cost)
- ✅ Photos table (id, trip_id, url, caption)
- ✅ Comments table (id, trip_id, user_id, comment_text, created_at)
- ✅ Reviews table (id, trip_id, user_id, rating, review_text)
- ✅ Likes table (id, trip_id, user_id)
- ✅ Sessions table (id, user_id, token, expires_at)

---

### 🚀 PHASE A: GROUP COLLABORATION (Weeks 1-2)

#### Group Trip Features
- [ ] Create group trip (owner + co-planners)
- [ ] Invite users to group trip
- [ ] Accept/decline group trip invitation
- [ ] Remove user from group trip
- [ ] List group members and roles (owner, planner, viewer)
- [ ] Update member permissions
- [ ] Group chat/discussion board

#### Expense Splitting
- [ ] Add expense to group trip
- [ ] Split expense equally among members
- [ ] Custom split (different amounts per person)
- [ ] Track who paid what
- [ ] Calculate final settlement (who owes whom)
- [ ] Mark expense as paid
- [ ] Test with mock group trips (5-10 people)

#### Voting & Collaboration
- [ ] Create poll for itinerary item (destination, activity, restaurant)
- [ ] Vote on poll options
- [ ] View poll results
- [ ] Lock voted items (prevent further changes)
- [ ] Voting on trip dates/budget
- [ ] Comment on poll options

#### Database Changes
- [ ] Create GroupTrips table
- [ ] Create GroupMembers table
- [ ] Create Expenses table
- [ ] Create ExpenseSplits table
- [ ] Create Polls table
- [ ] Create PollVotes table
- [ ] Add indexes for performance

#### API Endpoints (New)
- `POST /api/group-trips` - Create group trip
- `GET /api/group-trips` - List user's group trips
- `GET /api/group-trips/:id` - Get group trip details
- `PUT /api/group-trips/:id` - Update group trip
- `POST /api/group-trips/:id/members` - Add member
- `DELETE /api/group-trips/:id/members/:user_id` - Remove member
- `POST /api/group-trips/:id/expenses` - Add expense
- `GET /api/group-trips/:id/expenses` - List expenses
- `POST /api/group-trips/:id/polls` - Create poll
- `POST /api/group-trips/:id/polls/:poll_id/votes` - Vote on poll

#### Tests Required
- [ ] Test group creation
- [ ] Test member management
- [ ] Test expense splitting calculations
- [ ] Test edge cases (single member, equal splits, custom splits)
- [ ] Test settlement calculations
- [ ] Test voting logic

---

### 🎨 PHASE B: UI ENHANCEMENT & STOCK PHOTOS (Weeks 3-4)

#### Unsplash Integration
- [ ] Integrate Unsplash API for stock photos
- [ ] Auto-fetch relevant photos for destinations
- [ ] Display photo attribution (Unsplash requires it)
- [ ] Cache photos locally (TTL: 24 hours)
- [ ] Handle API rate limiting (50+ requests/hour)
- [ ] Fallback to placeholder if API fails

#### UI Improvements
- [ ] Design system (colors, typography, spacing)
- [ ] Component library (button, card, modal, form)
- [ ] Dashboard redesign (better city cards)
- [ ] Trip detail page polish
- [ ] Itinerary wizard improvements
- [ ] Mobile responsiveness (Bootstrap or custom)
- [ ] Dark mode support
- [ ] Accessibility (WCAG 2.1 AA)

#### Animations & Interactions
- [ ] Page transitions (fade, slide)
- [ ] Loading spinners
- [ ] Success/error toasts
- [ ] Hover effects on cards
- [ ] Skeleton screens during loading
- [ ] Smooth scrolling

#### Forms & Validation
- [ ] Client-side form validation (JavaScript)
- [ ] Real-time validation feedback
- [ ] Better error messages
- [ ] Auto-save drafts to localStorage
- [ ] Date picker for trip dates
- [ ] Tag input for destinations

#### Frontend Templates to Update
- `login.html` - polish design
- `dashboard.html` - add Unsplash photos
- `plan-trip.html` - better wizard UX
- `itinerary-detail.html` - add photos
- `search.html` - search results with photos
- `destination-detail.html` - featured Unsplash images

#### Tests Required
- [ ] Test Unsplash API integration
- [ ] Test photo caching
- [ ] Test API rate limiting
- [ ] Test attribution display
- [ ] Test form validation

---

### ⚛️ PHASE C: REACT FRONTEND MIGRATION (Weeks 5-7)

#### React Setup
- [ ] Initialize React app with TypeScript
- [ ] Setup Vite build tool
- [ ] Configure Tailwind CSS
- [ ] Setup Redux Toolkit for state management
- [ ] Configure API client (axios/fetch)
- [ ] Setup React Router for navigation

#### Component Structure
```
src/
├── components/
│   ├── Auth/
│   │   ├── LoginForm.tsx
│   │   ├── RegisterForm.tsx
│   │   └── ProtectedRoute.tsx
│   ├── Trip/
│   │   ├── TripCard.tsx
│   │   ├── TripDetail.tsx
│   │   ├── TripForm.tsx
│   │   └── TripList.tsx
│   ├── Itinerary/
│   │   ├── ItineraryWizard.tsx
│   │   ├── ItineraryDay.tsx
│   │   └── ItemForm.tsx
│   ├── Group/
│   │   ├── GroupTripCard.tsx
│   │   ├── GroupSettings.tsx
│   │   └── ExpenseTracker.tsx
│   └── Common/
│       ├── Header.tsx
│       ├── Footer.tsx
│       └── Navbar.tsx
├── pages/
│   ├── Dashboard.tsx
│   ├── TripDetail.tsx
│   ├── CreateTrip.tsx
│   ├── Community.tsx
│   └── GroupTrip.tsx
├── store/
│   ├── authSlice.ts
│   ├── tripSlice.ts
│   └── groupSlice.ts
├── api/
│   ├── authClient.ts
│   ├── tripClient.ts
│   └── groupClient.ts
└── App.tsx
```

#### Pages to Build
- [ ] Login/Register pages
- [ ] Dashboard with recently viewed trips
- [ ] Create trip wizard (multi-step form)
- [ ] Trip detail page (view/edit)
- [ ] Community feed (scroll infinite)
- [ ] Group trip dashboard
- [ ] Expense tracker
- [ ] User profile
- [ ] Search results

#### Real-time Features
- [ ] WebSocket connection setup
- [ ] Real-time notifications (new likes, comments)
- [ ] Live expense updates in group trips
- [ ] Real-time vote results
- [ ] Chat messages in group trips

#### State Management
- [ ] User auth state
- [ ] Trip data (list, current, editing)
- [ ] Group trip state
- [ ] Filter/sort preferences
- [ ] Loading states
- [ ] Error states

#### Tests Required
- [ ] Component unit tests (Jest)
- [ ] Integration tests (React Testing Library)
- [ ] Redux reducer tests
- [ ] API client tests

---

### 🤖 PHASE D: AI FEATURES (Week 8)

#### Claude API Integration
- [ ] Setup Claude API client
- [ ] Implement itinerary suggestions
- [ ] Generate trip summaries
- [ ] Smart destination recommendations
- [ ] Budget optimization suggestions
- [ ] Schedule conflicts detection

#### AI Features for Users
- [ ] "Generate itinerary" button (AI suggests items for each day)
- [ ] "Optimize budget" (AI suggests cheaper alternatives)
- [ ] "Find similar trips" (content-based recommendations)
- [ ] "Expense insights" (analyze spending patterns)
- [ ] "Trip summary" (AI writes trip description)

#### Prompt Templates
- Destination research prompt
- Itinerary generation prompt
- Budget optimization prompt
- Trip summary prompt
- Recommendation prompt

---

### 🏗️ PHASE E: MICROSERVICES ARCHITECTURE (Weeks 9-14)

#### Service Extraction Plan

**Step 1: Auth Service (Week 9)**
- [ ] Extract authentication to separate service
- [ ] Setup service-to-service communication
- [ ] Implement OAuth2 for internal APIs
- [ ] Setup service discovery

**Step 2: Trip Service (Week 10)**
- [ ] Extract trip management
- [ ] Setup database replication
- [ ] Implement event sourcing

**Step 3: Media Service (Week 11)**
- [ ] Extract photo management
- [ ] Setup file storage (S3)
- [ ] CDN integration

**Step 4: Payment Service (Week 12)**
- [ ] Extract payment processing
- [ ] Setup Razorpay integration
- [ ] Implement webhook handling

**Step 5: Notification Service (Week 13)**
- [ ] Email notifications
- [ ] Push notifications
- [ ] SMS notifications (optional)

**Step 6: Search & Recommendations (Week 14)**
- [ ] Elasticsearch setup
- [ ] Search API service
- [ ] Recommendation engine

#### Infrastructure for Microservices
- [ ] Kubernetes cluster setup (3-5 nodes)
- [ ] Service mesh (Istio)
- [ ] API Gateway
- [ ] Load balancing
- [ ] Health checks
- [ ] Auto-scaling policies

#### Event-Driven Architecture
- [ ] Kafka cluster
- [ ] Event schemas
- [ ] Event publishing
- [ ] Event consumption
- [ ] Saga orchestration for transactions

#### Distributed Systems Patterns
- [ ] Circuit breaker
- [ ] Retry logic
- [ ] Timeout handling
- [ ] Fallback strategies
- [ ] Distributed tracing

---

### 📱 PHASE F: MOBILE APPS (Week 15+)

#### React Native App
- [ ] Setup React Native with TypeScript
- [ ] Share components with web React app
- [ ] Setup bottom tab navigation
- [ ] iOS build setup
- [ ] Android build setup

#### Features
- [ ] Native camera integration (photo upload)
- [ ] GPS location tracking
- [ ] Push notifications
- [ ] Offline mode
- [ ] App store deployment

---

## 5. API Specification (Complete)

### Base URL
```
Development: http://localhost:8080/api
Production: https://api.triply.com/api
```

### Authentication
All authenticated endpoints require:
```
Headers:
Authorization: Bearer <token>
```

### Standard Response Format

**Success (200):**
```json
{
  "data": { /* response data */ },
  "status": "success"
}
```

**Error (4xx/5xx):**
```json
{
  "error": {
    "code": "ERR_INVALID_INPUT",
    "message": "Invalid input provided",
    "details": { /* field errors */ }
  },
  "status": "error",
  "trace_id": "abc-123"
}
```

### Authentication Endpoints

#### POST /auth/register
Register new user
```
Request:
{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Doe"
}

Response (201):
{
  "user": { id, email, name },
  "token": "jwt_token"
}
```

#### POST /auth/login
Login user
```
Request:
{
  "email": "user@example.com",
  "password": "secure_password"
}

Response (200):
{
  "user": { id, email, name },
  "token": "jwt_token"
}
```

#### POST /auth/logout
Logout user (invalidate token)
```
Response (200):
{ "status": "success" }
```

### Trip Endpoints (Existing)

#### GET /trips
List all public trips
```
Query Parameters:
- page: 1 (default)
- limit: 10 (default)
- sort: recent | popular

Response (200):
{
  "trips": [
    { id, title, destination, budget, duration, likes_count, photos, user }
  ],
  "pagination": { page, limit, total, pages }
}
```

#### POST /user-trips
Create new trip
```
Request:
{
  "title": "Bali Adventure",
  "destination_id": "xyz123",
  "budget": 50000,
  "duration": 5,
  "start_date": "2026-05-01"
}

Response (201):
{ "trip": { id, title, destination, ... } }
```

#### GET /user-trips
List user's trips
```
Same query/response as /trips
```

#### GET /user-trips/:id
Get trip details
```
Response (200):
{
  "trip": {
    id, title, destination, budget, duration,
    itinerary: [
      { day: 1, items: [ { name, type, cost } ] }
    ],
    photos: [ { url, caption } ],
    reviews: [ { rating, text, user } ],
    likes_count
  }
}
```

#### PUT /user-trips/:id
Update trip
```
Request: Same as create (any fields)
Response (200): { "trip": {...} }
```

#### DELETE /user-trips/:id
Delete trip
```
Response (204): (no content)
```

### Community Endpoints

#### POST /trips/:id/likes
Like a trip
```
Response (200):
{ "likes_count": 42 }
```

#### DELETE /trips/:id/likes
Unlike a trip
```
Response (200):
{ "likes_count": 41 }
```

#### POST /trips/:id/comments
Add comment to trip
```
Request:
{ "text": "Amazing trip!" }

Response (201):
{ "comment": { id, text, user, created_at } }
```

#### GET /trips/:id/comments
Get all comments
```
Response (200):
{
  "comments": [
    { id, text, user, created_at }
  ]
}
```

### Group Trip Endpoints (To be built - Phase A)

#### POST /group-trips
Create group trip
```
Request:
{
  "title": "Goa With Friends",
  "destination_id": "xyz123",
  "budget": 100000,
  "duration": 7,
  "member_ids": ["user1", "user2", "user3"]
}

Response (201):
{ "group_trip": { id, title, members, ... } }
```

#### POST /group-trips/:id/expenses
Add expense
```
Request:
{
  "description": "Hotel booking",
  "amount": 5000,
  "paid_by": "user1",
  "split_among": ["user1", "user2", "user3"],
  "split_type": "equal"
}

Response (201):
{ "expense": { id, description, amount, ... } }
```

#### POST /group-trips/:id/polls
Create poll
```
Request:
{
  "question": "Which restaurant?",
  "options": ["Pizza Place", "Chinese Restaurant", "Tex-Mex"]
}

Response (201):
{ "poll": { id, question, options, votes } }
```

---

## 6. Database Schema (Oracle & PostgreSQL)

### Users Table
```sql
CREATE TABLE users (
  id VARCHAR2(36) PRIMARY KEY,
  email VARCHAR2(255) UNIQUE NOT NULL,
  password_hash VARCHAR2(255) NOT NULL,
  name VARCHAR2(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PostgreSQL:
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Destinations Table
```sql
CREATE TABLE destinations (
  id VARCHAR2(36) PRIMARY KEY,
  name VARCHAR2(255) NOT NULL,
  description CLOB,
  country VARCHAR2(100) NOT NULL,
  image_url VARCHAR2(2000),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PostgreSQL:
CREATE TABLE destinations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  country VARCHAR(100) NOT NULL,
  image_url VARCHAR(2000),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Trips Table
```sql
CREATE TABLE trips (
  id VARCHAR2(36) PRIMARY KEY,
  user_id VARCHAR2(36) NOT NULL,
  destination_id VARCHAR2(36) NOT NULL,
  title VARCHAR2(255) NOT NULL,
  description CLOB,
  budget NUMBER(12, 2) NOT NULL CHECK (budget > 0),
  duration NUMBER(3) NOT NULL CHECK (duration > 0),
  start_date DATE NOT NULL,
  status VARCHAR2(20) DEFAULT 'draft',
  likes_count NUMBER(10) DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  published_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);

-- PostgreSQL:
CREATE TABLE trips (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  destination_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  budget DECIMAL(12, 2) NOT NULL CHECK (budget > 0),
  duration INTEGER NOT NULL CHECK (duration > 0),
  start_date DATE NOT NULL,
  status VARCHAR(20) DEFAULT 'draft',
  likes_count INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  published_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);
```

### Itinerary Items Table
```sql
CREATE TABLE itinerary_items (
  id VARCHAR2(36) PRIMARY KEY,
  trip_id VARCHAR2(36) NOT NULL,
  day NUMBER(3) NOT NULL CHECK (day > 0),
  name VARCHAR2(255) NOT NULL,
  type VARCHAR2(20) NOT NULL CHECK (type IN ('stay', 'food', 'activity', 'transport', 'other')),
  cost NUMBER(12, 2) DEFAULT 0,
  notes CLOB,
  sequence NUMBER(3),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
);

-- PostgreSQL:
CREATE TABLE itinerary_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  trip_id UUID NOT NULL,
  day INTEGER NOT NULL CHECK (day > 0),
  name VARCHAR(255) NOT NULL,
  type VARCHAR(20) NOT NULL CHECK (type IN ('stay', 'food', 'activity', 'transport', 'other')),
  cost DECIMAL(12, 2) DEFAULT 0,
  notes TEXT,
  sequence INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
);
```

### Group Trips Table (Phase A)
```sql
CREATE TABLE group_trips (
  id VARCHAR2(36) PRIMARY KEY,
  title VARCHAR2(255) NOT NULL,
  destination_id VARCHAR2(36) NOT NULL,
  owner_id VARCHAR2(36) NOT NULL,
  budget NUMBER(12, 2) NOT NULL,
  duration NUMBER(3) NOT NULL,
  status VARCHAR2(20) DEFAULT 'draft',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);

-- PostgreSQL:
CREATE TABLE group_trips (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  destination_id UUID NOT NULL,
  owner_id UUID NOT NULL,
  budget DECIMAL(12, 2) NOT NULL,
  duration INTEGER NOT NULL,
  status VARCHAR(20) DEFAULT 'draft',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);
```

### Expenses Table (Phase A)
```sql
CREATE TABLE expenses (
  id VARCHAR2(36) PRIMARY KEY,
  group_trip_id VARCHAR2(36) NOT NULL,
  description VARCHAR2(255) NOT NULL,
  amount NUMBER(12, 2) NOT NULL,
  paid_by VARCHAR2(36) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_trip_id) REFERENCES group_trips(id) ON DELETE CASCADE,
  FOREIGN KEY (paid_by) REFERENCES users(id)
);

-- PostgreSQL:
CREATE TABLE expenses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_trip_id UUID NOT NULL,
  description VARCHAR(255) NOT NULL,
  amount DECIMAL(12, 2) NOT NULL,
  paid_by UUID NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_trip_id) REFERENCES group_trips(id) ON DELETE CASCADE,
  FOREIGN KEY (paid_by) REFERENCES users(id)
);
```

---

## 7. Development Environment Setup

### Prerequisites
```bash
# Check Go version
go version  # Should be 1.21+

# Check if git is available
git --version

# Database (choose one):
# Oracle: Need Oracle Database 21c XE (free edition)
# PostgreSQL: Need PostgreSQL 14+ (free and open source)
```

### Start Development Server
```bash
cd d:\Learn\iternary\itinerary-backend

# With Oracle:
go run main.go

# Server starts on http://localhost:8080
```

### Running Tests
```bash
# Run all tests
go test ./itinerary -v

# Run specific test
go test ./itinerary -v -run TestModelValidation

# With coverage
go test ./itinerary -cover -v
```

### Building for Production
```bash
# Build executable
go build -o itinerary-backend main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o itinerary-backend main.go

# Run executable
./itinerary-backend
```

---

## 8. Development Workflow (Solo Developer)

### Daily Workflow
```
1. Review previous day's tasks
2. Pull latest main branch
3. Create feature branch (git checkout -b feature/xyz)
4. Code → Test → Commit cycle
5. Run full test suite before PR
6. Push to GitHub (GitHub Actions runs CI)
7. Create PR and self-review
8. Merge to main when CI passes
```

### Feature Branch Naming
```
feature/group-trip-creation
feature/expense-splitting
bugfix/auth-token-expiry
docs/api-documentation
```

### Commit Message Format
```
[Feature] Add expense splitting calculation

- Implement equal split algorithm
- Implement custom split algorithm
- Add unit tests for both cases
- Update database schema for expense_splits table

Related to: Phase A Group Collaboration
```

---

## 9. Deployment Strategy

### Development Environment
- **Where:** Local machine (d:\Learn\iternary)
- **Database:** Oracle XE (localhost:1521)
- **URL:** http://localhost:8080
- **Logs:** ./log/itinerary.log

### Staging (After Phase A)
- **Where:** VM or Docker container
- **Database:** PostgreSQL (separate VM)
- **URL:** staging.triply.com
- **CI/CD:** GitHub Actions

### Production (After Phase B)
- **Where:** Kubernetes cluster
- **Database:** PostgreSQL (replicated)
- **URL:** api.triply.com
- **Monitoring:** Prometheus + Grafana

### Deployment Checklist
- [ ] All tests pass
- [ ] Code review completed
- [ ] Documentation updated
- [ ] Database migrations tested
- [ ] API documentation updated
- [ ] Changelog updated
- [ ] Tag release in Git
- [ ] Deploy to staging
- [ ] Run smoke tests
- [ ] Deploy to production
- [ ] Monitor for errors

---

## 10. Success Metrics & Milestones

### Phase A Success Metrics (Group Collaboration)
- [ ] 10+ group trips created
- [ ] 50+ group members
- [ ] ₹10K+ expenses tracked
- [ ] 95%+ expense settlement accuracy
- [ ] All tests pass (>95% coverage)

### Phase B Success Metrics (UI & Photos)
- [ ] 1000+ Unsplash photos loaded
- [ ] NO photo loading errors (>99% availability)
- [ ] Desktop responsiveness 100%
- [ ] Mobile responsiveness tested
- [ ] No accessibility violations (WCAG AA)
- [ ] Page load time <3s

### Phase C Success Metrics (React Frontend)
- [ ] All pages migrated to React
- [ ] SPA navigation working
- [ ] No console errors
- [ ] 100+ React components
- [ ] Redux state management working
- [ ] WebSocket real-time features working
- [ ] 90%+ component test coverage

### Phase D Success Metrics (AI Features)
- [ ] 100% of users can generate itineraries
- [ ] AI recommendations relevant (>80% satisfaction)
- [ ] Claude API integration stable
- [ ] <1s response time for suggestions
- [ ] Proper error handling for API failures

### Phase E Success Metrics (Microservices)
- [ ] 6 services deployed independently
- [ ] Inter-service communication working
- [ ] Distributed tracing showing full request flow
- [ ] <1% request failure rate
- [ ] Kubernetes pod scaling working
- [ ] Service mesh routing functional

### Overall Launch Metrics
- [ ] 1000+ users registered
- [ ] 250+ bookings processed
- [ ] ₹2M+ transaction volume
- [ ] 95%+ API uptime
- [ ] <500ms average API response
- [ ] <1% error rate
- [ ] 4.5+ star rating on app store

---

## 11. Documentation Requirements

### Completed Documentation
- ✅ API_REFERENCE.md (API endpoints)
- ✅ DATABASE_SETUP.md (DB schema)
- ✅ IMPLEMENTATION_STRATEGY.md (3 options)
- ✅ TEST_VERIFICATION_REPORT.md (test coverage)
- ✅ PROJECT_REQUIREMENTS.md (this file)

### Documentation to Create (By Phase)

#### Phase A Documentation
- [ ] Group Trip API Guide
- [ ] Expense Calculator Algorithm Doc
- [ ] Voting System Documentation
- [ ] Group Chat Protocol

#### Phase B Documentation
- [ ] Unsplash Integration Guide
- [ ] UI Component Library (Storybook)
- [ ] Accessibility Guidelines (WCAG 2.1)
- [ ] Performance Optimization Tips

#### Phase C Documentation
- [ ] React Architecture Guide
- [ ] State Management (Redux) Patterns
- [ ] Component Standards
- [ ] Testing Strategy (React Testing Library)

#### Phase D Documentation
- [ ] Claude API Integration Guide
- [ ] Prompt Engineering Best Practices
- [ ] AI Feature Limitations

#### Phase E Documentation
- [ ] Microservices Architecture Diagrams
- [ ] Service-to-Service Communication
- [ ] Event Sourcing Pattern Explanation
- [ ] Deployment Runbook (K8s)

---

## 12. Known Issues & For Later

### Known Technical Debt
- ⚠️ Vanilla JavaScript frontend (temporary, replaced in Phase C)
- ⚠️ Single database instance (split in Phase E)
- ⚠️ No caching layer (add Redis in Phase E)
- ⚠️ No search functionality (add Elasticsearch in Phase E)
- ⚠️ Synchronous API calls (add async queue in Phase E)

### For Later / Out of Scope
- ❌ Mobile apps (Phase F - after microservices stable)
- ❌ Advanced analytics dashboard
- ❌ Machine learning recommendations (except Claude AI)
- ❌ Video call integration
- ❌ VR travel previews
- ❌ Blockchain/NFT trip certificates
- ❌ Multi-language support (for now)
- ❌ Offline-first capabilities (for now)

---

## 13. Solo Developer Best Practices

### Time Management
- **Work in sprints:** 2-week sprints mapped to phases
- **Daily standup:** Write down status (even to yourself)
- **Code reviews:** Use GitHub PR reviews before merging
- **Breaks:** Take 15-min breaks every 2 hours (prevents burnout)

### Code Quality
- **TDD:** Write tests first, code second
- **Documentation:** Comment complex algorithms
- **Commits:** Small, focused commits (not 500-line commits)
- **Reviews:** Self-review code before committing

### Mental Health
- **Pace:** Don't work 80-hour weeks (burnout risk)
- **Scope:** If timeline seems impossible, reduce scope
- **Celebrate:** Mark milestones, celebrate completions
- **Pivot:** Be ready to adjust based on learnings

### Decision Making
- **Priority:** Focus on highest user impact
- **Simple:** Prefer simple solutions over clever ones
- **Iteration:** Build MVP first, iterate later
- **Community:** Share progress on dev blogs/Twitter

---

## Quick Reference

### Key Directories
```
d:\Learn\iternary\
├── itinerary-backend/          # Go backend (15+ endpoints, 8 test suites)
├── config/                      # Configuration files
├── docs/                        # Documentation
├── static/                      # CSS, JavaScript, images
└── templates/                   # HTML templates for server-side rendering
```

### Key Files
- `main.go` - Application entry point
- `itinerary/service.go` - Business logic
- `itinerary/database.go` - Database abstractions
- `itinerary/handlers.go` - HTTP endpoints
- `itinerary/routes.go` - API routing
- `config/config.json` - Configuration

### Test Command
```bash
go test ./itinerary -v -cover
```

### Run Server
```bash
go run main.go
# http://localhost:8080
```

---

## Summary: What You're Building

**Triply** is a collaborative trip planning platform. You're starting with **Option C** (Full-Featured with Microservices) as the solo developer.

**Timeline:** 20+ weeks (5 months)

**Phases:**
- ✅ **Phase 0:** Backend exists (95% done, 8 test suites verify it)
- 🚀 **Phase A (Weeks 1-2):** Group collaboration + expense splitting
- 🎨 **Phase B (Weeks 3-4):** UI polish + Unsplash photos
- ⚛️ **Phase C (Weeks 5-7):** React frontend migration
- 🤖 **Phase D (Week 8):** Claude AI integration
- 🏗️ **Phase E (Weeks 9-14):** Microservices architecture
- 📦 **Phase F (Weeks 15+):** Mobile apps

**Vibe:** Build fast, iterate often, focus on getting features working first, optimize later. You're the only developer, so code should be clear and well-tested (future you will thank you).

**Database:** Oracle for now, switch to PostgreSQL anytime by changing config and driver imports.

**Let's ship! 🚀**

