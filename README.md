# The Nom Database - Restaurant Rating App

A production-ready, full-stack restaurant rating application with Google Maps integration, comprehensive monitoring, and automated testing.

## Tech Stack

- **Backend**: Go 1.23 with Gorilla Mux
- **Frontend**: React 18 + TypeScript + Vite + Tailwind CSS
- **Database**: PostgreSQL with automated migrations
- **Containerization**: Docker Compose
- **Logging**: Structured logging with zerolog
- **Monitoring**: Real-time metrics and request tracing
- **Testing**: Comprehensive test suite with 54.5% middleware coverage

## Features

### Core Functionality
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

### Production Features
- **Structured Logging**: zerolog-based logging with JSON and console formats
- **Request Tracing**: UUID-based request correlation across the entire request lifecycle
- **Real-time Metrics**: Track requests, errors, response times (p50, p95, p99)
- **Performance Monitoring**: Metrics endpoint at `/api/metrics` with detailed statistics
- **Security**: Rate limiting, CORS, input sanitization, security headers
- **Database Migrations**: Automated schema migrations with version control
- **API Documentation**: Interactive Swagger UI at `/api/docs`
- **Comprehensive Testing**: Unit tests, benchmarks, and coverage reports

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
   - Metrics Dashboard: http://localhost:8080/api/metrics

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
make install           # Install frontend dependencies
make test              # Run all tests (backend + frontend)
make test-coverage     # Generate test coverage reports
make benchmark         # Run performance benchmarks
make db-stop           # Stop the database
make clean             # Stop and remove database container and volume
make migrate-up        # Run database migrations
make migrate-version   # Show current migration version
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

### Restaurant Suggestions
- `GET /api/suggestions` - List suggestions
- `POST /api/suggestions` - Create suggestion
- `PATCH /api/suggestions/:id/status` - Update status
- `POST /api/suggestions/:id/convert` - Convert to restaurant
- `DELETE /api/suggestions/:id` - Delete suggestion

### Menu Photos
- `GET /api/restaurants/:id/photos` - Get menu photos
- `POST /api/restaurants/:id/photos` - Upload photo
- `PATCH /api/photos/:id` - Update photo caption
- `DELETE /api/photos/:id` - Delete photo

### Search & Monitoring
- `GET /api/search?q=query` - Global search across restaurants
- `GET /api/metrics` - Real-time performance metrics
- `GET /api/health` - Health check endpoint

## Monitoring & Observability

The application includes comprehensive monitoring and logging features for production use.

### Structured Logging

All HTTP requests are logged with structured fields:

```json
{
  "level": "info",
  "time": "2025-12-30T20:13:44Z",
  "message": "HTTP request completed",
  "request_id": "62d30079-6774-46f6-b623-83680479d9a7",
  "method": "GET",
  "path": "/api/restaurants",
  "ip": "192.168.65.1:55864",
  "duration": 4.333917,
  "status": 200,
  "bytes": 1124
}
```

**Configure logging:**
```bash
# Console format (development)
LOG_FORMAT=console

# JSON format (production)
LOG_FORMAT=json

# Enable debug logging
DEBUG=true
```

### Request Tracing

Every request gets a unique UUID for end-to-end tracing:
- Appears in response headers as `X-Request-ID`
- Included in all log entries
- Can be provided by client or auto-generated

**Example:**
```bash
# Client provides request ID
curl -H "X-Request-ID: my-trace-id" http://localhost:8080/api/restaurants

# Filter logs by request ID
docker compose logs backend | grep "my-trace-id"
```

### Real-time Metrics

Access live performance metrics at `/api/metrics`:

```json
{
  "total_requests": 1234,
  "total_errors": 5,
  "requests_by_method": {
    "GET": 800,
    "POST": 300
  },
  "requests_by_status": {
    "200": 1100,
    "500": 14
  },
  "avg_response_time": "2.5ms",
  "p50_response_time": "1.2ms",
  "p95_response_time": "8.5ms",
  "p99_response_time": "15.3ms",
  "uptime": "2h15m30s"
}
```

**Metrics are automatically logged every 5 minutes.**

For detailed monitoring documentation, see [MONITORING.md](MONITORING.md).

## Testing

The project includes comprehensive automated tests for both backend and frontend.

### Running Tests

```bash
# Run all tests
make test

# Backend tests only
make test-backend

# Frontend tests only
make test-frontend

# Generate coverage reports
make test-coverage

# Run performance benchmarks
make benchmark
```

### Test Coverage

**Backend:**
- Middleware: 54.5% coverage ✅
  - Logging middleware: Fully tested
  - Metrics collection: Fully tested
  - Request ID middleware: Fully tested
- Services: 17.4% coverage
  - Image processor: 78.6% coverage
- 30+ unit tests
- 3 benchmark tests

**Frontend:**
- StarRating component: 8 tests ✅
- ThemeToggle component: 6 tests ✅
- LazyImage component: Multiple tests ✅
- Vitest + React Testing Library

### Coverage Reports

After running `make test-coverage`, view reports at:
- Backend: `backend/coverage.html`
- Frontend: `frontend/coverage/index.html`

For detailed testing documentation, see [TESTING.md](TESTING.md).

## Database Migrations

The application uses automated database migrations for schema management.

### Migration Commands

```bash
# Run pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Show current version
make migrate-version

# Create new migration
make migrate-create NAME=add_new_feature

# Force specific version
make migrate-force VERSION=4
```

### Migration Files

Located in `db/migrations_new/`:
- `001_initial_schema.sql` - Base tables
- `002_suggestions.sql` - Suggestion system
- `003_menu_photos.sql` - Photo support
- `004_performance_indexes.sql` - Performance optimization
- `005_security_enhancements.sql` - Security improvements

Migrations run automatically on application startup.

## Security Features

The application implements multiple security layers:

### Middleware Stack

1. **Panic Recovery**: Graceful error handling
2. **Request ID**: Request correlation and tracing
3. **Security Headers**: XSS, clickjacking, MIME sniffing protection
4. **Rate Limiting**: 100 requests/minute per IP (burst: 20)
5. **Content-Type Validation**: Enforce correct content types
6. **Request Size Limits**: 10MB maximum request size
7. **Input Sanitization**: HTML/script tag stripping
8. **Compression**: gzip compression for responses
9. **CORS**: Configurable origin restrictions

### Configuration

```bash
# Allowed origins for CORS
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Rate limiting is automatic (100 req/min per IP)
```

## Performance Optimization

### Database Indexes

Performance indexes on:
- Restaurant lookups by category/food type
- Geospatial queries (lat/lng)
- Rating aggregations
- Suggestion status filtering

### Caching

- Image thumbnails are pre-generated
- Static assets served with compression
- Database connection pooling

### Benchmarks

```bash
make benchmark
```

Example results:
```
BenchmarkRequestIDMiddleware:  1437 ns/op
BenchmarkLoggingMiddleware:    Fast middleware chain
BenchmarkMetrics_RecordRequest: Efficient metric recording
```

## Documentation

Comprehensive documentation is available:

- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Complete API reference with examples
- **[MONITORING.md](MONITORING.md)** - Logging, metrics, and observability guide
- **[TESTING.md](TESTING.md)** - Testing strategy and writing tests
- **[MIGRATIONS.md](MIGRATIONS.md)** - Database migration guide
- **[OPTIMIZATIONS.md](OPTIMIZATIONS.md)** - Performance optimization strategies

## Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   ├── server/          # Application entry point
│   │   └── migrate/         # Migration tool
│   ├── internal/
│   │   ├── database/        # Database connection
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # HTTP middleware
│   │   ├── models/          # Data models
│   │   ├── services/        # Business logic
│   │   ├── logger/          # Structured logging
│   │   └── errors/          # Error handling
│   ├── docs/                # Swagger documentation
│   └── db/migrations_new/   # Database migrations
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── pages/           # Page components
│   │   ├── hooks/           # Custom hooks
│   │   ├── services/        # API client
│   │   └── test/            # Test utilities
│   └── coverage/            # Test coverage reports
├── db/                      # Database init scripts
├── docker-compose.yml       # Container orchestration
├── Makefile                 # Development automation
└── .env.example             # Environment template
```

## Contributing

When contributing to this project:

1. **Write tests** for new features
2. **Run tests** before committing: `make test`
3. **Check coverage**: `make test-coverage`
4. **Follow existing patterns** for consistency
5. **Update documentation** when adding features
6. **Use structured logging** for new log messages
7. **Add migrations** for schema changes

## Troubleshooting

### Common Issues

**Frontend can't connect to backend:**
- Check backend is running: `docker compose ps`
- Verify backend health: `curl http://localhost:8080/api/health`
- Check CORS configuration in backend

**Database connection failed:**
- Ensure database is healthy: `docker compose ps`
- Check connection string in `.env`
- Verify port 5432 is not in use

**Google Maps not loading:**
- Verify API key in `.env` file
- Check APIs are enabled in Google Cloud Console
- Review browser console for errors

**High memory usage:**
- Metrics are limited to prevent memory bloat
- Check log aggregation is configured
- Review Docker resource limits

### Debug Mode

Enable detailed logging to troubleshoot issues:

```bash
# In .env file
DEBUG=true
LOG_FORMAT=console
```

View logs:
```bash
# All backend logs
docker compose logs -f backend

# Filter by error level
docker compose logs backend | grep "ERR"

# Filter by request ID
docker compose logs backend | grep "request_id=abc123"
```

## Production Deployment

For production deployment:

1. **Set `LOG_FORMAT=json`** for log aggregation
2. **Configure `ALLOWED_ORIGINS`** for CORS
3. **Use strong database credentials**
4. **Enable SSL/TLS** for database connections
5. **Set up log aggregation** (CloudWatch, Elasticsearch, etc.)
6. **Configure alerts** based on metrics
7. **Use S3** for photo storage instead of local storage
8. **Set up monitoring dashboards** for metrics
9. **Run migrations** before deploying new versions
10. **Test with `make test`** before deployment

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or contributions:
- Create an issue on GitHub
- Check [TESTING.md](TESTING.md) for testing guidelines
- Review [MONITORING.md](MONITORING.md) for debugging help
- See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for API details
