-- Clear existing data
TRUNCATE table users CASCADE;
TRUNCATE table destinations CASCADE;
TRUNCATE table itineraries CASCADE;
TRUNCATE table itinerary_items CASCADE;
TRUNCATE table comments CASCADE;
TRUNCATE table likes CASCADE;
TRUNCATE table user_plans CASCADE;
TRUNCATE table user_trips CASCADE;
TRUNCATE table trip_segments CASCADE;
TRUNCATE table trip_photos CASCADE;
TRUNCATE table trip_reviews CASCADE;
TRUNCATE table user_trip_posts CASCADE;

-- Insert users with generated UUIDs
INSERT INTO users (id, username, email, created_at, updated_at)
SELECT gen_random_uuid()::text, username, email, COALESCE(created_at, NOW()), NOW()
FROM (
  SELECT 
    ROW_NUMBER() OVER () as rn,
    username,
    email,
    created_at
  FROM (VALUES 
    ('john_doe', 'john@example.com', NULL),
    ('jane_smith', 'jane@example.com', NULL),
    ('alice_johnson', 'alice@example.com', NULL)
  ) AS t(username, email, created_at)
) u;

-- Sample to verify migration worked
SELECT * FROM users;
