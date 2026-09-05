package service

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TodoRepo interface {
	FetchTasks(ctx context.Context, q model.TodoListQuery) ([]*model.Todo, error)
	FetchTask(ctx context.Context, taskID uuid.UUID) (*model.Todo, error)
	CreateTask(ctx context.Context, task model.TodoCreateRequest, dueDate *time.Time) error
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	UpdateTask(ctx context.Context, taskID uuid.UUID, taskReq model.TodoUpdateRequest, dueDate *time.Time) error
}

func (s *Services) GetTasks(ctx context.Context, q model.TodoListQuery) ([]*model.Todo, error) {
	if q.Status != nil {
		status := strings.ToLower(*q.Status)
		if status != "pending" && status != "done" {
			return nil, apperrors.ErrInvalidStatus
		}
		q.Status = &status
	}
	return s.repo.FetchTasks(ctx, q)
}

func (s *Services) GetTask(ctx context.Context, taskIDStr string) (*model.Todo, error) {
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidID
	}
	return s.repo.FetchTask(ctx, taskID)
}

func (s *Services) CreateTask(ctx context.Context, taskReq model.TodoCreateRequest) error {
	title := strings.TrimSpace(taskReq.Title)

	if title == "" {
		return apperrors.ErrTitleEmpty
	}

	if len(title) > 200 {
		return apperrors.ErrTitleTooLong
	}
	dueDate, err := parseDate(taskReq.DueDate)
	if err != nil {
		return err
	}

	taskReq.Title = title
	err = s.repo.CreateTask(ctx, taskReq, dueDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) DeleteTask(ctx context.Context, taskIDStr string) error {
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return apperrors.ErrInvalidID
	}
	return s.repo.DeleteTask(ctx, taskID)
}

func (s *Services) UpdateTask(ctx context.Context, taskIDStr string, taskReq model.TodoUpdateRequest) error {
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return apperrors.ErrInvalidID
	}

	title := strings.TrimSpace(taskReq.Title)
	status := strings.ToLower(taskReq.Status)

	if title == "" {
		return apperrors.ErrTitleEmpty
	}

	if len(title) > 200 {
		return apperrors.ErrTitleTooLong
	}

	if status != "pending" && status != "done" {
		return apperrors.ErrInvalidStatus
	}

	dueDate, err := parseDate(taskReq.DueDate)
	if err != nil {
		return err
	}

	taskReq.Title = title
	taskReq.Status = status

	err = s.repo.UpdateTask(ctx, taskID, taskReq, dueDate)
	if err != nil {
		return err
	}

	return nil
}
