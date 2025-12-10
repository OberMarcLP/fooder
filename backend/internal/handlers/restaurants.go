package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fooder/backend/internal/database"
	"github.com/fooder/backend/internal/models"
	"github.com/gorilla/mux"
)

func getFoodTypesForRestaurant(ctx context.Context, restaurantID int) ([]models.FoodType, error) {
	rows, err := database.GetPool().Query(ctx,
		`SELECT ft.id, ft.name, ft.created_at, ft.updated_at
		FROM food_types ft
		JOIN restaurant_food_types rft ON ft.id = rft.food_type_id
		WHERE rft.restaurant_id = $1
		ORDER BY ft.name`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foodTypes []models.FoodType
	for rows.Next() {
		var ft models.FoodType
		if err := rows.Scan(&ft.ID, &ft.Name, &ft.CreatedAt, &ft.UpdatedAt); err != nil {
			return nil, err
		}
		foodTypes = append(foodTypes, ft)
	}
	return foodTypes, nil
}

func setFoodTypesForRestaurant(ctx context.Context, restaurantID int, foodTypeIDs []int) error {
	// Delete existing food types
	_, err := database.GetPool().Exec(ctx,
		"DELETE FROM restaurant_food_types WHERE restaurant_id = $1", restaurantID)
	if err != nil {
		return err
	}

	// Insert new food types
	for _, ftID := range foodTypeIDs {
		_, err := database.GetPool().Exec(ctx,
			"INSERT INTO restaurant_food_types (restaurant_id, food_type_id) VALUES ($1, $2)",
			restaurantID, ftID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse query parameters for filtering
	queryParams := r.URL.Query()
	categoryID := queryParams.Get("category_id")
	foodTypeIDs := queryParams.Get("food_type_ids") // comma-separated
	lat := queryParams.Get("lat")
	lng := queryParams.Get("lng")
	radius := queryParams.Get("radius") // in kilometers

	// Build dynamic query with filters
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Category filter
	if categoryID != "" {
		if catID, err := strconv.Atoi(categoryID); err == nil {
			conditions = append(conditions, fmt.Sprintf("r.category_id = $%d", argIndex))
			args = append(args, catID)
			argIndex++
		}
	}

	// Food types filter (restaurants that have ANY of the specified food types)
	if foodTypeIDs != "" {
		ftIDs := strings.Split(foodTypeIDs, ",")
		var validIDs []int
		for _, idStr := range ftIDs {
			if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
				validIDs = append(validIDs, id)
			}
		}
		if len(validIDs) > 0 {
			placeholders := make([]string, len(validIDs))
			for i, id := range validIDs {
				placeholders[i] = fmt.Sprintf("$%d", argIndex)
				args = append(args, id)
				argIndex++
			}
			conditions = append(conditions, fmt.Sprintf(`r.id IN (
				SELECT DISTINCT restaurant_id FROM restaurant_food_types
				WHERE food_type_id IN (%s)
			)`, strings.Join(placeholders, ",")))
		}
	}

	// Location/radius filter using Haversine formula
	var distanceSelect string
	var distanceOrder string
	if lat != "" && lng != "" && radius != "" {
		latVal, latErr := strconv.ParseFloat(lat, 64)
		lngVal, lngErr := strconv.ParseFloat(lng, 64)
		radiusVal, radErr := strconv.ParseFloat(radius, 64)

		if latErr == nil && lngErr == nil && radErr == nil {
			// Haversine formula for distance in km
			distanceSelect = fmt.Sprintf(`,
				(6371 * acos(
					cos(radians($%d)) * cos(radians(r.latitude)) *
					cos(radians(r.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(r.latitude))
				)) as distance`, argIndex, argIndex+1, argIndex+2)
			args = append(args, latVal, lngVal, latVal)

			conditions = append(conditions, fmt.Sprintf(`r.latitude IS NOT NULL AND r.longitude IS NOT NULL AND
				(6371 * acos(
					cos(radians($%d)) * cos(radians(r.latitude)) *
					cos(radians(r.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(r.latitude))
				)) <= $%d`, argIndex+3, argIndex+4, argIndex+5, argIndex+6))
			args = append(args, latVal, lngVal, latVal, radiusVal)
			argIndex += 7

			distanceOrder = "distance ASC,"
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			r.id, r.name, r.description, r.address, r.phone, r.website, r.latitude, r.longitude,
			r.google_place_id, r.category_id, r.created_at, r.updated_at,
			c.id, c.name,
			COALESCE(AVG(rt.food_rating), 0) as avg_food,
			COALESCE(AVG(rt.service_rating), 0) as avg_service,
			COALESCE(AVG(rt.ambiance_rating), 0) as avg_ambiance,
			COUNT(rt.id) as rating_count
			%s
		FROM restaurants r
		LEFT JOIN categories c ON r.category_id = c.id
		LEFT JOIN ratings rt ON r.id = rt.restaurant_id
		%s
		GROUP BY r.id, c.id
		ORDER BY %s r.created_at DESC
	`, distanceSelect, whereClause, distanceOrder)

	rows, err := database.GetPool().Query(ctx, query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	hasDistance := distanceSelect != ""
	restaurants := []models.Restaurant{}
	for rows.Next() {
		var rest models.Restaurant
		var catID *int
		var catName *string
		var avgFood, avgService, avgAmbiance float64
		var ratingCount int
		var distance *float64

		var err error
		if hasDistance {
			err = rows.Scan(
				&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
				&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
				&catID, &catName,
				&avgFood, &avgService, &avgAmbiance, &ratingCount,
				&distance,
			)
		} else {
			err = rows.Scan(
				&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
				&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
				&catID, &catName,
				&avgFood, &avgService, &avgAmbiance, &ratingCount,
			)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if distance != nil {
			rest.Distance = distance
		}

		if catID != nil && catName != nil {
			rest.Category = &models.Category{ID: *catID, Name: *catName}
		}

		// Get food types for this restaurant
		foodTypes, err := getFoodTypesForRestaurant(ctx, rest.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rest.FoodTypes = foodTypes

		if ratingCount > 0 {
			overall := (avgFood + avgService + avgAmbiance) / 3
			rest.AvgRating = &models.AvgRating{
				Food:     avgFood,
				Service:  avgService,
				Ambiance: avgAmbiance,
				Overall:  overall,
				Count:    ratingCount,
			}
		}

		restaurants = append(restaurants, rest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	query := `
		SELECT
			r.id, r.name, r.description, r.address, r.phone, r.website, r.latitude, r.longitude,
			r.google_place_id, r.category_id, r.created_at, r.updated_at,
			c.id, c.name,
			COALESCE(AVG(rt.food_rating), 0) as avg_food,
			COALESCE(AVG(rt.service_rating), 0) as avg_service,
			COALESCE(AVG(rt.ambiance_rating), 0) as avg_ambiance,
			COUNT(rt.id) as rating_count
		FROM restaurants r
		LEFT JOIN categories c ON r.category_id = c.id
		LEFT JOIN ratings rt ON r.id = rt.restaurant_id
		WHERE r.id = $1
		GROUP BY r.id, c.id
	`

	var rest models.Restaurant
	var catID *int
	var catName *string
	var avgFood, avgService, avgAmbiance float64
	var ratingCount int

	err = database.GetPool().QueryRow(ctx, query, id).Scan(
		&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
		&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
		&catID, &catName,
		&avgFood, &avgService, &avgAmbiance, &ratingCount,
	)
	if err != nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	if catID != nil && catName != nil {
		rest.Category = &models.Category{ID: *catID, Name: *catName}
	}

	// Get food types
	foodTypes, err := getFoodTypesForRestaurant(ctx, rest.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rest.FoodTypes = foodTypes

	if ratingCount > 0 {
		overall := (avgFood + avgService + avgAmbiance) / 3
		rest.AvgRating = &models.AvgRating{
			Food:     avgFood,
			Service:  avgService,
			Ambiance: avgAmbiance,
			Overall:  overall,
			Count:    ratingCount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rest)
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	var rest models.Restaurant
	err := database.GetPool().QueryRow(ctx,
		`INSERT INTO restaurants (name, description, address, phone, website, latitude, longitude, google_place_id, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, description, address, phone, website, latitude, longitude, google_place_id, category_id, created_at, updated_at`,
		req.Name, req.Description, req.Address, req.Phone, req.Website, req.Latitude, req.Longitude, req.GooglePlaceID, req.CategoryID,
	).Scan(
		&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
		&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set food types
	if len(req.FoodTypeIDs) > 0 {
		if err := setFoodTypesForRestaurant(ctx, rest.ID, req.FoodTypeIDs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		foodTypes, _ := getFoodTypesForRestaurant(ctx, rest.ID)
		rest.FoodTypes = foodTypes
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rest)
}

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	var rest models.Restaurant
	err = database.GetPool().QueryRow(ctx,
		`UPDATE restaurants SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			address = COALESCE($3, address),
			phone = COALESCE($4, phone),
			website = COALESCE($5, website),
			latitude = COALESCE($6, latitude),
			longitude = COALESCE($7, longitude),
			google_place_id = COALESCE($8, google_place_id),
			category_id = COALESCE($9, category_id),
			updated_at = NOW()
		WHERE id = $10
		RETURNING id, name, description, address, phone, website, latitude, longitude, google_place_id, category_id, created_at, updated_at`,
		req.Name, req.Description, req.Address, req.Phone, req.Website, req.Latitude, req.Longitude, req.GooglePlaceID, req.CategoryID, id,
	).Scan(
		&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
		&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	// Update food types if provided
	if req.FoodTypeIDs != nil {
		if err := setFoodTypesForRestaurant(ctx, rest.ID, req.FoodTypeIDs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	foodTypes, _ := getFoodTypesForRestaurant(ctx, rest.ID)
	rest.FoodTypes = foodTypes

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rest)
}

func DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	result, err := database.GetPool().Exec(context.Background(),
		"DELETE FROM restaurants WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
