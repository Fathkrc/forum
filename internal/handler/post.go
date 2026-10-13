package handler

import (
	"cmd/forum/internal/model"
	"net/http"
	"strconv"
	"strings"
)

type createPostData struct {
	Error      string
	Title      string
	Content    string
	Categories []model.Category
}

type postData struct {
	Post     model.Post
	Comments []model.Comment
	Error    string
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories, err := h.svc.GetAllCategories(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "create-post", createPostData{Categories: categories})

	case http.MethodPost:
		user, ok := getUserFromContext(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		var categoryIDs []int
		for _, categoryStr := range r.Form["categories"] {
			id, err := strconv.Atoi(categoryStr)
			if err != nil {
				http.Error(w, "Invalid category", http.StatusBadRequest)
				return
			}
			categoryIDs = append(categoryIDs, id)
		}

		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			categories, _ := h.svc.GetAllCategories(r.Context())
			renderTemplate(w, "create-post", createPostData{
				Error:      "Title and content are required",
				Title:      title,
				Content:    content,
				Categories: categories,
			})
			return
		}

		err = h.svc.CreatePost(r.Context(), user.ID, title, content, categoryIDs)
		if err != nil {
			categories, _ := h.svc.GetAllCategories(r.Context())
			renderTemplate(w, "create-post", createPostData{
				Error:      err.Error(),
				Title:      title,
				Content:    content,
				Categories: categories,
			})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	post, err := h.svc.GetPostByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.svc.GetCommentsByPost(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "post", postData{Post: *post, Comments: comments})
}
