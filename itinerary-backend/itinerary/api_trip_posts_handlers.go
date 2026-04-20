package itinerary

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== CITIES API HANDLERS ====================

// GetCities handles GET /api/cities
func (h *Handlers) GetCities(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	h.logger.Debug("fetching_cities", "page", page, "page_size", pageSize)

	destinations, total, err := h.service.GetDestinations(page, pageSize)
	if err != nil {
		h.logger.Error("failed_to_get_cities", "error", err.Error())
		apiErr := NewDatabaseError("get_cities", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"data":        destinations,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetCityByID handles GET /api/cities/:cityId
func (h *Handlers) GetCityByID(c *gin.Context) {
	cityID := c.Param("cityId")

	if cityID == "" {
		apiErr := NewInvalidInputError("cityId", "city ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("fetching_city", "city_id", cityID)

	destination, err := h.service.GetDestinationByID(cityID)
	if err != nil || destination == nil {
		h.logger.Warn("city_not_found", "city_id", cityID)
		apiErr := NewNotFoundError("city", cityID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, destination)
}

// ==================== TRIP POSTS API HANDLERS ====================

// GetTripPosts handles GET /api/trip-posts
func (h *Handlers) GetTripPosts(c *gin.Context) {
	page := 1
	pageSize := 10
	cityID := c.Query("city_id")

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	h.logger.Debug("fetching_trip_posts", "city_id", cityID, "page", page)

	var posts []UserTripPost
	var total int
	var err error

	if cityID != "" {
		posts, total, err = h.service.GetTripPostsByCity(cityID, page, pageSize)
	} else {
		posts, total, err = h.service.GetAllTripPosts(page, pageSize)
	}

	if err != nil {
		h.logger.Error("failed_to_get_trip_posts", "error", err.Error())
		apiErr := NewDatabaseError("get_trip_posts", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"data":        posts,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetTripPostByID handles GET /api/trip-posts/:postId
func (h *Handlers) GetTripPostByID(c *gin.Context) {
	postID := c.Param("postId")

	if postID == "" {
		apiErr := NewInvalidInputError("postId", "post ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("fetching_trip_post", "post_id", postID)

	post, err := h.service.GetTripPostByID(postID)
	if err != nil || post == nil {
		h.logger.Warn("trip_post_not_found", "post_id", postID)
		apiErr := NewNotFoundError("trip_post", postID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Record view
	h.service.RecordTripPostView(postID)

	c.JSON(http.StatusOK, post)
}

// GetCityTripPosts handles GET /api/cities/:cityId/trip-posts
func (h *Handlers) GetCityTripPosts(c *gin.Context) {
	cityID := c.Param("cityId")
	page := 1
	pageSize := 10

	if cityID == "" {
		apiErr := NewInvalidInputError("cityId", "city ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	h.logger.Debug("fetching_city_trip_posts", "city_id", cityID, "page", page)

	posts, total, err := h.service.GetTripPostsByCity(cityID, page, pageSize)
	if err != nil {
		h.logger.Error("failed_to_get_city_trip_posts", "error", err.Error())
		apiErr := NewDatabaseError("get_trip_posts", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"data":        posts,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// LikeTripPost handles POST /api/trip-posts/:postId/like
func (h *Handlers) LikeTripPost(c *gin.Context) {
	postID := c.Param("postId")

	if postID == "" {
		apiErr := NewInvalidInputError("postId", "post ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("liking_trip_post", "post_id", postID)

	if err := h.service.IncrementTripPostLikes(postID); err != nil {
		h.logger.Error("failed_to_like_trip_post", "error", err.Error())
		apiErr := NewInternalServerError("like_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "liked successfully"})
}

// SaveTripPost handles POST /api/trip-posts/:postId/save
func (h *Handlers) SaveTripPost(c *gin.Context) {
	postID := c.Param("postId")
	userID, exists := c.Get("user_id")
	
	if !exists {
		apiErr := NewAuthenticationError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if postID == "" {
		apiErr := NewInvalidInputError("postId", "post ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("saving_trip_post", "post_id", postID, "user_id", userID)

	if err := h.service.SaveTripPost(userID.(string), postID); err != nil {
		h.logger.Error("failed_to_save_trip_post", "error", err.Error())
		apiErr := NewInternalServerError("save_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trip saved successfully"})
}

// AddTripPostToItinerary handles POST /api/trip-posts/:postId/add-to-itinerary
func (h *Handlers) AddTripPostToItinerary(c *gin.Context) {
	postID := c.Param("postId")
	userID, exists := c.Get("user_id")

	if !exists {
		apiErr := NewAuthenticationError("unauthorized", "user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	if postID == "" {
		apiErr := NewInvalidInputError("postId", "post ID is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("adding_trip_to_itinerary", "post_id", postID, "user_id", userID)

	// Get the trip post
	post, err := h.service.GetTripPostByID(postID)
	if err != nil || post == nil {
		apiErr := NewNotFoundError("trip_post", postID)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create a new user trip from the post
	newTrip := &UserTrip{
		ID:            uuid.New().String(),
		UserID:        userID.(string),
		Title:         post.Title,
		DestinationID: post.DestinationID,
		Budget:        post.TotalExpense,
		Duration:      post.Duration,
		Status:        "planning",
	}

	// Copy segments
	for _, segment := range post.Segments {
		newTrip.Segments = append(newTrip.Segments, segment)
	}

	if err := h.service.CreateUserTrip(newTrip); err != nil {
		h.logger.Error("failed_to_create_user_trip", "error", err.Error())
		apiErr := NewInternalServerError("trip_creation_failed")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "trip added to itinerary successfully",
		"trip_id": newTrip.ID,
	})
}
