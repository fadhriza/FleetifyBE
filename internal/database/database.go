package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"fleetify/internal/config"
	"fleetify/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Init
func Connect() error {
	cfg := config.AppConfig.Database

	dsn := cfg.GetDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		errors.LogError("Database Config Parse Error", err)
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	// Config
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30
	poolConfig.HealthCheckPeriod = time.Minute

	DB, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		errors.LogError("Database Connection Pool Creation Error", err)
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		errors.LogError("Database Ping Error", err)
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

// HealthCheck
func HealthCheck(ctx context.Context) error {
	if DB == nil {
		return fmt.Errorf("database connection pool is nil")
	}
	return DB.Ping(ctx)
}
