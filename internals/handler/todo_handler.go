package handler

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"errors"
	"net/http"
)

type TodoService interface {
	GetTasks(ctx context.Context) ([]*model.Todo, error)
	GetTask(ctx context.Context, taskID string) (*model.Todo, error)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetTasks(r.Context())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.PathValue("id")

	task, err := h.svc.GetTask(r.Context(), taskIDStr)

	switch {
	case errors.Is(err, apperrors.ErrTaskNotFound):
		writeJSON(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, apperrors.ErrInvalidID):
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		writeJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, task)
}
