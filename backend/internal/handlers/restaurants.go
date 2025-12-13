package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nomdb/backend/internal/database"
	"github.com/nomdb/backend/internal/logger"
	"github.com/nomdb/backend/internal/models"
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

func getFoodTypesForSuggestion(ctx context.Context, suggestionID int) ([]models.FoodType, error) {
	rows, err := database.GetPool().Query(ctx,
		`SELECT ft.id, ft.name, ft.created_at, ft.updated_at
		FROM food_types ft
		JOIN suggestion_food_types sft ON ft.id = sft.food_type_id
		WHERE sft.suggestion_id = $1
		ORDER BY ft.name`, suggestionID)
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

	// Always include suggestions in the restaurant list
	includeSuggestions := true

	// Build dynamic query with filters using UNION to include both restaurants and suggestions
	var args []interface{}

	// Build restaurant query
	var restaurantConditions []string
	var restaurantArgs []interface{}
	restaurantArgIndex := 1

	if categoryID != "" {
		if catID, err := strconv.Atoi(categoryID); err == nil {
			restaurantConditions = append(restaurantConditions, fmt.Sprintf("r.category_id = $%d", restaurantArgIndex))
			restaurantArgs = append(restaurantArgs, catID)
			restaurantArgIndex++
		}
	}

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
				placeholders[i] = fmt.Sprintf("$%d", restaurantArgIndex)
				restaurantArgs = append(restaurantArgs, id)
				restaurantArgIndex++
			}
			restaurantConditions = append(restaurantConditions, fmt.Sprintf(`r.id IN (
				SELECT DISTINCT restaurant_id FROM restaurant_food_types
				WHERE food_type_id IN (%s)
			)`, strings.Join(placeholders, ",")))
		}
	}

	// Location/radius filter
	var distanceSelect string
	var distanceOrder string
	var hasDistance bool

	if lat != "" && lng != "" && radius != "" {
		latVal, latErr := strconv.ParseFloat(lat, 64)
		lngVal, lngErr := strconv.ParseFloat(lng, 64)
		radiusVal, radErr := strconv.ParseFloat(radius, 64)

		if latErr == nil && lngErr == nil && radErr == nil {
			hasDistance = true
			distanceSelect = fmt.Sprintf(`,
				(6371 * acos(
					cos(radians($%d)) * cos(radians(r.latitude)) *
					cos(radians(r.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(r.latitude))
				)) as distance`, restaurantArgIndex, restaurantArgIndex+1, restaurantArgIndex+2)
			restaurantArgs = append(restaurantArgs, latVal, lngVal, latVal)

			restaurantConditions = append(restaurantConditions, fmt.Sprintf(`r.latitude IS NOT NULL AND r.longitude IS NOT NULL AND
				(6371 * acos(
					cos(radians($%d)) * cos(radians(r.latitude)) *
					cos(radians(r.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(r.latitude))
				)) <= $%d`, restaurantArgIndex+3, restaurantArgIndex+4, restaurantArgIndex+5, restaurantArgIndex+6))
			restaurantArgs = append(restaurantArgs, latVal, lngVal, latVal, radiusVal)
			restaurantArgIndex += 7

			distanceOrder = "distance ASC,"
		}
	}

	restaurantWhereClause := ""
	if len(restaurantConditions) > 0 {
		restaurantWhereClause = "WHERE " + strings.Join(restaurantConditions, " AND ")
	}

	restaurantQuery := fmt.Sprintf(`
		SELECT
			r.id, r.name, r.description, r.address, r.phone, r.website, r.latitude, r.longitude,
			r.google_place_id, r.category_id, r.created_at, r.updated_at,
			c.id, c.name,
			COALESCE(AVG(rt.food_rating), 0) as avg_food,
			COALESCE(AVG(rt.service_rating), 0) as avg_service,
			COALESCE(AVG(rt.ambiance_rating), 0) as avg_ambiance,
			COUNT(rt.id) as rating_count,
			false as is_suggestion,
			NULL::integer as suggestion_id,
			NULL::text as status
			%s
		FROM restaurants r
		LEFT JOIN categories c ON r.category_id = c.id
		LEFT JOIN ratings rt ON r.id = rt.restaurant_id
		%s
		GROUP BY r.id, c.id
	`, distanceSelect, restaurantWhereClause)

	args = restaurantArgs

	// Build suggestion query if requested
	var finalQuery string
	if includeSuggestions {
		// Build suggestion conditions similar to restaurant conditions
		var suggestionConditions []string
		suggestionArgIndex := len(args) + 1

		if categoryID != "" {
			if catID, err := strconv.Atoi(categoryID); err == nil {
				suggestionConditions = append(suggestionConditions, fmt.Sprintf("s.suggested_category_id = $%d", suggestionArgIndex))
				args = append(args, catID)
				suggestionArgIndex++
			}
		}

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
					placeholders[i] = fmt.Sprintf("$%d", suggestionArgIndex)
					args = append(args, id)
					suggestionArgIndex++
				}
				suggestionConditions = append(suggestionConditions, fmt.Sprintf(`s.id IN (
					SELECT DISTINCT suggestion_id FROM suggestion_food_types
					WHERE food_type_id IN (%s)
				)`, strings.Join(placeholders, ",")))
			}
		}

		// Add status filter to only show pending suggestions
		suggestionConditions = append(suggestionConditions, "s.status = 'pending'")

		var suggestionDistanceSelect string
		if hasDistance {
			latVal, _ := strconv.ParseFloat(lat, 64)
			lngVal, _ := strconv.ParseFloat(lng, 64)
			radiusVal, _ := strconv.ParseFloat(radius, 64)

			suggestionDistanceSelect = fmt.Sprintf(`,
				(6371 * acos(
					cos(radians($%d)) * cos(radians(s.latitude)) *
					cos(radians(s.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(s.latitude))
				)) as distance`, suggestionArgIndex, suggestionArgIndex+1, suggestionArgIndex+2)
			args = append(args, latVal, lngVal, latVal)

			suggestionConditions = append(suggestionConditions, fmt.Sprintf(`s.latitude IS NOT NULL AND s.longitude IS NOT NULL AND
				(6371 * acos(
					cos(radians($%d)) * cos(radians(s.latitude)) *
					cos(radians(s.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(s.latitude))
				)) <= $%d`, suggestionArgIndex+3, suggestionArgIndex+4, suggestionArgIndex+5, suggestionArgIndex+6))
			args = append(args, latVal, lngVal, latVal, radiusVal)
		}

		suggestionWhereClause := ""
		if len(suggestionConditions) > 0 {
			suggestionWhereClause = "WHERE " + strings.Join(suggestionConditions, " AND ")
		}

		suggestionQuery := fmt.Sprintf(`
			SELECT
				s.id, s.name, NULL::text as description, s.address, s.phone, s.website, s.latitude, s.longitude,
				s.google_place_id, s.suggested_category_id as category_id, s.created_at, s.updated_at,
				c.id, c.name,
				0.0 as avg_food,
				0.0 as avg_service,
				0.0 as avg_ambiance,
				0 as rating_count,
				true as is_suggestion,
				s.id as suggestion_id,
				s.status
				%s
			FROM restaurant_suggestions s
			LEFT JOIN categories c ON s.suggested_category_id = c.id
			%s
		`, suggestionDistanceSelect, suggestionWhereClause)

		finalQuery = fmt.Sprintf(`
			SELECT * FROM (
				%s
				UNION ALL
				%s
			) combined
			ORDER BY %s created_at DESC
		`, restaurantQuery, suggestionQuery, distanceOrder)
	} else {
		finalQuery = fmt.Sprintf(`
			%s
			ORDER BY %s r.created_at DESC
		`, restaurantQuery, distanceOrder)
	}

	rows, err := database.GetPool().Query(ctx, finalQuery, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

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
				&rest.IsSuggestion, &rest.SuggestionID, &rest.Status,
				&distance,
			)
		} else {
			err = rows.Scan(
				&rest.ID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
				&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
				&catID, &catName,
				&avgFood, &avgService, &avgAmbiance, &ratingCount,
				&rest.IsSuggestion, &rest.SuggestionID, &rest.Status,
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

		// Get food types for this restaurant or suggestion
		var foodTypes []models.FoodType
		if rest.IsSuggestion {
			foodTypes, err = getFoodTypesForSuggestion(ctx, rest.ID)
		} else {
			foodTypes, err = getFoodTypesForRestaurant(ctx, rest.ID)
		}
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
		// Check if it's a unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation
				logger.Warn("Duplicate restaurant creation attempt: %s", req.Name)
				if strings.Contains(pgErr.ConstraintName, "google_place_id") {
					http.Error(w, "A restaurant with this Google Place ID already exists", http.StatusConflict)
				} else if strings.Contains(pgErr.ConstraintName, "name_address") {
					http.Error(w, "A restaurant with this name and address already exists", http.StatusConflict)
				} else {
					http.Error(w, "This restaurant already exists", http.StatusConflict)
				}
				return
			}
		}
		logger.Error("Failed to create restaurant: %v", err)
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

// GlobalSearch searches both restaurants and suggestions by name
func GlobalSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	searchPattern := "%" + strings.ToLower(query) + "%"

	// Search restaurants
	restaurantsQuery := `
		SELECT DISTINCT
			r.id, r.name, r.description, r.address, r.phone, r.website, r.latitude, r.longitude,
			r.google_place_id, r.category_id, r.created_at, r.updated_at,
			c.id, c.name,
			COALESCE(AVG(rat.food_rating), 0) as avg_food,
			COALESCE(AVG(rat.service_rating), 0) as avg_service,
			COALESCE(AVG(rat.ambiance_rating), 0) as avg_ambiance,
			COUNT(rat.id) as rating_count,
			false as is_suggestion,
			NULL::integer as suggestion_id,
			NULL::text as status
		FROM restaurants r
		LEFT JOIN categories c ON r.category_id = c.id
		LEFT JOIN ratings rat ON r.id = rat.restaurant_id
		WHERE LOWER(r.name) LIKE $1
		GROUP BY r.id, r.name, r.description, r.address, r.phone, r.website, r.latitude, r.longitude,
			r.google_place_id, r.category_id, r.created_at, r.updated_at, c.id, c.name

		UNION ALL

		SELECT
			NULL::integer, s.name, NULL::text, s.address, s.phone, s.website, s.latitude, s.longitude,
			s.google_place_id, s.suggested_category_id, s.created_at, s.updated_at,
			c.id, c.name,
			0::float, 0::float, 0::float, 0::integer,
			true as is_suggestion,
			s.id as suggestion_id,
			s.status
		FROM restaurant_suggestions s
		LEFT JOIN categories c ON s.suggested_category_id = c.id
		WHERE LOWER(s.name) LIKE $1
			AND s.status = 'pending'

		ORDER BY 2
		LIMIT 20
	`

	rows, err := database.GetPool().Query(ctx, restaurantsQuery, searchPattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []models.Restaurant{}
	for rows.Next() {
		var rest models.Restaurant
		var restaurantID *int
		var catID *int
		var catName *string
		var avgFood, avgService, avgAmbiance float64
		var ratingCount int
		var isSuggestion bool
		var suggestionID *int
		var status *string

		err := rows.Scan(
			&restaurantID, &rest.Name, &rest.Description, &rest.Address, &rest.Phone, &rest.Website, &rest.Latitude, &rest.Longitude,
			&rest.GooglePlaceID, &rest.CategoryID, &rest.CreatedAt, &rest.UpdatedAt,
			&catID, &catName,
			&avgFood, &avgService, &avgAmbiance, &ratingCount,
			&isSuggestion, &suggestionID, &status,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set restaurant ID if it's not null (i.e., it's a real restaurant, not a suggestion)
		if restaurantID != nil {
			rest.ID = *restaurantID
		}

		if catID != nil && catName != nil {
			rest.Category = &models.Category{ID: *catID, Name: *catName}
		}

		rest.IsSuggestion = isSuggestion
		rest.SuggestionID = suggestionID
		rest.Status = status

		// Get food types
		if isSuggestion && suggestionID != nil {
			foodTypes, _ := getFoodTypesForSuggestion(ctx, *suggestionID)
			rest.FoodTypes = foodTypes
		} else if rest.ID > 0 {
			foodTypes, _ := getFoodTypesForRestaurant(ctx, rest.ID)
			rest.FoodTypes = foodTypes
		}

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

		results = append(results, rest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
