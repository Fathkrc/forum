package handler

import (
	"cmd/forum/internal/model"
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := h.svc.GetAllCategories(r.Context())
	if len(categories) == 0 {
		fmt.Fprintln(w, "No categories found.") // temporary
		return
	}

	if err != nil {
		log.Printf("GetAllCategories error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	for _, category := range categories {
		fmt.Fprintf(w, "ID: %d, Name: %s\n", category.ID, category.Name)
	}
}
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := model.Category{
		Name: r.FormValue("name"),
	}

	created, err := h.svc.CreateCategory(r.Context(), category)
	if err != nil {
		log.Printf("CreateCategory error: %v", err)
		http.Error(w, "Could not create category", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Created category: ID=%d, Name=%s\n", created.ID, created.Name)
}
