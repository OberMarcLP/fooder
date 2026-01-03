#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"
SSL_DIR="nginx/ssl"

# Helper functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    print_info "Checking Docker installation..."
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose v2."
        exit 1
    fi

    print_success "Docker and Docker Compose are installed"
}

check_env_file() {
    print_info "Checking environment configuration..."
    if [ ! -f "$ENV_FILE" ]; then
        print_error "Environment file '$ENV_FILE' not found."
        print_info "Please copy .env.production.example to .env.production and configure it."
        exit 1
    fi

    # Source the env file
    set -a
    source "$ENV_FILE"
    set +a

    # Check required variables
    local missing_vars=()

    if [ -z "$DB_PASSWORD" ] || [ "$DB_PASSWORD" == "CHANGE_ME_STRONG_PASSWORD_HERE" ]; then
        missing_vars+=("DB_PASSWORD")
    fi

    if [ -z "$JWT_SECRET_KEY" ] || [ "$JWT_SECRET_KEY" == "CHANGE_ME_GENERATE_WITH_openssl_rand_base64_64" ]; then
        missing_vars+=("JWT_SECRET_KEY")
    fi

    if [ -z "$GOOGLE_MAPS_API_KEY" ] || [ "$GOOGLE_MAPS_API_KEY" == "your_production_google_maps_api_key" ]; then
        missing_vars+=("GOOGLE_MAPS_API_KEY")
    fi

    if [ "$AUTH_MODE" == "oauth" ] || [ "$AUTH_MODE" == "both" ]; then
        if [ -z "$OIDC_ISSUER_URL" ] || [ "$OIDC_ISSUER_URL" == "https://auth.yourdomain.com/application/o/nom-database/" ]; then
            missing_vars+=("OIDC_ISSUER_URL")
        fi
        if [ -z "$OIDC_CLIENT_ID" ] || [ "$OIDC_CLIENT_ID" == "your_production_client_id" ]; then
            missing_vars+=("OIDC_CLIENT_ID")
        fi
        if [ -z "$OIDC_CLIENT_SECRET" ] || [ "$OIDC_CLIENT_SECRET" == "your_production_client_secret" ]; then
            missing_vars+=("OIDC_CLIENT_SECRET")
        fi
    fi

    if [ ${#missing_vars[@]} -gt 0 ]; then
        print_error "Missing or default values for required environment variables:"
        for var in "${missing_vars[@]}"; do
            echo "  - $var"
        done
        print_info "Please update $ENV_FILE with production values."
        exit 1
    fi

    print_success "Environment configuration is valid"
}

check_ssl_certificates() {
    print_info "Checking SSL certificates..."

    if [ ! -d "$SSL_DIR" ]; then
        print_warning "SSL directory not found. Creating it..."
        mkdir -p "$SSL_DIR"
    fi

    if [ ! -f "$SSL_DIR/fullchain.pem" ] || [ ! -f "$SSL_DIR/privkey.pem" ]; then
        print_warning "SSL certificates not found."
        read -p "Do you want to generate a self-signed certificate for testing? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            generate_self_signed_cert
        else
            print_error "SSL certificates are required. Please run certbot or place certificates in $SSL_DIR/"
            exit 1
        fi
    else
        print_success "SSL certificates found"
    fi
}

generate_self_signed_cert() {
    print_info "Generating self-signed certificate..."

    read -p "Enter your domain name (e.g., example.com): " domain

    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout "$SSL_DIR/privkey.pem" \
        -out "$SSL_DIR/fullchain.pem" \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=$domain"

    print_success "Self-signed certificate generated (valid for 365 days)"
    print_warning "Self-signed certificates should NOT be used in production!"
}

update_nginx_config() {
    print_info "Checking nginx configuration..."

    if grep -q "yourdomain.com" nginx/nginx.conf; then
        print_warning "nginx.conf still contains placeholder domain 'yourdomain.com'"
        read -p "Enter your actual domain name (e.g., example.com): " domain

        if [ -n "$domain" ]; then
            # Backup original
            cp nginx/nginx.conf nginx/nginx.conf.bak

            # Replace domain
            sed -i.tmp "s/yourdomain\.com/$domain/g" nginx/nginx.conf
            rm nginx/nginx.conf.tmp 2>/dev/null || true

            print_success "Updated nginx.conf with domain: $domain"
        fi
    else
        print_success "nginx configuration looks good"
    fi
}

build_images() {
    print_info "Building Docker images..."

    docker compose -f "$COMPOSE_FILE" build --no-cache

    print_success "Docker images built successfully"
}

start_services() {
    print_info "Starting services..."

    # Load environment variables
    export $(grep -v '^#' "$ENV_FILE" | xargs)

    docker compose -f "$COMPOSE_FILE" up -d

    print_success "Services started"
}

wait_for_healthy() {
    print_info "Waiting for services to become healthy..."

    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        local unhealthy=$(docker compose -f "$COMPOSE_FILE" ps --format json | jq -r '.[] | select(.Health != "healthy") | .Service' 2>/dev/null || echo "")

        if [ -z "$unhealthy" ]; then
            print_success "All services are healthy"
            return 0
        fi

        echo -n "."
        sleep 5
        attempt=$((attempt + 1))
    done

    echo ""
    print_warning "Some services are not healthy yet. Check with: docker compose -f $COMPOSE_FILE ps"
}

show_status() {
    print_info "Service status:"
    docker compose -f "$COMPOSE_FILE" ps

    echo ""
    print_info "To view logs: docker compose -f $COMPOSE_FILE logs -f"
    print_info "To stop services: docker compose -f $COMPOSE_FILE down"
}

# Main deployment flow
main() {
    echo "=========================================="
    echo "  The Nom Database - Production Deploy"
    echo "=========================================="
    echo ""

    check_docker
    check_env_file
    check_ssl_certificates
    update_nginx_config

    echo ""
    read -p "Ready to deploy? This will rebuild and restart all services. (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Deployment cancelled."
        exit 0
    fi

    build_images
    start_services
    wait_for_healthy

    echo ""
    print_success "Deployment completed!"
    echo ""
    show_status

    echo ""
    print_info "Next steps:"
    echo "  1. Test your application at https://your-domain.com"
    echo "  2. Create an admin user (see DEPLOYMENT.md)"
    echo "  3. Set up automated backups: crontab -e"
    echo "  4. Monitor logs: docker compose -f $COMPOSE_FILE logs -f"
}

# Handle script arguments
case "${1:-}" in
    --check-only)
        check_docker
        check_env_file
        check_ssl_certificates
        print_success "All checks passed!"
        ;;
    --build-only)
        check_env_file
        build_images
        ;;
    --start-only)
        check_env_file
        start_services
        ;;
    *)
        main
        ;;
esac
