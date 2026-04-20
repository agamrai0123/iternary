package itinerary

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== USER TRIPS CRUD HANDLERS ====================

// CreateUserTrip handles POST /api/user-trips
func (h *Handlers) CreateUserTrip(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req UserTrip
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_create_trip_request", "error", err.Error())
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.UserID = userID.(string)

	h.logger.Debug("creating_user_trip", "user_id", userID, "trip_title", req.Title)

	if err := h.service.CreateUserTrip(&req); err != nil {
		h.logger.Error("failed_to_create_user_trip", "error", err.Error())
		apiErr := NewInternalError("trip_creation_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GetUserTrips handles GET /api/user-trips
func (h *Handlers) GetUserTrips(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	h.logger.Debug("fetching_user_trips", "user_id", userID, "page", page)

	trips, total, err := h.service.GetUserTrips(userID.(string), page, pageSize)
	if err != nil {
		h.logger.Error("failed_to_get_user_trips", "error", err.Error())
		apiErr := NewDatabaseError("get_user_trips", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"data":        trips,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetUserTripByID handles GET /api/user-trips/:tripId
func (h *Handlers) GetUserTripByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	if tripID == "" {
		apiErr := NewInvalidInputError("tripId", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("fetching_user_trip", "trip_id", tripID, "user_id", userID)

	trip, err := h.service.GetUserTripByID(tripID)
	if err != nil || trip == nil {
		h.logger.Warn("user_trip_not_found", "trip_id", tripID)
		apiErr := NewNotFoundError("trip", tripID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify ownership
	if trip.UserID != userID.(string) {
		h.logger.Warn("unauthorized_trip_access", "trip_id", tripID, "user_id", userID)
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, trip)
}

// UpdateUserTrip handles PUT /api/user-trips/:tripId
func (h *Handlers) UpdateUserTrip(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	if tripID == "" {
		apiErr := NewInvalidInputError("tripId", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req UserTrip
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_update_trip_request", "error", err.Error())
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Get existing trip to verify ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("updating_user_trip", "trip_id", tripID, "user_id", userID)

	req.ID = tripID
	req.UserID = userID.(string)

	if err := h.service.UpdateUserTrip(&req); err != nil {
		h.logger.Error("failed_to_update_user_trip", "error", err.Error())
		apiErr := NewInternalError("trip_update_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteUserTrip handles DELETE /api/user-trips/:tripId
func (h *Handlers) DeleteUserTrip(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	if tripID == "" {
		apiErr := NewInvalidInputError("tripId", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("deleting_user_trip", "trip_id", tripID, "user_id", userID)

	if err := h.service.DeleteUserTrip(tripID); err != nil {
		h.logger.Error("failed_to_delete_user_trip", "error", err.Error())
		apiErr := NewInternalError("trip_deletion_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trip deleted successfully"})
}

// ==================== TRIP SEGMENTS HANDLERS ====================

// AddTripSegment handles POST /api/user-trips/:tripId/segments
func (h *Handlers) AddTripSegment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	if tripID == "" {
		apiErr := NewInvalidInputError("tripId", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripSegment
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_add_segment_request", "error", err.Error())
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify trip ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.UserTripID = tripID

	h.logger.Debug("adding_trip_segment", "trip_id", tripID, "segment_name", req.Name)

	if err := h.service.AddTripSegment(&req); err != nil {
		h.logger.Error("failed_to_add_trip_segment", "error", err.Error())
		apiErr := NewInternalError("segment_creation_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusCreated, req)
}

// UpdateTripSegment handles PUT /api/user-trips/:tripId/segments/:segmentId
func (h *Handlers) UpdateTripSegment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	segmentID := c.Param("segmentId")

	if tripID == "" || segmentID == "" {
		apiErr := NewInvalidInputError("params", "trip ID and segment ID are required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify trip ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripSegment
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_update_segment_request", "error", err.Error())
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = segmentID
	req.UserTripID = tripID

	h.logger.Debug("updating_trip_segment", "trip_id", tripID, "segment_id", segmentID)

	if err := h.service.UpdateTripSegment(&req); err != nil {
		h.logger.Error("failed_to_update_trip_segment", "error", err.Error())
		apiErr := NewInternalError("segment_update_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteTripSegment handles DELETE /api/user-trips/:tripId/segments/:segmentId
func (h *Handlers) DeleteTripSegment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	segmentID := c.Param("segmentId")

	if tripID == "" || segmentID == "" {
		apiErr := NewInvalidInputError("params", "trip ID and segment ID are required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify trip ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("deleting_trip_segment", "trip_id", tripID, "segment_id", segmentID)

	if err := h.service.DeleteTripSegment(segmentID); err != nil {
		h.logger.Error("failed_to_delete_trip_segment", "error", err.Error())
		apiErr := NewInternalError("segment_deletion_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "segment deleted successfully"})
}

// AddTripPhoto handles POST /api/user-trips/:tripId/segments/:segmentId/photos
func (h *Handlers) AddTripPhoto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	segmentID := c.Param("segmentId")

	if tripID == "" || segmentID == "" {
		apiErr := NewInvalidInputError("params", "trip ID and segment ID are required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify trip ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripPhoto
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.TripSegmentID = segmentID

	h.logger.Debug("adding_trip_photo", "segment_id", segmentID)

	if err := h.service.AddTripPhoto(&req); err != nil {
		h.logger.Error("failed_to_add_trip_photo", "error", err.Error())
		apiErr := NewInternalError("photo_upload_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusCreated, req)
}

// AddTripReview handles POST /api/user-trips/:tripId/segments/:segmentId/review
func (h *Handlers) AddTripReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	tripID := c.Param("tripId")
	segmentID := c.Param("segmentId")

	if tripID == "" || segmentID == "" {
		apiErr := NewInvalidInputError("params", "trip ID and segment ID are required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify trip ownership
	trip, _ := h.service.GetUserTripByID(tripID)
	if trip == nil || trip.UserID != userID.(string) {
		apiErr := NewAuthError("forbidden", "you do not have access to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripReview
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.TripSegmentID = segmentID

	h.logger.Debug("adding_trip_review", "segment_id", segmentID, "rating", req.Rating)

	if err := h.service.AddTripReview(&req); err != nil {
		h.logger.Error("failed_to_add_trip_review", "error", err.Error())
		apiErr := NewInternalError("review_submission_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusCreated, req)
}
