package api

import (
	"cmd/forum/internal/handler"
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
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", port)

	return server.ListenAndServe()
}
