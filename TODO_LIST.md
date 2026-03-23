# Triply Development - Complete Todo List

**Total Items:** 68 tasks  
**Phases:** 7  
**Estimated Timeline:** 8-10 weeks to MVP  

---

## 📊 Priority Legend
- 🔴 **CRITICAL** - Blocks other work, must do before MVP
- 🟠 **HIGH** - Important for MVP, impacts user experience
- 🟡 **MEDIUM** - Nice to have, improves polish
- 🟢 **LOW** - Post-launch features, technical debt

---
## 🆕 NEW PRIORITIES (User-Requested Features)

**Four new major features to integrate:**
1. **Group Trip Collaboration & Voting** - Join friends' trips, vote on places
2. **UI Enhancement & Stock Photos** - Modern design with internet photos
3. **React Frontend Migration** - Smooth animations and better UX
4. **Microservices Architecture** - Scale from monolith as complexity grows

---
## Phase 0: Fix & Validate Core (Week 1) 
**Goal:** Ensure existing code works end-to-end

### Backend Testing & Debugging
- [ ] 🔴 Test `LoginHandler()` - verify token generation works
- [ ] 🔴 Test `DashboardHandler()` - load cities correctly
- [ ] 🔴 Test `PlanTripPageHandler()` - render form without errors
- [ ] 🔴 Debug itinerary detail API 404 bug (from previous session)
- [ ] 🔴 Verify database connection on startup
- [ ] 🔴 Test middleware: RequireAuth() and OptionalAuth()
- [ ] 🔴 Verify tokens are correctly extracted from Authorization header
- [ ] 🔴 Test error handling returns proper JSON responses

### Frontend Integration
- [ ] 🔴 Verify login.html POSTs to /auth/login correctly
- [ ] 🔴 Test localStorage token storage and retrieval
- [ ] 🔴 Verify token is sent in subsequent API requests
- [ ] 🔴 Test redirect to /dashboard after login
- [ ] 🔴 Test logout clears localStorage
- [ ] 🟠 Add loading indicators during API calls
- [ ] 🟠 Add error messages for failed login

### Manual Testing Checklist
- [ ] 🔴 Start server: `go run main.go`
- [ ] 🔴 Open http://localhost:8080
- [ ] 🔴 Login with demo credentials (traveler@example.com / password123)
- [ ] 🔴 See dashboard with city cards
- [ ] 🔴 Click "Get Started with Planning a Trip"
- [ ] 🔴 See plan-trip form
- [ ] 🔴 Logout and verify redirect to login

### Metrics & Logging
- [ ] 🟡 Set up log file rotation (avoid unlimited growth)
- [ ] 🟡 Configure log levels (debug in dev, warn in prod)
- [ ] 🟡 Verify metrics are recorded correctly
- [ ] 🟡 Test error logging captures stack traces

**Deliverable:** Users can complete login → dashboard → plan-trip flow without errors

---

## Phase 1: Trip Creation & Storage (Week 2-3)
**Goal:** Enable users to create and save trip plans

### API Development
- [ ] 🔴 Implement `CreateUserTrip()` handler - validate input
- [ ] 🔴 Implement `database.CreateUserTrip()` - insert into DB
- [ ] 🔴 Test trip creation via POST /api/user-trips
- [ ] 🔴 Implement `GetUserTrips()` handler - list user's trips
- [ ] 🔴 Test trip list via GET /api/user-trips
- [ ] 🔴 Implement `GetUserTrip()` handler - single trip detail
- [ ] 🔴 Test trip detail via GET /api/user-trips/:id
- [ ] 🔴 Verify 404 returned for non-existent trips
- [ ] 🔴 Implement `UpdateUserTrip()` handler
- [ ] 🔴 Test trip update via PUT /api/user-trips/:id
- [ ] 🔴 Implement `DeleteUserTrip()` handler
- [ ] 🔴 Test trip deletion via DELETE /api/user-trips/:id

### Trip Segments (Places)
- [ ] 🔴 Implement `AddTripSegment()` handler
- [ ] 🔴 Implement `database.AddTripSegment()`
- [ ] 🔴 Test segment creation via POST /api/user-trips/:id/segments
- [ ] 🔴 Implement `GetTripSegments()` handler
- [ ] 🔴 Test segment list via GET /api/user-trips/:id/segments
- [ ] 🟡 Implement segment update (PUT)
- [ ] 🟡 Implement segment delete (DELETE)

### Frontend Integration
- [ ] 🔴 Wire wizard Step 1 form to POST /api/user-trips
- [ ] 🔴 Store returned trip ID for segment creation
- [ ] 🔴 Wire Step 2 "Add Place" button to POST /api/user-trips/:id/segments
- [ ] 🔴 Display added segments dynamically on form
- [ ] 🔴 Implement "Save Draft" button
- [ ] 🟠 Implement "Load Draft" from My Trips page
- [ ] 🟠 Implement segment edit/delete UI

### Data Validation
- [ ] 🔴 Budget must be > 0 (checked at model & DB level)
- [ ] 🔴 Duration must be > 0
- [ ] 🔴 Segment day must be > 0 and <= trip duration
- [ ] 🔴 Required fields: title, destination, budget, duration
- [ ] 🟠 Location string length validation
- [ ] 🟠 GPS coordinate validation (-90 to 90 for lat, -180 to 180 for long)

### Authorization
- [ ] 🔴 User can only create trips for themselves
- [ ] 🔴 User can only edit/delete their own trips
- [ ] 🔴 User can only add segments to their own trips
- [ ] 🟡 Implement proper ownership check in handlers

**Deliverable:** Users can create multi-day trips and add places via both UI and API

---

## Phase 2: Photo Upload (Week 4)
**Goal:** Enable users to upload images for each place

### Photo Upload Handler
- [ ] 🔴 Implement `AddTripPhoto()` handler with multipart form parsing
- [ ] 🔴 Implement file validation: max 5MB size
- [ ] 🔴 Validate image format (jpg, png only)
- [ ] 🔴 Generate unique filename to prevent collisions
- [ ] 🔴 Enforce 3-photo limit per segment

### File Storage
- [ ] 🔴 Choose storage backend (local disk vs S3)
- [ ] 🔴 If local: `static/uploads/photos/` directory structure
- [ ] 🔴 If S3: Set up AWS credentials and bucket
- [ ] 🔴 Implement file path generation algorithm
- [ ] 🔴 Save photo URL to database

### Frontend
- [ ] 🔴 Add file picker to Step 3 of wizard
- [ ] 🔴 Show image preview after selection
- [ ] 🔴 Track uploaded count (currently X/3)
- [ ] 🔴 POST multipart request to /api/trip-segments/:id/photos
- [ ] 🟠 Show upload progress bar
- [ ] 🟠 Drag-and-drop upload support
- [ ] 🟠 Edit photo captions

### Image Optimization
- [ ] 🟡 Generate thumbnail on upload (100x100px)
- [ ] 🟡 Compress full-size image (max 2MB)
- [ ] 🟡 EXIF data stripping (privacy)
- [ ] 🟡 Cache-busting via URL params

**Deliverable:** Users can upload 1-3 photos per place with validation

---

## Phase 3: Reviews & Ratings (Week 4)
**Goal:** Enable users to rate and review places

### Review Handler
- [ ] 🔴 Implement `AddTripReview()` handler
- [ ] 🔴 Implement `database.AddTripReview()` - upsert (one review per segment)
- [ ] 🔴 Validate rating: 1-5 stars only
- [ ] 🔴 Limit review text: max 1000 characters
- [ ] 🔴 Store review with timestamp

### Frontend
- [ ] 🔴 Add star rating selector to Step 3 of wizard
- [ ] 🔴 Add review text input (optional)
- [ ] 🔴 Show rating UI (⭐⭐⭐⭐)
- [ ] 🔴 POST review to /api/trip-segments/:id/review
- [ ] 🟡 Show average rating for published trips
- [ ] 🟡 Display individual reviews on community posts

### Community Display
- [ ] 🟠 Retrieve reviews when fetching community posts
- [ ] 🟠 Calculate average rating per segment
- [ ] 🟠 Display reviews in trip detail page

**Deliverable:** Users can rate and review places; reviews show on community posts

---

## Phase 4: Travel Maps Integration (Week 4)
**Goal:** Display places on Google Maps

### API Setup
- [ ] 🟠 Get Google Maps API key
- [ ] 🟠 Load Google Maps in templates
- [ ] 🟠 Configure API restrictions (domain whitelist)

### Frontend
- [ ] 🟠 Display map in plan-trip form Step 2
- [ ] 🟠 Show clickable markers for each segment
- [ ] 🟠 Allow map-based place selection
- [ ] 🟠 Show route between segments (optional)

### Geocoding
- [ ] 🟠 Convert address → lat/long using Google Places API
- [ ] 🟠 Cache geocoding results to avoid API quota overages
- [ ] 🟠 Handle geocoding failures gracefully

**Deliverable:** Users can see places on map; can select locations via map search

---

## Phase 5: Trip Publishing & Community (Week 5-6)
**Goal:** Enable users to publish trips and see community feed

### Publishing Flow
- [ ] 🔴 Implement `PublishUserTrip()` handler
- [ ] 🔴 Implement `database.PublishUserTrip()` - insert into user_trip_posts
- [ ] 🔴 Set published=true, published_at=NOW()
- [ ] 🔴 Update trip status to "published"
- [ ] 🔴 Generate cover image (first photo or placeholder)

### Community Feed API
- [ ] 🔴 Implement `GetCommunityPosts()` handler
- [ ] 🔴 Return paginated list (20 posts per page)
- [ ] 🔴 Include user info, likes count, views count
- [ ] 🔴 Sort by published_at DESC (newest first)

### Frontend - Community Feed
- [ ] 🔴 Create community.html template
- [ ] 🔴 Display trip cards with cover image
- [ ] 🔴 Show: title, author, duration, budget, likes count
- [ ] 🔴 Implement infinite scroll (load more)
- [ ] 🟠 Add sort options: newest, most liked, trending
- [ ] 🟠 Add search by destination

### Like System
- [ ] 🔴 Implement `LikeUserTripPost()` handler
- [ ] 🔴 Implement `database.LikeUserTripPost()` - insert into likes table
- [ ] 🔴 Prevent duplicate likes (unique constraint)
- [ ] 🔴 Increment likes count on user_trip_posts
- [ ] 🔴 Implement `UnlikeUserTripPost()` handler
- [ ] 🟠 Show "Liked" state on frontend (filled heart icon)
- [ ] 🟠 Update likes count dynamically

### View Tracking
- [ ] 🟡 Increment views count when post is opened
- [ ] 🟡 Prevent counting same user twice (add user_id check)
- [ ] 🟡 Display views count on post

**Deliverable:** Users can publish trips; others can see community posts and like them

---

## Phase 6: Comments & Notifications (Week 6)
**Goal:** Enable discussions on published trips

### Comment API
- [ ] 🟠 Implement `AddComment()` handler for community posts
- [ ] 🟠 Implement `GetComments()` handler
- [ ] 🟠 Implement `DeleteComment()` handler (author only)
- [ ] 🟠 Store comments with user_id, timestamp

### Frontend
- [ ] 🟠 Add comment input field on community post detail
- [ ] 🟠 Display comments (newest first)
- [ ] 🟠 Show comment author name and avatar
- [ ] 🟠 Implement "Delete" button for own comments

### Notifications (MVP = email only)
- [ ] 🟠 Email notification: "Your post got a new like"
- [ ] 🟠 Email notification: "Someone commented on your trip"
- [ ] 🟠 Batch notifications (send once per day, not per action)
- [ ] 🟠 User setting: enable/disable notifications

### In-App Notifications (Future)
- [ ] 🟡 WebSocket setup for real-time notifications
- [ ] 🟡 Notification bell icon showing unread count
- [ ] 🟡 Notification dropdown menu

**Deliverable:** Users can discuss trips via comments; receive email notifications

---

## Phase 7: Payment & Booking Integration (Week 7-8)
**Goal:** Enable users to book through platform

### Razorpay Setup
- [ ] 🟠 Create Razorpay merchant account
- [ ] 🟠 Get API keys (test & production)
- [ ] 🟠 Implement order creation endpoint
- [ ] 🟠 Implement payment verification

### Order Model & Endpoints
- [ ] 🟠 Create `orders` table: id, user_id, trip_id, total_amount, status, razorpay_id
- [ ] 🟠 Implement `CreateOrder()` - sum up trip budget
- [ ] 🟠 Implement `VerifyPayment()` - check Razorpay webhook
- [ ] 🟠 Update order status on payment success

### Frontend
- [ ] 🟠 Add "Book Now" button on trip detail
- [ ] 🟠 Show total cost breakdown (accommodation + food + activities + transport)
- [ ] 🟠 Redirect to Razorpay payment window
- [ ] 🟠 Handle payment success/failure responses
- [ ] 🟡 Show order confirmation with ticket

### Affiliate Booking
- [ ] 🟡 Booking.com affiliate links for hotels
- [ ] 🟡 MakeMyTrip affiliate link generation
- [ ] 🟡 Track affiliate commissions
- [ ] 🟡 Generate shareable booking links

**Deliverable:** Users can pay for trips via Razorpay; receive booking confirmation

---

## 🆕 Phase 5A: Group Trips & Collaborative Voting (Week 8-9)
**Goal:** Friends can join trips together and vote on places for better decision-making

### Core Features
- [ ] 🔴 Allow users to invite friends to join existing trips
- [ ] 🔴 Track who's part of each trip (trip members)
- [ ] 🔴 Implement voting system for trip segments (places)
- [ ] 🔴 Calculate most-voted places per day
- [ ] 🔴 Display voting results on shared trip view
- [ ] 🟠 Shared expense tracking (Splitwise-like)
- [ ] 🟠 Settlement calculations (who owes whom)
- [ ] 🟠 Payment settlement flow

### Database Changes
- [ ] Create `trip_members` table (user_id, trip_id, role: owner/contributor/viewer)
- [ ] Create `segment_votes` table (user_id, segment_id, vote: 👍/👎)
- [ ] Create `trip_expenses` table (expense_id, payer_id, items, amounts)
- [ ] Create `expense_shares` table (expense_id, user_id, share_amount)
- [ ] Add `status` column to trip_members (invited/accepted/rejected)
- [ ] Add `vote_count` to trip_segments for quick access

### API Endpoints
- [ ] POST /api/user-trips/:id/invite - Invite friend to trip
- [ ] GET /api/user-trips/:id/members - List trip members
- [ ] POST /api/trip-segments/:id/vote - Vote on a place
- [ ] GET /api/trip-segments/:id/votes - Get voting results
- [ ] POST /api/user-trips/:id/expenses - Add expense
- [ ] GET /api/user-trips/:id/expenses - Get all expenses
- [ ] GET /api/user-trips/:id/settlement - Get settlement info

### Frontend Features
- [ ] Add "Invite to Trip" button on trip detail
- [ ] Show trip members list with voting power
- [ ] Add voting UI next to each place (👍 Like / 👎 Dislike)
- [ ] Display vote counts for each segment
- [ ] Add expense tracker section
- [ ] Show settlement summary per person
- [ ] "Accept Trip Invite" flow for invited users

### Authorization
- [ ] Only trip owner can invite
- [ ] Only members can vote
- [ ] Only contributors can add expenses
- [ ] Only trip owner can finalize settlement

**Deliverable:** Users can invite friends, collaboratively vote on places, track shared expenses, and settle payments

---

## 🆕 Phase 5B: UI Enhancement & Stock Photos (Week 8-9, parallel with 5A)
**Goal:** Enhance visual appeal with images and design polish

### Stock Photo Integration
- [ ] 🔴 Integrate Unsplash API for destination images
- [ ] 🔴 Auto-fetch cover photos for destinations on dashboard
- [ ] 🔴 Add hero banners to trip detail page
- [ ] 🟠 Cache downloaded images locally
- [ ] 🟠 Fallback to placeholder images if API fails
- [ ] 🟠 Add image credit attribution

### UI/UX Improvements
- [ ] 🔴 Improve dashboard layout with cards and grid
- [ ] 🔴 Add CSS animations (fade-in, slide-up, hover effects)
- [ ] 🔴 Better form styling (rounded inputs, shadows)
- [ ] 🟠 Add loading spinners and skeleton screens
- [ ] 🟠 Improve mobile responsiveness
- [ ] 🟠 Add dark mode support (optional)
- [ ] 🟡 Smooth transitions between pages

### Design Assets
- [ ] Create color scheme (~5 primary colors)
- [ ] Define typography (fonts, sizes, weights)
- [ ] Create icon set for place types (hotel, restaurant, activity, etc.)
- [ ] Design button styles (primary, secondary, danger)
- [ ] Create card components (trip card, place card, post card)

### Performance Optimization
- [ ] 🟠 Lazy load images on community feed
- [ ] 🟠 Compress images before upload
- [ ] 🟠 Add image preloading for hero sections
- [ ] 🟡 Implement CDN for photo delivery (if on cloud)

**Deliverable:** Modern-looking interface with professional images and smooth animations

---

## Phase 6: React Frontend Migration (Week 10-12)
**Goal:** Rebuild frontend in React for better performance and maintainability

### Why React?
- Current vanilla JS lacks component structure
- Hard to manage state across pages
- Animations are cumbersome without framework
- Difficult to scale with new features
- React enables faster development & better UX

### Phase 6a: Setup & Prerequisites (Week 10)
- [ ] 🟡 Create React project (Vite or Create React App)
- [ ] 🟡 Set up TypeScript for type safety
- [ ] 🟡 Install UI library (Tailwind/Material-UI)
- [ ] 🟡 Set up state management (Context API or Redux)
- [ ] 🟡 Configure API client (Axios with interceptors)
- [ ] 🟡 Set up routing (React Router v6)
- [ ] 🟡 Set up build pipeline and dev server

### Phase 6b: Component Library (Week 10)
- [ ] React components for:
  - Button, Input, Form, Modal, Card, Navbar, Footer
  - Loading Spinner, Toast Notifications, Breadcrumbs
  - Tab, Accordion, Dropdown, Badge, Avatar
- [ ] Hooks for common operations:
  - useAuth (auth context)
  - useTrip (trip data fetching)
  - useApi (API calls with error handling)
  - useForm (form handling)

### Phase 6c: Page Migration (Week 11)
- [ ] Login page (React version)
- [ ] Dashboard page with city grid
- [ ] Trip detail page with segments
- [ ] My trips list page
- [ ] Plan trip wizard (4-step React form)
- [ ] Community feed (infinite scroll)
- [ ] User profile page

### Phase 6d: Advanced Features (Week 12)
- [ ] Real-time updates (WebSocket Optional)
- [ ] Smooth page transitions and animations
- [ ] Image galleries with lightbox
- [ ] Map integration (Google Maps React component)
- [ ] File upload with progress bar
- [ ] Notification system at top of page

### Migration Strategy
- [ ] Keep Go backend unchanged
- [ ] Deploy React app to separate domain or subdirectory
- [ ] Run A/B tests (React vs vanilla)
- [ ] Gradually migrate users to React version
- [ ] Keep vanilla JS version as fallback
- [ ] After stabilization, retire vanilla version

**Deliverable:** Modern React frontend with smooth animations, better state management, and improved UX

---

## Phase 7: Advanced Features (Week 13+)
**Goal:** Add AI and automation capabilities

### Claude API Integration
- [ ] 🟡 Set up Anthropic API credentials
- [ ] 🟡 Create `TripGenerator` service
- [ ] 🟡 Prompt engineering: "Based on ₹50,000 and 5 days in Goa, generate an itinerary"
- [ ] 🟡 Parse AI response and create trip segments

### Intelligent Features
- [ ] 🟡 Trip remixing: "Compress this 7-day trip to 3 days"
- [ ] 🟡 Price staleness detection: "Hotel prices may have changed"
- [ ] 🟡 Personalized recommendations: "Your budget allows luxury 3-star hotels"
- [ ] 🟡 Duplicate detection: "Similar trip already exists"

### Automation
- [ ] 🟡 Daily email: "Popular trips in your favorite destination"
- [ ] 🟡 Smart re-ranking: Sort by (likes × recency) / (age_in_days)
- [ ] 🟡 Trending notifications: "This destination is trending +50% this week"

**Deliverable:** AI-powered recommendations and one-click trip generation

---

## 🆕 Phase 8: Microservices Architecture Migration (Week 14+)
**Goal:** Transition from monolith to microservices as complexity increases

### Why Microservices?
- Different services grow at different rates
- Payment service needs separate scaling
- Notification service should be independent
- Auth service could be shared across products
- Allows teams to work independently

### Phase 8a: Architecture Design (Week 14)
- [ ] 🟡 Define service boundaries:
  - **Auth Service** - User authentication, JWT, profiles
  - **Trip Service** - Trip CRUD, segments, planning
  - **Payment Service** - Orders, Razorpay integration
  - **Notification Service** - Email, in-app, push notifications
  - **Media Service** - Photo upload, storage, CDN
  - **Community Service** - Posts, likes, comments
  - **Expense Service** - Expense tracking, settlement
- [ ] Design inter-service communication (REST or gRPC)
- [ ] Design database per service (polyglot persistence)
- [ ] Plan API Gateway (single entry point)
- [ ] Design service discovery
- [ ] Design circuit breakers and retry logic

### Phase 8b: Service Extraction (Week 15-17)
- [ ] Extract Auth Service
  - Move auth handlers to new Go service
  - Database with users table only
  - JWT generation and validation
  - Integration with API Gateway
- [ ] Extract Payment Service
  - Move payment handlers to new Go service
  - Database with orders table
  - Razorpay integration
  - Webhook handling
- [ ] Extract Notification Service
  - New Node.js/Go service
  - Email integration (SendGrid/Resend)
  - Push notification sending
  - Event listeners from other services
- [ ] Extract Media Service
  - Photo upload and storage
  - S3 or cloud storage integration
  - Image processing and CDN

### Phase 8c: Infrastructure (Week 18)
- [ ] 🟡 API Gateway setup (Kong or AWS API Gateway)
- [ ] Service mesh (optional: Istio or Consul)
- [ ] Container orchestration (Docker Swarm or Kubernetes)
- [ ] Inter-service communication:
  - Service A → Service B via REST
  - Event-driven messaging (RabbitMQ/Kafka) for async
  - Shared cache (Redis) for frequently accessed data
- [ ] Database migration:
  - Trip Service: PostgreSQL (same as before)
  - Auth Service: Separate PostgreSQL instance
  - Payment Service: Separate PostgreSQL instance
  - Cache: Redis for shared state
- [ ] Monitoring & logging:
  - Centralized logging (ELK or Datadog)
  - Service-level metrics (Prometheus)
  - Distributed tracing (Jaeger)

### Phase 8d: Deployment & Testing (Week 19)
- [ ] Integration tests between services
- [ ] Load tests on new architecture
- [ ] Chaos engineering tests (what if Service X goes down?)
- [ ] Gradual rollout (canary deployment)
- [ ] Monitoring and alerting
- [ ] Runbook for troubleshooting

### Timeline & Dependencies
```
Current: Monolith
└── Phase 8a: Design (1 week)
    └── Phase 8b: Extract Services (3 weeks)
        └── Phase 8c: Infrastructure (1 week)
            └── Phase 8d: Deploy (1 week)
            └── Future: New services easily added
```

### Benefits After Migration
- ✅ Auth Service can scale independently
- ✅ Payment Service isolated from main app
- ✅ Notifications sent asynchronously
- ✅ New team can work on Expense Service simultaneously
- ✅ Easier to adopt new technologies per service
- ✅ Faster deployments (smaller services = faster builds)

**Deliverable:** Scalable microservices architecture ready for growth

---

## Phase 9: Production & Deployment (Week 20)
**Goal:** Deploy to production environment

### Database & Infrastructure
- [ ] 🔴 Create PostgreSQL database on production database provider
- [ ] 🔴 Write database migration scripts
- [ ] 🔴 Test connection strings
- [ ] 🔴 Set up automated backups
- [ ] 🟠 Configure read replicas (for scale)

### Docker & Deployment
- [ ] 🔴 Create Dockerfile for Go app
- [ ] 🔴 Test Docker image locally
- [ ] 🔴 Deploy to Railway.app or AWS
- [ ] 🟠 Set up load balancer
- [ ] 🟠 Configure auto-scaling

### CI/CD Pipeline
- [ ] 🔴 GitHub Actions: run tests on every push
- [ ] 🔴 GitHub Actions: build Docker image
- [ ] 🔴 GitHub Actions: deploy on main branch push
- [ ] 🟠 Pre-deployment checks (lint, security scan)

### Monitoring & Logging
- [ ] 🟠 Set up Prometheus metrics dashboard
- [ ] 🟠 Configure alerts (high error rate, slow response time)
- [ ] 🟠 Centralized logging (Datadog or ELK)
- [ ] 🟠 Set up uptime monitoring

### Security & Compliance
- [ ] 🔴 Enable HTTPS/SSL certificates
- [ ] 🔴 Set up CORS properly (allow only your frontend domain)
- [ ] 🔴 Rate limiting: 100 requests per IP per minute
- [ ] 🟠 Security headers: X-Frame-Options, X-Content-Type-Options
- [ ] 🟠 Input sanitization review
- [ ] 🟠 Penetration testing

### Performance Optimization
- [ ] 🟠 Add database query indexes
- [ ] 🟠 Implement caching (Redis) for community feed
- [ ] 🟠 Lazy load images on community page
- [ ] 🟠 Minify CSS/JS files

**Deliverable:** Production-ready deployment with monitoring and security

---

## Phase 10: Analytics & Growth (Ongoing)
**Goal:** Understand user behavior and drive growth

### Analytics
- [ ] 🟡 Track funnel: signup → create trip → publish → receive like
- [ ] 🟡 Event tracking: which features are used most?
- [ ] 🟡 Cohort analysis: retention of newly registered users
- [ ] 🟡 Geographic breakdown: which cities most popular?

### Growth Hacking
- [ ] 🟡 Email sequences: onboarding, re-engagement, abandoned draft
- [ ] 🟡 Referral program: "Invite friends, get ₹500 credit"
- [ ] 🟡 Featured trips: highlight high-engagement trips on homepage
- [ ] 🟡 Creator program: Verified badge for prolific trip creators

### Retention & Engagement
- [ ] 🟡 Weekly digest: "5 trending trips in your favorite destination"
- [ ] 🟡 Push notifications: "New trip matching your interests"
- [ ] 🟡 Gamification: Badges for milestones (100 likes, 10 trips published)
- [ ] 🟡 Leaderboard: Top creators by likes/followers

**Deliverable:** Data-driven growth strategy and user engagement system

---

## 🔧 Additional Technical Tasks (Ongoing)

### Code Quality
- [ ] 🟡 Write unit tests (at least 80% coverage)
- [ ] 🟡 Implement code review process
- [ ] 🟡 Set up linting (golangci-lint)
- [ ] 🟡 Document API endpoints with swagger/OpenAPI

### Database
- [ ] 🟠 Add missing indexes for common queries
- [ ] 🟠 Performance tuning: EXPLAIN ANALYZE queries
- [ ] 🟠 Regular backups verification
- [ ] 🟡 Data migration from SQLite to PostgreSQL

### Frontend
- [ ] 🟠 Responsive design review (mobile, tablet, desktop)
- [ ] 🟠 Accessibility audit (WCAG 2.1 AA)
- [ ] 🟠 Browser compatibility testing
- [ ] 🟡 Migrate to React for better maintainability

### Documentation
- [ ] 🟡 API documentation (Swagger/OpenAPI spec)
- [ ] 🟡 Deployment runbook
- [ ] 🟡 Troubleshooting guide
- [ ] 🟡 Architecture decision records (ADRs)

---

## 📋 Task Dependencies Map (Updated with New Features)

```
┌─────────────────────┐
│  Phase 0: Testing   │  (Week 1-2)
│  & Validation       │
└──────────┬──────────┘
           │
    ┌──────▼──────────────────────────────┐
    │  Phase 1-4: Core Trip Features  │  (Week 2-6)
    │  • Trip creation                │
    │  • Photo upload                 │
    │  • Reviews & ratings            │
    │  • Community posts              │
    └──────┬───────────────────────────┘
           │
    ┌──────▼──────────────────────────────┐
    │  Phase 5: Payments (REQUIRED)    │  (Week 7-8)
    └──────┬───────────────────────────┘
           │
    ┌──────▴──────────────────────────────────────────┐
    │  DECISION POINT: Launch or Add Features?      │
    │                                               │
    │  ├─ PATH A (Basic MVP → Production)          │
    │  │  └─ Phase 9: Deploy (Week 10)             │
    │  │                                            │
    │  └─ PATH B (Enhanced MVP with new features)  │
    │     ├─ Phase 5A: Group Trips (Week 8-9)  ──┐│
    │     ├─ Phase 5B: UI Polish (Week 8-9)   ──┼→ Can run in parallel
    │     ├─ Phase 6: React Frontend (Week 10-12)││
    │     ├─ Phase 7: AI Features (Week 13+)     ││
    │     ├─ Phase 8: Microservices (Week 14+)   ││
    │     └─ Phase 9: Deploy (Week 20)           ││
    │                                            ││
    └────────────────────────────────────────────┘│
                                                  │
    ┌─────────────────────────────────────────────┘
    │
    ├─ Phase 10: Analytics & Growth (Ongoing)
    └─ Future: Mobile App, Advanced Analytics
```

**Three Implementation Paths:**

**1️⃣ FAST PATH (Basic MVP - 8 weeks)**
- Phases 0-5, then deploy
- Market test, get feedback
- Add features based on demand

**2️⃣ FEATURE-COMPLETE PATH (12 weeks)**
- Phases 0-5, then add 5A, 5B, & 6
- Launch with modern UX
- Group trips & voting included
- React frontend for smooth animations

**3️⃣ FULL STACK PATH (20+ weeks)**
- All phases including microservices
- Enterprise-ready architecture
- Requires 2-3 developers
- Highest complexity, highest scalability

---

## 🎯 MVP Checklist (Minimum to Launch - Option A: Basic MVP)

**Basic MVP (Weeks 1-8, ~365 hours, Core Features Only):**
To launch MVP, these items MUST be done:

- ✅ Phase 0: Core flow works (login → dashboard → wizard)
- ✅ Phase 1: Users can create trips with segments
- ✅ Phase 2: Photo upload works
- ✅ Phase 3: Reviews & ratings functional
- ✅ Phase 5: Community posts visible
- ✅ Phase 5: Like system working
- ✅ Phase 7: Payment integration (at least order creation)
- ✅ Phase 9: Deployed to production

**Not in Basic MVP (add later):**
- Comments system
- Group trips
- React frontend
- Microservices
- AI features

---

## 🎯 Enhanced MVP Checklist (Option B: Feature-Complete MVP - including user requests)

**Enhanced MVP (Weeks 1-12, ~485 hours, With User-Requested Features):**

**MUST HAVE:**
- ✅ Phase 0: Core flow works
- ✅ Phase 1: Trip creation
- ✅ Phase 2: Photo upload
- ✅ Phase 3: Reviews & ratings
- ✅ Phase 5: Community posts & likes
- ✅ Phase 5A: Group trips & voting (friends can join)
- ✅ Phase 5A: Shared expense tracking & settlement
- ✅ Phase 5B: Modern UI with stock photos
- ✅ Phase 5B: CSS animations & polish
- ✅ Phase 6: React frontend migration
- ✅ Phase 7: Payment integration working
- ✅ Phase 9: Deployed to production

**NOT in Enhanced MVP (post-launch):**
- React is FULL rewrite (might be feature-complete, but polish later)
- Microservices (use monolith initially)
- AI trip generation
- Advanced notifications
- Mobile app

---

### Recommendation: Start with Basic MVP (8-10 weeks)

Then:
1. **Launch to beta users** (get feedback)
2. **Add Phase 5A**: Group trips & voting (high user demand)
3. **Add Phase 5B**: UI polish & photos (improves perception)
4. **Evaluate**: If animations/performance poor → do React migration
5. **Post-MVP**: Plan microservices when complexity grows

---

## 📊 Effort Estimation (Including New Features)

| Phase | Tasks | Est. Hours | Priority | Timeline |
|-------|-------|-----------|----------|----------|
| 0 | 18 | 20 | 🔴 Critical | Week 1-2 |
| 1 | 17 | 35 | 🔴 Critical | Week 2-3 |
| 2 | 11 | 25 | 🔴 Critical | Week 4 |
| 3 | 11 | 18 | 🔴 Critical | Week 4 |
| 4 | 15 | 30 | 🟠 High | Week 5-6 |
| 5 | 13 | 40 | 🟠 High | Week 7-8 |
| **5A (NEW)** | **25** | **60** | 🟠 High | Week 8-9 |
| **5B (NEW)** | **20** | **45** | 🟠 High | Week 8-9 |
| **6 (React)** | **35** | **120** | 🟠 High | Week 10-12 |
| 7 (AI) | 8 | 30 | 🟡 Medium | Week 13+ |
| **8 (Microservices)** | **40** | **100** | 🟡 Medium | Week 14+ |
| 9 | 12 | 35 | 🔴 Critical | Week 20+ |
| 10 | 8 | 20 | 🟢 Low | Ongoing |
| **Misc** | 8 | 15 | 🟡 Medium | Ongoing |
| **TOTAL** | **241** | **633 hours** | | |

**Timeline Options:**
- **MVP (to Phase 5):** 8-10 weeks, ~320 hours (1 dev)
- **MVP + UI Polish (to Phase 5B):** 9-10 weeks, ~365 hours (1 dev)
- **MVP + React Frontend (to Phase 6):** 12-14 weeks, ~485 hours (could be 2 devs)
- **MVP + React + Microservices (to Phase 8):** 20+ weeks, ~585 hours (recommend 2-3 devs)

**Recommended Team Structure:**
- 1 Backend Dev (full 20 weeks)
- 1 Frontend Dev (weeks 5+ for React)
- 1 DevOps/Infrastructure (weeks 14+ for microservices)

---

## 🏁 Success Criteria

- ✅ Users can complete login → trip creation → publish flow without errors
- ✅ 99% of API endpoints return 2xx or expected error codes
- ✅ Zero data loss in payment processing
- ✅ Photos upload and display correctly
- ✅ Community feed loads in < 2 seconds
- ✅ Mobile responsive (works on phones)
- ✅ 0 critical security vulnerabilities
- ✅ 100+ users in beta

---

**Last Updated:** March 23, 2026  
**Status:** Ready for implementation  
**Next Action:** Start Phase 0 testing

