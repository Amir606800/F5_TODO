package service

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type MockTodoRepo struct {
	updateTaskCalled bool
	fetchTasksCalled bool
	fetchTaskCalled  bool
	deleteTaskCalled bool
	createTaskCalled bool

	deleteTaskError error
	createTaskError error
}

func (m *MockTodoRepo) UpdateTask(ctx context.Context, taskID uuid.UUID, taskReq model.TodoUpdateRequest, dueDate *time.Time) error {
	m.updateTaskCalled = true
	return nil
}

func (m *MockTodoRepo) FetchTasks(ctx context.Context, q model.TodoListQuery) ([]*model.Todo, error) {
	m.fetchTasksCalled = true
	return nil, nil
}

func (m *MockTodoRepo) FetchTask(ctx context.Context, taskID uuid.UUID) (*model.Todo, error) {
	m.fetchTaskCalled = true
	return nil, nil
}

func (m *MockTodoRepo) CreateTask(ctx context.Context, task model.TodoCreateRequest, dueDate *time.Time) error {
	m.createTaskCalled = true
	return m.createTaskError
}

func (m *MockTodoRepo) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	m.deleteTaskCalled = true
	return m.deleteTaskError
}

func TestUpdateTask_EmptyTitleDone(t *testing.T) {
	mockRepo := &MockTodoRepo{}
	service := &Services{
		repo: mockRepo,
	}

	req := model.TodoUpdateRequest{
		Title:  "",
		Status: "done",
	}

	err := service.UpdateTask(
		context.Background(),
		uuid.New().String(),
		req,
	)
	if !errors.Is(err, apperrors.ErrTitleEmpty) {
		t.Fatalf("expected ErrTitleEmpty, got %v", err)
	}
}

func TestUpdateTask_InvalidID(t *testing.T) {
	mockRepo := &MockTodoRepo{}
	service := &Services{repo: mockRepo}

	req := model.TodoUpdateRequest{
		Title:  "Buy milk",
		Status: "done",
	}

	err := service.UpdateTask(
		context.Background(),
		"InvalidID",
		req,
	)

	if !errors.Is(err, apperrors.ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestUpdateTask_InvalidStatus(t *testing.T) {
	mockRepo := &MockTodoRepo{}
	service := &Services{repo: mockRepo}

	req := model.TodoUpdateRequest{
		Title:  "Buy milk",
		Status: "invalid",
	}

	err := service.UpdateTask(
		context.Background(),
		uuid.New().String(),
		req,
	)

	if !errors.Is(err, apperrors.ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestGetTasks(t *testing.T) {
	pending := "pending"
	done := "done"
	invalid := "something"

	tests := []struct {
		name           string
		query          model.TodoListQuery
		wantErr        error
		repoShouldCall bool
	}{
		{
			name: "valid pending status",
			query: model.TodoListQuery{
				Status: &pending,
				Limit:  10,
				Offset: 0,
			},
			wantErr:        nil,
			repoShouldCall: true,
		},
		{
			name: "valid done status",
			query: model.TodoListQuery{
				Status: &done,
				Limit:  10,
				Offset: 0,
			},
			wantErr:        nil,
			repoShouldCall: true,
		},
		{
			name: "invalid status",
			query: model.TodoListQuery{
				Status: &invalid,
				Limit:  10,
				Offset: 0,
			},
			wantErr:        apperrors.ErrInvalidStatus,
			repoShouldCall: false,
		},
		{
			name: "no status filter",
			query: model.TodoListQuery{
				Status: nil,
				Limit:  10,
				Offset: 0,
			},
			wantErr:        nil,
			repoShouldCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTodoRepo{}

			service := &Services{
				repo: mockRepo,
			}

			_, err := service.GetTasks(
				context.Background(),
				tt.query,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}

			if mockRepo.fetchTasksCalled != tt.repoShouldCall {
				t.Fatalf(
					"expected repo called: %v, got: %v",
					tt.repoShouldCall,
					mockRepo.fetchTasksCalled,
				)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name           string
		taskID         string
		wantErr        error
		repoShouldCall bool
	}{
		{
			name:           "valid ID",
			taskID:         validID,
			wantErr:        nil,
			repoShouldCall: true,
		},
		{
			name:           "invalid ID",
			taskID:         "invalid-uuid",
			wantErr:        apperrors.ErrInvalidID,
			repoShouldCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTodoRepo{}

			service := &Services{
				repo: mockRepo,
			}

			_, err := service.GetTask(
				context.Background(),
				tt.taskID,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}

			if mockRepo.fetchTaskCalled != tt.repoShouldCall {
				t.Fatalf(
					"expected repo called: %v, got: %v",
					tt.repoShouldCall,
					mockRepo.fetchTaskCalled,
				)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name           string
		taskID         string
		wantErr        error
		repoError      error
		repoShouldCall bool
	}{
		{
			name:           "valid ID",
			taskID:         validID,
			wantErr:        nil,
			repoError:      nil,
			repoShouldCall: true,
		},
		{
			name:           "invalid ID",
			taskID:         "invalid-uuid",
			wantErr:        apperrors.ErrInvalidID,
			repoError:      nil,
			repoShouldCall: false,
		},
		{
			name:           "task not found",
			taskID:         validID,
			repoError:      apperrors.ErrTaskNotFound,
			wantErr:        apperrors.ErrTaskNotFound,
			repoShouldCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTodoRepo{
				deleteTaskError: tt.repoError,
			}

			service := &Services{repo: mockRepo}

			err := service.DeleteTask(
				context.Background(),
				tt.taskID,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}

			if mockRepo.deleteTaskCalled != tt.repoShouldCall {
				t.Fatalf(
					"expected repo called: %v, got: %v",
					tt.repoShouldCall,
					mockRepo.deleteTaskCalled,
				)
			}

		})
	}
}

func TestCreateTask(t *testing.T) {
	repoErr := errors.New("database error")

	tests := []struct {
		name           string
		taskReq        model.TodoCreateRequest
		repoError      error
		wantErr        error
		repoShouldCall bool
	}{
		{
			name: "valid task",
			taskReq: model.TodoCreateRequest{
				Title: "Buy groceries",
			},
			repoError:      nil,
			wantErr:        nil,
			repoShouldCall: true,
		},
		{
			name: "empty title",
			taskReq: model.TodoCreateRequest{
				Title: "",
			},
			repoError:      nil,
			wantErr:        apperrors.ErrTitleEmpty,
			repoShouldCall: false,
		},
		{
			name: "title too long",
			taskReq: model.TodoCreateRequest{
				Title: strings.Repeat("a", 201),
			},
			repoError:      nil,
			wantErr:        apperrors.ErrTitleTooLong,
			repoShouldCall: false,
		},
		{
			name: "invalid due date",
			taskReq: model.TodoCreateRequest{
				Title:   "Buy groceries",
				DueDate: new("invalid-date"),
			},
			repoError:      nil,
			wantErr:        apperrors.ErrInvalidDueDate,
			repoShouldCall: false,
		},
		{
			name: "repository error",
			taskReq: model.TodoCreateRequest{
				Title: "Buy groceries",
			},
			repoError:      repoErr,
			wantErr:        repoErr,
			repoShouldCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTodoRepo{
				createTaskError: tt.repoError,
			}

			service := &Services{
				repo: mockRepo,
			}

			err := service.CreateTask(
				context.Background(),
				tt.taskReq,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}

			if mockRepo.createTaskCalled != tt.repoShouldCall {
				t.Fatalf(
					"expected repo called: %v, got: %v",
					tt.repoShouldCall,
					mockRepo.createTaskCalled,
				)
			}
		})
	}
}
