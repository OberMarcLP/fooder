# Production Scripts

This directory contains helper scripts for deploying and managing The Nom Database in production.

## Scripts Overview

### deploy.sh

Automated deployment script that handles the complete deployment process.

**Usage:**

```bash
# Full deployment (interactive)
./scripts/deploy.sh

# Check configuration only
./scripts/deploy.sh --check-only

# Build images only
./scripts/deploy.sh --build-only

# Start services only (skip build)
./scripts/deploy.sh --start-only
```

**What it does:**
- ✅ Validates Docker installation
- ✅ Checks environment configuration (.env.production)
- ✅ Verifies SSL certificates exist
- ✅ Offers to generate self-signed cert if missing
- ✅ Updates nginx.conf with your domain
- ✅ Builds Docker images
- ✅ Starts all services
- ✅ Waits for health checks
- ✅ Shows service status

**Example:**

```bash
$ ./scripts/deploy.sh

==========================================
  The Nom Database - Production Deploy
==========================================

[INFO] Checking Docker installation...
[SUCCESS] Docker and Docker Compose are installed
[INFO] Checking environment configuration...
[SUCCESS] Environment configuration is valid
[INFO] Checking SSL certificates...
[SUCCESS] SSL certificates found
[INFO] Checking nginx configuration...
[SUCCESS] nginx configuration looks good

Ready to deploy? This will rebuild and restart all services. (y/n) y

[INFO] Building Docker images...
[SUCCESS] Docker images built successfully
[INFO] Starting services...
[SUCCESS] Services started
[INFO] Waiting for services to become healthy...
[SUCCESS] All services are healthy

[SUCCESS] Deployment completed!
```

### backup.sh

Creates automated backups of the PostgreSQL database.

**Usage:**

```bash
./scripts/backup.sh
```

**What it does:**
- ✅ Creates PostgreSQL dump
- ✅ Compresses backup with gzip
- ✅ Stores in `backups/` directory
- ✅ Removes backups older than 30 days
- ✅ Optionally uploads to S3 (if configured)
- ✅ Shows backup size and count

**Output:**

```
backups/
├── backup_20250130_020000.sql.gz
├── backup_20250129_020000.sql.gz
├── backup_20250128_020000.sql.gz
└── ...
```

**Automated Backups:**

Add to crontab for daily backups at 2 AM:

```bash
crontab -e
```

Add this line:

```cron
0 2 * * * cd /opt/nomdb && ./scripts/backup.sh >> /var/log/nomdb-backup.log 2>&1
```

**S3 Upload (Optional):**

Set environment variable to enable S3 uploads:

```bash
export S3_BACKUP_BUCKET=my-backup-bucket
```

### restore.sh

Restores the database from a backup file.

**Usage:**

```bash
# List available backups
./scripts/restore.sh

# Restore from specific backup
./scripts/restore.sh backups/backup_20250130_020000.sql.gz

# Or just provide filename
./scripts/restore.sh backup_20250130_020000.sql.gz
```

**What it does:**
- ⚠️ Stops dependent services (backend, frontend, nginx)
- ⚠️ Terminates all database connections
- ⚠️ Drops and recreates database
- ✅ Restores from backup file
- ✅ Restarts all services
- ✅ Waits for services to become healthy

**⚠️ WARNING:** This is a destructive operation! Always verify you're restoring the correct backup.

**Example:**

```bash
$ ./scripts/restore.sh backup_20250130_020000.sql.gz

==========================================
  DATABASE RESTORE - DESTRUCTIVE ACTION
==========================================

This will REPLACE the current database with the backup from:
  backups/backup_20250130_020000.sql.gz

Are you absolutely sure you want to continue? Type 'yes' to proceed: yes

[INFO] Stopping dependent services...
[INFO] Decompressing backup...
[INFO] Terminating existing database connections...
[INFO] Recreating database...
[INFO] Restoring database from backup...
[SUCCESS] Database restored successfully!
[INFO] Starting services...
[SUCCESS] Restore completed! Services are restarting...
```

## Setup Instructions

### Make Scripts Executable

```bash
chmod +x scripts/*.sh
```

### Environment Variables

Scripts use `.env.production` for configuration. Ensure it exists:

```bash
cp .env.production.example .env.production
nano .env.production
```

## Common Workflows

### Initial Deployment

```bash
# 1. Configure environment
cp .env.production.example .env.production
nano .env.production

# 2. Deploy
./scripts/deploy.sh
```

### Daily Operations

```bash
# View logs
docker compose -f docker-compose.prod.yml logs -f

# Restart service
docker compose -f docker-compose.prod.yml restart backend

# Check status
docker compose -f docker-compose.prod.yml ps
```

### Backup & Restore

```bash
# Manual backup
./scripts/backup.sh

# List backups
ls -lht backups/

# Restore
./scripts/restore.sh backups/backup_20250130_020000.sql.gz
```

### Updates

```bash
# Pull latest code
git pull origin main

# Rebuild and deploy
./scripts/deploy.sh
```

## Troubleshooting

### "Permission denied" Error

Make scripts executable:

```bash
chmod +x scripts/*.sh
```

### "Docker not found" Error

Install Docker:

```bash
sudo apt update
sudo apt install -y docker.io docker-compose-plugin
sudo usermod -aG docker $USER
newgrp docker
```

### "Database container not running" Error

Start the database:

```bash
docker compose -f docker-compose.prod.yml up -d db
```

### Backup Fails

Check disk space:

```bash
df -h
```

Clean up old Docker images:

```bash
docker system prune -a
```

## Best Practices

### Backups

- ✅ Run daily automated backups (cron job)
- ✅ Store backups off-site (S3, external server)
- ✅ Test restore procedure monthly
- ✅ Keep at least 30 days of backups
- ✅ Monitor backup success/failure

### Deployment

- ✅ Always test in staging first
- ✅ Review changes before deploying (`git diff`)
- ✅ Back up database before major updates
- ✅ Monitor logs during deployment
- ✅ Have rollback plan ready

### Security

- ✅ Never commit .env.production to git
- ✅ Restrict script permissions (chmod 700)
- ✅ Use strong database passwords
- ✅ Rotate secrets regularly
- ✅ Keep scripts and Docker updated

## Advanced Usage

### Custom Backup Retention

Edit `backup.sh` and change `RETENTION_DAYS`:

```bash
RETENTION_DAYS=60  # Keep 60 days instead of 30
```

### Backup to Multiple Locations

```bash
# Backup locally and to S3
export S3_BACKUP_BUCKET=my-backup-bucket
./scripts/backup.sh

# Also copy to remote server
scp backups/backup_*.sql.gz user@backup-server:/backups/
```

### Pre-deployment Checks

```bash
# Validate configuration
./scripts/deploy.sh --check-only

# Test build
./scripts/deploy.sh --build-only

# Check what will change
git diff origin/main
```

## Exit Codes

All scripts use standard exit codes:

- `0` - Success
- `1` - Error occurred

Check exit code in scripts:

```bash
./scripts/backup.sh
if [ $? -eq 0 ]; then
    echo "Backup successful"
else
    echo "Backup failed"
fi
```

## Logging

### View Script Logs

```bash
# Real-time logs
docker compose -f docker-compose.prod.yml logs -f

# Backup logs (if using cron)
tail -f /var/log/nomdb-backup.log
```

### Enable Debug Mode

Add `-x` to script shebang for debugging:

```bash
#!/bin/bash -x
```

## Support

For issues or questions:
- Check [DEPLOYMENT.md](../DEPLOYMENT.md)
- Review [PRODUCTION_CHECKLIST.md](../PRODUCTION_CHECKLIST.md)
- Open GitHub issue
