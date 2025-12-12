# Fooder - Restaurant Rating App

A full-stack restaurant rating application with Google Maps integration.

## Tech Stack

- **Backend**: Go with Gorilla Mux
- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Database**: PostgreSQL
- **Containerization**: Docker Compose

## Features

- Search restaurants via Google Maps API
- Rate restaurants on food, service, and ambiance
- Manage cultural categories (Italian, Asian, etc.)
- Manage food types (Pizza, Sushi, etc.)
- View restaurant locations on embedded map
- Get directions to restaurants
- Restaurant suggestion system (pending, approved, tested, rejected statuses)
- Include suggested restaurants in search results
- Upload menu photos (supports AWS S3 or local storage)
- Dark/Light theme toggle

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Google Maps API key with Places API enabled

### Setup

1. Create a `.env` file in the project root:

```bash
cp .env.example .env
```

2. Add your Google Maps API key to `.env`:

```
GOOGLE_MAPS_API_KEY=your_api_key_here

# Optional: Configure AWS S3 for photo storage (otherwise uses local storage)
# AWS_ACCESS_KEY_ID=your_aws_access_key_id
# AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
# AWS_REGION=us-east-1
# S3_BUCKET_NAME=your-bucket-name
```

3. Start all services:

```bash
docker compose up --build
```

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api

## Development

### Running locally without Docker

**Backend:**
```bash
cd backend
export DATABASE_URL="postgres://fooder:fooder_secret@localhost:5432/fooder?sslmode=disable"
export GOOGLE_MAPS_API_KEY="your_key"
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

**Database only:**
```bash
docker compose up db
```

## API Endpoints

### Restaurants
- `GET /api/restaurants` - List all restaurants
- `GET /api/restaurants/:id` - Get restaurant details
- `POST /api/restaurants` - Create restaurant
- `PUT /api/restaurants/:id` - Update restaurant
- `DELETE /api/restaurants/:id` - Delete restaurant

### Ratings
- `GET /api/restaurants/:id/ratings` - Get restaurant ratings
- `POST /api/ratings` - Create rating
- `DELETE /api/ratings/:id` - Delete rating

### Categories
- `GET /api/categories` - List categories
- `POST /api/categories` - Create category
- `PUT /api/categories/:id` - Update category
- `DELETE /api/categories/:id` - Delete category

### Food Types
- `GET /api/food-types` - List food types
- `POST /api/food-types` - Create food type
- `PUT /api/food-types/:id` - Update food type
- `DELETE /api/food-types/:id` - Delete food type

### Google Maps
- `GET /api/places/search?q=query` - Search places
- `GET /api/places/:placeId` - Get place details
