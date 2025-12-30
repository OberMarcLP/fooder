-- Categories (cultural categories like Italian, Asian, etc.)
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Food types (pizza, pasta, sushi, etc.)
CREATE TABLE food_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Restaurants
CREATE TABLE restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address VARCHAR(500),
    phone VARCHAR(50),
    website VARCHAR(500),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    google_place_id VARCHAR(255),
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Junction table for many-to-many relationship between restaurants and food types
CREATE TABLE restaurant_food_types (
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    food_type_id INTEGER NOT NULL REFERENCES food_types(id) ON DELETE CASCADE,
    PRIMARY KEY (restaurant_id, food_type_id)
);

-- Ratings
CREATE TABLE ratings (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    food_rating INTEGER NOT NULL CHECK (food_rating >= 1 AND food_rating <= 5),
    service_rating INTEGER NOT NULL CHECK (service_rating >= 1 AND service_rating <= 5),
    ambiance_rating INTEGER NOT NULL CHECK (ambiance_rating >= 1 AND ambiance_rating <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX idx_restaurants_category ON restaurants(category_id);
CREATE INDEX idx_restaurant_food_types_restaurant ON restaurant_food_types(restaurant_id);
CREATE INDEX idx_restaurant_food_types_food_type ON restaurant_food_types(food_type_id);
CREATE INDEX idx_ratings_restaurant ON ratings(restaurant_id);

-- Insert some default categories
INSERT INTO categories (name) VALUES
    ('Italian'),
    ('Asian'),
    ('Mexican'),
    ('American'),
    ('French'),
    ('Indian'),
    ('Mediterranean'),
    ('Japanese'),
    ('Chinese'),
    ('Thai');

-- Insert some default food types
INSERT INTO food_types (name) VALUES
    ('Pizza'),
    ('Pasta'),
    ('Sushi'),
    ('Burgers'),
    ('Tacos'),
    ('Curry'),
    ('Steak'),
    ('Seafood'),
    ('Salads'),
    ('Desserts');
