package handler

import (
	"cmd/forum/internal/model"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type profileData struct {
	User        model.User
	Posts       []model.Post
	CurrentUser *model.User
}

type postsData struct {
	Posts       []model.Post
	CurrentUser *model.User
}

type searchData struct {
	Posts       []model.Post
	Query       string
	CurrentUser *model.User
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The user ID is invalid.")
		return
	}

	user, err := h.svc.GetUserByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		renderNotFound(w)
		return
	}
	if err != nil {
		log.Printf("GetUserByID error: %v", err)
		renderInternalError(w)
		return
	}

	posts, err := h.svc.GetPostsByUser(r.Context(), id)
	if err != nil {
		log.Printf("GetPostsByUser error: %v", err)
		renderInternalError(w)
		return
	}
	currentUser, _ := getUserFromContext(r)

	renderTemplate(w, "profile", profileData{
		User:        *user,
		Posts:       posts,
		CurrentUser: currentUser,
	})
}

func (h *Handler) MyPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, err := h.svc.GetPostsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("MyPosts GetPostsByUser error: %v", err)
		renderInternalError(w)
		return
	}

	renderTemplate(w, "my-posts", postsData{Posts: posts, CurrentUser: user})
}

func (h *Handler) LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	posts, err := h.svc.GetLikedPostsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("LikedPosts error: %v", err)
		renderInternalError(w)
		return
	}

	renderTemplate(w, "liked-posts", postsData{Posts: posts, CurrentUser: user})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	user, _ := getUserFromContext(r)
	if query == "" {
		renderTemplate(w, "search", searchData{CurrentUser: user})
		return
	}

	posts, err := h.svc.SearchPosts(r.Context(), query)
	if err != nil {
		log.Printf("SearchPosts error: %v", err)
		renderInternalError(w)
		return
	}

	renderTemplate(w, "search", searchData{Posts: posts, Query: query, CurrentUser: user})
}
