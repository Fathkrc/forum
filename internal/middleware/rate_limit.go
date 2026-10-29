package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	tokens float64
	last   time.Time
}

func RateLimit(next http.Handler) http.Handler {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	const rate = 10.0
	const burst = 20.0

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid client address", http.StatusBadRequest)
			return
		}

		now := time.Now()

		mu.Lock()

		v, exists := visitors[ip]
		if !exists {
			v = &visitor{
				tokens: burst,
				last:   now,
			}
			visitors[ip] = v
		}

		elapsed := now.Sub(v.last).Seconds()
		v.tokens += elapsed * rate

		if v.tokens > burst {
			v.tokens = burst
		}

		v.last = now

		allowed := v.tokens >= 1
		if allowed {
			v.tokens--
		}

		mu.Unlock()

		if !allowed {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
