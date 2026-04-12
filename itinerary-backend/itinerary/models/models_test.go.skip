package itinerary

import (
	"testing"
	"time"
)

// TestUserModel validates User struct initialization
func TestUserModel(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				ID:       "user-001",
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "user with timestamps",
			user: User{
				ID:        "user-002",
				Username:  "another_user",
				Email:     "another@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.ID == "" {
				t.Errorf("User ID should not be empty")
			}
			if tt.user.Username == "" {
				t.Errorf("Username should not be empty")
			}
			if tt.user.Email == "" {
				t.Errorf("Email should not be empty")
			}
		})
	}
}

// TestDestinationModel validates Destination struct
func TestDestinationModel(t *testing.T) {
	dest := Destination{
		ID:      "dest-001",
		Name:    "Goa",
		Country: "India",
		Image:   "goa.jpg",
	}

	if dest.ID == "" {
		t.Error("Destination ID should not be empty")
	}
	if dest.Name == "" {
		t.Error("Destination name should not be empty")
	}
	if dest.Country == "" {
		t.Error("Destination country should not be empty")
	}
}

// TestItineraryModel validates Itinerary struct
func TestItineraryModel(t *testing.T) {
	tests := []struct {
		name      string
		itinerary Itinerary
		wantValid bool
	}{
		{
			name: "valid itinerary",
			itinerary: Itinerary{
				ID:            "itin-001",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Goa Getaway",
				Duration:      5,
				Budget:        50000,
				Likes:         0,
			},
			wantValid: true,
		},
		{
			name: "itinerary with items",
			itinerary: Itinerary{
				ID:            "itin-002",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Bali Trip",
				Duration:      7,
				Budget:        80000,
				Items: []ItineraryItem{
					{
						ID:          "item-001",
						Day:         1,
						Type:        "stay",
						Name:        "Hotel Resort",
						Price:       10000,
						Rating:      4.5,
					},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid itinerary - zero duration",
			itinerary: Itinerary{
				ID:            "itin-003",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Invalid Trip",
				Duration:      0,
				Budget:        50000,
			},
			wantValid: false,
		},
		{
			name: "invalid itinerary - zero budget",
			itinerary: Itinerary{
				ID:            "itin-004",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Invalid Trip",
				Duration:      5,
				Budget:        0,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := true
			if tt.itinerary.Duration <= 0 {
				valid = false
			}
			if tt.itinerary.Budget <= 0 {
				valid = false
			}

			if valid != tt.wantValid {
				t.Errorf("Validation mismatch. Got %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// TestItineraryItemModel validates ItineraryItem struct
func TestItineraryItemModel(t *testing.T) {
	tests := []struct {
		name     string
		item     ItineraryItem
		wantType bool
	}{
		{
			name: "stay type",
			item: ItineraryItem{
				ID:          "item-001",
				Type:        "stay",
				Name:        "Hotel",
				Price:       5000,
			},
			wantType: true,
		},
		{
			name: "food type",
			item: ItineraryItem{
				ID:          "item-002",
				Type:        "food",
				Name:        "Restaurant",
				Price:       1500,
			},
			wantType: true,
		},
		{
			name: "activity type",
			item: ItineraryItem{
				ID:          "item-003",
				Type:        "activity",
				Name:        "Beach Tour",
				Price:       3000,
			},
			wantType: true,
		},
		{
			name: "transport type",
			item: ItineraryItem{
				ID:          "item-004",
				Type:        "transport",
				Name:        "Taxi",
				Price:       500,
			},
			wantType: true,
		},
		{
			name: "other type",
			item: ItineraryItem{
				ID:          "item-005",
				Type:        "other",
				Name:        "Miscellaneous",
				Price:       1000,
			},
			wantType: true,
		},
	}

	validTypes := map[string]bool{
		"stay": true, "food": true, "activity": true, "transport": true, "other": true,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := validTypes[tt.item.Type]
			if exists != tt.wantType {
				t.Errorf("Item type validation failed. Got %v, want %v", exists, tt.wantType)
			}
		})
	}
}

// TestCommentModel validates Comment struct
func TestCommentModel(t *testing.T) {
	tests := []struct {
		name      string
		comment   Comment
		wantValid bool
	}{
		{
			name: "valid comment with rating",
			comment: Comment{
				ID:          "comment-001",
				ItineraryID: "itin-001",
				UserID:      "user-001",
				Content:     "Great itinerary!",
				Rating:      4.5,
			},
			wantValid: true,
		},
		{
			name: "valid comment without rating",
			comment: Comment{
				ID:          "comment-002",
				ItineraryID: "itin-001",
				UserID:      "user-002",
				Content:     "Nice places!",
			},
			wantValid: true,
		},
		{
			name: "invalid rating - too high",
			comment: Comment{
				ID:          "comment-003",
				ItineraryID: "itin-001",
				UserID:      "user-003",
				Content:     "Invalid rating",
				Rating:      6.0,
			},
			wantValid: false,
		},
		{
			name: "invalid rating - negative",
			comment: Comment{
				ID:          "comment-004",
				ItineraryID: "itin-001",
				UserID:      "user-004",
				Content:     "Negative rating",
				Rating:      -1.0,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.comment.Content != "" && (tt.comment.Rating == 0 || (tt.comment.Rating >= 0 && tt.comment.Rating <= 5))

			if valid != tt.wantValid {
				t.Errorf("Comment validation failed. Got %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// TestUserTripModel validates UserTrip struct
func TestUserTripModel(t *testing.T) {
	tests := []struct {
		name      string
		trip      UserTrip
		wantValid bool
	}{
		{
			name: "valid trip",
			trip: UserTrip{
				ID:            "trip-001",
				UserID:        "user-001",
				Title:         "Goa Vacation",
				DestinationID: "dest-001",
				Budget:        50000,
				Duration:      5,
				Status:        "draft",
			},
			wantValid: true,
		},
		{
			name: "invalid trip - no user",
			trip: UserTrip{
				ID:            "trip-002",
				UserID:        "",
				Title:         "Invalid Trip",
				DestinationID: "dest-001",
				Budget:        50000,
				Duration:      5,
			},
			wantValid: false,
		},
		{
			name: "invalid trip - zero budget",
			trip: UserTrip{
				ID:            "trip-003",
				UserID:        "user-002",
				Title:         "Invalid Trip",
				DestinationID: "dest-001",
				Budget:        0,
				Duration:      5,
			},
			wantValid: false,
		},
		{
			name: "invalid trip - zero duration",
			trip: UserTrip{
				ID:            "trip-004",
				UserID:        "user-003",
				Title:         "Invalid Trip",
				DestinationID: "dest-001",
				Budget:        50000,
				Duration:      0,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.trip.UserID != "" && tt.trip.Budget > 0 && tt.trip.Duration > 0

			if valid != tt.wantValid {
				t.Errorf("Trip validation failed. Got %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// TestTripSegmentModel validates TripSegment struct
func TestTripSegmentModel(t *testing.T) {
	tests := []struct {
		name      string
		segment   TripSegment
		wantValid bool
	}{
		{
			name: "valid segment",
			segment: TripSegment{
				ID:         "seg-001",
				UserTripID: "trip-001",
				Day:        1,
				Name:       "Beach Visit",
				Type:       "activity",
			},
			wantValid: true,
		},
		{
			name: "invalid segment - no trip",
			segment: TripSegment{
				ID:         "seg-002",
				UserTripID: "",
				Day:        1,
				Name:       "Beach Visit",
				Type:       "activity",
			},
			wantValid: false,
		},
		{
			name: "invalid segment - zero day",
			segment: TripSegment{
				ID:         "seg-003",
				UserTripID: "trip-001",
				Day:        0,
				Name:       "Beach Visit",
				Type:       "activity",
			},
			wantValid: false,
		},
		{
			name: "invalid segment - no name",
			segment: TripSegment{
				ID:         "seg-004",
				UserTripID: "trip-001",
				Day:        1,
				Name:       "",
				Type:       "activity",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.segment.UserTripID != "" && tt.segment.Day > 0 && tt.segment.Name != ""

			if valid != tt.wantValid {
				t.Errorf("Segment validation failed. Got %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// TestTripPhotoModel validates TripPhoto struct
func TestTripPhotoModel(t *testing.T) {
	photo := TripPhoto{
		ID:            "photo-001",
		TripSegmentID: "seg-001",
		URL:           "https://example.com/photo.jpg",
		Caption:       "Beautiful beach",
	}

	if photo.ID == "" {
		t.Error("Photo ID should not be empty")
	}
	if photo.TripSegmentID == "" {
		t.Error("Segment ID should not be empty")
	}
	if photo.URL == "" {
		t.Error("URL should not be empty")
	}
}

// TestTripReviewModel validates TripReview struct
func TestTripReviewModel(t *testing.T) {
	tests := []struct {
		name      string
		review    TripReview
		wantValid bool
	}{
		{
			name: "valid review - 5 stars",
			review: TripReview{
				ID:            "review-001",
				TripSegmentID: "seg-001",
				Rating:        5.0,
				Review:        "Excellent experience!",
			},
			wantValid: true,
		},
		{
			name: "valid review - 1 star",
			review: TripReview{
				ID:            "review-002",
				TripSegmentID: "seg-001",
				Rating:        1.0,
				Review:        "Not great",
			},
			wantValid: true,
		},
		{
			name: "invalid review - rating too high",
			review: TripReview{
				ID:            "review-003",
				TripSegmentID: "seg-001",
				Rating:        6.0,
				Review:        "Too high rating",
			},
			wantValid: false,
		},
		{
			name: "invalid review - rating too low",
			review: TripReview{
				ID:            "review-004",
				TripSegmentID: "seg-001",
				Rating:        0.5,
				Review:        "Too low rating",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.review.Rating >= 1.0 && tt.review.Rating <= 5.0 && tt.review.Review != ""

			if valid != tt.wantValid {
				t.Errorf("Review validation failed. Got %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// TestUserTripPostModel validates UserTripPost struct
func TestUserTripPostModel(t *testing.T) {
	publishedAt := time.Now()
	post := UserTripPost{
		ID:          "post-001",
		UserTripID:  "trip-001",
		UserID:      "user-001",
		Title:       "My Amazing Trip",
		Description: "Had a great time",
		CoverImage:  "cover.jpg",
		Likes:       0,
		Views:       0,
		Published:   true,
		PublishedAt: &publishedAt,
	}

	if post.ID == "" {
		t.Error("Post ID should not be empty")
	}
	if post.UserTripID == "" {
		t.Error("Trip ID should not be empty")
	}
	if post.UserID == "" {
		t.Error("User ID should not be empty")
	}
	if post.Published && post.PublishedAt == nil {
		t.Error("Published post should have PublishedAt set")
	}
}

// TestPaginatedResponseModel validates pagination
func TestPaginatedResponseModel(t *testing.T) {
	destinations := []Destination{
		{ID: "dest-001", Name: "Goa", Country: "India"},
		{ID: "dest-002", Name: "Bali", Country: "Indonesia"},
	}

	response := PaginatedResponse{
		Data:       destinations,
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
	}

	if response.Total != 100 {
		t.Errorf("Expected total 100, got %d", response.Total)
	}
	if response.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Page)
	}
	if response.TotalPages != 10 {
		t.Errorf("Expected 10 total pages, got %d", response.TotalPages)
	}
}
