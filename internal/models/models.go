package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "Не выполнено"
	StatusInProgress TaskStatus = "Выполняется"
	StatusDone       TaskStatus = "Выполнено"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Task struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	Name        string     `db:"name" json:"name"`
	Status      TaskStatus `db:"status" json:"status"`
	Description *string    `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
