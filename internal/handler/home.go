package handler

import (
	"cmd/forum/internal/model"
	"net/http"
	"strconv"
)

type HomeData struct {
	Posts      []model.Post
	Categories []model.Category
	SelectedID int
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.GetAllCategories()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	categoryParam := r.URL.Query().Get("category")

	var posts []model.Post
	selectedID := 0

	if categoryParam != "" {
		categoryID, err := strconv.Atoi(categoryParam)
		if err != nil {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}

		posts, err = h.svc.GetPostsByCategory(categoryID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		selectedID = categoryID
	} else {
		posts, err = h.svc.GetAllPosts()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	data := HomeData{
		Posts:      posts,
		Categories: categories,
		SelectedID: selectedID,
	}

	renderTemplate(w, "home", data)
}
