package handler

//type Services interface {
//	TodoService
//}

type Handler struct {
	svc TodoService
}

func NewHandler(svc TodoService) *Handler {
	return &Handler{svc: svc}
}
