package service

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"strings"

	"github.com/google/uuid"
)

type TodoRepo interface {
	FetchTasks(ctx context.Context) ([]*model.Todo, error)
	FetchTask(ctx context.Context, taskID uuid.UUID) (*model.Todo, error)
	CreateTask(ctx context.Context, task model.TodoCreateRequest) error
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	UpdateTask(ctx context.Context, taskID uuid.UUID, taskReq model.TodoUpdateRequest) error
}

func (s *Services) GetTasks(ctx context.Context) ([]*model.Todo, error) {
	return s.repo.FetchTasks(ctx)
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

	taskReq.Title = title

	err := s.repo.CreateTask(ctx, taskReq)
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

	if title == "" {
		return apperrors.ErrTitleEmpty
	}

	if len(title) > 200 {
		return apperrors.ErrTitleTooLong
	}

	taskReq.Title = title

	if strings.ToLower(taskReq.Status) != "pending" && strings.ToLower(taskReq.Status) != "done" {
		return apperrors.ErrInvalidStatus
	}

	err = s.repo.UpdateTask(ctx, taskID, taskReq)
	if err != nil {
		return err
	}

	return nil
}
