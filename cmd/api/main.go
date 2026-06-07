package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/qm3llz/tasksWebApi/internal/db"
	"github.com/qm3llz/tasksWebApi/internal/repository"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	conn := db.ConnectDB(ctx)
	defer conn.Close(ctx)
	taskRepo := repository.NewTaskRepository(conn)

	router := chi.NewRouter()

	http.ListenAndServe(":8080", router)
}
 