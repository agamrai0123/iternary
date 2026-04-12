-- Itinerary Application - PostgreSQL Schema Migration
-- Date: 2026-04-02
-- Version: 1.0

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Destinations table
CREATE TABLE IF NOT EXISTS destinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(512),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Itineraries table
CREATE TABLE IF NOT EXISTS itineraries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    destination_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL,
    budget NUMERIC(10, 2) NOT NULL,
    likes INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (destination_id) REFERENCES destinations(id)
);

-- Itinerary items table
CREATE TABLE IF NOT EXISTS itinerary_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_id UUID NOT NULL,
    day INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    duration INTEGER,
    location VARCHAR(255),
    rating NUMERIC(3, 1),
    image_url VARCHAR(512),
    booking_url VARCHAR(512),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
    CONSTRAINT valid_item_type CHECK (type IN ('stay', 'food', 'activity', 'transport', 'other'))
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    itinerary_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    rating NUMERIC(3, 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- User plans table
CREATE TABLE IF NOT EXISTS user_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    original_itinerary_id UUID NOT NULL,
    title VARCHAR(255),
    notes TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (original_itinerary_id) REFERENCES itineraries(id),
    CONSTRAINT valid_plan_status CHECK (status IN ('draft', 'planned', 'completed'))
);

-- Likes table
CREATE TABLE IF NOT EXISTS likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    itinerary_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, itinerary_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE
);

-- User Trips table (custom trip planning)
CREATE TABLE IF NOT EXISTS user_trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    destination_id UUID NOT NULL,
    budget NUMERIC(10, 2) NOT NULL,
    duration INTEGER NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (destination_id) REFERENCES destinations(id),
    CONSTRAINT valid_trip_status CHECK (status IN ('draft', 'planning', 'ongoing', 'completed'))
);

-- Trip Segments table (places/activities in a trip)
CREATE TABLE IF NOT EXISTS trip_segments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_trip_id UUID NOT NULL,
    day INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    location VARCHAR(255),
    latitude NUMERIC(10, 8),
    longitude NUMERIC(11, 8),
    notes TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_trip_id) REFERENCES user_trips(id) ON DELETE CASCADE
);

-- Trip Photos table (photos for trip segments)
CREATE TABLE IF NOT EXISTS trip_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_segment_id UUID NOT NULL,
    url VARCHAR(512) NOT NULL,
    caption TEXT,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (trip_segment_id) REFERENCES trip_segments(id) ON DELETE CASCADE
);

-- Trip Reviews table (reviews for completed segments)
CREATE TABLE IF NOT EXISTS trip_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_segment_id UUID NOT NULL,
    rating NUMERIC(3, 1) NOT NULL,
    review TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (trip_segment_id) REFERENCES trip_segments(id) ON DELETE CASCADE,
    CONSTRAINT valid_rating CHECK (rating >= 1 AND rating <= 5)
);

-- User Trip Posts table (published community posts)
CREATE TABLE IF NOT EXISTS user_trip_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_trip_id UUID NOT NULL,
    user_id UUID NOT NULL,
    title VARCHAR(255),
    description TEXT,
    cover_image VARCHAR(512),
    likes INTEGER DEFAULT 0,
    views INTEGER DEFAULT 0,
    published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_trip_id) REFERENCES user_trips(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_destinations_country ON destinations(country);
CREATE INDEX idx_destinations_name ON destinations(name);

CREATE INDEX idx_itineraries_user_id ON itineraries(user_id);
CREATE INDEX idx_itineraries_destination_id ON itineraries(destination_id);
CREATE INDEX idx_itineraries_likes ON itineraries(likes DESC);
CREATE INDEX idx_itineraries_created_at ON itineraries(created_at DESC);

CREATE INDEX idx_itinerary_items_itinerary_id ON itinerary_items(itinerary_id);
CREATE INDEX idx_itinerary_items_day ON itinerary_items(itinerary_id, day);
CREATE INDEX idx_itinerary_items_type ON itinerary_items(type);

CREATE INDEX idx_comments_itinerary_id ON comments(itinerary_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

CREATE INDEX idx_user_plans_user_id ON user_plans(user_id);
CREATE INDEX idx_user_plans_status ON user_plans(status);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_itinerary_id ON likes(itinerary_id);

CREATE INDEX idx_user_trips_user_id ON user_trips(user_id);
CREATE INDEX idx_user_trips_status ON user_trips(status);
CREATE INDEX idx_user_trips_destination_id ON user_trips(destination_id);

CREATE INDEX idx_trip_segments_user_trip_id ON trip_segments(user_trip_id);
CREATE INDEX idx_trip_segments_day ON trip_segments(user_trip_id, day);

CREATE INDEX idx_trip_photos_segment_id ON trip_photos(trip_segment_id);
CREATE INDEX idx_trip_reviews_segment_id ON trip_reviews(trip_segment_id);

CREATE INDEX idx_user_trip_posts_user_id ON user_trip_posts(user_id);
CREATE INDEX idx_user_trip_posts_published ON user_trip_posts(published);

-- Create migration tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    description VARCHAR(255),
    executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO schema_migrations (version, description) 
VALUES (1, 'Initial schema creation') 
ON CONFLICT DO NOTHING;
