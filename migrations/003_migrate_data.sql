-- Phase B Week 1 - Day 3 Database Migration
-- SQLite to PostgreSQL Direct Transfer

-- Clear existing data
TRUNCATE TABLE users, destinations, itineraries, itinerary_items, comments, likes,
                user_plans, user_trips, trip_segments, trip_photos, trip_reviews,
                user_trip_posts CASCADE;

-- ============================================================
-- USERS MIGRATION
-- ============================================================
INSERT INTO users (id, username, email, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'traveler1', 'traveler1@example.com', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'explorer2', 'explorer@example.com', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655440003', 'wanderer3', 'wanderer@example.com', '2026-03-24 08:07:36', NOW());

-- ============================================================
-- DESTINATIONS MIGRATION
-- ============================================================
INSERT INTO destinations (id, name, country, description, image_url, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655450001', 'Goa', 'India', 
   'Beautiful coastal state with beaches, churches and Portuguese heritage. Known for nightlife, water sports and seafood restaurants.', 
   '', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655450002', 'Paris', 'France',
   'The City of Light is known for the Eiffel Tower, museums, gardens, and romantic atmosphere. Famous for art, gastronomy, and fashion.',
   '', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655450003', 'Tokyo', 'Japan',
   'A vibrant metropolis blending ancient traditions with cutting-edge technology. Known for temples, gardens, shopping, and amazing cuisine.',
   '', '2026-03-24 08:07:36', NOW());

-- ============================================================
-- ITINERARIES MIGRATION
-- ============================================================
INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655460001', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655450001',
   'Goa Beach Escape', 'A relaxing 5-day beach vacation in beautiful Goa', 5, 50000, '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655460002', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655450002',
   'Paris 7 Day Tour', 'Experience the romance and culture of Paris over 7 days', 7, 120000, '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655460003', '550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655450003',
   'Tokyo Adventures', 'Explore the wonders of Tokyo with modern tech and ancient temples', 10, 200000, '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655460004', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655450002',
   'Paris Romantic Getaway', 'A 4-day romantic weekend in the heart of Paris', 4, 90000, '2026-03-24 08:07:36', NOW());

-- ============================================================
-- ITINERARY ITEMS MIGRATION
-- Map to new UUIDs and update day numbers format
-- ============================================================
INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, location, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655470001', '550e8400-e29b-41d4-a716-446655460001', 1, 'activity', 'Arrive in Goa', 'Arrive at Dabolim airport and check into resort', 'Dabolim', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470002', '550e8400-e29b-41d4-a716-446655460001', 2, 'activity', 'Beach Day 1', 'Relax at Baga Beach, enjoy local seafood', 'Baga Beach', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470003', '550e8400-e29b-41d4-a716-446655460001', 2, 'activity', 'Water Sports', 'Try parasailing and jet ski at Calangute', 'Calangute', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470004', '550e8400-e29b-41d4-a716-446655460001', 3, 'activity', 'Fort Exploration', 'Visit historic Aguada Fort', 'Aguada', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470005', '550e8400-e29b-41d4-a716-446655460002', 1, 'activity', 'Eiffel Tower', 'Visit iconic Eiffel Tower', 'Paris', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470006', '550e8400-e29b-41d4-a716-446655460002', 2, 'activity', 'Louvre Museum', 'Explore world famous Louvre Museum', 'Paris', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470007', '550e8400-e29b-41d4-a716-446655460003', 1, 'activity', 'Tokyo arrival', 'Land and check into hotel', 'Tokyo', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470008', '550e8400-e29b-41d4-a716-446655460003', 2, 'activity', 'Senso-ji Temple', 'Visit ancient Buddhist temple', 'Asakusa', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470009', '550e8400-e29b-41d4-a716-446655460003', 3, 'activity', 'Shibuya Crossing', 'Experience worlds busiest crossing', 'Shibuya', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655470010', '550e8400-e29b-41d4-a716-446655460004', 1, 'activity', 'Arrive Paris', 'Arrive and settle into hotel', 'Paris', '2026-03-24 08:07:36', NOW());

-- ============================================================
-- COMMENTS MIGRATION
-- ============================================================
INSERT INTO comments (id, user_id, itinerary_id, content, created_at, updated_at) VALUES
  ('550e8400-e29b-41d4-a716-446655480001', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655460001',
   'This Goa itinerary looks amazing! Love the beach activities.', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655480002', '550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655460002',
   'Paris is my dream destination. This 7-day plan is perfect!', '2026-03-24 08:07:36', NOW()),
  ('550e8400-e29b-41d4-a716-446655480003', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655460003',
   'Tokyo looks incredible! Can''t wait to do this trip.', '2026-03-24 08:07:36', NOW());

-- ============================================================
-- VERIFY MIGRATION
-- ============================================================
SELECT COUNT(*) as users_count FROM users;
SELECT COUNT(*) as destinations_count FROM destinations;
SELECT COUNT(*) as itineraries_count FROM itineraries;
SELECT COUNT(*) as items_count FROM itinerary_items;
SELECT COUNT(*) as comments_count FROM comments;

-- Check foreign key relationships
SELECT 'Users' as table_name, COUNT(*) FROM users
UNION ALL
SELECT 'Destinations', COUNT(*) FROM destinations
UNION ALL
SELECT 'Itineraries', COUNT(*) FROM itineraries
UNION ALL
SELECT 'Itinerary Items', COUNT(*) FROM itinerary_items
UNION ALL
SELECT 'Comments', COUNT(*) FROM comments;
