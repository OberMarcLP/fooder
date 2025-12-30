package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomdb/backend/internal/logger"
)

var pool *pgxpool.Pool

func Connect() error {
	logger.Info("🔌 Connecting to database...")

	// Try to get DATABASE_URL first, otherwise construct from individual components
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Build connection string from individual components
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")

		// Set defaults
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		if user == "" {
			user = "nomdb"
		}
		if password == "" {
			password = "nomdb_secret"
		}
		if dbname == "" {
			dbname = "nomdb"
		}
		if sslmode == "" {
			sslmode = "disable"
		}

		logger.Debug("Building connection string from components: host=%s port=%s user=%s dbname=%s sslmode=%s",
			host, port, user, dbname, sslmode)

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode)
	} else {
		logger.Debug("Using DATABASE_URL from environment")
	}

	// Parse the connection string and configure pool settings
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Error("❌ Failed to parse database URL: %v", err)
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Optimize connection pool settings for performance
	config.MaxConns = 25                              // Maximum number of connections
	config.MinConns = 5                               // Minimum number of idle connections
	config.MaxConnLifetime = time.Hour                // Max connection lifetime (1 hour)
	config.MaxConnIdleTime = 30 * time.Minute         // Max idle time (30 minutes)
	config.HealthCheckPeriod = time.Minute            // Health check every minute

	logger.Debug("Pool configuration: MaxConns=%d, MinConns=%d, MaxConnLifetime=%ds, MaxConnIdleTime=%ds",
		config.MaxConns, config.MinConns, int(config.MaxConnLifetime.Seconds()), int(config.MaxConnIdleTime.Seconds()))

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("❌ Failed to create database connection pool: %v", err)
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	logger.Debug("Testing database connection...")
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("❌ Database ping failed: %v", err)
		return fmt.Errorf("unable to ping database: %w", err)
	}

	logger.Info("✅ Database connected successfully")
	return nil
}

func GetPool() *pgxpool.Pool {
	return pool
}

func Close() {
	if pool != nil {
		logger.Info("🔌 Closing database connection...")
		pool.Close()
		logger.Info("✅ Database connection closed")
	}
}
