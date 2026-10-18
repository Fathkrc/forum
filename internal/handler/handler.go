package handler

import (
	"cmd/forum/internal/middleware"
	"cmd/forum/internal/model"
	"cmd/forum/internal/service"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := getUserFromContext(r); !ok {
			log.Println("AUTH: requireAuth could not find user in context")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func getUserFromContext(r *http.Request) (*model.User, bool) {
	user, ok := r.Context().Value(middleware.UserKey).(*model.User)
	return user, ok
}

func renderTemplate(w http.ResponseWriter, name string, data any) {
	files := []string{
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", name+".html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
