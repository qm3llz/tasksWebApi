// Package repository provides data access for tasks and users.
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/qm3llz/tasksWebApi/internal/models"
)

type TaskRepository struct {
	conn *pgx.Conn
}

func NewTaskRepository(conn *pgx.Conn) *TaskRepository {
	return &TaskRepository{conn: conn}
}

func (r *TaskRepository) Create(ctx context.Context, task models.Task) error {
	sql := `
	INSERT INTO tasks (user_id, name, status, description)
	VALUES ($1, $2, $3, $4);
	`
	_, err := r.conn.Exec(ctx, sql, task.UserID, task.Name, task.Status, task.Description)
	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Task, error) {
	var t models.Task

	sql := `
	SELECT id, user_id, name, status, description, created_at, updated_at FROM tasks
	WHERE id = $1;
	`
	err := r.conn.QueryRow(ctx, sql, id).Scan(&t.ID, &t.UserID, &t.Name, &t.Status, &t.Description, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

// GetAllByUser returns all tasks owned by the user with the given id.
// It returns an empty slice if the user has no tasks.
func (r *TaskRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Task, error) {
	var t []models.Task

	sql := `
	SELECT id, user_id, name, status, description, created_at, updated_at FROM tasks
	WHERE user_id = $1;
	`
	rows, err := r.conn.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.Status, &task.Description, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		t = append(t, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return t, err
}
