package handler

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type TodoService interface {
	GetTasks(ctx context.Context) ([]*model.Todo, error)
	GetTask(ctx context.Context, taskID string) (*model.Todo, error)
	CreateTask(ctx context.Context, taskReq model.TodoCreateRequest) error
	DeleteTask(ctx context.Context, taskIDStr string) error
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

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskReq model.TodoCreateRequest

	err := json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apperrors.ErrInvalidRequest)
		return
	}

	err = h.svc.CreateTask(r.Context(), taskReq)

	switch {
	case errors.Is(err, apperrors.ErrTitleEmpty):
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, apperrors.ErrTitleTooLong):
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		writeJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.PathValue("id")

	err := h.svc.DeleteTask(r.Context(), taskIDStr)

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

	w.WriteHeader(http.StatusNoContent)
}
