---
layout: default
title: The Nom Database Documentation
---

# The Nom Database

A full-stack restaurant rating application with Google Maps integration, multi-mode authentication, and automated Docker publishing.

## Quick Links

- [GitHub Repository](https://github.com/obermarclp/the-nom-database)
- [Docker Images](https://github.com/obermarclp/the-nom-database/pkgs/container/the-nom-database)
- [Latest Release](https://github.com/obermarclp/the-nom-database/releases/latest)

## Getting Started

- **[Quick Start](../README.md)** - Get up and running in 5 minutes
- **[Installation Guide](QUICKSTART_PRODUCTION.md)** - Detailed setup instructions

## Documentation

### Deployment

- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment instructions
- **[Deployment Summary](DEPLOYMENT_SUMMARY.md)** - Overview and architecture
- **[Production Checklist](PRODUCTION_CHECKLIST.md)** - Pre/post deployment checklist

### Configuration

- **[Authentication](AUTHENTICATION.md)** - Setup local auth or OIDC (Authentik, Keycloak, etc.)
- **[API Documentation](API_DOCUMENTATION.md)** - API endpoints and usage
- **[Database Migrations](MIGRATIONS.md)** - Database schema and migrations

### Development

- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Git Workflow](GIT_WORKFLOW.md)** - Branching strategy and release process
- **[GitHub Setup](GITHUB_SETUP.md)** - Repository configuration
- **[Testing](TESTING.md)** - Testing strategy and guidelines
- **[Monitoring](MONITORING.md)** - Application monitoring and metrics

### Advanced

- **[Optimizations](OPTIMIZATIONS.md)** - Performance optimization guide

## Features

### Core Functionality
- **Restaurant Management** - Create, update, delete restaurants
- **Rating System** - Multi-dimensional ratings (food, service, ambiance)
- **Google Maps Integration** - Search and locate restaurants
- **Photo Management** - Upload menu photos (S3 or local storage)
- **Suggestion Workflow** - Community restaurant suggestions

### Authentication
- **None Mode** - No authentication (testing/development)
- **Local Mode** - JWT-based with Argon2id password hashing
- **OAuth Mode** - OIDC integration (Authentik, Keycloak, Auth0, etc.)
- **Both Mode** - Combined local + OIDC authentication

### Infrastructure
- **Docker Publishing** - Automated multi-platform image builds
- **CI/CD Pipeline** - Automated testing and deployment
- **GitHub Pages** - Documentation hosting
- **Multi-platform** - Supports linux/amd64 and linux/arm64

## Tech Stack

### Backend
- **Go 1.24** - Fast, type-safe backend
- **PostgreSQL 16** - Robust relational database
- **Gorilla Mux** - HTTP routing
- **JWT & OIDC** - Flexible authentication

### Frontend
- **React 18** - Modern UI framework
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool
- **Tailwind CSS** - Utility-first styling

### Infrastructure
- **Docker** - Containerization
- **GitHub Actions** - CI/CD automation
- **Nginx** - Reverse proxy (production)

## Quick Start

### Using Pre-built Images (Recommended)

```bash
# Create .env file
cp .env.example .env

# Edit configuration (add your Google Maps API key)
nano .env

# Start services
docker compose up -d

# Access application
open http://localhost:3000
```

### Development Mode

```bash
# Build and run locally
docker compose -f docker-compose.dev.yml up --build
```

## Environment Variables

### Required
```bash
GOOGLE_MAPS_API_KEY  # Google Maps API key
DATABASE_URL          # PostgreSQL connection string
```

### Authentication (Optional)
```bash
AUTH_MODE            # none|local|oauth|both (default: none)
JWT_SECRET_KEY       # Required for local/both modes
OIDC_ISSUER_URL      # Required for oauth/both modes
OIDC_CLIENT_ID       # Required for oauth/both modes
OIDC_CLIENT_SECRET   # Required for oauth/both modes
```

### AWS S3 (Optional)
```bash
AWS_ACCESS_KEY_ID        # AWS access key
AWS_SECRET_ACCESS_KEY    # AWS secret key
S3_BUCKET_NAME           # S3 bucket name
```

## Docker Images

Pre-built images are automatically published on every release:

```bash
# Backend
docker pull ghcr.io/obermarclp/the-nom-database/backend:latest

# Frontend
docker pull ghcr.io/obermarclp/the-nom-database/frontend:latest
```

Available tags:
- `latest` - Latest stable release
- `develop` - Development version
- `vX.Y.Z` - Specific version

## Support

- **Issues**: [GitHub Issues](https://github.com/obermarclp/the-nom-database/issues)
- **Discussions**: [GitHub Discussions](https://github.com/obermarclp/the-nom-database/discussions)
- **Documentation**: [https://obermarclp.github.io/the-nom-database](https://obermarclp.github.io/the-nom-database)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](../LICENSE) file for details.

---

**Need help?** Check the documentation or open an issue on GitHub.
