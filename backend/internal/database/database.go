package database

import (
	"context"
	"fmt"
	"os"

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

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
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
