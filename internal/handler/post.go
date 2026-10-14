package handler

import (
	"cmd/forum/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
		categories, err := h.svc.GetAllCategories()
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
			categories, _ := h.svc.GetAllCategories()
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
			categories, _ := h.svc.GetAllCategories()
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

func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	post, err := h.svc.GetPostByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("GetPostByID error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	comments, err := h.svc.GetCommentsByPost(id)
	if err != nil {
		log.Printf("GetCommentsByPost error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "post", postData{Post: *post, Comments: comments})
}
func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.svc.GetAllPosts(r.Context())
	if err != nil {
		log.Printf("GetAllPosts error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		fmt.Fprintln(w, "No posts found.")
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
func (h *Handler) TestGetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.svc.GetPostByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("GetPostByID error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "ID: %d, Title: %s, Content: %s\n",
		post.ID, post.Title, post.Content)
}
