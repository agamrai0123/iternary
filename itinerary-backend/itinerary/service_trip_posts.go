package itinerary

import "fmt"

// ==================== TRIP POST COMMUNITY SERVICE METHODS ====================

// GetTripPostsByCity retrieves trip posts for a specific city/destination
func (s *Service) GetTripPostsByCity(cityID string, page, pageSize int) ([]UserTripPost, int, error) {
	if cityID == "" {
		return nil, 0, fmt.Errorf("city_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get all community posts and filter by destination
	allPosts, total, err := s.db.GetCommunityPosts(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get trip posts: %w", err)
	}

	// Filter posts by destination
	filteredPosts := make([]UserTripPost, 0)
	for _, post := range allPosts {
		if post.DestinationID == cityID {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts, len(filteredPosts), nil
}

// GetTripPostByID retrieves a single trip post by ID
func (s *Service) GetTripPostByID(tripPostID string) (*UserTripPost, error) {
	if tripPostID == "" {
		return nil, fmt.Errorf("trip_post_id is required")
	}

	// Get all community posts and find the one with matching ID
	posts, _, err := s.db.GetCommunityPosts(1, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip post: %w", err)
	}

	for _, post := range posts {
		if post.ID == tripPostID {
			return &post, nil
		}
	}

	return nil, fmt.Errorf("trip post not found")
}

// GetTripSegment retrieves a trip segment by ID
func (s *Service) GetTripSegment(segmentID string) (*TripSegment, error) {
	if segmentID == "" {
		return nil, fmt.Errorf("segment_id is required")
	}

	// Query database for segment by ID
	// This needs implementation in the database layer
	// For now, placeholder that will be implemented later
	return s.db.GetTripSegmentByID(segmentID)
}

// UpdateSegmentCompletionStatus updates whether a segment has been visited
func (s *Service) UpdateSegmentCompletionStatus(segmentID string, completed bool) error {
	if segmentID == "" {
		return fmt.Errorf("segment_id is required")
	}

	// Update segment completion status in database
	return s.db.UpdateTripSegmentCompletion(segmentID, completed)
}
