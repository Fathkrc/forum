package handler

import (
	"cmd/forum/internal/service"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The post ID is invalid.")
		return
	}

	content := r.FormValue("content")
	if strings.TrimSpace(content) == "" {
		h.renderPostWithError(w, r, postID, "Comment cannot be empty")
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = h.svc.CreateComment(r.Context(), postID, user.ID, content)
	if err != nil {
		log.Printf("CreateComment error: %v", err)
		renderInternalError(w)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func (h *Handler) renderPostWithError(w http.ResponseWriter, r *http.Request, postID int, msg string) {
	post, err := h.svc.GetPostByID(r.Context(), postID)
	if errors.Is(err, sql.ErrNoRows) {
		renderNotFound(w)
		return
	}
	if err != nil {
		log.Printf("GetPostByID error: %v", err)
		renderInternalError(w)
		return
	}

	comments, err := h.svc.GetCommentsByPost(r.Context(), postID)
	if err != nil {
		log.Printf("GetCommentsByPost error: %v", err)
		renderInternalError(w)
		return
	}
	user, _ := getUserFromContext(r)
	renderTemplate(w, "post", postData{Post: *post, Comments: comments, Error: msg, CurrentUser: user})
}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	// Get comment ID from the form
	commentID, err := strconv.ParseInt(r.FormValue("comment_id"), 10, 64)
	if err != nil || commentID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The comment ID is invalid.")
		return
	}

	// Get logged-in user
	user, ok := getUserFromContext(r)
	if !ok || user == nil {
		renderError(w, http.StatusUnauthorized,
			"Login Required",
			"Please log in to continue.")
		return
	}

	// Delete comment through service
	err = h.svc.DeleteComment(r.Context(), commentID, int64(user.ID))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			renderNotFound(w)
		case errors.Is(err, service.ErrUnauthorized):
			renderError(w, http.StatusForbidden,
				"Access Denied",
				"You don't have permission to delete this comment.")

		default:
			log.Printf("DeleteComment error: %v", err)
			renderInternalError(w)
		} // we should make a global switch case for errors : optional
		return
	}

	// Return to the post
	postID := r.FormValue("post_id")
	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
