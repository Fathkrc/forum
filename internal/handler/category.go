package handler

import (
	"cmd/forum/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.svc.GetCategoryByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("GetCategory error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "ID: %d, Name: %s\n", category.ID, category.Name)
}
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category := model.Category{
		ID:   id,
		Name: r.FormValue("name"),
	}

	err = h.svc.UpdateCategory(r.Context(), category)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("UpdateCategory error: %v", err)
		http.Error(w, "Could not update category", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Category updated successfully")
}
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.svc.DeleteCategory(r.Context(), id)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("DeleteCategory error: %v", err)
		http.Error(w, "Could not delete category", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Category deleted successfully")
}
