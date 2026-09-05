package service

import (
	"F5/internals/model"
	"context"
)

type TodoRepo interface {
	FetchTasks(ctx context.Context) ([]*model.Todo, error)
}

func (s *Services) GetTasks(ctx context.Context) ([]*model.Todo, error) {
	return s.repo.FetchTasks(ctx)
}
