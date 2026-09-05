package repository

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FetchTasks(ctx context.Context, q model.TodoListQuery) ([]*model.Todo, error) {
	var rows pgx.Rows
	var err error

	if q.Status == nil {
		// no filter — get everything
		rows, err = r.pool.Query(ctx,
			"SELECT id, title, status, created_at, due_date FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2",
			q.Limit, q.Offset)
	} else {
		// filter by status
		rows, err = r.pool.Query(ctx,
			"SELECT id, title, status, created_at, due_date FROM tasks WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
			*q.Status, q.Limit, q.Offset)
	}

	if err != nil {
		return nil, fmt.Errorf("fetch tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*model.Todo, 0)
	for rows.Next() {
		var task model.Todo
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Status,
			&task.CreatedAt,
			&task.DueDate,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (r *Repository) FetchTask(ctx context.Context, taskID uuid.UUID) (*model.Todo, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, title, status, created_at, due_date FROM tasks WHERE id = $1", taskID)
	var task model.Todo

	err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.DueDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrTaskNotFound
		}
		return nil, fmt.Errorf("scan task: %w", err)
	}

	return &task, nil
}

func (r *Repository) CreateTask(ctx context.Context, task model.TodoCreateRequest, dueDate *time.Time) error {
	fmt.Println(dueDate)
	_, err := r.pool.Exec(ctx, "INSERT INTO tasks(title, due_date) values($1, $2)", task.Title, dueDate)
	if err != nil {
		return fmt.Errorf("error occurred during task creation: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return fmt.Errorf("error occurred when deleting task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return apperrors.ErrTaskNotFound
	}

	return nil
}

func (r *Repository) UpdateTask(ctx context.Context, taskID uuid.UUID, taskReq model.TodoUpdateRequest, dueDate *time.Time) error {
	title := taskReq.Title
	status := taskReq.Status

	tag, err := r.pool.Exec(ctx,
		`UPDATE tasks 
			 SET title=$1, status=$2, due_date=$3 
			 WHERE id=$4`,
		title, status, dueDate, taskID)
	if err != nil {
		return fmt.Errorf("error occurred on updating task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return apperrors.ErrTaskNotFound
	}

	return nil
}
