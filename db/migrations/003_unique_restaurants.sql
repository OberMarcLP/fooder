-- Add unique constraint for Google Place ID to prevent duplicate restaurants from Google Maps
-- Allow NULL values since manually added restaurants won't have a google_place_id
CREATE UNIQUE INDEX idx_restaurants_google_place_id ON restaurants(google_place_id) WHERE google_place_id IS NOT NULL;

-- Add unique constraint for name and address combination to prevent duplicate manually added restaurants
-- This helps prevent the same restaurant being added twice with the same name and address
CREATE UNIQUE INDEX idx_restaurants_name_address ON restaurants(LOWER(name), LOWER(address)) WHERE address IS NOT NULL;
