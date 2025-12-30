-- Drop tables in reverse order (respecting foreign key constraints)
DROP INDEX IF EXISTS idx_ratings_restaurant;
DROP INDEX IF EXISTS idx_restaurant_food_types_food_type;
DROP INDEX IF EXISTS idx_restaurant_food_types_restaurant;
DROP INDEX IF EXISTS idx_restaurants_category;

DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS restaurant_food_types;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS food_types;
DROP TABLE IF EXISTS categories;
