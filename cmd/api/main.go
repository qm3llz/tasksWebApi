package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/qm3llz/tasksWebApi/internal/db"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	conn := db.ConnectDB(ctx)
	defer conn.Close(ctx)
}
