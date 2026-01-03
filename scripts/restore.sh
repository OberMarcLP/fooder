#!/bin/bash

set -e

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
BACKUP_DIR="backups"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# Check if backup file is provided
if [ -z "$1" ]; then
    print_error "Usage: $0 <backup_file>"
    echo ""
    print_info "Available backups:"
    ls -lht "$BACKUP_DIR"/backup_*.sql.gz 2>/dev/null | head -10 || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    # Try looking in backup directory
    if [ -f "$BACKUP_DIR/$BACKUP_FILE" ]; then
        BACKUP_FILE="$BACKUP_DIR/$BACKUP_FILE"
    else
        print_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
fi

print_warning "=========================================="
print_warning "  DATABASE RESTORE - DESTRUCTIVE ACTION"
print_warning "=========================================="
echo ""
print_warning "This will REPLACE the current database with the backup from:"
print_warning "  $BACKUP_FILE"
echo ""
read -p "Are you absolutely sure you want to continue? Type 'yes' to proceed: " confirm

if [ "$confirm" != "yes" ]; then
    print_info "Restore cancelled."
    exit 0
fi

# Check if database container is running
if ! docker compose -f "$COMPOSE_FILE" ps db | grep -q "Up"; then
    print_error "Database container is not running. Start it with: docker compose -f $COMPOSE_FILE up -d db"
    exit 1
fi

print_info "Stopping dependent services..."
docker compose -f "$COMPOSE_FILE" stop backend frontend nginx

# Decompress if needed
RESTORE_FILE="$BACKUP_FILE"
if [[ "$BACKUP_FILE" == *.gz ]]; then
    print_info "Decompressing backup..."
    RESTORE_FILE="${BACKUP_FILE%.gz}"
    gunzip -c "$BACKUP_FILE" > "$RESTORE_FILE"
fi

# Drop existing connections
print_info "Terminating existing database connections..."
docker compose -f "$COMPOSE_FILE" exec -T db psql -U nomdb -d postgres -c \
    "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'nomdb' AND pid <> pg_backend_pid();" || true

# Drop and recreate database
print_info "Recreating database..."
docker compose -f "$COMPOSE_FILE" exec -T db psql -U nomdb -d postgres -c "DROP DATABASE IF EXISTS nomdb;"
docker compose -f "$COMPOSE_FILE" exec -T db psql -U nomdb -d postgres -c "CREATE DATABASE nomdb;"

# Restore from backup
print_info "Restoring database from backup..."
docker compose -f "$COMPOSE_FILE" exec -T db psql -U nomdb nomdb < "$RESTORE_FILE"

if [ $? -eq 0 ]; then
    print_success "Database restored successfully!"
else
    print_error "Database restore failed"
    exit 1
fi

# Clean up decompressed file if we created it
if [[ "$BACKUP_FILE" == *.gz ]] && [ -f "$RESTORE_FILE" ]; then
    rm "$RESTORE_FILE"
fi

# Restart services
print_info "Starting services..."
docker compose -f "$COMPOSE_FILE" up -d

print_success "Restore completed! Services are restarting..."

# Wait for services to be healthy
print_info "Waiting for services to become healthy..."
sleep 10

docker compose -f "$COMPOSE_FILE" ps

print_success "Database restore completed successfully!"
