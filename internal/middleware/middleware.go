package middleware

import (
	"cmd/forum/internal/service"
	"context"
	"log"
	"net/http"
	"time"
)

type contextKey string

const UserKey contextKey = "user"

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("\033[32m  [LOG]\033[0m      %s localhost:8080%s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("\033[31m  [PANIC]\033[0m    %s localhost:8080%s %v", r.Method, r.URL.Path, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func LoadUser(svc *service.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			log.Printf("AUTH: user lookup failed: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		// fmt.Println(cookie.Value) // if there is no cookie it panics

		session, err := svc.GetSessionByToken(r.Context(), cookie.Value)
		// fmt.Println(session)
		if err != nil {
			log.Printf("AUTH: session lookup failed: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		user, err := svc.GetUserByID(r.Context(), int(session.UserID))
		if err != nil {
			log.Printf("AUTH: user lookup failed: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		log.Println("AUTH: user loaded successfully")
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
