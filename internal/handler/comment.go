package handler

import (
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
