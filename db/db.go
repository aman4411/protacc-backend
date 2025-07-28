package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development" // Default to development
	}
	if dbURL == "" {
		log.Fatalf("DATABASE_URL not set in %s environment", env)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	// Configure pool settings to avoid prepared statement conflicts
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute * 5

	// Disable prepared statements to avoid conflicts
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Printf("✅ Connected to %s PostgreSQL with connection pool (max: %d, min: %d)", env, config.MaxConns, config.MinConns)
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
		log.Printf("🔌 Disconnected from PostgreSQL")
	}
}
