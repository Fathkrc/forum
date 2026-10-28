package api

import (
	"cmd/forum/internal/handler"
	"cmd/forum/internal/middleware"
	"cmd/forum/internal/repository"
	"cmd/forum/internal/service"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func Run(port string, db *sql.DB) error {
	// TODO: wire repositories + service once implementations are ready.

	repository := repository.NewRepository(db)
	svc := service.New(repository)
	h := handler.New(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	var httpHandler http.Handler = mux

	csrf := http.NewCrossOriginProtection()
	httpHandler = csrf.Handler(httpHandler)

	httpHandler = middleware.LoadUser(svc, httpHandler)
	httpHandler = middleware.Logger(httpHandler)
	httpHandler = middleware.Recover(httpHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           httpHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", port)

	return server.ListenAndServe()
}
