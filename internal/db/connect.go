package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context) *pgx.Conn {
	dsn := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		panic(err)
	}
	return conn
}
