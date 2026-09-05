package service

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"

	"github.com/google/uuid"
)

type TodoRepo interface {
	FetchTasks(ctx context.Context) ([]*model.Todo, error)
	FetchTask(ctx context.Context, taskID uuid.UUID) (*model.Todo, error)
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
