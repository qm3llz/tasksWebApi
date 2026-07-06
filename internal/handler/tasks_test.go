package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/qm3llz/tasksWebApi/internal/models"
)

type fakeRepo struct {
	task models.Task
	err  error
}

func (f *fakeRepo) Create(ctx context.Context, task models.Task) error {
	return f.err
}

func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (models.Task, error) {
	return f.task, f.err
}

func (f *fakeRepo) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Task, error) {
	tasks := []models.Task{f.task}
	return tasks, f.err
}

func (f *fakeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return f.err
}

func (f *fakeRepo) Update(ctx context.Context, task models.Task) error {
	return f.err
}

func TestCreate(t *testing.T) {
	test := []struct {
		name       string
		repoErr    error
		wantStatus int
	}{
		{"succes", nil, http.StatusCreated},
		{"Interna err", errors.New("internal server error"), http.StatusInternalServerError},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{task: models.Task{Name: "task"}, err: tc.repoErr}
			h := NewTaskHandler(repo)

			task := models.Task{
				Name:        "task",
				Status:      "",
				Description: nil,
			}

			data, err := json.Marshal(task)
			if err != nil {
				panic("err")
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(data))

			h.Create(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("\"%s\": Ожидался статус %d, получен %d", tc.name, tc.wantStatus, rec.Code)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	tests := []struct {
		name       string
		repoErr    error
		wantStatus int
	}{
		{"success", nil, http.StatusOK},
		{"not found", errors.New("Not found"), http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{task: models.Task{Name: tc.name}, err: tc.repoErr}
			h := NewTaskHandler(repo)

			req := httptest.NewRequest(http.MethodGet, "/tasks/11111111-1111-4111-1111-111111111111", nil)

			r := chi.NewRouter()
			r.Get("/tasks/{id}", h.GetById)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("\"%s\"wait %d, get %d", tc.name, tc.wantStatus, rec.Code)
			}
		})
	}
}

func TestGetAllByUser(t *testing.T) {
	test := []struct {
		name       string
		repoErr    error
		wantStatus int
	}{
		{name: "succes", repoErr: nil, wantStatus: http.StatusOK},
		{name: "", repoErr: errors.New("BadRequest"), wantStatus: http.StatusBadRequest},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{task: models.Task{Name: tc.name}, err: tc.repoErr}
			h := NewTaskHandler(repo)

			body := strings.NewReader(`{"id":"11111111-1111-4111-1111-111111111111"}`)
			req := httptest.NewRequest(http.MethodGet, "/tasks/1", body)

			rec := httptest.NewRecorder()

			h.GetAllByUser(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("\"%s\": wait: %d, get: %d", tc.name, tc.wantStatus, rec.Code)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	test := []struct {
		name       string
		repoErr    error
		wantStatus int
	}{
		{name: "succes", repoErr: nil, wantStatus: 200},
		{name: "Server error", repoErr: errors.New("Status Internal Server Error"), wantStatus: 500},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{task: models.Task{Name: tc.name}, err: tc.repoErr}
			h := NewTaskHandler(repo)

			req := httptest.NewRequest(http.MethodDelete, "/tasks/11111111-1111-4111-1111-111111111111", nil)

			r := chi.NewRouter()
			r.Delete("/tasks/{id}", h.Delete)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("\"%s\": wait: %d, get: %d", tc.name, tc.wantStatus, rec.Code)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	test := []struct {
		name       string
		repoErr    error
		wantStatus int
	}{
		{name: "succes", repoErr: nil, wantStatus: 200},
		{name: "BadRequest", repoErr: errors.New("BadRequest"), wantStatus: 400}, // TODO: update mock for test
		// {name: "ID not found", repoErr: errors.New("ID not found"), wantStatus: 404},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{task: models.Task{Name: tc.name}, err: tc.repoErr}
			h := NewTaskHandler(repo)

			body := strings.NewReader(`{"id":"11111111-1111-4111-1111-111111111111"}`)
			req := httptest.NewRequest(http.MethodPut, "/tasks/1", body)

			rec := httptest.NewRecorder()

			h.Update(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("\"%s\": wait: %d, get: %d", tc.name, tc.wantStatus, rec.Code)
			}
		})
	}
}
