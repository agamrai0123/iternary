package itinerary

// This file has been disabled - Service, Database, and Logger types are not defined properly
// See service_*.go files for active service implementations
	db     *common.Database
	logger *common.Logger
}

// NewService creates a new service instance
func NewService(db *common.Database, logger *common.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}

// GetDestinations retrieves destinations with pagination
func (s *Service) GetDestinations(page, pageSize int) ([]Destination, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.db.GetDestinations(page, pageSize)
}

// GetItinerariesByDestination retrieves itineraries for a destination
func (s *Service) GetItinerariesByDestination(destinationID string, page, pageSize int) ([]Itinerary, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	itineraries, total, err := s.db.GetItinerariesByDestination(destinationID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get itineraries: %w", err)
	}

	// Enrich itineraries with items
	for i := range itineraries {
		items, err := s.db.GetItineraryItems(itineraries[i].ID)
		if err != nil {
			s.logger.Error("Failed to get items for itinerary: " + err.Error())
			continue
		}
		itineraries[i].Items = items
	}

	return itineraries, total, nil
}

// GetItineraryDetail retrieves a complete itinerary with all items
func (s *Service) GetItineraryDetail(itineraryID string) (*Itinerary, error) {
	return s.db.GetItineraryByID(itineraryID)
}

// CreateItinerary creates a new itinerary
func (s *Service) CreateItinerary(itinerary *Itinerary) error {
	if itinerary.Title == "" {
		return fmt.Errorf("title is required")
	}
	if itinerary.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	if itinerary.Budget <= 0 {
		return fmt.Errorf("budget must be greater than 0")
	}

	itinerary.CreatedAt = time.Now()
	itinerary.UpdatedAt = time.Now()
	itinerary.Likes = 0

	return s.db.CreateItinerary(itinerary)
}

// AddLikeToItinerary increments the like count for an itinerary
func (s *Service) AddLikeToItinerary(itineraryID string) error {
	if itineraryID == "" {
		return fmt.Errorf("itinerary ID is required")
	}
	return s.db.AddLikeToItinerary(itineraryID)
}

// AddComment adds a comment to an itinerary
func (s *Service) AddComment(comment *Comment) error {
	if comment.Content == "" {
		return fmt.Errorf("comment content is required")
	}

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return s.db.AddComment(comment)
}

// SearchItineraries searches itineraries by criteria
func (s *Service) SearchItineraries(query string, destinationID string, maxBudget float64, page, pageSize int) ([]Itinerary, int, error) {
	// TODO: Implement search logic
	return []Itinerary{}, 0, nil
}

// ==================== USER TRIP SERVICE METHODS ====================

// GetAllDestinations retrieves all destinations
func (s *Service) GetAllDestinations() ([]Destination, error) {
	destinations, _, err := s.db.GetDestinations(1, 100)
	return destinations, err
}

// CreateUserTrip creates a new user trip
func (s *Service) CreateUserTrip(trip *UserTrip) error {
	if trip.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if trip.Title == "" {
		return fmt.Errorf("title is required")
	}
	if trip.DestinationID == "" {
		return fmt.Errorf("destination_id is required")
	}
	if trip.Budget <= 0 {
		return fmt.Errorf("budget must be greater than 0")
	}
	if trip.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}

	trip.CreatedAt = time.Now()
	trip.UpdatedAt = time.Now()
	if trip.Status == "" {
		trip.Status = "draft"
	}

	return s.db.CreateUserTrip(trip)
}

// GetUserTrip retrieves a user trip by ID
func (s *Service) GetUserTrip(tripID string) (*UserTrip, error) {
	if tripID == "" {
		return nil, fmt.Errorf("trip_id is required")
	}
	return s.db.GetUserTripByID(tripID)
}

// GetUserTrips retrieves all trips for a user
func (s *Service) GetUserTrips(userID string) ([]UserTrip, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	return s.db.GetUserTripsByUserID(userID)
}

// UpdateUserTrip updates an existing user trip
func (s *Service) UpdateUserTrip(trip *UserTrip) error {
	if trip.ID == "" {
		return fmt.Errorf("trip_id is required")
	}
	if trip.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	trip.UpdatedAt = time.Now()
	return s.db.UpdateUserTrip(trip)
}

// DeleteUserTrip deletes a user trip
func (s *Service) DeleteUserTrip(tripID string) error {
	if tripID == "" {
		return fmt.Errorf("trip_id is required")
	}
	return s.db.DeleteUserTrip(tripID)
}

// AddTripSegment adds a segment to a trip
func (s *Service) AddTripSegment(segment *TripSegment) error {
	if segment.UserTripID == "" {
		return fmt.Errorf("user_trip_id is required")
	}
	if segment.Name == "" {
		return fmt.Errorf("name is required")
	}
	if segment.Day <= 0 {
		return fmt.Errorf("day must be greater than 0")
	}

	segment.CreatedAt = time.Now()
	segment.UpdatedAt = time.Now()

	return s.db.AddTripSegment(segment)
}

// AddTripPhoto adds a photo to a segment
func (s *Service) AddTripPhoto(photo *TripPhoto) error {
	if photo.TripSegmentID == "" {
		return fmt.Errorf("trip_segment_id is required")
	}
	if photo.URL == "" {
		return fmt.Errorf("url is required")
	}

	photo.UploadedAt = time.Now()
	return s.db.AddTripPhoto(photo)
}

// AddTripReview adds a review to a segment
func (s *Service) AddTripReview(review *TripReview) error {
	if review.TripSegmentID == "" {
		return fmt.Errorf("trip_segment_id is required")
	}
	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	return s.db.AddTripReview(review)
}

// PublishUserTrip publishes a trip as a community post
func (s *Service) PublishUserTrip(post *UserTripPost) error {
	if post.UserTripID == "" {
		return fmt.Errorf("user_trip_id is required")
	}
	if post.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if post.Title == "" {
		return fmt.Errorf("title is required")
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	now := time.Now()
	post.PublishedAt = &now
	post.Likes = 0
	post.Views = 0

	err := s.db.PublishUserTrip(post)
	if err == nil {
		// Update trip status to "published"
		s.db.UpdateUserTripStatus(post.UserTripID, "published")
	}
	return err
}

// GetCommunityPosts retrieves published community posts
func (s *Service) GetCommunityPosts(page, pageSize int) ([]UserTripPost, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.db.GetCommunityPosts(page, pageSize)
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(userID string) (*User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	s.logger.Debug("retrieving user", "user_id", userID)

	user, err := s.db.GetUserByID(userID)
	if err != nil {
		s.logger.Error("failed to get user", "user_id", userID, "error", err.Error())
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		s.logger.Warn("user not found", "user_id", userID)
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetDestination retrieves a destination by ID
func (s *Service) GetDestination(destinationID string) (*Destination, error) {
	if destinationID == "" {
		return nil, fmt.Errorf("destination_id is required")
	}

	s.logger.Debug("retrieving destination", "destination_id", destinationID)

	destination, err := s.db.GetDestinationByID(destinationID)
	if err != nil {
		s.logger.Error("failed to get destination", "destination_id", destinationID, "error", err.Error())
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	if destination == nil {
		s.logger.Warn("destination not found", "destination_id", destinationID)
		return nil, fmt.Errorf("destination not found")
	}

	return destination, nil
}
