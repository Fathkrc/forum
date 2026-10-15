package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) VotePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	isLike := r.URL.Path == "/like-post"

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = h.svc.VotePost(r.Context(), user.ID, postID, isLike)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func (h *Handler) VoteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	isLike := r.URL.Path == "/like-comment"

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = h.svc.VoteComment(r.Context(), user.ID, commentID, isLike)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}
