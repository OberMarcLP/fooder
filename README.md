# The Nom Database

[![CI](https://github.com/obermarclp/the-nom-database/actions/workflows/ci.yml/badge.svg)](https://github.com/obermarclp/the-nom-database/actions/workflows/ci.yml)
[![Docker](https://github.com/obermarclp/the-nom-database/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/obermarclp/the-nom-database/actions/workflows/docker-publish.yml)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://obermarclp.github.io/the-nom-database)

A modern, full-stack restaurant rating application with Google Maps integration, flexible authentication (local/OIDC), and automated Docker publishing.

## ✨ Features

- 🗺️ **Google Maps Integration** - Search and locate restaurants
- ⭐ **Multi-dimensional Ratings** - Rate food, service, and ambiance
- 📸 **Photo Management** - Upload menu photos (S3 or local storage)
- 🔐 **Flexible Authentication** - Support for local, OIDC, or no auth
- 🐳 **Pre-built Docker Images** - No build required, multi-platform support
- 🌓 **Dark/Light Mode** - Theme toggle with localStorage persistence

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Google Maps API key ([Get one here](https://developers.google.com/maps/documentation/javascript/get-api-key))

### Installation

```bash
# Clone repository
git clone https://github.com/obermarclp/the-nom-database.git
cd the-nom-database

# Create environment file
cp .env.example .env

# Edit .env and add your Google Maps API key
nano .env

# Start application (uses pre-built images)
docker compose up -d

# Access application
open http://localhost:3000
```

That's it! The application will be running on http://localhost:3000

### Development Mode

To build images locally and develop:

```bash
# Use development compose file
docker compose -f docker-compose.dev.yml up --build
```

## 📖 Documentation

**Complete documentation** is available at [obermarclp.github.io/the-nom-database](https://obermarclp.github.io/the-nom-database)

Quick links:
- [Deployment Guide](docs/DEPLOYMENT.md) - Production deployment
- [Authentication Setup](docs/AUTHENTICATION.md) - Configure OIDC/local auth
- [API Documentation](docs/API_DOCUMENTATION.md) - API reference
- [Contributing Guide](docs/CONTRIBUTING.md) - How to contribute

## 🐳 Docker Images

Pre-built images are automatically published to GitHub Container Registry:

```bash
# Pull latest images
docker pull ghcr.io/obermarclp/the-nom-database/backend:latest
docker pull ghcr.io/obermarclp/the-nom-database/frontend:latest
```

Available tags: `latest`, `develop`, `v1.0.0`, etc.

## 🛠️ Tech Stack

- **Backend**: Go 1.24, PostgreSQL 16, Gorilla Mux
- **Frontend**: React 18, TypeScript, Vite, Tailwind CSS
- **Auth**: JWT (Argon2id), OIDC (Authentik, Keycloak, Auth0, etc.)
- **Infrastructure**: Docker, GitHub Actions, Nginx

## ⚙️ Configuration

### Environment Variables

**Required:**
- `GOOGLE_MAPS_API_KEY` - Your Google Maps API key

**Authentication (Optional):**
- `AUTH_MODE` - Authentication mode: `none`, `local`, `oauth`, or `both` (default: `none`)
- `JWT_SECRET_KEY` - Required for `local` or `both` modes
- `OIDC_ISSUER_URL` - OIDC provider URL (for `oauth` or `both`)
- `OIDC_CLIENT_ID` - OIDC client ID
- `OIDC_CLIENT_SECRET` - OIDC client secret

**AWS S3 (Optional):**
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `S3_BUCKET_NAME`

See [.env.example](.env.example) for all options.

## 📊 API Endpoints

### Health
```
GET  /api/health       # Backend health check
GET  /api/health/db    # Database health check
```

### Restaurants
```
GET    /api/restaurants      # List restaurants
GET    /api/restaurants/:id  # Get restaurant
POST   /api/restaurants      # Create restaurant
PUT    /api/restaurants/:id  # Update restaurant
DELETE /api/restaurants/:id  # Delete restaurant
```

### Ratings
```
GET    /api/restaurants/:id/ratings  # Get ratings
POST   /api/ratings                  # Create rating
DELETE /api/ratings/:id              # Delete rating
```

See [API Documentation](docs/API_DOCUMENTATION.md) for complete reference.

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](docs/CONTRIBUTING.md) first.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...` for backend, `npm test` for frontend)
5. Commit (`git commit -m 'feat: add amazing feature'`)
6. Push (`git push origin feature/amazing-feature`)
7. Open a Pull Request

See [Git Workflow](docs/GIT_WORKFLOW.md) for branching strategy.

## 📝 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Google Maps Platform for location services
- PostgreSQL team for the excellent database
- Go and React communities for amazing tools and libraries

## 📞 Support

- **Documentation**: [obermarclp.github.io/the-nom-database](https://obermarclp.github.io/the-nom-database)
- **Issues**: [GitHub Issues](https://github.com/obermarclp/the-nom-database/issues)
- **Discussions**: [GitHub Discussions](https://github.com/obermarclp/the-nom-database/discussions)

---

**Built with ❤️ using Go, React, and Docker**
