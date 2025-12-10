# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Start all services (recommended)
docker compose up --build -d

# Reset database (required after schema changes)
docker compose down -v && docker compose up --build -d

# View logs
docker compose logs -f backend
docker compose logs -f frontend

# Run backend locally (requires DATABASE_URL and GOOGLE_MAPS_API_KEY env vars)
cd backend && go run ./cmd/server

# Run frontend locally
cd frontend && npm install && npm run dev
```

## Architecture

### Backend (Go)
- **Entry point**: `backend/cmd/server/main.go` - sets up Gorilla Mux router with CORS
- **Handlers**: `backend/internal/handlers/` - HTTP handlers for each resource (restaurants, ratings, categories, food_types, google_maps)
- **Models**: `backend/internal/models/models.go` - all domain types and request/response structs
- **Database**: `backend/internal/database/` - PostgreSQL connection using pgx pool
- **Services**: `backend/internal/services/` - external service integrations (Google Maps API)

### Frontend (React + TypeScript)
- **Entry point**: `frontend/src/main.tsx` → `frontend/src/App.tsx`
- **Pages**: `frontend/src/pages/` - main views (HomePage, RestaurantDetail, CategoriesPage, FoodTypesPage)
- **Components**: `frontend/src/components/` - reusable UI (RestaurantForm, RatingForm, StarRating, PlaceSearch, Modal, etc.)
- **API client**: `frontend/src/services/api.ts` - typed fetch wrappers for all backend endpoints
- **Styling**: Tailwind CSS with dark mode support via `useTheme` hook

### Database
- **Schema**: `db/migrations/001_init.sql` - PostgreSQL schema with categories, food_types, restaurants, ratings tables
- **Junction table**: `restaurant_food_types` for many-to-many relationship between restaurants and food types

## Key Data Flows

1. **Adding a restaurant**: PlaceSearch → Google Maps API → RestaurantForm (pre-fills name, address, phone, website, coords) → POST /api/restaurants
2. **Rating calculation**: Ratings stored individually, `avg_rating` computed via SQL aggregation in restaurant queries
3. **Food types**: Multi-select via junction table; handled by `setFoodTypesForRestaurant()` in restaurants handler

## Environment Variables

Required in `.env` file (see `.env.example`):
- `GOOGLE_MAPS_API_KEY` - needs Places API enabled in Google Cloud Console
