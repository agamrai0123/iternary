# Triply Development Milestones & Visual Mockups

---

## 🎯 Milestone Overview

```
CURRENT → M1 → M2 → M3 → M4 → M5 → M6 → MVP LAUNCH
(Done)   (Testing) (MVP Core) (Community) (Bookings) (Advanced) (Deploy)
```

**Timeline:** 8-10 weeks for MVP, 12+ weeks with advanced features

---

## 🏁 CURRENT STATE (March 23, 2026)

**Completed:**
- ✅ Backend API implementation (all handlers, services, database methods)
- ✅ 3 HTML templates (login, dashboard, plan-trip wizard)
- ✅ Database schema (12 tables with relationships)
- ✅ Authentication middleware
- ✅ Compilation: 0 errors

**Status:** Ready for end-to-end testing

---

## 📌 Milestone 1: Core Flow Testing & Fixes (Week 1-2)

**Goal:** Ensure login → dashboard → wizard flow works perfectly

### Key Features
- [x] User login with token storage
- [x] Dashboard displaying cities
- [x] Navigation to plan-trip wizard
- [x] Error handling & user feedback
- [x] Authorization checks

### UI Flow Mockup

```
┌─────────────────────────────────────────────────────────────────┐
│ MILESTONE 1: CORE FLOW                                          │
├─────────────────────────────────────────────────────────────────┤

[1] LOGIN PAGE
┌──────────────────────────┐
│  🌍 TRIPLY              │
│  Travel Planning       │
├──────────────────────────┤
│                         │
│ Email: [email]      ✓   │
│ Password: [pass]    ✓   │
│                         │
│    [LOGIN BUTTON]   ✓   │
│                         │
│  Demo: traveler@...  ✓  │
│        password123   ✓  │
└──────────────────────────┘
         ↓
    (FLOW TEST)
         ↓
[2] DASHBOARD
┌──────────────────────────────────────┐
│ Welcome User!  [Logout] [My Trips] ✓ │
├──────────────────────────────────────┤
│ CHOOSE DESTINATION:                  │
│                                      │
│ [GOA] [BALI] [AGRA]              ✓  │
│ [MANALI] [DELHI] [OOTY]          ✓  │
│                                      │
│ RIGHT SIDEBAR:                       │
│ ┌──────────────────────────────────┐│
│ │ [GET STARTED PLANNING BUTTON] ✓ ││
│ │ → My Trips Link             ✓  ││
│ │ → Community Link            ✓  ││
│ │ → Profile Settings          ✓  ││
│ └──────────────────────────────────┘│
└──────────────────────────────────────┘
         ↓
    (FLOW TEST)
         ↓
[3] PLAN TRIP WIZARD - STEP 1
┌──────────────────────────────────────┐
│ STEP 1 > TRIP BASICS               │
├──────────────────────────────────────┤
│ Destination:  [GOA dropdown]     ✓  │
│ Budget (₹):   [50000]           ✓  │
│ Duration:     [5 days]          ✓  │
│ Title:        [Trip Title]      ✓  │
│ Description:  [Text area]       ✓  │
│                                  │
│      [PREVIOUS]  [NEXT →]       ✓  │
└──────────────────────────────────────┘

✓ = TESTED & WORKING
```

### Testing Checklist
- [ ] Verify login with correct credentials
- [ ] Verify login fails with wrong password
- [ ] Check token storage in localStorage
- [ ] Check token is sent in Authorization header
- [ ] Test logout clears token
- [ ] Test redirect flow after login
- [ ] Verify cities load on dashboard
- [ ] Test "Get Started" button navigation
- [ ] Check form validation on wizard step 1
- [ ] Test back/next navigation

### Database Checks
- [ ] User record created on login
- [ ] Session/token properly stored
- [ ] Destinations data present

### Deliverable
✅ User can complete login → dashboard → wizard_step1 without errors

---

---

## 📌 Milestone 2: Trip Creation & Storage (Week 2-3)

**Goal:** Users can create multi-segment trips and save to database

### Key Features
- [x] Create new trip
- [x] Add multiple places (segments)
- [x] Save draft
- [x] Load saved trips
- [x] Edit trips

### UI Flow Enhancement

```
┌─────────────────────────────────────────────────────────────────┐
│ MILESTONE 2: TRIP CREATION & STORAGE                            │
├─────────────────────────────────────────────────────────────────┤

[WIZARD CONTINUATION]

STEP 1: TRIP BASICS (same as M1)
   ↓ [NEXT]
   
STEP 2: ADD PLACES
┌──────────────────────────────────────┐
│ DAY 1: Arriving in Goa             │
├──────────────────────────────────────┤
│                                     │
│ 🏨 PLACE 1:                         │
│ ├─ Name: Hotel Oceanview      ✓    │
│ ├─ Type: Hotel          ▼     ✓    │
│ ├─ Location: Baga Beach, Goa  ✓    │
│ ├─ Notes: 5-star resort       ✓    │
│ └─ [REMOVE]                   ✓    │
│                                     │
│ 🍽️ PLACE 2:                         │
│ ├─ Name: Spice Garden         ✓    │
│ ├─ Type: Restaurant    ▼      ✓    │
│ ├─ Location: Colaba, Goa      ✓    │
│ ├─ Notes: Best seafood        ✓    │
│ └─ [REMOVE]                   ✓    │
│                                     │
│ DAY 2: [showing 2 places added]    │
│ DAY 3: [empty]                      │
│                                     │
│ [+ ADD PLACE] [+ ADD NEW DAY]  ✓   │
│                                     │
│  [PREVIOUS] [NEXT →]           ✓   │
└──────────────────────────────────────┘
   ↓ Saves to DB on each step [SAVING...]

STEP 3: PHOTOS & REVIEWS
┌──────────────────────────────────────┐
│ DAY 1 > HOTEL OCEANVIEW             │
├──────────────────────────────────────┤
│                                     │
│ Photos (1-3):                       │
│ [📷 Click to Upload v] (1/3)   ✓    │
│ [📷 Click to Upload v] (2/3)   ✓    │
│                                     │
│ Review:                             │
│ Rating: ⭐⭐⭐⭐ (4/5)         ✓    │
│ Comment: [Great room view!]    ✓    │
│                                     │
│ ───────────────────────────────────│
│                                     │
│ DAY 1 > SPICE GARDEN                │
│ Photos (0/3)                        │
│ Rating: ⭐⭐⭐ (3/5)          ✓    │
│ Comment: [Decent seafood]     ✓    │
│                                     │
│  [PREVIOUS] [NEXT →]           ✓   │
└──────────────────────────────────────┘

STEP 4: REVIEW & PUBLISH
┌──────────────────────────────────────┐
│ TRIP SUMMARY                        │
├──────────────────────────────────────┤
│                                     │
│ ✓ Title: "Goa Getaway"        ✓    │
│ ✓ Budget: ₹50,000             ✓    │
│ ✓ Duration: 5 days            ✓    │
│ ✓ Places: 12 total            ✓    │
│ ✓ Photos: 24 uploaded         ✓    │
│ ✓ Reviews: 12 submitted       ✓    │
│                                     │
│ [x] Make Public Post          ✓    │
│                                     │
│  [SAVE AS DRAFT] [PUBLISH]    ✓    │
└──────────────────────────────────────┘

MY TRIPS PAGE (LOADED FROM DB)
┌──────────────────────────────────────┐
│ MY TRIPS                    [+ NEW] │
├──────────────────────────────────────┤
│                                     │
│ 🟡 Goa Getaway          (DRAFT)    │
│ └─ 5 days | ₹50,000    [EDIT][X]  │
│                                     │
│ 🟢 Bali Paradise       (PUBLISHED)  │
│ └─ 7 days | ₹75,000    [VIEW][X]  │
│                                     │
│ ⚪ Rajasthan Tour      (ONGOING)   │
│ └─ 10 days | ₹100,000  [EDIT][X]  │
│                                     │
│  [Load More...]                    │
└──────────────────────────────────────┘

✓ = SAVED TO DATABASE & RETRIEVED
```

### Database Operations
- [ ] CreateUserTrip stored in `user_trips` table
- [ ] Each segment inserted into `trip_segments` table
- [ ] GetUserTrips retrieves all user's trips
- [ ] UpdateUserTrip modifies existing trip
- [ ] DeleteUserTrip removes entire trip + cascade deletes segments

### Testing Checklist
- [ ] Create trip with all 4 wizard steps
- [ ] Verify data saved to database
- [ ] List trips on "My Trips" page
- [ ] Edit trip and verify changes persisted
- [ ] Delete trip and verify cascade delete of segments
- [ ] Test draft vs published status
- [ ] Test day-bound validation (can't add segment for day 6 if duration is 5)

### Deliverable
✅ Users can create complete trips with multiple places and retrieve them from database

---

---

## 📌 Milestone 3: Photo Upload & Reviews (Week 4)

**Goal:** Users can upload photos and write reviews for each place

### Key Features
- [x] Photo upload with validation
- [x] 3-photo limit per place
- [x] Quick preview
- [x] 1-5 star ratings
- [x] Review text submission

### UI Enhancement

```
┌─────────────────────────────────────────────────────────────────┐
│ MILESTONE 3: RICH MEDIA & REVIEWS                               │
├─────────────────────────────────────────────────────────────────┤

STEP 3 ENHANCED: PHOTOS & REVIEWS

PHOTO UPLOAD PREVIEW
┌──────────────────────────────────────┐
│ DAY 1 > HOTEL OCEANVIEW             │
├──────────────────────────────────────┤
│                                     │
│ UPLOAD 1/3:                         │
│ ┌────────────────────────────────┐ │
│ │ [Image Preview Here]   📷      │ │
│ │ Uploaded: beach.jpg            │ │
│ │ Size: 2.3 MB                   │ │
│ │ Caption: Beach room view       │ │
│ │ [REMOVE] [DOWNLOAD]            │ │
│ └────────────────────────────────┘ │
│                                     │
│ UPLOAD 2/3:                         │
│ ┌────────────────────────────────┐ │
│ │ [Image Preview Here]   📷      │ │
│ │ Uploaded: bathroom.jpg         │ │
│ │ Size: 1.8 MB                   │ │
│ │ Caption: Spacious bathroom     │ │
│ │ [REMOVE] [DOWNLOAD]            │ │
│ └────────────────────────────────┘ │
│                                     │
│ UPLOAD 3/3:                         │
│ [📷 CLICK TO UPLOAD NEW PHOTO]      │
│                                     │
│ ─────────────────────────────────   │
│                                     │
│ RATING & REVIEW:                    │
│                                     │
│ How would you rate?                 │
│ ⭐⭐⭐⭐⭐           (click to rate)   │
│                                     │
│ Your review (optional):             │
│ ┌────────────────────────────────┐ │
│ │ Amazing hotel! Very clean,     │ │
│ │ friendly staff, great location │ │
│ │ by the beach. Would return!    │ │
│ │ [SAVE REVIEW]                  │ │
│ └────────────────────────────────┘ │
│                                     │
│  [PREVIOUS] [NEXT →]           ✓   │
└──────────────────────────────────────┘

COMMUNITY POST WITH PHOTOS
┌──────────────────────────────────────────────┐
│ GIVEAWAY POST (After Publishing)            │
├──────────────────────────────────────────────┤
│                                              │
│ 👤 @TravelBug23    5 days ago    [Follow]  │
│ ⭐ Goa Getaway    ⭐ 5 days • ₹50,000      │
│                                              │
│ ┌────────────────────────────────────────┐  │
│ │     [COVER IMAGE - Large]       📷     │  │
│ │     (Beach sunset)                     │  │
│ └────────────────────────────────────────┘  │
│                                              │
│ DAY 1 - ACCOMMODATION                       │
│ 🏨 Hotel Oceanview (5★)                     │
│                                              │
│ ┌──────────┐┌──────────┐┌──────────┐       │
│ │ Photo 1  ││ Photo 2  ││ Photo 3  │       │
│ │ [Beach]  ││[Bathroom]││[Lobby]   │       │
│ └──────────┘└──────────┘└──────────┘       │
│                                              │
│ Review: "Amazing hotel! Very clean..." ▶    │
│                                              │
│ ───────────────────────────────────────────  │
│                                              │
│ DAY 1 - DINNER                              │
│ 🍽️ Spice Garden (3★)                        │
│                                              │
│ ┌──────────┐┌──────────┐                    │
│ │ Food     ││ Ambiance │                    │
│ │ [Seafood]││[Sunset]  │                    │
│ └──────────┘└──────────┘                    │
│                                              │
│ Review: "Decent seafood, good sunset..." ▶  │
│                                              │
│                                              │
│ 👍 342 Likes    ▶ Browse More       Comments│
│                                              │
│ [❤️ Like]  [📋 Copy Plan]  [💬 Comment]    │
└──────────────────────────────────────────────┘

✓ = PHOTO UPLOAD WORKING
✓ = REVIEWS VISIBLE ON POSTS
```

### File Handling
- [ ] Multipart form parsing implemented
- [ ] File size validation (< 5MB)
- [ ] Format validation (jpg/png only)
- [ ] File stored with unique name
- [ ] Photo URL saved to database
- [ ] 3-photo limit enforced

### Database Updates
- [ ] Photos inserted into `trip_photos` table
- [ ] Reviews upserted into `trip_reviews` table
- [ ] Photo count reflected in segment queries
- [ ] Average rating calculated for community display

### Testing Checklist
- [ ] Upload valid image (< 5MB, jpg/png) ✓
- [ ] Upload too-large image (rejected) ✓
- [ ] Upload non-image file (rejected) ✓
- [ ] Try uploading 4th photo (rejected) ✓
- [ ] Save review with 1-5 star rating ✓
- [ ] Edit existing review ✓
- [ ] Photos display on community post ✓
- [ ] Ratings show correctly ✓

### Deliverable
✅ Rich media posts with photos and user ratings work end-to-end

---

---

## 📌 Milestone 4: Community Feed & Engagement (Week 5-6)

**Goal:** Users can publish trips, see community, like posts, and comment

### Key Features
- [x] Publish trip to community
- [x] Community feed browse
- [x] Like/unlike posts
- [x] Like counts update dynamically
- [x] Comments on posts
- [x] Trending/sorting

### UI Complete

```
┌─────────────────────────────────────────────────────────────────┐
│ MILESTONE 4: COMMUNITY & ENGAGEMENT                             │
├─────────────────────────────────────────────────────────────────┤

COMMUNITY FEED PAGE
┌──────────────────────────────────────────────────────────────────┐
│ ✈️ TRIPLY COMMUNITY                              [Search...] 🔍 │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│ Sort: [Newest ▼] [Trending] [Most Liked]              [Pages] │
│                                                                  │
│ POST 1                                                           │
│ ┌─────────────────────────────────────────────────────────────┐│
│ │ 👤 @TravelBug23            5 DAYS AGO          [Follow]  ⋯ ││
│ │ 🏆 Goa Getaway   ⭐ 5 days | ₹50,000           [Reports]   ││
│ │                                                            ││
│ │ ┌─────────────────────────────────────────────────────────┐││
│ │ │ [COVER IMAGE - Beach Sunset]                    📷      │││
│ │ └─────────────────────────────────────────────────────────┘││
│ │                                                            ││
│ │ "Planning a relaxing beach vacation with family. This     ││
│ │  was one of the best trips ever! Highly recommend Goa." ││
│ │                                                            ││
│ │ Stats: 👍 342 Likes | 💬 28 Comments | 👁️ 1,250 Views    ││
│ │                                                            ││
│ │ [❤️ LIKE]  [📝 COMMENT]  [📋 COPY PLAN]  [↗️ SHARE]      ││
│ │                                                            ││
│ │ ─────────────────────────────────────────────────────────  ││
│ │                                                            ││
│ │ 💬 COMMENTS (TOP 3):                                      ││
│ │                                                            ││
│ │ @Explorer_Dev: "This looks amazing! When are you going?"  ││
│ │ ❤️ 12 replies • 5 days ago                                ││
│ │                                                            ││
│ │ @BudgetVan: "Can you share the hotel booking link?"       ││
│ │ ❤️ 8 replies • 4 days ago                                 ││
│ │                                                            ││
│ │ [View All 28 Comments...]                                 ││
│ │                                                            ││
│ └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│ POST 2                                                           │
│ ┌─────────────────────────────────────────────────────────────┐│
│ │ 👤 @JetsAway               3 WEEKS AGO         [Follow]  ⋯ ││
│ │ 🏆 Bali Paradise   ⭐ 7 days | ₹75,000                    ││
│ │ 👍 521 LIKES (you haven't liked this)                      ││
│ │                                                            ││
│ │ ┌─────────────────────────────────────────────────────────┐││
│ │ │ [COVER IMAGE - Temple at Sunrise]                ││
│ │ └─────────────────────────────────────────────────────────┘││
│ │                                                            ││
│ │ "Epic Bali adventure with amazing sunrises..."            ││
│ │                                                            ││
│ │ [❤️ LIKE]  [📝 COMMENT]  [📋 COPY PLAN]  [↗️ SHARE]      ││
│ │                                                            ││
│ └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│ POST 3                                                           │
│ [...loading more posts...]                                      │
│                                                                  │
│ [LOAD MORE POSTS...]                                            │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

LIKE INTERACTION (Dynamic Update)
Before:                          After Clicking Like:
👍 342 Likes                     👍 343 Likes
[❤️ LIKE]                        [❤️ LIKED] (filled heart)

COMMENT MODAL
┌───────────────────────────────────────┐
│ COMMENTS ON "GOA GETAWAY"             │
├───────────────────────────────────────┤
│                                       │
│ Comment 1:                            │
│ 👤 @Explorer_Dev  5 days ago         │
│ "This looks amazing! When are you    │
│  going?"                              │
│ ❤️ 12                                 │
│                                       │
│ ▼ 5 REPLIES:                          │
│   └─ @TravelBug23: "Heading May 1!"  │
│   └─ @Explorer_Dev: "Count me in!"   │
│                                       │
│ Comment 2:                            │
│ 👤 @BudgetVan  4 days ago            │
│ "Can you share the hotel booking     │
│  link?"                               │
│ ❤️ 8                                  │
│                                       │
│ ─────────────────────────────────────│
│                                       │
│ YOUR COMMENT:                         │
│ ┌────────────────────────────────────┐
│ │ [Type your comment...]             │
│ │                            [POST]  │
│ └────────────────────────────────────┘
│                                       │
└───────────────────────────────────────┘

✓ = COMMUNITY FEATURES LIVE
✓ = REAL-TIME LIKE COUNTS
✓ = COMMENT THREADS
```

### Database Operations
- [ ] `GetCommunityPosts()` retrieves published trips
- [ ] Pagination working (20 posts per page)
- [ ] Likes stored in `likes` join table
- [ ] Like count increments/decrements
- [ ] Comments retrieved with user info
- [ ] Sorting: newest first (published_at DESC)

### Features to Implement
- [ ] Like/unlike endpoints
- [ ] Like count management
- [ ] Comment submission
- [ ] Comment deletion (author only)
- [ ] Comment sorting (newest first)
- [ ] Pagination for comments
- [ ] Notification email when liked
- [ ] Notification email when commented

### Testing Checklist
- [ ] Like post (count updates) ✓
- [ ] Unlike post (count decreases) ✓
- [ ] View community feed ✓
- [ ] Filter by destination ✓
- [ ] Sort by newest/trending ✓
- [ ] Comment on post ✓
- [ ] Delete own comment ✓
- [ ] Can't delete others' comments ✓

### Deliverable
✅ Full community feed with engagement (likes, comments, browse)

---

---

## 📌 Milestone 5: Booking Integration (Week 7-8)

**Goal:** Users can pay and book through Razorpay

### Key Features
- [x] Order creation
- [x] Payment processing
- [x] Order confirmation
- [x] Booking confirmation email

### UI Addition

```
┌─────────────────────────────────────────────────────────────────┐
│ MILESTONE 5: BOOKING & PAYMENTS                                 │
├─────────────────────────────────────────────────────────────────┤

TRIP DETAIL PAGE - BOOKING SECTION
┌──────────────────────────────────────────────────────────────────┐
│ GOA GETAWAY - Trip Summary                                       │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│ COST BREAKDOWN:                                                  │
│ ├─ Accommodation (5 hotels)         ₹ 25,000                   │
│ ├─ Food & Dining (12 meals)         ₹ 10,000                   │
│ ├─ Activities (3 activities)         ₹ 8,000                    │
│ ├─ Transport (taxis, flights)        ₹ 5,000                    │
│ └─ Misc (tips, souvenirs)            ₹ 2,000                    │
│                                                                  │
│ ─────────────────────────────────────────────                   │
│ TOTAL COST:            ₹ 50,000                                 │
│ ─────────────────────────────────────────────                   │
│                                                                  │
│ [📋 COPY TO MY PLAN]  [🛒 BOOK NOW]                            │
│                                                                  │
│                                                                  │
│ PAYMENT PAGE (After Clicking "BOOK NOW"):                       │
│ ┌────────────────────────────────────────────────────┐         │
│ │ RAZORPAY PAYMENT GATEWAY                          │         │
│ ├────────────────────────────────────────────────────┤         │
│ │                                                    │         │
│ │ Order #: ORD_123456789                            │         │
│ │ Amount: ₹ 50,000                                  │         │
│ │ Trip: Goa Getaway                                 │         │
│ │                                                    │         │
│ │ Payment Method:                                    │         │
│ │ ◉ Debit/Credit Card                              │         │
│ │ ○ UPI (Google Pay, PhonePe)                       │         │
│ │ ○ Net Banking                                      │         │
│ │ ○ Wallet                                           │         │
│ │                                                    │         │
│ │ [PROCEED TO PAYMENT]                              │         │
│ │                                                    │         │
│ │ Secure: 🔒 256-bit SSL Encrypted                 │         │
│ │ Powered by Razorpay                               │         │
│ │                                                    │         │
│ └────────────────────────────────────────────────────┘         │
│                                                                  │
│ PAYMENT SUCCESS PAGE:                                            │
│ ┌────────────────────────────────────────────────────┐         │
│ │ ✅ PAYMENT SUCCESSFUL!                            │         │
│ ├────────────────────────────────────────────────────┤         │
│ │                                                    │         │
│ │ Order Confirmed: ORD_123456789                    │         │
│ │ Amount Paid: ₹ 50,000                             │         │
│ │ Payment Method: Debit Card (ending 1234)          │         │
│ │ Date & Time: 23 Mar 2026, 14:32:15 IST            │         │
│ │                                                    │         │
│ │ Receipt sent to: traveler@example.com             │         │
│ │                                                    │         │
│ │ MY BOOKINGS:                                       │         │
│ │ ├─ Hotel Oceanview - Confirmation #HZ123         │         │
│ │ ├─ Spice Garden Restaurant - Reservation OK      │         │
│ │ ├─ Scenic Tours Activity - Ticket #AT456         │         │
│ │ └─ Flight Transportation - E-ticket Ready        │         │
│ │                                                    │         │
│ │ [📩 DOWNLOAD INVOICE]  [🔖 SAVE EYE]             │         │
│ │ [📄 MY BOOKINGS]       [↗️ SHARE]                 │         │
│ │                                                    │         │
│ └────────────────────────────────────────────────────┘         │
│                                                                  │
│ EMAIL RECEIVED:                                                  │
│ From: bookings@triply.app                                       │
│ Subject: ✅ Booking Confirmed - Goa Getaway                    │
│                                                                  │
│ Hi Traveler,                                                     │
│                                                                  │
│ Your trip booking is confirmed!                                 │
│ Trip: Goa Getaway | Amount: ₹50,000                            │
│ Travel Dates: May 1-5, 2026                                     │
│                                                                  │
│ Confirmations:                                                   │
│ • Hotel Oceanview (Conf: HZ123)                                │
│ • Spice Garden Restaurant (Timing: 19:00)                      │
│ • Scenic Tours (Pickup: 09:00)                                 │
│                                                                  │
│ All attachments and e-tickets are available in your profile.   │
│                                                                  │
│ Questions? Contact us at support@triply.app                    │
│                                                                  │
│ Happy travels! 🌍                                               │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

✓ = ORDERS CREATED & TRACKING
✓ = PAYMENTS PROCESSED VIA RAZORPAY
✓ = INSTANT CONFIRMATION
```

### Payment Flow
- [ ] Order creation endpoint implemented
- [ ] Order stored with status 'pending'
- [ ] Razorpay merchant credentials configured
- [ ] Payment verification webhook implemented
- [ ] Order status updated to 'completed' on payment success
- [ ] Order status updated to 'failed' on payment failure

### Email Integration
- [ ] Booking confirmation email sent
- [ ] Email template with trip details
- [ ] Affiliate booking links included
- [ ] E-tickets generated (if applicable)

### Testing Checklist
- [ ] Create order on "Book Now" click ✓
- [ ] Razorpay payment window opens ✓
- [ ] Complete test payment ✓
- [ ] Webhook receives payment verification ✓
- [ ] Order status updates to 'completed' ✓
- [ ] Confirmation email received ✓
- [ ] Test payment failure scenario ✓
- [ ] Retry payment option works ✓

### Deliverable
✅ Complete payment flow from booking to confirmation

---

---

## 📌 Milestone 5A: Group Trips & Collaborative Voting (Week 8-9)

**Goal:** Enable users to join trips, vote on places, and share expenses

### Key Features
- [x] Trip member management (invite friends)
- [x] Group voting on places (democratic decision-making)
- [x] Shared expense tracking (Splitwise-style)
- [x] Automatic settlement calculations
- [x] Real-time expense notifications

### UI Enhancements

```
┌──────────────────────────────────────────────────────────────────┐
│ MILESTONE 5A: GROUP TRIPS & VOTING                               │
├──────────────────────────────────────────────────────────────────┤

TRIP DETAIL - INVITE FRIENDS
┌──────────────────────────────────────────────────────────────────┐
│ GOA GETAWAY - Group Planning                                     │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│ Trip Members (3/5):                                              │
│ ├─ 👤 You (Organizer)              🔑                           │
│ ├─ 👤 John Doe            [Admin]  👑                           │
│ ├─ 👤 Jane Smith          [Editor] ✏️                            │
│                                                                  │
│ [+ INVITE MORE MEMBERS]  [👥 MANAGE ROLES]                     │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ COLLABORATIVE VOTING - DAY 1 DINNER                             │
│ "Which restaurant should we go to?"                             │
│                                                                  │
│ Option 1: Spice Garden (Indian)      👍 3 votes (You, John, AI) │
│ Option 2: Beach Shack (Seafood)      👍 2 votes (Jane, +1 more) │
│ Option 3: Italian Bistro              👍 1 vote (Not decided)   │
│                                                                  │
│ [Vote] [View Details] [Add Option]                              │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ SHARED EXPENSES                                                  │
│ Total Trip Budget: ₹50,000                                      │
│ Per-Person Split: ₹16,666.67 (3 people)                         │
│                                                                  │
│ Expenses:                                                        │
│ ├─ Hotel (₹25,000)       Paid by: You       ✓                  │
│ │  Your share: ₹8,333    John: ₹8,333      Jane: ₹8,334        │
│ │                                                               │
│ ├─ Food (₹15,000)        Paid by: John      ✓                  │
│ │  Your share: ₹5,000    John: paid ₹15k    Jane: ₹5,000       │
│ │                                                               │
│ └─ Activities (₹10,000)  Paid by: You       ✓                  │
│    Your share: ₹3,333    John: ₹3,333      Jane: ₹3,334        │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ SETTLEMENT                                                       │
│                                                                  │
│ You paid:    ₹35,000                                            │
│ You owe:     ₹16,666.67                                         │
│ Net:         ₹18,333.33 (owed TO you)                           │
│                                                                  │
│ Breakdown:                                                       │
│ ├─ John owes you: ₹6,666.67                                    │
│ └─ Jane owes you: ₹11,666.67                                   │
│                                                                  │
│ [📱 SEND PAYMENT LINKS]  [📊 DETAILED REPORT]                  │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ REAL-TIME NOTIFICATIONS:                                        │
│ "John added expense: Dinner - ₹8,500 (3 participants)"         │
│ → Your share: ₹2,833.33 • Settlement changed                   │
│                                                                  │
│ "Jane voted for Spice Garden (2 votes total)"                  │
│ → Final decision coming soon (1 more vote needed)              │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

✓ = GROUP COLLABORATION ENABLED
✓ = DEMOCRATIC VOTING SYSTEM
✓ = AUTOMATIC EXPENSE SPLITTING
```

### Database Changes
- New table: `trip_members` (user_id, trip_id, role, joined_at)
- New table: `segment_votes` (vote_id, segment_id, member_id, choice, voted_at)
- New table: `trip_expenses` (expense_id, trip_id, paid_by, amount, description)
- New table: `expense_shares` (share_id, expense_id, member_id, amount)

### API Endpoints (New)
- `POST /trips/{id}/members` - Add member to trip
- `DELETE /trips/{id}/members/{member_id}` - Remove member
- `POST /segments/{id}/vote` - Cast vote on segment choice
- `GET /trips/{id}/votes` - Get all votes for trip
- `POST /expenses` - Create shared expense
- `GET /trips/{id}/settlement` - Calculate who owes whom

### Testing Checklist
- [ ] Invite friend to trip ✓
- [ ] Friend accepts/joins trip ✓
- [ ] Vote on place options ✓
- [ ] Vote count updates in real-time ✓
- [ ] Add expense for group ✓
- [ ] Expense automatically splits among members ✓
- [ ] Settlement calculations correct ✓
- [ ] Notifications sent on changes ✓

### Deliverable
✅ Group trip collaboration with voting and expense splitting

---

---

## 📌 Milestone 5B: UI Enhancement & Stock Photos (Week 8-9)

**Goal:** Polish UI with animations and integrate stock photo library

### Key Features
- [x] Unsplash API integration
- [x] Stock photo search & selection
- [x] Hero banners with gradient overlays
- [x] CSS animations and transitions
- [x] Design system (colors, typography, spacing)
- [x] Image caching and optimization

### UI Mockup

```
┌──────────────────────────────────────────────────────────────────┐
│ MILESTONE 5B: POLISHED UI WITH STOCK PHOTOS                      │
├──────────────────────────────────────────────────────────────────┤

HOMEPAGE HERO
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│   [Animated Background: Carousel of Travel Photos]   ✨        │
│                                                                  │
│   ┌────────────────────────────────────────────────────────┐   │
│   │ Dark gradient overlay                                  │   │
│   │                                                        │   │
│   │ ✈️ TRIPLY                                            │   │
│   │ Plan. Share. Explore. Together.                      │   │
│   │                                                        │   │
│   │ [START PLANNING NOW]                                 │   │
│   │                                                        │   │
│   └────────────────────────────────────────────────────────┘   │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ FEATURED TRIPS CAROUSEL (with Smooth Scroll)                    │
│                                                                  │
│ « [Trip Card 1]  [Trip Card 2]  [Trip Card 3]  [Trip 4] »     │
│     ↓ (hover shows animation)                                   │
│   ┌──────────────────────────────────────────────────┐         │
│   │ [Hero Photo from Unsplash]  📷                   │         │
│   │ (Beautiful Goa Beach)                            │         │
│   │                                                  │         │
│   │ Fade-in Animation on Load                        │         │
│   │ Scale on Hover (+5%)                             │         │
│   │                                                  │         │
│   │ 🏆 Goa Getaway                                  │         │
│   │ ⭐ 4.8 | 5 days | ₹50,000                       │         │
│   │                                                  │         │
│   │ "Amazing beach trip with family"                │         │
│   │                                                  │         │
│   │ 👍 342 Likes                                     │         │
│   │ [❤️ LIKE]  [📋 COPY]                           │         │
│   └──────────────────────────────────────────────────┘         │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ DESTINATION CARDS WITH IMAGES                                   │
│                                                                  │
│ [GOA]        [BALI]       [AGRA]      [MANALI]                 │
│ │Beautiful  │ Tropical    │ Historic  │ Mountain               │
│ │beaches    │ paradise    │ temples   │ views                  │
│ │📷 Unsplash│📷 Unsplash │📷 Unsplash│📷 Unsplash            │
│ │           │             │           │                        │
│ └Popular✓   └Trending✓    └Iconic✓   └Adventure✓             │
│                                                                  │
│ ─────────────────────────────────────────────────────────────  │
│                                                                  │
│ ANIMATIONS SHOWCASE                                             │
│                                                                  │
│ ✓ Fade-in on scroll (intersection observer)                    │
│ ✓ Slide-left/right card transitions                            │
│ ✓ Scale hover effects on buttons                               │
│ ✓ Smooth color transitions on links                            │
│ ✓ Skeleton loaders while images load                           │
│ ✓ Progress animations (vote bars, expense bars)                │
│ ✓ Notification toast animations                                │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

DESIGN SYSTEM (CSS Variables)
──────────────────────────────
Color Palette:
  Primary:      #2563EB (Vibrant Blue)
  Secondary:    #F59E0B (Warm Orange)
  Accent:       #10B981 (Fresh Green)
  Danger:       #EF4444 (Red)
  Text:         #1F2937 (Dark Gray)
  Background:   #F9FAFB (Off White)

Typography:
  Display:      48px, Bold (headings)
  Large:        32px, Semibold (section titles)
  Base:         16px, Regular (body text)
  Small:        14px, Regular (captions)

Spacing Scale: 4px, 8px, 12px, 16px, 24px, 32px, 48px

Border Radius: 4px (small), 8px (medium), 12px (large), 50% (circle)

Shadows:
  Small:   0 1px 3px rgba(0,0,0,0.1)
  Medium:  0 4px 6px rgba(0,0,0,0.1)
  Large:   0 10px 15px rgba(0,0,0,0.1)

✓ = COHESIVE DESIGN SYSTEM
```

### Implementation Tasks
- [ ] Unsplash API credentials configured
- [ ] Photo search endpoint created
- [ ] Image optimization (compression, resizing)
- [ ] CDN integration for fast loading
- [ ] CSS animations for smooth UX
- [ ] Design tokens in CSS custom properties
- [ ] Loading states (skeleton screens)
- [ ] Lazy loading for images
- [ ] Image caching strategy

### Image Categories (Unsplash Integration)
- Destination hero images (beaches, mountains, temples)
- Activity photos (hiking, water sports, dining)
- Hotel/accommodation previews
- Food & restaurant ambiance
- User-generated content thumbnails

### Testing Checklist
- [ ] Unsplash search api working ✓
- [ ] Images load without CORS issues ✓
- [ ] Images cached properly ✓
- [ ] Animations smooth at 60fps (Lighthouse) ✓
- [ ] Mobile responsive design ✓
- [ ] Lazy loading working ✓
- [ ] Design consistent across pages ✓

### Deliverable
✅ Polished, modern UI with stunning visuals and smooth animations

---

---

## 📌 Milestone 6: React Frontend Migration (Week 10-12)

**Goal:** Rebuild frontend with React for better performance and UX

### Key Features
- [x] React component library setup
- [x] State management (Context API or Redux)
- [x] Page-by-page migration from vanilla JS
- [x] Real-time updates (WebSockets)
- [x] Advanced animations (Framer Motion)
- [x] Mobile-responsive design

### Architecture Planning

```
┌──────────────────────────────────────────────────────────────────┐
│ MILESTONE 6: REACT FRONTEND MIGRATION                            │
├──────────────────────────────────────────────────────────────────┤

PHASE 6A: REACT SETUP & TOOLING (Week 10)
──────────────────────────────────────────
├─ Initialize Next.js with TypeScript
├─ Set up Vite for fast HMR
├─ Install component library (Material-UI, Shadcn/ui)
├─ Configure state management (Context + Redux)
├─ Set up API client (Axios/Fetch with retries)
├─ Configure ESLint + Prettier
├─ Set up Jest + React Testing Library
└─ GitHub Actions CI/CD for React builds

Technology Stack:
  Framework:        Next.js 14 (React 18)
  Language:         TypeScript
  Styling:          Tailwind CSS + CSS Modules
  State:            Zustand (lightweight) or Redux
  API Client:       Axios with interceptors
  Testing:          Vitest + React Testing Library
  Animations:       Framer Motion
  Forms:            React Hook Form + Zod validation
  HTTP:             TanStack Query (React Query)

PHASE 6B: COMPONENT LIBRARY (Week 10-11)
──────────────────────────────────────────
Shared Components:
├─ Button (variants: primary, secondary, danger)
├─ Input (text, email, password, number)
├─ Modal / Dialog
├─ Card / Container
├─ Navigation (Navbar, Sidebar)
├─ Footer
├─ Message (Alert, Toast, Snackbar)
├─ Loader / Skeleton
├─ Avatar / User Profile
├─ Rating / Star Component
├─ Pagination
├─ Tags / Chips
└─ Breadcrumbs

Page Components:
├─ AuthLayout (login, signup)
├─ DashboardLayout (main app)
├─ TripDetail / TripList
├─ CommunityFeed
├─ BookingFlow
├─ UserProfile
└─ Settings

PHASE 6C: PAGE MIGRATION (Week 11-12)
──────────────────────────────────────
Migration Strategy:
1. Keep Go backend running (no changes needed)
2. Build React pages using existing API
3. Deploy React to Vercel/Netlify (static host)
4. A/B test: 50% users on React, 50% on vanilla
5. Gradual rollout (90% → 100%) based on metrics
6. Decommission vanilla JS pages

Migration Order:
  Week 11:
  - Landing page (lowest risk)
  - Login/Signup pages
  - Dashboard (cities selection)
  
  Week 12:
  - Plan-Trip wizard (most complex)
  - Community feed
  - Trip detail page
  - Booking/payment flow

PHASE 6D: ADVANCED FEATURES (Week 12+)
───────────────────────────────────────
├─ Real-time updates via WebSockets
├─ Optimistic updates (instant UI feedback)
├─ Offline support (Service Workers)
├─ Push notifications
├─ Desktop version optimization
└─ Progressive Web App (PWA)

Migration Benefits:
  ✓ Faster perceived performance (SPA)
  ✓ Smoother animations (Framer Motion)
  ✓ Real-time collaboration (WebSockets)
  ✓ Better mobile experience (React Native future)
  ✓ Easier testing (component isolation)
  ✓ Better developer experience (HMR)
  ✓ Scalable architecture (monorepo ready)
```

### React Component Example
```typescript
// Example React Component
import React, { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';

const TripCard: React.FC<{ trip: Trip }> = ({ trip }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ scale: 1.05 }}
      className="card rounded-lg shadow-md p-4"
    >
      <img
        src={trip.heroImage}
        alt={trip.title}
        className="w-full h-48 object-cover rounded-md"
      />
      <h3 className="mt-4 text-xl font-bold">
        {trip.title}
      </h3>
      <p className="text-gray-600">
        ⭐ {trip.rating} | {trip.duration} days | ₹{trip.budget}
      </p>
      <button className="mt-4 btn btn-primary">
        View Trip
      </button>
    </motion.div>
  );
};
```

### Testing Checklist
- [ ] React app builds without errors ✓
- [ ] API calls working from React ✓
- [ ] Authentication persisting ✓
- [ ] Pages rendering correctly ✓
- [ ] Animations smooth (60fps) ✓
- [ ] Mobile responsive ✓
- [ ] Performance metrics good (Lighthouse) ✓
- [ ] A/B test metrics collected ✓

### Performance Targets
- Lighthouse Score: 90+
- First Contentful Paint: < 1.5s
- Time to Interactive: < 3s
- Cumulative Layout Shift: < 0.1

### Deliverable
✅ Modern React frontend with smooth UX and real-time features

---

---

## 📌 Milestone 8: Microservices Architecture (Week 14+)

**Goal:** Refactor monolith into scalable microservices

### Key Features
- [x] 7 independent services with separate databases
- [x] API Gateway for routing
- [x] Inter-service communication (gRPC/REST)
- [x] Event-driven architecture (message queues)
- [x] Service discovery and load balancing
- [x] Independent deployment and scaling

### Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│ MILESTONE 8: MICROSERVICES ARCHITECTURE                          │
├──────────────────────────────────────────────────────────────────┤

SYSTEM ARCHITECTURE
───────────────────

┌─ CLIENT (React App) ─────────────────────────────────────────┐
│                                                               │
│  Browser/Mobile Client                                       │
│  (Hosted on Vercel/Netlify)                                 │
└────────────────────────────────┬────────────────────────────┘
                                 │ HTTPS
                ┌────────────────▼────────────────┐
                │   🔐 API GATEWAY               │
                │   (Kong/Traefik/Express)       │
                │   - Rate limiting              │
                │   - Authentication             │
                │   - Request routing            │
                │   - Logging/Monitoring         │
                └────────┬─────┬────┬────┬───┬──┘
                         │     │    │    │   │
        ┌────────────────┴─┐   │    │    │   └──┬──────────────┐
        │                  │   │    │    │      │              │
        ▼                  ▼   ▼    ▼    │      ▼              ▼
    ┌─────────┐        ┌──────────┐    │  ┌────────┐     ┌──────────┐
    │ Auth    │        │Trip      │    │  │Payment │     │Notification
    │Service  │        │Service   │    │  │Service │     │Service
    ├─────────┤        ├──────────┤    │  ├────────┤     ├──────────┤
    │ Port    │        │ Port     │    │  │ Port   │     │ Port
    │ 3001    │        │ 3002     │    │  │ 3004   │     │ 3005
    │         │        │          │    │  │        │     │
    │PG DB:  │        │PG DB:   │    │  │PG DB: │     │PG DB:
    │auth_db │        │trip_db  │    │  │pay_db │     │notify_db
    └────┬────┘        └───┬─────┘    │  └───┬────┘     └────┬──────┘
         │                 │          │      │               │
         └────────┬────────┘          │      └───┬───────────┘
                  │                   │          │
        ┌─────────▼────────┐   ┌──────▼──┐   ┌──▼──────────────┐
        │ Message Queue    │   │Media    │   │Community Service
        │ (RabbitMQ/Kafka) │   │Service  │   ├──────────────────┤
        │ - Event stream   │   ├─────────┤   │ Port: 3006
        │ - Pub/Sub        │   │ Port    │   │
        └──────────────────┘   │ 3003    │   │PG DB:
                               │         │   │community_db
                               │S3/CDN   │   │
                               └─────────┘   └──────────────────┘

ALSO:
    - Expense Service (Port 3007, expense_db)
    - Discovery Service (Consul/Eureka)
    - Config Server (Centralized config)
    - Monitoring (Prometheus + Grafana)
    - Logging (ELK Stack or centralized logs)

SERVICE BREAKDOWN
─────────────────

1️⃣ AUTH SERVICE (Port 3001)
   Database: auth_db (PostgreSQL)
   Responsibilities:
   ├─ User registration & login
   ├─ Token generation (JWT)
   ├─ Token validation
   ├─ Session management
   ├─ Permission checks
   └─ OAuth integration
   
   API Endpoints:
   POST /auth/register       → Create user
   POST /auth/login          → Login & get token
   POST /auth/verify         → Verify token (internal)
   GET /auth/profile/{id}    → Get user profile
   PUT /auth/profile         → Update profile
   DELETE /auth/logout       → Logout

2️⃣ TRIP SERVICE (Port 3002)
   Database: trip_db (PostgreSQL)
   Responsibilities:
   ├─ Trip CRUD operations
   ├─ Segment management
   ├─ Trip cloning/remixing
   ├─ Draft vs published states
   └─ Trip history tracking
   
   API Endpoints:
   GET/POST /trips                → List/create trips
   GET/PUT/DELETE /trips/{id}    → Trip detail operations
   POST /trips/{id}/segments      → Add segment
   GET /trips/{id}/segments       → Get segments
   PUT /trips/{id}/publish        → Publish trip
   GET /community/trips           → Public trips (filtered)

3️⃣ MEDIA SERVICE (Port 3003)
   Storage: AWS S3 / GCS
   Database: media_db (PostgreSQL metadata)
   Responsibilities:
   ├─ Photo upload/storage
   ├─ Image optimization
   ├─ CDN integration
   ├─ Thumbnail generation
   ├─ Access control
   └─ Video streaming (future)
   
   API Endpoints:
   POST /media/upload/{trip-id}   → Upload photo
   GET /media/{photo-id}          → Download/retrieve
   DELETE /media/{photo-id}       → Delete photo
   POST /media/optimize           → Generate thumbnails
   GET /media/stats               → Storage usage

4️⃣ PAYMENT SERVICE (Port 3004)
   Database: payment_db (PostgreSQL)
   Responsibilities:
   ├─ Order creation
   ├─ Payment processing (Razorpay)
   ├─ Payment verification
   ├─ Refund handling
   ├─ Invoice generation
   └─ Subscription management
   
   API Endpoints:
   POST /payments/order           → Create order
   POST /payments/verify          → Verify payment
   POST /payments/refund          → Process refund
   GET /payments/history/{id}     → Payment history
   GET /payments/invoices/{id}    → Download invoice
   POST /payments/subscribe       → Create subscription

5️⃣ NOTIFICATION SERVICE (Port 3005)
   Database: notification_db (PostgreSQL)
   Responsibilities:
   ├─ Email notifications
   ├─ SMS notifications
   ├─ Push notifications
   ├─ In-app messaging
   ├─ Notification preferences
   └─ Scheduling
   
   API Endpoints:
   POST /notifications/send       → Send message
   GET /notifications/{user-id}   → Get notifications
   PUT /notifications/{id}/read   → Mark as read
   POST /notifications/preferences → Set preferences

6️⃣ COMMUNITY SERVICE (Port 3006)
   Database: community_db (PostgreSQL)
   Responsibilities:
   ├─ Posts/feed operations
   ├─ Comments & replies
   ├─ Likes/reactions
   ├─ Followers/following
   ├─ Trending algorithms
   └─ Search indexing
   
   API Endpoints:
   GET /feed                      → Community feed
   POST /posts/{id}/like          → Like a post
   DELETE /posts/{id}/like        → Unlike
   POST /posts/{id}/comment       → Add comment
   GET /posts/{id}/comments       → Get comments
   GET /trending                  → Trending trips

7️⃣ EXPENSE SERVICE (Port 3007)
   Database: expense_db (PostgreSQL)
   Responsibilities:
   ├─ Expense tracking
   ├─ Splitting logic
   ├─ Settlement calculations
   ├─ Payment reminders
   └─ Recurring expenses
   
   API Endpoints:
   POST /expenses                 → Create expense
   GET /trip/{id}/expenses        → Trip expenses
   POST /expenses/{id}/settle     → Settle payment
   GET /settlements/{user-id}     → Who owes whom
   PUT /expenses/{id}             → Update expense

DATA CONSISTENCY
────────────────
Challenges:
├─ Distributed transactions (2-phase commit avoided)
├─ Eventual consistency vs strong consistency
├─ Saga pattern for workflows
└─ Data duplication across services

Implementation:
├─ Event-driven architecture (message queue)
├─ Saga pattern for distributed transactions
├─ Change Data Capture (CDC) for sync
├─ Service-to-service retry logic
└─ Compensating transactions for rollbacks

DEPLOYMENT
──────────
Container Orchestration: Kubernetes (K8s)
├─ Docker images for each service
├─ Auto-scaling based on CPU/memory
├─ Rolling updates (zero downtime)
├─ Configuration management (ConfigMaps, Secrets)
├─ Service mesh (Istio/Linkerd) for advanced routing
└─ Monitoring & logging (Prometheus, ELK)

CI/CD Pipeline:
├─ Git push → GitHub Actions
├─ Build docker images for changed services
├─ Run tests in parallel
├─ Push to container registry
├─ Deploy to staging K8s cluster
├─ Smoke tests
├─ Deploy to production with canary rollout

SCALING STRATEGY
────────────────
With Microservices:
├─ Scale only busy services (e.g., Trip Service during peak)
├─ Separate databases prevent bottlenecks
├─ Caching layers (Redis) per service
├─ Message queues decouple dependencies
└─ Load balancers distribute traffic

Example:
  Auth Service: 2 pods (light traffic)
  Trip Service: 10 pods (heavy traffic during peak)
  Media Service: 5 pods (I/O intensive)
  Payment Service: 3 pods (critical, high availability)

✓ = ENTERPRISE-READY ARCHITECTURE
✓ = UNLIMITED HORIZONTAL SCALING
```

### Migration Phases
- **Phase 8A (Week 14):** Service design & planning, API contracts
- **Phase 8B (Week 15-16):** Extract services one by one (Auth → Trip → Media)
- **Phase 8C (Week 17-18):** Deploy with API Gateway, integrate services
- **Phase 8D (Week 19+):** Monitoring, optimization, scale to production

### Technology Stack
- **Service Framework:** Go (keep consistency), Node.js (rapid development)
- **Database:** PostgreSQL (separate DB per service)
- **Message Queue:** RabbitMQ or Kafka
- **API Gateway:** Kong, Traefik, or Express.js
- **Container:** Docker & Kubernetes
- **Service Mesh:** Istio (optional, for advanced routing)
- **Monitoring:** Prometheus + Grafana
- **Logging:** ELK Stack or Datadog

### Testing Checklist
- [ ] Services build and run independently ✓
- [ ] Inter-service communication working ✓
- [ ] API Gateway routing correctly ✓
- [ ] Message queue delivering events ✓
- [ ] Data consistency (eventual) maintained ✓
- [ ] Load testing at scale (1000+ QPS) ✓
- [ ] Chaos testing (service failures) ✓

### Deliverable
✅ Production-ready microservices architecture supporting 100K+ users

---

---

## 📌 Milestone 6 (Previous): AI & Advanced Features  [OPTIONAL]

**Goal:** Add AI-powered trip generation and recommendations

### Key Features
- [x] AI trip generator (Claude API)
- [x] Smart recommendations
- [x] Trip remixing
- [x] Price staleness alerts

### Example AI Interactions

```
┌──────────────────────────────────────────────────────┐
│ AI TRIP GENERATOR                                    │
├──────────────────────────────────────────────────────┤
│                                                      │
│ Help me plan!                                        │
│                                                      │
│ Destination: Goa                                     │
│ Budget: ₹ 50,000                                    │
│ Duration: 5 days                                     │
│                                                      │
│ [✨ GENERATE TRIP WITH AI]                          │
│                                                      │
│ ───────────────────────────────────────────────────  │
│                                                      │
│ 🤖 AI generating your perfect itinerary...         │
│                                                      │
│ ✓ Generated 12-place itinerary                      │
│ ✓ Balanced hotels, food, activities                │
│ ✓ Stays within ₹50,000 budget                      │
│ ✓ Optimal route for travel time                    │
│                                                      │
│ Generated Itinerary:                                 │
│ ├─ Day 1: Hotel + Sunset Beach Tour               │
│ ├─ Day 2: Water Sports + Beach Shack Dinner       │
│ ├─ Day 3: Local Village Tour + Market              │
│ ├─ Day 4: Spa + Portuguese Fort Exploration        │
│ └─ Day 5: Sunrise + Departure                       │
│                                                      │
│ Estimated Cost: ₹48,500 (within budget) ✓          │
│                                                      │
│ [✅ USE THIS ITINERARY] [❌ REGENERATE]             │
│ [🔧 CUSTOMIZE]                                      │
│                                                      │
│                                                      │
│ PRICE STALENESS ALERT:                              │
│ ⚠️  Hotel prices in Goa ↑ 12% since your trip      │
│     was published 30 days ago                        │
│                                                      │
│ Original: ₹3,500/night → Current: ₹3,920/night    │
│                                                      │
│ [⟳ UPDATE PRICES] [✓ ACKNOWLEDGE]                  │
│                                                      │
│                                                      │
│ SMART RECOMMENDATIONS:                              │
│ Because you visited Goa in May, you might enjoy:   │
│                                                      │
│ Similar: Kerala Backwaters (⭐ 4.8 rating)         │
│ Similar: Thai Islands (⭐ 4.9 rating)              │
│ Trending: Iceland in Summer (⭐ 4.7 rating)        │
│                                                      │
│ [VIEW TRIP] [ADD TO WISHLIST]                       │
│                                                      │
└──────────────────────────────────────────────────────┘

✓ = AI-POWERED RECOMMENDATIONS
✓ = AUTOMATED PRICE MONITORING
```

### Features
- [ ] Claude API integration
- [ ] Prompt engineering for trip generation
- [ ] Trip parsing from AI response
- [ ] Auto-creation of segments from generated itinerary
- [ ] Price comparison alerts
- [ ] Recommendation engine
- [ ] Trending destination detection

### Deliverable
✅ AI-generated trips and smart recommendations available

---

---

## 📌 Milestone 7: Production Deployment (Week 10)

**Goal:** Deploy to production with monitoring and security

### Infrastructure Setup
- [x] Docker containerization
- [x] PostgreSQL database
- [x] CI/CD pipeline (GitHub Actions)
- [x] Load balancing
- [x] HTTPS/SSL
- [x] Performance monitoring
- [x] Error tracking
- [x] Log aggregation

### Production Checklist
- [ ] Database migrated to PostgreSQL
- [ ] Environment variables secured
- [ ] All tests passing
- [ ] Load testing completed (1000 concurrent users)
- [ ] Security audit passed
- [ ] CORS configured
- [ ] Rate limiting enabled
- [ ] Backups automated
- [ ] CI/CD pipeline tested
- [ ] Monitoring dashboards set up
- [ ] Alert thresholds configured
- [ ] Documentation complete

### Deliverable
✅ Production-ready Triply platform live at https://triply.app

---

---

## 🎯 MVP Feature Comparison

```
                           M1    M2    M3    M4    M5    M6    M7
                          (Test)(Core)(Media)(Comm)(Pay) (AI)  (Prod)

LOGIN & AUTH              ✅    ✅    ✅    ✅    ✅    ✅    ✅
CREATE TRIP               ✅    ✅    ✅    ✅    ✅    ✅    ✅
ADD PLACES/SEGMENTS       ✅    ✅    ✅    ✅    ✅    ✅    ✅
PHOTO UPLOAD              ─     ─     ✅    ✅    ✅    ✅    ✅
RATINGS & REVIEWS         ─     ─     ✅    ✅    ✅    ✅    ✅
PUBLISH TO COMMUNITY      ─     ─     ✅    ✅    ✅    ✅    ✅
BROWSE COMMUNITY FEED     ─     ─     ✅    ✅    ✅    ✅    ✅
LIKE/UNLIKE POSTS         ─     ─     ─     ✅    ✅    ✅    ✅
COMMENTS ON POSTS         ─     ─     ─     ✅    ✅    ✅    ✅
BOOKING & PAYMENTS        ─     ─     ─     ─     ✅    ✅    ✅
AI TRIP GENERATION        ─     ─     ─     ─     ─     ✅    ✅
PRICE ALERTS              ─     ─     ─     ─     ─     ✅    ✅
NOTIFICATIONS             ─     ─     ─     ─     ─     ✅    ✅
PRODUCTION DEPLOY         ─     ─     ─     ─     ─     ─     ✅
```

---

---

## � Milestone 7: Production Deployment (Week 20)

**Goal:** Deploy final microservices architecture to production with monitoring

### Infrastructure Setup
- [x] Container orchestration (Kubernetes)
- [x] PostgreSQL databases (replicated)
- [x] Redis caching layer
- [x] Message queue (RabbitMQ/Kafka)
- [x] API Gateway configuration
- [x] CI/CD pipeline (GitHub Actions)
- [x] HTTPS/SSL certificates
- [x] Performance monitoring (Prometheus)
- [x] Error tracking (Sentry)
- [x] Log aggregation (ELK Stack)

### Production Checklist
- [ ] All services deployed to K8s
- [ ] Database migrations completed
- [ ] Redis cache warmed up
- [ ] Message queue topics created
- [ ] Monitoring dashboards active
- [ ] Alert thresholds configured
- [ ] Backup strategy implemented
- [ ] Disaster recovery plan tested
- [ ] Load testing (100K concurrent) passed
- [ ] Security audit completed
- [ ] GDPR compliance verified
- [ ] Documentation complete

### Deliverable
✅ Production-ready Triply platform handling massive scale

---

---

## 🎯 MVP Feature Matrix (3 Options)

```
                        OPTION A       OPTION B       OPTION C
                      (Basic MVP)    (Enhanced)      (Full Stack)
                       8-10 weeks    12+ weeks       20+ weeks
────────────────────────────────────────────────────────────────

LOGIN & AUTH              ✅            ✅            ✅
CREATE TRIP               ✅            ✅            ✅
ADD PLACES/SEGMENTS       ✅            ✅            ✅
PHOTO UPLOAD              ✅            ✅            ✅
RATINGS & REVIEWS         ✅            ✅            ✅
PUBLISH TO COMMUNITY      ✅            ✅            ✅
BROWSE COMMUNITY FEED     ✅            ✅            ✅
LIKE/UNLIKE POSTS         ✅            ✅            ✅
COMMENTS ON POSTS         ✅            ✅            ✅
BOOKING & PAYMENTS        ✅            ✅            ✅

GROUP TRIPS (NEW)         ─             ✅            ✅
GROUP VOTING (NEW)        ─             ✅            ✅
EXPENSE SPLITTING (NEW)   ─             ✅            ✅
UI POLISH & ANIMATIONS    ─             ✅            ✅
STOCK PHOTOS (Unsplash)   ─             ✅            ✅
REACT FRONTEND (NEW)      ─             ✅            ✅

AI TRIP GENERATION        ─             ─             ✅
PRICE ALERTS              ─             ─             ✅
MICROSERVICES (NEW)       ─             ─             ✅
KUBERNETES DEPLOYMENT     ─             ─             ✅
NOTIFICATIONS (Advanced)  ─             ─             ✅

────────────────────────────────────────────────────────────────
TEAM SIZE                 1-2 devs      2-3 devs      3-4 devs
COMPLEXITY              Low/Med        Medium/High    Very High
TIME TO MARKET           Fast           Balanced      Complete
USER EXPERIENCE          Good           Excellent     Premium
SCALABILITY              Good           Good          Enterprise
────────────────────────────────────────────────────────────────

RECOMMENDATION:
Start with OPTION A for fast market entry, then upgrade to B/C
based on user feedback and team capacity.
```

---

---

## 📊 Complete Development Timeline

```
OPTION A (Basic MVP) - 8-10 WEEKS
══════════════════════════════════

Week 1-2  │  M1: Core Flow Testing & Fixes
          │  Users can login → browse → plan
          │
Week 3-4  │  M2: Trip Creation & Storage
          │  Users create trips, save to database
          │
Week 5    │  M3: Photo Upload & Reviews
          │  Users add photos and ratings
          │  + M4: Community Feed & Engagement (parallel)
          │
Week 6-7  │  M5: Booking Integration (Razorpay)
          │  Users can purchase trips via payments
          │
Week 8-10 │  M7: Production Deployment
          │  Launch MVP to https://triply.app
          │
          │  STATUS: 🚀 LIVE!


OPTION B (Enhanced MVP) - 12+ WEEKS
════════════════════════════════════

Week 1-7  │  Same as OPTION A (Core features)
          │
Week 8-9  │  M5A: Group Trips & Collaborative Voting (PARALLEL)
          │  Friends can join, vote on places
          │  + M5B: UI Enhancement & Stock Photos (PARALLEL)
          │  Beautiful UI with Unsplash images
          │
Week 10-12│  M6: React Frontend Migration
          │  Rebuild with React for modern UX
          │  Real-time updates, smooth animations
          │
Week 13   │  M7: Production Deployment
          │  Launch with all new features
          │
          │  STATUS: 🚀 LIVE! (Feature-complete)


OPTION C (Full Stack) - 20+ WEEKS
══════════════════════════════════

Week 1-7  │  Core features (same as OPTION A)
          │
Week 8-9  │  M5A + M5B (parallel) - Group trips & UI polish
          │
Week 10-12│  M6: React Frontend Migration
          │
Week 13   │  M6.5: AI & Advanced Features
          │  AI trip generation, price alerts, recommendations
          │
Week 14-18│  M8: Microservices Architecture
          │  Extract into 7 independent services
          │  API Gateway, Service mesh, Event-driven
          │
Week 19-20│  M7: Full Production Deployment (K8s)
          │  Enterprise-ready infrastructure
          │
          │  STATUS: 🚀 LIVE! (Enterprise scale)
```

---

---

## 📈 Success Metrics Evolution

| Milestone | Target Users | Trips Published | Likes Total | Bookings | Expense Split | Group Trips |
|-----------|-------------|-----------------|-------------|----------|---------------|-------------|
| M1 | 0 | 0 | 0 | 0 | 0 | 0 |
| M2 | 10 | 0 | 0 | 0 | 0 | 0 |
| M3 | 25 | 5 | 0 | 0 | 0 | 0 |
| M4 | 50 | 25 | 500+ | 0 | 0 | 0 |
| M5 | 100 | 50 | 2000+ | 10+ | 0 | 0 |
| M5A🆕 | 150 | 75 | 3500+ | 25+ | ₹50K+ | 20+ |
| M5B🆕 | 200 | 100 | 5000+ | 40+ | ₹100K+ | 50+ |
| M6🆕 | 350 | 175 | 10000+ | 75+ | ₹250K+ | 100+ |
| M6.5 | 500 | 250 | 15000+ | 100+ | ₹500K+ | 150+ |
| M8🆕 | 1000+ | 500+ | 50000+ | 250+ | ₹2M+ | 500+ |

---

---

## 🎯 Decision Framework

**Choose OPTION A if:**
- You want to launch ASAP for market validation
- Your team is 1-2 developers
- Focus is on core platform features
- Fine with basic UI (no animations yet)
- Plan to iterate based on user feedback

**Choose OPTION B if:**
- You want a impressive first launch
- Your team is 2-3 developers
- Want to include collaboration features (group trips)
- Want polished UI/animations
- React frontend needed for future mobile app
- Timeline: 12+ weeks acceptable

**Choose OPTION C if:**
- You want enterprise-ready architecture
- Your team is 3-4 developers
- Plan to scale to 100K+ users
- Need AI-powered features
- Willing to invest 20+ weeks
- Long-term vision: global scaling

---

**Last Updated:** March 23, 2026  
**Current Status:** Planning Phase - Choose MVP Option (A, B, or C)  
**Next Checkpoint:** Decision on MVP level, then begin M1 implementation

