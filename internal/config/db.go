package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB initializes and returns a pgxpool.Pool for PostgreSQL
func InitDB(databaseURL string) *pgxpool.Pool {
	ctx := context.Background()
	pgxUrl := os.Getenv("DATABASE_URL")
	if pgxUrl == "" {
		pgxUrl = databaseURL
	}
	pool, err := pgxpool.New(ctx, pgxUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return pool
}
