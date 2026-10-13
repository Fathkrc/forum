package handler

import (
	"cmd/forum/internal/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type HomeData struct {
	Posts      []model.Post
	Categories []model.Category
	SelectedID int
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	categories, err := h.svc.GetAllCategories(r.Context())
	// this part is  temporary !!!!
	// this part is temporary !!!!
	fmt.Println(categories)
	if err != nil {
		log.Printf("category error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, category := range categories {
		fmt.Fprintf(w, "ID: %d, Name: %s\n", category.ID, category.Name)
	}
	//return // temporary: don't continue into posts
	//  ^
	// / \
	//  I
	// this is temporary !!!!

	categoryParam := r.URL.Query().Get("category")

	var posts []model.Post
	selectedID := 0

	if categoryParam != "" {
		categoryID, err := strconv.Atoi(categoryParam)
		if err != nil || categoryID <= 0 {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}

		posts, err = h.svc.GetPostsByCategory(r.Context(), categoryID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		selectedID = categoryID
	} else {
		posts, err = h.svc.GetAllPosts(r.Context())
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
