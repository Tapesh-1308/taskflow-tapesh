package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dbUrl string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("DB not reachable: %v", err)
	}

	log.Println("Connected to DB")

	return pool
}
