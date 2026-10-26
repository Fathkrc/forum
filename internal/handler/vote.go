package handler

import (
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) VotePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	isLike := r.URL.Path == "/like-post"

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The post ID is invalid.")
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = h.svc.VotePost(r.Context(), user.ID, postID, isLike)
	if err != nil {
		log.Printf("VotePost error: %v", err)
		renderInternalError(w)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func (h *Handler) VoteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed",
			"This request method is not supported.")
		return
	}

	isLike := r.URL.Path == "/like-comment"

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil || commentID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The comment ID is invalid.")
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		renderError(w, http.StatusBadRequest,
			"Invalid Request",
			"The post ID is invalid.")
		return
	}

	user, ok := getUserFromContext(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = h.svc.VoteComment(r.Context(), user.ID, commentID, isLike)
	if err != nil {
		log.Printf("VoteComment error: %v", err)
		renderInternalError(w)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}
