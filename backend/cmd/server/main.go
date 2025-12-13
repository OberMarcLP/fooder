package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nomdb/backend/internal/database"
	"github.com/nomdb/backend/internal/handlers"
	"github.com/nomdb/backend/internal/logger"
	"github.com/nomdb/backend/internal/middleware"
	"github.com/nomdb/backend/internal/services"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	logger.Info("🚀 Starting The Nom Database server...")
	if logger.IsDebugMode() {
		logger.Debug("🐛 Debug mode enabled - detailed logging active")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Google Maps service
	_ = services.NewGoogleMapsService()

	// Initialize S3 service (optional - falls back to local storage if not configured)
	if err := services.InitS3(); err != nil {
		logger.Debug("S3 initialization skipped: %v", err)
	}

	// Create router
	r := mux.NewRouter()

	// Create uploads directory and serve static files
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir+"/menu_photos", 0755); err != nil {
		logger.Warn("Failed to create uploads directory: %v", err)
	} else {
		logger.Debug("📁 Uploads directory ready: %s", uploadsDir)
	}
	r.PathPrefix("/api/uploads/").Handler(
		http.StripPrefix("/api/uploads/", http.FileServer(http.Dir(uploadsDir))))

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

	// Global Search
	api.HandleFunc("/search", handlers.GlobalSearch).Methods("GET")

	// Ratings
	api.HandleFunc("/restaurants/{restaurantId}/ratings", handlers.GetRatings).Methods("GET")
	api.HandleFunc("/ratings", handlers.CreateRating).Methods("POST")
	api.HandleFunc("/ratings/{id}", handlers.DeleteRating).Methods("DELETE")

	// Google Maps
	api.HandleFunc("/places/search", handlers.SearchPlaces).Methods("GET")
	api.HandleFunc("/places/{placeId}", handlers.GetPlaceDetails).Methods("GET")
	api.HandleFunc("/geocode/cities", handlers.GeocodeCities).Methods("GET")

	// Restaurant Suggestions
	api.HandleFunc("/suggestions", handlers.GetSuggestions).Methods("GET")
	api.HandleFunc("/suggestions/{id}", handlers.GetSuggestion).Methods("GET")
	api.HandleFunc("/suggestions", handlers.CreateSuggestion).Methods("POST")
	api.HandleFunc("/suggestions/{id}/status", handlers.UpdateSuggestionStatus).Methods("PATCH")
	api.HandleFunc("/suggestions/{id}/convert", handlers.ConvertSuggestion).Methods("POST")
	api.HandleFunc("/suggestions/{id}", handlers.DeleteSuggestion).Methods("DELETE")

	// Menu Photos
	api.HandleFunc("/restaurants/{restaurantId}/photos", handlers.GetMenuPhotos).Methods("GET")
	api.HandleFunc("/restaurants/{restaurantId}/photos", handlers.UploadMenuPhoto).Methods("POST")
	api.HandleFunc("/photos/{id}", handlers.UpdatePhotoCaption).Methods("PATCH")
	api.HandleFunc("/photos/{id}", handlers.DeleteMenuPhoto).Methods("DELETE")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Serve the swagger.yaml file first
	api.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "./docs/swagger.yaml")
	}).Methods("GET")

	// Swagger UI - serve at /api/docs (must be after swagger.yaml)
	api.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/swagger.yaml"),
	))

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Apply middleware chain
	handler := middleware.LoggingMiddleware(c.Handler(r))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("🌐 Server listening on http://localhost:%s", port)
	logger.Info("📡 API available at http://localhost:%s/api", port)
	logger.Info("📚 Swagger UI available at http://localhost:%s/api/docs", port)
	logger.Info("✅ Server ready to accept connections")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatal("Server failed to start: %v", err)
	}
}
