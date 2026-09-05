package handler

import (
	"F5/internals/apperrors"
	"F5/internals/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type TodoService interface {
	GetTasks(ctx context.Context, queries model.TodoListQuery) ([]*model.Todo, error)
	GetTask(ctx context.Context, taskID string) (*model.Todo, error)
	CreateTask(ctx context.Context, taskReq model.TodoCreateRequest) error
	DeleteTask(ctx context.Context, taskIDStr string) error
	UpdateTask(ctx context.Context, taskIDStr string, taskReq model.TodoUpdateRequest) error
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	q := model.TodoListQuery{
		Limit:  10,
		Offset: 0,
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		q.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, "Invalid offset")
			return
		}
		q.Offset = offset
	}

	if status := r.URL.Query().Get("status"); status != "" {
		q.Status = &status
	}

	tasks, err := h.svc.GetTasks(r.Context(), q)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "internal server error")
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

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var taskReq model.TodoUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apperrors.ErrInvalidRequest)
		return
	}

	taskIDStr := r.PathValue("id")

	err = h.svc.UpdateTask(r.Context(), taskIDStr, taskReq)

	switch {
	case errors.Is(err, apperrors.ErrInvalidID):
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, apperrors.ErrTaskNotFound):
		writeJSON(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, apperrors.ErrTitleEmpty),
		errors.Is(err, apperrors.ErrTitleTooLong),
		errors.Is(err, apperrors.ErrInvalidStatus):
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		writeJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
