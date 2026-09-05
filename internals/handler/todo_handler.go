package handler

import (
	"F5/internals/model"
	"context"
	"net/http"
)

type TodoService interface {
	GetTasks(ctx context.Context) ([]*model.Todo, error)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetTasks(r.Context())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}
