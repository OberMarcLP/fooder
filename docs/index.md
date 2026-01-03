---
layout: default
title: The Nom Database Documentation
---

# The Nom Database

A full-stack restaurant rating application with Google Maps integration.

## Quick Links

- [GitHub Repository](https://github.com/your-username/the-nom-database)
- [Docker Images](https://github.com/your-username/the-nom-database/pkgs/container/the-nom-database)
- [Latest Release](https://github.com/your-username/the-nom-database/releases/latest)

## Documentation

### Getting Started

- **[Quick Start Guide](QUICKSTART_PRODUCTION.md)** - Deploy in 5 minutes
- **[Deployment Summary](DEPLOYMENT_SUMMARY.md)** - Overview of deployment architecture

### Deployment Guides

- **[Full Deployment Guide](DEPLOYMENT.md)** - Comprehensive production deployment
- **[Production Checklist](PRODUCTION_CHECKLIST.md)** - Pre/post deployment checklist

### Features & Configuration

- **[Authentication Guide](AUTHENTICATION.md)** - Setup OIDC/Authentik authentication
- **[Git Workflow Guide](GIT_WORKFLOW.md)** - Branching strategy and release process
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute to the project

## Features

### Restaurant Management
- Create, read, update, delete restaurants
- Google Maps integration for restaurant search
- Embedded maps showing restaurant locations
- Get directions to restaurants
- Upload menu photos (AWS S3 or local storage)

### Rating System
- Multi-dimensional ratings: food, service, ambiance
- Star rating display
- Review comments
- Per-restaurant rating history

### Authentication
- **None Mode**: No authentication (testing/development)
- **Local Mode**: JWT-based authentication with Argon2id password hashing
- **OAuth Mode**: OIDC integration (Authentik, Keycloak, Auth0, etc.)
- **Both Mode**: Local + OIDC authentication

### Suggestion Workflow
- Users can suggest new restaurants
- Suggestions have statuses: pending, approved, tested, rejected
- Approved suggestions appear in search results

## Tech Stack

### Backend
- **Language**: Go 1.24
- **Router**: Gorilla Mux
- **Database**: PostgreSQL 16
- **Authentication**: JWT, OIDC
- **Storage**: AWS S3 / Local filesystem

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Features**: Dark/Light theme toggle, responsive design

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **CI/CD**: GitHub Actions

## Quick Start with Docker

### Using Pre-built Images (Recommended)

```bash
# Pull latest images
docker pull ghcr.io/your-username/the-nom-database/backend:latest
docker pull ghcr.io/your-username/the-nom-database/frontend:latest

# Or use docker-compose with pre-built images
docker compose -f docker-compose.prod.yml up -d
```

### Building Locally

```bash
# Clone repository
git clone https://github.com/your-username/the-nom-database.git
cd the-nom-database

# Start all services
docker compose up -d

# Access application
open http://localhost:3000
```

## Installation

### Prerequisites

- Docker 24.0+ and Docker Compose v2.0+
- PostgreSQL 16 (or use Docker)
- Google Maps API key
- (Optional) AWS S3 for photo storage
- (Optional) OIDC provider (Authentik, Keycloak, etc.)

### Environment Configuration

```bash
# Copy example environment file
cp .env.example .env

# Edit configuration
nano .env
```

Required variables:
- `DATABASE_URL` - PostgreSQL connection string
- `GOOGLE_MAPS_API_KEY` - Google Maps API key
- `AUTH_MODE` - Authentication mode (none/local/oauth/both)

See [Deployment Guide](DEPLOYMENT.md) for complete configuration.

## Docker Images

Pre-built Docker images are automatically published on every release:

```yaml
services:
  backend:
    image: ghcr.io/your-username/the-nom-database/backend:latest
    # Or use specific version
    # image: ghcr.io/your-username/the-nom-database/backend:v1.0.0

  frontend:
    image: ghcr.io/your-username/the-nom-database/frontend:latest
```

Available tags:
- `latest` - Latest stable release from main branch
- `develop` - Latest development version
- `vX.Y.Z` - Specific version (e.g., v1.0.0)
- `vX.Y` - Latest patch version (e.g., v1.0)
- `vX` - Latest minor version (e.g., v1)

## API Documentation

### Health Endpoints

```bash
# Backend health
curl http://localhost:8080/api/health

# Database health
curl http://localhost:8080/api/health/db
```

### Authentication Endpoints

```bash
# Register (local mode)
POST /api/auth/register

# Login (local mode)
POST /api/auth/login

# OIDC Login (oauth mode)
GET /api/auth/oidc/login

# Refresh token
POST /api/auth/refresh

# Logout
POST /api/auth/logout
```

### Restaurant Endpoints

```bash
# List restaurants
GET /api/restaurants

# Get restaurant
GET /api/restaurants/:id

# Create restaurant
POST /api/restaurants

# Update restaurant
PUT /api/restaurants/:id

# Delete restaurant
DELETE /api/restaurants/:id
```

See [API Documentation](API.md) for complete reference.

## Development

### Backend Development

```bash
cd backend
export DATABASE_URL="postgres://nomdb:nomdb_secret@localhost:5432/nomdb?sslmode=disable"
export AUTH_MODE=none
go run ./cmd/server
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

## Release Process

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards compatible)
- **PATCH** version for bug fixes (backwards compatible)

### Creating a Release

```bash
# Create and push tag
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build and publish Docker images
# 2. Create GitHub release
# 3. Generate changelog
# 4. Build platform binaries
```

See [Git Workflow Guide](GIT_WORKFLOW.md) for detailed release process.

## Support

- **Issues**: [GitHub Issues](https://github.com/your-username/the-nom-database/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/the-nom-database/discussions)
- **Documentation**: [https://your-username.github.io/the-nom-database](https://your-username.github.io/the-nom-database)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

### Development Setup

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...` and `npm test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## Acknowledgments

- Google Maps Platform for location services
- PostgreSQL team for the excellent database
- Go and React communities for amazing tools and libraries
- All contributors who help improve this project

---

**Built with ❤️ by the community**
