package repository

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FetchTasks(ctx context.Context) ([]*model.Todo, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, status, created_at FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("fetch tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Todo
	for rows.Next() {
		var task model.Todo
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Status,
			&task.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterate tasks: %w", err)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *Repository) FetchTask(ctx context.Context, taskID uuid.UUID) (*model.Todo, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, title, status, created_at FROM tasks WHERE id = $1", taskID)
	var task model.Todo

	err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrTaskNotFound
		}
		return nil, fmt.Errorf("scan task: %w", err)
	}

	return &task, nil
}

func (r *Repository) CreateTask(ctx context.Context, task model.TodoCreateRequest) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO tasks(title) values($1)", task.Title)
	if err != nil {
		return fmt.Errorf("error occurred during task creation")
	}
	return nil
}
