package handler

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/service"
	"context"
	"net/http"
)

type userContextKey string

const UserKey userContextKey = "user"

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := getUserFromContext(r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func getUserFromContext(r *http.Request) (*model.User, bool) {
	user, ok := r.Context().Value(UserKey).(*model.User)
	return user, ok
}
