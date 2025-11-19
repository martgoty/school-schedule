package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(cfg *config.Config) (*DB, error) {
	dbURL := cfg.GetDBURL()

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
        return nil, fmt.Errorf("unable to parse database config: %w", err)
    }

	poolConfig.MaxConns = 20
    poolConfig.MinConns = 5
    poolConfig.MaxConnLifetime = time.Hour
    poolConfig.MaxConnIdleTime = time.Minute * 30

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }

	if err := pool.Ping(ctx); err != nil {
		 return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
    return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}