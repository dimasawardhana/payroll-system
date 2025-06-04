package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	direction := "up"
	if len(os.Args) > 1 && (os.Args[1] == "down" || os.Args[1] == "up") {
		direction = os.Args[1]
	}
	_ = godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	migrationsDir := "./database/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && ((direction == "up" && filepath.Ext(file.Name()) == ".sql" && filepath.Ext(file.Name()[:len(file.Name())-4]) == ".up") ||
			(direction == "down" && filepath.Ext(file.Name()) == ".sql" && filepath.Ext(file.Name()[:len(file.Name())-4]) == ".down")) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	// Sort files to apply in order (reverse for down)
	sort.Strings(migrationFiles)
	if direction == "down" {
		for i, j := 0, len(migrationFiles)-1; i < j; i, j = i+1, j-1 {
			migrationFiles[i], migrationFiles[j] = migrationFiles[j], migrationFiles[i]
		}
	}

	fmt.Printf("Applying %s migrations from %d \n ", direction, len(migrationFiles))

	for _, fname := range migrationFiles {
		path := filepath.Join(migrationsDir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", fname, err)
		}
		fmt.Printf("Applying migration (%s): %s\n", direction, fname)
		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", fname, err)
		}
	}
	fmt.Printf("All %s migrations applied successfully.\n", direction)
}
