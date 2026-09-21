package handler

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/service"
	"html/template"
	"net/http"
	"path/filepath"
)

type userContextKey string

const UserKey userContextKey = "user"

type Handler struct {
	svc service.Service_ghost
}

func New(svc service.Service_ghost) *Handler {
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

func renderTemplate(w http.ResponseWriter, name string, data any) {
	files := []string{
		filepath.Join("templates", name+".html"),
		filepath.Join("templates", "layout.html"),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
