package itinerary

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Database represents the database connection
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new SQLite database connection
func NewDatabase(config *Config, logger *Logger) (*Database, error) {
	// SQLite connection - uses file-based database
	dbPath := "itinerary.db"

	logger.Debug("connecting to database", "path", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Error("failed to open database", "error", err.Error(), "path", dbPath)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err.Error())
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	logger.Info("SQLite database connection pool initialized", "max_open_conns", 25, "max_idle_conns", 5)

	// Initialize schema
	if err := initializeSchema(db, logger); err != nil {
		logger.Error("failed to initialize schema", "error", err.Error())
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Insert test data if not exists
	if err := insertTestData(db, logger); err != nil {
		logger.Error("failed to insert test data", "error", err.Error())
		return nil, fmt.Errorf("failed to insert test data: %w", err)
	}

	return &Database{conn: db}, nil
}

// initializeSchema creates all tables if they don't exist
func initializeSchema(db *sql.DB, logger *Logger) error {
	logger.Debug("initializing database schema")
	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Destinations table
	CREATE TABLE IF NOT EXISTS destinations (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		country TEXT NOT NULL,
		description TEXT,
		image_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Itineraries table
	CREATE TABLE IF NOT EXISTS itineraries (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		destination_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		duration INTEGER NOT NULL,
		budget REAL NOT NULL,
		likes INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (destination_id) REFERENCES destinations(id)
	);

	-- Itinerary items table
	CREATE TABLE IF NOT EXISTS itinerary_items (
		id TEXT PRIMARY KEY,
		itinerary_id TEXT NOT NULL,
		day INTEGER NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		duration INTEGER,
		location TEXT,
		rating REAL,
		image_url TEXT,
		booking_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
		CHECK (type IN ('stay', 'food', 'activity', 'transport', 'other'))
	);

	-- Comments table
	CREATE TABLE IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		itinerary_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		rating REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- User plans table
	CREATE TABLE IF NOT EXISTS user_plans (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		original_itinerary_id TEXT NOT NULL,
		title TEXT,
		notes TEXT,
		status TEXT DEFAULT 'draft',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (original_itinerary_id) REFERENCES itineraries(id),
		CHECK (status IN ('draft', 'planned', 'completed'))
	);

	-- Likes table
	CREATE TABLE IF NOT EXISTS likes (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		itinerary_id TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (user_id, itinerary_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE
	);

	-- User Trips table (custom trip planning)
	CREATE TABLE IF NOT EXISTS user_trips (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		destination_id TEXT NOT NULL,
		budget REAL NOT NULL,
		duration INTEGER NOT NULL,
		start_date TIMESTAMP,
		status TEXT DEFAULT 'draft',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (destination_id) REFERENCES destinations(id),
		CHECK (status IN ('draft', 'planning', 'ongoing', 'completed'))
	);

	-- Trip Segments table (places/activities in a trip)
	CREATE TABLE IF NOT EXISTS trip_segments (
		id TEXT PRIMARY KEY,
		user_trip_id TEXT NOT NULL,
		day INTEGER NOT NULL,
		name TEXT NOT NULL,
		type TEXT,
		location TEXT,
		latitude REAL,
		longitude REAL,
		notes TEXT,
		completed BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_trip_id) REFERENCES user_trips(id) ON DELETE CASCADE
	);

	-- Trip Photos table (photos for trip segments)
	CREATE TABLE IF NOT EXISTS trip_photos (
		id TEXT PRIMARY KEY,
		trip_segment_id TEXT NOT NULL,
		url TEXT NOT NULL,
		caption TEXT,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (trip_segment_id) REFERENCES trip_segments(id) ON DELETE CASCADE
	);

	-- Trip Reviews table (reviews for completed segments)
	CREATE TABLE IF NOT EXISTS trip_reviews (
		id TEXT PRIMARY KEY,
		trip_segment_id TEXT NOT NULL,
		rating REAL NOT NULL,
		review TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (trip_segment_id) REFERENCES trip_segments(id) ON DELETE CASCADE,
		CHECK (rating >= 1 AND rating <= 5)
	);

	-- User Trip Posts table (published community posts)
	CREATE TABLE IF NOT EXISTS user_trip_posts (
		id TEXT PRIMARY KEY,
		user_trip_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		title TEXT,
		description TEXT,
		cover_image TEXT,
		likes INTEGER DEFAULT 0,
		views INTEGER DEFAULT 0,
		published BOOLEAN DEFAULT 0,
		published_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_trip_id) REFERENCES user_trips(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);


	-- Create indexes
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_destinations_country ON destinations(country);
	CREATE INDEX IF NOT EXISTS idx_itineraries_user_id ON itineraries(user_id);
	CREATE INDEX IF NOT EXISTS idx_itineraries_destination_id ON itineraries(destination_id);
	CREATE INDEX IF NOT EXISTS idx_itineraries_likes ON itineraries(likes DESC);
	CREATE INDEX IF NOT EXISTS idx_itineraries_created_at ON itineraries(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_itinerary_items_day ON itinerary_items(itinerary_id, day);
	CREATE INDEX IF NOT EXISTS idx_itinerary_items_type ON itinerary_items(type);
	CREATE INDEX IF NOT EXISTS idx_comments_itinerary_id ON comments(itinerary_id);
	CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
	CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_user_plans_user_id ON user_plans(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_plans_status ON user_plans(status);
	CREATE INDEX IF NOT EXISTS idx_likes_itinerary_id ON likes(itinerary_id);

	-- Indexes for new trip planning tables
	CREATE INDEX IF NOT EXISTS idx_user_trips_user_id ON user_trips(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_trips_status ON user_trips(status);
	CREATE INDEX IF NOT EXISTS idx_user_trips_destination_id ON user_trips(destination_id);
	CREATE INDEX IF NOT EXISTS idx_trip_segments_user_trip_id ON trip_segments(user_trip_id);
	CREATE INDEX IF NOT EXISTS idx_trip_segments_day ON trip_segments(user_trip_id, day);
	CREATE INDEX IF NOT EXISTS idx_trip_photos_segment_id ON trip_photos(trip_segment_id);
	CREATE INDEX IF NOT EXISTS idx_trip_reviews_segment_id ON trip_reviews(trip_segment_id);
	CREATE INDEX IF NOT EXISTS idx_user_trip_posts_user_id ON user_trip_posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_trip_posts_published ON user_trip_posts(published);
	`

	_, err := db.Exec(schema)
	if err != nil {
		logger.Error("failed to execute schema", "error", err.Error())
		return err
	}
	logger.Info("database schema initialized successfully")
	return nil
}

// insertTestData inserts test data if tables are empty
func insertTestData(db *sql.DB, logger *Logger) error {
	logger.Debug("inserting test data into database")
	// Check if users already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Data already exists
	}

	// Insert users
	users := []struct {
		id, username, email string
	}{
		{"user-001", "traveler1", "traveler1@example.com"},
		{"user-002", "explorer2", "explorer@example.com"},
		{"user-003", "wanderer3", "wanderer@example.com"},
	}

	for _, u := range users {
		_, err := db.Exec("INSERT INTO users (id, username, email) VALUES (?, ?, ?)",
			u.id, u.username, u.email)
		if err != nil {
			return err
		}
	}

	// Insert destinations
	destinations := []struct {
		id, name, country, description string
	}{
		{"dest-001", "Goa", "India", "Beautiful coastal state with beaches, churches and Portuguese heritage. Known for nightlife, water sports and seafood restaurants."},
		{"dest-002", "Manali", "India", "Hill station in Himachal Pradesh. Perfect for adventure sports, trekking, and scenic mountain views. Great for families and adventure seekers."},
		{"dest-003", "Bali", "Indonesia", "Tropical island known for beaches, rice terraces, temples and vibrant culture. Budget-friendly and great for all types of travelers."},
	}

	for _, d := range destinations {
		_, err := db.Exec("INSERT INTO destinations (id, name, country, description) VALUES (?, ?, ?, ?)",
			d.id, d.name, d.country, d.description)
		if err != nil {
			return err
		}
	}

	// Insert itineraries
	itineraries := []struct {
		id, userId, destId, title, description string
		duration                               int
		budget                                 float64
		likes                                  int
	}{
		{"itin-001", "user-001", "dest-001", "5-Day Budget Goa Trip", "Perfect 5 days in Goa with beaches, culture, and food. Includes Baga Beach, Old Goa churches, and local cuisine. Includes 3-star accommodation.", 5, 15000, 45},
		{"itin-002", "user-002", "dest-001", "Luxury 7-Day Goa Experience", "High-end Goa itinerary with 5-star resorts, adventure activities and fine dining. Includes water sports, spa and sunset cruises.", 7, 45000, 32},
		{"itin-003", "user-003", "dest-002", "4-Day Manali Adventure", "Adventure-focused trip with trekking, paragliding, and camping. Includes visits to Rohtang Pass and local villages.", 4, 12000, 28},
		{"itin-004", "user-001", "dest-003", "6-Day Bali Paradise", "Comprehensive Bali experience with beaches, temples, rice paddies and nightlife. Perfect for first-time visitors.", 6, 18000, 67},
	}

	for _, i := range itineraries {
		_, err := db.Exec("INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, likes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			i.id, i.userId, i.destId, i.title, i.description, i.duration, i.budget, i.likes)
		if err != nil {
			return err
		}
	}

	// Insert itinerary items
	items := []struct {
		id, itinId, itemType, name, description, location string
		day, duration                                     int
		price, rating                                     float64
	}{
		{"item-001", "itin-001", "stay", "Beachside Hostel", "3-star hostel with AC rooms, common kitchen, and beach access", "Baga", 1, 24, 500, 4.5},
		{"item-002", "itin-001", "food", "Seafood Dinner Baga", "Fresh fish curry, prawn biryani and local wines at beachside restaurant", "Baga Beach", 1, 2, 350, 4.7},
		{"item-003", "itin-001", "activity", "Jet Ski Adventure", "Water sports activity including jet ski, banana boat ride", "Baga Beach", 2, 2, 800, 4.8},
		{"item-004", "itin-001", "food", "Beachside Lunch", "Grilled fish with rice and salad at beach shack", "Anjuna Beach", 2, 1, 250, 4.6},
		{"item-005", "itin-001", "activity", "Old Goa Heritage Tour", "Visit Se Cathedral, Basilica of Bom Jesus and other historic churches", "Old Goa", 3, 3, 300, 4.4},
		{"item-006", "itin-001", "transport", "Scooter Rental", "Rent a scooter to explore Goa independently", "Baga", 4, 48, 400, 4.3},
		{"item-101", "itin-004", "stay", "Luxury Resort Seminyak", "Beachfront 5-star resort with private pool and ocean views", "Seminyak", 1, 24, 3000, 4.9},
		{"item-102", "itin-004", "food", "Fine Dining Dinner", "Indonesian and international cuisine at beach club restaurant", "Seminyak Beach", 1, 2, 800, 4.8},
		{"item-103", "itin-004", "activity", "Ubud Rice Paddies Tour", "Explore scenic rice terraces and local villages", "Ubud", 2, 4, 400, 4.7},
		{"item-104", "itin-004", "activity", "Temple Tour", "Visit Tanah Lot and Uluwatu Temple at sunset", "Various", 3, 5, 350, 4.6},
	}

	for _, item := range items {
		_, err := db.Exec("INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			item.id, item.itinId, item.day, item.itemType, item.name, item.description, item.price, item.duration, item.location, item.rating)
		if err != nil {
			return err
		}
	}

	// Insert comments
	comments := []struct {
		id, itinId, userId, content string
		rating                      float64
	}{
		{"comment-001", "itin-001", "user-002", "Amazing itinerary! Did exactly this and had the best time. Jet ski was super fun and the food recommendations were spot on.", 5},
		{"comment-002", "itin-001", "user-003", "Great value for money. Would definitely recommend this plan to budget travelers. The beachside atmosphere is incredible.", 4.5},
		{"comment-003", "itin-004", "user-001", "Luxury experience but worth every penny. Bali is truly paradise. The temples are magnificent and people are super friendly.", 5},
	}

	for _, c := range comments {
		_, err := db.Exec("INSERT INTO comments (id, itinerary_id, user_id, content, rating) VALUES (?, ?, ?, ?, ?)",
			c.id, c.itinId, c.userId, c.content, c.rating)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.conn.Close()
}

// GetDestinations retrieves all destinations with pagination
func (d *Database) GetDestinations(page, pageSize int) ([]Destination, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := d.conn.QueryRow("SELECT COUNT(*) FROM destinations").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count destinations: %w", err)
	}

	// Get paginated results - SQLite uses LIMIT/OFFSET syntax
	rows, err := d.conn.Query(
		"SELECT id, name, country, description, image_url, created_at, updated_at FROM destinations ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query destinations: %w", err)
	}
	defer rows.Close()

	var destinations []Destination
	for rows.Next() {
		var d Destination
		var imageURL *string
		if err := rows.Scan(&d.ID, &d.Name, &d.Country, &d.Description, &imageURL, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if imageURL != nil {
			d.Image = *imageURL
		}
		destinations = append(destinations, d)
	}

	return destinations, total, rows.Err()
}

// GetItinerariesByDestination retrieves itineraries for a destination sorted by likes
func (d *Database) GetItinerariesByDestination(destinationID string, page, pageSize int) ([]Itinerary, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	query := "SELECT COUNT(*) FROM itineraries WHERE destination_id = ?"
	err := d.conn.QueryRow(query, destinationID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count itineraries: %w", err)
	}

	// Get paginated results sorted by likes - SQLite syntax
	rows, err := d.conn.Query(
		"SELECT id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at FROM itineraries WHERE destination_id = ? ORDER BY likes DESC, created_at DESC LIMIT ? OFFSET ?",
		destinationID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query itineraries: %w", err)
	}
	defer rows.Close()

	var itineraries []Itinerary
	for rows.Next() {
		var i Itinerary
		if err := rows.Scan(&i.ID, &i.UserID, &i.DestinationID, &i.Title, &i.Description, &i.Duration, &i.Budget, &i.Likes, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, 0, err
		}
		itineraries = append(itineraries, i)
	}

	return itineraries, total, rows.Err()
}

// CreateItinerary creates a new itinerary
func (d *Database) CreateItinerary(itinerary *Itinerary) error {
	query := `INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := d.conn.Exec(query, itinerary.ID, itinerary.UserID, itinerary.DestinationID, itinerary.Title, itinerary.Description, itinerary.Duration, itinerary.Budget)
	if err != nil {
		return fmt.Errorf("failed to create itinerary: %w", err)
	}
	return nil
}

// GetItineraryItems retrieves all items for an itinerary
func (d *Database) GetItineraryItems(itineraryID string) ([]ItineraryItem, error) {
	rows, err := d.conn.Query(
		"SELECT id, itinerary_id, day, type, name, description, price, duration, location, rating, image_url, booking_url, created_at, updated_at FROM itinerary_items WHERE itinerary_id = ? ORDER BY day, created_at",
		itineraryID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query itinerary items: %w", err)
	}
	defer rows.Close()

	var items []ItineraryItem
	for rows.Next() {
		var item ItineraryItem
		if err := rows.Scan(&item.ID, &item.ItineraryID, &item.Day, &item.Type, &item.Name, &item.Description, &item.Price, &item.Duration, &item.Location, &item.Rating, &item.ImageURL, &item.BookingURL, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetItineraryByID retrieves a complete itinerary with all items
func (d *Database) GetItineraryByID(itineraryID string) (*Itinerary, error) {
	query := `SELECT id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at 
			 FROM itineraries WHERE id = ?`

	var itinerary Itinerary
	err := d.conn.QueryRow(query, itineraryID).Scan(
		&itinerary.ID, &itinerary.UserID, &itinerary.DestinationID, &itinerary.Title,
		&itinerary.Description, &itinerary.Duration, &itinerary.Budget, &itinerary.Likes,
		&itinerary.CreatedAt, &itinerary.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get itinerary: %w", err)
	}

	// Get items
	items, err := d.GetItineraryItems(itineraryID)
	if err != nil {
		return nil, err
	}
	itinerary.Items = items

	return &itinerary, nil
}

// AddLikeToItinerary increments the likes count
func (d *Database) AddLikeToItinerary(itineraryID string) error {
	query := "UPDATE itineraries SET likes = likes + 1 WHERE id = ?"
	_, err := d.conn.Exec(query, itineraryID)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

// AddComment adds a comment to an itinerary
func (d *Database) AddComment(comment *Comment) error {
	query := `INSERT INTO comments (id, itinerary_id, user_id, content, rating, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := d.conn.Exec(query, comment.ID, comment.ItineraryID, comment.UserID, comment.Content, comment.Rating)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}

// GetCommentsByItinerary retrieves comments for an itinerary
func (d *Database) GetCommentsByItinerary(itineraryID string) ([]Comment, error) {
	rows, err := d.conn.Query(
		"SELECT id, itinerary_id, user_id, content, rating, created_at, updated_at FROM comments WHERE itinerary_id = ? ORDER BY created_at DESC",
		itineraryID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.ItineraryID, &c.UserID, &c.Content, &c.Rating, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}

// ==================== USER TRIP DATABASE METHODS ====================

// CreateUserTrip creates a new user trip
func (d *Database) CreateUserTrip(trip *UserTrip) error {
	query := `INSERT INTO user_trips (id, user_id, title, destination_id, budget, duration, start_date, status, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query, trip.ID, trip.UserID, trip.Title, trip.DestinationID, trip.Budget, trip.Duration, trip.StartDate, trip.Status, trip.CreatedAt, trip.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user trip: %w", err)
	}
	return nil
}

// GetUserTripByID retrieves a user trip by ID
func (d *Database) GetUserTripByID(tripID string) (*UserTrip, error) {
	query := `SELECT id, user_id, title, destination_id, budget, duration, start_date, status, created_at, updated_at
			 FROM user_trips WHERE id = ?`

	row := d.conn.QueryRow(query, tripID)
	var trip UserTrip
	err := row.Scan(&trip.ID, &trip.UserID, &trip.Title, &trip.DestinationID, &trip.Budget, &trip.Duration, &trip.StartDate, &trip.Status, &trip.CreatedAt, &trip.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user trip not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user trip: %w", err)
	}

	// Fetch segments for this trip
	segments, err := d.GetTripSegments(tripID)
	if err == nil {
		trip.Segments = segments
	}

	return &trip, nil
}

// GetUserTripsByUserID retrieves all trips for a user
func (d *Database) GetUserTripsByUserID(userID string) ([]UserTrip, error) {
	query := `SELECT id, user_id, title, destination_id, budget, duration, start_date, status, created_at, updated_at
			 FROM user_trips WHERE user_id = ? ORDER BY start_date DESC`

	rows, err := d.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user trips: %w", err)
	}
	defer rows.Close()

	var trips []UserTrip
	for rows.Next() {
		var trip UserTrip
		if err := rows.Scan(&trip.ID, &trip.UserID, &trip.Title, &trip.DestinationID, &trip.Budget, &trip.Duration, &trip.StartDate, &trip.Status, &trip.CreatedAt, &trip.UpdatedAt); err != nil {
			continue
		}

		// Fetch segments for each trip
		segments, err := d.GetTripSegments(trip.ID)
		if err == nil {
			trip.Segments = segments
		}

		trips = append(trips, trip)
	}

	return trips, rows.Err()
}

// UpdateUserTrip updates a user trip
func (d *Database) UpdateUserTrip(trip *UserTrip) error {
	query := `UPDATE user_trips SET title = ?, destination_id = ?, budget = ?, duration = ?, start_date = ?, status = ?, updated_at = ?
			 WHERE id = ? AND user_id = ?`

	result, err := d.conn.Exec(query, trip.Title, trip.DestinationID, trip.Budget, trip.Duration, trip.StartDate, trip.Status, trip.UpdatedAt, trip.ID, trip.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user trip: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user trip not found or permission denied")
	}

	return nil
}

// UpdateUserTripStatus updates the status of a user trip
func (d *Database) UpdateUserTripStatus(tripID, status string) error {
	query := `UPDATE user_trips SET status = ?, updated_at = ? WHERE id = ?`
	_, err := d.conn.Exec(query, status, time.Now(), tripID)
	if err != nil {
		return fmt.Errorf("failed to update trip status: %w", err)
	}
	return nil
}

// DeleteUserTrip deletes a user trip and cascading data
func (d *Database) DeleteUserTrip(tripID string) error {
	query := `DELETE FROM user_trips WHERE id = ?`
	_, err := d.conn.Exec(query, tripID)
	if err != nil {
		return fmt.Errorf("failed to delete user trip: %w", err)
	}
	return nil
}

// AddTripSegment adds a segment to a trip
func (d *Database) AddTripSegment(segment *TripSegment) error {
	query := `INSERT INTO trip_segments (id, user_trip_id, day, name, type, location, latitude, longitude, notes, completed, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query, segment.ID, segment.UserTripID, segment.Day, segment.Name, segment.Type, segment.Location, segment.Latitude, segment.Longitude, segment.Notes, segment.Completed, segment.CreatedAt, segment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to add trip segment: %w", err)
	}
	return nil
}

// GetTripSegments retrieves all segments for a trip
func (d *Database) GetTripSegments(tripID string) ([]TripSegment, error) {
	query := `SELECT id, user_trip_id, day, name, type, location, latitude, longitude, notes, completed, created_at, updated_at
			 FROM trip_segments WHERE user_trip_id = ? ORDER BY day ASC`

	rows, err := d.conn.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to query trip segments: %w", err)
	}
	defer rows.Close()

	var segments []TripSegment
	for rows.Next() {
		var segment TripSegment
		if err := rows.Scan(&segment.ID, &segment.UserTripID, &segment.Day, &segment.Name, &segment.Type, &segment.Location, &segment.Latitude, &segment.Longitude, &segment.Notes, &segment.Completed, &segment.CreatedAt, &segment.UpdatedAt); err != nil {
			continue
		}

		// Fetch photos for this segment
		photos, err := d.GetTripPhotos(segment.ID)
		if err == nil {
			segment.Photos = photos
		}

		// Fetch review for this segment
		review, err := d.GetTripReview(segment.ID)
		if err == nil {
			segment.Review = review
		}

		segments = append(segments, segment)
	}

	return segments, rows.Err()
}

// AddTripPhoto adds a photo to a segment
func (d *Database) AddTripPhoto(photo *TripPhoto) error {
	query := `INSERT INTO trip_photos (id, trip_segment_id, url, caption, uploaded_at)
			 VALUES (?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query, photo.ID, photo.TripSegmentID, photo.URL, photo.Caption, photo.UploadedAt)
	if err != nil {
		return fmt.Errorf("failed to add trip photo: %w", err)
	}
	return nil
}

// GetTripPhotos retrieves all photos for a segment
func (d *Database) GetTripPhotos(segmentID string) ([]TripPhoto, error) {
	query := `SELECT id, trip_segment_id, url, caption, uploaded_at
			 FROM trip_photos WHERE trip_segment_id = ? ORDER BY uploaded_at ASC`

	rows, err := d.conn.Query(query, segmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query trip photos: %w", err)
	}
	defer rows.Close()

	var photos []TripPhoto
	for rows.Next() {
		var photo TripPhoto
		if err := rows.Scan(&photo.ID, &photo.TripSegmentID, &photo.URL, &photo.Caption, &photo.UploadedAt); err != nil {
			continue
		}
		photos = append(photos, photo)
	}

	return photos, rows.Err()
}

// AddTripReview adds or updates a review for a segment
func (d *Database) AddTripReview(review *TripReview) error {
	// Check if review already exists
	var existingID string
	err := d.conn.QueryRow("SELECT id FROM trip_reviews WHERE trip_segment_id = ?", review.TripSegmentID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new review
		insertQuery := `INSERT INTO trip_reviews (id, trip_segment_id, rating, review, created_at, updated_at)
					   VALUES (?, ?, ?, ?, ?, ?)`
		_, err = d.conn.Exec(insertQuery, review.ID, review.TripSegmentID, review.Rating, review.Review, review.CreatedAt, review.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to add trip review: %w", err)
		}
	} else if err == nil {
		// Update existing review
		updateQuery := `UPDATE trip_reviews SET rating = ?, review = ?, updated_at = ? WHERE trip_segment_id = ?`
		_, err = d.conn.Exec(updateQuery, review.Rating, review.Review, review.UpdatedAt, review.TripSegmentID)
		if err != nil {
			return fmt.Errorf("failed to update trip review: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check existing review: %w", err)
	}

	return nil
}

// GetTripReview retrieves a review for a segment
func (d *Database) GetTripReview(segmentID string) (*TripReview, error) {
	query := `SELECT id, trip_segment_id, rating, review, created_at, updated_at
			 FROM trip_reviews WHERE trip_segment_id = ?`

	row := d.conn.QueryRow(query, segmentID)
	var review TripReview
	err := row.Scan(&review.ID, &review.TripSegmentID, &review.Rating, &review.Review, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get trip review: %w", err)
	}

	return &review, nil
}

// PublishUserTrip publishes a user trip as a community post
func (d *Database) PublishUserTrip(post *UserTripPost) error {
	query := `INSERT INTO user_trip_posts (id, user_trip_id, user_id, title, description, cover_image, likes, views, published, published_at, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query, post.ID, post.UserTripID, post.UserID, post.Title, post.Description, post.CoverImage, post.Likes, post.Views, post.Published, post.PublishedAt, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to publish user trip: %w", err)
	}
	return nil
}

// GetCommunityPosts retrieves published community posts with pagination
func (d *Database) GetCommunityPosts(page, pageSize int) ([]UserTripPost, error) {
	offset := (page - 1) * pageSize

	query := `SELECT id, user_trip_id, user_id, title, description, cover_image, likes, views, published, published_at, created_at, updated_at
			 FROM user_trip_posts WHERE published = 1 ORDER BY published_at DESC LIMIT ? OFFSET ?`

	rows, err := d.conn.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query community posts: %w", err)
	}
	defer rows.Close()

	var posts []UserTripPost
	for rows.Next() {
		var post UserTripPost
		if err := rows.Scan(&post.ID, &post.UserTripID, &post.UserID, &post.Title, &post.Description, &post.CoverImage, &post.Likes, &post.Views, &post.Published, &post.PublishedAt, &post.CreatedAt, &post.UpdatedAt); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(userID string) (*User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?`

	row := d.conn.QueryRow(query, userID)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetDestinationByID retrieves a destination by ID
func (d *Database) GetDestinationByID(destinationID string) (*Destination, error) {
	query := `SELECT id, name, country, description, image_url, created_at, updated_at FROM destinations WHERE id = ?`

	row := d.conn.QueryRow(query, destinationID)
	var dest Destination
	var imageURL *string
	err := row.Scan(&dest.ID, &dest.Name, &dest.Country, &dest.Description, &imageURL, &dest.CreatedAt, &dest.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("destination not found")
		}
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	if imageURL != nil {
		dest.Image = *imageURL
	}

	return &dest, nil
}
