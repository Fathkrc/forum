package handler

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/service"
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
	IsOwner     bool
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories, err := h.svc.GetAllCategories(r.Context())
		if err != nil {
			log.Printf("GetAllCategories error: %v", err)
			renderInternalError(w)
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
			var maxErr *http.MaxBytesError

			if errors.As(err, &maxErr) {
				renderError(w, http.StatusRequestEntityTooLarge, "Request body too large", "request body is more than 1 mb")
				return
			}
			renderError(w, http.StatusBadRequest,
				"Invalid Request",
				"The form could not be processed.")
			return
		}

		title := strings.TrimSpace(r.FormValue("title"))
		if len([]rune(title)) > 100 {
			renderError(w,
				http.StatusBadRequest,
				"Title Too Long",
				"Post titles cannot exceed 100 characters.")
			return
		}
		content := r.FormValue("content")

		var categoryIDs []int
		for _, categoryStr := range r.Form["categories"] {
			id, err := strconv.Atoi(categoryStr)
			if err != nil || id <= 0 {
				renderError(w, http.StatusBadRequest,
					"Invalid Request",
					"The selected category is invalid.")
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
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
	}
}

func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The post ID is invalid.")
		return
	}
	post, err := h.svc.GetPostByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("GetPostByID error: %v", err)
		renderNotFound(w)
		return
	}
	if err != nil {
		log.Printf("GetPostByID error: %v", err)
		renderInternalError(w)
		return
	}
	user, _ := getUserFromContext(r)
	comments, err := h.svc.GetCommentsByPost(r.Context(), id)
	if err != nil {
		log.Printf("GetCommentsByPost error: %v", err)
		renderInternalError(w)
		return
	}
	isOwner := user != nil && user.ID == post.UserID
	post.CommentCount = len(comments)
	// fmt.Println(post)
	// fmt.Println(&comments)

	renderTemplate(w, "post", postData{Post: *post, Comments: comments, CurrentUser: user, IsOwner: isOwner})
}
func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	posts, err := h.svc.GetAllPosts(r.Context())
	if err != nil {
		log.Printf("GetAllPosts error: %v", err)
		renderInternalError(w)
		return
	}
	if len(posts) == 0 {
		fmt.Fprintln(w, "No posts found.")
		return
	}

}
func (h *Handler) ListPostsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil || userID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The user ID is invalid.")
		return
	}

	posts, err := h.svc.GetPostsByUser(r.Context(), userID)
	if err != nil {
		log.Printf("GetPostsByUser error: %v", err)
		renderInternalError(w)
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
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	categoryID, err := strconv.Atoi(r.URL.Query().Get("category"))
	if err != nil || categoryID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The category ID is invalid.")
		return
	}

	posts, err := h.svc.GetPostsByCategory(r.Context(), categoryID)
	if err != nil {
		log.Printf("GetPostsByCategory error: %v", err)
		renderInternalError(w)
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
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	if query == "" {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"A search query is required.")
		return
	}

	posts, err := h.svc.SearchPosts(r.Context(), query)
	if err != nil {
		log.Printf("SearchPosts error: %v", err)
		renderInternalError(w)
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
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user"))
	if err != nil || userID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The user ID is invalid.")
		return
	}

	posts, err := h.svc.GetLikedPostsByUser(r.Context(), userID)
	if err != nil {
		log.Printf("GetLikedPosts error: %v", err)
		renderInternalError(w)
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

// func (h *Handler) TestGetPost(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		renderError(w, http.StatusMethodNotAllowed,
// 			"Method Not Allowed",
// 			"This request method is not supported.")
// 		return
// 	}
// 	id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err != nil || id <= 0 {
// 		renderError(w, http.StatusBadRequest,
// 			"Invalid Request",
// 			"The post ID is invalid.")
// 		return
// 	}
//
// 	post, err := h.svc.GetPostByID(r.Context(), id)
// 	if errors.Is(err, sql.ErrNoRows) {
// 		renderNotFound(w)
// 		return
// 	}
// 	if err != nil {
// 		log.Printf("GetPostByID error: %v", err)
// 		renderInternalError(w)
// 		return
// 	}
//
// 	fmt.Fprintf(w, "ID: %d, Title: %s, Content: %s\n",
// 		post.ID, post.Title, post.Content)
// }

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}
	postid := r.FormValue("post_id")
	id, err := strconv.ParseInt(postid, 10, 64)

	if err != nil || id <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The post ID is invalid.")
		return
	}

	user, ok := getUserFromContext(r)
	if !ok || user == nil {
		renderError(w, http.StatusUnauthorized,
			"Login Required",
			"Please log in to continue.")
		return
	}
	log.Printf("DELETE DEBUG: postID=%d userID=%d", id, user.ID)
	err = h.svc.DeletePost(
		r.Context(), id, int64(user.ID),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			renderNotFound(w)
		case errors.Is(err, service.ErrUnauthorized):
			renderError(w, http.StatusForbidden,
				"Access Denied",
				"You don't have permission to delete this post.")
		default:
			log.Printf("DeletePost error: %v", err)
			renderInternalError(w)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
