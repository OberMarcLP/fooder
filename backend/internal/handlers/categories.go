package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fooder/backend/internal/database"
	"github.com/fooder/backend/internal/models"
	"github.com/gorilla/mux"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := database.GetPool().Query(context.Background(),
		"SELECT id, name, created_at, updated_at FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var c models.Category
	err = database.GetPool().QueryRow(context.Background(),
		"SELECT id, name, created_at, updated_at FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var c models.Category
	err := database.GetPool().QueryRow(context.Background(),
		"INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at",
		req.Name).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var c models.Category
	err = database.GetPool().QueryRow(context.Background(),
		"UPDATE categories SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING id, name, created_at, updated_at",
		req.Name, id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	result, err := database.GetPool().Exec(context.Background(),
		"DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
