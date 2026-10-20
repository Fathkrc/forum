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
	Error       string
	Title       string
	Content     string
	Categories  []model.Category
	CurrentUser *model.User
}

type postData struct {
	Post        model.Post
	Comments    []*model.Comment
	Error       string
	CurrentUser *model.User
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories, err := h.svc.GetAllCategories(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		user, _ := getUserFromContext(r)
		renderTemplate(w, "create-post", createPostData{Categories: categories, CurrentUser: user})

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
				Error:       "Title and content are required",
				Title:       title,
				Content:     content,
				Categories:  categories,
				CurrentUser: user,
			})
			return
		}

		err = h.svc.CreatePost(r.Context(), user.ID, title, content, categoryIDs)
		if err != nil {
			categories, _ := h.svc.GetAllCategories(r.Context())
			renderTemplate(w, "create-post", createPostData{
				Error:       err.Error(),
				Title:       title,
				Content:     content,
				Categories:  categories,
				CurrentUser: user,
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
	user, _ := getUserFromContext(r)
	comments, err := h.svc.GetCommentsByPost(r.Context(), id)
	if err != nil {
		log.Printf("GetCommentsByPost error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	for _, comment := range comments {
		comment.Likes, comment.Dislikes, err = h.svc.GetCommentVoteCounts(r.Context(), int64(comment.ID))
		if err != nil {
			log.Printf("GetCommentVoteCounts %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	post.Likes, post.Dislikes, err = h.svc.VoteService.GetPostVoteCounts(r.Context(), int64(post.ID))
	if err != nil {
		log.Printf("GetPostVoteCounts error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(post)
	fmt.Println(&comments)

	renderTemplate(w, "post", postData{Post: *post, Comments: comments, CurrentUser: user})
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
	// for _, post := range posts {
	// 	post.Likes, post.Dislikes, err = h.svc.VoteService.GetPostVoteCounts(r.Context(), int64(post.ID))
	// 	if err != nil {
	// 		log.Printf("GetPostVoteCounts in ListPosts error: %v", err)
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	if len(posts) == 0 {
		fmt.Fprintln(w, "No posts found.")
		return
	}
	// todo : delete later
	for _, post := range posts {
		fmt.Fprintf(w, "ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
func (h *Handler) ListPostsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	posts, err := h.svc.GetPostsByUser(r.Context(), userID)
	if err != nil {
		log.Printf("GetPostsByUser error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(posts) == 0 {
		fmt.Fprintln(w, "No posts found for this user.")
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
func (h *Handler) ListPostsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryID, err := strconv.Atoi(r.URL.Query().Get("category"))
	if err != nil || categoryID <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	posts, err := h.svc.GetPostsByCategory(r.Context(), categoryID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		fmt.Fprintln(w, "No posts found in this category.")
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
func (h *Handler) SearchPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	posts, err := h.svc.SearchPosts(r.Context(), query)
	if err != nil {
		log.Printf("SearchPosts error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(posts) == 0 {
		fmt.Fprintln(w, "No matching posts found.")
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
func (h *Handler) ListLikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	posts, err := h.svc.GetLikedPostsByUser(r.Context(), userID)
	if err != nil {
		log.Printf("GetLikedPosts error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(posts) == 0 {
		fmt.Fprintln(w, "No liked posts found.")
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
func (h *Handler) TestGetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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
