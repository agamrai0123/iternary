package itinerary

import (
	"testing"
	"time"
)

// MockDatabase provides a mock implementation of database operations
type MockDatabase struct {
	destinations   []Destination
	itineraries    []Itinerary
	userTrips      []UserTrip
	tripSegments   []TripSegment
	communityPosts []UserTripPost
}

// GetDestinations mock implementation
func (m *MockDatabase) GetDestinations(page, pageSize int) ([]Destination, int, error) {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(m.destinations) {
		return []Destination{}, len(m.destinations), nil
	}

	if end > len(m.destinations) {
		end = len(m.destinations)
	}

	return m.destinations[start:end], len(m.destinations), nil
}

// GetItinerariesByDestination mock implementation
func (m *MockDatabase) GetItinerariesByDestination(destID string, page, pageSize int) ([]Itinerary, int, error) {
	var filtered []Itinerary
	for _, it := range m.itineraries {
		if it.DestinationID == destID {
			filtered = append(filtered, it)
		}
	}

	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(filtered) {
		return []Itinerary{}, len(filtered), nil
	}

	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], len(filtered), nil
}

// GetItineraryItems mock implementation
func (m *MockDatabase) GetItineraryItems(itineraryID string) ([]ItineraryItem, error) {
	return []ItineraryItem{}, nil
}

// GetItineraryByID mock implementation
func (m *MockDatabase) GetItineraryByID(id string) (*Itinerary, error) {
	for i := range m.itineraries {
		if m.itineraries[i].ID == id {
			return &m.itineraries[i], nil
		}
	}
	return nil, nil
}

// CreateItinerary mock implementation
func (m *MockDatabase) CreateItinerary(it *Itinerary) error {
	m.itineraries = append(m.itineraries, *it)
	return nil
}

// AddLikeToItinerary mock implementation
func (m *MockDatabase) AddLikeToItinerary(itineraryID string) error {
	for i := range m.itineraries {
		if m.itineraries[i].ID == itineraryID {
			m.itineraries[i].Likes++
			return nil
		}
	}
	return nil
}

// AddComment mock implementation
func (m *MockDatabase) AddComment(comment *Comment) error {
	return nil
}

// CreateUserTrip mock implementation
func (m *MockDatabase) CreateUserTrip(trip *UserTrip) error {
	m.userTrips = append(m.userTrips, *trip)
	return nil
}

// GetUserTripByID mock implementation
func (m *MockDatabase) GetUserTripByID(id string) (*UserTrip, error) {
	for i := range m.userTrips {
		if m.userTrips[i].ID == id {
			return &m.userTrips[i], nil
		}
	}
	return nil, nil
}

// GetUserTripsByUserID mock implementation
func (m *MockDatabase) GetUserTripsByUserID(userID string) ([]UserTrip, error) {
	var trips []UserTrip
	for _, trip := range m.userTrips {
		if trip.UserID == userID {
			trips = append(trips, trip)
		}
	}
	return trips, nil
}

// UpdateUserTrip mock implementation
func (m *MockDatabase) UpdateUserTrip(trip *UserTrip) error {
	for i := range m.userTrips {
		if m.userTrips[i].ID == trip.ID {
			m.userTrips[i] = *trip
			return nil
		}
	}
	return nil
}

// DeleteUserTrip mock implementation
func (m *MockDatabase) DeleteUserTrip(tripID string) error {
	for i := range m.userTrips {
		if m.userTrips[i].ID == tripID {
			m.userTrips = append(m.userTrips[:i], m.userTrips[i+1:]...)
			return nil
		}
	}
	return nil
}

// AddTripSegment mock implementation
func (m *MockDatabase) AddTripSegment(segment *TripSegment) error {
	m.tripSegments = append(m.tripSegments, *segment)
	return nil
}

// AddTripPhoto mock implementation
func (m *MockDatabase) AddTripPhoto(photo *TripPhoto) error {
	return nil
}

// AddTripReview mock implementation
func (m *MockDatabase) AddTripReview(review *TripReview) error {
	return nil
}

// PublishUserTrip mock implementation
func (m *MockDatabase) PublishUserTrip(post *UserTripPost) error {
	m.communityPosts = append(m.communityPosts, *post)
	return nil
}

// UpdateUserTripStatus mock implementation
func (m *MockDatabase) UpdateUserTripStatus(tripID, status string) error {
	for i := range m.userTrips {
		if m.userTrips[i].ID == tripID {
			m.userTrips[i].Status = status
			return nil
		}
	}
	return nil
}

// GetCommunityPosts mock implementation
func (m *MockDatabase) GetCommunityPosts(page, pageSize int) ([]UserTripPost, error) {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(m.communityPosts) {
		return []UserTripPost{}, nil
	}

	if end > len(m.communityPosts) {
		end = len(m.communityPosts)
	}

	return m.communityPosts[start:end], nil
}

// TestServiceGetDestinations verifies destination retrieval
func TestServiceGetDestinations(t *testing.T) {
	mockDB := &MockDatabase{
		destinations: []Destination{
			{ID: "dest-001", Name: "Goa", Country: "India"},
			{ID: "dest-002", Name: "Bali", Country: "Indonesia"},
			{ID: "dest-003", Name: "Paris", Country: "France"},
		},
	}

	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantCount  int
		wantTotal  int
	}{
		{
			name:      "first page with default size",
			page:      1,
			pageSize:  10,
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "invalid page defaults to 1",
			page:      0,
			pageSize:  10,
			wantCount: 3,
			wantTotal: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dests, total, err := service.GetDestinations(tt.page, tt.pageSize)

			if err != nil {
				t.Errorf("GetDestinations() error = %v", err)
			}

			if len(dests) != tt.wantCount {
				t.Errorf("Expected %d destinations, got %d", tt.wantCount, len(dests))
			}

			if total != tt.wantTotal {
				t.Errorf("Expected total %d, got %d", tt.wantTotal, total)
			}
		})
	}
}

// TestServiceCreateItinerary verifies itinerary creation validation
func TestServiceCreateItinerary(t *testing.T) {
	mockDB := &MockDatabase{}
	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name      string
		itinerary *Itinerary
		wantErr   bool
	}{
		{
			name: "valid itinerary",
			itinerary: &Itinerary{
				ID:            "itin-001",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Goa Trip",
				Duration:      5,
				Budget:        50000,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			itinerary: &Itinerary{
				ID:            "itin-002",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "",
				Duration:      5,
				Budget:        50000,
			},
			wantErr: true,
		},
		{
			name: "zero duration",
			itinerary: &Itinerary{
				ID:            "itin-003",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Trip",
				Duration:      0,
				Budget:        50000,
			},
			wantErr: true,
		},
		{
			name: "zero budget",
			itinerary: &Itinerary{
				ID:            "itin-004",
				UserID:        "user-001",
				DestinationID: "dest-001",
				Title:         "Trip",
				Duration:      5,
				Budget:        0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateItinerary(tt.itinerary)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateItinerary() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify timestamps are set on success
			if err == nil && tt.itinerary.CreatedAt.IsZero() {
				t.Error("CreatedAt should be set")
			}
		})
	}
}

// TestServiceAddLikeToItinerary verifies like functionality
func TestServiceAddLikeToItinerary(t *testing.T) {
	mockDB := &MockDatabase{
		itineraries: []Itinerary{
			{ID: "itin-001", Likes: 0},
		},
	}

	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name          string
		itineraryID   string
		wantErr       bool
	}{
		{
			name:        "add like to existing itinerary",
			itineraryID: "itin-001",
			wantErr:     false,
		},
		{
			name:        "empty itinerary ID",
			itineraryID: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddLikeToItinerary(tt.itineraryID)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddLikeToItinerary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestServiceCreateUserTrip verifies user trip creation
func TestServiceCreateUserTrip(t *testing.T) {
	mockDB := &MockDatabase{}
	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name    string
		trip    *UserTrip
		wantErr bool
	}{
		{
			name: "valid user trip",
			trip: &UserTrip{
				ID:            "trip-001",
				UserID:        "user-001",
				Title:         "Goa Vacation",
				DestinationID: "dest-001",
				Budget:        50000,
				Duration:      5,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			trip: &UserTrip{
				ID:            "trip-002",
				UserID:        "",
				Title:         "Trip",
				DestinationID: "dest-001",
				Budget:        50000,
				Duration:      5,
			},
			wantErr: true,
		},
		{
			name: "missing title",
			trip: &UserTrip{
				ID:            "trip-003",
				UserID:        "user-001",
				Title:         "",
				DestinationID: "dest-001",
				Budget:        50000,
				Duration:      5,
			},
			wantErr: true,
		},
		{
			name: "zero budget",
			trip: &UserTrip{
				ID:            "trip-004",
				UserID:        "user-001",
				Title:         "Trip",
				DestinationID: "dest-001",
				Budget:        0,
				Duration:      5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateUserTrip(tt.trip)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserTrip() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify default status is set
			if err == nil && tt.trip.Status != "draft" {
				t.Errorf("Default status should be 'draft', got %q", tt.trip.Status)
			}
		})
	}
}

// TestServiceGetUserTrips verifies user trips retrieval
func TestServiceGetUserTrips(t *testing.T) {
	mockDB := &MockDatabase{
		userTrips: []UserTrip{
			{ID: "trip-001", UserID: "user-001", Title: "Trip 1"},
			{ID: "trip-002", UserID: "user-001", Title: "Trip 2"},
			{ID: "trip-003", UserID: "user-002", Title: "Trip 3"},
		},
	}

	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name      string
		userID    string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get trips for user with multiple trips",
			userID:    "user-001",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "get trips for user with one trip",
			userID:    "user-002",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "empty user ID",
			userID:    "",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trips, err := service.GetUserTrips(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTrips() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(trips) != tt.wantCount {
				t.Errorf("Expected %d trips, got %d", tt.wantCount, len(trips))
			}
		})
	}
}

// TestServiceAddTripSegment verifies segment creation validation
func TestServiceAddTripSegment(t *testing.T) {
	mockDB := &MockDatabase{}
	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name    string
		segment *TripSegment
		wantErr bool
	}{
		{
			name: "valid segment",
			segment: &TripSegment{
				ID:         "seg-001",
				UserTripID: "trip-001",
				Day:        1,
				Name:       "Hotel",
			},
			wantErr: false,
		},
		{
			name: "missing trip ID",
			segment: &TripSegment{
				ID:         "seg-002",
				UserTripID: "",
				Day:        1,
				Name:       "Hotel",
			},
			wantErr: true,
		},
		{
			name: "zero day",
			segment: &TripSegment{
				ID:         "seg-003",
				UserTripID: "trip-001",
				Day:        0,
				Name:       "Hotel",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			segment: &TripSegment{
				ID:         "seg-004",
				UserTripID: "trip-001",
				Day:        1,
				Name:       "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddTripSegment(tt.segment)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddTripSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestServiceAddTripReview verifies review validation
func TestServiceAddTripReview(t *testing.T) {
	mockDB := &MockDatabase{}
	logger := &Logger{}
	service := NewService(mockDB, logger)

	tests := []struct {
		name    string
		review  *TripReview
		wantErr bool
	}{
		{
			name: "valid 5-star review",
			review: &TripReview{
				ID:            "review-001",
				TripSegmentID: "seg-001",
				Rating:        5.0,
				Review:        "Excellent!",
			},
			wantErr: false,
		},
		{
			name: "valid 1-star review",
			review: &TripReview{
				ID:            "review-002",
				TripSegmentID: "seg-001",
				Rating:        1.0,
				Review:        "Poor",
			},
			wantErr: false,
		},
		{
			name: "invalid rating - too high",
			review: &TripReview{
				ID:            "review-003",
				TripSegmentID: "seg-001",
				Rating:        6.0,
				Review:        "Review",
			},
			wantErr: true,
		},
		{
			name: "invalid rating - too low",
			review: &TripReview{
				ID:            "review-004",
				TripSegmentID: "seg-001",
				Rating:        0.5,
				Review:        "Review",
			},
			wantErr: true,
		},
		{
			name: "missing segment ID",
			review: &TripReview{
				ID:            "review-005",
				TripSegmentID: "",
				Rating:        3.0,
				Review:        "Review",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddTripReview(tt.review)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddTripReview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestServicePublishUserTrip verifies trip publication
func TestServicePublishUserTrip(t *testing.T) {
	mockDB := &MockDatabase{
		userTrips: []UserTrip{
			{ID: "trip-001", UserID: "user-001", Status: "draft"},
		},
	}

	logger := &Logger{}
	service := NewService(mockDB, logger)

	now := time.Now()
	post := &UserTripPost{
		ID:         "post-001",
		UserTripID: "trip-001",
		UserID:     "user-001",
		Title:      "My Trip",
	}

	err := service.PublishUserTrip(post)

	if err != nil {
		t.Errorf("PublishUserTrip() error = %v", err)
	}

	if post.PublishedAt == nil {
		t.Error("PublishedAt should be set")
	}

	if post.PublishedAt.Before(now) {
		t.Error("PublishedAt should be current time or later")
	}

	// Verify trip status was updated
	trip, _ := mockDB.GetUserTripByID("trip-001")
	if trip.Status != "published" {
		t.Errorf("Trip status should be 'published', got %q", trip.Status)
	}
}
