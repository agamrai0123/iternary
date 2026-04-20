package itinerary

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" `
	Username  string    `json:"username" binding:"required,min=3,max=50"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Destination represents a travel destination
type Destination struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Country     string    `json:"country" binding:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Itinerary represents a complete travel itinerary plan
type Itinerary struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id" binding:"required"`
	DestinationID string          `json:"destination_id" binding:"required"`
	Title         string          `json:"title" binding:"required"`
	Description   string          `json:"description"`
	Duration      int             `json:"duration" binding:"required,min=1"`
	Budget        float64         `json:"budget" binding:"required,gt=0"`
	Items         []ItineraryItem `json:"items"`
	Likes         int             `json:"likes"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ItineraryItem represents a single item (stay, food, activity) in an itinerary
type ItineraryItem struct {
	ID          string    `json:"id"`
	ItineraryID string    `json:"itinerary_id"`
	Day         int       `json:"day" binding:"required,min=1"`
	Type        string    `json:"type" binding:"required,oneof=stay food activity transport other"` // stay, food, activity, transport, etc.
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gte=0"`
	Duration    int       `json:"duration"` // in hours
	Location    string    `json:"location"`
	Rating      float64   `json:"rating"` // 0-5
	ImageURL    string    `json:"image_url"`
	BookingURL  string    `json:"booking_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Comment represents a user comment on an itinerary
type Comment struct {
	ID          string    `json:"id"`
	ItineraryID string    `json:"itinerary_id" binding:"required"`
	UserID      string    `json:"user_id" binding:"required"`
	Content     string    `json:"content" binding:"required,min=1,max=1000"`
	Rating      float64   `json:"rating"` // optional: 0-5
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserPlan represents a user's saved/copied itinerary plan
type UserPlan struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id" binding:"required"`
	OriginalID string    `json:"original_itinerary_id" binding:"required"`
	Title      string    `json:"title"`
	Notes      string    `json:"notes"`
	Status     string    `json:"status"` // draft, planned, completed
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
// ==================== NEW: USER TRIP PLANNING MODELS ====================

// UserTrip represents a user's custom trip plan
type UserTrip struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id" binding:"required"`
	Title         string        `json:"title" binding:"required"`
	DestinationID string        `json:"destination_id" binding:"required"`
	Budget        float64       `json:"budget" binding:"required,gt=0"`
	Duration      int           `json:"duration" binding:"required,min=1"` // in days
	StartDate     time.Time     `json:"start_date"`
	Segments      []TripSegment `json:"segments"`
	Status        string        `json:"status"` // draft, planning, ongoing, completed
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// TripSegment represents a place/activity in a user trip
type TripSegment struct {
	ID              string        `json:"id"`
	UserTripID      string        `json:"user_trip_id" binding:"required"`
	Day             int           `json:"day" binding:"required,min=1"`
	TimeOfDay       string        `json:"time_of_day"` // morning, afternoon, evening, night
	Name            string        `json:"name" binding:"required"`
	Type            string        `json:"type"` // stay, food, activity, transport, other
	Location        string        `json:"location"`
	Latitude        float64       `json:"latitude"`  // For Google Maps
	Longitude       float64       `json:"longitude"` // For Google Maps
	Expense         float64       `json:"expense" binding:"gte=0"` // Cost of the place/activity
	BestTimeToVisit string        `json:"best_time_to_visit"` // e.g., "Spring", "Morning", "Peak hours 10am-12pm"
	Notes           string        `json:"notes"`
	Photos          []TripPhoto   `json:"photos"`
	Review          *TripReview   `json:"review"`
	Completed       bool          `json:"completed"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// TripPhoto represents a photo for a trip segment
type TripPhoto struct {
	ID           string    `json:"id"`
	TripSegmentID string   `json:"trip_segment_id" binding:"required"`
	URL          string    `json:"url" binding:"required"`
	Caption      string    `json:"caption"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// TripReview represents a review for a completed trip segment
type TripReview struct {
	ID           string    `json:"id"`
	TripSegmentID string   `json:"trip_segment_id" binding:"required"`
	Rating       float64   `json:"rating" binding:"required,min=1,max=5"`
	Review       string    `json:"review" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserTripPost represents a published trip post in the community
type UserTripPost struct {
	ID            string        `json:"id"`
	UserTripID    string        `json:"user_trip_id" binding:"required"`
	UserID        string        `json:"user_id" binding:"required"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	CoverImage    string        `json:"cover_image"`
	DestinationID string        `json:"destination_id"` // Which city/destination this trip is about
	Duration      int           `json:"duration"` // Number of days
	TotalExpense  float64       `json:"total_expense"` // Total cost of the trip
	Segments      []TripSegment `json:"segments"` // All places/activities in this trip post
	Likes         int           `json:"likes"`
	Views         int           `json:"views"`
	Published     bool          `json:"published"`
	PublishedAt   *time.Time    `json:"published_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// TripPostResponse is the response format for trip posts in API
type TripPostResponse struct {
	ID            string                `json:"id"`
	UserID        string                `json:"user_id"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	CoverImage    string                `json:"cover_image"`
	DestinationID string                `json:"destination_id"`
	Duration      int                   `json:"duration"`
	TotalExpense  float64               `json:"total_expense"`
	Places        []TripPlaceResponse   `json:"places"` // Renamed from segments for clarity
	Likes         int                   `json:"likes"`
	Views         int                   `json:"views"`
	PublishedAt   *time.Time            `json:"published_at"`
	CreatedAt     time.Time             `json:"created_at"`
}

// TripPlaceResponse represents a place in a trip post
type TripPlaceResponse struct {
	ID              string              `json:"id"`
	Day             int                 `json:"day"`
	TimeOfDay       string              `json:"time_of_day"`
	Name            string              `json:"name"`
	Type            string              `json:"type"`
	Location        string              `json:"location"`
	Latitude        float64             `json:"latitude"`
	Longitude       float64             `json:"longitude"`
	Expense         float64             `json:"expense"`
	BestTimeToVisit string              `json:"best_time_to_visit"`
	Photos          []TripPhoto         `json:"photos"`
	Review          *TripReviewResponse `json:"review"`
}

// TripReviewResponse represents a review in a trip post
type TripReviewResponse struct {
	Rating float64 `json:"rating"`
	Review string  `json:"review"`
}

// ==================== API Request/Response Types ====================

// AddTripPostToItineraryRequest requests adding a shared trip post to user's itinerary
type AddTripPostToItineraryRequest struct {
	TripPostID string `json:"trip_post_id" binding:"required"`
}

// MarkSegmentVisitedRequest requests marking a trip segment as visited
type MarkSegmentVisitedRequest struct {
	Completed bool `json:"completed" binding:"required"`
}

// SubmitReviewRequest represents a review submission for a visited place
type SubmitReviewRequest struct {
	SegmentID string  `json:"segment_id" binding:"required"`
	Rating    float64 `json:"rating" binding:"required,min=1,max=5"`
	Review    string  `json:"review" binding:"required"`
}

// ListCitiesResponse represents the response for listing cities/destinations
type ListCitiesResponse struct {
	Data       []Destination `json:"data"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}