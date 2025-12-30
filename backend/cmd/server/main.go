package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nomdb/backend/internal/database"
	"github.com/nomdb/backend/internal/handlers"
	"github.com/nomdb/backend/internal/logger"
	"github.com/nomdb/backend/internal/middleware"
	"github.com/nomdb/backend/internal/services"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"

	_ "github.com/nomdb/backend/docs" // Import generated docs
)

// @title The Nom Database API
// @version 1.0
// @description Restaurant rating and discovery API with Google Maps integration
// @description
// @description This API provides endpoints for managing restaurants, ratings, categories, and food types.
// @description It integrates with Google Maps for restaurant search and location data.
//
// @contact.name API Support
// @contact.url https://github.com/your-username/the-nom-database
// @contact.email support@nomdb.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api
//
// @schemes http https
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @tag.name Restaurants
// @tag.description Restaurant management endpoints
//
// @tag.name Ratings
// @tag.description Restaurant rating endpoints
//
// @tag.name Categories
// @tag.description Cultural category management
//
// @tag.name Food Types
// @tag.description Food type/cuisine management
//
// @tag.name Suggestions
// @tag.description Restaurant suggestion workflow
//
// @tag.name Google Maps
// @tag.description Google Places API integration
//
// @tag.name Photos
// @tag.description Menu photo management
//
// @tag.name Search
// @tag.description Global search functionality
//
// @tag.name Health
// @tag.description Health check endpoints
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

	// Run database migrations
	databaseURL := os.Getenv("DATABASE_URL")
	migrationsPath := "/app/db/migrations_new"
	if _, err := os.Stat("./db/migrations_new"); err == nil {
		// Local development
		migrationsPath = "./db/migrations_new"
	}
	if err := database.RunMigrations(databaseURL, migrationsPath); err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}

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
	api.HandleFunc("/restaurants/paginated", handlers.GetRestaurantsPaginated).Methods("GET")
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

	// Health check (support both GET and HEAD for Docker healthcheck)
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == "GET" {
			w.Write([]byte("OK"))
		}
	}).Methods("GET", "HEAD")

	// Metrics endpoint
	api.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		stats := middleware.GetMetrics().GetStats()
		json.NewEncoder(w).Encode(stats)
	}).Methods("GET")

	// Serve the swagger.yaml file first
	api.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "./docs/swagger.yaml")
	}).Methods("GET")

	// Redirect /docs to /docs/ for Swagger UI
	api.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/", http.StatusMovedPermanently)
	}).Methods("GET")

	// Swagger UI - serve at /api/docs/ (must be after swagger.yaml)
	api.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/swagger.yaml"),
	))

	// Initialize rate limiter
	// Allow 100 requests per minute per IP, with burst of 20
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Minute/100), 20)
	// Start cleanup task to prevent memory leaks (run every 10 minutes)
	rateLimiter.StartCleanupTask(10 * time.Minute)
	logger.Info("🔒 Rate limiting enabled: 100 req/min per IP, burst: 20")

	// CORS middleware - more restrictive configuration
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:5173"}
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		// In production, set ALLOWED_ORIGINS to your actual domain
		allowedOrigins = []string{envOrigins}
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300, // Cache preflight for 5 minutes
	})
	logger.Info("🌐 CORS configured for origins: %v", allowedOrigins)

	// Apply middleware chain (order matters)
	// Recovery -> RequestID -> Security headers -> Rate limiting -> Request validation -> Max bytes -> Sanitization -> Compression -> Logging -> CORS -> Router
	handler := middleware.RecoveryMiddleware(
		middleware.RequestIDMiddleware(
			middleware.SecurityHeadersMiddleware(
				middleware.RateLimitMiddleware(rateLimiter)(
					middleware.ValidateContentTypeMiddleware(
						middleware.MaxBytesMiddleware(10 * 1024 * 1024)( // 10MB max request size
							middleware.SanitizeInputMiddleware(
								middleware.CompressionMiddleware(
									middleware.LoggingMiddleware(
										c.Handler(r))))))))))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("🌐 Server listening on http://localhost:%s", port)
	logger.Info("📡 API available at http://localhost:%s/api", port)
	logger.Info("📚 Swagger UI available at http://localhost:%s/api/docs", port)
	logger.Info("🛡️  Security features enabled:")
	logger.Info("   ✓ Panic recovery and error handling")
	logger.Info("   ✓ Rate limiting (100 req/min per IP)")
	logger.Info("   ✓ Request size limits (10MB max)")
	logger.Info("   ✓ Content-Type validation")
	logger.Info("   ✓ Input sanitization")
	logger.Info("   ✓ Security headers (XSS, clickjacking, MIME sniffing protection)")
	logger.Info("   ✓ CORS restrictions")
	logger.Info("✅ Server ready to accept connections")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatal("Server failed to start: %v", err)
	}
}
