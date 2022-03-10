package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect() *pgxpool.Pool {
	connStr := "postgres://postgres:postgres@localhost:5432/es_go_local?sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	return pool
}
