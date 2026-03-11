package middleware

import (
	"net/http"
)

func Authenticate(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			w.Write([]byte("Error: Authorization token is empty"))
			return
		}

		next.ServeHTTP(w, r)
	})

}
