package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nomdb/backend/internal/database"
	"github.com/nomdb/backend/internal/models"
)

func GetFoodTypes(w http.ResponseWriter, r *http.Request) {
	rows, err := database.GetPool().Query(context.Background(),
		"SELECT id, name, created_at, updated_at FROM food_types ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	foodTypes := []models.FoodType{}
	for rows.Next() {
		var ft models.FoodType
		if err := rows.Scan(&ft.ID, &ft.Name, &ft.CreatedAt, &ft.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		foodTypes = append(foodTypes, ft)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foodTypes)
}

func GetFoodType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid food type ID", http.StatusBadRequest)
		return
	}

	var ft models.FoodType
	err = database.GetPool().QueryRow(context.Background(),
		"SELECT id, name, created_at, updated_at FROM food_types WHERE id = $1", id).
		Scan(&ft.ID, &ft.Name, &ft.CreatedAt, &ft.UpdatedAt)
	if err != nil {
		http.Error(w, "Food type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ft)
}

func CreateFoodType(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFoodTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var ft models.FoodType
	err := database.GetPool().QueryRow(context.Background(),
		"INSERT INTO food_types (name) VALUES ($1) RETURNING id, name, created_at, updated_at",
		req.Name).Scan(&ft.ID, &ft.Name, &ft.CreatedAt, &ft.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ft)
}

func UpdateFoodType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid food type ID", http.StatusBadRequest)
		return
	}

	var req models.CreateFoodTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var ft models.FoodType
	err = database.GetPool().QueryRow(context.Background(),
		"UPDATE food_types SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING id, name, created_at, updated_at",
		req.Name, id).Scan(&ft.ID, &ft.Name, &ft.CreatedAt, &ft.UpdatedAt)
	if err != nil {
		http.Error(w, "Food type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ft)
}

func DeleteFoodType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid food type ID", http.StatusBadRequest)
		return
	}

	result, err := database.GetPool().Exec(context.Background(),
		"DELETE FROM food_types WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Food type not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
