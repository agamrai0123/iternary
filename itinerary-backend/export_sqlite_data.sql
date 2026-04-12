.mode csv
.headers on

.output users.csv
SELECT id, username, email, created_at FROM users;

.output destinations.csv
SELECT id, name, country, description, image_url FROM destinations;

.output itineraries.csv
SELECT id, user_id, destination_id, title, description, duration, budget, created_at FROM itineraries;

.output itinerary_items.csv
SELECT id, itinerary_id, day_number, title, description, location, start_time, end_time, created_at FROM itinerary_items;

.output comments.csv
SELECT id, user_id, itinerary_id, comment_text, created_at FROM comments;

.output likes.csv
SELECT id, user_id, itinerary_id, created_at FROM likes;

.output user_plans.csv
SELECT id, user_id, plan_name, plan_data, created_at FROM user_plans;

.output user_trips.csv
SELECT id, user_id, trip_name, trip_description, start_date, end_date, created_at FROM user_trips;

.output trip_segments.csv
SELECT id, user_trip_id, segment_number, location, description, start_time, end_time, created_at FROM trip_segments;

.output trip_photos.csv
SELECT id, user_trip_id, photo_url, description, created_at FROM trip_photos;

.output trip_reviews.csv
SELECT id, user_trip_id, rating, review_text, created_at FROM trip_reviews;

.output user_trip_posts.csv
SELECT id, user_trip_id, user_id, post_text, created_at FROM user_trip_posts;

.output
quit
