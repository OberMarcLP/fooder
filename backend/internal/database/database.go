package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable not set")
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

func GetPool() *pgxpool.Pool {
	return pool
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}
