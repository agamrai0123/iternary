# Itinerary Hub MVP Enhancement Roadmap

## Phase 1: Community Posts System (Reddit-like)

### 1.1 New Features

#### Posts/Reviews Model
- **Post Types:**
  - Destination Reviews (tips, experience, what I liked/disliked)
  - Itinerary Feedback (suggestions, modifications)
  - Travel Stories (personal experiences, photos)

```go
type Post struct {
  ID                string
  UserID            string
  Type              string              // "destination_review", "itinerary_feedback", "travel_story"
  DestinationID     *string            // Optional - for destination posts
  ItineraryID       *string            // Optional - for itinerary posts
  Title             string
  Content           string
  Images            []PostImage
  BulletPoints      []string           // Key takeaways, tips
  Likes             int
  Comments          []PostComment
  WhatILiked        []string           // Positive experiences
  WhatIDidntLike    []string           // Negative experiences  
  Suggestions       []string           // Improvements & tips
  Tags              []string           // budget, adventure, family, etc
  CreatedAt         time.Time
  UpdatedAt         time.Time
}

type PostComment struct {
  ID        string
  UserID    string
  Content   string
  Likes     int
  CreatedAt time.Time
}
```

#### Image Upload System
- Support multiple images per post
- Image compression/optimization
- Gallery view

#### Engagement Features
- Like/Unlike posts
- Comment on posts
- Share posts
- Save posts to reading list

---

## Phase 2: Enhanced Customizable Travel Plans

### 2.1 User Travel Plans (Fork & Customize)

When a user "adds to travel plans", create an editable copy:

```go
type UserTravelPlan struct {
  ID               string
  UserID           string
  OriginalItineraryID  string
  Title            string
  Description      string
  Status           string              // "draft", "planned", "ongoing", "completed"
  
  // Customizations
  StartDate        *time.Time
  EndDate          *time.Time
  CustomItems      []CustomPlanItem   // User-added/modified items
  Budget           float64            // User's budget (can override)
  Notes            string             // Personal notes
  Collaborators    []string           // Share with friends
  
  CreatedAt        time.Time
  UpdatedAt        time.Time
}

type CustomPlanItem struct {
  ID               string
  Day              int
  Type             string             // stay, food, activity, transport
  Name             string
  Description      string
  Price            float64
  Duration         int
  Location         string
  BookingURL       string
  Status           string             // not_booked, pending, confirmed
  Notes            string
  IsFromOriginal   bool               // True if from original itinerary
}
```

---

## Phase 3: UX/UI Improvements

### 3.1 Discovery & Browsing
- Feed-style homepage showing popular posts
- Trending destinations/experiences
- Filter by tags (budget, adventure, family, romantic, etc)
- Sort by rating, recency, most helpful

### 3.2 Post Creation Interface
- Rich text editor with markdown support
- Image upload with drag-and-drop
- Structured input:
  - What I Liked (checklist)
  - What I Didn't Like (checklist)
  - Suggestions (add multiple)
  - Bullet points generator
  - Budget breakdown

### 3.3 Itinerary Customization
- Drag-and-drop day reordering
- Add/remove days easily
- Inline editing of items
- Cost calculator
- Timeline view with booking status

### 3.4 Social Features
- User profiles with contribution history
- "Helpful" voting on posts
- Following users/destinations
- Notifications for replies

---

## Implementation Priority

### High Priority (MVP Enhancement - Week 1-2)
1. ✅ Posts/Reviews model and API
2. ✅ Image upload (basic)
3. ✅ Like/Comment on posts
4. ✅ User travel plans (fork & save)
5. ✅ Basic customization UI

### Medium Priority (Polish - Week 3)
1. Rich text editor
2. Better image gallery
3. User profiles
4. Trending/Popular section
5. Search improvements

### Lower Priority (Future)
1. Collaborative planning
2. Social sharing
3. Notifications system
4. Mobile app
5. AI recommendations

---

## Database Schema Changes

### New Tables Required

```sql
-- Posts/Reviews
CREATE TABLE posts (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  type TEXT NOT NULL,
  destination_id TEXT,
  itinerary_id TEXT,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  likes INTEGER DEFAULT 0,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id),
  FOREIGN KEY (itinerary_id) REFERENCES itineraries(id)
);

-- Post Images
CREATE TABLE post_images (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  image_url TEXT NOT NULL,
  display_order INTEGER,
  created_at TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- Post Metadata (likes, dislikes, suggestions)
CREATE TABLE post_metadata (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  type TEXT NOT NULL,  -- "liked", "disliked", "suggestion"
  content TEXT NOT NULL,
  display_order INTEGER,
  created_at TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- Post Comments
CREATE TABLE post_comments (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  content TEXT NOT NULL,
  likes INTEGER DEFAULT 0,
  created_at TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- User Travel Plans (personalized copies)
CREATE TABLE user_travel_plans (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  original_itinerary_id TEXT,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'draft',
  start_date DATE,
  end_date DATE,
  budget REAL,
  notes TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (original_itinerary_id) REFERENCES itineraries(id)
);

-- Custom Plan Items (personalized day-by-day plan)
CREATE TABLE custom_plan_items (
  id TEXT PRIMARY KEY,
  plan_id TEXT NOT NULL,
  day INTEGER NOT NULL,
  type TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  price REAL,
  duration INTEGER,
  location TEXT,
  booking_url TEXT,
  status TEXT DEFAULT 'not_booked',
  notes TEXT,
  is_from_original BOOLEAN DEFAULT 1,
  display_order INTEGER,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (plan_id) REFERENCES user_travel_plans(id) ON DELETE CASCADE
);

-- Post Tags
CREATE TABLE post_tags (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL,
  tag_name TEXT NOT NULL,
  created_at TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
```

---

## API Endpoints to Add

### Posts
- `POST /api/posts` - Create post
- `GET /api/posts` - List posts (feed)
- `GET /api/posts/:id` - Get post details
- `PUT /api/posts/:id` - Update post
- `DELETE /api/posts/:id` - Delete post
- `POST /api/posts/:id/like` - Like post
- `POST /api/posts/:id/comments` - Comment on post
- `GET /api/destinations/:id/posts` - Posts for destination
- `GET /api/posts/trending` - Trending posts

### User Travel Plans
- `POST /api/travel-plans` - Create plan (from itinerary)
- `GET /api/travel-plans` - List user's plans
- `GET /api/travel-plans/:id` - Get plan details
- `PUT /api/travel-plans/:id` - Update plan
- `DELETE /api/travel-plans/:id` - Delete plan
- `POST /api/travel-plans/:id/items` - Add item to plan
- `PUT /api/travel-plans/:id/items/:itemId` - Update item
- `DELETE /api/travel-plans/:id/items/:itemId` - Remove item
- `POST /api/travel-plans/:id/duplicate` - Duplicate a shared plan

### User Profiles
- `GET /api/users/:id/profile` - User profile
- `GET /api/users/:id/posts` - User's posts
- `GET /api/users/:id/plans` - User's travel plans
- `PUT /api/users/:id/profile` - Update profile

---

## Frontend Component Changes

### New Pages/Sections
1. **Feed/Discover Page** (`/discover`)
   - Filter by destination, budget, tags
   - Search posts
   - Sort by trending, recent, helpful

2. **Post Viewer** (`/posts/:id`)
   - Full post with images
   - Comments section
   - "Add to Plan" button

3. **Create Post** (`/posts/create`)
   - Rich editor
   - Image uploader
   - Structured input fields
   - Preview

4. **My Travel Plans** (`/plans`)
   - List of saved/forked plans
   - Status badges
   - Quick edit links

5. **Plan Editor** (`/plans/:id/edit`)
   - Day-by-day builder
   - Drag-drop ordering
   - Add custom items
   - Cost breakdown
   - Booking status tracking

6. **User Profile** (`/users/:id`)
   - User info
   - Posts created
   - Travel plans shared
   - Follower/Following

---

## Benefits to MVP

| Current | Enhanced |
|---------|----------|
| Browse itineraries | **Create community posts about destinations** |
| View pre-made plans | **Share personal experiences & tips** |
| Static plans | **Fully customizable personal travel plans** |
| One way info | **Two-way community feedback** |
| Limited context | **Real user insights with images** |
| Can't modify | **Fork, customize, and personalize anything** |

---

## Implementation Complexity

- **Posts System**: Medium (database + API + UI)
- **Customizable Plans**: Medium (new models + state management)
- **Image Upload**: Low-Medium (storage + UI)
- **Social Features**: Low (likes/comments pattern)
- **UI Improvements**: Medium (new pages + components)

**Estimated Effort**: 3-4 weeks for full Phase 1 & 2 implementation

---

## Technology Stack Recommendations

- **Image Storage**: AWS S3 / Azure Blob / Local with CDN
- **Rich Editor**: TinyMCE / Quill.js / Markdown Editor
- **Image Upload**: Dropzone.js / React Dropzone
- **State Management**: LocalStorage for drafts
- **Real-time**: WebSockets for notifications (future)

