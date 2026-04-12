package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PageData represents common page data for templates
type PageData struct {
	Title string
	Error string
}

// Handlers represents HTTP request handlers
type Handlers struct {
	service *Service
	logger  *Logger
	metrics *Metrics
}

// NewHandlers creates new handlers
func NewHandlers(service *Service, logger *Logger, metrics *Metrics) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
		metrics: metrics,
	}
}

// GetDestinations handles GET /api/destinations
func (h *Handlers) GetDestinations(c *gin.Context) {
	page := 1
	pageSize := 10

	// Parse query parameters
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		} else {
			h.metrics.RecordValidationError()
			h.logger.Warn("invalid_page_parameter", "value", p)
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		} else {
			h.metrics.RecordValidationError()
			h.logger.Warn("invalid_page_size_parameter", "value", ps)
		}
	}

	h.logger.Debug("fetching_destinations", "page", page, "page_size", pageSize)

	destinations, total, err := h.service.GetDestinations(page, pageSize)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetDestinations"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_get_destinations",
			"error", err.Error(),
			"page", page,
			"page_size", pageSize,
		)

		apiErr := NewDatabaseError("get_destinations", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	h.logger.Debug("destinations_fetched",
		"total", total,
		"page", page,
		"total_pages", totalPages,
	)

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       destinations,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetItinerariesByDestination handles GET /api/destinations/:destinationId/itineraries
func (h *Handlers) GetItinerariesByDestination(c *gin.Context) {
	destinationID := c.Param("destinationId")

	if destinationID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("destinationId", "destination ID is required")
		h.logger.Warn("missing_destination_id")
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

	h.logger.Debug("fetching_itineraries_by_destination",
		"destination_id", destinationID,
		"page", page,
		"page_size", pageSize,
	)

	itineraries, total, err := h.service.GetItinerariesByDestination(destinationID, page, pageSize)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetItinerariesByDestination"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_get_itineraries",
			"error", err.Error(),
			"destination_id", destinationID,
		)

		apiErr := NewDatabaseError("get_itineraries", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	h.logger.Debug("itineraries_fetched",
		"destination_id", destinationID,
		"total", total,
	)

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       itineraries,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetItineraryDetail handles GET /api/itineraries/:itineraryId
func (h *Handlers) GetItineraryDetail(c *gin.Context) {
	itineraryID := c.Param("itineraryId")

	if itineraryID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("itineraryId", "itinerary ID is required")
		h.logger.Warn("missing_itinerary_id")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("fetching_itinerary_detail", "itinerary_id", itineraryID)

	itinerary, err := h.service.GetItineraryDetail(itineraryID)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetItineraryDetail"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_get_itinerary_detail",
			"error", err.Error(),
			"itinerary_id", itineraryID,
		)

		apiErr := NewNotFoundError("Itinerary", itineraryID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// CreateItinerary handles POST /api/itineraries
func (h *Handlers) CreateItinerary(c *gin.Context) {
	var req Itinerary
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = uuid.New().String()

	if err := h.service.CreateItinerary(&req); err != nil {
		h.logger.Error("Failed to create itinerary: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create itinerary"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// LikeItinerary handles POST /api/itineraries/:itineraryId/like
func (h *Handlers) LikeItinerary(c *gin.Context) {
	itineraryID := c.Param("itineraryId")

	if itineraryID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("itineraryId", "itinerary ID is required")
		h.logger.Warn("missing_itinerary_id_for_like")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("adding_like_to_itinerary", "itinerary_id", itineraryID)

	if err := h.service.AddLikeToItinerary(itineraryID); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["AddLikeToItinerary"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_add_like",
			"error", err.Error(),
			"itinerary_id", itineraryID,
		)

		apiErr := NewDatabaseError("add_like", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.metrics.RecordLike()
	h.logger.Info("like_added_successfully", "itinerary_id", itineraryID)

	c.JSON(http.StatusOK, gin.H{"message": "Like added successfully"})
}

// CommentOnItinerary handles POST /api/itineraries/:itineraryId/comments
func (h *Handlers) CommentOnItinerary(c *gin.Context) {
	itineraryID := c.Param("itineraryId")

	if itineraryID == "" {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("itineraryId", "itinerary ID is required")
		h.logger.Warn("missing_itinerary_id_for_comment")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		h.metrics.RecordValidationError()
		h.logger.Warn("invalid_comment_payload", "error", err.Error())
		apiErr := NewValidationError("Invalid comment format", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	comment.ID = uuid.New().String()
	comment.ItineraryID = itineraryID

	h.logger.Debug("adding_comment_to_itinerary",
		"comment_id", comment.ID,
		"itinerary_id", itineraryID,
		"user_id", comment.UserID,
	)

	if err := h.service.AddComment(&comment); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["AddComment"]++
		h.metrics.mu.Unlock()

		h.logger.Error("failed_to_add_comment",
			"error", err.Error(),
			"comment_id", comment.ID,
			"itinerary_id", itineraryID,
		)

		apiErr := NewDatabaseError("add_comment", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.metrics.RecordCommentCreated()
	h.logger.Info("comment_added_successfully",
		"comment_id", comment.ID,
		"itinerary_id", itineraryID,
	)

	c.JSON(http.StatusCreated, comment)
}

// ==================== HTML Template Handlers ====================

// Index handles GET / - Home page with destinations
func (h *Handlers) Index(c *gin.Context) {
	page := 1
	pageSize := 12

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	destinations, total, err := h.service.GetDestinations(page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get destinations: " + err.Error())
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Failed to load destinations",
		})
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.HTML(http.StatusOK, "index.html", gin.H{
		"destinations": destinations,
		"currentPage":  page,
		"totalPages":   totalPages,
		"total":        total,
	})
}

// DestinationDetail handles GET /destination/:id - Show itineraries for destination
func (h *Handlers) DestinationDetail(c *gin.Context) {
	destinationID := c.Param("id")
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	itineraries, total, err := h.service.GetItinerariesByDestination(destinationID, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get itineraries: " + err.Error())
		c.HTML(http.StatusInternalServerError, "destination-detail.html", gin.H{
			"error": "Failed to load itineraries",
		})
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.HTML(http.StatusOK, "destination-detail.html", gin.H{
		"destinationID": destinationID,
		"itineraries":   itineraries,
		"currentPage":   page,
		"totalPages":    totalPages,
		"total":         total,
	})
}

// ItineraryDetail handles GET /itinerary/:id - Show complete itinerary with all items
func (h *Handlers) ItineraryDetail(c *gin.Context) {
	itineraryID := c.Param("id")

	itinerary, err := h.service.GetItineraryDetail(itineraryID)
	if err != nil {
		h.logger.Error("Failed to get itinerary: " + err.Error())
		c.HTML(http.StatusNotFound, "itinerary-detail.html", gin.H{
			"error": "Itinerary not found",
		})
		return
	}

	// Calculate total price
	totalPrice := 0.0
	for _, item := range itinerary.Items {
		totalPrice += item.Price
	}

	c.HTML(http.StatusOK, "itinerary-detail.html", gin.H{
		"itinerary":  itinerary,
		"totalPrice": totalPrice,
	})
}

// CreateItineraryPage handles GET /create - Show create itinerary form
func (h *Handlers) CreateItineraryPage(c *gin.Context) {
	destinations, _, err := h.service.GetDestinations(1, 50)
	if err != nil {
		h.logger.Error("Failed to get destinations: " + err.Error())
		destinations = []Destination{}
	}

	c.HTML(http.StatusOK, "create-itinerary.html", gin.H{
		"destinations": destinations,
	})
}

// CreateItinerarySubmit handles POST /create - Process uploaded itinerary
func (h *Handlers) CreateItinerarySubmit(c *gin.Context) {
	var req Itinerary
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "create-itinerary.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	req.ID = uuid.New().String()

	if err := h.service.CreateItinerary(&req); err != nil {
		h.logger.Error("Failed to create itinerary: " + err.Error())
		c.HTML(http.StatusInternalServerError, "create-itinerary.html", gin.H{
			"error": "Failed to create itinerary",
		})
		return
	}

	// Redirect to the new itinerary
	c.Redirect(http.StatusSeeOther, "/itinerary/"+req.ID)
}

// SearchPage handles GET /search - Search form and results
func (h *Handlers) SearchPage(c *gin.Context) {
	query := c.Query("q")
	destinationID := c.Query("destination")
	maxBudget := c.Query("max_budget")

	var results []Itinerary
	if query != "" {
		// TODO: Implement search
		results = []Itinerary{}
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"query":       query,
		"destination": destinationID,
		"maxBudget":   maxBudget,
		"results":     results,
	})
}

// ==================== NEW USER TRIP HANDLERS ====================

// LoginPage handles GET /login and GET /
func (h *Handlers) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// Dashboard handles GET /dashboard - Cities and trip planning hub
func (h *Handlers) Dashboard(c *gin.Context) {
	token := c.GetString("token")
	if token == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get destinations list
	destinations, err := h.service.GetAllDestinations()
	if err != nil {
		h.logger.Error("dashboard_fetch_destinations_error", "error", err.Error())
		destinations = []Destination{}
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":        "Trip Planning Dashboard",
		"destinations": destinations,
	})
}

// PlanTripPage handles GET /plan-trip - Trip planning wizard
func (h *Handlers) PlanTripPage(c *gin.Context) {
	token := c.GetString("token")
	if token == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	destinationID := c.Query("destination")
	destinations, err := h.service.GetAllDestinations()
	if err != nil {
		h.logger.Error("plan_trip_fetch_destinations_error", "error", err.Error())
		destinations = []Destination{}
	}

	c.HTML(http.StatusOK, "plan-trip.html", gin.H{
		"title":        "Plan Your Trip",
		"destination":  destinationID,
		"destinations": destinations,
	})
}

// MyTripsPage handles GET /my-trips - List user's trips
func (h *Handlers) MyTripsPage(c *gin.Context) {
	token := c.GetString("token")
	if token == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userID := c.GetString("user_id")
	trips, err := h.service.GetUserTrips(userID)
	if err != nil {
		h.logger.Error("my_trips_fetch_error", "error", err.Error(), "user_id", userID)
		trips = []UserTrip{}
	}

	c.HTML(http.StatusOK, "my-trips.html", gin.H{
		"title": "My Trips",
		"trips": trips,
	})
}

// MyTripDetail handles GET /my-trips/:id - Show specific trip
func (h *Handlers) MyTripDetail(c *gin.Context) {
	token := c.GetString("token")
	if token == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	tripID := c.Param("id")
	trip, err := h.service.GetUserTrip(tripID)
	if err != nil {
		h.logger.Error("trip_detail_fetch_error", "error", err.Error(), "trip_id", tripID)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Trip not found",
		})
		return
	}

	c.HTML(http.StatusOK, "trip-detail.html", gin.H{
		"title": trip.Title,
		"trip":  trip,
	})
}

// CommunityPage handles GET /community - Published trips feed
func (h *Handlers) CommunityPage(c *gin.Context) {
	posts, err := h.service.GetCommunityPosts(1, 20)
	if err != nil {
		h.logger.Error("community_fetch_posts_error", "error", err.Error())
		posts = []UserTripPost{}
	}

	c.HTML(http.StatusOK, "community.html", gin.H{
		"title": "Community Trips",
		"posts": posts,
	})
}

// ==================== USER TRIP API HANDLERS ====================

// CreateUserTrip handles POST /api/user-trips
func (h *Handlers) CreateUserTrip(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		apiErr := NewAuthenticationError("user_id not found in token")
		h.logger.Warn("create_user_trip_missing_token")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req UserTrip
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		h.logger.Warn("create_user_trip_invalid_input", "error", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.UserID = userID
	req.Status = "draft"

	if err := h.service.CreateUserTrip(&req); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["CreateUserTrip"]++
		h.metrics.mu.Unlock()

		h.logger.Error("create_user_trip_error", "error", err.Error(), "user_id", userID)
		apiErr := NewDatabaseError("create_user_trip", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("user_trip_created", "trip_id", req.ID, "user_id", userID)
	c.JSON(http.StatusCreated, req)
}

// GetUserTrip handles GET /api/user-trips/:id
func (h *Handlers) GetUserTrip(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	if tripID == "" {
		apiErr := NewInvalidInputError("id", "trip ID is required")
		h.logger.Warn("get_user_trip_missing_id")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	trip, err := h.service.GetUserTrip(tripID)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["GetUserTrip"]++
		h.metrics.mu.Unlock()

		h.logger.Error("get_user_trip_error", "error", err.Error(), "trip_id", tripID)
		apiErr := NewNotFoundError("trip", tripID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Check if user owns this trip
	if trip.UserID != userID && trip.Status != "published" {
		apiErr := NewAuthorizationError("you do not have permission to view this trip")
		h.logger.Warn("unauthorized_trip_access", "trip_id", tripID, "user_id", userID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, trip)
}

// UpdateUserTrip handles PUT /api/user-trips/:id
func (h *Handlers) UpdateUserTrip(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	if tripID == "" {
		apiErr := NewInvalidInputError("id", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req UserTrip
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify ownership
	trip, err := h.service.GetUserTrip(tripID)
	if err != nil || trip.UserID != userID {
		apiErr := NewAuthorizationError("you do not have permission to update this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = tripID
	req.UserID = userID

	if err := h.service.UpdateUserTrip(&req); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["UpdateUserTrip"]++
		h.metrics.mu.Unlock()

		h.logger.Error("update_user_trip_error", "error", err.Error(), "trip_id", tripID)
		apiErr := NewDatabaseError("update_user_trip", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("user_trip_updated", "trip_id", tripID, "user_id", userID)
	c.JSON(http.StatusOK, req)
}

// DeleteUserTrip handles DELETE /api/user-trips/:id
func (h *Handlers) DeleteUserTrip(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	if tripID == "" {
		apiErr := NewInvalidInputError("id", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify ownership
	trip, err := h.service.GetUserTrip(tripID)
	if err != nil || trip.UserID != userID {
		apiErr := NewAuthorizationError("you do not have permission to delete this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if err := h.service.DeleteUserTrip(tripID); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["DeleteUserTrip"]++
		h.metrics.mu.Unlock()

		h.logger.Error("delete_user_trip_error", "error", err.Error(), "trip_id", tripID)
		apiErr := NewDatabaseError("delete_user_trip", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("user_trip_deleted", "trip_id", tripID, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Trip deleted successfully"})
}

// ListUserTrips handles GET /api/user-trips
func (h *Handlers) ListUserTrips(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		apiErr := NewAuthenticationError("user_id not found")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	trips, err := h.service.GetUserTrips(userID)
	if err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["ListUserTrips"]++
		h.metrics.mu.Unlock()

		h.logger.Error("list_user_trips_error", "error", err.Error(), "user_id", userID)
		apiErr := NewDatabaseError("list_user_trips", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  trips,
		"count": len(trips),
	})
}

// AddTripSegment handles POST /api/user-trips/:id/segments
func (h *Handlers) AddTripSegment(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	if tripID == "" {
		apiErr := NewInvalidInputError("id", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify ownership
	trip, err := h.service.GetUserTrip(tripID)
	if err != nil || trip.UserID != userID {
		apiErr := NewAuthorizationError("you do not have permission to add segments to this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripSegment
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.UserTripID = tripID

	if err := h.service.AddTripSegment(&req); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["AddTripSegment"]++
		h.metrics.mu.Unlock()

		h.logger.Error("add_trip_segment_error", "error", err.Error(), "trip_id", tripID)
		apiErr := NewDatabaseError("add_trip_segment", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("trip_segment_added", "segment_id", req.ID, "trip_id", tripID)
	c.JSON(http.StatusCreated, req)
}

// AddTripPhoto handles POST /api/trip-segments/:id/photos
func (h *Handlers) AddTripPhoto(c *gin.Context) {
	segmentID := c.Param("id")
	userID := c.GetString("user_id")

	if segmentID == "" {
		apiErr := NewInvalidInputError("id", "segment ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripPhoto
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.TripSegmentID = segmentID

	if err := h.service.AddTripPhoto(&req); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["AddTripPhoto"]++
		h.metrics.mu.Unlock()

		h.logger.Error("add_trip_photo_error", "error", err.Error(), "segment_id", segmentID)
		apiErr := NewDatabaseError("add_trip_photo", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("trip_photo_added", "photo_id", req.ID, "segment_id", segmentID, "user_id", userID)
	c.JSON(http.StatusCreated, req)
}

// AddTripReview handles POST /api/trip-segments/:id/review
func (h *Handlers) AddTripReview(c *gin.Context) {
	segmentID := c.Param("id")
	userID := c.GetString("user_id")

	if segmentID == "" {
		apiErr := NewInvalidInputError("id", "segment ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req TripReview
	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		apiErr := NewInvalidInputError("rating", "rating must be between 1 and 5")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	req.ID = uuid.New().String()
	req.TripSegmentID = segmentID

	if err := h.service.AddTripReview(&req); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["AddTripReview"]++
		h.metrics.mu.Unlock()

		h.logger.Error("add_trip_review_error", "error", err.Error(), "segment_id", segmentID)
		apiErr := NewDatabaseError("add_trip_review", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("trip_review_added", "review_id", req.ID, "segment_id", segmentID, "user_id", userID)
	c.JSON(http.StatusCreated, req)
}

// PublishUserTrip handles POST /api/user-trips/:id/publish
func (h *Handlers) PublishUserTrip(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	if tripID == "" {
		apiErr := NewInvalidInputError("id", "trip ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify ownership
	trip, err := h.service.GetUserTrip(tripID)
	if err != nil || trip.UserID != userID {
		apiErr := NewAuthorizationError("you do not have permission to publish this trip")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		CoverImage  string `json:"cover_image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	post := UserTripPost{
		ID:          uuid.New().String(),
		UserTripID:  tripID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Published:   true,
	}

	if err := h.service.PublishUserTrip(&post); err != nil {
		h.metrics.mu.Lock()
		h.metrics.DatabaseQueryErrors["PublishUserTrip"]++
		h.metrics.mu.Unlock()

		h.logger.Error("publish_user_trip_error", "error", err.Error(), "trip_id", tripID)
		apiErr := NewDatabaseError("publish_user_trip", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("user_trip_published", "trip_id", tripID, "post_id", post.ID, "user_id", userID)
	c.JSON(http.StatusCreated, post)
}
