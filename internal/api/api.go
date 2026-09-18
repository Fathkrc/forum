package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func Run(port string, db *sql.DB) error {
	// TODO: wire repositories + service once implementations are ready.

	// svc := service.NewService( /* dependencies */ )
	// h := handler.New(svc)
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server running on http://localhost:%s", port)

	return server.ListenAndServe()
}
