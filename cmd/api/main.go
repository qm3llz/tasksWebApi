package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/qm3llz/tasksWebApi/internal/db"
	"github.com/qm3llz/tasksWebApi/internal/handler"
	mw "github.com/qm3llz/tasksWebApi/internal/middleware"
	"github.com/qm3llz/tasksWebApi/internal/repository"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	conn := db.ConnectDB(ctx)
	defer conn.Close(ctx)
	taskRepo := repository.NewTaskRepository(conn)
	userRepo := repository.NewUserRepository(conn)

	th := handler.NewTaskHandler(taskRepo)
	uh := handler.NewUserHandler(userRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/register", uh.Register)
	r.Post("/login", uh.Login)

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)

		r.Post("/tasks", th.Create)
		r.Get("/tasks/{id}", th.GetByID)
		r.Get("/tasks", th.GetAllByUser)
		r.Delete("/tasks/{id}", th.Delete)
		r.Put("/tasks/{id}", th.Update)
	})

	http.ListenAndServe(":8080", r)
}
