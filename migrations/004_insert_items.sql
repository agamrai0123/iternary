INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, location, price, created_at, updated_at) 
VALUES 
  ('550e8400-e29b-41d4-a716-446655470001', '550e8400-e29b-41d4-a716-446655460001', 1, 'activity', 'Arrive in Goa', 'Arrive at airport', 'Dabolim', 0.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470002', '550e8400-e29b-41d4-a716-446655460001', 2, 'activity', 'Beach Day 1', 'Relax at beach', 'Baga Beach', 0.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470003', '550e8400-e29b-41d4-a716-446655460001', 2, 'activity', 'Water Sports', 'Parasailing and jet ski', 'Calangute', 2500.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470004', '550e8400-e29b-41d4-a716-446655460001', 3, 'activity', 'Fort Exploration', 'Visit historic fort', 'Aguada', 0.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470005', '550e8400-e29b-41d4-a716-446655460002', 1, 'activity', 'Eiffel Tower', 'Visit Eiffel Tower', 'Paris', 25.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470006', '550e8400-e29b-41d4-a716-446655460002', 2, 'activity', 'Louvre Museum', 'Visit Louvre', 'Paris', 20.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470007', '550e8400-e29b-41d4-a716-446655460003', 1, 'activity', 'Tokyo arrival', 'Land and check in', 'Tokyo', 0.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470008', '550e8400-e29b-41d4-a716-446655460003', 2, 'activity', 'Senso-ji Temple', 'Visit temple', 'Asakusa', 0.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470009', '550e8400-e29b-41d4-a716-446655460003', 3, 'activity', 'Shibuya Crossing', 'Experience crossing', 'Shibuya', 0.00, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655470010', '550e8400-e29b-41d4-a716-446655460004', 1, 'activity', 'Arrive Paris', 'Check into hotel', 'Paris', 0.00, NOW(), NOW());

SELECT COUNT(*) FROM itinerary_items;
