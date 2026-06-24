package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/qm3llz/tasksWebApi/internal/models"
	"github.com/qm3llz/tasksWebApi/internal/repository"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

// NewTaskHandler
func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// Create
func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	json.NewDecoder(r.Body).Decode(&task)
	err := t.repo.Create(r.Context(), task)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
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

func (t *TaskHandler) GetById(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)

	task, err := t.repo.GetByID(r.Context(), task.ID)
	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
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
	var body struct {
		UserID uuid.UUID `json:"user_id"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	tasks, err := t.repo.GetAllByUser(r.Context(), body.UserID)
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
	var task models.Task

	json.NewDecoder(r.Body).Decode(&task)

	err := t.repo.Delete(r.Context(), task.ID)
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
	var task models.Task

	json.NewDecoder(r.Body).Decode(&task)

	err := t.repo.Update(r.Context(), task)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	newTask, err := t.repo.GetByID(r.Context(), task.ID)
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
