package database

import (
	"context"
	"log"
	"time"
	"waitlist-golang/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBManager struct {
	pool *pgxpool.Pool
}

var DBPool *pgxpool.Pool

func NewDBManager(connUrl string) (*DBManager, error) {
	config, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, err
	}

	config.MinConns = int32(configs.GetMinConns())
	config.MaxConns = int32(configs.GetMaxConns())
	config.HealthCheckPeriod = configs.GetMaxConnLifetime()

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	DBPool = pool // Save the pool to a global variable
	return &DBManager{pool: pool}, nil
}

func (d *DBManager) GetPool() *pgxpool.Pool {
	return d.pool
}

func (db *DBManager) RefreshConnection(connURL string) {
	ticker := time.NewTicker(configs.GetMaxConnLifetime())
	go func() {
		for range ticker.C {
			log.Println("Refreshing database connection pool...")
			newPool, err := pgxpool.New(context.Background(), connURL)
			if err != nil {
				log.Printf("Failed to refresh connection pool: %v\n", err)
				continue
			}
			oldPool := db.pool
			db.pool = newPool

			// set the global pool to the new pool
			DBPool = newPool
			oldPool.Close()
		}
	}()
}

func (db *DBManager) Close() {
	db.pool.Close()
}