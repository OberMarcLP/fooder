package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fooder/backend/internal/database"
	"github.com/fooder/backend/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Categories
	api.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	api.HandleFunc("/categories/{id}", handlers.GetCategory).Methods("GET")
	api.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")
	api.HandleFunc("/categories/{id}", handlers.UpdateCategory).Methods("PUT")
	api.HandleFunc("/categories/{id}", handlers.DeleteCategory).Methods("DELETE")

	// Food Types
	api.HandleFunc("/food-types", handlers.GetFoodTypes).Methods("GET")
	api.HandleFunc("/food-types/{id}", handlers.GetFoodType).Methods("GET")
	api.HandleFunc("/food-types", handlers.CreateFoodType).Methods("POST")
	api.HandleFunc("/food-types/{id}", handlers.UpdateFoodType).Methods("PUT")
	api.HandleFunc("/food-types/{id}", handlers.DeleteFoodType).Methods("DELETE")

	// Restaurants
	api.HandleFunc("/restaurants", handlers.GetRestaurants).Methods("GET")
	api.HandleFunc("/restaurants/{id}", handlers.GetRestaurant).Methods("GET")
	api.HandleFunc("/restaurants", handlers.CreateRestaurant).Methods("POST")
	api.HandleFunc("/restaurants/{id}", handlers.UpdateRestaurant).Methods("PUT")
	api.HandleFunc("/restaurants/{id}", handlers.DeleteRestaurant).Methods("DELETE")

	// Ratings
	api.HandleFunc("/restaurants/{restaurantId}/ratings", handlers.GetRatings).Methods("GET")
	api.HandleFunc("/ratings", handlers.CreateRating).Methods("POST")
	api.HandleFunc("/ratings/{id}", handlers.DeleteRating).Methods("DELETE")

	// Google Maps
	api.HandleFunc("/places/search", handlers.SearchPlaces).Methods("GET")
	api.HandleFunc("/places/{placeId}", handlers.GetPlaceDetails).Methods("GET")
	api.HandleFunc("/geocode/cities", handlers.GeocodeCities).Methods("GET")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
