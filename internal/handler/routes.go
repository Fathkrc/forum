package handler

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/logout", h.requireAuth(h.Logout))

	mux.HandleFunc("/create-post", h.requireAuth(h.CreatePost))
	mux.HandleFunc("/post", h.ViewPost)
	mux.HandleFunc("/like-post", h.requireAuth(h.VotePost))
	mux.HandleFunc("/dislike-post", h.requireAuth(h.VotePost))

	mux.HandleFunc("/like-comment", h.requireAuth(h.VoteComment))
	mux.HandleFunc("/dislike-comment", h.requireAuth(h.VoteComment))
	mux.HandleFunc("/comment", h.requireAuth(h.AddComment))

	mux.HandleFunc("/categories", h.ListCategories)
	mux.HandleFunc("/categories/create", h.CreateCategory)
	mux.HandleFunc("/categories/view", h.GetCategory)
	mux.HandleFunc("/categories/update", h.UpdateCategory)
	mux.HandleFunc("/categories/delete", h.DeleteCategory)

	mux.HandleFunc("/profile", h.Profile)
	mux.HandleFunc("/my-posts", h.requireAuth(h.MyPosts))
	mux.HandleFunc("/liked-posts", h.requireAuth(h.LikedPosts))
	mux.HandleFunc("/search", h.Search)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
