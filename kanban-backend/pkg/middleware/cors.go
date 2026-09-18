package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
			if allowedOriginsEnv == "" {
				allowedOriginsEnv = "http://localhost:8080"
			}

			origins := strings.Split(allowedOriginsEnv, ",")
			incomingOrigin := r.Header.Get("Origin")

			isAllowed := false
			for _, o := range origins {
				if strings.TrimSpace(o) == incomingOrigin {
					isAllowed = true
					break
				}
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", incomingOrigin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			next.ServeHTTP(w, r)
		},
	)
}