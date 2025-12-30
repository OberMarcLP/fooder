-- Additional performance indexes for optimized queries

-- Index for searching restaurants by name (used in GlobalSearch)
CREATE INDEX IF NOT EXISTS idx_restaurants_name_lower ON restaurants(LOWER(name));

-- Composite index for geospatial queries (used in location-based filtering)
CREATE INDEX IF NOT EXISTS idx_restaurants_location ON restaurants(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Index for suggestion searches by name
CREATE INDEX IF NOT EXISTS idx_suggestions_name_lower ON restaurant_suggestions(LOWER(name));

-- Composite index for suggestion geospatial queries
CREATE INDEX IF NOT EXISTS idx_suggestions_location ON restaurant_suggestions(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Composite index for suggestion category filtering
CREATE INDEX IF NOT EXISTS idx_suggestions_category_status ON restaurant_suggestions(suggested_category_id, status);

-- Index for created_at ordering (frequently used in ORDER BY clauses)
CREATE INDEX IF NOT EXISTS idx_restaurants_created_at ON restaurants(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_suggestions_created_at ON restaurant_suggestions(created_at DESC);

-- Index for ratings aggregation queries
CREATE INDEX IF NOT EXISTS idx_ratings_restaurant_ratings ON ratings(restaurant_id, food_rating, service_rating, ambiance_rating);

-- Index for menu photos lookup by restaurant
CREATE INDEX IF NOT EXISTS idx_menu_photos_restaurant_created ON menu_photos(restaurant_id, created_at DESC);
