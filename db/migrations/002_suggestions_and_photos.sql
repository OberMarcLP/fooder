-- Restaurant Suggestions System
CREATE TABLE restaurant_suggestions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500),
    phone VARCHAR(50),
    website VARCHAR(500),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    google_place_id VARCHAR(255),
    suggested_category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'tested', 'rejected')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_suggestions_status ON restaurant_suggestions(status);

-- Suggestion Food Types Junction Table
CREATE TABLE suggestion_food_types (
    suggestion_id INTEGER NOT NULL REFERENCES restaurant_suggestions(id) ON DELETE CASCADE,
    food_type_id INTEGER NOT NULL REFERENCES food_types(id) ON DELETE CASCADE,
    PRIMARY KEY (suggestion_id, food_type_id)
);

CREATE INDEX idx_suggestion_food_types_suggestion ON suggestion_food_types(suggestion_id);
CREATE INDEX idx_suggestion_food_types_food_type ON suggestion_food_types(food_type_id);

-- Menu Photos Table
CREATE TABLE menu_photos (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255),
    caption VARCHAR(255) NOT NULL,
    file_size INTEGER,
    mime_type VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_menu_photos_restaurant ON menu_photos(restaurant_id);
