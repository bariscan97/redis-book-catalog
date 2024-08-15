package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Pool() *pgxpool.Pool {
	connString := fmt.Sprintf("%s://%s:%s@postgres:%s/%s?sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbPool, err := pgxpool.Connect(context.Background(), connString)
	
	if err != nil {
		panic(err)
	}
	return dbPool
}
