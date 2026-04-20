package common

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Database represents the database connection
// Moved from package itinerary
// Only types and methods that do not depend on itinerary.Config or itinerary.Logger should remain here.
type Database struct {
	Conn *sql.DB
}
