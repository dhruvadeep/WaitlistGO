package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) {
	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
    		_id SERIAL PRIMARY KEY,  -- Auto-incrementing primary key
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Defaults to the current timestamp
			ipaddress TEXT,
			ref_by TEXT
		);
		`,

		// Referrals table
		`CREATE TABLE IF NOT EXISTS referrals (
			_id SERIAL PRIMARY KEY,  -- Auto-incrementing primary key
			ref_email VARCHAR(255) NOT NULL,  -- No foreign key reference
			ref_to_email VARCHAR(255) NOT NULL,  -- No foreign key reference
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Defaults to the current timestamp
			ipaddress TEXT
		);`,
	}

	// Execute each query
	for _, query := range queries {
		_, err := pool.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to run migration query: %v\nQuery: %s", err, query)
		}
	}

	log.Println("Migrations executed successfully.")
}