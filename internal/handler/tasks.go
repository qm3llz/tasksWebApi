package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/qm3llz/tasksWebApi/internal/models"
)

type TaskRepo interface {
	Create(ctx context.Context, task models.Task) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (models.Task, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Task, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, task models.Task, id uuid.UUID) error
}

type TaskHandler struct {
	repo TaskRepo
}

// NewTaskHandler
func NewTaskHandler(repo TaskRepo) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// Create
func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	json.NewDecoder(r.Body).Decode(&task)

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	task.UserID = userID

	err := t.repo.Create(r.Context(), task)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "task created",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (t *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	taskID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid  task ID format", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	
	task, err := t.repo.GetByID(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response := map[string]any{
		"status": "success",
		"task":   task,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAllByUser
func (t *TaskHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := t.repo.GetAllByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"status":  "success",
		"message": "getted",
		"tasks":   tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (t *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	taskID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid  task ID format", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = t.repo.Delete(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "task deleted (or the task does not exist)",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	StrID := chi.URLParam(r, "id")

	ID, err := uuid.Parse(StrID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}
 
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)
	task.UserID = userID

	err = t.repo.Update(r.Context(), task, ID)
	if err != nil {
		http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
		return
	}

	newTask, err := t.repo.GetByID(r.Context(), task.ID, task.UserID)
	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	response := map[string]any{
		"status":  "success",
		"message": "task updated",
		"task":    newTask,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
