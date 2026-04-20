package itinerary

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/itinerary-backend/itinerary/common"
)

// Handlers handles HTTP requests
type Handlers struct {
	service *Service
	logger  *common.Logger
	metrics *Metrics
}

// NewHandlers creates new handlers
func NewHandlers(service *Service, logger *common.Logger, metrics *Metrics) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
		metrics: metrics,
	}
}

// Page handlers - stubs
func (h *Handlers) LoginPage(c *gin.Context) { c.JSON(200, gin.H{"page": "login"}) }
func (h *Handlers) Dashboard(c *gin.Context) { c.JSON(200, gin.H{"page": "dashboard"}) }
func (h *Handlers) PlanTripPage(c *gin.Context) { c.JSON(200, gin.H{"page": "plan-trip"}) }
func (h *Handlers) MyTripsPage(c *gin.Context) { c.JSON(200, gin.H{"page": "my-trips"}) }
func (h *Handlers) MyTripDetail(c *gin.Context) { c.JSON(200, gin.H{"page": "trip-detail"}) }
func (h *Handlers) CommunityPage(c *gin.Context) { c.JSON(200, gin.H{"page": "community"}) }
func (h *Handlers) DestinationDetail(c *gin.Context) { c.JSON(200, gin.H{"page": "destination-detail"}) }
func (h *Handlers) ItineraryDetail(c *gin.Context) { c.JSON(200, gin.H{"page": "itinerary-detail"}) }
func (h *Handlers) CreateItineraryPage(c *gin.Context) { c.JSON(200, gin.H{"page": "create-itinerary"}) }
func (h *Handlers) CreateItinerarySubmit(c *gin.Context) { c.JSON(200, gin.H{"message": "created"}) }
func (h *Handlers) SearchPage(c *gin.Context) { c.JSON(200, gin.H{"page": "search"}) }
func (h *Handlers) Index(c *gin.Context) { c.JSON(200, gin.H{"page": "index"}) }

// API handlers - implementations in api_*.go files
func (h *Handlers) GetDestinations(c *gin.Context) { c.JSON(200, []Destination{}) }
func (h *Handlers) GetItinerariesByDestination(c *gin.Context) { c.JSON(200, []Itinerary{}) }
func (h *Handlers) GetItineraryDetail(c *gin.Context) { c.JSON(200, nil) }
func (h *Handlers) CreateItinerary(c *gin.Context) { c.JSON(201, gin.H{"message": "created"}) }
func (h *Handlers) LikeItinerary(c *gin.Context) { c.JSON(200, gin.H{"message": "liked"}) }
func (h *Handlers) CommentOnItinerary(c *gin.Context) { c.JSON(201, gin.H{"message": "comment created"}) }
func (h *Handlers) GetUserTrip(c *gin.Context) { c.JSON(200, nil) }
func (h *Handlers) ListUserTrips(c *gin.Context) { c.JSON(200, []UserTrip{}) }

