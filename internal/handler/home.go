package handler

import (
	"cmd/forum/internal/model"
	"log"
	"net/http"
	"strconv"
)

type HomeData struct {
	Posts       []model.Post
	Categories  []model.Category
	SelectedID  int
	CurrentUser *model.User
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderNotFound(w)
		return
	}
	categories, err := h.svc.GetAllCategories(r.Context())
	// this part is  temporary !!!!
	// this part is temporary !!!!
	// fmt.Println(categories)
	if err != nil {
		log.Printf("category error: %v", err)
		renderInternalError(w)
		return
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
			renderError(w, http.StatusBadRequest,
				"Invalid Request",
				"The category ID is invalid.")
			return
		}

		posts, err = h.svc.GetPostsByCategory(r.Context(), categoryID)
		if err != nil {
			log.Printf("GetPostsByCategory error: %v", err)
			renderInternalError(w)
			return
		}
		selectedID = categoryID
	} else {
		posts, err = h.svc.GetAllPosts(r.Context())
		if err != nil {
			log.Printf("GetAllPosts error: %v", err)
			renderInternalError(w)
			return
		}
	}
	user, _ := getUserFromContext(r)
	data := HomeData{
		Posts:       posts,
		Categories:  categories,
		SelectedID:  selectedID,
		CurrentUser: user,
	}

	renderTemplate(w, "home", data)
}
