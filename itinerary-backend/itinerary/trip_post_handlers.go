package itinerary

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== TRIP POST COMMUNITY FEED HANDLERS ====================

// GetCitiesList handles GET /api/cities - Returns all destinations/cities
func (h *Handlers) GetCitiesList(c *gin.Context) {
	page := 1
	pageSize := 20

	// Parse query parameters
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	h.logger.Debug("fetching_cities_list", "page", page, "page_size", pageSize)

	destinations, total, err := h.service.GetDestinations(page, pageSize)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetCitiesList"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_get_cities_list", "error", err.Error())
		apiErr := NewDatabaseError("get_cities", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, ListCitiesResponse{
		Data:       destinations,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetTripPostsByCity handles GET /api/cities/:cityId/trip-posts
// Returns all trip posts (published community posts) for a specific city
func (h *Handlers) GetTripPostsByCity(c *gin.Context) {
	cityID := c.Param("cityId")

	if cityID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("cityId", "city ID is required")
		h.logger.Warn("missing_city_id")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	h.logger.Debug("fetching_trip_posts_by_city",
		"city_id", cityID,
		"page", page,
		"page_size", pageSize,
	)

	// Get published trip posts for this city
	posts, total, err := h.service.GetTripPostsByCity(cityID, page, pageSize)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetTripPostsByCity"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_get_trip_posts",
			"error", err.Error(),
			"city_id", cityID,
		)

		apiErr := NewDatabaseError("get_trip_posts", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	h.logger.Debug("trip_posts_fetched",
		"city_id", cityID,
		"total", total,
	)

	// Convert to response format with places
	responses := make([]TripPostResponse, len(posts))
	for i, post := range posts {
		responses[i] = convertTripPostToResponse(post)
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// AddTripPostToItinerary handles POST /api/user-trips/add-from-post
// User adds a published trip post to their own itinerary
func (h *Handlers) AddTripPostToItinerary(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		apiErr := NewAuthenticationError("user_id not found in token")
		h.logger.Warn("add_trip_post_missing_token")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req AddTripPostToItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		h.logger.Warn("add_trip_post_invalid_input", "error", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if req.TripPostID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("trip_post_id", "trip_post_id is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("adding_trip_post_to_itinerary",
		"user_id", userID,
		"trip_post_id", req.TripPostID,
	)

	// Get the trip post to copy
	tripPost, err := h.service.GetTripPostByID(req.TripPostID)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetTripPostByID"]++
		h.metrics.mu.Unlock()

		h.logger.Error("trip_post_not_found", "error", err.Error(), "trip_post_id", req.TripPostID)
		apiErr := NewNotFoundError("trip_post", req.TripPostID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create new user trip based on the post
	newTrip := &UserTrip{
		ID:            uuid.New().String(),
		UserID:        userID,
		Title:         tripPost.Title,
		DestinationID: tripPost.DestinationID,
		Budget:        tripPost.TotalExpense,
		Duration:      tripPost.Duration,
		Status:        "planning",
		Segments:      tripPost.Segments,
	}

	// Save the new trip
	if err := h.service.CreateUserTrip(newTrip); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["CreateUserTrip"]++
		h.metrics.mu.Unlock()

		h.logger.Error("create_user_trip_from_post_error",
			"error", err.Error(),
			"user_id", userID,
			"trip_post_id", req.TripPostID,
		)

		apiErr := NewDatabaseError("create_user_trip", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("trip_post_added_to_itinerary",
		"user_id", userID,
		"trip_id", newTrip.ID,
		"trip_post_id", req.TripPostID,
	)

	c.JSON(http.StatusCreated, newTrip)
}

// MarkSegmentVisited handles POST /api/trip-segments/:segmentId/mark-visited
// User marks a segment as visited and gets prompted for review
func (h *Handlers) MarkSegmentVisited(c *gin.Context) {
	segmentID := c.Param("segmentId")
	userID := c.GetString("user_id")

	if segmentID == "" {
		apiErr := NewInvalidInputError("segmentId", "segment ID is required")
		h.logger.Warn("mark_visited_missing_segment_id")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req MarkSegmentVisitedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("marking_segment_visited",
		"segment_id", segmentID,
		"user_id", userID,
		"completed", req.Completed,
	)

	// Get the segment to verify ownership
	segment, err := h.service.GetTripSegment(segmentID)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetTripSegment"]++
		h.metrics.mu.Unlock()

		h.logger.Error("segment_not_found", "error", err.Error(), "segment_id", segmentID)
		apiErr := NewNotFoundError("segment", segmentID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify user owns the trip
	trip, err := h.service.GetUserTrip(segment.UserTripID)
	if err != nil || trip.UserID != userID {
		apiErr := NewAuthorizationError("you do not have permission to update this segment")
		h.logger.Warn("unauthorized_segment_access",
			"segment_id", segmentID,
			"user_id", userID,
		)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Update segment completion status
	if err := h.service.UpdateSegmentCompletionStatus(segmentID, req.Completed); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["UpdateSegmentCompletionStatus"]++
		h.metrics.mu.Unlock()

		h.logger.Error("update_segment_completion_error",
			"error", err.Error(),
			"segment_id", segmentID,
		)

		apiErr := NewDatabaseError("update_segment", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("segment_marked_visited",
		"segment_id", segmentID,
		"user_id", userID,
	)

	// Return the segment with prompt to review
	c.JSON(http.StatusOK, gin.H{
		"message": "Place marked as visited",
		"segment": segment,
		"prompt":  "Please review this place to share your experience",
	})
}

// SubmitReview handles POST /api/reviews
// User submits a review for a visited place, which auto-generates a post
func (h *Handlers) SubmitReview(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		apiErr := NewAuthenticationError("user_id not found in token")
		h.logger.Warn("submit_review_missing_token")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		h.logger.Warn("submit_review_invalid_input", "error", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if req.SegmentID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("segment_id", "segment_id is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("submitting_review",
		"user_id", userID,
		"segment_id", req.SegmentID,
		"rating", req.Rating,
	)

	// Get the segment
	segment, err := h.service.GetTripSegment(req.SegmentID)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetTripSegment"]++
		h.metrics.mu.Unlock()

		h.logger.Error("segment_not_found_for_review",
			"error", err.Error(),
			"segment_id", req.SegmentID,
		)

		apiErr := NewNotFoundError("segment", req.SegmentID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Get the trip to verify ownership
	trip, err := h.service.GetUserTrip(segment.UserTripID)
	if err != nil || trip.UserID != userID {
		apiErr := NewAuthorizationError("you do not have permission to review this segment")
		h.logger.Warn("unauthorized_review_access",
			"segment_id", req.SegmentID,
			"user_id", userID,
		)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create and save the review
	review := &TripReview{
		ID:            uuid.New().String(),
		TripSegmentID: req.SegmentID,
		Rating:        req.Rating,
		Review:        req.Review,
	}

	if err := h.service.AddTripReview(review); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["AddTripReview"]++
		h.metrics.mu.Unlock()

		h.logger.Error("add_review_error",
			"error", err.Error(),
			"segment_id", req.SegmentID,
		)

		apiErr := NewDatabaseError("add_review", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("review_submitted",
		"review_id", review.ID,
		"segment_id", req.SegmentID,
		"user_id", userID,
	)

	// The review is now automatically part of the trip post
	// No need for separate "publish" - reviews are part of the segments in the trip post
	c.JSON(http.StatusCreated, gin.H{
		"message": "Review submitted successfully and added to your trip post",
		"review":  review,
	})
}

// ==================== HELPER FUNCTIONS ====================

// convertTripPostToResponse converts a UserTripPost to TripPostResponse format
func convertTripPostToResponse(post UserTripPost) TripPostResponse {
	places := make([]TripPlaceResponse, len(post.Segments))

	for i, segment := range post.Segments {
		placeResp := TripPlaceResponse{
			ID:              segment.ID,
			Day:             segment.Day,
			TimeOfDay:       segment.TimeOfDay,
			Name:            segment.Name,
			Type:            segment.Type,
			Location:        segment.Location,
			Latitude:        segment.Latitude,
			Longitude:       segment.Longitude,
			Expense:         segment.Expense,
			BestTimeToVisit: segment.BestTimeToVisit,
			Photos:          segment.Photos,
		}

		// Convert review if exists
		if segment.Review != nil {
			placeResp.Review = &TripReviewResponse{
				Rating: segment.Review.Rating,
				Review: segment.Review.Review,
			}
		}

		places[i] = placeResp
	}

	return TripPostResponse{
		ID:            post.ID,
		UserID:        post.UserID,
		Title:         post.Title,
		Description:   post.Description,
		CoverImage:    post.CoverImage,
		DestinationID: post.DestinationID,
		Duration:      post.Duration,
		TotalExpense:  post.TotalExpense,
		Places:        places,
		Likes:         post.Likes,
		Views:         post.Views,
		PublishedAt:   post.PublishedAt,
		CreatedAt:     post.CreatedAt,
	}
}
