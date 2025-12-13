-- Add unique constraint for Google Place ID on suggestions to prevent duplicate suggestion submissions
-- Allow NULL values since manually added suggestions won't have a google_place_id
CREATE UNIQUE INDEX idx_suggestions_google_place_id ON restaurant_suggestions(google_place_id) WHERE google_place_id IS NOT NULL;

-- Add unique constraint for name and address combination on suggestions
-- This helps prevent the same suggestion being submitted multiple times
CREATE UNIQUE INDEX idx_suggestions_name_address ON restaurant_suggestions(LOWER(name), LOWER(address)) WHERE address IS NOT NULL;
