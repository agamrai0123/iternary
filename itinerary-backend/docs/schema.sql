-- Itinerary Platform Oracle Database Schema
-- For use with Oracle 12c+

-- ==================== DROP EXISTING TABLES ====================
-- Uncomment if you need to drop existing tables
/*
BEGIN
  FOR rec IN (SELECT table_name FROM user_tables WHERE table_name IN ('USERS', 'DESTINATIONS', 'ITINERARIES', 'ITINERARY_ITEMS', 'COMMENTS', 'USER_PLANS', 'LIKES'))
  LOOP
    EXECUTE IMMEDIATE 'DROP TABLE ' || rec.table_name || ' CASCADE CONSTRAINTS';
  END LOOP;
END;
/
*/

-- ==================== CREATE TABLES ====================

-- Create Users Table
CREATE TABLE users (
  id VARCHAR2(36) PRIMARY KEY,
  username VARCHAR2(50) UNIQUE NOT NULL,
  email VARCHAR2(100) UNIQUE NOT NULL,
  password_hash VARCHAR2(255),
  created_at TIMESTAMP DEFAULT SYSDATE,
  updated_at TIMESTAMP DEFAULT SYSDATE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- Create Destinations Table
CREATE TABLE destinations (
  id VARCHAR2(36) PRIMARY KEY,
  name VARCHAR2(100) NOT NULL,
  country VARCHAR2(100) NOT NULL,
  description CLOB,
  image_url VARCHAR2(500),
  created_at TIMESTAMP DEFAULT SYSDATE,
  updated_at TIMESTAMP DEFAULT SYSDATE
);

CREATE INDEX idx_destinations_country ON destinations(country);

-- Create Itineraries Table
CREATE TABLE itineraries (
  id VARCHAR2(36) PRIMARY KEY,
  user_id VARCHAR2(36) NOT NULL,
  destination_id VARCHAR2(36) NOT NULL,
  title VARCHAR2(200) NOT NULL,
  description CLOB,
  duration NUMBER NOT NULL,
  budget NUMBER(10,2) NOT NULL,
  likes NUMBER DEFAULT 0,
  created_at TIMESTAMP DEFAULT SYSDATE,
  updated_at TIMESTAMP DEFAULT SYSDATE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);

CREATE INDEX idx_itineraries_user_id ON itineraries(user_id);
CREATE INDEX idx_itineraries_destination_id ON itineraries(destination_id);
CREATE INDEX idx_itineraries_likes ON itineraries(likes DESC);
CREATE INDEX idx_itineraries_created_at ON itineraries(created_at DESC);

-- Create Itinerary Items Table
CREATE TABLE itinerary_items (
  id VARCHAR2(36) PRIMARY KEY,
  itinerary_id VARCHAR2(36) NOT NULL,
  day NUMBER NOT NULL,
  type VARCHAR2(20) NOT NULL,
  name VARCHAR2(200) NOT NULL,
  description CLOB,
  price NUMBER(10,2) NOT NULL,
  duration NUMBER,
  location VARCHAR2(200),
  rating NUMBER(3,2),
  image_url VARCHAR2(500),
  booking_url VARCHAR2(500),
  created_at TIMESTAMP DEFAULT SYSDATE,
  updated_at TIMESTAMP DEFAULT SYSDATE,
  FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
  CHECK (type IN ('stay', 'food', 'activity', 'transport', 'other'))
);

CREATE INDEX idx_itinerary_items_day ON itinerary_items(itinerary_id, day);
CREATE INDEX idx_itinerary_items_type ON itinerary_items(type);

-- Create Comments Table
CREATE TABLE comments (
  id VARCHAR2(36) PRIMARY KEY,
  itinerary_id VARCHAR2(36) NOT NULL,
  user_id VARCHAR2(36) NOT NULL,
  content CLOB NOT NULL,
  rating NUMBER(3,2),
  created_at TIMESTAMP DEFAULT SYSDATE,
  updated_at TIMESTAMP DEFAULT SYSDATE,
  FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_itinerary_id ON comments(itinerary_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

-- Create User Plans Table (Saved/copied itineraries)
CREATE TABLE user_plans (
  id VARCHAR2(36) PRIMARY KEY,
  user_id VARCHAR2(36) NOT NULL,
  original_itinerary_id VARCHAR2(36) NOT NULL,
  title VARCHAR2(200),
  notes CLOB,
  status VARCHAR2(20) DEFAULT 'draft',
  created_at TIMESTAMP DEFAULT SYSDATE,
  updated_at TIMESTAMP DEFAULT SYSDATE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (original_itinerary_id) REFERENCES itineraries(id),
  CHECK (status IN ('draft', 'planned', 'completed'))
);

CREATE INDEX idx_user_plans_user_id ON user_plans(user_id);
CREATE INDEX idx_user_plans_status ON user_plans(status);

-- Create Likes Table
CREATE TABLE likes (
  id VARCHAR2(36) PRIMARY KEY,
  user_id VARCHAR2(36) NOT NULL,
  itinerary_id VARCHAR2(36) NOT NULL,
  created_at TIMESTAMP DEFAULT SYSDATE,
  UNIQUE (user_id, itinerary_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE
);

CREATE INDEX idx_likes_itinerary_id ON likes(itinerary_id);

-- ==================== INSERT TEST DATA ====================

-- Insert test users
INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES ('user-001', 'traveler1', 'traveler1@example.com', 'hash_placeholder_1', SYSDATE, SYSDATE);

INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES ('user-002', 'explorer2', 'explorer@example.com', 'hash_placeholder_2', SYSDATE, SYSDATE);

INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES ('user-003', 'wanderer3', 'wanderer@example.com', 'hash_placeholder_3', SYSDATE, SYSDATE);

-- Insert test destinations
INSERT INTO destinations (id, name, country, description, created_at, updated_at)
VALUES ('dest-001', 'Goa', 'India', 'Beautiful coastal state with beaches, churches and Portuguese heritage. Known for nightlife, water sports and seafood restaurants.', SYSDATE, SYSDATE);

INSERT INTO destinations (id, name, country, description, created_at, updated_at)
VALUES ('dest-002', 'Manali', 'India', 'Hill station in Himachal Pradesh. Perfect for adventure sports, trekking, and scenic mountain views. Great for families and adventure seekers.', SYSDATE, SYSDATE);

INSERT INTO destinations (id, name, country, description, created_at, updated_at)
VALUES ('dest-003', 'Bali', 'Indonesia', 'Tropical island known for beaches, rice terraces, temples and vibrant culture. Budget-friendly and great for all types of travelers.', SYSDATE, SYSDATE);

-- Insert test itineraries
INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at)
VALUES ('itin-001', 'user-001', 'dest-001', '5-Day Budget Goa Trip', 'Perfect 5 days in Goa with beaches, culture, and food. Includes Baga Beach, Old Goa churches, and local cuisine. Includes 3-star accommodation.', 5, 15000, 45, SYSDATE, SYSDATE);

INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at)
VALUES ('itin-002', 'user-002', 'dest-001', 'Luxury 7-Day Goa Experience', 'High-end Goa itinerary with 5-star resorts, adventure activities and fine dining. Includes water sports, spa and sunset cruises.', 7, 45000, 32, SYSDATE, SYSDATE);

INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at)
VALUES ('itin-003', 'user-003', 'dest-002', '4-Day Manali Adventure', 'Adventure-focused trip with trekking, paragliding, and camping. Includes visits to Rohtang Pass and local villages.', 4, 12000, 28, SYSDATE, SYSDATE);

INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at)
VALUES ('itin-004', 'user-001', 'dest-003', '6-Day Bali Paradise', 'Comprehensive Bali experience with beaches, temples, rice paddies and nightlife. Perfect for first-time visitors.', 6, 18000, 67, SYSDATE, SYSDATE);

-- Insert itinerary items for itin-001
INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-001', 'itin-001', 1, 'stay', 'Beachside Hostel', '3-star hostel with AC rooms, common kitchen, and beach access', 500, 24, 'Baga', 4.5, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-002', 'itin-001', 1, 'food', 'Seafood Dinner Baga', 'Fresh fish curry, prawn biryani and local wines at beachside restaurant', 350, 2, 'Baga Beach', 4.7, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-003', 'itin-001', 2, 'activity', 'Jet Ski Adventure', 'Water sports activity including jet ski, banana boat ride', 800, 2, 'Baga Beach', 4.8, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-004', 'itin-001', 2, 'food', 'Beachside Lunch', 'Grilled fish with rice and salad at beach shack', 250, 1.5, 'Anjuna Beach', 4.6, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-005', 'itin-001', 3, 'activity', 'Old Goa Heritage Tour', 'Visit Se Cathedral, Basilica of Bom Jesus and other historic churches', 300, 3, 'Old Goa', 4.4, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-006', 'itin-001', 4, 'transport', 'Scooter Rental', 'Rent a scooter to explore Goa independently', 400, 48, 'Baga', 4.3, SYSDATE, SYSDATE);

-- Insert itinerary items for itin-004
INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-101', 'itin-004', 1, 'stay', 'Luxury Resort Seminyak', 'Beachfront 5-star resort with private pool and ocean views', 3000, 24, 'Seminyak', 4.9, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-102', 'itin-004', 1, 'food', 'Fine Dining Dinner', 'Indonesian and international cuisine at beach club restaurant', 800, 2, 'Seminyak Beach', 4.8, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-103', 'itin-004', 2, 'activity', 'Ubud Rice Paddies Tour', 'Explore scenic rice terraces and local villages', 400, 4, 'Ubud', 4.7, SYSDATE, SYSDATE);

INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, rating, created_at, updated_at)
VALUES ('item-104', 'itin-004', 3, 'activity', 'Temple Tour', 'Visit Tanah Lot and Uluwatu Temple at sunset', 350, 5, 'Various', 4.6, SYSDATE, SYSDATE);

-- Insert comments
INSERT INTO comments (id, itinerary_id, user_id, content, rating, created_at, updated_at)
VALUES ('comment-001', 'itin-001', 'user-002', 'Amazing itinerary! Did exactly this and had the best time. Jet ski was super fun and the food recommendations were spot on.', 5, SYSDATE, SYSDATE);

INSERT INTO comments (id, itinerary_id, user_id, content, rating, created_at, updated_at)
VALUES ('comment-002', 'itin-001', 'user-003', 'Great value for money. Would definitely recommend this plan to budget travelers. The beachside atmosphere is incredible.', 4.5, SYSDATE, SYSDATE);

INSERT INTO comments (id, itinerary_id, user_id, content, rating, created_at, updated_at)
VALUES ('comment-003', 'itin-004', 'user-001', 'Luxury experience but worth every penny. Bali is truly paradise. The temples are magnificent and people are super friendly.', 5, SYSDATE, SYSDATE);

-- ==================== COMMIT CHANGES ====================
COMMIT;

-- ==================== VERIFICATION QUERIES ====================
-- Run these to verify data was inserted correctly:
/*
SELECT COUNT(*) as user_count FROM users;
SELECT COUNT(*) as destination_count FROM destinations;
SELECT COUNT(*) as itinerary_count FROM itineraries;
SELECT COUNT(*) as item_count FROM itinerary_items;
SELECT COUNT(*) as comment_count FROM comments;

-- View sample data
SELECT * FROM destinations;
SELECT * FROM itineraries;
SELECT * FROM itinerary_items WHERE itinerary_id = 'itin-001';
*/

