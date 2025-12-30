package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nomdb/backend/internal/models"
)

const (
	DefaultPageLimit = 20
	MaxPageLimit     = 100
)

// ParsePaginationParams extracts pagination parameters from request
func ParsePaginationParams(r *http.Request) models.PaginationParams {
	params := models.PaginationParams{
		Limit:  DefaultPageLimit,
		Cursor: "",
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			if limit > 0 && limit <= MaxPageLimit {
				params.Limit = limit
			} else if limit > MaxPageLimit {
				params.Limit = MaxPageLimit
			}
		}
	}

	// Parse cursor
	params.Cursor = r.URL.Query().Get("cursor")

	return params
}

// EncodeCursor creates a base64-encoded cursor from an ID
func EncodeCursor(id int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", id)))
}

// DecodeCursor decodes a base64-encoded cursor to an ID
func DecodeCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0, err
	}

	return id, nil
}

// BuildPaginatedResponse creates a paginated response
func BuildPaginatedResponse(data interface{}, hasMore bool, nextID int) models.PaginatedResponse {
	response := models.PaginatedResponse{
		Data:    data,
		HasMore: hasMore,
	}

	if hasMore && nextID > 0 {
		cursor := EncodeCursor(nextID)
		response.NextCursor = &cursor
	}

	return response
}
