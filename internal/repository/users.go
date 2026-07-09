// Package repository provides data access for tasks and users.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qm3llz/tasksWebApi/internal/models"
)

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) Create(ctx context.Context, username string, hash []byte) error {
	sql := `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	`

	_, err := r.conn.Exec(ctx, sql, username, hash)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	sql := `
	SELECT id, username, password, created_at FROM users
	WHERE username = $1
	`

	var t models.User
	err := r.conn.QueryRow(ctx, sql, username).Scan(&t.ID, &t.Username, &t.Password, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
