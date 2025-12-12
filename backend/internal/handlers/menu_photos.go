package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fooder/backend/internal/database"
	"github.com/fooder/backend/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	maxUploadSize = 10 << 20 // 10MB
	uploadsDir    = "./uploads/menu_photos"
)

func init() {
	// Create uploads directory if it doesn't exist
	os.MkdirAll(uploadsDir, 0755)
}

// GetMenuPhotos retrieves all photos for a restaurant
func GetMenuPhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantID, err := strconv.Atoi(vars["restaurantId"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	rows, err := database.GetPool().Query(ctx,
		`SELECT id, restaurant_id, filename, original_filename, caption, file_size, mime_type, created_at, updated_at
		FROM menu_photos
		WHERE restaurant_id = $1
		ORDER BY created_at DESC`, restaurantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	photos := []models.MenuPhoto{}
	for rows.Next() {
		var photo models.MenuPhoto
		if err := rows.Scan(
			&photo.ID, &photo.RestaurantID, &photo.Filename, &photo.OriginalFilename,
			&photo.Caption, &photo.FileSize, &photo.MimeType, &photo.CreatedAt, &photo.UpdatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Compute URL
		photo.URL = fmt.Sprintf("/api/uploads/menu_photos/%s", photo.Filename)
		photos = append(photos, photo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

// UploadMenuPhoto handles file upload for menu photos
func UploadMenuPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantID, err := strconv.Atoi(vars["restaurantId"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	// Parse multipart form
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Get caption
	caption := r.FormValue("caption")
	if caption == "" {
		http.Error(w, "Caption is required", http.StatusBadRequest)
		return
	}

	// Get file
	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Only image files are allowed", http.StatusBadRequest)
		return
	}

	// Validate specific image types
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !validTypes[contentType] {
		http.Error(w, "Only JPEG, PNG, and WebP images are allowed", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext

	// Create file on disk
	filePath := filepath.Join(uploadsDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	fileSize, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Save to database
	ctx := context.Background()
	var photo models.MenuPhoto
	err = database.GetPool().QueryRow(ctx,
		`INSERT INTO menu_photos (restaurant_id, filename, original_filename, caption, file_size, mime_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, restaurant_id, filename, original_filename, caption, file_size, mime_type, created_at, updated_at`,
		restaurantID, filename, header.Filename, caption, int(fileSize), contentType,
	).Scan(
		&photo.ID, &photo.RestaurantID, &photo.Filename, &photo.OriginalFilename,
		&photo.Caption, &photo.FileSize, &photo.MimeType, &photo.CreatedAt, &photo.UpdatedAt,
	)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Compute URL
	photo.URL = fmt.Sprintf("/api/uploads/menu_photos/%s", photo.Filename)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.UploadPhotoResponse{Photo: photo})
}

// UpdatePhotoCaption updates the caption of a photo
func UpdatePhotoCaption(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Caption string `json:"caption"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Caption == "" {
		http.Error(w, "Caption is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var photo models.MenuPhoto
	err = database.GetPool().QueryRow(ctx,
		`UPDATE menu_photos SET caption = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, restaurant_id, filename, original_filename, caption, file_size, mime_type, created_at, updated_at`,
		req.Caption, id,
	).Scan(
		&photo.ID, &photo.RestaurantID, &photo.Filename, &photo.OriginalFilename,
		&photo.Caption, &photo.FileSize, &photo.MimeType, &photo.CreatedAt, &photo.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	photo.URL = fmt.Sprintf("/api/uploads/menu_photos/%s", photo.Filename)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photo)
}

// DeleteMenuPhoto deletes a photo
func DeleteMenuPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Get filename before deleting from DB
	var filename string
	err = database.GetPool().QueryRow(ctx,
		"SELECT filename FROM menu_photos WHERE id = $1", id).Scan(&filename)
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Delete from database
	result, err := database.GetPool().Exec(ctx,
		"DELETE FROM menu_photos WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Delete file from disk (non-fatal if fails)
	filePath := filepath.Join(uploadsDir, filename)
	os.Remove(filePath)

	w.WriteHeader(http.StatusNoContent)
}
