package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nomdb/backend/internal/database"
	"github.com/nomdb/backend/internal/logger"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Fatal("DATABASE_URL environment variable is required")
	}

	migrationsPath := "./db/migrations_new"
	command := os.Args[1]

	switch command {
	case "up":
		if err := database.RunMigrations(databaseURL, migrationsPath); err != nil {
			logger.Fatal("Migration failed: %v", err)
		}

	case "down":
		if err := database.MigrateDown(databaseURL, migrationsPath); err != nil {
			logger.Fatal("Rollback failed: %v", err)
		}

	case "version":
		version, dirty, err := database.MigrateVersion(databaseURL, migrationsPath)
		if err != nil {
			logger.Fatal("Failed to get version: %v", err)
		}
		dirtyStatus := ""
		if dirty {
			dirtyStatus = " (dirty)"
		}
		logger.Info("Current migration version: %d%s", version, dirtyStatus)

	case "force":
		if len(os.Args) < 3 {
			logger.Fatal("Version number required for force command")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.Fatal("Invalid version number: %v", err)
		}
		if err := database.MigrateForce(databaseURL, migrationsPath, version); err != nil {
			logger.Fatal("Force migration failed: %v", err)
		}

	case "create":
		if len(os.Args) < 3 {
			logger.Fatal("Migration name required for create command")
		}
		name := os.Args[2]
		createMigration(migrationsPath, name)

	default:
		logger.Fatal("Unknown command: %s", command)
		printUsage()
	}
}

func createMigration(path, name string) {
	timestamp := time.Now().Unix()
	upFile := fmt.Sprintf("%s/%d_%s.up.sql", path, timestamp, name)
	downFile := fmt.Sprintf("%s/%d_%s.down.sql", path, timestamp, name)

	// Create up migration
	if err := os.WriteFile(upFile, []byte("-- Write your migration here\n"), 0644); err != nil {
		logger.Fatal("Failed to create up migration: %v", err)
	}

	// Create down migration
	if err := os.WriteFile(downFile, []byte("-- Write your rollback here\n"), 0644); err != nil {
		logger.Fatal("Failed to create down migration: %v", err)
	}

	logger.Info("✓ Created migration files:")
	logger.Info("  %s", upFile)
	logger.Info("  %s", downFile)
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  go run cmd/migrate/main.go <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  up              Run all pending migrations")
	fmt.Println("  down            Rollback the last migration")
	fmt.Println("  version         Show current migration version")
	fmt.Println("  force <version> Force set migration version (use with caution)")
	fmt.Println("  create <name>   Create a new migration file")
	fmt.Println("\nExamples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go down")
	fmt.Println("  go run cmd/migrate/main.go create add_user_table")
	fmt.Println("  go run cmd/migrate/main.go force 3")
}
