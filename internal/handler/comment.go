package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func (h *Handler) renderPostWithError(w http.ResponseWriter, r *http.Request, postID int, msg string) {
	post, err := h.svc.GetPostByID(r.Context(), postID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.svc.GetCommentsByPost(r.Context(), postID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user, _ := getUserFromContext(r)
	renderTemplate(w, "post", postData{Post: *post, Comments: comments, Error: msg, CurrentUser: user})
}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get comment ID from the form
	commentID, err := strconv.ParseInt(r.FormValue("comment_id"), 10, 64)
	if err != nil || commentID <= 0 {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Get logged-in user
	user, ok := getUserFromContext(r)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete comment through service
	err = h.svc.DeleteComment(r.Context(), commentID, int64(user.ID))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			http.Error(w, "Comment not found", http.StatusNotFound)

		case errors.Is(err, errors.New("unauthorized")): //change to custom
			http.Error(w, "Forbidden", http.StatusForbidden)

		default:
			log.Printf("DeleteComment error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} // we should make a global switch case for errors : optional
		return
	}

	// Return to the post
	postID := r.FormValue("post_id")
	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
