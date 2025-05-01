package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	env := os.Getenv("ENVIRONMENT")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL not set in %s environment", env)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Printf("✅ Connected to %s PostgreSQL", env)
}

func CloseDB() {
	Pool.Close()
	log.Printf("🔌 Disconnected from PostgreSQL")
}
