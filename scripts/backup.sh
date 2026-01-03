#!/bin/bash

set -e

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
BACKUP_DIR="backups"
RETENTION_DAYS=30
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="backup_${DATE}.sql"

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

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

print_info "Starting database backup..."

# Check if database container is running
if ! docker compose -f "$COMPOSE_FILE" ps db | grep -q "Up"; then
    print_error "Database container is not running"
    exit 1
fi

# Create database backup
print_info "Creating PostgreSQL dump..."
docker compose -f "$COMPOSE_FILE" exec -T db pg_dump -U nomdb nomdb > "$BACKUP_DIR/$BACKUP_FILE"

if [ $? -eq 0 ]; then
    print_success "Database dumped to $BACKUP_DIR/$BACKUP_FILE"
else
    print_error "Database backup failed"
    exit 1
fi

# Compress backup
print_info "Compressing backup..."
gzip "$BACKUP_DIR/$BACKUP_FILE"
BACKUP_FILE="${BACKUP_FILE}.gz"

print_success "Backup compressed: $BACKUP_DIR/$BACKUP_FILE"

# Get backup size
BACKUP_SIZE=$(du -h "$BACKUP_DIR/$BACKUP_FILE" | cut -f1)
print_info "Backup size: $BACKUP_SIZE"

# Clean up old backups
print_info "Removing backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "backup_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete

REMAINING_BACKUPS=$(ls -1 "$BACKUP_DIR"/backup_*.sql.gz 2>/dev/null | wc -l)
print_info "Remaining backups: $REMAINING_BACKUPS"

# Optional: Upload to S3 (if AWS CLI is configured)
if command -v aws &> /dev/null && [ -n "${S3_BACKUP_BUCKET:-}" ]; then
    print_info "Uploading backup to S3..."
    aws s3 cp "$BACKUP_DIR/$BACKUP_FILE" "s3://$S3_BACKUP_BUCKET/nomdb-backups/$BACKUP_FILE"

    if [ $? -eq 0 ]; then
        print_success "Backup uploaded to S3: s3://$S3_BACKUP_BUCKET/nomdb-backups/$BACKUP_FILE"
    else
        print_error "S3 upload failed (backup is still saved locally)"
    fi
fi

print_success "Backup completed successfully!"

# Show latest backups
print_info "Latest backups:"
ls -lht "$BACKUP_DIR"/backup_*.sql.gz | head -5
