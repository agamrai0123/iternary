# Triply API Reference

**Base URL:** `http://localhost:8080` (development)  
**Version:** 1.0 (Beta)  
**Authentication:** Token-based (Bearer Token in Authorization header)  

---

## 📌 Table of Contents

1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Destinations](#destinations)
4. [User Trips](#user-trips)
5. [Trip Segments](#trip-segments)
6. [Photos](#photos)
7. [Reviews](#reviews)
8. [Community Posts](#community-posts)
9. [Error Codes](#error-codes)
10. [Rate Limiting](#rate-limiting)

---

## Authentication

All API requests (except login and public endpoints) require a Bearer token in the Authorization header.

### Login

Create an authentication session and receive a token.

```
POST /auth/login
Content-Type: application/json

{
  "email": "traveler@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "usr_123abc",
    "email": "traveler@example.com",
    "username": "TravelBug23",
    "avatar": "https://api.gravatar.com/avatar/abc123"
  },
  "expires_at": "2026-04-20T10:30:00Z"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "InvalidCredentials",
  "message": "Email or password is incorrect"
}
```

### Logout

```
POST /auth/logout
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "message": "Successfully logged out"
}
```

### Get Current User Profile

```
GET /auth/profile
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "id": "usr_123abc",
  "email": "traveler@example.com",
  "username": "TravelBug23",
  "avatar": "https://api.gravatar.com/avatar/abc123",
  "bio": "Love exploring new destinations",
  "created_at": "2026-01-15T14:32:00Z"
}
```

### Update User Profile

```
PUT /auth/profile
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "JetsAway",
  "bio": "Budget traveler exploring Asia",
  "avatar": "https://domain.com/avatar.jpg"
}
```

**Response (200 OK):**
```json
{
  "id": "usr_123abc",
  "email": "traveler@example.com",
  "username": "JetsAway",
  "bio": "Budget traveler exploring Asia",
  "avatar": "https://domain.com/avatar.jpg",
  "updated_at": "2026-03-23T10:15:00Z"
}
```

---

## User Management

### Register New User

```
POST /auth/signup
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "SecurePass123",
  "username": "ExplorerNow"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "usr_456def",
    "email": "newuser@example.com",
    "username": "ExplorerNow"
  }
}
```

**Response (400 Bad Request):**
```json
{
  "error": "ValidationError",
  "message": "Email already registered",
  "details": {
    "email": "This email is already in use"
  }
}
```

---

## Destinations

### List All Destinations

```
GET /api/destinations
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "dest_goa",
      "name": "Goa",
      "country": "India",
      "description": "Beautiful beaches and Portuguese heritage",
      "image_url": "https://images.example.com/goa.jpg",
      "trips_count": 245,
      "created_at": "2026-01-10T12:00:00Z"
    },
    {
      "id": "dest_bali",
      "name": "Bali",
      "country": "Indonesia",
      "description": "Tropical island with temples and rice paddies",
      "image_url": "https://images.example.com/bali.jpg",
      "trips_count": 412,
      "created_at": "2026-01-10T12:00:00Z"
    }
  ],
  "total": 28,
  "page": 1,
  "page_size": 20
}
```

### Get Destination Detail

```
GET /api/destinations/{destination_id}
```

**Response (200 OK):**
```json
{
  "id": "dest_goa",
  "name": "Goa",
  "country": "India",
  "description": "Beautiful beaches and Portuguese heritage",
  "image_url": "https://images.example.com/goa.jpg",
  "tips": [
    "Best time to visit: November-February",
    "Average budget: ₹40,000-₹70,000 per week",
    "Visa: Indian tourist visa required"
  ],
  "popular_itineraries": 245,
  "created_at": "2026-01-10T12:00:00Z"
}
```

---

## User Trips

### Create New Trip

Authenticated users only. Creates a new draft trip.

```
POST /api/user-trips
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Family Goa Getaway",
  "destination_id": "dest_goa",
  "budget": 500000,
  "duration": 5,
  "start_date": "2026-05-01T00:00:00Z",
  "description": "Planning a relaxing beach vacation with family"
}
```

**Request Validation:**
- `title`: Required, min 3 chars, max 100 chars
- `destination_id`: Required, must exist
- `budget`: Required, must be > 0 (in paise: ₹1 = 100 paise)
- `duration`: Required, must be > 0
- `start_date`: Optional, must be future date if provided

**Response (201 Created):**
```json
{
  "id": "trip_789xyz",
  "user_id": "usr_123abc",
  "title": "Family Goa Getaway",
  "destination_id": "dest_goa",
  "budget": 500000,
  "duration": 5,
  "start_date": "2026-05-01T00:00:00Z",
  "description": "Planning a relaxing beach vacation with family",
  "status": "draft",
  "segments": [],
  "created_at": "2026-03-23T10:15:00Z",
  "updated_at": "2026-03-23T10:15:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "ValidationError",
  "message": "Budget must be greater than 0",
  "details": {
    "budget": "Required and must be positive"
  }
}
```

### Get User's Trips

```
GET /api/user-trips?page=1&page_size=10
Authorization: Bearer {token}
```

**Query Parameters:**
- `page`: Integer, default 1
- `page_size`: Integer, default 10, max 50
- `status`: Optional filter (draft, planning, ongoing, completed, published)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "trip_789xyz",
      "title": "Family Goa Getaway",
      "destination_id": "dest_goa",
      "budget": 500000,
      "duration": 5,
      "status": "draft",
      "segments_count": 12,
      "photo_count": 28,
      "created_at": "2026-03-23T10:15:00Z",
      "updated_at": "2026-03-23T10:15:00Z"
    }
  ],
  "total": 5,
  "page": 1,
  "page_size": 10
}
```

### Get Single Trip Detail

Includes all segments, photos, and reviews (nested).

```
GET /api/user-trips/{trip_id}
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "id": "trip_789xyz",
  "user_id": "usr_123abc",
  "title": "Family Goa Getaway",
  "destination_id": "dest_goa",
  "budget": 500000,
  "duration": 5,
  "start_date": "2026-05-01T00:00:00Z",
  "description": "Planning a relaxing beach vacation with family",
  "status": "draft",
  "segments": [
    {
      "id": "seg_001",
      "day": 1,
      "name": "Hotel Oceanview",
      "type": "hotel",
      "location": "Baga Beach, Goa",
      "latitude": 15.5897,
      "longitude": 73.7997,
      "notes": "5-star resort with ocean view",
      "photos": [
        {
          "id": "photo_001",
          "url": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
          "caption": "Main lobby",
          "uploaded_at": "2026-03-23T10:20:00Z"
        }
      ],
      "review": {
        "id": "rev_001",
        "rating": 5,
        "review": "Excellent service, highly recommended!",
        "created_at": "2026-03-23T10:25:00Z"
      },
      "completed": false,
      "created_at": "2026-03-23T10:15:00Z"
    }
  ],
  "created_at": "2026-03-23T10:15:00Z",
  "updated_at": "2026-03-23T10:15:00Z"
}
```

**Response (404 Not Found):**
```json
{
  "error": "NotFound",
  "message": "Trip not found"
}
```

**Response (403 Forbidden):**
```json
{
  "error": "Forbidden",
  "message": "You don't have permission to access this trip"
}
```

### Update Trip

```
PUT /api/user-trips/{trip_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Goa Getaway",
  "budget": 600000,
  "duration": 6,
  "description": "Extended trip with island hopping"
}
```

**Response (200 OK):**
```json
{
  "id": "trip_789xyz",
  "title": "Updated Goa Getaway",
  "budget": 600000,
  "duration": 6,
  "description": "Extended trip with island hopping",
  "updated_at": "2026-03-23T11:00:00Z"
}
```

### Delete Trip

```
DELETE /api/user-trips/{trip_id}
Authorization: Bearer {token}
```

**Response (204 No Content):**
```
(No response body)
```

**Response (403 Forbidden):**
```json
{
  "error": "Forbidden",
  "message": "Cannot delete published trips. Unpublish first."
}
```

---

## Trip Segments

### Add Segment (Place/Activity)

```
POST /api/user-trips/{trip_id}/segments
Authorization: Bearer {token}
Content-Type: application/json

{
  "day": 1,
  "name": "Hotel Oceanview",
  "type": "hotel",
  "location": "Baga Beach, Goa",
  "latitude": 15.5897,
  "longitude": 73.7997,
  "notes": "5-star resort with ocean view"
}
```

**Request Validation:**
- `day`: Required, must be > 0 and <= trip duration
- `name`: Required, min 1 char, max 100 chars
- `type`: Optional (hotel, restaurant, activity, transport)
- `location`: Optional
- `latitude`: Optional, must be between -90 and 90
- `longitude`: Optional, must be between -180 and 180

**Response (201 Created):**
```json
{
  "id": "seg_001",
  "user_trip_id": "trip_789xyz",
  "day": 1,
  "name": "Hotel Oceanview",
  "type": "hotel",
  "location": "Baga Beach, Goa",
  "latitude": 15.5897,
  "longitude": 73.7997,
  "notes": "5-star resort with ocean view",
  "photos": [],
  "review": null,
  "completed": false,
  "created_at": "2026-03-23T10:15:00Z"
}
```

### Get Trip Segments

```
GET /api/user-trips/{trip_id}/segments
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "seg_001",
      "day": 1,
      "name": "Hotel Oceanview",
      "type": "hotel",
      "location": "Baga Beach, Goa",
      "photos_count": 3,
      "review": {
        "id": "rev_001",
        "rating": 5,
        "review": "Great place!"
      },
      "completed": false
    }
  ],
  "total": 12
}
```

### Update Segment

```
PUT /api/user-trips/{trip_id}/segments/{segment_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Hotel Paradise Updated",
  "notes": "Modified notes",
  "completed": true
}
```

**Response (200 OK):**
```json
{
  "id": "seg_001",
  "name": "Hotel Paradise Updated",
  "notes": "Modified notes",
  "completed": true,
  "updated_at": "2026-03-23T11:00:00Z"
}
```

### Delete Segment

```
DELETE /api/user-trips/{trip_id}/segments/{segment_id}
Authorization: Bearer {token}
```

**Response (204 No Content):**
```
(No response body - photos and reviews auto-deleted via CASCADE)
```

---

## Photos

### Upload Photo

Multipart form upload. Maximum 5MB per photo, max 3 photos per segment.

```
POST /api/trip-segments/{segment_id}/photos
Authorization: Bearer {token}
Content-Type: multipart/form-data

[Form Data]
file: (binary image data)
caption: "Beautiful ocean view from room"
```

**Accepted Formats:** jpg, jpeg, png
**Max Size:** 5MB
**Max Per Segment:** 3

**Response (201 Created):**
```json
{
  "id": "photo_001",
  "trip_segment_id": "seg_001",
  "url": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
  "caption": "Beautiful ocean view from room",
  "uploaded_at": "2026-03-23T10:20:00Z"
}
```

**Response (400 Bad Request - File too large):**
```json
{
  "error": "ValidationError",
  "message": "File size exceeds 5MB limit"
}
```

**Response (400 Bad Request - Too many photos):**
```json
{
  "error": "ValidationError",
  "message": "Maximum 3 photos per segment allowed"
}
```

### Get Segment Photos

```
GET /api/trip-segments/{segment_id}/photos
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "photo_001",
      "trip_segment_id": "seg_001",
      "url": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
      "caption": "Beautiful ocean view from room",
      "thumbnail": "https://storage.example.com/trip_789xyz/seg_001/thumb_photo_001.jpg",
      "uploaded_at": "2026-03-23T10:20:00Z"
    }
  ],
  "total": 3
}
```

### Delete Photo

```
DELETE /api/trip-segments/{segment_id}/photos/{photo_id}
Authorization: Bearer {token}
```

**Response (204 No Content):**
```
(No response body)
```

---

## Reviews

### Add/Update Review

One review per segment (upsert).

```
POST /api/trip-segments/{segment_id}/review
Authorization: Bearer {token}
Content-Type: application/json

{
  "rating": 5,
  "review": "Excellent service and friendly staff. Highly recommended!"
}
```

**Request Validation:**
- `rating`: Required, must be 1-5
- `review`: Optional, max 1000 characters

**Response (201 Created / 200 Updated):**
```json
{
  "id": "rev_001",
  "trip_segment_id": "seg_001",
  "rating": 5,
  "review": "Excellent service and friendly staff. Highly recommended!",
  "created_at": "2026-03-23T10:25:00Z",
  "updated_at": "2026-03-23T10:25:00Z"
}
```

**Response (400 Bad Request - Invalid rating):**
```json
{
  "error": "ValidationError",
  "message": "Rating must be between 1 and 5"
}
```

### Get Review

```
GET /api/trip-segments/{segment_id}/review
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "id": "rev_001",
  "trip_segment_id": "seg_001",
  "rating": 5,
  "review": "Excellent service and friendly staff. Highly recommended!",
  "created_at": "2026-03-23T10:25:00Z",
  "updated_at": "2026-03-23T10:25:00Z"
}
```

**Response (404 Not Found):**
```json
{
  "error": "NotFound",
  "message": "No review found for this segment"
}
```

---

## Community Posts

### Publish Trip

Convert a draft trip to public community post.

```
POST /api/user-trips/{trip_id}/publish
Authorization: Bearer {token}
Content-Type: application/json

{
  "cover_image": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg"
}
```

**Response (200 OK):**
```json
{
  "id": "post_abc123",
  "user_trip_id": "trip_789xyz",
  "user_id": "usr_123abc",
  "title": "Family Goa Getaway",
  "description": "Planning a relaxing beach vacation with family",
  "cover_image": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
  "likes": 0,
  "views": 0,
  "published": true,
  "published_at": "2026-03-23T11:05:00Z"
}
```

### Get Community Posts (Feed)

```
GET /api/community/posts?page=1&sort=newest&limit=20
```

**Query Parameters:**
- `page`: Integer, default 1
- `limit`: Integer, default 20, max 50
- `sort`: newest | trending | most_liked (default: newest)
- `destination_id`: Optional filter by destination

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "post_abc123",
      "user_trip_id": "trip_789xyz",
      "user": {
        "id": "usr_123abc",
        "username": "TravelBug23",
        "avatar": "https://api.gravatar.com/avatar/abc123"
      },
      "title": "Family Goa Getaway",
      "description": "Planning a relaxing beach vacation with family",
      "cover_image": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
      "destination": {
        "id": "dest_goa",
        "name": "Goa"
      },
      "duration": 5,
      "budget": 500000,
      "likes": 342,
      "views": 1250,
      "comments_count": 28,
      "liked_by_user": false,
      "published_at": "2026-03-23T11:05:00Z"
    }
  ],
  "total": 1250,
  "page": 1,
  "page_size": 20
}
```

### Get Post Detail

```
GET /api/community/posts/{post_id}
```

**Response (200 OK):**
```json
{
  "id": "post_abc123",
  "user_trip_id": "trip_789xyz",
  "user": {
    "id": "usr_123abc",
    "username": "TravelBug23",
    "avatar": "https://api.gravatar.com/avatar/abc123"
  },
  "title": "Family Goa Getaway",
  "description": "Planning a relaxing beach vacation with family",
  "cover_image": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
  "destination": {
    "id": "dest_goa",
    "name": "Goa",
    "country": "India"
  },
  "duration": 5,
  "budget": 500000,
  "segments": [
    {
      "id": "seg_001",
      "day": 1,
      "name": "Hotel Oceanview",
      "type": "hotel",
      "location": "Baga Beach, Goa",
      "rating": 5,
      "photos": [
        {
          "url": "https://storage.example.com/trip_789xyz/seg_001/photo_001.jpg",
          "caption": "Main lobby"
        }
      ]
    }
  ],
  "likes": 342,
  "views": 1250,
  "comments_count": 28,
  "liked_by_user": false,
  "published_at": "2026-03-23T11:05:00Z"
}
```

### Like Post

```
POST /api/community/posts/{post_id}/like
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "liked": true,
  "likes_count": 343
}
```

### Unlike Post

```
DELETE /api/community/posts/{post_id}/like
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "liked": false,
  "likes_count": 342
}
```

### Add Comment

```
POST /api/community/posts/{post_id}/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "This looks amazing! When are you going?",
  "rating": 5
}
```

**Request Validation:**
- `content`: Required, min 1 char, max 500 chars
- `rating`: Optional, 1-5 stars

**Response (201 Created):**
```json
{
  "id": "comment_001",
  "post_id": "post_abc123",
  "user": {
    "id": "usr_456def",
    "username": "ExplorerNow",
    "avatar": "https://api.gravatar.com/avatar/def456"
  },
  "content": "This looks amazing! When are you going?",
  "rating": 5,
  "created_at": "2026-03-23T12:30:00Z"
}
```

### Get Comments

```
GET /api/community/posts/{post_id}/comments?page=1&limit=20
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "comment_001",
      "user": {
        "id": "usr_456def",
        "username": "ExplorerNow",
        "avatar": "https://api.gravatar.com/avatar/def456"
      },
      "content": "This looks amazing! When are you going?",
      "rating": 5,
      "created_at": "2026-03-23T12:30:00Z"
    }
  ],
  "total": 28,
  "page": 1
}
```

---

## Error Codes

### Standard Error Response Format

```json
{
  "error": "ErrorType",
  "message": "Human-readable error message",
  "details": {
    "field_name": "Specific error for this field"
  },
  "timestamp": "2026-03-23T10:15:00Z"
}
```

### Error Types

| Code | Status | Description |
|------|--------|-------------|
| ValidationError | 400 | Request validation failed (missing/invalid fields) |
| InvalidCredentials | 401 | Email/password incorrect |
| Unauthorized | 401 | Missing or invalid token |
| Forbidden | 403 | User lacks permission (e.g., editing other's trip) |
| NotFound | 404 | Resource doesn't exist |
| Conflict | 409 | Resource already exists (e.g., email taken) |
| TooManyRequests | 429 | Rate limit exceeded |
| InternalServerError | 500 | Server error (our fault) |
| DatabaseError | 500 | Database operation failed |

---

## Rate Limiting

**Limits:**
- Anonymous users: 30 requests per minute per IP
- Authenticated users: 100 requests per minute per user
- Photo uploads: 10 per minute per user

**Response Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1711265900
```

**When limit exceeded (429 Too Many Requests):**
```json
{
  "error": "TooManyRequests",
  "message": "Rate limit exceeded. Try again in 52 seconds",
  "retry_after": 52
}
```

---

## CORS & Security Headers

**Allowed Origins (Production):** https://triply.app, https://www.triply.app

**Security Headers Returned:**
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
```

---

## Webhook Events (Future)

Future versions will support webhooks for:
- `trip.published` - User published a new trip
- `post.liked` - Someone liked user's post
- `post.commented` - Someone commented on user's post
- `payment.completed` - Booking payment succeeded
- `payment.failed` - Payment failed

---

## Pagination

All list endpoints support pagination:

**Query Parameters:**
- `page`: Page number (1-indexed), default 1
- `limit` or `page_size`: Items per page, default 20, max 50

**Response:**
```json
{
  "data": [...],
  "total": 1250,
  "page": 1,
  "page_size": 20,
  "has_more": true,
  "total_pages": 63
}
```

---

## Sorting & Filtering

**Sorting (in list endpoints):**
```
GET /api/community/posts?sort=newest
GET /api/community/posts?sort=likes:desc
GET /api/community/posts?sort=views:asc
```

**Filtering:**
```
GET /api/user-trips?status=published
GET /api/community/posts?destination_id=dest_goa
GET /api/community/posts?min_rating=4
```

---

## Timestamps

All timestamps in responses use ISO 8601 format with UTC timezone:

```
"created_at": "2026-03-23T10:15:00Z"
"updated_at": "2026-03-23T11:05:00Z"
"published_at": "2026-03-23T11:05:00Z"
```

---

## Request/Response Content-Type

- **Standard requests/responses:** `application/json`
- **Photo upload:** `multipart/form-data`
- **All responses have:** `Content-Type: application/json; charset=utf-8`

---

## API Version Management

Current API version: `1.0` (Beta)

Version indicated in:
1. GitHub releases
2. Documentation
3. May eventually add `Accept: application/vnd.triply.v1+json` header support

---

## Deprecation Policy

Endpoints marked as deprecated will:
1. Continue working for 90 days
2. Return `Deprecation: true` in response headers
3. Include `X-API-Warn: This endpoint is deprecated` header
4. Have `deprecated: true` in OpenAPI spec

---

## Support & Feedback

- **Status Page:** https://status.triply.app
- **API Status:** Check status page for incidents
- **Documentation:** https://api.triply.app/docs
- **Support Email:** api-support@triply.app

---

**Last Updated:** March 23, 2026  
**API Version:** 1.0 (Beta)  
**Status:** Ready for Development

