package model

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	DueDate   *time.Time `json:"due_date"`
	CreatedAt time.Time  `json:"created_at"`
}

type TodoCreateRequest struct {
	Title   string  `json:"title"`
	DueDate *string `json:"due_date"`
}

type TodoUpdateRequest struct {
	Title   string  `json:"title"`
	Status  string  `json:"status"`
	DueDate *string `json:"due_date"`
}

type TodoListQuery struct {
	Status *string
	Limit  int
	Offset int
}
