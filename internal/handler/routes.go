package handler

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/register", h.Register)

	mux.HandleFunc("/logout", h.requireAuth(h.Logout))

	mux.HandleFunc("/create-post", h.requireAuth(h.CreatePost)) //this is real route
	// mux.HandleFunc("/create-post", h.CreatePost) // this is for test todo: remove
	mux.HandleFunc("/posts/user", h.ListPostsByUser)
	mux.HandleFunc("/posts/search", h.SearchPosts)
	mux.HandleFunc("/post", h.GetPostById)
	mux.HandleFunc("/like-post", h.requireAuth(h.VotePost))
	mux.HandleFunc("/dislike-post", h.requireAuth(h.VotePost))
	mux.HandleFunc("/posts/liked", h.requireAuth(h.ListLikedPosts))
	mux.HandleFunc("/posts/list", h.ListPosts)
	mux.HandleFunc("/posts/test-view", h.TestGetPost)
	mux.HandleFunc("/posts/category", h.ListPostsByCategory)
	mux.HandleFunc("/posts/delete", h.requireAuth(h.DeletePost))

	mux.HandleFunc("/like-comment", h.requireAuth(h.VoteComment))
	mux.HandleFunc("/dislike-comment", h.requireAuth(h.VoteComment))
	mux.HandleFunc("/comment", h.requireAuth(h.AddComment))
	mux.HandleFunc("/comments/delete", h.requireAuth(h.DeleteComment))

	mux.HandleFunc("/categories", h.ListCategories)
	mux.HandleFunc("/categories/create", h.CreateCategory)
	mux.HandleFunc("/categories/view", h.GetCategory)
	mux.HandleFunc("/categories/update", h.UpdateCategory)
	mux.HandleFunc("/categories/delete", h.DeleteCategory)

	mux.HandleFunc("/profile", h.Profile)
	mux.HandleFunc("/my-posts", h.requireAuth(h.MyPosts))
	mux.HandleFunc("/liked-posts", h.requireAuth(h.LikedPosts))
	mux.HandleFunc("/search", h.Search)

	//error checkers
	// Error page testing (remove before submission)
	mux.HandleFunc("/test-400", func(w http.ResponseWriter, r *http.Request) {
		renderError(w, http.StatusBadRequest,
			"Invalid Request", "The request is invalid.")
	})

	mux.HandleFunc("/test-401", func(w http.ResponseWriter, r *http.Request) {
		renderError(w, http.StatusUnauthorized,
			"Login Required", "Please log in to continue.")
	})

	mux.HandleFunc("/test-403", func(w http.ResponseWriter, r *http.Request) {
		renderError(w, http.StatusForbidden,
			"Access Denied", "You don't have permission to access this page.")
	})

	mux.HandleFunc("/test-404", func(w http.ResponseWriter, r *http.Request) {
		renderNotFound(w)
	})

	mux.HandleFunc("/test-405", func(w http.ResponseWriter, r *http.Request) {
		renderError(w, http.StatusMethodNotAllowed,
			"Method Not Allowed", "This request method is not supported.")
	})
	mux.HandleFunc("/test-500", func(w http.ResponseWriter, r *http.Request) {
		renderInternalError(w)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
