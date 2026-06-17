package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/qm3llz/tasksWebApi/internal/db"
	"github.com/qm3llz/tasksWebApi/internal/handler"
	"github.com/qm3llz/tasksWebApi/internal/repository"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	conn := db.ConnectDB(ctx)
	defer conn.Close(ctx)
	taskRepo := repository.NewTaskRepository(conn)
	h := handler.NewTaskHandler(taskRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/tasks", h.Create)
	r.Get("/tasks", h.GetAllByUser)
	r.Get("/tasks/{id}", h.GetById)
	r.Put("/tasks/{id}", h.Update)
	r.Delete("/tasks/{id}", h.Delete)

	http.ListenAndServe(":8080", r)
}
