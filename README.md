# The Nom Database - Restaurant Rating App

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
- Google Maps API key with Places API enabled (see setup guide below)

### Google Maps API Setup

The application requires a Google Maps API key with specific APIs enabled. Follow these steps:

#### 1. Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click "Select a project" → "New Project"
3. Enter a project name (e.g., "The Nom Database")
4. Click "Create"

#### 2. Enable Required APIs

Your API key needs these three APIs enabled:

1. **Maps JavaScript API** (for frontend map display)
   - Go to: https://console.cloud.google.com/apis/library/maps-backend.googleapis.com
   - Click "Enable"

2. **Places API** (for restaurant search)
   - Go to: https://console.cloud.google.com/apis/library/places-backend.googleapis.com
   - Click "Enable"

3. **Geocoding API** (for address lookups)
   - Go to: https://console.cloud.google.com/apis/library/geocoding-backend.googleapis.com
   - Click "Enable"

**Quick link:** You can also enable all Maps APIs at once: https://console.cloud.google.com/google/maps-apis/start

#### 3. Create API Key

1. Go to [API Credentials](https://console.cloud.google.com/apis/credentials)
2. Click "Create Credentials" → "API Key"
3. Copy the generated API key
4. (Optional but recommended) Click "Restrict Key" to add security:
   - Under "API restrictions", select "Restrict key"
   - Check: Maps JavaScript API, Places API, Geocoding API
   - Click "Save"

#### 4. Add API Key to Project

1. Create a `.env` file in the project root:

```bash
cp .env.example .env
```

2. Edit the `.env` file and configure:

```bash
# Database Configuration (for local development without Docker)
DB_HOST=localhost
DB_PORT=5432
DB_USER=nomdb
DB_PASSWORD=nomdb_secret
DB_NAME=nomdb
DB_SSLMODE=disable

# Or use a connection string (overrides individual settings if set)
# DATABASE_URL=postgres://nomdb:nomdb_secret@localhost:5432/nomdb?sslmode=disable

# Google Maps API Configuration (optional - leave empty to disable)
GOOGLE_MAPS_API_KEY=your_api_key_here

# AWS S3 Configuration (optional - uses local storage if not configured)
# AWS_ACCESS_KEY_ID=your_aws_access_key_id
# AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
# AWS_REGION=us-east-1
# S3_BUCKET_NAME=your-bucket-name

# Debug Mode (optional - set to true for detailed logging)
DEBUG=true
```

#### Logging and Debug Mode

The backend includes a comprehensive logging system with color-coded output and emojis for easy readability:

- **INFO** (green): General information about server operations
- **DEBUG** (cyan): Detailed information for debugging (only shown when DEBUG=true)
- **WARN** (yellow): Warning messages for non-critical issues
- **ERROR** (red): Error messages for failures

**Enable Debug Mode:**
Set `DEBUG=true` in your `.env` file to see detailed logging including:
- Database connection details
- Google Maps API calls and responses
- S3 upload/delete operations
- HTTP request/response details with timing

**Example logs:**
```
[2025-12-13 09:35:06] INFO  🚀 Starting The Nom Database server...
[2025-12-13 09:35:06] INFO  🔌 Connecting to database...
[2025-12-13 09:35:06] INFO  ✅ Database connected successfully
[2025-12-13 09:35:06] INFO  🗺️  Google Maps service initialized
[2025-12-13 09:35:06] WARN  ⚠️  AWS S3 not configured - using local storage for photos
[2025-12-13 09:35:06] INFO  🌐 Server listening on http://localhost:8080
[2025-12-13 09:35:13] INFO  ← ✅ GET /api/health 200 (53.334µs) 2 bytes
```

**Note:** Google Maps features are optional. If you don't have an API key, you can:
- Leave `GOOGLE_MAPS_API_KEY` empty or remove it
- Manually add restaurants without using the search feature
- Maps won't display, but all other features will work

### Running the Application

1. Start all services:

```bash
docker compose up --build
```

2. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api
   - Swagger UI (API Documentation): http://localhost:8080/api/docs

**Note:** The first build may take a few minutes. Subsequent starts will be faster.

3. Check service health:

```bash
docker compose ps
```

All services include health checks:
- **Database**: Checks PostgreSQL is ready to accept connections
- **Backend**: Checks `/api/health` endpoint responds
- **Frontend**: Checks nginx is serving the application

Services will show as `(healthy)` when ready. The frontend waits for the backend to be healthy, and the backend waits for the database to be healthy, ensuring proper startup order.

## Development

### Running locally without Docker

You can use the Makefile for convenient local development commands:

#### Using Makefile (Recommended)

**View all available commands:**
```bash
make help
```

**Start local development (3 separate terminals):**

Terminal 1 - Start database:
```bash
make db
```

Terminal 2 - Start backend:
```bash
make backend
```

Terminal 3 - Start frontend:
```bash
make frontend
```

**Other useful commands:**
```bash
make install      # Install frontend dependencies
make db-stop      # Stop the database
make clean        # Stop and remove database container and volume
```

**Note:** The Makefile automatically loads environment variables from your `.env` file.

#### Manual Setup (Without Makefile)

**Database only:**
```bash
docker compose up db
```

**Backend:**
```bash
cd backend
export DATABASE_URL="postgres://nomdb:nomdb_secret@localhost:5432/nomdb?sslmode=disable"
export GOOGLE_MAPS_API_KEY="your_key"
export DEBUG=true  # Optional: enable debug logging
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## API Documentation

The API includes interactive documentation via Swagger UI, accessible at:
- **Swagger UI**: http://localhost:8080/api/docs

The Swagger UI provides:
- Complete API endpoint documentation
- Request/response schemas with examples
- Interactive "Try it out" functionality to test endpoints directly
- Authentication support (if configured)

For a quick reference, see the endpoints below, or use Swagger UI for detailed schemas and testing.

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
