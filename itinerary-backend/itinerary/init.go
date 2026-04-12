package itinerary

// This file re-exports functions from subpackages to provide a clean API from the itinerary package

import (
	authpkg "github.com/yourusername/itinerary-backend/itinerary/auth"
	"github.com/yourusername/itinerary-backend/itinerary/config"
	"github.com/yourusername/itinerary-backend/itinerary/middleware"
	"github.com/yourusername/itinerary-backend/itinerary/routes"
	"github.com/yourusername/itinerary-backend/itinerary/service"
	"github.com/yourusername/itinerary-backend/itinerary/utils"
)

// Re-export types and functions from config package
type Config = config.Config

var LoadConfig = config.LoadConfig

// Re-export types and functions from utils package
type Logger = utils.Logger

var NewLogger = utils.NewLogger

// Re-export types and functions from middleware package
type Metrics = middleware.Metrics

var NewMetrics = middleware.NewMetrics

// Re-export types and functions from service package
type Database = service.Database
type Service = service.Service

var NewDatabase = service.NewDatabase
var NewService = service.NewService

// Re-export types and functions from auth package
type AuthService = authpkg.AuthService

var NewAuthService = authpkg.NewAuthService

// Re-export types and functions from routes package
var SetupRoutes = routes.SetupRoutes
