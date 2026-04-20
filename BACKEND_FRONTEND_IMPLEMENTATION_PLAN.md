# Backend & Frontend Implementation Plan

## Backend API Endpoints Created ✅

### Cities/Destinations
- `GET /api/cities` - List all cities with pagination
- `GET /api/cities/:cityId/trip-posts` - Get trip posts (community posts) for a city

### Trip Community Posts
- `POST /api/user-trips/add-from-post` - Add a published trip post to user's itinerary
- `POST /api/trip-segments/:segmentId/mark-visited` - Mark a place as visited
- `POST /api/reviews` - Submit a review (auto-posts to trip)

### Supporting Endpoints
- `GET /api/user-trips` - List user's trips
- `POST /api/user-trips` - Create new trip
- `POST /api/user-trips/:id/segments` - Add segment to trip
- `POST /api/user-trips/:id/publish` - Publish trip as community post

## Backend Models Updated ✅

**TripSegment** - Added fields:
- `TimeOfDay` - "morning", "afternoon", "evening", "night"
- `Expense` - Cost of the place
- `BestTimeToVisit` - Best season/time to visit
- Fields: name, location, map coordinates, review, photos

**UserTripPost** - Enhanced:
- `DestinationID` - Which city the trip is for
- `Segments` - All places/activities in the trip
- `TotalExpense` - Total cost
- Response includes structured places data with reviews

**Request/Response Types**:
- `AddTripPostToItineraryRequest`
- `MarkSegmentVisitedRequest`
- `SubmitReviewRequest`
- `TripPostResponse` - Formatted response for API clients
- `TripPlaceResponse` - Formatted place data

## Frontend Pages & Components Needed

### 1. **Login Page** ✅
- Only login form (no signup/OAuth buttons)
- Simple email/password form
- "Forgot password" link (future)
- Submit to POST /auth/login

### 2. **Dashboard Page** (after login)
**URL**: `/dashboard`
**Components**:
- Header with user profile & logout
- Search/filter bar for cities
- **City Grid/List**: Cards showing all cities
  - City name, image, tagline
  - Click → Navigate to city page
  - GET /api/cities

### 3. **City Page** (view trip posts)
**URL**: `/city/:cityId`
**Components**:
- City header (name, description, image)
- **Trip Post Feed**: List of community posts
  - GET `/api/cities/:cityId/trip-posts`
  - Each post card shows:
    - Title, cover image
    - Trip duration (X days)
    - Total expense
    - Number of places/segments
    - User profile (optional)
    - "View Details" button
- Pagination controls

### 4. **Trip Post Detail Page** (view all places in a post)
**URL**: `/trip-posts/:postId`
**Components**:
- **Post Header**:
  - Title, description
  - Duration, total cost
  - Author info
- **Places List** (chronological by day/time):
  - For each place show:
    - **Place Card**:
      - Name, type (stay/food/activity)
      - Day & time of day (Morning, Lunch, Evening, Night)
      - Location with map preview (Google Maps embed)
      - Expense/cost
      - Best time to visit
      - Photos/images carousel
      - User's review & rating (if exists)
      - "View on Map" link
- **CTA Button**: "Add This Trip to My Itinerary"
  - POST `/api/user-trips/add-from-post` with tripPostId
  - Redirects to user's itinerary page

### 5. **My Itinerary Page** (user's trip planning)
**URL**: `/my-itinerary/:tripId`
**Components**:
- **Trip Header**:
  - Title, destination, duration
  - Total budget
  - Status badge (Draft, Planning, Ongoing, Completed)
- **Itinerary Timeline** (grouped by day/time):
  - Each day section showing:
    - Morning places
    - Lunch places
    - Afternoon/Evening places
    - Night places
- **Place Card** (in itinerary):
  - Name, location, expense
  - Status: Pending / Visited / Completed
  - **Button**: "Mark as Visited" (if pending)
    - POST `/api/trip-segments/:segmentId/mark-visited`
- **Actions**:
  - Add another place
  - Edit place details
  - Remove place
  - Publish trip

### 6. **Review Prompt Modal** (after marking visited)
**Trigger**: After POST `/api/trip-segments/:segmentId/mark-visited`
**Components**:
- Place name & photo
- **Review Form**:
  - Star rating (1-5)
  - Text review/comment field
  - Submit button
- **Submit**:
  - POST `/api/reviews` with rating, review, segmentId
  - Success message: "Review submitted! It's now part of your trip post"
- Auto-close after 2 seconds or manual close

### 7. **My Trips Page** (list all user trips)
**URL**: `/my-trips`
**Components**:
- **Trip Cards**:
  - Trip title, destination
  - Status badge
  - Number of places
  - Start date
  - Click to view details
  - GET `/api/user-trips`
- **Tab Navigation**:
  - All Trips
  - Planning
  - Ongoing
  - Completed
- **Button**: "Create New Trip"

### 8. **Community Feed Page** (optional)
**URL**: `/community`
**Components**:
- Recent published trip posts from all users
- Filter by destination
- Search trips
- Sort by likes, recent, trending

---

## Frontend User Flow

```
Login Page
    ↓
Dashboard (Select City)
    ↓
City Page (View Trip Posts)
    ↓
Trip Post Detail (See All Places)
    ↓
Add to Itinerary
    ↓
My Itinerary Page
    ↓
[Mark Place as Visited]
    ↓
Review Prompt Modal
    ↓
Review Submitted
    ↓
Continue Trip Planning
    ↓
Publish Trip as Post
    ↓
Back to Dashboard (Others can see your post)
```

---

## Frontend Technologies

- **Framework**: React / Vue (existing)
- **State Management**: Redux / Vuex / Context API
- **HTTP Client**: Axios / Fetch
- **Maps**: Google Maps API embed
- **Styling**: Tailwind CSS / Bootstrap (existing)
- **Forms**: React Hook Form / Formik
- **Toast/Alerts**: Toastr / SweetAlert
- **Image Carousel**: Swiper / Slick Carousel

---

## Database Schema Updates Needed

**Tables to modify/add**:

1. `trip_segments` - Add fields:
   - `time_of_day` VARCHAR(20) - morning, afternoon, evening, night
   - `expense` DECIMAL(10,2)
   - `best_time_to_visit` TEXT
   - `completed` BOOLEAN DEFAULT false

2. `user_trip_posts` - Add field:
   - `destination_id` UUID (foreign key to destinations)
   - `total_expense` DECIMAL(10,2)

3. New index:
   - `idx_user_trip_posts_destination` on (destination_id, published)
   - `idx_trip_segments_trip_id` on (user_trip_id)

---

## Next Steps (Priority Order)

**Backend**:
1. ✅ Add API endpoints (Done)
2. ✅ Update models (Done)
3. Add missing DB methods:
   - `GetTripSegmentByID()`
   - `UpdateTripSegmentCompletion()`
   - `GetTripPostsByDestination()` (optimized query)
4. Test all endpoints with Postman

**Frontend**:
1. Create Login page component
2. Create Dashboard page (city list)
3. Create City page (trip posts feed)
4. Create Trip Post Detail page
5. Create My Itinerary page
6. Create Review Prompt modal
7. Integrate all API calls
8. Add responsive design
9. Add error handling & loading states
10. Test user flows

