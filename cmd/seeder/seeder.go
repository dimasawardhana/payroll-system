package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
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

	seederDir := "./database/seeder"
	files, err := os.ReadDir(seederDir)
	if err != nil {
		log.Fatalf("Failed to read seeder directory: %v", err)
	}

	direction := "up"
	if len(os.Args) > 1 && (os.Args[1] == "down" || os.Args[1] == "up") {
		direction = os.Args[1]
	}

	var seederFiles []string
	for _, file := range files {
		if !file.IsDir() && ((direction == "up" && strings.HasSuffix(file.Name(), ".up.sql")) || (direction == "down" && strings.HasSuffix(file.Name(), ".down.sql"))) {
			seederFiles = append(seederFiles, file.Name())
		}
	}
	sort.Strings(seederFiles)
	if direction == "down" {
		for i, j := 0, len(seederFiles)-1; i < j; i, j = i+1, j-1 {
			seederFiles[i], seederFiles[j] = seederFiles[j], seederFiles[i]
		}
	}

	fmt.Printf("Applying %s for %d seeder files...\n", direction, len(seederFiles))
	for _, fname := range seederFiles {
		path := filepath.Join(seederDir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read seeder file %s: %v", fname, err)
		}
		fmt.Printf("Applying seeder (%s): %s\n", direction, fname)
		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("Failed to execute seeder %s: %v", fname, err)
		}
	}
	fmt.Printf("All %s seeders applied successfully.\n", direction)
}
