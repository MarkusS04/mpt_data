// Package middleware provides middlewares for http mux server
package middleware

import (
	"mpt_data/database/auth"
	"mpt_data/helper/config"
	"net/http"
	"strings"
)

// CheckAuthentication middleware checks if a user is correctly authenticated
func CheckAuthentication(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !config.Config.API.AuthenticationRequired {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing auth token", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "wrong authheader type", http.StatusUnauthorized)
			return
		}

		_, err := auth.ValidateJWT(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
